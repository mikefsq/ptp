package fuji

import (
	"github.com/mikefsq/ptp"
	"testing"
	"time"
)

// The X-T5 carries shutter speed as a plain microsecond count. The wire values
// are exact powers of two where the dial shows the conventional rounded number,
// so 32000000 is 32 seconds and reads "30s" on the body.
func TestShutterMicrosecondEncoding(t *testing.T) {
	for _, tc := range []struct {
		wire uint64
		want time.Duration
		dial string
	}{
		{32000000, 32 * time.Second, "30s"}, // observed on a real X-T5
		{16000000, 16 * time.Second, "15s"},
		{64000000, 64 * time.Second, "60s"},
		{10079368, 10079368 * time.Microsecond, "10s"},
		{1000, time.Millisecond, "1/1000"},
		{250, 250 * time.Microsecond, "1/4000"},
	} {
		if got := DecodeShutter(tc.wire); got != tc.want {
			t.Errorf("DecodeShutter(%d) = %v, want %v (dial shows %s)", tc.wire, got, tc.want, tc.dial)
		}
		if got := EncodeShutter(tc.want); got != tc.wire {
			t.Errorf("EncodeShutter(%v) = %d, want %d", tc.want, got, tc.wire)
		}
	}
}

// The range a client needs for an eclipse: 1/4000 up to tens of seconds.
func TestShutterRangeRoundTrip(t *testing.T) {
	for _, d := range []time.Duration{
		250 * time.Microsecond, time.Millisecond, 10 * time.Millisecond,
		time.Second / 2, time.Second, 10 * time.Second, 30 * time.Second,
	} {
		if err := ValidateShutter(d); err != nil {
			t.Errorf("%v should be valid: %v", d, err)
			continue
		}
		if got := DecodeShutter(EncodeShutter(d)); got != d {
			t.Errorf("round trip of %v gave %v", d, got)
		}
	}
	if err := ValidateShutter(500 * time.Nanosecond); err == nil {
		t.Error("sub-microsecond should be rejected, not silently truncated")
	}
	if err := ValidateShutter(2 * time.Hour); err == nil {
		t.Error("beyond the uint32 microsecond range should be rejected")
	}
}

// A body with its shutter dial on a fixed speed reports exactly one accepted
// value. That is normal, and snapping must return it rather than the request.
func TestNearestShutterSnapsToDialLockedValue(t *testing.T) {
	d := &ptp.StdPropDesc{Code: ptp.PropExposureTime, Type: ptp.TypeUint32, Enum: []uint64{32000000}}
	if got := NearestShutter(d, time.Millisecond); got != 32000000 {
		t.Errorf("got %d, want the single advertised value 32000000", got)
	}
}

func TestNearestShutterPicksByRatio(t *testing.T) {
	d := &ptp.StdPropDesc{Code: ptp.PropExposureTime, Type: ptp.TypeUint32,
		Enum: []uint64{250, 1000, 4000, 16000, 1000000, 16000000, 32000000}}
	for _, tc := range []struct {
		want   time.Duration
		expect time.Duration
	}{
		{time.Millisecond, time.Millisecond},
		{900 * time.Microsecond, time.Millisecond},
		{20 * time.Second, 16 * time.Second}, // nearer by ratio than 32s
		{time.Hour, 32 * time.Second},
		{time.Nanosecond, 250 * time.Microsecond},
	} {
		if got := DecodeShutter(NearestShutter(d, tc.want)); got != tc.expect {
			t.Errorf("NearestShutter(%v) = %v, want %v", tc.want, got, tc.expect)
		}
	}
}

func TestShutterSpeedsSorted(t *testing.T) {
	d := &ptp.StdPropDesc{Enum: []uint64{32000000, 250, 1000000, 1000}}
	got := ShutterSpeeds(d)
	if len(got) != 4 || got[0] != 250*time.Microsecond || got[3] != 32*time.Second {
		t.Errorf("ShutterSpeeds = %v, want fastest first", got)
	}
}

// Aperture uses the same F×100 convention as Sony; f/4 read back as 400 from a
// real X-T5.
func TestApertureEncoding(t *testing.T) {
	if got := EncodeAperture(4.0); got != 400 {
		t.Errorf("EncodeAperture(4.0) = %d, want 400 (observed on an X-T5)", got)
	}
	if got := DecodeAperture(400); got != 4.0 {
		t.Errorf("DecodeAperture(400) = %v, want 4.0", got)
	}
	if got := EncodeAperture(2.8); got != 280 {
		t.Errorf("EncodeAperture(2.8) = %d, want 280", got)
	}
}
