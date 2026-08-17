// fujiprobe is the bring-up tool for the fujicam driver.
//
// It enumerates attached Fujifilm bodies, opens a PTP session, and dumps what
// the camera reports. Like sonyprobe, a run writes a capture bundle holding
// every raw payload, so a decode problem can be worked offline.
//
// It is read-only.
package main

import (
	"crypto/md5"
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
	"github.com/mikefsq/ptp/fuji"
	"github.com/mikefsq/ptp/usb"
)

func main() {
	var (
		serial = flag.String("serial", "", "USB serial of the body to open (required when more than one is attached)")
		list   = flag.Bool("list", false, "list attached Fujifilm devices and exit")
		dump   = flag.String("dump", "fujiprobe-capture", "directory to write the capture bundle to; empty disables")
		all    = flag.Bool("all", false, "describe every property the camera reports, not just the named ones")
		shoot  = flag.Bool("capture", false, "take a photograph (S1 half press, then S2)")
		get    = flag.String("get", "", "download captured frames to this directory")
		card   = flag.String("card", "", "also record tethered frames to the SD card: rawjpeg, raw, jpeg or off")
		qual   = flag.String("quality", "", "what the camera produces per frame: raw, fine, normal, raw+fine or raw+normal")
		rawc   = flag.String("raw", "", "RAW recording: uncompressed, lossless or lossy")

		// Exposure, matching sonyprobe's flags.
		sh      = flag.Duration("shutter", 0, "set the exposure time, e.g. 10s or 1ms (shutter dial must be on T)")
		ap      = flag.Float64("aperture", 0, "set the f-number, e.g. 5.6 (aperture ring must be on A)")
		is      = flag.Int("iso", 0, "set ISO, e.g. 1600 (ISO dial must be on C)")
		snap    = flag.Bool("exposure", false, "print the current exposure triangle and exit")
		rel     = flag.Bool("release", false, "let the shutter button up; recovers a body answering DeviceBusy")
		mf      = flag.Bool("mf", false, "switch the body to manual focus, so the half press does not hunt")
		mode    = flag.String("mode", "", "exposure program: manual, aperture, shutter or auto")
		prop    = flag.String("prop", "", "read one property by name or 0xNNNN code")
		setstr  = flag.String("setstr", noWrite, "with -prop, write this string value (pass an empty string to clear the property)")
		snap2   = flag.String("snapshot", "", "read every property and write it to this file")
		diff    = flag.String("diff", "", "read every property and report what changed since this snapshot")
		rawp    = flag.String("rawprop", "", "fetch a property's raw GetDevicePropValue payload and hex-dump it")
		syncCl  = flag.Bool("syncclock", false, "set the camera's clock from this host")
		setv    = flag.Int64("set", -1, "with -prop, write this value (validated against what the camera offers)")
		hold    = flag.Duration("hold", 0, "hold the shutter down for this long; in a continuous drive mode this bursts")
		bulb    = flag.Duration("bulb", 0, "hold the shutter open for this long; the body's shutter dial must be on B")
		live    = flag.Int("liveview", 0, "start live view, grab N preview frames to -get, and stop")
		bracket = flag.String("bracket", "", "with -burst, cycle these shutter speeds across the frames, e.g. 1ms,4ms,15ms,60ms")
		burst   = flag.Int("burst", 0, "shoot N frames back to back in one session, reporting cadence and buffer depth")
		lsAll   = flag.Bool("ls", false, "list every storage volume the camera exposes, and what is on it")
		settle  = flag.Duration("settle", 0, "extra wait after a frame appears before downloading it. MEASURED as unnecessary on an X-T5 (see doCapture); kept because a body that does stall would need it")
		discard = flag.Bool("discard", false, "card-only shooting: after capture, delete the frame from the camera's buffer WITHOUT downloading it. Needs -card raw, or the frame is lost — the card copy is what survives")
		drain   = flag.Bool("drain", false, "DESTRUCTIVE: delete frames from the camera's volatile store after a verified download")
	)
	flag.Parse()

	ex := exposureOpts{rawprop: *rawp, prop: *prop, setv: *setv, setstr: *setstr, syncClock: *syncCl, shutter: *sh, aperture: *ap, iso: *is, show: *snap, release: *rel, mf: *mf, mode: *mode}
	if err := run(*serial, *list, *dump, *all, *shoot, *get, *card, *qual, *rawc, ex, *drain, *discard, *lsAll, *burst, *bracket, *live, *bulb, *snap2, *diff, *hold, *settle); err != nil {
		fmt.Fprintln(os.Stderr, "\nerror:", err)
		if *dump != "" {
			fmt.Fprintf(os.Stderr, "a capture bundle was still written to %s — it holds the raw\n"+
				"bytes needed to diagnose this without the camera.\n", *dump)
		}
		os.Exit(1)
	}
}

// tee writes to stdout and the capture log at once.
type tee struct{ w []io.Writer }

func (t *tee) Write(p []byte) (int, error) {
	for _, w := range t.w {
		w.Write(p)
	}
	return len(p), nil
}

// noWrite is the -setstr default, distinct from the empty string so that
// clearing a property and not writing one are different requests. A sentinel is
// needed because "" is a legitimate value.
const noWrite = "\x00<unset>"

// cardModes maps the -card flag to the camera's MediaRecord values.
var cardModes = map[string]uint64{
	"rawjpeg": fuji.MediaRecordRawJPEG,
	"raw":     fuji.MediaRecordRaw,
	"jpeg":    fuji.MediaRecordJPEG,
	"off":     fuji.MediaRecordOff,
}

func cardModeName(v uint64) string {
	for n, x := range cardModes {
		if x == v {
			return n
		}
	}
	return fmt.Sprintf("0x%04X", v)
}

