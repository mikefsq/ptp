//go:build !darwin && !linux

package usb

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/mikefsq/ptp"
)

// There is no transport for this platform yet. Windows would need WinUSB, and
// no camera here has been driven from it — the stub exists so the module still
// builds and its pure-Go parts stay testable everywhere.

// ErrNoBackend means this platform has no USB transport. It is a sentinel so
// cross-platform callers can branch on it rather than matching error text.
var ErrNoBackend = errors.New("usb: no USB backend for " + runtime.GOOS + " (macOS and Linux only)")

func enumerate(vids []uint16) ([]DeviceInfo, error) { return nil, ErrNoBackend }

func openDevice(want DeviceInfo) (ptp.Transport, error) {
	return nil, fmt.Errorf("opening %s: %w", want, ErrNoBackend)
}
