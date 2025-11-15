package text

import "github.com/neurlang/wayland/wl"

type BaseProxy = wl.BaseProxy
type Context = wl.Context
type Event = wl.Event
type Surface = wl.Surface
type Seat = wl.Seat

func SafeCast[T any](p wl.Proxy) T {
	return wl.SafeCast[T](p)
}
