# Fujifilm test fixtures — making your own

This directory is empty on a fresh clone, and `.gitignore` keeps it that way:
these are bytes off a camera, and one RAW frame is 25–80 MB. Every test that
reads them calls `t.Skipf` when the file is absent, so `go test ./fuji/` passes
green having checked the DeviceInfo decode against nothing and the RAF decoders
against nothing at all.

**A green run on a clean checkout is not a green light.** If you have a
Fujifilm body, half an hour with it turns those skips into real coverage. This
is how.

## What the tests are looking for

| put it here | read by | what it proves |
|---|---|---|
| `xt5-deviceinfo.bin` | `TestParseRealDeviceInfo` | `ptp.ParseDeviceInfo` against a real reply |
| `raf/xt5-uncompressed.raf` | `TestDecodeRAFMatchesLibRaw` | the uncompressed RAF path |
| `raf/xt5-lossless.raf` | `TestDecodeRAFMatchesLibRaw` | the lossless decoder |
| `raf/xt5-lossless-lit.raf` | `TestDecodeRAFMatchesLibRaw`, `TestBlocksAreIndependent` | the decoder, and its per-block parallelism |
| `raf/*.raf.tiff` | `TestDecodeRAFMatchesLibRaw` | the libraw reference each frame is compared to |

The filenames are literal — the tests open exactly these paths. Read
[If your body is not an X-T5](#if-your-body-is-not-an-x-t5) before you start if
it isn't.

## The tool

`cmd/fujiprobe` does all of it. Every run writes a capture bundle — a
`session.log`, a `deviceinfo.bin`, and one `NNN-opXXXX-in.bin` per transaction
that carried data — into the directory given by `-dump`, default
`fujiprobe-capture/`.

Set the camera's USB mode to PC Shoot / USB Tether, not card reader, and switch
it on. `fujiprobe` writes nothing to the camera unless you pass an exposure
flag.

    go run ./cmd/fujiprobe -list        # confirm it is seen

## `xt5-deviceinfo.bin` — one GetDeviceInfo reply

    go run ./cmd/fujiprobe -dump /tmp/probe
    cp /tmp/probe/deviceinfo.bin fuji/testdata/xt5-deviceinfo.bin

That is the whole job — 771 bytes on an X-T5. The test checks the model string,
the operation and property lists, and the capture capability flags, so it fails
on any decode that is subtly wrong about the layout.

## `raf/` — three RAW frames

Shoot the *same scene at the same settings* in each RAW recording mode. The
frames only isolate the compression if nothing else changes between them.

The exposure flags need the dials handed to the camera first: the shutter dial
on **T**, the aperture ring on **A**, and the ISO dial on **C**. `-mf` stops the
half press hunting on a static subject.

    go run ./cmd/fujiprobe -mf -quality raw -raw uncompressed \
        -shutter 397ms -aperture 4 -iso 640 -capture -get fuji/testdata/raf

    go run ./cmd/fujiprobe -mf -quality raw -raw lossless \
        -shutter 397ms -aperture 4 -iso 640 -capture -get fuji/testdata/raf

Rename what lands to `xt5-uncompressed.raf` and `xt5-lossless.raf`. `-raw lossy`
is the third mode and is not used — the decoder does not implement it.

**Point the camera at something with highlights in it.** The third frame,
`xt5-lossless-lit.raf`, exists because the first two were nearly black, and
against a dark frame a mosaic-phase error produced deltas of a few counts and
read as noise rather than as the bug it was. Against a frame using the full
0–16383 range the same error was obvious immediately. Shoot it lossless, well
exposed, with real detail in it:

    go run ./cmd/fujiprobe -mf -quality raw -raw lossless \
        -shutter 397ms -aperture 4 -iso 640 -capture -get fuji/testdata/raf

A degenerate frame is worth keeping too — that is why `xt5-lossless.raf` is
still in the list. It is the case that hides errors, so a decoder that handles
both is the one you want.

## The libraw references

`TestDecodeRAFMatchesLibRaw` does not decode a frame and look at it. It compares
every sample against what libraw produces, because a wrong decoder yields a
plausible image and "it looks right" is not evidence. It needs a reference TIFF
beside each frame, and skips without one.

Install libraw (`brew install libraw`, or your distribution's `libraw-bin`),
then, for each frame:

    unprocessed_raw -T fuji/testdata/raf/xt5-lossless-lit.raf

which writes `xt5-lossless-lit.raf.tiff` next to it. They are 82 MB each and
regenerate exactly, so there is no reason to keep them anywhere else.

Then:

    go test -run TestDecodeRAFMatchesLibRaw -v ./fuji/

A pass means every one of the 40,902,912 samples matched.

## If your body is not an X-T5

Everything above still works — `fujiprobe` is not X-T5-specific — but the tests
open hardcoded paths, so a frame from an X-H2 named for an X-H2 is a frame no
test looks at, and you get the same silent skip you started with.

Two options. Name your files as the table above, which is the quick way and
costs the honesty of the name. Or add your own cases: the frame list in
`raf_lossless_verify_test.go` is three strings, and another body's frames
alongside the X-T5's make the suite stronger rather than merely different.

The decoder itself reads the X-Trans phase from tag 0x0131 in the file rather
than looking the model up in a table, precisely so an unrecognised body works.
An X-T5's phase is `GRBGBR BGGRGG RGGBGG GBRGRB RGGBGG BGGRGG`, which appears in
none of the published model-to-pattern tables.

## Descriptor sweeps, if you want a property table

`XT5Props` in `props_xt5_gen.go` is a snapshot of what one X-T5 advertised. No
test reads it and nothing here regenerates it — the generator is not part of
this repository — but the raw material is one command, and it is the way to see
what *your* body offers:

    go run ./cmd/fujiprobe -all -dump /tmp/sweep

Pass no capture flag. `-capture`, `-get`, `-card`, `-quality` and `-raw` each
divert the run into the capture path, and the descriptor sweep never happens —
which is the one way to come away with an empty directory and no error.

`-all` does not change what is dumped; the sweep asks for every property either
way. It changes the log, which is the only record of which properties the body
*declined* to describe, because a declined property produces no file at all.

Two runs cover a Fujifilm body, and the second needs a hand on the camera: move
the STILL/MOVIE collar to movie and repeat. The movie properties cannot be
selected over USB, because DriveMode does not offer `XSDK_DRIVE_MODE_MOVIE`. The
order matters if you ever merge them — movie mode does not merely add
properties, it changes what some codes report, so stills should win every code
it describes.

A run that answers `GeneralError` for a block of properties is healthy: the
camera is declining to describe those, which is a fact about the properties. The
declines cluster, because the unsupported codes are consecutive. Only a
transport failure stops the sweep, and `fujiprobe` gives up after five in a row.

## Keeping your captures out of git

`.gitignore` excludes everything in this directory except this README, by
pattern rather than by filename, so whatever you capture and whatever you call
it stays local. You do not need to add anything.
