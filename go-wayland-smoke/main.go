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
import wl "github.com/neurlang/wayland/wl"
import window "github.com/neurlang/wayland/window"

import "fmt"

type smoke struct {
	display       *window.Display
	window        *window.Window
	widget        *window.Widget
	width, height int32
	current       int
	bb            [2]struct {
		d []float32
		u []float32
		v []float32
	}
}

func diffuse(smoke *smoke, time uint32, source []float32, dest []float32) {
	var s, d []float32
	var x, y, k, stride int32
	var t float32
	var a float32 = 0.0002

	stride = smoke.width

	for k = 0; k < 5; k++ {
		for y = 1; y < smoke.height-1; y++ {
			s = source[y*stride:]
			d = dest[y*stride-stride:]
			for x = 1; x < smoke.width-1; x++ {
				t = d[x-1+stride] + d[x+1+stride] +
					d[x] + d[x+stride+stride]
				d[x+stride] = (s[x] + a*t) / (1 + 4*a) * 0.995
			}
		}
	}
}

func advect(smoke *smoke, time uint32, uu []float32, vv []float32, source []float32, dest []float32) {

	var s, d []float32
	var u, v []float32
	var x, y, stride int
	var i, j int
	var px, py, fx, fy float32

	stride = int(smoke.width)

	for y = 1; y < int(smoke.height)-1; y++ {
		d = dest[y*stride:]
		u = uu[y*stride:]
		v = vv[y*stride:]

		for x = 1; x < int(smoke.width)-1; x++ {
			px = float32(x) - u[x]
			py = float32(y) - v[x]
			if px < 0.5 {
				px = 0.5
			}
			if py < 0.5 {
				py = 0.5
			}
			if px > float32(smoke.width)-1.5 {
				px = float32(smoke.width) - 1.5
			}
			if py > float32(smoke.height)-1.5 {
				py = float32(smoke.height) - 1.5
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

func project(smoke *smoke, time uint32, u []float32, v []float32, p []float32, div []float32) {
	var x, y, k, l, s int
	var h = 1.0 / float32(smoke.width)

	s = int(smoke.width)
	for i := 0; i < int(smoke.height*smoke.width); i++ {
		p[i] = 0.
	}
	for y = 1; y < int(smoke.height)-1; y++ {
		l = y * s
		for x = 1; x < int(smoke.width)-1; x++ {
			div[l+x] = -0.5 * h * (u[l+x+1] - u[l+x-1] +
				v[l+x+s] - v[l+x-s])
			p[l+x] = 0
		}
	}

	for k = 0; k < 5; k++ {
		for y = 1; y < int(smoke.height)-1; y++ {
			l = y * s
			for x = 1; x < int(smoke.width)-1; x++ {
				p[l+x] = (div[l+x] +
					p[l+x-1] +
					p[l+x+1] +
					p[l+x-s] +
					p[l+x+s]) / 4
			}
		}
	}

	for y = 1; y < int(smoke.height)-1; y++ {
		l = y * s
		for x = 1; x < int(smoke.width)-1; x++ {
			u[l+x] -= 0.5 * (p[l+x+1] - p[l+x-1]) / h
			v[l+x] -= 0.5 * (p[l+x+s] - p[l+x-s]) / h
		}
	}
}

func (*smoke) Resize(widget *window.Widget, width int32, height int32) {
}

func render(smoke *smoke, surface cairo.Surface) {
	var x, y int
	var width, height, stride int
	var c, a uint32
	var dst8 []byte

	dst8 = surface.ImageSurfaceGetData()
	width = surface.ImageSurfaceGetWidth()
	height = surface.ImageSurfaceGetHeight()
	stride = surface.ImageSurfaceGetStride()

	data := smoke.bb[int(smoke.current)].d

	for y = 1; y < int(height)-1; y++ {
		for x = 1; x < int(width)-1; x++ {

			c = uint32(data[x+y*int(smoke.height)] * 800.)
			if c > 255 {
				c = 255
			}
			a = c
			if a < 0x33 {
				a = 0x33
			}
			if dst8 != nil {
				dst8[4*x+y*int(stride)+0] = byte(c)
				dst8[4*x+y*int(stride)+1] = byte(c)
				dst8[4*x+y*int(stride)+2] = byte(c)
				dst8[4*x+y*int(stride)+3] = byte(a)
			}
		}
	}
}

func (smoke *smoke) Redraw(widget *window.Widget) {

	var time = (uint32)(smoke.widget.WidgetGetLastTime())

	diffuse(smoke, time/30, smoke.bb[0].u, smoke.bb[1].u)
	diffuse(smoke, time/30, smoke.bb[0].v, smoke.bb[1].v)
	project(smoke, time/30, smoke.bb[1].u, smoke.bb[1].v, smoke.bb[0].u, smoke.bb[0].v)
	advect(smoke, time/30, smoke.bb[1].u, smoke.bb[1].v, smoke.bb[1].u, smoke.bb[0].u)
	advect(smoke, time/30, smoke.bb[1].u, smoke.bb[1].v, smoke.bb[1].v, smoke.bb[0].v)
	project(smoke, time/30, smoke.bb[0].u, smoke.bb[0].v, smoke.bb[1].u, smoke.bb[1].v)
	diffuse(smoke, time/30, smoke.bb[0].d, smoke.bb[1].d)
	advect(smoke, time/30, smoke.bb[0].u, smoke.bb[0].v, smoke.bb[1].d, smoke.bb[0].d)

	var surface = smoke.window.WindowGetSurface()

	if surface != nil {

		render(smoke, surface)

	}

	surface.Destroy()

	smoke.widget.WidgetScheduleRedraw()
}
func smoke_motion_handler(smoke *smoke, x float32, y float32) {
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
			smoke.bb[0].u[k] += float32(256 - (int(rand.Int()) & 512))
			smoke.bb[0].v[k] += float32(256 - (int(rand.Int()) & 512))
			smoke.bb[0].d[k] += float32(1)
		}
	}
}

func (*smoke) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
}
func (*smoke) Leave(widget *window.Widget, input *window.Input) {
}
func (s *smoke) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	smoke_motion_handler(s, x, y)

	return window.CURSOR_HAND1
}
func (*smoke) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
}
func (*smoke) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (*smoke) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (s *smoke) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {

	smoke_motion_handler(s, x, y)

}
func (*smoke) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (*smoke) TouchCancel(widget *window.Widget, width int32, height int32) {
}
func (*smoke) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value wl.Fixed) {
}
func (*smoke) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (*smoke) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (*smoke) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (*smoke) PointerFrame(widget *window.Widget, input *window.Input) {
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
	smoke.window.SetBufferType(window.WINDOW_BUFFER_TYPE_SHM)
	rand.Seed(int64(time.Now().Nanosecond()))

	smoke.current = 0
	var size = int(smoke.height * smoke.width)

	smoke.bb[0].d = make([]float32, size)
	smoke.bb[0].u = make([]float32, size)
	smoke.bb[0].v = make([]float32, size)
	smoke.bb[1].d = make([]float32, size)
	smoke.bb[1].u = make([]float32, size)
	smoke.bb[1].v = make([]float32, size)

	smoke.widget.Userdata = &smoke

	smoke.widget.ScheduleResize(smoke.width, smoke.height)

	window.DisplayRun(d)

	smoke.widget.Destroy()
	smoke.window.Destroy()
	d.Destroy()

}
