// sonyprobe is the bring-up and diagnosis tool for ptp/sony.
//
// It is built for testing by proxy: someone with the camera runs one command,
// and it writes a capture bundle holding every raw payload the body sent. That
// bundle is enough to re-parse the whole session offline with -replay, so a
// decode bug can be found and fixed without the hardware present.
//
// It is read-only. Nothing is written to the camera unless -set is passed.
package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mikefsq/ptp"
	"github.com/mikefsq/ptp/sony"
	"github.com/mikefsq/ptp/usb"
)

// opts is every command-line setting, gathered so the call chain does not grow
// a parameter per flag.
type opts struct {
	serial   string
	list     bool
	reset    bool
	liveview int
	lvDir    string
	dump     string
	replay   string
	all      bool
	setSpec  string
	ls       bool
	get      string
	maxGet   int
	fast     bool
	thumbs   string
	events   time.Duration
	capture  bool
	bulb     time.Duration
	del      string
	yes      bool

	halfPress   bool
	halfRelease bool
	fullPress   bool
	fullRelease bool
	releaseAll  bool
	af          time.Duration

	// Exposure and focus. Pointers so "not given" is distinguishable from a
	// deliberate zero — -ev 0 is a real request.
	aperture *float64
	shutter  *time.Duration
	iso      *int
	isoAuto  bool
	ev       *float64
	focus    *int
	snap     bool
}

func main() {
	var o opts
	flag.StringVar(&o.serial, "serial", "", "USB serial of the body to open (required when more than one is attached)")
	flag.BoolVar(&o.list, "list", false, "list attached Sony devices and exit")
	flag.IntVar(&o.liveview, "liveview", 0, "fetch this many live view frames")
	flag.StringVar(&o.lvDir, "liveview-dir", "", "write live view frames here as JPEG (default: report only)")
	flag.BoolVar(&o.reset, "reset", false, "clear the endpoints and reset the camera's PTP stack, then exit; recovers a body wedged by an abandoned transfer")
	flag.StringVar(&o.dump, "dump", "sonyprobe-capture", "directory to write the capture bundle to; empty disables")
	flag.StringVar(&o.replay, "replay", "", "re-parse a capture bundle offline; no camera needed")
	flag.BoolVar(&o.all, "all", false, "print every property, not just those the camera reports as supported")
	flag.StringVar(&o.setSpec, "set", "", "write a raw property: 0xNNNN=value")
	flag.BoolVar(&o.ls, "ls", false, "list the images on the camera's card")
	flag.StringVar(&o.get, "get", "", "download objects to this directory")
	flag.IntVar(&o.maxGet, "max", 3, "with -get, how many objects to download")
	flag.BoolVar(&o.fast, "fast", false, "use MTP GetObjPropList (0x9805) for bulk metadata")
	flag.StringVar(&o.thumbs, "thumbs", "", "download thumbnails to this directory")
	flag.DurationVar(&o.events, "events", 0, "watch the interrupt endpoint for this long (e.g. 10s)")
	flag.BoolVar(&o.capture, "capture", false, "take a photograph: S2 full press and release")
	flag.BoolVar(&o.halfPress, "half-press", false, "hold S1 down (autofocus and metering); stays held")
	flag.BoolVar(&o.halfRelease, "half-release", false, "let S1 up")
	flag.BoolVar(&o.fullPress, "full-press", false, "hold S2 down (fires the shutter); stays held")
	flag.BoolVar(&o.fullRelease, "full-release", false, "let S2 up")
	flag.BoolVar(&o.releaseAll, "release-all", false, "let both shutter buttons up")
	flag.DurationVar(&o.af, "af", 0, "autofocus then shoot: hold S1 for this long, then S2 (e.g. 500ms)")
	flag.DurationVar(&o.bulb, "bulb", 0, "hold the shutter open for this long; body must be in Bulb mode")
	flag.StringVar(&o.del, "delete", "", "DESTRUCTIVE: delete object 0xNNNNNNNN; requires -yes")
	flag.BoolVar(&o.yes, "yes", false, "confirm a destructive operation")

	// Exposure and focus, in real-world units. The driver does the encoding.
	ap := flag.Float64("aperture", 0, "set the f-number, e.g. 5.6 (needs A or M mode)")
	sh := flag.Duration("shutter", 0, "set the exposure time, e.g. 30s or 4ms (needs S or M mode)")
	is := flag.Int("iso", 0, "set ISO, e.g. 1600")
	evv := flag.Float64("ev", 0, "set exposure compensation in EV, e.g. -0.3")
	fo := flag.Int("focus", 0, "nudge focus by -7..7; negative is nearer")
	flag.BoolVar(&o.isoAuto, "iso-auto", false, "hand ISO back to the camera")
	flag.BoolVar(&o.snap, "exposure", false, "print the current exposure triangle and exit")
	flag.Parse()

	// Only treat a value as requested if the flag was actually given, so that
	// -ev 0 and "no -ev" are different things.
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "aperture":
			o.aperture = ap
		case "shutter":
			o.shutter = sh
		case "iso":
			o.iso = is
		case "ev":
			o.ev = evv
		case "focus":
			o.focus = fo
		}
	})

	if err := run(o); err != nil {
		fmt.Fprintln(os.Stderr, "\nerror:", err)
		if o.dump != "" && o.replay == "" {
			fmt.Fprintf(os.Stderr, "a capture bundle was still written to %s — send that directory back,\n"+
				"it holds the raw bytes needed to diagnose this without the camera.\n", o.dump)
		}
		os.Exit(1)
	}
}

