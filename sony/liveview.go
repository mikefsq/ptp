package sony

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Live view.
//
// Sony does NOT use PTP's open-capture operations for this, and it has no
// vendor operation of its own either: the preview is an ordinary PTP object at
// a fixed, magic handle, fetched with the standard GetObjectInfo and GetObject.
//
// That is recovered from Sony's own USB transport adapter, not inferred.
// DeviceConnectionPTP::GetLiveViewImage issues operation 0x1009 with one
// parameter, and the instruction that builds that parameter is:
//
//	mov  w8, #-0x3ffe     ; = 0xFFFFC002 as a uint32
//	str  w8, [sp, #0x14]
//	add  x2, sp, #0x14    ; -> the parameter array
//
// DeviceConnectionPTP::GetLiveViewImageSize does the same with 0x1008 and a
// 116-byte ObjectInfo buffer. gphoto2 uses the same handle, which is a second
// source agreeing with Sony's binary.
//
// UNVERIFIED on hardware. The handle and the operations are pinned; what the
// returned payload looks like byte-for-byte is not, because no body supporting
// live view has been driven. See LiveFrame for how that uncertainty is handled.
const (
	// LiveViewObject is the magic handle the preview frame lives at.
	LiveViewObject uint32 = 0xFFFFC002

	// Two neighbouring handles appear in the same adapter binary and are not
	// decoded. They are named so a probe can try them rather than rediscover
	// that they exist.
	LiveViewObjectAlt1 uint32 = 0xFFFFC001
	LiveViewObjectAlt2 uint32 = 0xFFFFC003
)

// LiveViewStatus values are not enumerated by the SDK headers; treat non-zero
// as running.
const LiveViewRunning uint64 = 1

// LiveViewAvailable reports whether the camera says live view is running.
//
// Two properties bear on it and they mean different things: LiveViewStatus
// (0xD221) is whether the body has live view up at all, and
// MonitoringIsDelivering (0xE099) is whether frames are actually flowing —
// which is what the SDK's own polling loop watches before it fetches. A body
// can report the first without the second while it is still starting up.
func (c *Camera) LiveViewAvailable() (running, delivering bool, err error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return false, false, err
	}
	if p := FindProp(props, PropLiveViewStatus); p != nil {
		running = p.Current != 0
	}
	if p := FindProp(props, PropMonitoringIsDelivering); p != nil {
		delivering = p.Current != 0
	}
	return running, delivering, nil
}

// StartLiveView asks the camera to begin delivering preview frames and waits
// for it to say it is.
//
// There is no "start live view" operation. Sony bodies bring the preview up as
// a side effect of being in a shooting state, so this polls rather than
// commands — and reports what it saw if the body never starts, because "no
// frames" and "camera refused" need different fixes.
func (c *Camera) StartLiveView(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var running bool
	for {
		r, delivering, err := c.LiveViewAvailable()
		if err != nil {
			return err
		}
		running = r
		if delivering {
			return nil
		}
		if time.Now().After(deadline) {
			if running {
				return fmt.Errorf("sony: live view is up but the camera is not delivering frames after %v; "+
					"the body may be in playback or a menu", timeout)
			}
			return fmt.Errorf("sony: camera did not start live view within %v; "+
				"check it is in a shooting mode and not in playback", timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// StopLiveView is a no-op that exists to keep the pairing obvious.
//
// Nothing in the SDK or the adapter stops the preview explicitly: it is a
// property of the camera's state, not something the host switched on, so there
// is nothing to switch off. Provided so callers can write the symmetric
// start/stop they expect and so this note has somewhere to live.
func (c *Camera) StopLiveView() error { return nil }

// LiveFrame fetches one preview frame as JPEG bytes.
//
// The payload is NOT bare JPEG. Sony prefixes the object with a header whose
// size is not pinned — the adapter reads an ObjectInfo first and the SDK copies
// out of a cached block, so neither reveals the layout. Rather than guess an
// offset, this finds the JPEG start-of-image marker and returns from there,
// which is correct whatever the header length turns out to be.
//
// It returns an error rather than a truncated image if no SOI is present, since
// a body that answers with something else entirely should not be reported as a
// successful frame.
func (c *Camera) LiveFrame() ([]byte, error) {
	data, err := c.GetObject(LiveViewObject)
	if err != nil {
		return nil, fmt.Errorf("sony: fetching live view object 0x%08X: %w", LiveViewObject, err)
	}
	return trimToJPEG(data)
}

// LiveFrameInfo reads the preview object's ObjectInfo, which carries its
// dimensions and size without transferring the frame.
func (c *Camera) LiveFrameInfo() (*ptp.ObjectInfo, error) {
	return c.GetObjectInfo(LiveViewObject)
}

// jpegSOI and jpegEOI bracket a JPEG.
var (
	jpegSOI = []byte{0xFF, 0xD8}
	jpegEOI = []byte{0xFF, 0xD9}
)

// trimToJPEG finds the JPEG inside a live view payload.
//
// It also trims after the end-of-image marker. A Sony pads object transfers to
// a block boundary — confirmed on a NEX-6, where a 1,611,005-byte JPEG arrived
// declared as 1,638,400 — so trailing bytes are expected, not a fault.
func trimToJPEG(b []byte) ([]byte, error) {
	start := bytes.Index(b, jpegSOI)
	if start < 0 {
		return nil, fmt.Errorf("sony: live view payload of %d bytes contains no JPEG start marker; "+
			"the frame format is not what this driver expects", len(b))
	}
	img := b[start:]
	if end := bytes.LastIndex(img, jpegEOI); end >= 0 {
		img = img[:end+2]
	}
	return img, nil
}
