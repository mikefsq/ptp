package sony

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Sony has NO InitiateCapture.
//
// The shutter is two virtual buttons pressed through vendor operation 0x9207,
// and they STAY DOWN until released — the camera has no timeout. Down and up
// are distinct calls:
//
//	S1 (half press)  autofocus and metering
//	S2 (full press)  the shutter itself
//
// A caller that presses must release, including on its own error paths.
// Camera.Close is the safety net, and it only releases what this session
// pressed.

// PressButton sends a button-type control down or up (0x9207). It is the
// primitive the shutter buttons are built from; callers usually want the named
// methods below.
func (c *Camera) PressButton(ctrl ControlCode, down bool) error {
	v := ButtonUp
	if down {
		v = ButtonDown
	}
	err := c.SendControl(ctrl, ptp.TypeUint16, v)
	// Remember what is held so Close can let it up. A press that reported an
	// error may still have landed, so a failed press counts as held: releasing
	// something already up is harmless, leaving a shutter held is not.
	c.noteButton(ctrl, down || err != nil)
	return err
}

// HalfPress holds S1 down: autofocus and metering, as a half-press of the
// shutter does on the body. It stays held until HalfRelease.
func (c *Camera) HalfPress() error { return c.PressButton(CtrlS1Button, true) }

// HalfRelease lets S1 up.
func (c *Camera) HalfRelease() error { return c.PressButton(CtrlS1Button, false) }

// FullPress holds S2 down: this fires the shutter. It stays held until
// FullRelease, which in a continuous drive mode means the camera keeps shooting.
func (c *Camera) FullPress() error { return c.PressButton(CtrlS2Button, true) }

// FullRelease lets S2 up.
func (c *Camera) FullRelease() error { return c.PressButton(CtrlS2Button, false) }

// ReleaseAll lets both shutter buttons up, whatever state they were in. Worth
// deferring around any sequence that presses: a body left with S2 held keeps
// firing in continuous drive.
func (c *Camera) ReleaseAll() error {
	// S2 first: releasing S1 while S2 is held is not a state the body expects.
	errS2 := c.FullRelease()
	errS1 := c.HalfRelease()
	if errS2 != nil {
		return errS2
	}
	return errS1
}

// Shoot takes one photograph: S2 down, then S2 up.
//
// It does NOT autofocus first. For an eclipse or observatory sequence that is
// what you want, since focus was set deliberately and an AF hunt would ruin the
// frame — and on a dark or filtered subject the hunt may never finish. Use
// ShootWithAF where autofocus is wanted.
//
// Returning does not mean the image is ready: the camera writes to card and
// announces the object separately. Use WaitForCapture, or poll the card.
func (c *Camera) Shoot() error {
	if err := c.FullPress(); err != nil {
		// Release anyway: a press that reported an error may still have landed,
		// and a held S2 is worse than a spurious release.
		c.FullRelease()
		return fmt.Errorf("sony: shutter press: %w", err)
	}
	if err := c.FullRelease(); err != nil {
		return fmt.Errorf("sony: shutter release: %w", err)
	}
	return nil
}

// Capture takes one frame and waits for the camera to announce it, satisfying
// ptp.Capturer.
//
// The timeout must exceed the exposure. A body that does not raise ObjectAdded
// — which is normal when shooting to card — returns ptp.ErrTimeout here even
// though the frame was taken; use Shoot directly in that case.
func (c *Camera) Capture(timeout time.Duration) error {
	if err := c.Shoot(); err != nil {
		return err
	}
	_, err := c.WaitForCapture(timeout)
	return err
}

// ShootWithAF half-presses for autofocus, waits settle for it to lock, takes
// the shot, then releases both buttons. A zero settle fires immediately, which
// will usually beat the AF.
func (c *Camera) ShootWithAF(settle time.Duration) error {
	if err := c.HalfPress(); err != nil {
		return fmt.Errorf("sony: half press: %w", err)
	}
	defer c.HalfRelease()
	if settle > 0 {
		time.Sleep(settle)
	}
	return c.Shoot()
}

