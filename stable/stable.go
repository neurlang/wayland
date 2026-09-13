package stable

import (
	"github.com/neurlang/wayland/wl"

	dmabuf "github.com/neurlang/wayland/stable/linux-dmabuf-v1"
	"github.com/neurlang/wayland/stable/viewporter"
)

// GetNewFunc returns a constructor function for stable protocol interfaces
func GetNewFunc(iface string) func(*wl.Context) wl.Proxy {
	switch iface {
	case "zwp_linux_dmabuf_v1":
		return func(ctx *wl.Context) wl.Proxy {
			return dmabuf.NewZwpDmabufV1(ctx)
		}
	case "wp_viewporter":
		return func(ctx *wl.Context) wl.Proxy {
			return viewporter.NewWpViewporter(ctx)
		}
	default:
		return nil
	}
}
