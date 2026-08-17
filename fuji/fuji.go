// Package fuji drives Fujifilm X and GFX bodies over PTP.
//
// Fujifilm's X Series SDK is a thin wrapper over PTP and says so plainly: its
// transport library exports FTL_PTP_OpenSession, FTL_PTP_GetObject,
// FTL_PTP_InitiateCapture and the rest of the standard set by name. What is
// vendor-specific — and what lives here — is the capture gesture, the exposure
// encodings, and handing control back when the host is finished.
//
// Hardware-validated on an X-T5: capture, download, exposure control, and
// teardown.
package fuji

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
	"github.com/mikefsq/ptp/usb"
)

func init() {
	ptp.Register(ptp.Vendor{
		ID:   ptp.Fujifilm,
		Name: "Fujifilm",
		// Bodies report a generic USB product name — an X-T5 calls itself
		// "USB PTP Camera" — so the product ID is the only way to say which
		// model is attached.
		Models: map[uint16]string{
			0x02FC: "X-T5",
		},
	})
}

// StillStore is the virtual RAM store a Fujifilm body delivers captured frames
// through while tethered. It is not an SD card: the storage type reports fixed
// RAM, and the camera names it "Still".
//
// The card is not reachable over PTP at all while tethered, so this is the only
// place a frame can be read from. Frames stay here until deleted, the store
// holds only about five, and the camera will not hand control back while it
// still has any — see Camera.Close.
const (
	StillStore uint32 = 0x10000001
	LiveStore  uint32 = 0x10000002
)

// Card-only shooting is possible, but not the way it looks.
//
// MediaRecord adds a card copy; it never removes the RAM one. Verified on an
// X-T5: with MediaRecord set to RAW, a captured frame STILL appeared in the
// volatile store, and GetStorageIDs reports only two volumes, both fixed RAM
// with access "read and delete, no write" — the card is not addressable at all
// while tethered.
//
// What works instead is to DELETE the frame without downloading it. DeleteObject
// does not require GetObject, so the camera keeps its card copy and the host
// never moves the file. **Confirmed on hardware: five frames shot this way all
// landed on the card.** It costs 25ms against roughly 800ms for the transfer,
// and it also removes most of the settle wait, because that existed only
// because the newest frame is not READABLE for several seconds — discarding
// does not read it. The capture cycle drops from about 10.9s to 2.35s.
//
// The trade is that nothing on the card can be checked from the host. Retrieval
// afterwards means the body's card-reader USB mode, which cannot shoot, or
// pulling the card.
//
// MediaRecord adds a card copy; it never removes the RAM one. Verified on an
// X-T5: with MediaRecord set to RAW, a captured frame STILL appeared in the
// volatile store. The SDK says the same in its storage model — every image
// taken while tethered goes to in-camera volatile storage, and the card is not
// reachable through the interface at all. The two properties that might have
// gated the transfer, 0xD186 and 0xD187, are not advertised by this body.
//
// Either way the volatile buffer must be drained: it holds only about five
// frames, and while it is occupied the camera refuses to hand control back.

// Camera is an open Fujifilm body. It embeds the PTP session, so every standard
// operation — object handles, downloads, property descriptors — is available
// directly; the methods defined in this package are the ones Fujifilm does
// differently.
type Camera struct {
	*ptp.Session
	t ptp.Transport

	// seen remembers which frame handles NewFrames has already reported, so a
	// caller polling in a loop is told about each frame once.
	seen map[uint32]bool
}

// Open opens a Fujifilm body by USB serial. An empty serial opens the only
// attached Fujifilm camera.
//
// The body must be in tethered shooting mode (MENU, CONNECTION SETTING, USB,
// USB TETHER SHOOTING AUTO). In a file-transfer mode it enumerates and opens
// perfectly well but reports four device properties and no capture support,
// which looks like a driver fault and is not.
func Open(serial string) (*Camera, error) {
	t, err := usb.OpenVendor(ptp.Fujifilm, serial)
	if err != nil {
		return nil, err
	}
	return New(t)
}

// New wraps an already-open transport, for tests and for callers that do their
// own device selection.
func New(t ptp.Transport) (*Camera, error) {
	c := &Camera{Session: ptp.NewSession(t), t: t, seen: map[uint32]bool{}}
	c.Session.Teardown = c.handBack
	c.Session.ResponseNames = ResponseName
	if err := c.Session.Open(); err != nil {
		t.Close()
		return nil, err
	}
	return c, nil
}

// Close ends the session, hands the camera back to its owner and releases the
// USB device.
//
// The handback is not optional politeness. While the host holds PC Priority the
// body's dials, buttons and shutter release are all dead — everything except
// the lens manual focus ring — and closing the PTP session does not undo it. A
// host that simply exits leaves the camera inert in its owner's hands.
func (c *Camera) Close() error {
	err := c.Session.Close()
	if c.t != nil {
		c.t.Close()
		c.t = nil
	}
	return err
}

