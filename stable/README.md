# Stable Wayland Protocols

This directory contains Go bindings for stable Wayland protocols.

## Available Protocols

### linux-dmabuf-v1

The Linux DMA-BUF protocol provides a way to create `wl_buffer` objects from Linux DMA-BUF file descriptors. This is useful for zero-copy buffer sharing between GPU and compositor.

**Interface name:** `zwp_linux_dmabuf_v1`

**Usage:**

```go
import (
    "github.com/neurlang/wayland/stable"
    dmabuf "github.com/neurlang/wayland/stable/linux-dmabuf-v1"
    "github.com/neurlang/wayland/wl"
)

// Method 1: Direct import
dmabufMgr := dmabuf.NewZwpDmabufV1(ctx)

// Method 2: Using GetNewFunc (for dynamic binding)
newFunc := stable.GetNewFunc("zwp_linux_dmabuf_v1")
if newFunc != nil {
    dmabufMgr := newFunc(ctx).(*dmabuf.ZwpDmabufV1)
}

// Bind to registry
registry.Bind(name, "zwp_linux_dmabuf_v1", version, dmabufMgr)
```

## Structure

Each protocol is organized in its own subdirectory:

```
stable/
├── stable.go                    # GetNewFunc for dynamic protocol loading
└── linux-dmabuf-v1/
    ├── doc.go                   # Package documentation and go:generate directive
    ├── types.go                 # Type aliases for Wayland core types
    ├── linux-dmabuf-v1.xml      # Protocol specification
    └── linux-dmabuf-v1.xml.go   # Generated Go bindings (auto-generated)
```

## Regenerating Bindings

To regenerate the Go bindings from XML:

```bash
cd stable/linux-dmabuf-v1
go generate
```

## Differences from Unstable Protocols

Stable protocols have reached maturity and their interfaces are guaranteed not to change in backwards-incompatible ways. The usage pattern is identical to unstable protocols:

- Both support direct import and usage
- Both support dynamic loading via `GetNewFunc`
- Both follow the same code generation and structure patterns

## Adding New Stable Protocols

1. Create a new subdirectory under `stable/`
2. Add the protocol XML file
3. Create `doc.go` with `//go:generate` directive
4. Create `types.go` with necessary type aliases
5. Run `go generate` to create bindings
6. Update `stable/stable.go` to add the new protocol to `GetNewFunc`