// tee writes to stdout and to the capture log at once, so the operator sees
// progress and the bundle keeps a verbatim record.
type tee struct {
	w []io.Writer
}

func (t *tee) Write(p []byte) (int, error) {
	for _, w := range t.w {
		w.Write(p)
	}
	return len(p), nil
}

func run(o opts) error {
	serial, list, dumpDir, replayDir := o.serial, o.list, o.dump, o.replay
	all, setSpec, ls, getDir, maxGet := o.all, o.setSpec, o.ls, o.get, o.maxGet
	fast, thumbDir, watch := o.fast, o.thumbs, o.events
	capture, bulb, del, yes := o.capture, o.bulb, o.del, o.yes

	if replayDir != "" {
		return doReplay(replayDir, all)
	}

	var out io.Writer = os.Stdout
	var logFile *os.File
	if dumpDir != "" {
		if err := os.MkdirAll(dumpDir, 0o755); err != nil {
			return fmt.Errorf("creating capture directory: %w", err)
		}
		f, err := os.Create(filepath.Join(dumpDir, "session.log"))
		if err != nil {
			return fmt.Errorf("creating session log: %w", err)
		}
		logFile = f
		defer logFile.Close()
		out = &tee{w: []io.Writer{os.Stdout, logFile}}
	}

	p := func(format string, a ...any) { fmt.Fprintf(out, format, a...) }

	p("sonyprobe — Sony bring-up\n")
	p("time: %s\nhost: %s/%s, go %s\n\n", time.Now().Format(time.RFC3339), runtime.GOOS, runtime.GOARCH, runtime.Version())

	devs, err := usb.EnumerateVendor(ptp.Sony)
	if err != nil {
		return err
	}
	p("%d Sony device(s) attached:\n", len(devs))
	for _, d := range devs {
		p("  %s\n", d)
	}
	if len(devs) == 0 {
		return fmt.Errorf("no Sony USB device found (is the body powered on and set to PC Remote?)")
	}
	if list {
		return nil
	}
	p("\n")

	t, err := usb.OpenVendor(ptp.Sony, serial)
	if err != nil {
		return err
	}
	p("opened; bulk max packet size %d\n", t.MaxPacketSize())

	if o.reset {
		r, ok := t.(ptp.Resetter)
		if !ok {
			t.Close()
			return fmt.Errorf("this USB backend cannot reset a device")
		}
		if err := r.Reset(); err != nil {
			t.Close()
			return fmt.Errorf("device reset: %w", err)
		}
		p("device reset accepted: endpoints cleared and the camera's PTP stack told to go idle\n")
		t.Close()
		return nil
	}

	// New opens the PTP session but NOT the vendor handshake: the probe runs
	// against bodies that have no vendor surface at all, and it reports what it
	// found rather than refusing them.
	s, err := sony.New(t)
	if err != nil {
		return fmt.Errorf("opening PTP session: %w", err)
	}
	defer s.Close()

	// Trace every transaction into the log and, for anything with a data
	// phase, a raw .bin beside it. This is what makes one run sufficient.
	seq := 0
	s.Trace = func(ev ptp.TraceEvent) {
		seq++
		status := "ok"
		if ev.Err != nil {
			status = "ERR " + ev.Err.Error()
		} else if ev.RespCode != 0 && ev.RespCode != 0x2001 {
			status = ev.RespCode.String()
		}
		p("  [%03d] op=0x%04X params=%v tx=%d in=%dB out=%dB %v %s\n",
			seq, uint16(ev.Op), ev.Params, ev.TxID, len(ev.DataIn), len(ev.DataOut), ev.Duration.Round(time.Millisecond), status)
		if dumpDir != "" && len(ev.DataIn) > 0 {
			name := fmt.Sprintf("%03d-op%04X-in.bin", seq, uint16(ev.Op))
			os.WriteFile(filepath.Join(dumpDir, name), ev.DataIn, 0o644)
		}
	}

	// Standard PTP device info names the model and firmware, which pins down
	// which body a capture came from.
	p("\nGetDeviceInfo (0x1001)...\n")
	var di *ptp.DevInfo
	raw, _, derr := s.Do(ptp.OpGetDeviceInfo, nil, nil, 5*time.Second)
	if derr != nil {
		p("  GetDeviceInfo failed: %v (continuing)\n", derr)
	} else {
		if dumpDir != "" {
			os.WriteFile(filepath.Join(dumpDir, "deviceinfo.bin"), raw, 0o644)
		}
		var perr error
		di, perr = ptp.ParseDeviceInfo(raw)
		if perr != nil {
			p("  %d bytes, but parse failed: %v\n", len(raw), perr)
		} else {
			p("  %s\n", di)
			p("  operations (%d): %s\n", len(di.Operations), fmtOps(di.Operations))
			p("  device properties (%d): %s\n", len(di.DeviceProps), fmtProps(di.DeviceProps))
			p("  capture formats: %d, image formats: %d\n", len(di.CaptureFormats), len(di.ImageFormats))
			if sony.SupportsSDIO(di) {
				p("  -> exposes the Sony SDIO remote-control operations\n")
			} else {
				p("  -> NO Sony SDIO operations: this body is file-transfer only over this interface\n")
			}
		}
	}

	if o.snap || o.aperture != nil || o.shutter != nil || o.iso != nil || o.isoAuto || o.ev != nil || o.focus != nil {
		if err := doExposure(out, s, di, o); err != nil {
			return err
		}
	}
	if o.halfPress || o.halfRelease || o.fullPress || o.fullRelease || o.releaseAll {
		if err := doButtons(out, s, di, o); err != nil {
			return err
		}
	}
	if capture || bulb > 0 || o.af > 0 {
		if err := doCapture(out, s, di, bulb, o.af); err != nil {
			return err
		}
	}
	if del != "" {
		if err := doDelete(out, s, del, yes); err != nil {
			return err
		}
	}
	if fast {
		if err := browseFast(out, s); err != nil {
			p("  bulk metadata failed: %v\n", err)
		}
	}
	if ls || getDir != "" || thumbDir != "" {
		if err := browse(out, s, getDir, maxGet, thumbDir); err != nil {
			return err
		}
		if !ls {
			return nil
		}
	}
	if watch > 0 {
		watchEvents(out, s, watch)
	}
	if o.liveview > 0 {
		if err := doLiveView(out, s, o.liveview, o.lvDir); err != nil {
			return err
		}
	}

	// A body with no vendor surface can still be exercised over standard PTP.
	// That is worth doing: it validates the transport and the descriptor reader
	// against real data even when the Sony layer is unavailable.
	if di != nil && !sony.SupportsSDIO(di) {
		p("\nno vendor surface; reading standard PTP properties (0x1014)...\n")
		for _, code := range di.DeviceProps {
			d, err := s.GetPropDesc(code)
			if err != nil {
				p("  0x%04X: %v\n", uint16(code), err)
				continue
			}
			w := "read-only"
			if d.Writable() {
				w = "writable"
			}
			extra := ""
			switch {
			case d.Form == 1:
				extra = fmt.Sprintf(" range %d..%d step %d", int64(d.Min), int64(d.Max), int64(d.Step))
			case len(d.Enum) > 0:
				extra = fmt.Sprintf(" %d values %v", len(d.Enum), trim(d.Enum))
			}
			val := fmt.Sprintf("%d", int64(d.Current))
			if d.Type == 0xFFFF {
				val = fmt.Sprintf("%q", d.CurrentStr)
			}
			p("  0x%04X %-28s type=0x%04X cur=%s %s%s\n",
				uint16(d.Code), d.Code, uint16(d.Type), val, w, extra)
		}
		return fmt.Errorf("%s %s exposes no Sony remote-control operations; "+
			"standard PTP results above", di.Manufacturer, di.Model)
	}

	p("\nSony vendor handshake (0x9201 x3 + 0x9202)...\n")
	info, err := s.Connect()
	if info != nil && info.Raw != nil && dumpDir != "" {
		os.WriteFile(filepath.Join(dumpDir, "extdeviceinfo-9202.bin"), info.Raw, 0o644)
	}
	if err != nil {
		if info != nil && len(info.Raw) > 0 {
			p("  handshake payload did not parse; raw bytes follow\n%s\n", hexDump(info.Raw, 512))
		}
		return err
	}
	p("  protocol 0x%04X (mode %d) after %d attempt(s): %d properties, %d controls\n",
		info.Version, info.ModeVer(), info.Attempts, len(info.Props), len(info.Controls))

	supported := make(map[ptp.Prop]bool, len(info.Props))
	for _, c := range info.Props {
		supported[c] = true
	}

	p("\nGetAllDevicePropData (0x9209)...\n")
	props, perr := s.GetAllDevicePropData()
	if perr != nil {
		p("  PARSE FAILED after %d properties: %v\n", len(props), perr)
		p("  the raw blob is in the bundle as the 0x9209 .bin file; this is\n")
		p("  recoverable offline with -replay once the parser is corrected\n\n")
	} else {
		p("  parsed %d properties cleanly\n\n", len(props))
	}
	dumpProps(out, props, supported, all)

	if perr != nil {
		return fmt.Errorf("property blob parse failed (partial results above): %w", perr)
	}
	if setSpec != "" {
		return doSet(out, s, props, setSpec)
	}
	if dumpDir != "" {
		p("\ncapture bundle written to %s — send that directory back.\n", dumpDir)
	}
	return nil
}

