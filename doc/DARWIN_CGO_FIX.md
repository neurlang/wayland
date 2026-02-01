# Darwin CGO Duplicate Symbol Fix

## Problem

When compiling on macOS, clang reported duplicate symbol errors:

```
duplicate symbol '_setWindowMaximized' for architecture x86_64
duplicate symbol '_setWindowTitle' for architecture x86_64
... (15 duplicate symbols total)
```

## Root Cause

The C functions defined in the `import "C"` block of `window_darwin.go` had **external linkage** by default. When CGO compiles multiple Go files in the same package, these functions were being compiled multiple times, creating duplicate symbols during the linking phase.

### Why This Happens

1. CGO processes each `.go` file with C code separately
2. C functions without `static` have external linkage (visible across compilation units)
3. The linker sees multiple definitions of the same function
4. Result: "duplicate symbol" error

## Solution

Two fixes were required:

### Fix 1: Static Functions

Mark all C functions as `static` to give them **internal linkage**, making them local to each compilation unit.

All C function definitions in `window_darwin.go` were changed from:
```c
void functionName(...) {
```

To:
```c
static void functionName(...) {
```

### Fix 2: Header Guards for Objective-C Classes

Wrap Objective-C class definitions in header guards to prevent multiple definitions.

Added header guards around `@interface` and `@implementation` blocks:
```objective-c
#ifndef WINDOW_DARWIN_OBJC_CLASSES_H
#define WINDOW_DARWIN_OBJC_CLASSES_H

@interface WindowDelegate : NSObject <NSWindowDelegate>
// ...
@end

@implementation WindowDelegate
// ...
@end

@interface WindowView : NSView
// ...
@end

@implementation WindowView
// ...
@end

#endif // WINDOW_DARWIN_OBJC_CLASSES_H
```

### C Functions Fixed (9 total)

1. `createWindow` - Window creation
2. `destroyWindow` - Window destruction
3. `setWindowTitle` - Set window title
4. `setWindowCallbacks` - Setup event callbacks
5. `drawBitmap` - Render bitmap to window
6. `runMainLoop` - Start event loop
7. `stopMainLoop` - Stop event loop
8. `setWindowFullscreen` - Toggle fullscreen
9. `setWindowMaximized` - Toggle maximize

### Objective-C Classes Fixed (2 total)

1. `WindowDelegate` - Window event delegate
2. `WindowView` - Custom view for input handling

**Duplicate symbols resolved:**
- `_OBJC_IVAR_$_WindowDelegate._goWindow`
- `_OBJC_IVAR_$_WindowView._goWindow`
- `_OBJC_CLASS_$_WindowDelegate`
- `_OBJC_CLASS_$_WindowView`
- `_OBJC_METACLASS_$_WindowDelegate`
- `_OBJC_METACLASS_$_WindowView`

## Technical Details

### Static Linkage

```c
static void myFunction() {
    // This function is only visible within this compilation unit
    // Each .go file gets its own copy (if needed)
    // No symbol conflicts during linking
}
```

### External Linkage (Original - Problematic)

```c
void myFunction() {
    // This function is visible globally
    // If defined in multiple files -> duplicate symbol error
}
```

## Benefits of Static Functions

1. ✅ **No duplicate symbols** - Each compilation unit has its own copy
2. ✅ **Better optimization** - Compiler can inline static functions
3. ✅ **Encapsulation** - Functions are private to the file
4. ✅ **Smaller binary** - Unused static functions can be eliminated

## Verification

After the fix, the code should compile successfully on macOS:

```bash
go build ./window
# Should complete without duplicate symbol errors
```

## Best Practices for CGO

### DO:
- ✅ Use `static` for C functions that are only called from Go
- ✅ Use `static inline` for small helper functions
- ✅ Keep C code minimal in Go files
- ✅ Use header guards if splitting C code into separate files

### DON'T:
- ❌ Define C functions without `static` in CGO blocks
- ❌ Duplicate C function definitions across multiple Go files
- ❌ Use global C variables without `static`

## Alternative Solutions

If you need to share C functions across multiple Go files:

### Option 1: Header File
```c
// window_darwin.h
#ifndef WINDOW_DARWIN_H
#define WINDOW_DARWIN_H

void sharedFunction();

#endif
```

### Option 2: Single CGO File
Move all C code to one Go file and use Go functions to expose functionality.

### Option 3: Separate C Files
Use `#cgo CFLAGS` to include separate `.c` files:
```go
/*
#cgo CFLAGS: -I./include
#include "window_impl.c"
*/
import "C"
```

## Related Issues

This is a common CGO issue when:
- Multiple Go files in the same package use CGO
- C functions are defined inline in `import "C"` blocks
- Functions have external linkage (no `static` keyword)

## References

- [CGO Documentation](https://pkg.go.dev/cmd/cgo)
- [C Static Functions](https://en.cppreference.com/w/c/language/static)
- [Go Issue #13467](https://github.com/golang/go/issues/13467) - CGO and duplicate symbols


---

## Update: Refactored to Separate CGO Layer (Latest Fix)

The code has been further refactored to completely separate CGO code from the high-level Go implementation.

### New File Structure

#### 1. `window_cgo_darwin.go` (CGO Layer)
- Build tag: `// +build darwin,cgo`
- Contains ALL Objective-C code and CGO bindings
- Provides clean Go wrapper functions:
  - `darwin_createWindow()` - Returns `*darwinWindowHandle`
  - `darwin_destroyWindow()`
  - `darwin_setTitle()`
  - `darwin_runMainLoop()`
  - `darwin_stopMainLoop()`
  - `darwin_getWindowSize()`
  - `darwin_setFullscreen()`
  - `darwin_drawBitmap()`

#### 2. `window_darwin.go` (High-Level Go)
- Build tag: `// +build darwin`
- Pure Go code - NO CGO, NO Objective-C
- Implements Window and Display types
- Calls functions from `window_cgo_darwin.go`

### Benefits of This Approach

1. **Zero Duplicate Symbols** - Objective-C classes exist in only one file
2. **Clean Separation** - CGO complexity isolated from business logic
3. **Easier Testing** - Can mock CGO layer for tests
4. **Better Maintainability** - Clear boundary between native and Go code
5. **Compile Safety** - Build tags ensure correct file selection

### Migration Summary

**Before:**
- `window_darwin.go` - Mixed CGO and Go code with Objective-C classes
- `window_darwin_cgo.go` - Duplicate CGO implementation

**After:**
- `window_cgo_darwin.go` - All CGO code (renamed from `window_darwin_cgo.go`)
- `window_darwin.go` - Pure Go code calling CGO functions

This architecture follows Go best practices for CGO integration and eliminates all duplicate symbol issues.
