package sony

import (
	"fmt"

	"github.com/mikefsq/ptp"
)

// Names for property VALUES, not just property codes.
// Where no enum exists, a value stays a number.
var valueNames = map[Prop]map[uint64]string{
	// CrFileType. Note the ordering.
	PropFileType: {
		FileTypeNone:    "None",
		FileTypeJPEG:    "JPEG",
		FileTypeRAW:     "RAW",
		FileTypeRAWJPEG: "RAW+JPEG",
		FileTypeRAWHEIF: "RAW+HEIF",
		FileTypeHEIF:    "HEIF",
	},
	// CrImageQuality.
	PropStillImageQuality: {
		QualityUnknown:   "Unknown",
		QualityLight:     "Light",
		QualityStandard:  "Standard",
		QualityFine:      "Fine",
		QualityExtraFine: "ExtraFine",
	},
	// CrImageSize.
	PropImageSize: {
		ImageSizeLarge:  "L",
		ImageSizeMedium: "M",
		ImageSizeSmall:  "S",
		ImageSizeVGA:    "VGA",
	},
	// CrAspectRatioIndex.
	PropAspectRatio: {
		Aspect3to2:  "3:2",
		Aspect16to9: "16:9",
		Aspect4to3:  "4:3",
		Aspect1to1:  "1:1",
	},
	// CrRAWFileCompressionType. The lossless sizes are Sony's S/M/L variants,
	// which trade resolution for file size.
	PropRAWFileCompressionType: {
		RAWUncompressed:  "Uncompressed",
		RAWCompressed:    "Compressed",
		RAWLossless:      "Lossless",
		RAWLosslessSmall: "LosslessS",
		RAWLosslessMed:   "LosslessM",
		RAWLosslessLarge: "LosslessL",
		RAWCompressedHQ:  "CompressedHQ",
	},
	// CrStillImageStoreDestination.
	PropStillImageStoreDestination: {
		StoreHostPC:      "HostPC",
		StoreMemoryCard:  "MemoryCard",
		StoreHostAndCard: "HostPC+Card",
	},
	// CrRecordingMedia. Simultaneous and Sort jump to 0x0101, which is exactly
	// the kind of gap that defeats naming by ordinal.
	PropRecordingMedia: {
		RecordingMediaSlot1:        "Slot1",
		RecordingMediaSlot2:        "Slot2",
		RecordingMediaSimultaneous: "Simultaneous",
		RecordingMediaSort:         "SortByType",
	},
	// CrDriveMode. Single is 0x00000001 and the continuous family starts at
	// 0x00010001 — a 64K gap, again not ordinal.
	PropDriveMode: {
		DriveSingle:              "Single",
		DriveContinuousHi:        "ContinuousHi",
		DriveContinuousHiPlus:    "ContinuousHi+",
		DriveContinuousHiLive:    "ContinuousHiLive",
		DriveContinuousLo:        "ContinuousLo",
		DriveContinuous:          "Continuous",
		DriveContinuousSpeedPrio: "ContinuousSpeedPriority",
		DriveContinuousMid:       "ContinuousMid",
		DriveContinuousMidLive:   "ContinuousMidLive",
		DriveContinuousLoLive:    "ContinuousLoLive",
	},
	// CrFocusMode.
	PropFocusMode: {
		FocusManual: "MF",
		FocusAFS:    "AF-S",
		FocusAFC:    "AF-C",
		FocusAFA:    "AF-A",
		FocusAFD:    "AF-D",
		FocusDMF:    "DMF",
		FocusPF:     "PF",
	},
	// CrExposureProgram. Only the four the host can meaningfully drive are
	// named; the body offers a long tail of scene modes that take control away.
	PropExposureProgramMode: {
		ExposureManual:           "Manual",
		ExposureProgramAuto:      "Program",
		ExposureAperturePriority: "AperturePriority",
		ExposureShutterPriority:  "ShutterPriority",
	},
	// CrWriteCopyrightInfo. Off is 1, not 0 — zero is not a value it takes.
	PropWriteCopyrightInfo: {
		CopyrightInfoOff: "Off",
		CopyrightInfoOn:  "On",
	},
	// CrMediaSlotWritingState.
	PropMediaSLOT1WritingState: writingStateNames,
	PropMediaSLOT2WritingState: writingStateNames,
	// CrSlotStatus.
	PropMediaSLOT1Status: slotStatusNames,
	PropMediaSLOT2Status: slotStatusNames,
}

