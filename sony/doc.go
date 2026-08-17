// Package sony drives Sony Alpha bodies over PTP.
//
// # Status
//
// The transport, framing, object and property layers are hardware-validated on
// a NEX-6. The vendor layer — the SDIO handshake, exposure, focus and capture —
// is decoded but has NEVER touched a camera that supports it: a NEX-6 has no
// SDIO operations, so it cannot exercise any of it. Validation on an A7R V or
// A7R VI is the outstanding work.
//
// Treat every 0x92xx path here as decoded, tested against the SDK's own
// structures, and unverified on hardware.
//
// # Nothing here is standard capture
//
// Sony has NO InitiateCapture. The shutter is two virtual buttons, S1 (half
// press, autofocus and metering) and S2 (the shutter), pressed through vendor
// operation 0x9207 — and they STAY DOWN until released, because the camera has
// no timeout. A host that exits mid-capture leaves the shutter held, and in a
// continuous drive mode that means the body keeps shooting.
//
// Camera tracks what this session pressed and lets it up on Close. It releases
// only what it pressed: a body without the vendor surface stalls its bulk pipe
// on an unsupported operation, so a speculative release is not free.
//
// # Live view
//
// Sony has no live view operation. The preview is an ordinary PTP object at the
// magic handle 0xFFFFC002, fetched with the standard GetObjectInfo and
// GetObject. See liveview.go. The payload carries a leading header of unpinned
// length, so LiveFrame finds the JPEG rather than assuming an offset.
//
// # Reading and writing settings
//
// The primary read path is one operation. GetAllDevicePropData (0x9209) returns
// every property's value, valid set and enable flag in a single blob, and the
// SDK polls that rather than reading properties one at a time. Take a snapshot,
// then work from it — the per-property accessors cost a round trip each.
//
// A property's valid set is authoritative and it MOVES: it changes with lens,
// exposure mode and drive mode. A value outside the current set is commonly
// accepted and then ignored rather than refused, so a write must be read back.
// The Nearest helpers snap a request into the set the camera is offering.
//
// # Traps that cost real time
//
//   - Packed shutter values sort BACKWARDS. In the fraction form 1/N encodes as
//     0x0001_000N, so a LARGER wire value is a SHORTER exposure and 1/1000
//     sorts above 1/2. Use NearestShutter, which compares decoded durations;
//     Nearest compares raw values and will pick badly at the fast end.
//   - Aperture is F×100 and ISO packs a mode into bits 24-27. Neither is the
//     number you would expect to send.
//   - A body changes its USB product ID with its mode. Only the still-image or
//     PC Remote PID can be driven; the same camera charging or mounted as mass
//     storage is a different PID and has no vendor surface.
//   - Cameras pad transfers, and ObjectInfo reports the PADDED size. A NEX-6
//     rounds CompressedSize up to a 4 KiB boundary and zero-fills the
//     remainder: a 1,611,005-byte JPEG is declared and sent as 1,638,400 bytes
//     (1600 KiB) with 27,395 trailing zeros after the EOI marker. The file
//     still decodes, but it is not byte-exact with what is on the card, and a
//     transfer length matching ObjectInfo proves nothing. Confirmed on
//     hardware; trim at the format's own end marker if an exact copy matters.
//   - OK does not mean applied. A NEX-6 answers OK to a write it then ignores,
//     while stalling the pipe for a genuinely illegal one.
//
// # Naming
//
// PropName consults Sony's own table first, then the standard PTP names. Sony's
// wins because the same code means different things in the vendor surface than
// it does as a standard property.
//
// ValueName names a property's VALUE, which is the part that matters when
// reading a camera: an enum printed as an ordinal says nothing about which mode
// the body is in, and the exposure triangle is encoded rather than enumerated —
// a shutter speed of 1/1000 packs to 0x000103E8 and prints as 66536 unless it is
// decoded. Every name is transcribed from an enum in CrDeviceProperty.h; where
// no enum exists the value stays a number, which is the honest answer.
//
// # String metadata
//
// Six properties are documented CrDataType_STR, and that is the whole string
// surface. Sony has NO Artist and NO Comment field — the credit is
// SetPhotographer — so accessors named for them would be inventing fields the
// body does not have. Both credit strings are gated by WriteCopyrightInfo,
// whose off value is 1 rather than 0; setting the text without that switch
// changes nothing in the files.
//
// Supported reports what Sony's DOCUMENTATION says a model offers. ExtDeviceInfo,
// filled in by Connect, is what the body in front of you actually reported. When
// they disagree, believe the camera.
package sony
