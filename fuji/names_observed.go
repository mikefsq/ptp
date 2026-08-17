package fuji

import "github.com/mikefsq/ptp"

// Names established by experiment on a real body, where the plugin binary has
// no accessor to recover one from.
//
// These are a separate layer from names_xt5_gen.go because the evidence is a
// different kind: the generated names come from the camera's own code, while
// these come from changing one thing on the camera and seeing which property
// moved. That is strong evidence — stronger than inference from a value's
// shape — but it is observation, not the vendor's own word, so it is kept
// apart and each entry records what was seen.
//
// The method is `fujiprobe -snapshot F` then `-diff F` with one control moved
// in between.
var observedNames = map[ptp.Prop]string{
	// 0xD033 is CommandDialStatus in the plugin, and 0xD034..0xD036 sit
	// immediately after it as read-only uint16s. They are per-dial position
	// indicators, and two are confirmed:
	//
	//   EV dial moved -3 -> C: 0xD034 went 2 -> 1
	//   ISO dial moved A -> C: 0xD035 went 2 -> 1
	//
	// So the encoding is DialOnCommand (1) when the dial is on C and the host
	// owns the setting, DialOnPosition (2) when a physical position does.
	0xD034: "EVDialStatus",
	0xD035: "ISODialStatus",
}

// The dial-status family, as confirmed by moving each control:
//
//	0xD033 CommandDialStatus  shutter dial: T -> 1/8s gave 1 -> 2
//	0xD034 EVDialStatus       EV dial:      -3 -> C   gave 2 -> 1
//	0xD035 ISODialStatus      ISO dial:     A -> C    gave 2 -> 1
//
// The plugin already named 0xD033 CommandDialStatus, and moving the SHUTTER
// dial is what changes it — so "command dial" means the shutter dial's T
// position, where the command dial sets the speed. 0xD036 did not move for any
// of the three and remains unidentified.

// Dial status values, for the *DialStatus properties.
//
// This is a direct answer to "does the host own this setting?", which is
// otherwise inferred from the property advertising a single value. The dial
// status says so without the inference.
const (
	DialOnCommand  uint64 = 1 // dial on C: the host can set it
	DialOnPosition uint64 = 2 // dial on a marked position: the camera owns it
)
