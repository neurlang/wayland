# Display Event Loop: deferredListNew Race Fix

## Problem

Two related issues exist in `window/window_linux.go`:

### Issue 1 — Data race on `Display.deferredListNew` (FIXED)

`displayDefer()` appends to `Display.deferredListNew` without synchronization,
while `DisplayRun()` reads and reslices it on the display thread. Any goroutine
calling `ScheduleRedraw()` → `windowScheduleRedrawTask()` → `displayDefer()`
concurrently triggers a data race detectable by `go test -race`.

### Issue 2 — Blocking event loop with no wakeup (KNOWN LIMITATION, not yet fixed)

After draining `deferredListNew`, `DisplayRun()` calls `wlclient.DisplayRun()`
→ `wl.Context.Run()` → `ctx.readEvent()`, which blocks on `ReadMsgUnix` with
no deadline. Tasks enqueued by background goroutines after this point are not
processed until the compositor sends an event (e.g. mouse move, key press).

This means asynchronous UI updates (progress bars, background worker callbacks,
timer-driven redraws) remain visually pending until the next external input event.

A `wl_display.sync` wakeup mechanism (Option B from the original bug report) was
attempted but caused application freezes because `Sync()` was being called from
within the display thread's own drain loop — the sync request reached the
compositor but the client-side event loop could not process the response while
it was itself blocked waiting for that response. This issue is left for a future
fix.

## Fix Applied (Issue 1)

Added `deferredMu sync.Mutex` to the `Display` struct to guard all accesses to
`deferredListNew`:

**`displayDefer`** — wraps the slice append under the mutex:

```go
func displayDefer(Display *Display, fun runner) {
    Display.deferredMu.Lock()
    Display.deferredListNew = append(Display.deferredListNew, fun)
    Display.deferredMu.Unlock()
}
```

**`DisplayRun`** — drains the queue under the mutex, dropping it before each
`task.Run(0)` call to avoid holding the lock during arbitrary user callbacks
(which may themselves call `displayDefer`):

```go
func DisplayRun(Display *Display) {
    Display.running = true
    for {
        Display.deferredMu.Lock()
        for len(Display.deferredListNew) > 0 {
            task := Display.deferredListNew[0]
            Display.deferredListNew = Display.deferredListNew[1:]
            Display.deferredMu.Unlock()
            task.Run(0)
            Display.deferredMu.Lock()
        }
        Display.deferredMu.Unlock()

        if !Display.running {
            break
        }

        if err := wlclient.DisplayRun(Display.Display); err != nil {
            fmt.Println(err)
            return
        }
    }
}
```

## Test Coverage

`window/window_linux_test.go` (build tag `//go:build linux`) contains:

- **`TestDisplayDeferRace`** — spawns 50 goroutines each calling `displayDefer()`
  concurrently. Passes cleanly under `go test -race` with the fix; triggers the
  race detector without it.

- **`TestDisplayRunWakesOnDefer`** — drives the drain loop manually and verifies
  that a task enqueued from a background goroutine is processed within a bounded
  timeout. This test passes both before and after the fix since it does not
  exercise the blocking-readEvent path.

## Files Changed

| File | Change |
|------|--------|
| `window/window_linux.go` | Added `deferredMu sync.Mutex` to `Display`; mutex-guarded `displayDefer` and `DisplayRun` drain loop |
| `window/window_linux_test.go` | New file: `TestDisplayDeferRace`, `TestDisplayRunWakesOnDefer` |
