# Fujifilm compressed RAF

Notes on the compressed RAF modes. Everything below is **measured from sample
files**, and says so where it is inference instead.

## RAF file can use three modes

A body offers **Uncompressed**, **Lossless Compressed** and **Compressed**. The
last two are different coders, not two settings of one — a decoder written for
one will produce confident nonsense on the other.

| mode | status |
|---|---|
| Uncompressed | **done**
| Lossless Compressed | **done** |
| Compressed (lossy) | **not implemented** |

## Layout: transposed, and why

The 18 line buffers are stored slot-major, `[slot][buffer]`, NOT buffer-major.
Every predictor read stays inside one colour group, so a slot's buffers being
contiguous puts a sample's four reads within ~76 bytes instead of two rows a
kilobyte apart.

**The six passes cannot be reordered.** Pass 2's bits begin where pass 1's end
and the gradient state carries forward, so interleaving them by column
desynchronises immediately. 

## Strip header — verified

16 bytes at `StripOffsets` (2048) inside the CFA block, big-endian:

    0  uint16  signature 0x4953
    2  uint8   version              1
    3  uint8   raw type             16 = X-Trans, 0 = Bayer (GFX)
    4  uint8   bits per sample      14
    5  uint16  height               5196
    7  uint16  rounded width        8448
    9  uint16  width                7872
    11 uint16  block size           768   (columns per block)
    13 uint8   blocks in row        11
    14 uint16  total lines          866

Then `blocks_in_row` big-endian uint32 block sizes, then padding to a 16-byte
boundary — for 11 blocks the data starts at strip+0x40.

Two invariants hold on every sample seen, and are worth asserting:

    total_lines * 6 == height          866 * 6 == 5196
    blocks_in_row * block_size >= width, with the last block partial
                                       11 * 768 = 8448, last covers 192

`line_width = block_size * 2 / 3` = 512 for X-Trans (`block_size / 2` for
Bayer). A "line" is **six rows** — one X-Trans period — so a block is decoded
866 times, each pass emitting a 6 x 768 tile.

## The coding scheme

Blocks are independent: block *k* covers columns `k*768 .. k*768+767` for the
whole frame height, with its byte range from the size table. They decode
concurrently, one goroutine each.

Per sample: predict from already-decoded neighbours, then Golomb-Rice code the
difference.

- **Prediction** is a branching weighted average. Even positions read only the
  previous line (Rb above, Rc above-left, Rd above-right, Rf two lines up) and
  drop whichever neighbour is furthest from Rb; odd positions read both sides.
  The coefficients are in `sampleEven`/`sampleOdd`.
- **Buckets.** Two neighbour deltas quantise to nine levels each (-4..+4) and
  combine as `9*a + b`, giving 81 signed values; the sign picks the correction's
  direction and the magnitude indexes 41 adaptive slots. The 81 and the 41 are
  the same table before and after `abs`.
- **Adaptation.** Each bucket keeps a running sum of magnitudes and a count; the
  fixed-length width comes from their ratio, so the coder tracks local noise.
- **Codes.** Unary zeros terminated by a 1, then a fixed remainder.
- **Sign** is signed-magnitude (`zigzag`), NOT two's complement.
- **Failsafe.** A prefix at or past `maxBits-rawBits-1` means the value follows
  verbatim.
- **Quantiser breakpoints** are 0x12, 0x43, 0x114 — deliberately not powers of
  two. The table covers only +/-512 because everything past the last breakpoint
  saturates, which keeps it in L1.

**Order.** Six passes per line, each filling two buffers, even slots running
ahead of odd ones which start once the evens pass 8. Some slots are interpolated
rather than coded, on a positional rule that differs per pass.
