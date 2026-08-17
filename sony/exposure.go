package sony

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/mikefsq/ptp"
)

// Typed access to the exposure triangle and focus.
//
// Every encoding here comes from Sony's own API reference, not from guesswork.
// The wire values are not what you would expect: aperture is F×100, shutter
// speed is a packed numerator/denominator pair, and ISO packs a mode into the
// high bits.
//
// A camera only accepts values from the set it currently advertises. Aperture,
// shutter and ISO are all array-typed properties whose valid values arrive in
// the 0x9209 blob, and that set changes with lens, exposure mode and drive
// mode. Encode* produces the wire value for a given real-world setting; the
// Nearest* helpers pick the closest value the camera will actually accept.

// Aperture special values (CrFnumberSet).
const (
	FNumberIrisClose uint64 = 0xFFFD
	FNumberUnknown   uint64 = 0xFFFE
	FNumberNothing   uint64 = 0xFFFF
)

// Shutter speed special values (CrShutterSpeedSet).
const (
	ShutterBulb    uint64 = 0x00000000
	ShutterNothing uint64 = 0xFFFFFFFF
)

// ISO mode bits (CrISOMode), held in bits 24-27, with bits 0-23 the value.
const (
	ISOModeNormal                = 0x00
	ISOModeMultiFrameNR          = 0x01
	ISOModeMultiFrameNRHi        = 0x02
	ISOModeExtended              = 0x10
	ISOValueAuto          uint32 = 0xFFFFFF
)

// EncodeAperture converts an f-number to its wire value: F×100, so f/4 is 400
// and f/9.5 is 950.
func EncodeAperture(f float64) uint64 { return uint64(math.Round(f * 100)) }

// DecodeAperture converts a wire value back to an f-number. ok is false for the
// special values (iris closed, unknown, nothing to display).
func DecodeAperture(v uint64) (f float64, ok bool) {
	switch v {
	case FNumberIrisClose, FNumberUnknown, FNumberNothing:
		return 0, false
	}
	return float64(v) / 100, true
}

// Shutter speed limits imposed by the wire format. Both halves of the packed
// value are 16 bits, so the fastest representable time is 1/65535 and the
// slowest is 65535/10.
//
// The fast limit is well beyond any camera: the quickest Sony electronic
// shutter is 1/32000 (31.25us) and the quickest mechanical is 1/8000 (125us).
const (
	MinShutterSpeed = time.Second / 65535        // ~15.26us
	MaxShutterSpeed = 6553500 * time.Millisecond // 6553.5s, about 109 minutes
)

// ValidateShutterSpeed reports whether d can be represented on the wire.
// A zero or negative duration is valid and means Bulb.
func ValidateShutterSpeed(d time.Duration) error {
	if d <= 0 {
		return nil // Bulb
	}
	if d < MinShutterSpeed {
		return fmt.Errorf("sony: shutter speed %v is faster than the wire format allows "+
			"(minimum %v, 1/65535); no camera is this fast either", d, MinShutterSpeed)
	}
	if d > MaxShutterSpeed {
		return fmt.Errorf("sony: shutter speed %v is longer than the wire format allows "+
			"(maximum %v); use Bulb for longer exposures", d, MaxShutterSpeed)
	}
	return nil
}

// EncodeShutterSpeed converts an exposure time to its wire value.
//
// Times outside MinShutterSpeed..MaxShutterSpeed are CLAMPED, not rejected, so
// call ValidateShutterSpeed first when the duration comes from outside the
// program. SetShutterSpeed does that already.
//
// The value packs a numerator in the upper 16 bits and a denominator in the
// lower 16. Sony uses two conventions, matching what the camera displays:
// times of a second or more are written in tenths (denominator fixed at 10), so
// 1.5s is 0x000F000A; shorter times are written as a unit fraction (numerator
// fixed at 1), so 1/1000 is 0x000103E8.
//
// A zero or negative duration encodes as Bulb.
func EncodeShutterSpeed(d time.Duration) uint64 {
	if d <= 0 {
		return ShutterBulb
	}
	secs := d.Seconds()
	if secs < 1 {
		// Prefer the unit fraction, but only when it actually represents the
		// time: 1/round(1/0.7) is 1/1, which is a whole second, not 0.7. Fall
		// through to tenths when the fraction is a poor fit.
		den := math.Round(1 / secs)
		if den >= 1 && math.Abs(1/den-secs) <= secs*0.02 {
			if den > 0xFFFF {
				den = 0xFFFF
			}
			return 1<<16 | uint64(den)
		}
	}
	tenths := uint64(math.Round(secs * 10))
	if tenths > 0xFFFF {
		tenths = 0xFFFF
	}
	if tenths == 0 {
		tenths = 1
	}
	return tenths<<16 | 10
}

