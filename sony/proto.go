package sony

import (
	"fmt"

	"github.com/mikefsq/ptp"
)

// Prop is a Sony device property code. It is an alias rather than a distinct
// type so that a Sony code passes straight to the standard operations on
// ptp.Session: the vendor surface and the standard one share a numeric space,
// and a body answers 0x5007 through either path.
//
// The consequence is that Prop cannot carry Sony methods — a package may not
// define methods on another package's type. Naming and per-model support are
// package functions here (PropName, Supported) instead, which matches how fuji
// does it and is the right shape anyway: vendor property codes COLLIDE, so only
// code that knows which body is attached can interpret one.
type Prop = ptp.Prop

// Sony vendor operations.
const (
	// OpSDIOConnect performs the three-phase handshake that unlocks the vendor
	// operations. It is issued three times with param1 = 1, 2, then 3.
	OpSDIOConnect ptp.OpCode = 0x9201

	// OpSDIOGetExtDeviceInfo returns the list of supported vendor property and
	// control codes. Must follow phase 2 of the connect handshake.
	OpSDIOGetExtDeviceInfo ptp.OpCode = 0x9202

	// OpGetDevicePropDesc and OpGetDevicePropertyValue are the VENDOR property
	// accessors. They are distinct from the standard 0x1014/0x1015, and the
	// descriptor they return has a different layout — see DeviceProperty.
	OpGetDevicePropDesc      ptp.OpCode = 0x9203
	OpGetDevicePropertyValue ptp.OpCode = 0x9204

	// OpSetControlDeviceA sets a stored setting (aperture, ISO, shutter).
	// gphoto2 calls this SDIO_SetExtDevicePropValue.
	OpSetControlDeviceA ptp.OpCode = 0x9205

	OpGetControlDeviceDesc ptp.OpCode = 0x9206

	// OpSetControlDeviceB issues a momentary control (shutter release, AF
	// trigger, zoom). gphoto2 calls this SDIO_ControlDevice.
	OpSetControlDeviceB ptp.OpCode = 0x9207

	// OpGetAllDevicePropData returns every property's current value, valid set
	// and enable flag in a single blob. This is the primary read path; the SDK
	// polls it rather than reading properties individually.
	OpGetAllDevicePropData ptp.OpCode = 0x9209

	OpSDIOOpenSession           ptp.OpCode = 0x9210
	OpSDIOGetPartialLargeObject ptp.OpCode = 0x9211

	// OpSetContentsTransferMode selects how captured images are delivered.
	OpSetContentsTransferMode ptp.OpCode = 0x9212

	OpGetDisplayStringList       ptp.OpCode = 0x9215
	OpGetLensInformation         ptp.OpCode = 0x9223
	OpOperationsResultsSupported ptp.OpCode = 0x922F
	OpGetLicenseInfoList         ptp.OpCode = 0x924D
)

// ControlCode is a Sony control ("device B") code — a momentary action rather
// than a stored setting.
//
// These are 32-bit, unlike device property codes: most sit in the 0xD2xx range
// that looks like a property code, but a few (Release, CancelContentsTransfer)
// are above 0xFFFF, so they cannot share the Prop type.
type ControlCode uint32

func (c ControlCode) String() string {
	if n, ok := controlNames[c]; ok {
		return n
	}
	return fmt.Sprintf("ControlCode(0x%08X)", uint32(c))
}

// Button values for button-type controls. Note these differ from the SDK's
// CrCommandParam (Up=0, Down=1), which belongs to its higher-level command API;
// these are what goes on the wire.
const (
	ButtonUp   uint64 = 0x0001
	ButtonDown uint64 = 0x0002
)

// PropName returns a Sony property's name.
//
// Sony's own table wins over the standard PTP names, because the same code can
// mean something different in the vendor surface than it does as a standard
// property.
func PropName(p Prop) string {
	if info, ok := PropTable[p]; ok {
		return info.Name
	}
	return p.String() // standard PTP names, then the raw code
}

// Supported reports whether a property is supported on the given model, per
// Sony's own per-model matrix. An unknown property reports false.
//
// This is what the SDK's documentation says, not what the body in front of you
// says. Ask the camera with Camera.Connect, whose ExtDeviceInfo lists what this
// particular body actually offers.
func Supported(p Prop, m Model) bool {
	info, ok := PropTable[p]
	if !ok {
		return false
	}
	switch m {
	case ModelA7RV:
		return info.A7RV
	case ModelA7RVI:
		return info.A7RVI
	}
	return false
}

// Model identifies a camera body. Only the bodies this driver targets are
// enumerated; the SDK's matrix covers 33.
type Model int

const (
	ModelUnknown Model = iota
	ModelA7RV          // ILCE-7RM5
	ModelA7RVI         // ILCE-7RM6
)

func (m Model) String() string {
	switch m {
	case ModelA7RV:
		return "ILCE-7RM5"
	case ModelA7RVI:
		return "ILCE-7RM6"
	}
	return "unknown"
}

// Wire returns the wire property code for an SDK property code, and whether the
// SDK code is known.
func (s SDKProp) Wire() (Prop, bool) {
	p, ok := sdkToWire[s]
	return p, ok
}
