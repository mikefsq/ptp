package sony

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/mikefsq/ptp"
)

// ARW, Sony's RAW container, decoded to an undemosaiced Bayer readout.
//
// This shares no code with the Fujifilm decoder and should not: an ARW is an
// ordinary little-endian TIFF/EP whose raw plane hangs off a SubIFD, while a RAF
// is a big-endian wrapper around a second, block-relative TIFF. A parser
// generalised over both would be worse at each.
//
// # Variants
//
// The SubIFD's Compression tag selects the path:
//
//	1        uncompressed, BitsPerSample 12 or 14, packed little-endian
//	7        lossless JPEG. What an A7R V and A7R VI write, and NOT a Sony
//	         codec — see ljpeg.go.
//	32767    Sony's own ARW2 lossy scheme, decoded below.
//
// # ARW2 is lossy, and that matters here
//
// The lossy scheme keeps 16 same-colour pixels per 16 bytes as an 11-bit min
// and max plus 7-bit residuals, then maps the result through a piecewise tone
// curve. It is visually excellent and photometrically compromised: the curve is
// non-linear and the residuals are quantised in bright tones. For calibration
// work set the body to uncompressed or lossless RAW. This path exists because a
// frame that already exists in this format is better read than refused.

// SubIFD tags. Sony's private tags are in the 0x70xx range.
const (
	tagNewSubfileType  = 0x00FE
	tagImageWidth      = 0x0100
	tagImageLength     = 0x0101
	tagBitsPerSample   = 0x0102
	tagCompression     = 0x0103
	tagPhotometric     = 0x0106
	tagModel           = 0x0110
	tagStripOffsets    = 0x0111
	tagStripByteCounts = 0x0117
	tagSubIFDs         = 0x014A
	tagCFAPatternDim   = 0x828D
	tagCFAPattern      = 0x828E
	tagSonyToneCurve   = 0x7010
	tagSonyBlackLevel  = 0x7310
	tagWhiteLevel      = 0xC61D
	tagCropOrigin      = 0xC61F // DefaultCropOrigin, the DNG-standard active area
	tagCropSize        = 0xC620 // DefaultCropSize
	tagTileWidth       = 0x0142
	tagTileLength      = 0x0143
	tagTileOffsets     = 0x0144
	tagTileByteCounts  = 0x0145

	photometricCFA  = 32803 // PhotometricInterpretation for a colour filter array
	compressionNone = 1
	compressionJPEG = 7 // TIFF "JPEG"; in a raw IFD that means LOSSLESS JPEG
	compressionSony = 32767
)

// ErrLosslessARW is returned for a strip that claims Sony compression but is
// neither the ARW2 lossy scheme nor a lossless-JPEG stream.
//
// Both variants carry Compression 32767 and nothing else in the container
// distinguishes them, so the JPEG SOI is what tells them apart. A strip with
// neither shape is something this decoder has not seen.
var ErrLosslessARW = errors.New("sony: this ARW uses a compression this decoder does not implement")

// DecodeARW decodes a Sony ARW into an undemosaiced readout.
//
// The returned CFA is the full sensor readout with the mosaic described rather
// than applied. Samples sit in the low bits of each uint16 exactly as recorded;
// for an ARW2 frame that means after the tone curve, which is the only form in
// which those values exist.
func DecodeARW(data []byte) (*ptp.CFA, error) { return DecodeARWInto(data, nil) }

