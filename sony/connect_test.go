package sony

import (
	"encoding/binary"
	"testing"

	"github.com/mikefsq/ptp"
)

// mkExtDeviceInfo builds a 0x9202 payload: a uint16 version, then two arrays,
// each a uint32 count followed by that many uint16 codes.
func mkExtDeviceInfo(version uint16, props, controls []uint16) []byte {
	b := binary.LittleEndian.AppendUint16(nil, version)
	for _, list := range [][]uint16{props, controls} {
		b = binary.LittleEndian.AppendUint32(b, uint32(len(list)))
		for _, v := range list {
			b = binary.LittleEndian.AppendUint16(b, v)
		}
	}
	return b
}

func TestParseExtDeviceInfo(t *testing.T) {
	payload := mkExtDeviceInfo(Protocol300,
		[]uint16{0x5007, 0xD20D, 0xD21E},
		[]uint16{0xD2C1, 0xD2C2})

	info, err := parseExtDeviceInfo(payload)
	if err != nil {
		t.Fatalf("parseExtDeviceInfo: %v", err)
	}
	if info.Version != Protocol300 {
		t.Errorf("version = 0x%04X, want 0x%04X", info.Version, Protocol300)
	}
	if info.ModeVer() != 3 {
		t.Errorf("ModeVer = %d, want 3", info.ModeVer())
	}
	if len(info.Props) != 3 || info.Props[1] != PropShutterSpeed {
		t.Errorf("props = %v, want 3 entries with ShutterSpeed second", info.Props)
	}
	if len(info.Controls) != 2 {
		t.Errorf("controls = %v, want 2 entries", info.Controls)
	}
}

// A body speaking the older protocol must be reported as generation 2: it
// selects a different parser in the SDK, and treating it as 3 misreads the
// property blob rather than failing cleanly.
func TestExtDeviceInfoOldProtocolIsModeVer2(t *testing.T) {
	info, err := parseExtDeviceInfo(mkExtDeviceInfo(Protocol200, nil, nil))
	if err != nil {
		t.Fatalf("parseExtDeviceInfo: %v", err)
	}
	if info.ModeVer() != 2 {
		t.Errorf("ModeVer = %d, want 2", info.ModeVer())
	}
}

// A truncated payload must fail rather than return a half-built list. A count
// the buffer cannot hold is the shape a corrupt or misparsed reply takes.
func TestParseExtDeviceInfoRejectsShortPayload(t *testing.T) {
	b := binary.LittleEndian.AppendUint16(nil, Protocol300)
	b = binary.LittleEndian.AppendUint32(b, 1000) // claims 1000 codes
	b = binary.LittleEndian.AppendUint16(b, 0x5007)
	if _, err := parseExtDeviceInfo(b); err == nil {
		t.Fatal("accepted a payload claiming 1000 codes with room for one")
	}
}

// Connect must record what the body reported, so Ready and any later capability
// check see it. A handshake that succeeded but left Ext nil reads as "never
// connected".
func TestConnectRecordsExtDeviceInfo(t *testing.T) {
	payload := mkExtDeviceInfo(Protocol300, []uint16{0x5007}, []uint16{0xD2C1})
	f := &fakeTransport{in: [][]byte{
		// Connect asks GetDeviceInfo first, to avoid stalling a body with no
		// vendor surface. A bare OK with no data phase fails to parse, and an
		// unreadable DeviceInfo means "carry on" rather than "give up".
		ok(0),
		ok(0), // phase 1
		ok(0), // phase 2
		mkContainer(ptp.ContainerData, uint16(OpSDIOGetExtDeviceInfo), 0, nil, payload),
		ok(0), // response to 0x9202
		ok(0), // phase 3
	}}
	c := openCamera(f)
	info, err := c.Connect()
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if c.Ext == nil {
		t.Fatal("Connect succeeded but left Camera.Ext nil, so Ready would report no handshake")
	}
	if c.Ext != info {
		t.Error("Camera.Ext is not the ExtDeviceInfo Connect returned")
	}
	if err := c.Ready(); err != nil && err.Error() == "sony: vendor handshake has not run; call Connect (the body must be in PC Remote USB mode)" {
		t.Error("Ready still reports no handshake after a successful Connect")
	}
}
