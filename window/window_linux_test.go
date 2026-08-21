//go:build linux

package window

import (
	"github.com/neurlang/wayland/wl"

	"sync"
	"testing"
	"time"
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
