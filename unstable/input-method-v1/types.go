package input

import "github.com/neurlang/wayland/wl"

type BaseProxy = wl.BaseProxy
type Context = wl.Context
type Event = wl.Event
type Seat = wl.Seat
type Surface = wl.Surface
type Keyboard = wl.Keyboard
type Output = wl.Output

var NewKeyboard = wl.NewKeyboard

func SafeCast[T any](p wl.Proxy) T {
	return wl.SafeCast[T](p)
}
