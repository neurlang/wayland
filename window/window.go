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

// Package window implements a convenient wayland windowing
package window

import "syscall"

//import zwp "github.com/neurlang/wayland/wayland"
import wlclient "github.com/neurlang/wayland/wlclient"
import wlcursor "github.com/neurlang/wayland/wlcursor"
import wl "github.com/neurlang/wayland/wl"
import zxdg "github.com/neurlang/wayland/xdg"
import cairo "github.com/neurlang/wayland/cairoshim"

import "github.com/neurlang/wayland/os"

import "errors"

import "fmt"

type runner interface {
	Run(uint32)
}

const SURFACE_OPAQUE = 0x01
const SURFACE_SHM = 0x02

const SURFACE_HINT_RESIZE = 0x10
const SURFACE_HINT_RGB565 = 0x100

const WINDOW_PREFERRED_FORMAT_NONE = 0
const WINDOW_PREFERRED_FORMAT_RGB565 = 1

const WINDOW_BUFFER_TYPE_EGL_WINDOW = 0
const WINDOW_BUFFER_TYPE_SHM = 1

const CURSOR_BOTTOM_LEFT = 0
const CURSOR_BOTTOM_RIGHT = 1
const CURSOR_BOTTOM = 2
const CURSOR_DRAGGING = 3
const CURSOR_LEFT_PTR = 4
const CURSOR_LEFT = 5
const CURSOR_RIGHT = 6
const CURSOR_TOP_LEFT = 7
const CURSOR_TOP_RIGHT = 8
const CURSOR_TOP = 9
const CURSOR_IBEAM = 10
const CURSOR_HAND1 = 11
const CURSOR_WATCH = 12
const CURSOR_DND_MOVE = 13
const CURSOR_DND_COPY = 14
const CURSOR_DND_FORBIDDEN = 15
const CURSOR_BLANK = 16

const ZWP_RELATIVE_POINTER_MANAGER_V1_VERSION = 1
const ZWP_POINTER_CONSTRAINTS_V1_VERSION = 1

type global struct {
	name    uint32
	iface   string
	version uint32
}

type Display struct {
	Display              *wl.Display
	registry             *wl.Registry
	compositor           *wl.Compositor
	subcompositor        *wl.Subcompositor
	shm                  *wl.Shm
	data_device_manager  *wl.DataDeviceManager
	text_cursor_position *struct{}
	xdg_shell            *zxdg.WmBase
	serial               uint32

	//display_fd        int32
	display_fd_events uint32

	//display_task task
	//	pad4		uint64
	//	pad5		uint64
	//	pad6		uint64

	epoll_fd      int32
	deferred_list [2]uintptr
	//	pad7		uint64
	//	pad8		uint64

	running int32

	global_list []*global
	//	pad9		uint64
	//	pada		uint64
	window_list [2]*Window
	//	padb		uint64
	//	padc		uint64
	input_list []*Input
	//	padd		uint64
	//	pade		uint64
	output_list [2]*output
	//	padf		uint64
	//	padg		uint64

	theme        *theme
	cursor_theme *wlcursor.Theme
	cursors      *[lengthCursors]*wlcursor.Cursor

	xkb_context *struct{}

	/* A hack to get text extents for tooltips */
	dummy_surface *cairo.Surface

	has_rgb565                  int32
	data_device_manager_version uint32

	deferred_list_new []runner

	//display_task_new os.Runner
	surface2window map[*wl.Surface]*Window

	global_handler GlobalHandler

	user_data interface{}

	seat_handler SeatHandler
}

type rectangle struct {
	x      int32
	y      int32
	width  int32
	height int32
}

type toysurface interface {
	prepare(dx int, dy int, width int32, height int32, flags uint32,
		buffer_transform uint32, buffer_scale int32) cairo.Surface
	swap(buffer_transform uint32, buffer_scale int32, server_allocation *rectangle)
	acquire(ctx *struct{}) int
	release()
	destroy()
}

type surface struct {
	Window *Window

	surface_             *wl.Surface
	subsurface           *wl.Subsurface
	synchronized         int32
	synchronized_default int32
	toysurface           *toysurface
	Widget               *Widget
	redraw_needed        int32

	frame_cb  *wl.Callback
	last_time uint32
	//	pad1	uint32

	allocation        rectangle
	server_allocation rectangle

	input_region  *wl.Region
	opaque_region *wl.Region

	buffer_type      int32
	buffer_transform int32
	buffer_scale     int32

	cairo_surface cairo.Surface
}

func (s *surface) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	s.CallbackDone(ev.C, ev.CallbackData)
}

type Window struct {
	Display            *Display
	window_output_list [2]uintptr

	title string

	saved_allocation   rectangle
	min_allocation     rectangle
	pending_allocation rectangle
	last_geometry      rectangle

	x, y int32

	redraw_inhibited      int32
	redraw_needed         int32
	redraw_task_scheduled int32

	//redraw_task task

	//	pad1	uint64
	//	pad2	uint64

	resize_needed int32
	custom        int32
	focused       int32

	resizing int32

	fullscreen int32
	maximized  int32

	preferred_format int

	main_surface *surface
	xdg_surface  *zxdg.Surface
	xdg_toplevel *zxdg.Toplevel
	xdg_popup    *zxdg.Popup

	parent      *Window
	last_parent *Window

	/* struct surface::link, contains also main_surface */
	subsurface_list [2]*surface

	pointer_locked bool

	confined bool

	link [2]*Window

	Userdata WidgetHandler

	redraw_runner runner

	subsurface_list_new []*surface

	keyboard_handler KeyboardHandler

	frame *window_frame
}

func (Window *Window) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {
	Window.SurfaceConfigure(Window.xdg_surface, ev.Serial)
}

func (Window *Window) SurfaceConfigure(zxdg_surface_v6 *zxdg.Surface, serial uint32) {

	Window.xdg_surface.AckConfigure(serial)

	window_uninhibit_redraw(Window)

}

func (Window *Window) SetKeyboardHandler(handler KeyboardHandler) {

	Window.keyboard_handler = handler

}

func (Window *Window) HandleToplevelConfigure(ev zxdg.ToplevelConfigureEvent) {
	Window.ToplevelConfigure(Window.xdg_toplevel, ev.Width, ev.Height, ev.States)
}

func (Window *Window) ToplevelConfigure(zxdg_toplevel_v6 *zxdg.Toplevel, width int32, height int32, states []int32) {

	Window.maximized = 0
	Window.fullscreen = 0
	Window.resizing = 0
	Window.focused = 0

	for i := range states {
		switch states[i] {
		case zxdg.ToplevelStateMaximized:
			Window.maximized = 1
		case zxdg.ToplevelStateFullscreen:
			Window.fullscreen = 1
		case zxdg.ToplevelStateResizing:
			Window.resizing = 1
		case zxdg.ToplevelStateActivated:
			Window.focused = 1
		default:
			/* Unknown state */
		}
	}

	if (width > 0) && (height > 0) {
		/* The width / height params are for Window geometry,
		 * but window_schedule_resize takes allocation. Add
		 * on the shadow margin to get the difference. */
		var margin int32 = 0

		Window.ScheduleResize(width+margin*2, height+margin*2)
	} else if (Window.saved_allocation.width > 0) &&
		(Window.saved_allocation.height > 0) {
		Window.ScheduleResize(Window.saved_allocation.width, Window.saved_allocation.height)
	}

}

func (Window *Window) HandleToplevelClose(ev zxdg.ToplevelCloseEvent) {
	Window.ToplevelClose(Window.xdg_toplevel)
}

func (w *Window) ToplevelClose(zxdg_toplevel_v6 *zxdg.Toplevel) {

	w.Display.Exit()
}

