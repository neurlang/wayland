# Darwin Window Features

This document describes the macOS window implementation features.

## Overview

The Darwin window implementation provides native macOS window support with:
- Window creation and management
- Automatic redraw loop (~60 FPS)
- Mouse hover tracking
- Window resizing
- Bitmap rendering

## Architecture

### File Structure

1. **window_cgo_darwin.go** - CGO layer with Objective-C code
   - Build tag: `// +build darwin,cgo`
   - Contains all Objective-C/Cocoa code
   - Provides C functions wrapped in Go

2. **window_darwin.go** - High-level Go implementation
   - Build tag: `// +build darwin`
   - Pure Go code that calls CGO functions
   - Implements Window and Display types

## Features

### 1. Window Creation

Windows are created with a default size of 200x200 pixels:

```go
w := Create(display)
// Creates a 200x200 window
```

### 2. Automatic Redraw Loop

A timer runs at ~60 FPS (every 16ms) to automatically trigger redraws:

```c
NSTimer* timer = [NSTimer scheduledTimerWithTimeInterval:0.016
                                repeats:YES
                                  block:^(NSTimer * _Nonnull t) {
    goScheduleRedraw(windowPtr);
}];
```

The redraw callback calls `ScheduleRedraw()` and `Redraw()` on the window.

### 3. Mouse Hover Tracking

Mouse tracking is implemented using NSView's event system with a custom `MouseTrackingView` class:

```objc
@interface MouseTrackingView : BitmapView
- (void)mouseMoved:(NSEvent *)event;
- (void)mouseEntered:(NSEvent *)event;
- (void)mouseExited:(NSEvent *)event;
- (void)updateTrackingAreas;
@end
```

The view creates an NSTrackingArea that captures mouse movement:

```objc
NSTrackingArea* trackingArea = [[NSTrackingArea alloc]
    initWithRect:[self bounds]
    options:(NSTrackingMouseMoved | NSTrackingMouseEnteredAndExited | 
             NSTrackingActiveInKeyWindow | NSTrackingInVisibleRect)
    owner:self
    userInfo:nil];
```

When mouse events occur, they're forwarded to Go:

```objc
- (void)mouseMoved:(NSEvent *)event {
    if (goWindowPtr) {
        NSPoint location = [self convertPoint:[event locationInWindow] fromView:nil];
        goMouseMotion(goWindowPtr, location.x, location.y);
    }
}
```

Motion events are dispatched to all widgets:

```go
//export goMouseMotion
func goMouseMotion(windowPtr unsafe.Pointer, x, y C.float) {
    for widget := range window.widgets {
        if widget.handler != nil {
            widget.handler.Motion(widget, window.input, timestamp, float32(x), float32(y))
        }
    }
}
```

### 4. Window Resizing

Windows can be resized programmatically:

```go
window.ScheduleResize(200, 200)
```

This:
- Updates the NSWindow frame
- Resizes all widgets
- Triggers resize callbacks
- Invalidates drawn content

### 5. Bitmap Rendering

Rendering uses a custom `BitmapView` class that properly handles Cairo's BGRA format:

```objc
@interface BitmapView : NSView {
    @public
    CGImageRef currentImage;
    NSLock* imageLock;
}
- (void)updateImage:(CGImageRef)newImage;
- (void)drawRect:(NSRect)dirtyRect;
@end
```

The view stores the CGImage and draws it in `drawRect:`:

```objc
- (void)drawRect:(NSRect)dirtyRect {
    [imageLock lock];
    if (currentImage) {
        CGContextRef context = [[NSGraphicsContext currentContext] CGContext];
        CGContextSaveGState(context);
        
        // Flip coordinate system (CGImage is top-down, NSView is bottom-up)
        CGContextTranslateCTM(context, 0, self.bounds.size.height);
        CGContextScaleCTM(context, 1.0, -1.0);
        
        // Draw the image
        CGContextDrawImage(context, self.bounds, currentImage);
        
        CGContextRestoreGState(context);
    }
    [imageLock unlock];
}
```

Bitmap data is converted from Cairo's BGRA format:

```c
CGImageRef cgImage = CGImageCreate(
    width, height, 8, 32, width * 4,
    colorSpace,
    kCGImageAlphaPremultipliedFirst | kCGBitmapByteOrder32Little,
    provider, NULL, false, kCGRenderingIntentDefault
);
```

Key improvements:
- **Thread-safe**: Image updates are protected by NSLock
- **Proper coordinate system**: Flips Y-axis to match Cairo's top-down format
- **Modern API**: Uses `CGContextDrawImage` instead of deprecated `lockFocus`
- **Memory management**: Copies bitmap data to ensure it stays valid

## CGO Functions

### Window Management
- `darwin_createWindow()` - Create NSWindow with callbacks
- `darwin_destroyWindow()` - Close and cleanup window
- `darwin_setTitle()` - Set window title
- `darwin_resizeWindow()` - Resize window content

### Event Loop
- `darwin_runMainLoop()` - Start NSApp main loop
- `darwin_stopMainLoop()` - Stop NSApp main loop

### Rendering
- `darwin_drawBitmap()` - Draw BGRA bitmap to window
- `darwin_startRedrawTimer()` - Start 60 FPS timer
- `darwin_enableMouseTracking()` - Enable mouse tracking

### Input
- `darwin_getMousePosition()` - Get current mouse position
- `goMouseMotion()` - Go callback for mouse motion
- `goScheduleRedraw()` - Go callback for redraw timer

## Window Registry

A global registry maps Go window pointers to window handles:

```go
var (
    windowRegistry = make(map[unsafe.Pointer]*darwinWindowHandle)
    windowMutex    sync.RWMutex
)
```

This allows C callbacks to find the corresponding Go window object.

## Benefits

1. **No Duplicate Symbols** - All Objective-C code in one file
2. **Automatic Updates** - 60 FPS redraw loop
3. **Smooth Mouse Tracking** - Polled every frame
4. **Clean Separation** - CGO isolated from business logic
5. **Thread Safe** - Registry protected by mutex

## Usage Example

```go
// Create display
display, _ := DisplayCreate(os.Args)

// Create window (200x200)
window := Create(display)
window.SetTitle("My App")

// Add widget with handler
widget := window.AddWidget(myHandler)

// Run event loop (blocks)
DisplayRun(display)
```

The window will automatically:
- Redraw at 60 FPS
- Track mouse movement
- Dispatch motion events to widgets
- Handle window resize events
