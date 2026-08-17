package ptp

import "time"

// PTP uses three endpoints on interface class 0x06 (Still Image):
// bulk OUT for commands and outbound data, bulk IN for inbound data and
// responses, and interrupt IN for asynchronous events.
type Transport interface {
	// BulkOut writes one transfer to the bulk-OUT endpoint. Implementations
	// must append a zero-length packet when len(p) is an exact multiple of the
	// endpoint's max packet size, so the device sees the transfer end.
	BulkOut(p []byte, timeout time.Duration) error

	// BulkIn reads one transfer from the bulk-IN endpoint and returns the byte
	// count. A short read terminates the transfer.
	BulkIn(p []byte, timeout time.Duration) (int, error)

	// InterruptIn reads one event packet. It returns ErrTimeout if none
	// arrives within the timeout; callers poll it and must treat a timeout as
	// normal, not as a failure.
	InterruptIn(p []byte, timeout time.Duration) (int, error)

	// MaxPacketSize reports the bulk endpoints' wMaxPacketSize, which decides
	// where zero-length packets are required.
	MaxPacketSize() int

	Close() error
}

// Resetter is an optional Transport capability: recovering a device left
// wedged, without a power cycle.
//
// A PTP transaction abandoned mid-transfer leaves the camera waiting for the
// rest of a data phase, after which it answers nothing and every request costs
// a full timeout. That is not rare on macOS — ptpcamerad is demand-launched
// whenever a still-image device is enumerated, SIP prevents removing it, and it
// can appear part-way through a session and abort a transfer in flight.
//
// A backend that cannot do this simply does not implement it, and Open then
// reports the failure as before.
type Resetter interface {
	// Reset clears the endpoints' halt state at both ends and asks the camera's
	// own PTP stack to return to idle. It does not re-enumerate the device, so
	// the Transport stays usable.
	Reset() error
}

// Timeouts. PTP itself sets no deadlines; these are what the SDK's adapter
// uses in practice.
const (
	// DefaultTimeout covers ordinary command/response round trips.
	DefaultTimeout = 5 * time.Second

	// CaptureTimeout covers operations that wait on the shutter, which can
	// legitimately take as long as the exposure.
	CaptureTimeout = 120 * time.Second
)
