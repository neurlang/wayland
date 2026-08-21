//go:build linux

package window

import (
	"sync"
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
	steps []int32
}

func (h *discreteAxisHandler) AxisDiscrete(_ *Widget, _ *Input, _ uint32, discrete int32) {
	h.steps = append(h.steps, discrete)
}

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
	}
	input.PointerAxisValue120(nil, wl.PointerAxisVerticalScroll, -240)

	want := []int32{1, -2}
	if len(handler.steps) != len(want) {
		t.Fatalf("discrete events = %v, want %v", handler.steps, want)
	}
	for i := range want {
		if handler.steps[i] != want[i] {
			t.Fatalf("discrete events = %v, want %v", handler.steps, want)
		}
	}
}
