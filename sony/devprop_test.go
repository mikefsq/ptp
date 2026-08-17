package sony

import (
	"encoding/binary"
	"testing"

	"github.com/mikefsq/ptp"
)

// blobBuilder assembles a 0x9209 payload the way a camera would.
type blobBuilder struct {
	entries []byte
	n       uint32
}

func (b *blobBuilder) u8(v uint8)   { b.entries = append(b.entries, v) }
func (b *blobBuilder) u16(v uint16) { b.entries = binary.LittleEndian.AppendUint16(b.entries, v) }
func (b *blobBuilder) u32(v uint32) { b.entries = binary.LittleEndian.AppendUint32(b.entries, v) }

// i16 encodes a signed value without tripping Go's constant-overflow check.
func i16(v int16) uint16 { return uint16(v) }

// head writes the fixed part of an entry: code, type, getset, isEnabled.
func (b *blobBuilder) head(code Prop, t ptp.DataType, getset, enabled uint8) {
	b.n++
	b.u16(uint16(code))
	b.u16(uint16(t))
	b.u8(getset)
	b.u8(enabled)
}

func (b *blobBuilder) bytes() []byte {
	out := make([]byte, 0, 8+len(b.entries))
	out = binary.LittleEndian.AppendUint32(out, b.n)
	out = binary.LittleEndian.AppendUint32(out, 0)
	return append(out, b.entries...)
}

func TestParseScalarUint16Enum(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropFNumber, ptp.TypeUint16, 0x01, 1)
	b.u16(400) // default: F4.0
	b.u16(560) // current: F5.6
	b.u8(uint8(ptp.FormEnum))
	b.u16(3)
	b.u16(280)
	b.u16(400)
	b.u16(560)

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(props) != 1 {
		t.Fatalf("got %d properties, want 1", len(props))
	}
	p := props[0]
	if p.Code != PropFNumber {
		t.Errorf("code = %v, want FNumber", p.Code)
	}
	if p.Default != 400 || p.Current != 560 {
		t.Errorf("default/current = %d/%d, want 400/560", p.Default, p.Current)
	}
	if p.Form != ptp.FormEnum {
		t.Errorf("form = %d, want enum", p.Form)
	}
	want := []uint64{280, 400, 560}
	if len(p.Enum) != len(want) {
		t.Fatalf("enum has %d values, want %d", len(p.Enum), len(want))
	}
	for i := range want {
		if p.Enum[i] != want[i] {
			t.Errorf("enum[%d] = %d, want %d", i, p.Enum[i], want[i])
		}
	}
	if !p.Writable() {
		t.Error("isEnabled=1 should be writable")
	}
}

// Exposure compensation is signed; a negative value must survive the round trip.
func TestParseSignedValueIsSignExtended(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropExposureBiasCompensation, ptp.TypeInt16, 0x01, 1)
	b.u16(0)
	b.u16(i16(-1000)) // -3.3 EV in Sony's 1/1000 units
	b.u8(uint8(ptp.FormRange))
	b.u16(i16(-3000))
	b.u16(3000)
	b.u16(333)

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	p := props[0]
	if int64(p.Current) != -1000 {
		t.Errorf("current = %d, want -1000", int64(p.Current))
	}
	if int64(p.Min) != -3000 {
		t.Errorf("min = %d, want -3000", int64(p.Min))
	}
	if int64(p.Max) != 3000 || int64(p.Step) != 333 {
		t.Errorf("max/step = %d/%d, want 3000/333", int64(p.Max), int64(p.Step))
	}
}