// DecodeShutterSpeed converts a wire value to a duration. ok is false for Bulb
// and for "nothing to display", neither of which is a time.
func DecodeShutterSpeed(v uint64) (d time.Duration, ok bool) {
	switch uint32(v) {
	case uint32(ShutterBulb), uint32(ShutterNothing):
		return 0, false
	}
	num := (v >> 16) & 0xFFFF
	den := v & 0xFFFF
	if den == 0 {
		return 0, false
	}
	return time.Duration(float64(num) / float64(den) * float64(time.Second)), true
}

// EncodeISO converts an ISO value and mode to its wire value. Bits 0-23 hold
// the sensitivity, bits 24-27 the mode.
func EncodeISO(iso uint32, mode uint8) uint64 {
	return uint64(iso&0xFFFFFF) | uint64(mode&0x0F)<<24
}

// ISOAuto is the wire value for automatic sensitivity.
func ISOAuto() uint64 { return EncodeISO(ISOValueAuto, ISOModeNormal) }

// DecodeISO splits a wire value into its sensitivity and mode. auto reports
// whether the camera is choosing.
func DecodeISO(v uint64) (iso uint32, mode uint8, auto bool) {
	iso = uint32(v & 0xFFFFFF)
	mode = uint8((v >> 24) & 0x0F)
	return iso, mode, iso == ISOValueAuto
}

// EncodeExposureCompensation converts EV to its wire value, which is EV×1000,
// so +0.7 EV is 700 and -1.3 EV is -1300. The property is signed.
func EncodeExposureCompensation(ev float64) uint64 {
	return uint64(int64(math.Round(ev * 1000)))
}

// DecodeExposureCompensation converts a wire value back to EV.
func DecodeExposureCompensation(v uint64) float64 { return float64(int64(v)) / 1000 }

// Nearest returns the value in the camera's advertised set closest to want,
// comparing raw wire values. It exists because a camera rejects any value
// outside that set, and the set changes with lens, exposure mode and drive
// mode.
//
// DO NOT use it for shutter speed. Packed shutter values do not sort in time
// order: in the fraction form 1/N is 0x0001_000N, so a larger wire value means
// a SHORTER exposure, and 1/1000 sorts above 1/2. Use NearestShutter, which
// compares decoded durations.
//
// If the property advertises no set, want is returned unchanged.
func Nearest(d *DeviceProperty, want uint64) uint64 {
	if d == nil || len(d.Enum) == 0 {
		return want
	}
	best, bestDist := d.Enum[0], absDiff(d.Enum[0], want)
	for _, v := range d.Enum[1:] {
		if dist := absDiff(v, want); dist < bestDist {
			best, bestDist = v, dist
		}
	}
	return best
}

