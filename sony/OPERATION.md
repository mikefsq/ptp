# Driving a Sony body over PTP

Practical notes for `ptp/sony`. Unlike [fuji/OPERATION.md](../fuji/OPERATION.md),
most of this is **not yet confirmed on hardware** — the only body driven so far
is a NEX-6, which has no vendor surface. What has and has not met a camera is
marked throughout, and the decode evidence is in [PROTOCOL.md](PROTOCOL.md).

## Before connecting

Set the body's USB connection mode to **PC Remote**. In a file-transfer mode it
enumerates and opens perfectly well but exposes no vendor operations at all, and
every 0x92xx call then refuses by stalling the bulk pipe — which looks like a
driver fault and is not.

Disable the power-save timeout if the body allows it. A sleeping camera drops
off the bus mid-run and produces failures that look exactly like protocol bugs.

## macOS: ptpcamerad has to be killed, every run

`USBInterfaceOpenSeize` does **not** wrest a still-image device from macOS's
`ptpcamerad` — both the open and the seize return `kIOReturnExclusiveAccess`
(`0xE00002C5`). It has to be terminated, and `SIGSTOP` is no help because a
stopped process keeps its claim.

Worse, killing it is not durable. It is demand-launched whenever a still-image
device is enumerated, so it returns about a second into any run; `launchctl
disable` does not stop that, and removing it outright is blocked by SIP:

    $ launchctl bootout gui/$(id -u)/com.apple.ptpcamerad
    Boot-out failed: 150: Operation not permitted while System Integrity Protection is engaged

So the working pattern is: **build first**, then kill and grab in one step, so a
compile does not sit between them.

    go build -o sonyprobe ./cmd/sonyprobe
    ./sonyprobe -list
    pkill -9 ptpcamerad; ./sonyprobe -serial <serial from that list>

Restore normal camera handling when finished:

    launchctl enable gui/$(id -u)/com.apple.ptpcamerad

Measured 2026-08-07: a respawn part-way through a run does **not** by itself
kill a live session — several runs completed straight through one. But it is the
most likely explanation for a transfer aborted mid-session with
`kIOReturnAborted` (`0xE00002EB`).

## When the camera stops answering

A transaction abandoned mid-transfer leaves the body waiting for the rest of a
data phase. It then answers nothing, every request costs a full timeout, and it
stays that way across process restarts.

This no longer needs a power cycle. `Session.Open` clears both bulk endpoints
and sends the still-image class Device Reset (`bRequest 0x66`) before giving up,
which a NEX-6 accepts. To do it explicitly:

    pkill -9 ptpcamerad; ./sonyprobe -reset

| symptom | cause | fix |
|---|---|---|
| claim fails `0xE00002C5` | ptpcamerad holds the interface | `pkill -9 ptpcamerad`, then open promptly |
| claim succeeds, first write times out | body wedged by an abandoned transfer | `-reset`, or just retry — Open resets automatically |
| `kIOReturnAborted` mid-session | something reset the device underneath us | re-run; suspect ptpcamerad |
| opens, but 4-5 properties and no vendor ops | body is not in PC Remote mode | change the USB connection setting |
| nothing on the bus | asleep, or charging-only PID | wake it; check the PID is the still-image one |

## Testing by proxy

Built for the case where whoever has the camera is not whoever is debugging the
driver. Every run writes a **capture bundle** holding each raw payload the body
sent:

    sonyprobe-capture/
      session.log               full transcript, including every transaction
      deviceinfo.bin            standard PTP DeviceInfo (model, firmware)
      extdeviceinfo-9202.bin    the vendor handshake reply
      NNN-op9209-in.bin         the raw property blob
      NNN-opXXXX-in.bin         every other data phase, in order

That bundle re-parses the whole session with no camera attached:

    go run ./cmd/sonyprobe -replay sonyprobe-capture

