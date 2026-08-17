package ptp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrTimeout is returned by a Transport when no data arrived in time. Callers
// polling for events must treat it as normal.
var ErrTimeout = errors.New("ptp: timeout")

// ErrStalled means the camera halted the bulk endpoint rather than answering.
var ErrStalled = errors.New("ptp: camera stalled the endpoint (operation likely unsupported)")

// ErrNotResponding means the camera has stopped answering entirely
// e.g. asleep, unplugged, or wedged.
var ErrNotResponding = errors.New("ptp: camera is not responding")

// ResponseCode is a PTP response code.
type ResponseCode uint16

// PTP response codes. Sony adds none of its own; a vendor operation that fails
// reports one of these.
const (
	RespUndefined                ResponseCode = 0x2000
	RespOK                       ResponseCode = 0x2001
	RespGeneralError             ResponseCode = 0x2002
	RespSessionNotOpen           ResponseCode = 0x2003
	RespInvalidTransactionID     ResponseCode = 0x2004
	RespOperationNotSupported    ResponseCode = 0x2005
	RespParameterNotSupported    ResponseCode = 0x2006
	RespIncompleteTransfer       ResponseCode = 0x2007
	RespInvalidStorageID         ResponseCode = 0x2008
	RespInvalidObjectHandle      ResponseCode = 0x2009
	RespDevicePropNotSupported   ResponseCode = 0x200A
	RespInvalidObjectFormat      ResponseCode = 0x200B
	RespStoreFull                ResponseCode = 0x200C
	RespObjectWriteProtected     ResponseCode = 0x200D
	RespStoreReadOnly            ResponseCode = 0x200E
	RespAccessDenied             ResponseCode = 0x200F
	RespNoThumbnailPresent       ResponseCode = 0x2010
	RespSelfTestFailed           ResponseCode = 0x2011
	RespPartialDeletion          ResponseCode = 0x2012
	RespStoreNotAvailable        ResponseCode = 0x2013
	RespSpecByFormatUnsupported  ResponseCode = 0x2014
	RespNoValidObjectInfo        ResponseCode = 0x2015
	RespInvalidCodeFormat        ResponseCode = 0x2016
	RespUnknownVendorCode        ResponseCode = 0x2017
	RespCaptureAlreadyTerminated ResponseCode = 0x2018
	RespDeviceBusy               ResponseCode = 0x2019
	RespInvalidParentObject      ResponseCode = 0x201A
	RespInvalidDevicePropFormat  ResponseCode = 0x201B
	RespInvalidDevicePropValue   ResponseCode = 0x201C
	RespInvalidParameter         ResponseCode = 0x201D
	RespSessionAlreadyOpen       ResponseCode = 0x201E
	RespTransactionCancelled     ResponseCode = 0x201F
	RespSpecOfDestUnsupported    ResponseCode = 0x2020
)

func (r ResponseCode) String() string {
	if n, ok := respNames[r]; ok {
		return n
	}
	return fmt.Sprintf("Response(0x%04X)", uint16(r))
}

var respNames = map[ResponseCode]string{
	RespUndefined: "Undefined", RespOK: "OK", RespGeneralError: "GeneralError",
	RespSessionNotOpen: "SessionNotOpen", RespInvalidTransactionID: "InvalidTransactionID",
	RespOperationNotSupported: "OperationNotSupported", RespParameterNotSupported: "ParameterNotSupported",
	RespIncompleteTransfer: "IncompleteTransfer", RespInvalidStorageID: "InvalidStorageID",
	RespInvalidObjectHandle: "InvalidObjectHandle", RespDevicePropNotSupported: "DevicePropNotSupported",
	RespInvalidObjectFormat: "InvalidObjectFormat", RespStoreFull: "StoreFull",
	RespObjectWriteProtected: "ObjectWriteProtected", RespStoreReadOnly: "StoreReadOnly",
	RespAccessDenied: "AccessDenied", RespNoThumbnailPresent: "NoThumbnailPresent",
	RespSelfTestFailed: "SelfTestFailed", RespPartialDeletion: "PartialDeletion",
	RespStoreNotAvailable: "StoreNotAvailable", RespSpecByFormatUnsupported: "SpecificationByFormatUnsupported",
	RespNoValidObjectInfo: "NoValidObjectInfo", RespInvalidCodeFormat: "InvalidCodeFormat",
	RespUnknownVendorCode: "UnknownVendorCode", RespCaptureAlreadyTerminated: "CaptureAlreadyTerminated",
	RespDeviceBusy: "DeviceBusy", RespInvalidParentObject: "InvalidParentObject",
	RespInvalidDevicePropFormat: "InvalidDevicePropFormat", RespInvalidDevicePropValue: "InvalidDevicePropValue",
	RespInvalidParameter: "InvalidParameter", RespSessionAlreadyOpen: "SessionAlreadyOpen",
	RespTransactionCancelled: "TransactionCancelled", RespSpecOfDestUnsupported: "SpecificationOfDestinationUnsupported",
}

