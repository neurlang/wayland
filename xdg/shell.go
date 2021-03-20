package xdg

import "github.com/neurlang/wayland/wl"

type BaseProxy = wl.BaseProxy
type Event = wl.Event
type Context = wl.Context
type Proxy = wl.Proxy
type Surface = XdgSurface

func (s *Surface) AddListener(_ interface{}) {
	return
}



type Seat = wl.Seat
type Output = wl.Output

func NewShell(ctx *Context) *Shell {
	ret := new(Shell)
	ctx.Register(ret)
	return ret
}

type Shell = XdgWmBase


func WmBaseAddListener(s *Shell, _ interface{}) {
	return
}

func ToplevelAddListener(*Toplevel, interface{}) {
	return
}
