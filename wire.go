package ptp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
	"unicode/utf16"
)

// Generic PTP wire decoding

// ErrShortBlob dataset ended mid-field.
var ErrShortBlob = errors.New("ptp: dataset truncated")

// TypeSize returns the encoded width of a scalar PTP data type.
func TypeSize(t DataType) (int, bool) {
	switch t.Elem() {
	case TypeInt8, TypeUint8:
		return 1, true
	case TypeInt16, TypeUint16:
		return 2, true
	case TypeInt32, TypeUint32:
		return 4, true
	case TypeInt64, TypeUint64:
		return 8, true
	}
	return 0, false
}

// signed reports whether the type is signed, so values are sign-extended.
func signed(t DataType) bool {
	switch t.Elem() {
	case TypeInt8, TypeInt16, TypeInt32, TypeInt64:
		return true
	}
	return false
}

// Reader walks a PTP dataset with bounds checks. Every read is checked, because
// a malformed or unexpected entry must fail the parse
type Reader struct {
	b   []byte
	off int
}

// NewReader reads a dataset from b.
func NewReader(b []byte) *Reader { return &Reader{b: b} }

// Remaining reports how many bytes are left unread.
func (r *Reader) Remaining() int { return len(r.b) - r.off }

// Bytes consumes n raw bytes.
func (r *Reader) Bytes(n int) ([]byte, error) {
	if n < 0 || r.Remaining() < n {
		return nil, ErrShortBlob
	}
	b := r.b[r.off : r.off+n]
	r.off += n
	return b, nil
}

// Peek16 reads the next uint16 WITHOUT consuming it
func (r *Reader) Peek16() (uint16, bool) {
	if r.Remaining() < 2 {
		return 0, false
	}
	return binary.LittleEndian.Uint16(r.b[r.off:]), true
}

// U8 reads one byte.
func (r *Reader) U8() (uint8, error) {
	if r.Remaining() < 1 {
		return 0, ErrShortBlob
	}
	v := r.b[r.off]
	r.off++
	return v, nil
}

// U16 reads a little-endian uint16.
func (r *Reader) U16() (uint16, error) {
	if r.Remaining() < 2 {
		return 0, ErrShortBlob
	}
	v := binary.LittleEndian.Uint16(r.b[r.off:])
	r.off += 2
	return v, nil
}

// U32 reads a little-endian uint32.
func (r *Reader) U32() (uint32, error) {
	if r.Remaining() < 4 {
		return 0, ErrShortBlob
	}
	v := binary.LittleEndian.Uint32(r.b[r.off:])
	r.off += 4
	return v, nil
}

// Scalar reads one value of type t, sign-extending if t is signed.
func (r *Reader) Scalar(t DataType) (uint64, error) {
	n, ok := TypeSize(t)
	if !ok {
		return 0, fmt.Errorf("ptp: unsupported data type 0x%04X", uint16(t))
	}
	if r.Remaining() < n {
		return 0, ErrShortBlob
	}
	var v uint64
	switch n {
	case 1:
		v = uint64(r.b[r.off])
	case 2:
		v = uint64(binary.LittleEndian.Uint16(r.b[r.off:]))
	case 4:
		v = uint64(binary.LittleEndian.Uint32(r.b[r.off:]))
	case 8:
		v = binary.LittleEndian.Uint64(r.b[r.off:])
	}
	r.off += n
	if signed(t) && n < 8 {
		shift := uint(64 - n*8)
		v = uint64(int64(v<<shift) >> shift)
	}
	return v, nil
}

// Str reads a PTP string
func (r *Reader) Str() (string, error) {
	n, err := r.U8()
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil
	}
	if r.Remaining() < int(n)*2 {
		return "", ErrShortBlob
	}
	u := make([]uint16, n)
	for i := range u {
		u[i] = binary.LittleEndian.Uint16(r.b[r.off+i*2:])
	}
	r.off += int(n) * 2
	// Sony NUL-terminates within the declared count; drop it if present.
	if len(u) > 0 && u[len(u)-1] == 0 {
		u = u[:len(u)-1]
	}
	return string(utf16.Decode(u)), nil
}

// Array reads a PTP array
func (r *Reader) Array(t DataType) ([]uint64, error) {
	n, err := r.U32()
	if err != nil {
		return nil, err
	}
	sz, ok := TypeSize(t)
	if !ok {
		return nil, fmt.Errorf("ptp: unsupported array element type 0x%04X", uint16(t))
	}
	if int(n) > r.Remaining()/sz {
		return nil, ErrShortBlob
	}
	out := make([]uint64, n)
	for i := range out {
		if out[i], err = r.Scalar(t); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// FindProp returns the property with the given code from a snapshot, or nil.
func FindProp(props []StdPropDesc, code Prop) *StdPropDesc {
	for i := range props {
		if props[i].Code == code {
			return &props[i]
		}
	}
	return nil
}

// Value holds a decoded property value of any PTP type.
type Value struct {
	Type DataType
	Num  uint64
	Str  string
	Arr  []uint64
}

func (v Value) String() string {
	switch {
	case v.Type == TypeString:
		return fmt.Sprintf("%q", v.Str)
	case v.Type.IsArray():
		return fmt.Sprintf("%v", v.Arr)
	default:
		return fmt.Sprintf("%d", int64(v.Num))
	}
}

// Value reads one value of type t
func (r *Reader) Value(t DataType) (Value, error) {
	v := Value{Type: t}
	var err error
	switch {
	case t == TypeString:
		v.Str, err = r.Str()
	case t.IsArray():
		v.Arr, err = r.Array(t)
	default:
		v.Num, err = r.Scalar(t)
	}
	return v, err
}

// EncodeString packs a PTP string
func EncodeString(s string) []byte {
	if s == "" {
		return []byte{0}
	}
	u := utf16.Encode([]rune(s))
	out := make([]byte, 0, 1+(len(u)+1)*2)
	out = append(out, byte(len(u)+1))
	for _, c := range u {
		out = binary.LittleEndian.AppendUint16(out, c)
	}
	return binary.LittleEndian.AppendUint16(out, 0)
}

// EncodeValue packs a scalar into a property payload.
func EncodeValue(t DataType, v uint64) ([]byte, error) {
	n, ok := TypeSize(t)
	if !ok {
		return nil, fmt.Errorf("ptp: cannot encode data type 0x%04X", uint16(t))
	}
	b := make([]byte, n)
	switch n {
	case 1:
		b[0] = byte(v)
	case 2:
		binary.LittleEndian.PutUint16(b, uint16(v))
	case 4:
		binary.LittleEndian.PutUint32(b, uint32(v))
	case 8:
		binary.LittleEndian.PutUint64(b, v)
	}
	return b, nil
}

// PTP carries a date-time as a string: "YYYYMMDDThhmmss", optionally with a
// fractional part and a Z or +/-hhmm zone. Bodies vary in how much of that they
// send, so parsing accepts the shorter forms.
const dateTimeLayout = "20060102T150405"

// ParseDateTime decodes a PTP date-time string.
func ParseDateTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("ptp: empty date-time")
	}
	for _, layout := range []string{
		dateTimeLayout + "Z0700",
		dateTimeLayout + "Z",
		dateTimeLayout + ".0",
		dateTimeLayout,
		"20060102T1504",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("ptp: cannot parse date-time %q", s)
}

// FormatDateTime encodes a time for a PTP date-time property.
func FormatDateTime(t time.Time) string { return t.Format(dateTimeLayout) }
