# ptp

USB PTP Camera control implemented in Go.

```go
import (
    "github.com/mikefsq/ptp"
    "github.com/mikefsq/ptp/fuji"
    //"github.com/mikefsq/ptp/sony"
)

cam, err := fuji.Open("")        // empty serial: the only attached body
defer cam.Close()                // hands the camera back to its owner

cam.SetManualFocus()
cam.SetShutter(time.Second / 1000)
cam.Capture(120 * time.Second)
```

## Layout

| | |
|---|---|
| `ptp` | framing, sessions, objects, property descriptors, events. |
| `ptp/usb` | the per-OS USB transports, enumeration with a per-attachment identity, and `Hotplug` (attach and detach notifications: IOKit on macOS, uevent netlink on Linux). |
| `ptp/fuji` | Fujifilm X and GFX |
| `ptp/sony` | Sony Alpha |
| `ptp/canon`, `ptp/nikon` | not implemented yet |

## Status

**`ptp/fuji`** is validated with an X-T5.

**`ptp/sony`** The control layer is implemented but untested. 

## RAW decoding

These vendor libs can decode and return the raw sensor readout as an undemosaiced
`ptp.CFA` array.

## Extending it

This module provides a standardized interface to PTP cameras. Since there are many vendor
differences, the vendor package supplies its own tables, capture semantics and response-code
names. 

See [fuji/OPERATION.md](fuji/OPERATION.md) and [sony/OPERATION.md](sony/OPERATION.md)
for the practical side of driving a camera body: which dial positions the host needs,
the two shooting modes and their measured rates, what to do when the camera stops 
answering, and what does not work. 

## Building

    make            build every package
    make help       list the targets
    make tools      build the two command-line tools into ./bin
    make test       go test ./...
    make cross      check it still builds for linux and windows

Everything builds on a clean checkout. `go.mod` requires nothing, no vendor SDK
is linked on any platform, and the protocol and camera packages — `ptp`,
`ptp/fuji`, `ptp/sony` — are pure Go, so their tests run and cross-compile
anywhere.

The one exception is the USB transport. `ptp/usb` on macOS is cgo against
IOKit; on Linux it drives usbfs ioctls directly and is pure Go; Windows is a
stub. That is the only cgo in the module.

## Tools

    cmd/fujiprobe    bring-up and diagnosis: capture, download, live view,
                     property get/set, burst timing, snapshot/diff
    cmd/sonyprobe    the same for Sony: handshake, property dump, capture,
                     exposure, live view, browse, -reset. Writes a capture
                     bundle every run, and -replay re-parses one offline with
                     no camera attached

`fujiprobe -snapshot F` then `-diff F`, with one control changed in between,
identifies which property a camera control owns.

## License

MIT — see [LICENSE](LICENSE).
