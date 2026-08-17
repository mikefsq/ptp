package sony

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

// Lossless JPEG (ITU-T T.81 Annex H), the codec behind Sony's lossless ARW.
//
// This is NOT a Sony invention. libraw's sony_ljpeg_load_raw is a dozen lines
// that call the same generic ljpeg_start/ljpeg_row/ljpeg_end used for Canon
// CR2, Adobe DNG, Kodak and Nikon — Sony reused a published standard where
// Fujifilm invented a codec. That is why this file is written from the spec
// rather than reverse engineered: every parameter that Fujifilm leaves implied
// and adaptive is carried explicitly here, in the frame and scan headers.
//
// Only the lossless mode is implemented, SOF3, which is all a RAW file uses.
// There is no DCT, no quantisation and no colour transform: each sample is
// predicted from its already-decoded neighbours and the difference is Huffman
// coded.

var (
	errNotLJPEG   = errors.New("sony: not a lossless JPEG stream")
	errBadHuffman = errors.New("sony: a Huffman code ran past 16 bits")
)

// JPEG markers, only the ones a lossless stream uses.
const (
	mSOI  = 0xD8
	mEOI  = 0xD9
	mSOF3 = 0xC3 // lossless, Huffman
	mDHT  = 0xC4
	mSOS  = 0xDA
	mDRI  = 0xDD
	mRST0 = 0xD0
	mRST7 = 0xD7
)

// ljpegComp is one component of the frame.
type ljpegComp struct {
	id      int
	h, v    int // sampling factors; RAW uses 1x1 throughout
	tableID int // which Huffman table the scan assigns it
}

// ljpegFrame is a decoded SOF3 plus the scan parameters.
type ljpegFrame struct {
	precision int // P, bits per sample
	height    int // Y, lines
	width     int // X, samples per line PER COMPONENT
	comps     []ljpegComp

	predictor  int // Ss, selects one of the seven predictors in Table H.1
	pointTrans int // Al, a right shift applied to every sample
	restart    int // DRI interval in MCUs, 0 for none

	huff [4]*huffTable
	scan []byte // the entropy-coded segment
}

// huffLookupBits is the width of the flat decode table.
//
// Nine bits resolves the overwhelming majority of codes in one indexed read.
// Longer codes fall back to the canonical length-by-length walk, which is rare
// enough not to matter — and a wider table would cost cache for no gain.
const huffLookupBits = 9

// huffTable is a canonical JPEG Huffman table with a flat lookup in front.
//
// The spec's decoder (T.81 Annex F) compares against a per-length maxcode, one
// bit at a time, up to sixteen times per symbol. That was 30% of this decoder
// and another 20% went into the per-bit calls it made. The flat table turns the
// common case into one shift and one array read.
type huffTable struct {
	// fast maps the next huffLookupBits of the stream to (length<<8 | value).
	// Zero means the code is longer than the table and needs the slow path;
	// a real entry always has length >= 1, so it can never be zero.
	fast [1 << huffLookupBits]uint16

	vals    []byte
	mincode [17]int32
	maxcode [18]int32
	valptr  [17]int
}

func newHuffTable(counts [17]int, vals []byte) *huffTable {
	t := &huffTable{vals: vals}
	code, k := int32(0), 0
	for l := 1; l <= 16; l++ {
		t.valptr[l] = k
		t.mincode[l] = code
		// Fill every table slot whose leading l bits are this code.
		if l <= huffLookupBits {
			for i := 0; i < counts[l]; i++ {
				c := int(code) + i
				base := c << (huffLookupBits - l)
				entry := uint16(l)<<8 | uint16(vals[k+i])
				for j := 0; j < 1<<(huffLookupBits-l); j++ {
					t.fast[base+j] = entry
				}
			}
		}
		code += int32(counts[l])
		k += counts[l]
		t.maxcode[l] = code - 1
		if counts[l] == 0 {
			t.maxcode[l] = -1
		}
		code <<= 1
	}
	t.maxcode[17] = 0x7FFFFFFF
	return t
}

func (t *huffTable) decode(b *jbits) (int, error) {
	b.fill()
	if e := t.fast[b.acc>>(64-huffLookupBits)]; e != 0 {
		n := uint(e >> 8)
		b.acc <<= n
		b.n -= n
		return int(e & 0xFF), nil
	}
	return t.decodeLong(b)
}

