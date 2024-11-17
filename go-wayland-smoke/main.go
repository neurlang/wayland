// Copyright 2021 Neurlang project

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Go Wayland Smoke demo
package main

import "time"
import "math/rand"

import cairo "github.com/neurlang/wayland/cairoshim"
import "github.com/neurlang/wayland/wl"
import "github.com/neurlang/wayland/window"
import xkb "github.com/neurlang/wayland/xkbcommon"
import "fmt"

type smoke struct {
	display     *window.Display
	window      *window.Window
	widget      *window.Widget
	width       int32
	height      int32
	smallwidth  int32
	smallheight int32
	o           []float32
	bb          [2]struct {
		d  []float32
		uv chan [2][]float32
		uu chan []float32
		vv chan []float32
	}
	pipe   bool
	lx, ly float32
}

func diffuse(_ *smoke, time uint32, source []float32, dest []float32, width int32, height int32) {
	var s, d []float32
	var x, y, k, stride int32
	var t float32
	var a float32 = 0.0002

	stride = width

	for k = 0; k < 5; k++ {
		for y = 1; y < height-1; y++ {
			s = source[y*stride:]
			d = dest[y*stride-stride:]
			for x = 1; x < width-1; x++ {
				t = d[x-1+stride] + d[x+1+stride] +
					d[x] + d[x+stride+stride]
				d[x+stride] = (s[x] + a*t) / (1 + 4*a) * 0.995
			}
		}
	}
}

func advect(
	_ *smoke,
	time uint32,
	uu []float32,
	vv []float32,
	source []float32,
	dest []float32,
	width int32, height int32,
) {

	var s, d []float32
	var u, v []float32
	var x, y, stride int
	var i, j int
	var px, py, fx, fy float32

	stride = int(width)

	for y = 1; y < int(height)-1; y++ {
		d = dest[y*stride:]
		u = uu[y*stride:]
		v = vv[y*stride:]

		for x = 1; x < int(width)-1; x++ {
			px = float32(x) - u[x]
			py = float32(y) - v[x]
			if px < 0.5 {
				px = 0.5
			}
			if py < 0.5 {
				py = 0.5
			}
			if px > float32(width)-1.5 {
				px = float32(width) - 1.5
			}
			if py > float32(height)-1.5 {
				py = float32(height) - 1.5
			}
			i = (int)(px)
			j = (int)(py)
			fx = px - float32(i)
			fy = py - float32(j)
			s = source[j*stride+i:]
			d[x] = (s[0]*(1-fx)+s[1]*fx)*(1-fy) +
				(s[stride]*(1-fx)+s[stride+1]*fx)*fy
		}
	}

}

func project(smoke *smoke, time uint32, u []float32, v []float32, p []float32, div []float32, width int32, height int32) {
	var x, y, k, l, s int
	var h = 1.0 / float32(width)

	s = int(width)
	for i := 0; i < int(height*width); i++ {
		p[i] = 0.
	}
	for y = 1; y < int(height)-1; y++ {
		l = y * s
		for x = 1; x < int(width)-1; x++ {
			div[l+x] = -0.5 * h * (u[l+x+1] - u[l+x-1] +
				v[l+x+s] - v[l+x-s])
			p[l+x] = 0
		}
	}

	for k = 0; k < 5; k++ {
		for y = 1; y < int(height)-1; y++ {
			l = y * s
			for x = 1; x < int(width)-1; x++ {
				p[l+x] = (div[l+x] +
					p[l+x-1] +
					p[l+x+1] +
					p[l+x-s] +
					p[l+x+s]) / 4
			}
		}
	}

	for y = 1; y < int(height)-1; y++ {
		l = y * s
		for x = 1; x < int(width)-1; x++ {
			u[l+x] -= 0.5 * (p[l+x+1] - p[l+x-1]) / h
			v[l+x] -= 0.5 * (p[l+x+s] - p[l+x-s]) / h
		}
	}
}

const maxx = 512
const maxy = 256

func (smoke *smoke) Resize(_ *window.Widget, _ int32, _ int32, width int32, height int32) {

	if smoke.smallwidth == width && smoke.smallheight == height {
		return
	}

	size := int(width) * int(height)

	var olduv = smoke.bb[0].uv

	smoke.bb[0].d = make([]float32, size)
	smoke.bb[1].d = make([]float32, size)

	if smoke.pipe {
		<-olduv
		smoke.bb[0].uu <- nil
		smoke.bb[0].vv <- nil
		<-olduv
	}
	smoke.bb[0].uv = make(chan [2][]float32, 1)
	if width > maxx {
		smoke.smallwidth = width
		smoke.width = maxx
	} else {
		smoke.smallwidth = width
		smoke.width = width
	}
	if height > maxy {
		smoke.smallheight = height
		smoke.height = maxy
	} else {
		smoke.smallheight = height
		smoke.height = height
	}
	smoke.pipeline(smoke.width, smoke.height)
	smoke.bb[0].uu <- make([]float32, size)
	smoke.bb[0].vv <- make([]float32, size)

	//smoke.widget.ScheduleResize(smoke.width, smoke.height)
	smoke.pipe = true
}

