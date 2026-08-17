package fuji

import "math/bits"

// A native Go decoder for lossless-compressed RAF, built bottom-up so each
// layer can be verified before the next depends on it.
//
// Layers, innermost first:
//
//	rafBits       MSB-first bit reader           <- implemented, tested
//	quantise      gradient -> one of 9 levels    <- implemented, tested
//	gradTable     adaptive per-bucket statistics <- implemented, tested
//	predict       spatial prediction             <- in raf_lossless_decode.go
//	decodeLine    the 6-row band                 <- in raf_lossless_decode.go
//	decodeBlock   one vertical stripe            <- in raf_lossless_decode.go
//
// The coder is adaptive, so a subtle error desynchronises the bitstream
// rather than shifting a value.

// rafBits reads MSB-first across byte boundaries, which is the order the
// bitstream uses — block 0 of a real frame opens 0x80 0x40 0x40 0x40 0x81,
// i.e. 1 / 0000000 01 / 000000 01 / 000000 01, a zero-run terminated by a 1.
type rafBits struct {
	data []byte
	pos  int    // next byte to load into acc
	acc  uint64 // pending bits, left-aligned: the next bit is the MSB
	n    uint   // how many bits of acc are valid
}

func newRAFBits(data []byte) *rafBits { return &rafBits{data: data} }

// fill tops the accumulator up to at least 57 bits, so any single read below —
// at most 41 zeros plus a 14-bit remainder — is satisfied without refilling
// mid-value.
func (b *rafBits) fill() {
	// Byte at a time on purpose. Replacing this with one unaligned 64-bit load
	// MEASURED 10% SLOWER (257ms vs 232ms, back to back): the load needs its
	// own bounds check, and when n is not a multiple of eight the surplus bits
	// have to be masked off, which costs more than the few byte loads it saves.
	// fill also returns immediately whenever the accumulator is still full, so
	// this loop is usually short.
	for b.n <= 56 && b.pos < len(b.data) {
		b.acc |= uint64(b.data[b.pos]) << (56 - b.n)
		b.n += 8
		b.pos++
	}
}

// exhausted reports that the reader has run past the end of its block. A
// truncated block must be detectable rather than silently reading zeroes: the
// coder is adaptive, so reading past the end corrupts every later sample.
func (b *rafBits) exhausted() bool { return b.pos >= len(b.data) && b.n == 0 }

// bit1 reads one bit.
func (b *rafBits) bit1() uint32 {
	if b.n == 0 {
		b.fill()
		if b.n == 0 {
			return 0
		}
	}
	v := uint32(b.acc >> 63)
	b.acc <<= 1
	b.n--
	return v
}

// bits reads k bits, most significant first.
func (b *rafBits) bits(k uint) uint32 {
	if k == 0 {
		return 0
	}
	b.fill()
	v := uint32(b.acc >> (64 - k))
	if b.n < k {
		// Past the end: the absent bits read as zero, which is what the
		// shift already produced.
		b.acc, b.n = 0, 0
		return v
	}
	b.acc <<= k
	b.n -= k
	return v
}

// zeros counts the zero bits before the next 1, consuming the 1. This is the
// unary prefix of the Golomb-Rice code.
//
// Counting is one LeadingZeros64 rather than a loop over bits: the run can be
// 41 long, and this was the single hottest thing in the decoder.
func (b *rafBits) zeros(limit int) int {
	total := 0
	for {
		b.fill()
		if b.n == 0 {
			return total
		}
		z := uint(bits.LeadingZeros64(b.acc))
		if z >= b.n {
			// Every valid bit left is a zero; count them and refill.
			total += int(b.n)
			b.acc, b.n = 0, 0
			if total > limit {
				return total
			}
			continue
		}
		total += int(z)
		b.acc <<= z + 1
		b.n -= z + 1
		return total
	}
}