// decodeLong walks the canonical table for codes too long for the flat one.
func (t *huffTable) decodeLong(b *jbits) (int, error) {
	code := int32(0)
	for l := 1; l <= 16; l++ {
		code = code<<1 | int32(b.bit())
		if t.maxcode[l] >= 0 && code <= t.maxcode[l] {
			i := t.valptr[l] + int(code-t.mincode[l])
			if i >= len(t.vals) {
				return 0, errBadHuffman
			}
			return int(t.vals[i]), nil
		}
	}
	return 0, errBadHuffman
}

// jbits reads the entropy-coded segment MSB-first, undoing JPEG's byte
// stuffing: a 0xFF in the data is written as 0xFF 0x00, and a 0xFF followed by
// anything else is a marker, which ends the segment.
//
// Bits are held left-aligned in a 64-bit accumulator so a caller can take
// several at once. Handing them out one at a time cost 20% of the decoder,
// because every Huffman symbol wanted up to sixteen.
type jbits struct {
	data   []byte
	pos    int
	acc    uint64 // pending bits, the next one is the MSB
	n      uint   // how many bits of acc are valid
	marker byte   // non-zero once a marker has been reached
}

// fill tops the accumulator up so any single read below — at most 16 bits of
// Huffman code plus 16 of magnitude — is satisfied without refilling midway.
func (b *jbits) fill() {
	for b.n <= 56 {
		if b.pos >= len(b.data) {
			return
		}
		c := b.data[b.pos]
		if c == 0xFF {
			if b.pos+1 >= len(b.data) {
				return
			}
			if nx := b.data[b.pos+1]; nx != 0 {
				// A marker ends the scan. Leave it in place and stop; the
				// accumulator's low bits are already zero, so a read past here
				// yields zeroes rather than the marker's bytes.
				b.marker = nx
				return
			}
			b.pos++ // stuffed: consume the 0x00, keep the 0xFF
		}
		b.pos++
		b.acc |= uint64(c) << (56 - b.n)
		b.n += 8
	}
}

