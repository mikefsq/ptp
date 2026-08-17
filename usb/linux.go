//go:build linux

// Linux USB transport over usbfs (/dev/bus/usb), driven straight through
// ioctls. No cgo and no libusb.
//
// PTP needs three endpoints on the Still Image interface (class 6): bulk OUT,
// bulk IN and interrupt IN. They are found by walking the device's
// configuration descriptor, because their addresses are not fixed across
// vendors — a Sony body and a Fujifilm body do not agree on them.
//
// STATUS: NOT YET RUN AGAINST A CAMERA. The macOS backend is the validated one.
//
// Access needs a udev rule granting the camera's vendor ID to the user, or the
// open fails with EACCES. Something like:
//
//	SUBSYSTEM=="usb", ATTR{idVendor}=="04cb", MODE="0666"
package usb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/mikefsq/ptp"
)

// usbfs ioctl encoding (asm-generic _IOC, 64-bit). dir: 0 none, 1 write,
// 2 read, 3 read+write. layout: (dir<<30)|(size<<16)|(type<<8)|nr, type 'U'.
func ioc(dir, nr, size uintptr) uintptr { return dir<<30 | size<<16 | 0x55<<8 | nr }

var (
	usbdevfsBulk            = ioc(3, 2, 24)     // _IOWR('U',2,usbdevfs_bulktransfer)
	usbdevfsClaimInterface  = ioc(2, 15, 4)     // _IOR('U',15,unsigned int)
	usbdevfsReleaseIface    = ioc(2, 16, 4)     // _IOR('U',16,unsigned int)
	usbdevfsClearHalt       = ioc(2, 21, 4)     // _IOR('U',21,unsigned int)
	usbdevfsDisconnectClaim = ioc(2, 27, 8+256) // _IOR('U',27,usbdevfs_disconnect_claim)
)

// bulkTransfer is <linux/usbdevice_fs.h> struct usbdevfs_bulktransfer.
type bulkTransfer struct {
	ep      uint32
	len     uint32
	timeout uint32 // milliseconds; 0 means wait forever
	_       uint32 // padding to align the pointer on 64-bit
	data    unsafe.Pointer
}

// disconnectClaim is struct usbdevfs_disconnect_claim. Claiming with
// USBDEVFS_DISCONNECT_CLAIM detaches whatever kernel driver holds the
// interface, which on a desktop Linux is how gphoto2's or the desktop
// environment's MTP handler is displaced — the equivalent of having to kill
// ptpcamerad on macOS, except that here the kernel will do it for us.
type disconnectClaim struct {
	iface  uint32
	flags  uint32
	driver [256]byte
}

const disconnectClaimExceptDriver = 0x02 // USBDEVFS_DISCONNECT_CLAIM_EXCEPT_DRIVER

// sysfsRead returns a trimmed sysfs attribute, or "" if it is unreadable.
func sysfsRead(dir, name string) string {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func sysfsHex(dir, name string) uint32 {
	v, err := strconv.ParseUint(sysfsRead(dir, name), 16, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}

// enumerate lists attached devices whose vendor ID is one of vids.
//
// sysfs is read rather than every /dev/bus/usb node opened: opening a device to
// ask what it is would disturb cameras that are in use, and enumeration must
// stay safe to call at any time.
func enumerate(vids []uint16) ([]DeviceInfo, error) {
	want := map[uint32]bool{}
	for _, v := range vids {
		want[uint32(v)] = true
	}

	entries, err := os.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return nil, fmt.Errorf("usb: reading /sys/bus/usb/devices: %w", err)
	}

	var out []DeviceInfo
	for _, e := range entries {
		dir := filepath.Join("/sys/bus/usb/devices", e.Name())
		// Interfaces appear as bus-port:config.iface; skip them, we want devices.
		if strings.Contains(e.Name(), ":") {
			continue
		}
		vid := sysfsHex(dir, "idVendor")
		if !want[vid] {
			continue
		}
		bus, err1 := strconv.Atoi(sysfsRead(dir, "busnum"))
		dev, err2 := strconv.Atoi(sysfsRead(dir, "devnum"))
		if err1 != nil || err2 != nil {
			continue
		}
		d := DeviceInfo{
			VID:      uint16(vid),
			PID:      uint16(sysfsHex(dir, "idProduct")),
			Location: uint32(bus)<<16 | uint32(dev&0xFFFF),
			Serial:   sysfsRead(dir, "serial"),
			Name:     sysfsRead(dir, "product"),
		}
		// bDeviceClass is 0 when the class lives on the interface, which is the
		// normal case for cameras; read the first interface's class so a body
		// sitting in mass-storage mode is reported as such rather than just
		// failing to open.
		d.IfaceClass = firstInterfaceClass(dir, e.Name())
		out = append(out, d)
	}
	return out, nil
}

// firstInterfaceClass reports bInterfaceClass of the device's first interface.
func firstInterfaceClass(dir, name string) uint32 {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), name+":") {
			continue
		}
		if c := sysfsRead(filepath.Join(dir, e.Name()), "bInterfaceClass"); c != "" {
			if v, err := strconv.ParseUint(c, 16, 32); err == nil {
				return uint32(v)
			}
		}
	}
	return 0
}

