package fuji

import (
	"errors"
	"fmt"
	"time"

	"github.com/mikefsq/ptp"
)

// Capture on a Fujifilm body, and the exposure encodings that go with it.
//
// Unlike Sony, Fujifilm uses the standard PTP InitiateCapture. But a bare
// InitiateCapture is not enough: the body must first be told the host has
// priority, and afterwards its autofocus status must be polled, because the
// camera reports focus failure there rather than in the operation's response.
//
// The sequence is taken from gphoto2's working Fuji implementation.

// Fujifilm capture-related properties.
const (
	// PropPriorityMode hands control to the host. Without it the body ignores
	// remote capture.
	PropPriorityModeCode ptp.Prop = 0xD207

	// PropReleaseGesture (0xD208) carries the shutter gesture. It is NOT
	// gphoto2's PropReleaseMode (0xD201), which is a different property; this
	// is the one whose 17 advertised values match the SDK's XSDK_RELEASE_*
	// constants exactly, and the one gphoto2 writes 0x0200/0x0004 to.
	PropReleaseGesture ptp.Prop = 0xD208

	// PropAFStatusCode reports focus progress during a capture.
	PropAFStatusCode ptp.Prop = 0xD209
)

// Priority modes for PropPriorityModeCode, the SDK's XSDK_PRIORITY_* values.
//
// PC Priority is not merely advisory: while it is set, the camera's dials,
// buttons and shutter release are LOCKED OUT — everything except the lens
// manual focus ring. It therefore has to be given back. Closing the PTP session
// does not do it, so a host that simply exits leaves the body inert in the
// user's hands and held awake. ptp.Session.Close restores camera priority.
const (
	PriorityModeCamera uint64 = 0x0001 // body controls work
	PriorityModeHost   uint64 = 0x0002 // host drives; body controls dead
)

// Release modes for PropReleaseMode (0xD208).
//
// The value is a bitfield: the high byte is the press action, the low byte the
// release action, so a whole shot is one value combining both. Sony models the
// shutter as two buttons you press and release separately; Fujifilm packs the
// whole gesture into a single write.
//
// Confirmed against the SDK's XSDK_RELEASE_* constants and against the 17
// values a real X-T5 advertises for this property.
const (
	// Press actions (high byte).
	ReleaseShoot     uint64 = 0x0100 // shoot
	ReleaseS1On      uint64 = 0x0200 // half press: AF and metering
	ReleaseS2        uint64 = 0x0300 // full press from half-pressed
	ReleaseBulbOn    uint64 = 0x0400
	ReleaseBulbS2On  uint64 = 0x0500 // start a bulb exposure
	ReleaseRecStart  uint64 = 0x0800 // start movie
	ReleasePixShift  uint64 = 0x4000
	ReleaseCustomWB  uint64 = 0x8000
	ReleaseAEOn      uint64 = 0x9000
	ReleaseAFOn      uint64 = 0x9100
	ReleaseAFAEOn    uint64 = 0x9200
	ReleaseAF        uint64 = 0x9300
	ReleaseInstantAF uint64 = 0xA000

	// Release actions (low byte).
	ReleaseNAFOff   uint64 = 0x0001
	ReleaseNAEOff   uint64 = 0x0002
	ReleaseNS1Off   uint64 = 0x0004
	ReleaseNBulbOff uint64 = 0x0008
	ReleaseNInstAF  uint64 = 0x0010
	ReleaseNAF      uint64 = 0x0020
	ReleaseRecStop  uint64 = 0x0080

	// ReleaseShootS1Off is a complete single frame: full press from the
	// half-pressed state, then release. THIS is what takes a picture.
	//
	// Sending only ReleaseNS1Off releases a half-press that was never made and
	// fires nothing, and ReleaseS1On alone only focuses — both are accepted by
	// the camera and neither produces a frame.
	ReleaseShootS1Off uint64 = ReleaseS2 | ReleaseNS1Off // 0x0304

	// ReleaseBulbEnd closes a bulb exposure and lets the half-press up.
	ReleaseBulbEnd uint64 = ReleaseNBulbOff | ReleaseNS1Off // 0x000C
)

// AF status values reported by PropAFStatusCode.
const (
	AFStatusBusy   uint64 = 0x0001 // still working; keep polling
	AFStatusOK     uint64 = 0x0002 // focused, shot taken
	AFStatusFailed uint64 = 0x0003 // could not focus; the shot did not happen
)

