// Copyright 2021 Neurlang project
// SPDX-License-Identifier: MIT

// Test program for GTK-style window decorations
package main

import (
	"github.com/neurlang/wayland/window"
	"log"
)

type testWindow struct{}

func (t *testWindow) Resize(w *window.Widget, width, height, pwidth, pheight int32) {}
func (t *testWindow) Redraw(w *window.Widget) {
	// Simple content drawing would go here
	// For now, just a placeholder since we're focusing on decorations
}
func (t *testWindow) Enter(w *window.Widget, i *window.Input, x, y float32)                {}
func (t *testWindow) Leave(w *window.Widget, i *window.Input)                              {}
func (t *testWindow) Motion(w *window.Widget, i *window.Input, time uint32, x, y float32) int {
	return window.CursorLeftPtr
}
func (t *testWindow) Button(w *window.Widget, i *window.Input, time, button uint32, 
	state interface{}, data window.WidgetHandler) {}
func (t *testWindow) TouchUp(w *window.Widget, i *window.Input, serial, time uint32, id int32) {}
func (t *testWindow) TouchDown(w *window.Widget, i *window.Input, serial, time uint32, 
	id int32, x, y float32) {}
func (t *testWindow) TouchMotion(w *window.Widget, i *window.Input, time uint32, 
	id int32, x, y float32) {}
func (t *testWindow) TouchFrame(w *window.Widget, i *window.Input)   {}
func (t *testWindow) TouchCancel(w *window.Widget, w2, h int32)      {}
func (t *testWindow) Axis(w *window.Widget, i *window.Input, time, axis uint32, value float32) {}
func (t *testWindow) AxisSource(w *window.Widget, i *window.Input, source uint32)              {}
func (t *testWindow) AxisStop(w *window.Widget, i *window.Input, time, axis uint32)            {}
func (t *testWindow) AxisDiscrete(w *window.Widget, i *window.Input, axis uint32, discrete int32) {}
func (t *testWindow) PointerFrame(w *window.Widget, i *window.Input) {}

func main() {
	log.Println("Starting decoration test...")
	// Implementation will be added
}