// doReplay re-parses a capture bundle with no camera attached.
func doReplay(dir string, all bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading capture bundle: %w", err)
	}
	var blob string
	for _, e := range entries {
		if strings.Contains(e.Name(), "op9209") && strings.HasSuffix(e.Name(), "-in.bin") {
			blob = filepath.Join(dir, e.Name())
		}
	}
	if blob == "" {
		return fmt.Errorf("no 0x9209 payload found in %s (expected a file like NNN-op9209-in.bin)", dir)
	}
	b, err := os.ReadFile(blob)
	if err != nil {
		return err
	}
	fmt.Printf("replaying %s (%d bytes)\n\n", blob, len(b))

	props, perr := sony.ParseAllDevicePropData(b)
	if perr != nil {
		fmt.Printf("PARSE FAILED after %d properties: %v\n\n", len(props), perr)
	} else {
		fmt.Printf("parsed %d properties cleanly\n\n", len(props))
	}
	dumpProps(os.Stdout, props, nil, all)
	return perr
}

func dumpProps(out io.Writer, props []sony.DeviceProperty, supported map[ptp.Prop]bool, all bool) {
	sort.Slice(props, func(i, j int) bool { return props[i].Code < props[j].Code })
	unnamed := 0
	for _, pr := range props {
		_, known := sony.PropTable[pr.Code]
		if !known {
			unnamed++
		}
		if !all && supported != nil && !supported[pr.Code] {
			continue
		}
		flags := []string{}
		if pr.IsControl() {
			flags = append(flags, "control")
		}
		if pr.Writable() {
			flags = append(flags, "writable")
		}
		if !known {
			flags = append(flags, "NOT-IN-SDK-TABLE")
		}
		// Name the value, not just the code. A shutter speed printed as its
		// packed wire value reads as 66536, and an enum printed as an ordinal
		// tells you nothing about which mode the camera is in.
		cur := sony.ValueName(pr.Code, pr.Current)
		if pr.Type == ptp.TypeString {
			cur = fmt.Sprintf("%q", pr.CurrentStr)
		}
		fmt.Fprintf(out, "  0x%04X %-40s type=0x%04X cur=%s %s\n",
			uint16(pr.Code), sony.PropName(pr.Code), uint16(pr.Type), cur, strings.Join(flags, ","))
		switch {
		case pr.Form == 1:
			fmt.Fprintf(out, "         range %d..%d step %d\n", int64(pr.Min), int64(pr.Max), int64(pr.Step))
		case len(pr.Enum) > 0:
			fmt.Fprintf(out, "         %d values: %s\n", len(pr.Enum), sony.DescribeValues(pr.Code, pr.Enum))
		}
	}
	if unnamed > 0 {
		fmt.Fprintf(out, "\n%d properties are not in the SDK's table (newer than SDK 2.02.00, or private)\n", unnamed)
	}
}

