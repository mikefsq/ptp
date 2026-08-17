package sony

import (
	"os"
	"testing"

	"github.com/mikefsq/ptp"
)

// Both fixtures are verbatim 0x1014 replies captured from a real Sony NEX-6 on
// 2026-08-06.
//
// These are the STANDARD descriptor layout — no isEnabled byte — parsed by the
// parent package. Sony's own 0x9209 entries are a different shape and are
// parsed in devprop.go; keeping both sets of golden bytes in this package is
// what makes that contrast checkable.
func TestParseStdPropDescBatteryLevel(t *testing.T) {
	raw, err := os.ReadFile("testdata/nex6-propdesc-5001.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	d, err := ptp.ParseStdPropDesc(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if d.Code != 0x5001 {
		t.Errorf("code = 0x%04X, want 0x5001", uint16(d.Code))
	}
	if d.Type != ptp.TypeUint8 {
		t.Errorf("type = 0x%04X, want UInt8", uint16(d.Type))
	}
	if d.Writable() {
		t.Error("battery level should be read-only")
	}
	if d.Form != ptp.FormEnum {
		t.Errorf("form = %d, want enum", d.Form)
	}
	// The camera advertises a 62-entry enum, one per battery step.
	if len(d.Enum) != 62 {
		t.Errorf("got %d enum values, want 62", len(d.Enum))
	}
}

// The string case is the one most likely to be got wrong: PTP strings are a
// uint8 character count followed by UTF-16, with the NUL inside the count.
func TestParseStdPropDescFriendlyName(t *testing.T) {
	raw, err := os.ReadFile("testdata/nex6-propdesc-d402.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	d, err := ptp.ParseStdPropDesc(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if d.Code != 0xD402 {
		t.Errorf("code = 0x%04X, want 0xD402", uint16(d.Code))
	}
	if d.Type != ptp.TypeString {
		t.Errorf("type = 0x%04X, want string", uint16(d.Type))
	}
	if d.CurrentStr != "NEX-6" {
		t.Errorf("current = %q, want %q — the trailing NUL must be stripped", d.CurrentStr, "NEX-6")
	}
	if d.DefaultStr != "" {
		t.Errorf("default = %q, want empty", d.DefaultStr)
	}
	if got := d.Code.String(); got != "DeviceFriendlyName" {
		t.Errorf("name = %q, want DeviceFriendlyName", got)
	}
}

func TestParseStdPropDescTruncated(t *testing.T) {
	for _, f := range []string{"testdata/nex6-propdesc-5001.bin", "testdata/nex6-propdesc-d402.bin"} {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Skipf("fixture missing: %v", err)
		}
		// Cutting inside the fixed header must always error. Past the values a
		// short read is legal, since the form block is optional.
		for n := range 5 {
			if _, err := ptp.ParseStdPropDesc(raw[:n]); err == nil {
				t.Errorf("%s truncated to %d bytes parsed without error", f, n)
			}
		}
	}
}
