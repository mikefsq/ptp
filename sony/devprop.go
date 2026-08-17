package sony

import (
	"fmt"

	"github.com/mikefsq/ptp"
)

// DeviceProperty is one entry of the 0x9209 GetAllDevicePropData blob.
//
//	+0  uint16  property code
//	+2  uint16  data type (standard PTP codes 1..8, 0xFFFF for string)
//	+4  uint8   get/set
//	+5  uint8   isEnabled
//	+6  default value, sized by data type
//	    current value, sized by data type
//	    uint8 form flag, then Range{min,max,step} or Enum{uint16 n, values}
//
// This is NOT ptp.StdPropDesc. The standard layout has no isEnabled byte and
// starts its values at +5, so parsing one with the other's offsets yields
// plausible-looking rubbish rather than an error. That is why the vendor
// descriptor is parsed here and the standard one in the parent package.
type DeviceProperty struct {
	Code    Prop
	Type    ptp.DataType
	GetSet  uint8
	Enabled uint8

	// Default and Current hold integer values zero-extended to 64 bits. For
	// signed types the value is sign-extended first, so a negative exposure
	// compensation reads back negative when cast to int64.
	Default uint64
	Current uint64

	// DefaultStr and CurrentStr hold string-typed values (Type == ptp.TypeString).
	DefaultStr string
	CurrentStr string

	// DefaultArr and CurrentArr hold array-typed values (Type.IsArray()).
	DefaultArr []uint64
	CurrentArr []uint64

	Form ptp.FormFlag

	Min, Max, Step uint64 // valid when Form == ptp.FormRange
	Enum           []uint64
}

// Writable reports whether the camera says this property can be set.
//
// The get/set byte carries two different meanings. Values with the 0x80 bit set
// mark a control rather than a stored setting (0x81 button, 0x82 notch,
// 0x83 lock, 0x84 variable); those are always actionable. Otherwise the
// isEnabled byte governs: 1 means enabled, 0 means greyed out, 2 means display
// only.
func (d *DeviceProperty) Writable() bool {
	if d.GetSet&0x80 != 0 {
		return true
	}
	return d.Enabled == 1
}

// IsControl reports whether this entry is a control (driven with SendControl)
// rather than a stored setting (driven with SetProperty).
func (d *DeviceProperty) IsControl() bool { return d.GetSet&0x80 != 0 }

func (d *DeviceProperty) String() string {
	return fmt.Sprintf("%s(0x%04X) type=0x%04X getset=0x%02X en=%d cur=%d",
		PropName(d.Code), uint16(d.Code), uint16(d.Type), d.GetSet, d.Enabled, d.Current)
}

// ParseAllDevicePropData parses the payload of a 0x9209 GetAllDevicePropData
// transfer.
//
// The blob opens with a uint32 entry count and a uint32 zero, then that many
// property entries.
func ParseAllDevicePropData(b []byte) ([]DeviceProperty, error) {
	r := ptp.NewReader(b)
	count, err := r.U32()
	if err != nil {
		return nil, fmt.Errorf("sony: reading property count: %w", err)
	}
	if _, err = r.U32(); err != nil { // reserved, always zero
		return nil, fmt.Errorf("sony: reading blob header: %w", err)
	}

	out := make([]DeviceProperty, 0, count)
	for i := range count {
		d, err := parseOne(r)
		if err != nil {
			return out, fmt.Errorf("sony: property %d of %d: %w", i, count, err)
		}
		out = append(out, d)
	}
	return out, nil
}

