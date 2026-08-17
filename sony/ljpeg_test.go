package sony

import (
	"os"
	"path/filepath"
	"testing"
)

// The lossless-JPEG decoder is validated against a REAL stream from another
// vendor, because the codec is a published standard: an Adobe DNG, a Canon CR2
// and a Sony lossless ARW all go through the same decoder.
//
// Regenerate the fixture from any lossless-JPEG raw with:
//
//	(extract a tile from the raw SubIFD, Compression == 7)
func TestLJPEGDecodesARealStream(t *testing.T) {
	path := filepath.Join("testdata", "ljpeg", "tile.ljpeg")
	d, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("%s is not present", path)
	}
	f, err := parseLJPEG(d)
	if err != nil {
		t.Fatalf("parseLJPEG: %v", err)
	}
	if f.precision != 10 || f.width != 504 || f.height != 504 {
		t.Errorf("SOF3 says %d-bit %dx%d, want 10-bit 504x504", f.precision, f.width, f.height)
	}
	if len(f.comps) != 3 {
		t.Fatalf("%d components, want 3", len(f.comps))
	}
	if f.predictor != 7 {
		t.Errorf("predictor %d, want 7", f.predictor)
	}
	px, w, h, err := decodeLJPEG(f)
	if err != nil {
		t.Fatalf("decodeLJPEG: %v", err)
	}
	if w != 504*3 || h != 504 {
		t.Fatalf("decoded %dx%d, want %dx504", w, h, 504*3)
	}
	// Every sample must be inside the declared precision. A desynchronised
	// entropy decoder produces values all over the range, so this is a real
	// check rather than a formality.
	for i, v := range px {
		if v >= 1<<10 {
			t.Fatalf("sample %d is %d, past the 10-bit range", i, v)
		}
	}
	// The first row of an iPhone ProRAW tile is smooth: neighbouring pixels
	// differ by a few counts, not hundreds. Noise would fail this.
	big := 0
	for i := 3; i < 3*200; i++ {
		d := int(px[i]) - int(px[i-3])
		if d < 0 {
			d = -d
		}
		if d > 32 {
			big++
		}
	}
	if big > 5 {
		t.Errorf("%d large jumps in the first 200 pixels; the decode looks desynchronised", big)
	}
}

// The seven predictors are Table H.1 and must be exact: picking the wrong one
// yields an image that looks almost right and is wrong everywhere.
func TestLJPEGPredictors(t *testing.T) {
	const ra, rb, rc = 100, 80, 60
	for sel, want := range map[int]int32{
		1: ra,
		2: rb,
		3: rc,
		4: ra + rb - rc,
		5: ra + (rb-rc)>>1,
		6: rb + (ra-rc)>>1,
		7: (ra + rb) >> 1,
	} {
		if got := predictLJPEG(sel, ra, rb, rc); got != want {
			t.Errorf("predictor %d gave %d, want %d", sel, got, want)
		}
	}
}

// extend turns a magnitude into a signed difference (T.81 F.2.2.1). Getting it
// wrong flips the sign of half of all differences.
func TestLJPEGExtend(t *testing.T) {
	for _, tc := range []struct {
		v    int32
		s    int
		want int32
	}{
		{0, 0, 0},
		{0, 1, -1}, {1, 1, 1},
		{0, 2, -3}, {1, 2, -2}, {2, 2, 2}, {3, 3, -4},
	} {
		if got := extend(tc.v, tc.s); got != tc.want {
			t.Errorf("extend(%d,%d) = %d, want %d", tc.v, tc.s, got, tc.want)
		}
	}
}

// Byte stuffing: a literal 0xFF is written 0xFF 0x00, and 0xFF followed by
// anything else is a marker that ends the scan.
func TestLJPEGByteStuffing(t *testing.T) {
	b := &jbits{data: []byte{0xFF, 0x00, 0x55}}
	var got byte
	for i := 0; i < 8; i++ {
		got = got<<1 | byte(b.bit())
	}
	if got != 0xFF {
		t.Errorf("stuffed byte read as %#02x, want 0xFF", got)
	}
	b = &jbits{data: []byte{0xFF, 0xD9}}
	b.bit()
	if b.marker != 0xD9 {
		t.Errorf("marker %#02x not reported, so the scan would run past its end", b.marker)
	}
}

// A malformed stream must be refused, not decoded into noise.
func TestLJPEGRejectsRubbish(t *testing.T) {
	for _, tc := range []struct {
		name string
		d    []byte
	}{
		{"empty", nil},
		{"no SOI", []byte{0x00, 0x01, 0x02, 0x03}},
		{"SOI only", []byte{0xFF, 0xD8}},
		{"baseline JPEG", []byte{0xFF, 0xD8, 0xFF, 0xC0, 0x00, 0x02}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseLJPEG(tc.d); err == nil {
				t.Error("a malformed stream parsed without error")
			}
		})
	}
}

// isLJPEG is what tells the lossless variant from the ARW2 lossy one: both
// carry Compression 32767 and nothing else in the container distinguishes them.
func TestLosslessIsDetectedBySOI(t *testing.T) {
	if !isLJPEG([]byte{0xFF, 0xD8, 0xFF, 0xC3}) {
		t.Error("a lossless-JPEG strip was not recognised")
	}
	if isLJPEG([]byte{0x00, 0x11, 0x22}) {
		t.Error("an ARW2 lossy strip was mistaken for lossless")
	}
}

// The shared-prelude cache is an optimisation with a per-tile fallback, and the
// fallback must produce identical pixels.
//
// Every Sony frame here carries one Huffman table set across all its tiles —
// verified by hashing the prelude of all 280 and 247 of them — but the encoder
// picks those tables from content, so a frame with varying tiles is legal. If
// one ever arrives, tiles decoded against the wrong tables would look plausible
// and be wrong, so this proves the escape hatch works before it is needed.
func TestTiledFallbackMatchesSharedPrelude(t *testing.T) {
	arw := loadARWSample(t, "testdata/ljpeg/_DSC4922.ARW")

	fast, err := DecodeARW(arw)
	if err != nil {
		t.Fatalf("shared-prelude decode: %v", err)
	}

	forceTileParse = true
	defer func() { forceTileParse = false }()
	slow, err := DecodeARW(arw)
	if err != nil {
		t.Fatalf("per-tile decode: %v", err)
	}

	if len(fast.Pixels) != len(slow.Pixels) {
		t.Fatalf("%d samples vs %d", len(fast.Pixels), len(slow.Pixels))
	}
	for i := range fast.Pixels {
		if fast.Pixels[i] != slow.Pixels[i] {
			t.Fatalf("sample %d (%d,%d) differs: shared %d, per-tile %d",
				i, i%fast.Width, i/fast.Width, fast.Pixels[i], slow.Pixels[i])
		}
	}
	t.Logf("both paths agree on all %d samples", len(fast.Pixels))
}

func loadARWSample(tb testing.TB, path string) []byte {
	tb.Helper()
	d, err := os.ReadFile(path)
	if err != nil {
		tb.Skipf("no sample at %s", path)
	}
	return d
}
