package fuji

import (
	"encoding/binary"
	"time"

	"github.com/mikefsq/ptp"
)

// fakeTransport replays a scripted sequence of bulk-IN transfers and records
// every bulk-OUT write, so a camera's side of an exchange can be staged without
// a camera. The core has its own copy; duplicating the few lines here keeps the
// core from having to export test scaffolding.
type fakeTransport struct {
	in    [][]byte
	out   [][]byte
	evt   [][]byte
	inErr error
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
		return 0, ptp.ErrTimeout
	}
	next := f.in[0]
	f.in = f.in[1:]
	return copy(p, next), nil
}

func (f *fakeTransport) InterruptIn(p []byte, _ time.Duration) (int, error) {
	if len(f.evt) == 0 {
		return 0, ptp.ErrTimeout
	}
	next := f.evt[0]
	f.evt = f.evt[1:]
	return copy(p, next), nil
}

func (f *fakeTransport) MaxPacketSize() int { return 512 }
func (f *fakeTransport) Close() error       { return nil }

// mkContainer builds a PTP container for the fake camera to return.
func mkContainer(typ, code uint16, txID uint32, params []uint32, payload []byte) []byte {
	n := ptp.ContainerHeaderLen + len(params)*4 + len(payload)
	b := make([]byte, n)
	binary.LittleEndian.PutUint32(b[0:], uint32(n))
	binary.LittleEndian.PutUint16(b[4:], typ)
	binary.LittleEndian.PutUint16(b[6:], code)
	binary.LittleEndian.PutUint32(b[8:], txID)
	for i, p := range params {
		binary.LittleEndian.PutUint32(b[ptp.ContainerHeaderLen+i*4:], p)
	}
	copy(b[ptp.ContainerHeaderLen+len(params)*4:], payload)
	return b
}

// ok is a success response for the given transaction.
func ok(txID uint32) []byte {
	return mkContainer(ptp.ContainerResponse, uint16(ptp.RespOK), txID, nil, nil)
}

// refused is the camera's "not right now" — vendor code 0xA002, which a
// Fujifilm body returns while it is still busy.
func refused(txID uint32) []byte {
	return mkContainer(ptp.ContainerResponse, 0xA002, txID, nil, nil)
}

// refusedCode is a vendor refusal with a specific code.
func refusedCode(txID uint32, code uint16) []byte {
	return mkContainer(ptp.ContainerResponse, code, txID, nil, nil)
}

// openCamera builds a Camera over a fake whose session is already open.
func openCamera(f *fakeTransport) *Camera {
	f.in = append([][]byte{ok(1)}, f.in...) // answer OpenSession
	c, err := New(f)
	if err != nil {
		panic(err)
	}
	return c
}
