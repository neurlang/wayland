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
import "os/exec"
import "strings"
import "bytes"
import "os"
import "io"

const navigateUp = 1
const navigateDown = 2
const navigateLeft = 4
const navigateRight = 8

type textarea struct {
	display *window.Display
	window  *window.Window
	widget  *window.Widget
	src     *window.DataSource
	width   int32
	height  int32
	StringGrid
	controls     StringGrid
	scrolls      [1]Scrollbar
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

func (s *surface) PutRGB(position ObjectPosition, texture_rgb [][3]byte, texture_width int, Bg, Fg [3]byte, flip bool) {
	dst8 := s.ImageSurfaceGetData()
	width := s.ImageSurfaceGetWidth()
	height := s.ImageSurfaceGetHeight()
	stride := s.ImageSurfaceGetStride()
	var texture_height = len(texture_rgb) / texture_width

	var pos = position

	if pos.X < 0 {
		pos.X = width + pos.X
	}
	if pos.Y < 0 {
		pos.Y = height + pos.Y
	}

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

	content, err := load_content(ContentRequest{
		Xpos:   t.StringGrid.FilePosition.X,
		Ypos:   t.StringGrid.FilePosition.Y,
		Width:  xcells,
		Height: ycells,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	t.StringGrid.XCells = xcells
	t.StringGrid.YCells = ycells
	t.StringGrid.DoLineNumbers()
	var scroll = t.scrolls[0].SyncWith(&t.StringGrid)
	t.handleContent(content, scroll)
}

func render(textarea *textarea, s cairo.Surface, time uint32) {

	textarea.mutex.RLock()
	defer textarea.mutex.RUnlock()

	textarea.StringGrid.Render(&surface{s, time})
	textarea.controls.Render(&surface{s, time})

	for i := range textarea.scrolls {
		textarea.scrolls[i].Render(&surface{s, time})
	}
}

func (textarea *textarea) Redraw(_ *window.Widget) {

	var lastTime = textarea.widget.WidgetGetLastTime()

	var surface = textarea.window.WindowGetSurface()

	if surface != nil {

		render(textarea, surface, lastTime)
		surface.Destroy()
	}

	textarea.widget.ScheduleRedraw()
}

func (s *textarea) Enter(_ *window.Widget, _ *window.Input, x float32, y float32) {

	println("enter")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.StringGrid.Selecting = false
	s.StringGrid.Motion(ObjectPosition{int((x + float32(s.StringGrid.CellWidth)*0.5) / float32(s.StringGrid.CellWidth)), int(y / float32(s.StringGrid.CellHeight))})

}
func (s *textarea) Leave(_ *window.Widget, _ *window.Input) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.StringGrid.Selecting = false

}
func (s *textarea) Motion(_ *window.Widget, _ *window.Input, time uint32, x float32, y float32) int {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.scrolls {
		if s.scrolls[i].IsHover(x, y, s.width, s.height) {
			return window.CursorHand1
		}
	}

	if s.controls.IsHover(x, y, s.width, s.height) {
		s.controls.Motion(ObjectPosition{1 + int((x-float32(s.width)-float32(s.controls.Pos.X))/float32(s.controls.CellWidth)), int(y / float32(s.controls.CellHeight))})
		return window.CursorLeftPtr
	} else {
		s.controls.Hover.X = 0
	}

	s.StringGrid.Motion(ObjectPosition{int((x + float32(s.StringGrid.CellWidth)*0.5) / float32(s.StringGrid.CellWidth)), int(y / float32(s.StringGrid.CellHeight))})

	if s.StringGrid.IsHover(x, y, s.width, s.height) {
		return window.CursorIbeam
	}

	return window.CursorHand1
}

var lastCommand *exec.Cmd

func (s *textarea) Button(_ *window.Widget, _ *window.Input, time uint32, button uint32, state wl.PointerButtonState, _ window.WidgetHandler) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if button == 272 {
		if state == wl.PointerButtonStatePressed {
			for i := range s.scrolls {
				if s.scrolls[i].IsHoverButton() {
					s.scrolls[i].Scroll()
					s.scrolls[i].SyncTo(&s.StringGrid)
					content, err := load_content(ContentRequest{
						Xpos:   s.StringGrid.FilePosition.X,
						Ypos:   s.StringGrid.FilePosition.Y,
						Width:  s.StringGrid.XCells,
						Height: s.StringGrid.YCells,
					})
					if err != nil {
						fmt.Println(err)
						return
					}
					s.StringGrid.DoLineNumbers()
					s.handleContent(content, false)
					return
				}
			}
		}
		if s.controls.IsHoverButton() {
			if state == wl.PointerButtonStateReleased {

				switch s.controls.Hover.X {
				case 0:
					println("no hover")
				case 1, 2, 3, 4:

					file := s.openedFile
					go s.saveToFileOperation(file)

				case 5, 6, 7, 8:
					var sep [8]byte
					rand.Read(sep[:])
					separator := fmt.Sprintf("%x", sep)
					cmd := exec.Command("zenity", "--multiple", "--file-selection", "--separator", separator)
					go func(cmd *exec.Cmd) {
						s.mutex.Lock()
						if lastCommand != nil {
							lastCommand.Process.Kill()
						}
						lastCommand = cmd
						s.mutex.Unlock()

						var outb, errb bytes.Buffer
						cmd.Stdout = &outb
						cmd.Stderr = &errb
						err := cmd.Run()
						if err == nil {
							files := strings.Split(outb.String(), separator)
							for _, srcPath := range files {
								if srcPath[len(srcPath)-1] == '\n' {
									srcPath = srcPath[0 : len(srcPath)-1]
								}
								s.mutex.Lock()
								s.openedFile = srcPath
								s.mutex.Unlock()
								// Open the source file
								sourceFile, err := os.Open(srcPath)
								if err != nil {
									println(err.Error())
									continue
								}
								dest := &Paste{Textarea: s, all: true}
								_, err = io.Copy(dest, sourceFile)
								if err != nil {
									println(err.Error())
								}
								sourceFile.Close()
								dest.Close()
								break // singletab load
							}
						}
					}(cmd)
				case 9, 10, 11, 12:
					println("new file")
				case 13, 14, 15, 16:
					s.window.SetMinimized()
				case 17, 18, 19, 20:
					s.window.ToggleMaximized()
				case 21, 22, 23, 24:
					if lastCommand != nil {
						lastCommand.Process.Kill()
					}
					s.display.Exit()
				}

			}
			s.controls.IbeamCursor.Y = 10000
		}
		if s.StringGrid.Button(state == wl.PointerButtonStateReleased) {
			content, err := load_content(ContentRequest{
				Xpos:   s.StringGrid.FilePosition.X,
				Ypos:   s.StringGrid.FilePosition.Y,
				Width:  s.StringGrid.XCells,
				Height: s.StringGrid.YCells,
				MultiClick: &MultiClickRequest{
					Double: s.StringGrid.WasDoubleClick(),
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			s.StringGrid.DoLineNumbers()
			s.handleContent(content, false)
			return
		}
	} else {

	}
}
func (*textarea) TouchUp(_ *window.Widget, _ *window.Input, serial uint32, time uint32, id int32) {
}
func (*textarea) TouchDown(_ *window.Widget, _ *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
	println(x, y)
}
func (s *textarea) TouchMotion(_ *window.Widget, _ *window.Input, time uint32, id int32, x float32, y float32) {
	println(x, y)
}
func (*textarea) TouchFrame(_ *window.Widget, _ *window.Input) {
}
func (*textarea) TouchCancel(_ *window.Widget, _ int32, height int32) {
}
func (*textarea) Axis(_ *window.Widget, _ *window.Input, time uint32, axis uint32, value float32) {
	println("axis", axis, value)
}
func (*textarea) AxisSource(_ *window.Widget, _ *window.Input, source uint32) {
	println("axis source", source)
}
func (*textarea) AxisStop(_ *window.Widget, _ *window.Input, time uint32, axis uint32) {
	println("axis stop", axis)
}
func (t *textarea) AxisDiscrete(_ *window.Widget, _ *window.Input, axis uint32, discrete int32) {
	t.axisDiscrete(4 * discrete)
}

func (t *textarea) axisDiscrete(discrete int32) {

	if (t.StringGrid.FilePosition.Y + int(discrete)) < 0 {
		discrete = -int32(t.StringGrid.FilePosition.Y)
		if discrete == 0 {
			return
		}
	}

	t.StringGrid.FilePosition.Y += int(discrete)
	t.StringGrid.IbeamCursor.Y -= int(discrete)
	t.StringGrid.SelectionCursor.Y -= int(discrete)
	t.StringGrid.DoLineNumbers()
	t.StringGrid.ReMotion()
	t.scrolls[0].SyncWith(&t.StringGrid)

	content, err := load_content(ContentRequest{
		Xpos:   t.StringGrid.FilePosition.X,
		Ypos:   t.StringGrid.FilePosition.Y,
		Width:  t.StringGrid.XCells,
		Height: t.StringGrid.YCells,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	t.handleContent(content, false)

	t.StringGrid.ReMotion()

}
func (*textarea) PointerFrame(_ *window.Widget, _ *window.Input) {
}
func (textarea *textarea) Key(
	win *window.Window,
	input *window.Input,
	time uint32,
	key uint32,
	notUnicode uint32,
	state wl.KeyboardKeyState,
	_ window.WidgetHandler,
) {

	win.UninhibitRedraw()

	textarea.mutex.Lock()
	defer textarea.mutex.Unlock()

	var control_key = ((state == wl.KeyboardKeyStatePressed) &&
		(notUnicode == xkb.KeyControlL ||
			notUnicode == xkb.KeyControlR ||
			(input.GetModifiers()&window.ModControlMask) != 0))

	println("Control:", control_key)

	textarea.StringGrid.Control(control_key)

	var entered = input.GetRune(&notUnicode, key)

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
		case 'c', 'v', 'x', 'a':

			if input.GetModifiers() == window.ModControlMask {

				if notUnicode == 'a' {
					println("CTRL A")

					textarea.selectAll()
					break
				}

				if notUnicode == 'c' || notUnicode == 'x' {
					println("CTRL C/X")

					if !textarea.StringGrid.IsSelection() {
						break

					}

					textarea.copyOperation(input, notUnicode == 'x')

				}
				if notUnicode == 'v' {
					println("CTRL V")
					err := input.ReceiveSelectionData("text/plain;charset=utf-8", &Paste{Textarea: textarea})
					if err != nil {
						fmt.Println(err)
					}
				}

				break
			}

			fallthrough

		default:
			println(string(input.GetRune(&notUnicode, key)))
			textarea.KeyReloadNoMutex(string(input.GetRune(&notUnicode, key)), 0, time)
		}

	} else if state == wl.KeyboardKeyStatePressed {

		switch notUnicode {
		case 65366:
			KeyRepeatSubscribe(textarea, "", notUnicode, time)
		case 65365:
			KeyRepeatSubscribe(textarea, "", notUnicode, time)
		case 65505:
			fallthrough
		case 65506:
			println("start")
			if !textarea.StringGrid.IsSelection() {
				textarea.StringGrid.SelectionCursor = textarea.StringGrid.IbeamCursor
				textarea.StringGrid.SelectionCursorAbs = textarea.StringGrid.IbeamCursorAbs
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

		fmt.Println("input.GetRune(&notUnicode, key)=", string(input.GetRune(&notUnicode, key)),
			"input.GetUtf8()=", string(input.GetUtf8()),
			"key=", key, "notUnicode=", notUnicode)
	}
}
func (textarea *textarea) Focus(window *window.Window, device *window.Input) {

	if device == nil {

		textarea.StringGrid.Control(false)

		window.InhibitRedraw()
	} else {
		window.UninhibitRedraw()
	}

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
	case 65366:
		textarea.axisDiscrete(int32(textarea.StringGrid.YCells))
		return true
	case 65365:
		textarea.axisDiscrete(-int32(textarea.StringGrid.YCells))
		return true
	case xkb.KeyHome:
		textarea.StringGrid.IbeamCursorAbs.X = 0
		textarea.StringGrid.IbeamCursor.X = 0
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyEnd:
		if len(textarea.StringGrid.LineLens) > textarea.StringGrid.IbeamCursor.Y {
			textarea.StringGrid.IbeamCursorAbs.X = textarea.StringGrid.LineLens[textarea.StringGrid.IbeamCursor.Y]
			textarea.StringGrid.IbeamCursor.X = textarea.StringGrid.IbeamCursorAbs.X - textarea.StringGrid.FilePosition.X
		}
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyDown:
		textarea.navigateHeld |= navigateDown
		if textarea.StringGrid.LineCount > textarea.StringGrid.IbeamCursorAbs.Y+1 {
			textarea.StringGrid.IbeamCursor.Y++
			textarea.StringGrid.IbeamCursorAbs.Y++
		}
		var l = textarea.StringGrid.IbeamCursorAbs.X - textarea.StringGrid.LineLen(textarea.StringGrid.IbeamCursor.Y)
		if 0 < l {
			textarea.StringGrid.IbeamCursor.X -= l
			textarea.StringGrid.IbeamCursorAbs.X -= l
		}

		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyUp:
		textarea.navigateHeld |= navigateUp
		if textarea.StringGrid.IbeamCursorAbs.Y > 0 {
			textarea.StringGrid.IbeamCursor.Y--
			textarea.StringGrid.IbeamCursorAbs.Y--
		}
		var l = textarea.StringGrid.IbeamCursorAbs.X - textarea.StringGrid.LineLen(textarea.StringGrid.IbeamCursor.Y)
		if 0 < l {
			textarea.StringGrid.IbeamCursor.X -= l
			textarea.StringGrid.IbeamCursorAbs.X -= l
		}

		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyLeft:
		textarea.navigateHeld |= navigateLeft
		if textarea.StringGrid.IbeamCursorAbs.X > 0 {
			textarea.StringGrid.IbeamCursor.X--
			textarea.StringGrid.IbeamCursorAbs.X--
			if "\t" == textarea.StringGrid.GetContent(textarea.StringGrid.IbeamCursor.X,
				textarea.StringGrid.IbeamCursor.Y) {
				for (len(textarea.StringGrid.GetContent(textarea.StringGrid.IbeamCursor.X-1,
					textarea.StringGrid.IbeamCursor.Y)) == 0) &&
					(textarea.StringGrid.IbeamCursorAbs.X > 0) {
					textarea.StringGrid.IbeamCursor.X--
					textarea.StringGrid.IbeamCursorAbs.X--
				}
			}
		}
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	case xkb.KeyRight:
		textarea.navigateHeld |= navigateRight
		for textarea.StringGrid.IbeamCursorAbs.X < textarea.StringGrid.LineLen(textarea.StringGrid.IbeamCursor.Y) {
			textarea.StringGrid.IbeamCursor.X++
			textarea.StringGrid.IbeamCursorAbs.X++
			if len(textarea.StringGrid.GetContent(textarea.StringGrid.IbeamCursor.X-1,
				textarea.StringGrid.IbeamCursor.Y)) != 0 {
				break
			}
		}
		textarea.StringGrid.IbeamCursorBlinkFix = (time)
		return true

	default:

		/*

			fmt.Println("input.GetRune(&notUnicode, key)=", string(input.GetRune(&notUnicode, key)),
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
func (textarea *textarea) expandContent(content [][]string, width int) {
	textarea.StringGrid.Content = nil
	for i := range content {
		textarea.StringGrid.Content = append(textarea.StringGrid.Content, content[i]...)
		if width > len(content[i]) {
			textarea.StringGrid.Content = append(textarea.StringGrid.Content, make([]string, width-len(content[i]))...)
		}
	}
}
func (textarea *textarea) handleContent(content *ContentResponse, needScrollBar bool) {

	textarea.expandContent(content.Content, textarea.StringGrid.XCells-content.Xpos)
	textarea.StringGrid.ContentFgColor = make(map[[2]int][3]byte)
	textarea.StringGrid.LineLens = content.LineLens
	textarea.StringGrid.EndLineLen = content.EndLineLen

	if textarea.StringGrid.LineCount != content.LineCount {
		textarea.StringGrid.LineCount = content.LineCount
		textarea.StringGrid.DoLineNumbers()
	}

	for _, v := range content.FgColor {
		textarea.StringGrid.ContentFgColor[[2]int{-textarea.StringGrid.FilePosition.X + v[0], -textarea.StringGrid.FilePosition.Y + v[1]}] = [3]byte{byte(v[2]), byte(v[3]), byte(v[4])}
	}

	if content.Erase != nil && content.Erase.Erased {
		needScrollBar = true

		if textarea.StringGrid.SelectionCursor.Less(&textarea.StringGrid.IbeamCursor) {

			textarea.StringGrid.IbeamCursor = textarea.StringGrid.SelectionCursor
			textarea.StringGrid.IbeamCursorAbs = textarea.StringGrid.SelectionCursorAbs
		} else {

			textarea.StringGrid.SelectionCursor = textarea.StringGrid.IbeamCursor
			textarea.StringGrid.SelectionCursorAbs = textarea.StringGrid.IbeamCursorAbs
		}
		textarea.StringGrid.Selecting = false
	}
	if content.Write != nil {
		needScrollBar = true
		textarea.StringGrid.IbeamCursor.X += content.Write.MoveX
		textarea.StringGrid.IbeamCursor.Y += content.Write.MoveY
		textarea.StringGrid.IbeamCursorAbs.X += content.Write.MoveX
		textarea.StringGrid.IbeamCursorAbs.Y += content.Write.MoveY

		textarea.StringGrid.SelectionCursor = textarea.StringGrid.IbeamCursor
		textarea.StringGrid.SelectionCursorAbs = textarea.StringGrid.IbeamCursorAbs
		textarea.StringGrid.Selecting = false
	}
	if content.Paste != nil {
		needScrollBar = true
	}
	if needScrollBar {
		ScrollbarSync(&(textarea.scrolls[0]), []patchScrollbar{{scrollTestFilename, ObjectPosition{0, 0}}}, content.LineCount)
	}
}

func (textarea *textarea) KeyReloadNoMutex(key string, notUnicode, time uint32) {

	if key == "" {
		textarea.KeyNavigate(key, notUnicode, time)
	} else {

		var erase = &EraseRequest{
			X0: textarea.StringGrid.IbeamCursorAbsolute().X,     /*+ textarea.StringGrid.FilePosition.X*/
			Y0: textarea.StringGrid.IbeamCursorAbsolute().Y,     /*+ textarea.StringGrid.FilePosition.Y*/
			X1: textarea.StringGrid.SelectionCursorAbsolute().X, /*+ textarea.StringGrid.FilePosition.X*/
			Y1: textarea.StringGrid.SelectionCursorAbsolute().Y, /*+ textarea.StringGrid.FilePosition.Y*/
		}
		var write = &WriteRequest{
			X:      textarea.StringGrid.IbeamCursorAbsolute().X, /*+ textarea.StringGrid.FilePosition.X*/
			Y:      textarea.StringGrid.IbeamCursorAbsolute().Y, /*+ textarea.StringGrid.FilePosition.Y*/
			Key:    key,
			Insert: true,
		}
		if key == "Delete" || key == "Backspace" {
			if !textarea.StringGrid.IsSelection() {
				erase = nil
			}
		} else {
			if !(textarea.StringGrid.IsSelection() && textarea.StringGrid.IsSelectionStrict()) {
				erase = nil
			} else {
				var writeErase = &WriteRequest{
					X:      (&textarea.StringGrid).IbeamCursorAbsolute().Lesser(textarea.StringGrid.SelectionCursorAbsolute()).X, /*+ textarea.StringGrid.FilePosition.X*/
					Y:      (&textarea.StringGrid).IbeamCursorAbsolute().Lesser(textarea.StringGrid.SelectionCursorAbsolute()).Y, /*+ textarea.StringGrid.FilePosition.Y*/
					Key:    key,
					Insert: true,
				}
				write = writeErase
			}
		}

		content, err := load_content(ContentRequest{
			Xpos:   textarea.StringGrid.FilePosition.X, /*+ textarea.StringGrid.FilePosition.X*/
			Ypos:   textarea.StringGrid.FilePosition.Y, /*+ textarea.StringGrid.FilePosition.Y*/
			Width:  textarea.StringGrid.XCells,
			Height: textarea.StringGrid.YCells,
			Erase:  erase,
			Write:  write,
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		textarea.handleContent(content, false)
	}
}

func (textarea *textarea) HandleDataSourceSend(ev wl.DataSourceSendEvent) {
	println("HandleDataSourceSend", ev.MimeType, ev.Fd, ev.FdError)
	var c Copy

	c.Textarea = textarea

	c.Receive(ev.Fd, ev.MimeType)

}
func (textarea *textarea) HandleDataSourceAction(_ wl.DataSourceActionEvent) {
	println("HandleDataSourceActions")
}
func (textarea *textarea) HandleDataSourceTarget(_ wl.DataSourceTargetEvent) {
	println("HandleDataSourceTarget")
}
func (textarea *textarea) HandleDataSourceCancelled(_ wl.DataSourceCancelledEvent) {

	println("HandleDataSourceCancelled")

	textarea.src = nil

}
func (textarea *textarea) HandleDataSourceDndDropPerformed(_ wl.DataSourceDndDropPerformedEvent) {
	println("HandleDataSourceDndDropPerformed")
}
func (textarea *textarea) HandleDataSourceDndFinished(_ wl.DataSourceDndFinishedEvent) {
	println("HandleDataSourceDndFinished")
}

func (textarea *textarea) Fullscreen(w *window.Window, _ window.WidgetHandler) {

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
	ScrollbarSync(&(textarea.scrolls[0]), []patchScrollbar{{scrollTestFilename, ObjectPosition{0, 0}}}, content.LineCount)

	textarea.controls.Font = &ControlFont
	textarea.controls.Content = []string{" ", "Ctr5l", "Ctr5r", " ", " ", "Ctr3l", "Ctr3r", " ", " ", "Ctr4l", "Ctr4r", " ", " ", "Ctr0l", "Ctr0r", " ", " ", "Ctr1l", "Ctr1r", " ", " ", "Ctr2l", "Ctr2r", " "}
	textarea.controls.XCells = len(textarea.controls.Content)
	textarea.controls.YCells = 1
	textarea.controls.CellWidth = 12
	textarea.controls.CellHeight = 24
	textarea.controls.IbeamCursor.Y = 10000
	textarea.controls.LineNumbers = 0
	textarea.controls.LineCount = 1
	textarea.controls.LineLens = []int{10000}
	textarea.controls.LastColHint = 10000
	textarea.controls.Pos = ObjectPosition{-48 * 6, 0}
	textarea.controls.BgColor = [3]byte{0, 13, 26}
	textarea.controls.FgColor = [3]byte{255, 255, 255}
	textarea.controls.Control(false)

	textarea.StringGrid.Font = &UnicodeFont
	textarea.StringGrid.XCells = 30
	textarea.StringGrid.YCells = 10
	textarea.StringGrid.CellWidth = 12
	textarea.StringGrid.CellHeight = 24
	//textarea.StringGrid.IbeamCursor.X = 1
	//textarea.StringGrid.IbeamCursor.Y = 1
	textarea.StringGrid.DoLineNumbers()
	textarea.StringGrid.LastColHint = 80
	textarea.StringGrid.FlipColor = false
	textarea.StringGrid.BgColor = [3]byte{0, 59, 112}
	textarea.StringGrid.FgColor = [3]byte{255, 255, 255}
	textarea.expandContent(content.Content, textarea.StringGrid.XCells-content.Xpos)

	const scrollWidth = 96

	textarea.scrolls[0].Pos = ObjectPosition{-scrollWidth, 24}
	textarea.scrolls[0].Width = scrollWidth
	textarea.scrolls[0].RGB = make([][3]byte, scrollWidth*2000)
	textarea.scrolls[0].BgRGB = [3]byte{0, 13, 26}
	textarea.scrolls[0].FgRGB = [3]byte{255, 255, 255}
	textarea.scrolls[0].SyncWith(&textarea.StringGrid)

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