// Error is a failed PTP transaction.
type Error struct {
	Op   OpCode
	Code ResponseCode

	// Name is the vendor's name for Code, when the session had a
	// ResponseNames hook and it recognised the code. Vendor response codes
	// share one numeric space and mean different things to different
	// manufacturers, so only the vendor package can name them.
	Name string
}

func (e *Error) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("ptp: operation 0x%04X failed: %s (0x%04X)",
			uint16(e.Op), e.Name, uint16(e.Code))
	}
	return fmt.Sprintf("ptp: operation 0x%04X failed: %s", uint16(e.Op), e.Code)
}

// container is a decoded PTP container header.
type container struct {
	length uint32
	typ    uint16
	code   uint16
	txID   uint32
}

// TraceEvent describes one completed PTP transaction.
type TraceEvent struct {
	Op       OpCode
	Params   []uint32
	TxID     uint32
	DataOut  []byte // data sent to the camera, if any
	DataIn   []byte // data received from the camera, if any
	RespCode ResponseCode
	Params2  []uint32 // response parameters
	Duration time.Duration
	Err      error
}

// Session is a PTP session over a Transport. PTP has one outstanding transaction
// per session, so every exported method holds the mutex for a whole
// command/data/response exchange.
type Session struct {
	t Transport

	// Trace, if set, is called after every transaction. It is called with the
	// session lock held, so it must not call back into the Session.
	Trace func(TraceEvent)

	// ResponseNames, if set, names a vendor response code. Without it a vendor
	// error prints as a bare number.
	ResponseNames func(ResponseCode) (string, bool)

	// Teardown, if set, runs just before the session closes, to hand the camera
	// back to its owner.
	Teardown func(Tx)

	mu   sync.Mutex
	txID uint32
	open bool

	// residue holds bytes read past the end of a container.
	//
	// A bulk read is not container-aligned. IOKit's ReadPipeTO completes on a
	// short packet, so when the camera has already queued both the data and the
	// response container the host controller hands them over in ONE transfer.
	// Discarding the overshoot loses the response, and the next read then blocks
	// forever waiting for something already thrown away — which surfaces as
	// "camera is not responding". Whatever runs past the current container is
	// the start of the next one, so it is kept here.
	residue []byte
}

// NewSession wraps a Transport. It does not talk to the device; call Open.
func NewSession(t Transport) *Session { return &Session{t: t} }

// Open starts PTP session 1. Sony bodies accept only one session at a time.
func (s *Session) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.open {
		return errors.New("ptp: session already open")
	}
	// The session ID is a transaction parameter, and the transaction counter
	// restarts at 1 for each new session.
	s.txID = 0
	_, _, err := s.transact(OpOpenSession, []uint32{1}, nil, DefaultTimeout)

	// A session left open by a previous host, usually macOS's ptpcamerad
	// makes the camera refuse a new session. Close the stale session and
	// try again rather than making the user power-cycle the body.
	var pe *Error
	if errors.As(err, &pe) && pe.Code == RespSessionAlreadyOpen {
		if _, _, cerr := s.transact(OpCloseSession, nil, nil, DefaultTimeout); cerr != nil {
			return fmt.Errorf("ptp: a session was already open and closing it failed: %w", cerr)
		}
		s.txID = 0
		_, _, err = s.transact(OpOpenSession, []uint32{1}, nil, DefaultTimeout)
	}
	// Resetting the endpoints and the camera's PTP stack may recover a body left
	// in an undetermined state by an abandoned transaction. Retry once before
	// telling the user to unplug the camera.
	if errors.Is(err, ErrNotResponding) || errors.Is(err, ErrTimeout) {
		if r, ok := s.t.(Resetter); ok {
			if rerr := r.Reset(); rerr == nil {
				s.resync()
				s.txID = 0
				_, _, err = s.transact(OpOpenSession, []uint32{1}, nil, DefaultTimeout)
			}
		}
	}
	if err != nil {
		return err
	}
	s.open = true
	return nil
}

// Close ends the PTP session. It does not close the Transport.
func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.open {
		return nil
	}
	// Calls the vendor's method to return camera to the user in a useable state.
	if s.Teardown != nil {
		s.Teardown(s.transact)
	}

	_, _, err := s.transact(OpCloseSession, nil, nil, DefaultTimeout)
	s.open = false
	return err
}