func SurfaceEnter(wl_surface *wl.Surface, wl_output *wl.Output) {
}
func SurfaceLeave(wl_surface *wl.Surface, wl_output *wl.Output) {
}

type Widget struct {
	Window     *Window
	surface    *surface
	tooltip    *struct{}
	child_list *widget_list
	allocation rectangle

	opaque         int32
	tooltip_count  int32
	default_cursor int32

	/* If this is set to false then no cairo surface will be
	 * created before redrawing the surface. This is useful if the
	 * redraw handler is going to do completely custom rendering
	 * such as using EGL directly */
	use_cairo int32

	Userdata WidgetHandler
}

type widget_list struct {
	l []*Widget
}

func (l *widget_list) Add(w *Widget) {
	(l.l) = append((l.l), w)
}

func (l *widget_list) Remove(w *Widget) {
	if len(l.l) > 0 {
		if (l.l)[0] == w {
			(l.l) = (l.l)[1:]
			return
		}
		if (l.l)[len(l.l)-1] == w {
			(l.l) = (l.l)[0 : len(l.l)-1]
			return
		}
	}

	for i, v := range l.l {
		if v == w {
			(l.l) = append((l.l)[0:i], (l.l)[i+1:]...)
		}
	}
}
func (l *widget_list) Insert(w *Widget) {
	w.child_list = l
	l.Add(w)
}

type WidgetHandler interface {
	Resize(Widget *Widget, width int32, height int32)
	Redraw(Widget *Widget)
	Enter(Widget *Widget, Input *Input, x float32, y float32)
	Leave(Widget *Widget, Input *Input)
	Motion(Widget *Widget, Input *Input, time uint32, x float32, y float32) int
	Button(Widget *Widget, Input *Input, time uint32, button uint32, state wl.PointerButtonState, data WidgetHandler)
	TouchUp(Widget *Widget, Input *Input, serial uint32, time uint32, id int32)
	TouchDown(Widget *Widget, Input *Input, serial uint32, time uint32, id int32, x float32, y float32)
	TouchMotion(Widget *Widget, Input *Input, time uint32, id int32, x float32, y float32)
	TouchFrame(Widget *Widget, Input *Input)
	TouchCancel(Widget *Widget, width int32, height int32)
	Axis(Widget *Widget, Input *Input, time uint32, axis uint32, value wl.Fixed)
	AxisSource(Widget *Widget, Input *Input, source uint32)
	AxisStop(Widget *Widget, Input *Input, time uint32, axis uint32)
	AxisDiscrete(Widget *Widget, Input *Input, axis uint32, discrete int32)
	PointerFrame(Widget *Widget, Input *Input)
}

type xkb_mod_mask_t uint32

type Input struct {
	Display              *Display
	seat                 *wl.Seat
	pointer              *wl.Pointer
	keyboard             *wl.Keyboard
	touch                *wl.Touch
	touch_point_list     [2]uintptr
	pointer_focus        *Window
	keyboard_focus       *Window
	touch_focus          int32
	current_cursor       int32
	cursor_anim_start    uint32
	cursor_frame_cb      *wl.Callback
	cursor_timer_start   uint32
	cursor_anim_current  uint32
	cursor_delay_fd      int32
	cursor_timer_running bool
	//cursor_task          task
	pointer_surface      *wl.Surface
	modifiers            uint32
	pointer_enter_serial uint32
	cursor_serial        uint32
	sx                   float32
	sy                   float32

	focus_widget *Widget
	grab         *Widget
	grab_button  uint32

	data_device       *wl.DataDevice
	touch_grab        uint32
	touch_grab_id     int32
	drag_x            float32
	drag_y            float32
	drag_focus        *Window
	drag_enter_serial uint32

	xkb struct {
		control_mask xkb_mod_mask_t
		alt_mask     xkb_mod_mask_t
		shift_mask   xkb_mod_mask_t
	}

	repeat_rate_sec   int32
	repeat_rate_nsec  int32
	repeat_delay_sec  int32
	repeat_delay_nsec int32

	//repeat_task     task
	repeat_sym   uint32
	repeat_key   uint32
	repeat_time  uint32
	seat_version int32
}

func (i *Input) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	i.CallbackDone(ev.C, ev.CallbackData)
}

type KeyboardHandler interface {
	Key(window *Window, input *Input, time uint32, key uint32, unicode uint32, state wl.KeyboardKeyState, data WidgetHandler)
	Focus(window *Window, device *Input)
}

func input_remove_keyboard_focus(input *Input) {
	var window = input.keyboard_focus

	if window == nil {
		return
	}

	if window.keyboard_handler != nil {
		window.keyboard_handler.Focus(window, nil)
	}

	input.keyboard_focus = nil
}

type output struct {
	Display          *Display
	output           *wl.Output
	server_output_id uint32
	allocation       rectangle
	link             [2]*output
	transform        int32
	scale            int32
	maker            string
	model            string
}

type shm_pool struct {
	pool *wl.ShmPool
	size uintptr
	used uintptr
	data []byte
}

const CURSOR_DEFAULT = 100
const CURSOR_UNSET = 101

//line 509
func surface_to_buffer_size(buffer_transform uint32, buffer_scale int32, width *int32, height *int32) {

	switch buffer_transform {
	case wl.OutputTransform90:
		fallthrough
	case wl.OutputTransform270:
		fallthrough
	case wl.OutputTransformFlipped90:
		fallthrough
	case wl.OutputTransformFlipped270:
		*width, *height = *height, *width
	}

	*width *= buffer_scale
	*height *= buffer_scale
}

//line 532
func buffer_to_surface_size(buffer_transform uint32, buffer_scale int32, width *int32, height *int32) {
	switch buffer_transform {
	case wl.OutputTransform90:
		fallthrough
	case wl.OutputTransform270:
		fallthrough
	case wl.OutputTransformFlipped90:
		fallthrough
	case wl.OutputTransformFlipped270:
		*width, *height = *height, *width

	}

	*width /= buffer_scale
	*height /= buffer_scale
}

func (Input *Input) HandlePointerEnter(ev wl.PointerEnterEvent) {
	Input.PointerEnter(nil, ev.Serial, ev.Surface, ev.SurfaceX, ev.SurfaceY)
}

func (Input *Input) PointerEnter(wl_pointer *wl.Pointer, serial uint32, surface *wl.Surface, sx float32, sy float32) {

	if nil == surface {
		/* enter event for a Window we've just destroyed */
		return
	}

	var Window = Input.Display.surface2window[surface]

	if surface != Window.main_surface.surface_ {
		//		DBG("Ignoring Input event from subsurface %p\n", surface);
		return
	}

	Input.Display.serial = serial
	Input.pointer_enter_serial = serial
	Input.pointer_focus = Window

	Input.sx = sx
	Input.sy = sy

}

func (Input *Input) HandlePointerLeave(ev wl.PointerLeaveEvent) {
	Input.PointerLeave(nil, ev.Serial, ev.Surface)
}

func (Input *Input) PointerLeave(wl_pointer *wl.Pointer, serial uint32, wl_surface *wl.Surface) {

	Input.Display.serial = serial
	input_remove_pointer_focus(Input)

}

func (Input *Input) HandlePointerMotion(ev wl.PointerMotionEvent) {
	Input.PointerMotion(ev.P, ev.Time, ev.SurfaceX, ev.SurfaceY)
}

func (Input *Input) PointerMotion(wl_pointer *wl.Pointer, time uint32, surface_x float32, surface_y float32) {

	pointer_handle_motion(Input, wl_pointer, time, surface_x, surface_y)
}

