package ptp

import (
	"bytes"
	"testing"
	"time"
)

// A USB bulk read can return the data container and the response container in
// one transfer. IOKit's ReadPipeTO completes on a short packet, and when both
// containers are already queued on the endpoint the host controller hands them
// over together — so a single BulkIn yields data+response concatenated.
//
// This is not a hypothetical: it is what hung the X-T5. The reader trimmed the
// overshoot and discarded the response, then blocked waiting for a response it
// had already thrown away, and reported "camera is not responding".
func TestCoalescedDataAndResponse(t *testing.T) {
	payload := []byte("descriptor-body")
	data := mkContainer(ContainerData, 0x1014, 1, nil, payload)
	resp := mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)

	f := &fakeTransport{in: [][]byte{append(append([]byte{}, data...), resp...)}}
	s := NewSession(f)
	s.txID = 0

	got, _, err := s.Do(OpGetDevicePropDesc, []uint32{0x5003}, nil, time.Second)
	if err != nil {
		t.Fatalf("coalesced data+response must not hang the transaction: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload = %q, want %q", got, payload)
	}
}

// The same coalescing across two transactions: the residue from the first read
// is the start of the second transaction's data, and must not be dropped.
func TestResidueCarriesToNextTransaction(t *testing.T) {
	d1 := mkContainer(ContainerData, 0x1014, 1, nil, []byte("first"))
	r1 := mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)
	d2 := mkContainer(ContainerData, 0x1014, 2, nil, []byte("second"))
	r2 := mkContainer(ContainerResponse, uint16(RespOK), 2, nil, nil)

	all := bytes.Join([][]byte{d1, r1, d2, r2}, nil)
	f := &fakeTransport{in: [][]byte{all}}
	s := NewSession(f)
	s.txID = 0

	for i, want := range []string{"first", "second"} {
		got, _, err := s.Do(OpGetDevicePropDesc, []uint32{0x5003}, nil, time.Second)
		if err != nil {
			t.Fatalf("transaction %d: %v", i+1, err)
		}
		if string(got) != want {
			t.Fatalf("transaction %d payload = %q, want %q", i+1, got, want)
		}
	}
}

// A response carrying the wrong transaction ID means the stream is one
// transaction behind. Accepting it silently returns the PREVIOUS exchange's
// data as though it answered the current request.
func TestStaleTransactionIDRejected(t *testing.T) {
	data := mkContainer(ContainerData, 0x1014, 41, nil, []byte("stale"))
	resp := mkContainer(ContainerResponse, uint16(RespOK), 41, nil, nil)
	f := &fakeTransport{in: [][]byte{data, resp}}
	s := NewSession(f)
	s.txID = 0

	if _, _, err := s.Do(OpGetDevicePropDesc, []uint32{0x5003}, nil, time.Second); err == nil {
		t.Fatal("a response for transaction 41 must not satisfy transaction 1")
	}
}

// A late reply from an EARLIER transaction — seen on a real X-T5, which
// delivered transaction 28's containers while transaction 29 was running. The
// session must skip past them and find its own answer, not fail: an eclipse
// sequence cannot afford to stop dead because one reply arrived late.
func TestLateReplyIsSkipped(t *testing.T) {
	late := mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)
	lateData := mkContainer(ContainerData, 0x1014, 1, nil, []byte("previous"))
	mine := mkContainer(ContainerData, 0x1014, 2, nil, []byte("mine"))
	myResp := mkContainer(ContainerResponse, uint16(RespOK), 2, nil, nil)

	f := &fakeTransport{in: [][]byte{lateData, late, mine, myResp}}
	s := NewSession(f)
	s.txID = 1 // next transaction is 2

	got, _, err := s.Do(OpGetDevicePropDesc, []uint32{0x5003}, nil, time.Second)
	if err != nil {
		t.Fatalf("must resynchronise past a late reply: %v", err)
	}
	if string(got) != "mine" {
		t.Fatalf("payload = %q, want %q (stale data was accepted)", got, "mine")
	}
}