// TakePriority tells the camera the host is in charge.
//
// A camera holding undownloaded frames REFUSES, with vendor code 0xA001. That
// is the same rule the SDK documents for the other direction, where its own
// sample loops "waiting for the in-camera buffer empty": pending frames block
// priority transitions BOTH ways. The refusal is opaque on its own — the
// property's descriptor says the value is allowed, and it is — so the frame
// count is looked up and reported, because "download the frames first" is
// actionable where "0xA001" is not.
func (c *Camera) TakePriority() error {
	// Read first: an X-T5 answers a redundant write with vendor error 0xA002,
	// so setting a mode it is already in is worse than doing nothing.
	if v, err := c.GetPropValue(PropPriorityModeCode, ptp.TypeUint16); err == nil && v.Num == PriorityModeHost {
		return nil
	}
	err := c.SetPropValue(PropPriorityModeCode, ptp.TypeUint16, PriorityModeHost)
	if err == nil {
		return nil
	}
	if n, lerr := c.pendingFrames(); lerr == nil && n > 0 {
		return fmt.Errorf("fuji: the camera will not hand control to the host while it "+
			"holds %d undownloaded frame(s) — download and delete them first: %w", n, err)
	}
	return err
}

// pendingFrames counts frames sitting unread in the volatile store.
func (c *Camera) pendingFrames() (int, error) {
	h, err := c.GetObjectHandles(StillStore, 0, 0)
	if err != nil {
		return 0, err
	}
	return len(h), nil
}

// PropMediaRecordCode governs whether a tethered frame is ALSO copied to the SD
// card. It does not affect the USB transfer: the SDK's storage model says every
// frame taken while tethered lands in the in-camera volatile store and reaches
// the host from there, and this setting only asks the camera to additionally
// write a copy to the card. With it off, nothing is recorded to the card at all
// and the frame exists only until it is downloaded.
const PropMediaRecordCode ptp.Prop = 0xD20C

// Values for PropMediaRecordCode (the SDK's XSDK_MEDIAREC_* constants).
const (
	MediaRecordRawJPEG uint64 = 0x0001 // record both RAW and JPEG to the card
	MediaRecordRaw     uint64 = 0x0002 // record RAW only
	MediaRecordJPEG    uint64 = 0x0003 // record JPEG only
	MediaRecordOff     uint64 = 0x0004 // record nothing; USB transfer only
)

