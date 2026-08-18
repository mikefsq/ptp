//go:build darwin

// macOS USB transport over IOKit's IOUSBLib. This is the one part that needs
// cgo — there is no pure-Go path to IOKit.
//
// PTP needs three endpoints on the Still Image interface (class 6): bulk OUT,
// bulk IN, and interrupt IN.
//
// STATUS: this backend is ported from sonycam, where it is hardware-validated
// against a real camera. Only the vendor ID differs. It has NOT been run
// against a Fujifilm body.
//
// Three findings carried over from that bring-up, all of which cost real time:
//
//   - A top-level idVendor entry in the matching dictionary matches NOTHING on
//     the modern IOUSBHost stack, and fails silently rather than erroring. The
//     class is matched and the property compared by hand. See ptp_enumerate.
//   - USBInterfaceOpenSeize does NOT take a still-image device away from
//     macOS's ptpcamerad: both the open and the seize return
//     kIOReturnExclusiveAccess. ptpcamerad has to be killed, and SIGSTOP is no
//     use because a stopped process keeps its claim.
//   - A camera that has gone to sleep produces failures that look exactly like
//     protocol bugs. ptp.ErrNotResponding distinguishes them.

package usb

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <IOKit/IOKitLib.h>
#include <IOKit/IOCFPlugIn.h>
#include <IOKit/usb/IOUSBLib.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>
#include <string.h>

// ptp_dev is one opened camera: the device, its claimed still-image
// interface, and the three pipe references PTP needs.
typedef struct {
    IOUSBDeviceInterface**    dev;
    IOUSBInterfaceInterface** intf;
    UInt8  outPipe;
    UInt8  inPipe;
    UInt8  evtPipe;
    UInt16 maxPacket;
} ptp_dev;

// ptp_info describes an enumerated device without opening it.
typedef struct {
    UInt16   vid;
    UInt16   pid;
    UInt32   location;
    uint32_t ifaceClass;   // USB interface class: 6 = still image (PTP), 8 = mass storage
    uint64_t entry;        // IORegistry entry id: new on every plugging-in
    char     serial[64];
    char     name[128];
} ptp_info;

static IOUSBDeviceInterface** device_interface(io_service_t svc) {
    IOCFPlugInInterface** plugin = NULL; SInt32 score;
    if (IOCreatePlugInInterfaceForService(svc, kIOUSBDeviceUserClientTypeID,
            kIOCFPlugInInterfaceID, &plugin, &score) != KERN_SUCCESS || !plugin)
        return NULL;
    IOUSBDeviceInterface** dev = NULL;
    (*plugin)->QueryInterface(plugin, CFUUIDGetUUIDBytes(kIOUSBDeviceInterfaceID), (LPVOID*)&dev);
    (*plugin)->Release(plugin);
    return dev;
}

// reg_str reads a string property from the IORegistry.
static void reg_str(io_service_t svc, CFStringRef key, char* out, size_t n) {
    out[0] = 0;
    CFTypeRef s = IORegistryEntrySearchCFProperty(svc, kIOServicePlane,
        key, kCFAllocatorDefault, kIORegistryIterateRecursively);
    if (!s) return;
    if (CFGetTypeID(s) == CFStringGetTypeID())
        CFStringGetCString((CFStringRef)s, out, n, kCFStringEncodingUTF8);
    CFRelease(s);
}

// reg_u32 reads a numeric property from the IORegistry. Returns 0 on success.
static int reg_u32(io_service_t svc, CFStringRef key, uint32_t* out) {
    CFTypeRef v = IORegistryEntrySearchCFProperty(svc, kIOServicePlane,
        key, kCFAllocatorDefault, kIORegistryIterateRecursively);
    if (!v) return -1;
    int rc = -1;
    if (CFGetTypeID(v) == CFNumberGetTypeID()) {
        SInt64 n = 0;
        if (CFNumberGetValue((CFNumberRef)v, kCFNumberSInt64Type, &n)) {
            *out = (uint32_t)n;
            rc = 0;
        }
    }
    CFRelease(v);
    return rc;
}