// handBack runs as the session's teardown hook, with the session lock held, so
// it is given a ptp.Tx rather than calling back through Camera.
//
// Errors are ignored throughout: this runs on the way out, and a camera that
// has already stopped answering cannot be handed anything.
func (c *Camera) handBack(tx ptp.Tx) {
	set := func(p ptp.Prop, v uint64) error {
		payload, err := ptp.EncodeValue(ptp.TypeUint16, v)
		if err != nil {
			return err
		}
		_, _, err = tx(ptp.OpSetDevicePropValue, []uint32{uint32(p)}, payload, ptp.DefaultTimeout)
		return err
	}

	// Let S1 up first. Handing control back while the shutter is still held
	// would leave the body half-pressed with no host left to release it.
	if err := set(PropReleaseGesture, ReleaseNS1Off); err == nil {
		tx(ptp.OpInitiateCapture, []uint32{0, 0}, nil, ptp.DefaultTimeout)
	}

	// The camera REFUSES to give control back while it is still busy, answering
	// vendor error 0xA002 — which here means "not yet", not "bad value".
	// Fujifilm's own sample loops on this, commenting "waiting for the in-camera
	// buffer empty": frames left undeleted in the volatile store are exactly
	// what keeps it busy. A single write leaves the body locked out.
	//
	// Bounded, unlike the SDK's unbounded loop: a camera that never frees up
	// must not hang the caller's Close forever.
	deadline := time.Now().Add(HandBackTimeout)
	for {
		if err := set(PropPriorityModeCode, PriorityModeCamera); err == nil {
			return
		}
		if time.Now().After(deadline) {
			return
		}
		// Slower than the SDK's 50ms: it refuses because it is busy, so
		// hammering it competes with the very work being waited on.
		time.Sleep(250 * time.Millisecond)
	}
}

// HandBackTimeout bounds how long Close waits for the camera to accept control
// back. It only accepts once its internal buffer has drained, which takes as
// long as the last frame needs to finish.
var HandBackTimeout = 10 * time.Second

// Info reports the USB device this camera is bound to, when the transport
// knows. It is the only stable per-body identifier across replugs.
func (c *Camera) Info() (usb.DeviceInfo, bool) {
	if d, ok := c.t.(interface{ Info() usb.DeviceInfo }); ok {
		return d.Info(), true
	}
	return usb.DeviceInfo{}, false
}

// PropName returns a Fujifilm property's name.
//
// This is a package-level function rather than a method on ptp.Prop because
// vendor property codes collide: 0xD001 is film simulation here and something
// else entirely on a Sony body, so only code that knows which camera is
// attached can name them.
func PropName(p ptp.Prop) string {
	// Names established by experiment first: where they exist, the plugin had
	// no accessor to recover one from, so there is nothing to conflict with.
	if n, ok := observedNames[p]; ok {
		return n
	}
	// The body's own plugin first: it is the only authoritative source for this
	// model, and where it disagrees with gphoto2 it is right. gphoto2 names
	// 0xD028 CommandDialMode; the X-T5 itself calls it DOFScale.
	if n, ok := xt5Names[p]; ok {
		return n
	}
	// gphoto2 next — generic across Fujifilm bodies, and it covers codes the
	// X-T5 plugin has no function for.
	if n, ok := propNames[p]; ok {
		return n
	}
	return p.String() // standard PTP names, then the raw code
}

// String describes the camera for logs and errors.
func (c *Camera) String() string {
	if d, ok := c.Info(); ok {
		return d.String()
	}
	return fmt.Sprintf("Fujifilm camera")
}

// Fujifilm vendor response codes. The numeric space is shared across vendors
// and the meanings are not, so only this package can name these.
var responseNames = map[ptp.ResponseCode]string{
	// Seen refusing SetPriorityMode(PC) persistently, while the property's own
	// descriptor said the value was allowed — so it is a refusal about the
	// camera's STATE, not about the value. The exact precondition is not known;
	// the SDK only says these calls need "State S3". Both this driver and the
	// one it replaced get it identically, so it is the body, not the host.
	0xA001: "RefusedInThisState",

	// Not "invalid value". The camera returns this when it will not do
	// something YET — it is still metering, focusing, or holding frames in its
	// volatile store. Fujifilm's own sample loops on it. Reading it as a hard
	// error cost most of a bring-up session.
	0xA002: "RefusedRightNow",
}

// ResponseName names a Fujifilm vendor response code.
func ResponseName(c ptp.ResponseCode) (string, bool) {
	n, ok := responseNames[c]
	return n, ok
}
