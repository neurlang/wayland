package unstable

import "github.com/neurlang/wayland/wl"
import tiv3 "github.com/neurlang/wayland/unstable/text-input-v3"
import imv1 "github.com/neurlang/wayland/unstable/input-method-v1"
import xdgd1 "github.com/neurlang/wayland/unstable/xdg-decoration-v1"

func GetNewFunc(iface string) func(*wl.Context) wl.Proxy {
	switch iface {
	case "zwp_text_input_manager_v3":
		return func(ctx *wl.Context) wl.Proxy {
			return tiv3.NewZwpTextInputManagerV3(ctx)
		}
	case "zwp_input_method_v1":
		return func(ctx *wl.Context) wl.Proxy {
			return imv1.NewZwpInputMethodV1(ctx)
		}
	case "zxdg_decoration_manager_v1":
		return func(ctx *wl.Context) wl.Proxy {
			return xdgd1.NewZxdgDecorationManagerV1(ctx)
		}
	// TODO: add more
	default:
		return nil
	}
}