So a decode bug found on someone else's A7R VI can be fixed and re-verified
locally without another trip to the hardware. The bundle is written **even when
the run fails** — a failed parse is when the bytes matter most — and the tool
reports which property index it died on.

## Instructions to hand to whoever has the camera

1. Set the body's USB connection mode to **PC Remote**.
2. Connect it by USB and power it on. Disable the power-save timeout if you can.
3. Build first, so the kill and the open are not separated by a compile:

   ```
   git clone https://github.com/mikefsq/ptp && cd ptp
   go build -o sonyprobe ./cmd/sonyprobe
   ./sonyprobe -list
   pkill -9 ptpcamerad; ./sonyprobe -serial <serial from that list>
   ```
4. Send back the whole `sonyprobe-capture` directory, whether or not it worked.

With two bodies, run once per serial into separate directories with `-dump`.

## Shooting

**Unverified — no body supporting these operations has been driven.**

There is no `InitiateCapture`. A shot is S2 down then S2 up, and the buttons
stay down until released:

    cam.SetManualFocus()          // before any unattended sequence
    cam.SetExposureProgram(sony.ExposureManual)
    cam.SetShutter(time.Second / 500)
    cam.Shoot()

Manual focus matters for the same reason it does on a Fujifilm body: autofocus
on a subject the lens cannot lock — a dark sky, a filtered sun — hunts
indefinitely.

For a sequence, take ONE property snapshot and reuse it. Each
`GetAllDevicePropData` is a full round trip, and `SetShutter` fetches one
implicitly; `SetShutterFrom` does not.

In a continuous drive mode, `Burst` holds S2 for a duration and the camera fires
at its own rate — two PTP transactions for the whole burst rather than two per
frame. `StoreMemoryCard` is the fastest destination; `StoreHostAndCard` is the
slowest, since every frame also crosses USB.

## Live view

**Unverified.** There is no live view operation — the preview is a normal PTP
object at handle `0xFFFFC002`, read with `GetObjectInfo`/`GetObject`. 

    pkill -9 ptpcamerad; ./sonyprobe -liveview 20 -liveview-dir /tmp/lv

Two properties bear on readiness and mean different things: `LiveViewStatus`
(`0xD221`) is whether the preview is up, `MonitoringIsDelivering` (`0xE099`) is
whether frames are actually flowing — the latter is what the SDK's own poll
loop waits on. `StopLiveView` is a deliberate no-op: nothing turns the preview
off, because nothing turned it on.

The payload is not bare JPEG; a header of unpinned length precedes it, so
`LiveFrame` locates the SOI marker rather than assuming an offset, and trims the
block padding off the end.

## Clock

    cam.SyncClockAtMinute()   // safe form
    cam.SyncClock()           // may be up to 59s slow if the body truncates

The clock is a string property (`DateTimeSettings`, `0xD223`), written through
the vendor setter. Whether a Sony truncates or rounds the seconds is **not
known** — an X-T5 truncates — so prefer the minute-boundary form when the
timestamps matter.

## Metadata written into the files

    cam.SetPhotographer("M. Furlotti")
    cam.SetCopyright("(c) 2026")
    cam.SetCopyrightInfo(true)      // without this the strings reach nothing

Sony's string surface is six properties and no more — there is **no Artist and
no Comment field**. The credit is `SetPhotographer`. `WriteCopyrightInfo` is the
switch that embeds both strings, and its off value is `1`, not `0`.

`Metadata()` reads all of them in one round trip; the individual accessors each
cost a full property snapshot.

## Two cautions from the decode

- **`Min`/`Max` in `PropTable` are meaningless for array-typed properties**,
  whose valid set comes from the camera at runtime. `FNumber` is a UInt16Array
  and reports 2/2 — a placeholder, not an aperture range. Use the `Enum` the
  camera reports.
- **The A7R V and A7R VI are not a superset relationship.** 12 properties exist
  on the V and were dropped on the VI, so running both cannot assume one
  capability set keyed off the newer body.
