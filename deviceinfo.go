package ptp

import (
	"encoding/binary"
	"fmt"
)

// DevInfo is the standard PTP DeviceInfo dataset (operation 0x1001)
type DevInfo struct {
	StandardVersion    uint16 // in hundredths, so 100 = PTP 1.00
	VendorExtensionID  uint32
	VendorExtensionVer uint16
	VendorExtensionsc  string
	FunctionalMode     uint16

	Operations     []OpCode
	Events         []uint16
	DeviceProps    []Prop
	CaptureFormats []uint16
	ImageFormats   []uint16

	Manufacturer  string
	Model         string
	DeviceVersion string
	SerialNumber  string
}

// Supports reports whether the camera lists op in its DeviceInfo.
func (d *DevInfo) Supports(op OpCode) bool {
	for _, o := range d.Operations {
		if o == op {
			return true
		}
	}
	return false
}

// SupportsCapture reports whether the body exposes the standard PTP capture
// operation.
func (d *DevInfo) SupportsCapture() bool { return d.Supports(OpInitiateCapture) }

// SupportsBulb reports whether the body offers open-ended (bulb) capture.
func (d *DevInfo) SupportsBulb() bool { return d.Supports(OpInitiateOpenCapture) }

func (d *DevInfo) String() string {
	return fmt.Sprintf("%s %s (firmware %s, serial %s), PTP %.2f, %d operations",
		d.Manufacturer, d.Model, d.DeviceVersion, d.SerialNumber,
		float64(d.StandardVersion)/100, len(d.Operations))
}

// ParseDeviceInfo decodes the payload of operation 0x1001.
func ParseDeviceInfo(b []byte) (*DevInfo, error) {
	r := NewReader(b)
	d := &DevInfo{}
	var err error

	if d.StandardVersion, err = r.U16(); err != nil {
		return nil, fmt.Errorf("ptp: device info: standard version: %w", err)
	}
	if d.VendorExtensionID, err = r.U32(); err != nil {
		return nil, fmt.Errorf("ptp: device info: vendor extension id: %w", err)
	}
	if d.VendorExtensionVer, err = r.U16(); err != nil {
		return nil, fmt.Errorf("ptp: device info: vendor extension version: %w", err)
	}
	if d.VendorExtensionsc, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: device info: vendor extension description: %w", err)
	}
	if d.FunctionalMode, err = r.U16(); err != nil {
		return nil, fmt.Errorf("ptp: device info: functional mode: %w", err)
	}

	ops, err := r.U16Array()
	if err != nil {
		return nil, fmt.Errorf("ptp: device info: operations: %w", err)
	}
	d.Operations = make([]OpCode, len(ops))
	for i, o := range ops {
		d.Operations[i] = OpCode(o)
	}
	if d.Events, err = r.U16Array(); err != nil {
		return nil, fmt.Errorf("ptp: device info: events: %w", err)
	}
	props, err := r.U16Array()
	if err != nil {
		return nil, fmt.Errorf("ptp: device info: device properties: %w", err)
	}
	d.DeviceProps = make([]Prop, len(props))
	for i, p := range props {
		d.DeviceProps[i] = Prop(p)
	}
	if d.CaptureFormats, err = r.U16Array(); err != nil {
		return nil, fmt.Errorf("ptp: device info: capture formats: %w", err)
	}
	if d.ImageFormats, err = r.U16Array(); err != nil {
		return nil, fmt.Errorf("ptp: device info: image formats: %w", err)
	}
	if d.Manufacturer, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: device info: manufacturer: %w", err)
	}
	if d.Model, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: device info: model: %w", err)
	}
	if d.DeviceVersion, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: device info: device version: %w", err)
	}
	if d.SerialNumber, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: device info: serial number: %w", err)
	}
	return d, nil
}

// u16Array reads a PTP array: uint32 count then that many uint16 values.
// U16Array reads a PTP array of uint16: a uint32 count, then that many values.
func (r *Reader) U16Array() ([]uint16, error) {
	n, err := r.U32()
	if err != nil {
		return nil, err
	}
	if int(n) > r.Remaining()/2 {
		return nil, fmt.Errorf("array of %d values: %w", n, ErrShortBlob)
	}
	out := make([]uint16, n)
	for i := range out {
		out[i] = binary.LittleEndian.Uint16(r.b[r.off:])
		r.off += 2
	}
	return out, nil
}

// GetDeviceInfo issues 0x1001 and parses the reply.
func (s *Session) GetDeviceInfo() (*DevInfo, error) {
	data, _, err := s.Do(OpGetDeviceInfo, nil, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	return ParseDeviceInfo(data)
}
