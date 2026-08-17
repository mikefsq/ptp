package sony

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Sony extension protocol versions, as reported by 0x9202.
const (
	Protocol200 uint16 = 0x00C8
	Protocol300 uint16 = 0x012C
)

// ExtDeviceInfo is the reply to 0x9202 SDIO_GetExtDeviceInfo: which property
// and control codes this particular body actually supports.
//
// This, not the SDK's documentation matrix, is what the camera in front of you
// offers. Supported reports the former; this reports the latter.
type ExtDeviceInfo struct {
	Version  uint16 // Protocol200 or Protocol300
	Props    []Prop // supported device property codes
	Controls []Prop // supported control codes

	// Raw is the undecoded 0x9202 payload, and Attempts is how many tries it
	// took before the camera returned one. Both exist so a body that answers in
	// an unexpected shape can be diagnosed from a single run.
	Raw      []byte
	Attempts int
}

// ModeVer reports the parser generation the body speaks: 3 for the newer
// protocol, 2 otherwise.
func (e *ExtDeviceInfo) ModeVer() int {
	if e.Version == Protocol300 {
		return 3
	}
	return 2
}

// Connect performs the Sony vendor handshake that unlocks the 0x92xx
// operations. Open calls it; it is exported for callers driving the sequence
// themselves.
//
// The sequence is three phases of 0x9201 with 0x9202 in the middle:
//
//	0x9201(1,0,0)  0x9201(2,0,0)  0x9202  0x9201(3,0,0)
//
// The 0x9202 step is retried: the camera commonly answers with an empty list
// for the first few attempts while it is still bringing the session up, and a
// single try will often come back empty on a cold connect.
func (c *Camera) Connect() (*ExtDeviceInfo, error) {
	// Ask the camera what it supports before issuing a vendor operation. A body
	// that does not implement one refuses by STALLING the bulk pipe, which costs
	// a stall recovery and produces a far less useful error than this check.
	if di, err := c.GetDeviceInfo(); err == nil && !SupportsSDIO(di) {
		return nil, fmt.Errorf("sony: %s %s does not support the Sony remote-control "+
			"operations: its DeviceInfo lists %d operations, without 0x9201/0x9202. "+
			"This body exposes file transfer only over this interface — check it has a "+
			"PC Remote USB mode, and that it is set to it",
			di.Manufacturer, di.Model, len(di.Operations))
	}

	if _, _, err := c.Do(OpSDIOConnect, []uint32{1, 0, 0}, nil, ptp.DefaultTimeout); err != nil {
		return nil, fmt.Errorf("sony: connect phase 1: %w", err)
	}
	if _, _, err := c.Do(OpSDIOConnect, []uint32{2, 0, 0}, nil, ptp.DefaultTimeout); err != nil {
		return nil, fmt.Errorf("sony: connect phase 2: %w", err)
	}

	var info *ExtDeviceInfo
	const tries = 20
	attempts := 0
	for range tries {
		attempts++
		data, _, err := c.Do(OpSDIOGetExtDeviceInfo, []uint32{uint32(Protocol300), 1}, nil, ptp.DefaultTimeout)
		if err != nil {
			return nil, fmt.Errorf("sony: getting extended device info (attempt %d): %w", attempts, err)
		}
		if len(data) > 0 {
			info, err = parseExtDeviceInfo(data)
			if err != nil {
				// Hand back the raw bytes with the error: an unexpected shape
				// here is exactly the case where the payload is the evidence.
				return &ExtDeviceInfo{Raw: data, Attempts: attempts}, err
			}
			info.Raw, info.Attempts = data, attempts
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if info == nil {
		return nil, fmt.Errorf("sony: camera returned an empty extended device info list after %d attempts; "+
			"check the body is in PC Remote mode", tries)
	}

	if _, _, err := c.Do(OpSDIOConnect, []uint32{3, 0, 0}, nil, ptp.DefaultTimeout); err != nil {
		return nil, fmt.Errorf("sony: connect phase 3: %w", err)
	}
	c.Ext = info
	return info, nil
}

// SupportsSDIO reports whether a body advertises the vendor handshake at all.
//
// Worth checking before any 0x92xx call: a body without the vendor surface
// refuses by stalling the bulk pipe, which costs a stall recovery and says far
// less than this does. A NEX-6 is the case in point — it speaks PTP perfectly
// and has none of this.
func SupportsSDIO(di *ptp.DevInfo) bool {
	return di.Supports(OpSDIOConnect) && di.Supports(OpSDIOGetExtDeviceInfo)
}

// parseExtDeviceInfo decodes the 0x9202 payload: a uint16 protocol version,
// then two uint16 arrays (each a uint32 count followed by that many codes) —
// device properties, then controls.
func parseExtDeviceInfo(b []byte) (*ExtDeviceInfo, error) {
	r := ptp.NewReader(b)
	v, err := r.U16()
	if err != nil {
		return nil, fmt.Errorf("sony: reading protocol version: %w", err)
	}
	info := &ExtDeviceInfo{Version: v}
	if info.Props, err = readCodeArray(r); err != nil {
		return nil, fmt.Errorf("sony: reading supported property codes: %w", err)
	}
	if info.Controls, err = readCodeArray(r); err != nil {
		return nil, fmt.Errorf("sony: reading supported control codes: %w", err)
	}
	return info, nil
}

func readCodeArray(r *ptp.Reader) ([]Prop, error) {
	n, err := r.U32()
	if err != nil {
		return nil, err
	}
	if int(n) > r.Remaining()/2 {
		return nil, fmt.Errorf("array of %d codes: %w", n, ptp.ErrShortBlob)
	}
	out := make([]Prop, n)
	for i := range out {
		v, err := r.U16()
		if err != nil {
			return nil, err
		}
		out[i] = Prop(v)
	}
	return out, nil
}

// SetProperty sets a stored setting via 0x9205 SetControlDeviceA. The value is
// encoded in the width t implies.
//
// The camera accepts only values from the set it currently advertises, and that
// set changes with lens, exposure mode and drive mode. A write outside it is
// commonly accepted and then ignored rather than refused, so read back.
func (c *Camera) SetProperty(p Prop, t ptp.DataType, v uint64) error {
	payload, err := ptp.EncodeValue(t, v)
	if err != nil {
		return err
	}
	_, _, err = c.Do(OpSetControlDeviceA, []uint32{uint32(p)}, payload, ptp.DefaultTimeout)
	return err
}

// SendControl issues a momentary action via 0x9207 SetControlDeviceB — shutter
// release, AF trigger and the like, as opposed to a stored setting.
func (c *Camera) SendControl(ctrl ControlCode, t ptp.DataType, v uint64) error {
	payload, err := ptp.EncodeValue(t, v)
	if err != nil {
		return err
	}
	_, _, err = c.Do(OpSetControlDeviceB, []uint32{uint32(ctrl)}, payload, ptp.DefaultTimeout)
	return err
}
