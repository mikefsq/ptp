package sony

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Stored settings beyond the exposure triangle: the clock, and the image
// format surface.
//
// All of it goes through the vendor write path (0x9205), so all of it is
// UNVERIFIED — see the package doc. The wire codes come from Sony's own SDK
// table and the types from its headers.

// SetPropertyString writes a string-typed setting via 0x9205.
//
// SetProperty encodes a scalar of a fixed width; a string is a uint8 character
// count followed by UTF-16, so it needs its own path.
func (c *Camera) SetPropertyString(p Prop, v string) error {
	_, _, err := c.Do(OpSetControlDeviceA, []uint32{uint32(p)}, ptp.EncodeString(v), ptp.DefaultTimeout)
	return err
}

// PropString reads a string-typed setting out of a property snapshot.
func PropString(props []DeviceProperty, p Prop) (string, bool) {
	d := FindProp(props, p)
	if d == nil || d.Type != ptp.TypeString {
		return "", false
	}
	return d.CurrentStr, true
}

// DateTime reads the camera's clock.
//
// The property is a string in PTP's date-time form. Sony's own struct for it,
// CrDateTimeSetting, carries char dateTime[18] alongside a "time zone exists"
// flag, so a body may report a zone this deliberately ignores — the camera
// keeps local time and the host has no business reinterpreting it.
func (c *Camera) DateTime() (time.Time, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return time.Time{}, err
	}
	s, ok := PropString(props, PropDateTimeSettings)
	if !ok {
		return time.Time{}, fmt.Errorf("sony: this body does not report %s as a string",
			PropName(PropDateTimeSettings))
	}
	return ptp.ParseDateTime(s)
}

// SetDateTime sets the camera's clock.
func (c *Camera) SetDateTime(t time.Time) error {
	return c.SetPropertyString(PropDateTimeSettings, ptp.FormatDateTime(t))
}

// SyncClock sets the camera's clock to the host's current time.
//
// Whether the body truncates or rounds the seconds is NOT known — a Fujifilm
// X-T5 truncates, leaving it up to 59 seconds slow, and no Sony has been
// measured. Use SyncClockAtMinute where the timestamps matter, which is correct
// either way.
func (c *Camera) SyncClock() error { return c.SetDateTime(time.Now()) }

// SyncClockAtMinute waits for the next minute boundary and sets the clock then.
//
// This is the safe form: if the body truncates seconds, writing at any other
// moment bakes in the remainder as error, and writing on the boundary leaves
// only the round trip. It blocks for up to a minute, which is why it is a
// separate call.
func (c *Camera) SyncClockAtMinute() error {
	now := time.Now()
	next := now.Truncate(time.Minute).Add(time.Minute)
	// Aim marginally early so the write lands as the minute turns rather than
	// just after it, which would truncate to the minute just gone.
	if d := time.Until(next) - 50*time.Millisecond; d > 0 {
		time.Sleep(d)
	}
	return c.SetDateTime(time.Now().Add(60 * time.Millisecond))
}

// Still image file formats (CrFileType, CrInt16u).
//
// The order is Sony's and it is NOT the one you would guess: JPEG comes before
// RAW. Values are taken verbatim from CrDeviceProperty.h rather than assumed —
// getting these two the wrong way round means shooting JPEG for a whole
// sequence while believing it is RAW.
//
// RawJpeg and RawHeif write two files per frame, which doubles the buffer
// drain — see burst.go.
const (
	FileTypeNone    uint64 = 0x0000
	FileTypeJPEG    uint64 = 0x0001
	FileTypeRAW     uint64 = 0x0002
	FileTypeRAWJPEG uint64 = 0x0003
	FileTypeRAWHEIF uint64 = 0x0004
	FileTypeHEIF    uint64 = 0x0005
)

// SetFileType selects RAW, JPEG, or both (0xD253).
func (c *Camera) SetFileType(t uint64) error {
	return c.SetProperty(PropFileType, ptp.TypeUint16, t)
}

