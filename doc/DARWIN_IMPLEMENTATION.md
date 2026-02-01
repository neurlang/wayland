# macOS (Darwin) Window Implementation

## Overview

The `window_darwin.go` file provides native macOS window support using Cocoa/AppKit frameworks, matching the pattern used in `window_windows.go` and `window_linux.go`.

## Architecture

### Display Management
- **Display struct**: Manages multiple windows and the application lifecycle
- **DisplayCreate()**: Initializes the display
- **DisplayRun()**: Starts the NSApplication event loop
- **Exit()**: Stops the event loop

### Window Creation
- **Create(d *Display)**: Creates a new NSWindow with standard decorations
- Uses Cocoa's NSWindow with title bar, close, minimize, and resize buttons
- Default size: 800x600 pixels

### Event Handling

The implementation uses Objective-C classes to bridge between Cocoa events and Go:

1. **WindowDelegate**: Handles window-level events
   - Window close
   - Window resize

2. **WindowView**: Custom NSView subclass for input events
   - Keyboard events (keyDown, keyUp)
   - Mouse events (mouseDown, mouseUp, mouseMoved)
   - Right-click support

### Callback Mechanism

Events flow from Objective-C to Go through exported functions:

```
Cocoa Event → Objective-C Method → Go Callback → Widget Handler
```

**Exported Go Functions:**
- `goCloseCallback`: Window close event
- `goResizeCallback`: Window resize event
- `goKeyCallback`: Keyboard events
- `goMouseCallback`: Mouse events (buttons and motion)

### Drawing

- **drawBitmap()**: Renders BGRA pixel data to the window
- Uses CGBitmapContext for efficient bitmap rendering
- Supports automatic view refresh

## Key Features

✅ Window creation and destruction
✅ Title bar management
✅ Fullscreen support
✅ Maximize/restore
✅ Keyboard input
✅ Mouse input (buttons and motion)
✅ Window resizing
✅ Multiple windows support
✅ Proper cleanup on exit

## Differences from Linux/Windows

### vs Linux (Wayland)
- No Wayland protocol - uses native Cocoa APIs
- Direct window management instead of compositor communication
- Native macOS window decorations

### vs Windows
- Uses Objective-C instead of Win32 API
- Cocoa event loop instead of Windows message loop
- CGBitmapContext instead of StretchDIBits

## Build Requirements

```bash
# Requires macOS with Xcode Command Line Tools
xcode-select --install

# Build tags automatically select darwin implementation
GOOS=darwin go build ./window
```

## CGO Configuration

```go
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore
```

## Usage Example

```go
import "github.com/neurlang/wayland/window"

// Create display
display, err := window.DisplayCreate(os.Args)
if err != nil {
    log.Fatal(err)
}
defer display.Destroy()

// Create window
win := window.Create(display)
win.SetTitle("My macOS App")

// Add widget
widget := win.AddWidget(myHandler)

// Run event loop
window.DisplayRun(display)
```

## Implementation Notes

### Memory Management
- NSWindow is retained using CFBridgingRetain
- Released with CFBridgingRelease in destroyWindow
- Go window pointer passed to Objective-C for callbacks

### Thread Safety
- Display uses sync.RWMutex for window list access
- All Cocoa calls wrapped in @autoreleasepool

### Coordinate System
- Cocoa uses bottom-left origin
- Mouse coordinates converted to match expected top-left origin

### Key Mapping
- macOS key codes passed directly to handlers
- Applications should map to virtual key codes as needed

## Future Enhancements

Potential improvements:
- [ ] Retina display support (high DPI)
- [ ] Menu bar integration
- [ ] Dock icon customization
- [ ] Native file dialogs
- [ ] Drag and drop support
- [ ] Multi-monitor support
- [ ] Window shadows and effects
- [ ] Keyboard modifier state tracking
