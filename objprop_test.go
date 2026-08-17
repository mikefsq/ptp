package ptp

import (
	"encoding/binary"
	"testing"
	"time"
)

// oplBuilder assembles a GetObjPropList reply.
type oplBuilder struct {
	rows []byte
	n    uint32
}

func (b *oplBuilder) row(handle uint32, p ObjProp, t DataType, enc func()) {
	b.n++
	b.rows = binary.LittleEndian.AppendUint32(b.rows, handle)
	b.rows = binary.LittleEndian.AppendUint16(b.rows, uint16(p))
	b.rows = binary.LittleEndian.AppendUint16(b.rows, uint16(t))
	enc()
}

func (b *oplBuilder) u32(v uint32) { b.rows = binary.LittleEndian.AppendUint32(b.rows, v) }
func (b *oplBuilder) u64(v uint64) { b.rows = binary.LittleEndian.AppendUint64(b.rows, v) }
func (b *oplBuilder) str(s string) { b.rows = append(b.rows, EncodeString(s)...) }

func (b *oplBuilder) bytes() []byte {
	out := binary.LittleEndian.AppendUint32(nil, b.n)
	return append(out, b.rows...)
}

// The mixed-type case is the point: a single reply interleaves integers and
// strings of different widths, and the reader must stay aligned across them.
func TestParseObjPropListMixedTypes(t *testing.T) {
	b := &oplBuilder{}
	b.row(0x0C0005, OPCObjectFormat, TypeUint16, func() {
		b.rows = binary.LittleEndian.AppendUint16(b.rows, FormatEXIFJPEG)
	})
	b.row(0x0C0005, OPCObjectSize, TypeUint64, func() { b.u64(3080192) })
	b.row(0x0C0005, OPCObjectFileName, TypeString, func() { b.str("_DSC0449.JPG") })
	b.row(0x0C0009, OPCObjectFileName, TypeString, func() { b.str("_DSC0450.JPG") })
	b.row(0x0C0009, OPCWidth, TypeUint32, func() { b.u32(4912) })
	b.row(0x0C0009, OPCHeight, TypeUint32, func() { b.u32(3264) })

	entries, err := ParseObjPropList(b.bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(entries) != 6 {
		t.Fatalf("got %d entries, want 6", len(entries))
	}

	byH := ByHandle(entries)
	if len(byH) != 2 {
		t.Fatalf("got %d objects, want 2", len(byH))
	}
	first := byH[0x0C0005]
	if got := first[OPCObjectFileName].Str; got != "_DSC0449.JPG" {
		t.Errorf("filename = %q", got)
	}
	if got := first[OPCObjectSize].Num; got != 3080192 {
		t.Errorf("size = %d, want 3080192", got)
	}
	second := byH[0x0C0009]
	if second[OPCWidth].Num != 4912 || second[OPCHeight].Num != 3264 {
		t.Errorf("dimensions = %dx%d, want 4912x3264", second[OPCWidth].Num, second[OPCHeight].Num)
	}
}

func TestParseObjPropListAbsurdCount(t *testing.T) {
	b := binary.LittleEndian.AppendUint32(nil, 0xFFFFFFF)
	if _, err := ParseObjPropList(b); err == nil {
		t.Fatal("expected an error for a count the buffer cannot hold")
	}
}

func TestParseObjPropListTruncatedReportsProgress(t *testing.T) {
	b := &oplBuilder{}
	b.row(1, OPCObjectSize, TypeUint64, func() { b.u64(100) })
	raw := b.bytes()
	binary.LittleEndian.PutUint32(raw[0:], 2) // claim 2 rows, supply 1

	entries, err := ParseObjPropList(raw)
	if err == nil {
		t.Fatal("expected an error when rows are missing")
	}
	if len(entries) != 1 {
		t.Errorf("got %d parsed entries, want the 1 good row", len(entries))
	}
}

func TestEncodeStringRoundTrip(t *testing.T) {
	for _, s := range []string{"", "NEX-6", "a", "ILCE-7RM5"} {
		enc := EncodeString(s)
		r := NewReader(enc)
		got, err := r.Str()
		if err != nil {
			t.Fatalf("%q: decode: %v", s, err)
		}
		if got != s {
			t.Errorf("round trip of %q gave %q", s, got)
		}
	}
}

func TestPollEventDecodesContainer(t *testing.T) {
	ev := mkContainer(ContainerEvent, uint16(EventObjectAdded), 7, []uint32{0x0C0005}, nil)
	f := &fakeTransport{evt: [][]byte{ev}}
	s := NewSession(f)

	got, err := s.PollEvent(time.Second)
	if err != nil {
		t.Fatalf("PollEvent: %v", err)
	}
	if got.Code != EventObjectAdded {
		t.Errorf("code = %v, want ObjectAdded", got.Code)
	}
	if len(got.Params) != 1 || got.Params[0] != 0x0C0005 {
		t.Errorf("params = %v, want [0xC0005]", got.Params)
	}
	if got.TxID != 7 {
		t.Errorf("txID = %d, want 7", got.TxID)
	}
}

// An idle camera produces timeouts continuously; that must not read as failure.
func TestPollEventTimeoutIsNotAnError(t *testing.T) {
	s := NewSession(&fakeTransport{})
	_, err := s.PollEvent(10 * time.Millisecond)
	if err != ErrTimeout {
		t.Fatalf("err = %v, want ErrTimeout", err)
	}
}

// A device overstating its container length must not walk the reader past what
// actually arrived.
func TestPollEventOverstatedLength(t *testing.T) {
	ev := mkContainer(ContainerEvent, uint16(EventDevicePropChanged), 1, []uint32{0x5001}, nil)
	binary.LittleEndian.PutUint32(ev[0:], 999)
	f := &fakeTransport{evt: [][]byte{ev}}
	s := NewSession(f)

	got, err := s.PollEvent(time.Second)
	if err != nil {
		t.Fatalf("PollEvent: %v", err)
	}
	if got.Code != EventDevicePropChanged {
		t.Errorf("code = %v", got.Code)
	}
	if len(got.Params) != 1 {
		t.Errorf("params = %v, want exactly the one that arrived", got.Params)
	}
}