func fmtOps(ops []ptp.OpCode) string {
	parts := make([]string, 0, len(ops))
	for _, o := range ops {
		parts = append(parts, fmt.Sprintf("0x%04X", uint16(o)))
	}
	return strings.Join(parts, " ")
}

func fmtProps(ps []ptp.Prop) string {
	parts := make([]string, 0, len(ps))
	for _, x := range ps {
		parts = append(parts, fmt.Sprintf("0x%04X", uint16(x)))
	}
	return strings.Join(parts, " ")
}

func trim(v []uint64) []int64 {
	const max = 16
	out := make([]int64, 0, max)
	for i, x := range v {
		if i == max {
			break
		}
		out = append(out, int64(x))
	}
	return out
}

// hexDump renders up to limit bytes, for pasting into a bug report when the
// bundle itself cannot be sent.
func hexDump(b []byte, limit int) string {
	if len(b) > limit {
		b = b[:limit]
	}
	return hex.Dump(b)
}

// asciiStrings pulls printable runs out of a blob — enough to read the model
// and firmware out of a PTP DeviceInfo without decoding it properly.
func asciiStrings(b []byte) []string {
	var out []string
	var cur []byte
	for _, c := range b {
		if c >= 0x20 && c < 0x7F {
			cur = append(cur, c)
			continue
		}
		if len(cur) >= 4 {
			out = append(out, string(cur))
		}
		cur = nil
	}
	if len(cur) >= 4 {
		out = append(out, string(cur))
	}
	return out
}