func absDiff(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

// SetAperture sets the f-number, e.g. 5.6.
//
// The camera must be in an exposure mode that gives the host aperture control
// (A or M); in P or S it will refuse or silently ignore the write.
func (c *Camera) SetAperture(f float64) error {
	return c.SetProperty(PropFNumber, ptp.TypeUint16, EncodeAperture(f))
}

// SetShutterSpeed sets the exposure time. A zero or negative duration selects
// Bulb.
//
// Requires an exposure mode with host shutter control (S or M). For exposures
// longer than the body's slowest discrete speed, use Bulb and BulbCapture.
//
// The duration is validated rather than clamped: a request the format cannot
// represent is an error, not a silently different exposure.
func (c *Camera) SetShutterSpeed(d time.Duration) error {
	if err := ValidateShutterSpeed(d); err != nil {
		return err
	}
	return c.SetProperty(PropShutterSpeed, ptp.TypeUint32, EncodeShutterSpeed(d))
}

// SetISO sets the sensitivity, e.g. 1600.
func (c *Camera) SetISO(iso uint32) error {
	return c.SetProperty(PropIsoSensitivity, ptp.TypeUint32, EncodeISO(iso, ISOModeNormal))
}

// SetISOAuto hands sensitivity back to the camera.
func (c *Camera) SetISOAuto() error {
	return c.SetProperty(PropIsoSensitivity, ptp.TypeUint32, ISOAuto())
}

// SetExposureCompensation sets exposure compensation in EV.
func (c *Camera) SetExposureCompensation(ev float64) error {
	return c.SetProperty(PropExposureBiasCompensation, ptp.TypeInt16, EncodeExposureCompensation(ev))
}

// FocusNearFar nudges focus. Negative moves nearer, positive further, and the
// magnitude is the step size — the camera accepts -7 to +7.
//
// This is the manual-focus primitive: for astronomical work it is how you walk
// focus around a star's minimum while the body stays in manual focus mode.
func (c *Camera) FocusNearFar(steps int16) error {
	if steps < -7 || steps > 7 {
		return fmt.Errorf("sony: focus step %d out of range, want -7..7", steps)
	}
	return c.SendControl(CtrlNearFar, ptp.TypeInt16, uint64(int64(steps)))
}

// DriveFocus starts a continuous focus drive and keeps it running until called
// again with 0. Negative drives wide, positive drives tele, and the magnitude
// is speed. The permitted range comes from the camera's focus speed range
// property, not from here.
func (c *Camera) DriveFocus(speed int8) error {
	return c.SendControl(CtrlFocusOperation, ptp.TypeInt8, uint64(int64(speed)))
}

// StopFocus halts a drive started by DriveFocus.
func (c *Camera) StopFocus() error { return c.DriveFocus(0) }

// Exposure is a snapshot of the exposure triangle, decoded.
type Exposure struct {
	Aperture      float64       // f-number; 0 if not reported
	ShutterSpeed  time.Duration // 0 when Bulb
	Bulb          bool
	ISO           uint32
	ISOAuto       bool
	Compensation  float64 // EV
	HasAperture   bool
	HasShutter    bool
	HasISO        bool
	HasCompensate bool
}

func (e Exposure) String() string {
	shutter := "?"
	switch {
	case e.Bulb:
		shutter = "Bulb"
	case e.HasShutter:
		shutter = e.ShutterSpeed.String()
	}
	iso := "?"
	if e.ISOAuto {
		iso = "AUTO"
	} else if e.HasISO {
		iso = fmt.Sprintf("%d", e.ISO)
	}
	ap := "?"
	if e.HasAperture {
		ap = fmt.Sprintf("f/%.1f", e.Aperture)
	}
	return fmt.Sprintf("%s  %s  ISO %s  %+.1f EV", ap, shutter, iso, e.Compensation)
}

// ReadExposure pulls the exposure triangle out of a property snapshot, so the
// caller does not repeat the decoding for each field.
func ReadExposure(props []DeviceProperty) Exposure {
	var e Exposure
	for i := range props {
		p := &props[i]
		switch p.Code {
		case PropFNumber:
			if f, ok := DecodeAperture(p.Current); ok {
				e.Aperture, e.HasAperture = f, true
			}
		case PropShutterSpeed:
			if d, ok := DecodeShutterSpeed(p.Current); ok {
				e.ShutterSpeed, e.HasShutter = d, true
			} else if uint32(p.Current) == uint32(ShutterBulb) {
				e.Bulb, e.HasShutter = true, true
			}
		case PropIsoSensitivity:
			e.ISO, _, e.ISOAuto = DecodeISO(p.Current)
			e.HasISO = true
		case PropExposureBiasCompensation:
			e.Compensation, e.HasCompensate = DecodeExposureCompensation(p.Current), true
		}
	}
	return e
}

// ShutterSpeeds returns the exposure times the camera currently accepts, in
// order from fastest to slowest. Values that are not times — Bulb, "nothing to
// display" — are omitted; HasBulb reports whether Bulb was among them.
//
// The camera's set is authoritative. It changes with exposure mode and drive
// mode, so read it after putting the body in the mode you intend to shoot in.
func ShutterSpeeds(d *DeviceProperty) (speeds []time.Duration, hasBulb bool) {
	if d == nil {
		return nil, false
	}
	for _, v := range d.Enum {
		if uint32(v) == uint32(ShutterBulb) {
			hasBulb = true
			continue
		}
		if t, ok := DecodeShutterSpeed(v); ok {
			speeds = append(speeds, t)
		}
	}
	sort.Slice(speeds, func(i, j int) bool { return speeds[i] < speeds[j] })
	return speeds, hasBulb
}

// NearestShutter returns the wire value for the advertised speed closest to
// want, compared as durations.
//
// This is the one to use for shutter speed. Comparing packed values instead
// would pick badly at the fast end, where the encoding runs backwards, and it
// also sidesteps the two display conventions: whether the body wants 1/2 as
// 0x00010002 or as 0x0005000A is its business, not the caller's, because the
// value comes from its own list.
//
// Closeness is measured in log space, so a request is snapped by ratio rather
// than by absolute difference — at 1/1000 a millisecond matters, at 10s it does
// not. If the property advertises no set, want is encoded directly.
func NearestShutter(d *DeviceProperty, want time.Duration) uint64 {
	if d == nil || len(d.Enum) == 0 || want <= 0 {
		return EncodeShutterSpeed(want)
	}
	best, bestDist := uint64(0), math.Inf(1)
	found := false
	target := math.Log(want.Seconds())
	for _, v := range d.Enum {
		t, ok := DecodeShutterSpeed(v)
		if !ok || t <= 0 {
			continue
		}
		if dist := math.Abs(math.Log(t.Seconds()) - target); dist < bestDist {
			best, bestDist, found = v, dist, true
		}
	}
	if !found {
		return EncodeShutterSpeed(want)
	}
	return best
}

// SetShutterSpeedFrom sets the exposure time, snapping to the nearest speed the
// camera advertises in d.
//
// Prefer this to SetShutterSpeed whenever a property snapshot is to hand: an
// exact value the body does not list is rejected, and the lists differ between
// bodies and modes.
func (c *Camera) SetShutterSpeedFrom(d *DeviceProperty, want time.Duration) error {
	return c.SetProperty(PropShutterSpeed, ptp.TypeUint32, NearestShutter(d, want))
}

// SetBulb puts the shutter speed into Bulb, for exposures beyond the longest
// discrete speed. Pair it with BulbCapture.
func (c *Camera) SetBulb() error {
	return c.SetProperty(PropShutterSpeed, ptp.TypeUint32, ShutterBulb)
}

// SetShutterFrom sets the exposure time from a property snapshot, snapping to
// the nearest speed the camera accepts.
//
// This is the one to reach for in a sequence. It does the three things that are
// easy to get wrong: it finds the property, it snaps by duration rather than by
// wire value, and it fails loudly when the body will not take the write. The
// snap is silent — the camera's list is authoritative, and ShutterOf reports
// what a request will become before committing to it.
//
// Camera.SetShutter is the same thing for a caller with no snapshot to hand, at
// the cost of fetching one.
func (c *Camera) SetShutterFrom(props []DeviceProperty, d time.Duration) error {
	if err := ValidateShutterSpeed(d); err != nil {
		return err
	}
	p := FindProp(props, PropShutterSpeed)
	if p == nil {
		return fmt.Errorf("sony: this body does not report shutter speed")
	}
	if !p.Writable() {
		return fmt.Errorf("sony: camera reports shutter speed as not settable right now; " +
			"the mode dial needs to be on S or M")
	}
	return c.SetProperty(PropShutterSpeed, ptp.TypeUint32, NearestShutter(p, d))
}

// ShutterOf returns the exposure time the camera would actually use for d,
// given a snapshot — the value SetShutter would send. Use it to check what a
// request will become before committing to it.
func ShutterOf(props []DeviceProperty, d time.Duration) (time.Duration, bool) {
	p := FindProp(props, PropShutterSpeed)
	if p == nil {
		return 0, false
	}
	return DecodeShutterSpeed(NearestShutter(p, d))
}
