package fuji

import (
	"strings"
	"testing"

	"github.com/mikefsq/ptp"
)

// The generated table is real hardware, so it should agree with what the driver
// learned the hard way.
func TestGeneratedTableMatchesKnownFacts(t *testing.T) {
	cases := []struct {
		code ptp.Prop
		name string
		typ  ptp.DataType
	}{
		{PropQuality, "Quality", ptp.TypeUint16},
		{PropRawCompression, "RawCompression", ptp.TypeUint16},
		{PropMediaRecordCode, "MediaRecord", ptp.TypeUint16},
		// ExposureIndex is UINT16 on this body, not the UINT32 gphoto2 implies.
		// Getting it wrong means the camera rejects the write.
		{ptp.PropISO, "ExposureIndex", ptp.TypeInt32},
		{ptp.PropExposureTime, "ExposureTime", ptp.TypeUint32},
	}
	for _, tc := range cases {
		info, ok := Known(tc.code)
		if !ok {
			t.Errorf("0x%04X (%s) missing from the captured table", uint16(tc.code), tc.name)
			continue
		}
		if info.Type != tc.typ {
			t.Errorf("%s type = 0x%04X, want 0x%04X", tc.name, uint16(info.Type), uint16(tc.typ))
		}
	}
}

// A single advertised value means the camera owns the setting. Saying "invalid
// value" there sends the user looking for a different number, when the fix is a
// dial on the body.
func TestSingleValuedPropertyReportsCameraControl(t *testing.T) {
	d := &ptp.StdPropDesc{Code: ptp.PropFNumber, GetSet: 1, Enum: []uint64{400}}
	err := checkValue(ptp.PropFNumber, d, 800)
	if err == nil {
		t.Fatal("writing an unoffered value must be refused")
	}
	if !strings.Contains(err.Error(), "camera-controlled") {
		t.Errorf("error should explain the camera owns it, got: %v", err)
	}
	if err := checkValue(ptp.PropFNumber, d, 400); err != nil {
		t.Errorf("the one offered value must be accepted: %v", err)
	}
}

func TestEnumAndRangeValidation(t *testing.T) {
	enum := &ptp.StdPropDesc{Code: PropQuality, GetSet: 1, Enum: []uint64{1, 2, 3, 4, 5}}
	if err := checkValue(PropQuality, enum, 3); err != nil {
		t.Errorf("3 is offered: %v", err)
	}
	if err := checkValue(PropQuality, enum, 9); err == nil {
		t.Error("9 is not offered and must be refused")
	}

	rng := &ptp.StdPropDesc{Code: 0xD00B, GetSet: 1, Form: ptp.FormRange, Min: 0, Max: 9, Step: 1}
	if err := checkValue(0xD00B, rng, 5); err != nil {
		t.Errorf("5 is in range: %v", err)
	}
	if err := checkValue(0xD00B, rng, 20); err == nil {
		t.Error("20 is out of range and must be refused")
	}
}

func TestPropByName(t *testing.T) {
	p, ok := PropByName("rawcompression")
	if !ok || p != PropRawCompression {
		t.Errorf("PropByName(rawcompression) = 0x%04X, %v", uint16(p), ok)
	}
	if _, ok := PropByName("NoSuchProperty"); ok {
		t.Error("an unknown name must not resolve")
	}
}

// Every captured property must name a real data type; a zero would mean the
// generator emitted something the parser did not understand.
func TestGeneratedTableIsWellFormed(t *testing.T) {
	if len(XT5Props) < 100 {
		t.Fatalf("only %d properties in the table", len(XT5Props))
	}
	seen := map[ptp.Prop]bool{}
	for _, info := range XT5Props {
		if seen[info.Code] {
			t.Errorf("0x%04X appears twice", uint16(info.Code))
		}
		seen[info.Code] = true
		if info.Type == 0 && len(info.Values) > 0 {
			t.Errorf("0x%04X has values but no type", uint16(info.Code))
		}
	}
}

// Property names must resolve from the body's own plugin table, not only from
// gphoto2's. gphoto2 calls 0xD201 ReleaseMode; the X-T5 calls it DriveMode, and
// a tool that only knew the gphoto2 name would fail to find it.
func TestPropByNameUsesPluginNames(t *testing.T) {
	p, ok := PropByName("DriveMode")
	if !ok {
		t.Fatal("DriveMode should resolve: it is the name the X-T5's own plugin uses")
	}
	if p != PropDriveMode {
		t.Errorf("DriveMode = 0x%04X, want 0x%04X", uint16(p), uint16(PropDriveMode))
	}
	// The gphoto2 name for the same code must still work.
	if q, ok := PropByName("ReleaseMode"); !ok || q != PropDriveMode {
		t.Errorf("ReleaseMode = 0x%04X, %v; want the same code", uint16(q), ok)
	}
}

// The plugin's name wins where the two sources disagree: gphoto2 calls 0xD028
// CommandDialMode, the X-T5 itself calls it DOFScale.
func TestPluginNameOutranksGphoto2(t *testing.T) {
	if got := PropName(0xD028); got != "DOFScale" {
		t.Errorf("PropName(0xD028) = %q, want DOFScale (the body's own name)", got)
	}
	if got := PropName(0x500F); got != "Sensitivity" {
		t.Errorf("PropName(0x500F) = %q, want Sensitivity", got)
	}
}