// The mirror of TestStaleTransactionIDRejected. A reply carrying a transaction
// ID from the FUTURE is just as wrong as one from the past, and is what a
// camera sends when the host and device have lost sync in the other direction —
// after a session was reopened without the counter restarting, say. Accepting
// it would hand the caller another transaction's data.
func TestFutureTransactionIDRejected(t *testing.T) {
	data := mkContainer(ContainerData, 0x1014, 41, nil, []byte("not ours"))
	resp := mkContainer(ContainerResponse, uint16(RespOK), 41, nil, nil)
	f := &fakeTransport{in: [][]byte{data, resp}}
	s := NewSession(f)
	s.txID = 0

	if _, _, err := s.Do(OpGetDevicePropDesc, []uint32{0x5003}, nil, time.Second); err == nil {
		t.Fatal("a response for transaction 41 must not satisfy transaction 1")
	}
}

// A frame larger than the read chunk takes many reads, and those land DIRECTLY
// in the destination rather than being staged and copied. A 61 MP RAW is 38 MB,
// so staging it costs a second 38 MB memcpy per frame.
//
// The small-payload tests above never reach that loop: they complete in the
// first read.
func TestLargeContainerAssembledAcrossReads(t *testing.T) {
	// Big enough to need several 512 KB reads, and not a multiple of the chunk
	// so the final partial read is exercised too.
	payload := make([]byte, 3*512*1024+1234)
	for i := range payload {
		payload[i] = byte(i*7 + i/251)
	}
	data := mkContainer(ContainerData, 0x1009, 1, nil, payload)
	resp := mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)

	f := &fakeTransport{in: [][]byte{data, resp}}
	s := NewSession(f)
	s.txID = 0

	got, _, err := s.Do(OpGetObject, []uint32{1}, nil, time.Second)
	if err != nil {
		t.Fatalf("large transfer: %v", err)
	}
	if len(got) != len(payload) {
		t.Fatalf("got %d bytes, want %d", len(got), len(payload))
	}
	if !bytes.Equal(got, payload) {
		for i := range got {
			if got[i] != payload[i] {
				t.Fatalf("first difference at byte %d: got %#02x want %#02x", i, got[i], payload[i])
			}
		}
	}
}

// The same, with the response coalesced onto the end of the final read. The
// last read is the one that can overshoot, and losing the response there is
// what hung the X-T5.
func TestLargeContainerWithCoalescedResponse(t *testing.T) {
	payload := make([]byte, 2*512*1024+777)
	for i := range payload {
		payload[i] = byte(i ^ 0x5A)
	}
	data := mkContainer(ContainerData, 0x1009, 1, nil, payload)
	resp := mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)

	// One continuous stream: the reader must find the response inside the tail
	// of the data transfer.
	f := &fakeTransport{in: [][]byte{append(append([]byte{}, data...), resp...)}}
	s := NewSession(f)
	s.txID = 0

	got, _, err := s.Do(OpGetObject, []uint32{1}, nil, time.Second)
	if err != nil {
		t.Fatalf("coalesced response after a large transfer: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload corrupted: %d bytes, want %d", len(got), len(payload))
	}
}

// Trailing zeros after a large frame are device padding, not a container, and
// must be dropped — a NEX-6 sends thousands after a JPEG. Keeping them makes
// the next read try to parse zeros as a header.
func TestLargeContainerWithTrailingPadding(t *testing.T) {
	payload := make([]byte, 512*1024+64)
	for i := range payload {
		payload[i] = byte(i)
	}
	data := mkContainer(ContainerData, 0x1009, 1, nil, payload)
	padded := append(append([]byte{}, data...), make([]byte, 4096)...)
	resp := mkContainer(ContainerResponse, uint16(RespOK), 1, nil, nil)

	f := &fakeTransport{in: [][]byte{padded, resp}}
	s := NewSession(f)
	s.txID = 0

	got, _, err := s.Do(OpGetObject, []uint32{1}, nil, time.Second)
	if err != nil {
		t.Fatalf("padded transfer: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload corrupted or padding kept: got %d bytes, want %d", len(got), len(payload))
	}
}
