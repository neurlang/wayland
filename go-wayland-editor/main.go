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

// Go Wayland Editor demo
package main

import cairo "github.com/neurlang/wayland/cairoshim"
import "github.com/neurlang/wayland/wl"
import "github.com/neurlang/wayland/window"
import zwptiv3 "github.com/neurlang/wayland/unstable/text-input-v3"

import "fmt"

type editor struct {
	display          *window.Display
	window           *window.Window
	widget           *window.Widget
	width            int32
	height           int32
	textInputManager *zwptiv3.ZwpTextInputManagerV3

	seat *wl.Seat

	entry       *textEntry
	editor      *textEntry
	activeEntry *textEntry

	selection    *wl.DataSource
	selectedText string
}

type textEntry struct {
	window *window.Window
	widget *window.Widget

	text           string
	active         int
	panelVisible   bool
	cursorPosition int
	anchorPosition int

	textInput *zwptiv3.ZwpTextInputV3

	clickToShow       bool
	preferredLanguage string
	contentPurpose    int
}

func textEntryCreate(editor *editor, text string) *textEntry {
	var entry *textEntry

	entry = new(textEntry)

	entry.widget = editor.widget.AddWidget(entry)
	entry.window = editor.window
	entry.text = text
	entry.active = 0
	entry.panelVisible = false
	entry.cursorPosition = len(text)

	editor.display.SetSeatHandler(editor)

	//editor.widget.Userdata = entry

	//zwp_text_input_v1_add_listener(entry.text_input,
	//			       &text_input_listener, entry);

	return entry
}

func (*textEntry) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value float32) {
}
func (*textEntry) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (*textEntry) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (*textEntry) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (*textEntry) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
}
func (*textEntry) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (*textEntry) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (s *textEntry) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {

}
func (*textEntry) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (*textEntry) TouchCancel(widget *window.Widget, width int32, height int32) {
}

func (*textEntry) Resize(widget *window.Widget, width int32, height int32, totalwidth int32, totalheight int32) {
}
func (*textEntry) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
}
func (*textEntry) Leave(widget *window.Widget, input *window.Input) {
}
func (s *textEntry) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	return window.CursorIbeam
}
func (*textEntry) PointerFrame(widget *window.Widget, input *window.Input) {
}
func (*textEntry) Redraw(widget *window.Widget) {
}

func render(editor *editor, surface cairo.Surface) {

	dst8 := surface.ImageSurfaceGetData()
	width := surface.ImageSurfaceGetWidth()
	height := surface.ImageSurfaceGetHeight()
	stride := surface.ImageSurfaceGetStride()

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

			b := y * 255 / height
			g := (height - y) * 255 / height
			r := x * 255 / width
			a := (width - x) * 255 / width

			if dst8 != nil {
				dst8[4*x+y*stride+0] = byte(b)
				dst8[4*x+y*stride+1] = byte(g)
				dst8[4*x+y*stride+2] = byte(r)
				dst8[4*x+y*stride+3] = byte(a)
			}
		}
	}
}

func (editor *editor) Redraw(widget *window.Widget) {

	var surface = editor.window.WindowGetSurface()

	if surface != nil {

		render(editor, surface)
		surface.Destroy()
	}

	editor.widget.ScheduleRedraw()
}

func (*editor) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
}
func (*editor) Leave(widget *window.Widget, input *window.Input) {
}
func (editor *editor) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	return window.CursorHand1
}
func (*editor) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
}
func (*editor) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (*editor) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (editor *editor) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {

}
func (*editor) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (*editor) TouchCancel(widget *window.Widget, width int32, height int32) {
}
func (*editor) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value float32) {
}
func (*editor) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (*editor) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (*editor) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (*editor) PointerFrame(widget *window.Widget, input *window.Input) {
}

func textEntryAllocate(textEntry *textEntry, x int32, y int32,
	width int32, height int32) {
	textEntry.widget.SetAllocation(x, y, width, height)
}

func (editor *editor) Resize(Widget *window.Widget, width int32, height int32, totalwidth int32, totalheight int32) {

	var allocation = editor.widget.GetAllocation()

	textEntryAllocate(editor.entry,
		allocation.X+20, allocation.Y+20,
		width-40, height/2-40)
	textEntryAllocate(editor.editor,
		allocation.X+20, allocation.Y+height/2+20,
		width-40, height/2-40)

}
func (editor *editor) HandleGlobal(display *window.Display, name uint32,
	iface string, version uint32, data interface{}) {

	if iface == "zwp_text_input_manager_v3" {
		if tim, ok := display.BindUnstableInterface(name, iface, 1).(*zwptiv3.ZwpTextInputManagerV3); ok {

			editor.textInputManager = tim
		}
	}

}

func (editor *editor) Key(window *window.Window, input *window.Input, time uint32, key uint32, notUnicode uint32, state wl.KeyboardKeyState, data window.WidgetHandler) {
}

func (editor *editor) Focus(window *window.Window, device *window.Input) {
	editor.window.ScheduleRedraw()
}

func (editor *editor) Capabilities(input *window.Input, seat *wl.Seat, caps uint32) {

}

func (editor *editor) Name(input *window.Input, wlSeat *wl.Seat, name string) {
	if name == "default" {
		editor.seat = wlSeat
	}
	if editor.textInputManager != nil {

		input, err := editor.textInputManager.GetTextInput(editor.seat)
		if err != nil {
			fmt.Println(err)
			return
		}
		editor.entry.textInput = input
		editor.editor.textInput = input
	}
}

func textEntryDestroy(entry *textEntry) {
	if entry.widget != nil {
		entry.widget.Destroy()
	}
	if entry.textInput != nil {
		entry.textInput.Destroy()
	}
	entry.text = ""
	entry.preferredLanguage = ""
}

func main() {
	var editor editor

	optClickToShow := true
	optPreferredLanguage := "en_US"

	d, err := window.DisplayCreate([]string{})
	if err != nil {
		fmt.Println(err)
		return
	}

	editor.width = 500
	editor.height = 400
	editor.display = d

	d.SetUserData(&editor)
	d.SetGlobalHandler(&editor)
	editor.window = window.Create(d)
	editor.widget = editor.window.FrameCreate(&editor)

	editor.entry = textEntryCreate(&editor, "Entry")

	editor.entry.clickToShow = optClickToShow
	if optPreferredLanguage != "" {
		editor.entry.preferredLanguage = optPreferredLanguage
	}
	editor.editor = textEntryCreate(&editor, "Numeric")
	editor.editor.contentPurpose = zwptiv3.ZwpTextInputV3ContentPurposeNumber
	editor.editor.clickToShow = optClickToShow
	editor.selection = nil
	editor.selectedText = ""

	editor.window.SetTitle("Text Editor")
	editor.window.SetBufferType(window.BufferTypeShm)
	editor.window.SetKeyboardHandler(&editor)
	editor.window.Userdata = &editor

	editor.widget.Userdata = &editor

	editor.window.ScheduleResize(editor.width, editor.height)

	window.DisplayRun(d)

	if editor.selectedText != "" {
		editor.selectedText = ""
	}
	if editor.selection != nil {
		editor.selection.Destroy()
	}
	textEntryDestroy(editor.entry)
	editor.entry = nil
	textEntryDestroy(editor.editor)
	editor.editor = nil
	editor.activeEntry = nil
	editor.widget.Destroy()
	editor.window.Destroy()
	d.Destroy()
}