func (Input *Input) HandlePointerButton(ev wl.PointerButtonEvent) {
	Input.PointerButton(ev.P, ev.Serial, ev.Time, ev.Button, ev.State)
}

func (input *Input) PointerButton(wl_pointer *wl.Pointer, serial uint32, time uint32, button uint32, state_w uint32) {
	var widget *Widget
	var state = wl.PointerButtonState(state_w)

	input.Display.serial = serial
	if input.focus_widget != nil && input.grab == nil &&
		state == wl.PointerButtonStatePressed {
		input_grab(input, input.focus_widget, button)
	}

	widget = input.grab
	if widget != nil && widget.Userdata != nil {
		widget.Userdata.Button(widget,
			input, time,
			button, state,
			input.grab.Userdata)
	}

	if input.grab != nil && input.grab_button == button &&
		state == wl.PointerButtonStateReleased {
		input_ungrab(input)
	}

}

func (Input *Input) HandlePointerAxis(ev wl.PointerAxisEvent) {

}

func (*Input) PointerAxis(wl_pointer *wl.Pointer, time uint32, axis uint32, value wl.Fixed) {
}

func (Input *Input) HandlePointerFrame(ev wl.PointerFrameEvent) {

}

func (*Input) PointerFrame(wl_pointer *wl.Pointer) {
}

func (Input *Input) HandlePointerAxisSource(ev wl.PointerAxisSourceEvent) {

}

func (*Input) PointerAxisSource(wl_pointer *wl.Pointer, axis_source uint32) {
}

func (Input *Input) HandlePointerAxisStop(ev wl.PointerAxisStopEvent) {

}

func (*Input) PointerAxisStop(wl_pointer *wl.Pointer, time uint32, axis uint32) {
}

func (Input *Input) HandlePointerAxisDiscrete(ev wl.PointerAxisDiscreteEvent) {

}

func (*Input) PointerAxisDiscrete(wl_pointer *wl.Pointer, axis uint32, discrete int32) {
}

type SeatHandler interface {
	Capabilities(i *Input, seat *wl.Seat, caps uint32)
	Name(i *Input, seat *wl.Seat, name string)
}

func (input_ *Input) HandleSeatCapabilities(ev wl.SeatCapabilitiesEvent) {
	input_.SeatCapabilities(input_.seat, ev.Capabilities)
}

func (input_ *Input) HandleSeatName(ev wl.SeatNameEvent) {
	input_.SeatName(input_.seat, ev.Name)
}

func (input_ *Input) SeatCapabilities(seat *wl.Seat, caps uint32) {

	if ((caps & wl.SeatCapabilityPointer) != 0) && (input_.pointer == nil) {
		var err error
		input_.pointer, err = seat.GetPointer()
		if err != nil {
			fmt.Println(err)
			return
		}
		wlclient.PointerSetUserData(input_.pointer, input_)
		wlclient.PointerAddListener(input_.pointer, input_)

	} else if ((caps & wl.SeatCapabilityPointer) == 0) && (nil != input_.pointer) {
		if input_.seat_version >= wl.POINTER_RELEASE_SINCE_VERSION {
			input_.pointer.Release()
		} else {
			wlclient.PointerDestroy(input_.pointer)
			input_.pointer = nil
		}
	}

	if input_.Display.seat_handler != nil {
		input_.Display.seat_handler.Capabilities(input_, seat, caps)
	}

}
func (i *Input) SeatName(wl_seat *wl.Seat, name string) {
	if i.Display.seat_handler != nil {
		i.Display.seat_handler.Name(i, wl_seat, name)
	}
}

// line 2697
func input_grab(input *Input, widget *Widget, button uint32) {
	input.grab = widget
	input.grab_button = button

	input_set_focus_widget(input, widget, input.sx, input.sy)
}

// line 2706
func input_ungrab(input *Input) {

	input.grab = nil
	if input.pointer_focus != nil {
		var widget = window_find_widget(input.pointer_focus,
			int32(input.sx), int32(input.sy))
		input_set_focus_widget(input, widget, (input.sx), (input.sy))
	}
}

type shm_surface_data struct {
	buffer *wl.Buffer
	pool   *shm_pool
}

//line 734
func shm_surface_data_destroy(data *shm_surface_data) {
	data.buffer.Destroy()
	if data.pool != nil {
		shm_pool_destroy(data.pool)
	}
}