// Gradient quantisation.
//
// A gradient is reduced to one of nine levels, -4..+4, by comparing its
// magnitude against four thresholds. The thresholds are NOT powers of two.
const (
	qLevels  = 9  // -4..+4
	qBuckets = 41 // after taking the absolute value of the combined index
)

// quantise maps a gradient to -4..+4 using the supplied ascending thresholds.
func quantise(g int32, t [4]int32) int32 {
	neg := g < 0
	if neg {
		g = -g
	}
	var lvl int32
	switch {
	case g < t[0]:
		lvl = 0
	case g < t[1]:
		lvl = 1
	case g < t[2]:
		lvl = 2
	case g < t[3]:
		lvl = 3
	default:
		lvl = 4
	}
	if neg {
		return -lvl
	}
	return lvl
}

// gradBucket combines two quantised levels into one context.
//
// The combined value is SIGNED, spanning -40..+40; its sign selects the
// direction the correction is applied in, and its magnitude indexes the 41
// adaptive slots. Collapsing to an unsigned 0..80 index would discard the
// direction and get every correction's sign wrong.
func gradBucket(a, b int32) (idx int, negative bool) {
	v := a*qLevels + b
	if v < 0 {
		return int(-v), true
	}
	return int(v), false
}

// gradTable is the adaptive statistic behind the variable-length codes: each
// bucket accumulates the magnitudes that landed in it and how many there were,
// and the code width is derived from their ratio.
type gradTable struct {
	sum   [qBuckets]int32
	count [qBuckets]int32
	// w caches width(bucket). It only changes when the bucket is observed, and
	// recomputing the shift loop on every sample cost 8% of the decoder.
	w [qBuckets]uint
	// rescaleAt halves both once count reaches it, so the statistic tracks
	// recent data instead of the whole block.
	rescaleAt int32
}

func newGradTable(initSum, rescaleAt int32) *gradTable {
	t := &gradTable{rescaleAt: rescaleAt}
	for i := range t.sum {
		t.sum[i] = initSum
		t.count[i] = 1
		t.recompute(i)
	}
	return t
}

// recompute refreshes a bucket's cached code width.
func (t *gradTable) recompute(bucket int) {
	t.w[bucket] = codeWidth(t.sum[bucket], t.count[bucket])
}

// codeWidth is the smallest w such that count<<w >= sum, capped at 15, or 0
// when count already covers sum.
//
// Derived from the bit lengths rather than doubling in a loop: this runs on
// every decoded sample, and the loop form was 10% of the decoder.
func codeWidth(sum, count int32) uint {
	if count >= sum || count <= 0 {
		return 0
	}
	w := uint(bits.Len32(uint32(sum)) - bits.Len32(uint32(count)))
	// Len32 gives the answer to within one either way; correct it.
	if w > 0 && count<<(w-1) >= sum {
		w--
	} else if count<<w < sum {
		w++
	}
	if w > 15 {
		w = 15
	}
	return w
}

// width is how many bits the fixed-length remainder occupies for this bucket:
// the number of doublings of count needed to reach sum. Zero when the average
// magnitude is at most one.
func (t *gradTable) width(bucket int) uint { return t.w[bucket] }

// observe folds a decoded magnitude into a bucket's statistic.
func (t *gradTable) observe(bucket int, magnitude int32) {
	if magnitude < 0 {
		magnitude = -magnitude
	}
	t.sum[bucket] += magnitude
	if t.count[bucket] == t.rescaleAt {
		t.sum[bucket] >>= 1
		t.count[bucket] >>= 1
	}
	t.count[bucket]++
	t.recompute(bucket)
}

// zigzag undoes the signed-magnitude mapping: the coder emits a non-negative
// code whose low bit carries the sign, so 0,1,2,3,4 -> 0,-1,1,-2,2. This is NOT
// two's complement, and treating it as such negates every odd value.
func zigzag(code int32) int32 {
	// Branchless, and the code is always non-negative here so the shift is the
	// same as the division: 0,1,2,3,4 -> 0,-1,1,-2,2.
	return (code >> 1) ^ -(code & 1)
}
