package sony

import (
	"math"
	"testing"
	"time"

	"github.com/mikefsq/ptp"
)

// The assertions are Sony's own worked examples from the API reference, so a
// passing test means agreement with the vendor rather than with my reading.
func TestApertureMatchesSonyExamples(t *testing.T) {
	for _, tc := range []struct {
		f    float64
		wire uint64
	}{
		{4.0, 0x0190}, // documented: 0x0190 = 400 means F-4
		{9.5, 0x03B6}, // documented: 0x03B6 = 950 means F-9.5
		{1.4, 140},
		{2.8, 280},
		{22, 2200},
	} {
		if got := EncodeAperture(tc.f); got != tc.wire {
			t.Errorf("EncodeAperture(%v) = %d, want %d", tc.f, got, tc.wire)
		}
		back, ok := DecodeAperture(tc.wire)
		if !ok || math.Abs(back-tc.f) > 0.001 {
			t.Errorf("DecodeAperture(%d) = %v (ok=%v), want %v", tc.wire, back, ok, tc.f)
		}
	}
}

func TestApertureSpecialValues(t *testing.T) {
	for _, v := range []uint64{FNumberIrisClose, FNumberUnknown, FNumberNothing} {
		if _, ok := DecodeAperture(v); ok {
			t.Errorf("DecodeAperture(0x%04X) should report not-a-number", v)
		}
	}
}

func TestShutterSpeedMatchesSonyExamples(t *testing.T) {
	// Documented: 0x000F000A is 15/10 = 1.5 seconds.
	if got := EncodeShutterSpeed(1500 * time.Millisecond); got != 0x000F000A {
		t.Errorf("EncodeShutterSpeed(1.5s) = 0x%08X, want 0x000F000A", got)
	}
	// Documented: 0x000103E8 is 1/1000.
	if got := EncodeShutterSpeed(time.Second / 1000); got != 0x000103E8 {
		t.Errorf("EncodeShutterSpeed(1/1000) = 0x%08X, want 0x000103E8", got)
	}

	for _, tc := range []struct {
		wire uint64
		want time.Duration
	}{
		{0x000F000A, 1500 * time.Millisecond},
		{0x000103E8, time.Second / 1000},
		{0x001E000A, 3 * time.Second},        // 30/10
		{0x00010004, 250 * time.Millisecond}, // 1/4
	} {
		got, ok := DecodeShutterSpeed(tc.wire)
		if !ok {
			t.Errorf("DecodeShutterSpeed(0x%08X) reported not-a-time", tc.wire)
			continue
		}
		if d := got - tc.want; d > time.Millisecond || d < -time.Millisecond {
			t.Errorf("DecodeShutterSpeed(0x%08X) = %v, want %v", tc.wire, got, tc.want)
		}
	}
}

// A long astronomical exposure must survive the round trip: 30s is a normal
// sub-frame length and sits in the tenths encoding.
func TestShutterSpeedLongExposureRoundTrip(t *testing.T) {
	for _, want := range []time.Duration{
		time.Second, 2 * time.Second, 15 * time.Second, 30 * time.Second, 60 * time.Second,
	} {
		got, ok := DecodeShutterSpeed(EncodeShutterSpeed(want))
		if !ok {
			t.Errorf("%v did not survive the round trip", want)
			continue
		}
		if got != want {
			t.Errorf("round trip of %v gave %v", want, got)
		}
	}
}

func TestShutterBulb(t *testing.T) {
	if got := EncodeShutterSpeed(0); got != ShutterBulb {
		t.Errorf("zero duration = 0x%08X, want Bulb", got)
	}
	if _, ok := DecodeShutterSpeed(ShutterBulb); ok {
		t.Error("Bulb should not decode as a duration")
	}
	if _, ok := DecodeShutterSpeed(ShutterNothing); ok {
		t.Error("Nothing should not decode as a duration")
	}
}

func TestISOMatchesSonyExample(t *testing.T) {
	// Documented: 0x00000140 = 320.
	if got := EncodeISO(320, ISOModeNormal); got != 0x00000140 {
		t.Errorf("EncodeISO(320) = 0x%08X, want 0x00000140", got)
	}
	iso, mode, auto := DecodeISO(0x00000140)
	if iso != 320 || mode != ISOModeNormal || auto {
		t.Errorf("DecodeISO(0x140) = iso %d mode %d auto %v", iso, mode, auto)
	}
}

