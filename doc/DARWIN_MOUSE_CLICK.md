# Darwin Mouse Click Implementation

## Overview
This document describes the implementation of mouse click support for the macOS (Darwin) platform in the `window` package, enabling link clicking and button interactions in the go-wayland-web-browser demo.

## Changes Made

### 1. Mouse Button Event Monitoring
Added NSEvent local monitor to capture mouse button events:

- **Location**: `window/window_cgo_darwin.go`
- **Event Types Captured**:
  - `NSEventTypeLeftMouseDown` - Left button pressed
  - `NSEventTypeLeftMouseUp` - Left button released
  - `NSEventTypeRightMouseDown` - Right button pressed
  - `NSEventTypeRightMouseUp` - Right button released

### 2. Button Code Mapping
Mapped macOS mouse events to Linux button codes for compatibility:

- **Left Button**: 272 (BTN_LEFT)
- **Right Button**: 273 (BTN_RIGHT)
- **State**: 1 for pressed, 0 for released

### 3. Go Callback for Button Events
Added a new Go callback function that is invoked from the Objective-C event monitor:

```go
//export goMouseButton
func goMouseButton(windowPtr unsafe.Pointer, button, state C.int, x, y C.float)
```

This callback:
- Retrieves the window from the registry
- Converts button state to `wl.PointerButtonState`
- Calls the widget handler's `Button()` method with proper parameters
- Provides timestamp for event ordering

### 4. DarwinWindow Structure Update
Extended the `DarwinWindow` C struct to include:
- `buttonMonitor` field to store the button event monitor instance

### 5. Window Creation Update
Modified `darwin_createWindow()` to:
- Create and register the button event monitor
- Store the monitor reference for cleanup

### 6. Window Destruction Update
Modified `darwin_destroyWindow()` to:
- Properly remove the button event monitor
- Clean up the monitor reference

## How It Works

1. **User Clicks Mouse**: The user clicks left or right mouse button in the window
2. **NSEvent Notification**: macOS sends a mouse button event
3. **Event Monitor Callback**: The local event monitor captures the event
4. **Coordinate Conversion**: Mouse coordinates are converted to Cairo's coordinate system (Y-axis flip)
5. **Go Callback**: `goMouseButton()` is called with button code, state, and position
6. **Widget Handler**: Each widget's handler `Button()` method is called
7. **Browser Processing**: The browser processes the click (e.g., navigating to a link)

## Button Event Flow in Browser

When a link is clicked:
1. Button pressed (state=1) → `clickSerial` is captured
2. Button released (state=0) → `ProcessPointerClick()` is called
3. Browser checks if cursor is over a link
4. If yes, navigates to the link URL

## Benefits

- **Interactive UI**: Users can now click buttons and links
- **Full Browser Functionality**: Link navigation works as expected
- **Context Menus**: Right-click support enables context menus
- **Native Feel**: Uses native macOS mouse event handling
- **Proper Event Ordering**: Timestamps ensure correct event sequence

## Testing

The implementation was tested by building and running the go-wayland-web-browser:

```bash
go build -tags darwin,cgo ./go-wayland-web-browser/browser
./go-wayland-web-browser/browser/browser
```

**Test Results:**
- ✅ Left mouse clicks detected correctly
- ✅ Right mouse clicks detected correctly
- ✅ Links can be clicked and navigated
- ✅ Button widgets respond to clicks
- ✅ Context menus work with right-click
- ✅ No crashes or event handling issues

## Compatibility

- **Platform**: macOS (Darwin) only
- **Build Tags**: Requires `darwin,cgo` build tags
- **Dependencies**: Cocoa framework
- **Go Version**: Compatible with Go 1.16+

## Future Enhancements

Potential improvements:
- Add middle mouse button support (button 274)
- Implement double-click detection
- Add mouse button modifier key support (Shift, Ctrl, etc.)
- Support for additional mouse buttons (back/forward buttons)
- Implement drag-and-drop with mouse buttons
