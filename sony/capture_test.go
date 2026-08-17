package sony

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/mikefsq/ptp"
)

// A Sony shot is S2 down then S2 up. Getting the button values or the ordering
// wrong is the difference between a photo and a body stuck in continuous drive.
func TestCaptureSendsPressThenRelease(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)
	if err := c.Shoot(); err != nil {
		t.Fatalf("Shoot: %v", err)
	}
	// Two transactions, each a command container plus a data container.
	if len(f.out) != 4 {
		t.Fatalf("wrote %d transfers, want 4 (command+data twice)", len(f.out))
	}
	for i, want := range []struct {
		code uint16
		val  uint16
	}{
		{uint16(OpSetControlDeviceB), uint16(ButtonDown)},
		{uint16(OpSetControlDeviceB), uint16(ButtonUp)},
	} {
		cmd := f.out[i*2]
		data := f.out[i*2+1]
		if got := binary.LittleEndian.Uint16(cmd[6:]); got != want.code {
			t.Errorf("transaction %d opcode = 0x%04X, want 0x%04X", i, got, want.code)
		}
		if got := binary.LittleEndian.Uint32(cmd[12:]); got != uint32(CtrlS2Button) {
			t.Errorf("transaction %d control = 0x%08X, want S2Button 0x%08X", i, got, uint32(CtrlS2Button))
		}
		if got := binary.LittleEndian.Uint16(data[ptp.ContainerHeaderLen:]); got != want.val {
			t.Errorf("transaction %d button value = %d, want %d", i, got, want.val)
		}
	}
}

func TestHalfPressUsesS1(t *testing.T) {
	f := &fakeTransport{in: okRun(2)}
	c := openCamera(f)
	if err := c.HalfPress(); err != nil {
		t.Fatalf("HalfPress: %v", err)
	}
	if got := binary.LittleEndian.Uint32(f.out[0][12:]); got != uint32(CtrlS1Button) {
		t.Errorf("control = 0x%08X, want S1Button 0x%08X", got, uint32(CtrlS1Button))
	}
}

// A press that reports an error must still be followed by a release, or the
// shutter stays held.
func TestBulbCaptureHoldsForDuration(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)
	start := time.Now()
	if err := c.BulbCapture(120 * time.Millisecond); err != nil {
		t.Fatalf("BulbCapture: %v", err)
	}
	if el := time.Since(start); el < 100*time.Millisecond {
		t.Errorf("held for %v, want at least ~120ms", el)
	}
	if len(f.out) != 4 {
		t.Fatalf("wrote %d transfers, want 4", len(f.out))
	}
	if got := binary.LittleEndian.Uint16(f.out[1][ptp.ContainerHeaderLen:]); got != uint16(ButtonDown) {
		t.Errorf("first value = %d, want Down", got)
	}
	if got := binary.LittleEndian.Uint16(f.out[3][ptp.ContainerHeaderLen:]); got != uint16(ButtonUp) {
		t.Errorf("second value = %d, want Up", got)
	}
}

func TestControlCodeIs32Bit(t *testing.T) {
	// Release sits above 0xFFFF, which is why ControlCode cannot be a Prop.
	if uint32(CtrlRelease) != 0x00010001 {
		t.Errorf("CtrlRelease = 0x%08X, want 0x00010001", uint32(CtrlRelease))
	}
	if uint32(CtrlS2Button) != 0xD2C2 {
		t.Errorf("CtrlS2Button = 0x%08X, want 0xD2C2", uint32(CtrlS2Button))
	}
	if got := CtrlS1Button.String(); got != "S1Button" {
		t.Errorf("name = %q, want S1Button", got)
	}
}

func TestWaitForCaptureReturnsHandle(t *testing.T) {
	f := &fakeTransport{evt: [][]byte{
		mkContainer(ptp.ContainerEvent, uint16(ptp.EventObjectAdded), 1, []uint32{0x000C0011}, nil),
	}}
	c := openCamera(f)
	h, err := c.WaitForCapture(2 * time.Second)
	if err != nil {
		t.Fatalf("WaitForCapture: %v", err)
	}
	if h != 0x000C0011 {
		t.Errorf("handle = 0x%08X, want 0x000C0011", h)
	}
}

func TestWaitForCaptureTimesOut(t *testing.T) {
	c := openCamera(&fakeTransport{})
	if _, err := c.WaitForCapture(200 * time.Millisecond); err != ptp.ErrTimeout {
		t.Fatalf("err = %v, want ptp.ErrTimeout", err)
	}
}

func TestParseLicenseInfoList(t *testing.T) {
	var b []byte
	b = append(b, 2) // two licenses
	b = binary.LittleEndian.AppendUint32(b, InfiniteLicense)
	b = append(b, 4)
	b = append(b, []byte("PERM")...)
	b = binary.LittleEndian.AppendUint32(b, 720)
	b = append(b, 5)
	b = append(b, []byte("TRIAL")...)

	f := &fakeTransport{in: [][]byte{
		mkContainer(ptp.ContainerData, uint16(OpGetLicenseInfoList), 1, nil, b),
		mkContainer(ptp.ContainerResponse, uint16(ptp.RespOK), 1, nil, nil),
	}}
	c := openCamera(f)
	got, err := c.GetLicenseInfoList()
	if err != nil {
		t.Fatalf("GetLicenseInfoList: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d licenses, want 2", len(got))
	}
	if got[0].ID != "PERM" || got[0].RemainingHours != InfiniteLicense {
		t.Errorf("license 0 = %+v", got[0])
	}
	if got[1].ID != "TRIAL" || got[1].RemainingHours != 720 {
		t.Errorf("license 1 = %+v", got[1])
	}
}

