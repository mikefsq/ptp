package sony

import (
	"strings"
	"testing"
	"time"

	"github.com/mikefsq/ptp"
)

// CrFileType puts JPEG at 1 and RAW at 2 — the reverse of the obvious order.
// This was got wrong once already, and getting it wrong means shooting JPEG for
// a whole sequence while the driver reports RAW. The header is the authority.
func TestFileTypeConstantsMatchTheSDK(t *testing.T) {
	for _, tc := range []struct {
		v    uint64
		want string
	}{
		{0x0000, "None"},
		{0x0001, "JPEG"},
		{0x0002, "RAW"},
		{0x0003, "RAW+JPEG"},
		{0x0004, "RAW+HEIF"},
		{0x0005, "HEIF"},
	} {
		if got := ValueName(PropFileType, tc.v); got != tc.want {
			t.Errorf("FileType 0x%04X = %q, want %q", tc.v, got, tc.want)
		}
	}
	if FileTypeJPEG != 0x0001 || FileTypeRAW != 0x0002 {
		t.Errorf("JPEG=%d RAW=%d — the SDK has JPEG=1, RAW=2", FileTypeJPEG, FileTypeRAW)
	}
}

// Enums whose size is easy to over-guess. Both have exactly four values, and a
// name for a fifth would mean one was invented.
func TestNoInventedEnumValues(t *testing.T) {
	for _, tc := range []struct {
		name  string
		p     Prop
		valid []uint64
	}{
		{"ImageQuality", PropStillImageQuality, []uint64{0, 1, 2, 3, 4}},
		{"AspectRatio", PropAspectRatio, []uint64{1, 2, 3, 4}},
		{"ImageSize", PropImageSize, []uint64{1, 2, 3, 4}},
	} {
		named := valueNames[tc.p]
		if len(named) != len(tc.valid) {
			t.Errorf("%s names %d values, want %d — an extra name means one was invented",
				tc.name, len(named), len(tc.valid))
		}
		for _, v := range tc.valid {
			if _, ok := named[v]; !ok {
				t.Errorf("%s value %d has no name", tc.name, v)
			}
		}
	}
}

// Drive mode and recording media have large gaps in their numbering — Single is
// 0x00000001 while the continuous family starts at 0x00010001. Naming by
// position would be wrong for every continuous mode.
func TestSparseEnumsAreNamedByValue(t *testing.T) {
	if got := ValueName(PropDriveMode, DriveSingle); got != "Single" {
		t.Errorf("DriveSingle = %q", got)
	}
	if got := ValueName(PropDriveMode, DriveContinuousHi); got != "ContinuousHi" {
		t.Errorf("DriveContinuousHi = %q", got)
	}
	if got := ValueName(PropRecordingMedia, RecordingMediaSimultaneous); got != "Simultaneous" {
		t.Errorf("RecordingMediaSimultaneous = %q", got)
	}
}

// The exposure triangle is ENCODED, not enumerated. Printing the raw wire value
// is actively misleading: 1/1000 packs to 0x000103E8, which reads as 66536.
func TestEncodedValuesAreDecodedNotPrinted(t *testing.T) {
	if got := ValueName(PropShutterSpeed, EncodeShutterSpeed(time.Second/1000)); got != "1ms" {
		t.Errorf("shutter 1/1000 = %q, want 1ms", got)
	}
	if got := ValueName(PropShutterSpeed, ShutterBulb); got != "Bulb" {
		t.Errorf("bulb = %q", got)
	}
	if got := ValueName(PropFNumber, EncodeAperture(5.6)); got != "f/5.6" {
		t.Errorf("aperture = %q, want f/5.6", got)
	}
	if got := ValueName(PropIsoSensitivity, ISOAuto()); got != "AUTO" {
		t.Errorf("iso auto = %q", got)
	}
	if got := ValueName(PropIsoSensitivity, EncodeISO(1600, ISOModeNormal)); got != "1600" {
		t.Errorf("iso 1600 = %q", got)
	}
	if got := ValueName(PropExposureBiasCompensation, EncodeExposureCompensation(-0.3)); got != "-0.3 EV" {
		t.Errorf("exposure comp = %q, want -0.3 EV", got)
	}
}

// A signed value printed unsigned comes out as 18446744073709551615 and looks
// like corruption rather than a negative number.
func TestNegativeValuesPrintSigned(t *testing.T) {
	if got := ValueName(Prop(0xDEAD), uint64(0xFFFFFFFFFFFFFFFF)); got != "-1" {
		t.Errorf("got %q, want -1", got)
	}
}

// An unknown value must come back as a number, not a guess.
func TestUnknownValuesStayNumeric(t *testing.T) {
	if got := ValueName(PropFileType, 0x0099); got != "153" {
		t.Errorf("got %q, want the bare number", got)
	}
}

func TestDescribeValuesTruncatesLongSets(t *testing.T) {
	vs := make([]uint64, 100)
	got := DescribeValues(PropImageSize, vs)
	if !strings.Contains(got, "more)") {
		t.Errorf("a 100-value set was not truncated: %q", got)
	}
}

// Describe must show a string property's text, not its numeric Current, which
// for a string entry is meaningless.
func TestDescribeRendersStringProperties(t *testing.T) {
	d := &DeviceProperty{
		Code: PropSetCopyright, Type: ptp.TypeString,
		CurrentStr: "M. Furlotti", GetSet: 0x01, Enabled: 1,
	}
	got := Describe(d)
	if !strings.Contains(got, `"M. Furlotti"`) {
		t.Errorf("got %q, want the string value", got)
	}
	if !strings.Contains(got, "writable") {
		t.Errorf("got %q, want the access noted", got)
	}
}

// SlotStatusName and ValueName must agree, because they read one table. The
// register is prose because SlotStatusName's output lands in a user-facing
// error.
func TestSlotStatusNamingIsShared(t *testing.T) {
	if got := SlotStatusName(SlotNoCard); got != "no card" {
		t.Errorf("SlotStatusName = %q", got)
	}
	if got := ValueName(PropMediaSLOT1Status, SlotNoCard); got != "no card" {
		t.Errorf("ValueName = %q; the two must read the same table", got)
	}
}
