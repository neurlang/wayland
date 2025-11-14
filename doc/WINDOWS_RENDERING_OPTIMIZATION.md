# Windows Rendering Optimization

## Problem
The original `redrawer` function in `window/window_windows.go` was extremely slow because it:
- Divided the texture into 512x512 patches
- Analyzed each patch to find rectangles of the same color
- Drew hundreds or thousands of individual rectangles using `FillRect`
- Used complex goroutine-based parallelization to process patches

This approach was inefficient for texture rendering on Windows.

## Solution
Replaced the entire rendering pipeline with a single Windows GDI API call: **StretchDIBits**

### Key Changes

1. **Direct Bitmap Transfer**: Instead of drawing many rectangles, we now push the entire RGBA buffer to the screen in one API call.

2. **Simplified Code**: Removed ~200 lines of complex rectangle-finding logic and replaced with ~30 lines of straightforward bitmap rendering.

3. **Windows API Used**: `StretchDIBits` from gdi32.dll
   - Transfers Device Independent Bitmap (DIB) data directly to device context
   - Single syscall instead of thousands of GDI operations
   - Hardware-accelerated when available

### Technical Details

The new implementation:
- Uses `BITMAPINFO` structure to describe the 32-bit RGBA format
- Sets negative height (-h) for top-down bitmap orientation
- Uses `SRCCOPY` raster operation for direct pixel copy
- Maintains the hash-based change detection to avoid redundant redraws

### Performance Impact

Expected improvements:
- **10-100x faster** rendering depending on texture complexity
- Reduced CPU usage (no goroutine overhead, no rectangle analysis)
- Lower memory usage (no rectangle tracking structures)
- Smoother frame rates for animated content

### Files Modified

1. `window/window_windows.go`
   - Added Windows GDI syscall declarations
   - Replaced `redrawer()` function
   - Removed helper functions: `reduceDrawRectsNew`, `findMaxRectNew`, `iterateCoordinates`, `sameColor`, `getColor`

2. `window/widget_windows.go`
   - Removed `drawnRects` field from Widget struct
   - Updated related functions to remove rectangle tracking

## Testing

Build verification:
```bash
go build -v ./window
```

The code compiles successfully with no errors.
