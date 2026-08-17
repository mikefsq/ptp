package sony

import (
	"encoding/binary"
	"errors"
	"strings"
	"testing"
)

// A synthetic ARW: an ordinary little-endian TIFF whose raw plane hangs off a
// SubIFD, which is the shape a real one has. Built here so the decoder is
// covered without a 43 MB sample in the repository.
type arwBuilder struct {
	model       string
	w, h        int
	bits        int
	compression int
	pattern     []byte
	black       []uint16
	white       uint16
	strip       []byte
	omitCFA     bool
}

func (b arwBuilder) build() []byte {
	le := binary.LittleEndian

	type entry struct {
		tag, typ uint16
		cnt      uint32
		inline   uint32
		data     []byte // when longer than four bytes
	}
	// Values that do not fit inline are appended after both IFDs, so their
	// offsets are known only once the directory sizes are.
	sub := []entry{
		{tagImageWidth, 4, 1, uint32(b.w), nil},
		{tagImageLength, 4, 1, uint32(b.h), nil},
		{tagBitsPerSample, 3, 1, uint32(b.bits), nil},
		{tagCompression, 3, 1, uint32(b.compression), nil},
		{tagPhotometric, 3, 1, photometricCFA, nil},
		{tagCFAPatternDim, 3, 2, 2 | 2<<16, nil},
		{tagCropOrigin, 4, 2, 0, u32s(le, 4, 2)},
		{tagCropSize, 4, 2, 0, u32s(le, uint32(b.w-8), uint32(b.h-2))},
	}
	if !b.omitCFA {
		sub = append(sub, entry{tagCFAPattern, 1, 4, le32(b.pattern), nil})
	}
	if b.white != 0 {
		sub = append(sub, entry{tagWhiteLevel, 3, 1, uint32(b.white), nil})
	}
	if len(b.black) > 0 {
		d := make([]byte, len(b.black)*2)
		for i, v := range b.black {
			le.PutUint16(d[i*2:], v)
		}
		sub = append(sub, entry{tagSonyBlackLevel, 3, uint32(len(b.black)), 0, d})
	}
	if b.compression == compressionSony {
		// The tone curve the lossy scheme is quantised against, as a real body
		// records it.
		sub = append(sub, entry{tagSonyToneCurve, 3, 4, 0,
			u16s(le, 8000, 10400, 12900, 14100)})
	}
	// StripOffsets is patched once the layout is fixed.
	sub = append(sub,
		entry{tagStripOffsets, 4, 1, 0, nil},
		entry{tagStripByteCounts, 4, 1, uint32(len(b.strip)), nil})
	stripIdx := len(sub) - 2

	modelBytes := append([]byte(b.model), 0)
	root := []entry{
		{tagModel, 2, uint32(len(modelBytes)), 0, modelBytes},
		{tagSubIFDs, 4, 1, 0, nil}, // patched below
	}

	rootAt := 8
	rootSize := 2 + len(root)*12 + 4
	subAt := rootAt + rootSize
	subSize := 2 + len(sub)*12 + 4
	heapAt := subAt + subSize

	// Lay out the oversized values on the heap.
	heap := []byte{}
	place := func(e *entry) {
		if len(e.data) > 4 {
			e.inline = uint32(heapAt + len(heap))
			heap = append(heap, e.data...)
		} else if len(e.data) > 0 {
			e.inline = le32(e.data)
		}
	}
	for i := range root {
		place(&root[i])
	}
	for i := range sub {
		place(&sub[i])
	}
	root[1].inline = uint32(subAt)
	stripAt := heapAt + len(heap)
	sub[stripIdx].inline = uint32(stripAt)

	out := make([]byte, stripAt)
	copy(out, "II*\x00")
	le.PutUint32(out[4:], uint32(rootAt))

	write := func(at int, es []entry) {
		le.PutUint16(out[at:], uint16(len(es)))
		for i, e := range es {
			o := at + 2 + i*12
			le.PutUint16(out[o:], e.tag)
			le.PutUint16(out[o+2:], e.typ)
			le.PutUint32(out[o+4:], e.cnt)
			le.PutUint32(out[o+8:], e.inline)
		}
		le.PutUint32(out[at+2+len(es)*12:], 0)
	}
	write(rootAt, root)
	write(subAt, sub)
	copy(out[heapAt:], heap)
	return append(out, b.strip...)
}

func u32s(le binary.ByteOrder, v ...uint32) []byte {
	b := make([]byte, len(v)*4)
	for i, x := range v {
		le.PutUint32(b[i*4:], x)
	}
	return b
}

