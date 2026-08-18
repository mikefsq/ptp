package usb

import (
	"context"
	"errors"

	"github.com/mikefsq/ptp"
)

// HotplugEvent reports a USB device on a registered vendor's id appearing on
// or leaving the bus. Attachment is the plugging-in's identity where the
// platform gives one (DeviceInfo.Attachment); on a detach it names the
// attachment that ended, so a holder of an open handle can match it.
type HotplugEvent struct {
	Attached   bool
	VID, PID   uint16
	Location   uint32
	Attachment uint64
}

// ErrNoHotplug is returned by Hotplug on a platform with no notification
// source; the caller keeps polling.
var ErrNoHotplug = errors.New("usb: no hotplug notification source on this platform")

// Hotplug subscribes to attach and detach notifications for devices on the
// registered vendors' ids and returns the channel they arrive on. The
// subscription ends when ctx does, and the channel is closed then. It is the
// interrupt a supervisor loop selects on beside its slow poll: the OS reports
// the change within milliseconds, where a poll sees it on its next pass.
//
// Sources: IOKit matching notifications on macOS (first-match and terminated
// on the USB device class) and the kernel uevent netlink socket on Linux;
// elsewhere ErrNoHotplug. On macOS the first-match notification also reports
// every device already attached when the subscription starts.
func Hotplug(ctx context.Context) (<-chan HotplugEvent, error) {
	return hotplug(ctx)
}

// registeredVID reports whether vid belongs to a vendor with a driver linked in.
func registeredVID(vid uint16) bool {
	for _, v := range ptp.Registered() {
		if uint16(v.ID) == vid {
			return true
		}
	}
	return false
}
