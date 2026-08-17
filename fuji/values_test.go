package fuji

import (
	"testing"

	"github.com/mikefsq/ptp"
)

// The values named here are cross-checked against the SDK's numerically
// documented constants. For these three properties every value the X-T5
// advertises is one of them — 9 of 9 — which is what makes the table
// trustworthy rather than plausible.
func TestDocumentedValueNames(t *testing.T) {
	cases := []struct {
		p    ptp.Prop
		v    uint64
		want string
	}{
		{PropPriorityModeCode, 1, "Camera"},
		{PropPriorityModeCode, 2, "PC"},
		{PropMediaRecordCode, 1, "RAW+JPEG"},
		{PropMediaRecordCode, 4, "Off"},
		// Sparse and non-sequential: naming by ordinal position would put
		// "ContinuousH" here. The camera reported 8 with its collar on movie.
		{PropDriveMode, 4, "Single"},
		{PropDriveMode, 8, "Movie"},
		{PropDriveMode, 16, "PixelShift"},
		{PropRawCompression, 2, "LosslessCompressed"},
		{PropQuality, 1, "RAW"},
		{ptp.PropFocusMode, 1, "Manual"},
		{ptp.PropFocusMode, 0x8001, "AF-S"},
		{ptp.PropExposureProgram, 1, "Manual"},
	}
	for _, c := range cases {
		if got := ValueName(c.p, c.v); got != c.want {
			t.Errorf("ValueName(0x%04X, %d) = %q, want %q", uint16(c.p), c.v, got, c.want)
		}
	}
}

// An unestablished value must print as a number. A plausible wrong name is
// worse than an honest number, and the SDK's enum ORDER does not give the
// values — checked against hardware, that guess disagrees for 10 of 25
// properties.
func TestUnknownValuesStayNumeric(t *testing.T) {
	if got := ValueName(PropDriveMode, 0x99); got != "153" {
		t.Errorf("an undocumented DriveMode value should print as a number, got %q", got)
	}
	if got := ValueName(0xD999, 7); got != "7" {
		t.Errorf("an unknown property's value should print as a number, got %q", got)
	}
}

func TestDescribeValues(t *testing.T) {
	got := DescribeValues(PropDriveMode, []uint64{4, 5, 16})
	want := "Single, MultiExposure, PixelShift"
	if got != want {
		t.Errorf("DescribeValues = %q, want %q", got, want)
	}
	// Unknown values still render, as numbers.
	if got := DescribeValues(0xD999, []uint64{1, 2}); got != "1, 2" {
		t.Errorf("unknown property = %q, want %q", got, "1, 2")
	}
}