//line 744
func make_shm_pool(Display *Display, size uintptr, data *[]byte) (pool *wl.ShmPool) {
	fd, err := os.CreateAnonymousFile(int64(size))
	if err != nil {
		println("creating a buffer file failed")
		return nil
	}

	*data, err = syscall.Mmap(int(fd.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		println("mmap failed")
		fd.Close()
		return nil
	}

	pool, err = Display.shm.CreatePool(uintptr(fd.Fd()), int32(size))
	if err != nil {
		println("create pool failed")
		fd.Close()
		return nil
	}

	fd.Close()

	return pool
}

//line 772
func shm_pool_create(Display *Display, size uintptr) *shm_pool {
	var pool *shm_pool = &shm_pool{}

	pool.pool = make_shm_pool(Display, size, &pool.data)

	if pool.pool == nil {
		return nil
	}

	pool.size = size
	pool.used = 0

	return pool
}

//line 792
func shm_pool_allocate(pool *shm_pool, size uintptr, offset *int) (ret []byte) {

	if pool.used+size > uintptr(pool.size) {
		return nil
	}

	*offset = int(pool.used)
	ret = pool.data[uintptr(pool.used):]
	pool.used += size
	pool.data = pool.data[0:uintptr(pool.used)]

	return ret
}

//line 804
/* destroy the pool. this does not unmap the memory though */
func shm_pool_destroy(pool *shm_pool) {

	err := syscall.Munmap(pool.data)
	if err != nil {
		println(err)
	}

	pool.pool.Destroy()
	pool.data = nil
	pool.pool = nil
	pool.size = 0
	pool.used = 0
}

//line 820
func data_length_for_shm_surface(rect *rectangle) uintptr {
	var stride = int32(cairo.FormatStrideForWidth(cairo.FORMAT_ARGB32, int(rect.width)))
	return uintptr(int(stride * rect.height))
}

func shm_pool_reset(pool *shm_pool) {
	pool.used = 0
}

//line 829
func display_create_shm_surface_from_pool(Display *Display,
	rectangle *rectangle,
	flags uint32, pool *shm_pool) (*cairo.Surface, *shm_surface_data) {
	var data *shm_surface_data = &shm_surface_data{}
	var format uint32
	var surface cairo.Surface
	var cairo_format cairo.Format
	var stride, length int
	var offset int
	var map_ []byte
	var err error

	if (flags&uint32(SURFACE_HINT_RGB565) != 0) && Display.has_rgb565 != 0 {
		cairo_format = cairo.FORMAT_RGB16_565
	} else {
		cairo_format = cairo.FORMAT_ARGB32
	}

	stride = cairo.FormatStrideForWidth(cairo_format, int(rectangle.width))

	length = stride * int(rectangle.height)
	data.pool = nil

	map_ = shm_pool_allocate(pool, uintptr(length), &offset)

	if map_ == nil {
		return nil, nil
	}

	surface = cairo.ImageSurfaceCreateForData(map_,
		cairo_format,
		int(rectangle.width),
		int(rectangle.height),
		int(stride))

	surface.SetUserData(func() {

		shm_surface_data_destroy((data))
	})

	if (flags&uint32(SURFACE_HINT_RGB565) != 0) && Display.has_rgb565 != 0 {
		format = wl.ShmFormatRgb565
	} else {
		if flags&SURFACE_OPAQUE != 0 {
			format = wl.ShmFormatXrgb8888
		} else {
			format = wl.ShmFormatArgb8888
		}
	}

	data.buffer, err = pool.pool.CreateBuffer(int32(offset),
		rectangle.width,
		rectangle.height,
		int32(stride), format)
	if err != nil {
		return nil, nil
	}

	return &surface, data
}

//line 886
func display_create_shm_surface(Display *Display,
	rectangle *rectangle, flags uint32,
	alternate_pool *shm_pool,
	data_ret **shm_surface_data) *cairo.Surface {
	var data *shm_surface_data
	var pool *shm_pool
	var surface *cairo.Surface

	if alternate_pool != nil {
		shm_pool_reset(alternate_pool)

		surface, data = display_create_shm_surface_from_pool(Display, rectangle, flags, alternate_pool)

		if surface != nil {
			goto out
		}
	}

	pool = shm_pool_create(Display, data_length_for_shm_surface(rectangle))

	if pool == nil {
		return nil
	}

	surface, data =
		display_create_shm_surface_from_pool(Display, rectangle, flags, pool)

	if surface == nil {
		shm_pool_destroy(pool)
		return nil
	}

	/* make sure we destroy the pool when the surface is destroyed */
	data.pool = pool

out:
	if data_ret != nil {
		*data_ret = data
	}

	return surface
}

type shm_surface_leaf struct {
	cairo_surface *cairo.Surface
	/* 'data' is automatically destroyed, when 'cairo_surface' is */
	data *shm_surface_data

	resize_pool *shm_pool
	busy        int32
}

func shm_surface_leaf_release(leaf *shm_surface_leaf) {
	if leaf.cairo_surface != nil {
		(*leaf.cairo_surface).Destroy()
	}
	/* leaf.data already destroyed via cairo private */
}

const MAX_LEAVES = 3

//line 983
type shm_surface struct {
	Display *Display
	surface *wl.Surface
	flags   uint32
	dx      int32
	dy      int32

	leaf    [MAX_LEAVES]shm_surface_leaf
	current *shm_surface_leaf
}

func shm_surface_buffer_release(surface *shm_surface, buffer *wl.Buffer) {
	var leaf *shm_surface_leaf
	var i int
	var free_found int

	for i = 0; i < MAX_LEAVES; i++ {
		leaf = &surface.leaf[i]
		if leaf.data != nil && leaf.data.buffer == buffer {
			leaf.busy = 0
			break
		}
	}
	if i >= MAX_LEAVES {
		panic("unknown buffer released")
	}

	/* Leave one free leaf with storage, release others */
	free_found = 0
	for i = 0; i < MAX_LEAVES; i++ {
		leaf = &surface.leaf[i]

		if (leaf.cairo_surface == nil) || (leaf.busy != 0) {
			continue
		}

		if free_found == 0 {
			free_found = 1
		} else {
			shm_surface_leaf_release(leaf)

		}
	}
}

func (s *shm_surface) HandleBufferRelease(ev wl.BufferReleaseEvent) {
	s.BufferRelease(ev.B)
}

func (s *shm_surface) BufferRelease(buf *wl.Buffer) {
	shm_surface_buffer_release(s, buf)

}

func (s *shm_surface) prepare(dx int, dy int, width int32, height int32, flags uint32,
	buffer_transform uint32, buffer_scale int32) cairo.Surface {

	var resize_hint bool = (flags & SURFACE_HINT_RESIZE) != 0
	surface := s
	var rect rectangle
	var leaf *shm_surface_leaf
	var i int

	surface.dx = int32(dx)
	surface.dy = int32(dy)

	for i = 0; i < MAX_LEAVES; i++ {
		if surface.leaf[i].busy != 0 {
			continue
		}

		if leaf == nil || surface.leaf[i].cairo_surface != nil {
			leaf = &surface.leaf[i]
		}
	}

	if nil == leaf {
		panic("all buffers are held by the server.\n")

	}

	if !resize_hint && (leaf.resize_pool != nil) {
		(*leaf.cairo_surface).Destroy()
		leaf.cairo_surface = nil
		shm_pool_destroy(leaf.resize_pool)
		leaf.resize_pool = nil
	}

	surface_to_buffer_size(buffer_transform, (buffer_scale), (&width), (&height))

	if (leaf.cairo_surface != nil) &&
		(int32((*leaf.cairo_surface).ImageSurfaceGetWidth()) == width) &&
		(int32((*leaf.cairo_surface).ImageSurfaceGetHeight()) == height) {
		goto out
	}

	if leaf.cairo_surface != nil {
		(*leaf.cairo_surface).Destroy()
	}

	rect.width = width
	rect.height = height

	leaf.cairo_surface = display_create_shm_surface(surface.Display, &rect, surface.flags, leaf.resize_pool, &leaf.data)

	if leaf.cairo_surface == nil {
		return nil
	}

	wlclient.BufferAddListener(leaf.data.buffer, surface)

out:
	surface.current = leaf

	return (*leaf.cairo_surface).Reference()
}

//line 1146
func shm_surface_swap(surface *shm_surface, buffer_transform uint32,
	buffer_scale int32, server_allocation *rectangle) {
	var leaf *shm_surface_leaf = surface.current

	server_allocation.width =
		int32((*leaf.cairo_surface).ImageSurfaceGetWidth())
	server_allocation.height =
		int32((*leaf.cairo_surface).ImageSurfaceGetHeight())

	buffer_to_surface_size(buffer_transform, (buffer_scale),
		(&server_allocation.width),
		(&server_allocation.height))

	surface.surface.Attach(leaf.data.buffer,
		int32(surface.dx), int32(surface.dy))
	surface.surface.Damage(0, 0,
		server_allocation.width, server_allocation.height)
	surface.surface.Commit()

	leaf.busy = 1
	surface.current = nil
}

func (s *shm_surface) swap(buffer_transform uint32, buffer_scale int32, server_allocation *rectangle) {
	shm_surface_swap(s, buffer_transform, buffer_scale, server_allocation)

}

func (*shm_surface) acquire(ctx *struct{}) int {
	return -1
}

func (*shm_surface) release() {
}

func shm_surface_destroy(surface *shm_surface) {
	var i int

	for i = 0; i < MAX_LEAVES; i++ {
		shm_surface_leaf_release(&surface.leaf[i])
	}
}

func (s *shm_surface) destroy() {

	shm_surface_destroy(s)
}

//line 1199
func shm_surface_create(Display *Display, wl_surface *wl.Surface,
	flags uint32, rectangle *rectangle) toysurface {
	var surface = &shm_surface{}

	surface.Display = Display
	surface.surface = wl_surface
	surface.flags = flags

	return surface
}

const lengthCursors = 16

//line 1343
func create_cursors(Display *Display) (err error) {

	//line 1323
	var Cursors = [lengthCursors][]string{
		{"bottom_left_corner\000", "sw-resize\000", "size_bdiag\000"},
		{"bottom_right_corner\000", "se-resize\000", "size_fdiag\000"},
		{"bottom_side\000", "s-resize\000", "size_ver\000"},
		{"grabbing\000", "closedhand\000", "208530c400c041818281048008011002\000"},
		{"left_ptr\000", "default\000", "top_left_arrow\000", "left-arrow\000"},
		{"left_side\000", "w-resize\000", "size_hor\000"},
		{"right_side\000", "e-resize\000", "size_hor\000"},
		{"top_left_corner\000", "nw-resize\000", "size_fdiag\000"},
		{"top_right_corner\000", "ne-resize\000", "size_bdiag\000"},
		{"top_side\000", "n-resize\000", "size_ver\000"},
		{"xterm\000", "ibeam\000", "text\000"},
		{"hand1\000", "pointer\000", "pointing_hand\000", "e29285e634086352946a0e7090d73106\000"},
		{"watch\000", "wait\000", "0426c94ea35c87780ff01dc239897213\000"},
		{"dnd-move\000"},
		{"dnd-copy\000"},
		{"dnd-none\000", "dnd-no-drop\000"},
	}

	var cursor *wlcursor.Cursor

	theme, err := wlcursor.LoadTheme(32, Display.shm)
	if err != nil {
		return err
	}
	Display.cursor_theme = theme

	var wlCursors = [lengthCursors]*wlcursor.Cursor{}

	Display.cursors = &wlCursors

	for i := range Cursors {
		for j := range Cursors[i] {

			var str = string(Cursors[i][j])

			str = str[:len(str)-1]

			cursor, err = Display.cursor_theme.GetCursor(str)
			if err != nil {
				println("could not get cursor")

			} else if cursor != nil {

				(*Display.cursors)[i] = cursor
				break
			}
		}

		if (*Display.cursors)[i] == nil {
			println("could not load cursor")
		}
	}
	return nil
}

//line 1386
func destroy_cursors(Display *Display) {
	Display.cursor_theme.Destroy()
}

//line 1402
func surface_flush(surface *surface) {
	if surface.cairo_surface == nil {
		return
	}

	if surface.opaque_region != nil {
		surface.surface_.SetOpaqueRegion(surface.opaque_region)
		surface.opaque_region.Destroy()
		surface.opaque_region = nil
	}

	if surface.input_region != nil {
		surface.surface_.SetInputRegion(surface.input_region)
		surface.input_region.Destroy()
		surface.input_region = nil
	}

	(*surface.toysurface).swap(uint32(surface.buffer_transform), surface.buffer_scale,
		&surface.server_allocation)

	surface.cairo_surface.Destroy()
	surface.cairo_surface = nil
}

//line 1462
func surface_create_surface(surface *surface, flags uint32) {
	var Display *Display = surface.Window.Display
	var allocation rectangle = surface.allocation

	if surface.toysurface == nil {
		var toy = shm_surface_create(Display, surface.surface_, flags, &allocation)

		surface.toysurface = &toy
	}

	surface.cairo_surface = (*surface.toysurface).prepare(
		0, 0,
		allocation.width, allocation.height, flags,
		uint32(surface.buffer_transform), surface.buffer_scale)

}

//line 1488
func window_create_main_surface(Window *Window) {
	var surface *surface = Window.main_surface
	var flags uint32 = 0

	if Window.resizing != 0 {
		flags |= SURFACE_HINT_RESIZE
	}

	if Window.preferred_format == WINDOW_PREFERRED_FORMAT_RGB565 {
		flags |= SURFACE_HINT_RGB565
	}

	surface_create_surface(surface, flags)

}

//line 1552
func surface_destroy(surface *surface) {
	if surface.frame_cb != nil {
		wlclient.CallbackDestroy((surface.frame_cb))
	}

	if surface.input_region != nil {
		wlclient.RegionDestroy((surface.input_region))
	}

	if surface.opaque_region != nil {
		wlclient.RegionDestroy((surface.opaque_region))
	}

	if surface.subsurface != nil {
		wlclient.SubsurfaceDestroy((surface.subsurface))
	}

	surface.surface_.Destroy()

	if surface.toysurface != nil {
		(*surface.toysurface).destroy()
	}

}

//line 1577
func (Window *Window) Destroy() {

	if Window.xdg_toplevel != nil {
		Window.xdg_toplevel.Destroy()
	}
	if Window.xdg_popup != nil {
		Window.xdg_popup.Destroy()
	}
	if Window.xdg_surface != nil {
		Window.xdg_surface.Destroy()
	}

	surface_destroy(Window.main_surface)

}

//line 1624
func widget_find_widget(Widget *Widget, x int32, y int32) *Widget {

	if Widget.allocation.x <= x &&
		x < Widget.allocation.x+Widget.allocation.width &&
		Widget.allocation.y <= y &&
		y < Widget.allocation.y+Widget.allocation.height {
		return Widget
	}

	return nil
}

//line 1645
func window_find_widget(Window *Window, x int32, y int32) *Widget {
	var surface *surface
	var Widget *Widget

	for _, surface = range Window.subsurface_list_new {

		Widget = widget_find_widget(surface.Widget, x, y)

		if Widget != nil {
			return Widget
		}
	}

	return nil
}

//line 1655
func widget_create(Window *Window, surface *surface, data WidgetHandler) *Widget {
	var w = new(Widget)
	w.Window = Window
	w.surface = surface
	w.Userdata = data
	w.allocation = surface.allocation
	w.child_list = new(widget_list)
	w.opaque = 0
	w.tooltip = nil
	w.tooltip_count = 0
	w.default_cursor = CURSOR_LEFT_PTR
	w.use_cairo = 1

	return w
}

//line 1675
func (Window *Window) AddWidget(data WidgetHandler) *Widget {
	var w = widget_create(Window, Window.main_surface, data)

	Window.main_surface.Widget = w

	return w
}

//line 1702
func (parent *Widget) AddWidget(data WidgetHandler) *Widget {
	widget := widget_create(parent.Window, parent.surface, data)

	parent.child_list.Insert(widget)

	return widget
}

//line 1701
func (Widget *Widget) Destroy() {

	var surface *surface = Widget.surface

	/* Destroy the sub-surface along with the root Widget */
	if (surface.Widget == Widget) && (surface.subsurface != nil) {
		surface_destroy(Widget.surface)
	}

}

func (d *Display) SetSeatHandler(h SeatHandler) {
	d.seat_handler = h
}

func (d *Display) HandleWmBasePing(ev zxdg.WmBasePingEvent) {
	d.ShellPing(d.xdg_shell, ev.Serial)
}

func (d *Display) ShellPing(shell *zxdg.WmBase, serial uint32) {
	shell.Pong(serial)
}

func min_u32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
func (d *Display) HandleRegistryGlobal(e wl.RegistryGlobalEvent) {
	d.RegistryGlobal(d.registry, e.Name, e.Interface, e.Version)
}
func (d *Display) HandleRegistryGlobalRemove(e wl.RegistryGlobalRemoveEvent) {
	d.RegistryGlobalRemove(d.registry, e.Name)
}
func (d *Display) RegistryGlobal(registry *wl.Registry, id uint32, iface string, version uint32) {
	var global = &global{}

	global.name = id
	global.iface = iface
	global.version = version

	d.global_list = append(d.global_list, global)

	switch iface {

	case "wl_compositor":
		d.compositor = wlclient.RegistryBindCompositorInterface(d.registry, id, 1)

	case "wl_output":

		display_add_output(d, id)
		// TODO
	case "wl_seat":

		display_add_input(d, id, int(version))

	case "wl_shm":
		d.shm = wlclient.RegistryBindShmInterface(d.registry, id, 1)
		wlclient.ShmAddListener(d.shm, d)
	case "wl_data_device_manager":
		d.data_device_manager_version = min_u32(version, 3)

		wlclient.RegistryBindDataDeviceManagerInterface(d.registry, id,
			d.data_device_manager_version)

	//case "zxdg_shell_v6":
	case "xdg_wm_base":

		d.xdg_shell = wlclient.RegistryBindWmBaseInterface(d.registry, id, 1)

		zxdg.WmBaseAddListener(d.xdg_shell, d)

	case "text_cursor_position":
	case "wl_subcompositor":

	default:

	}
	if d.global_handler != nil {
		d.global_handler.HandleGlobal(d, id, iface, version, d.user_data)
	}

}
func (d *Display) RegistryGlobalRemove(wl_registry *wl.Registry, name uint32) {

}

type GlobalHandler interface {
	HandleGlobal(d *Display, id uint32, iface string, version uint32, data interface{})
}

func (d *Display) SetGlobalHandler(gh GlobalHandler) {
	d.global_handler = gh
	if gh == nil {
		return
	}
	for _, v := range d.global_list {
		d.global_handler.HandleGlobal(d, v.name, v.iface, v.version, d.user_data)
	}
}

func (d *Display) HandleShmFormat(e wl.ShmFormatEvent) {
	d.ShmFormat(nil, e.Format)
}
func (d *Display) ShmFormat(wl_shm *wl.Shm, format uint32) {
	print("SHM FORMAT: ")
	println(format)
}

//line 1733
func widget_set_size(Widget *Widget, width int32, height int32) {
	Widget.allocation.width = width
	Widget.allocation.height = height
}

//line 1740
func widget_set_allocation(Widget *Widget, x int32, y int32, width int32, height int32) {
	Widget.allocation.x = x
	Widget.allocation.y = y
	widget_set_size(Widget, width, height)
}

// line 1763
func widget_get_cairo_surface(Widget *Widget) cairo.Surface {
	var surface *surface = Widget.surface
	var Window *Window = Widget.Window

	if Widget.use_cairo == 0 {
		panic("assert")
	}

	if nil == surface.cairo_surface {
		if surface == Window.main_surface {
			window_create_main_surface(Window)

		} else {
			surface_create_surface(surface, 0)

		}
	}

	return surface.cairo_surface
}

// line 1887
func (w *Widget) WidgetGetLastTime() uint32 {
	return w.surface.last_time
}

//line 2013
func (Widget *Widget) WidgetScheduleRedraw() {
	Widget.surface.redraw_needed = 1
	window_schedule_redraw_task(Widget.Window)
}

//line 2036
func (Window *Window) WindowGetSurface() cairo.Surface {
	var cairo_surface = widget_get_cairo_surface(Window.main_surface.Widget)

	return cairo_surface.Reference()
}

func (window *Window) FrameCreate(data WidgetHandler) *Widget {
	var buttons uint32

	if window.custom != 0 {
		buttons = FRAME_BUTTON_NONE
	} else {
		buttons = FRAME_BUTTON_ALL
	}

	var frame = new(window_frame)
	frame.frame = frame_create(window.Display.theme, 0, 0, buttons, window.title, nil)
	if frame.frame == nil {
		frame = nil
		return nil
	}

	frame.widget = window.AddWidget(frame)
	frame.child = frame.widget.AddWidget(data)

	window.frame = frame

	return frame.child
}

//line 2614
func input_set_focus_widget(Input *Input, focus *Widget,
	x float32, y float32) {
	var old, Widget *Widget
	var cursor int

	if focus == Input.focus_widget {
		return
	}

	old = Input.focus_widget
	if old != nil {
		Input.focus_widget = nil
	}

	if focus != nil {
		Widget = focus
		if Input.grab != nil {
			Widget = Input.grab
		}
		Input.focus_widget = focus
		cursor = int(Widget.default_cursor)

		input_set_pointer_image(Input, cursor)
	}
}

//line 2714
func cancel_pointer_image_update(Input *Input) {

}

// line 2718
func input_remove_pointer_focus(input_ *Input) {
	var Window = input_.pointer_focus

	if nil == Window {
		return
	}

	input_set_focus_widget(input_, nil, 0, 0)

	input_.pointer_focus = nil
	input_.current_cursor = CURSOR_UNSET

	cancel_pointer_image_update(input_)
}

// line 2776
func pointer_handle_motion(data *Input, pointer *wl.Pointer,
	time uint32, sx float32, sy float32) {
	var Input *Input = data
	_ = Input
	var Window *Window = Input.pointer_focus
	var Widget *Widget
	var cursor int

	if Window == nil {
		return
	}

	Input.sx = float32(sx)
	Input.sy = float32(sy)

	// when making the Window smaller - e.g. after an unmaximise we might
	// * still have a pending motion event that the compositor has picked
	// * based on the old surface dimensions. However, if we have an active
	// * grab, we expect to see Input from outside the Window anyway.

	if nil == Input.grab && (sx < float32(Window.main_surface.allocation.x) ||
		sy < float32(Window.main_surface.allocation.y) ||
		sx > float32(Window.main_surface.allocation.width) ||
		sy > float32(Window.main_surface.allocation.height)) {
		return
	}

	if !(Input.grab != nil && Input.grab_button != 0) {
		Widget = window_find_widget(Window, int32(sx), int32(sy))
		input_set_focus_widget(Input, Widget, sx, sy)

	}

	if Input.grab != nil {
		Widget = Input.grab
	} else {
		Widget = Input.focus_widget
	}
	if Widget != nil {
		if Widget.Userdata != nil {
			cursor = Widget.Userdata.Motion(Input.focus_widget,
				Input, time, float32(sx), float32(sy))
		} else {
			cursor = int(Widget.default_cursor)
		}
	} else {
		cursor = int(CURSOR_LEFT_PTR)
	}
	_ = cursor

	input_set_pointer_image(Input, cursor)
}

//line 3552
func input_get_seat(Input *Input) *wl.Seat {
	return Input.seat
}

//line 3754
func input_set_pointer_image_index(Input *Input, index int) {
	var buffer *wl.Buffer
	var cursor *wlcursor.Cursor
	var image wlcursor.Image

	if Input.pointer == nil {
		return
	}

	cursor = Input.Display.cursors[Input.current_cursor]
	if cursor == nil {
		return
	}

	image = cursor.GetCursorImage(index)
	if image == nil {
		print("cursor index out of range\n")
		return
	}

	buffer = image.GetBuffer()
	if buffer == nil {
		return
	}

	Input.pointer_surface.Attach(buffer, 0, 0)
	Input.pointer_surface.Damage(0, 0,
		int32(image.GetWidth()), int32(image.GetHeight()))
	Input.pointer_surface.Commit()
	wlcursor.PointerSetCursor(Input.pointer, Input.pointer_enter_serial, Input.pointer_surface,
		int32(image.GetHotspotX()), int32(image.GetHotspotY()))

}

//line 3789
func input_set_pointer_special(Input *Input) bool {
	if Input.current_cursor == CURSOR_BLANK {
		wlcursor.PointerSetCursor((Input.pointer),
			(Input.pointer_enter_serial),
			nil, 0, 0)
		return true
	}

	if Input.current_cursor == CURSOR_UNSET {
		return true
	}

	return false
}

//line 3805
func schedule_pointer_image_update(Input *Input,
	cursor *wlcursor.Cursor,
	duration uint32,
	force_frame bool) {
	/* Some silly cursor sets have enormous pauses in them.  In these
	 * cases it's better to use a timer even if it results in less
	 * accurate presentation, since it will save us having to set the
	 * same cursor image over and over again.
	 *
	 * This is really not the way we're supposed to time any kind of
	 * animation, but we're pretending it's OK here because we don't
	 * want animated cursors with long delays to needlessly hog CPU.
	 *
	 * We use force_frame to ensure we don't accumulate large timing
	 * errors by running off the wrong clock.
	 */
	if !force_frame && (duration > 100) {
		return
	}

	/* for short durations we'll just spin on frame callbacks for
	 * accurate timing - the way any kind of timing sensitive animation
	 * should really be done. */
	cb, err := Input.pointer_surface.Frame()
	if err != nil {
		fmt.Println(err)
		return
	}

	Input.cursor_frame_cb = cb

	wlclient.CallbackAddListener(Input.cursor_frame_cb, Input)

}

func (Input *Input) CallbackDone(wl_callback *wl.Callback, callback_data uint32) {
	pointer_surface_frame_callback(Input, wl_callback, callback_data)
}

//line 3842
func pointer_surface_frame_callback(Input *Input, callback *wl.Callback, time uint32) {
	var cursor *wlcursor.Cursor
	var i int
	var duration uint32
	var force_frame = true

	cancel_pointer_image_update(Input)

	if callback != nil {
		if callback != Input.cursor_frame_cb {
			panic("assert")
		}
		wlclient.CallbackDestroy(callback)
		Input.cursor_frame_cb = nil
		force_frame = false
	}

	if Input.pointer == nil {
		return
	}

	if input_set_pointer_special(Input) {
		return
	}

	cursor = Input.Display.cursors[Input.current_cursor]
	if cursor == nil {
		return
	}

	/* FIXME We don't have the current time on the first call so we set
	 * the animation start to the time of the first frame callback. */
	if time == 0 {
		Input.cursor_anim_start = 0
	} else if Input.cursor_anim_start == 0 {
		Input.cursor_anim_start = time
	}

	Input.cursor_anim_current = time

	if time == 0 || Input.cursor_anim_start == 0 {
		duration = 0
		i = 0
	} else {
		frame_duration := cursor.FrameAndDuration(time - Input.cursor_anim_start)

		i, duration = frame_duration.FrameIndex, frame_duration.FrameDuration
	}

	if cursor.ImageCount() > 1 {
		schedule_pointer_image_update(Input, cursor, duration,
			force_frame)
	}

	input_set_pointer_image_index(Input, i)
}

//line 3925
func input_set_pointer_image(Input *Input, pointer int) {
	var force bool

	if Input.pointer == nil {
		return
	}

	if Input.pointer_enter_serial > Input.cursor_serial {
		force = true
	}

	if !force && pointer == int(Input.current_cursor) {
		return
	}

	Input.current_cursor = int32(pointer)
	Input.cursor_serial = Input.pointer_enter_serial
	if Input.cursor_frame_cb == nil {
		pointer_surface_frame_callback(Input, nil, 0)
	} else if force && (!input_set_pointer_special(Input)) {
		/* The current frame callback may be stuck if, for instance,
		 * the set cursor request was processed by the server after
		 * this client lost the focus. In this case the cursor surface
		 * might not be mapped and the frame callback wouldn't ever
		 * complete. Send a set_cursor and attach to try to map the
		 * cursor surface again so that the callback will finish */

		input_set_pointer_image_index(Input, 0)
	}
}

// line 4104
func surface_resize(surface *surface) {
	var Widget *Widget = surface.Widget

	if (surface.allocation.width != Widget.allocation.width) ||
		(surface.allocation.height != Widget.allocation.height) {
		window_schedule_redraw(Widget.Window)

	}

	surface.allocation = Widget.allocation

}

//line 4144
func window_do_resize(Window *Window) {
	widget_set_allocation(Window.main_surface.Widget,
		Window.pending_allocation.x,
		Window.pending_allocation.y,
		Window.pending_allocation.width,
		Window.pending_allocation.height)

	surface_resize(Window.main_surface)

	if (Window.fullscreen != 0) && (Window.maximized != 0) {
		Window.saved_allocation = Window.pending_allocation
	}
}

//line 4191
func idle_resize(Window *Window) {
	Window.resize_needed = 0
	Window.redraw_needed = 1

	window_do_resize(Window)
}

//line 4223
func (Window *Window) ScheduleResize(width int32, height int32) {
	/* We should probably get these numbers from the theme. */
	const min_width = 200
	const min_height = 200

	Window.pending_allocation.x = 0
	Window.pending_allocation.y = 0
	Window.pending_allocation.width = width
	Window.pending_allocation.height = height

	if Window.min_allocation.width == 0 {
		if width < min_width {
			Window.min_allocation.width = min_width
		} else {
			Window.min_allocation.width = width
		}
		if height < min_height {
			Window.min_allocation.height = min_height
		} else {
			Window.min_allocation.height = height
		}
	}

	if Window.pending_allocation.width < Window.min_allocation.width {
		Window.pending_allocation.width = Window.min_allocation.width
	}
	if Window.pending_allocation.height < Window.min_allocation.height {
		Window.pending_allocation.height = Window.min_allocation.height
	}

	Window.resize_needed = 1
	window_schedule_redraw(Window)
}

//line 4254
func (Widget *Widget) ScheduleResize(width int32, height int32) {
	Widget.Window.ScheduleResize(width, height)
}

//line 4269
func window_inhibit_redraw(Window *Window) {
	Window.redraw_inhibited = 1
	Window.redraw_task_scheduled = 0
}

// line 4284
func window_uninhibit_redraw(Window *Window) {
	Window.redraw_inhibited = 0
	if (Window.redraw_needed != 0) || (Window.resize_needed != 0) {
		window_schedule_redraw_task(Window)
	}
}

//line 4521
func window_get_allocation(Window *Window, allocation *rectangle) {
	*allocation = Window.main_surface.allocation
}

//line 4445
func window_get_geometry(Window *Window, geometry *rectangle) {
	if Window.fullscreen != 0 {
		window_get_allocation(Window, geometry)
	}
}

//line 4458
func window_sync_geometry(Window *Window) {
	var geometry rectangle

	if Window.xdg_surface == nil {
		return
	}

	window_get_geometry(Window, &geometry)

	if geometry.x == Window.last_geometry.x &&
		geometry.y == Window.last_geometry.y &&
		geometry.width == Window.last_geometry.width &&
		geometry.height == Window.last_geometry.height {
		return
	}

	Window.xdg_surface.SetWindowGeometry(
		geometry.x,
		geometry.y,
		geometry.width,
		geometry.height)
	Window.last_geometry = geometry
}

// line 4480
func window_flush(Window *Window) {

	if Window.redraw_inhibited != 0 {
		panic("assert\n")
	}

	if Window.custom == 0 {
		if Window.xdg_surface != nil {
			window_sync_geometry(Window)

		}

	}

	surface_flush(Window.main_surface)

}

// line 4505
func widget_redraw(Widget *Widget) {
	if Widget.Userdata != nil {
		Widget.Userdata.Redraw(Widget)
	}
}

//line 4517
func (surface *surface) CallbackDone(callback *wl.Callback, time uint32) {
	wlclient.CallbackDestroy(callback)
	surface.frame_cb = nil

	surface.last_time = time

	if (surface.redraw_needed != 0) || (surface.Window.redraw_needed != 0) {

		window_schedule_redraw_task(surface.Window)
	}
}

//line 4545
func surface_redraw(surface *surface) int {

	if (surface.Window.redraw_needed == 0) && (surface.redraw_needed == 0) {
		return 0
	}

	// Whole-Window redraw forces a redraw even if the previous has
	// not yet hit the screen
	if nil != surface.frame_cb {
		if surface.Window.redraw_needed == 0 {
			return 0
		}

		wlclient.CallbackDestroy(surface.frame_cb)
	}

	cb, err := surface.surface_.Frame()
	if err != nil {
		fmt.Println(err)
	} else {
		surface.frame_cb = cb

		// add listener here
		wlclient.CallbackAddListener(surface.frame_cb, surface)
	}

	surface.redraw_needed = 0

	widget_redraw(surface.Widget)

	return 0
}

// This is the alternative to idle_redraw
// line 4617
func (Window *Window) Run(events uint32) {

	Window.redraw_task_scheduled = 0

	if Window.resize_needed != 0 {
		if nil != Window.main_surface.frame_cb {
			return
		}

		idle_resize(Window)

	}

	surface_redraw(Window.main_surface)

	Window.redraw_needed = 0
	window_flush(Window)

}

//line 4619

func window_schedule_redraw_task(Window *Window) {
	if Window.redraw_inhibited != 0 {
		return
	}

	if Window.redraw_task_scheduled == 0 {

		Window.redraw_runner = Window
		display_defer(Window.Display /*&Window.redraw_task,*/, Window)
		Window.redraw_task_scheduled = 1
	}
}

// line 4636
func window_schedule_redraw(Window *Window) {
	window_schedule_redraw_task(Window)
}

// line 4793
func (Window *Window) SetTitle(title string) {

	if Window.xdg_toplevel != nil {
		Window.xdg_toplevel.SetTitle(title)
	}
}

// line 5178
func surface_create(Window *Window) *surface {
	var Display *Display = Window.Display
	var surface = &surface{}
	surface.Window = Window
	surf, err := Display.compositor.CreateSurface()
	if err != nil {
		panic(err.Error())
		return nil
	}
	surface.surface_ = surf

	surface.buffer_scale = 1
	wlclient.SurfaceAddListener(surface.surface_, SurfaceEnter, SurfaceLeave)

	Window.subsurface_list_new = append(Window.subsurface_list_new, surface)

	return surface
}

// line 5219
func window_create_internal(Display *Display, custom int) *Window {

	var Window = &Window{}
	var surface_ *surface

	Window.Display = Display
	surface_ = surface_create(Window)

	Window.main_surface = (*surface)(surface_)

	if (custom > 0) || (Display.xdg_shell != nil) {
	} else {
		panic("assertion failed")
	}
	Window.custom = (int32)(custom)
	Window.preferred_format = WINDOW_PREFERRED_FORMAT_NONE

	surface_.buffer_type = WINDOW_BUFFER_TYPE_SHM

	wlclient.SurfaceSetUserData(surface_.surface_, uintptr(0))
	Display.surface2window[surface_.surface_] = Window

	return Window
}

//line 5250
func Create(Display *Display) *Window {
	var Window = window_create_internal(Display, 0)

	if Window.Display.xdg_shell != nil {
		surf, err := Window.Display.xdg_shell.GetSurface(Window.main_surface.surface_)
		if err != nil {
			fmt.Println(err)
			return nil
		} else {
			Window.xdg_surface = surf
		}

		Window.xdg_surface.AddListener(Window)

		tl, err := Window.xdg_surface.GetToplevel()
		if err != nil {
			fmt.Println(err)
			return nil
		} else {
			Window.xdg_toplevel = tl
		}

		zxdg.ToplevelAddListener(Window.xdg_toplevel, Window)

		window_inhibit_redraw(Window)

		Window.main_surface.surface_.Commit()
	}

	return Window
}

// line 5592
func (Window *Window) SetBufferType(t int32) {
	Window.main_surface.buffer_type = t
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (o *output) HandleOutputDone(ev wl.OutputDoneEvent) {
	o.OutputDone(o.output)
}
func (o *output) HandleOutputGeometry(ev wl.OutputGeometryEvent) {
	o.OutputGeometry(o.output, int(ev.X), int(ev.Y), int(ev.PhysicalWidth),
		int(ev.PhysicalHeight), int(ev.Subpixel), ev.Make, ev.Model, int(ev.Transform))
}
func (o *output) HandleOutputMode(ev wl.OutputModeEvent) {
	o.OutputMode(o.output, ev.Flags, int(ev.Width), int(ev.Height), int(ev.Refresh))
}
func (o *output) HandleOutputScale(ev wl.OutputScaleEvent) {
	o.OutputScale(o.output, ev.Factor)
}

func (o *output) OutputGeometry(wl_output *wl.Output, x int, y int, physical_width int,
	physical_height int, subpixel int, maker string, model string, transform int) {

	o.maker = maker
	o.model = model

}
func (o *output) OutputDone(wl_output *wl.Output) {

}
func (o *output) OutputScale(wl_output *wl.Output, factor int32) {

}
func (o *output) OutputMode(wl_output *wl.Output, flags uint32, width int, height int, refresh int) {

}

// line 5771
func display_add_output(d *Display, id uint32) {

	var output *output = &output{}

	output.Display = d
	output.scale = 1
	output.output = wlclient.RegistryBindOutputInterface(d.registry, id, 2)

	output.server_output_id = id

	wlclient.OutputAddListener(output.output, output)

}

//line 5925
func display_add_input(d *Display, id uint32, display_seat_version int) {

	var input_ *Input
	var seat_version = min(display_seat_version, 7)

	_ = seat_version

	input_ = new(Input)

	input_.Display = d
	input_.seat = wlclient.RegistryBindSeatInterface(d.registry, id, uint32(seat_version))
	input_.touch_focus = 0
	input_.pointer_focus = nil
	input_.keyboard_focus = nil
	input_.seat_version = int32(seat_version)

	d.input_list = append(d.input_list, input_)

	wlclient.SeatAddListener(input_.seat, input_)

	ps, err := d.compositor.CreateSurface()
	if err != nil {
		fmt.Println(err)
	} else {
		input_.pointer_surface = ps
	}

}

// line 6237
func DisplayCreate(argv []string) (d *Display, e error) {

	d = &Display{}

	d.Display, e = wlclient.DisplayConnect(nil)
	if e != nil {
		return nil, fmt.Errorf("failed to connect to Wayland Display: %w", e)
	}

	//d.display_fd = (int32)(wlclient.DisplayGetFd(d.Display))

	//d.display_task_new = d

	//display_watch_fd(d, int(d.display_fd), uint32(syscall.EPOLLIN|syscall.EPOLLERR|syscall.EPOLLHUP),
	//	os.DoFlagRunner(&d.display_task_new))

	d.surface2window = make(map[*wl.Surface]*Window)

	d.registry, e = wlclient.DisplayGetRegistry(d.Display)
	if e != nil {
		return nil, fmt.Errorf("failed to get Registry: %w", e)
	}
	wlclient.RegistryAddListener(d.registry, d)

	if wlclient.DisplayRoundtrip(d.Display) != nil {
		return nil, errors.New("failed to process Wayland connection")
	}

	create_cursors(d)

	return d, nil
}

func (Display *Display) BindUnstableInterface(name uint32, iface string, version uint32) wl.Proxy {
	return wlclient.RegistryBindUnstableInterface(Display.registry, name, iface, version)
}

func (Display *Display) SetUserData(data interface{}) {
	Display.user_data = data
}

//line 6387
func (Display *Display) Destroy() {

	if Display.dummy_surface != nil {
		(*Display.dummy_surface).Destroy()
	}

	destroy_cursors(Display)

	if Display.xdg_shell != nil {
		Display.xdg_shell.Destroy()
	}

	if Display.shm != nil {
		wlclient.ShmDestroy((Display.shm))
	}

	wlclient.RegistryDestroy((Display.registry))

	syscall.Close(int(Display.epoll_fd))

	wlclient.DisplayDisconnect((Display.Display))
}

//line 6478
func display_defer(Display *Display /*task *task,*/, fun runner) {

	Display.deferred_list_new = append(Display.deferred_list_new, (fun))
}

//line 6501
func DisplayRun(Display *Display) {

	Display.running = 1
	for {

		for len(Display.deferred_list_new) > 0 {

			Display.deferred_list_new[0].Run(0)

			Display.deferred_list_new = Display.deferred_list_new[1:]

		}

		if Display.running == 0 {
			break
		}

		if wlclient.DisplayRun(Display.Display) != nil {
			return
		}

	}
}

func (Display *Display) Exit() {
	Display.running = 0
}