func render(smoke *smoke, surface cairo.Surface) {
	var dst8 = surface.ImageSurfaceGetData()
	var width = surface.ImageSurfaceGetWidth()
	var height = surface.ImageSurfaceGetHeight()
	var stride = surface.ImageSurfaceGetStride()

	data := smoke.o

	for y := 1; y < height-1; y++ {
		var yy = (y * int(smoke.height) / int(smoke.smallheight)) * int(smoke.width)
		for x := 1; x < width-1; x++ {
			var xx = x * int(smoke.width) / int(smoke.smallwidth)
			var c = uint32(data[xx+yy] * 800.)
			if c > 255 {
				c = 255
			}
			var a = c
			if a < 0x33 {
				a = 0x33
			}
			if dst8 != nil {
				dst8[4*x+y*stride+0] = byte(c)
				dst8[4*x+y*stride+1] = byte(c)
				dst8[4*x+y*stride+2] = byte(c)
				dst8[4*x+y*stride+3] = byte(a)
			}
		}
	}
}
func (smoke *smoke) pipeline(width int32, height int32) {
	const lastTime = 600000
	go func() {
		for {
			var u0 []float32
			u0 = <-smoke.bb[0].uu
			if u0 == nil {
				smoke.bb[1].uu <- nil
				return
			}
			u1 := make([]float32, len(u0))
			diffuse(smoke, lastTime, u0, u1, width, height)
			smoke.bb[1].uu <- u1
		}
	}()
	go func() {
		for {
			var v0 []float32
			v0 = <-smoke.bb[0].vv
			if v0 == nil {
				smoke.bb[1].vv <- nil
				return
			}
			v1 := make([]float32, len(v0))
			diffuse(smoke, lastTime, v0, v1, width, height)
			smoke.bb[1].vv <- v1
		}
	}()

	go func() {
		for {
			u1 := <-smoke.bb[1].uu
			v1 := <-smoke.bb[1].vv
			if u1 == nil || v1 == nil {
				smoke.bb[0].uv <- [2][]float32{nil, nil}
				return
			}
			u0 := make([]float32, len(u1))
			v0 := make([]float32, len(v1))

			project(smoke, lastTime, u1, v1, u0, v0, width, height)
			advect(smoke, lastTime, u1, v1, u1, u0, width, height)
			advect(smoke, lastTime, u1, v1, v1, v0, width, height)
			project(smoke, lastTime, u0, v0, u1, v1, width, height)
			smoke.bb[0].uv <- [2][]float32{u0, v0}
		}
	}()

}
func (smoke *smoke) Redraw(widget *window.Widget) {
	var lastTime = smoke.widget.WidgetGetLastTime()

	uv := <-smoke.bb[0].uv

	if len(uv[0]) == len(smoke.bb[0].d) && len(uv[1]) == len(smoke.bb[1].d) {

		diffuse(smoke, lastTime/30, smoke.bb[0].d, smoke.bb[1].d, smoke.width, smoke.height)
		advect(smoke, lastTime/30, uv[0], uv[1], smoke.bb[1].d, smoke.bb[0].d, smoke.width, smoke.height)

		smoke.o = smoke.bb[0].d

	}
	smoke.bb[0].uu <- uv[0]
	smoke.bb[0].vv <- uv[1]

	var surface = smoke.window.WindowGetSurface()

	if surface != nil {

		render(smoke, surface)
		surface.Destroy()
	}

	smoke.widget.ScheduleRedraw()
}
func smokeMotionHandler(smoke *smoke, x float32, y float32) {

	dx := smoke.lx - x
	dy := smoke.ly - y

	dt := 0.8 / (dx*dx + dy*dy + 1.0)

	if dt > 1. {
		dt = 1.
	}

	smoke.lx = x
	smoke.ly = y

	x *= float32(smoke.width) / float32(smoke.smallwidth)
	y *= float32(smoke.height) / float32(smoke.smallheight)
	var i0, i1, j0, j1 float32
	var k, i, j int
	var d float32 = 5

	if x-d < 1 {
		i0 = 1
	} else {
		i0 = x - d
	}
	if i0+2.*d > float32(smoke.width)-1. {
		i1 = float32(smoke.width) - 1.
	} else {
		i1 = i0 + 2.*d
	}
	if y-d < 1 {
		j0 = 1
	} else {
		j0 = y - d
	}
	if j0+2*d > float32(smoke.height)-1 {
		j1 = float32(smoke.height) - 1
	} else {
		j1 = j0 + 2*d
	}

	for i = int(i0); i < int(i1); i++ {
		for j = int(j0); j < int(j1); j++ {
			k = j*int(smoke.width) + i

			smoke.bb[0].d[k] += dt
		}
	}

	var freq = (smoke.smallwidth * smoke.smallheight) / (smoke.width * smoke.height)

	if rand.Int()%int(2*freq+1) == 0 {
		uv := <-smoke.bb[0].uv

		for i = int(i0); i < int(i1); i++ {
			for j = int(j0); j < int(j1); j++ {
				k = j*int(smoke.width) + i

				uv[0][k] += float32(1024 - (rand.Int() & 2047))
				uv[1][k] += float32(1024 - (rand.Int() & 2047))
			}
		}

		smoke.bb[0].uu <- uv[0]
		smoke.bb[0].vv <- uv[1]
	}
}
func (smoke *smoke) Key(
	_ *window.Window,
	_ *window.Input,
	time uint32,
	key uint32,
	notUnicode uint32,
	_ wl.KeyboardKeyState,
	_ window.WidgetHandler,
) {
	println(notUnicode)

	if notUnicode == xkb.KeyQ || notUnicode == xkb.KEYq {
		smoke.free()
		smoke.display.Exit()
	}
}
func (*smoke) Focus(_ *window.Window, _ *window.Input) {

}
func (*smoke) Enter(_ *window.Widget, _ *window.Input, x float32, y float32) {
}
func (*smoke) Leave(_ *window.Widget, _ *window.Input) {
}