// responseName asks the vendor hook to name a code. It assumes s.mu is held.
func (s *Session) responseName(c ResponseCode) string {
	if s.ResponseNames == nil {
		return ""
	}
	if n, ok := s.ResponseNames(c); ok {
		return n
	}
	return ""
}

// Tx runs one PTP transaction with the session lock.
type Tx func(op OpCode, params []uint32, dataOut []byte, timeout time.Duration) ([]byte, []uint32, error)

// Do runs one PTP transaction. dataOut is sent in a data phase if non-nil.
// It returns the device's data phase, or nil if none.
func (s *Session) Do(op OpCode, params []uint32, dataOut []byte, timeout time.Duration) ([]byte, []uint32, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.transact(op, params, dataOut, timeout)
}

// transact assumes s.mu is held.
func (s *Session) transact(op OpCode, params []uint32, dataOut []byte, timeout time.Duration) (dataIn []byte, respParams []uint32, err error) {
	if len(params) > 5 {
		return nil, nil, fmt.Errorf("ptp: %d parameters, PTP allows at most 5", len(params))
	}
	s.txID++
	txID := s.txID

	if s.Trace != nil {
		start := time.Now()
		defer func() {
			ev := TraceEvent{
				Op: op, Params: params, TxID: txID,
				DataOut: dataOut, DataIn: dataIn, Params2: respParams,
				Duration: time.Since(start), Err: err,
				RespCode: RespOK,
			}
			var pe *Error
			if errors.As(err, &pe) {
				ev.RespCode = pe.Code
			}
			s.Trace(ev)
		}()
	}

	if err := s.writeContainer(ContainerCommand, uint16(op), txID, params, nil, timeout); err != nil {
		return nil, nil, fmt.Errorf("ptp: command phase: %w", err)
	}
	if dataOut != nil {
		if err := s.writeContainer(ContainerData, uint16(op), txID, nil, dataOut, timeout); err != nil {
			return nil, nil, fmt.Errorf("ptp: data-out phase: %w", err)
		}
	}

	// maxStale bounds how many late replies we will discard while hunting for
	// this transaction's answer.
	const maxStale = 8
	stale := 0
	for {
		hdr, body, rerr := s.readContainer(timeout)
		if rerr != nil {
			// dataIn is left as-is. if a data phase already arrived, the trace
			// keeps it, which is what tells us where the exchange went wrong.
			s.resync()
			return dataIn, nil, fmt.Errorf("ptp: response phase: %w", rerr)
		}
		// A container for an earlier transaction is a late reply that arrived
		// after we gave up on it. Accepting it would answer this request with
		// the previous one's data. Drop it and keep reading the stream for our
		//  response so the session recovers by itself.
		if hdr.txID != txID {
			if hdr.txID < txID && stale < maxStale {
				stale++
				continue
			}
			s.resync()
			return dataIn, nil, fmt.Errorf(
				"ptp: out of step with the camera: got a container for transaction %d "+
					"while running %d, and could not resynchronise", hdr.txID, txID)
		}
		switch hdr.typ {
		case ContainerData:
			// A device may only send one data phase, but read defensively:
			// a second would mean we lost sync with the transaction stream.
			if dataIn != nil {
				return dataIn, nil, errors.New("ptp: device sent two data containers in one transaction")
			}
			dataIn = body
		case ContainerResponse:
			if code := ResponseCode(hdr.code); code != RespOK {
				return dataIn, nil, &Error{Op: op, Code: code, Name: s.responseName(code)}
			}
			return dataIn, decodeParams(body), nil
		default:
			return dataIn, nil, fmt.Errorf("ptp: unexpected container type %d", hdr.typ)
		}
	}
}

func (s *Session) writeContainer(typ, code uint16, txID uint32, params []uint32, payload []byte, timeout time.Duration) error {
	n := ContainerHeaderLen + len(params)*4 + len(payload)
	buf := make([]byte, n)
	binary.LittleEndian.PutUint32(buf[0:], uint32(n))
	binary.LittleEndian.PutUint16(buf[4:], typ)
	binary.LittleEndian.PutUint16(buf[6:], code)
	binary.LittleEndian.PutUint32(buf[8:], txID)
	for i, p := range params {
		binary.LittleEndian.PutUint32(buf[ContainerHeaderLen+i*4:], p)
	}
	copy(buf[ContainerHeaderLen+len(params)*4:], payload)
	return s.t.BulkOut(buf, timeout)
}

