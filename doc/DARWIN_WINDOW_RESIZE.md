# Darwin Window Resizing Implementation

## Overview
This document describes the implementation of window resizing support for the macOS (Darwin) platform in the `window` package, specifically for the go-wayland-web-browser demo.

## Changes Made

### 1. Window Delegate for Resize Events
Added a custom NSWindowDelegate class to handle window resize events:

- **Location**: `window/window_cgo_darwin.go`
- **Implementation**: Dynamic Objective-C class creation using the runtime API
- **Class Name**: `DarwinWindowDelegate`
- **Key Methods**:
  - `windowDidResize:` - Called when the window is resized by the user
  - `windowShouldClose:` - Allows the window to close properly

### 2. Go Callback for Resize Events
Added a new Go callback function that is invoked from the Objective-C delegate:

```go
//export goWindowResize
func goWindowResize(windowPtr unsafe.Pointer, width, height C.int)
```

This callback:
- Updates the window's internal width and height
- Notifies all widgets of the new size via `SetAllocation()`
- Calls the widget handler's `Resize()` method
- Triggers a redraw to reflect the new size

### 3. DarwinWindow Structure Update
Extended the `DarwinWindow` C struct to include:
- `windowDelegate` field to store the delegate instance

### 4. Window Creation Update
Modified `darwin_createWindow()` to:
- Create an instance of the delegate class
- Associate the DarwinWindow pointer with the delegate
- Set the delegate on the NSWindow

### 5. Window Destruction Update
Modified `darwin_destroyWindow()` to:
- Properly clean up the delegate instance
- Clear the delegate reference before closing the window

## How It Works

1. **User Resizes Window**: The user drags the window edge or corner to resize
2. **NSWindow Notification**: macOS sends a `windowDidResize:` notification
3. **Delegate Callback**: The delegate's `windowDidResize:` method is called
4. **Go Callback**: The delegate calls `goWindowResize()` with the new dimensions
5. **Widget Update**: All widgets are notified of the new size
6. **Handler Resize**: Each widget's handler `Resize()` method is called
7. **Redraw**: The window is redrawn with the new content

## Benefits

- **Responsive UI**: The web browser now properly responds to window resizing
- **Proper Layout**: Widgets are notified and can reflow their content
- **Native Feel**: Uses native macOS window resizing mechanisms
- **Thread Safe**: All UI operations are dispatched to the main thread

## Testing

The implementation was tested by building and running the go-wayland-web-browser:

```bash
go build -tags darwin,cgo ./go-wayland-web-browser/browser
./go-wayland-web-browser/browser/browser
```

Build succeeded with only minor warnings about `__bridge_transfer` casts (expected when not using ARC).

**Test Results:**
- ✅ Window resizing works correctly
- ✅ Content reflows properly when window is resized
- ✅ All widgets receive resize notifications
- ✅ Redraw is triggered after resize
- ✅ No crashes or memory leaks observed

## Compatibility

- **Platform**: macOS (Darwin) only
- **Build Tags**: Requires `darwin,cgo` build tags
- **Dependencies**: Cocoa, QuartzCore, CoreVideo frameworks
- **Go Version**: Compatible with Go 1.16+

## Future Enhancements

Potential improvements:
- Add minimum/maximum window size constraints
- Implement aspect ratio preservation
- Add resize throttling to reduce CPU usage during rapid resizing
- Support for live resize (updating during drag vs. after drag completes)