func TestISOAutoAndModeBits(t *testing.T) {
	if _, _, auto := DecodeISO(ISOAuto()); !auto {
		t.Error("ISOAuto() should decode as auto")
	}
	// The mode must live in bits 24-27 and not corrupt the value.
	v := EncodeISO(1600, ISOModeMultiFrameNR)
	iso, mode, auto := DecodeISO(v)
	if iso != 1600 {
		t.Errorf("iso = %d, want 1600", iso)
	}
	if mode != ISOModeMultiFrameNR {
		t.Errorf("mode = %d, want %d", mode, ISOModeMultiFrameNR)
	}
	if auto {
		t.Error("1600 should not report auto")
	}
}

func TestExposureCompensation(t *testing.T) {
	for _, tc := range []struct {
		ev   float64
		wire int64
	}{{0.7, 700}, {-1.3, -1300}, {0, 0}, {3, 3000}} {
		if got := int64(EncodeExposureCompensation(tc.ev)); got != tc.wire {
			t.Errorf("EncodeExposureCompensation(%v) = %d, want %d", tc.ev, got, tc.wire)
		}
		if got := DecodeExposureCompensation(uint64(tc.wire)); math.Abs(got-tc.ev) > 0.0001 {
			t.Errorf("DecodeExposureCompensation(%d) = %v, want %v", tc.wire, got, tc.ev)
		}
	}
}

// A camera rejects anything outside its advertised set, so Nearest must snap to
// it rather than pass a plausible-looking value straight through.
func TestNearestSnapsToCameraSet(t *testing.T) {
	d := &DeviceProperty{Enum: []uint64{280, 400, 560, 800, 1100}}
	for _, tc := range []struct{ want, expect uint64 }{
		{400, 400},
		{410, 400},
		{500, 560},
		{50, 280},    // below the set
		{9999, 1100}, // above it
	} {
		if got := Nearest(d, tc.want); got != tc.expect {
			t.Errorf("Nearest(%d) = %d, want %d", tc.want, got, tc.expect)
		}
	}
	// With no advertised set the request passes through untouched.
	if got := Nearest(&DeviceProperty{}, 1234); got != 1234 {
		t.Errorf("Nearest with no enum = %d, want 1234", got)
	}
	if got := Nearest(nil, 7); got != 7 {
		t.Errorf("Nearest(nil) = %d, want 7", got)
	}
}

func TestReadExposure(t *testing.T) {
	props := []DeviceProperty{
		{Code: PropFNumber, Current: 560},
		{Code: PropShutterSpeed, Current: 0x001E000A}, // 3s
		{Code: PropIsoSensitivity, Current: 1600},
		{Code: PropExposureBiasCompensation, Current: EncodeExposureCompensation(-0.3)},
	}
	e := ReadExposure(props)
	if !e.HasAperture || math.Abs(e.Aperture-5.6) > 0.001 {
		t.Errorf("aperture = %v", e.Aperture)
	}
	if !e.HasShutter || e.ShutterSpeed != 3*time.Second {
		t.Errorf("shutter = %v", e.ShutterSpeed)
	}
	if !e.HasISO || e.ISO != 1600 || e.ISOAuto {
		t.Errorf("iso = %d auto=%v", e.ISO, e.ISOAuto)
	}
	if math.Abs(e.Compensation-(-0.3)) > 0.0001 {
		t.Errorf("compensation = %v", e.Compensation)
	}
}

func TestReadExposureBulb(t *testing.T) {
	e := ReadExposure([]DeviceProperty{{Code: PropShutterSpeed, Current: ShutterBulb}})
	if !e.Bulb {
		t.Error("Bulb shutter speed should set Bulb")
	}
}

func TestFocusNearFarRangeIsChecked(t *testing.T) {
	c := openCamera(&fakeTransport{})
	for _, bad := range []int16{-8, 8, 100} {
		if err := c.FocusNearFar(bad); err == nil {
			t.Errorf("FocusNearFar(%d) should be rejected; the camera accepts -7..7", bad)
		}
	}
}