func doSet(out io.Writer, s *sony.Camera, props []sony.DeviceProperty, spec string) error {
	code, valStr, ok := strings.Cut(spec, "=")
	if !ok {
		return fmt.Errorf("-set wants CODE=VALUE, got %q", spec)
	}
	c, err := strconv.ParseUint(strings.TrimPrefix(strings.TrimSpace(code), "0x"), 16, 16)
	if err != nil {
		return fmt.Errorf("parsing property code %q: %w", code, err)
	}
	v, err := strconv.ParseInt(strings.TrimSpace(valStr), 0, 64)
	if err != nil {
		return fmt.Errorf("parsing value %q: %w", valStr, err)
	}

	// Take the data type from what the camera just reported rather than
	// assuming: writing the wrong width is how a body gets wedged.
	var target *sony.DeviceProperty
	for i := range props {
		if uint16(props[i].Code) == uint16(c) {
			target = &props[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("camera did not report property 0x%04X, refusing to write it", c)
	}
	if !target.Writable() {
		return fmt.Errorf("camera reports 0x%04X (%s) as not writable right now", c, target.Code)
	}

	fmt.Fprintf(out, "\nsetting %s (0x%04X) type 0x%04X to %d\n", target.Code, c, uint16(target.Type), v)
	if target.IsControl() {
		return s.SendControl(sony.ControlCode(target.Code), target.Type, uint64(v))
	}
	return s.SetProperty(target.Code, target.Type, uint64(v))
}

// browse lists the camera's storage and objects, and optionally downloads some.
func browse(out io.Writer, s *sony.Camera, getDir string, maxGet int, thumbDir string) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }

	p("\nstorage (0x1004)...\n")
	ids, err := s.GetStorageIDs()
	if err != nil {
		return fmt.Errorf("listing storage: %w", err)
	}
	for _, id := range ids {
		si, err := s.GetStorageInfo(id)
		if err != nil {
			p("  0x%08X: %v\n", id, err)
			continue
		}
		p("  0x%08X %q %q  %.1f GB total, %.1f GB free, %d images free\n",
			id, si.StorageDescription, si.VolumeLabel,
			float64(si.MaxCapacity)/1e9, float64(si.FreeSpaceInBytes)/1e9, si.FreeSpaceInImages)
	}

	p("\nobjects (0x1007)...\n")
	handles, err := s.GetObjectHandles(ptp.AllStorages, ptp.AllFormats, 0)
	if err != nil {
		return fmt.Errorf("listing objects: %w", err)
	}
	p("  %d object(s)\n", len(handles))

	// Describing every object on a full card is slow and rarely useful; cap it.
	const describe = 20
	shown := handles
	if len(shown) > describe {
		shown = shown[:describe]
	}
	var files []struct {
		h  uint32
		oi *ptp.ObjectInfo
	}
	fails := 0
	for _, h := range shown {
		oi, err := s.GetObjectInfo(h)
		if err != nil {
			p("  0x%08X: %v\n", h, err)
			// A sleeping camera answers nothing further, and each dead request
			// costs a full USB timeout. Give up rather than grind through them.
			if errors.Is(err, ptp.ErrNotResponding) {
				return fmt.Errorf("camera stopped responding — wake it, or re-seat the "+
					"cable if it has wedged: %w", err)
			}
			if fails++; fails >= 3 {
				return fmt.Errorf("camera stopped responding after %d consecutive failures "+
					"(it has probably gone to sleep) — wake it and retry", fails)
			}
			continue
		}
		fails = 0
		p("  0x%08X %s\n", h, oi)
		// Only real images are download candidates. A card root typically also
		// holds stubs like REGISTER.URL, and grabbing those instead of a photo
		// is not what anyone means by -get.
		if isImage(oi.ObjectFormat) {
			files = append(files, struct {
				h  uint32
				oi *ptp.ObjectInfo
			}{h, oi})
		}
	}
	if len(handles) > describe {
		p("  ... %d more not described (cap %d)\n", len(handles)-describe, describe)
	}

	if thumbDir != "" {
		if err := os.MkdirAll(thumbDir, 0o755); err != nil {
			return fmt.Errorf("creating thumbnail directory: %w", err)
		}
		p("\nthumbnails (0x100A)...\n")
		for i, f := range files {
			if i >= maxGet {
				break
			}
			th, err := s.GetThumb(f.h)
			if err != nil {
				p("  %s: %v\n", f.oi.Filename, err)
				continue
			}
			dst := filepath.Join(thumbDir, "thumb_"+f.oi.Filename)
			if err := os.WriteFile(dst, th, 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", dst, err)
			}
			p("  %s — %d bytes (%dx%d)\n", dst, len(th), f.oi.ThumbPixWidth, f.oi.ThumbPixHeight)
		}
	}

	if getDir == "" {
		return nil
	}
	if err := os.MkdirAll(getDir, 0o755); err != nil {
		return fmt.Errorf("creating download directory: %w", err)
	}
	n := 0
	for _, f := range files {
		if n >= maxGet {
			p("\nstopping after %d download(s); raise -max for more\n", maxGet)
			break
		}
		p("\ndownloading %s (%d bytes)...\n", f.oi.Filename, f.oi.CompressedSize)
		start := time.Now()
		data, err := s.GetObject(f.h)
		if err != nil {
			p("  failed: %v\n", err)
			continue
		}
		dst := filepath.Join(getDir, f.oi.Filename)
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", dst, err)
		}
		el := time.Since(start)
		p("  wrote %s — %d bytes in %v (%.1f MB/s)\n", dst, len(data), el.Round(time.Millisecond),
			float64(len(data))/el.Seconds()/1e6)
		if uint32(len(data)) != f.oi.CompressedSize {
			p("  WARNING: size mismatch — ObjectInfo said %d, got %d\n", f.oi.CompressedSize, len(data))
		}
		n++
	}
	return nil
}