func u16s(le binary.ByteOrder, v ...uint16) []byte {
	b := make([]byte, len(v)*2)
	for i, x := range v {
		le.PutUint16(b[i*2:], x)
	}
	return b
}

func le32(b []byte) uint32 {
	var v uint32
	for i := 0; i < len(b) && i < 4; i++ {
		v |= uint32(b[i]) << (8 * uint(i))
	}
	return v
}

// pack14 packs samples the way an uncompressed ARW does: little-endian, no row
// padding, each value in the low bits.
func pack14(px []uint16, bits int) []byte {
	out := make([]byte, (len(px)*bits+7)/8)
	var acc uint32
	var have uint
	pos := 0
	for _, p := range px {
		acc |= uint32(p) << have
		have += uint(bits)
		for have >= 8 {
			out[pos] = byte(acc)
			acc >>= 8
			have -= 8
			pos++
		}
	}
	if have > 0 {
		out[pos] = byte(acc)
	}
	return out
}

func sampleARW(bits int) (arwBuilder, []uint16) {
	const w, h = 16, 4
	px := make([]uint16, w*h)
	for i := range px {
		px[i] = uint16(i * 37 % (1 << uint(bits)))
	}
	return arwBuilder{
		model: "ILCE-7RM5", w: w, h: h, bits: bits,
		compression: compressionNone,
		pattern:     []byte{0, 1, 1, 2}, // RGGB
		black:       []uint16{512, 512, 512, 512},
		white:       16380,
		strip:       pack14(px, bits),
	}, px
}

func TestDecodeARWUncompressed(t *testing.T) {
	for _, bits := range []int{12, 14} {
		b, want := sampleARW(bits)
		cfa, err := DecodeARW(b.build())
		if err != nil {
			t.Fatalf("%d-bit: DecodeARW: %v", bits, err)
		}
		if cfa.Model != "ILCE-7RM5" {
			t.Errorf("model %q", cfa.Model)
		}
		if cfa.Width != 16 || cfa.Height != 4 {
			t.Errorf("readout %dx%d, want 16x4", cfa.Width, cfa.Height)
		}
		if cfa.BitDepth != bits {
			t.Errorf("bit depth %d, want %d", cfa.BitDepth, bits)
		}
		for i, got := range cfa.Pixels {
			if got != want[i] {
				t.Fatalf("%d-bit: pixel %d is %d, want %d", bits, i, got, want[i])
			}
		}
	}
}

func TestDecodeARWMosaicAndLevels(t *testing.T) {
	b, _ := sampleARW(14)
	cfa, err := DecodeARW(b.build())
	if err != nil {
		t.Fatal(err)
	}
	if !cfa.IsBayer() {
		t.Error("a 2x2 Sony mosaic does not report as Bayer, so ASCOM cannot describe it")
	}
	if got := cfa.PatternString(); got != "RGGB" {
		t.Errorf("pattern %q, want RGGB", got)
	}
	if cfa.WhiteLevel != 16380 {
		t.Errorf("white level %d, want 16380", cfa.WhiteLevel)
	}
	if got := cfa.BlackAt(1, 1); got != 512 {
		t.Errorf("BlackAt = %d, want 512", got)
	}
	// The active area is reported, not applied: calibration needs the
	// optical-black margin that cropping would discard.
	if cfa.CropX != 4 || cfa.CropY != 2 || cfa.CropWidth != 8 || cfa.CropHeight != 2 {
		t.Errorf("crop %d,%d %dx%d, want 4,2 8x2", cfa.CropX, cfa.CropY, cfa.CropWidth, cfa.CropHeight)
	}
	if len(cfa.Pixels) != cfa.Width*cfa.Height {
		t.Error("the frame was cropped, losing the optical-black margin")
	}
}

