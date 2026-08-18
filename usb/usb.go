// The package comment is in doc.go.
package usb

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// DeviceInfo describes an attached camera before it is opened.
type DeviceInfo struct {
	VID, PID uint16

	// Location identifies the device well enough to reopen it: the IOKit
	// location ID on macOS, bus<<16|address on Linux. It is stable only while
	// the device stays plugged in — use Serial to recognise a body across
	// replugs.
	Location uint32

	Serial string // USB serial; the only stable per-body identifier
	Name   string // USB product name

	// IfaceClass is the USB interface class: 6 is still image (PTP), which is
	// what this package needs; 8 is mass storage, meaning the body's USB mode
	// needs changing on the camera itself.
	IfaceClass uint32

	// Attachment identifies this plugging-in of the body: a new value every
	// time it enumerates, where Serial and (on macOS) Location stay the same.
	// A probe comparing it against the open transport's Info().Attachment
	// tells continued presence from a replug that left the handle dead. macOS:
	// the IORegistry entry id; Linux: bus<<16|address, the same as Location; 0
	// where the platform offers none.
	Attachment uint64
}

// IsPTP reports whether the body is presenting a still-image interface. A body
// in mass storage or MTP mode enumerates fine but cannot be driven.
func (d DeviceInfo) IsPTP() bool { return d.IfaceClass == 6 }

// ModeName describes the USB mode the body is currently in.
func (d DeviceInfo) ModeName() string {
	switch d.IfaceClass {
	case 6:
		return "still image (PTP)"
	case 8:
		return "mass storage"
	case 0:
		return "unknown"
	}
	return fmt.Sprintf("USB class %d", d.IfaceClass)
}

// Vendor names the manufacturer, if a driver for it is linked into this binary.
func (d DeviceInfo) Vendor() string {
	if v, ok := ptp.Lookup(ptp.VendorID(d.VID)); ok {
		return v.Name
	}
	return fmt.Sprintf("vendor 0x%04X", d.VID)
}

// Model names the body, preferring what the vendor package knows over the USB
// product string, which is often generic — several Fujifilm bodies report only
// "USB PTP Camera".
func (d DeviceInfo) Model() string {
	if v, ok := ptp.Lookup(ptp.VendorID(d.VID)); ok {
		if m, known := v.Models[d.PID]; known {
			return m
		}
	}
	if d.Name != "" {
		return d.Name
	}
	return fmt.Sprintf("0x%04X", d.PID)
}

func (d DeviceInfo) String() string {
	return fmt.Sprintf("%s %s %04x:%04x serial=%s loc=0x%08x mode=%s",
		d.Vendor(), d.Model(), d.VID, d.PID, d.Serial, d.Location, d.ModeName())
}

// Enumerate lists attached cameras from every vendor with a driver linked into
// this binary. It does not open them, so it is safe to call while cameras are
// in use.
//
// A build that imports no vendor package finds nothing, which is deliberate:
// matching a camera this binary cannot drive helps nobody.
func Enumerate() ([]DeviceInfo, error) {
	reg := ptp.Registered()
	if len(reg) == 0 {
		return nil, fmt.Errorf("usb: no camera vendors registered — import a vendor " +
			"package such as github.com/mikefsq/ptp/fuji")
	}
	vids := make([]uint16, 0, len(reg))
	for _, v := range reg {
		vids = append(vids, uint16(v.ID))
	}
	return enumerate(vids)
}

// EnumerateVendor lists attached cameras from one vendor only.
func EnumerateVendor(vid ptp.VendorID) ([]DeviceInfo, error) {
	return enumerate([]uint16{uint16(vid)})
}

// Open opens the camera with the given serial number. An empty serial opens the
// only attached camera and fails if more than one is present: with two bodies
// on one host, being explicit is the point.
func Open(serial string) (ptp.Transport, error) {
	devs, err := Enumerate()
	if err != nil {
		return nil, err
	}
	return openFrom(devs, serial)
}

// OpenVendor is Open restricted to one vendor, for a host running bodies from
// more than one manufacturer.
func OpenVendor(vid ptp.VendorID, serial string) (ptp.Transport, error) {
	devs, err := EnumerateVendor(vid)
	if err != nil {
		return nil, err
	}
	return openFrom(devs, serial)
}

func openFrom(devs []DeviceInfo, serial string) (ptp.Transport, error) {
	if len(devs) == 0 {
		return nil, fmt.Errorf("usb: no camera found (is the body powered on, and its " +
			"USB mode set to tethered shooting rather than card reader?)")
	}

	var want DeviceInfo
	if serial == "" {
		if len(devs) > 1 {
			return nil, fmt.Errorf("usb: %d cameras attached (%s); pass a serial to choose one",
				len(devs), joinSerials(devs))
		}
		want = devs[0]
	} else {
		found := false
		for _, d := range devs {
			if d.Serial == serial {
				want, found = d, true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("usb: no camera with serial %q (attached: %s)",
				serial, joinSerials(devs))
		}
	}

	// Refuse before opening if the body is not presenting a still-image
	// interface. Enumeration already told us the class, and opening — let alone
	// seizing — a device whose volume is mounted risks disturbing the mount.
	if want.IfaceClass != 0 && !want.IsPTP() {
		return nil, fmt.Errorf("usb: %s is in %s mode; set the body's USB connection "+
			"mode to tethered shooting and reconnect", want.Model(), want.ModeName())
	}
	return openDevice(want)
}

func joinSerials(devs []DeviceInfo) string {
	s := ""
	for i, d := range devs {
		if i > 0 {
			s += ", "
		}
		if d.Serial == "" {
			s += fmt.Sprintf("%04x:%04x@0x%08x", d.VID, d.PID, d.Location)
		} else {
			s += d.Serial
		}
	}
	return s
}

// ms converts a timeout to whole milliseconds, which is what both backends
// take. A zero or negative timeout would mean "wait forever" to the kernel, so
// it is clamped to the smallest real wait instead.
func ms(d time.Duration) uint32 {
	if d <= 0 {
		return 1
	}
	return uint32(d / time.Millisecond)
}
