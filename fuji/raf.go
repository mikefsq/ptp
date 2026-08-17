package fuji

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/mikefsq/ptp"
)

// RAF, Fujifilm's RAW container, decoded to an undemosaiced X-Trans readout.
//
// Nothing here is shared with the Sony decoder. The two containers agree only
// that they are both broadly TIFF-shaped, and even that is misleading: a RAF is
// a big-endian wrapper whose CFA block contains a SECOND, little-endian TIFF
// whose offsets are relative to the block rather than the file. Factoring that
// against Sony's ordinary TIFF would produce a parser that fits neither.
//
// # Layout
//
// Big-endian outer container:
//
//	0    "FUJIFILMCCD-RAW"
//	28   model string
//	84   six uint32: jpegOff jpegLen cfaHdrOff cfaHdrLen cfaOff cfaLen
//
// The CFA header at cfaHdrOff is a count followed by {uint16 tag, uint16 len,
// data} records — this is where the 6x6 mosaic lives. The CFA block at cfaOff
// opens with the inner TIFF, whose sole entry points at a FujiIFD carrying the
// geometry, bit depth and strip location.
//
// # Why the mosaic is read from the file
//
// Tag 0x0131 carries the actual 6x6 pattern for the actual body. The published
// tables map camera model to one of two known phases, which means every new
// model is unsupported until somebody adds it — an X-T5 appears in none of
// them. Reading it per-frame makes the question not arise.
//
// The pattern is phased to the full readout, which is what DecodeRAF returns,
// so it applies as-is. Crop the frame and it must be re-phased.

// Container and tag constants.
const (
	rafMagic    = "FUJIFILMCCD-RAW"
	rafOffsets  = 84 // where the six uint32 block offsets start
	rafModelAt  = 28
	rafModelLen = 32

	// CFA-header tags (big-endian records).
	tagRawDims    = 0x0100 // (height, width) of the full readout
	tagActiveOrig = 0x0110 // (y, x) origin of the vendor's active area
	tagActiveDims = 0x0111 // (height, width) of the vendor's active area
	tagXTransPat  = 0x0131 // 36 bytes, 0=R 1=G 2=B, row-major 6x6

	// FujiIFD tags (inside the CFA block's little-endian TIFF).
	tagFujiIFD         = 0xF000
	tagRawFullWidth    = 0xF001
	tagRawFullHeight   = 0xF002
	tagBitsPerSample   = 0xF003
	tagStripOffsets    = 0xF007
	tagStripByteCounts = 0xF008
	tagBlackLevel      = 0xF00A

	// rawTypeXTrans is the compressed strip's raw-type byte for an X-Trans
	// sensor; the Bayer GFX bodies use 0.
	rawTypeXTrans = 16
)

// A Fujifilm body offers THREE RAW recording modes, not two, and the last two
// are different codecs rather than two settings of one:
//
//	UNCOMPRESSED         two bytes per photosite. Decoded natively.
//	LOSSLESS COMPRESSED  entropy-coded, bit-exact. Decoded natively.
//	COMPRESSED           lossy. A DIFFERENT coder. Not implemented, and there is
//	                     still no sample of it here to work from.
//
// Lumping the last two together would be a real error rather than a naming one:
// a decoder written for one produces confident nonsense on the other. There is
// still no sample of the lossy mode here, so it cannot even be told apart from
// the lossless one by inspection yet — capture one and diff the strip header.

// ErrCompressedRAF is returned for a compressed RAF this decoder cannot unpack.
//
// The lossless mode decodes natively. This is for the lossy mode and for any
// variant whose parameters fall outside what the native decoder handles: a
// half-right entropy decoder yields a plausible image with wrong values, which
// is precisely the failure calibration cannot tolerate, so it refuses rather
// than approximates.
var ErrCompressedRAF = errors.New("fuji: this RAF uses a compression this decoder does not implement")

// compressedSignature marks the start of a compressed strip.
const compressedSignature = 0x4953

// compressedHeader is the 16-byte descriptor at the head of a compressed strip,
// followed by one big-endian uint32 size per block.
//
// Verified against an X-Pro2 frame: total_lines*6 == height, blocks*blockSize ==
// roundedWidth, and the block sizes sum to exactly the remaining strip bytes.
type compressedHeader struct {
	Version     int
	RawType     int // 16 for X-Trans, 0 for the Bayer GFX bodies
	Bits        int
	Height      int
	RoundedWide int
	Width       int
	BlockSize   int
	BlocksInRow int
	TotalLines  int

	// BlockSizes is the compressed length of each block, and DataOffset where
	// the first one starts, both relative to the strip.
	BlockSizes []int
	DataOffset int
}

