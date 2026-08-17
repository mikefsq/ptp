package sony

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// The specialised decode path for what a Sony body actually writes.
//
// Every lossless tile from an A7R V or A7R VI carries the same scan
// parameters: four components, predictor 1, point transform 0, no restart
// interval. Measured on real tiles from three frames across both bodies —
// not assumed from the spec, which permits seven predictors and any of it to
// vary per tile.
//
// Predictor 1 is Ra, the previous sample of the same component. That makes
// the generic decoder's two line buffers almost entirely dead weight: the
// interior of a row predicts from a value that was computed four samples ago
// and is still in a register, and column 0 — the only place the previous row
// is consulted — predicts from the sample directly above, which is four
// values, not a buffer. So this path carries left[4] and above[4] and no
// lines at all.
//
// The generic decoder in ljpeg.go remains the reference. It handles every
// stream the standard allows, and the two are held to sample-for-sample
// agreement by TestFastPathMatchesGeneric on real frames and by the
// known-plaintext test on a synthesised stream.

// destuffPad is how many zero bytes destuff appends past the scan, so the bit
// reader can always load a whole word without a bounds branch.
const destuffPad = 8

// destuff copies a scan into dst with the JPEG byte stuffing removed — a 0xFF
// data byte is written 0xFF 0x00, and 0xFF followed by anything else is a
// marker, which ends the scan. A lone trailing 0xFF is likewise the start of
// a truncated marker rather than data, because data 0xFF is always stuffed.
//
// Stuffed bytes are rare in real tiles — 877 in a 271 KB A7R V scan, 0.3% —
// so this runs at the speed of the IndexByte sweeps, which the runtime
// vectorises. The result carries destuffPad zero bytes of padding.
func destuff(scan, dst []byte) []byte {
	dst = dst[:0]
	for {
		i := bytes.IndexByte(scan, 0xFF)
		if i < 0 {
			dst = append(dst, scan...)
			break
		}
		if i+1 >= len(scan) || scan[i+1] != 0x00 {
			dst = append(dst, scan[:i]...)
			break
		}
		dst = append(dst, scan[:i+1]...) // keep the 0xFF, drop the 0x00
		scan = scan[i+2:]
	}
	return append(dst, 0, 0, 0, 0, 0, 0, 0, 0)
}

// dbits reads a destuffed entropy segment MSB-first.
//
// jbits interleaves destuffing with reading: one byte at a time, each
// compared against 0xFF. A wide out-of-order core hides that loop; the
// Raspberry Pi this targets does not. With the stuffing already removed and
// the segment padded, the refill is one unaligned load and three arithmetic
// ops, no branches in the common case.
type dbits struct {
	data []byte // destuffed scan plus destuffPad zero bytes
	acc  uint64 // pending bits, left-aligned: the next bit is the MSB
	n    uint   // how many bits of acc are valid
	pos  int    // next byte to load
}

// fill tops the accumulator up to at least 56 bits. The invariant is the
// standard lookahead refill: pos*8 - n bits have been consumed, so the word
// at pos begins exactly where acc's valid bits end, and re-reading overlapped
// bytes ORs in the same bits they already contributed.
func (b *dbits) fill() {
	if b.pos+8 <= len(b.data) {
		b.acc |= binary.BigEndian.Uint64(b.data[b.pos:]) >> b.n
		b.pos += int((63 - b.n) >> 3)
		b.n |= 56
		return
	}
	// Inside the final word: byte at a time, and bits past the end read as
	// zero, the same as jbits.
	for b.n <= 56 && b.pos < len(b.data) {
		b.acc |= uint64(b.data[b.pos]) << (56 - b.n)
		b.n += 8
		b.pos++
	}
}

// bit reads one bit, for the long-code fallback only.
func (b *dbits) bit() uint32 {
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

// decodeLongD walks the canonical table for a code longer than the flat
// lookup. The twin of decodeLong, on the destuffed reader.
func (t *huffTable) decodeLongD(b *dbits) (int, error) {
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

// decodeLJPEGFast is decodeLJPEGInto for the shape every Sony tile takes:
// four components, predictor 1, point transform 0, no restarts, and a tile
// lying entirely inside the frame, so no store needs a clip check. data is
// the destuffed scan.
func decodeLJPEGFast(f *ljpegFrame, data []byte, dst []uint16, stride, tx, ty int) error {
	tables := [4]*huffTable{
		f.huff[f.comps[0].tableID],
		f.huff[f.comps[1].tableID],
		f.huff[f.comps[2].tableID],
		f.huff[f.comps[3].tableID],
	}
	b := &dbits{data: data}

	// Out-of-range samples are ORed together and checked once at the end,
	// exactly as in the generic path.
	var seen int32
	limit := int32(1)<<uint(f.precision) - 1
	start := int32(1) << (f.precision - 1)

	var above, left [4]int32
	w := f.width
	for y := 0; y < f.height; y++ {
		o0 := (ty+2*y)*stride + tx
		row0 := dst[o0 : o0+2*w]
		row1 := dst[o0+stride : o0+stride+2*w]
		for x := 0; x < w; x++ {
			xx := 2 * x
			for c := 0; c < 4; c++ {
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
					if sbits, err = t.decodeLongD(b); err != nil {
						return err
					}
				}

				// Ordered by frequency; the symbol-16 escape (difference is
				// -32768, no magnitude bits — see decodeLJPEGInto) comes last
				// because no real sample has ever hit it.
				var diff int32
				switch {
				case sbits == 0:
				case sbits < 16:
					k := uint(sbits)
					diff = int32(b.acc >> (64 - k))
					b.acc <<= k
					b.n -= k
					if diff < 1<<(k-1) {
						diff -= 1<<k - 1
					}
				default:
					diff = -32768
				}

				var pred int32
				if x != 0 {
					pred = left[c]
				} else if y != 0 {
					pred = above[c]
				} else {
					pred = start
				}
				v := pred + diff
				left[c] = v
				if x == 0 {
					above[c] = v
				}
				seen |= v

				// The 2x2 scatter, with the component index static enough for
				// the branch predictor: c<2 alternates in a fixed period-4
				// pattern.
				if c < 2 {
					row0[xx+c] = uint16(v)
				} else {
					row1[xx+(c&1)] = uint16(v)
				}
			}
		}
	}
	if seen < 0 || seen > limit {
		return fmt.Errorf("sony: a decoded sample fell outside %d bits, so the stream desynchronised",
			f.precision)
	}
	return nil
}
