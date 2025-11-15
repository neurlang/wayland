package stable

import (
	"github.com/neurlang/wayland/wl"

	dmabuf "github.com/neurlang/wayland/stable/linux-dmabuf-v1"
)

// GetNewFunc returns a constructor function for stable protocol interfaces
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
