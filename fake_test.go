package ptp

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"
)

// fakeTransport replays a scripted sequence of bulk-IN transfers and records
// every bulk-OUT write, so the framing layer can be exercised without a camera.
type fakeTransport struct {
	in     [][]byte // queued bulk-IN transfers, consumed in order
	out    [][]byte // captured bulk-OUT writes
	evt    [][]byte // queued interrupt-IN packets
	maxPkt int
	inErr  error // returned once the queue is empty
}

func (f *fakeTransport) BulkOut(p []byte, _ time.Duration) error {
	f.out = append(f.out, append([]byte(nil), p...))
	return nil
}

func (f *fakeTransport) BulkIn(p []byte, _ time.Duration) (int, error) {
	if len(f.in) == 0 {
		if f.inErr != nil {
			return 0, f.inErr
		}
		return 0, ErrTimeout
	}
	next := f.in[0]
	n := copy(p, next)
	// A short destination does not lose the rest of the transfer: the endpoint
	// still holds it and the next read continues from there. Dropping it here
	// would hide every bug in a reader that sizes its destination.
	if n < len(next) {
		f.in[0] = next[n:]
	} else {
		f.in = f.in[1:]
	}
	return n, nil
}

func (f *fakeTransport) InterruptIn(p []byte, _ time.Duration) (int, error) {
	if len(f.evt) == 0 {
		return 0, ErrTimeout
	}
	next := f.evt[0]
	f.evt = f.evt[1:]
	return copy(p, next), nil
}

func (f *fakeTransport) MaxPacketSize() int {
	if f.maxPkt == 0 {
		return 512
	}
	return f.maxPkt
}

func (f *fakeTransport) Close() error { return nil }

// mkContainer builds a PTP container for the fake device to return.
func mkContainer(typ, code uint16, txID uint32, params []uint32, payload []byte) []byte {
	n := ContainerHeaderLen + len(params)*4 + len(payload)
	b := make([]byte, n)
	binary.LittleEndian.PutUint32(b[0:], uint32(n))
	binary.LittleEndian.PutUint16(b[4:], typ)
	binary.LittleEndian.PutUint16(b[6:], code)
	binary.LittleEndian.PutUint32(b[8:], txID)
	for i, p := range params {
		binary.LittleEndian.PutUint32(b[ContainerHeaderLen+i*4:], p)
	}
	copy(b[ContainerHeaderLen+len(params)*4:], payload)
	return b
}

func TestOpenSessionSendsCorrectCommand(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
	}}
	s := NewSession(f)
	if err := s.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	if len(f.out) != 1 {
		t.Fatalf("wrote %d transfers, want 1", len(f.out))
	}
	got := f.out[0]
	if n := binary.LittleEndian.Uint32(got[0:]); n != 16 {
		t.Errorf("container length = %d, want 16 (header + 1 param)", n)
	}
	if typ := binary.LittleEndian.Uint16(got[4:]); typ != ContainerCommand {
		t.Errorf("container type = %d, want %d", typ, ContainerCommand)
	}
	if code := binary.LittleEndian.Uint16(got[6:]); code != uint16(OpOpenSession) {
		t.Errorf("opcode = 0x%04X, want 0x%04X", code, uint16(OpOpenSession))
	}
	if tx := binary.LittleEndian.Uint32(got[8:]); tx != 1 {
		t.Errorf("transaction ID = %d, want 1", tx)
	}
	if sid := binary.LittleEndian.Uint32(got[12:]); sid != 1 {
		t.Errorf("session ID = %d, want 1", sid)
	}
}

func TestTransactionIDIncrements(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
		mkContainer(ContainerResponse, uint16(RespOK), 2, nil, nil),
		mkContainer(ContainerResponse, uint16(RespOK), 3, nil, nil),
	}}
	s := NewSession(f)
	if err := s.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	for i := range 2 {
		if _, _, err := s.Do(OpGetDeviceInfo, nil, nil, DefaultTimeout); err != nil {
			t.Fatalf("Do %d: %v", i, err)
		}
	}
	for i, w := range f.out {
		if tx := binary.LittleEndian.Uint32(w[8:]); tx != uint32(i+1) {
			t.Errorf("transfer %d has transaction ID %d, want %d", i, tx, i+1)
		}
	}
}

