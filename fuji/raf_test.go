package fuji

import (
	"encoding/binary"
	"errors"
	"github.com/mikefsq/ptp"
	"strings"
	"testing"
)

// A synthetic RAF, so the decoder is covered without a 50 MB sample in the
// repository. The layout mirrors a real file exactly: big-endian outer
// container, a tag-record CFA header, and a CFA block opening with its own
// little-endian TIFF whose offsets are relative to the block.
type rafBuilder struct {
	model       string
	w, h        int
	bits        int
	pattern     []byte // 36 entries, 0=R 1=G 2=B
	black       []uint32
	pixels      []uint16
	stripBytes  []byte // overrides pixels, to fake a compressed strip
	omitPattern bool
}

func (b rafBuilder) build() []byte {
	be, le := binary.BigEndian, binary.LittleEndian

	// --- CFA header: uint32 count, then {uint16 tag, uint16 len, data}.
	var hdr []byte
	put := func(tag uint16, data []byte) {
		var rec [4]byte
		be.PutUint16(rec[0:], tag)
		be.PutUint16(rec[2:], uint16(len(data)))
		hdr = append(hdr, rec[:]...)
		hdr = append(hdr, data...)
	}
	nTags := 4
	if b.omitPattern {
		nTags = 3
	}
	hdr = make([]byte, 4)
	be.PutUint32(hdr, uint32(nTags))
	dims := make([]byte, 4)
	be.PutUint16(dims[0:], uint16(b.h))
	be.PutUint16(dims[2:], uint16(b.w))
	put(tagRawDims, dims)
	orig := make([]byte, 4)
	be.PutUint16(orig[0:], 2) // y
	be.PutUint16(orig[2:], 4) // x
	put(tagActiveOrig, orig)
	act := make([]byte, 4)
	be.PutUint16(act[0:], uint16(b.h-4))
	be.PutUint16(act[2:], uint16(b.w-8))
	put(tagActiveDims, act)
	if !b.omitPattern {
		put(tagXTransPat, b.pattern)
	}

	// --- CFA block: inner TIFF, FujiIFD, then the strip at 2048.
	const (
		fujiIFDAt = 26
		stripAt   = 2048
	)
	strip := b.stripBytes
	if strip == nil {
		strip = make([]byte, len(b.pixels)*2)
		for i, p := range b.pixels {
			le.PutUint16(strip[i*2:], p)
		}
	}
	block := make([]byte, stripAt)
	copy(block, "II*\x00")
	le.PutUint32(block[4:], 8)
	le.PutUint16(block[8:], 1) // one entry in the root IFD
	le.PutUint16(block[10:], tagFujiIFD)
	le.PutUint16(block[12:], 13) // Fujifilm's IFD-pointer type
	le.PutUint32(block[14:], 1)
	le.PutUint32(block[18:], fujiIFDAt)
	le.PutUint32(block[22:], 0) // no next IFD

	entries := []struct {
		tag uint16
		val uint32
	}{
		{tagRawFullWidth, uint32(b.w)},
		{tagRawFullHeight, uint32(b.h)},
		{tagBitsPerSample, uint32(b.bits)},
		{tagStripOffsets, stripAt},
		{tagStripByteCounts, uint32(len(strip))},
	}
	n := len(entries)
	blackAt := fujiIFDAt + 2 + (n+1)*12 + 4
	le.PutUint16(block[fujiIFDAt:], uint16(n+1))
	for i, e := range entries {
		o := fujiIFDAt + 2 + i*12
		le.PutUint16(block[o:], e.tag)
		le.PutUint16(block[o+2:], 4) // LONG
		le.PutUint32(block[o+4:], 1)
		le.PutUint32(block[o+8:], e.val)
	}
	o := fujiIFDAt + 2 + n*12
	le.PutUint16(block[o:], tagBlackLevel)
	le.PutUint16(block[o+2:], 4)
	le.PutUint32(block[o+4:], uint32(len(b.black)))
	le.PutUint32(block[o+8:], uint32(blackAt))
	for i, v := range b.black {
		le.PutUint32(block[blackAt+i*4:], v)
	}
	block = append(block, strip...)

	// --- Outer container.
	const cfaHdrAt = 256
	cfaAt := cfaHdrAt + len(hdr)
	out := make([]byte, cfaAt)
	copy(out, rafMagic)
	copy(out[rafModelAt:], b.model)
	be.PutUint32(out[rafOffsets+0:], 0) // jpeg offset
	be.PutUint32(out[rafOffsets+4:], 0) // jpeg length
	be.PutUint32(out[rafOffsets+8:], cfaHdrAt)
	be.PutUint32(out[rafOffsets+12:], uint32(len(hdr)))
	be.PutUint32(out[rafOffsets+16:], uint32(cfaAt))
	be.PutUint32(out[rafOffsets+20:], uint32(len(block)))
	copy(out[cfaHdrAt:], hdr)
	return append(out, block...)
}