// isImage reports whether an object format is a photo worth downloading.
func isImage(f uint16) bool {
	switch f {
	case ptp.FormatEXIFJPEG, ptp.FormatSonyARW, ptp.FormatSonyARWLegacy,
		ptp.FormatTIFF, ptp.FormatMPO, ptp.FormatUndefinedImg:
		return true
	}
	return false
}

// browseFast reads metadata for every object in one transaction, instead of a
// GetObjectInfo round trip each.
func browseFast(out io.Writer, s *sony.Camera) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }
	p("\nbulk metadata (0x9805 GetObjPropList)...\n")
	// One round trip per property, covering every object — not AllProps, which
	// is the form bodies refuse.
	want := []ptp.ObjProp{
		ptp.OPCObjectFileName, ptp.OPCObjectFormat, ptp.OPCObjectSize,
		ptp.OPCDateCreated, ptp.OPCWidth, ptp.OPCHeight,
	}
	start := time.Now()
	byH, err := s.BulkMetadata(want)
	if err != nil {
		return err
	}
	p("  %d objects x %d properties in %v (%d round trips)\n",
		len(byH), len(want), time.Since(start).Round(time.Millisecond), len(want))

	shown := 0
	for h, props := range byH {
		if shown >= 5 {
			p("  ... %d more objects\n", len(byH)-shown)
			break
		}
		p("  0x%08X", h)
		for _, k := range []ptp.ObjProp{
			ptp.OPCObjectFileName, ptp.OPCObjectSize,
			ptp.OPCWidth, ptp.OPCHeight, ptp.OPCDateCreated,
		} {
			if v, ok := props[k]; ok {
				p("  %s=%s", k, v)
			}
		}
		p("\n")
		shown++
	}
	return nil
}

// watchEvents polls the interrupt endpoint. An idle camera yields nothing,
// which is not a fault.
func watchEvents(out io.Writer, s *sony.Camera, d time.Duration) {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }
	p("\nwatching for events for %v (interrupt endpoint)...\n", d)
	deadline := time.Now().Add(d)
	n := 0
	for time.Now().Before(deadline) {
		ev, err := s.PollEvent(500 * time.Millisecond)
		if err != nil {
			if err == ptp.ErrTimeout {
				continue
			}
			p("  event poll failed: %v\n", err)
			return
		}
		p("  %v  %s\n", time.Now().Format("15:04:05.000"), ev)
		n++
	}
	if n == 0 {
		p("  no events in %v (normal for an idle camera)\n", d)
	}
}

// doCapture takes a photograph. Sony bodies drive the shutter with the S2
// button rather than PTP InitiateCapture, so that is the path used unless the
// camera advertises the standard operation and not the vendor one.
func doCapture(out io.Writer, s *sony.Camera, di *ptp.DevInfo, bulb, af time.Duration) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }

	sdio := di != nil && sony.SupportsSDIO(di)
	std := di != nil && di.Supports(ptp.OpInitiateCapture)
	if di != nil && !sdio && !std {
		return fmt.Errorf("%s %s advertises neither the Sony control surface nor "+
			"InitiateCapture (0x100E) — it cannot be told to take a picture over USB",
			di.Manufacturer, di.Model)
	}

	if bulb > 0 {
		p("\nbulb exposure for %v (body must already be in Bulb mode)...\n", bulb)
		start := time.Now()
		if err := s.BulbCapture(bulb); err != nil {
			return fmt.Errorf("bulb capture: %w", err)
		}
		p("  shutter held %v\n", time.Since(start).Round(time.Millisecond))
	} else if af > 0 {
		p("\ncapture: S1 held %v for autofocus, then S2 press/release...\n", af)
		if err := s.ShootWithAF(af); err != nil {
			return fmt.Errorf("capture with autofocus: %w", err)
		}
	} else if sdio {
		p("\ncapture: S2 full press/release (0x9207)...\n")
		if err := s.Shoot(); err != nil {
			return fmt.Errorf("capture: %w", err)
		}
	} else {
		p("\ncapture: standard PTP InitiateCapture (0x100E)...\n")
		if err := s.InitiateCapture(0, 0); err != nil {
			return fmt.Errorf("capture: %w", err)
		}
	}

	p("  waiting for the camera to announce the new image...\n")
	h, err := s.WaitForCapture(30 * time.Second)
	if err != nil {
		p("  no ObjectAdded event within 30s — the body may not raise one; re-list the card\n")
		return nil
	}
	p("  new object 0x%08X\n", h)
	if oi, err := s.GetObjectInfo(h); err == nil {
		p("  %s\n", oi)
	}
	return nil
}

