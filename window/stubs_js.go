//go:build js && wasm
// +build js,wasm

package window

import (
	cairo "github.com/neurlang/wayland/cairoshim"
)

func (w *Widget) ImageSurfaceGetData() []byte { return nil }
func (w *Widget) ImageSurfaceGetWidth() int   { return 0 }
func (w *Widget) ImageSurfaceGetHeight() int  { return 0 }
func (w *Widget) ImageSurfaceGetStride() int  { return 0 }
func (w *Widget) Reference() cairo.Surface    { return nil }
func (w *Widget) SetDestructor(f func())      {}
func (w *Widget) SetUserData(f func())        {}

func (w *Window) SeMaximized(maximized bool) error  { return nil }
func (w *Window) SetMaximized(maximized bool) error { return nil }