// ptp_diag carries where-did-it-fail detail out of the open path, so an
// opaque IOReturn becomes a precise message on the Go side.
typedef struct {
    int nInterfaces;   // interfaces the iterator returned
    int classes[8];    // their bInterfaceClass values
    int stillImageIf;  // 1 if a class-6 interface was seen
    int openKR;        // IOReturn from USBInterfaceOpen
    int seizeKR;       // IOReturn from USBInterfaceOpenSeize
    int numEndpoints;
    int outPipe, inPipe, evtPipe;
} ptp_diag;

// claim_interface finds the still-image (PTP) interface and claims it,
// seizing it from ptpcamerad if necessary, then locates the three pipes.
// Returns 0 on success.
static int claim_interface_d(IOUSBDeviceInterface** dev, ptp_dev* out, ptp_diag* dg);

static int claim_interface(IOUSBDeviceInterface** dev, ptp_dev* out) {
    ptp_diag dg = {0};
    return claim_interface_d(dev, out, &dg);
}

static int claim_interface_d(IOUSBDeviceInterface** dev, ptp_dev* out, ptp_diag* dg) {
    IOUSBFindInterfaceRequest req;
    req.bInterfaceClass    = kIOUSBFindInterfaceDontCare;
    req.bInterfaceSubClass = kIOUSBFindInterfaceDontCare;
    req.bInterfaceProtocol = kIOUSBFindInterfaceDontCare;
    req.bAlternateSetting  = kIOUSBFindInterfaceDontCare;
    io_iterator_t it;
    if ((*dev)->CreateInterfaceIterator(dev, &req, &it) != kIOReturnSuccess) return -1;

    int rc = -2;
    io_service_t usbIf;
    while ((usbIf = IOIteratorNext(it))) {
        IOCFPlugInInterface** pl = NULL; SInt32 score;
        if (IOCreatePlugInInterfaceForService(usbIf, kIOUSBInterfaceUserClientTypeID,
                kIOCFPlugInInterfaceID, &pl, &score) != KERN_SUCCESS || !pl) {
            IOObjectRelease(usbIf);
            continue;
        }
        IOUSBInterfaceInterface** intf = NULL;
        (*pl)->QueryInterface(pl, CFUUIDGetUUIDBytes(kIOUSBInterfaceInterfaceID), (LPVOID*)&intf);
        (*pl)->Release(pl);
        IOObjectRelease(usbIf);
        if (!intf) continue;

        UInt8 cls = 0, sub = 0, proto = 0;
        (*intf)->GetInterfaceClass(intf, &cls);
        (*intf)->GetInterfaceSubClass(intf, &sub);
        (*intf)->GetInterfaceProtocol(intf, &proto);
        if (dg->nInterfaces < 8) dg->classes[dg->nInterfaces] = cls;
        dg->nInterfaces++;
        // 6/1/1 is PTP still image. Anything else on the body (mass storage,
        // audio on the video-capable bodies) is not ours.
        if (cls != 6) { (*intf)->Release(intf); continue; }
        dg->stillImageIf = 1;

        IOReturn kr = (*intf)->USBInterfaceOpen(intf);
        dg->openKR = kr;
        if (kr != kIOReturnSuccess) {
            kr = (*intf)->USBInterfaceOpenSeize(intf);
            dg->seizeKR = kr;
        }
        if (kr != kIOReturnSuccess) { (*intf)->Release(intf); rc = -4; continue; }

        UInt8 n = 0;
        (*intf)->GetNumEndpoints(intf, &n);
        dg->numEndpoints = n;
        out->outPipe = out->inPipe = out->evtPipe = 0;
        for (UInt8 i = 1; i <= n; i++) {
            UInt8 dir, num, tt, interval; UInt16 maxp;
            if ((*intf)->GetPipeProperties(intf, i, &dir, &num, &tt, &maxp, &interval) != kIOReturnSuccess)
                continue;
            if (tt == kUSBBulk && dir == kUSBOut && !out->outPipe) {
                out->outPipe = i;
            } else if (tt == kUSBBulk && dir == kUSBIn && !out->inPipe) {
                out->inPipe = i;
                out->maxPacket = maxp;
            } else if (tt == kUSBInterrupt && dir == kUSBIn && !out->evtPipe) {
                out->evtPipe = i;
            }
        }
        dg->outPipe = out->outPipe; dg->inPipe = out->inPipe; dg->evtPipe = out->evtPipe;
        if (!out->outPipe || !out->inPipe) {
            (*intf)->USBInterfaceClose(intf);
            (*intf)->Release(intf);
            rc = -5;
            continue;
        }
        out->intf = intf;
        IOObjectRelease(it);
        return 0;
    }
    IOObjectRelease(it);
    return rc;
}