// xtransPattern is the X-T2's, as read from a real file: two of the six rows
// carry the red and blue pairs.
var xtransPattern = []byte{
	2, 1, 1, 0, 1, 1,
	0, 1, 1, 2, 1, 1,
	1, 2, 0, 1, 0, 2,
	0, 1, 1, 2, 1, 1,
	2, 1, 1, 0, 1, 1,
	1, 0, 2, 1, 2, 0,
}

func sampleRAF() rafBuilder {
	const w, h = 12, 12
	px := make([]uint16, w*h)
	for i := range px {
		px[i] = uint16(1000 + i)
	}
	black := make([]uint32, 36)
	for i := range black {
		black[i] = 1022
	}
	return rafBuilder{
		model: "X-T2", w: w, h: h, bits: 14,
		pattern: xtransPattern, black: black, pixels: px,
	}
}

func TestDecodeRAF(t *testing.T) {
	b := sampleRAF()
	cfa, err := DecodeRAF(b.build())
	if err != nil {
		t.Fatalf("DecodeRAF: %v", err)
	}
	if cfa.Model != "X-T2" {
		t.Errorf("model %q", cfa.Model)
	}
	if cfa.Width != 12 || cfa.Height != 12 {
		t.Errorf("readout %dx%d, want 12x12", cfa.Width, cfa.Height)
	}
	if cfa.BitDepth != 14 {
		t.Errorf("bit depth %d, want 14", cfa.BitDepth)
	}
	if cfa.WhiteLevel != 16383 {
		t.Errorf("white level %d, want 16383", cfa.WhiteLevel)
	}
	// The active area is reported, NOT applied: a caller doing calibration
	// needs the optical-black margin that the crop would discard.
	if cfa.CropX != 4 || cfa.CropY != 2 || cfa.CropWidth != 4 || cfa.CropHeight != 8 {
		t.Errorf("crop %d,%d %dx%d, want 4,2 4x8", cfa.CropX, cfa.CropY, cfa.CropWidth, cfa.CropHeight)
	}
	if len(cfa.Pixels) != cfa.Width*cfa.Height {
		t.Error("the frame was cropped, losing the optical-black margin")
	}
	if len(cfa.Pixels) != 144 {
		t.Fatalf("%d pixels, want 144", len(cfa.Pixels))
	}
	for i, got := range cfa.Pixels {
		if want := uint16(1000 + i); got != want {
			t.Fatalf("pixel %d is %d, want %d", i, got, want)
		}
	}
}

// The mosaic must come from the file, not a model lookup: the published tables
// cover neither an X-T5 nor anything newer.
func TestDecodeRAFReadsTheMosaicFromTheFile(t *testing.T) {
	b := sampleRAF()
	cfa, err := DecodeRAF(b.build())
	if err != nil {
		t.Fatal(err)
	}
	if cfa.PatternWidth != 6 || cfa.PatternHeight != 6 {
		t.Fatalf("mosaic %dx%d, want 6x6", cfa.PatternWidth, cfa.PatternHeight)
	}
	const want = "BGGRGGRGGBGGGBRGRBRGGBGGBGGRGGGRBGBR"
	if got := cfa.PatternString(); got != want {
		t.Errorf("pattern %q, want %q", got, want)
	}
	// X-Trans is 6x6, so it must never claim to be Bayer — ASCOM's SensorType
	// and the FITS BAYERPAT convention can only describe 2x2.
	if cfa.IsBayer() {
		t.Error("a 6x6 X-Trans mosaic reports itself as Bayer")
	}
	if got := cfa.ColorAt(0, 0); got.String() != "B" {
		t.Errorf("ColorAt(0,0) = %s, want B", got)
	}
	if got := cfa.ColorAt(3, 1); got.String() != "B" {
		t.Errorf("ColorAt(3,1) = %s, want B", got)
	}
	// Phase must repeat every six, in both axes.
	if cfa.ColorAt(6, 7) != cfa.ColorAt(0, 1) {
		t.Error("the mosaic does not repeat with period 6")
	}
}

func TestDecodeRAFBlackLevelIsPerCell(t *testing.T) {
	cfa, err := DecodeRAF(sampleRAF().build())
	if err != nil {
		t.Fatal(err)
	}
	if len(cfa.BlackLevel) != 36 {
		t.Fatalf("%d black levels, want one per 6x6 cell", len(cfa.BlackLevel))
	}
	if got := cfa.BlackAt(7, 3); got != 1022 {
		t.Errorf("BlackAt = %d, want 1022", got)
	}
}

