# Darwin Keyboard Support Implementation

## Overview

Keyboard support has been added to the Darwin (macOS) window implementation, allowing applications to receive and handle keyboard events.

## Implementation Details

### Files Modified

1. **window/window_cgo_darwin.go**
   - Added `goKeyPress` callback declaration
   - Added `keyMonitor` field to `DarwinWindow` struct
   - Implemented NSEvent keyboard event monitoring for `NSEventMaskKeyDown`, `NSEventMaskKeyUp`, and `NSEventMaskFlagsChanged`
   - Captures key code, state, modifier flags, and unicode character
   - Properly cleans up keyboard monitor on window destruction

2. **window/widget_darwin.go**
   - Added `modifiers` field to `Input` struct to track modifier key state
   - Implemented `updateModifiers()` method to convert NSEvent modifier flags to window.ModType
   - Updated `GetModifiers()` to return actual modifier state instead of 0
   - Maps NSEvent modifier flags to existing constants: `ModShiftMask`, `ModControlMask`, `ModAltMask`

3. **window/window_cgo_darwin.go (Go side)**
   - Implemented `goKeyPress` export function
   - Calls keyboard handler with proper parameters
   - Updates modifier state on each key event
   - Converts NSEvent state to `wl.KeyboardKeyState` (Pressed/Released)

### Key Mapping

The implementation uses macOS virtual key codes directly, which are already defined in `xkbcommon/keysyms_darwin.go`. These include:

- Letter keys (A-Z)
- Number keys (0-9)
- Function keys (F1-F12)
- Navigation keys (arrows, Home, End, Page Up/Down)
- Special keys (Return, Backspace, Delete, Tab, Escape, Space)
- Modifier keys (Shift, Control, Alt/Option, Command)

### Modifier Support

The following modifiers are tracked and reported:

- **Shift** → `ModShiftMask`
- **Control** → `ModControlMask`
- **Alt/Option** → `ModAltMask`
- **Command** (tracked but no constant in window/constants.go)
- **CapsLock** (tracked but no constant in window/constants.go)

### KeyboardHandler Interface

Applications implement the `KeyboardHandler` interface defined in `window/widget_handler_darwin.go`:

```go
type KeyboardHandler interface {
    Key(Window *Window, Input *Input, time uint32, vKey uint32, code uint32, state wl.KeyboardKeyState, data WidgetHandler)
    Focus(Window *Window, Input *Input)
}
```

Parameters:
- `vKey`: macOS virtual key code (matches constants in xkbcommon/keysyms_darwin.go)
- `code`: Unicode character value (if available)
- `state`: Key state (Pressed or Released)
- `Input.GetModifiers()`: Returns current modifier key state

## Usage Example

```go
type MyApp struct {
    window *window.Window
    widget *window.Widget
}

// Implement KeyboardHandler interface
func (app *MyApp) Key(win *window.Window, input *window.Input, time uint32, vKey uint32, code uint32, state wl.KeyboardKeyState, data window.WidgetHandler) {
    if state != wl.KeyboardKeyStatePressed {
        return
    }
    
    mods := input.GetModifiers()
    
    // Check for Cmd+Q (using Control mask as Command equivalent)
    if vKey == xkbcommon.KeyQ && mods&window.ModControlMask != 0 {
        app.window.Destroy()
    }
    
    // Handle other keys
    switch vKey {
    case xkbcommon.KeyLeft:
        // Handle left arrow
    case xkbcommon.KeyRight:
        // Handle right arrow
    case xkbcommon.KeyReturn:
        // Handle return/enter
    }
}

func (app *MyApp) Focus(win *window.Window, input *window.Input) {
    // Handle focus gained
}

func main() {
    display, _ := window.DisplayCreate(nil)
    app := &MyApp{}
    app.window = window.Create(display)
    
    // Set keyboard handler
    app.window.SetKeyboardHandler(app)
    
    app.widget = app.window.AddWidget(app)
    app.window.ScheduleResize(800, 600)
    
    window.DisplayRun(display)
}
```

## Test Application

A test application is provided in `go-wayland-keyboard-test/main.go` that demonstrates keyboard event handling:

```bash
./run-keyboard-test.sh
```

The test app displays:
- Last key pressed
- Key code (virtual key code)
- Active modifiers
- Unicode character (if printable)

## Technical Notes

### NSEvent Key Codes vs Unicode

The implementation provides both:
1. **vKey**: macOS virtual key code (hardware-independent, layout-independent)
2. **code**: Unicode character (layout-dependent, affected by modifiers)

Applications should typically use `vKey` for command shortcuts and game controls, and `code` for text input.

### Modifier Flag Mapping

NSEvent modifier flags are mapped as follows:

| NSEvent Flag | Bit Position | window.ModType |
|--------------|--------------|----------------|
| NSEventModifierFlagShift | 1 << 17 | ModShiftMask (0x01) |
| NSEventModifierFlagControl | 1 << 18 | ModControlMask (0x04) |
| NSEventModifierFlagOption | 1 << 19 | ModAltMask (0x02) |
| NSEventModifierFlagCommand | 1 << 20 | (not mapped) |
| NSEventModifierFlagCapsLock | 1 << 16 | (not mapped) |

### Thread Safety

Keyboard events are captured on the main thread via NSEvent local monitor and dispatched to Go callbacks. The window registry uses RWMutex for thread-safe access.

### Event Monitoring

The implementation uses `addLocalMonitorForEventsMatchingMask:` which only captures events for the application's own windows. This is appropriate for application keyboard input and doesn't interfere with system-wide keyboard handling.

## Compatibility

This implementation is compatible with the existing window package API and follows the same patterns used in:
- Linux implementation (window/window_linux.go)
- Windows implementation (window/window_windows.go)

Applications using the `KeyboardHandler` interface will work across all platforms with minimal platform-specific code.

## Future Enhancements

Potential improvements:
1. Add constants for Command and CapsLock modifiers
2. Support for IME (Input Method Editor) text input
3. Dead key handling for international keyboards
4. Key repeat rate configuration
5. Modifier-only key events (NSEventTypeFlagsChanged handling)
