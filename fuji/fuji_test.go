package fuji

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mikefsq/ptp"
)

// A real X-T5's GetDeviceInfo response, captured over USB. It is the only
// hardware artefact the parser is checked against, so it guards the decode
// against a body that is not on the desk.
func TestParseRealDeviceInfo(t *testing.T) {
	raw, err := os.ReadFile("testdata/xt5-deviceinfo.bin")
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	di, err := ptp.ParseDeviceInfo(raw)
	if err != nil {
		t.Fatalf("parsing a known-good X-T5 DeviceInfo: %v", err)
	}
	if di.Model == "" {
		t.Error("no model parsed")
	}
	// The body advertises 263 device properties in tethered shooting mode; a
	// handful means it is in a file-transfer USB mode instead.
	if len(di.DeviceProps) < 200 {
		t.Errorf("got %d device properties, want the tethered-mode count (~263)",
			len(di.DeviceProps))
	}
	if !di.SupportsCapture() {
		t.Error("SupportsCapture false for a body that captures")
	}
}

// PC Priority locks out the camera's dials, buttons and shutter release, and
// closing the PTP session does not undo it. A host that exits without handing
// control back leaves the body inert in its owner's hands, so Close must write
// camera priority — and must keep trying, because the camera answers 0xA002
// ("not yet") while it is still busy.
func TestCloseHandsControlBackToTheCamera(t *testing.T) {
	f := &fakeTransport{}
	for i := uint32(2); i < 18; i++ {
		f.in = append(f.in, ok(i))
	}
	c := openCamera(f)

	if err := c.Session.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	var gaveBack bool
	for _, w := range f.out {
		if len(w) == ptp.ContainerHeaderLen+2 &&
			binary.LittleEndian.Uint16(w[4:]) == ptp.ContainerData &&
			ptp.OpCode(binary.LittleEndian.Uint16(w[6:])) == ptp.OpSetDevicePropValue &&
			uint64(binary.LittleEndian.Uint16(w[ptp.ContainerHeaderLen:])) == PriorityModeCamera {
			gaveBack = true
		}
	}
	if !gaveBack {
		t.Error("Close did not return the camera to Camera Priority: the body's " +
			"dials, buttons and shutter would stay dead after the host exits")
	}
}

// 0xA002 means "refused right now", not "bad value". The camera returns it
// while its buffer still holds frames, so a single write is not enough — the
// SDK's own sample loops here.
func TestHandBackRetriesWhileTheCameraIsBusy(t *testing.T) {
	f := &fakeTransport{}
	// The S1-off write and its InitiateCapture, then two refusals before the
	// camera finally accepts control back.
	f.in = append(f.in, ok(2), ok(3), refused(4), refused(5), ok(6), ok(7))
	c := openCamera(f)

	old := HandBackTimeout
	HandBackTimeout = 3 * time.Second
	defer func() { HandBackTimeout = old }()

	done := make(chan struct{})
	go func() { c.Session.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("handback did not give up within its bound")
	}

	// Three attempts at the priority write: two refused, one accepted.
	writes := 0
	for _, w := range f.out {
		if len(w) == ptp.ContainerHeaderLen+2 &&
			binary.LittleEndian.Uint16(w[4:]) == ptp.ContainerData &&
			uint64(binary.LittleEndian.Uint16(w[ptp.ContainerHeaderLen:])) == PriorityModeCamera {
			writes++
		}
	}
	if writes < 2 {
		t.Errorf("priority write attempted %d times, want retries after 0xA002", writes)
	}
}

// A vendor error must say what it means. 0xA002 printing as a bare number is
// exactly what made "the camera refuses, retry later" look like "bad value",
// and cost most of a bring-up session.
func TestVendorResponseCodeIsNamed(t *testing.T) {
	f := &fakeTransport{}
	f.in = append(f.in, refused(2), ok(3), ok(4), ok(5))
	c := openCamera(f)

	_, _, err := c.Do(ptp.OpGetDevicePropValue, []uint32{0xD207}, nil, time.Second)
	if err == nil {
		t.Fatal("the refusal should surface as an error")
	}
	var pe *ptp.Error
	if !errors.As(err, &pe) {
		t.Fatalf("want *ptp.Error, got %T", err)
	}
	if pe.Name != "RefusedRightNow" {
		t.Errorf("Name = %q, want RefusedRightNow — an unnamed vendor code is a "+
			"number the reader has to go and look up", pe.Name)
	}
	if !strings.Contains(err.Error(), "RefusedRightNow") {
		t.Errorf("the message should name the code, got: %v", err)
	}
}

// A camera holding undownloaded frames refuses to hand control to the host,
// and the raw refusal says nothing about why: the property's own descriptor
// says the value is allowed. Reporting the frame count turns an opaque vendor
// code into an instruction.
func TestTakePriorityExplainsPendingFrames(t *testing.T) {
	f := &fakeTransport{}
	f.in = append(f.in,
		// PriorityMode read says "camera", so the write is attempted.
		mkContainer(ptp.ContainerData, uint16(ptp.OpGetDevicePropValue), 2, nil, []byte{1, 0}),
		ok(2),
		refusedCode(3, 0xA001), // the write is refused
		// GetObjectHandles reports two frames still waiting.
		mkContainer(ptp.ContainerData, uint16(ptp.OpGetObjectHandles), 4,
			nil, []byte{2, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0}),
		ok(4),
		ok(5), ok(6), ok(7),
	)
	c := openCamera(f)

	err := c.TakePriority()
	if err == nil {
		t.Fatal("expected the refusal to surface")
	}
	if !strings.Contains(err.Error(), "undownloaded frame") {
		t.Errorf("error should name the cause and the fix, got: %v", err)
	}
	if !strings.Contains(err.Error(), "2") {
		t.Errorf("error should report how many frames are waiting, got: %v", err)
	}
}