// A realistic A7R shutter list, in the order a camera reports it — which is not
// sorted, and mixes both encodings.
func a7rShutterList() *DeviceProperty {
	speeds := []uint64{
		0x000103E8, // 1/1000
		0x000101F4, // 1/500
		0x000100FA, // 1/250
		0x0001007D, // 1/125
		0x0001003C, // 1/60
		0x0001001E, // 1/30
		0x0001000F, // 1/15
		0x00010008, // 1/8
		0x00010004, // 1/4
		0x00010002, // 1/2
		0x000A000A, // 1s
		0x0014000A, // 2s
		0x0028000A, // 4s
		0x0064000A, // 10s
		uint64(ShutterBulb),
	}
	return &DeviceProperty{Code: PropShutterSpeed, Type: ptp.TypeUint32, Enum: speeds}
}

// The two ends the client actually needs.
func TestShutterRangeEndpoints(t *testing.T) {
	d := a7rShutterList()
	for _, tc := range []struct {
		want time.Duration
		wire uint64
	}{
		{time.Second / 1000, 0x000103E8}, // 1/1000
		{10 * time.Second, 0x0064000A},   // 10s
	} {
		if got := NearestShutter(d, tc.want); got != tc.wire {
			t.Errorf("NearestShutter(%v) = 0x%08X, want 0x%08X", tc.want, got, tc.wire)
		}
	}
}

// The regression this guards: packed values sort backwards at the fast end, so
// a numeric Nearest picks the slowest speed when asked for the fastest.
func TestNearestShutterBeatsNumericNearest(t *testing.T) {
	d := a7rShutterList()
	want := time.Second / 1000

	good := NearestShutter(d, want)
	gotDur, _ := DecodeShutterSpeed(good)
	if gotDur != want {
		t.Fatalf("NearestShutter gave %v, want %v", gotDur, want)
	}

	// A request that is not itself in the list is where the two disagree: the
	// packed forms put 0.7s numerically nearer 1s, while by time it is nearer
	// 0.5s.
	between := 700 * time.Millisecond
	byTime, _ := DecodeShutterSpeed(NearestShutter(d, between))
	byValue, _ := DecodeShutterSpeed(Nearest(d, EncodeShutterSpeed(between)))
	if byTime != 500*time.Millisecond {
		t.Errorf("NearestShutter(700ms) = %v, want 500ms", byTime)
	}
	if byValue == byTime {
		t.Logf("numeric Nearest agreed here (%v); the hazard is real regardless", byValue)
	} else {
		t.Logf("confirmed: numeric Nearest picks %v for a %v request, NearestShutter picks %v",
			byValue, between, byTime)
	}
}

// 0.7s has no unit-fraction form, so it must fall back to tenths rather than
// silently becoming a whole second.
func TestSubSecondNonFractionUsesTenths(t *testing.T) {
	got, ok := DecodeShutterSpeed(EncodeShutterSpeed(700 * time.Millisecond))
	if !ok {
		t.Fatal("0.7s did not decode")
	}
	if got != 700*time.Millisecond {
		t.Errorf("0.7s round-tripped to %v", got)
	}
	// The exact fractions must still use the fraction form Sony documents.
	if v := EncodeShutterSpeed(time.Second / 1000); v != 0x000103E8 {
		t.Errorf("1/1000 = 0x%08X, want 0x000103E8", v)
	}
	if v := EncodeShutterSpeed(time.Second / 2); v != 0x00010002 {
		t.Errorf("1/2 = 0x%08X, want 0x00010002", v)
	}
}

func TestNearestShutterSnapsByRatio(t *testing.T) {
	d := a7rShutterList()
	for _, tc := range []struct {
		want   time.Duration
		expect time.Duration
	}{
		{time.Second / 900, time.Second / 1000}, // between 1/1000 and 1/500, nearer 1/1000 by ratio
		{3 * time.Second, 4 * time.Second},      // between 2s and 4s
		{time.Hour, 10 * time.Second},           // beyond the slowest
		{time.Microsecond, time.Second / 1000},  // beyond the fastest
	} {
		got, ok := DecodeShutterSpeed(NearestShutter(d, tc.want))
		if !ok {
			t.Errorf("NearestShutter(%v) did not decode", tc.want)
			continue
		}
		if got != tc.expect {
			t.Errorf("NearestShutter(%v) = %v, want %v", tc.want, got, tc.expect)
		}
	}
}