// doDelete removes an object. Guarded behind -yes because it destroys data.
func doDelete(out io.Writer, s *sony.Camera, spec string, yes bool) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }
	h, err := strconv.ParseUint(strings.TrimPrefix(strings.TrimSpace(spec), "0x"), 16, 32)
	if err != nil {
		return fmt.Errorf("parsing object handle %q: %w", spec, err)
	}
	oi, err := s.GetObjectInfo(uint32(h))
	if err != nil {
		return fmt.Errorf("object 0x%08X not found: %w", h, err)
	}
	if !yes {
		return fmt.Errorf("refusing to delete %s (0x%08X, %d bytes) without -yes",
			oi.Filename, h, oi.CompressedSize)
	}
	p("\ndeleting %s (0x%08X, %d bytes)...\n", oi.Filename, h, oi.CompressedSize)
	if err := s.DeleteObject(uint32(h), 0); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	// A camera may answer OK having deleted nothing, so confirm.
	if _, err := s.GetObjectInfo(uint32(h)); err == nil {
		return fmt.Errorf("camera reported success but %s is still present", oi.Filename)
	}
	p("  deleted; the handle no longer resolves\n")
	return nil
}

// doExposure reads the current exposure triangle and applies any requested
// changes. Values are snapped to the set the camera currently advertises,
// because anything outside it is rejected.
func doExposure(out io.Writer, s *sony.Camera, di *ptp.DevInfo, o opts) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }

	if di != nil && !sony.SupportsSDIO(di) {
		return fmt.Errorf("%s %s has no Sony vendor surface, so exposure and focus "+
			"cannot be set over USB on this body", di.Manufacturer, di.Model)
	}

	props, err := s.GetAllDevicePropData()
	if err != nil {
		return fmt.Errorf("reading properties: %w", err)
	}
	byCode := make(map[ptp.Prop]*sony.DeviceProperty, len(props))
	for i := range props {
		byCode[props[i].Code] = &props[i]
	}

	p("\ncurrent exposure: %s\n", sony.ReadExposure(props))
	if o.snap {
		for _, c := range []ptp.Prop{
			sony.PropFNumber, sony.PropShutterSpeed,
			sony.PropIsoSensitivity, sony.PropExposureBiasCompensation,
		} {
			d, ok := byCode[c]
			if !ok {
				p("  %-26s not reported by this body\n", c)
				continue
			}
			w := "read-only"
			if d.Writable() {
				w = "writable"
			}
			p("  %-26s %s, %d accepted values\n", c, w, len(d.Enum))
		}
		if d, ok := byCode[sony.PropShutterSpeed]; ok {
			speeds, bulb := sony.ShutterSpeeds(d)
			if len(speeds) > 0 {
				p("\n  shutter speeds accepted (fastest first), Bulb=%v:\n   ", bulb)
				for i, sp := range speeds {
					if i > 0 && i%8 == 0 {
						p("\n   ")
					}
					p(" %v", sp)
				}
				p("\n  range: %v .. %v\n", speeds[0], speeds[len(speeds)-1])
			}
		}
		return nil
	}

	// set applies one property, snapping to the camera's list and reporting
	// when the snap moved the value.
	set := func(code ptp.Prop, t ptp.DataType, want uint64, label string, show func(uint64) string) error {
		d, ok := byCode[code]
		if !ok {
			return fmt.Errorf("%s: this body does not report %s", label, code)
		}
		if !d.Writable() {
			return fmt.Errorf("%s: camera reports %s as not settable right now "+
				"(check the exposure mode dial)", label, code)
		}
		final := sony.Nearest(d, want)
		if final != want {
			p("  %s: %s is not an accepted value; using %s\n", label, show(want), show(final))
		}
		if err := s.SetProperty(code, t, final); err != nil {
			return fmt.Errorf("%s: %w", label, err)
		}
		p("  %s -> %s\n", label, show(final))
		return nil
	}

	if o.aperture != nil {
		f := func(v uint64) string { a, _ := sony.DecodeAperture(v); return fmt.Sprintf("f/%.1f", a) }
		if err := set(sony.PropFNumber, ptp.TypeUint16,
			sony.EncodeAperture(*o.aperture), "aperture", f); err != nil {
			return err
		}
	}
	if o.shutter != nil {
		// Shutter speed must snap by duration, not by wire value: the packed
		// form runs backwards at the fast end, so a numeric comparison picks
		// the wrong speed.
		d, ok := byCode[sony.PropShutterSpeed]
		if !ok {
			return fmt.Errorf("shutter: this body does not report %s", sony.PropShutterSpeed)
		}
		if !d.Writable() {
			return fmt.Errorf("shutter: camera reports it as not settable right now " +
				"(the mode dial needs to be on S or M)")
		}
		if err := sony.ValidateShutterSpeed(*o.shutter); err != nil {
			return err
		}
		wire := sony.NearestShutter(d, *o.shutter)
		got, _ := sony.DecodeShutterSpeed(wire)
		if got != *o.shutter {
			p("  shutter: %v is not an accepted speed; using %v\n", *o.shutter, got)
		}
		if err := s.SetProperty(sony.PropShutterSpeed, ptp.TypeUint32, wire); err != nil {
			return fmt.Errorf("shutter: %w", err)
		}
		p("  shutter -> %v\n", got)
	}
	if o.isoAuto {
		if err := set(sony.PropIsoSensitivity, ptp.TypeUint32,
			sony.ISOAuto(), "iso", func(uint64) string { return "AUTO" }); err != nil {
			return err
		}
	} else if o.iso != nil {
		f := func(v uint64) string { i, _, _ := sony.DecodeISO(v); return fmt.Sprintf("%d", i) }
		if err := set(sony.PropIsoSensitivity, ptp.TypeUint32,
			sony.EncodeISO(uint32(*o.iso), sony.ISOModeNormal), "iso", f); err != nil {
			return err
		}
	}
	if o.ev != nil {
		f := func(v uint64) string { return fmt.Sprintf("%+.1f EV", sony.DecodeExposureCompensation(v)) }
		if err := set(sony.PropExposureBiasCompensation, ptp.TypeInt16,
			sony.EncodeExposureCompensation(*o.ev), "exposure compensation", f); err != nil {
			return err
		}
	}
	if o.focus != nil {
		p("  focus nudge %+d\n", *o.focus)
		if err := s.FocusNearFar(int16(*o.focus)); err != nil {
			return fmt.Errorf("focus: %w", err)
		}
	}

	// Read back: a camera can answer OK and not apply the change.
	if after, err := s.GetAllDevicePropData(); err == nil {
		p("  now: %s\n", sony.ReadExposure(after))
	}
	return nil
}