// usbfsTransport is the Linux Transport implementation. Every method holds the
// mutex: overlapping ioctls on one interface is not a thing usbfs promises to
// handle, and each camera gets its own transport, so this serialises per device
// rather than globally.
type usbfsTransport struct {
	mu                 sync.Mutex
	f                  *os.File
	inf                DeviceInfo
	epIn, epOut, epEvt uint8
	maxPacket          int
	iface              uint32
}

// openDevice opens one enumerated device. The portable front-end has already
// chosen it and checked its USB mode.
func openDevice(want DeviceInfo) (ptp.Transport, error) {
	bus, dev := want.Location>>16, want.Location&0xFFFF
	path := fmt.Sprintf("/dev/bus/usb/%03d/%03d", bus, dev)
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Errorf("usb: %s: permission denied opening %s — a udev rule "+
				"is needed, e.g. SUBSYSTEM==\"usb\", ATTR{idVendor}==\"%04x\", MODE=\"0666\"",
				want.Model(), path, want.VID)
		}
		return nil, fmt.Errorf("usb: opening %s: %w", path, err)
	}

	t := &usbfsTransport{f: f, inf: want}
	if err := t.findEndpoints(); err != nil {
		f.Close()
		return nil, err
	}
	// Detach whatever kernel driver holds the interface and claim it in one
	// step. Falling back to a plain claim keeps older kernels working, but then
	// a bound usbfs/MTP driver will refuse us.
	dc := disconnectClaim{iface: t.iface, flags: disconnectClaimExceptDriver}
	copy(dc.driver[:], "usbfs")
	if err := t.ioctl(usbdevfsDisconnectClaim, unsafe.Pointer(&dc)); err != nil {
		if err := t.ioctl(usbdevfsClaimInterface, unsafe.Pointer(&t.iface)); err != nil {
			f.Close()
			return nil, fmt.Errorf("usb: claiming interface %d of %s: %w (another driver, "+
				"such as a desktop MTP handler, may hold it)", t.iface, want.Model(), err)
		}
	}
	return t, nil
}

// findEndpoints walks the active configuration descriptor for the Still Image
// interface and records its three endpoints.
//
// The addresses are not fixed across vendors, so they cannot be hardcoded the
// way a single-vendor driver can get away with.
func (t *usbfsTransport) findEndpoints() error {
	// The device node begins with the device descriptor followed by the active
	// configuration. Reading from offset 0 returns exactly that.
	raw, err := os.ReadFile(fmt.Sprintf("/dev/bus/usb/%03d/%03d",
		t.inf.Location>>16, t.inf.Location&0xFFFF))
	if err != nil {
		return fmt.Errorf("usb: reading descriptors: %w", err)
	}
	if len(raw) < 18 {
		return fmt.Errorf("usb: descriptor block too short (%d bytes)", len(raw))
	}

	// Walk the descriptor chain: each entry is [length, type, ...].
	var inIface bool
	found := 0
	for i := 18; i+1 < len(raw); {
		l, typ := int(raw[i]), raw[i+1]
		if l < 2 || i+l > len(raw) {
			break
		}
		switch typ {
		case 0x04: // interface
			if i+6 < len(raw) {
				class := raw[i+5]
				inIface = class == 6 // still image
				if inIface {
					t.iface = uint32(raw[i+2])
					found = 0
				}
			}
		case 0x05: // endpoint
			if inIface && i+6 < len(raw) {
				addr, attr := raw[i+2], raw[i+3]&0x03
				mps := int(raw[i+4]) | int(raw[i+5])<<8
				switch {
				case attr == 2 && addr&0x80 != 0: // bulk IN
					t.epIn, t.maxPacket = addr, mps
					found++
				case attr == 2: // bulk OUT
					t.epOut = addr
					found++
				case attr == 3 && addr&0x80 != 0: // interrupt IN
					t.epEvt = addr
					found++
				}
			}
		}
		i += l
	}
	if t.epIn == 0 || t.epOut == 0 {
		return fmt.Errorf("usb: %s has no still-image interface with bulk endpoints "+
			"(found %d); is the body in tethered mode?", t.inf.Model(), found)
	}
	if t.maxPacket == 0 {
		t.maxPacket = 512
	}
	return nil
}