// A body that records no active area must still report a usable rectangle: the
// ILCE-7 carries neither crop tag, and a zero rectangle is something a caller
// might act on.
func TestDecodeARWWithoutCropTagsUsesTheFullReadout(t *testing.T) {
	b, _ := sampleARW(14)
	cfa, err := DecodeARW(b.build())
	if err != nil {
		t.Fatal(err)
	}
	_ = cfa
	b2 := b
	b2.strip = b.strip
	// Rebuild without the crop tags by zeroing the sizes the builder writes.
	raw := b2.build()
	le := binary.LittleEndian
	// Find and neuter DefaultCropSize in the SubIFD by setting its count to 0.
	for i := 0; i+12 <= len(raw); i += 2 {
		if le.Uint16(raw[i:]) == tagCropSize && le.Uint16(raw[i+2:]) == 4 {
			le.PutUint16(raw[i:], 0xFFFF) // an unknown tag is ignored
			break
		}
	}
	got, err := DecodeARW(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.CropWidth != got.Width || got.CropHeight != got.Height {
		t.Errorf("crop %dx%d, want the full readout %dx%d",
			got.CropWidth, got.CropHeight, got.Width, got.Height)
	}
}

// The tone curve is what makes ARW2 non-linear, so its construction is pinned:
// five segments whose step doubles, from the four stored knots.
func TestSonyToneCurve(t *testing.T) {
	curve := sonyToneCurve([]uint16{8000, 10400, 12900, 14100})
	if len(curve) != 4096 {
		t.Fatalf("curve has %d entries, want 4096", len(curve))
	}
	// Knots are the stored values >> 2.
	for _, tc := range []struct{ in, want uint16 }{
		{0, 0},
		{1, 1},
		{2000, 2000},  // end of the unit-step segment
		{2600, 3200},  // +2 per step for 600 steps
		{3225, 5700},  // +4 per step for 625
		{3525, 8100},  // +8 per step for 300
		{4095, 17220}, // +16 per step for 570
	} {
		if got := curve[tc.in]; got != tc.want {
			t.Errorf("curve[%d] = %d, want %d", tc.in, got, tc.want)
		}
	}
	// Monotonic, or the mapping is not a curve.
	for i := 1; i < len(curve); i++ {
		if curve[i] < curve[i-1] {
			t.Fatalf("the curve decreases at %d", i)
		}
	}
	// With no tag it must be the identity rather than all zeroes, which would
	// silently blank a frame.
	id := sonyToneCurve(nil)
	for _, i := range []int{0, 1, 2047, 4095} {
		if id[i] != uint16(i) {
			t.Fatalf("the fallback curve is not the identity at %d", i)
		}
	}
}

func TestDecodeARWRejectsBadInput(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want string
	}{
		{"empty", nil, "too short"},
		{"big-endian TIFF", append([]byte("MM\x00*"), make([]byte, 64)...), "little-endian"},
		{"no SubIFDs", append([]byte("II*\x00\x08\x00\x00\x00\x00\x00"), make([]byte, 32)...), "SubIFDs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeARW(tc.data)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %v, want one mentioning %q", err, tc.want)
			}
		})
	}
}

// A strip claiming Sony compression but matching neither shape is refused.
//
// Both variants carry Compression 32767, so the JPEG SOI is what tells them
// apart: ARW2 lossy is exactly one byte per pixel, lossless opens FFD8. This
// fixture is neither, and decoding it as either would produce noise.
func TestDecodeARWRejectsUnknownCompression(t *testing.T) {
	b, _ := sampleARW(14)
	b.compression = compressionSony
	b.bits = 14
	_, err := DecodeARW(b.build())
	if !errors.Is(err, ErrLosslessARW) {
		t.Fatalf("error %v, want ErrLosslessARW", err)
	}
	// The refusal must describe the strip, so a real file can be told apart
	// from a corrupt one.
	for _, want := range []string{"14 bits", "16x4"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not mention %q: %v", want, err)
		}
	}
}

// A lossless strip is recognised by its JPEG SOI and routed to the LJPEG
// decoder rather than refused. The fixture is not a valid stream, so the
// failure must come from the DECODER — proving the routing happened — not from
// the container rejecting it outright.
func TestDecodeARWRoutesLosslessToLJPEG(t *testing.T) {
	b, _ := sampleARW(14)
	b.compression = compressionSony
	b.bits = 14
	b.strip = append([]byte{0xFF, 0xD8, 0xFF, 0xC3}, make([]byte, 108)...)
	_, err := DecodeARW(b.build())
	if err == nil {
		t.Fatal("a truncated lossless-JPEG stream decoded without error")
	}
	if errors.Is(err, ErrLosslessARW) {
		t.Fatalf("a lossless strip was refused instead of decoded: %v", err)
	}
}

func TestDecodeARWTruncatedDoesNotPanic(t *testing.T) {
	b, _ := sampleARW(14)
	full := b.build()
	for n := 0; n < len(full); n += 13 {
		if _, err := DecodeARW(full[:n]); err == nil {
			t.Fatalf("a %d-byte prefix of a %d-byte ARW decoded without error", n, len(full))
		}
	}
}