// ptp_enumerate fills up to max entries with Sony USB devices. Returns the
// count, or negative on failure.
// The main port argument is 0 (kIOMainPortDefault); the old
// kIOMasterPortDefault spelling is deprecated from macOS 12.
static int ptp_enumerate(const uint32_t* vids, int nvids, ptp_info* list, int max) {
    // Match every USB device and filter by reading idVendor ourselves, rather
    // than putting idVendor in the match dictionary. On the modern IOUSBHost
    // stack a class + top-level idVendor dictionary matches NOTHING — it fails
    // silently, returning zero results rather than an error. (Putting it under
    // IOPropertyMatch does work, but reading the property back is what asicam's
    // darwin backend settled on, and it does not depend on that subtlety.)
    CFMutableDictionaryRef m = IOServiceMatching(kIOUSBDeviceClassName);
    if (!m) return -1;

    io_iterator_t it;
    if (IOServiceGetMatchingServices(0, m, &it) != KERN_SUCCESS) return -2;

    int n = 0;
    io_service_t svc;
    while ((svc = IOIteratorNext(it))) {
        uint32_t gotVID = 0;
        int wanted = 0;
        if (reg_u32(svc, CFSTR("idVendor"), &gotVID) == 0) {
            for (int i = 0; i < nvids; i++) if (vids[i] == gotVID) { wanted = 1; break; }
        }
        if (wanted && n < max) {
            uint32_t pid = 0, loc = 0;
            reg_u32(svc, CFSTR("idProduct"), &pid);
            reg_u32(svc, CFSTR("locationID"), &loc);
            list[n].vid = (UInt16)gotVID;
            list[n].pid = (UInt16)pid;
            list[n].location = loc;
            reg_str(svc, CFSTR("USB Serial Number"), list[n].serial, sizeof(list[n].serial));
            reg_str(svc, CFSTR("USB Product Name"), list[n].name, sizeof(list[n].name));
            // bDeviceClass 0 means the class lives on the interface; record the
            // interface class so a body sitting in Mass Storage mode can be
            // reported as such instead of just failing to open.
            list[n].ifaceClass = 0;
            reg_u32(svc, CFSTR("bInterfaceClass"), &list[n].ifaceClass);
            list[n].entry = 0;
            IORegistryEntryGetRegistryEntryID(svc, &list[n].entry);
            n++;
        }
        IOObjectRelease(svc);
    }
    IOObjectRelease(it);
    return n;
}

// ptp_open opens the Sony device at the given location ID.
// Returns 0 on success; -100 means the device was found but is exclusively
// held even after a seize attempt.
static int ptp_open(uint32_t vid, UInt32 location, ptp_dev* out, ptp_diag* dg) {
    // Same caveat as ptp_enumerate: filter on the property we read back, not
    // in the match dictionary.
    CFMutableDictionaryRef m = IOServiceMatching(kIOUSBDeviceClassName);
    if (!m) return -1;

    io_iterator_t it;
    if (IOServiceGetMatchingServices(0, m, &it) != KERN_SUCCESS) return -2;

    int rc = -3;
    io_service_t svc;
    while ((svc = IOIteratorNext(it))) {
        uint32_t gotVID = 0, gotLoc = 0;
        if (reg_u32(svc, CFSTR("idVendor"), &gotVID) != 0 || gotVID != vid ||
            reg_u32(svc, CFSTR("locationID"), &gotLoc) != 0 || gotLoc != location) {
            IOObjectRelease(svc);
            continue;
        }
        IOUSBDeviceInterface** dev = device_interface(svc);
        IOObjectRelease(svc);
        if (!dev) continue;

        IOReturn kr = (*dev)->USBDeviceOpen(dev);
        if (kr != kIOReturnSuccess)
            kr = (*dev)->USBDeviceOpenSeize(dev);
        if (kr != kIOReturnSuccess) { (*dev)->Release(dev); rc = -100; break; }

        if (claim_interface_d(dev, out, dg) != 0) {
            (*dev)->USBDeviceClose(dev);
            (*dev)->Release(dev);
            rc = -101;
            break;
        }
        out->dev = dev;
        rc = 0;
        break;
    }
    IOObjectRelease(it);
    return rc;
}

