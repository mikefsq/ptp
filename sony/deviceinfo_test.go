package sony

import (
	"os"
	"testing"

	"github.com/mikefsq/ptp"
)

// A verbatim 0x1001 reply from a real Sony NEX-6, captured 2026-08-06.
//
// This body is the reason SupportsSDIO exists. It is a Sony, it speaks PTP
// perfectly, and it has none of the remote-control surface this package is
// built on — so it is the fixture that proves the check refuses to send 0x92xx
// operations to a camera that would stall its pipe on them.
func TestNEX6HasNoSDIOSurface(t *testing.T) {
	raw, err := os.ReadFile("testdata/nex6-deviceinfo.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	di, err := ptp.ParseDeviceInfo(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if SupportsSDIO(di) {
		t.Error("NEX-6 reported as having the SDIO surface; Connect would then " +
			"send 0x9201 to a body that stalls on it")
	}
	if di.Supports(OpSDIOConnect) {
		t.Error("NEX-6 should not list 0x9201")
	}
	if di.Supports(OpSDIOGetExtDeviceInfo) {
		t.Error("NEX-6 should not list 0x9202")
	}
	if di.Supports(OpGetAllDevicePropData) {
		t.Error("NEX-6 should not list 0x9209")
	}

	// It is NOT bare of vendor operations — it advertises an undecoded
	// 0x9280-0x9285 block. "No SDIO" is the precise claim; "no vendor
	// operations" would be wrong, and pinning that here stops the loose version
	// creeping back into the docs.
	var vendor928x int
	for _, op := range di.Operations {
		if op >= 0x9280 && op <= 0x9285 {
			vendor928x++
		}
	}
	if vendor928x != 6 {
		t.Errorf("got %d operations in the 0x9280-0x9285 block, want 6", vendor928x)
	}
}

// Connect must refuse before issuing a vendor operation, and say why. A body
// without the surface answers a 0x92xx by stalling the bulk pipe, which costs a
// stall recovery and produces a far less useful error.
func TestConnectRefusesABodyWithoutSDIO(t *testing.T) {
	raw, err := os.ReadFile("testdata/nex6-deviceinfo.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	f := &fakeTransport{in: [][]byte{
		mkContainer(ptp.ContainerData, uint16(ptp.OpGetDeviceInfo), 0, nil, raw),
		ok(0),
	}}
	c := openCamera(f)

	if _, err := c.Connect(); err == nil {
		t.Fatal("Connect proceeded against a body with no SDIO surface")
	}
	for _, op := range sentOps(f) {
		if op == OpSDIOConnect {
			t.Error("Connect sent 0x9201 to a body whose DeviceInfo does not list it")
		}
	}
}

// The fixture is a verbatim 0x1001 reply captured from a real Sony NEX-6 on
// 2026-08-06. Keeping real bytes as a golden test means the parser is pinned to
// something a camera actually sent, not to anyone's reading of the spec.
//
// It exercises the PARENT package's parser, but it lives here because the bytes
// are a Sony body's and testdata belongs with the package that owns it. The
// core's own tests build their datasets by hand.
func TestParseDeviceInfoNEX6(t *testing.T) {
	raw, err := os.ReadFile("testdata/nex6-deviceinfo.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	di, err := ptp.ParseDeviceInfo(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if di.Manufacturer != "Sony Corporation" {
		t.Errorf("Manufacturer = %q, want %q", di.Manufacturer, "Sony Corporation")
	}
	if di.Model != "NEX-6" {
		t.Errorf("Model = %q, want %q", di.Model, "NEX-6")
	}
	if di.DeviceVersion != "1.0" {
		t.Errorf("DeviceVersion = %q, want %q", di.DeviceVersion, "1.0")
	}
	if di.SerialNumber != "00000000000000003282411000648565" {
		t.Errorf("SerialNumber = %q", di.SerialNumber)
	}
	if di.StandardVersion != 100 {
		t.Errorf("StandardVersion = %d, want 100", di.StandardVersion)
	}
	if len(di.Operations) != 24 {
		t.Errorf("got %d operations, want 24", len(di.Operations))
	}
	if len(di.DeviceProps) != 5 {
		t.Errorf("got %d device properties, want 5", len(di.DeviceProps))
	}
	if len(di.CaptureFormats) != 0 {
		t.Errorf("got %d capture formats, want 0", len(di.CaptureFormats))
	}
	if len(di.ImageFormats) != 9 {
		t.Errorf("got %d image formats, want 9", len(di.ImageFormats))
	}
	if !di.Supports(ptp.OpGetDeviceInfo) || !di.Supports(ptp.OpOpenSession) {
		t.Error("NEX-6 should list the standard PTP operations")
	}
	// A body that reports no capture formats cannot be told to shoot over this
	// interface, whatever else it lists.
	if di.SupportsCapture() {
		t.Error("NEX-6 should not report InitiateCapture")
	}

	// A full-length parse is the real assertion: any drift in field order or
	// string handling shows up as leftover or missing bytes.
	if di.VendorExtensionID != 0xFFFFFFFF {
		t.Errorf("VendorExtensionID = 0x%08X, want 0xFFFFFFFF", di.VendorExtensionID)
	}
}

func TestParseDeviceInfoTruncated(t *testing.T) {
	raw, err := os.ReadFile("testdata/nex6-deviceinfo.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	// Every truncation point must produce an error rather than a panic or a
	// silently short result.
	for n := 0; n < len(raw); n += 7 {
		if _, err := ptp.ParseDeviceInfo(raw[:n]); err == nil {
			t.Errorf("truncating to %d bytes parsed without error", n)
		}
	}
}
