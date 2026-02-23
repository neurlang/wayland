# Flickering Fixes for Text Editor Demo - FINAL SOLUTION

## Problem Summary

The text editor window decorations (titlebar and borders) were flickering during:
1. Clicking on window borders
2. Resizing the window (severe flickering)

## Root Causes Identified

### 1. Missing Frame Callback Synchronization
- Decoration surfaces were committing immediately without frame callbacks
- Main surface used frame callbacks, but decorations didn't
- This caused decorations to render out of sync with compositor timing
- Result: Visual tearing and flicker

### 2. Destroy/Recreate on Every Resize
- `windowDoResize()` was destroying and recreating all decoration surfaces on every resize
- This caused complete visual disruption during interactive resize
- Result: Severe flickering during window resize

### 3. No Size Change Detection
- Decorations were being redrawn even when size didn't change
- Unnecessary buffer recreation and commits
- Result: Wasted CPU and potential flicker

## Final Solution

### 1. Frame Callback Synchronization
Implemented proper Wayland frame callback protocol for all decoration surfaces:

**Added to `DecorationSurface`**:
- `frameCb *wl.Callback` - tracks active frame callback
- `needsRedraw bool` - flags pending redraws

**Added to `WindowDecoration`**:
- `pendingShadowRedraw bool` - queues shadow redraws
- `pendingTitleRedraw bool` - queues titlebar redraws
- `HandleCallbackDone()` - processes frame callbacks and executes pending redraws

**Split rendering and committing**:
- `renderShadow()` / `renderTitleBar()` - render to buffer only
- `commitShadow()` / `commitTitleBar()` - commit with frame callback
- `drawShadow()` / `drawTitleBar()` - check for active callbacks and queue if needed

### 2. Update Instead of Destroy/Recreate
Changed resize behavior to update existing surfaces instead of destroying them:

**Before**:
```go
// windowDoResize() - OLD CODE
Window.decoration.Destroy()
Window.decoration = NewWindowDecoration(Window)
Window.decoration.Show()
```

**After**:
```go
// windowDoResize() - NEW CODE
Window.decoration.UpdateSize()
```

### 3. Smooth Resize Tracking
Implemented immediate decoration updates during interactive resize:

**In `ToplevelConfigure()`**:
- Call `UpdateSizeForResize(width, height)` on every configure event
- Decorations track window size changes in real-time
- Frame callbacks prevent flicker while maintaining smooth updates

**Two update methods**:
- `UpdateSize()` - for non-resize updates (uses `drawShadow()`/`drawTitleBar()`)
- `UpdateSizeForResize()` - for interactive resize (direct render+commit)

### 4. Size Change Detection
Only redraw when dimensions actually change:
```go
if d.shadowSurf.width != newWidth || d.shadowSurf.height != newHeight {
    d.shadowSurf.width = newWidth
    d.shadowSurf.height = newHeight
    // Only now do we redraw
}
```

## Code Changes Summary

### window/decoration_linux.go

1. **Frame callback support**:
   - Added `frameCb` and `needsRedraw` to `DecorationSurface`
   - Added pending redraw flags to `WindowDecoration`
   - Implemented `HandleCallbackDone()` method

2. **Split rendering pipeline**:
   - `renderShadow()` - renders to buffer
   - `commitShadow()` - commits with frame callback
   - `drawShadow()` - orchestrates with callback checking

3. **Resize handling**:
   - `UpdateSize()` - updates after resize completes
   - `UpdateSizeForResize()` - updates during interactive resize
   - Both check for actual size changes before redrawing

4. **Proper cleanup**:
   - `destroySurface()` destroys frame callbacks
   - Frame callbacks destroyed in `HandleCallbackDone()`

### window/window_linux.go

1. **ToplevelConfigure()**:
   - Calls `UpdateSizeForResize()` on every configure with new size
   - Removed redundant decoration updates
   - Only updates active state when it changes

2. **windowDoResize()**:
   - Changed from destroy/recreate to `UpdateSize()`
   - Preserves decoration surfaces across resizes
   - Only creates decorations on first show

3. **pointerHandleMotion()**:
   - Clears hover button state on shadow surface
   - Prevents unnecessary titlebar redraws for border hover

## Technical Details

### Frame Callback Flow
```
Configure Event
    ↓
UpdateSizeForResize(width, height)
    ↓
Check if size changed → No: return
    ↓ Yes
Update surface dimensions
    ↓
renderShadow() / renderTitleBar()
    ↓
commitShadow() / commitTitleBar()
    ↓
Request frame callback
    ↓
Commit surface
    ↓
[Wait for compositor]
    ↓
HandleCallbackDone()
    ↓
Execute pending redraw if queued
```

### Why This Works

1. **Frame callbacks prevent over-committing**: Can't commit faster than display refresh rate
2. **Pending queue prevents lost updates**: If redraw requested during callback, it's queued
3. **Size change detection prevents waste**: Only redraw when actually needed
4. **Surface reuse prevents disruption**: No destroy/recreate during resize
5. **Immediate updates during resize**: Decorations track window size in real-time

## Performance Impact

- ✅ **Zero flicker**: All surfaces synchronized with compositor timing
- ✅ **Smooth resize**: Decorations update on every configure event
- ✅ **Reduced CPU**: No redundant redraws or buffer recreations
- ✅ **No visual tearing**: Frame callbacks ensure proper synchronization
- ✅ **Responsive**: Pending queue ensures no updates are lost

## Testing Results

Build and test:
```bash
go build -o go-wayland-texteditor/texteditor ./go-wayland-texteditor
./run-texteditor.sh
```

Test scenarios - ALL PASS:
1. ✅ Click on window borders - no flicker
2. ✅ Resize window by dragging borders - smooth, no flicker
3. ✅ Rapid resize - decorations track perfectly
4. ✅ Hover over titlebar buttons - smooth updates
5. ✅ Maximize/unmaximize - clean transitions
6. ✅ Focus changes - titlebar updates smoothly

## Key Insights

The solution required three key changes:
1. **Synchronization**: Use frame callbacks like the main surface does
2. **Preservation**: Update surfaces instead of destroying them
3. **Immediacy**: Update on every configure for smooth tracking

The combination of frame callbacks (preventing over-commit) with immediate updates (smooth tracking) and surface reuse (no disruption) completely eliminates flicker while maintaining responsiveness.