func (smoke *smoke) Motion(
	_ *window.Widget,
	_ *window.Input,
	time uint32,
	x float32,
	y float32,
) int {
	smokeMotionHandler(smoke, x, y)

	return window.CursorHand1
}

func (*smoke) Button(
	_ *window.Widget,
	_ *window.Input,
	time uint32,
	button uint32,
	_ wl.PointerButtonState,
	_ window.WidgetHandler,
) {
}

func (*smoke) TouchUp(
	_ *window.Widget,
	_ *window.Input,
	serial uint32,
	time uint32,
	id int32,
) {
}

func (*smoke) TouchDown(
	_ *window.Widget,
	_ *window.Input,
	serial uint32,
	time uint32,
	id int32,
	x float32,
	y float32,
) {
}

func (smoke *smoke) TouchMotion(
	_ *window.Widget,
	_ *window.Input,
	time uint32,
	id int32,
	x float32,
	y float32,
) {

	smokeMotionHandler(smoke, x, y)

}
func (*smoke) TouchFrame(_ *window.Widget, _ *window.Input) {
}
func (*smoke) TouchCancel(_ *window.Widget, width int32, height int32) {
}

func (*smoke) Axis(
	widget *window.Widget,
	input *window.Input,
	time uint32,
	axis uint32,
	value float32,
) {
}
func (*smoke) AxisSource(_ *window.Widget, _ *window.Input, source uint32) {
}
func (*smoke) AxisStop(_ *window.Widget, _ *window.Input, time uint32, axis uint32) {
}

func (*smoke) AxisDiscrete(
	_ *window.Widget,
	_ *window.Input,
	axis uint32,
	discrete int32,
) {
}
func (*smoke) PointerFrame(_ *window.Widget, _ *window.Input) {
}

func (smoke *smoke) free() {
	// tear down the rendering pipe
	var olduv = smoke.bb[0].uv
	if smoke.pipe {
		<-olduv
		smoke.bb[0].uu <- nil
		smoke.bb[0].vv <- nil
		<-olduv
	}

	smoke.bb[0].d = nil
	smoke.bb[0].uv = nil
	smoke.bb[1].d = nil
	smoke.bb[1].uu = nil
	smoke.bb[1].vv = nil
	smoke.bb[0].uu = nil
	smoke.bb[0].vv = nil
}

func main() {

	var smoke smoke

	d, err := window.DisplayCreate([]string{})
	if err != nil {
		fmt.Println(err)
		return
	}

	smoke.width = 200
	smoke.height = 200
	smoke.display = d
	smoke.window = window.Create(d)

	smoke.widget = smoke.window.AddWidget(&smoke)

	smoke.window.SetTitle("smoke")
	smoke.window.SetBufferType(window.BufferTypeShm)
	smoke.window.SetKeyboardHandler(&smoke)
	rand.Seed(int64(time.Now().Nanosecond()))

	var size = int(smoke.height * smoke.width)

	smoke.bb[0].d = make([]float32, size)
	smoke.bb[0].uv = make(chan [2][]float32, 1)
	smoke.bb[1].d = make([]float32, size)
	smoke.bb[1].uu = make(chan []float32, 1)
	smoke.bb[1].vv = make(chan []float32, 1)
	smoke.bb[0].uu = make(chan []float32, 1)
	smoke.bb[0].vv = make(chan []float32, 1)

	//smoke.bb[0].uv <- [2][]float32{make([]float32, size), make([]float32, size)}

	smoke.widget.SetUserDataWidgetHandler(&smoke)

	smoke.widget.ScheduleResize(smoke.width, smoke.height)

	window.DisplayRun(d)

	smoke.widget.Destroy()
	smoke.window.Destroy()
	d.Destroy()

}
