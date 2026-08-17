package ptp

import (
	"fmt"
	"strings"
)

// PTP container types
const (
	ContainerCommand  uint16 = 1
	ContainerData     uint16 = 2
	ContainerResponse uint16 = 3
	ContainerEvent    uint16 = 4
)

// ContainerHeaderLen is the fixed PTP container header
const ContainerHeaderLen = 12

// OpCode is a PTP operation code
type OpCode uint16

// Standard PTP operations
const (
	OpGetDeviceInfo    OpCode = 0x1001
	OpOpenSession      OpCode = 0x1002
	OpCloseSession     OpCode = 0x1003
	OpGetStorageIDs    OpCode = 0x1004
	OpGetStorageInfo   OpCode = 0x1005
	OpGetNumObjects    OpCode = 0x1006
	OpGetObjectHandles OpCode = 0x1007
	OpGetObjectInfo    OpCode = 0x1008
	OpGetObject        OpCode = 0x1009
	OpGetThumb         OpCode = 0x100A
	OpDeleteObject     OpCode = 0x100B
	OpSendObjectInfo   OpCode = 0x100C
	OpSendObject       OpCode = 0x100D
	OpInitiateCapture  OpCode = 0x100E
	OpFormatStore      OpCode = 0x100F
	OpResetDevice      OpCode = 0x1010
	OpSetObjectProt    OpCode = 0x1012
	OpPowerDown        OpCode = 0x1013

	OpGetDevicePropDesc    OpCode = 0x1014
	OpGetDevicePropValue   OpCode = 0x1015
	OpSetDevicePropValue   OpCode = 0x1016
	OpResetDevicePropValue OpCode = 0x1017
	OpTerminateOpenCapture OpCode = 0x1018
	OpMoveObject           OpCode = 0x1019
	OpCopyObject           OpCode = 0x101A
	OpGetPartialObj        OpCode = 0x101B
	OpInitiateOpenCapture  OpCode = 0x101C
)

// MTP object-property operations, for bodies that expose them.
const (
	OpMTPGetObjectPropsSupported OpCode = 0x9801
	OpMTPGetObjectPropDesc       OpCode = 0x9802
	OpMTPGetObjectPropValue      OpCode = 0x9803
	OpMTPSetObjectPropValue      OpCode = 0x9804
	OpMTPGetObjPropList          OpCode = 0x9805
)

// Prop is a PTP device property code as it appears on the wire.
type Prop uint16

// DataType identifies how a property value is encoded. Values are the standard
// PTP datatype codes.
type DataType uint16

const (
	// TypeUndefined means the property carries no typed value
	TypeUndefined DataType = 0x0000

	TypeInt8   DataType = 0x0001
	TypeUint8  DataType = 0x0002
	TypeInt16  DataType = 0x0003
	TypeUint16 DataType = 0x0004
	TypeInt32  DataType = 0x0005
	TypeUint32 DataType = 0x0006
	TypeInt64  DataType = 0x0007
	TypeUint64 DataType = 0x0008
	TypeString DataType = 0xFFFF

	// Array variants set the 0x4000 bit.
	typeArrayBit DataType = 0x4000
)

// IsArray reports whether t is an array of its element type.
func (t DataType) IsArray() bool { return t&typeArrayBit != 0 }

// Elem returns the element type of an array type, or t itself.
func (t DataType) Elem() DataType { return t &^ typeArrayBit }

// FormFlag describes the shape of a property's valid-value set.
type FormFlag uint8

const (
	FormNone  FormFlag = 0x00 // no constraint
	FormRange FormFlag = 0x01 // min, max, step
	FormEnum  FormFlag = 0x02 // count followed by that many values
)

// String returns the property's name, or its hex code if it is not a standard
// one.
//
// This deliberately names ONLY standard properties. Vendor codes collide: 0xD001
// is film simulation on a Fujifilm body and something else entirely on a Sony.
func (p Prop) String() string {
	if n, ok := stdPropNames[p]; ok {
		return n
	}
	return fmt.Sprintf("Prop(0x%04X)", uint16(p))
}

// PropByName resolves a STANDARD property name, case-insensitively. Vendor
// packages layer their own names over this; it is here so that a name every
// camera shares does not have to be repeated in each of them.
func PropByName(name string) (Prop, bool) {
	for p, n := range stdPropNames {
		if strings.EqualFold(n, name) {
			return p, true
		}
	}
	return 0, false
}

// Standard PTP device properties the X-T5 reports and this driver uses by name.
const (
	PropImageSize       Prop = 0x5003
	PropWhiteBalance    Prop = 0x5005
	PropFNumber         Prop = 0x5007
	PropFocusMode       Prop = 0x500A
	PropMeteringMode    Prop = 0x500B
	PropExposureTime    Prop = 0x500D
	PropExposureProgram Prop = 0x500E

	// PropISO is the STANDARD ExposureIndex, and the one that actually carries
	// sensitivity. Do not confuse this with a vendor property of the same name
	PropISO Prop = 0x500F

	PropExposureBias Prop = 0x5010
	PropDateTime     Prop = 0x5011
	PropCaptureDelay Prop = 0x5012
	PropSharpness    Prop = 0x5015
)
