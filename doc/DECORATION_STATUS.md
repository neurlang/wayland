# Window Decoration Status

## Current Status: Debugging Compositor Disconnect

The GTK-style window decoration implementation is complete and compiles successfully, but the compositor disconnects immediately after decorations are created.

## What Works

✅ All decoration rendering code is implemented
✅ Shadow rendering with Gaussian blur
✅ Title bar with buttons (minimize, maximize, close)
✅ Subsurface creation and positioning
✅ SHM buffer management and drawing with `gg` library
✅ Buffers are created and attached successfully

## The Problem

❌ Compositor closes connection with "connection reset by peer" immediately after parent surface commit
❌ Error occurs during the roundtrip after decoration creation
❌ All protocol calls appear correct but compositor rejects them

## Investigation Summary

We've tried multiple approaches:
1. ✅ Set subsurfaces to desync mode - no change
2. ✅ Don't commit subsurfaces individually - no change  
3. ✅ Keep SHM pools alive with buffers - no change
4. ✅ Start in sync mode, commit parent, then set desync - no change
5. ✅ Don't commit subsurfaces at all (let parent commit handle it) - no change

The compositor disconnects specifically when we commit the parent surface after:
- Creating two subsurfaces (shadow and titlebar)
- Attaching buffers to them
- Positioning them relative to parent

## Debug Output

```
Creating subsurfaces...
Creating decoration surface at (-24,-48) size 448x372
Parent surface: 0xc0001a4000
Decoration surface created
Creating decoration surface at (0,-24) size 400x24  
Parent surface: 0xc0001a4000
Decoration surface created
Drawing decorations...
[... buffer creation successful ...]
Attaching shadow surface (no commit yet)...
Shadow attached successfully
Attaching titlebar surface (no commit yet)...
Titlebar attached successfully
Committing parent surface 0xc0001a4000...
Setting subsurfaces to desync mode...
Decorations created successfully
Doing roundtrip after creating decorations...
Roundtrip failed: connection reset by peer
```

## Possible Remaining Issues

1. **Timing**: Creating decorations during Redraw callback might interfere with event loop
2. **Surface state**: Parent surface might be in a state that doesn't allow subsurfaces
3. **Compositor bug**: Some compositors might not handle subsurfaces on xdg_toplevel correctly
4. **Missing protocol step**: We might be missing a required protocol call
5. **Buffer format**: The compositor might not like ARGB8888 format for subsurfaces
6. **Size limits**: The shadow buffer is quite large (448x372 = 666KB)

## Next Steps to Try

1. Test with a minimal subsurface (just 10x10 solid color rectangle)
2. Try creating decorations BEFORE the first window commit
3. Test on different compositors (currently testing on unknown compositor)
4. Add more detailed protocol-level debugging
5. Check if subsurfaces need to be created in a specific order
6. Try using XRGB8888 format for both buffers (no alpha)

## Workaround

For now, windows work fine without decorations. The compositor's server-side decorations (SSD) are used instead.

## Files

- `window/decoration_linux.go` - Main implementation (~800 lines)
- `window/window_linux.go` - Integration (EnableDecorations method)
- `test-decorations-trace.go` - Test with windowtrace debugging
- `doc/DECORATION_IMPLEMENTATION.md` - Technical documentation
- `doc/DECORATION_SUMMARY.md` - Implementation summary
- `doc/DECORATION_INTEGRATION.md` - Integration guide
