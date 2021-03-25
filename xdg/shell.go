// Package xdg implements the stable XDG Window Manager Base protocol
package xdg

//go:generate ../../../../../bin/go-wayland-scanner -pkg xdg -i xdg-shell.xml -o xdg-shell.xml.go

import "github.com/neurlang/wayland/wl"

type BaseProxy = wl.BaseProxy
type Event = wl.Event
type Context = wl.Context
type Proxy = wl.Proxy

func (s *Surface) AddListener(h SurfaceConfigureHandler) {
	s.AddConfigureHandler(h)
}

type Seat = wl.Seat
type Output = wl.Output

func NewShell(ctx *Context) *Shell {
	ret := new(Shell)
	ctx.Register(ret)
	return ret
}

// TODO: remove
type Shell = WmBase

func WmBaseAddListener(s *Shell, h WmBasePingHandler) {
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
