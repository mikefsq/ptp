package fuji

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Typed access to the settings a tethered session is most likely to want, over
// and above capture and exposure.
//
// Each is a thin wrapper over GetProp/SetProp, which validate against a live
// descriptor. The value is in naming the property and its values, and in the
// string ones being reachable at all — see GetPropString.

// Property codes not already declared from gphoto2's table. Recovered from the
// X-T5's own plugin binary; where both sources have a code they agree, which is
// worth knowing given gphoto2 gets some NAMES wrong for this body.
const (
	PropRAWOutputDepth ptp.Prop = 0xD188
	PropAutoPowerOff   ptp.Prop = 0xD381
	PropArtist         ptp.Prop = 0xD237
	PropFilenamePrefix ptp.Prop = 0xD366
	PropCropZoom       ptp.Prop = 0xD1B3
	PropImageFormat    ptp.Prop = 0xD1B2
	PropApertureUnit   ptp.Prop = 0xD239
)

// GetPropString reads a string-typed property.
//
// This exists because GetProp returns the descriptor's numeric field, which is
// zero for a string property — so DateTime read back as 0 and looked broken.
// PTP carries strings in their own slot, and nothing about the property code
// says which kind it is; only the descriptor does.
func (c *Camera) GetPropString(p ptp.Prop) (string, error) {
	d, err := c.GetPropDesc(p)
	if err != nil {
		return "", fmt.Errorf("fuji: reading %s: %w", PropName(p), err)
	}
	if d.Type != ptp.TypeString {
		return "", fmt.Errorf("fuji: %s is not a string property (type 0x%04X)",
			PropName(p), uint16(d.Type))
	}
	return d.CurrentStr, nil
}

// SetPropStringChecked writes a string-typed property, refusing one the camera
// will not accept rather than writing something it silently ignores.
func (c *Camera) SetPropStringChecked(p ptp.Prop, v string) error {
	d, err := c.GetPropDesc(p)
	if err != nil {
		return fmt.Errorf("fuji: reading %s before writing it: %w", PropName(p), err)
	}
	if d.Type != ptp.TypeString {
		return fmt.Errorf("fuji: %s is not a string property (type 0x%04X)",
			PropName(p), uint16(d.Type))
	}
	if !d.Writable() {
		return fmt.Errorf("fuji: %s is read-only", PropName(p))
	}
	if d.CurrentStr == v {
		return nil // a redundant write draws vendor error 0xA002
	}
	if err := c.SetPropString(p, v); err != nil {
		return fmt.Errorf("fuji: setting %s to %q: %w", PropName(p), v, err)
	}
	return nil
}

// Clock.

// DateTime reads the camera's clock.
//
// Worth setting for any sequence whose frames will be correlated with anything
// else afterwards: the timestamp in each file comes from here, and a body that
// has been in a drawer for a year is not close.
func (c *Camera) DateTime() (time.Time, error) {
	s, err := c.GetPropString(ptp.PropDateTime)
	if err != nil {
		return time.Time{}, err
	}
	return ptp.ParseDateTime(s)
}

// SetDateTime sets the camera's clock.
func (c *Camera) SetDateTime(t time.Time) error {
	return c.SetPropStringChecked(ptp.PropDateTime, ptp.FormatDateTime(t))
}

// SyncClock sets the camera's clock to the host's current time.
//
// The X-T5 stores whole minutes and TRUNCATES the seconds rather than rounding:
// a sync at 14:41:56 leaves the camera reading 14:41:00. So this can be up to 59
// seconds slow. Use SyncClockAtMinute where the timestamps matter.
func (c *Camera) SyncClock() error { return c.SetDateTime(time.Now()) }

// SyncClockAtMinute waits for the next minute boundary and sets the clock then,
// which is the only way to get the camera within a second of the host.
//
// Because the body truncates seconds, writing at any other moment bakes in the
// remainder as error. Writing exactly on the boundary leaves only the round
// trip, a few milliseconds. It blocks for up to a minute, which is why it is a
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

// Metadata written into every frame.

func (c *Camera) Artist() (string, error)     { return c.GetPropString(PropArtist) }
func (c *Camera) SetArtist(v string) error    { return c.SetPropStringChecked(PropArtist, v) }
func (c *Camera) Comment() (string, error)    { return c.GetPropString(PropComment) }
func (c *Camera) SetComment(v string) error   { return c.SetPropStringChecked(PropComment, v) }
func (c *Camera) Copyright() (string, error)  { return c.GetPropString(PropCopyright) }
func (c *Camera) SetCopyright(v string) error { return c.SetPropStringChecked(PropCopyright, v) }

