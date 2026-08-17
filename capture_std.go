package ptp

// Standard PTP capture operations.
//
// Not every camera has them. Sony bodies generally do NOT — they drive the
// shutter through vendor control codes, and a shot is a button press and
// release rather than an operation. Fujifilm uses InitiateCapture, but gated by
// a release-gesture property and a mandatory half press.
//
// Check DevInfo.Supports first. A body that lacks an operation refuses by
// stalling the bulk pipe.

// InitiateCapture triggers a standard capture (0x100E).
//
// Pass 0 for storage to let the camera choose, and 0 for format to use its
// current setting. The captured object is announced later by an ObjectAdded
// event, so this returning does not mean the image is ready.
func (s *Session) InitiateCapture(storage uint32, format uint16) error {
	_, _, err := s.Do(OpInitiateCapture, []uint32{storage, uint32(format)}, nil, CaptureTimeout)
	return err
}

// InitiateOpenCapture starts an open-ended capture — a bulb exposure (0x101C).
// It runs until TerminateOpenCapture. The response parameter is the transaction
// ID that TerminateOpenCapture must be given.
func (s *Session) InitiateOpenCapture(storage uint32, format uint16) (uint32, error) {
	_, params, err := s.Do(OpInitiateOpenCapture, []uint32{storage, uint32(format)}, nil, CaptureTimeout)
	if err != nil {
		return 0, err
	}
	if len(params) > 0 {
		return params[0], nil
	}
	return 0, nil
}

// TerminateOpenCapture ends an open capture started by InitiateOpenCapture
// (0x1018).
func (s *Session) TerminateOpenCapture(txID uint32) error {
	_, _, err := s.Do(OpTerminateOpenCapture, []uint32{txID}, nil, CaptureTimeout)
	return err
}
