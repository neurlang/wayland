package window

import "github.com/neurlang/wayland/wl"

type KeyboardHandler interface {
	Key(
		window *Window,
		input *Input,
		time uint32,
		key uint32,
		notUnicode uint32,
		state wl.KeyboardKeyState,
		data WidgetHandler,
	)
	Focus(window *Window, device *Input)
}
