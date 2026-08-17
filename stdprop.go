package ptp

import "fmt"

// StdPropDesc is a standard PTP DevicePropDesc (reply to operation 0x1014).
//
//	+0  uint16  property code
//	+2  uint16  data type
//	+4  uint8   get/set
//	+5  default value, sized by data type
//	    current value, sized by data type
//	    uint8 form flag, then Range{min,max,step} or Enum{uint16 n, values}
//
// Note that Sony parsing is nonstandard because they insert isEnabled at +5 and
// starts values at +6.
type StdPropDesc struct {
	Code    Prop
	Type    DataType
	GetSet  uint8
	Default uint64
	Current uint64

	DefaultStr string
	CurrentStr string

	Form           FormFlag
	Min, Max, Step uint64

	// MinStr, MaxStr and StepStr hold the bounds of a STRING-typed range, set
	// instead of Min/Max/Step when Type is TypeString.
	MinStr, MaxStr, StepStr string
	Enum                    []uint64
	EnumStr                 []string // set instead of Enum for string-typed properties
}

// Writable reports whether the camera says the property can be set. In the
// standard encoding get/set is a plain 0/1, with no Sony control bit.
func (d *StdPropDesc) Writable() bool { return d.GetSet == 1 }

func (d *StdPropDesc) String() string {
	return fmt.Sprintf("%s(0x%04X) type=0x%04X getset=%d cur=%d",
		d.Code, uint16(d.Code), uint16(d.Type), d.GetSet, d.Current)
}

// ParseStdPropDesc decodes a standard PTP DevicePropDesc payload.
func ParseStdPropDesc(b []byte) (*StdPropDesc, error) {
	r := NewReader(b)
	d := &StdPropDesc{}

	code, err := r.U16()
	if err != nil {
		return nil, fmt.Errorf("ptp: property desc: code: %w", err)
	}
	typ, err := r.U16()
	if err != nil {
		return nil, fmt.Errorf("ptp: property desc: data type: %w", err)
	}
	d.Code, d.Type = Prop(code), DataType(typ)
	if d.GetSet, err = r.U8(); err != nil {
		return nil, fmt.Errorf("ptp: property desc: get/set: %w", err)
	}

	// A property with no datatype has no default or current value either: the
	// dataset ends after get/set, possibly with a form flag.
	if d.Type == TypeUndefined {
		if f, err := r.U8(); err == nil {
			d.Form = FormFlag(f)
		}
		return d, nil
	}

	readValue := func() (uint64, string, error) {
		if d.Type == TypeString {
			s, err := r.Str()
			return 0, s, err
		}
		v, err := r.Scalar(d.Type)
		return v, "", err
	}
	if d.Default, d.DefaultStr, err = readValue(); err != nil {
		return nil, fmt.Errorf("ptp: property desc: default value: %w", err)
	}
	if d.Current, d.CurrentStr, err = readValue(); err != nil {
		return nil, fmt.Errorf("ptp: property desc: current value: %w", err)
	}

	// The form is optional: a device may end the dataset after the values.
	form, err := r.U8()
	if err != nil {
		return d, nil
	}
	d.Form = FormFlag(form)
	switch d.Form {
	case FormRange:
		// A range can be string-typed. An X-T5 describes five properties that
		// way, with bounds like min "-6,-4,1" max "6,4,6" step "1,1,1" — a
		// structured triple carried as text rather than a number. Reading those
		// as scalars fails, and failing the whole descriptor loses the property
		// entirely.
		if d.Type == TypeString {
			if d.MinStr, err = r.Str(); err != nil {
				return nil, fmt.Errorf("ptp: property desc: range min string: %w", err)
			}
			if d.MaxStr, err = r.Str(); err != nil {
				return nil, fmt.Errorf("ptp: property desc: range max string: %w", err)
			}
			if d.StepStr, err = r.Str(); err != nil {
				return nil, fmt.Errorf("ptp: property desc: range step string: %w", err)
			}
			break
		}
		if d.Min, err = r.Scalar(d.Type); err != nil {
			return nil, fmt.Errorf("ptp: property desc: range min: %w", err)
		}
		if d.Max, err = r.Scalar(d.Type); err != nil {
			return nil, fmt.Errorf("ptp: property desc: range max: %w", err)
		}
		if d.Step, err = r.Scalar(d.Type); err != nil {
			return nil, fmt.Errorf("ptp: property desc: range step: %w", err)
		}
	case FormEnum:
		n, err := r.U16()
		if err != nil {
			return nil, fmt.Errorf("ptp: property desc: enum count: %w", err)
		}
		// A string-typed property has a string enum: one body reports ImageSize
		// as a list of "4096x3072"-style strings, not numbers. Sizing the
		// allocation from a scalar width would reject it outright.
		if d.Type == TypeString {
			if int(n) > r.Remaining() {
				return nil, fmt.Errorf("ptp: property desc: string enum of %d values: %w", n, ErrShortBlob)
			}
			d.EnumStr = make([]string, n)
			for i := range d.EnumStr {
				if d.EnumStr[i], err = r.Str(); err != nil {
					return nil, fmt.Errorf("ptp: property desc: string enum value %d: %w", i, err)
				}
			}
			break
		}
		sz, ok := TypeSize(d.Type)
		if !ok {
			return nil, fmt.Errorf("ptp: property desc: unsupported enum type 0x%04X", uint16(d.Type))
		}
		if int(n) > r.Remaining()/sz {
			return nil, fmt.Errorf("ptp: property desc: enum of %d values: %w", n, ErrShortBlob)
		}
		d.Enum = make([]uint64, n)
		for i := range d.Enum {
			if d.Enum[i], err = r.Scalar(d.Type); err != nil {
				return nil, fmt.Errorf("ptp: property desc: enum value %d: %w", i, err)
			}
		}
	}
	return d, nil
}

