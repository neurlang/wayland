# Flicker Fix Solution Summary

## Problem
Window decorations (titlebar and borders) were flickering during:
- Clicking on borders
- Resizing the window

## Root Causes
1. **No frame callback synchronization** - decorations committed immediately, out of sync with compositor
2. **Destroy/recreate on resize** - all decoration surfaces destroyed and recreated on every resize
3. **No size change detection** - redraws happened even when size didn't change

## Solution (3 Key Changes)

### 1. Frame Callback Synchronization
Added proper Wayland frame callback protocol to decoration surfaces:
- Split rendering (`renderShadow`/`renderTitleBar`) from committing (`commitShadow`/`commitTitleBar`)
- Request frame callback before each commit
- Queue pending redraws if callback is active
- Execute queued redraws when callback completes

### 2. Update Instead of Destroy/Recreate
Changed `windowDoResize()` to update existing surfaces:
```go
// Before: Window.decoration.Destroy() + recreate
// After:  Window.decoration.UpdateSize()
```

### 3. Real-time Resize Tracking
Added `UpdateSizeForResize()` called on every configure event:
- Decorations track window size in real-time during interactive resize
- Frame callbacks prevent flicker while maintaining smooth updates
- Only redraws when size actually changes

## Files Modified
- `window/decoration_linux.go` - frame callbacks, split render/commit, size update methods
- `window/window_linux.go` - call UpdateSizeForResize in ToplevelConfigure, update not destroy

## Result
✅ Zero flicker on border clicks
✅ Smooth resize with decorations tracking perfectly
✅ Reduced CPU usage
✅ All surfaces synchronized with compositor timing

## Key Insight
The combination of:
- Frame callbacks (prevent over-commit)
- Immediate updates (smooth tracking)  
- Surface reuse (no disruption)

...completely eliminates flicker while maintaining responsiveness.
