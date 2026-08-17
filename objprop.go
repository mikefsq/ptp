package ptp

import "fmt"

// MTP object property codes
type ObjProp uint16

const (
	OPCStorageID        ObjProp = 0xDC01
	OPCObjectFormat     ObjProp = 0xDC02
	OPCProtectionStatus ObjProp = 0xDC03
	OPCObjectSize       ObjProp = 0xDC04
	OPCAssociationType  ObjProp = 0xDC05
	OPCAssociationDesc  ObjProp = 0xDC06
	OPCObjectFileName   ObjProp = 0xDC07
	OPCDateCreated      ObjProp = 0xDC08
	OPCDateModified     ObjProp = 0xDC09
	OPCKeywords         ObjProp = 0xDC0A
	OPCParentObject     ObjProp = 0xDC0B
	OPCHidden           ObjProp = 0xDC0D
	OPCSystemObject     ObjProp = 0xDC0E
	OPCPersistentUID    ObjProp = 0xDC41
	OPCName             ObjProp = 0xDC44
	OPCArtist           ObjProp = 0xDC46
	OPCDescription      ObjProp = 0xDC48
	OPCCopyright        ObjProp = 0xDC4B
	OPCWidth            ObjProp = 0xDC87
	OPCHeight           ObjProp = 0xDC88
	OPCDuration         ObjProp = 0xDC89
)

var objPropNames = map[ObjProp]string{
	OPCStorageID: "StorageID", OPCObjectFormat: "ObjectFormat",
	OPCProtectionStatus: "ProtectionStatus", OPCObjectSize: "ObjectSize",
	OPCAssociationType: "AssociationType", OPCAssociationDesc: "AssociationDesc",
	OPCObjectFileName: "ObjectFileName", OPCDateCreated: "DateCreated",
	OPCDateModified: "DateModified", OPCKeywords: "Keywords",
	OPCParentObject: "ParentObject", OPCHidden: "Hidden",
	OPCSystemObject: "SystemObject", OPCPersistentUID: "PersistentUID",
	OPCName: "Name", OPCArtist: "Artist", OPCDescription: "Description",
	OPCCopyright: "Copyright", OPCWidth: "Width", OPCHeight: "Height",
	OPCDuration: "Duration",
}

func (p ObjProp) String() string {
	if n, ok := objPropNames[p]; ok {
		return n
	}
	return fmt.Sprintf("ObjProp(0x%04X)", uint16(p))
}

// ObjPropEntry is one row of a GetObjPropList reply.
type ObjPropEntry struct {
	Handle uint32
	Prop   ObjProp
	Value  Value
}

