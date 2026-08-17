package sony

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Camera satisfies the parent package's capability interfaces. These assertions
// are compile-time: if an interface drifts away from what this driver can
// honour, the build breaks here rather than at a caller's type assertion.
var (
	_ ptp.Camera          = (*Camera)(nil)
	_ ptp.Capturer        = (*Camera)(nil)
	_ ptp.ExposureControl = (*Camera)(nil)
	_ ptp.Downloader      = (*Camera)(nil)
	_ ptp.FocusControl    = (*Camera)(nil)
	_ ptp.LiveViewer      = (*Camera)(nil)
)

// Focus modes (CrFocusMode).
const (
	FocusManual uint64 = 0x0001
	FocusAFS    uint64 = 0x0002
	FocusAFC    uint64 = 0x0003
	FocusAFA    uint64 = 0x0004
	FocusAFD    uint64 = 0x0005
	FocusDMF    uint64 = 0x0006
	FocusPF     uint64 = 0x0007
)

// Exposure programs (CrExposureProgram). Only the four that matter to a host
// are named; the body offers a long tail of scene modes that give it control.
const (
	ExposureManual           uint64 = 0x00000001
	ExposureProgramAuto      uint64 = 0x00000002
	ExposureAperturePriority uint64 = 0x00000003
	ExposureShutterPriority  uint64 = 0x00000004
)

// SetFocusMode sets the focus mode, e.g. FocusManual.
func (c *Camera) SetFocusMode(mode uint64) error {
	return c.SetProperty(PropFocusMode, ptp.TypeUint16, mode)
}

// SetManualFocus puts the body in manual focus, satisfying ptp.FocusControl.
//
// Worth doing before any unattended sequence. Autofocus on a subject the lens
// cannot lock — a dark sky, a filtered sun — hunts indefinitely, and a
// half-press that never resolves is how a tethered body stops answering.
func (c *Camera) SetManualFocus() error { return c.SetFocusMode(FocusManual) }

// SetExposureProgram sets the exposure program, e.g. ExposureManual.
//
// This governs which of the exposure triangle the host may set at all: in P the
// camera owns shutter and aperture, and writes to them are refused or ignored.
func (c *Camera) SetExposureProgram(mode uint64) error {
	return c.SetProperty(PropExposureProgramMode, ptp.TypeUint32, mode)
}

// SetShutter sets the exposure time, snapping to the nearest speed the camera
// currently advertises. It satisfies ptp.ExposureControl.
//
// This costs a full property snapshot. In a burst, take one snapshot and use
// SetShutterFrom rather than paying for one per frame.
func (c *Camera) SetShutter(d time.Duration) error {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return err
	}
	return c.SetShutterFrom(props, d)
}

// ExposureState reads the exposure triangle as Sony reports it — the vendor
// view, with each component's own "is it present" flag.
//
// Exposure is the portable view of the same data.
func (c *Camera) ExposureState() (Exposure, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return Exposure{}, err
	}
	return ReadExposure(props), nil
}

// Exposure reports the exposure triangle in the parent package's shape,
// satisfying ptp.ExposureControl.
//
// The Settable flags come from the camera's own enable byte, so a false means
// the body currently owns that component — usually because the mode dial is not
// on M, or a scene mode has taken it. Writing it anyway is commonly ACCEPTED
// and then ignored, which is why this is worth checking rather than inferring
// from a write's response.
func (c *Camera) Exposure() (*ptp.Exposure, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return nil, err
	}
	state := ReadExposure(props)
	out := &ptp.Exposure{
		Shutter:  state.ShutterSpeed,
		Aperture: state.Aperture,
		ISO:      uint64(state.ISO),
	}
	settable := func(code Prop) bool {
		p := FindProp(props, code)
		return p != nil && p.Writable()
	}
	out.ShutterSettable = settable(PropShutterSpeed)
	out.ApertureSettable = settable(PropFNumber)
	out.ISOSettable = settable(PropIsoSensitivity)
	return out, nil
}

// Ready reports whether the camera can be driven right now, and says what is
// wrong when it cannot.
//
// The three things that stop a host shooting are different problems needing
// different fixes, so they are distinguished rather than collapsed into one
// error: no vendor handshake, no usable card, or an exposure the host does not
// control.
func (c *Camera) Ready() error {
	if c.Ext == nil {
		return fmt.Errorf("sony: vendor handshake has not run; call Connect (the body must be in PC Remote USB mode)")
	}
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return err
	}
	if r := ReadReadiness(props); !r.Ready() {
		return fmt.Errorf("sony: camera is not ready to shoot: %s", r)
	}
	return nil
}
