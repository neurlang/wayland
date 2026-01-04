# Darwin Bitmap Rendering and Mouse Hover Fix

## Issues Fixed

### 1. Bitmap Rendering Not Working

**Problem**: Bitmap content was not displaying in the window despite the window showing up at the correct size.

**Root Causes**:
- Using deprecated `lockFocus`/`unlockFocus` API which doesn't work properly on modern macOS
- Incorrect pixel format handling for Cairo's BGRA format
- No proper view invalidation after drawing
- Coordinate system mismatch between Cairo (top-down) and NSView (bottom-up)

**Solution**: Created a custom `BitmapView` class that:
- Stores the CGImage in an instance variable
- Implements `drawRect:` to properly render using Core Graphics
- Flips the coordinate system to match Cairo's format
- Uses thread-safe locking for image updates
- Properly triggers view refresh with `setNeedsDisplay:`

### 2. Mouse Hover Events Not Working

**Problem**: Mouse hover/motion events were not being captured or dispatched to widgets.

**Root Causes**:
- NSTrackingArea was created but events weren't being handled
- No override of `mouseMoved:` method to capture events
- Polling-based approach in redraw timer was inefficient

**Solution**: Created a custom `MouseTrackingView` class that:
- Extends `BitmapView` with mouse event handling
- Overrides `mouseMoved:`, `mouseEntered:`, and `mouseExited:`
- Implements `updateTrackingAreas` to maintain tracking area on resize
- Directly calls Go callback `goMouseMotion()` when events occur
- Uses event-driven approach instead of polling

## Implementation Details

### Custom View Hierarchy

```
NSView
  └── BitmapView (handles bitmap rendering)
        └── MouseTrackingView (adds mouse event handling)
```

### BitmapView Class

```objc
@interface BitmapView : NSView {
    @public
    CGImageRef currentImage;
    NSLock* imageLock;
}
- (void)updateImage:(CGImageRef)newImage;
- (void)drawRect:(NSRect)dirtyRect;
- (BOOL)isFlipped;
@end
```

Key features:
- Thread-safe image storage with NSLock
- Proper coordinate system handling with `isFlipped`
- Efficient rendering using `CGContextDrawImage`
- Automatic memory management with retain/release

### MouseTrackingView Class

```objc
@interface MouseTrackingView : BitmapView {
    @public
    void* goWindowPtr;
}
- (void)mouseMoved:(NSEvent *)event;
- (void)mouseEntered:(NSEvent *)event;
- (void)mouseExited:(NSEvent *)event;
- (void)updateTrackingAreas;
@end
```

Key features:
- Stores Go window pointer for callbacks
- Converts mouse coordinates to view space
- Maintains tracking area on view resize
- Forwards events to Go immediately

### Bitmap Format Handling

Cairo uses BGRA premultiplied format, which maps to:
```c
kCGImageAlphaPremultipliedFirst | kCGBitmapByteOrder32Little
```

This ensures correct color interpretation:
- B (blue) at byte 0
- G (green) at byte 1
- R (red) at byte 2
- A (alpha) at byte 3

### Memory Management

The bitmap data is copied before creating the CGDataProvider:
```c
void* dataCopy = malloc(dataSize);
memcpy(dataCopy, data, dataSize);

CGDataProviderRef provider = CGDataProviderCreateWithData(
    NULL, dataCopy, dataSize,
    ^(void *info, const void *data, size_t size) {
        free((void*)data);  // Cleanup when done
    }
);
```

This ensures the data remains valid even after the Go buffer is garbage collected.

## Compilation Fixes

Several Objective-C syntax issues were also resolved:

1. **Class Declaration Order**: BitmapView must be declared before MouseTrackingView since the latter inherits from the former
2. **Type Casting**: Added explicit `(float)` casts for NSPoint coordinates when calling Go callbacks
3. **Method Calls**: Changed `[self bounds]` to `self.bounds` in appropriate contexts
4. **Block vs Function Pointer**: Replaced block syntax with a C function pointer for CGDataProviderReleaseDataCallback to avoid compatibility issues

## Testing

To test the fixes on macOS:

```bash
# Build one of the example applications
go build -tags darwin,cgo ./go-wayland-simple-shm

# Or build the smoke demo
go build -tags darwin,cgo ./go-wayland-smoke

# Run it
./go-wayland-simple-shm
```

Expected behavior:
- Window appears with correct size
- Bitmap content is visible and updates smoothly at 60 FPS
- Mouse hover events trigger visual feedback in real-time
- No crashes or memory leaks
- Proper coordinate handling (no flipped/inverted rendering)

## Benefits

1. **Proper Rendering**: Bitmap content now displays correctly
2. **Smooth Updates**: 60 FPS redraw loop with efficient rendering
3. **Responsive Mouse**: Event-driven mouse tracking with no lag
4. **Thread Safe**: Proper locking prevents race conditions
5. **Modern API**: Uses current macOS APIs, not deprecated ones
6. **Memory Safe**: Proper memory management prevents leaks

## Files Modified

- `window/window_cgo_darwin.go` - Added BitmapView and MouseTrackingView classes
- `doc/DARWIN_WINDOW_FEATURES.md` - Updated documentation

## Compatibility

- Requires macOS 10.10+ (for modern Cocoa APIs)
- Works with both Intel and Apple Silicon Macs
- Compatible with Go 1.16+ with CGO enabled