// Bodies from 2024 on may append a second enum that supersedes the first.
func TestSecondaryEnumSupersedesPrimary(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropIsoSensitivity, ptp.TypeUint32, 0x01, 1)
	b.u32(100)
	b.u32(6400)
	b.u8(uint8(ptp.FormEnum))
	b.u16(2) // primary list
	b.u32(100)
	b.u32(200)
	b.u16(3) // secondary list, count < 0x200
	b.u32(100)
	b.u32(3200)
	b.u32(6400)

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := props[0].Enum
	want := []uint64{100, 3200, 6400}
	if len(got) != len(want) {
		t.Fatalf("enum has %d values, want %d (secondary list should win)", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("enum[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

// A following property code (0x5xxx / 0xDxxx) must not be mistaken for a
// secondary enum count.
func TestNextPropertyCodeIsNotTakenAsSecondaryEnum(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropFNumber, ptp.TypeUint16, 0x01, 1)
	b.u16(400)
	b.u16(400)
	b.u8(uint8(ptp.FormEnum))
	b.u16(1)
	b.u16(400)
	// Second entry follows immediately; its code 0xD20D is >= 0x200.
	b.head(PropShutterSpeed, ptp.TypeUint32, 0x01, 1)
	b.u32(0)
	b.u32(0x000A0001)
	b.u8(uint8(ptp.FormNone))

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(props) != 2 {
		t.Fatalf("got %d properties, want 2", len(props))
	}
	if len(props[0].Enum) != 1 || props[0].Enum[0] != 400 {
		t.Errorf("first enum = %v, want [400]", props[0].Enum)
	}
	if props[1].Code != PropShutterSpeed {
		t.Errorf("second code = %v, want ShutterSpeed", props[1].Code)
	}
	if props[1].Current != 0x000A0001 {
		t.Errorf("shutter current = 0x%X, want 0xA0001", props[1].Current)
	}
}

func TestControlPropertyIsWritableRegardlessOfEnabled(t *testing.T) {
	b := &blobBuilder{}
	// S1 (shutter half-press) lives in the control-code space, not the device
	// property table, so it is written as a literal here.
	b.head(Prop(0xD2C1), ptp.TypeUint16, 0x81, 0) // 0x81 = button control, isEnabled 0
	b.u16(0)
	b.u16(0)
	b.u8(uint8(ptp.FormNone))

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	p := props[0]
	if !p.IsControl() {
		t.Error("getset 0x81 should be a control")
	}
	if !p.Writable() {
		t.Error("a control should be writable even with isEnabled = 0")
	}
}

func TestNonWritableWhenDisplayOnly(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropBatteryRemain, ptp.TypeUint16, 0x00, 2) // 2 = display only
	b.u16(0)
	b.u16(87)
	b.u8(uint8(ptp.FormNone))

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if props[0].Writable() {
		t.Error("isEnabled = 2 (display only) should not be writable")
	}
	if props[0].Current != 87 {
		t.Errorf("current = %d, want 87", props[0].Current)
	}
}

func TestParseArrayValue(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropFNumber, ptp.TypeUint16|0x4000, 0x01, 1)
	b.u32(2) // default array
	b.u16(280)
	b.u16(400)
	b.u32(1) // current array
	b.u16(560)
	b.u8(uint8(ptp.FormNone))

	props, err := ParseAllDevicePropData(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	p := props[0]
	if !p.Type.IsArray() {
		t.Error("type should report as an array")
	}
	if len(p.DefaultArr) != 2 || p.DefaultArr[1] != 400 {
		t.Errorf("default array = %v, want [280 400]", p.DefaultArr)
	}
	if len(p.CurrentArr) != 1 || p.CurrentArr[0] != 560 {
		t.Errorf("current array = %v, want [560]", p.CurrentArr)
	}
}

func TestTruncatedBlobIsAnError(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropFNumber, ptp.TypeUint16, 0x01, 1)
	b.u16(400)
	// Entry stops mid-way: no current value, no form flag.
	if _, err := ParseAllDevicePropData(b.bytes()); err == nil {
		t.Fatal("expected an error for a truncated entry")
	}
}

// A corrupt enum count must not drive a huge allocation.
func TestAbsurdEnumCountRejected(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropFNumber, ptp.TypeUint16, 0x01, 1)
	b.u16(400)
	b.u16(400)
	b.u8(uint8(ptp.FormEnum))
	b.u16(0xFFFF) // claims 65535 values, supplies none

	if _, err := ParseAllDevicePropData(b.bytes()); err == nil {
		t.Fatal("expected an error for an enum count the buffer cannot hold")
	}
}

func TestCountMismatchReportsProgress(t *testing.T) {
	b := &blobBuilder{}
	b.head(PropFNumber, ptp.TypeUint16, 0x01, 1)
	b.u16(400)
	b.u16(400)
	b.u8(uint8(ptp.FormNone))
	raw := b.bytes()
	binary.LittleEndian.PutUint32(raw[0:], 5) // claim 5 entries, supply 1

	props, err := ParseAllDevicePropData(raw)
	if err == nil {
		t.Fatal("expected an error when the blob has fewer entries than declared")
	}
	// The successfully parsed prefix should still come back, for diagnosis.
	if len(props) != 1 {
		t.Errorf("got %d parsed properties, want the 1 good entry", len(props))
	}
}
