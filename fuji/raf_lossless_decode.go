package fuji

import (
	"fmt"
	"sync"

	"github.com/mikefsq/ptp"
)

// The lossless codec: prediction, the 6-row line band, and the block loop.
//
// # Line buffers
//
// A "line" is six rows — one X-Trans period. Samples are not decoded in raster
// order; they are collated into 18 per-colour buffers, five red, eight green,
// five blue, each lineWidth+2 wide. The two extra slots are the left and right
// edges, so a predictor at either end always has neighbours to read.
//
// Of the 18, only _R2.._R4, _G2.._G7 and _B2.._B4 are decoded each line; the
// _x0/_x1 pairs hold the tail of the PREVIOUS line, which is what makes
// prediction work across the band boundary.
//
// # Why the layout is transposed
//
// The buffers are stored slot-major, [slot][buffer], NOT buffer-major. Every
// read a predictor makes stays inside one colour group — R2 predicts from R1
// and R0, G3 from G2 and G1 — so storing a slot's 18 buffers contiguously puts
// a sample's four reads within about 76 bytes, one or two cache lines, instead
// of two rows a kilobyte apart. carry benefits for the same reason: the six
// values it moves per slot are adjacent.
//
//
//	 1 core    buffer-major 1226ms   transposed 1170ms   (4.6% faster)
//	 4 cores   buffer-major  328ms   transposed  320ms   (2.5% faster)
//	16 cores   inconclusive, lost in noise
//
//
// # The six passes cannot be reordered
//
// Pass 2's bits begin where pass 1's end and the gradient state carries
// forward, so interleaving them by column desynchronises the stream
// immediately. Layout is the only locality lever available.

// minLineWidth is the narrowest block this decoder accepts.
//
// Odd samples only begin once the even ones have run past 8, so they have
// decoded neighbours to predict from. A line narrower than that gate can never
// satisfy it, and the passes below would spin. Real blocks are 512 wide (768
// columns * 2/3), so this only ever rejects a malformed frame — but a frame
// arriving over USB must not be able to wedge the driver.
const minLineWidth = 16

const (
	lR0 = iota
	lR1
	lR2
	lR3
	lR4
	lG0
	lG1
	lG2
	lG3
	lG4
	lG5
	lG6
	lG7
	lB0
	lB1
	lB2
	lB3
	lB4
	lTotal
)

// Quantiser breakpoints. Deliberately not powers of two.
var qPoints = [4]int32{0x12, 0x43, 0x114, 0} // the fourth is maxValue, filled in

// The quantiser table covers only the range where the answer varies.
//
// The breakpoints are at +/-18, +/-67 and +/-276, so every difference of 276 or
// more quantises to +/-4 and needs no lookup. A table spanning +/-512 is
// therefore complete, and at 1 KB it stays in L1 — where a full +/-16383 table
// is 32 KB PER BLOCK, and with a goroutine per block that is 352 KB of
// randomly-accessed lookup evicting everything else. Measured identical on a
// large-cache desktop (196ms either way); the point is the Raspberry Pi this
// runs on, where L1 is 32-64 KB per core.
const (
	qTableHalf = 1 << 9 // covers -512..511
	qTableSize = 2 * qTableHalf
	qTableMask = qTableSize - 1
)

const (
	gradInitDivisor = 6    // max_diff = max(2, (totalValues+0x20) >> 6)
	gradRescaleAt   = 0x40 // halve the statistic once count reaches this
	qGradMult       = 9    // two levels combine as 9*a + b
)

// buildQTable precomputes the nine-level quantiser over the whole difference
// range, so the hot path is one array read rather than four comparisons.
func buildQTable(maxValue int32) *[qTableSize]int8 {
	t := new([qTableSize]int8)
	p := qPoints
	p[3] = maxValue
	for v := int32(-qTableHalf); v < qTableHalf; v++ {
		var q int8
		switch {
		case v <= -p[2]:
			q = -4
		case v <= -p[1]:
			q = -3
		case v <= -p[0]:
			q = -2
		case v < 0:
			q = -1
		case v == 0:
			q = 0
		case v < p[0]:
			q = 1
		case v < p[1]:
			q = 2
		case v < p[2]:
			q = 3
		default:
			q = 4
		}
		t[(v+qTableHalf)&qTableMask] = q
	}
	return t
}

