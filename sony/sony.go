package sony

import (
	"fmt"
	"sync"

	"github.com/mikefsq/ptp"
	"github.com/mikefsq/ptp/usb"
)

func init() {
	ptp.Register(ptp.Vendor{
		ID:   ptp.Sony,
		Name: "Sony",
		// A body changes its product ID with its USB mode, so only the PIDs that
		// can actually be driven belong here. A NEX-6 reports 0x0677 as mass
		// storage and 0x05B2 while merely charging; neither can be controlled.
		Models: map[uint16]string{
			0x0678: "Alpha (still image)",
			0x0C34: "Alpha (PC Remote)",
		},
	})
}

// Camera is an open Sony body. It embeds the PTP session, so every standard
// operation — object handles, downloads, storage info — is available directly;
// the methods defined in this package are the ones Sony does differently, which
// on this vendor is most of the interesting surface.
type Camera struct {
	*ptp.Session
	t ptp.Transport

	// Ext is what the body said it supports during the vendor handshake. It is
	// nil if the handshake was skipped.
	Ext *ExtDeviceInfo

	// mu guards held. The teardown hook runs with the session lock held, not
	// this one, so the two locks never nest.
	mu sync.Mutex

	// held tracks which control buttons this session has pressed and not
	// released. Sony's shutter is a pair of virtual buttons that STAY DOWN
	// until released — the camera has no timeout — so a host that exits
	// mid-capture leaves the shutter held and, in a continuous drive mode, the
	// camera still shooting.
	held map[ControlCode]bool

	// seen remembers which frame handles NewFrames has already reported.
	seen map[uint32]bool
}

// Open opens a Sony body by USB serial and performs the vendor handshake. An
// empty serial opens the only attached Sony camera.
//
// The body must be in PC Remote USB mode. In a file-transfer mode it enumerates
// and opens perfectly well but exposes no vendor operations at all, and every
// 0x92xx call then refuses by stalling the bulk pipe.
func Open(serial string) (*Camera, error) {
	t, err := usb.OpenVendor(ptp.Sony, serial)
	if err != nil {
		return nil, err
	}
	c, err := New(t)
	if err != nil {
		return nil, err
	}
	if _, err := c.Connect(); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

// New wraps an already-open transport and opens a PTP session, WITHOUT the
// vendor handshake. For tests, for callers doing their own device selection,
// and for talking to a body that has no vendor surface. Call Connect to unlock
// the 0x92xx operations.
func New(t ptp.Transport) (*Camera, error) {
	c := &Camera{
		Session: ptp.NewSession(t),
		t:       t,
		held:    map[ControlCode]bool{},
		seen:    map[uint32]bool{},
	}
	c.Session.Teardown = c.releaseHeld
	c.Session.ResponseNames = ResponseName
	if err := c.Session.Open(); err != nil {
		t.Close()
		return nil, err
	}
	return c, nil
}

// Close lets go of any button this session is still holding, ends the session
// and releases the USB device.
//
// The release is not politeness. A Sony body models its shutter as a physical
// button with no timeout, so a session that exits with S2 down leaves the
// camera firing.
func (c *Camera) Close() error {
	err := c.Session.Close()
	if c.t != nil {
		c.t.Close()
		c.t = nil
	}
	return err
}

// releaseHeld runs as the session's teardown hook, with the session lock held,
// so it is given a ptp.Tx rather than calling back through Camera.
//
// Only buttons this session actually pressed are released: sending vendor
// operations to a body that never had any is how a driver turns a clean exit
// into a stalled pipe.
//
// Errors are ignored throughout — this runs on the way out, and a camera that
// has already stopped answering cannot be handed anything.
func (c *Camera) releaseHeld(tx ptp.Tx) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.held) == 0 {
		return
	}

	up := func(ctrl ControlCode) {
		payload, err := ptp.EncodeValue(ptp.TypeUint16, ButtonUp)
		if err != nil {
			return
		}
		tx(OpSetControlDeviceB, []uint32{uint32(ctrl)}, payload, ptp.DefaultTimeout)
		delete(c.held, ctrl)
	}

	// S2 before S1: releasing S1 while S2 is held is not a state the body
	// expects. Any other buttons follow.
	for _, ctrl := range []ControlCode{CtrlS2Button, CtrlS1Button, CtrlS1AndS2Button} {
		if c.held[ctrl] {
			up(ctrl)
		}
	}
	for ctrl := range c.held {
		up(ctrl)
	}
}

// noteButton records whether a control button is currently held.
func (c *Camera) noteButton(ctrl ControlCode, down bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if down {
		c.held[ctrl] = true
	} else {
		delete(c.held, ctrl)
	}
}

// Held lists the buttons this session is currently holding down.
func (c *Camera) Held() []ControlCode {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]ControlCode, 0, len(c.held))
	for ctrl := range c.held {
		out = append(out, ctrl)
	}
	return out
}

// Model names the body. The USB product ID is all that is available before the
// vendor handshake, and Sony reuses it across the range, so this identifies the
// interface mode rather than the model.
func (c *Camera) Model() string {
	if d, ok := c.Info(); ok {
		return d.Name
	}
	return "Sony camera"
}

// Info reports the USB device this camera is bound to, when the transport
// knows. It is the only stable per-body identifier across replugs.
func (c *Camera) Info() (usb.DeviceInfo, bool) {
	if d, ok := c.t.(interface{ Info() usb.DeviceInfo }); ok {
		return d.Info(), true
	}
	return usb.DeviceInfo{}, false
}

// String describes the camera for logs and errors.
func (c *Camera) String() string {
	if d, ok := c.Info(); ok {
		return d.String()
	}
	return fmt.Sprintf("Sony camera")
}

// ResponseName names a Sony vendor response code.
//
// Sony adds none of its own: a vendor operation that fails reports a standard
// PTP response code, and a body that does not implement one stalls the pipe
// instead of answering at all. The hook is still installed, because that fact
// is worth stating in one place rather than rediscovering it.
func ResponseName(c ptp.ResponseCode) (string, bool) { return "", false }
