package sony

import (
	"encoding/binary"
	"testing"

	"github.com/mikefsq/ptp"
)

// The credit fields are strings and must go out as PTP strings through the
// vendor setter, not as scalars.
func TestSetPhotographerWritesAString(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)

	if err := c.SetPhotographer("M. Furlotti"); err != nil {
		t.Fatalf("SetPhotographer: %v", err)
	}
	var payload []byte
	for _, w := range f.out {
		if len(w) > ptp.ContainerHeaderLen && binary.LittleEndian.Uint16(w[4:]) == ptp.ContainerData {
			payload = w[ptp.ContainerHeaderLen:]
		}
	}
	if string(payload) != string(ptp.EncodeString("M. Furlotti")) {
		t.Errorf("payload = %x", payload)
	}
}

// WriteCopyrightInfo is the switch that makes the credit strings reach the
// files. Its OFF value is 1, not 0 — sending 0 would be a value the property
// does not take, and setting the text without the switch changes nothing in the
// output, which is the obvious way to lose an hour.
func TestCopyrightInfoOffIsOneNotZero(t *testing.T) {
	if CopyrightInfoOff != 0x01 || CopyrightInfoOn != 0x02 {
		t.Fatalf("off=%d on=%d, want 1 and 2", CopyrightInfoOff, CopyrightInfoOn)
	}
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)
	if err := c.SetCopyrightInfo(false); err != nil {
		t.Fatalf("SetCopyrightInfo: %v", err)
	}
	for _, w := range f.out {
		if len(w) == ptp.ContainerHeaderLen+1 && binary.LittleEndian.Uint16(w[4:]) == ptp.ContainerData {
			if got := w[ptp.ContainerHeaderLen]; got != byte(CopyrightInfoOff) {
				t.Errorf("wrote %d for off, want %d", got, CopyrightInfoOff)
			}
			return
		}
	}
	t.Error("no single-byte data container was written; the property is UInt8")
}

// ReadMetadata pulls every string field from one snapshot. The individual
// accessors each cost a full property fetch, so a caller reading several should
// not be paying for four.
func TestReadMetadataFromOneSnapshot(t *testing.T) {
	props := []DeviceProperty{
		{Code: PropSetPhotographer, Type: ptp.TypeString, CurrentStr: "M. Furlotti"},
		{Code: PropSetCopyright, Type: ptp.TypeString, CurrentStr: "(c) 2026"},
		{Code: PropFileSettingsTitleNameSettings, Type: ptp.TypeString, CurrentStr: "ECLIPSE"},
		{Code: PropWriteCopyrightInfo, Type: ptp.TypeUint8, Current: CopyrightInfoOn},
	}
	m := ReadMetadata(props)

	if m.Photographer != "M. Furlotti" || m.Copyright != "(c) 2026" || m.TitleName != "ECLIPSE" {
		t.Errorf("metadata = %+v", m)
	}
	if !m.HasEmbedding || !m.Embedding {
		t.Error("embedding should be reported as on")
	}
}

// A body that never mentions the switch must not read as "embedding off" — the
// same absent-versus-zero distinction the rest of the package makes.
func TestMetadataDistinguishesAbsentFromOff(t *testing.T) {
	m := ReadMetadata(nil)
	if m.HasEmbedding {
		t.Error("embedding reported as known though the camera never mentioned it")
	}
	off := ReadMetadata([]DeviceProperty{
		{Code: PropWriteCopyrightInfo, Type: ptp.TypeUint8, Current: CopyrightInfoOff},
	})
	if !off.HasEmbedding || off.Embedding {
		t.Errorf("explicit off not reported correctly: %+v", off)
	}
}

// Sony has no Artist and no Comment field. Accessors named for them would be
// inventing properties the body does not have, so the string surface is exactly
// the six the SDK marks CrDataType_STR.
func TestOnlyDocumentedStringPropertiesAreExposed(t *testing.T) {
	for _, p := range []Prop{
		PropSetPhotographer, PropSetCopyright,
		PropFileSettingsTitleNameSettings, PropRecordingSettingFileName,
	} {
		if _, ok := PropTable[p]; !ok {
			t.Errorf("0x%04X is not in the SDK property table", uint16(p))
		}
	}
}