// parseCompressedHeader reads the strip descriptor. It reports whether the
// strip is compressed at all; an uncompressed strip has no signature.
func parseCompressedHeader(strip []byte) (*compressedHeader, bool) {
	if len(strip) < 16 || binary.BigEndian.Uint16(strip) != compressedSignature {
		return nil, false
	}
	be := binary.BigEndian
	h := &compressedHeader{
		Version:     int(strip[2]),
		RawType:     int(strip[3]),
		Bits:        int(strip[4]),
		Height:      int(be.Uint16(strip[5:])),
		RoundedWide: int(be.Uint16(strip[7:])),
		Width:       int(be.Uint16(strip[9:])),
		BlockSize:   int(be.Uint16(strip[11:])),
		BlocksInRow: int(strip[13]),
		TotalLines:  int(be.Uint16(strip[14:])),
	}
	if h.BlocksInRow <= 0 || 16+4*h.BlocksInRow > len(strip) {
		return nil, false
	}
	h.BlockSizes = make([]int, h.BlocksInRow)
	for i := range h.BlockSizes {
		h.BlockSizes[i] = int(be.Uint32(strip[16+i*4:]))
	}
	// The block table is padded so the compressed data starts 16-byte aligned.
	off := 16 + 4*h.BlocksInRow
	if off&0xC != 0 {
		off += 0x10 - (off & 0xC)
	}
	h.DataOffset = off
	return h, true
}

// DecodeRAF decodes a Fujifilm RAF into an undemosaiced readout.
//
// The returned CFA is the FULL sensor readout, with the vendor's active area
// reported in the Crop fields rather than applied. Samples are 14-bit in the
// low bits of each uint16, exactly as the sensor reported them. See ptp.CFA for
// what the margin actually contains — it is padding, not optical black.
func DecodeRAF(data []byte) (*ptp.CFA, error) {
	if len(data) < rafOffsets+24 {
		return nil, fmt.Errorf("fuji: not a RAF: %d bytes is too short for the header", len(data))
	}
	if string(data[:len(rafMagic)]) != rafMagic {
		return nil, errors.New("fuji: not a RAF: the magic is missing")
	}

	be := binary.BigEndian
	cfaHdrOff := be.Uint32(data[rafOffsets+8:])
	cfaHdrLen := be.Uint32(data[rafOffsets+12:])
	cfaOff := be.Uint32(data[rafOffsets+16:])
	cfaLen := be.Uint32(data[rafOffsets+20:])

	hdr, err := slice(data, cfaHdrOff, cfaHdrLen, "CFA header")
	if err != nil {
		return nil, err
	}
	block, err := slice(data, cfaOff, cfaLen, "CFA block")
	if err != nil {
		return nil, err
	}

	out := &ptp.CFA{Model: rafModel(data)}
	if err := readCFAHeader(hdr, out); err != nil {
		return nil, err
	}

	strip, bits, black, err := readFujiIFD(block, out)
	if err != nil {
		return nil, err
	}
	out.BitDepth = bits
	if len(black) > 0 {
		out.BlackLevel = black
	}
	// The saturation point is not recorded in the file, so it is the depth's
	// full scale. A caller doing photometry should treat it as an upper bound
	// and find the true clipping level from the data.
	out.WhiteLevel = uint16(1<<uint(bits) - 1)

	if h, ok := parseCompressedHeader(strip); ok {
		px, err := decodeCompressed(strip, h, out)
		if err != nil {
			return nil, err
		}
		out.Pixels = px
		if err := out.Validate(); err != nil {
			return nil, err
		}
		return out, nil
	}
	if n := out.Width * out.Height * 2; len(strip) != n {
		return nil, fmt.Errorf("fuji: the strip holds %d bytes, want %d for an uncompressed %dx%d frame",
			len(strip), n, out.Width, out.Height)
	}

	out.Pixels = make([]uint16, out.Width*out.Height)
	for i := range out.Pixels {
		// Little-endian inside the CFA block, though the container around it is
		// big-endian.
		out.Pixels[i] = binary.LittleEndian.Uint16(strip[i*2:])
	}
	if err := out.Validate(); err != nil {
		return nil, err
	}
	return out, nil
}

