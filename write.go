package ptp

import (
	"encoding/binary"
	"fmt"
)

// DeleteObject deletes one object from the camera (0x100B).
func (s *Session) DeleteObject(handle uint32, format uint16) error {
	_, _, err := s.Do(OpDeleteObject, []uint32{handle, uint32(format)}, nil, CaptureTimeout)
	return err
}

// DeleteAll for DeleteObject.
const DeleteAll uint32 = 0xFFFFFFFF

// FormatStore erases a storage volume (0x100F).
func (s *Session) FormatStore(storage uint32, format uint32) error {
	_, _, err := s.Do(OpFormatStore, []uint32{storage, format}, nil, CaptureTimeout)
	return err
}

// SetObjectProtection sets or clears an object's read-only flag (0x1012).
// status is 0 for unprotected, 1 for read-only.
func (s *Session) SetObjectProtection(handle uint32, status uint16) error {
	_, _, err := s.Do(OpSetObjectProt, []uint32{handle, uint32(status)}, nil, DefaultTimeout)
	return err
}

// MoveObject moves an object to another storage or folder (0x1019).
func (s *Session) MoveObject(handle, storage, parent uint32) error {
	_, _, err := s.Do(OpMoveObject, []uint32{handle, storage, parent}, nil, CaptureTimeout)
	return err
}

// CopyObject copies an object (0x101A). The response parameter is the new
// object's handle.
func (s *Session) CopyObject(handle, storage, parent uint32) (uint32, error) {
	_, params, err := s.Do(OpCopyObject, []uint32{handle, storage, parent}, nil, CaptureTimeout)
	if err != nil {
		return 0, err
	}
	if len(params) > 0 {
		return params[0], nil
	}
	return 0, nil
}

// SendObject uploads a file to the camera.
// Returns the storage ID, parent handle and object handle the camera assigned.
func (s *Session) SendObject(storage, parent uint32, oi *ObjectInfo, data []byte) (uint32, uint32, uint32, error) {
	oi.CompressedSize = uint32(len(data))
	_, params, err := s.Do(OpSendObjectInfo, []uint32{storage, parent}, EncodeObjectInfo(oi), DefaultTimeout)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ptp: sending object info: %w", err)
	}
	if len(params) < 3 {
		return 0, 0, 0, fmt.Errorf("ptp: SendObjectInfo returned %d parameters, want 3", len(params))
	}
	if _, _, err := s.Do(OpSendObject, nil, data, CaptureTimeout); err != nil {
		return params[0], params[1], params[2],
			fmt.Errorf("ptp: sending %d bytes of object data: %w", len(data), err)
	}
	return params[0], params[1], params[2], nil
}

// EncodeObjectInfo packs an ObjectInfo dataset for SendObjectInfo. It is the
// inverse of ParseObjectInfo.
func EncodeObjectInfo(o *ObjectInfo) []byte {
	var b []byte
	u32 := func(v uint32) { b = binary.LittleEndian.AppendUint32(b, v) }
	u16 := func(v uint16) { b = binary.LittleEndian.AppendUint16(b, v) }

	u32(o.StorageID)
	u16(o.ObjectFormat)
	u16(o.ProtectionStatus)
	u32(o.CompressedSize)
	u16(o.ThumbFormat)
	u32(o.ThumbCompressedSize)
	u32(o.ThumbPixWidth)
	u32(o.ThumbPixHeight)
	u32(o.ImagePixWidth)
	u32(o.ImagePixHeight)
	u32(o.ImageBitDepth)
	u32(o.ParentObject)
	u16(o.AssociationType)
	u32(o.AssociationDesc)
	u32(o.SequenceNumber)
	b = append(b, EncodeString(o.Filename)...)
	b = append(b, EncodeString(o.CaptureDate)...)
	b = append(b, EncodeString(o.ModificationDate)...)
	b = append(b, EncodeString(o.Keywords)...)
	return b
}

// PowerDown asks the camera to power off (0x1013). Sony bodies normally use the
// PowerOff control code instead; this is the standard PTP operation.
func (s *Session) PowerDown() error {
	_, _, err := s.Do(OpPowerDown, nil, nil, DefaultTimeout)
	return err
}

// ResetDevice resets the camera to its power-on state (0x1010). The session
// does not survive; reopen afterwards.
func (s *Session) ResetDevice() error {
	_, _, err := s.Do(OpResetDevice, nil, nil, DefaultTimeout)
	return err
}