func TestDataInPhase(t *testing.T) {
	payload := []byte{1, 2, 3, 4, 5, 6, 7}
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerData, uint16(OpGetDeviceInfo), 1, nil, payload),
		mkContainer(ContainerResponse, uint16(RespOK), 1, []uint32{0xAA}, nil),
	}}
	s := NewSession(f)
	data, params, err := s.Do(OpGetDeviceInfo, nil, nil, DefaultTimeout)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if string(data) != string(payload) {
		t.Errorf("data = %v, want %v", data, payload)
	}
	if len(params) != 1 || params[0] != 0xAA {
		t.Errorf("response params = %v, want [0xAA]", params)
	}
}

func TestDataOutPhase(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
	}}
	s := NewSession(f)
	body := []byte{0xDE, 0xAD}
	if _, _, err := s.Do(OpSetDevicePropValue, []uint32{0x5007}, body, DefaultTimeout); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(f.out) != 2 {
		t.Fatalf("wrote %d transfers, want 2 (command + data)", len(f.out))
	}
	d := f.out[1]
	if typ := binary.LittleEndian.Uint16(d[4:]); typ != ContainerData {
		t.Errorf("second transfer type = %d, want %d", typ, ContainerData)
	}
	if got := d[ContainerHeaderLen:]; string(got) != string(body) {
		t.Errorf("data payload = %v, want %v", got, body)
	}
	// Command and data must share a transaction ID.
	if a, b := binary.LittleEndian.Uint32(f.out[0][8:]), binary.LittleEndian.Uint32(d[8:]); a != b {
		t.Errorf("command txID %d != data txID %d", a, b)
	}
}

func TestErrorResponse(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespDeviceBusy), 1, nil, nil),
	}}
	s := NewSession(f)
	_, _, err := s.Do(OpGetObject, nil, nil, DefaultTimeout)
	if err == nil {
		t.Fatal("expected an error for a DeviceBusy response")
	}
	var pe *Error
	if !errors.As(err, &pe) {
		t.Fatalf("error is %T, want *ptp.Error", err)
	}
	if pe.Code != RespDeviceBusy {
		t.Errorf("code = %v, want DeviceBusy", pe.Code)
	}
	if pe.Op != OpGetObject {
		t.Errorf("op = 0x%04X, want 0x%04X", uint16(pe.Op), uint16(OpGetObject))
	}
}

// A large data phase arrives across several bulk transfers; the reader must
// reassemble it using the declared container length.
func TestMultiTransferReassembly(t *testing.T) {
	payload := make([]byte, 3000)
	for i := range payload {
		payload[i] = byte(i)
	}
	full := mkContainer(ContainerData, uint16(OpGetObject), 1, nil, payload)
	f := &fakeTransport{in: [][]byte{
		full[:1024], full[1024:2048], full[2048:],
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
	}}
	s := NewSession(f)
	data, _, err := s.Do(OpGetObject, nil, nil, DefaultTimeout)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(data) != len(payload) {
		t.Fatalf("got %d bytes, want %d", len(data), len(payload))
	}
	for i := range data {
		if data[i] != payload[i] {
			t.Fatalf("byte %d = %d, want %d", i, data[i], payload[i])
		}
	}
}

func TestTruncatedContainerIsAnError(t *testing.T) {
	payload := make([]byte, 500)
	full := mkContainer(ContainerData, uint16(OpGetObject), 1, nil, payload)
	// Deliver only the first half and then run dry.
	f := &fakeTransport{in: [][]byte{full[:200]}}
	s := NewSession(f)
	if _, _, err := s.Do(OpGetObject, nil, nil, DefaultTimeout); err == nil {
		t.Fatal("expected an error when the container is truncated")
	}
}

func TestTooManyParametersRejected(t *testing.T) {
	s := NewSession(&fakeTransport{})
	_, _, err := s.Do(OpGetDeviceInfo, []uint32{1, 2, 3, 4, 5, 6}, nil, DefaultTimeout)
	if err == nil {
		t.Fatal("expected an error for 6 parameters; PTP allows 5")
	}
}
