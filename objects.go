package ptp

import (
	"fmt"
	"time"
)

// Object format codes. Only the ones that turn up on a camera are listed.
const (
	FormatUndefined    uint16 = 0x3000
	FormatAssociation  uint16 = 0x3001 // a folder
	FormatScript       uint16 = 0x3002
	FormatDPOF         uint16 = 0x3006
	FormatUndefinedImg uint16 = 0x3800
	FormatEXIFJPEG     uint16 = 0x3801
	FormatTIFF         uint16 = 0x380D
	FormatMPO          uint16 = 0xB301

	// Vendor RAW formats. These sit here rather than in the vendor packages
	// because a caller sorting a directory listing has to recognise a RAW file
	// whoever made it, and the codes do not collide.
	FormatFujiRAF       uint16 = 0xB103 // Fujifilm RAW
	FormatSonyARWLegacy uint16 = 0xB101 // Sony RAW, older bodies
	FormatSonyARW       uint16 = 0xB310 // Sony RAW
)

// FormatName gives a readable name for an object format code.
func FormatName(f uint16) string {
	switch f {
	case FormatAssociation:
		return "folder"
	case FormatEXIFJPEG:
		return "JPEG"
	case FormatFujiRAF:
		return "RAF (Fujifilm RAW)"
	case FormatSonyARW, FormatSonyARWLegacy:
		return "ARW (Sony RAW)"
	case FormatMPO:
		return "MPO"
	case FormatTIFF:
		return "TIFF"
	case FormatScript:
		return "script"
	case FormatDPOF:
		return "DPOF"
	case FormatUndefined, FormatUndefinedImg:
		return "undefined"
	}
	return fmt.Sprintf("0x%04X", f)
}

// StorageInfo describes one storage volume (operation 0x1005).
type StorageInfo struct {
	StorageType        uint16
	FilesystemType     uint16
	AccessCapability   uint16
	MaxCapacity        uint64
	FreeSpaceInBytes   uint64
	FreeSpaceInImages  uint32
	StorageDescription string
	VolumeLabel        string
}

// Storage types (PTP). A camera's memory card is removable; a tethered
// Fujifilm body's capture buffer is fixed RAM, which is how the two are told
// apart without guessing from the name.
const (
	StorageFixedROM     uint16 = 0x0001
	StorageRemovableROM uint16 = 0x0002
	StorageFixedRAM     uint16 = 0x0003
	StorageRemovableRAM uint16 = 0x0004
)

// Access capabilities (PTP).
const (
	AccessReadWrite            uint16 = 0x0000
	AccessReadOnly             uint16 = 0x0001
	AccessReadOnlyWithDeletion uint16 = 0x0002
)

// IsRemovable reports whether the volume is a card rather than internal memory.
func (s *StorageInfo) IsRemovable() bool {
	return s.StorageType == StorageRemovableROM || s.StorageType == StorageRemovableRAM
}

func (s *StorageInfo) String() string {
	kind := fmt.Sprintf("type 0x%04X", s.StorageType)
	switch s.StorageType {
	case StorageFixedROM:
		kind = "fixed ROM"
	case StorageRemovableROM:
		kind = "removable ROM (card)"
	case StorageFixedRAM:
		kind = "fixed RAM (internal buffer)"
	case StorageRemovableRAM:
		kind = "removable RAM (card)"
	}
	access := fmt.Sprintf("access 0x%04X", s.AccessCapability)
	switch s.AccessCapability {
	case AccessReadWrite:
		access = "read-write"
	case AccessReadOnly:
		access = "read-only"
	case AccessReadOnlyWithDeletion:
		access = "read + delete, no write"
	}
	size := "capacity unknown"
	if s.MaxCapacity != ^uint64(0) {
		size = fmt.Sprintf("%d of %d bytes free", s.FreeSpaceInBytes, s.MaxCapacity)
	}
	name := s.StorageDescription
	if name == "" {
		name = s.VolumeLabel
	}
	return fmt.Sprintf("%q %s, %s, %s", name, kind, access, size)
}

// ObjectInfo describes one object on the camera (operation 0x1008).
type ObjectInfo struct {
	StorageID           uint32
	ObjectFormat        uint16
	ProtectionStatus    uint16
	CompressedSize      uint32
	ThumbFormat         uint16
	ThumbCompressedSize uint32
	ThumbPixWidth       uint32
	ThumbPixHeight      uint32
	ImagePixWidth       uint32
	ImagePixHeight      uint32
	ImageBitDepth       uint32
	ParentObject        uint32
	AssociationType     uint16
	AssociationDesc     uint32
	SequenceNumber      uint32
	Filename            string
	CaptureDate         string
	ModificationDate    string
	Keywords            string
}