// Copyright2 is the body's second copyright line — the X-T5 exposes two.
func (c *Camera) Copyright2() (string, error)  { return c.GetPropString(PropCopyright2) }
func (c *Camera) SetCopyright2(v string) error { return c.SetPropStringChecked(PropCopyright2, v) }

// FilenamePrefix is the four characters that begin each filename, "DSCF" by
// default. Worth setting when two bodies shoot the same sequence, so their
// frames do not collide by name.
func (c *Camera) FilenamePrefix() (string, error) { return c.GetPropString(PropFilenamePrefix) }

// SetFilenamePrefix sets that prefix.
func (c *Camera) SetFilenamePrefix(v string) error {
	return c.SetPropStringChecked(PropFilenamePrefix, v)
}

// ImageSize reports the frame dimensions, as the camera words them ("7728x5152").
func (c *Camera) ImageSize() (string, error) { return c.GetPropString(ptp.PropImageSize) }

// Numeric settings.

// RAWOutputDepth is the RAW bit depth, 14 or 16 bit on an X-T5.
func (c *Camera) RAWOutputDepth() (uint64, error) { return c.GetProp(PropRAWOutputDepth) }

// SetRAWOutputDepth selects the RAW bit depth.
func (c *Camera) SetRAWOutputDepth(v uint64) error { return c.SetProp(PropRAWOutputDepth, v) }

// LiveViewImageQuality and LiveViewImageSize govern the preview stream rather
// than captured frames — smaller and coarser previews cost less to transfer,
// which matters when focusing against a live image.
func (c *Camera) LiveViewImageQuality() (uint64, error) { return c.GetProp(PropLiveViewImageQuality) }

// SetLiveViewImageQuality sets the preview quality.
func (c *Camera) SetLiveViewImageQuality(v uint64) error {
	return c.SetProp(PropLiveViewImageQuality, v)
}

// LiveViewImageSize reports the preview frame size.
func (c *Camera) LiveViewImageSize() (uint64, error) { return c.GetProp(PropLiveViewImageSize) }

// SetLiveViewImageSize sets the preview frame size.
func (c *Camera) SetLiveViewImageSize(v uint64) error { return c.SetProp(PropLiveViewImageSize, v) }

// AutoPowerOff is the body's idle timeout.
//
// Relevant beyond battery life: a camera that sleeps mid-sequence drops the USB
// link, and recovering from that costs far more than the power saved.
func (c *Camera) AutoPowerOff() (uint64, error) { return c.GetProp(PropAutoPowerOff) }

// SetAutoPowerOff sets the idle timeout.
func (c *Camera) SetAutoPowerOff(v uint64) error { return c.SetProp(PropAutoPowerOff, v) }

// CropMode and CropZoom control in-camera cropping.
func (c *Camera) CropMode() (uint64, error)     { return c.GetProp(PropCropMode) }
func (c *Camera) SetCropMode(v uint64) error    { return c.SetProp(PropCropMode, v) }
func (c *Camera) CropZoom() (uint64, error)     { return c.GetProp(PropCropZoom) }
func (c *Camera) SetCropZoom(v uint64) error    { return c.SetProp(PropCropZoom, v) }
func (c *Camera) ImageFormat() (uint64, error)  { return c.GetProp(PropImageFormat) }
func (c *Camera) SetImageFormat(v uint64) error { return c.SetProp(PropImageFormat, v) }

// ApertureUnit selects how the camera expresses aperture — F-stops or T-stops
// on a cine lens.
func (c *Camera) ApertureUnit() (uint64, error)  { return c.GetProp(PropApertureUnit) }
func (c *Camera) SetApertureUnit(v uint64) error { return c.SetProp(PropApertureUnit, v) }

// ExposurePreview is the body's WYSIWYG preview. Turning it off is usual for a
// dark subject, where a faithful preview shows almost nothing.
func (c *Camera) ExposurePreview() (uint64, error)  { return c.GetProp(PropExposurePreview) }
func (c *Camera) SetExposurePreview(v uint64) error { return c.SetProp(PropExposurePreview, v) }

// PropHistogramData carries the live histogram. Its descriptor declares type
// "undefined", so the value is NOT in the descriptor — a plain property read
// returns 0, and the data only comes back from GetDevicePropValue as a blob.
const PropHistogramData ptp.Prop = 0xD22F

// Histogram is the camera's own histogram of the current scene: four channels
// of 256 bins.
//
// The channel order is not documented and has not been confirmed — on a neutral
// subject all four are near-identical, which is no help. Three are presumably
// red, green and blue, and the fourth luminance. Bins are counts, and all four
// channels sum to the same total, being four views of one image.
type Histogram struct {
	Channel [4][256]uint32
}

