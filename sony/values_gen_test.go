package sony

import "testing"

// The curated table and the generated enums must agree wherever they overlap.
//
// This is the check that would have caught the CrFileType transcription error:
// the hand-written constants said RAW=1, the header says RAW=2, and a value
// named in one but absent from the other is exactly that shape of mistake.
func TestCuratedNamesAgreeWithGeneratedEnums(t *testing.T) {
	for p, curated := range valueNames {
		enum, ok := boundEnums[p]
		if !ok {
			enum, ok = autoBoundEnums[p]
		}
		if !ok {
			continue // no generated enum covers this property
		}
		gen := enumValues[enum]
		if len(gen) == 0 {
			t.Errorf("%s is bound to %q, which the generated tables do not define",
				PropName(p), enum)
			continue
		}
		for v := range curated {
			if _, inGen := gen[v]; !inGen {
				t.Errorf("%s: the curated table names value %d, but %s does not define it — "+
					"one of the two is wrong", PropName(p), v, enum)
			}
		}
	}
}

// Every hand-written binding must name an enum that exists. A typo here would
// silently disable naming for that property rather than fail.
func TestHandBoundEnumsExist(t *testing.T) {
	for p, enum := range boundEnums {
		if _, ok := enumValues[enum]; !ok {
			t.Errorf("%s is bound to %q, which is not in the generated tables",
				PropName(p), enum)
		}
		if auto, ok := autoBoundEnums[p]; ok && auto == enum {
			t.Errorf("%s is bound by hand to %q, which the generator already binds automatically",
				PropName(p), enum)
		}
	}
}

// A floor on coverage, so a regeneration that silently produced far less is
// noticed. The generated tables reach roughly 300 properties more than the
// curated table alone.
func TestGeneratedCoverageFloor(t *testing.T) {
	reachable := map[Prop]bool{}
	for p := range valueNames {
		reachable[p] = true
	}
	for p := range autoBoundEnums {
		reachable[p] = true
	}
	for p := range boundEnums {
		reachable[p] = true
	}
	if len(reachable) < 300 {
		t.Errorf("only %d of %d properties have named values; expected 300+ — "+
			"has values_gen.go been regenerated?", len(reachable), len(PropTable))
	}
	if len(enumValues) < 500 {
		t.Errorf("only %d value enums; the SDK header defines about 515", len(enumValues))
	}
}

// The generated property table must still describe a whole camera. A parse that
// silently produced a fraction of the rows would otherwise look like success.
func TestGeneratedPropTableIsComplete(t *testing.T) {
	if len(PropTable) != 810 {
		t.Errorf("PropTable has %d entries, want 810", len(PropTable))
	}
	// Spot-check the first rows of the SDK's own table, which are the ones
	// pinned in sony/PROTOCOL.md.
	for _, tc := range []struct {
		p   Prop
		sdk SDKProp
	}{
		{PropFNumber, 0x0100},
		{PropExposureBiasCompensation, 0x0101},
		{PropFlashCompensation, 0x0102},
		{PropShutterSpeed, 0x0103},
		{PropIsoSensitivity, 0x0104},
	} {
		info, ok := PropTable[tc.p]
		if !ok {
			t.Errorf("0x%04X missing from PropTable", uint16(tc.p))
			continue
		}
		if info.SDK != tc.sdk {
			t.Errorf("%s SDK code = 0x%04X, want 0x%04X", PropName(tc.p), info.SDK, tc.sdk)
		}
	}
}