// Still image quality (CrImageQuality, CrInt16u) — the JPEG/HEIF compression
// level, independent of FileType.
//
// Four values, and that is all of them. The SDK also carries a deprecated
// CrJpegQuality with the same numbering, marked "Do not use".
const (
	QualityUnknown   uint64 = 0x0000
	QualityLight     uint64 = 0x0001
	QualityStandard  uint64 = 0x0002
	QualityFine      uint64 = 0x0003
	QualityExtraFine uint64 = 0x0004
)

// SetStillImageQuality sets the JPEG/HEIF quality (0xD252).
func (c *Camera) SetStillImageQuality(q uint64) error {
	return c.SetProperty(PropStillImageQuality, ptp.TypeUint16, q)
}

// Image sizes (CrImageSize).
const (
	ImageSizeLarge  uint64 = 0x0001
	ImageSizeMedium uint64 = 0x0002
	ImageSizeSmall  uint64 = 0x0003
	ImageSizeVGA    uint64 = 0x0004
)

// SetImageSize sets the recorded pixel size (0xD203).
//
// On a 61 MP body this is the difference between a 60 MB and a 15 MB RAW, so it
// is one of the few settings that changes sustained burst depth as much as the
// compression type does.
func (c *Camera) SetImageSize(s uint64) error {
	return c.SetProperty(PropImageSize, ptp.TypeUint16, s)
}

// Aspect ratios (CrAspectRatioIndex, CrInt16u). Four values; the SDK defines no
// others for stills.
const (
	Aspect3to2  uint64 = 0x0001
	Aspect16to9 uint64 = 0x0002
	Aspect4to3  uint64 = 0x0003
	Aspect1to1  uint64 = 0x0004
)

// SetAspectRatio sets the recorded aspect ratio (0xD211).
func (c *Camera) SetAspectRatio(a uint64) error {
	return c.SetProperty(PropAspectRatio, ptp.TypeUint16, a)
}

// Settings is a decoded snapshot of the image-format surface, with a "was it
// reported" flag per field for the same reason Exposure has them: a body that
// does not expose a setting is not the same as one reporting zero.
type Settings struct {
	FileType    uint64
	HasFileType bool

	Quality    uint64
	HasQuality bool

	ImageSize    uint64
	HasImageSize bool

	AspectRatio    uint64
	HasAspectRatio bool

	RAWCompression    uint64
	HasRAWCompression bool

	StoreDestination    uint64
	HasStoreDestination bool

	DriveMode    uint64
	HasDriveMode bool
}

// ReadSettings pulls the image-format surface out of a property snapshot, so a
// caller checking several does not pay a round trip each.
func ReadSettings(props []DeviceProperty) Settings {
	var s Settings
	get := func(code Prop, dst *uint64, has *bool) {
		if p := FindProp(props, code); p != nil {
			*dst, *has = p.Current, true
		}
	}
	get(PropFileType, &s.FileType, &s.HasFileType)
	get(PropStillImageQuality, &s.Quality, &s.HasQuality)
	get(PropImageSize, &s.ImageSize, &s.HasImageSize)
	get(PropAspectRatio, &s.AspectRatio, &s.HasAspectRatio)
	get(PropRAWFileCompressionType, &s.RAWCompression, &s.HasRAWCompression)
	get(PropStillImageStoreDestination, &s.StoreDestination, &s.HasStoreDestination)
	get(PropDriveMode, &s.DriveMode, &s.HasDriveMode)
	return s
}

// Settings fetches a snapshot and decodes the image-format surface.
func (c *Camera) Settings() (Settings, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return Settings{}, err
	}
	return ReadSettings(props), nil
}

// Settable reports whether the host can currently write a property.
//
// A Sony hands control over through its exposure program and its physical
// controls much as a Fujifilm does through its dials, and the symptom is the
// same: the write is ACCEPTED and then ignored. The camera's own enable byte is
// the only reliable answer, so ask it rather than inferring from a write's
// response.
func (c *Camera) Settable(p Prop) (bool, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return false, err
	}
	d := FindProp(props, p)
	if d == nil {
		return false, fmt.Errorf("sony: this body does not report %s", PropName(p))
	}
	return d.Writable(), nil
}