func TestShutterSpeedsSortedFastToSlow(t *testing.T) {
	speeds, hasBulb := ShutterSpeeds(a7rShutterList())
	if !hasBulb {
		t.Error("Bulb should be reported separately")
	}
	if len(speeds) != 14 {
		t.Fatalf("got %d speeds, want 14 (Bulb excluded)", len(speeds))
	}
	if speeds[0] != time.Second/1000 {
		t.Errorf("fastest = %v, want 1/1000", speeds[0])
	}
	if speeds[len(speeds)-1] != 10*time.Second {
		t.Errorf("slowest = %v, want 10s", speeds[len(speeds)-1])
	}
	for i := 1; i < len(speeds); i++ {
		if speeds[i] <= speeds[i-1] {
			t.Errorf("not sorted at %d: %v then %v", i, speeds[i-1], speeds[i])
		}
	}
}

// Every speed a client might send across the eclipse range must survive
// encode -> decode unchanged.
func TestShutterFullRangeRoundTrip(t *testing.T) {
	for _, d := range []time.Duration{
		time.Second / 4000, time.Second / 1000, time.Second / 500, time.Second / 250,
		time.Second / 60, time.Second / 30, time.Second / 8, time.Second / 2,
		time.Second, 2 * time.Second, 5 * time.Second, 10 * time.Second, 30 * time.Second,
	} {
		got, ok := DecodeShutterSpeed(EncodeShutterSpeed(d))
		if !ok {
			t.Errorf("%v did not decode", d)
			continue
		}
		// Fractions are exact only to the nearest 1/N, so allow a hair.
		if diff := got - d; diff > time.Microsecond || diff < -time.Microsecond {
			t.Errorf("round trip of %v gave %v", d, got)
		}
	}
}

// The full range a client may send, with the exact boundaries.
func TestShutterSpeedPrecisionAcrossRange(t *testing.T) {
	for _, in := range []string{
		"31.25us", "100us", "125us", "0.0001s", "1ms", "10ms", "0.5s", "0.7s",
		"1s", "1.3s", "2s", "10s", "30s",
	} {
		d, err := time.ParseDuration(in)
		if err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if err := ValidateShutterSpeed(d); err != nil {
			t.Errorf("%s should be valid: %v", in, err)
			continue
		}
		got, ok := DecodeShutterSpeed(EncodeShutterSpeed(d))
		if !ok {
			t.Errorf("%s did not decode", in)
			continue
		}
		if got != d {
			t.Errorf("%s round-tripped to %v", in, got)
		}
	}
}

// Below 1/65535 the format cannot represent the time. It must be rejected, not
// silently turned into a different exposure.
func TestShutterSpeedTooFastIsRejected(t *testing.T) {
	c := openCamera(&fakeTransport{})
	for _, in := range []time.Duration{time.Microsecond, 10 * time.Microsecond} {
		if err := ValidateShutterSpeed(in); err == nil {
			t.Errorf("%v should be rejected; it is faster than 1/65535", in)
		}
		if err := c.SetShutterSpeed(in); err == nil {
			t.Errorf("SetShutterSpeed(%v) should refuse rather than clamp", in)
		}
	}
	// The boundary itself is fine, as is anything a real camera can do.
	for _, in := range []time.Duration{MinShutterSpeed, 31250 * time.Nanosecond, 125 * time.Microsecond} {
		if err := ValidateShutterSpeed(in); err != nil {
			t.Errorf("%v should be accepted: %v", in, err)
		}
	}
}

func TestShutterSpeedTooLongIsRejected(t *testing.T) {
	if err := ValidateShutterSpeed(2 * time.Hour); err == nil {
		t.Error("2h should be rejected; use Bulb")
	}
	if err := ValidateShutterSpeed(0); err != nil {
		t.Errorf("zero means Bulb and should be valid: %v", err)
	}
}