func (t *usbfsTransport) ioctl(req uintptr, arg unsafe.Pointer) error {
	for {
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, t.f.Fd(), req, uintptr(arg))
		if errno == 0 {
			return nil
		}
		// The Go runtime preempts with SIGURG, which interrupts ioctls
		// routinely. Retrying is correct; treating EINTR as failure would make
		// transfers fail at random under load.
		if errno == syscall.EINTR {
			continue
		}
		return errno
	}
}

// bulk runs one usbfs bulk transfer and returns the byte count.
func (t *usbfsTransport) bulk(ep uint8, p []byte, timeout time.Duration) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	bt := bulkTransfer{
		ep:      uint32(ep),
		len:     uint32(len(p)),
		timeout: ms(timeout),
		data:    unsafe.Pointer(&p[0]),
	}
	for {
		n, _, errno := syscall.Syscall(syscall.SYS_IOCTL, t.f.Fd(), usbdevfsBulk,
			uintptr(unsafe.Pointer(&bt)))
		switch errno {
		case 0:
			return int(n), nil
		case syscall.EINTR:
			continue
		case syscall.ETIMEDOUT:
			return 0, ptp.ErrTimeout
		case syscall.EPIPE:
			// A PTP device refuses an unsupported operation by stalling the
			// pipe rather than answering. Clear it, or every later transfer on
			// this endpoint fails too.
			e := uint32(ep)
			t.ioctl(usbdevfsClearHalt, unsafe.Pointer(&e))
			return 0, ptp.ErrStalled
		case syscall.ENODEV, syscall.ESHUTDOWN:
			return 0, ptp.ErrNotResponding
		default:
			return 0, fmt.Errorf("usb: bulk transfer on endpoint 0x%02X: %w", ep, errno)
		}
	}
}

func (t *usbfsTransport) BulkOut(p []byte, timeout time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	sent := 0
	for sent < len(p) {
		n, err := t.bulk(t.epOut, p[sent:], timeout)
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("usb: bulk write stalled at %d/%d bytes", sent, len(p))
		}
		sent += n
	}
	// A transfer that is an exact multiple of the packet size needs a
	// zero-length packet, or the device never sees it end.
	if len(p)%t.maxPacket == 0 {
		t.bulk(t.epOut, []byte{}, timeout)
	}
	return nil
}

func (t *usbfsTransport) BulkIn(p []byte, timeout time.Duration) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.bulk(t.epIn, p, timeout)
}

func (t *usbfsTransport) InterruptIn(p []byte, timeout time.Duration) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.epEvt == 0 {
		return 0, ptp.ErrTimeout // no event endpoint; polling it is not an error
	}
	return t.bulk(t.epEvt, p, timeout)
}

func (t *usbfsTransport) MaxPacketSize() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.maxPacket
}

// Info reports which device this transport is bound to.
func (t *usbfsTransport) Info() DeviceInfo { return t.inf }

func (t *usbfsTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.f == nil {
		return nil
	}
	t.ioctl(usbdevfsReleaseIface, unsafe.Pointer(&t.iface))
	err := t.f.Close()
	t.f = nil
	return err
}