// qualities maps the -quality flag to the camera's Quality values.
var qualities = map[string]uint64{
	"raw":        fuji.QualityRaw,
	"fine":       fuji.QualityFine,
	"normal":     fuji.QualityNormal,
	"raw+fine":   fuji.QualityRawFine,
	"raw+normal": fuji.QualityRawNormal,
}

func qualityName(v uint64) string {
	for n, x := range qualities {
		if x == v {
			return n
		}
	}
	return fmt.Sprintf("0x%04X", v)
}

// rawModes maps the -raw flag to the camera's RawCompression values.
var rawModes = map[string]uint64{
	"uncompressed": fuji.RawUncompressed,
	"lossless":     fuji.RawLossless,
	"lossy":        fuji.RawLossy,
}

func rawModeName(v uint64) string {
	for n, x := range rawModes {
		if x == v {
			return n
		}
	}
	return fmt.Sprintf("0x%04X", v)
}

// programs maps the -mode flag to the camera's exposure program codes.
var programs = map[string]uint64{
	"manual":   fuji.ProgramManual,
	"aperture": fuji.ProgramAperturePriority,
	"shutter":  fuji.ProgramShutterPriority,
	"auto":     fuji.ProgramAuto,
}

// exposureOpts groups the exposure flags.
type exposureOpts struct {
	shutter   time.Duration
	aperture  float64
	iso       int
	show      bool
	release   bool
	mf        bool
	mode      string
	prop      string
	setv      int64
	setstr    string
	rawprop   string
	syncClock bool
}

func (e exposureOpts) any() bool {
	return e.shutter != 0 || e.aperture != 0 || e.iso != 0 || e.show || e.release || e.mf ||
		e.mode != "" || e.prop != "" || e.syncClock || e.rawprop != ""
}

// doSnapshot reads every property the camera advertises and either records it
// or compares it against a previous reading.
//
// This is how an unnamed property gets identified: change ONE thing on the
// camera and see which code moves. Guessing from a value's shape narrows it
// down; moving a dial settles it.
func doSnapshot(p func(string, ...any), s *fuji.Camera, di *ptp.DevInfo, snapPath, diffPath string) error {
	cur := map[ptp.Prop]string{}
	for _, c := range di.DeviceProps {
		d, err := s.GetPropDesc(c)
		if err != nil {
			continue
		}
		v := fmt.Sprintf("%d", int64(d.Current))
		if d.Type == ptp.TypeString {
			v = strconv.Quote(d.CurrentStr)
		}
		cur[c] = v
	}
	p("\nread %d of %d properties\n", len(cur), len(di.DeviceProps))

	if diffPath != "" {
		raw, err := os.ReadFile(diffPath)
		if err != nil {
			return fmt.Errorf("reading the snapshot: %w", err)
		}
		old := map[ptp.Prop]string{}
		for _, line := range strings.Split(string(raw), "\n") {
			var code uint32
			var val string
			if n, _ := fmt.Sscanf(line, "0x%04X %s", &code, &val); n >= 1 {
				if i := strings.Index(line, " "); i >= 0 {
					old[ptp.Prop(code)] = strings.TrimSpace(line[i+1:])
				}
			}
		}
		changed := 0
		for _, c := range di.DeviceProps {
			a, hadOld := old[c]
			b, hadNew := cur[c]
			if !hadOld || !hadNew || a == b {
				continue
			}
			// The clock always differs and drowns the signal.
			if c == ptp.PropDateTime {
				continue
			}
			changed++
			p("  CHANGED 0x%04X %-28s %s -> %s\n", uint16(c), fuji.PropName(c), a, b)
		}
		if changed == 0 {
			p("  nothing changed\n")
		}
		return nil
	}

	var b strings.Builder
	for _, c := range di.DeviceProps {
		if v, ok := cur[c]; ok {
			fmt.Fprintf(&b, "0x%04X %s\n", uint16(c), v)
		}
	}
	if err := os.WriteFile(snapPath, []byte(b.String()), 0o644); err != nil {
		return err
	}
	p("snapshot written to %s\n", snapPath)
	return nil
}

// doLiveView starts the preview stream, grabs frames, and stops it.
func doLiveView(p func(string, ...any), s *fuji.Camera, n int, getDir string) error {
	if ok, why := s.Ready(); !ok {
		return fmt.Errorf("not ready: %s", why)
	}
	if err := s.TakePriority(); err != nil {
		return err
	}
	if err := s.StartLiveView(); err != nil {
		return err
	}
	defer s.StopLiveView()

	p("\nlive view started; grabbing %d frame(s)\n", n)
	got, dup := 0, 0
	var last [16]byte
	var lastAt time.Time
	deadline := time.Now().Add(30 * time.Second)
	for got < n && time.Now().Before(deadline) {
		t0 := time.Now()
		data, err := s.LiveFrame()
		if err != nil {
			return err
		}
		if data == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// The store hands back the same object until the camera refreshes it,
		// so polling faster than the stream just re-fetches one frame. Count
		// only frames that actually changed, which also measures the rate.
		sum := md5.Sum(data)
		if sum == last {
			dup++
			// Poll gently. Live view refreshes on the order of tens of
			// milliseconds, and hammering GetObjectHandles faster than that
			// stopped the camera answering altogether.
			time.Sleep(100 * time.Millisecond)
			continue
		}
		gap := ""
		if !lastAt.IsZero() {
			gap = fmt.Sprintf(", %v since the last", time.Since(lastAt).Round(time.Millisecond))
		}
		last, lastAt = sum, time.Now()
		_ = gap
		got++
		kind := "?"
		if len(data) > 3 && data[0] == 0xFF && data[1] == 0xD8 {
			kind = "JPEG"
		}
		hs := ""
		if h, err := s.Histogram(); err == nil {
			hs = fmt.Sprintf("  hist total=%d clipped=%.2f%%", h.Total(0), h.Clipped(0)*100)
		}
		p("  frame %d: %d bytes %s, fetched in %v%s%s\n", got, len(data), kind,
			time.Since(t0).Round(time.Millisecond), gap, hs)
		if getDir != "" {
			os.MkdirAll(getDir, 0o755)
			os.WriteFile(filepath.Join(getDir, fmt.Sprintf("live%02d.jpg", got)), data, 0o644)
		}
	}
	if got == 0 {
		return fmt.Errorf("live view produced no frames in 30s")
	}
	p("  %d distinct frames, %d repeat fetches\n", got, dup)
	return nil
}