// doButtons drives the shutter buttons individually, so a caller can hold one
// across several operations rather than only firing a whole capture.
func doButtons(out io.Writer, s *sony.Camera, di *ptp.DevInfo, o opts) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }
	if di != nil && !sony.SupportsSDIO(di) {
		return fmt.Errorf("%s %s has no Sony control surface, so the shutter buttons "+
			"cannot be driven over USB on this body", di.Manufacturer, di.Model)
	}
	// Order matters: a press and a release given together should do the obvious
	// thing, and releases must come after presses.
	if o.halfPress {
		p("\nS1 down (autofocus, metering) — held until released\n")
		if err := s.HalfPress(); err != nil {
			return fmt.Errorf("half press: %w", err)
		}
	}
	if o.fullPress {
		p("S2 down (shutter) — held until released\n")
		if err := s.FullPress(); err != nil {
			return fmt.Errorf("full press: %w", err)
		}
	}
	if o.fullRelease {
		p("S2 up\n")
		if err := s.FullRelease(); err != nil {
			return fmt.Errorf("full release: %w", err)
		}
	}
	if o.halfRelease {
		p("S1 up\n")
		if err := s.HalfRelease(); err != nil {
			return fmt.Errorf("half release: %w", err)
		}
	}
	if o.releaseAll {
		p("releasing both shutter buttons\n")
		if err := s.ReleaseAll(); err != nil {
			return fmt.Errorf("release all: %w", err)
		}
	}
	return nil
}

// doLiveView fetches preview frames and reports the cadence.
//
// The whole path is UNVERIFIED: the object handle and the operations come from
// Sony's own USB adapter, but no body supporting live view has been driven. So
// this reports what it actually got — sizes, dimensions, whether a JPEG was
// found at all — rather than asserting success, because the first real run is
// the experiment.
func doLiveView(out io.Writer, s *sony.Camera, n int, dir string) error {
	p := func(format string, a ...any) { fmt.Fprintf(out, format, a...) }

	p("\nlive view: checking whether the camera is delivering...\n")
	running, delivering, err := s.LiveViewAvailable()
	if err != nil {
		return fmt.Errorf("reading live view status: %w", err)
	}
	p("  LiveViewStatus(0xD221)=%v  MonitoringIsDelivering(0xE099)=%v\n", running, delivering)
	if !delivering {
		if err := s.StartLiveView(5 * time.Second); err != nil {
			p("  %v\n", err)
			p("  fetching anyway — the status properties may not mean what we think\n")
		}
	}

	if oi, err := s.LiveFrameInfo(); err != nil {
		p("  ObjectInfo on the live view handle failed: %v\n", err)
	} else {
		p("  frame object: %s\n", oi)
	}

	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	var total int
	start := time.Now()
	for i := range n {
		t0 := time.Now()
		frame, err := s.LiveFrame()
		if err != nil {
			p("  frame %d: %v\n", i, err)
			return err
		}
		total += len(frame)
		p("  frame %d: %d bytes in %v\n", i, len(frame), time.Since(t0).Round(time.Millisecond))
		if dir != "" {
			name := filepath.Join(dir, fmt.Sprintf("live-%03d.jpg", i))
			if err := os.WriteFile(name, frame, 0o644); err != nil {
				return err
			}
		}
	}
	if el := time.Since(start); n > 0 && el > 0 {
		p("  %d frames, %d bytes, %.1f fps\n", n, total, float64(n)/el.Seconds())
	}
	return nil
}
