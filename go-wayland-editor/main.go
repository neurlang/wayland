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
import wl "github.com/neurlang/wayland/wl"
import window "github.com/neurlang/wayland/window"
import zwp_ti_v3 "github.com/neurlang/wayland/unstable/text-input-v3"

import "fmt"

type editor struct {
	display            *window.Display
	window             *window.Window
	widget             *window.Widget
	width              int32
	height             int32
	text_input_manager *zwp_ti_v3.ZwpTextInputManagerV3

	seat *wl.Seat

	entry        *text_entry
	editor       *text_entry
	active_entry *text_entry

	selection     *wl.DataSource
	selected_text string
}

type text_entry struct {
	window *window.Window
	widget *window.Widget

	text            string
	active          int
	panel_visible   bool
	cursor_position int
	anchor_position int

	text_input *zwp_ti_v3.ZwpTextInputV3

	click_to_show      bool
	preferred_language string
	content_purpose    int
}

func text_entry_create(editor *editor, text string) *text_entry {
	var entry *text_entry

	entry = new(text_entry)

	entry.widget = editor.widget.AddWidget(entry)
	entry.window = editor.window
	entry.text = text
	entry.active = 0
	entry.panel_visible = false
	entry.cursor_position = len(text)
	entry.anchor_position = entry.anchor_position
	editor.display.SetSeatHandler(editor)

	//zwp_text_input_v1_add_listener(entry.text_input,
	//			       &text_input_listener, entry);

	return entry
}

func (*text_entry) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value wl.Fixed) {
}
func (*text_entry) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (*text_entry) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (*text_entry) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (*text_entry) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
}
func (*text_entry) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (*text_entry) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (s *text_entry) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {

}
func (*text_entry) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (*text_entry) TouchCancel(widget *window.Widget, width int32, height int32) {
}

func (*text_entry) Resize(widget *window.Widget, width int32, height int32) {
}
func (*text_entry) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
}
func (*text_entry) Leave(widget *window.Widget, input *window.Input) {
}
func (s *text_entry) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	return window.CursorIbeam
}
func (*text_entry) PointerFrame(widget *window.Widget, input *window.Input) {
}
func (*text_entry) Redraw(widget *window.Widget) {
}

func render(editor *editor, surface cairo.Surface) {

	dst8 := surface.ImageSurfaceGetData()
	width := surface.ImageSurfaceGetWidth()
	height := surface.ImageSurfaceGetHeight()
	stride := surface.ImageSurfaceGetStride()

	for y := 1; y < int(height)-1; y++ {
		for x := 1; x < int(width)-1; x++ {

			b := y * 255 / int(height)
			g := (int(height) - y) * 255 / int(height)
			r := x * 255 / int(width)
			a := (int(width) - x) * 255 / int(width)

			if dst8 != nil {
				dst8[4*x+y*int(stride)+0] = byte(b)
				dst8[4*x+y*int(stride)+1] = byte(g)
				dst8[4*x+y*int(stride)+2] = byte(r)
				dst8[4*x+y*int(stride)+3] = byte(a)
			}
		}
	}
}

func (editor *editor) Redraw(widget *window.Widget) {

	var surface = editor.window.WindowGetSurface()

	if surface != nil {

		render(editor, surface)

	}

	surface.Destroy()

	editor.widget.WidgetScheduleRedraw()
}

func (*editor) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
}
func (*editor) Leave(widget *window.Widget, input *window.Input) {
}
func (s *editor) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	return window.CursorHand1
}
func (*editor) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
}
func (*editor) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (*editor) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (s *editor) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {

}
func (*editor) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (*editor) TouchCancel(widget *window.Widget, width int32, height int32) {
}
func (*editor) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value wl.Fixed) {
}
func (*editor) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (*editor) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (*editor) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (*editor) PointerFrame(widget *window.Widget, input *window.Input) {
}
func (*editor) Resize(Widget *window.Widget, width int32, height int32) {
}
func (editor *editor) HandleGlobal(display *window.Display, name uint32,
	iface string, version uint32, data interface{}) {

	if iface == "zwp_text_input_manager_v3" {
		if tim, ok := display.BindUnstableInterface(name, iface, 1).(*zwp_ti_v3.ZwpTextInputManagerV3); ok {

			editor.text_input_manager = tim
		}
	}

}

func (editor *editor) Key(window *window.Window, input *window.Input, time uint32, key uint32, unicode uint32, state wl.KeyboardKeyState, data window.WidgetHandler) {

}

func (editor *editor) Focus(window *window.Window, device *window.Input) {

}

func (editor *editor) Capabilities(input *window.Input, seat *wl.Seat, caps uint32) {

}

func (editor *editor) Name(input *window.Input, wl_seat *wl.Seat, name string) {
	if name == "default" {
		editor.seat = wl_seat
	}
	if editor.text_input_manager != nil {

		input, err := editor.text_input_manager.GetTextInput(editor.seat)
		if err != nil {
			fmt.Println(err)
			return
		}
		editor.entry.text_input = input
		editor.editor.text_input = input
	}
}

func text_entry_destroy(entry *text_entry) {
	if entry.widget != nil {
		entry.widget.Destroy()
	}
	if entry.text_input != nil {
		entry.text_input.Destroy()
	}
	entry.text = ""
	entry.preferred_language = ""
}

func main() {
	var editor editor

	opt_click_to_show := true
	opt_preferred_language := "en_US"

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

	editor.entry = text_entry_create(&editor, "Entry")

	editor.entry.click_to_show = opt_click_to_show
	if opt_preferred_language != "" {
		editor.entry.preferred_language = opt_preferred_language
	}
	editor.editor = text_entry_create(&editor, "Numeric")
	editor.editor.content_purpose = zwp_ti_v3.ZwpTextInputV3ContentPurposeNumber
	editor.editor.click_to_show = opt_click_to_show
	editor.selection = nil
	editor.selected_text = ""

	editor.window.SetTitle("Text Editor")
	editor.window.SetBufferType(window.WindowBufferTypeShm)
	editor.window.SetKeyboardHandler(&editor)
	editor.window.Userdata = &editor

	editor.window.ScheduleResize(editor.width, editor.height)

	window.DisplayRun(d)

	if editor.selected_text != "" {
		editor.selected_text = ""
	}
	if editor.selection != nil {
		editor.selection.Destroy()
	}
	text_entry_destroy(editor.entry)
	editor.entry = nil
	text_entry_destroy(editor.editor)
	editor.editor = nil
	editor.active_entry = nil
	editor.widget.Destroy()
	editor.window.Destroy()
	d.Destroy()
}