static int ptp_write(ptp_dev* d, void* buf, UInt32 len, UInt32 ms) {
    return (*d->intf)->WritePipeTO(d->intf, d->outPipe, buf, len, ms, ms);
}

static int ptp_read(ptp_dev* d, void* buf, UInt32* len, UInt32 ms) {
    return (*d->intf)->ReadPipeTO(d->intf, d->inPipe, buf, len, ms, ms);
}

static int ptp_read_event(ptp_dev* d, void* buf, UInt32* len, UInt32 ms) {
    if (!d->evtPipe) return kIOReturnNoDevice;
    return (*d->intf)->ReadPipeTO(d->intf, d->evtPipe, buf, len, ms, ms);
}

// ptp_clear_stall clears a halted bulk endpoint. A camera refuses an
// unsupported operation by stalling the pipe, and until the halt is cleared
// every later transfer on that endpoint fails too — including the CloseSession
// that would otherwise tidy up.
static int ptp_clear_stall(ptp_dev* d, int which) {
    if (!d->intf) return -1;
    UInt8 pipe = which == 0 ? d->outPipe : (which == 1 ? d->inPipe : d->evtPipe);
    if (!pipe) return -1;
    // ClearPipeStallBothEnds also sends CLEAR_FEATURE(ENDPOINT_HALT) to the
    // device, so the camera's own endpoint state is reset, not just the host's.
    return (*d->intf)->ClearPipeStallBothEnds(d->intf, pipe);
}

// ptp_reset recovers a device left in a bad state — typically a transaction
// abandoned mid-transfer, after which the camera is still waiting for the rest
// of a data phase and answers nothing.
//
// Two steps, in order. Clearing both bulk endpoints resets the halt state at
// BOTH ends. Then the still-image class Device Reset (bRequest 0x66, PIMA
// 15740 clause 13.2.2) tells the camera's own PTP stack to abandon whatever it
// thought was in flight and return to an idle, session-less state. The halts
// have to go first: the control request is useless while the pipes the reply
// would travel on are still halted.
static int ptp_reset(ptp_dev* d, UInt32 ms) {
    if (!d->intf) return -1;
    if (d->inPipe)  (*d->intf)->ClearPipeStallBothEnds(d->intf, d->inPipe);
    if (d->outPipe) (*d->intf)->ClearPipeStallBothEnds(d->intf, d->outPipe);

    IOUSBDevRequestTO req;
    req.bmRequestType = USBmakebmRequestType(kUSBOut, kUSBClass, kUSBInterface);
    req.bRequest      = 0x66; // Device Reset Request
    req.wValue        = 0;
    req.wIndex        = 0;    // IOKit fills in the interface number
    req.wLength       = 0;
    req.pData         = NULL;
    req.completionTimeout = ms;
    req.noDataTimeout     = ms;
    return (*d->intf)->ControlRequestTO(d->intf, 0, &req);
}

static void ptp_close(ptp_dev* d) {
    if (d->intf) {
        (*d->intf)->USBInterfaceClose(d->intf);
        (*d->intf)->Release(d->intf);
        d->intf = NULL;
    }
    if (d->dev) {
        (*d->dev)->USBDeviceClose(d->dev);
        (*d->dev)->Release(d->dev);
        d->dev = NULL;
    }
}
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/mikefsq/ptp"
)

// kIOReturnExclusiveAccess, the error macOS returns when ptpcamerad already
// holds the device.
const errExclusiveAccess = 0xE00002C5

// kIOUSBPipeStalled. A PTP device refuses an operation it does not support by
// stalling the bulk pipe rather than returning a response code.
const errPipeStalled = 0xE000404F

// The device has stopped answering: it went to sleep, was unplugged, or its
// firmware wedged. Every later request costs a full timeout, so these are
// mapped to ptp.ErrNotResponding and callers are expected to give up promptly.
const (
	errNotResponding      = 0xE00002EB // kIOReturnNotResponding
	errTransactionTimeout = 0xE0004051 // kIOUSBTransactionTimeout
	errNoDevice           = 0xE00002C0 // kIOReturnNoDevice
)