// ShootS1AndS2 fires Sony's combined button, which presses both at once — the
// body's own "focus and shoot" action, in one round trip rather than four.
func (c *Camera) ShootS1AndS2() error {
	if err := c.PressButton(CtrlS1AndS2Button, true); err != nil {
		c.PressButton(CtrlS1AndS2Button, false)
		return fmt.Errorf("sony: S1+S2 press: %w", err)
	}
	if err := c.PressButton(CtrlS1AndS2Button, false); err != nil {
		return fmt.Errorf("sony: S1+S2 release: %w", err)
	}
	return nil
}

// BulbCapture holds the shutter open for d, using the S2 button rather than
// PTP's open-capture operation.
//
// The camera must already be in Bulb exposure mode; this does not set it. Call
// SetBulb first. The duration is host-timed, so it is only as accurate as the
// USB round trips at each end — fine for minutes-long exposures, poor for short
// ones, where a discrete shutter speed is both easier and more accurate.
func (c *Camera) BulbCapture(d time.Duration) error {
	if err := c.FullPress(); err != nil {
		c.FullRelease()
		return fmt.Errorf("sony: bulb press: %w", err)
	}
	time.Sleep(d)
	if err := c.FullRelease(); err != nil {
		return fmt.Errorf("sony: bulb release after %v: %w", d, err)
	}
	return nil
}

// StartMovie starts movie recording.
func (c *Camera) StartMovie() error { return c.PressButton(CtrlMovieRecButton, true) }

// StopMovie stops movie recording.
func (c *Camera) StopMovie() error { return c.PressButton(CtrlMovieRecButton, false) }

// AELock holds or releases the AE lock button.
func (c *Camera) AELock(down bool) error { return c.PressButton(CtrlAELButton, down) }

// AWBLock holds or releases the AWB lock button.
func (c *Camera) AWBLock(down bool) error { return c.PressButton(CtrlAWBLButton, down) }

// WaitForCapture waits for the camera to announce a new object, returning its
// handle. It exists because Shoot returns as soon as the button is released,
// long before the image is written.
//
// Returns ptp.ErrTimeout if nothing arrives, which on a body that does not
// raise ObjectAdded means you must poll the card instead.
func (c *Camera) WaitForCapture(timeout time.Duration) (uint32, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ev, err := c.PollEvent(500 * time.Millisecond)
		if err != nil {
			if err == ptp.ErrTimeout {
				continue
			}
			return 0, err
		}
		switch ev.Code {
		case ptp.EventObjectAdded, ptp.EventSonyObjectAdded:
			if len(ev.Params) > 0 {
				return ev.Params[0], nil
			}
			return 0, nil
		case ptp.EventCaptureComplete:
			// Capture finished, but the object handle comes separately; keep
			// waiting for the ObjectAdded that follows.
			continue
		}
	}
	return 0, ptp.ErrTimeout
}

// NewFrames lists frames that have appeared since the last call, satisfying
// ptp.Downloader.
func (c *Camera) NewFrames() ([]uint32, error) {
	handles, err := c.GetObjectHandles(0xFFFFFFFF, 0, 0)
	if err != nil {
		return nil, err
	}
	var fresh []uint32
	for _, h := range handles {
		if !c.seen[h] {
			c.seen[h] = true
			fresh = append(fresh, h)
		}
	}
	return fresh, nil
}

// Download fetches one frame and its filename.
func (c *Camera) Download(handle uint32) ([]byte, string, error) {
	oi, err := c.GetObjectInfo(handle)
	if err != nil {
		return nil, "", err
	}
	data, err := c.GetObject(handle)
	if err != nil {
		return nil, oi.Filename, err
	}
	return data, oi.Filename, nil
}

// Delete removes a frame from the camera.
func (c *Camera) Delete(handle uint32) error { return c.DeleteObject(handle, 0) }
