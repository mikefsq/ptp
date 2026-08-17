package ptp

import (
	"errors"
	"testing"
	"time"
)

// resettableTransport is a fake whose device is wedged until Reset is called:
// every transfer fails with ErrNotResponding, exactly as a camera left waiting
// for the rest of an abandoned data phase does.
type resettableTransport struct {
	fakeTransport
	wedged  bool
	resets  int
	failing bool // Reset itself fails
}

func (r *resettableTransport) BulkOut(p []byte, d time.Duration) error {
	if r.wedged {
		return ErrNotResponding
	}
	return r.fakeTransport.BulkOut(p, d)
}

func (r *resettableTransport) BulkIn(p []byte, d time.Duration) (int, error) {
	if r.wedged {
		return 0, ErrNotResponding
	}
	return r.fakeTransport.BulkIn(p, d)
}

func (r *resettableTransport) Reset() error {
	r.resets++
	if r.failing {
		return errors.New("reset failed")
	}
	r.wedged = false
	return nil
}

// A wedged camera must be recovered in place. The alternative is telling the
// user to unplug the body, which during a timed sequence is no alternative.
func TestOpenRecoversWedgedDeviceByResetting(t *testing.T) {
	tr := &resettableTransport{wedged: true}
	tr.in = [][]byte{mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)}
	s := NewSession(tr)

	if err := s.Open(); err != nil {
		t.Fatalf("Open did not recover a wedged device: %v", err)
	}
	if tr.resets != 1 {
		t.Errorf("Reset called %d times, want exactly 1", tr.resets)
	}
}

// The retry must restart the transaction counter. A session that opened after a
// reset but kept counting would have every later reply rejected as belonging to
// a different transaction.
func TestOpenAfterResetRestartsTransactionIDs(t *testing.T) {
	tr := &resettableTransport{wedged: true}
	tr.in = [][]byte{mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)}
	s := NewSession(tr)

	if err := s.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	if s.txID != 1 {
		t.Errorf("txID = %d after a reset-and-reopen, want 1", s.txID)
	}
}

// A failed reset must not be dressed up as success: the original error is what
// tells the user the camera is genuinely gone.
func TestOpenReportsFailureWhenResetFails(t *testing.T) {
	tr := &resettableTransport{wedged: true, failing: true}
	s := NewSession(tr)

	err := s.Open()
	if err == nil {
		t.Fatal("Open reported success though the device stayed wedged")
	}
	if !errors.Is(err, ErrNotResponding) {
		t.Errorf("err = %v, want it to wrap ErrNotResponding", err)
	}
}

// A transport with no Reset must still work, and must not be retried blindly.
func TestOpenWithoutResetterIsUnchanged(t *testing.T) {
	f := &fakeTransport{}
	s := NewSession(f)
	if err := s.Open(); err == nil {
		t.Fatal("expected an error from a transport that answers nothing")
	}
}

// A healthy camera must not be reset. Resetting one that is answering perfectly
// well would abandon whatever it was doing.
func TestOpenDoesNotResetAHealthyDevice(t *testing.T) {
	tr := &resettableTransport{}
	tr.in = [][]byte{mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)}
	s := NewSession(tr)

	if err := s.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	if tr.resets != 0 {
		t.Errorf("Reset called %d times on a healthy device, want 0", tr.resets)
	}
}
