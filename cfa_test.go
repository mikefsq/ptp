package ptp

import "testing"

// bayer is a 4x2 RGGB readout, small enough to reason about by hand.
func bayer() *CFA {
	return &CFA{
		Model: "test", Width: 4, Height: 2,
		Pixels:       []uint16{1, 2, 3, 4, 5, 6, 7, 8},
		BitDepth:     14,
		Pattern:      []CFAColor{CFARed, CFAGreen, CFAGreen, CFABlue},
		PatternWidth: 2, PatternHeight: 2,
		BlackLevel: []uint16{512, 513, 514, 515},
		WhiteLevel: 16383,
	}
}

func TestCFAColorAtRepeatsThePattern(t *testing.T) {
	c := bayer()
	want := [2][4]CFAColor{
		{CFARed, CFAGreen, CFARed, CFAGreen},
		{CFAGreen, CFABlue, CFAGreen, CFABlue},
	}
	for y := 0; y < 2; y++ {
		for x := 0; x < 4; x++ {
			if got := c.ColorAt(x, y); got != want[y][x] {
				t.Errorf("ColorAt(%d,%d) = %s, want %s", x, y, got, want[y][x])
			}
		}
	}
}

func TestCFABlackAtIsPerCell(t *testing.T) {
	c := bayer()
	if got := c.BlackAt(0, 0); got != 512 {
		t.Errorf("BlackAt(0,0) = %d, want 512", got)
	}
	if got := c.BlackAt(3, 1); got != 515 {
		t.Errorf("BlackAt(3,1) = %d, want 515", got)
	}

	// One level covers the frame.
	c.BlackLevel = []uint16{100}
	if got := c.BlackAt(3, 1); got != 100 {
		t.Errorf("BlackAt with a single level = %d, want 100", got)
	}
	// None recorded must read as zero, not panic: the ILCE-7 carries no black
	// level tag at all.
	c.BlackLevel = nil
	if got := c.BlackAt(3, 1); got != 0 {
		t.Errorf("BlackAt with no level = %d, want 0", got)
	}
}

func TestCFAPatternString(t *testing.T) {
	if got := bayer().PatternString(); got != "RGGB" {
		t.Errorf("PatternString = %q, want RGGB", got)
	}
}

// The Bayer test gates how a frame may be described to ASCOM and to FITS: only
// a 2x2 mosaic can be, and calling a 6x6 X-Trans frame Bayer would make a
// client debayer it with the wrong kernel.
func TestCFAIsBayerRejectsXTrans(t *testing.T) {
	if !bayer().IsBayer() {
		t.Error("a 2x2 mosaic does not report as Bayer")
	}
	xt := &CFA{PatternWidth: 6, PatternHeight: 6}
	if xt.IsBayer() {
		t.Error("a 6x6 X-Trans mosaic reports as Bayer, which ASCOM cannot describe")
	}
}

func TestCFAValidate(t *testing.T) {
	if err := bayer().Validate(); err != nil {
		t.Fatalf("a well-formed CFA does not validate: %v", err)
	}
	for _, tc := range []struct {
		name   string
		mangle func(*CFA)
	}{
		{"no dimensions", func(c *CFA) { c.Width = 0 }},
		{"wrong pixel count", func(c *CFA) { c.Pixels = c.Pixels[:3] }},
		{"no pattern", func(c *CFA) { c.PatternWidth = 0 }},
		{"pattern size mismatch", func(c *CFA) { c.Pattern = c.Pattern[:3] }},
		{"black level count", func(c *CFA) { c.BlackLevel = []uint16{1, 2} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := bayer()
			tc.mangle(c)
			if err := c.Validate(); err == nil {
				t.Error("Validate accepted an inconsistent CFA")
			}
		})
	}
}

func TestCFAColorString(t *testing.T) {
	for c, want := range map[CFAColor]string{CFARed: "R", CFAGreen: "G", CFABlue: "B", 9: "?"} {
		if got := c.String(); got != want {
			t.Errorf("CFAColor(%d) = %q, want %q", c, got, want)
		}
	}
}
