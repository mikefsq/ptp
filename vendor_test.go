package ptp

import "testing"

func TestRegisterAndLookup(t *testing.T) {
	v := Vendor{ID: 0xFFFE, Name: "Test", Models: map[uint16]string{0x0001: "T-1"}}
	Register(v)

	got, ok := Lookup(0xFFFE)
	if !ok {
		t.Fatal("registered vendor not found")
	}
	if got.Name != "Test" {
		t.Errorf("Name = %q, want %q", got.Name, "Test")
	}
	if _, ok := Lookup(0xFFFD); ok {
		t.Error("unregistered vendor was found")
	}
}

// A build must only offer cameras it can actually drive, so enumeration reads
// the registry rather than a fixed list of vendor IDs.
func TestRegisteredIsSortedAndReflectsLinkedDrivers(t *testing.T) {
	Register(Vendor{ID: 0xFF02, Name: "B"})
	Register(Vendor{ID: 0xFF01, Name: "A"})

	all := Registered()
	var seen []VendorID
	for _, v := range all {
		if v.ID == 0xFF01 || v.ID == 0xFF02 {
			seen = append(seen, v.ID)
		}
	}
	if len(seen) != 2 || seen[0] != 0xFF01 || seen[1] != 0xFF02 {
		t.Errorf("Registered() = %v, want ascending 0xFF01, 0xFF02", seen)
	}
}

func TestModelFallsBackToRawID(t *testing.T) {
	v := Vendor{ID: 0xFF03, Name: "Acme", Models: map[uint16]string{0x0010: "A-10"}}
	if got := v.Model(0x0010); got != "A-10" {
		t.Errorf("Model(0x0010) = %q, want %q", got, "A-10")
	}
	// An unknown product must still identify itself: a body can change its
	// product ID with its USB mode, and "unknown" is not a useful diagnostic.
	if got := v.Model(0x0099); got != "Acme device 0x0099" {
		t.Errorf("Model(0x0099) = %q, want the raw ID", got)
	}
}

func TestDuplicateRegistrationPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("registering the same vendor twice must panic: two packages " +
				"disagreeing about the same hardware is a build error, not a runtime one")
		}
	}()
	Register(Vendor{ID: 0xFF04, Name: "One"})
	Register(Vendor{ID: 0xFF04, Name: "Two"})
}
