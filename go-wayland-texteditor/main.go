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

// Go Wayland Textarea demo
package main

import "time"
import "math/rand"

import cairo "github.com/neurlang/wayland/cairoshim"
import wl "github.com/neurlang/wayland/wl"
import xkb "github.com/neurlang/wayland/xkbcommon"
import window "github.com/neurlang/wayland/window"
import "fmt"
import "sync"

const navigateUp = 1
const navigateDown = 2
const navigateLeft = 4
const navigateRight = 8

type textarea struct {
	display *window.Display
	window  *window.Window
	widget  *window.Widget
	width   int32
	height  int32
	StringGrid
	mutex        sync.RWMutex
	navigateHeld byte
	fullscreen   bool
}

type surface struct {
	cairo.Surface
	time uint32
}

func (s *surface) GetTime() uint32 {
	return s.time
}

func minColor(a, b [3]byte) (o [3]byte) {
	o = b
	if a[0] < b[0] {
		o[0] = a[0]
	}
	if a[1] < b[1] {
		o[1] = a[1]
	}
	if a[2] < b[2] {
		o[2] = a[2]
	}
	return o
}

func maxColor(a, b [3]byte) (o [3]byte) {
	o = b
	if a[0] > b[0] {
		o[0] = a[0]
	}
	if a[1] > b[1] {
		o[1] = a[1]
	}
	if a[2] > b[2] {
		o[2] = a[2]
	}
	return o
}

func (s *surface) PutRGB(pos ObjectPosition, texture_rgb [][3]byte, texture_width int, Bg, Fg [3]byte, flip bool) {
	dst8 := s.ImageSurfaceGetData()
	width := s.ImageSurfaceGetWidth()
	height := s.ImageSurfaceGetHeight()
	stride := s.ImageSurfaceGetStride()
	var texture_height = len(texture_rgb) / texture_width

	//println(pos.X, pos.Y, width, height, stride, texture_width, texture_height, texture_rgb)

	if texture_rgb == nil {
		return
	}

	var j int
	for y := pos.X; y < width; y++ {
		if j >= texture_width {
			continue
		}
		var i int
		for x := pos.Y; x < height; x++ {
			if i >= texture_height {
				continue
			}
			var dstpos = x*stride + y*4
			var srcpos = i*texture_width + j
			if flip {
				dst8[dstpos] = 255 - texture_rgb[srcpos][2]
				dst8[dstpos+1] = 255 - texture_rgb[srcpos][1]
				dst8[dstpos+2] = 255 - texture_rgb[srcpos][0]
			} else {
				dst8[dstpos] = texture_rgb[srcpos][2]
				dst8[dstpos+1] = texture_rgb[srcpos][1]
				dst8[dstpos+2] = texture_rgb[srcpos][0]
			}
			dst8[dstpos+3] = 255

			if dst8[dstpos] < Bg[2] {
				dst8[dstpos] = Bg[2]
			}
			if dst8[dstpos+1] < Bg[1] {
				dst8[dstpos+1] = Bg[1]
			}
			if dst8[dstpos+2] < Bg[0] {
				dst8[dstpos+2] = Bg[0]
			}
			if dst8[dstpos] > Fg[2] {
				dst8[dstpos] = Fg[2]
			}
			if dst8[dstpos+1] > Fg[1] {
				dst8[dstpos+1] = Fg[1]
			}
			if dst8[dstpos+2] > Fg[0] {
				dst8[dstpos+2] = Fg[0]
			}
			i++
		}
		j++
	}

}

func (t *textarea) Resize(widget *window.Widget, width int32, height int32, pwidth int32, pheight int32) {

	t.width = pwidth
	t.height = pheight
	if t.width != pwidth || t.height != pheight {
		widget.SetAllocation(0, 0, pwidth, pheight)
	}

	xcells := 1 + int(pwidth)/t.StringGrid.CellWidth
	ycells := 1 + int(pheight)/t.StringGrid.CellHeight

	content, err := load_content(ContentRequest{Width: xcells, Height: ycells})
	if err != nil {
		fmt.Println(err)
		return
	}

	t.StringGrid.XCells = xcells
	t.StringGrid.YCells = ycells
	t.handleContent(content)
}

func render(textarea *textarea, s cairo.Surface, time uint32) {
	textarea.mutex.RLock()
	defer textarea.mutex.RUnlock()

	textarea.StringGrid.Render(&surface{s, time})
}

func (textarea *textarea) Redraw(widget *window.Widget) {

	var time = (uint32)(textarea.widget.WidgetGetLastTime())

	var surface = textarea.window.WindowGetSurface()

	if surface != nil {

		render(textarea, surface, time)
		surface.Destroy()
	}

	textarea.widget.ScheduleRedraw()
}