// readContainer reads one complete container, reassembling it across as many
// bulk transfers as its header length calls for.
// readSome returns the next bytes of the bulk-IN stream, preferring anything
// left over from the previous container read.
func (s *Session) readSome(buf []byte, timeout time.Duration) (int, error) {
	if len(s.residue) > 0 {
		n := copy(buf, s.residue)
		s.residue = s.residue[n:]
		if len(s.residue) == 0 {
			s.residue = nil
		}
		return n, nil
	}
	return s.t.BulkIn(buf, timeout)
}

// resync drops any buffered stream state. Call it when the exchange is known to
// be out of step, so the damage stops at one transaction instead of every
// later read returning the previous one's answer.
func (s *Session) resync() { s.residue = nil }

func (s *Session) readContainer(timeout time.Duration) (container, []byte, error) {
	const chunk = 512 * 1024
	buf := make([]byte, chunk)

	n, err := s.readSome(buf, timeout)
	if err != nil {
		return container{}, nil, err
	}
	if n < ContainerHeaderLen {
		return container{}, nil, fmt.Errorf("ptp: short container header: %d bytes", n)
	}
	hdr := container{
		length: binary.LittleEndian.Uint32(buf[0:]),
		typ:    binary.LittleEndian.Uint16(buf[4:]),
		code:   binary.LittleEndian.Uint16(buf[6:]),
		txID:   binary.LittleEndian.Uint32(buf[8:]),
	}
	if hdr.length < ContainerHeaderLen {
		return container{}, nil, fmt.Errorf("ptp: container length %d below header size", hdr.length)
	}
	// PTP allows 0xFFFFFFFF to mean "length unknown", and a desynchronised read
	// can produce any garbage here. Either way, sizing an allocation from it
	// directly would try to reserve gigabytes.
	const maxContainer = 512 << 20
	if hdr.length > maxContainer {
		return container{}, nil, fmt.Errorf(
			"ptp: container declares %d bytes, above the %d-byte sanity limit "+
				"(lost sync with the device, or an unknown-length transfer)", hdr.length, maxContainer)
	}

	want := int(hdr.length) - ContainerHeaderLen
	body := make([]byte, want)

	// Read the bulk of the transfer STRAIGHT into body.
	have := copy(body, buf[ContainerHeaderLen:n])
	var over []byte
	if n-ContainerHeaderLen > want {
		// The whole container, and more, arrived in the first read.
		over = append(over, buf[ContainerHeaderLen+want:n]...)
	}
	for have < want {
		remaining := want - have
		if remaining >= chunk {
			// Room to spare: a transfer cannot overrun the space left, so it is
			// safe to land directly in body.
			n, err = s.readSome(body[have:have+chunk], timeout)
		} else {
			// Near the end the device may hand over this container's tail and
			// the head of the next one together. Stage that read so the excess
			// survives instead of being clipped by a too-small destination.
			n, err = s.readSome(buf, timeout)
			if n > 0 {
				take := n
				if take > remaining {
					take = remaining
				}
				copy(body[have:], buf[:take])
				if n > take {
					over = append(over, buf[take:n]...)
				}
				n = take
			}
		}
		if err != nil {
			return container{}, nil, fmt.Errorf("reading container body at %d/%d bytes: %w",
				have+ContainerHeaderLen, hdr.length, err)
		}
		if n == 0 {
			return container{}, nil, fmt.Errorf("ptp: container truncated at %d/%d bytes",
				have+ContainerHeaderLen, hdr.length)
		}
		have += n
	}

	// A container header distinguishes between the start of the NEXT container
	// and padding when a device rounds the final packet up
	if len(over) > 0 && looksLikeContainer(over) {
		s.residue = append(over, s.residue...)
	}
	return hdr, body, nil
}

// looksLikeContainer reports whether b plausibly starts a PTP container, as
// opposed to a device's trailing padding.
func looksLikeContainer(b []byte) bool {
	for _, c := range b {
		if c != 0 {
			goto nonzero
		}
	}
	return false // all zeros: padding
nonzero:
	if len(b) < 6 {
		// Too short to check, but non-zero, so it is not padding. Keep it and
		// let the next read complete the header.
		return true
	}
	if binary.LittleEndian.Uint32(b) < ContainerHeaderLen {
		return false
	}
	switch binary.LittleEndian.Uint16(b[4:]) {
	case ContainerCommand, ContainerData, ContainerResponse, ContainerEvent:
		return true
	}
	return false
}

func decodeParams(b []byte) []uint32 {
	if len(b) < 4 {
		return nil
	}
	out := make([]uint32, 0, len(b)/4)
	for i := 0; i+4 <= len(b); i += 4 {
		out = append(out, binary.LittleEndian.Uint32(b[i:]))
	}
	return out
}
