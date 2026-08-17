package sony

import (
	"math/bits"
	"testing"
)

// destuff must keep a stuffed 0xFF, drop the 0x00, and stop at markers —
// getting any of these wrong shifts the bitstream and desynchronises the
// whole tile.
func TestDestuff(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   []byte
		want []byte
	}{
		{"plain", []byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03}},
		{"stuffed", []byte{0x01, 0xFF, 0x00, 0x02}, []byte{0x01, 0xFF, 0x02}},
		{"two stuffed", []byte{0xFF, 0x00, 0xFF, 0x00}, []byte{0xFF, 0xFF}},
		{"marker ends the scan", []byte{0x01, 0xFF, 0xD9, 0x02}, []byte{0x01}},
		{"trailing lone FF", []byte{0x01, 0x02, 0xFF}, []byte{0x01, 0x02}},
		{"empty", nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := destuff(tc.in, nil)
			if len(got) != len(tc.want)+destuffPad {
				t.Fatalf("got %d bytes, want %d data + %d padding", len(got), len(tc.want), destuffPad)
			}
			for i, b := range tc.want {
				if got[i] != b {
					t.Fatalf("byte %d is %#02x, want %#02x", i, got[i], b)
				}
			}
			for i := len(tc.want); i < len(got); i++ {
				if got[i] != 0 {
					t.Fatalf("padding byte %d is %#02x, want 0", i, got[i])
				}
			}
		})
	}
}

// The destuffed reader and the stuffing-aware one must produce the same bits
// from the same stream. The stream is chosen to hold stuffed bytes at awkward
// alignments.
func TestDBitsMatchesJBits(t *testing.T) {
	raw := []byte{0xA5, 0xFF, 0x00, 0x3C, 0xFF, 0x00, 0xFF, 0x00, 0x01, 0x80, 0x7E}
	j := &jbits{data: raw}
	d := &dbits{data: destuff(raw, nil)}
	for i := 0; i < (len(raw)-3)*8; i++ {
		if jb, db := j.bit(), d.bit(); jb != db {
			t.Fatalf("bit %d: jbits %d, dbits %d", i, jb, db)
		}
	}
}

// --- a tiny lossless-JPEG encoder, so the fast path has known plaintext ---
//
// The real-frame tests skip on a clean checkout because the samples are
// gitignored, and every real tile takes the fast path — which would leave it
// with no untagged coverage at all. This encoder produces the exact shape a
// Sony tile takes (four components, predictor 1, point transform 0, no
// restarts) from samples the test chose, so a decode is checked against known
// values, not just against another decoder.

type bitWriter struct {
	buf []byte
	acc uint32
	n   uint
}

func (w *bitWriter) put(v uint32, k uint) {
	w.acc = w.acc<<k | v&(1<<k-1)
	w.n += k
	for w.n >= 8 {
		w.n -= 8
		w.buf = append(w.buf, byte(w.acc>>w.n))
	}
}

// flush pads the final byte with 1s, which is what T.81 specifies.
func (w *bitWriter) flush() {
	if w.n > 0 {
		w.put((1<<(8-w.n))-1, 8-w.n)
	}
}

// encodeLJPEG builds a four-component predictor-1 SOF3 stream around samples,
// which is indexed [y*w+x][c]. One Huffman table serves all components: every
// sbits symbol 0..14 coded in four bits, which is valid (15 of 16 slots) and
// trivially canonical — symbol s IS code s.
func encodeLJPEG(samples [][4]int32, w, h, precision int) []byte {
	bw := &bitWriter{}
	start := int32(1) << (precision - 1)
	var above, left [4]int32
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 4; c++ {
				var pred int32
				switch {
				case x != 0:
					pred = left[c]
				case y != 0:
					pred = above[c]
				default:
					pred = start
				}
				v := samples[y*w+x][c]
				diff := v - pred
				left[c] = v
				if x == 0 {
					above[c] = v
				}
				var sbits uint
				mag := diff
				if diff < 0 {
					sbits = uint(bits.Len32(uint32(-diff)))
					mag = diff + 1<<sbits - 1
				} else if diff > 0 {
					sbits = uint(bits.Len32(uint32(diff)))
				}
				bw.put(uint32(sbits), 4)
				if sbits > 0 {
					bw.put(uint32(mag), sbits)
				}
			}
		}
	}
	bw.flush()

	var d []byte
	d = append(d, 0xFF, mSOI)
	// SOF3: precision, height, width, 4 components all 1x1 on table 0.
	d = append(d, 0xFF, mSOF3, 0, 20, byte(precision),
		byte(h>>8), byte(h), byte(w>>8), byte(w), 4)
	for c := 0; c < 4; c++ {
		d = append(d, byte(c), 0x11, 0)
	}
	// DHT: class 0 table 0, 15 codes of length four, values 0..14.
	d = append(d, 0xFF, mDHT, 0, 34, 0x00)
	for l := 1; l <= 16; l++ {
		if l == 4 {
			d = append(d, 15)
		} else {
			d = append(d, 0)
		}
	}
	for s := byte(0); s < 15; s++ {
		d = append(d, s)
	}
	// SOS: 4 components on table 0, predictor 1, point transform 0.
	d = append(d, 0xFF, mSOS, 0, 14)
	d = append(d, 4)
	for c := 0; c < 4; c++ {
		d = append(d, byte(c), 0x00)
	}
	d = append(d, 1, 0, 0)
	// The entropy segment, stuffed.
	for _, b := range bw.buf {
		d = append(d, b)
		if b == 0xFF {
			d = append(d, 0x00)
		}
	}
	return append(d, 0xFF, mEOI)
}

