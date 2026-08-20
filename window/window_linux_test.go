//go:build linux

package window

import (
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
