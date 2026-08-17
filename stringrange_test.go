package ptp

import "testing"

// An X-T5 describes five properties as a STRING-typed RANGE — bounds like
// min "-6,-4,1" max "6,4,6" step "1,1,1", a structured triple carried as text.
// Reading those bounds as scalars fails, and failing the descriptor loses the
// property from the sweep entirely.
func TestStringTypedRangeDescriptor(t *testing.T) {
	// 0xD347, string type, writable, default/current strings, form=range,
	// then three string bounds.
	raw := []byte{
		0x47, 0xD3, // code
		0xFF, 0xFF, // type: string
		0x01, // get/set: writable
	}
	str := func(s string) []byte {
		out := []byte{byte(len(s) + 1)}
		for _, r := range s {
			out = append(out, byte(r), 0)
		}
		return append(out, 0, 0)
	}
	raw = append(raw, str("a,b,c")...)   // default
	raw = append(raw, str("0,0,4")...)   // current
	raw = append(raw, 0x01)              // form: range
	raw = append(raw, str("-6,-4,1")...) // min
	raw = append(raw, str("6,4,6")...)   // max
	raw = append(raw, str("1,1,1")...)   // step

	d, err := ParseStdPropDesc(raw)
	if err != nil {
		t.Fatalf("a string-typed range must parse: %v", err)
	}
	if d.Code != 0xD347 || d.Type != TypeString {
		t.Fatalf("code/type = 0x%04X/0x%04X", uint16(d.Code), uint16(d.Type))
	}
	if d.CurrentStr != "0,0,4" {
		t.Errorf("CurrentStr = %q, want %q", d.CurrentStr, "0,0,4")
	}
	if d.MinStr != "-6,-4,1" || d.MaxStr != "6,4,6" || d.StepStr != "1,1,1" {
		t.Errorf("bounds = %q..%q step %q, want -6,-4,1..6,4,6 step 1,1,1",
			d.MinStr, d.MaxStr, d.StepStr)
	}
}
