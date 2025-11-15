package linux_test

import (
	"log"

	"github.com/neurlang/wayland/stable"
	dmabuf "github.com/neurlang/wayland/stable/linux-dmabuf-v1"
	"github.com/neurlang/wayland/wl"
)

// Example demonstrates how to use the linux-dmabuf-v1 stable protocol
func Example() {
	// Connect to Wayland display
	display, err := wl.Connect("")
	if err != nil {
		log.Fatalf("unable to connect to wayland: %v", err)
	}
	defer display.Context().Close()

	// Get registry
	registry, err := display.GetRegistry()
	if err != nil {
		log.Fatalf("unable to get registry: %v", err)
	}

	// Method 1: Direct import and usage
	dmabufDirect := dmabuf.NewZwpDmabufV1(display.Context())
	_ = dmabufDirect

	// Method 2: Using GetNewFunc (similar to unstable protocols)
	newFunc := stable.GetNewFunc("zwp_linux_dmabuf_v1")
	if newFunc != nil {
		dmabufProxy := newFunc(display.Context())
		dmabufFromFunc := dmabufProxy.(*dmabuf.ZwpDmabufV1)
		_ = dmabufFromFunc
	}

	// Example: Bind the interface when discovered in registry
	// registry.Bind(name, "zwp_linux_dmabuf_v1", version, dmabufDirect)

	// Example: Create buffer parameters for DMA-BUF import
	// params, err := dmabufDirect.CreateParams()
	// if err != nil {
	//     log.Fatalf("unable to create params: %v", err)
	// }

	// Example: Add plane to buffer parameters
	// err = params.Add(fd, planeIdx, offset, stride, modifierHi, modifierLo)

	// Example: Create buffer from parameters
	// buffer, err := params.CreateImmed(width, height, format, flags)

	_ = registry
}