// The fast path against known plaintext: every decoded sample must be the one
// that was encoded, at the 2x2-scattered position, and the generic path must
// say exactly the same thing. This is the coverage a clean checkout gets.
func TestFastPathDecodesKnownPlaintext(t *testing.T) {
	const w, h, precision = 32, 32, 14
	// A deterministic full-range walk: large differences force long magnitude
	// codes and 0xFF bytes into the stream, so destuffing is exercised too.
	rng := uint32(1)
	samples := make([][4]int32, w*h)
	for i := range samples {
		for c := 0; c < 4; c++ {
			rng = rng*1664525 + 1013904223
			samples[i][c] = int32(rng>>18) & 0x3FFF
		}
	}
	stream := encodeLJPEG(samples, w, h, precision)

	stuffed := 0
	for i := 0; i+1 < len(stream); i++ {
		if stream[i] == 0xFF && stream[i+1] == 0x00 {
			stuffed++
		}
	}
	if stuffed == 0 {
		t.Fatal("the fixture has no stuffed bytes, so it does not exercise destuffing")
	}

	f, err := parseLJPEG(stream)
	if err != nil {
		t.Fatalf("parseLJPEG: %v", err)
	}
	if f.predictor != 1 || f.pointTrans != 0 || f.restart != 0 || len(f.comps) != 4 {
		t.Fatalf("the fixture parsed as predictor %d, pt %d, restart %d, %d comps — not the Sony shape",
			f.predictor, f.pointTrans, f.restart, len(f.comps))
	}

	decode := func() []uint16 {
		dst := make([]uint16, 2*w*2*h)
		if err := decodeLJPEGInto(f, dst, 2*w, 0, 0, nil); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return dst
	}
	fast := decode()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for c := 0; c < 4; c++ {
				row := 2*y + c>>1
				col := 2*x + c&1
				if got, want := fast[row*2*w+col], uint16(samples[y*w+x][c]); got != want {
					t.Fatalf("sample (%d,%d) comp %d: decoded %d, encoded %d", x, y, c, got, want)
				}
			}
		}
	}

	forceGenericLJPEG = true
	defer func() { forceGenericLJPEG = false }()
	generic := decode()
	for i := range fast {
		if fast[i] != generic[i] {
			t.Fatalf("sample %d: fast %d, generic %d", i, fast[i], generic[i])
		}
	}
}

// Both decode paths on the real frames. Every real tile qualifies for the
// fast path, so without this the generic loop would never see a Sony stream
// again — and the fast path is only trustworthy while the two agree.
func TestFastPathMatchesGeneric(t *testing.T) {
	for _, path := range []string{
		"testdata/ljpeg/_DSC4922.ARW",
		"testdata/ljpeg/_DSC2429.ARW",
		"testdata/ljpeg/Bias.ARW",
	} {
		t.Run(path, func(t *testing.T) {
			arw := loadARWSample(t, path)
			fast, err := DecodeARW(arw)
			if err != nil {
				t.Fatalf("fast decode: %v", err)
			}
			forceGenericLJPEG = true
			defer func() { forceGenericLJPEG = false }()
			generic, err := DecodeARW(arw)
			if err != nil {
				t.Fatalf("generic decode: %v", err)
			}
			if len(fast.Pixels) != len(generic.Pixels) {
				t.Fatalf("%d samples vs %d", len(fast.Pixels), len(generic.Pixels))
			}
			for i := range fast.Pixels {
				if fast.Pixels[i] != generic.Pixels[i] {
					t.Fatalf("sample %d (%d,%d) differs: fast %d, generic %d",
						i, i%fast.Width, i/fast.Width, fast.Pixels[i], generic.Pixels[i])
				}
			}
			t.Logf("both paths agree on all %d samples", len(fast.Pixels))
		})
	}
}