func isDead(kr uint32) bool {
	return kr == errNotResponding || kr == errTransactionTimeout || kr == errNoDevice
}

// enumerate lists attached devices whose vendor ID is one of vids.
func enumerate(vids []uint16) ([]DeviceInfo, error) {
	const max = 16
	if len(vids) == 0 {
		return nil, nil
	}
	cv := make([]C.uint32_t, len(vids))
	for i, v := range vids {
		cv[i] = C.uint32_t(v)
	}
	list := make([]C.ptp_info, max)
	n := C.ptp_enumerate(&cv[0], C.int(len(cv)), &list[0], C.int(max))
	if n < 0 {
		return nil, fmt.Errorf("usb: enumerating USB devices failed (%d)", int(n))
	}
	out := make([]DeviceInfo, 0, int(n))
	for i := 0; i < int(n); i++ {
		out = append(out, DeviceInfo{
			VID:        uint16(list[i].vid),
			PID:        uint16(list[i].pid),
			Location:   uint32(list[i].location),
			Serial:     C.GoString(&list[i].serial[0]),
			Name:       C.GoString(&list[i].name[0]),
			IfaceClass: uint32(list[i].ifaceClass),
			Attachment: uint64(list[i].entry),
		})
	}
	return out, nil
}

// usbTransport is the darwin Transport implementation. Every method holds the
// mutex: IOUSBLib pipe calls on one interface are not safe to overlap, and two
// cameras each get their own transport, so this serialises per device rather
// than globally.
type usbTransport struct {
	mu  sync.Mutex
	d   C.ptp_dev
	inf DeviceInfo
}

// openDevice opens one enumerated device. The portable front-end in usb.go has
// already chosen it and checked its USB mode.
func openDevice(want DeviceInfo) (ptp.Transport, error) {
	t := &usbTransport{inf: want}
	var dg C.ptp_diag
	switch rc := C.ptp_open(C.uint32_t(want.VID), C.UInt32(want.Location), &t.d, &dg); int(rc) {
	case 0:
		return t, nil
	case -100:
		return nil, fmt.Errorf("usb: %s is held exclusively and could not be seized "+
			"(IOReturn 0x%08X). macOS's ptpcamerad claims still-image devices on "+
			"enumeration. SIGSTOP will not help — a stopped process keeps its claim; "+
			"it must be killed: pkill -9 ptpcamerad", want, errExclusiveAccess)
	case -101:
		return nil, fmt.Errorf("usb: %s: %s", want, describeClaimFailure(&dg))
	default:
		return nil, fmt.Errorf("usb: opening %s failed (%d); %s", want, int(rc), describeClaimFailure(&dg))
	}
}

// Info reports which device this transport is bound to.
func (t *usbTransport) Info() DeviceInfo { return t.inf }

func (t *usbTransport) BulkOut(p []byte, timeout time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(p) == 0 {
		return nil
	}
	kr := C.ptp_write(&t.d, unsafe.Pointer(&p[0]), C.UInt32(len(p)), C.UInt32(ms(timeout)))
	if kr != 0 {
		if uint32(kr) == errPipeStalled {
			C.ptp_clear_stall(&t.d, 0)
			return fmt.Errorf("usb: bulk write of %d bytes stalled; endpoint cleared: %w", len(p), ptp.ErrStalled)
		}
		if isDead(uint32(kr)) {
			return fmt.Errorf("usb: bulk write of %d bytes (IOReturn 0x%08X): %w",
				len(p), uint32(kr), ptp.ErrNotResponding)
		}
		return fmt.Errorf("usb: bulk write of %d bytes failed (IOReturn 0x%08X)", len(p), uint32(kr))
	}
	// IOUSBLib does not append the terminating zero-length packet itself. PTP
	// containers that land on an exact multiple of the packet size would
	// otherwise leave the device waiting for more of the transfer.
	if mp := t.MaxPacketSizeLocked(); mp > 0 && len(p)%mp == 0 {
		if kr := C.ptp_write(&t.d, nil, 0, C.UInt32(ms(timeout))); kr != 0 {
			return fmt.Errorf("usb: writing terminating zero-length packet failed (IOReturn 0x%08X)", uint32(kr))
		}
	}
	return nil
}

