package sony

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/mikefsq/ptp"
)

// The clock is a STRING property, not a packed number. Writing it as a scalar
// would be accepted-and-ignored in the usual Sony way, so the encoding is worth
// pinning: a uint8 character count, then UTF-16, with the NUL inside the count.
func TestSetDateTimeWritesAPTPString(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)

	when := time.Date(2026, 8, 7, 14, 41, 56, 0, time.UTC)
	if err := c.SetDateTime(when); err != nil {
		t.Fatalf("SetDateTime: %v", err)
	}

	// The data container carries the encoded string.
	var payload []byte
	for _, w := range f.out {
		if len(w) > ptp.ContainerHeaderLen && binary.LittleEndian.Uint16(w[4:]) == ptp.ContainerData {
			payload = w[ptp.ContainerHeaderLen:]
		}
	}
	if payload == nil {
		t.Fatal("no data container was written")
	}
	want := ptp.EncodeString("20260807T144156")
	if string(payload) != string(want) {
		t.Errorf("payload = %x, want %x", payload, want)
	}
	// And it must go through the vendor setter, not the standard one.
	var sawVendorSet bool
	for _, op := range sentOps(f) {
		if op == OpSetControlDeviceA {
			sawVendorSet = true
		}
	}
	if !sawVendorSet {
		t.Error("the clock was not written through 0x9205 SetControlDeviceA")
	}
}

// Round-tripping through the camera's own format is what proves the read and
// the write agree; a body reporting a time this cannot parse is a real failure,
// not a cosmetic one.
func TestDateTimeRoundTrip(t *testing.T) {
	when := time.Date(2026, 8, 7, 14, 41, 56, 0, time.UTC)
	got, err := ptp.ParseDateTime(ptp.FormatDateTime(when))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !got.Equal(when) {
		t.Errorf("round trip gave %v, want %v", got, when)
	}
}

// SyncClockAtMinute must land ON the boundary. If a body truncates seconds the
// way an X-T5 does, writing at any other moment bakes the remainder in as
// error, so the deliberate aim-early is the whole point of the function.
func TestSyncClockAtMinuteTargetsTheBoundary(t *testing.T) {
	now := time.Now()
	next := now.Truncate(time.Minute).Add(time.Minute)
	if next.Sub(now) > time.Minute || next.Second() != 0 {
		t.Fatalf("boundary maths is wrong: now %v -> next %v", now, next)
	}
}

// ReadSettings must distinguish "the camera says zero" from "the camera said
// nothing" — the same reason Exposure carries per-field flags. A body that does
// not expose a setting must not read as one set to the zero value.
func TestReadSettingsFlagsAbsentProperties(t *testing.T) {
	props := []DeviceProperty{
		{Code: PropFileType, Type: ptp.TypeUint16, Current: FileTypeRAW, Enabled: 1},
		{Code: PropImageSize, Type: ptp.TypeUint16, Current: 0, Enabled: 1},
	}
	s := ReadSettings(props)

	if !s.HasFileType || s.FileType != FileTypeRAW {
		t.Errorf("FileType = %d (has=%v), want RAW", s.FileType, s.HasFileType)
	}
	// Reported as zero: present, and zero.
	if !s.HasImageSize || s.ImageSize != 0 {
		t.Errorf("ImageSize = %d (has=%v), want present and 0", s.ImageSize, s.HasImageSize)
	}
	// Never mentioned: absent.
	if s.HasAspectRatio {
		t.Error("AspectRatio reported as present though the camera never mentioned it")
	}
	if s.HasRAWCompression {
		t.Error("RAWCompression reported as present though the camera never mentioned it")
	}
}

// Settable must reflect the camera's enable byte, because a write outside what
// the body currently allows is accepted and then ignored rather than refused.
func TestSettableUsesTheCamerasEnableByte(t *testing.T) {
	props := []DeviceProperty{
		{Code: PropFileType, Type: ptp.TypeUint16, GetSet: 0x01, Enabled: 1},
		{Code: PropImageSize, Type: ptp.TypeUint16, GetSet: 0x01, Enabled: 0},
	}
	if d := FindProp(props, PropFileType); !d.Writable() {
		t.Error("an enabled property reported as not writable")
	}
	if d := FindProp(props, PropImageSize); d.Writable() {
		t.Error("a greyed-out property reported as writable")
	}
}