// The two buttons must be independently controllable: an eclipse client holds
// S1 across a sequence and fires S2 repeatedly.
func TestButtonsAreIndependent(t *testing.T) {
	for _, tc := range []struct {
		name string
		call func(*Camera) error
		ctrl ControlCode
		val  uint16
	}{
		{"HalfPress", (*Camera).HalfPress, CtrlS1Button, uint16(ButtonDown)},
		{"HalfRelease", (*Camera).HalfRelease, CtrlS1Button, uint16(ButtonUp)},
		{"FullPress", (*Camera).FullPress, CtrlS2Button, uint16(ButtonDown)},
		{"FullRelease", (*Camera).FullRelease, CtrlS2Button, uint16(ButtonUp)},
	} {
		f := &fakeTransport{in: okRun(2)}
		c := openCamera(f)
		if err := tc.call(c); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if len(f.out) != 2 {
			t.Fatalf("%s wrote %d transfers, want 2", tc.name, len(f.out))
		}
		if got := binary.LittleEndian.Uint32(f.out[0][12:]); got != uint32(tc.ctrl) {
			t.Errorf("%s control = 0x%08X, want 0x%08X", tc.name, got, uint32(tc.ctrl))
		}
		if got := binary.LittleEndian.Uint16(f.out[1][ptp.ContainerHeaderLen:]); got != tc.val {
			t.Errorf("%s value = %d, want %d", tc.name, got, tc.val)
		}
	}
}

// ReleaseAll must let S2 up before S1: releasing the half-press while the full
// press is still held is not a state the body expects.
func TestReleaseAllOrdersS2BeforeS1(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)
	if err := c.ReleaseAll(); err != nil {
		t.Fatalf("ReleaseAll: %v", err)
	}
	if len(f.out) != 4 {
		t.Fatalf("wrote %d transfers, want 4", len(f.out))
	}
	if got := binary.LittleEndian.Uint32(f.out[0][12:]); got != uint32(CtrlS2Button) {
		t.Errorf("first control = 0x%08X, want S2 released first", got)
	}
	if got := binary.LittleEndian.Uint32(f.out[2][12:]); got != uint32(CtrlS1Button) {
		t.Errorf("second control = 0x%08X, want S1 second", got)
	}
}

// A press whose response never arrives must still be followed by a release: a
// body left with S2 held keeps firing in continuous drive.
func TestCaptureReleasesEvenWhenPressFails(t *testing.T) {
	// Only enough responses for the press to fail and the release to be tried.
	f := &fakeTransport{in: nil}
	c := openCamera(f)
	if err := c.Shoot(); err == nil {
		t.Fatal("expected an error when the press fails")
	}
	// Two command containers: the failed press and the release attempt.
	presses := 0
	for _, w := range f.out {
		if binary.LittleEndian.Uint16(w[4:]) == ptp.ContainerCommand &&
			binary.LittleEndian.Uint16(w[6:]) == uint16(OpSetControlDeviceB) {
			presses++
		}
	}
	if presses < 2 {
		t.Errorf("saw %d control commands, want the release to be attempted after a failed press", presses)
	}
}

func TestShootWithAFHoldsS1AroundTheShot(t *testing.T) {
	f := &fakeTransport{in: okRun(8)}
	c := openCamera(f)
	if err := c.ShootWithAF(10 * time.Millisecond); err != nil {
		t.Fatalf("ShootWithAF: %v", err)
	}
	// S1 down, S2 down, S2 up, S1 up — four transactions, eight transfers.
	if len(f.out) != 8 {
		t.Fatalf("wrote %d transfers, want 8", len(f.out))
	}
	want := []ControlCode{CtrlS1Button, CtrlS2Button, CtrlS2Button, CtrlS1Button}
	for i, w := range want {
		if got := binary.LittleEndian.Uint32(f.out[i*2][12:]); got != uint32(w) {
			t.Errorf("transaction %d control = 0x%08X, want 0x%08X", i, got, uint32(w))
		}
	}
}

// A burst is two transactions for the whole sequence, not two per frame. That
// is the point: frame timing becomes the camera's, not the USB link'c.
func TestBurstIsTwoTransactions(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)
	start := time.Now()
	if err := c.Burst(80 * time.Millisecond); err != nil {
		t.Fatalf("Burst: %v", err)
	}
	if el := time.Since(start); el < 70*time.Millisecond {
		t.Errorf("held %v, want ~80ms", el)
	}
	if len(f.out) != 4 {
		t.Fatalf("wrote %d transfers, want 4 (two transactions for the whole burst)", len(f.out))
	}
	if got := binary.LittleEndian.Uint16(f.out[1][ptp.ContainerHeaderLen:]); got != uint16(ButtonDown) {
		t.Errorf("first value = %d, want Down", got)
	}
	if got := binary.LittleEndian.Uint16(f.out[3][ptp.ContainerHeaderLen:]); got != uint16(ButtonUp) {
		t.Errorf("last value = %d, want Up", got)
	}
}

func TestBurstUntilStopsEarly(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)
	stop := make(chan struct{})
	go func() { time.Sleep(30 * time.Millisecond); close(stop) }()

	start := time.Now()
	if err := c.BurstUntil(stop, 5*time.Second); err != nil {
		t.Fatalf("BurstUntil: %v", err)
	}
	if el := time.Since(start); el > time.Second {
		t.Errorf("took %v; it should stop on the channel, not the timeout", el)
	}
	// The shutter must be released whichever way it ended.
	if got := binary.LittleEndian.Uint16(f.out[3][ptp.ContainerHeaderLen:]); got != uint16(ButtonUp) {
		t.Errorf("final value = %d, want Up", got)
	}
}
