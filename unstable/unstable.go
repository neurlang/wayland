package unstable

import (
	"github.com/neurlang/wayland/wl"
	
	fullscreen "github.com/neurlang/wayland/unstable/fullscreen-shell-v1"
	imv1 "github.com/neurlang/wayland/unstable/input-method-v1"
	tiv3 "github.com/neurlang/wayland/unstable/text-input-v3"
	xdgd1 "github.com/neurlang/wayland/unstable/xdg-decoration-v1"
)

func GetNewFunc(iface string) func(*wl.Context) wl.Proxy {
	switch iface {
	case "zwp_text_input_manager_v3":
		return func(ctx *wl.Context) wl.Proxy {
			return tiv3.NewZwpInputManagerV3(ctx)
		}
	case "zwp_input_method_v1":
		return func(ctx *wl.Context) wl.Proxy {
			return imv1.NewZwpMethodV1(ctx)
		}
	case "zxdg_decoration_manager_v1":
		return func(ctx *wl.Context) wl.Proxy {
			return xdgd1.NewZxdgDecorationManagerV1(ctx)
		}
	case "zwp_fullscreen_shell_v1":
		return func(ctx *wl.Context) wl.Proxy {
			return fullscreen.NewZwpShellV1(ctx)
		}
	default:
		return nil
	}
}
