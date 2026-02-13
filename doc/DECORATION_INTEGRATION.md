# Window Decoration Integration

## Overview

GTK-style window decorations are now automatically enabled for all Wayland windows created with the `window` package.

## What Was Changed

### Automatic Decoration Creation

Decorations are automatically created and managed for every window:

1. **Window Creation** (`window.Create()`)
   - Decorations are created if `wl_subcompositor` is available
   - Initially hidden until first configure event

2. **Window Configuration** (`ToplevelConfigure`)
   - Decorations shown/hidden based on window state
   - Active state updated based on focus
   - Fullscreen windows hide decorations
   - Normal/maximized windows show decorations

3. **Window Resize** (`windowDoResize`)
   - Decorations recreated with new window size
   - Shadow and titlebar surfaces adjusted automatically

4. **Window Destruction** (`Window.Destroy()`)
   - Decorations cleaned up properly
   - All subsurfaces and buffers freed

## Behavior

### Decoration States

- **Normal Window**: Shadow + titlebar with buttons
- **Maximized Window**: Titlebar only (no shadow)
- **Fullscreen Window**: No decorations
- **Focused Window**: Active colors (bright)
- **Unfocused Window**: Inactive colors (dimmed)

### Visual Elements

- **Shadow**: Blurred shadow surrounding the window
- **Title Bar**: 24px height with window title
- **Buttons**: Minimize (yellow), Maximize (green), Close (red)
- **Title Text**: Centered, using DejaVu Sans font

## Testing

Run any of the demo programs to see decorations:

```bash
# Build and run smoke demo
go build -o go-wayland-smoke/smoke ./go-wayland-smoke
./go-wayland-smoke/smoke

# Or use the run script
./run-smoke.sh
```

You should see:
- Window with shadow around it
- Title bar at the top with "smoke" title
- Three colored buttons on the right (min/max/close)
- Active/inactive color changes when focusing/unfocusing

## Current Limitations

1. **Buttons Not Interactive**: Clicking buttons doesn't do anything yet
2. **No Resize Handles**: Can't resize by dragging edges/corners
3. **No Titlebar Drag**: Can't move window by dragging titlebar
4. **Fixed Font**: Hardcoded to DejaVu Sans
5. **No HiDPI**: Scale factor not implemented yet

## Next Steps

To make decorations fully functional:

1. **Add Input Handling**
   - Detect pointer events on decoration surfaces
   - Handle button clicks (minimize, maximize, close)
   - Handle titlebar drag for window move
   - Handle edge/corner drag for resize

2. **Improve Rendering**
   - Better font fallback
   - HiDPI support
   - Smooth animations
   - Configurable themes

3. **Polish**
   - Respect compositor hints
   - Handle tiled window states
   - Better shadow rendering
   - Accessibility support

## Code Structure

```
window/
├── decoration_linux.go      # Decoration implementation
├── window_linux.go           # Window management (modified)
│   ├── Create()             # Creates decorations
│   ├── ToplevelConfigure()  # Shows/hides decorations
│   ├── windowDoResize()     # Updates decoration size
│   └── Destroy()            # Cleans up decorations
└── ...

wlclient/
└── wlclient.go              # Added subcompositor binding
```

## Disabling Decorations

If you want to disable decorations for a specific window, you can:

```go
window := window.Create(display)

// Disable decorations
if window.decoration != nil {
    window.decoration.Destroy()
    window.decoration = nil
}
```

Or check for `nil` before creating:

```go
// In window.Create(), comment out:
// if Display.subcompositor != nil {
//     Window.decoration = NewWindowDecoration(Window)
// }
```

## Debugging

If decorations don't appear:

1. Check if `wl_subcompositor` is available:
   ```go
   if display.subcompositor == nil {
       println("No subcompositor support")
   }
   ```

2. Check for errors in decoration creation:
   ```go
   if err := window.decoration.Show(); err != nil {
       println("Decoration error:", err)
   }
   ```

3. Verify font file exists:
   ```bash
   ls /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
   ```

## Conclusion

Window decorations are now fully integrated and work automatically for all Wayland windows. The implementation provides GTK-style decorations in pure Go without any C dependencies.