// A compressed RAF must say so, and say what it is. Returning noise, or a
// half-decoded frame, would be far worse than refusing: the result feeds
// calibration.
func TestDecodeRAFRejectsCompressed(t *testing.T) {
	b := sampleRAF()
	strip := []byte{
		0x49, 0x53, 0x01, 0x10, 0x0e, // signature, version 1, X-Trans, 14-bit
		0x00, 0x0c, // height 12
		0x00, 0x0c, // rounded width 12
		0x00, 0x0c, // width 12
		0x00, 0x0c, // block size 12
		0x01,       // one block in a row
		0x00, 0x02, // two lines of six rows
		0, 0, 0, 32, // the single block's compressed size
	}
	b.stripBytes = append(strip, make([]byte, 64)...)

	_, err := DecodeRAF(b.build())
	if !errors.Is(err, ErrCompressedRAF) {
		t.Fatalf("error %v, want ErrCompressedRAF", err)
	}
	// The refusal must name the parameters, so a user can tell a compressed
	// frame from a corrupt one.
	for _, want := range []string{"14-bit", "12x12", "1 block"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not mention %q: %v", want, err)
		}
	}
}

// A strip that is neither compressed nor the right size for an uncompressed
// frame is corrupt, and must not be reported as compressed.
func TestDecodeRAFRejectsShortStrip(t *testing.T) {
	b := sampleRAF()
	b.stripBytes = make([]byte, 12*12) // half of what uncompressed would need
	_, err := DecodeRAF(b.build())
	if err == nil || errors.Is(err, ErrCompressedRAF) {
		t.Fatalf("error %v, want a plain size complaint", err)
	}
	if !strings.Contains(err.Error(), "want 288") {
		t.Errorf("the error does not say what size was expected: %v", err)
	}
}

func TestDecodeRAFRejectsBadInput(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want string
	}{
		{"empty", nil, "too short"},
		{"short", make([]byte, 20), "too short"},
		{"not a RAF", make([]byte, 512), "magic"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeRAF(tc.data)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %v, want one mentioning %q", err, tc.want)
			}
		})
	}

	b := sampleRAF()
	b.omitPattern = true
	if _, err := DecodeRAF(b.build()); err == nil ||
		!strings.Contains(err.Error(), "0x0131") {
		t.Fatalf("error %v, want one naming the missing pattern tag", err)
	}
}

// A truncated file must not panic: frames arrive over USB and a short read is a
// normal failure, not a programming error.
func TestDecodeRAFTruncatedDoesNotPanic(t *testing.T) {
	full := sampleRAF().build()
	for n := 0; n < len(full); n += 37 {
		if _, err := DecodeRAF(full[:n]); err == nil && n < len(full) {
			t.Fatalf("a %d-byte prefix of a %d-byte RAF decoded without error", n, len(full))
		}
	}
}

// The compressed framing is parsed even though the entropy coder is not, so the
// refusal names what the file actually is. Header bytes are an X-Pro2's.
func TestParseCompressedHeader(t *testing.T) {
	strip := []byte{
		0x49, 0x53, 0x01, 0x10, 0x0e, // signature, version 1, X-Trans, 14-bit
		0x0f, 0xc6, // height 4038
		0x18, 0x00, // rounded width 6144
		0x17, 0xa0, // width 6048
		0x03, 0x00, // block size 768
		0x08,       // 8 blocks in a row
		0x02, 0xa1, // 673 lines
	}
	sizes := []uint32{2697264, 2807360, 2915472, 3057552, 3008320, 2889344, 2814320, 2436912}
	for _, s := range sizes {
		strip = append(strip, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
	}
	h, ok := parseCompressedHeader(strip)
	if !ok {
		t.Fatal("a compressed strip was not recognised")
	}
	if h.Bits != 14 || h.RawType != 16 {
		t.Errorf("raw type %d at %d bits, want X-Trans at 14", h.RawType, h.Bits)
	}
	// The two internal consistency checks a real file satisfies.
	if h.TotalLines*6 != h.Height {
		t.Errorf("%d lines of 6 rows is %d, not the stated height %d",
			h.TotalLines, h.TotalLines*6, h.Height)
	}
	if h.BlocksInRow*h.BlockSize != h.RoundedWide {
		t.Errorf("%d blocks of %d is %d, not the rounded width %d",
			h.BlocksInRow, h.BlockSize, h.BlocksInRow*h.BlockSize, h.RoundedWide)
	}
	if h.DataOffset != 48 {
		t.Errorf("data starts at %d, want 48 (16-byte aligned past the block table)", h.DataOffset)
	}
	// An uncompressed strip must not be mistaken for a compressed one.
	if _, ok := parseCompressedHeader(make([]byte, 4096)); ok {
		t.Error("an uncompressed strip was read as compressed")
	}
}

// The vendor packages must satisfy the decode capability, or a vendor-neutral
// driver cannot get pixels out of them.
var _ ptp.RawDecoder = (*Camera)(nil)