func parseOne(r *ptp.Reader) (DeviceProperty, error) {
	var d DeviceProperty
	code, err := r.U16()
	if err != nil {
		return d, err
	}
	typ, err := r.U16()
	if err != nil {
		return d, err
	}
	d.Code, d.Type = Prop(code), ptp.DataType(typ)
	if d.GetSet, err = r.U8(); err != nil {
		return d, err
	}
	if d.Enabled, err = r.U8(); err != nil {
		return d, err
	}

	readValue := func() (uint64, string, []uint64, error) {
		v, err := r.Value(d.Type)
		return v.Num, v.Str, v.Arr, err
	}

	if d.Default, d.DefaultStr, d.DefaultArr, err = readValue(); err != nil {
		return d, fmt.Errorf("default value: %w", err)
	}
	// A string property with nothing after the default has no current value or
	// form; gphoto2 relies on the same short-circuit.
	if d.Type == ptp.TypeString && r.Remaining() == 0 {
		return d, nil
	}
	if d.Current, d.CurrentStr, d.CurrentArr, err = readValue(); err != nil {
		return d, fmt.Errorf("current value: %w", err)
	}

	form, err := r.U8()
	if err != nil {
		return d, fmt.Errorf("form flag: %w", err)
	}
	d.Form = ptp.FormFlag(form)

	switch d.Form {
	case ptp.FormRange:
		if d.Min, err = r.Scalar(d.Type); err != nil {
			return d, fmt.Errorf("range min: %w", err)
		}
		if d.Max, err = r.Scalar(d.Type); err != nil {
			return d, fmt.Errorf("range max: %w", err)
		}
		if d.Step, err = r.Scalar(d.Type); err != nil {
			return d, fmt.Errorf("range step: %w", err)
		}
	case ptp.FormEnum:
		if d.Enum, err = readEnum(r, d.Type); err != nil {
			return d, err
		}
		// Bodies from 2024 on (which includes the A7R V and A7R VI) may follow
		// the enum with a second list that supersedes it. It is distinguished
		// by its count being below 0x200 — a real next property code is always
		// 0x5xxx or 0xDxxx and so cannot be confused with one.
		if next, ok := r.Peek16(); ok && next < 0x200 {
			r.U16() // consume the count Peek16 only looked at
			second, err := readEnumOf(r, d.Type, int(next))
			if err != nil {
				return d, fmt.Errorf("secondary %w", err)
			}
			d.Enum = second
		}
	}
	return d, nil
}

// readEnum reads a uint16 count and that many values of type t.
func readEnum(r *ptp.Reader, t ptp.DataType) ([]uint64, error) {
	n, err := r.U16()
	if err != nil {
		return nil, fmt.Errorf("enum count: %w", err)
	}
	return readEnumOf(r, t, int(n))
}

func readEnumOf(r *ptp.Reader, t ptp.DataType, n int) ([]uint64, error) {
	sz, ok := ptp.TypeSize(t)
	if !ok {
		return nil, fmt.Errorf("enum of unsupported type 0x%04X", uint16(t))
	}
	// Reject a count the buffer cannot hold before allocating, so a corrupt
	// length cannot drive a huge allocation.
	if n > r.Remaining()/sz {
		return nil, fmt.Errorf("enum of %d values: %w", n, ptp.ErrShortBlob)
	}
	out := make([]uint64, n)
	for i := range out {
		v, err := r.Scalar(t)
		if err != nil {
			return nil, fmt.Errorf("enum value %d: %w", i, err)
		}
		out[i] = v
	}
	return out, nil
}

// GetAllDevicePropData issues 0x9209 and parses the result. This is the primary
// read path: the SDK polls it wholesale rather than reading properties one at a
// time.
func (c *Camera) GetAllDevicePropData() ([]DeviceProperty, error) {
	data, _, err := c.Do(OpGetAllDevicePropData, nil, nil, ptp.DefaultTimeout)
	if err != nil {
		return nil, err
	}
	return ParseAllDevicePropData(data)
}

// FindProp returns the property with the given code from a snapshot, or nil.
//
// GetAllDevicePropData returns a slice because that is what the camera sends;
// callers almost always want to look one up by code.
func FindProp(props []DeviceProperty, code Prop) *DeviceProperty {
	for i := range props {
		if props[i].Code == code {
			return &props[i]
		}
	}
	return nil
}

// PropMap indexes a snapshot by code, for callers touching several properties.
func PropMap(props []DeviceProperty) map[Prop]*DeviceProperty {
	m := make(map[Prop]*DeviceProperty, len(props))
	for i := range props {
		m[props[i].Code] = &props[i]
	}
	return m
}