func (t *usbTransport) BulkIn(p []byte, timeout time.Duration) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(p) == 0 {
		return 0, nil
	}
	n := C.UInt32(len(p))
	kr := C.ptp_read(&t.d, unsafe.Pointer(&p[0]), &n, C.UInt32(ms(timeout)))
	if kr != 0 {
		if uint32(kr) == errPipeStalled {
			// The camera stalled to refuse the operation. Clear both endpoints:
			// leaving either halted wedges every later transaction.
			C.ptp_clear_stall(&t.d, 1)
			C.ptp_clear_stall(&t.d, 0)
			return 0, ptp.ErrStalled
		}
		if isDead(uint32(kr)) {
			return 0, fmt.Errorf("usb: bulk read (IOReturn 0x%08X): %w", uint32(kr), ptp.ErrNotResponding)
		}
		return 0, fmt.Errorf("usb: bulk read failed (IOReturn 0x%08X)", uint32(kr))
	}
	return int(n), nil
}

func (t *usbTransport) InterruptIn(p []byte, timeout time.Duration) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(p) == 0 {
		return 0, nil
	}
	n := C.UInt32(len(p))
	kr := C.ptp_read_event(&t.d, unsafe.Pointer(&p[0]), &n, C.UInt32(ms(timeout)))
	if kr != 0 {
		// A quiet event pipe is the normal case, not a fault.
		return 0, ptp.ErrTimeout
	}
	return int(n), nil
}

func (t *usbTransport) MaxPacketSize() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.MaxPacketSizeLocked()
}

// MaxPacketSizeLocked assumes t.mu is held.
func (t *usbTransport) MaxPacketSizeLocked() int { return int(t.d.maxPacket) }

func (t *usbTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	C.ptp_close(&t.d)
	return nil
}

// describeClaimFailure turns the C-side diagnostics into a specific reason,
// rather than one message covering three unrelated failures.
func describeClaimFailure(dg *C.ptp_diag) string {
	classes := make([]int, 0, 8)
	for i := 0; i < int(dg.nInterfaces) && i < 8; i++ {
		classes = append(classes, int(dg.classes[i]))
	}
	switch {
	case dg.nInterfaces == 0:
		return "the device exposed no USB interfaces at all"
	case dg.stillImageIf == 0:
		return fmt.Sprintf("no still-image (class 6) interface; the device offers %d interface(s) of class %v — "+
			"set the body's USB connection mode to PC Remote or MTP", int(dg.nInterfaces), classes)
	case dg.seizeKR != 0:
		hint := ""
		if uint32(dg.seizeKR) == errExclusiveAccess {
			hint = " — macOS's ptpcamerad already holds it. SIGSTOP will not help, a " +
				"stopped process keeps its claim; it must be killed: pkill -9 ptpcamerad"
		}
		return fmt.Sprintf("its still-image interface could not be claimed (open 0x%08X, seize 0x%08X)%s",
			uint32(dg.openKR), uint32(dg.seizeKR), hint)
	case dg.outPipe == 0 || dg.inPipe == 0:
		return fmt.Sprintf("its still-image interface was claimed but the bulk pipes were not found "+
			"(%d endpoints; out=%d in=%d evt=%d)",
			int(dg.numEndpoints), int(dg.outPipe), int(dg.inPipe), int(dg.evtPipe))
	}
	return fmt.Sprintf("its still-image interface could not be claimed (%d interfaces, classes %v)",
		int(dg.nInterfaces), classes)
}

// Reset recovers a camera left wedged by an abandoned transaction, without a
// power cycle.
//
// This is not hypothetical housekeeping. On macOS, ptpcamerad is demand-launched
// whenever a still-image device is enumerated, and SIP prevents removing it, so
// it can appear part-way through a session and abort a transfer in flight. The
// camera is then left waiting for the rest of a data phase and answers nothing —
// every later request costs a full timeout, and the only other cure is
// unplugging the body.
func (t *usbTransport) Reset() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if kr := C.ptp_reset(&t.d, C.UInt32(ms(ptp.DefaultTimeout))); kr != 0 {
		return fmt.Errorf("usb: device reset failed (IOReturn 0x%08X)", uint32(kr))
	}
	return nil
}
