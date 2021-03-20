package xdg

import "github.com/neurlang/wayland/wl"

type BaseProxy = wl.BaseProxy
type Event = wl.Event
type Context = wl.Context
type Proxy = wl.Proxy
type Surface = XdgSurface

func (s *Surface) AddListener(h XdgSurfaceConfigureHandler) {
	s.AddConfigureHandler(h)
}

type Seat = wl.Seat
type Output = wl.Output

func NewShell(ctx *Context) *Shell {
	ret := new(Shell)
	ctx.Register(ret)
	return ret
}

type Shell = XdgWmBase

func WmBaseAddListener(s *Shell, h XdgWmBasePingHandler) {
	s.AddPingHandler(h)
}

type ToplevelListener interface {
	ToplevelConfigureHandler
	ToplevelCloseHandler
}

func ToplevelAddListener(tl *Toplevel, h ToplevelListener) {
	tl.AddConfigureHandler(h)
	tl.AddCloseHandler(h)
}
