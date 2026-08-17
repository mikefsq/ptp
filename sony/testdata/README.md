# Sony test fixtures — making your own

This directory is empty on a fresh clone, and `.gitignore` keeps it that way:
these are bytes off a camera, and one ARW is 65–75 MB. Every test that reads
them calls `t.Skipf` when the file is absent, so `go test ./sony/` passes green
having checked neither `ptp.ParseDeviceInfo` nor `ptp.ParseStdPropDesc` against
a real reply, and having exercised the lossless decoder not at all.

**A green run on a clean checkout is not a green light.** If you have a Sony
body, this is how to turn those skips into real coverage.

## What the tests are looking for

| put it here | read by | what it proves |
|---|---|---|
| `nex6-deviceinfo.bin` | `TestParseDeviceInfoNEX6`, `TestNEX6HasNoSDIOSurface`, `TestConnectRefusesABodyWithoutSDIO`, `TestParseDeviceInfoTruncated` | `ptp.ParseDeviceInfo`, and that `Connect` refuses a body with no vendor surface |
| `nex6-propdesc-5001.bin` | `TestParseStdPropDescBatteryLevel`, `TestParseStdPropDescTruncated` | the numeric descriptor layout |
| `nex6-propdesc-d402.bin` | `TestParseStdPropDescFriendlyName` | the string descriptor layout |
| `ljpeg/tile.ljpeg` | `TestLJPEGDecodesARealStream` | the lossless-JPEG codec on its own |
| `ljpeg/_DSC4922.ARW` | `TestFastPathMatchesGeneric`, `TestTiledFallbackMatchesSharedPrelude` | the tiled ARW path |
| `ljpeg/_DSC2429.ARW` | `TestFastPathMatchesGeneric` | a second geometry |
| `ljpeg/Bias.ARW` | `TestFastPathMatchesGeneric` | pure read noise — the worst case for the coder |

The filenames are literal; the tests open exactly these paths. See
[If your body is not a NEX-6](#if-your-body-is-not-a-nex-6).

## Wire captures — all three from one run

    go run ./cmd/sonyprobe -dump /tmp/probe

Set the body to PC Remote and switch it on. `sonyprobe` writes a `session.log`,
a `deviceinfo.bin`, and one `NNN-opXXXX-in.bin` per transaction that carried
data. It is read-only unless you pass `-set`.

On a body with no vendor surface — which is most older Alphas — `sonyprobe`
falls through to its standard-PTP branch and describes every property
`DeviceInfo` lists using operation 0x1014. That one run produces all three
fixtures.

    cp /tmp/probe/deviceinfo.bin sony/testdata/nex6-deviceinfo.bin

The two descriptors need picking out of the bundle. Their sequence numbers are
not predictable — the sweep follows whatever order the body lists its properties
in — so identify them by their first two bytes, which are the property code
little-endian:

    01 50 02 00 00 64 ...        -> 0x5001  BatteryLevel
    02 d4 ff ff 00 00 06 4e ...  -> 0xD402  DeviceFriendlyName

    cp /tmp/probe/NNN-op1014-in.bin sony/testdata/nex6-propdesc-5001.bin
    cp /tmp/probe/NNN-op1014-in.bin sony/testdata/nex6-propdesc-d402.bin

The run then exits non-zero, reporting that the body exposes no Sony
remote-control operations. On such a body that is the expected result and not a
failure — the bundle is written as the transactions happen, so every file is
already on disk before the message prints.

**Take one of each kind.** The pair is deliberate: 0x5001 is numeric with a
62-value enum, 0xD402 is string-typed. PTP strings are a uint8 character count
of UTF-16 with the NUL counted, which is the part most easily got wrong, and a
parser tested only against numbers will not catch it. Any two properties of
those two shapes will do.

## `ljpeg/` — RAW frames

Any lossless-compressed ARW from a recent body will do. Set it to compressed
RAW, then shoot two things:

- **a scene** — ordinary subject, real detail
- **a bias frame** — lens cap on, fastest shutter, base ISO

Shoot both on the same body so they share a geometry. That is the point of the
pair: it separates content from geometry, and it is what showed that a bias
frame is the *slowest* of the three to decode rather than the fastest. Pure
noise is incompressible, so every tile runs long.

Pull them over USB:

    go run ./cmd/sonyprobe -ls
    go run ./cmd/sonyprobe -get sony/testdata/ljpeg -max 3

or take the card out and copy them, which is quicker for a batch. Rename to the
names in the table.

## `tile.ljpeg` — one tile, from anything

A single lossless-JPEG stream: JFIF header, SOF3, 10-bit, 504x504, three
components, predictor 7. `TestLJPEGDecodesARealStream` asserts all of that, so a
replacement of different geometry needs those numbers updated with it.

Take it from a **different manufacturer** on purpose. Sony's lossless is plain
lossless JPEG, so a non-Sony stream is what proves the codec is implemented from
the spec rather than fitted to one vendor's files. The committed one came from
an iPhone ProRAW DNG.

No tool here extracts it. Open a raw DNG, find the SubIFD whose `Compression`
(0x0103) is 7, and slice out the byte range given by one entry of `TileOffsets`
(0x0144) and the matching `TileByteCounts` (0x0145). That range is a complete
JPEG stream on its own, SOI to EOI, and is all the file holds.

Take the tile from a **flat part** of the frame. The test also checks that the
first row is smooth — neighbouring samples differing by a few counts, not
hundreds — which a noise tile would fail.

## If your body is not a NEX-6

The capture procedure is not body-specific, but the test paths are hardcoded, so
a file named for your camera is a file no test opens. Either name yours as the
table above, or add cases — the frame list in `ljpeg_fast_test.go` is three
strings, and a third geometry alongside the existing two makes the suite
stronger.

One thing does depend on the body. `TestNEX6HasNoSDIOSurface` and
`TestConnectRefusesABodyWithoutSDIO` need a DeviceInfo from a camera with **no**
vendor operations, because what they prove is that `Connect` refuses such a body
rather than stalling its bulk pipe on an unsupported operation. A modern Alpha
has the full 0x92xx surface and will make both tests fail rather than skip. If
that is all you have, capture the frames and leave the DeviceInfo fixture out.

## Two descriptor layouts, and why they are parsed separately

Worth knowing before you trust a capture. `nex6-propdesc-*.bin` are the STANDARD
descriptor layout: no `isEnabled` byte, values starting at +5. Sony's own 0x9209
entries insert `isEnabled` and start values at +6. Parsing either with the
other's offsets yields plausible-looking rubbish rather than an error, which is
why `ptp.ParseStdPropDesc` and this package's `devprop.go` are deliberately two
parsers. Keeping golden bytes of the standard shape here is what makes that
contrast checkable.

## Keeping your captures out of git

`.gitignore` excludes everything in this directory except this README, by
pattern rather than by filename, so whatever you capture and whatever you call
it stays local. You do not need to add anything.