// rafModel reads the NUL-padded model string from the outer header.
func rafModel(data []byte) string {
	end := rafModelAt + rafModelLen
	if end > len(data) {
		end = len(data)
	}
	s := data[rafModelAt:end]
	for i, b := range s {
		if b == 0 {
			return string(s[:i])
		}
	}
	return string(s)
}

// readCFAHeader walks the big-endian {tag, len, data} records, filling in the
// geometry and the mosaic.
func readCFAHeader(hdr []byte, out *ptp.CFA) error {
	if len(hdr) < 4 {
		return errors.New("fuji: the CFA header is truncated")
	}
	be := binary.BigEndian
	n := be.Uint32(hdr)
	off := 4
	for i := uint32(0); i < n; i++ {
		if off+4 > len(hdr) {
			break
		}
		tag := be.Uint16(hdr[off:])
		ln := int(be.Uint16(hdr[off+2:]))
		off += 4
		if off+ln > len(hdr) {
			break
		}
		v := hdr[off : off+ln]
		off += ln

		switch tag {
		case tagRawDims:
			if ln >= 4 {
				// Stored (height, width), which is the opposite of the order
				// every other field in this file uses.
				out.Height = int(be.Uint16(v))
				out.Width = int(be.Uint16(v[2:]))
			}
		case tagActiveOrig:
			if ln >= 4 {
				out.CropY = int(be.Uint16(v))
				out.CropX = int(be.Uint16(v[2:]))
			}
		case tagActiveDims:
			if ln >= 4 {
				out.CropHeight = int(be.Uint16(v))
				out.CropWidth = int(be.Uint16(v[2:]))
			}
		case tagXTransPat:
			if ln != 36 {
				return fmt.Errorf("fuji: the X-Trans pattern is %d bytes, want 36", ln)
			}
			out.PatternWidth, out.PatternHeight = 6, 6
			out.Pattern = make([]ptp.CFAColor, 36)
			for j, b := range v {
				if b > 2 {
					return fmt.Errorf("fuji: the X-Trans pattern has an unknown colour %d at %d", b, j)
				}
				out.Pattern[j] = ptp.CFAColor(b)
			}
		}
	}
	if out.Pattern == nil {
		return errors.New("fuji: the RAF carries no X-Trans pattern (tag 0x0131)")
	}
	return nil
}

// readFujiIFD parses the little-endian TIFF at the head of the CFA block and
// returns the pixel strip along with the bit depth and per-cell black levels.
// All offsets in that TIFF are relative to the block, not the file.
func readFujiIFD(block []byte, out *ptp.CFA) (strip []byte, bits int, black []uint16, err error) {
	if len(block) < 8 || string(block[:4]) != "II*\x00" {
		return nil, 0, nil, errors.New("fuji: the CFA block does not open with a little-endian TIFF")
	}
	le := binary.LittleEndian
	root := int(le.Uint32(block[4:]))

	sub, ok := findFujiSubIFD(block, root)
	if !ok {
		return nil, 0, nil, errors.New("fuji: the CFA block has no FujiIFD (tag 0xF000)")
	}

	var stripOff, stripLen int
	err = eachTIFFEntry(block, sub, func(tag uint16, typ uint16, cnt int, val []byte, at int) error {
		switch tag {
		case tagRawFullWidth:
			out.Width = int(le.Uint32(val))
		case tagRawFullHeight:
			out.Height = int(le.Uint32(val))
		case tagBitsPerSample:
			bits = int(le.Uint32(val))
		case tagStripOffsets:
			stripOff = int(le.Uint32(val))
		case tagStripByteCounts:
			stripLen = int(le.Uint32(val))
		case tagBlackLevel:
			// 36 entries: one per cell of the 6x6 mosaic.
			if cnt == 36 && len(val) >= 36*4 {
				black = make([]uint16, 36)
				for i := range black {
					black[i] = uint16(le.Uint32(val[i*4:]))
				}
			} else if cnt >= 1 && len(val) >= 4 {
				black = []uint16{uint16(le.Uint32(val))}
			}
		}
		return nil
	})
	if err != nil {
		return nil, 0, nil, err
	}
	if bits == 0 {
		return nil, 0, nil, errors.New("fuji: the FujiIFD gives no bit depth")
	}
	if stripOff <= 0 || stripLen <= 0 {
		return nil, 0, nil, errors.New("fuji: the FujiIFD gives no pixel strip")
	}
	if stripOff+stripLen > len(block) {
		return nil, 0, nil, fmt.Errorf("fuji: the pixel strip (%d+%d) runs past the CFA block (%d)",
			stripOff, stripLen, len(block))
	}
	return block[stripOff : stripOff+stripLen], bits, black, nil
}

