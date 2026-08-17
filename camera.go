package ptp

import (
	"time"
)

// Capability interfaces.
//
// A vendor package builds its own camera type. The interfaces here are
// deliberately small.
//
// Ask for a capability with a type assertion:
//
//	if e, ok := cam.(ptp.ExposureControl); ok {
//		e.SetShutter(time.Second / 1000)
//	}
//

// Camera is the least a driver must provide.
type Camera interface {
	// Model names the body for logs and errors.
	Model() string

	// Close hands the camera back to its owner and releases it.
	Close() error
}

// Capturer takes photographs.
type Capturer interface {
	// Capture takes one frame and waits for the camera to finish.
	Capture(timeout time.Duration) error
}

// Exposure is the exposure triangle as a camera reports it.
type Exposure struct {
	Shutter  time.Duration
	Aperture float64
	ISO      uint64

	// Settable records, per component, whether the host can currently write it.
	// A camera hands control over through its physical dials and its exposure
	// program, so a false here usually means something on the body needs
	// moving, not that a different value would work.
	ShutterSettable, ApertureSettable, ISOSettable bool
}

// ExposureControl sets the exposure triangle.
type ExposureControl interface {
	SetShutter(d time.Duration) error
	SetAperture(f float64) error
	SetISO(iso uint32) error
	Exposure() (*Exposure, error)
}

// Downloader moves frames off the camera.
type Downloader interface {
	// NewFrames lists frames that have appeared since the last call.
	NewFrames() ([]uint32, error)

	// Download fetches one frame and its filename.
	Download(handle uint32) (data []byte, filename string, err error)

	// Delete removes a frame from the camera. On a body with a small volatile
	// store this is required.
	Delete(handle uint32) error
}

// RawDecoder turns a frame this camera produced into an undemosaiced readout.
// The bytes are what Download returned. Nothing is read from disk.
type RawDecoder interface {
	DecodeCFA(raw []byte) (*CFA, error)

	// SensorInfo describes what a capture WILL produce.
	SensorInfo() (*CFA, error)
}

// BufferReporter reports remaining space in the camera's capture buffer.
type BufferReporter interface {
	FreeBuffer() (uint64, error)
}

// FocusControl sets how the body focuses.
type FocusControl interface {
	SetManualFocus() error
}

// LiveViewer streams the camera's live preview.
type LiveViewer interface {
	// LiveFrame returns the most recent preview frame as JPEG, or nil if none
	// is waiting.
	LiveFrame() ([]byte, error)
}
