package sony

import (
	"bytes"
	"testing"
)

// The live view payload is not bare JPEG — Sony prefixes it with a header whose
// length is not pinned. Finding the SOI is what makes that safe: a hardcoded
// offset would silently produce a corrupt image if the header were a different
// size on another body.
func TestTrimToJPEGSkipsLeadingHeader(t *testing.T) {
	img := append(append([]byte{0xFF, 0xD8}, []byte("payload")...), 0xFF, 0xD9)
	header := []byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x10, 0x20}

	got, err := trimToJPEG(append(header, img...))
	if err != nil {
		t.Fatalf("trimToJPEG: %v", err)
	}
	if !bytes.Equal(got, img) {
		t.Errorf("got %x, want %x", got, img)
	}
}

// Sony pads object transfers to a block boundary — measured on a NEX-6, a
// 1,611,005-byte JPEG arrived declared as 1,638,400. Those trailing bytes must
// come off, or every frame carries junk.
func TestTrimToJPEGDropsTrailingPadding(t *testing.T) {
	img := append(append([]byte{0xFF, 0xD8}, []byte("payload")...), 0xFF, 0xD9)
	padded := append(append([]byte{}, img...), make([]byte, 4096)...)

	got, err := trimToJPEG(padded)
	if err != nil {
		t.Fatalf("trimToJPEG: %v", err)
	}
	if len(got) != len(img) {
		t.Errorf("got %d bytes, want %d — trailing padding was not trimmed", len(got), len(img))
	}
}

// A payload with no JPEG in it must be an error, not a truncated frame. If the
// format is not what this driver expects, saying so is far more useful than
// handing back bytes that will not decode.
func TestTrimToJPEGRejectsNonJPEG(t *testing.T) {
	if _, err := trimToJPEG(bytes.Repeat([]byte{0x41}, 512)); err == nil {
		t.Fatal("accepted a payload containing no JPEG start marker")
	}
}

// An SOI with no EOI is returned as-is rather than rejected: a frame cut short
// by a transfer is still worth handing to a decoder, which will show what
// arrived.
func TestTrimToJPEGKeepsUnterminatedImage(t *testing.T) {
	partial := append([]byte{0xFF, 0xD8}, []byte("truncated")...)
	got, err := trimToJPEG(append([]byte{0x00, 0x00}, partial...))
	if err != nil {
		t.Fatalf("trimToJPEG: %v", err)
	}
	if !bytes.Equal(got, partial) {
		t.Errorf("got %x, want %x", got, partial)
	}
}

// The handle to the liveview object
func TestLiveViewObjectHandle(t *testing.T) {
	if LiveViewObject != 0xFFFFC002 {
		t.Errorf("LiveViewObject = 0x%08X, want 0xFFFFC002", LiveViewObject)
	}
	// Via a variable: converting the constant directly is a compile error,
	// which is itself the point — this value only makes sense reinterpreted.
	h := LiveViewObject
	if int32(h) != -0x3ffe {
		t.Errorf("0x%08X is not the -0x3ffe the adapter loads", h)
	}
}
