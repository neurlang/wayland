# Linux DMA-BUF v1 Implementation Notes

## Overview

The linux-dmabuf-v1 protocol has been successfully added to the `stable/` directory and follows the same patterns as unstable protocols.

## Key Implementation Details

### 1. Package Structure

```
stable/
├── stable.go                          # GetNewFunc for dynamic loading
├── README.md                          # Documentation
├── IMPLEMENTATION_NOTES.md           # This file
└── linux-dmabuf-v1/
    ├── doc.go                        # Package doc + go:generate
    ├── types.go                      # Type definitions and wrappers
    ├── example_test.go               # Usage examples
    ├── linux-dmabuf-v1.xml           # Protocol specification
    └── linux-dmabuf-v1.xml.go        # Generated bindings
```

### 2. Type Wrappers

The `types.go` file provides necessary type aliases and wrappers:

- **BaseProxy, Context, Event, Surface**: Direct aliases to `wl` package types
- **Buffer**: Custom wrapper struct that mimics `wl.Buffer` behavior

#### Why Buffer is a Wrapper

The `Buffer` type cannot be a simple alias because:

1. The generated code calls `initBuffer()` which is a private method in `wl.Buffer`
2. The generated code creates Buffer instances inline: `new(Buffer)`
3. We need to maintain the same interface as `wl.Buffer` for compatibility

The wrapper provides:
- Same initialization behavior (`initBuffer()`)
- Same methods (`Destroy()`, `AddReleaseHandler()`, etc.)
- Compatible event dispatching
- Proper handler management

### 3. GetNewFunc Integration

The `stable.GetNewFunc()` function in `stable/stable.go` provides dynamic protocol loading:

```go
func GetNewFunc(iface string) func(*wl.Context) wl.Proxy {
    switch iface {
    case "zwp_linux_dmabuf_v1":
        return func(ctx *wl.Context) wl.Proxy {
            return dmabuf.NewZwpDmabufV1(ctx)
        }
    default:
        return nil
    }
}
```

This matches the pattern used in `unstable/unstable.go`.

### 4. Usage Patterns

Two methods are supported (identical to unstable protocols):

**Method 1: Direct Import**
```go
import dmabuf "github.com/neurlang/wayland/stable/linux-dmabuf-v1"

dmabufMgr := dmabuf.NewZwpDmabufV1(ctx)
```

**Method 2: Dynamic Loading**
```go
import "github.com/neurlang/wayland/stable"

newFunc := stable.GetNewFunc("zwp_linux_dmabuf_v1")
dmabufMgr := newFunc(ctx).(*dmabuf.ZwpDmabufV1)
```

### 5. Protocol Features

The linux-dmabuf-v1 protocol provides:

- **ZwpDmabufV1**: Main factory interface for creating DMA-BUF based buffers
- **ZwpBufferParamsV1**: Parameter object for buffer creation
- **ZwpDmabufFeedbackV1**: Feedback mechanism for format/modifier support

Key methods:
- `CreateParams()`: Create buffer parameters
- `GetDefaultFeedback()`: Get default format feedback
- `GetSurfaceFeedback()`: Get surface-specific feedback
- `CreateImmed()`: Immediately create a wl_buffer from DMA-BUF

## Testing

The implementation has been verified to:
- ✅ Compile without errors
- ✅ Pass all diagnostics checks
- ✅ Follow the same patterns as unstable protocols
- ✅ Provide proper type safety
- ✅ Support both direct and dynamic loading

## Future Additions

When adding new stable protocols:

1. Create subdirectory under `stable/`
2. Add XML specification
3. Create `doc.go` with `//go:generate` directive
4. Create `types.go` with necessary wrappers
5. Run `go generate` to create bindings
6. Add protocol to `stable.GetNewFunc()`
7. Create example usage in `example_test.go`
8. Update `stable/README.md`

## Differences from wl.Buffer

The `linux.Buffer` wrapper differs from `wl.Buffer` in:
- Package location (stable/linux-dmabuf-v1 vs wl)
- Created by DMA-BUF protocol, not wl_shm
- Same interface and behavior otherwise

This allows DMA-BUF buffers to be used interchangeably with regular wl_buffers in surface operations.