// doBurst shoots N frames in ONE session, which is the only way to learn the
// sustained rate. Shooting one frame per process hides the two things that
// actually bound a sequence: the camera's ~5-frame buffer, and whether it keeps
// up once it cannot rest between shots.
func doBurst(p func(string, ...any), s *fuji.Camera, n int, discard bool, getDir, bracket string) error {
	if ok, why := s.Ready(); !ok {
		return fmt.Errorf("not ready to shoot: %s", why)
	}
	if err := s.TakePriority(); err != nil {
		return err
	}
	// An eclipse sequence changes exposure every frame, so a cadence measured
	// at a constant shutter speed does not tell you what you need to know.
	var speeds []time.Duration
	if bracket != "" {
		for _, tok := range strings.Split(bracket, ",") {
			d, err := time.ParseDuration(strings.TrimSpace(tok))
			if err != nil {
				return fmt.Errorf("-bracket %q: %w", tok, err)
			}
			speeds = append(speeds, d)
		}
	}
	mode := "downloading"
	if discard {
		mode = "discarding"
	}
	p("\nburst: %d frames, %s each\n", n, mode)

	start := time.Now()
	var slowest time.Duration
	var pending []uint32      // frames shot but not yet downloaded
	held := map[uint32]bool{} // which handles are already in pending
	var setTotal time.Duration
	for i := 1; i <= n; i++ {
		t0 := time.Now()
		if len(speeds) > 0 {
			want := speeds[(i-1)%len(speeds)]
			ts := time.Now()
			if err := s.SetShutter(want); err != nil {
				return fmt.Errorf("frame %d: setting shutter to %v: %w", i, want, err)
			}
			setTotal += time.Since(ts)
		}
		if err := s.Capture(120 * time.Second); err != nil {
			return fmt.Errorf("frame %d of %d: %w", i, n, err)
		}
		// Wait for a NEW frame, not merely a non-empty store: when earlier
		// frames are still held the store is never empty, and testing for that
		// returns instantly with a stale handle.
		var h []uint32
		want := len(pending) + 1
		deadline := time.Now().Add(30 * time.Second)
		for time.Now().Before(deadline) {
			h, _ = s.GetObjectHandles(fuji.StillStore, 0, 0)
			if len(h) >= want {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if len(h) < want {
			return fmt.Errorf("frame %d of %d never appeared", i, n)
		}

		var note string
		if !discard {
			// Shoot now, download later. A one-frame lag is not enough: the
			// settle a frame needs is measured in seconds and a fast burst
			// produces one every 650ms, so keeping up would need ~15 frames of
			// lag against a buffer that holds 5. Downloads simply cannot track
			// a fast burst. Fill the buffer, then drain it once the frames have
			// aged — which is what an exposure bracket wants anyway.
			// GetObjectHandles returns the WHOLE store, not just what is new,
			// so appending it blindly counts frames several times and the
			// drain then tries to fetch a handle it already deleted.
			for _, x := range h {
				if !held[x] {
					held[x] = true
					pending = append(pending, x)
				}
			}
			note = fmt.Sprintf("held (%d in buffer)", len(pending))
			el := time.Since(t0)
			if el > slowest {
				slowest = el
			}
			free, _ := s.FreeBuffer()
			p("  %2d/%d  %-28s %6v  buffer %d free\n", i, n, note, el.Round(time.Millisecond), free)
			continue
		}
		if false {
			// Download the PREVIOUS frame, not this one. The newest frame is
			// not readable for several seconds — reading it early stalls
			// part-way and desynchronises the session — but older frames are
			// never affected. Lagging by one turns that wait into work, so the
			// settle happens while the next exposure is being taken instead of
			// while the host sits idle.
			if len(pending) > 0 {
				x := pending[0]
				pending = pending[1:]
				data, name, err := s.Download(x)
				if err != nil {
					return fmt.Errorf("frame %d: downloading the previous frame: %w", i, err)
				}
				if getDir != "" {
					os.WriteFile(filepath.Join(getDir, name), data, 0o644)
				}
				s.DeleteObject(x, 0)
				note = fmt.Sprintf("%s %d bytes", name, len(data))
			} else {
				note = "(priming)"
			}
			pending = append(pending, h...)
			el := time.Since(t0)
			if el > slowest {
				slowest = el
			}
			free, _ := s.FreeBuffer()
			p("  %2d/%d  %-28s %6v  buffer %d free\n", i, n, note, el.Round(time.Millisecond), free)
			continue
		}
		if discard {
			for _, x := range h {
				if err := s.DeleteObject(x, 0); err != nil {
					return fmt.Errorf("frame %d: discarding: %w", i, err)
				}
			}
			note = fmt.Sprintf("discarded %d", len(h))
		}
		el := time.Since(t0)
		if el > slowest {
			slowest = el
		}
		free, _ := s.FreeBuffer()
		p("  %2d/%d  %-28s %6v  buffer %d free\n", i, n, note, el.Round(time.Millisecond), free)
	}
	if len(pending) > 0 {
		shot := time.Since(start)
		// The frames have to age before they can be read. The oldest is
		// already `shot` old; wait out the remainder rather than a flat delay.
		const settle = 10 * time.Second
		if shot < settle {
			time.Sleep(settle - shot)
		}
		p("\n  shooting done in %v; draining %d frame(s)\n",
			shot.Round(time.Millisecond), len(pending))
	}
	for _, x := range pending {
		data, name, err := s.Download(x)
		if err != nil {
			return fmt.Errorf("draining the last frame: %w", err)
		}
		if getDir != "" {
			os.WriteFile(filepath.Join(getDir, name), data, 0o644)
		}
		s.DeleteObject(x, 0)
		p("  drain  %-28s %d bytes\n", name, len(data))
	}

	total := time.Since(start)
	p("\n  %d frames in %v — %.2fs/frame sustained, slowest %v\n",
		n, total.Round(time.Millisecond), total.Seconds()/float64(n),
		slowest.Round(time.Millisecond))
	if len(speeds) > 0 {
		p("  exposure changes cost %v total, %v per frame\n",
			setTotal.Round(time.Millisecond), (setTotal / time.Duration(n)).Round(time.Millisecond))
	}
	return nil
}

// listStores enumerates every storage volume and its contents. The question it
// answers is whether the SD card is reachable at all while tethered — the SDK
// says it is not, and this checks that against the camera rather than the
// manual.
func listStores(p func(string, ...any), s *fuji.Camera) error {
	ids, err := s.GetStorageIDs()
	if err != nil {
		return fmt.Errorf("listing storage volumes: %w", err)
	}
	p("\n%d storage volume(s):\n", len(ids))
	for _, id := range ids {
		info, err := s.GetStorageInfo(id)
		if err != nil {
			p("  0x%08X: %v\n", id, err)
			continue
		}
		p("  0x%08X %s\n", id, info)
		h, err := s.GetObjectHandles(id, 0, 0)
		if err != nil {
			p("      objects: %v\n", err)
			continue
		}
		p("      %d object(s)\n", len(h))
		// Summarise by format first: on a card with hundreds of files, the
		// counts are the useful thing and the names are noise.
		counts := map[string]int{}
		var newest []string
		for _, x := range h {
			if oi, err := s.GetObjectInfo(x); err == nil {
				counts[ptp.FormatName(oi.ObjectFormat)]++
				if oi.Filename != "" {
					newest = append(newest, oi.Filename)
				}
			}
		}
		for k, v := range counts {
			p("        %-28s %d\n", k, v)
		}
		if n := len(newest); n > 0 {
			from := n - 8
			if from < 0 {
				from = 0
			}
			p("        newest: %s\n", strings.Join(newest[from:], " "))
		}
		for i, x := range h {
			if i >= 0 {
				break
			}
			oi, err := s.GetObjectInfo(x)
			if err != nil {
				p("      0x%08X: %v\n", x, err)
				continue
			}
			p("      0x%08X %s\n", x, oi)
		}
	}
	return nil
}

// doProp reads or writes one property by name or code.
func (e exposureOpts) doProp(p func(string, ...any), s *fuji.Camera) error {
	code, ok := fuji.PropByName(e.prop)
	if !ok {
		// Parse as a code only if the whole string is one. Sscanf("%x") would
		// happily read "DriveMode" as 0x000D and then query a property nobody
		// asked for.
		hex := strings.TrimPrefix(strings.ToLower(e.prop), "0x")
		v, err := strconv.ParseUint(hex, 16, 16)
		if err != nil {
			// With 263 properties, a near miss is far more likely than a typo
			// of something that does not exist. ExposureProgram vs
			// ExposureProgramMode is the sort of thing worth catching.
			if near := fuji.PropsMatching(e.prop); len(near) > 0 {
				return fmt.Errorf("-prop %q: no such property. Did you mean: %s",
					e.prop, strings.Join(near, ", "))
			}
			return fmt.Errorf("-prop %q: not a known property name, and not a 0xNNNN code", e.prop)
		}
		code = ptp.Prop(v)
	}

	// A string property has to be read from its own slot; the numeric field is
	// zero for one, which is how DateTime came back as 0.
	if d, err := s.GetPropDesc(code); err == nil && d.Type == ptp.TypeString {
		if e.setstr != noWrite {
			if err := s.SetPropStringChecked(code, e.setstr); err != nil {
				return err
			}
		}
		v, err := s.GetPropString(code)
		if err != nil {
			return err
		}
		settable, why := s.Settable(code)
		note := ""
		if !settable {
			note = "  (" + why + ")"
		}
		p("  %s (0x%04X) = %q%s\n", fuji.PropName(code), uint16(code), v, note)
		return nil
	}

	if e.setv >= 0 {
		if err := s.SetProp(code, uint64(e.setv)); err != nil {
			return err
		}
	}
	cur, err := s.GetProp(code)
	if err != nil {
		return err
	}
	settable, why := s.Settable(code)
	note := ""
	if !settable {
		note = "  (" + why + ")"
	}
	p("  %s (0x%04X) = %s%s\n", fuji.PropName(code), uint16(code),
		fuji.ValueName(code, cur), note)
	if d, err := s.GetPropDesc(code); err == nil && len(d.Enum) > 1 {
		p("    offers: %s\n", fuji.DescribeValues(code, d.Enum))
	}
	return nil
}

// apply writes whichever exposure components were asked for, then reports the
// resulting triangle.
func (e exposureOpts) apply(p func(string, ...any), s *fuji.Camera) error {
	if e.rawprop != "" {
		var v uint64
		if _, err := fmt.Sscanf(strings.TrimPrefix(strings.ToLower(e.rawprop), "0x"), "%x", &v); err != nil {
			return fmt.Errorf("-rawprop %q: want a 0xNNNN code", e.rawprop)
		}
		code := ptp.Prop(v)
		data, _, err := s.Do(ptp.OpGetDevicePropValue, []uint32{uint32(code)}, nil, ptp.DefaultTimeout)
		if err != nil {
			return fmt.Errorf("fetching %s raw: %w", fuji.PropName(code), err)
		}
		p("  %s (0x%04X): %d bytes\n", fuji.PropName(code), uint16(code), len(data))
		if n := len(data); n > 0 {
			show := data
			if n > 256 {
				show = data[:256]
			}
			p("%s", hex.Dump(show))
			if n > 256 {
				p("  ... %d more bytes\n", n-256)
			}
		}
		return nil
	}
	if e.prop != "" {
		if err := e.doProp(p, s); err != nil {
			return err
		}
		return nil
	}
	if e.syncClock {
		before, _ := s.DateTime()
		p("  waiting for the next minute boundary — the body truncates seconds\n")
		if err := s.SyncClockAtMinute(); err != nil {
			return fmt.Errorf("syncing the clock: %w", err)
		}
		after, _ := s.DateTime()
		p("  clock: was %s, now %s (host %s)\n",
			before.Format("15:04:05"), after.Format("15:04:05"), time.Now().Format("15:04:05"))
	}
	if e.release {
		s.ReleaseAll()
		p("  shutter button released\n")
	}
	if e.mf {
		if err := s.SetFocusMode(fuji.FocusManual); err != nil {
			return fmt.Errorf("switching to manual focus: %w", err)
		}
		p("  focus mode: manual\n")
	}
	// The program mode comes first: it decides whether the writes below are
	// honoured at all.
	if e.mode != "" {
		m, ok := programs[strings.ToLower(e.mode)]
		if !ok {
			return fmt.Errorf("-mode %q: want manual, aperture, shutter or auto", e.mode)
		}
		if err := s.SetExposureProgram(m); err != nil {
			return fmt.Errorf("setting exposure program to %s: %w", e.mode, err)
		}
		p("  exposure program: %s\n", e.mode)
	}
	if e.shutter != 0 {
		d, err := s.GetPropDesc(ptp.PropExposureTime)
		if err != nil {
			return fmt.Errorf("reading the shutter speed range: %w", err)
		}
		if err := s.SetShutterFrom(d, e.shutter); err != nil {
			return fmt.Errorf("setting shutter to %v: %w", e.shutter, err)
		}
		if got := fuji.DecodeShutter(fuji.NearestShutter(d, e.shutter)); got != e.shutter {
			p("  shutter %v is not offered; snapped to %v\n", e.shutter, got)
		}
	}
	if e.aperture != 0 {
		if err := s.SetAperture(e.aperture); err != nil {
			return fmt.Errorf("setting aperture to f/%.1f: %w", e.aperture, err)
		}
	}
	if e.iso != 0 {
		if err := s.SetISO(uint32(e.iso)); err != nil {
			return fmt.Errorf("setting ISO to %d: %w", e.iso, err)
		}
	}

	// Always read back: a write that the camera accepts and ignores is the
	// normal failure here, not an error response.
	cur, err := s.ReadExposure()
	if err != nil {
		return err
	}
	p("  exposure: %s\n", cur)
	if !cur.ShutterSettable || !cur.ApertureSettable || !cur.ISOSettable {
		p("  a locked component is on a fixed dial position — the camera offers only\n" +
			"  one value and will ignore writes. Shutter dial to T, ISO dial to C,\n" +
			"  aperture ring to A to hand each one to the host.\n")
	}
	return nil
}

func run(serial string, list bool, dumpDir string, all bool, shoot bool, getDir, card, qual, rawc string, ex exposureOpts, drain, discard, lsAll bool, burst int, bracket string, live int, bulb time.Duration, snapPath, diffPath string, hold, settleFor time.Duration) error {
	if rawc != "" {
		if _, ok := rawModes[strings.ToLower(rawc)]; !ok {
			return fmt.Errorf("-raw %q: want uncompressed, lossless or lossy", rawc)
		}
	}
	if card != "" {
		if _, ok := cardModes[strings.ToLower(card)]; !ok {
			return fmt.Errorf("-card %q: want rawjpeg, raw, jpeg or off", card)
		}
	}
	if qual != "" {
		if _, ok := qualities[strings.ToLower(qual)]; !ok {
			return fmt.Errorf("-quality %q: want raw, fine, normal, raw+fine or raw+normal", qual)
		}
	}
	var out io.Writer = os.Stdout
	if dumpDir != "" {
		if err := os.MkdirAll(dumpDir, 0o755); err != nil {
			return fmt.Errorf("creating capture directory: %w", err)
		}
		f, err := os.Create(filepath.Join(dumpDir, "session.log"))
		if err != nil {
			return fmt.Errorf("creating session log: %w", err)
		}
		defer f.Close()
		out = &tee{w: []io.Writer{os.Stdout, f}}
	}
	p := func(format string, a ...any) { fmt.Fprintf(out, format, a...) }

	p("fujiprobe — ptp/fuji bring-up\n")
	p("time: %s\nhost: %s/%s, go %s\n\n", time.Now().Format(time.RFC3339),
		runtime.GOOS, runtime.GOARCH, runtime.Version())

	devs, err := usb.Enumerate()
	if err != nil {
		return err
	}
	p("%d Fujifilm device(s) attached:\n", len(devs))
	for _, d := range devs {
		p("  %s\n", d)
	}
	if len(devs) == 0 {
		return fmt.Errorf("no Fujifilm USB device found (is the body powered on, " +
			"and its USB mode set to PC Shoot / USB Tether rather than card reader?)")
	}
	if list {
		return nil
	}
	p("\n")

	// fuji.Open returns a camera with its PTP session already open and its
	// teardown hook wired, so there is no separate session to construct.
	s, err := fuji.Open(serial)
	if err != nil {
		return err
	}
	defer s.Close()
	if d, ok := s.Info(); ok {
		p("opened %s\n", d)
	}

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
			seq, uint16(ev.Op), ev.Params, ev.TxID, len(ev.DataIn), len(ev.DataOut),
			ev.Duration.Round(time.Millisecond), status)
		if dumpDir != "" && len(ev.DataIn) > 0 {
			os.WriteFile(filepath.Join(dumpDir,
				fmt.Sprintf("%03d-op%04X-in.bin", seq, uint16(ev.Op))), ev.DataIn, 0o644)
		}
	}

	p("\nGetDeviceInfo (0x1001)...\n")
	raw, _, err := s.Do(ptp.OpGetDeviceInfo, nil, nil, 5*time.Second)
	if err != nil {
		return fmt.Errorf("GetDeviceInfo: %w", err)
	}
	if dumpDir != "" {
		os.WriteFile(filepath.Join(dumpDir, "deviceinfo.bin"), raw, 0o644)
	}
	di, err := ptp.ParseDeviceInfo(raw)
	if err != nil {
		p("  %d bytes, but parse failed: %v\n%s\n", len(raw), err, hex.Dump(raw[:min(len(raw), 256)]))
		return err
	}
	p("  %s\n", di)
	p("  operations (%d): %s\n", len(di.Operations), fmtOps(di.Operations))
	p("  device properties (%d): %s\n", len(di.DeviceProps), fmtProps(di.DeviceProps))
	p("  events (%d): %v\n", len(di.Events), di.Events)
	p("  capture formats: %d, image formats: %d\n", len(di.CaptureFormats), len(di.ImageFormats))
	p("  capture: InitiateCapture=%v, open/bulb capture=%v\n", di.SupportsCapture(), di.SupportsBulb())

	// Which of the camera's properties does the decode already name, and which
	// are new? The unnamed ones are where the remaining work is.
	named, unnamed := 0, []ptp.Prop{}
	for _, c := range di.DeviceProps {
		if strings.HasPrefix(fuji.PropName(c), "Prop(0x") {
			unnamed = append(unnamed, c)
		} else {
			named++
		}
	}
	p("  -> %d of %d properties are named by the decode; %d unknown\n", named, len(di.DeviceProps), len(unnamed))
	if len(unnamed) > 0 {
		p("     unknown: %s\n", fmtProps(unnamed))
	}

	if ok, why := s.Ready(); !ok {
		p("  NOT READY: %s\n", why)
	}

	if lsAll {
		return listStores(p, s)
	}
	if snapPath != "" || diffPath != "" {
		return doSnapshot(p, s, di, snapPath, diffPath)
	}
	if hold > 0 {
		if ok, why := s.Ready(); !ok {
			return fmt.Errorf("not ready: %s", why)
		}
		if err := s.TakePriority(); err != nil {
			return err
		}
		dm, _ := s.DriveMode()
		p("\nholding the shutter for %v (drive mode %s)\n", hold, fuji.ValueName(fuji.PropDriveMode, dm))
		st := time.Now()
		if err := s.CaptureHold(hold); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		h, _ := s.GetObjectHandles(fuji.StillStore, 0, 0)
		p("  %d frame(s) in %v\n", len(h), time.Since(st).Round(time.Millisecond))
		for _, x := range h {
			s.DeleteObject(x, 0)
		}
		p("  discarded\n")
		return nil
	}
	if bulb > 0 {
		if ok, why := s.Ready(); !ok {
			return fmt.Errorf("not ready: %s", why)
		}
		if err := s.TakePriority(); err != nil {
			return err
		}
		p("\nbulb: holding the shutter open for %v\n", bulb)
		st := time.Now()
		if err := s.BulbCapture(bulb); err != nil {
			return err
		}
		p("  shutter cycle took %v\n", time.Since(st).Round(time.Millisecond))
		deadline := time.Now().Add(60 * time.Second)
		for time.Now().Before(deadline) {
			h, _ := s.GetObjectHandles(fuji.StillStore, 0, 0)
			if len(h) > 0 {
				oi, err := s.GetObjectInfo(h[0])
				if err == nil {
					p("  frame appeared %v after the shutter opened: %s\n",
						time.Since(st).Round(time.Millisecond), oi)
				}
				return nil
			}
			time.Sleep(100 * time.Millisecond)
		}
		return fmt.Errorf("bulb exposure produced no frame")
	}
	if live > 0 {
		return doLiveView(p, s, live, getDir)
	}
	if burst > 0 {
		return doBurst(p, s, burst, discard, getDir, bracket)
	}

	// Where do frames go? Always worth stating: the USB transfer happens either
	// way, this is only the SD card copy.
	if v, err := s.MediaRecord(); err == nil {
		p("  tethered frames: transferred over USB; SD card copy = %s\n", cardModeName(v))
	}
	if v, err := s.Quality(); err == nil {
		p("  quality: %s\n", qualityName(v))
	}
	if v, err := s.RawCompression(); err == nil {
		p("  RAW recording: %s\n", rawModeName(v))
	}
	if v, err := s.FreeBuffer(); err == nil {
		p("  volatile buffer: room for %d more frame(s)\n", v)
	}
	if v, err := s.GetPropValue(fuji.PropPriorityModeCode, ptp.TypeUint16); err == nil {
		who := fmt.Sprintf("0x%04X", v.Num)
		switch v.Num {
		case fuji.PriorityModeCamera:
			who = "camera (body controls live)"
		case fuji.PriorityModeHost:
			who = "PC (body controls LOCKED OUT)"
		}
		p("  priority: %s\n", who)
	}

	// Exposure before capture, so -shutter 1ms -capture is one command.
	if ex.any() {
		if err := ex.apply(p, s); err != nil {
			return err
		}
		if !shoot && getDir == "" {
			return nil
		}
	}

	if shoot || getDir != "" || card != "" || qual != "" || rawc != "" {
		return doCapture(out, s, shoot, getDir, card, qual, rawc, drain, discard, settleFor)
	}

	p("\nreading property descriptors (0x1014)...\n")
	sort.Slice(di.DeviceProps, func(i, j int) bool { return di.DeviceProps[i] < di.DeviceProps[j] })
	fails, declined := 0, 0
	for _, c := range di.DeviceProps {
		d, err := s.GetPropDesc(c)
		if err != nil {
			p("  0x%04X %-30s %v\n", uint16(c), fuji.PropName(c), err)
			if errors.Is(err, ptp.ErrNotResponding) {
				return fmt.Errorf("camera stopped responding — wake it and retry: %w", err)
			}
			// A camera answering GeneralError is healthy: it is declining to
			// describe THAT property, which is a fact about the property, not
			// about the connection. Only a transport failure means the sweep
			// cannot continue — and a run of declines is normal, since the
			// unsupported ones cluster (0xD229..0xD22D are consecutive).
			var pe *ptp.Error
			if errors.As(err, &pe) {
				declined++
				continue
			}
			if fails++; fails >= 5 {
				return fmt.Errorf("giving up after %d consecutive transport failures", fails)
			}
			continue
		}
		fails = 0
		w := "read-only"
		if d.Writable() {
			w = "writable"
		}
		extra := ""
		switch {
		case d.Form == 1:
			extra = fmt.Sprintf("  range %d..%d step %d", int64(d.Min), int64(d.Max), int64(d.Step))
		case len(d.Enum) > 0:
			extra = fmt.Sprintf("  %d values %v", len(d.Enum), trim(d.Enum))
		case len(d.EnumStr) > 0:
			extra = fmt.Sprintf("  %d values %v", len(d.EnumStr), trimStr(d.EnumStr))
		}
		val := fmt.Sprintf("%d", int64(d.Current))
		if d.Type == 0xFFFF {
			val = fmt.Sprintf("%q", d.CurrentStr)
		}
		if !all && !d.Writable() && extra == "" {
			continue
		}
		p("  0x%04X %-30s type=0x%04X cur=%s %s%s\n", uint16(c), fuji.PropName(c), uint16(d.Type), val, w, extra)
	}

	p("\n%d properties described, %d declined by the camera\n",
		len(di.DeviceProps)-declined, declined)
	if dumpDir != "" {
		p("capture bundle written to %s\n", dumpDir)
	}
	return nil
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

func trimStr(v []string) []string {
	const max = 8
	if len(v) > max {
		return v[:max]
	}
	return v
}

func trim(v []uint64) []int64 {
	const max = 12
	out := make([]int64, 0, max)
	for i, x := range v {
		if i == max {
			break
		}
		out = append(out, int64(x))
	}
	return out
}

// doCapture takes a frame and optionally downloads what appeared.
func doCapture(out io.Writer, s *fuji.Camera, shoot bool, getDir, card, qual, rawc string, drain, discard bool, settleFor time.Duration) error {
	p := func(f string, a ...any) { fmt.Fprintf(out, f, a...) }

	// Only take priority if we are going to shoot. Downloading and deleting
	// need none, and asking for it here deadlocks the one case that matters:
	// the camera refuses priority WHILE frames are pending, so a drain that
	// demanded priority first could never clear them.
	if shoot {
		if err := s.TakePriority(); err != nil {
			return fmt.Errorf("handing priority to the host: %w", err)
		}
	}

	// Every tethered frame reaches the host from the camera's volatile store;
	// this only decides whether a copy is left on the SD card as well.
	if card != "" {
		if err := s.SetMediaRecord(cardModes[strings.ToLower(card)]); err != nil {
			return fmt.Errorf("setting the SD card copy to %s: %w", card, err)
		}
	}
	if qual != "" {
		if err := s.SetQuality(qualities[strings.ToLower(qual)]); err != nil {
			return fmt.Errorf("setting quality to %s: %w", qual, err)
		}
	}
	if rawc != "" {
		if err := s.SetRawCompression(rawModes[strings.ToLower(rawc)]); err != nil {
			return fmt.Errorf("setting RAW recording to %s: %w", rawc, err)
		}
	}
	if card != "" || qual != "" || rawc != "" {
		v, _ := s.MediaRecord()
		q, _ := s.Quality()
		r, _ := s.RawCompression()
		p("now: quality %s, RAW %s, SD card copy %s\n",
			qualityName(q), rawModeName(r), cardModeName(v))
	}
	if !shoot && getDir == "" {
		return nil
	}
	before, _ := s.GetObjectHandles(fuji.StillStore, 0, 0)

	if shoot {
		p("\ncapture: S1 half press, settle, S2 full press...\n")
		st := time.Now()
		if err := s.Capture(120 * time.Second); err != nil {
			return fmt.Errorf("capture: %w", err)
		}
		p("  shutter released after %v; waiting for the frame\n", time.Since(st).Round(time.Millisecond))
		_ = st

		// The frame appears in the Still store a moment after the shutter.
		var newest uint32
		deadline := time.Now().Add(2 * time.Minute)
		for time.Now().Before(deadline) {
			h, _ := s.GetObjectHandles(fuji.StillStore, 0, 0)
			if len(h) > len(before) {
				newest = h[len(h)-1]
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		// Time from the full press to the frame existing. For a long exposure
		// this is dominated by the exposure itself, so it is the direct check
		// that a shutter speed actually reached the sensor.
		p("  frame appeared %v after the shutter\n", time.Since(st).Round(10*time.Millisecond))
		// The frame appears in the store several seconds before it can be read.
		// Downloading too early stalls part-way and desynchronises the session
		// (seen at 2, 19 and 31 MB of 37 MB on an X-T5); the same frame then
		// downloads perfectly once the camera has settled. Older frames are
		// unaffected — only the newest one is still being written.
		// This used to hold a fixed ten seconds, on the theory that a frame is
		// not READABLE for several seconds after it appears and that
		// downloading early stalls part-way.
		//
		// MEASURED FALSE on an X-T5 (2026-08-07). Downloading the instant the
		// handle appears produced a complete 82,341,888-byte uncompressed frame
		// that decodes PIXEL-FOR-PIXEL IDENTICAL to libraw across all
		// 40,902,912 samples, plus a lossless frame libraw reads end to end. It
		// took a capture from 11.6s to 2.7s.
		//
		// What matters is waiting for the handle to APPEAR, which the loop above
		// already does — the frame showed up 650ms after the shutter. The fixed
		// wait was belt-and-braces on top of that.
		//
		// The knob survives at 0 because another body or firmware might really
		// stall, and this is then the first thing to try. If it is ever needed,
		// do NOT wait it out with a bare sleep: ten seconds of silence and the
		// body drops the link. Poll, as below, to keep the bus alive.
		if minSettle := minSettleWanted(settleFor, discard); minSettle > 0 {
			p("  letting the camera finish writing the frame...\n")
		}
		minSettle := minSettleWanted(settleFor, discard)
		start := time.Now()
		for time.Since(start) < minSettle {
			time.Sleep(500 * time.Millisecond)
			s.GetObjectInfo(newest) // keep-alive; failures are expected early
		}
	}

	// Card-only shooting: drop the frame from the volatile store without ever
	// transferring it. The camera keeps its card copy, and the host pays a few
	// milliseconds for a delete instead of most of a second for 35 MB.
	if discard {
		h, err := s.GetObjectHandles(fuji.StillStore, 0, 0)
		if err != nil {
			return fmt.Errorf("listing the Still store: %w", err)
		}
		st := time.Now()
		for _, x := range h {
			// Name the frame before dropping it. The card copy cannot be
			// checked from here, so the filename is the only handle the user
			// has for finding it in the camera's own playback.
			name := fmt.Sprintf("0x%08X", x)
			if oi, err := s.GetObjectInfo(x); err == nil && oi.Filename != "" {
				name = oi.Filename
			}
			if err := s.DeleteObject(x, 0); err != nil {
				return fmt.Errorf("discarding frame %s: %w", name, err)
			}
			p("  discarded %s from the buffer without transferring it\n", name)
		}
		p("  %d frame(s) in %v\n", len(h), time.Since(st).Round(time.Millisecond))
		return nil
	}

	after, err := s.GetObjectHandles(fuji.StillStore, 0, 0)
	if err != nil {
		return fmt.Errorf("listing the Still store: %w", err)
	}
	if shoot && len(after) <= len(before) {
		return fmt.Errorf("no new frame appeared: the camera accepted every command " +
			"but produced nothing")
	}
	p("\nStill store: %d frame(s)\n", len(after))

	seen := map[uint32]bool{}
	for _, h := range before {
		seen[h] = true
	}
	for _, h := range after {
		oi, err := s.GetObjectInfo(h)
		if err != nil {
			p("  0x%08X: %v\n", h, err)
			continue
		}
		tag := " "
		if !seen[h] {
			tag = "*"
		}
		p("  %s %s\n", tag, oi)
		if getDir == "" || (shoot && seen[h]) {
			continue
		}
		if err := os.MkdirAll(getDir, 0o755); err != nil {
			return err
		}
		st := time.Now()
		data, err := s.GetObject(h)
		if err != nil {
			p("    download failed: %v\n", err)
			continue
		}
		dst := filepath.Join(getDir, oi.Filename)
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
		el := time.Since(st)
		p("    -> %s  %d bytes in %v (%.1f MB/s)\n", dst, len(data),
			el.Round(time.Millisecond), float64(len(data))/el.Seconds()/1e6)
		if uint32(len(data)) != oi.CompressedSize {
			p("    WARNING: got %d bytes, ObjectInfo said %d\n", len(data), oi.CompressedSize)
			continue // size mismatch: do not delete the camera's copy
		}
		// Only now is deleting safe, and only now does the buffer drain. Frames
		// left in the volatile store fill its ~5 slots AND stop the camera
		// handing control back to the user on close.
		if drain {
			if err := s.DeleteObject(h, 0); err != nil {
				p("    could not delete from the camera: %v\n", err)
			} else {
				p("    deleted from the camera's buffer\n")
			}
		}
	}
	return nil
}

// minSettleWanted caps the post-capture wait. Discarding never reads the frame,
// so it needs less of one than a download does.
func minSettleWanted(want time.Duration, discard bool) time.Duration {
	if want < 0 {
		return 0
	}
	if discard && want > 1500*time.Millisecond {
		return 1500 * time.Millisecond
	}
	return want
}
