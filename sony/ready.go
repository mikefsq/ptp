package sony

import (
	"fmt"
	"time"
)

// Media slot writing state (CrMediaSlotWritingState). The camera reports this
// while it is flushing its buffer to the card.
const (
	WritingStateNotWriting uint64 = 0x01
	WritingStateWriting    uint64 = 0x02
)

// Media slot status (CrSlotStatus).
const (
	SlotOK                   uint64 = 0x0000
	SlotNoCard               uint64 = 0x0001
	SlotCardError            uint64 = 0x0002
	SlotRecognizingOrLocked  uint64 = 0x0003
	SlotDBError              uint64 = 0x0004
	SlotCardRecognizing      uint64 = 0x0005
	SlotCardLockedAndDBError uint64 = 0x0006
	SlotDBErrorNeedFormat    uint64 = 0x0007
	SlotCardErrorReadOnly    uint64 = 0x0008
)

// SlotStatusName gives a readable name for a slot status.
//
// It reads the same table ValueName uses, rather than repeating the list: two
// spellings of the same nine values would drift, and the one in the error
// message is the one a user acts on.
func SlotStatusName(v uint64) string {
	if n, ok := slotStatusNames[v]; ok {
		return n
	}
	return fmt.Sprintf("status 0x%04X", v)
}

// Slot identifies a card slot. Both the A7R V and A7R VI are dual-slot, and
// report the two independently.
type Slot int

const (
	Slot1 Slot = 1
	Slot2 Slot = 2
)

func (s Slot) String() string { return fmt.Sprintf("slot %d", int(s)) }

// props returns the three readiness property codes for this slot.
func (s Slot) props() (status, remaining, writing Prop) {
	if s == Slot2 {
		return PropMediaSLOT2Status, PropMediaSLOT2RemainingNumber, PropMediaSLOT2WritingState
	}
	return PropMediaSLOT1Status, PropMediaSLOT1RemainingNumber, PropMediaSLOT1WritingState
}

// SlotState is one card slot's condition.
type SlotState struct {
	Slot Slot

	// Reported is false when the camera said nothing about this slot — a
	// single-slot body, or a property this model does not expose. An
	// unreported slot is not a fault; it is simply absent.
	Reported bool

	Status uint64
	OK     bool // Status == SlotOK

	RemainingShots uint64
	HasRemaining   bool

	Writing    bool // still flushing buffer to this card
	HasWriting bool
}

// Usable reports that this slot has a healthy card with room on it.
func (s SlotState) Usable() bool {
	if !s.Reported || !s.OK {
		return false
	}
	if s.HasRemaining && s.RemainingShots == 0 {
		return false
	}
	return true
}

func (s SlotState) String() string {
	if !s.Reported {
		return fmt.Sprintf("%s: not reported", s.Slot)
	}
	out := fmt.Sprintf("%s: %s", s.Slot, SlotStatusName(s.Status))
	if s.HasRemaining {
		out += fmt.Sprintf(", %d shots left", s.RemainingShots)
	}
	if s.HasWriting {
		if s.Writing {
			out += ", writing"
		} else {
			out += ", idle"
		}
	}
	return out
}

// Readiness describes whether the camera can take another frame right now.
//
// It covers both slots, because either can be the one recording: a body with an
// empty slot 1 and a good card in slot 2 is perfectly able to shoot, and
// looking only at slot 1 would call that not ready.
type Readiness struct {
	Slot1, Slot2 SlotState
}

// Slots returns both slot states, for iterating.
func (r Readiness) Slots() [2]SlotState { return [2]SlotState{r.Slot1, r.Slot2} }

// Ready reports whether it is sensible to fire the shutter: at least one slot
// holds a healthy card with room.
//
// It deliberately does NOT require the buffer to be empty — a camera writes
// while it shoots, and gating each frame on an idle buffer would throttle a
// burst to card speed. Use Settled when you need the buffer actually flushed.
//
// A camera that reports no slot properties at all is treated as ready: absent
// is not the same as bad, and refusing to shoot because a body is quiet would
// be worse than trying.
func (r Readiness) Ready() bool {
	if !r.Slot1.Reported && !r.Slot2.Reported {
		return true
	}
	return r.Slot1.Usable() || r.Slot2.Usable()
}

// Settled reports that nothing is left to write on either slot. This is what to
// wait for before disconnecting, powering down, or formatting.
func (r Readiness) Settled() bool {
	for _, s := range r.Slots() {
		if s.HasWriting && s.Writing {
			return false
		}
	}
	return true
}

// Writing reports whether either slot is still flushing to card.
func (r Readiness) Writing() bool { return !r.Settled() }

// RemainingShots is the room on the fullest usable slot — how many more frames
// can be taken before the camera has nowhere to put them.
func (r Readiness) RemainingShots() (n uint64, known bool) {
	for _, s := range r.Slots() {
		if s.Usable() && s.HasRemaining {
			if s.RemainingShots > n {
				n = s.RemainingShots
			}
			known = true
		}
	}
	return n, known
}

func (r Readiness) String() string {
	head := "ready"
	if !r.Ready() {
		head = "NOT ready"
	}
	if n, ok := r.RemainingShots(); ok {
		head += fmt.Sprintf(" (%d shots)", n)
	}
	out := head
	for _, s := range r.Slots() {
		if s.Reported {
			out += "; " + s.String()
		}
	}
	return out
}

// ReadSlot pulls one slot's state out of a property snapshot.
func ReadSlot(props []DeviceProperty, slot Slot) SlotState {
	st := SlotState{Slot: slot}
	statusCode, remainingCode, writingCode := slot.props()

	if p := FindProp(props, statusCode); p != nil {
		st.Reported = true
		st.Status = p.Current
		st.OK = p.Current == SlotOK
	}
	if p := FindProp(props, remainingCode); p != nil {
		st.Reported = true
		st.HasRemaining = true
		st.RemainingShots = p.Current
	}
	if p := FindProp(props, writingCode); p != nil {
		st.Reported = true
		st.HasWriting = true
		st.Writing = p.Current == WritingStateWriting
	}
	return st
}

// ReadReadiness pulls the readiness picture for both slots out of a property
// snapshot.
func ReadReadiness(props []DeviceProperty) Readiness {
	return Readiness{
		Slot1: ReadSlot(props, Slot1),
		Slot2: ReadSlot(props, Slot2),
	}
}

// Readiness fetches a fresh property snapshot and reports readiness.
//
// Each call is a full 0x9209 round trip, which is not free. In a burst, read
// the snapshot once and use ReadReadiness rather than calling this per frame.
func (c *Camera) Readiness() (Readiness, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return Readiness{}, err
	}
	return ReadReadiness(props), nil
}

// WaitSettled blocks until neither slot is writing, or the timeout expires.
//
// Call this before disconnecting, powering down, or formatting — not between
// frames of a burst, where it would throttle shooting to card speed.
//
// A body that reports no writing state returns immediately: there is nothing to
// wait on, and blocking for the full timeout would be worse than proceeding.
func (c *Camera) WaitSettled(timeout, poll time.Duration) error {
	if poll <= 0 {
		poll = 250 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for {
		r, err := c.Readiness()
		if err != nil {
			return err
		}
		if r.Settled() {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("sony: camera still writing to card after %v (%s)", timeout, r)
		}
		time.Sleep(poll)
	}
}
