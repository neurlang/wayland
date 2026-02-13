# GTK-Style Window Decorations Implementation

## Overview

This implementation provides client-side window decorations (CSD) for Wayland in pure Go, inspired by libdecor's GTK plugin but without requiring the C libwayland library.

The implementation uses the `gg` (Go Graphics) library to draw decorations directly to RGBA buffers, which are then shared with the Wayland compositor via shared memory (SHM).

## Architecture

### Core Components

1. **decoration_linux.go** - Main decoration implementation
   - `WindowDecoration` - Manages decoration surfaces and rendering
   - `DecorationSurface` - Represents shadow and titlebar subsurfaces
   - Drawing functions for shadows, titlebar, and buttons using `gg`

2. **Integration with window_linux.go**
   - Added `decoration *WindowDecoration` field to `Window` struct
   - Bound `wl_subcompositor` in registry handler
   - Decorations use subsurfaces positioned behind the main window surface

3. **wlclient/wlclient.go**
   - Added `RegistryBindSubcompositorInterface` function

## Features Implemented

### Shadow Rendering
- Gaussian blur algorithm ported from libdecor
- Pre-rendered shadow tile with blur effect
- Shadow surface surrounds the window with proper margins
- Content area is masked out from shadow

### Title Bar
- GTK-style title bar with window title text
- Three buttons: Minimize, Maximize, Close
- Active/inactive color states
- Hover effects on buttons
- Proper button symbols (minimize line, maximize rectangle, close X)

### Color Scheme
Matches libdecor's GTK theme:
- Active title: Dark gray (#080706)
- Inactive title: Medium gray (#303030)
- Minimize button: Yellow (#FFBB00)
- Maximize button: Green (#238823)
- Close button: Red (#FB6542)
- Symbols: Light (#F4F4EF) / Dark (#20322A) / Inactive (#909090)

## Technical Details

### Drawing Library

Uses `github.com/fogleman/gg` for 2D graphics:
- Direct drawing to RGBA buffers
- Anti-aliased lines and shapes
- Font rendering support
- Image manipulation for blur effects

### Subsurface Architecture
```
Main Window Surface
├── Shadow Subsurface (below, surrounds window)
└── Title Bar Subsurface (below, at top)
```

### Buffer Management
- Shared memory (SHM) buffers for decoration surfaces
- ARGB8888 format for transparency support
- XRGB8888 for opaque titlebar (optimization)
- Proper buffer lifecycle with release events

### Drawing Pipeline
1. Create cairo surface from SHM buffer data
2. Clear with transparent background
3. Draw decoration content (shadow/titlebar/buttons)
4. Attach buffer to wayland surface
5. Commit surface

## Usage

```go
// Create decoration for a window
decoration := window.NewWindowDecoration(myWindow)

// Show decorations
err := decoration.Show()

// Update active state
decoration.SetActive(true)

// Update hover state
decoration.SetHoverButton(window.ComponentButtonClose)

// Redraw when needed
decoration.Redraw()

// Clean up
decoration.Destroy()
```

## Differences from libdecor

1. **Pure Go** - No C dependencies, uses Go `gg` library for drawing
2. **Simplified** - Direct font rendering with `gg`, no Pango needed
3. **Direct** - Directly manages subsurfaces without plugin system
4. **Integrated** - Part of window package, not separate library
5. **RGBA Buffers** - Draws to standard Go image.RGBA, then shares via SHM

## Dependencies

- `github.com/fogleman/gg` - 2D graphics library
- `github.com/golang/freetype` - Font rendering (via gg)
- Standard Go image libraries

## Future Enhancements

1. **Pango Integration** - Better text rendering with proper font support
2. **Resize Handles** - Interactive window resizing from decoration edges
3. **Button Interaction** - Click handlers for min/max/close buttons
4. **Themes** - Configurable color schemes
5. **Animations** - Smooth transitions for hover effects
6. **HiDPI** - Proper scaling for high-resolution displays
7. **Accessibility** - Screen reader support

## Testing

A test program is provided in `go-wayland-decoration-test/` to demonstrate the decorations.

## References

- libdecor: https://gitlab.freedesktop.org/libdecor/libdecor
- Wayland protocols: https://wayland.freedesktop.org/
- Cairo graphics: https://www.cairographics.org/