func (s *textarea) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {

	println("enter")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.StringGrid.Selecting = false
	s.StringGrid.Motion(ObjectPosition{int((x + float32(s.StringGrid.CellWidth)*0.5) / float32(s.StringGrid.CellWidth)), int(y / float32(s.StringGrid.CellHeight))})

}
func (s *textarea) Leave(widget *window.Widget, input *window.Input) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.StringGrid.Selecting = false

}
func (s *textarea) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.StringGrid.Motion(ObjectPosition{int((x + float32(s.StringGrid.CellWidth)*0.5) / float32(s.StringGrid.CellWidth)), int(y / float32(s.StringGrid.CellHeight))})

	return window.CursorIbeam
}
func (s *textarea) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if button == 272 {
		s.StringGrid.Button(state == wl.PointerButtonStateReleased)

	} else {

	}
}
func (*textarea) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (*textarea) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
	println(x, y)
}
func (s *textarea) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {
	println(x, y)
}
func (*textarea) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (*textarea) TouchCancel(widget *window.Widget, width int32, height int32) {
}
func (*textarea) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value wl.Fixed) {
}
func (*textarea) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (*textarea) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (*textarea) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (*textarea) PointerFrame(widget *window.Widget, input *window.Input) {
}
func (textarea *textarea) Key(
	win *window.Window,
	input *window.Input,
	time uint32,
	key uint32,
	notUnicode uint32,
	state wl.KeyboardKeyState,
	data window.WidgetHandler,
) {
	textarea.mutex.Lock()
	defer textarea.mutex.Unlock()

	var entered = input.GetRune(notUnicode)

	if state == wl.KeyboardKeyStatePressed && entered != 0 {

		switch notUnicode {
		case xkb.KeyKpEnter:
			fallthrough
		case xkb.KeyReturn:
			textarea.KeyReloadNoMutex("Enter", 0, time)
			KeyRepeatSubscribe(textarea, "Enter", 0, time)
		case xkb.KeyBackspace:
			textarea.KeyReloadNoMutex("Backspace", 0, time)
			KeyRepeatSubscribe(textarea, "Backspace", 0, time)
		case xkb.KeyDelete:
			textarea.KeyReloadNoMutex("Delete", 0, time)
			KeyRepeatSubscribe(textarea, "Delete", 0, time)

		case 'c', 'v':

			if input.GetModifiers() == window.ModControlMask {

				/*
					src, err := textarea.display.CreateDataSource()
					if err != nil {
						fmt.Println(err)
					}
					_ = src
				*/
				if notUnicode == 'c' {

				}
				if notUnicode == 'v' {

					input.ReceiveSelectionData("text/plain;charset=utf-8", &Paste{Textarea: textarea})

				}

				println("CTRL C/V")
				break
			}

			fallthrough

		default:
			println(string(input.GetRune(notUnicode)))
			textarea.KeyReloadNoMutex(string(input.GetRune(notUnicode)), 0, time)
		}

	} else if state == wl.KeyboardKeyStatePressed {

		switch notUnicode {
		case 65505:
			fallthrough
		case 65506:
			println("start")
			if !(textarea.StringGrid.IsSelected || textarea.StringGrid.Selecting) {
				textarea.StringGrid.SelectionCursor = textarea.StringGrid.IbeamCursor
			}
			textarea.StringGrid.Selecting = true
		}
		KeyRepeatSubscribe(textarea, "", notUnicode, time)
		if textarea.KeyNavigate("", notUnicode, time) {
			if input.GetModifiers()&window.ModShiftMask == 0 {
				textarea.StringGrid.Selecting = false
				textarea.StringGrid.IsSelected = false
			}
		}

	} else {

		if textarea.KeyUnNavigate("", notUnicode, time) {
			println("un repeat")
			KeyRepeatSubscribe(nil, "", 0, time)
		}

		switch notUnicode {
		case 65505:
			fallthrough
		case 65506:
			println("stop")
			textarea.StringGrid.IsSelected = textarea.StringGrid.Selecting
			textarea.StringGrid.Selecting = false
		}

		fmt.Println("input.GetRune(notUnicode)=", string(input.GetRune(notUnicode)),
			"input.GetUtf8()=", string(input.GetUtf8()),
			"key=", key, "notUnicode=", notUnicode)
	}
}
func (*textarea) Focus(window *window.Window, device *window.Input) {
	print(device)
	println("Focus")

}

func (textarea *textarea) KeyUnNavigate(key string, notUnicode, time uint32) bool {
	switch notUnicode {
	case xkb.KeyDown:
		textarea.navigateHeld &= 0xf ^ navigateDown
		return textarea.navigateHeld == 0
	case xkb.KeyUp:
		textarea.navigateHeld &= 0xf ^ navigateUp
		return textarea.navigateHeld == 0
	case xkb.KeyLeft:
		textarea.navigateHeld &= 0xf ^ navigateLeft
		return textarea.navigateHeld == 0
	case xkb.KeyRight:
		textarea.navigateHeld &= 0xf ^ navigateRight
		return textarea.navigateHeld == 0
	}
	return true
}