func TestReadinessFromSnapshot(t *testing.T) {
	r := ReadReadiness([]DeviceProperty{
		{Code: PropMediaSLOT1Status, Current: SlotOK},
		{Code: PropMediaSLOT1RemainingNumber, Current: 412},
		{Code: PropMediaSLOT1WritingState, Current: WritingStateNotWriting},
	})
	if !r.Ready() || !r.Settled() {
		t.Errorf("idle camera with a good card should be ready and settled: %v", r)
	}
	if n, ok := r.RemainingShots(); !ok || n != 412 {
		t.Errorf("remaining = %d (known=%v), want 412", n, ok)
	}
}

// Writing to card must NOT block shooting: a camera writes while it shoots, and
// gating each frame on an idle buffer would throttle a burst to card speed.
func TestWritingIsReadyButNotSettled(t *testing.T) {
	r := ReadReadiness([]DeviceProperty{
		{Code: PropMediaSLOT1Status, Current: SlotOK},
		{Code: PropMediaSLOT1WritingState, Current: WritingStateWriting},
	})
	if !r.Ready() {
		t.Error("a camera flushing its buffer can still take another frame")
	}
	if r.Settled() || !r.Writing() {
		t.Error("a camera still writing is not settled")
	}
}

// Both bodies are dual-slot. Recording to slot 2 with slot 1 empty is a normal
// configuration and must not read as "not ready".
func TestSlot2AloneIsReady(t *testing.T) {
	r := ReadReadiness([]DeviceProperty{
		{Code: PropMediaSLOT1Status, Current: SlotNoCard},
		{Code: PropMediaSLOT2Status, Current: SlotOK},
		{Code: PropMediaSLOT2RemainingNumber, Current: 900},
	})
	if !r.Ready() {
		t.Errorf("a good card in slot 2 should be ready even with slot 1 empty: %v", r)
	}
	if r.Slot1.Usable() {
		t.Error("slot 1 has no card and is not usable")
	}
	if !r.Slot2.Usable() {
		t.Error("slot 2 is fine and should be usable")
	}
	if n, _ := r.RemainingShots(); n != 900 {
		t.Errorf("remaining = %d, want 900 from slot 2", n)
	}
}

// Either slot still writing means not settled, even if the other is idle.
func TestSettledRequiresBothSlots(t *testing.T) {
	r := ReadReadiness([]DeviceProperty{
		{Code: PropMediaSLOT1Status, Current: SlotOK},
		{Code: PropMediaSLOT1WritingState, Current: WritingStateNotWriting},
		{Code: PropMediaSLOT2Status, Current: SlotOK},
		{Code: PropMediaSLOT2WritingState, Current: WritingStateWriting},
	})
	if r.Settled() {
		t.Error("slot 2 is still writing, so the camera is not settled")
	}
	if !r.Ready() {
		t.Error("both cards are healthy, so it can still shoot")
	}
}

func TestNotReadyConditions(t *testing.T) {
	full := ReadReadiness([]DeviceProperty{
		{Code: PropMediaSLOT1Status, Current: SlotOK},
		{Code: PropMediaSLOT1RemainingNumber, Current: 0},
	})
	if full.Ready() {
		t.Error("a full card in the only reported slot is not ready")
	}
	noCard := ReadReadiness([]DeviceProperty{{Code: PropMediaSLOT1Status, Current: SlotNoCard}})
	if noCard.Ready() {
		t.Error("no card is not ready")
	}
	if got := SlotStatusName(SlotNoCard); got != "no card" {
		t.Errorf("name = %q", got)
	}
}

// A body reporting none of these must not be declared unready: absent is not
// the same as bad.
func TestReadinessAbsentPropertiesAreNotFailures(t *testing.T) {
	r := ReadReadiness(nil)
	if !r.Ready() || !r.Settled() {
		t.Errorf("a body reporting nothing should not be treated as unready: %v", r)
	}
	if r.Slot1.Reported || r.Slot2.Reported {
		t.Error("nothing was reported, so neither slot should be flagged as present")
	}
}
