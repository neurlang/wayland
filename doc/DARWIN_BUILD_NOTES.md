# Darwin Build Notes

## Building on macOS

The Darwin window implementation uses CGO with Objective-C and Cocoa frameworks. It can only be built on macOS systems.

### Prerequisites

```bash
# Ensure Xcode Command Line Tools are installed
xcode-select --install

# Verify CGO is enabled
go env CGO_ENABLED  # should be "1"
```

### Building

```bash
# On macOS, build normally
go build ./window

# The build system will automatically include:
# - window_darwin.go
# - widget_darwin.go  
# - widget_handler_darwin.go
```

### Cross-Compilation Limitation

**Important**: You cannot cross-compile the Darwin window package from Linux/Windows because:

1. CGO requires the target platform's headers and libraries
2. Cocoa/AppKit frameworks are only available on macOS
3. The Objective-C runtime is macOS-specific

When building on Linux with `GOOS=darwin`, the `window_darwin.go` file will be ignored because the Cocoa headers are not available.

## File Structure

### window_darwin.go
- Contains CGO Objective-C code
- Defines Window and Display structs
- Implements Cocoa event loop
- Handles window creation and management

### widget_darwin.go
- Defines Widget struct
- Implements cairo.Surface interface
- Handles buffer management and rendering
- Provides widget lifecycle methods

### widget_handler_darwin.go
- Defines WidgetHandler interface
- Defines KeyboardHandler interface
- Defines other event handler interfaces

## Verification on Non-macOS Systems

To verify the code structure without building:

```bash
# Check syntax (will fail on imports, but shows structure)
go tool compile -I . window/window_darwin.go

# List files that would be included on Darwin
GOOS=darwin go list -f '{{.GoFiles}}' ./window

# Check for ignored files
GOOS=darwin go list -f '{{.IgnoredGoFiles}}' ./window
```

## Testing on macOS

Once on a macOS system:

```bash
# Build the window package
go build ./window

# Run tests (if any)
go test ./window

# Build example applications
go build ./go-wayland-simple-shm
go build ./go-wayland-smoke
```

## Common Issues

### Issue: window_darwin.go is ignored
**Cause**: Building on non-macOS system
**Solution**: Build on actual macOS hardware

### Issue: Cocoa/Cocoa.h not found
**Cause**: Xcode Command Line Tools not installed
**Solution**: Run `xcode-select --install`

### Issue: CGO_ENABLED=0
**Cause**: CGO is disabled
**Solution**: Set `export CGO_ENABLED=1`

### Issue: Undefined symbols
**Cause**: Missing framework linkage
**Solution**: Verify `#cgo LDFLAGS` includes `-framework Cocoa -framework QuartzCore`

## Implementation Status

✅ Display management
✅ Window creation/destruction
✅ Event handling (keyboard, mouse)
✅ Widget support
✅ Buffer rendering
✅ Fullscreen/maximize
✅ Multiple windows

## Future Enhancements

- [ ] High DPI (Retina) support
- [ ] Menu bar integration
- [ ] Native dialogs
- [ ] Drag and drop
- [ ] Clipboard integration
- [ ] Multi-monitor support