func (textarea *textarea) KeyNavigate(key string, notUnicode, time uint32) bool {

	switch notUnicode {
	case xkb.KeyHome:
		textarea.StringGrid.IbeamCursor.X = 0
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyEnd:
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyDown:
		textarea.navigateHeld |= navigateDown
		textarea.StringGrid.IbeamCursor.Y++
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyUp:
		textarea.navigateHeld |= navigateUp
		if textarea.StringGrid.IbeamCursor.Y > 0 {
			textarea.StringGrid.IbeamCursor.Y--
		}
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyLeft:
		textarea.navigateHeld |= navigateLeft
		if textarea.StringGrid.IbeamCursor.X > 0 {
			textarea.StringGrid.IbeamCursor.X--
		}
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyRight:
		textarea.navigateHeld |= navigateRight
		textarea.StringGrid.IbeamCursor.X++
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	default:

		/*

			fmt.Println("input.GetRune(notUnicode)=", string(input.GetRune(notUnicode)),
				"input.GetUtf8()=", string(input.GetUtf8()),
				"key=", key, "notUnicode=", notUnicode)*/
	}
	return false
}

func (textarea *textarea) KeyReload(key string, notUnicode, time uint32) {
	textarea.mutex.Lock()
	defer textarea.mutex.Unlock()
	textarea.KeyReloadNoMutex(key, notUnicode, time)
}

func (textarea *textarea) handleContent(content *ContentResponse) {

	textarea.StringGrid.Content = content.Content
	textarea.StringGrid.ContentFgColor = make(map[[2]int][3]byte)
	for _, v := range content.FgColor {
		textarea.StringGrid.ContentFgColor[[2]int{v[0], v[1]}] = [3]byte{byte(v[2]), byte(v[3]), byte(v[4])}
	}

	if content.Write != nil {
		textarea.StringGrid.IbeamCursor.X += content.Write.MoveX
		textarea.StringGrid.IbeamCursor.Y += content.Write.MoveY

		textarea.StringGrid.SelectionCursor = textarea.StringGrid.IbeamCursor
		textarea.StringGrid.Selecting = false
	}
}

func (textarea *textarea) KeyReloadNoMutex(key string, notUnicode, time uint32) {

	if key == "" {
		textarea.KeyNavigate(key, notUnicode, time)
	} else {

		content, err := load_content(ContentRequest{
			Width:  textarea.StringGrid.XCells,
			Height: textarea.StringGrid.YCells,
			Write: &WriteRequest{
				X:      textarea.StringGrid.IbeamCursor.X,
				Y:      textarea.StringGrid.IbeamCursor.Y,
				Key:    key,
				Insert: true,
			}})
		if err != nil {
			fmt.Println(err)
			return
		}

		textarea.handleContent(content)
	}
}

func (textarea *textarea) Fullscreen(w *window.Window, wh window.WidgetHandler) {
	textarea.fullscreen = !textarea.fullscreen
	w.SetFullscreen(textarea.fullscreen)
}

func main() {

	var textarea textarea

	d, err := window.DisplayCreate([]string{})
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := load_content(ContentRequest{Width: 90, Height: 30})
	if err != nil {
		fmt.Println(err)
		return
	}

	textarea.StringGrid.Font = &UnicodeFont
	textarea.StringGrid.Content = content.Content
	textarea.StringGrid.XCells = 90
	textarea.StringGrid.YCells = 30
	textarea.StringGrid.CellWidth = 12
	textarea.StringGrid.CellHeight = 24
	textarea.StringGrid.IbeamCursor.X = 1
	textarea.StringGrid.IbeamCursor.Y = 1

	textarea.width = int32(textarea.StringGrid.XCells * textarea.StringGrid.CellWidth)
	textarea.height = int32(textarea.StringGrid.YCells * textarea.StringGrid.CellHeight)

	textarea.display = d
	textarea.window = window.Create(d)

	textarea.widget = textarea.window.AddWidget(&textarea)

	textarea.window.SetTitle("textarea")
	textarea.window.SetBufferType(window.BufferTypeShm)

	textarea.window.SetKeyboardHandler(&textarea)
	textarea.window.SetFullscreenHandler(&textarea)

	rand.Seed(int64(time.Now().Nanosecond()))

	textarea.widget.Userdata = &textarea

	textarea.widget.ScheduleResize(textarea.width, textarea.height)

	window.DisplayRun(d)

	textarea.widget.Destroy()
	textarea.window.Destroy()
	d.Destroy()

}