// IsFolder reports whether the object is a directory rather than a file.
func (o *ObjectInfo) IsFolder() bool { return o.ObjectFormat == FormatAssociation }

// Captured parses the PTP capture date, which is "YYYYMMDDThhmmss" with an
// optional fractional part and trailing Z. A zero time means it did not parse.
func (o *ObjectInfo) Captured() time.Time {
	for _, layout := range []string{"20060102T150405", "20060102T150405.0", "20060102T150405Z"} {
		if t, err := time.Parse(layout, o.CaptureDate); err == nil {
			return t
		}
	}
	return time.Time{}
}

func (o *ObjectInfo) String() string {
	if o.IsFolder() {
		return fmt.Sprintf("%s/", o.Filename)
	}
	return fmt.Sprintf("%-16s %-14s %8d bytes  %dx%d",
		o.Filename, FormatName(o.ObjectFormat), o.CompressedSize, o.ImagePixWidth, o.ImagePixHeight)
}

// GetStorageIDs lists the camera's storage volumes (operation 0x1004).
func (s *Session) GetStorageIDs() ([]uint32, error) {
	data, _, err := s.Do(OpGetStorageIDs, nil, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	r := NewReader(data)
	n, err := r.U32()
	if err != nil {
		return nil, fmt.Errorf("ptp: storage IDs: count: %w", err)
	}
	if int(n) > r.Remaining()/4 {
		return nil, fmt.Errorf("ptp: storage IDs: %d ids: %w", n, ErrShortBlob)
	}
	out := make([]uint32, n)
	for i := range out {
		if out[i], err = r.U32(); err != nil {
			return nil, fmt.Errorf("ptp: storage IDs: %w", err)
		}
	}
	return out, nil
}

// GetStorageInfo describes one volume (operation 0x1005).
func (s *Session) GetStorageInfo(id uint32) (*StorageInfo, error) {
	data, _, err := s.Do(OpGetStorageInfo, []uint32{id}, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	r := NewReader(data)
	si := &StorageInfo{}
	if si.StorageType, err = r.U16(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: type: %w", err)
	}
	if si.FilesystemType, err = r.U16(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: filesystem: %w", err)
	}
	if si.AccessCapability, err = r.U16(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: access: %w", err)
	}
	if si.MaxCapacity, err = r.U64(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: capacity: %w", err)
	}
	if si.FreeSpaceInBytes, err = r.U64(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: free bytes: %w", err)
	}
	if si.FreeSpaceInImages, err = r.U32(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: free images: %w", err)
	}
	if si.StorageDescription, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: description: %w", err)
	}
	if si.VolumeLabel, err = r.Str(); err != nil {
		return nil, fmt.Errorf("ptp: storage info: volume label: %w", err)
	}
	return si, nil
}

// AllStorages is the wildcard storage ID, meaning "every volume".
const AllStorages uint32 = 0xFFFFFFFF

// AllFormats matches every object format in GetObjectHandles.
const AllFormats uint32 = 0

// RootFolder is the parent handle meaning the storage root.
const RootFolder uint32 = 0xFFFFFFFF

// GetObjectHandles lists object handles (operation 0x1007). Pass AllStorages
// and AllFormats to list everything, and 0 as parent to recurse the whole tree.
func (s *Session) GetObjectHandles(storage, format, parent uint32) ([]uint32, error) {
	data, _, err := s.Do(OpGetObjectHandles, []uint32{storage, format, parent}, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	r := NewReader(data)
	n, err := r.U32()
	if err != nil {
		return nil, fmt.Errorf("ptp: object handles: count: %w", err)
	}
	if int(n) > r.Remaining()/4 {
		return nil, fmt.Errorf("ptp: object handles: %d handles: %w", n, ErrShortBlob)
	}
	out := make([]uint32, n)
	for i := range out {
		if out[i], err = r.U32(); err != nil {
			return nil, fmt.Errorf("ptp: object handles: %w", err)
		}
	}
	return out, nil
}

// GetObjectInfo describes one object (operation 0x1008).
func (s *Session) GetObjectInfo(handle uint32) (*ObjectInfo, error) {
	data, _, err := s.Do(OpGetObjectInfo, []uint32{handle}, nil, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	return ParseObjectInfo(data)
}

// ParseObjectInfo decodes an ObjectInfo dataset.
func ParseObjectInfo(b []byte) (*ObjectInfo, error) {
	r := NewReader(b)
	o := &ObjectInfo{}
	var err error
	u32 := func(dst *uint32, what string) bool {
		if err != nil {
			return false
		}
		if *dst, err = r.U32(); err != nil {
			err = fmt.Errorf("ptp: object info: %s: %w", what, err)
			return false
		}
		return true
	}
	u16 := func(dst *uint16, what string) bool {
		if err != nil {
			return false
		}
		if *dst, err = r.U16(); err != nil {
			err = fmt.Errorf("ptp: object info: %s: %w", what, err)
			return false
		}
		return true
	}
	str := func(dst *string, what string) bool {
		if err != nil {
			return false
		}
		if *dst, err = r.Str(); err != nil {
			err = fmt.Errorf("ptp: object info: %s: %w", what, err)
			return false
		}
		return true
	}

	u32(&o.StorageID, "storage id")
	u16(&o.ObjectFormat, "object format")
	u16(&o.ProtectionStatus, "protection status")
	u32(&o.CompressedSize, "compressed size")
	u16(&o.ThumbFormat, "thumb format")
	u32(&o.ThumbCompressedSize, "thumb size")
	u32(&o.ThumbPixWidth, "thumb width")
	u32(&o.ThumbPixHeight, "thumb height")
	u32(&o.ImagePixWidth, "image width")
	u32(&o.ImagePixHeight, "image height")
	u32(&o.ImageBitDepth, "bit depth")
	u32(&o.ParentObject, "parent")
	u16(&o.AssociationType, "association type")
	u32(&o.AssociationDesc, "association desc")
	u32(&o.SequenceNumber, "sequence number")
	str(&o.Filename, "filename")
	str(&o.CaptureDate, "capture date")
	str(&o.ModificationDate, "modification date")
	str(&o.Keywords, "keywords")
	if err != nil {
		return nil, err
	}
	return o, nil
}

// TransferTimeout bounds a single bulk read during an object download. It is
// per-transfer, not per-file, so it does not need to scale with file size.
const TransferTimeout = 20 * time.Second

// GetObject downloads an object whole (operation 0x1009).
func (s *Session) GetObject(handle uint32) ([]byte, error) {
	data, _, err := s.Do(OpGetObject, []uint32{handle}, nil, TransferTimeout)
	return data, err
}

// GetThumb downloads an object's embedded thumbnail (operation 0x100A).
func (s *Session) GetThumb(handle uint32) ([]byte, error) {
	data, _, err := s.Do(OpGetThumb, []uint32{handle}, nil, DefaultTimeout)
	return data, err
}

// GetPartialObject downloads a byte range (operation 0x101B). Useful for large
// RAW files, and for pulling just a header without the whole frame.
func (s *Session) GetPartialObject(handle, offset, count uint32) ([]byte, error) {
	data, _, err := s.Do(OpGetPartialObj, []uint32{handle, offset, count}, nil, CaptureTimeout)
	return data, err
}

// u64 reads a little-endian uint64.
// U64 reads a little-endian uint64.
func (r *Reader) U64() (uint64, error) {
	lo, err := r.U32()
	if err != nil {
		return 0, err
	}
	hi, err := r.U32()
	if err != nil {
		return 0, err
	}
	return uint64(hi)<<32 | uint64(lo), nil
}

// GetNumObjects counts objects without enumerating their handles
// (operation 0x1006). Pass AllStorages and AllFormats for a total.
func (s *Session) GetNumObjects(storage, format, parent uint32) (uint32, error) {
	_, params, err := s.Do(OpGetNumObjects, []uint32{storage, format, parent}, nil, DefaultTimeout)
	if err != nil {
		return 0, err
	}
	if len(params) == 0 {
		return 0, fmt.Errorf("ptp: GetNumObjects returned no response parameter")
	}
	return params[0], nil
}

// GetPropValue reads a standard device property's current value
// (operation 0x1015). The type must be known in advance, from the descriptor.
func (s *Session) GetPropValue(p Prop, t DataType) (Value, error) {
	data, _, err := s.Do(OpGetDevicePropValue, []uint32{uint32(p)}, nil, DefaultTimeout)
	if err != nil {
		return Value{}, err
	}
	r := NewReader(data)
	return r.Value(t)
}

// SetPropValue writes a standard device property (operation 0x1016).
func (s *Session) SetPropValue(p Prop, t DataType, v uint64) error {
	payload, err := EncodeValue(t, v)
	if err != nil {
		return err
	}
	_, _, err = s.Do(OpSetDevicePropValue, []uint32{uint32(p)}, payload, DefaultTimeout)
	return err
}

// SetPropString writes a string-typed standard device property.
func (s *Session) SetPropString(p Prop, v string) error {
	_, _, err := s.Do(OpSetDevicePropValue, []uint32{uint32(p)}, EncodeString(v), DefaultTimeout)
	return err
}
