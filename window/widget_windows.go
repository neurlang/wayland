package window

import cairo "github.com/neurlang/wayland/cairoshim"

type Widget struct {
	Userdata interface{}
	//canvas                     *winc.Canvas
	buffer []byte
	drawn  map[int][4]byte
	drawn2 map[int][4]byte

	allocation_x, allocation_y int

	parent_window *Window
	handler       WidgetHandler

	allocation_width  int
	allocation_height int
	destroyed         bool
}

const BUFFER_BYTES = 4

func (w *Widget) ImageSurfaceGetData() []byte {

	if len(w.buffer) == 0 {

		w.buffer = make([]byte, BUFFER_BYTES*w.allocation_width*w.allocation_height, BUFFER_BYTES*w.allocation_width*w.allocation_height)
	}
	return w.buffer
}

func (w *Widget) ImageSurfaceGetWidth() int {
	return w.allocation_width
}

func (w *Widget) ImageSurfaceGetHeight() int {
	return w.allocation_height
}

func (w *Widget) ImageSurfaceGetStride() int {
	return w.allocation_width * BUFFER_BYTES
}

func (w *Widget) Reference() cairo.Surface {
	return w
}

func (w *Widget) SetDestructor(f func()) {
}

func (w *Widget) SetUserData(f func()) {
}

type mustResize struct {
	w      *Widget
	width  int32
	height int32
}

func (w *Widget) ScheduleResize(width int32, height int32) {
	println("ScheduleResize", width, height)

	var bx = (w.parent_window.form.Width() - w.parent_window.form.ClientWidth())
	var by = w.parent_window.form.Height() - w.parent_window.form.ClientHeight()

	// simple impl
	w.parent_window.form.SetSize(int(width)+bx, int(height)+by)
	w.parent_window.ScheduleResize(width, height)

}

func (w *Widget) Destroy() {
	w.destroyed = true
}

func (w *Widget) SetAllocation(x int, y int, pwidth int32, pheight int32) {

	w.allocation_x = int(x)
	w.allocation_y = int(y)
	w.allocation_width = int(pwidth)
	w.allocation_height = int(pheight)
	w.buffer = nil
	w.drawn = make(map[int][4]byte)
	w.drawn2 = make(map[int][4]byte)
}

func (w *Widget) WidgetGetLastTime() uint32 {
	return 0
}

func (w *Widget) ScheduleRedraw() {
	if !w.parent_window.inhibited && !w.destroyed {
		w.parent_window.form.Invalidate(false)

	}
}

type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

func (w *Widget) GetAllocation() Rectangle {
	return Rectangle{}
}
