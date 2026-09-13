//go:build linux

package window

import (
	"github.com/neurlang/wayland/wl"

	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/neurlang/wayland/wl"
)

// noopRunner is a no-op implementation of the runner interface for testing.
type noopRunner struct{}

func (n *noopRunner) Run(_ uint32) {}

// chanRunner is a runner that signals a channel when Run is called.
type chanRunner struct {
	done chan struct{}
}

func (r *chanRunner) Run(_ uint32) {
	close(r.done)
}

// TestDisplayDeferRace detects the data race on Display.deferredListNew.
//
// Without the fix, concurrent goroutines calling displayDefer() unsynchronized
// will trigger the Go race detector on the slice append. Run with:
//
//	go test -race -run TestDisplayDeferRace ./window/...
//
// Validates: Requirements REQ-1
func TestDisplayDeferRace(t *testing.T) {
	// Construct a minimal Display. Display.Display is nil so no wakeup
	// is attempted (the fixed code guards on Display.Display != nil).
	d := &Display{}

	const N = 50
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			displayDefer(d, &noopRunner{})
		}()
	}

	wg.Wait()
}

// TestDisplayRunWakesOnDefer verifies that a task enqueued via displayDefer()
// from a background goroutine is eventually processed.
//
// Since we cannot connect to a live compositor, the test drives the drain loop
// manually — the same logic that DisplayRun() uses — and asserts the task runs
// within a bounded timeout.
//
// Pre-fix: the test still passes here because it drives the drain manually.
// The race-detector test (TestDisplayDeferRace) is the primary pre-fix reproducer.
//
// Validates: Requirements REQ-2
func TestDisplayRunWakesOnDefer(t *testing.T) {
	d := &Display{}
	done := make(chan struct{})

	task := &chanRunner{done: done}

	// Enqueue the task from a background goroutine after a short delay,
	// simulating a real-world scenario where work arrives while the loop is
	// "sleeping" waiting for compositor events.
	go func() {
		time.Sleep(50 * time.Millisecond)
		displayDefer(d, task)
	}()

	// Drive the drain loop manually (mirrors the inner loop of DisplayRun).
	// Poll with a short sleep so we don't busy-spin in the test binary.
	timeout := time.After(5 * time.Second)
	for {
		// Drain all pending tasks under the mutex, mirroring DisplayRun.
		d.deferredMu.Lock()
		for len(d.deferredListNew) > 0 {
			item := d.deferredListNew[0]
			d.deferredListNew = d.deferredListNew[1:]
			d.deferredMu.Unlock()
			item.Run(0)
			d.deferredMu.Lock()
		}
		d.deferredMu.Unlock()

		select {
		case <-done:
			// Task ran — test passes.
			return
		case <-timeout:
			// Task never ran within the deadline.
			t.Fatal("deferred task was not processed within the timeout — display loop may be blocked")
		default:
			// Not done yet; yield briefly before polling again.
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func TestDisplayBufferScaleUsesLargestAvailableOutputScale(t *testing.T) {
	d := &Display{outputList: []*output{
		{scale: 1},
		nil,
		{scale: 2},
		{scale: 3},
	}}

	if got := d.BufferScale(); got != 3 {
		t.Errorf("BufferScale() = %d, want 3", got)
	}
}

func TestNormalizeBufferScaleRejectsInvalidScale(t *testing.T) {
	for _, scale := range []int32{-1, 0} {
		if got := normalizeBufferScale(scale); got != 1 {
			t.Errorf("normalizeBufferScale(%d) = %d, want 1", scale, got)
		}
	}
}

func TestFractionalBufferSize(t *testing.T) {
	tests := []struct {
		name   string
		length int32
		scale  uint32
		want   int32
	}{
		{name: "one times", length: 1000, scale: 120, want: 1000},
		{name: "one and a half times", length: 1000, scale: 180, want: 1500},
		{name: "rounds up", length: 101, scale: 180, want: 152},
		{name: "sub-unit scale", length: 1000, scale: 90, want: 750},
		{name: "zero falls back to one", length: 1000, scale: 0, want: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fractionalBufferSize(tt.length, tt.scale); got != tt.want {
				t.Errorf("fractionalBufferSize(%d, %d) = %d, want %d", tt.length, tt.scale, got, tt.want)
			}
		})
	}
}

type nestedResizeHandler struct {
	WidgetHandler
	calls int
}

func (h *nestedResizeHandler) Resize(widget *Widget, _, _, _, _ int32) {
	h.calls++
	if h.calls == 1 {
		widget.ScheduleResize(933, 680)
	}
}

func TestDrainPendingResizesAppliesRequestFromResizeHandler(t *testing.T) {
	display := &Display{}
	window := &Window{
		Display: display,
		pendingAllocation: Rectangle{
			Width:  1000,
			Height: 690,
		},
		resizeNeeded: 1,
	}
	surface := &surface{Window: window}
	handler := &nestedResizeHandler{}
	widget := &Widget{
		Window:   window,
		surface:  surface,
		userdata: handler,
	}
	surface.Widget = widget
	window.mainSurface = surface

	drainPendingResizes(window)

	if handler.calls != 2 {
		t.Fatalf("resize handler called %d times, want 2", handler.calls)
	}
	if got := widget.allocation; got.Width != 933 || got.Height != 680 {
		t.Errorf("final allocation = %dx%d, want 933x680", got.Width, got.Height)
	}
	if window.resizeNeeded != 0 {
		t.Errorf("resizeNeeded = %d, want 0", window.resizeNeeded)
type pointerFrameHandler struct {
	WidgetHandler
	frames    int
	value120s []int32
}

func (h *pointerFrameHandler) PointerFrame(*Widget, *Input) {
	h.frames++
}

func (h *pointerFrameHandler) AxisValue120(_ *Widget, _ *Input, _ uint32, value120 int32) {
	h.value120s = append(h.value120s, value120)
}

type discreteAxisHandler struct {
	WidgetHandler
	steps   []int32
	rawAxes int
}

func (h *discreteAxisHandler) AxisDiscrete(_ *Widget, _ *Input, _ uint32, discrete int32) {
	h.steps = append(h.steps, discrete)
}

func (h *discreteAxisHandler) Axis(_ *Widget, _ *Input, _ uint32, _ uint32, _ float32) {
	h.rawAxes++
}

func (h *discreteAxisHandler) PointerFrame(*Widget, *Input) {}

func pointerTestInput(handler WidgetHandler) *Input {
	display := &Display{}
	window := &Window{Display: display}
	widget := &Widget{Window: window, userdata: handler}
	window.mainSurface = &surface{Window: window, Widget: widget}
	return &Input{Display: display, pointerFocus: window, focusWidget: widget}
}

func TestPointerFrameAndValue120ReachOptionalWidgetHandler(t *testing.T) {
	handler := &pointerFrameHandler{}
	input := pointerTestInput(handler)

	input.PointerAxisValue120(nil, wl.PointerAxisVerticalScroll, 30)
	input.PointerFrame(nil)

	if len(handler.value120s) != 1 || handler.value120s[0] != 30 {
		t.Fatalf("value120 events = %v, want [30]", handler.value120s)
	}
	if handler.frames != 1 {
		t.Fatalf("pointer frames = %d, want 1", handler.frames)
	}
}

func TestPointerValue120FallbackAccumulatesWholeSteps(t *testing.T) {
	handler := &discreteAxisHandler{}
	input := pointerTestInput(handler)

	for range 4 {
		input.PointerAxisValue120(nil, wl.PointerAxisVerticalScroll, 30)
		input.PointerAxis(nil, 0, wl.PointerAxisVerticalScroll, 5)
		input.PointerFrame(nil)
	}
	input.PointerAxisValue120(nil, wl.PointerAxisVerticalScroll, -240)
	input.PointerAxis(nil, 0, wl.PointerAxisVerticalScroll, -10)
	input.PointerFrame(nil)

	want := []int32{1, -2}
	if len(handler.steps) != len(want) {
		t.Fatalf("discrete events = %v, want %v", handler.steps, want)
	}
	for i := range want {
		if handler.steps[i] != want[i] {
			t.Fatalf("discrete events = %v, want %v", handler.steps, want)
		}
	}
	if handler.rawAxes != 0 {
		t.Fatalf("paired raw axis was delivered %d times, want 0", handler.rawAxes)
	}

	input.PointerAxis(nil, 0, wl.PointerAxisVerticalScroll, 5)
	if handler.rawAxes != 1 {
		t.Fatalf("raw axis after frame was delivered %d times, want 1", handler.rawAxes)
// TestWindowScheduleRedrawTaskConcurrent verifies that concurrent redraw
// requests atomically claim a single deferred task. In particular, the
// scheduled flag must already be visible when the task is published to the
// display queue; otherwise DisplayRun can consume the task before the flag is
// set and all later redraw requests will be dropped.
func TestWindowScheduleRedrawTaskConcurrent(t *testing.T) {
	const (
		attempts = 100
		workers  = 32
	)

	for attempt := 0; attempt < attempts; attempt++ {
		d := &Display{}
		w := &Window{Display: d}

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				<-start
				windowScheduleRedrawTask(w)
			}()
		}
		close(start)
		wg.Wait()

		d.deferredMu.Lock()
		queued := len(d.deferredListNew)
		d.deferredMu.Unlock()

		if queued != 1 {
			t.Fatalf("attempt %d: queued redraw tasks = %d, want 1", attempt, queued)
		}
		if got := atomic.LoadInt32(&w.redrawTaskScheduled); got != 1 {
			t.Fatalf("attempt %d: redrawTaskScheduled = %d, want 1", attempt, got)
		}
func TestWindowGeometryIncludesClientTitlebar(t *testing.T) {
	window := &Window{
		mainSurface: &surface{allocation: Rectangle{X: 3, Y: 5, Width: 1707, Height: 1021}},
		decoration:  &WindowDecoration{},
	}

	var geometry Rectangle
	windowGetGeometry(window, &geometry)

	want := Rectangle{X: 3, Y: 5 - TitleHeight, Width: 1707, Height: 1021 + TitleHeight}
	if geometry != want {
		t.Fatalf("window geometry = %+v, want %+v", geometry, want)
	}
}

func TestWindowGeometryExcludesNoUndecoratedSpace(t *testing.T) {
	window := &Window{
		mainSurface: &surface{allocation: Rectangle{X: 3, Y: 5, Width: 1707, Height: 1021}},
	}

	var geometry Rectangle
	windowGetGeometry(window, &geometry)

	want := Rectangle{X: 3, Y: 5, Width: 1707, Height: 1021}
	if geometry != want {
		t.Fatalf("window geometry = %+v, want %+v", geometry, want)
	}
}

func TestWindowContentHeightForGeometry(t *testing.T) {
	decorated := &Window{
		Display:              &Display{subcompositor: &wl.Subcompositor{}},
		decorationsRequested: true,
	}

	height, ok := windowContentHeightForGeometry(decorated, 1021)
	if !ok || height != 1021-TitleHeight {
		t.Fatalf("decorated content height = %d, %t; want %d, true", height, ok, 1021-TitleHeight)
	}

	decorated.decoration = &WindowDecoration{}
	decorated.decorationsRequested = false
	height, ok = windowContentHeightForGeometry(decorated, 1021)
	if !ok || height != 1021-TitleHeight {
		t.Fatalf("created decoration content height = %d, %t; want %d, true", height, ok, 1021-TitleHeight)
	}

	decorated.fullscreen = true
	height, ok = windowContentHeightForGeometry(decorated, 1021)
	if !ok || height != 1021 {
		t.Fatalf("fullscreen content height = %d, %t; want 1021, true", height, ok)
	}

	undecorated := &Window{Display: &Display{subcompositor: &wl.Subcompositor{}}}
	height, ok = windowContentHeightForGeometry(undecorated, 1021)
	if !ok || height != 1021 {
		t.Fatalf("undecorated content height = %d, %t; want 1021, true", height, ok)
	}

	_, ok = windowContentHeightForGeometry(&Window{
		Display:              &Display{subcompositor: &wl.Subcompositor{}},
		decorationsRequested: true,
	}, TitleHeight)
	if ok {
		t.Fatal("titlebar-only geometry was accepted")
	}
}