// findFujiSubIFD returns the offset of the FujiIFD referenced by the root IFD.
func findFujiSubIFD(block []byte, root int) (int, bool) {
	le := binary.LittleEndian
	var sub int
	var found bool
	// The pointer's type is 13, a non-standard "IFD" that is a uint32 in all but
	// name, so it is read as one.
	eachTIFFEntry(block, root, func(tag uint16, typ uint16, cnt int, val []byte, at int) error {
		if tag == tagFujiIFD && len(val) >= 4 {
			sub, found = int(le.Uint32(val)), true
		}
		return nil
	})
	return sub, found
}

// eachTIFFEntry walks one little-endian IFD, handing each entry's value bytes to
// fn. Values of four bytes or fewer live inline; longer ones are at an offset
// relative to the start of block.
func eachTIFFEntry(block []byte, ifd int, fn func(tag, typ uint16, cnt int, val []byte, at int) error) error {
	if ifd < 0 || ifd+2 > len(block) {
		return fmt.Errorf("fuji: the IFD at %d is out of range", ifd)
	}
	le := binary.LittleEndian
	n := int(le.Uint16(block[ifd:]))
	if ifd+2+n*12 > len(block) {
		return fmt.Errorf("fuji: the IFD at %d claims %d entries, which do not fit", ifd, n)
	}
	for i := 0; i < n; i++ {
		e := ifd + 2 + i*12
		tag := le.Uint16(block[e:])
		typ := le.Uint16(block[e+2:])
		cnt := int(le.Uint32(block[e+4:]))
		size := tiffTypeSize(typ)
		if size == 0 || cnt < 0 {
			continue
		}
		total := size * cnt
		at := e + 8
		if total > 4 {
			at = int(le.Uint32(block[e+8:]))
		}
		if at < 0 || at+total > len(block) || total < 0 {
			continue
		}
		if err := fn(tag, typ, cnt, block[at:at+total], at); err != nil {
			return err
		}
	}
	return nil
}

// tiffTypeSize is the byte width of one TIFF value of the given type, or 0 when
// the type is unknown. Type 13 is Fujifilm's IFD pointer, a uint32.
func tiffTypeSize(typ uint16) int {
	switch typ {
	case 1, 2, 6, 7:
		return 1
	case 3, 8:
		return 2
	case 4, 9, 11, 13:
		return 4
	case 5, 10, 12:
		return 8
	}
	return 0
}

// slice bounds-checks one of the container's declared blocks.
func slice(data []byte, off, ln uint32, what string) ([]byte, error) {
	if uint64(off)+uint64(ln) > uint64(len(data)) {
		return nil, fmt.Errorf("fuji: the %s (%d+%d) runs past the file (%d)", what, off, ln, len(data))
	}
	return data[off : off+ln], nil
}

// DecodeCFA decodes a frame this camera produced, satisfying ptp.RawDecoder so
// a vendor-neutral driver can get pixels without knowing the container.
func (c *Camera) DecodeCFA(raw []byte) (*ptp.CFA, error) { return DecodeRAF(raw) }

// SensorInfo describes what a capture will produce, before one has been taken.
//
// Every current Fujifilm body reads out at 14 bits. The mosaic is reported as
// 6x6 X-Trans, which is CONSERVATIVE rather than certain: the X series is
// X-Trans but the GFX and entry-level lines are Bayer, and no PTP property
// distinguishes them — the raw type byte lives in the compressed strip, which
// only exists once a frame does.
//
// Guessing X-Trans is the safe direction. A driver turns a 6x6 mosaic into
// "monochrome", which claims nothing about colour; guessing Bayer would have it
// announce RGGB, and a client would then debayer X-Trans data with a 2x2 kernel
// and get confidently wrong colour. The first decoded frame replaces this with
// the truth either way.
func (c *Camera) SensorInfo() (*ptp.CFA, error) {
	pat := make([]ptp.CFAColor, 36)
	for i := range pat {
		pat[i] = ptp.CFAGreen
	}
	return &ptp.CFA{
		Model:         c.Model(),
		BitDepth:      14,
		WhiteLevel:    1<<14 - 1,
		Pattern:       pat,
		PatternWidth:  6,
		PatternHeight: 6,
	}, nil
}
