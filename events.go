package ptp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// EventCode is a PTP event code, delivered on the interrupt endpoint.
type EventCode uint16

const (
	EventCancelTransaction  EventCode = 0x4001
	EventObjectAdded        EventCode = 0x4002
	EventObjectRemoved      EventCode = 0x4003
	EventStoreAdded         EventCode = 0x4004
	EventStoreRemoved       EventCode = 0x4005
	EventDevicePropChanged  EventCode = 0x4006
	EventObjectInfoChanged  EventCode = 0x4007
	EventDeviceInfoChanged  EventCode = 0x4008
	EventRequestObjTransfer EventCode = 0x4009
	EventStoreFull          EventCode = 0x400A
	EventDeviceReset        EventCode = 0x400B
	EventStorageInfoChanged EventCode = 0x400C
	EventCaptureComplete    EventCode = 0x400D
	EventUnreportedStatus   EventCode = 0x400E

	// Vendor events. A body signals a new frame with its own ObjectAdded
	// rather than the standard 0x4002.
	//
	// These live here, unlike vendor PROPERTY codes, because the event space has
	// not collided in practice: each vendor took a different 0xC0xx block, and a
	// session polls one interrupt endpoint whoever made the camera. If two
	// vendors ever do disagree about a code, these move out to their packages.
	EventFujiPreviewAvailable EventCode = 0xC001
	EventFujiObjectAdded      EventCode = 0xC004

	// Sony's. The SDK polls property state rather than relying on the first,
	// so treat it as a hint that something changed, not as a complete report.
	EventSonyPropertyChanged EventCode = 0xC201
	EventSonyObjectAdded     EventCode = 0xC202
)

var eventNames = map[EventCode]string{
	EventCancelTransaction: "CancelTransaction", EventObjectAdded: "ObjectAdded",
	EventObjectRemoved: "ObjectRemoved", EventStoreAdded: "StoreAdded",
	EventStoreRemoved: "StoreRemoved", EventDevicePropChanged: "DevicePropChanged",
	EventObjectInfoChanged: "ObjectInfoChanged", EventDeviceInfoChanged: "DeviceInfoChanged",
	EventRequestObjTransfer: "RequestObjectTransfer", EventStoreFull: "StoreFull",
	EventDeviceReset: "DeviceReset", EventStorageInfoChanged: "StorageInfoChanged",
	EventCaptureComplete: "CaptureComplete", EventUnreportedStatus: "UnreportedStatus",
	EventFujiPreviewAvailable: "FujiPreviewAvailable", EventFujiObjectAdded: "FujiObjectAdded",
	EventSonyPropertyChanged: "SonyPropertyChanged", EventSonyObjectAdded: "SonyObjectAdded",
}

func (e EventCode) String() string {
	if n, ok := eventNames[e]; ok {
		return n
	}
	return fmt.Sprintf("Event(0x%04X)", uint16(e))
}

// Event is one asynchronous notification from the camera.
type Event struct {
	Code   EventCode
	TxID   uint32
	Params []uint32
}

func (e Event) String() string {
	if len(e.Params) == 0 {
		return e.Code.String()
	}
	return fmt.Sprintf("%s %v", e.Code, e.Params)
}

// PollEvent waits for one event on the interrupt endpoint.
//
// It returns ErrTimeout when nothing arrives in time, which is the normal case
// for an idle camera and must not be treated as a failure. Events use the same
// container framing as commands, with type 4.
func (s *Session) PollEvent(timeout time.Duration) (Event, error) {
	buf := make([]byte, 64)
	n, err := s.t.InterruptIn(buf, timeout)
	if err != nil {
		return Event{}, err
	}
	if n == 0 {
		return Event{}, ErrTimeout
	}
	if n < ContainerHeaderLen {
		return Event{}, fmt.Errorf("ptp: short event packet: %d bytes", n)
	}
	length := binary.LittleEndian.Uint32(buf[0:])
	typ := binary.LittleEndian.Uint16(buf[4:])
	if typ != ContainerEvent {
		return Event{}, fmt.Errorf("ptp: container type %d on the event endpoint, want %d",
			typ, ContainerEvent)
	}
	// Trust the smaller of the declared length and what actually arrived, so a
	// device overstating its length cannot walk us off the buffer.
	if int(length) < n {
		n = int(length)
	}
	ev := Event{
		Code: EventCode(binary.LittleEndian.Uint16(buf[6:])),
		TxID: binary.LittleEndian.Uint32(buf[8:]),
	}
	for off := ContainerHeaderLen; off+4 <= n; off += 4 {
		ev.Params = append(ev.Params, binary.LittleEndian.Uint32(buf[off:]))
	}
	return ev, nil
}

// WatchEvents delivers events to fn until the context-free stop channel closes
// or a non-timeout error occurs. Timeouts are swallowed, since an idle camera
// produces them continuously.
func (s *Session) WatchEvents(stop <-chan struct{}, fn func(Event)) error {
	for {
		select {
		case <-stop:
			return nil
		default:
		}
		ev, err := s.PollEvent(500 * time.Millisecond)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				continue
			}
			return err
		}
		fn(ev)
	}
}