// bit reads one bit. Kept for the long-code path and for tests; the hot paths
// take several bits at a time.
func (b *jbits) bit() uint32 {
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

// receive reads s bits as an unsigned value.
func (b *jbits) receive(s int) int32 {
	if s == 0 {
		return 0
	}
	b.fill()
	v := int32(b.acc >> (64 - uint(s)))
	if uint(s) > b.n {
		b.acc, b.n = 0, 0
		return v
	}
	b.acc <<= uint(s)
	b.n -= uint(s)
	return v
}

// extend turns the magnitude-coded difference into a signed one (T.81 F.2.2.1).
func extend(v int32, s int) int32 {
	if s == 0 {
		return 0
	}
	if v < 1<<(s-1) {
		return v - (1 << s) + 1
	}
	return v
}

// restart resets the bit reader at a restart marker: the coder realigns to a
// byte boundary and the predictors start over.
func (b *jbits) restart() bool {
	b.acc, b.n = 0, 0
	// Skip to just past the RSTn marker.
	for b.pos+1 < len(b.data) {
		if b.data[b.pos] == 0xFF {
			m := b.data[b.pos+1]
			if m >= mRST0 && m <= mRST7 {
				b.pos += 2
				b.marker = 0
				return true
			}
			if m != 0 {
				b.marker = m
				return false
			}
		}
		b.pos++
	}
	return false
}

// parseLJPEG reads the headers up to and including SOS, leaving scan pointing
// at the entropy-coded data.
func parseLJPEG(d []byte) (*ljpegFrame, error) {
	if len(d) < 4 || d[0] != 0xFF || d[1] != mSOI {
		return nil, errNotLJPEG
	}
	f := &ljpegFrame{}
	p := 2
	for p+4 <= len(d) {
		if d[p] != 0xFF {
			return nil, fmt.Errorf("sony: expected a JPEG marker at %d, found %#02x", p, d[p])
		}
		marker := d[p+1]
		p += 2
		if marker == mEOI {
			break
		}
		if p+2 > len(d) {
			return nil, errNotLJPEG
		}
		seglen := int(binary.BigEndian.Uint16(d[p:]))
		if seglen < 2 || p+seglen > len(d) {
			return nil, fmt.Errorf("sony: JPEG segment %#02x has a bad length %d", marker, seglen)
		}
		seg := d[p+2 : p+seglen]

		switch marker {
		case mSOF3:
			if len(seg) < 6 {
				return nil, errNotLJPEG
			}
			f.precision = int(seg[0])
			f.height = int(binary.BigEndian.Uint16(seg[1:]))
			f.width = int(binary.BigEndian.Uint16(seg[3:]))
			n := int(seg[5])
			if len(seg) < 6+3*n {
				return nil, errNotLJPEG
			}
			for i := 0; i < n; i++ {
				c := seg[6+3*i:]
				f.comps = append(f.comps, ljpegComp{
					id: int(c[0]), h: int(c[1] >> 4), v: int(c[1] & 0xF),
				})
			}
		case mDHT:
			for q := 0; q+17 <= len(seg); {
				tc, th := seg[q]>>4, seg[q]&0xF
				var counts [17]int
				total := 0
				for i := 1; i <= 16; i++ {
					counts[i] = int(seg[q+i])
					total += counts[i]
				}
				if q+17+total > len(seg) {
					return nil, errNotLJPEG
				}
				vals := append([]byte(nil), seg[q+17:q+17+total]...)
				if tc != 0 {
					return nil, errors.New("sony: an AC Huffman table in a lossless stream")
				}
				if th > 3 {
					return nil, errNotLJPEG
				}
				f.huff[th] = newHuffTable(counts, vals)
				q += 17 + total
			}
		case mDRI:
			if len(seg) >= 2 {
				f.restart = int(binary.BigEndian.Uint16(seg))
			}
		case mSOS:
			if len(seg) < 1 {
				return nil, errNotLJPEG
			}
			ns := int(seg[0])
			if len(seg) < 1+2*ns+3 {
				return nil, errNotLJPEG
			}
			for i := 0; i < ns; i++ {
				cs, tt := int(seg[1+2*i]), int(seg[2+2*i])
				for j := range f.comps {
					if f.comps[j].id == cs {
						f.comps[j].tableID = tt >> 4
					}
				}
			}
			tail := seg[1+2*ns:]
			f.predictor = int(tail[0])
			f.pointTrans = int(tail[2] & 0xF)
			f.scan = d[p+seglen:]
			return f, f.validate()
		}
		p += seglen
	}
	return nil, errors.New("sony: the lossless JPEG stream has no scan")
}

func (f *ljpegFrame) validate() error {
	if f.width <= 0 || f.height <= 0 || len(f.comps) == 0 {
		return errNotLJPEG
	}
	if f.precision < 8 || f.precision > 16 {
		return fmt.Errorf("sony: %d-bit lossless JPEG is out of range", f.precision)
	}
	if f.predictor < 1 || f.predictor > 7 {
		return fmt.Errorf("sony: predictor %d is not one of the seven in T.81 Table H.1", f.predictor)
	}
	for _, c := range f.comps {
		if c.h != 1 || c.v != 1 {
			return fmt.Errorf("sony: component %d is subsampled %dx%d; RAW is always 1x1",
				c.id, c.h, c.v)
		}
		if f.huff[c.tableID] == nil {
			return fmt.Errorf("sony: component %d names Huffman table %d, which was never defined",
				c.id, c.tableID)
		}
	}
	return nil
}

// tileScratch is the per-tile working state: the generic path's two predictor
// line buffers and resolved component tables, and the fast path's destuffed
// scan. A frame has hundreds of tiles and they are all the same shape, so a
// worker allocates this once and reuses it.
type tileScratch struct {
	prev, cur []int32
	tables    []*huffTable
	destuffed []byte
}

func (t *tileScratch) reset(outW, n int) {
	if cap(t.prev) < outW {
		t.prev = make([]int32, outW)
		t.cur = make([]int32, outW)
	}
	t.prev, t.cur = t.prev[:outW], t.cur[:outW]
	clear(t.prev)
	clear(t.cur)
	if cap(t.tables) < n {
		t.tables = make([]*huffTable, n)
	}
	t.tables = t.tables[:n]
}

// decodeLJPEGInto decodes a four-component tile DIRECTLY onto the Bayer plane.
//
// The obvious shape — decode to an interleaved buffer, scatter that into a tile
// buffer, copy the tile into the frame — writes every sample three times and
// allocates two 512 KB buffers per tile. The scatter is a fixed function of the
// component index, so it can be applied as the samples are produced:
//
//	component 0 -> (y*2,   x*2)
//	component 1 -> (y*2,   x*2+1)
//	component 2 -> (y*2+1, x*2)
//	component 3 -> (y*2+1, x*2+1)
//
// Only the predictor's two line buffers remain, and those are 4 KB each rather
// than per-tile allocations. dst is the whole frame; tx and ty are the tile's
// origin in it, and a tile overhanging the right or bottom edge is clipped.
//
// A tile in the shape every Sony body writes — predictor 1, no point
// transform, no restarts, wholly inside the frame — takes the specialised path
// in ljpeg_fast.go instead; everything else runs the general loop below, which
// implements the full standard and is the reference the fast path is held to.
func decodeLJPEGInto(f *ljpegFrame, dst []uint16, stride, tx, ty int, sc *tileScratch) error {
	n := len(f.comps)
	if n != 4 {
		return fmt.Errorf("sony: a lossless ARW tile should carry 4 components, this has %d", n)
	}
	if sc == nil {
		sc = &tileScratch{}
	}
	if height := len(dst) / stride; !forceGenericLJPEG &&
		f.predictor == 1 && f.pointTrans == 0 && f.restart == 0 &&
		tx+2*f.width <= stride && ty+2*f.height <= height {
		sc.destuffed = destuff(f.scan, sc.destuffed)
		return decodeLJPEGFast(f, sc.destuffed, dst, stride, tx, ty)
	}
	outW := f.width * n

	b := &jbits{data: f.scan}
	sc.reset(outW, n)
	prev, cur := sc.prev, sc.cur

	tables := sc.tables
	for c := range f.comps {
		tables[c] = f.huff[f.comps[c].tableID]
	}
	pt := uint(f.pointTrans)
	sel := f.predictor
	height := len(dst) / stride

	// Out-of-range samples are ORed together and checked once at the end.
	// libraw branches per sample; accumulating costs 1.2% and a branch cost
	// more when measured. A sample past the declared precision means the
	// stream desynchronised, and saying so beats handing back plausible-looking
	// wrong pixels — which is the one failure calibration cannot absorb.
	var seen int32
	limit := int32(1)<<uint(f.precision) - 1

	start := int32(1) << (f.precision - f.pointTrans - 1)
	mcu := 0
	atRestart := false

	for y := 0; y < f.height; y++ {
		// The two sensor rows this decoded line lands on. Either may fall off
		// the bottom of a padded frame.
		y0, y1 := ty+y*2, ty+y*2+1
		var row0, row1 []uint16
		if y0 < height {
			row0 = dst[y0*stride : (y0+1)*stride]
		}
		if y1 < height {
			row1 = dst[y1*stride : (y1+1)*stride]
		}

		col := 0
		for x := 0; x < f.width; x++ {
			if f.restart > 0 && mcu > 0 && mcu%f.restart == 0 {
				if !b.restart() {
					return errors.New("sony: a restart marker was expected and not found")
				}
				// Reset the FIRST-COLUMN predictors only, leaving the line
				// buffers alone. That is what libraw does, and it is the
				// behaviour proven against real files; zeroing the buffers as
				// well — which this used to do — would wreck every later
				// sample in the interval.
				//
				// Untestable either way: every frame seen has restart == 0.
				// Matching the reference beats keeping a guess.
				atRestart = true
			}
			x0, x1 := tx+x*2, tx+x*2+1
			for c := 0; c < n; c++ {
				b.fill()

				t := tables[c]
				var sbits int
				if e := t.fast[b.acc>>(64-huffLookupBits)]; e != 0 {
					k := uint(e >> 8)
					b.acc <<= k
					b.n -= k
					sbits = int(e & 0xFF)
				} else {
					var err error
					if sbits, err = t.decodeLong(b); err != nil {
						return err
					}
				}

				// A symbol of 16 is an ESCAPE, not a 16-bit magnitude: the
				// difference is -32768 and no magnitude bits follow. Reading
				// sixteen bits here instead would consume the next symbol's
				// code and desynchronise the rest of the tile.
				//
				// Neither sample frame exercises it, so this comes from
				// libraw's ljpeg_diff rather than from measurement. libraw
				// conditions it on DNG version; an ARW is not a DNG, so the
				// escape always applies here.
				var diff int32
				switch {
				case sbits == 16:
					diff = -32768
				case sbits > 0:
					k := uint(sbits)
					diff = int32(b.acc >> (64 - k))
					b.acc <<= k
					b.n -= k
					if diff < 1<<(k-1) {
						diff -= 1<<k - 1
					}
				}

				// Ordered by frequency: almost every sample is interior, so
				// that case is tested first. Measured at 1.7% — small, but the
				// only one of four profile-guided changes here that helped.
				// Guarding fill() with "if b.n < 32" cost 9%: its loop already
				// short-circuits when the accumulator is full, so the guard was
				// a branch in front of a free check.
				var pred int32
				switch {
				case x != 0 && y != 0 && !atRestart:
					pred = predictLJPEG(sel, cur[col-n], prev[col], prev[col-n])
				case atRestart || (x == 0 && y == 0):
					pred = start
				case x == 0:
					pred = prev[col]
				default:
					pred = cur[col-n]
				}
				v := pred + diff
				cur[col] = v
				seen |= v
				col++

				// Scatter straight to the sensor grid.
				tr, tc := row0, x0
				if c >= 2 {
					tr = row1
				}
				if c&1 != 0 {
					tc = x1
				}
				if tr != nil && tc < stride {
					tr[tc] = uint16(v << pt)
				}
			}
			mcu++
			atRestart = false
		}
		prev, cur = cur, prev
	}
	if seen < 0 || seen > limit {
		return fmt.Errorf("sony: a decoded sample fell outside %d bits, so the stream desynchronised",
			f.precision)
	}
	return nil
}

// decodeLJPEG expands a scan into one interleaved plane: with N components a
// decoded line is width*N samples, component c of MCU m landing at column
// m*N+c. Used for streams that are not a Sony Bayer tile — the verification
// fixture from another vendor.
func decodeLJPEG(f *ljpegFrame) ([]uint16, int, int, error) {
	n := len(f.comps)
	outW := f.width * n
	out := make([]uint16, outW*f.height)

	b := &jbits{data: f.scan}
	prev := make([]int32, outW)
	cur := make([]int32, outW)

	tables := make([]*huffTable, n)
	for c := range f.comps {
		tables[c] = f.huff[f.comps[c].tableID]
	}
	pt := uint(f.pointTrans)
	sel := f.predictor
	start := int32(1) << (f.precision - f.pointTrans - 1)
	mcu := 0
	atRestart := false

	for y := 0; y < f.height; y++ {
		row := out[y*outW : (y+1)*outW]
		col := 0
		for x := 0; x < f.width; x++ {
			if f.restart > 0 && mcu > 0 && mcu%f.restart == 0 {
				if !b.restart() {
					return nil, 0, 0, errors.New("sony: a restart marker was expected and not found")
				}
				// See decodeLJPEGInto: first-column predictors only.
				atRestart = true
			}
			for c := 0; c < n; c++ {
				b.fill()
				t := tables[c]
				var sbits int
				if e := t.fast[b.acc>>(64-huffLookupBits)]; e != 0 {
					k := uint(e >> 8)
					b.acc <<= k
					b.n -= k
					sbits = int(e & 0xFF)
				} else {
					var err error
					if sbits, err = t.decodeLong(b); err != nil {
						return nil, 0, 0, err
					}
				}
				// A symbol of 16 is an ESCAPE, not a 16-bit magnitude: the
				// difference is -32768 and no magnitude bits follow. Reading
				// sixteen bits here instead would consume the next symbol's
				// code and desynchronise the rest of the tile.
				//
				// Neither sample frame exercises it, so this comes from
				// libraw's ljpeg_diff rather than from measurement. libraw
				// conditions it on DNG version; an ARW is not a DNG, so the
				// escape always applies here.
				var diff int32
				switch {
				case sbits == 16:
					diff = -32768
				case sbits > 0:
					k := uint(sbits)
					diff = int32(b.acc >> (64 - k))
					b.acc <<= k
					b.n -= k
					if diff < 1<<(k-1) {
						diff -= 1<<k - 1
					}
				}
				var pred int32
				switch {
				case x != 0 && y != 0 && !atRestart:
					pred = predictLJPEG(sel, cur[col-n], prev[col], prev[col-n])
				case atRestart || (x == 0 && y == 0):
					pred = start
				case x == 0:
					pred = prev[col]
				default:
					pred = cur[col-n]
				}
				v := pred + diff
				cur[col] = v
				row[col] = uint16(v << pt)
				col++
			}
			mcu++
			atRestart = false
		}
		prev, cur = cur, prev
	}
	return out, outW, f.height, nil
}

// predictLJPEG is Table H.1: Ra is the sample to the left, Rb the one above,
// Rc the one above-left.
func predictLJPEG(sel int, ra, rb, rc int32) int32 {
	switch sel {
	case 1:
		return ra
	case 2:
		return rb
	case 3:
		return rc
	case 4:
		return ra + rb - rc
	case 5:
		return ra + (rb-rc)>>1
	case 6:
		return rb + (ra-rc)>>1
	case 7:
		return (ra + rb) >> 1
	}
	return 0
}

// DebugDecodeLJPEG exposes the lossless-JPEG decoder for verification against
// files from other vendors — the codec is a standard, so a Canon CR2 or an
// Adobe DNG exercises exactly the same path a Sony lossless ARW will.
func DebugDecodeLJPEG(d []byte) ([]uint16, int, int, string, error) {
	f, err := parseLJPEG(d)
	if err != nil {
		return nil, 0, 0, "", err
	}
	info := fmt.Sprintf("SOF3 %d-bit %dx%d, %d components, predictor %d, point transform %d, restart %d",
		f.precision, f.width, f.height, len(f.comps), f.predictor, f.pointTrans, f.restart)
	px, w, h, err := decodeLJPEG(f)
	return px, w, h, info, err
}

// decodeLosslessARW expands a single-tile lossless strip onto the Bayer plane.
//
// See decodeLJPEGInto for the component-to-position mapping and how it was
// established.
func decodeLosslessARW(strip []byte, width, height int, dst []uint16) ([]uint16, error) {
	f, err := parseLJPEG(strip)
	if err != nil {
		return nil, err
	}
	if f.width*2 != width || f.height*2 != height {
		return nil, fmt.Errorf("sony: the lossless JPEG is %dx%d, which does not tile %dx%d",
			f.width, f.height, width, height)
	}
	// The geometry check above guarantees the single tile covers the frame
	// exactly, so every sample is written or an error comes back.
	out := reuseAll(dst, width*height)
	if err := decodeLJPEGInto(f, out, width, 0, 0, nil); err != nil {
		return nil, err
	}
	return out, nil
}

// isLJPEG reports whether a strip opens with a JPEG SOI, which is what tells
// the lossless variant apart from the ARW2 lossy one: both carry Compression
// 32767, and nothing else in the container distinguishes them.
func isLJPEG(strip []byte) bool {
	return len(strip) >= 2 && strip[0] == 0xFF && strip[1] == mSOI
}

// decodeTiledLossless assembles a frame from its lossless-JPEG tiles.
//
// The A7R V and A7R VI split the readout into 512x512 tiles, each a separate
// four-component stream covering a 256x256 MCU grid. Tiles run left to right,
// top to bottom, and the readout is padded to whole tiles in both axes — which
// is why these sensors report dimensions divisible by 512.
//
// Tiles are independent, so they decode concurrently.
func decodeTiledLossless(data []byte, offs, lens []int, tileW, tileH, width, height int, dst []uint16) ([]uint16, error) {
	if len(offs) != len(lens) {
		return nil, fmt.Errorf("sony: %d tile offsets but %d lengths", len(offs), len(lens))
	}
	across := (width + tileW - 1) / tileW
	down := (height + tileH - 1) / tileH
	if want := across * down; want != len(offs) {
		return nil, fmt.Errorf("sony: %d tiles for a %dx%d frame of %dx%d tiles, want %d",
			len(offs), width, height, tileW, tileH, want)
	}

	// Tiles cover the whole padded frame — the count is checked above, a tile
	// of the wrong geometry is refused, and any tile that fails to decode
	// fails the frame — so every sample of a returned buffer has been
	// written and the zeroing reuse() does is skipped. It used to run here
	// on the grounds of being cheap next to page faults; that is true on
	// first use and false in the pooled steady state, where it is a
	// 147 MB memset per A7R VI frame.
	out := reuseAll(dst, width*height)
	// Bounds-check every tile before any worker slices into data.
	for i := range offs {
		if offs[i] < 0 || lens[i] < 0 || offs[i]+lens[i] > len(data) {
			return nil, fmt.Errorf("sony: tile %d (%d+%d) runs past the file (%d)",
				i, offs[i], lens[i], len(data))
		}
	}

	// forceTileParse makes every tile take the per-tile path, so the fallback
	// below is exercised rather than assumed. See TestTiledFallbackMatches.
	//
	// Every tile in a frame carries a BYTE-IDENTICAL prelude: same SOF3, same
	// Huffman tables, same scan header. Verified across all 247 tiles of an
	// A7R V frame. So the tables are parsed once and shared; rebuilding them
	// per tile meant four table allocations and a 1 KB lookup fill, hundreds of
	// times over, for the same answer.
	first, err := parseLJPEG(data[offs[0] : offs[0]+lens[0]])
	if err != nil {
		return nil, fmt.Errorf("tile 0: %w", err)
	}
	if first.width*2 != tileW || first.height*2 != tileH {
		return nil, fmt.Errorf("tile 0 is %dx%d, which does not tile %dx%d",
			first.width, first.height, tileW, tileH)
	}
	preludeLen := len(data[offs[0]:offs[0]+lens[0]]) - len(first.scan)
	prelude := data[offs[0] : offs[0]+preludeLen]

	workers := runtime.GOMAXPROCS(0)
	if workers > len(offs) {
		workers = len(offs)
	}
	var wg sync.WaitGroup
	errs := make([]error, len(offs))
	next := make(chan int, len(offs))
	for i := range offs {
		next <- i
	}
	close(next)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// One scratch per worker, reused across every tile it handles.
			sc := &tileScratch{}
			f := *first // a copy, so each worker can retarget scan independently
			for i := range next {
				tile := data[offs[i] : offs[i]+lens[i]]
				// The prelude is shared, so only confirm it and take the scan.
				if forceTileParse || len(tile) <= preludeLen || !bytes.Equal(tile[:preludeLen], prelude) {
					// A tile that differs must be parsed on its own terms
					// rather than decoded against the wrong tables.
					tf, err := parseLJPEG(tile)
					if err != nil {
						errs[i] = fmt.Errorf("tile %d: %w", i, err)
						continue
					}
					if tf.width*2 != tileW || tf.height*2 != tileH {
						errs[i] = fmt.Errorf("tile %d is %dx%d, which does not tile %dx%d",
							i, tf.width, tf.height, tileW, tileH)
						continue
					}
					tx, ty := (i%across)*tileW, (i/across)*tileH
					if err := decodeLJPEGInto(tf, out, width, tx, ty, sc); err != nil {
						errs[i] = fmt.Errorf("tile %d: %w", i, err)
					}
					continue
				}
				f.scan = tile[preludeLen:]
				tx, ty := (i%across)*tileW, (i/across)*tileH
				if err := decodeLJPEGInto(&f, out, width, tx, ty, sc); err != nil {
					errs[i] = fmt.Errorf("tile %d: %w", i, err)
				}
			}
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// forceGenericLJPEG disables the specialised predictor-1 path, so tests can
// hold the two decoders to sample-for-sample agreement on the same stream.
// Every real Sony tile qualifies for the fast path, which would otherwise
// leave the generic loop exercised only by non-Sony fixtures.
var forceGenericLJPEG = false

// forceTileParse disables the shared-prelude fast path.
//
// The encoder chooses Huffman tables, so it is free to emit different ones per
// tile even though every Sony frame seen here uses one set throughout. The
// fallback that handles that must not be dead code: a tile decoded against the
// wrong tables yields plausible pixels rather than an error, which is the one
// failure this decoder cannot afford. Tests flip this to prove both paths agree.
var forceTileParse = false
