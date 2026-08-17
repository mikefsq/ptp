package fuji

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mikefsq/ptp"
)

// Generic access to every property a Fujifilm body exposes.
//
// The typed helpers elsewhere in this package cover the settings a capture
// sequence needs — shutter, aperture, ISO, quality, focus mode. This is for the
// rest: an X-T5 advertises 263 properties, and hand-writing an accessor for
// each would be a lot of code to say the same thing 263 times.

// PropByName finds a property by name, case-insensitively. It is for tools that
// take a property on the command line.
func PropByName(name string) (ptp.Prop, bool) {
	// Names established by experiment first, matching PropName's order.
	for p, n := range observedNames {
		if strings.EqualFold(n, name) {
			return p, true
		}
	}
	// The body's own plugin names next, for the same reason PropName prefers
	// them: they are authoritative for this model, and they cover properties
	// gphoto2 names differently or not at all.
	for p, n := range xt5Names {
		if strings.EqualFold(n, name) {
			return p, true
		}
	}
	for p, n := range propNames {
		if strings.EqualFold(n, name) {
			return p, true
		}
	}
	// Finally the standard PTP names — ExposureProgram and friends belong to
	// every camera, so neither vendor table carries them.
	return ptp.PropByName(name)
}

// Known returns what a real X-T5 advertised for a property, if it was captured.
//
// This is reference material for tooling and documentation. Do not validate
// against it — see SetProp for why.
func Known(p ptp.Prop) (PropInfo, bool) {
	for _, info := range XT5Props {
		if info.Code == p {
			return info, true
		}
	}
	return PropInfo{}, false
}

// GetProp reads any property, using the data type the camera declares rather
// than one the caller has to know.
//
// Getting the width wrong is a real failure mode, not a theoretical one: the
// X-T5 declares ExposureIndex as UINT16 where gphoto2 and a plain reading of
// the PTP spec both suggest UINT32, and the camera rejects the wrong width.
func (c *Camera) GetProp(p ptp.Prop) (uint64, error) {
	d, err := c.GetPropDesc(p)
	if err != nil {
		return 0, fmt.Errorf("fuji: reading %s: %w", PropName(p), err)
	}
	return d.Current, nil
}

// SetProp writes any property, validating the value against what the camera
// says it will accept, and using the type the camera declares.
//
// Validation reads a live descriptor rather than consulting the generated
// table, because what a body accepts depends on the state it is in. With the
// shutter dial on T an X-T5 offers 76 shutter speeds; with it on a marked
// position, exactly one — and a write of anything else is accepted and silently
// ignored. A descriptor read costs a few milliseconds and is always right,
// where the table is a snapshot of one session.
//
// A property advertising a single value is reported as camera-controlled, which
// is the useful thing to say: the fix is a dial on the body, not a different
// value.
func (c *Camera) SetProp(p ptp.Prop, v uint64) error {
	d, err := c.GetPropDesc(p)
	if err != nil {
		return fmt.Errorf("fuji: reading %s before writing it: %w", PropName(p), err)
	}
	if !d.Writable() {
		return fmt.Errorf("fuji: %s is read-only", PropName(p))
	}
	if err := checkValue(p, d, v); err != nil {
		return err
	}
	if d.Current == v {
		// A redundant write draws vendor error 0xA002, so a no-op write is
		// worse than doing nothing.
		return nil
	}
	if err := c.SetPropValue(p, d.Type, v); err != nil {
		return fmt.Errorf("fuji: setting %s to %d: %w", PropName(p), v, err)
	}
	return nil
}

// checkValue reports whether the camera will accept v for p, given a live
// descriptor.
func checkValue(p ptp.Prop, d *ptp.StdPropDesc, v uint64) error {
	switch {
	case len(d.Enum) == 1 && d.Enum[0] != v:
		return fmt.Errorf("fuji: %s is camera-controlled right now — it offers only %s, "+
			"so a write of %s would be accepted and ignored. A dial or the exposure "+
			"program decides this, not the host",
			PropName(p), ValueName(p, d.Enum[0]), ValueName(p, v))

	case len(d.Enum) > 1:
		for _, allowed := range d.Enum {
			if allowed == v {
				return nil
			}
		}
		return fmt.Errorf("fuji: %s does not accept %d; it offers %s",
			PropName(p), v, formatValues(d.Enum))

	case d.Form == ptp.FormRange:
		if v < d.Min || v > d.Max {
			return fmt.Errorf("fuji: %s does not accept %d; its range is %d..%d step %d",
				PropName(p), v, d.Min, d.Max, d.Step)
		}
	}
	return nil
}

// formatValues renders an advertised value set, abbreviating a long one.
func formatValues(vs []uint64) string {
	const max = 12
	var b strings.Builder
	for i, v := range vs {
		if i == max {
			fmt.Fprintf(&b, " ... (%d more)", len(vs)-max)
			break
		}
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%d", v)
	}
	return b.String()
}

// Settable reports whether the host can currently write a property, and if not,
// why. A body hands control to the host through its dials and its exposure
// program, so "no" is usually about the camera's physical state.
func (c *Camera) Settable(p ptp.Prop) (bool, string) {
	d, err := c.GetPropDesc(p)
	if err != nil {
		return false, fmt.Sprintf("cannot be read: %v", err)
	}
	if !d.Writable() {
		return false, "read-only"
	}
	if len(d.Enum) == 1 {
		return false, fmt.Sprintf("camera-controlled: only %s is offered",
			ValueName(p, d.Enum[0]))
	}
	return true, ""
}

// PropsMatching returns property names containing the given substring, for
// suggesting alternatives when a lookup fails. With 263 properties on one body,
// a near miss is much more likely than a name that does not exist at all.
func PropsMatching(frag string) []string {
	frag = strings.ToLower(frag)
	seen := map[string]bool{}
	var out []string
	add := func(n string) {
		if !seen[n] && strings.Contains(strings.ToLower(n), frag) {
			seen[n] = true
			out = append(out, n)
		}
	}
	for _, n := range xt5Names {
		add(n)
	}
	for _, n := range propNames {
		add(n)
	}
	for _, info := range XT5Props {
		if info.Name != "" {
			add(info.Name)
		}
	}
	sort.Strings(out)
	if len(out) > 8 {
		out = out[:8]
	}
	return out
}
