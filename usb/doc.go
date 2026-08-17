// Package usb carries PTP over USB bulk endpoints.
//
// The cgo lives here and nowhere else. macOS needs IOKit, so darwin.go is cgo;
// Linux is driven straight through usbfs ioctls and needs neither cgo nor
// libusb. The parent package stays pure Go either way, which is what keeps its
// tests runnable without a camera or a Mac.
//
// Enumeration finds only vendors whose driver is linked into the binary —
// importing ptp/fuji is what makes Fujifilm bodies visible. Matching a camera
// the build cannot drive helps nobody.
//
// # Status
//
// The macOS backend is hardware-validated on a Fujifilm X-T5 and a Sony NEX-6.
// The Linux backend has NEVER been run against a camera. Windows is a stub.
//
// # Traps the macOS backend already solves
//
//   - A top-level idVendor entry in an IOKit matching dictionary matches
//     NOTHING on the modern IOUSBHost stack. No error, no results. Match the
//     device class and compare the property read back instead.
//   - USBInterfaceOpenSeize does NOT take a camera away from macOS's
//     ptpcamerad. Both open and seize return kIOReturnExclusiveAccess. The
//     daemon must be killed; SIGSTOP does not work, because a stopped process
//     keeps its claim.
//   - BulkOut must append a zero-length packet when the payload is an exact
//     multiple of the endpoint's max packet size, or the device never sees the
//     transfer end.
//
// # Linux
//
// Access needs a udev rule granting the camera's vendor ID to the user, or the
// open fails with EACCES:
//
//	SUBSYSTEM=="usb", ATTR{idVendor}=="04cb", MODE="0666"
//
// The interface is claimed with USBDEVFS_DISCONNECT_CLAIM, which detaches
// whatever kernel driver holds it — the equivalent of having to kill
// ptpcamerad on macOS, except the kernel does it for us. Endpoints are found by
// walking the configuration descriptor rather than hardcoded, because their
// addresses are not fixed across vendors.
//
// # Recovering a wedged camera
//
// A transaction abandoned mid-transfer leaves the camera waiting for the rest of
// a data phase: it then answers nothing, and every request costs a full timeout.
// The darwin backend implements ptp.Resetter, which clears both bulk endpoints'
// halt state at BOTH ends and then issues the still-image class Device Reset
// (bRequest 0x66) to put the camera's own PTP stack back to idle. Session.Open
// calls it automatically before giving up, so a wedged body recovers without
// being unplugged. Hardware-confirmed accepted by a NEX-6.
//
// This matters more on macOS than it should, because ptpcamerad is
// demand-launched whenever a still-image device is enumerated and SIP forbids
// booting it out of the domain, so it can reappear underneath a live session.
//
// # Not this package's problem
//
// A bulk read is not container-aligned and one read can return several
// containers. That belongs to the framing layer, and is why BulkIn returns a
// byte count rather than a container.
package usb