// Total returns the number of samples in a channel, which is the same for all
// four and is a useful sanity check on a decode.
func (h *Histogram) Total(ch int) uint64 {
	var n uint64
	for _, v := range h.Channel[ch] {
		n += uint64(v)
	}
	return n
}

// Clipped reports the fraction of samples in the topmost bin — the highlights
// that have hit the ceiling. For a filtered sun that is the number worth
// watching, since it says the exposure has run out of headroom without needing
// the frame downloaded.
func (h *Histogram) Clipped(ch int) float64 {
	total := h.Total(ch)
	if total == 0 {
		return 0
	}
	return float64(h.Channel[ch][255]) / float64(total)
}

// Histogram fetches the camera's histogram of what it is seeing now.
//
// It works with live view running or stopped, needs no frame downloaded, and
// costs a few milliseconds — so exposure can be watched continuously, where a
// captured RAW is 35 MB and most of a second.
//
// UNRESOLVED: whether this describes the CAPTURE exposure or the LIVE VIEW
// image. It does respond to shutter changes — 60ms produced a far brighter
// distribution than 244us, and clipping appeared — but the response was not
// monotonic across three settings, so either it lags a change or the camera
// gains its preview independently. Before trusting it to verify a capture
// exposure, check it against a frame actually taken.
func (c *Camera) Histogram() (*Histogram, error) {
	data, _, err := c.Do(ptp.OpGetDevicePropValue, []uint32{uint32(PropHistogramData)},
		nil, ptp.DefaultTimeout)
	if err != nil {
		return nil, fmt.Errorf("fuji: reading the histogram: %w", err)
	}
	const want = 4 * 256 * 4
	if len(data) != want {
		return nil, fmt.Errorf("fuji: histogram is %d bytes, want %d", len(data), want)
	}
	h := &Histogram{}
	for ch := 0; ch < 4; ch++ {
		for b := 0; b < 256; b++ {
			off := (ch*256 + b) * 4
			h.Channel[ch][b] = uint32(data[off]) | uint32(data[off+1])<<8 |
				uint32(data[off+2])<<16 | uint32(data[off+3])<<24
		}
	}
	return h, nil
}

// ExposureBias is exposure compensation in EV x 1000: -3 EV reads -3000.
//
// Confirmed on hardware by moving the body's EV dial to -3 and watching the
// property. The encoding matches Sony's, which is a rare case of two vendors
// agreeing on something.
//
// It is signed, so the parser sign-extends and a raw read of a negative value
// looks enormous; these helpers return EV directly.
func (c *Camera) ExposureBias() (float64, error) {
	d, err := c.GetPropDesc(ptp.PropExposureBias)
	if err != nil {
		return 0, fmt.Errorf("fuji: reading exposure compensation: %w", err)
	}
	return float64(int64(d.Current)) / 1000, nil
}

// SetExposureBias sets exposure compensation in EV.
func (c *Camera) SetExposureBias(ev float64) error {
	return c.SetProp(ptp.PropExposureBias, uint64(int64(ev*1000)))
}

// EVDialStatus and ISODialStatus report whether the body's exposure
// compensation and ISO dials are on C — that is, whether the host may set
// those values or a physical position has claimed them.
//
// Both were identified by moving the dial and diffing every property; see
// names_observed.go. They answer the ownership question directly, where
// Settable has to infer it from the property advertising a single value.
func (c *Camera) EVDialStatus() (uint64, error)  { return c.GetProp(0xD034) }
func (c *Camera) ISODialStatus() (uint64, error) { return c.GetProp(0xD035) }

// HostOwnsExposureBias reports whether the EV dial is on C.
func (c *Camera) HostOwnsExposureBias() (bool, error) {
	v, err := c.EVDialStatus()
	return v == DialOnCommand, err
}

// HostOwnsISO reports whether the ISO dial is on C.
func (c *Camera) HostOwnsISO() (bool, error) {
	v, err := c.ISODialStatus()
	return v == DialOnCommand, err
}

// Exposure preview modes for PropExposurePreview — the body's PREVIEW EXP./WB
// IN MANUAL MODE setting.
//
// ExposurePreviewExposureAndWB stops the lens down to show the exposure as it
// will be recorded. That is what you want when judging a normal scene, and
// exactly what you do not want at a small aperture on a dark subject, where it
// leaves the finder black and looks like the camera is stuck.
const (
	ExposurePreviewExposureAndWB uint64 = 1
	ExposurePreviewWBOnly        uint64 = 2 // inferred from the menu, untested
	ExposurePreviewOff           uint64 = 3
)
