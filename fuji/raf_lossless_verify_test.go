package fuji

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// Decode the real sample frames and compare every sample against libraw.
//
// This is the only test that can prove the lossless decoder: the codec is
// adaptive, so a subtle error does not produce a slightly wrong image, it
// desynchronises the bitstream and produces noise from that point on. Nothing
// short of an exact match means anything.
//
// The references are regenerated rather than committed — they are 82 MB each:
//
//	unprocessed_raw -T fuji/testdata/raf/xt5-lossless-lit.raf
func TestDecodeRAFMatchesLibRaw(t *testing.T) {
	for _, name := range []string{
		"xt5-lossless-lit.raf", // compressed, full tonal range
		"xt5-lossless.raf",     // compressed, nearly black
		"xt5-uncompressed.raf", // the native uncompressed path
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join("testdata", "raf", name)
			raf, err := os.ReadFile(path)
			if err != nil {
				t.Skipf("%s is not present", path)
			}
			ref, w, h, err := readRefTIFF(path + ".tiff")
			if err != nil {
				t.Skipf("no libraw reference: %v (run unprocessed_raw -T %s)", err, path)
			}

			cfa, err := DecodeRAF(raf)
			if err != nil {
				t.Fatalf("DecodeRAF: %v", err)
			}
			if cfa.Width != w || cfa.Height != h {
				t.Fatalf("decoded %dx%d, libraw says %dx%d", cfa.Width, cfa.Height, w, h)
			}
			for i, got := range cfa.Pixels {
				if got != ref[i] {
					t.Fatalf("sample %d (%d,%d) is %d, libraw says %d",
						i, i%w, i/w, got, ref[i])
				}
			}
		})
	}
}

// readRefTIFF reads the single-strip 16-bit TIFF unprocessed_raw writes.
func readRefTIFF(path string) ([]uint16, int, int, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, 0, err
	}
	le := binary.LittleEndian
	ifd := int(le.Uint32(d[4:]))
	n := int(le.Uint16(d[ifd:]))
	var w, h, off int
	for i := 0; i < n; i++ {
		e := ifd + 2 + i*12
		tag, typ := le.Uint16(d[e:]), le.Uint16(d[e+2:])
		v := int(le.Uint32(d[e+8:]))
		if typ == 3 {
			v = int(le.Uint16(d[e+8:]))
		}
		switch tag {
		case 256:
			w = v
		case 257:
			h = v
		case 273:
			off = v
		}
	}
	px := make([]uint16, w*h)
	for i := range px {
		px[i] = le.Uint16(d[off+i*2:])
	}
	return px, w, h, nil
}
