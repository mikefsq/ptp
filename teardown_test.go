package ptp

import (
	"testing"
	"time"
)

// Closing a PTP session does not undo whatever the host did to take control,
// and what that means is vendor-specific: a Fujifilm body left in PC Priority
// has its dials and shutter dead, a Sony body can be left with its shutter
// physically held. The core must therefore give the vendor a chance to hand the
// camera back — before the session goes away, while transactions still work.
func TestCloseRunsTeardownBeforeClosingSession(t *testing.T) {
	f := &fakeTransport{}
	for i := 0; i < 4; i++ {
		f.in = append(f.in, mkContainer(ContainerResponse, uint16(RespOK), uint32(i+1), nil, nil))
	}
	s := NewSession(f)
	s.open = true

	var order []OpCode
	s.Teardown = func(tx Tx) {
		tx(OpSetDevicePropValue, []uint32{0xD207}, []byte{1, 0}, time.Second)
		order = append(order, OpSetDevicePropValue)
	}
	s.Trace = func(ev TraceEvent) { order = append(order, ev.Op) }

	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Expect the teardown transaction, its trace, then CloseSession.
	var sawTeardown, sawClose bool
	for _, op := range order {
		switch op {
		case OpSetDevicePropValue:
			sawTeardown = true
		case OpCloseSession:
			if !sawTeardown {
				t.Fatal("CloseSession ran before the vendor teardown: the camera " +
					"would be left in whatever state the host put it in")
			}
			sawClose = true
		}
	}
	if !sawTeardown || !sawClose {
		t.Fatalf("teardown=%v close=%v, want both", sawTeardown, sawClose)
	}
}

// A teardown hook must be able to run transactions without deadlocking. It is
// handed a Tx precisely because the session lock is already held; this test
// fails by hanging if that ever regresses to passing the Session.
func TestTeardownHookCanTransactWithoutDeadlock(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
		mkContainer(ContainerResponse, uint16(RespOK), 2, nil, nil),
	}}
	s := NewSession(f)
	s.open = true

	done := make(chan struct{})
	s.Teardown = func(tx Tx) {
		tx(OpInitiateCapture, []uint32{0, 0}, nil, time.Second)
		close(done)
	}
	go func() { s.Close() }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("teardown hook deadlocked against the session lock")
	}
}

// No hook, no extra traffic: a core that guessed at a teardown sequence would
// stall the pipe on a body without that vendor surface.
func TestCloseWithoutTeardownSendsOnlyCloseSession(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
	}}
	s := NewSession(f)
	s.open = true

	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if len(f.out) != 1 {
		t.Errorf("Close sent %d transfers, want 1 (CloseSession only)", len(f.out))
	}
}
