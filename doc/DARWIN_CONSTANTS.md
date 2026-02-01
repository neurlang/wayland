# macOS (Darwin) Window Constants

## Overview

This document describes the Darwin-specific constants for cursor types and buffer types in the window package.

## Files Created

### cursor_darwin.go
Defines cursor type constants that map to NSCursor types in Cocoa.

### buffer_darwin.go
Defines buffer type constants for memory management.

## Cursor Constants

### Required Cursors (Cross-Platform)

| Constant | Value | macOS Equivalent | Description |
|----------|-------|------------------|-------------|
| `CursorHand1` | 1 | `NSCursor.pointingHandCursor` | Pointing hand (clickable) |
| `CursorLeftPtr` | 2 | `NSCursor.arrowCursor` | Default arrow pointer |
| `CursorIbeam` | 3 | `NSCursor.IBeamCursor` | Text selection cursor |

### Resize Cursors

| Constant | Value | Description |
|----------|-------|-------------|
| `CursorBottomLeft` | 4 | Resize bottom-left corner |
| `CursorBottomRight` | 5 | Resize bottom-right corner |
| `CursorBottom` | 6 | Resize bottom edge |
| `CursorLeft` | 8 | Resize left edge |
| `CursorRight` | 9 | Resize right edge |
| `CursorTopLeft` | 10 | Resize top-left corner |
| `CursorTopRight` | 11 | Resize top-right corner |
| `CursorTop` | 12 | Resize top edge |

### Interaction Cursors

| Constant | Value | macOS Equivalent | Description |
|----------|-------|------------------|-------------|
| `CursorDragging` | 7 | `NSCursor.openHandCursor` | Open hand (draggable) |
| `CursorWatch` | 13 | `NSCursor.operationNotAllowedCursor` | Wait/busy |
| `CursorDndMove` | 14 | Drag and drop move | Move operation |
| `CursorDndCopy` | 15 | Drag and drop copy | Copy operation |
| `CursorDndForbidden` | 16 | `NSCursor.operationNotAllowedCursor` | Not allowed |
| `CursorBlank` | 17 | Hidden cursor | Invisible cursor |

### macOS-Specific Cursors

| Constant | Value | macOS Equivalent | Description |
|----------|-------|------------------|-------------|
| `CursorCrosshair` | 18 | `NSCursor.crosshairCursor` | Crosshair for precision |
| `CursorClosedHand` | 19 | `NSCursor.closedHandCursor` | Closed hand (grabbing) |
| `CursorDisappearingItem` | 20 | `NSCursor.disappearingItemCursor` | Item being removed |
| `CursorResizeLeftRight` | 21 | `NSCursor.resizeLeftRightCursor` | Horizontal resize |
| `CursorResizeUpDown` | 22 | `NSCursor.resizeUpDownCursor` | Vertical resize |

## Buffer Type Constants

### BufferTypeShm

```go
const BufferTypeShm = 0
```

**Description**: Shared memory buffer type for macOS.

**Usage**: On macOS, we use shared memory buffers similar to Wayland's `wl_shm` protocol. This is the default and only buffer type currently supported.

**Implementation**: Buffers are allocated in memory and rendered using `CGBitmapContext` to the window's view.

## Cross-Platform Comparison

### Cursor Values

| Cursor | Linux | Windows | macOS |
|--------|-------|---------|-------|
| `CursorHand1` | 11 | 1 | 1 |
| `CursorLeftPtr` | 4 | 2 | 2 |
| `CursorIbeam` | 10 | 3 | 3 |

**Note**: Values differ across platforms because they map to different underlying cursor systems (X11, Win32, Cocoa).

### Buffer Type Values

| Buffer Type | Linux | Windows | macOS |
|-------------|-------|---------|-------|
| `BufferTypeShm` | 1 | 0 | 0 |

## Usage Example

### Setting Cursor Type

```go
// In a widget handler
func (h *MyHandler) Enter(widget *window.Widget, input *window.Input, x, y float32) {
    // Set cursor to hand when hovering over clickable area
    widget.SetCursor(window.CursorHand1)
}

func (h *MyHandler) Leave(widget *window.Widget, input *window.Input) {
    // Reset to default cursor
    widget.SetCursor(window.CursorLeftPtr)
}
```

### Buffer Type

```go
// Buffer type is automatically set when creating a window
// On macOS, BufferTypeShm is the default
window.SetBufferType(window.BufferTypeShm)
```

## Implementation Notes

### Cursor Implementation

To implement cursor changes on macOS, you would use:

```objective-c
// In Cocoa/Objective-C
switch (cursorType) {
    case 1: // CursorHand1
        [[NSCursor pointingHandCursor] set];
        break;
    case 2: // CursorLeftPtr
        [[NSCursor arrowCursor] set];
        break;
    case 3: // CursorIbeam
        [[NSCursor IBeamCursor] set];
        break;
    // ... etc
}
```

### Buffer Implementation

Buffers on macOS are:
1. Allocated as byte arrays in Go
2. Filled with BGRA pixel data
3. Rendered using `CGBitmapContext`
4. Displayed in the window's `NSView`

## Future Enhancements

- [ ] Implement cursor changing in window_darwin.go
- [ ] Add custom cursor support (from image data)
- [ ] Animated cursor support
- [ ] Hardware cursor support
- [ ] Additional buffer types (if needed)

## References

- [NSCursor Documentation](https://developer.apple.com/documentation/appkit/nscursor)
- [CGBitmapContext Documentation](https://developer.apple.com/documentation/coregraphics/cgbitmapcontext)
