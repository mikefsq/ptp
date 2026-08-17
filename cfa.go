package ptp

import "fmt"

// A color filter array readout, undemosaiced.
// CFAColor is the filter over one photosite.
type CFAColor uint8

// The filter colours. Fujifilm's X-Trans and Sony's Bayer both draw from these
// three; no vendor here uses a CMYG array.
const (
	CFARed CFAColor = iota
	CFAGreen
	CFABlue
)

func (c CFAColor) String() string {
	switch c {
	case CFARed:
		return "R"
	case CFAGreen:
		return "G"
	case CFABlue:
		return "B"
	}
	return "?"
}

// CFA is one sensor readout with its mosaic described rather than applied.
type CFA struct {
	// Model is the camera as the file reports it
	Model string

	// Width and Height are the FULL readout, padding included.
	Width, Height int

	// Pixels is row-major, len == Width*Height, one sample per photosite.
	// Values are left exactly as the sensor reported them: not black-subtracted,
	// not scaled, not white-balanced.
	Pixels []uint16

	// BitDepth is the significant bits per sample, e.g. 14. Values occupy the
	// low bits of each uint16; they are not shifted up to fill it.
	BitDepth int

	// Pattern is the mosaic, row-major over PatternWidth x PatternHeight, and
	// is phased to (0,0) of the FULL readout
	// 2x2 for Bayer, 6x6 for X-Trans.
	Pattern                     []CFAColor
	PatternWidth, PatternHeight int

	// BlackLevel is the sensor pedestal, either one value for the whole frame
	// or one per pattern cell (len == PatternWidth*PatternHeight). WhiteLevel is
	// the saturation point.
	BlackLevel []uint16
	WhiteLevel uint16

	// Crop is the vendor's active area within the full readout.
	CropX, CropY, CropWidth, CropHeight int
}

// At returns the sample at (x, y). It panics out of range, like a slice.
func (c *CFA) At(x, y int) uint16 { return c.Pixels[y*c.Width+x] }

// ColorAt returns the filter colour over (x, y) in the full readout.
func (c *CFA) ColorAt(x, y int) CFAColor {
	return c.Pattern[(y%c.PatternHeight)*c.PatternWidth+(x%c.PatternWidth)]
}

// BlackAt returns the pedestal for (x, y), honouring a per-cell black level.
func (c *CFA) BlackAt(x, y int) uint16 {
	switch len(c.BlackLevel) {
	case 0:
		return 0
	case 1:
		return c.BlackLevel[0]
	}
	return c.BlackLevel[(y%c.PatternHeight)*c.PatternWidth+(x%c.PatternWidth)]
}

// IsBayer reports whether the mosaic is a plain 2x2, which is the only shape
// ASCOM's SensorType and the FITS BAYERPAT convention can describe. A 6x6
// X-Trans frame is treaded as mono and not labelled as Bayer.
func (c *CFA) IsBayer() bool { return c.PatternWidth == 2 && c.PatternHeight == 2 }

// PatternString renders the mosaic in reading order, e.g. "RGGB" for Bayer or a
// 36-character string for X-Trans.
func (c *CFA) PatternString() string {
	s := make([]byte, 0, len(c.Pattern))
	for _, p := range c.Pattern {
		s = append(s, p.String()[0])
	}
	return string(s)
}

// Validate reports the first structural inconsistency, so a decoder bug
// surfaces here rather than as a plausible-looking image with wrong colours.
func (c *CFA) Validate() error {
	if c.Width <= 0 || c.Height <= 0 {
		return fmt.Errorf("ptp: CFA has non-positive dimensions %dx%d", c.Width, c.Height)
	}
	if len(c.Pixels) != c.Width*c.Height {
		return fmt.Errorf("ptp: CFA has %d pixels, want %d for %dx%d",
			len(c.Pixels), c.Width*c.Height, c.Width, c.Height)
	}
	if c.PatternWidth <= 0 || c.PatternHeight <= 0 {
		return fmt.Errorf("ptp: CFA has no mosaic pattern")
	}
	if len(c.Pattern) != c.PatternWidth*c.PatternHeight {
		return fmt.Errorf("ptp: CFA pattern is %d entries, want %d for %dx%d",
			len(c.Pattern), c.PatternWidth*c.PatternHeight, c.PatternWidth, c.PatternHeight)
	}
	if n := len(c.BlackLevel); n != 0 && n != 1 && n != c.PatternWidth*c.PatternHeight {
		return fmt.Errorf("ptp: CFA has %d black levels, want 0, 1 or %d",
			n, c.PatternWidth*c.PatternHeight)
	}
	return nil
}
