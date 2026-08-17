package sony

import (
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Settings that govern shot-to-shot delay, and a burst helper.
//
// The dominant cost in a continuous sequence is NOT the USB protocol — it is
// how much data each frame makes and where it has to go. In rough order of
// effect:
//
//  1. Drive mode. In a continuous mode the camera fires at its own rate for as
//     long as S2 is held, with no host involvement per frame. Issuing Capture
//     per frame instead costs two PTP transactions each, and any property write
//     between frames costs another.
//  2. File format and RAW compression. Lossless-compressed RAW is roughly half
//     the bytes of uncompressed, which is the buffer draining twice as fast.
//  3. Store destination. Writing to both host and card is the slowest option,
//     because every frame goes over USB as well as to the card.
//  4. Recording media. See RecordingMedia: the two slots do NOT share a burst.

// Drive modes (CrDriveMode). Continuous modes are what make a burst a burst.
const (
	DriveSingle              uint64 = 0x00000001
	DriveContinuousHi        uint64 = 0x00010001
	DriveContinuousHiPlus    uint64 = 0x00010002
	DriveContinuousHiLive    uint64 = 0x00010003
	DriveContinuousLo        uint64 = 0x00010004
	DriveContinuous          uint64 = 0x00010005
	DriveContinuousSpeedPrio uint64 = 0x00010006
	DriveContinuousMid       uint64 = 0x00010007
	DriveContinuousMidLive   uint64 = 0x00010008
	DriveContinuousLoLive    uint64 = 0x00010009
)

// Recording media routing for stills (CrRecordingMedia).
//
// Note what is NOT here: there is no mode that alternates frames between the
// two slots. Sony offers single-slot, the same image to both, or a split by
// file type. Only the split puts different data on different cards, and only
// then if two formats are being written.
const (
	// RecordingMediaSlot1 and Slot2 write everything to one card.
	RecordingMediaSlot1 uint64 = 0x0001
	RecordingMediaSlot2 uint64 = 0x0002

	// RecordingMediaSimultaneous writes the SAME image to both slots. This is
	// slower than one slot, not faster: it is twice the data, and the burst
	// runs at the slower card's pace.
	RecordingMediaSimultaneous uint64 = 0x0101

	// RecordingMediaSort splits by file type — RAW to one slot, JPEG to the
	// other. This is the only routing that spreads a single burst across both
	// cards, and it only helps when both formats are being recorded.
	RecordingMediaSort uint64 = 0x0102
)

// Still image store destination (CrStillImageStoreDestination).
const (
	StoreHostPC      uint64 = 0x0001 // over USB only, nothing written to card
	StoreMemoryCard  uint64 = 0x0002 // card only: fastest for a sustained burst
	StoreHostAndCard uint64 = 0x0003 // both: slowest, every frame also crosses USB
)

// RAW compression (CrRAWFileCompressionType). Fewer bytes per frame is a deeper
// buffer and a shorter drain.
const (
	RAWUncompressed  uint64 = 0x0000
	RAWCompressed    uint64 = 0x0001
	RAWLossless      uint64 = 0x0002
	RAWLosslessSmall uint64 = 0x0003
	RAWLosslessMed   uint64 = 0x0004
	RAWLosslessLarge uint64 = 0x0005
	RAWCompressedHQ  uint64 = 0x0006
)

// SetDriveMode sets the drive mode, e.g. DriveContinuousHi.
func (c *Camera) SetDriveMode(mode uint64) error {
	return c.SetProperty(PropDriveMode, ptp.TypeUint32, mode)
}

// SetRecordingMedia routes stills to a slot, both slots, or split by type.
func (c *Camera) SetRecordingMedia(mode uint64) error {
	return c.SetProperty(PropRecordingMedia, ptp.TypeUint16, mode)
}

// SetStoreDestination selects where captured stills go.
func (c *Camera) SetStoreDestination(dest uint64) error {
	return c.SetProperty(PropStillImageStoreDestination, ptp.TypeUint16, dest)
}

// SetRAWCompression sets the RAW compression type.
func (c *Camera) SetRAWCompression(t uint64) error {
	return c.SetProperty(PropRAWFileCompressionType, ptp.TypeUint16, t)
}

// Burst holds the shutter down for d, then releases it.
//
// In a continuous drive mode this is the fastest way to shoot a sequence: the
// camera fires at its own rate for as long as the button is held, and the host
// sends exactly two PTP transactions for the whole burst rather than two per
// frame. Frame timing is then the camera's, not the USB link's.
//
// In single drive mode it takes one frame regardless of d.
//
// The number of frames is whatever the camera manages in d, which depends on
// drive mode, shutter speed, buffer depth and card speed. Count the objects
// afterwards rather than assuming.
func (c *Camera) Burst(d time.Duration) error {
	if err := c.FullPress(); err != nil {
		c.FullRelease()
		return fmt.Errorf("sony: burst press: %w", err)
	}
	time.Sleep(d)
	if err := c.FullRelease(); err != nil {
		return fmt.Errorf("sony: burst release after %v: %w", d, err)
	}
	return nil
}

// BurstUntil holds the shutter down until stop is closed or the timeout
// expires, whichever comes first. The shutter is always released.
//
// Use this where the burst length is decided by something other than the clock
// — a keypress, an external trigger, a contact time.
func (c *Camera) BurstUntil(stop <-chan struct{}, timeout time.Duration) error {
	if err := c.FullPress(); err != nil {
		c.FullRelease()
		return fmt.Errorf("sony: burst press: %w", err)
	}
	select {
	case <-stop:
	case <-time.After(timeout):
	}
	if err := c.FullRelease(); err != nil {
		return fmt.Errorf("sony: burst release: %w", err)
	}
	return nil
}
