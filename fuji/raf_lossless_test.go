package fuji

import "testing"

// The bit reader must be MSB-first across byte boundaries. Pinned against the
// real opening bytes of block 0 in xt5-lossless.raf, which read as a zero-run
// terminated by a 1 — the unary prefix of the Golomb-Rice code.
func TestRAFBitsIsMSBFirst(t *testing.T) {
	b := newRAFBits([]byte{0x80, 0x40, 0x40, 0x40, 0x81})
	// 1000 0000 0100 0000 0100 0000 0100 0000 1000 0001
	// ^1        ^8 zeros then 1
	//              ^7 zeros then 1  ^7 zeros then 1
	//                                    ^6 zeros then 1
	for i, want := range []int{0, 8, 7, 7, 6} {
		if got := b.zeros(41); got != want {
			t.Fatalf("run %d: %d zeros, want %d", i, got, want)
		}
	}
}

func TestRAFBitsReadsAcrossByteBoundaries(t *testing.T) {
	b := newRAFBits([]byte{0xAB, 0xCD, 0xEF})
	if got := b.bits(4); got != 0xA {
		t.Errorf("first nibble %#x, want 0xA", got)
	}
	if got := b.bits(8); got != 0xBC { // spans a byte boundary
		t.Errorf("straddling byte %#x, want 0xBC", got)
	}
	if got := b.bits(12); got != 0xDEF {
		t.Errorf("remaining 12 bits %#x, want 0xDEF", got)
	}
}

// Running off the end must be visible. The coder is adaptive, so silently
// reading zeroes past a truncated block corrupts every later sample instead of
// failing where the damage is.
func TestRAFBitsExhaustion(t *testing.T) {
	b := newRAFBits([]byte{0x00})
	b.bits(8)
	if !b.exhausted() {
		t.Fatal("the reader does not report running out of data")
	}
	if got := b.zeros(41); got != 0 {
		t.Errorf("zeros past the end returned %d, want 0 rather than spinning", got)
	}
}

func TestQuantise(t *testing.T) {
	th := [4]int32{18, 67, 276, 1024}
	for _, tc := range []struct {
		g    int32
		want int32
	}{
		{0, 0}, {17, 0}, {18, 1}, {66, 1}, {67, 2}, {275, 2},
		{276, 3}, {1023, 3}, {1024, 4}, {9999, 4},
		{-17, 0}, {-18, -1}, {-276, -3}, {-9999, -4},
	} {
		if got := quantise(tc.g, th); got != tc.want {
			t.Errorf("quantise(%d) = %d, want %d", tc.g, got, tc.want)
		}
	}
}

// The combined index must keep its SIGN: it selects the direction the
// correction is applied in. Collapsing to an unsigned 0..80 would invert half
// the corrections while still producing a plausible-looking image.
func TestGradBucketKeepsSign(t *testing.T) {
	if idx, neg := gradBucket(4, 4); idx != 40 || neg {
		t.Errorf("gradBucket(4,4) = %d,%v want 40,false", idx, neg)
	}
	if idx, neg := gradBucket(-4, -4); idx != 40 || !neg {
		t.Errorf("gradBucket(-4,-4) = %d,%v want 40,true", idx, neg)
	}
	if idx, _ := gradBucket(0, 0); idx != 0 {
		t.Errorf("gradBucket(0,0) = %d, want 0", idx)
	}
	// Every combination must land inside the table.
	for a := int32(-4); a <= 4; a++ {
		for b := int32(-4); b <= 4; b++ {
			if idx, _ := gradBucket(a, b); idx < 0 || idx >= qBuckets {
				t.Fatalf("gradBucket(%d,%d) = %d, outside 0..%d", a, b, idx, qBuckets-1)
			}
		}
	}
}

func TestGradTableAdapts(t *testing.T) {
	tab := newGradTable(256, 64)
	if got := tab.width(0); got == 0 {
		t.Error("a fresh bucket asks for 0 bits, so nothing would be read")
	}
	// A run of small magnitudes must narrow the code.
	before := tab.width(0)
	for i := 0; i < 40; i++ {
		tab.observe(0, 0)
	}
	if after := tab.width(0); after >= before {
		t.Errorf("width went %d -> %d; small values should narrow it", before, after)
	}
	// Rescaling must keep the counter bounded.
	for i := 0; i < 500; i++ {
		tab.observe(1, 3)
	}
	if tab.count[1] > tab.rescaleAt {
		t.Errorf("count reached %d, past the rescale gate %d", tab.count[1], tab.rescaleAt)
	}
}

// Signed-magnitude, NOT two's complement. Getting this wrong negates every odd
// value, which looks like noise rather than an obvious failure.
func TestZigzag(t *testing.T) {
	for code, want := range map[int32]int32{0: 0, 1: -1, 2: 1, 3: -2, 4: 2, 5: -3, 6: 3} {
		if got := zigzag(code); got != want {
			t.Errorf("zigzag(%d) = %d, want %d", code, got, want)
		}
	}
}

// codeWidth replaced a doubling loop with bit arithmetic. It must agree with
// that loop everywhere, not just on the values a sample frame happens to hit.
func TestCodeWidthMatchesTheDoublingLoop(t *testing.T) {
	ref := func(sum, count int32) uint {
		var w uint
		if count < sum {
			for w <= 14 {
				count <<= 1
				w++
				if count >= sum {
					break
				}
			}
		}
		return w
	}
	for _, sum := range []int32{0, 1, 2, 3, 255, 256, 257, 1023, 4096, 16383, 65535, 1 << 20} {
		for _, count := range []int32{1, 2, 3, 7, 15, 16, 63, 64, 255, 4096, 65535} {
			if got, want := codeWidth(sum, count), ref(sum, count); got != want {
				t.Errorf("codeWidth(%d,%d) = %d, the loop gives %d", sum, count, got, want)
			}
		}
	}
}
