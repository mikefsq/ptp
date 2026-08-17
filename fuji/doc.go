// Package fuji drives Fujifilm X and GFX bodies over PTP.
//
// Hardware-validated on an X-T5 (firmware 1.00): capture, download, live view,
// histogram, exposure from 5us to 64s, bracketed bursts, clock sync, string
// metadata, and a teardown that hands the camera back to its owner.
//
// # The one rule that explains most surprises
//
// A physical control in an automatic position takes ownership of its setting,
// and the property then advertises exactly ONE value. Writes to it are accepted
// and silently ignored — no error, no effect. This holds for every control:
//
//	shutter dial on T ......... host sets the shutter, 76 values, 5us..64s
//	shutter dial on a speed ... camera owns it, one value
//	shutter dial on B ......... camera owns it; its own bulb timer, not the host's
//	aperture ring on A ........ exposure program is S; the camera picks the f-number
//	aperture ring off A ....... exposure program is M; the ring picks it
//	ISO dial on C ............. host sets ISO, 24 values including three AUTO presets
//	ISO dial on A ............. camera owns it, locked to auto1
//	focus lever on M .......... host may read and set FocusPos
//	focus lever on S or C ..... camera owns focus; FocusPos reads 0
//	drive dial ................ owns DriveMode outright; the host cannot change it
//
// Camera.Settable reports this, and Camera.Ready checks the preconditions that
// stop a host taking control at all. Both exist because the alternative is a
// capture that fails obscurely.
//
// # Taking control
//
// Capture needs PC Priority, which LOCKS OUT the body's own dials and shutter.
// It has to be given back — Camera.Close does that, and a host that merely
// exits leaves the camera inert in its owner's hands.
//
// The camera refuses to grant or return priority while it holds undownloaded
// frames. Downloading and deleting need no priority, so a drain is always
// possible; demanding priority first would deadlock.
//
// # Where frames go
//
// Every tethered frame lands in a volatile RAM store holding about five, and
// stays there until deleted. The SD card is not addressable at all in this
// mode. MediaRecord adds a card copy but never removes the RAM one, so the
// buffer must be drained either way.
//
// Two workable modes:
//
//	USB only     MediaRecordOff, download each frame. ~800ms per frame, and a
//	             dropped download loses it permanently.
//	Card only    MediaRecordRaw, then DELETE without downloading. 25ms per
//	             frame; the card copy survives, verified on hardware. Nothing is
//	             checkable from the host until afterwards.
//
// Card-only sustains 0.65s per frame, or 0.76s bracketed. Downloading cannot
// keep pace with a burst at any lag: a frame needs seconds before it is
// readable and a burst produces one every 0.65s, so keeping up would need ~15
// frames of lag against a buffer of five.
//
// Afterwards the card IS reachable: the body's card-reader USB mode exposes it
// as a normal removable volume, though capture is unavailable there.
//
// # Naming
//
// PropName consults, in order: names established by experiment
// (names_observed.go), names recovered from the body's own plugin binary
// (names_xt5_gen.go), gphoto2's table, then the standard PTP names. The plugin
// outranks gphoto2 because gphoto2 is sometimes WRONG for this body — it calls
// 0xD241 ImageAspectRatio where the camera reports shutter speeds.
//
// XT5Props (generated from captured descriptors) records what a real body
// advertised. Treat it as reference: validate against a live descriptor,
// because the advertised set changes with the camera's state.
package fuji