// GetObjectPropsSupported lists the metadata properties a given object format
// supports (operation 0x9801).
func (s *Session) GetObjectPropsSupported(format uint16) ([]ObjProp, error) {
	data, _, err := s.Do(OpMTPGetObjectPropsSupported, []uint32{uint32(format)}, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	r := NewReader(data)
	codes, err := r.U16Array()
	if err != nil {
		return nil, fmt.Errorf("ptp: object props supported: %w", err)
	}
	out := make([]ObjProp, len(codes))
	for i, c := range codes {
		out[i] = ObjProp(c)
	}
	return out, nil
}

// ObjPropDesc describes one object property (operation 0x9802).
type ObjPropDesc struct {
	Prop     ObjProp
	Type     DataType
	Writable bool
	Default  Value
	GroupUID uint32
	Form     FormFlag
}

// GetObjectPropDesc describes one property for one object format (0x9802).
func (s *Session) GetObjectPropDesc(p ObjProp, format uint16) (*ObjPropDesc, error) {
	data, _, err := s.Do(OpMTPGetObjectPropDesc, []uint32{uint32(p), uint32(format)}, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	r := NewReader(data)
	d := &ObjPropDesc{}
	code, err := r.U16()
	if err != nil {
		return nil, fmt.Errorf("ptp: object prop desc: code: %w", err)
	}
	typ, err := r.U16()
	if err != nil {
		return nil, fmt.Errorf("ptp: object prop desc: type: %w", err)
	}
	d.Prop, d.Type = ObjProp(code), DataType(typ)
	gs, err := r.U8()
	if err != nil {
		return nil, fmt.Errorf("ptp: object prop desc: get/set: %w", err)
	}
	d.Writable = gs == 1
	if d.Default, err = r.Value(d.Type); err != nil {
		return nil, fmt.Errorf("ptp: object prop desc: default: %w", err)
	}
	if d.GroupUID, err = r.U32(); err != nil {
		// The group code and form are optional in practice; a short reply here
		// is not an error worth failing the whole call over.
		return d, nil
	}
	if f, err := r.U8(); err == nil {
		d.Form = FormFlag(f)
	}
	return d, nil
}

// GetObjectPropValue reads one metadata property of one object (0x9803).
//
// The type must be known in advance, from GetObjectPropDesc or the format's
// supported list.
func (s *Session) GetObjectPropValue(handle uint32, p ObjProp, t DataType) (Value, error) {
	data, _, err := s.Do(OpMTPGetObjectPropValue, []uint32{handle, uint32(p)}, nil, DefaultTimeout)
	if err != nil {
		return Value{}, err
	}
	r := NewReader(data)
	return r.Value(t)
}

// GetObjPropList fetches one metadata property for many objects in a single
// transaction (operation 0x9805).
func (s *Session) GetObjPropList(handle uint32, format uint16, prop ObjProp, depth uint32) ([]ObjPropEntry, error) {
	data, _, err := s.Do(OpMTPGetObjPropList,
		[]uint32{handle, uint32(format), uint32(prop), 0, depth}, nil, CaptureTimeout)
	if err != nil {
		return nil, err
	}
	return ParseObjPropList(data)
}

// AllObjects covers every object in GetObjPropList.
const AllObjects uint32 = 0xFFFFFFFF

// AllProps fetches every property in GetObjPropList.
const AllProps ObjProp = 0xFFFF

// MaxDepth recurses the whole tree.
const MaxDepth uint32 = 0xFFFFFFFF

// ParseObjPropList decodes a GetObjPropList reply: a uint32 count, then that
// many {uint32 handle, uint16 prop, uint16 type, value} rows.
func ParseObjPropList(b []byte) ([]ObjPropEntry, error) {
	r := NewReader(b)
	n, err := r.U32()
	if err != nil {
		return nil, fmt.Errorf("ptp: object prop list: count: %w", err)
	}
	// Each row is at least 8 bytes, so a count beyond that is nonsense and must
	// not be used to size an allocation.
	if int(n) > r.Remaining()/8 {
		return nil, fmt.Errorf("ptp: object prop list: %d rows will not fit in %d bytes: %w",
			n, r.Remaining(), ErrShortBlob)
	}
	out := make([]ObjPropEntry, 0, n)
	for i := uint32(0); i < n; i++ {
		var e ObjPropEntry
		if e.Handle, err = r.U32(); err != nil {
			return out, fmt.Errorf("ptp: object prop list: row %d: handle: %w", i, err)
		}
		code, err := r.U16()
		if err != nil {
			return out, fmt.Errorf("ptp: object prop list: row %d: code: %w", i, err)
		}
		typ, err := r.U16()
		if err != nil {
			return out, fmt.Errorf("ptp: object prop list: row %d: type: %w", i, err)
		}
		e.Prop = ObjProp(code)
		if e.Value, err = r.Value(DataType(typ)); err != nil {
			return out, fmt.Errorf("ptp: object prop list: row %d (%s): %w", i, e.Prop, err)
		}
		out = append(out, e)
	}
	return out, nil
}

// ByHandle groups object property entries by object handle, which is how a
// caller usually wants them — one map entry per file.
func ByHandle(entries []ObjPropEntry) map[uint32]map[ObjProp]Value {
	out := make(map[uint32]map[ObjProp]Value)
	for _, e := range entries {
		m, ok := out[e.Handle]
		if !ok {
			m = make(map[ObjProp]Value)
			out[e.Handle] = m
		}
		m[e.Prop] = e.Value
	}
	return out
}

// BulkMetadata fetches the given properties for every object on the camera,
// using one GetObjPropList round trip per property rather than one per object.
func (s *Session) BulkMetadata(props []ObjProp) (map[uint32]map[ObjProp]Value, error) {
	out := make(map[uint32]map[ObjProp]Value)
	var firstErr error
	got := 0
	for _, p := range props {
		entries, err := s.GetObjPropList(AllObjects, 0, p, 0)
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("property %s: %w", p, err)
			}
			continue
		}
		for _, e := range entries {
			m, ok := out[e.Handle]
			if !ok {
				m = make(map[ObjProp]Value)
				out[e.Handle] = m
			}
			m[e.Prop] = e.Value
		}
		got++
	}
	if got == 0 && firstErr != nil {
		return out, fmt.Errorf("ptp: bulk metadata returned nothing: %w", firstErr)
	}
	return out, nil
}
