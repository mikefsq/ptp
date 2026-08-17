# Operating a Fujifilm body over PTP

Practical notes from bringing up an X-T5. Everything here was established on
hardware; where something is inferred or unverified it says so.

## Before connecting

Set the USB mode to **tethered shooting** (`NETWORK/USB SETTING` → `SELECT
CONNECTION SETTING` → `USB TETHER SHOOTING AUTO`). In a file-transfer mode the
body enumerates and opens perfectly well but reports 4 properties instead of
263 and `InitiateCapture=false`, which looks like a driver fault and is not.

## Dial positions

The host can only set what the body has not claimed. Check with
`fujiprobe -exposure`; anything marked `(locked)` is the camera's.

| control | for host control | what it costs |
|---|---|---|
| shutter dial | **T** | `B` gives the camera's own bulb timer instead; a marked speed locks it |
| aperture ring | **A**, then force `-mode manual` | left off `A` you get Manual, but the ring owns the aperture |
| ISO dial | **C** | `A` locks it to auto1 |
| focus lever | **M** | required — see below |
| drive dial | **Single** | continuous gives no burst over PTP; HDR reserves buffer |

**Manual focus is not optional.** The capture sequence needs a half press, and
in an AF mode that starts a focus hunt. On a subject the lens cannot lock —
a dark sky, a filtered sun — the hunt never finishes, the camera stops
answering, and it stays wedged across process restarts. MF also keeps
`FocusPos` readable, so a focus position can be recorded and restored.

`PREVIEW EXP./WB IN MANUAL MODE` stops the lens down to show the true exposure.
At a small aperture on a dark subject that leaves the finder black and looks
like a fault. `SetExposurePreview(ExposurePreviewOff)`.

## Shooting

Two modes work. Card-only is roughly four times faster.

    # card-only: the camera keeps the frame, the host never transfers it
    fujiprobe -card raw -discard -burst 20 -bracket 1ms,4ms,15ms,60ms,250ms

    # USB: the frame comes to the host and nothing is left on the card
    fujiprobe -card off -capture -get ./frames -drain

Measured, card-only: **0.65 s/frame**, or **0.76 s/frame** bracketed. Per frame
that is ~0.66 s fixed, plus the exposure, plus 60 ms for a shutter change. A
25 ms delete replaces an ~800 ms transfer.

Downloading cannot keep pace with a burst, at any lag. A frame needs several
seconds before it is readable and a burst produces one every 0.65 s, so keeping
up would need ~15 frames of lag against a buffer that holds five. Shoot
card-only and retrieve afterwards.

Retrieval afterwards: switch the body to **card-reader USB mode**, where the SD
card appears as an ordinary removable volume. Capture is unavailable in that
mode, so it is an after-the-fact step.

## When it goes wrong

`Camera.Ready()` distinguishes the three causes, which need three different
fixes:

| symptom | cause | fix |
|---|---|---|
| `PriorityMode` unreadable, `0xA001` | body in playback or a menu | half-press the shutter |
| priority write refused, `0xA002` | undownloaded frames in the buffer | `-get DIR -drain` (needs no priority) |
| `DeviceBusy` on the full press | shutter left half-pressed by a failed capture | `-release` |
| transport timeouts | asleep, or unplugged | wake it |
| `AFStatus` stuck at 4 | a latch PTP cannot clear | **power cycle** — no gesture clears it |

## Thermal

The body overheated and froze during sustained tethered work, and afterwards
repeated bursts failed in quick succession until it was rested. Ten frames is
the longest clean run recorded here; a multi-minute sequence has never been
tested. Tethering holds the camera awake and draws bus power even between
frames, so it runs warmer than it would shooting alone.

Treat 0.65 s/frame as a measured rate, not a sustained one, until a long run
has been done on a cool body.

## What does not work

- **Bulb.** Both documented gestures fail: `0x0500` is accepted and then the
  capture returns `GeneralError`; `0x0400`, which the SDK's own `StartBulb`
  sends, is refused outright. Use the shutter ladder — it reaches 64 s.
- **Continuous burst.** The body advertises no bare S2 press, only the combined
  press-and-release, so the shutter cannot be held down.
- **Card-only without deleting.** `MediaRecord` adds a card copy; it never
  removes the RAM one. The buffer fills either way.
