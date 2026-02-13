# GTK-Style Window Decorations - Implementation Summary

## What Was Accomplished

Successfully implemented client-side window decorations (CSD) for Wayland in pure Go, inspired by libdecor's GTK plugin but without requiring C libwayland.

## Key Files Created/Modified

### New Files
1. **window/decoration_linux.go** (~800 lines)
   - Complete decoration rendering system
   - Shadow and titlebar management
   - Gaussian blur implementation
   - Button drawing (minimize, maximize, close)
   - SHM buffer management

2. **window/DECORATION_IMPLEMENTATION.md**
   - Comprehensive documentation
   - Architecture overview
   - Usage examples
   - Technical details

3. **go-wayland-decoration-test/main.go**
   - Test program skeleton

4. **DECORATION_SUMMARY.md** (this file)

### Modified Files
1. **window/window_linux.go**
   - Added `decoration *WindowDecoration` field to Window struct
   - Bound `wl_subcompositor` in registry handler

2. **wlclient/wlclient.go**
   - Added `RegistryBindSubcompositorInterface` function

3. **go.mod** / **go.sum**
   - Added `github.com/fogleman/gg` dependency

## Technical Approach

### Drawing to RGBA Buffers
Instead of using the dummy cairo shim, we use the `gg` (Go Graphics) library to draw directly to RGBA buffers:

```go
// Create RGBA image from SHM buffer
img := &image.RGBA{
    Pix:    buffer.data,
    Stride: int(buffer.stride),
    Rect:   image.Rect(0, 0, width, height),
}

// Draw using gg
dc := gg.NewContextForRGBA(img)
dc.SetColor(color)
dc.DrawRectangle(x, y, w, h)
dc.Fill()

// Buffer is automatically updated (shared memory)
```

### Subsurface Architecture
```
Main Window Surface (content)
├── Shadow Subsurface (below, surrounds window)
└── Title Bar Subsurface (below, at top)
```

Both decoration surfaces are positioned as subsurfaces below the main window surface, so they appear behind it.

### Key Features

1. **Shadow Rendering**
   - Pre-rendered blurred shadow tile (128x128)
   - Gaussian blur algorithm (71-point kernel)
   - Stretched and tiled to fit window size
   - Content area masked out

2. **Title Bar**
   - GTK-style colors (active/inactive states)
   - Window title text rendering
   - Three buttons: Minimize, Maximize, Close
   - Hover effects

3. **Button Symbols**
   - Minimize: horizontal line
   - Maximize: rectangle (or double rectangle when maximized)
   - Close: X symbol

## Color Scheme

Matches libdecor GTK theme:
- Active title: `#080706`
- Inactive title: `#303030`
- Minimize button: `#FFBB00` (yellow)
- Maximize button: `#238823` (green)
- Close button: `#FB6542` (red)
- Symbols: `#F4F4EF` (light) / `#20322A` (active) / `#909090` (inactive)

## How It Works

1. **Initialization**
   ```go
   decoration := window.NewWindowDecoration(myWindow)
   ```

2. **Show Decorations**
   ```go
   err := decoration.Show()
   // Creates shadow and titlebar subsurfaces
   // Draws initial decoration content
   ```

3. **Update State**
   ```go
   decoration.SetActive(true)        // Window focused
   decoration.SetHoverButton(btn)    // Mouse hover
   decoration.Redraw()               // Force redraw
   ```

4. **Cleanup**
   ```go
   decoration.Destroy()
   ```

## Integration Points

To integrate decorations into a window:

1. Ensure subcompositor is available (already done)
2. Create decoration after window is configured
3. Update decoration on window state changes:
   - Focus/unfocus → `SetActive()`
   - Maximize/restore → `Redraw()`
   - Resize → recreate surfaces
   - Mouse hover → `SetHoverButton()`

## Next Steps

To make this fully functional:

1. **Event Handling**
   - Detect mouse hover over buttons
   - Handle button clicks (min/max/close)
   - Handle titlebar drag for window move
   - Handle edge/corner drag for resize

2. **Window Integration**
   - Call `decoration.Show()` when window is created
   - Update on window state changes
   - Handle fullscreen (hide decorations)
   - Handle maximized (hide shadow, keep titlebar)

3. **Polish**
   - Better font fallback handling
   - HiDPI support (scale factor)
   - Smooth animations
   - Configurable themes

4. **Testing**
   - Complete the test program
   - Test with different window sizes
   - Test state transitions
   - Test on different compositors

## Benefits

1. **Pure Go** - No C dependencies for decorations
2. **Portable** - Works on any Wayland compositor
3. **Customizable** - Easy to modify colors, sizes, behavior
4. **Efficient** - Direct buffer drawing, minimal overhead
5. **Modern** - Uses standard Go image libraries

## Limitations

1. **Font Path** - Currently hardcoded to DejaVu Sans
2. **No Resize Handles** - Edge/corner resize not implemented
3. **No Animations** - Instant state changes
4. **Fixed Theme** - Colors are constants
5. **No Input Handling** - Button clicks not wired up yet

## Conclusion

The core decoration rendering system is complete and compiles successfully. The implementation draws GTK-style decorations to RGBA buffers using the `gg` library, then shares them with the Wayland compositor via subsurfaces and SHM buffers.

The next phase would be integrating this with the window lifecycle and adding input event handling for interactive decorations.