// stdPropNames covers the standard PTP and MTP device properties. These are not
// in Sony's table — that only describes the vendor remote-control surface — but
// they are what a body exposes when it is acting purely as a file-transfer
// device, so naming them keeps that path readable.
var stdPropNames = map[Prop]string{
	0x5001: "BatteryLevel",
	0x5002: "FunctionalMode",
	0x5003: "ImageSize",
	0x5004: "CompressionSetting",
	0x5005: "WhiteBalance",
	0x5006: "RGBGain",
	0x5007: "FNumber",
	0x5008: "FocalLength",
	0x5009: "FocusDistance",
	0x500A: "FocusMode",
	0x500B: "ExposureMeteringMode",
	0x500C: "FlashMode",
	0x500D: "ExposureTime",
	0x500E: "ExposureProgramMode",
	0x500F: "ExposureIndex",
	0x5010: "ExposureBiasCompensation",
	0x5011: "DateTime",
	0x5012: "CaptureDelay",
	0x5013: "StillCaptureMode",
	0x5014: "Contrast",
	0x5015: "Sharpness",
	0x5016: "DigitalZoom",
	0x5017: "EffectMode",
	0x5018: "BurstNumber",
	0x5019: "BurstInterval",
	0x501A: "TimelapseNumber",
	0x501B: "TimelapseInterval",
	0x501C: "FocusMeteringMode",
	0x501D: "UploadURL",
	0x501E: "Artist",
	0x501F: "CopyrightInfo",
	// MTP device properties.
	0xD401: "SynchronizationPartner",
	0xD402: "DeviceFriendlyName",
	0xD403: "Volume",
	0xD405: "DeviceIcon",
	0xD406: "SessionInitiatorInfo",
	0xD407: "PerceivedDeviceType",
}

// GetPropDesc issues 0x1014 GetDevicePropDesc for one property. This is the
// standard PTP path, which works on bodies that expose no Sony vendor
// operations at all.
func (s *Session) GetPropDesc(p Prop) (*StdPropDesc, error) {
	data, _, err := s.Do(OpGetDevicePropDesc, []uint32{uint32(p)}, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	return ParseStdPropDesc(data)
}