// DecodeARWInto decodes into dst when it is large enough, so a capture loop can
// reuse one buffer instead of allocating a frame's worth every time.
//
// A 61 MP readout is 129 MB. Faulting that many fresh pages in from the OS
// measured 18% of a decode — more than any single line of the codec — and a
// streaming camera pays it on every frame.
//
// The buffer is the CFA's Pixels, so the caller owns its lifetime: do not reuse
// it while anything still holds the previous frame. ptpcam aliases Pixels into
// the ImageFrame it serves, and overwriting that mid-download would corrupt a
// response in flight. Pass nil to allocate, which is what DecodeARW does.
func DecodeARWInto(data []byte, dst []uint16) (*ptp.CFA, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("sony: not an ARW: %d bytes is too short", len(data))
	}
	if string(data[:4]) != "II*\x00" {
		// Sony has shipped only little-endian ARWs. A big-endian one would be a
		// new variant, and reading it with the wrong order silently yields
		// byte-swapped nonsense, so it is refused.
		return nil, errors.New("sony: not an ARW: expected a little-endian TIFF header")
	}
	le := binary.LittleEndian
	root := int(le.Uint32(data[4:]))

	model := ""
	var subIFDs []int
	if err := eachARWEntry(data, root, func(tag, typ uint16, cnt int, val []byte) error {
		switch tag {
		case tagModel:
			model = cstring(val)
		case tagSubIFDs:
			for i := 0; i+4 <= len(val); i += 4 {
				subIFDs = append(subIFDs, int(le.Uint32(val[i:])))
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(subIFDs) == 0 {
		return nil, errors.New("sony: the ARW has no SubIFDs, so there is no raw plane to find")
	}

	for _, sub := range subIFDs {
		out, err := decodeARWSubIFD(data, sub, model, dst)
		if err != nil {
			return nil, err
		}
		if out != nil {
			return out, nil
		}
	}
	return nil, errors.New("sony: no SubIFD holds a colour filter array")
}

// decodeARWSubIFD decodes one SubIFD if it is the raw plane, and returns nil
// without error if it is merely a thumbnail or preview.
func decodeARWSubIFD(data []byte, sub int, model string, dst []uint16) (*ptp.CFA, error) {
	le := binary.LittleEndian
	var (
		width, height int
		bits          int
		compression   int
		photometric   int
		stripOff      int
		stripLen      int
		patW, patH    = 2, 2
		pattern       []ptp.CFAColor
		curvePts      []uint16
		black         []uint16
		white         uint16
		cropX, cropY  int
		cropW, cropH  int
		tileW, tileH  int
		tileOff       []int
		tileLen       []int
	)
	err := eachARWEntry(data, sub, func(tag, typ uint16, cnt int, val []byte) error {
		u16 := func() int {
			if len(val) >= 2 {
				return int(le.Uint16(val))
			}
			return 0
		}
		u32 := func() int {
			if typ == 3 {
				return u16()
			}
			if len(val) >= 4 {
				return int(le.Uint32(val))
			}
			return 0
		}
		switch tag {
		case tagImageWidth:
			width = u32()
		case tagImageLength:
			height = u32()
		case tagBitsPerSample:
			bits = u16()
		case tagCompression:
			compression = u16()
		case tagPhotometric:
			photometric = u16()
		case tagStripOffsets:
			stripOff = u32()
		case tagStripByteCounts:
			stripLen = u32()
		case tagCFAPatternDim:
			if len(val) >= 4 {
				patW, patH = int(le.Uint16(val)), int(le.Uint16(val[2:]))
			}
		case tagCFAPattern:
			pattern = make([]ptp.CFAColor, 0, len(val))
			for _, b := range val {
				// TIFF/EP CFA colours: 0=R 1=G 2=B, the same coding the CFA
				// type uses, but validated rather than assumed.
				if b > 2 {
					return fmt.Errorf("sony: the CFA pattern has an unknown colour %d", b)
				}
				pattern = append(pattern, ptp.CFAColor(b))
			}
		case tagSonyToneCurve:
			for i := 0; i+2 <= len(val); i += 2 {
				curvePts = append(curvePts, le.Uint16(val[i:]))
			}
		case tagSonyBlackLevel:
			for i := 0; i+2 <= len(val); i += 2 {
				black = append(black, le.Uint16(val[i:]))
			}
		case tagWhiteLevel:
			if len(val) >= 2 {
				white = le.Uint16(val)
			}
		case tagCropOrigin:
			if len(val) >= 8 {
				cropX, cropY = int(le.Uint32(val)), int(le.Uint32(val[4:]))
			}
		case tagCropSize:
			if len(val) >= 8 {
				cropW, cropH = int(le.Uint32(val)), int(le.Uint32(val[4:]))
			}
		case tagTileWidth:
			tileW = u32()
		case tagTileLength:
			tileH = u32()
		case tagTileOffsets:
			tileOff = u32slice(le, val, typ)
		case tagTileByteCounts:
			tileLen = u32slice(le, val, typ)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if photometric != photometricCFA {
		return nil, nil // a preview or thumbnail, not the raw plane
	}
	if width <= 0 || height <= 0 {
		return nil, errors.New("sony: the raw SubIFD gives no dimensions")
	}
	tiled := tileW > 0 && tileH > 0 && len(tileOff) > 0
	if !tiled && (stripOff <= 0 || stripLen <= 0 || stripOff+stripLen > len(data)) {
		return nil, fmt.Errorf("sony: the raw strip (%d+%d) runs past the file (%d)",
			stripOff, stripLen, len(data))
	}
	if len(pattern) != patW*patH {
		return nil, fmt.Errorf("sony: the CFA pattern is %d entries, want %d",
			len(pattern), patW*patH)
	}
	var strip []byte
	if !tiled {
		strip = data[stripOff : stripOff+stripLen]
	}

	out := &ptp.CFA{
		Model:         model,
		Width:         width,
		Height:        height,
		BitDepth:      bits,
		Pattern:       pattern,
		PatternWidth:  patW,
		PatternHeight: patH,
		BlackLevel:    black,
		WhiteLevel:    white,
		CropX:         cropX,
		CropY:         cropY,
		CropWidth:     cropW,
		CropHeight:    cropH,
	}
	if cropW == 0 || cropH == 0 {
		// Older bodies record no active area — the ILCE-7 carries neither crop
		// tag — so the whole readout stands as the crop rather than a zero
		// rectangle a caller might act on.
		out.CropX, out.CropY = 0, 0
		out.CropWidth, out.CropHeight = width, height
	}
	if len(black) > 1 && len(black) != patW*patH {
		// Sony stores four levels in mosaic order; anything else is collapsed
		// rather than mismatched against the pattern.
		out.BlackLevel = black[:1]
	}

	switch {
	case compression == compressionNone:
		out.Pixels, err = unpackARW(strip, width, height, bits, dst)
	case compression == compressionSony && bits == 12 && stripLen == width*height:
		out.Pixels, err = unpackARW2Lossy(strip, width, height, curvePts, dst)
		// The curve output exceeds 12 bits, so the depth the file declares is
		// the depth of the STORED value, not of the result.
		out.BitDepth = 14
		if out.WhiteLevel == 0 {
			out.WhiteLevel = 16383
		}
	case tiled && compression == compressionJPEG:
		// An A7R V or A7R VI writes 512x512 tiles, each its own four-component
		// lossless-JPEG stream. That is why the readout is padded to a multiple
		// of 512 in both axes.
		out.Pixels, err = decodeTiledLossless(data, tileOff, tileLen, tileW, tileH, width, height, dst)
	case compression == compressionJPEG || (compression == compressionSony && isLJPEG(strip)):
		// Lossless: a four-component lossless-JPEG stream, scattered back onto
		// the Bayer plane. See decodeLosslessARW — the codec is verified
		// against real lossless-JPEG data, the Sony wrapper around it is not,
		// for want of a body that produces one.
		//
		out.Pixels, err = decodeLosslessARW(strip, width, height, dst)
	case compression == compressionSony:
		return nil, fmt.Errorf("%w (%d bits, %d bytes for %dx%d)",
			ErrLosslessARW, bits, stripLen, width, height)
	default:
		return nil, fmt.Errorf("sony: unknown ARW compression %d", compression)
	}
	if err != nil {
		return nil, err
	}
	if err := out.Validate(); err != nil {
		return nil, err
	}
	return out, nil
}

// unpackARW expands an uncompressed strip. Sony writes 12- and 14-bit samples
// little-endian and packed with no row padding.
func unpackARW(strip []byte, width, height, bits int, dst []uint16) ([]uint16, error) {
	n := width * height
	switch bits {
	case 16:
		if len(strip) < n*2 {
			return nil, fmt.Errorf("sony: the strip holds %d bytes, want %d", len(strip), n*2)
		}
		// Both unpack loops below write every sample, so neither needs the
		// zeroing reuse() does.
		px := reuseAll(dst, n)
		for i := range px {
			px[i] = binary.LittleEndian.Uint16(strip[i*2:])
		}
		return px, nil
	case 12, 14:
		need := (n*bits + 7) / 8
		if len(strip) < need {
			return nil, fmt.Errorf("sony: the strip holds %d bytes, want %d for %d-bit", len(strip), need, bits)
		}
		px := reuseAll(dst, n)
		var acc uint32
		var have uint
		pos := 0
		for i := range px {
			for have < uint(bits) {
				acc |= uint32(strip[pos]) << have
				have += 8
				pos++
			}
			px[i] = uint16(acc & (1<<uint(bits) - 1))
			acc >>= uint(bits)
			have -= uint(bits)
		}
		return px, nil
	}
	return nil, fmt.Errorf("sony: unsupported uncompressed depth %d", bits)
}

// unpackARW2Lossy expands the ARW2 lossy scheme.
//
// Each 16 bytes carry 16 pixels of ONE colour — every other column — as an
// 11-bit max and min with their positions, plus 7-bit residuals for the other
// fourteen, shifted to cover the local range. A 32-column span is therefore
// two groups: the even columns, then the odd ones.
func unpackARW2Lossy(strip []byte, width, height int, curvePts []uint16, dst []uint16) ([]uint16, error) {
	if len(strip) < width*height {
		return nil, fmt.Errorf("sony: the strip holds %d bytes, want %d", len(strip), width*height)
	}
	curve := sonyToneCurve(curvePts)
	// This one KEEPS the zeroing: the group loop stops at width-30, so the
	// rightmost columns of a row can go unwritten.
	px := reuse(dst, width*height)

	var pix [16]uint16
	for row := 0; row < height; row++ {
		line := strip[row*width : (row+1)*width]
		col := 0
		for dp := 0; col < width-30 && dp+16 <= len(line); dp += 16 {
			val := binary.LittleEndian.Uint32(line[dp:])
			max := uint16(val & 0x7FF)
			min := uint16((val >> 11) & 0x7FF)
			imax := (val >> 22) & 0x0F
			imin := (val >> 26) & 0x0F

			// The residuals are 7-bit, so a range wider than 128 is stored
			// shifted; sh is the smallest shift that spans max-min.
			sh := 0
			for sh < 4 && (0x80<<uint(sh)) <= int(max)-int(min) {
				sh++
			}

			bit := 30
			for i := 0; i < 16; i++ {
				switch uint32(i) {
				case imax:
					pix[i] = max
				case imin:
					pix[i] = min
				default:
					// The final residual of a row starts at bit 121, so the
					// 16-bit window reaches one byte past the row. Those bits
					// are never used — a 7-bit field at offset 1 ends at bit
					// 127 — so the absent byte reads as zero rather than
					// widening every row by one, which is how dcraw handles it.
					b := dp + bit>>3
					if b >= len(line) {
						return nil, fmt.Errorf("sony: the ARW2 row %d is truncated at byte %d", row, b)
					}
					w := uint16(line[b])
					if b+1 < len(line) {
						w |= uint16(line[b+1]) << 8
					}
					v := (w >> uint(bit&7)) & 0x7F
					p := uint32(v)<<uint(sh) + uint32(min)
					if p > 0x7FF {
						p = 0x7FF
					}
					pix[i] = uint16(p)
					bit += 7
				}
			}
			for i := 0; i < 16; i++ {
				if col < width {
					px[row*width+col] = curve[pix[i]<<1]
				}
				col += 2
			}
			// Having done one parity across a 32-column span, step back to the
			// other. This is why col advances by 2 above and lands on 1, then 32.
			if col&1 != 0 {
				col--
			} else {
				col -= 31
			}
		}
	}
	return px, nil
}

// sonyToneCurve builds the 4096-entry lookup the lossy scheme is quantised
// against. The four stored points split the range into five segments whose
// steps double: 1, 2, 4, 8, 16. Absent the tag the mapping is the identity.
func sonyToneCurve(pts []uint16) []uint16 {
	curve := make([]uint16, 4096)
	if len(pts) < 4 {
		for i := range curve {
			curve[i] = uint16(i)
		}
		return curve
	}
	knots := [6]int{0, 0, 0, 0, 0, 4095}
	for i := 0; i < 4; i++ {
		knots[i+1] = int(pts[i]>>2) & 0xFFF
	}
	for i := 0; i < 5; i++ {
		for j := knots[i] + 1; j <= knots[i+1] && j < len(curve); j++ {
			curve[j] = curve[j-1] + uint16(1<<uint(i))
		}
	}
	return curve
}

// eachARWEntry walks one little-endian IFD, handing each entry's value bytes to
// fn. Offsets are relative to the start of the file.
func eachARWEntry(data []byte, ifd int, fn func(tag, typ uint16, cnt int, val []byte) error) error {
	if ifd < 0 || ifd+2 > len(data) {
		return fmt.Errorf("sony: the IFD at %d is out of range", ifd)
	}
	le := binary.LittleEndian
	n := int(le.Uint16(data[ifd:]))
	if ifd+2+n*12 > len(data) {
		return fmt.Errorf("sony: the IFD at %d claims %d entries, which do not fit", ifd, n)
	}
	for i := 0; i < n; i++ {
		e := ifd + 2 + i*12
		tag := le.Uint16(data[e:])
		typ := le.Uint16(data[e+2:])
		cnt := int(le.Uint32(data[e+4:]))
		size := arwTypeSize(typ)
		if size == 0 || cnt < 0 {
			continue
		}
		total := size * cnt
		at := e + 8
		if total > 4 {
			at = int(le.Uint32(data[e+8:]))
		}
		if at < 0 || total < 0 || at+total > len(data) {
			continue
		}
		if err := fn(tag, typ, cnt, data[at:at+total]); err != nil {
			return err
		}
	}
	return nil
}

func arwTypeSize(typ uint16) int {
	switch typ {
	case 1, 2, 6, 7:
		return 1
	case 3, 8:
		return 2
	case 4, 9, 11:
		return 4
	case 5, 10, 12:
		return 8
	}
	return 0
}

func cstring(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

// DecodeCFA decodes a frame this camera produced, satisfying ptp.RawDecoder so
// a vendor-neutral driver can get pixels without knowing the container.
func (c *Camera) DecodeCFA(raw []byte) (*ptp.CFA, error) { return DecodeARW(raw) }

// SensorInfo describes what a capture will produce, before one has been taken.
//
// Unlike Fujifilm this is not a guess: every Sony ILCE body is a 2x2 RGGB Bayer
// sensor reading out at 14 bits. Both sample files confirm the pattern, and no
// Sony interchangeable-lens body has shipped with anything else.
func (c *Camera) SensorInfo() (*ptp.CFA, error) {
	return &ptp.CFA{
		Model:         c.Model(),
		BitDepth:      14,
		WhiteLevel:    1<<14 - 1,
		Pattern:       []ptp.CFAColor{ptp.CFARed, ptp.CFAGreen, ptp.CFAGreen, ptp.CFABlue},
		PatternWidth:  2,
		PatternHeight: 2,
	}, nil
}

// u32slice reads a TIFF LONG or SHORT array.
func u32slice(le binary.ByteOrder, val []byte, typ uint16) []int {
	step := 4
	if typ == 3 {
		step = 2
	}
	out := make([]int, 0, len(val)/step)
	for i := 0; i+step <= len(val); i += step {
		if step == 2 {
			out = append(out, int(le.Uint16(val[i:])))
		} else {
			out = append(out, int(le.Uint32(val[i:])))
		}
	}
	return out
}

// reuse hands back dst sized to n when it is big enough, and a fresh slice
// otherwise. A reused buffer is ZEROED, because a decoder that does not write
// every sample would otherwise leak the previous frame into the gaps.
func reuse(dst []uint16, n int) []uint16 {
	if cap(dst) < n {
		return make([]uint16, n)
	}
	dst = dst[:n]
	clear(dst)
	return dst
}

// reuseAll is reuse for a decode path that writes EVERY sample or fails: the
// zeroing is skipped. That is not a micro-optimisation on the machine this
// targets — clearing a 147 MB A7R VI buffer every frame is tens of
// milliseconds of pure memset at Raspberry Pi memory bandwidth, in the pooled
// capture loop that is this package's whole reason to reuse buffers.
//
// Only safe when a partial decode returns an error rather than pixels;
// otherwise the gaps would leak the previous frame.
func reuseAll(dst []uint16, n int) []uint16 {
	if cap(dst) < n {
		return make([]uint16, n)
	}
	return dst[:n]
}
