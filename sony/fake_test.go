package sony

import (
	"encoding/binary"
	"time"

	"github.com/mikefsq/ptp"
)

// fakeTransport replays a scripted sequence of bulk-IN transfers and records
// every bulk-OUT write, so a camera's side of an exchange can be staged without
// a camera.
type fakeTransport struct {
	in    [][]byte
	out   [][]byte
	evt   [][]byte
	inErr error

	// lastTx is the transaction ID of the most recent command container. Replies
	// are stamped with it on the way out, which is what a real camera does and
	// what the session now checks — so a scripted exchange does not have to know
	// how many transactions ran before it.
	lastTx uint32
}

func (f *fakeTransport) BulkOut(p []byte, _ time.Duration) error {
	f.out = append(f.out, append([]byte(nil), p...))
	if len(p) >= ptp.ContainerHeaderLen && binary.LittleEndian.Uint16(p[4:]) == ptp.ContainerCommand {
		f.lastTx = binary.LittleEndian.Uint32(p[8:])
	}
	return nil
}

func (f *fakeTransport) BulkIn(p []byte, _ time.Duration) (int, error) {
	if len(f.in) == 0 {
		if f.inErr != nil {
			return 0, f.inErr
		}
		return 0, ptp.ErrTimeout
	}
	next := append([]byte(nil), f.in[0]...)
	f.in = f.in[1:]
	if len(next) >= ptp.ContainerHeaderLen {
		if typ := binary.LittleEndian.Uint16(next[4:]); typ == ptp.ContainerResponse || typ == ptp.ContainerData {
			binary.LittleEndian.PutUint32(next[8:], f.lastTx)
		}
	}
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

// okRun queues n success responses, for a sequence whose exact length is not
// the point of the test.
func okRun(n int) [][]byte {
	out := make([][]byte, n)
	for i := range out {
		out[i] = ok(uint32(i + 1))
	}
	return out
}

// openCamera builds a Camera over a fake whose session is already open, without
// the vendor handshake. The transfers the open itself made are discarded, so a
// test counting writes sees only its own.
func openCamera(f *fakeTransport) *Camera {
	f.in = append([][]byte{ok(1)}, f.in...) // answer OpenSession
	c, err := New(f)
	if err != nil {
		panic(err)
	}
	f.out = nil
	return c
}

// sentOps lists the operation codes written as command containers.
func sentOps(f *fakeTransport) []ptp.OpCode {
	var ops []ptp.OpCode
	for _, w := range f.out {
		if len(w) < ptp.ContainerHeaderLen {
			continue
		}
		if binary.LittleEndian.Uint16(w[4:]) != ptp.ContainerCommand {
			continue
		}
		ops = append(ops, ptp.OpCode(binary.LittleEndian.Uint16(w[6:])))
	}
	return ops
}

// buttonWrites lists the (control, value) pairs sent through 0x9207, in order.
// A control write is a data container carrying a uint16, whose command
// container named the control in parameter 1.
func buttonWrites(f *fakeTransport) []struct {
	Ctrl ControlCode
	Val  uint16
} {
	var out []struct {
		Ctrl ControlCode
		Val  uint16
	}
	var pending ControlCode
	var armed bool
	for _, w := range f.out {
		if len(w) < ptp.ContainerHeaderLen {
			continue
		}
		typ := binary.LittleEndian.Uint16(w[4:])
		code := ptp.OpCode(binary.LittleEndian.Uint16(w[6:]))
		switch typ {
		case ptp.ContainerCommand:
			armed = code == OpSetControlDeviceB
			if armed && len(w) >= ptp.ContainerHeaderLen+4 {
				pending = ControlCode(binary.LittleEndian.Uint32(w[ptp.ContainerHeaderLen:]))
			}
		case ptp.ContainerData:
			if armed && len(w) >= ptp.ContainerHeaderLen+2 {
				out = append(out, struct {
					Ctrl ControlCode
					Val  uint16
				}{pending, binary.LittleEndian.Uint16(w[ptp.ContainerHeaderLen:])})
			}
		}
	}
	return out
}
