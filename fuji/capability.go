package fuji

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Camera satisfies the shared capability interfaces, so a client can drive a
// Fujifilm body through the same code as any other vendor's.
//
// These assertions are the point: if a core interface drifts away from what a
// real driver can honour, this file stops compiling. That is cheaper than
// discovering it when a second vendor is written.
var (
	_ ptp.Camera          = (*Camera)(nil)
	_ ptp.Capturer        = (*Camera)(nil)
	_ ptp.ExposureControl = (*Camera)(nil)
	_ ptp.Downloader      = (*Camera)(nil)
	_ ptp.BufferReporter  = (*Camera)(nil)
	_ ptp.FocusControl    = (*Camera)(nil)
	_ ptp.LiveViewer      = (*Camera)(nil)
)

// Model names the body.
func (c *Camera) Model() string {
	if d, ok := c.Info(); ok {
		return d.Model()
	}
	return "Fujifilm camera"
}

// Capture takes one frame. Focus behaviour comes from the body's focus mode,
// not from here — see SetManualFocus.
func (c *Camera) Capture(timeout time.Duration) error {
	return c.CaptureFrame(false, timeout)
}

// SetManualFocus stops the mandatory half press from hunting.
func (c *Camera) SetManualFocus() error { return c.SetFocusMode(FocusManual) }

// SetShutter snaps to the nearest speed the camera offers.
//
// The descriptor is read here rather than asked of the caller: the offered set
// changes with the shutter dial's position, so a cached one goes stale the
// moment somebody touches the body.
func (c *Camera) SetShutter(d time.Duration) error {
	desc, err := c.GetPropDesc(ptp.PropExposureTime)
	if err != nil {
		return fmt.Errorf("fuji: reading the shutter speed range: %w", err)
	}
	return c.SetShutterFrom(desc, d)
}

// Exposure reports the current triangle and what the host may write.
func (c *Camera) Exposure() (*ptp.Exposure, error) {
	e, err := c.ReadExposure()
	if err != nil {
		return nil, err
	}
	return &ptp.Exposure{
		Shutter: e.Shutter, Aperture: e.Aperture, ISO: e.ISO,
		ShutterSettable:  e.ShutterSettable,
		ApertureSettable: e.ApertureSettable,
		ISOSettable:      e.ISOSettable,
	}, nil
}

// NewFrames lists frames that have appeared in the volatile store since the
// last call.
func (c *Camera) NewFrames() ([]uint32, error) {
	handles, err := c.GetObjectHandles(StillStore, 0, 0)
	if err != nil {
		return nil, err
	}
	var fresh []uint32
	for _, h := range handles {
		if !c.seen[h] {
			fresh = append(fresh, h)
		}
	}
	if c.seen == nil {
		c.seen = map[uint32]bool{}
	}
	for _, h := range handles {
		c.seen[h] = true
	}
	return fresh, nil
}

// Download fetches one frame and its filename, checking the length against
// what the camera said it would send.
func (c *Camera) Download(handle uint32) ([]byte, string, error) {
	oi, err := c.GetObjectInfo(handle)
	if err != nil {
		return nil, "", err
	}
	data, err := c.GetObject(handle)
	if err != nil {
		return nil, "", err
	}
	if uint32(len(data)) != oi.CompressedSize {
		return data, oi.Filename, fmt.Errorf(
			"fuji: got %d bytes for %s, ObjectInfo said %d",
			len(data), oi.Filename, oi.CompressedSize)
	}
	return data, oi.Filename, nil
}

// Delete removes a frame from the volatile store.
//
// On a tethered Fujifilm body this discards the camera's ONLY copy unless
// MediaRecord is also writing to the card, so call it after checking a download
// against ObjectInfo's size — which Download does.
func (c *Camera) Delete(handle uint32) error {
	if err := c.DeleteObject(handle, 0); err != nil {
		return err
	}
	delete(c.seen, handle)
	return nil
}
