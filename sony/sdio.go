package sony

import (
	"fmt"

	"github.com/mikefsq/ptp"
)

// The rest of the Sony SDIO vendor surface, beyond the handshake and the
// property blob. None of it has been seen on a camera: the only body available
// during development was a NEX-6, which has no SDIO operations at all. Treat
// every function here as decoded but UNVERIFIED.

// SDIOOpenSession opens a vendor session (0x9210). Distinct from the standard
// PTP OpenSession, and used by newer bodies after the connect handshake.
func (c *Camera) SDIOOpenSession(sessionID uint32) error {
	_, _, err := c.Do(OpSDIOOpenSession, []uint32{sessionID}, nil, ptp.DefaultTimeout)
	return err
}

// GetSonyPropDesc reads one vendor property descriptor (0x9203).
//
// The reply uses the same entry layout as one row of the 0x9209 blob, so it is
// parsed with the same code.
func (c *Camera) GetSonyPropDesc(p Prop) (*DeviceProperty, error) {
	data, _, err := c.Do(OpGetDevicePropDesc, []uint32{uint32(p)}, nil, ptp.DefaultTimeout)
	if err != nil {
		return nil, err
	}
	d, err := parseOne(ptp.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("sony: vendor property descriptor for %s: %w", PropName(p), err)
	}
	return &d, nil
}

// GetSonyPropValue reads one vendor property's current value (0x9204). The type
// must be known in advance.
func (c *Camera) GetSonyPropValue(p Prop, t ptp.DataType) (ptp.Value, error) {
	data, _, err := c.Do(OpGetDevicePropertyValue, []uint32{uint32(p)}, nil, ptp.DefaultTimeout)
	if err != nil {
		return ptp.Value{}, err
	}
	return ptp.NewReader(data).Value(t)
}

// GetControlDeviceDesc describes a control code (0x9206) — which values it
// accepts, and whether it is currently actionable.
func (c *Camera) GetControlDeviceDesc(ctrl ControlCode) ([]byte, error) {
	data, _, err := c.Do(OpGetControlDeviceDesc, []uint32{uint32(ctrl)}, nil, ptp.DefaultTimeout)
	return data, err
}

// ContentsTransferMode selects how the camera delivers captured images.
type ContentsTransferMode uint32

const (
	// TransferModeCameraOnly writes to the card and leaves it there.
	TransferModeCameraOnly ContentsTransferMode = 0
	// TransferModeHost hands captures to the host as they are taken, which is
	// what a tethered observatory loop wants.
	TransferModeHost ContentsTransferMode = 1
)

// SetContentsTransferMode selects capture delivery (0x9212).
func (c *Camera) SetContentsTransferMode(m ContentsTransferMode) error {
	_, _, err := c.Do(OpSetContentsTransferMode, []uint32{uint32(m)}, nil, ptp.DefaultTimeout)
	return err
}

// GetPartialLargeObject reads a byte range of a large object using 64-bit
// offsets (0x9211).
//
// This exists because the standard GetPartialObject takes a 32-bit offset,
// which cannot address past 4 GB. Sony splits the offset across two parameters,
// low word first — inferred from the operation's name and parameter count, and
// unverified.
func (c *Camera) GetPartialLargeObject(handle uint32, offset uint64, count uint32) ([]byte, error) {
	data, _, err := c.Do(OpSDIOGetPartialLargeObject,
		[]uint32{handle, uint32(offset & 0xFFFFFFFF), uint32(offset >> 32), count}, nil, ptp.CaptureTimeout)
	return data, err
}

// LensInfo is the reply to GetLensInformation (0x9223), returned raw. Its
// layout is not decoded: no body was available to produce a sample, and
// guessing at field offsets from nothing would be worse than handing back the
// bytes.
type LensInfo []byte

// GetLensInformation reads attached-lens details (0x9223).
func (c *Camera) GetLensInformation() (LensInfo, error) {
	data, _, err := c.Do(OpGetLensInformation, nil, nil, ptp.DefaultTimeout)
	return data, err
}

// GetDisplayStringList reads the camera's own human-readable strings for a
// property's values (0x9215) — how the body would label them on screen.
func (c *Camera) GetDisplayStringList(kind uint32) ([]byte, error) {
	data, _, err := c.Do(OpGetDisplayStringList, []uint32{kind}, nil, ptp.DefaultTimeout)
	return data, err
}

// OperationsResultsSupported asks whether the body reports per-operation
// results (0x922F).
func (c *Camera) OperationsResultsSupported(op ptp.OpCode) (bool, error) {
	_, params, err := c.Do(OpOperationsResultsSupported, []uint32{uint32(op)}, nil, ptp.DefaultTimeout)
	if err != nil {
		return false, err
	}
	return len(params) > 0 && params[0] != 0, nil
}

// LicenseInfo is one entry of the license list.
type LicenseInfo struct {
	RemainingHours uint32
	ID             string
}

// InfiniteLicense is the RemainingHours value meaning no expiry.
const InfiniteLicense uint32 = 0xFFFFFFFF

// GetLicenseInfoList reads installed licenses (0x924D). This is the only vendor
// operation Sony documents in the public SDK headers.
//
// The reply layout — a uint8 count, then per entry a uint32 remaining-hours,
// a uint8 id length and that many bytes — follows CrLicenseInfo in
// CrOperationCode.h. Unverified against a camera.
func (c *Camera) GetLicenseInfoList() ([]LicenseInfo, error) {
	data, _, err := c.Do(OpGetLicenseInfoList, nil, nil, ptp.DefaultTimeout)
	if err != nil {
		return nil, err
	}
	r := ptp.NewReader(data)
	n, err := r.U8()
	if err != nil {
		return nil, fmt.Errorf("sony: license list: count: %w", err)
	}
	out := make([]LicenseInfo, 0, n)
	for i := range int(n) {
		var li LicenseInfo
		if li.RemainingHours, err = r.U32(); err != nil {
			return out, fmt.Errorf("sony: license %d: remaining hours: %w", i, err)
		}
		idLen, err := r.U8()
		if err != nil {
			return out, fmt.Errorf("sony: license %d: id length: %w", i, err)
		}
		id, err := r.Bytes(int(idLen))
		if err != nil {
			return out, fmt.Errorf("sony: license %d: id of %d bytes: %w", i, idLen, err)
		}
		li.ID = string(id)
		out = append(out, li)
	}
	return out, nil
}