// emitCell is one output pixel's precomputed source: which line buffer, and
// which slot within it.
type emitCell struct{ src int32 }

func abs32(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

// decodeLosslessBlock decodes one vertical stripe into dst, which is the whole
// frame; blockCol is the stripe's first column.
func decodeLosslessBlock(data []byte, h *compressedHeader, lineWidth, blockCol, width int,
	dst []uint16, stride int, pattern []ptp.CFAColor) error {
	if width <= 0 {
		return fmt.Errorf("fuji: block at column %d has no width", blockCol)
	}
	d := newLosslessDecoder(lineWidth, h.Bits, data)
	plan := d.planEmit(width, pattern)
	for line := 0; line < h.TotalLines; line++ {
		d.decodeLine()
		d.emit(dst, stride, blockCol, width, plan)
		d.carry()
		dst = dst[6*stride:]
	}
	return nil
}

// idx2 is the flat index of a slot within a buffer. Slot 0 is the left edge,
// 1..lineWidth the samples, lineWidth+1 the right edge.
func (d *losslessDecoder) at(slot, buf int) int { return slot*lTotal + buf }

type losslessDecoder struct {
	lineWidth   int
	maxValue    int32
	totalValues int32
	rawBits     uint
	maxBits     int
	qTable      *[qTableSize]int8
	lines       []uint16 // (lineWidth+2) * lTotal, slot-major
	even        [3]*gradTable
	odd         [3]*gradTable
	br          *rafBits
}

func newLosslessDecoder(lineWidth, bits int, data []byte) *losslessDecoder {
	maxValue := int32(1)<<uint(bits) - 1
	d := &losslessDecoder{
		lineWidth:   lineWidth,
		maxValue:    maxValue,
		totalValues: maxValue + 1,
		rawBits:     uint(bits),
		maxBits:     4 * bits,
		lines:       make([]uint16, (lineWidth+2)*lTotal),
		br:          newRAFBits(data),
	}
	d.qTable = buildQTable(maxValue)
	initSum := (d.totalValues + 0x20) >> gradInitDivisor
	if initSum < 2 {
		initSum = 2
	}
	for i := 0; i < 3; i++ {
		d.even[i] = newGradTable(initSum, gradRescaleAt)
		d.odd[i] = newGradTable(initSum, gradRescaleAt)
	}
	return d
}

func (d *losslessDecoder) q(v int32) int32 {
	if uint32(v+qTableHalf) < qTableSize {
		return int32(d.qTable[(v+qTableHalf)&qTableMask])
	}
	if v < 0 {
		return -4
	}
	return 4
}

func (d *losslessDecoder) readCode(g *gradTable, bucket int) int32 {
	sample := d.br.zeros(d.maxBits)
	if sample < d.maxBits-int(d.rawBits)-1 {
		w := g.width(bucket)
		return int32(d.br.bits(w)) + int32(sample)<<w
	}
	return int32(d.br.bits(d.rawBits)) + 1
}

func (d *losslessDecoder) finish(at int, pred, code int32, negCtx bool, g *gradTable, bucket int) {
	delta := zigzag(code)
	g.observe(bucket, delta)
	if negCtx {
		pred -= delta
	} else {
		pred += delta
	}
	if pred < 0 {
		pred += d.totalValues
	} else if pred > d.maxValue {
		pred -= d.totalValues
	}
	if pred < 0 {
		pred = 0
	} else if pred > d.maxValue {
		pred = d.maxValue
	}
	d.lines[at] = uint16(pred)
}

// sampleEven decodes an even position. base is the slot being written; the four
// neighbours sit within two slots of it.
func (d *losslessDecoder) sampleEven(buf, pos int, g *gradTable) {
	base := (pos+1)*lTotal + buf
	lines := d.lines
	Rf := int32(lines[base-2])
	Rb := int32(lines[base-1])
	Rc := int32(lines[base-1-lTotal])
	Rd := int32(lines[base-1+lTotal])

	var pred int32
	dcb, dfb, ddb := abs32(Rc-Rb), abs32(Rf-Rb), abs32(Rd-Rb)
	switch {
	case dcb > dfb && dcb > ddb:
		pred = Rf + Rd + 2*Rb
	case ddb > dcb && ddb > dfb:
		pred = Rf + Rc + 2*Rb
	default:
		pred = Rd + Rc + 2*Rb
	}
	pred >>= 2

	grad := qGradMult*d.q(Rb-Rf) + d.q(Rc-Rb)
	bucket, neg := int(grad), false
	if grad < 0 {
		bucket, neg = int(-grad), true
	}
	d.finish(base, pred, d.readCode(g, bucket), neg, g, bucket)
}

func (d *losslessDecoder) sampleOdd(buf, pos int, g *gradTable) {
	base := (pos+1)*lTotal + buf
	lines := d.lines
	Rc := int32(lines[base-1-lTotal])
	Rb := int32(lines[base-1])
	Rd := int32(lines[base-1+lTotal])
	Ra := int32(lines[base-lTotal])
	Rg := int32(lines[base+lTotal])

	var pred int32
	if (Rb > Rc && Rb > Rd) || (Rb < Rc && Rb < Rd) {
		pred = (Rg + Ra + 2*Rb) >> 2
	} else {
		pred = (Ra + Rg) >> 1
	}

	grad := qGradMult*d.q(Rb-Rc) + d.q(Rc-Ra)
	bucket, neg := int(grad), false
	if grad < 0 {
		bucket, neg = int(-grad), true
	}
	d.finish(base, pred, d.readCode(g, bucket), neg, g, bucket)
}

func (d *losslessDecoder) interpolateEven(buf, pos int) {
	base := (pos+1)*lTotal + buf
	lines := d.lines
	Rf := int32(lines[base-2])
	Rb := int32(lines[base-1])
	Rc := int32(lines[base-1-lTotal])
	Rd := int32(lines[base-1+lTotal])

	var pred int32
	dcb, dfb, ddb := abs32(Rc-Rb), abs32(Rf-Rb), abs32(Rd-Rb)
	switch {
	case dcb > dfb && dcb > ddb:
		pred = Rf + Rd + 2*Rb
	case ddb > dcb && ddb > dfb:
		pred = Rf + Rc + 2*Rb
	default:
		pred = Rd + Rc + 2*Rb
	}
	d.lines[base] = uint16(pred >> 2)
}

func (d *losslessDecoder) extend(from, to int) {
	lw := d.lineWidth
	for i := from; i <= to; i++ {
		d.lines[d.at(0, i)] = d.lines[d.at(1, i-1)]
		d.lines[d.at(lw+1, i)] = d.lines[d.at(lw, i-1)]
	}
}

func (d *losslessDecoder) extendRed()   { d.extend(lR2, lR4) }
func (d *losslessDecoder) extendGreen() { d.extend(lG2, lG7) }
func (d *losslessDecoder) extendBlue()  { d.extend(lB2, lB4) }

// carry moves each colour's last two buffers into its history slots. In this
// layout the six values are adjacent within a slot, so it walks memory once
// instead of doing six strided copies.
func (d *losslessDecoder) carry() {
	for slot := 0; slot < d.lineWidth+2; slot++ {
		b := slot * lTotal
		d.lines[b+lR0] = d.lines[b+lR3]
		d.lines[b+lR1] = d.lines[b+lR4]
		d.lines[b+lG0] = d.lines[b+lG6]
		d.lines[b+lG1] = d.lines[b+lG7]
		d.lines[b+lB0] = d.lines[b+lB3]
		d.lines[b+lB1] = d.lines[b+lB4]
	}
}

// planEmit precomputes the buffer-and-slot for every pixel of a 6-row band.
//
// The mapping depends only on the column and the mosaic, so it is identical for
// all 866 lines of a block. Computing it per pixel meant three integer
// divisions 44 million times per frame; computing it once costs 4608 entries.
func (d *losslessDecoder) planEmit(width int, pattern []ptp.CFAColor) []emitCell {
	plan := make([]emitCell, 6*width)
	for row := 0; row < 6; row++ {
		for col := 0; col < width; col++ {
			var buf int
			switch pattern[((row+1)%6)*6+(col+1)%6] {
			case ptp.CFARed:
				buf = lR2 + row>>1
			case ptp.CFABlue:
				buf = lB2 + row>>1
			default:
				buf = lG2 + row
			}
			idx := (((col * 2 / 3) & 0x7FFFFFFE) | ((col % 3) & 1)) + ((col % 3) >> 1)
			plan[row*width+col] = emitCell{src: int32(d.at(idx+1, buf))}
		}
	}
	return plan
}

func (d *losslessDecoder) emit(dst []uint16, stride, blockCol, width int, plan []emitCell) {
	lines := d.lines
	for row := 0; row < 6; row++ {
		out := dst[row*stride+blockCol:]
		p := plan[row*width : row*width+width]
		out = out[:len(p)]
		for col, c := range p {
			out[col] = lines[c.src]
		}
	}
}

// decodeLine decodes one 6-row band: six passes, each filling two buffers,
// even slots running ahead of odd ones, which start once the evens pass 8.
// Which slots are coded and which interpolated follows a positional rule that
// differs per pass — all of it derived against known plaintext.
func (d *losslessDecoder) decodeLine() {
	lw := d.lineWidth
	d.extendRed()
	d.extendGreen()
	d.extendBlue()

	rEven, rOdd, gEven, gOdd := 0, 1, 0, 1
	for gEven < lw || gOdd < lw {
		if gEven < lw {
			d.interpolateEven(lR2, rEven)
			rEven += 2
			d.sampleEven(lG2, gEven, d.even[0])
			gEven += 2
		}
		if gEven > 8 {
			d.sampleOdd(lR2, rOdd, d.odd[0])
			rOdd += 2
			d.sampleOdd(lG2, gOdd, d.odd[0])
			gOdd += 2
		}
	}
	d.extendRed()
	d.extendGreen()

	gEven, gOdd = 0, 1
	bEven, bOdd := 0, 1
	for gEven < lw || gOdd < lw {
		if gEven < lw {
			d.sampleEven(lG3, gEven, d.even[1])
			gEven += 2
			d.interpolateEven(lB2, bEven)
			bEven += 2
		}
		if gEven > 8 {
			d.sampleOdd(lG3, gOdd, d.odd[1])
			gOdd += 2
			d.sampleOdd(lB2, bOdd, d.odd[1])
			bOdd += 2
		}
	}
	d.extendGreen()
	d.extendBlue()

	rEven, rOdd, gEven, gOdd = 0, 1, 0, 1
	for gEven < lw || gOdd < lw {
		if gEven < lw {
			if rEven&3 != 0 {
				d.sampleEven(lR3, rEven, d.even[2])
			} else {
				d.interpolateEven(lR3, rEven)
			}
			rEven += 2
			d.interpolateEven(lG4, gEven)
			gEven += 2
		}
		if gEven > 8 {
			d.sampleOdd(lR3, rOdd, d.odd[2])
			rOdd += 2
			d.sampleOdd(lG4, gOdd, d.odd[2])
			gOdd += 2
		}
	}
	d.extendRed()
	d.extendGreen()

	gEven, gOdd = 0, 1
	bEven, bOdd = 0, 1
	for gEven < lw || gOdd < lw {
		if gEven < lw {
			d.sampleEven(lG5, gEven, d.even[0])
			gEven += 2
			if bEven&3 == 2 {
				d.interpolateEven(lB3, bEven)
			} else {
				d.sampleEven(lB3, bEven, d.even[0])
			}
			bEven += 2
		}
		if gEven > 8 {
			d.sampleOdd(lG5, gOdd, d.odd[0])
			gOdd += 2
			d.sampleOdd(lB3, bOdd, d.odd[0])
			bOdd += 2
		}
	}
	d.extendGreen()
	d.extendBlue()

	rEven, rOdd, gEven, gOdd = 0, 1, 0, 1
	for gEven < lw || gOdd < lw {
		if gEven < lw {
			if rEven&3 == 2 {
				d.interpolateEven(lR4, rEven)
			} else {
				d.sampleEven(lR4, rEven, d.even[1])
			}
			rEven += 2
			d.sampleEven(lG6, gEven, d.even[1])
			gEven += 2
		}
		if gEven > 8 {
			d.sampleOdd(lR4, rOdd, d.odd[1])
			rOdd += 2
			d.sampleOdd(lG6, gOdd, d.odd[1])
			gOdd += 2
		}
	}
	d.extendRed()
	d.extendGreen()

	gEven, gOdd = 0, 1
	bEven, bOdd = 0, 1
	for gEven < lw || gOdd < lw {
		if gEven < lw {
			d.interpolateEven(lG7, gEven)
			gEven += 2
			if bEven&3 != 0 {
				d.sampleEven(lB4, bEven, d.even[2])
			} else {
				d.interpolateEven(lB4, bEven)
			}
			bEven += 2
		}
		if gEven > 8 {
			d.sampleOdd(lG7, gOdd, d.odd[2])
			gOdd += 2
			d.sampleOdd(lB4, bOdd, d.odd[2])
			bOdd += 2
		}
	}
	d.extendGreen()
	d.extendBlue()
}

// decodeCompressed unpacks a compressed strip into the frame's samples.
//
// Blocks are independent — each owns a byte range and a column range, and
// shares no coder state — so they decode concurrently. That is not an
// optimisation bolted on afterwards; it is why the format is shaped this way.
//
// The native decoder handles the lossless mode. Anything else is refused: the
// lossy mode is a different coder, and running it through this one would
// produce a confident wrong answer.
func decodeCompressed(strip []byte, h *compressedHeader, out *ptp.CFA) ([]uint16, error) {
	if h.RawType != rawTypeXTrans {
		return nil, unsupported(h, fmt.Errorf("raw type %d is not X-Trans", h.RawType))
	}
	if h.BlockSize%3 != 0 {
		return nil, unsupported(h, fmt.Errorf("block size %d is not a multiple of 3", h.BlockSize))
	}

	lineWidth := h.BlockSize * 2 / 3
	if lineWidth < minLineWidth {
		return nil, unsupported(h, fmt.Errorf("block size %d gives a %d-sample line, below the %d minimum",
			h.BlockSize, lineWidth, minLineWidth))
	}
	px := make([]uint16, out.Width*out.Height)

	offset := h.DataOffset
	type job struct {
		data     []byte
		col, wid int
	}
	jobs := make([]job, 0, h.BlocksInRow)
	for i := 0; i < h.BlocksInRow; i++ {
		size := h.BlockSizes[i]
		if offset+size > len(strip) {
			return nil, fmt.Errorf("fuji: block %d (%d+%d) runs past the strip (%d)",
				i, offset, size, len(strip))
		}
		col := i * h.BlockSize
		wid := h.BlockSize
		if col+wid > out.Width {
			// The last block is partial: 11 blocks of 768 cover 8448 columns of
			// a 7872-wide frame, so the final one carries only 192.
			wid = out.Width - col
		}
		if wid > 0 {
			jobs = append(jobs, job{strip[offset : offset+size], col, wid})
		}
		offset += size
	}

	var wg sync.WaitGroup
	errs := make([]error, len(jobs))
	for i, j := range jobs {
		wg.Add(1)
		go func(i int, j job) {
			defer wg.Done()
			errs[i] = decodeLosslessBlock(j.data, h, lineWidth, j.col, j.wid,
				px, out.Width, out.Pattern)
		}(i, j)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return px, nil
}

// unsupported explains a compressed variant this decoder does not implement.
//
// The lossless mode decodes natively. What lands here is the lossy mode, or a
// Bayer GFX body, or a frame whose parameters fall outside what the native
// decoder handles — and a decoder written for one of those produces confident
// nonsense on the others, so it refuses rather than approximates.
func unsupported(h *compressedHeader, why error) error {
	return fmt.Errorf("%w: %v (%d-bit, %dx%d, %d blocks of %d columns, %d lines)",
		ErrCompressedRAF, why, h.Bits, h.Width, h.Height, h.BlocksInRow, h.BlockSize, h.TotalLines)
}
