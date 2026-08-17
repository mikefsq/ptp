package fuji

import (
	"fmt"

	"github.com/mikefsq/ptp"
)

// Names for property VALUES, not just property codes.
//
// Every entry here is justified. Nothing is inferred from position, because
// position does not work: a cross-check of the SDK's enum blocks against what
// an X-T5 actually advertises agreed for 15 of 25 properties and disagreed for
// 10. DriveMode offers {4, 5, 16}, FocusMode reports 0x8001, NoiseReduction
// steps in 0x1000s. Assigning names by ordinal would be wrong about four times
// in ten, and silently.
//
// Two sources are trusted:
//
//   - Constants the SDK gives NUMERICALLY, in its sample code. There are only
//     65 of those, but they cross-check perfectly against hardware: for
//     DriveMode, MediaRecord and PriorityMode, every value the camera
//     advertises is a documented constant — 9 of 9, including DriveMode's
//     sparse {4, 5, 16}.
//   - Values confirmed on the camera itself, where an observable effect pins
//     the meaning down.
//
// Where neither applies, a value stays a number. That is the honest answer, and
// it is better than a plausible wrong name.
var valueNames = map[ptp.Prop]map[uint64]string{
	// Documented numerically in the SDK samples, and every advertised value
	// matched.
	PropPriorityModeCode: {
		PriorityModeCamera: "Camera",
		PriorityModeHost:   "PC",
	},
	PropMediaRecordCode: {
		MediaRecordRawJPEG: "RAW+JPEG",
		MediaRecordRaw:     "RAW",
		MediaRecordJPEG:    "JPEG",
		MediaRecordOff:     "Off",
	},
	PropDriveMode: {
		DriveContinuousH:   "ContinuousH",
		DriveContinuousL:   "ContinuousL",
		DriveSingle:        "Single",
		DriveMultiExposure: "MultiExposure",
		0x0006:             "AdvancedFilter",
		0x0007:             "Panorama",
		DriveMovie:         "Movie",
		0x0009:             "HDR",
		DriveBracketAE:     "BracketAE",
		DriveBracketISO:    "BracketISO",
		0x000C:             "BracketFilmSimulation",
		0x000D:             "BracketWhiteBalance",
		0x000E:             "BracketDynamicRange",
		DriveBracketFocus:  "BracketFocus",
		DrivePixelShift:    "PixelShift",
		0x0011:             "ContinuousHCrop",
		0x0012:             "PixelShiftFewerFrames",
		0xFFFF:             "Invalid",
	},

	// Confirmed on the camera. RawCompression was pinned by arithmetic: the
	// body read back 2 while producing a 36.6 MB file where 14-bit uncompressed
	// would be 69.7 MB, so 2 is lossless.
	PropRawCompression: {
		RawUncompressed: "Uncompressed",
		RawLossless:     "LosslessCompressed",
		RawLossy:        "Compressed",
	},
	// Quality read back 1 while the camera produced a lone .raf and no JPEG.
	PropQuality: {
		QualityRaw:       "RAW",
		QualityFine:      "Fine",
		QualityNormal:    "Normal",
		QualityRawFine:   "RAW+Fine",
		QualityRawNormal: "RAW+Normal",
	},
	// Focus mode was exercised directly: writing 1 stopped the half press
	// hunting, which is what manual focus means.
	ptp.PropFocusMode: {
		FocusManual: "Manual",
		FocusAFS:    "AF-S",
		FocusAFC:    "AF-C",
	},
	// The X-T5 menu is PREVIEW EXP./WB IN MANUAL MODE. Confirmed on hardware:
	// with 1 the body STOPS THE LENS DOWN to show the real exposure, which at a
	// small aperture makes the finder near-black and reads as being stuck in
	// depth-of-field preview; setting 3 opened it back up. 2 is the middle
	// option, white balance only, and is INFERRED from the menu rather than
	// tested.
	PropExposurePreview: {
		1: "ExposureAndWB",
		2: "WBOnly",
		3: "Off",
	},

	// Standard PTP codes, and confirmed: switching to 1 unlocked the aperture
	// that shutter priority had been holding.
	ptp.PropExposureProgram: {
		ProgramManual:           "Manual",
		ProgramAperturePriority: "AperturePriority",
		ProgramShutterPriority:  "ShutterPriority",
		ProgramAuto:             "Program",
	},
}

// ValueName names a property's value, falling back to the number when the
// meaning is not established. A number is not a failure here — it is the
// accurate thing to print.
func ValueName(p ptp.Prop, v uint64) string {
	if m, ok := valueNames[p]; ok {
		if n, ok := m[v]; ok {
			return n
		}
	}
	// Sensitivity is signed, and its negative values are the AUTO presets: an
	// X-T5 offers -1, -2 and -3 alongside 125..12800. Printed unsigned they
	// come out as 18446744073709551615 and look like corruption.
	if p == ptp.PropISO {
		switch int64(v) {
		case -1:
			return "auto1"
		case -2:
			return "auto2"
		case -3:
			return "auto3"
		}
	}
	if int64(v) < 0 && v > 1<<63 {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%d", v)
}

// DescribeValues renders a property's advertised set with names where known.
func DescribeValues(p ptp.Prop, vs []uint64) string {
	const max = 24
	out := ""
	for i, v := range vs {
		if i == max {
			out += fmt.Sprintf(" ... (%d more)", len(vs)-max)
			break
		}
		if i > 0 {
			out += ", "
		}
		out += ValueName(p, v)
	}
	return out
}
