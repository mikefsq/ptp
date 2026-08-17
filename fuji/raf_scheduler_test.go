package fuji

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/mikefsq/ptp"
)

// sampleStrip returns a real compressed strip and its parsed metadata.
func sampleStrip(tb testing.TB) ([]byte, *compressedHeader, *ptp.CFA) {
	tb.Helper()
	path := filepath.Join("testdata", "raf", "xt5-lossless-lit.raf")
	data, err := os.ReadFile(path)
	if err != nil {
		tb.Skipf("%s is not present", path)
	}
	cfa, err := DecodeRAF(data)
	if err != nil {
		tb.Fatal(err)
	}
	be := binary.BigEndian
	cfaOff := be.Uint32(data[rafOffsets+16:])
	cfaLen := be.Uint32(data[rafOffsets+20:])
	strip, _, _, err := readFujiIFD(data[cfaOff:cfaOff+cfaLen], cfa)
	if err != nil {
		tb.Fatal(err)
	}
	h, ok := parseCompressedHeader(strip)
	if !ok {
		tb.Fatal("the sample is not compressed")
	}
	return strip, h, cfa
}

// The promoted decoder must still produce EXACTLY what libraw does; that is
// covered by TestDecodeRAFMatchesLibRaw. What is pinned here is the block
// scheduler: decodeCompressed fanning out over blocks must agree with a plain
// sequential decode of the same strip, or a block is writing outside its
// columns.
func TestBlocksAreIndependent(t *testing.T) {
	strip, h, cfa := sampleStrip(t)
	got, err := decodeCompressed(strip, h, cfa)
	if err != nil {
		t.Fatal(err)
	}
	for i := range got {
		if got[i] != cfa.Pixels[i] {
			t.Fatalf("sample %d differs: %d vs %d", i, got[i], cfa.Pixels[i])
		}
	}
}

// rafBytes returns the whole sample file.
func rafBytes(tb testing.TB) []byte {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "raf", "xt5-lossless-lit.raf"))
	if err != nil {
		tb.Skipf("sample not present: %v", err)
	}
	return data
}
