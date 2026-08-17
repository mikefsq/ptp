package ptp

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// mkObject builds a GetObject reply: a data container of size bytes followed by
// an OK response, split into transfers of the given size the way a real bulk
// endpoint delivers them.
func mkObject(size, transfer int, txID uint32) ([][]byte, []byte) {
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte(i * 7)
	}
	full := mkContainer(ContainerData, uint16(OpGetObject), txID, nil, payload)
	var chunks [][]byte
	for off := 0; off < len(full); off += transfer {
		end := min(off+transfer, len(full))
		chunks = append(chunks, full[off:end])
	}
	chunks = append(chunks, mkContainer(ContainerResponse, uint16(RespOK), txID, nil, nil))
	return chunks, payload
}

// A multi-megabyte object must reassemble across many bulk transfers. This is
// the path a real photo download takes, which synthetic small payloads miss.
func TestGetObjectLargeReassembly(t *testing.T) {
	const size = 3 * 1024 * 1024 // a NEX-6 JPEG is about this
	chunks, want := mkObject(size, 512*1024, 1)
	f := &fakeTransport{in: chunks}
	s := NewSession(f)

	got, err := s.GetObject(0x000C0005)
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	if len(got) != size {
		t.Fatalf("got %d bytes, want %d", len(got), size)
	}
	if !bytes.Equal(got, want) {
		t.Error("payload does not round-trip")
	}
}

// The final transfer may overshoot the declared length; that must be trimmed,
// not treated as a failure, or a completed download is thrown away.
func TestGetObjectTrailingBytesAreTrimmed(t *testing.T) {
	const size = 4096
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte(i)
	}
	full := mkContainer(ContainerData, uint16(OpGetObject), 1, nil, payload)
	padded := append(append([]byte(nil), full...), make([]byte, 64)...)

	f := &fakeTransport{in: [][]byte{
		padded,
		mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil),
	}}
	s := NewSession(f)

	got, err := s.GetObject(1)
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	if len(got) != size {
		t.Fatalf("got %d bytes, want %d (padding should be trimmed)", len(got), size)
	}
	if !bytes.Equal(got, payload) {
		t.Error("payload corrupted by trimming")
	}
}

// A desynchronised read can produce an absurd length. Sizing an allocation from
// it would try to reserve gigabytes, so it must be rejected up front.
func TestAbsurdContainerLengthRejected(t *testing.T) {
	bad := mkContainer(ContainerData, uint16(OpGetObject), 1, nil, make([]byte, 32))
	binary.LittleEndian.PutUint32(bad[0:], 0xFFFFFFFF)

	f := &fakeTransport{in: [][]byte{bad}}
	s := NewSession(f)
	if _, err := s.GetObject(1); err == nil {
		t.Fatal("expected an error for a 0xFFFFFFFF container length")
	}
}

// A download cut short must report how far it got, not hand back a truncated
// file as if it were complete.
func TestTruncatedDownloadIsAnError(t *testing.T) {
	chunks, _ := mkObject(1024*1024, 256*1024, 1)
	// Deliver only the first two transfers, then run dry.
	f := &fakeTransport{in: chunks[:2]}
	s := NewSession(f)
	if _, err := s.GetObject(1); err == nil {
		t.Fatal("expected an error when the download is cut short")
	}
}

func TestParseObjectInfoRoundTrip(t *testing.T) {
	// Build an ObjectInfo the way a camera lays it out.
	var b []byte
	u32 := func(v uint32) { b = binary.LittleEndian.AppendUint32(b, v) }
	u16 := func(v uint16) { b = binary.LittleEndian.AppendUint16(b, v) }
	str := func(s string) {
		if s == "" {
			b = append(b, 0)
			return
		}
		r := []rune(s)
		b = append(b, byte(len(r)+1))
		for _, c := range r {
			b = binary.LittleEndian.AppendUint16(b, uint16(c))
		}
		b = binary.LittleEndian.AppendUint16(b, 0)
	}
	u32(0x00010001) // storage
	u16(FormatEXIFJPEG)
	u16(0) // protection
	u32(3080192)
	u16(FormatEXIFJPEG)
	u32(12288)
	u32(160)
	u32(120)
	u32(4912)
	u32(3264)
	u32(24)
	u32(0x00080000) // parent
	u16(0)
	u32(0)
	u32(449)
	str("_DSC0449.JPG")
	str("20200831T142530")
	str("20200831T142530")
	str("")

	oi, err := ParseObjectInfo(b)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if oi.Filename != "_DSC0449.JPG" {
		t.Errorf("filename = %q", oi.Filename)
	}
	if oi.CompressedSize != 3080192 {
		t.Errorf("size = %d", oi.CompressedSize)
	}
	if oi.ImagePixWidth != 4912 || oi.ImagePixHeight != 3264 {
		t.Errorf("dimensions = %dx%d, want 4912x3264", oi.ImagePixWidth, oi.ImagePixHeight)
	}
	if oi.IsFolder() {
		t.Error("a JPEG should not report as a folder")
	}
	if got := oi.Captured(); got.Year() != 2020 || got.Month() != 8 || got.Day() != 31 {
		t.Errorf("captured = %v, want 2020-08-31", got)
	}
}
