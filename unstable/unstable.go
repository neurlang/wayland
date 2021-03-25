package unstable

import wl "github.com/neurlang/wayland/wl"
import ti_v3 "github.com/neurlang/wayland/unstable/text-input-v3"

func GetNewFunc(iface string) func(*wl.Context) wl.Proxy {
	switch iface {
	case "zwp_text_input_manager_v3":
		return func(ctx *wl.Context) wl.Proxy {
			return ti_v3.NewZwpTextInputManagerV3(ctx)
		}
	// TODO: add more
	default:
		return nil
	}
}
