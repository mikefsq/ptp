package sony

import (
	"testing"

	"github.com/mikefsq/ptp"
)

// Sony's shutter buttons STAY DOWN until released. A host that exits after
// pressing one leaves the body with its shutter held and nobody to lift it;
// closing the PTP session does not undo that.
func TestCloseReleasesHeldButtons(t *testing.T) {
	f := &fakeTransport{in: okRun(8)}
	c := openCamera(f)

	if err := c.FullPress(); err != nil {
		t.Fatalf("FullPress: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	var released bool
	for _, w := range buttonWrites(f) {
		if w.Ctrl == CtrlS2Button && uint64(w.Val) == ButtonUp {
			released = true
		}
	}
	if !released {
		t.Error("Close left the shutter button held: the body would sit with S2 " +
			"pressed after the host exits")
	}
}

// A session that pressed nothing must send no vendor control ops on close.
// Bodies without the SDIO surface — a NEX-6, for instance — stall the pipe on
// an unsupported vendor op, so speculative releases are not harmless.
func TestCloseSendsNoControlOpsWhenNothingPressed(t *testing.T) {
	f := &fakeTransport{in: okRun(4)}
	c := openCamera(f)

	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	for _, op := range sentOps(f) {
		if op == OpSetControlDeviceB {
			t.Fatal("Close sent a vendor control op without having pressed anything")
		}
	}
}

// Releasing a button clears it, so a clean shot leaves nothing for Close.
func TestReleasedButtonIsNotReReleased(t *testing.T) {
	f := &fakeTransport{in: okRun(8)}
	c := openCamera(f)

	if err := c.Shoot(); err != nil { // press then release
		t.Fatalf("Shoot: %v", err)
	}
	if held := c.Held(); len(held) != 0 {
		t.Fatalf("after a clean shot the session still thinks %v are held", held)
	}
	before := len(f.out)
	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	// Only the CloseSession command should follow.
	if got := len(f.out) - before; got != 1 {
		t.Errorf("Close sent %d transfers, want 1 (CloseSession only)", got)
	}
}

// A press whose response never arrives must still count as held. The camera may
// have acted on it, and a shutter left down is far worse than a spurious
// release on the way out.
func TestFailedPressCountsAsHeld(t *testing.T) {
	// No queued replies beyond the one openCamera adds for OpenSession, so the
	// press itself goes unanswered.
	f := &fakeTransport{inErr: ptp.ErrNotResponding}
	c := openCamera(f)

	if err := c.FullPress(); err == nil {
		t.Fatal("FullPress reported success though the camera never answered")
	}
	held := c.Held()
	if len(held) != 1 || held[0] != CtrlS2Button {
		t.Errorf("held = %v, want [S2Button]: a press that errored may still have "+
			"landed, so it must be released on close", held)
	}
}

// S2 must be released before S1: letting S1 up while S2 is still down is not a
// state the body expects.
func TestReleaseOrderS2BeforeS1(t *testing.T) {
	f := &fakeTransport{in: okRun(10)}
	c := openCamera(f)

	if err := c.HalfPress(); err != nil {
		t.Fatalf("HalfPress: %v", err)
	}
	if err := c.FullPress(); err != nil {
		t.Fatalf("FullPress: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	var order []ControlCode
	for _, w := range buttonWrites(f) {
		if uint64(w.Val) == ButtonUp {
			order = append(order, w.Ctrl)
		}
	}
	if len(order) < 2 {
		t.Fatalf("expected both buttons released, got %v", order)
	}
	if order[0] != CtrlS2Button || order[1] != CtrlS1Button {
		t.Errorf("release order = %v, want S2 then S1", order)
	}
}