var writingStateNames = map[uint64]string{
	WritingStateNotWriting: "Idle",
	WritingStateWriting:    "Writing",
}

// slotStatusNames is deliberately written as PROSE, not in the identifier style
// the other tables use. It is the one table that reaches a user directly:
// SlotStatusName feeds "sony: camera is not ready to shoot: slot 1: no card",
// where "NoCard" would read as a symbol rather than a sentence. Sharing the one
// table across both callers matters more than a uniform register — two spellings
// of the same nine states would drift, and the drift would show up in an error
// message someone acts on.
var slotStatusNames = map[uint64]string{
	SlotOK:                   "OK",
	SlotNoCard:               "no card",
	SlotCardError:            "card error",
	SlotRecognizingOrLocked:  "recognizing, or locked with a DB error",
	SlotDBError:              "database error",
	SlotCardRecognizing:      "recognizing card",
	SlotCardLockedAndDBError: "card locked, database error",
	SlotDBErrorNeedFormat:    "database error, needs formatting",
	SlotCardErrorReadOnly:    "card error, read-only",
}

// boundEnums binds a property to an SDK enum where the names do NOT match, so
// the generator cannot bind it automatically.
//
// Each of these was checked by hand against the header. The mismatches are real:
// StillImageQuality's values live in CrImageQuality, not a CrStillImageQuality
// that does not exist. Getting a binding wrong puts confident, incorrect names
// on a camera's settings, which is worse than leaving them as numbers.
var boundEnums = map[Prop]string{
	PropStillImageQuality:   "CrImageQuality",
	PropAspectRatio:         "CrAspectRatioIndex",
	PropExposureProgramMode: "CrExposureProgram",
}

// enumName looks up a value through the generated enum tables.
func enumName(p Prop, v uint64) (string, bool) {
	name, ok := boundEnums[p]
	if !ok {
		name, ok = autoBoundEnums[p]
	}
	if !ok {
		return "", false
	}
	n, ok := enumValues[name][v]
	return n, ok
}

// ValueName names a property's value, falling back to the number when the
// meaning is not established. A number is not a failure — it is the accurate
// thing to print.
//
// Order matters. The curated table above wins because its entries were checked
// against the header by hand and carry the register a reader wants — prose for
// slot status, "RAW+JPEG" rather than "RawJpeg". The generated tables cover
// everything else, roughly 300 further properties.
func ValueName(p Prop, v uint64) string {
	if m, ok := valueNames[p]; ok {
		if n, ok := m[v]; ok {
			return n
		}
	}

	// The exposure triangle is encoded, not enumerated, so it is decoded rather
	// than looked up. Printing the raw wire value here is actively misleading:
	// a shutter speed of 1/1000 is 0x000103E8, which reads as 66536.
	switch p {
	case PropShutterSpeed:
		if uint32(v) == uint32(ShutterBulb) {
			return "Bulb"
		}
		if d, ok := DecodeShutterSpeed(v); ok {
			return d.String()
		}
	case PropFNumber:
		if f, ok := DecodeAperture(v); ok {
			return fmt.Sprintf("f/%.1f", f)
		}
	case PropIsoSensitivity:
		iso, mode, auto := DecodeISO(v)
		if auto {
			return "AUTO"
		}
		if mode != ISOModeNormal {
			return fmt.Sprintf("%d (mode %d)", iso, mode)
		}
		return fmt.Sprintf("%d", iso)
	case PropExposureBiasCompensation:
		return fmt.Sprintf("%+.1f EV", DecodeExposureCompensation(v))
	}

	// The generated enum tables: everything the SDK header defines, which is
	// roughly 300 more properties than the curated table above covers.
	if n, ok := enumName(p, v); ok {
		return n
	}

	// Signed properties printed unsigned come out as 18446744073709551615 and
	// look like corruption.
	if int64(v) < 0 {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%d", v)
}

// DescribeValues renders a property's advertised set with names where known.
func DescribeValues(p Prop, vs []uint64) string {
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

// Describe renders one property as the camera currently reports it: its name,
// its value, and whether the host may write it.
func Describe(d *DeviceProperty) string {
	if d == nil {
		return "<nil>"
	}
	val := ValueName(d.Code, d.Current)
	if d.Type == ptp.TypeString {
		val = fmt.Sprintf("%q", d.CurrentStr)
	}
	access := "read-only"
	if d.Writable() {
		access = "writable"
	}
	return fmt.Sprintf("%s(0x%04X) = %s [%s]", PropName(d.Code), uint16(d.Code), val, access)
}
