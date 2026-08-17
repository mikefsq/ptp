package ptp

import (
	"fmt"
	"sort"
	"sync"
)

// VendorID is a USB vendor ID.
type VendorID uint16

// The camera vendors this module knows by name.
const (
	Canon    VendorID = 0x04A9
	Nikon    VendorID = 0x04B0
	Fujifilm VendorID = 0x04CB
	Sony     VendorID = 0x054C
)

// Vendor describes one manufacturer's cameras well enough to find and identify
// them on the bus. Vendor packages register themselves from an init function.
type Vendor struct {
	ID VendorID

	// Name is the manufacturer as it should appear to a user.
	Name string

	// Models maps USB product ID to a human-readable model.
	Models map[uint16]string
}

// Model names a device, falling back to the raw ID when the product is not one
// this package knows.
func (v Vendor) Model(pid uint16) string {
	if m, ok := v.Models[pid]; ok {
		return m
	}
	return fmt.Sprintf("%s device 0x%04X", v.Name, pid)
}

var (
	vendorMu sync.RWMutex
	vendors  = map[VendorID]Vendor{}
)

// Vendor packages call Register it from init.
func Register(v Vendor) {
	vendorMu.Lock()
	defer vendorMu.Unlock()
	if _, dup := vendors[v.ID]; dup {
		panic(fmt.Sprintf("ptp: vendor 0x%04X registered twice", uint16(v.ID)))
	}
	vendors[v.ID] = v
}

// Lookup returns the registered vendor for a USB vendor ID.
func Lookup(id VendorID) (Vendor, bool) {
	vendorMu.RLock()
	defer vendorMu.RUnlock()
	v, ok := vendors[id]
	return v, ok
}

// Registered lists every vendor with a driver linked into this binary.
func Registered() []Vendor {
	vendorMu.RLock()
	defer vendorMu.RUnlock()
	out := make([]Vendor, 0, len(vendors))
	for _, v := range vendors {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