// MediaRecord reports whether tethered frames are also written to the SD card.
func (c *Camera) MediaRecord() (uint64, error) {
	v, err := c.GetPropValue(PropMediaRecordCode, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// SetMediaRecord chooses what a tethered frame leaves behind on the SD card.
//
// Recording to the card costs shot-to-shot time, so an eclipse sequence that
// downloads every frame anyway should use MediaRecordOff — at the price of
// having no in-camera copy if the download fails.
func (c *Camera) SetMediaRecord(mode uint64) error {
	// Same read-before-write rule as TakePriority: a redundant write draws
	// vendor error 0xA002 from an X-T5.
	if v, err := c.MediaRecord(); err == nil && v == mode {
		return nil
	}
	return c.SetPropValue(PropMediaRecordCode, ptp.TypeUint16, mode)
}

// Values for PropQuality (0xD018), which decides what the camera produces and
// therefore how many bytes cross USB per frame — the dominant cost in a
// continuous sequence, far larger than the SD card copy.
//
// An X-T5 advertises exactly five, matching the first five of the SDK's seven
// XT5_IMAGEQUALITY_* constants (it has no SUPERFINE). The numbering is
// gphoto2's and is INFERRED, not confirmed on the wire: it is consistent with
// the body reading back 1 while producing a lone .raf and no JPEG pair.
const (
	QualityRaw       uint64 = 0x0001
	QualityFine      uint64 = 0x0002
	QualityNormal    uint64 = 0x0003
	QualityRawFine   uint64 = 0x0004
	QualityRawNormal uint64 = 0x0005
)

// Quality reports what the camera produces per frame.
func (c *Camera) Quality() (uint64, error) {
	v, err := c.GetPropValue(PropQuality, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// SetQuality chooses RAW, JPEG or both.
//
// This is the biggest shot-to-shot lever: a 40 MP RAW is ~37 MB and takes most
// of a second to transfer, where a JPEG is a fraction of that. RAW+JPEG is the
// slowest of all, since both files cross the wire.
func (c *Camera) SetQuality(q uint64) error {
	if v, err := c.Quality(); err == nil && v == q {
		return nil
	}
	return c.SetPropValue(PropQuality, ptp.TypeUint16, q)
}

// Values for PropRawCompression (0xD022) — the camera's RAW RECORDING setting.
// RAW is not one size: the X-T5 can write it uncompressed, losslessly
// compressed, or lossily compressed, and the choice roughly halves and halves
// again the bytes that cross USB without touching image quality in the lossless
// case.
//
// Ordering follows the SDK's XT5_RAW_COMPRESSION_{OFF,LOSSLESS,LOSSY} enum. The
// 1-based numbering was originally inferred from the sibling XSDK_MEDIAREC_*
// constants and is now HARDWARE-CONFIRMED on an X-T5 (2026-08-07): writing 1
// then 2 produced an 82,378,240-byte uncompressed RAF and a 25,206,944-byte
// losslessly compressed one, and the body read the setting back each time.
const (
	RawUncompressed uint64 = 0x0001 // UNCOMPRESSED
	RawLossless     uint64 = 0x0002 // LOSSLESS COMPRESSED — full quality, ~half the size
	RawLossy        uint64 = 0x0003 // COMPRESSED — smallest, and lossy
)

// RawCompression reports how the camera packs its RAW files.
func (c *Camera) RawCompression() (uint64, error) {
	v, err := c.GetPropValue(PropRawCompression, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// SetRawCompression chooses uncompressed, lossless or lossy RAW.
//
// RawLossless is the one to want for a continuous sequence: it is mathematically
// identical to uncompressed once decoded, but roughly half the bytes, and the
// transfer is the bottleneck. RawUncompressed buys nothing but time.
func (c *Camera) SetRawCompression(mode uint64) error {
	if v, err := c.RawCompression(); err == nil && v == mode {
		return nil
	}
	return c.SetPropValue(PropRawCompression, ptp.TypeUint16, mode)
}

// FreeBuffer reports how many more frames the camera's volatile store can hold
// before it is full.
//
// This is the real limit on a continuous sequence, and it is small: an X-T5
// reports 5. Every tethered frame lands in that store and stays there until the
// host downloads it, so shooting faster than you download drains this to zero
// and the camera stops accepting captures. Poll it between frames; if it is
// falling, the downloader is behind the shutter.
func (c *Camera) FreeBuffer() (uint64, error) {
	v, err := c.GetPropValue(PropFreeSDRAMImages, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// Exposure program modes for ptp.PropExposureProgram (0x500E), the standard PTP
// codes. An X-T5 advertises 1, 3, 4 and 6.
//
// This is what decides whether the HOST or the CAMERA owns each side of the
// exposure. Fujifilm derives it from the physical dials — shutter dial on T with
// the aperture ring on A gives shutter priority, and in that mode the camera
// picks the f-number and FNumber advertises exactly one value, so writes to it
// are accepted and ignored. Setting ProgramManual hands both back to the host
// without touching the dials.
const (
	ProgramManual           uint64 = 0x0001 // M: host sets shutter AND aperture
	ProgramAperturePriority uint64 = 0x0003 // A: host sets aperture
	ProgramShutterPriority  uint64 = 0x0004 // S: host sets shutter
	ProgramAuto             uint64 = 0x0006 // P: camera sets both
)

// ExposureProgram reports which side of the exposure the host controls.
func (c *Camera) ExposureProgram() (uint64, error) {
	v, err := c.GetPropValue(ptp.PropExposureProgram, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// SetExposureProgram selects manual, aperture priority, shutter priority or
// program. ProgramManual is what an eclipse sequence wants: every component
// under host control, nothing re-metered between frames.
func (c *Camera) SetExposureProgram(mode uint64) error {
	if v, err := c.ExposureProgram(); err == nil && v == mode {
		return nil
	}
	return c.SetPropValue(ptp.PropExposureProgram, ptp.TypeUint16, mode)
}

// Focus modes for ptp.PropFocusMode (0x500A). An X-T5 advertises exactly these
// three.
const (
	FocusManual uint64 = 0x0001 // MF — the half press does not hunt
	FocusAFS    uint64 = 0x8001 // AF-S, single
	FocusAFC    uint64 = 0x8002 // AF-C, continuous
)

// FocusMode reports how the body focuses on a half press.
func (c *Camera) FocusMode() (uint64, error) {
	v, err := c.GetPropValue(ptp.PropFocusMode, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// SetFocusMode selects manual or autofocus for subsequent captures.
//
// FocusManual matters for more than focus: in an AF mode the mandatory half
// press starts a focus hunt, and on a subject the lens cannot lock onto — a
// dark sky, or the sun through a solar filter — that hunt never finishes. The
// camera then reports AF-busy indefinitely and answers DeviceBusy or nothing at
// all to the full press. Eclipse work is manually focused anyway, so set MF and
// the half press becomes the metering-only step it needs to be.
func (c *Camera) SetFocusMode(mode uint64) error {
	if v, err := c.FocusMode(); err == nil && v == mode {
		return nil
	}
	return c.SetPropValue(ptp.PropFocusMode, ptp.TypeUint16, mode)
}

// SetReleaseMode writes the shutter gesture that the next InitiateCapture will
// perform.
func (c *Camera) SetReleaseMode(mode uint64) error {
	return c.SetPropValue(PropReleaseGesture, ptp.TypeUint16, mode)
}

// AFStatus reads the camera's focus state.
func (c *Camera) AFStatus() (uint64, error) {
	v, err := c.GetPropValue(PropAFStatusCode, ptp.TypeUint16)
	if err != nil {
		return 0, err
	}
	return v.Num, nil
}

// Capture takes one photograph and waits for the camera to finish.
//
// The sequence is always S1 (half press) then S2 (full press and release):
// Fujifilm's release gesture for a frame is defined relative to the half-pressed
// state, so the half press cannot be skipped. Whether that half press drives
// autofocus is governed by the body's focus mode (set MF on the camera, or via
// ptp.PropFocusMode), not by this call.
//
// It polls AF status rather than returning immediately, because that is where
// the camera reports a focus failure: InitiateCapture itself answers OK and
// then takes no picture. timeout bounds the poll, and must exceed the exposure
// — a 30-second shutter needs at least that.
func (c *Camera) CaptureFrame(autofocus bool, timeout time.Duration) error {
	_ = autofocus // focus behaviour comes from the body's focus mode, not from here
	// The half-press is MANDATORY, not an autofocus nicety: ReleaseShootS1Off
	// (0x0304) means "full press FROM THE HALF-PRESSED STATE", so a camera that
	// never saw S1 accepts the write, reports success and takes no picture.
	// Confirmed on an X-T5 — without this step nothing is ever produced.
	if err := c.SetReleaseMode(ReleaseS1On); err != nil {
		return fmt.Errorf("fuji: half press: %w", err)
	}
	if _, _, err := c.Do(ptp.OpInitiateCapture, []uint32{0, 0}, nil, timeout); err != nil {
		return fmt.Errorf("fuji: half press capture: %w", err)
	}
	// From here the shutter button is HELD DOWN. Every failure path below must
	// let it up again, or the body stays half-pressed and answers DeviceBusy to
	// the next capture — a wedge that survives the session and looks like a
	// protocol bug on the following run.
	fail := func(format string, err error) error {
		c.ReleaseAll()
		return fmt.Errorf(format, err)
	}
	// The body needs a moment between the half press and the full press: sent
	// back-to-back the camera accepts both and takes no picture.
	//
	// This used to be a flat 300ms on every frame. It does not need to be: the
	// waitIdle below already polls AF status, so the sleep is only guarding the
	// race where the body has not yet reported BUSY and waitIdle returns
	// immediately. Watching for that transition costs a few milliseconds
	// instead of 300, and 300ms per frame is a third of the whole capture on a
	// short exposure.
	c.settleAfterS1()
	c.waitIdle(5 * time.Second)

	if err := c.SetReleaseMode(ReleaseShootS1Off); err != nil {
		return fail("fuji: setting release mode: %w", err)
	}
	// Even after AF reports idle the body can still be briefly busy, so treat
	// DeviceBusy as "not yet" rather than as a failure.
	if err := c.initiateWhenReady(timeout); err != nil {
		return fail("fuji: initiating capture: %w", err)
	}
	if err := c.waitAF(timeout, autofocus); err != nil {
		return fail("%w", err)
	}
	return nil
}

// settleAfterS1 waits just long enough for the half press to take effect.
//
// It watches for the body to report BUSY, which is the signal that S1 has
// actually landed; once it has, waitIdle takes over and waits for the cycle to
// finish. A body that never reports busy — manual focus, or one that does not
// expose AF status — is already settled, so this returns in a few milliseconds
// rather than sleeping through a fixed delay that exists for the AF case.
//
// The floor is not zero: back-to-back writes are accepted and produce no
// photograph, which is a silent failure and the one outcome to avoid.
func (c *Camera) settleAfterS1() {
	const (
		floor = 20 * time.Millisecond
		bound = 300 * time.Millisecond
		poll  = 10 * time.Millisecond
	)
	start := time.Now()
	time.Sleep(floor)
	for time.Since(start) < bound {
		if st, err := c.AFStatus(); err != nil || st == AFStatusBusy {
			// Either the cycle has started — waitIdle will see it through — or
			// the body does not answer, in which case waiting longer for a
			// signal that is not coming buys nothing.
			return
		}
		time.Sleep(poll)
	}
}

// waitIdle blocks until the camera stops reporting AF-busy, or the bound
// expires. It never fails: a body that does not expose AF status, or one in
// manual focus, simply has nothing to wait for.
func (c *Camera) waitIdle(bound time.Duration) {
	deadline := time.Now().Add(bound)
	for time.Now().Before(deadline) {
		st, err := c.AFStatus()
		if err != nil || st != AFStatusBusy {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// initiateWhenReady sends InitiateCapture, retrying while the camera answers
// DeviceBusy.
//
// DeviceBusy is not a failure, it is backpressure: the body is still metering,
// focusing or writing the previous frame. Treating it as fatal is what made
// capture look intermittent — it succeeded only when the camera happened to be
// ready at the instant the full press went out.
func (c *Camera) initiateWhenReady(timeout time.Duration) error {
	deadline := time.Now().Add(10 * time.Second)
	for {
		_, _, err := c.Do(ptp.OpInitiateCapture, []uint32{0, 0}, nil, timeout)
		if err == nil {
			return nil
		}
		var pe *ptp.Error
		if !errors.As(err, &pe) || pe.Code != ptp.RespDeviceBusy || time.Now().After(deadline) {
			return err
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// ReleaseAll lets go of everything the host may have latched.
//
// AFStatus (0xD209) is a BITFIELD, not an enum: the SDK documents S1, BULB, AF,
// AEL, AFL, WBL and SHOOTING bits. Releasing only N_S1Off therefore clears the
// half press and leaves anything else engaged — an X-T5 was found holding
// AFStatus 4 (AF in progress) with its lens stopped down in depth-of-field
// preview, which no amount of releasing S1 would clear.
//
// So this releases AF and S1 together, then AE and bulb. Both combinations are
// advertised by the body; the full 0x0007 is not, which is why it is two
// writes rather than one.
//
// Errors are ignored: this runs on recovery paths, and a camera that cannot be
// talked to cannot be released either.
func (c *Camera) ReleaseAll() {
	for _, gesture := range []uint64{
		ReleaseNAFOff | ReleaseNS1Off,                   // 0x0005 AF + half press
		ReleaseNAEOff | ReleaseNS1Off | ReleaseNBulbOff, // 0x000E AE + bulb
	} {
		if err := c.SetReleaseMode(gesture); err != nil {
			continue
		}
		c.Do(ptp.OpInitiateCapture, []uint32{0, 0}, nil, ptp.DefaultTimeout)
	}
}

// waitAF polls until the camera stops reporting busy.
func (c *Camera) waitAF(timeout time.Duration, autofocus bool) error {
	deadline := time.Now().Add(timeout)
	for {
		st, err := c.AFStatus()
		if err != nil {
			// A body that does not expose AF status has nothing to poll; the
			// capture already went out, so this is not a failure.
			return nil
		}
		switch st {
		case AFStatusOK:
			return nil
		case AFStatusFailed:
			if autofocus {
				return fmt.Errorf("fuji: capture failed: the camera could not focus. " +
					"Pass autofocus=false to shoot at the current focus setting")
			}
			return fmt.Errorf("fuji: capture failed (AF status 3)")
		case AFStatusBusy:
			if time.Now().After(deadline) {
				return fmt.Errorf("fuji: camera still busy after %v; "+
					"the timeout must exceed the exposure time", timeout)
			}
			time.Sleep(50 * time.Millisecond)
		default:
			// Anything else means it is no longer working on the shot.
			return nil
		}
	}
}

// BulbCapture holds the shutter open for d.
//
// DOES NOT WORK YET ON AN X-T5.
//
// USE THE SHUTTER LADDER INSTEAD. With the dial on T rather than B, ExposureTime
// is host-settable across 76 values from 5us to 64s, which is verified working.
// B hands the timing to the camera's own bulb timer — ExposureTime then reports
// a single locked value, and the SDK sets bulb through that same property, so
// there is nothing for a host to write.
func (c *Camera) BulbCapture(d time.Duration) error {
	fail := func(format string, err error) error {
		c.ReleaseAll()
		return fmt.Errorf(format, err)
	}
	// Open with BULB_ON (0x0400) and StopBulb loads 0x0C. An earlier attempt
	// used BULBS2_ON (0x0500) on the strength of the manual's wording and the
	// camera answered GeneralError.
	//
	if err := c.SetReleaseMode(ReleaseBulbOn); err != nil {
		return fmt.Errorf("fuji: arming bulb: %w", err)
	}
	opened := time.Now()
	if err := c.initiateWhenReady(ptp.DefaultTimeout); err != nil {
		return fail("fuji: opening the shutter: %w", err)
	}

	// Hold. Subtract what the open already cost so the exposure is the length
	// asked for rather than that plus the round trips.
	if rest := d - time.Since(opened); rest > 0 {
		time.Sleep(rest)
	}

	if err := c.SetReleaseMode(ReleaseBulbEnd); err != nil {
		return fail("fuji: closing the shutter: %w", err)
	}
	if _, _, err := c.Do(ptp.OpInitiateCapture, []uint32{0, 0}, nil, d+30*time.Second); err != nil {
		return fail("fuji: closing the shutter: %w", err)
	}
	return nil
}

// ---- Exposure ----------------------------------------------------------

// Shutter speed is carried in ptp.PropExposureTime (the standard PTP 0x500D) as a
// plain microsecond count.
//
// The wire values are exact powers of two where the camera's dial shows the
// conventional rounded number: 32000000 is 32 seconds and reads "30s" on the
// body, 16000000 is 16 seconds and reads "15s". Do not expect the displayed
// number back.
const (
	// MinShutter and MaxShutter bound what a uint32 microsecond count holds.
	MinShutter = time.Microsecond
	MaxShutter = time.Duration(4294967295) * time.Microsecond // ~71 minutes
)

// EncodeShutter converts an exposure time to its wire value in microseconds.
func EncodeShutter(d time.Duration) uint64 { return uint64(d.Microseconds()) }

// DecodeShutter converts a wire value back to a duration.
func DecodeShutter(v uint64) time.Duration { return time.Duration(v) * time.Microsecond }

// ValidateShutter reports whether d can be represented.
func ValidateShutter(d time.Duration) error {
	if d < MinShutter {
		return fmt.Errorf("fuji: shutter speed %v is below the 1us resolution of the wire format", d)
	}
	if d > MaxShutter {
		return fmt.Errorf("fuji: shutter speed %v exceeds the format maximum of %v; use Bulb", d, MaxShutter)
	}
	return nil
}

// ShutterSpeeds returns the exposure times the camera currently accepts,
// fastest first. The set is dial- and mode-dependent: with the shutter dial on
// a fixed speed a body reports exactly one value, which is not a fault.
func ShutterSpeeds(d *ptp.StdPropDesc) []time.Duration {
	if d == nil {
		return nil
	}
	out := make([]time.Duration, 0, len(d.Enum))
	for _, v := range d.Enum {
		out = append(out, DecodeShutter(v))
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// NearestShutter returns the wire value for the advertised speed closest to
// want, compared as durations by ratio.
//
// Unlike Sony's packed format this one is monotonic, so a numeric comparison
// would also work — but ratio matching is still what you want, since a
// millisecond matters at 1/1000 and does not at 30 seconds.
func NearestShutter(d *ptp.StdPropDesc, want time.Duration) uint64 {
	if d == nil || len(d.Enum) == 0 || want <= 0 {
		return EncodeShutter(want)
	}
	best, bestRatio := d.Enum[0], 0.0
	first := true
	for _, v := range d.Enum {
		t := DecodeShutter(v)
		if t <= 0 {
			continue
		}
		r := float64(t) / float64(want)
		if r < 1 {
			r = 1 / r
		}
		if first || r < bestRatio {
			best, bestRatio, first = v, r, false
		}
	}
	return best
}

// SetShutter sets the exposure time, snapping to the nearest speed the camera
// advertises in d.
func (c *Camera) SetShutterFrom(d *ptp.StdPropDesc, want time.Duration) error {
	if err := ValidateShutter(want); err != nil {
		return err
	}
	if d != nil && !d.Writable() {
		return fmt.Errorf("fuji: camera reports shutter speed as not settable; " +
			"the shutter dial needs to be on a host-controlled position")
	}
	return c.SetPropValue(ptp.PropExposureTime, ptp.TypeUint32, NearestShutter(d, want))
}

// Aperture uses the same F×100 convention Sony does: f/4 is 400.
func EncodeAperture(f float64) uint64 { return uint64(f*100 + 0.5) }

// DecodeAperture converts a wire value to an f-number.
func DecodeAperture(v uint64) float64 { return float64(v) / 100 }

// SetAperture sets the f-number.
func (c *Camera) SetAperture(f float64) error {
	return c.SetPropValue(ptp.PropFNumber, ptp.TypeUint16, EncodeAperture(f))
}

// ISOAuto is the value the Fuji vendor ExposureIndex (0xD02A) reports for
// automatic sensitivity. That property is informational only — see ptp.PropISO.
const ISOAuto uint64 = 0x8001

// SetISO sets sensitivity via the standard ExposureIndex property (ptp.PropISO).
//
// The width comes from the camera rather than from a guess: an X-T5 declares
// this one as INT32. Writing the wrong width is rejected, and any property
// whose declared type we can read, we should.
func (c *Camera) SetISO(iso uint32) error {
	t := ptp.TypeInt32
	if d, err := c.GetPropDesc(ptp.PropISO); err == nil && d.Type != 0 {
		t = d.Type
	}
	return c.SetPropValue(ptp.PropISO, t, uint64(iso))
}

// ISO reads the current sensitivity.
func (c *Camera) ISO() (uint64, error) {
	d, err := c.GetPropDesc(ptp.PropISO)
	if err != nil {
		return 0, err
	}
	return d.Current, nil
}

// Exposure is the current exposure triangle as the camera reports it.
type Exposure struct {
	Shutter  time.Duration
	Aperture float64
	ISO      uint64

	// Settable records, per component, whether the camera will accept a write.
	// An X-T5 with its dials on fixed positions advertises exactly one value
	// for each and refuses host control; the dials must be on A/T/C first.
	ShutterSettable, ApertureSettable, ISOSettable bool
}

func (e Exposure) String() string {
	// Sensitivity is a SIGNED property (INT32 on an X-T5), and the parser
	// sign-extends, so -1 arrives as all-ones. Printing that unsigned gives
	// 18446744073709551615, which is not a helpful way to say "the camera is
	// not reporting an ISO" — which is what it means in Bulb.
	// Sensitivity is SIGNED (INT32 on an X-T5) and the parser sign-extends, so
	// its negative values arrive as all-ones. They are not errors: -1, -2 and
	// -3 are the three AUTO presets, which is what the body reports with its
	// ISO dial on A.
	iso := ValueName(ptp.PropISO, e.ISO)
	if e.ISO == ISOAuto {
		iso = "auto"
	}
	mark := func(ok bool) string {
		if ok {
			return ""
		}
		return " (locked)"
	}
	return fmt.Sprintf("%v%s  f/%.1f%s  ISO %s%s",
		e.Shutter, mark(e.ShutterSettable),
		e.Aperture, mark(e.ApertureSettable),
		iso, mark(e.ISOSettable))
}

// ReadExposure reports the current shutter, aperture and ISO, and whether each
// is under host control.
func (c *Camera) ReadExposure() (*Exposure, error) {
	e := &Exposure{}
	sh, err := c.GetPropDesc(ptp.PropExposureTime)
	if err != nil {
		return nil, fmt.Errorf("reading shutter speed: %w", err)
	}
	e.Shutter, e.ShutterSettable = DecodeShutter(sh.Current), sh.Writable() && len(sh.Enum) > 1

	ap, err := c.GetPropDesc(ptp.PropFNumber)
	if err != nil {
		return nil, fmt.Errorf("reading aperture: %w", err)
	}
	e.Aperture, e.ApertureSettable = DecodeAperture(ap.Current), ap.Writable() && len(ap.Enum) > 1

	is, err := c.GetPropDesc(ptp.PropISO)
	if err != nil {
		return nil, fmt.Errorf("reading ISO: %w", err)
	}
	e.ISO, e.ISOSettable = is.Current, is.Writable() && len(is.Enum) > 1
	return e, nil
}

// Drive modes for PropDriveMode (0xD201), the SDK's XSDK_DRIVE_MODE_* values.
//
// An X-T5 in stills mode offers only Single, MultiExposure and PixelShift.
// Movie is NOT selectable here: the body has a physical STILL/MOVIE collar
// under the drive dial, and the whole movie property surface — 82 of the
// camera's 263 properties — is only describable once that is moved.
const (
	DriveSingle        uint64 = 0x0004
	DriveMultiExposure uint64 = 0x0005
	DriveContinuousH   uint64 = 0x0002
	DriveContinuousL   uint64 = 0x0003
	DriveMovie         uint64 = 0x0008
	DriveBracketAE     uint64 = 0x000A
	DriveBracketISO    uint64 = 0x000B
	DriveBracketFocus  uint64 = 0x000F
	DrivePixelShift    uint64 = 0x0010
)

// PropDriveMode selects single, continuous, bracketing or pixel-shift.
const PropDriveMode ptp.Prop = 0xD201

// DriveMode reports the current drive mode.
func (c *Camera) DriveMode() (uint64, error) { return c.GetProp(PropDriveMode) }

// SetDriveMode selects a drive mode, refusing one the body does not currently
// offer rather than writing a value it will silently ignore.
func (c *Camera) SetDriveMode(m uint64) error { return c.SetProp(PropDriveMode, m) }

// Ready reports whether the body is in a state where a host can drive it, and
// if not, why.
//
// This exists because the alternative is to attempt a capture and interpret the
// wreckage. Three different conditions stop a host taking control, and they
// need three different fixes:
//
//   - The body is in playback or a menu. PriorityMode is then not even
//     READABLE, answering vendor code 0xA001. Half-press the shutter.
//   - The volatile store holds undownloaded frames. PriorityMode reads fine but
//     the write is refused. Download or discard them.
//   - The camera is asleep or gone, which reads as a transport failure.
//
// Checking beforehand costs one property read.
func (c *Camera) Ready() (bool, string) {
	if _, err := c.GetPropValue(PropPriorityModeCode, ptp.TypeUint16); err != nil {
		var pe *ptp.Error
		if errors.As(err, &pe) {
			return false, "the body is not in shooting mode — it is showing a menu or " +
				"playback. Half-press the shutter to return it"
		}
		return false, fmt.Sprintf("the camera is not answering: %v", err)
	}
	if n, err := c.pendingFrames(); err == nil && n > 0 {
		return false, fmt.Sprintf("%d undownloaded frame(s) are in the buffer; the camera "+
			"will not hand over control until they are downloaded or discarded", n)
	}
	return true, ""
}

// Live view.
//
// The SDK's StartLiveView is the standard PTP InitiateOpenCapture (0x101C) and
// StopLiveView is TerminateOpenCapture (0x1018) — recovered from the transport
// library, where CPTPController::initiateOpenCapture loads 0x101C and
// terminateOpenCapture loads 0x1018. The per-model plugin only names a command
// class, so the opcodes are not visible there.
//
// While it runs, preview frames appear in LiveStore rather than the Still
// store, and are fetched the same way as any other object.

// StartLiveView begins streaming preview frames into LiveStore.
//
// It needs host priority, like any other capture operation.
func (c *Camera) StartLiveView() error {
	if _, _, err := c.Do(ptp.OpInitiateOpenCapture, []uint32{0, 0}, nil, ptp.DefaultTimeout); err != nil {
		return fmt.Errorf("fuji: starting live view: %w", err)
	}
	return nil
}

// StopLiveView ends the stream. It is safe to call when live view is not
// running; a camera that was not streaming simply refuses.
func (c *Camera) StopLiveView() error {
	if _, _, err := c.Do(ptp.OpTerminateOpenCapture, []uint32{0}, nil, ptp.DefaultTimeout); err != nil {
		return fmt.Errorf("fuji: stopping live view: %w", err)
	}
	return nil
}

// LiveFrame fetches the most recent preview frame, or nil if none is waiting.
//
// The frame is a JPEG. Unlike a captured RAW it is small and immediately
// readable — the settle that a full frame needs does not apply, because the
// camera is not committing it anywhere.
func (c *Camera) LiveFrame() ([]byte, error) {
	h, err := c.GetObjectHandles(LiveStore, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("fuji: listing live frames: %w", err)
	}
	if len(h) == 0 {
		return nil, nil
	}
	// The newest is last; older previews are of no interest.
	newest := h[len(h)-1]
	data, err := c.GetObject(newest)
	if err != nil {
		return nil, fmt.Errorf("fuji: fetching a live frame: %w", err)
	}
	// Consume it. The live store behaves like the still store: an object sits
	// there until it is deleted, and leaving previews to accumulate stops the
	// camera answering after the first one.
	for _, x := range h {
		c.DeleteObject(x, 0)
	}
	return data, nil
}

// CaptureHold holds the shutter down for d, then releases it.
//
// DOES NOT WORK ON AN X-T5: the body refuses ReleaseS2 (0x0300) with
// InvalidDevicePropValue, because it does not advertise a bare press. Its
// gesture set offers only the combined ReleaseShootS1Off (0x0304), which
// presses and releases together — so ONE trigger yields ONE frame however the
// drive dial is set. Verified: in ContinuousH, a normal capture produced a
// single frame.
//
// There is therefore no host-driven burst on this body, and the camera's
// continuous drive rate is unreachable over PTP. The SDK does expose separate
// XSDK_RELEASE_EX_S2_ON / S2_OFF through its ReleaseEx call, so the capability
// exists in the family, but it does not map onto this property and the X-T5's
// advertised gesture set has no equivalent.
//
// Kept because the gesture split is right in principle and another body may
// advertise 0x0300.
func (c *Camera) CaptureHold(d time.Duration) error {
	fail := func(format string, err error) error {
		c.ReleaseAll()
		return fmt.Errorf(format, err)
	}
	if err := c.SetReleaseMode(ReleaseS1On); err != nil {
		return fmt.Errorf("fuji: half press: %w", err)
	}
	if _, _, err := c.Do(ptp.OpInitiateCapture, []uint32{0, 0}, nil, ptp.DefaultTimeout); err != nil {
		return fail("fuji: half press capture: %w", err)
	}
	time.Sleep(300 * time.Millisecond)
	c.waitIdle(5 * time.Second)

	if err := c.SetReleaseMode(ReleaseS2); err != nil {
		return fail("fuji: holding the shutter: %w", err)
	}
	if err := c.initiateWhenReady(ptp.DefaultTimeout); err != nil {
		return fail("fuji: holding the shutter: %w", err)
	}
	time.Sleep(d)

	if err := c.SetReleaseMode(ReleaseNS1Off); err != nil {
		return fail("fuji: releasing the shutter: %w", err)
	}
	if _, _, err := c.Do(ptp.OpInitiateCapture, []uint32{0, 0}, nil, ptp.DefaultTimeout); err != nil {
		return fail("fuji: releasing the shutter: %w", err)
	}
	return nil
}
