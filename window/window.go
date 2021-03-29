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

//import zwp "github.com/neurlang/wayland/wayland"
import "github.com/neurlang/wayland/wlclient"
import "github.com/neurlang/wayland/wlcursor"
import "github.com/neurlang/wayland/wl"
import zxdg "github.com/neurlang/wayland/xdg"
import cairo "github.com/neurlang/wayland/cairoshim"

import "github.com/neurlang/wayland/os"

import "errors"

import "fmt"

type runner interface {
	Run(uint32)
}

const SurfaceOpaque = 0x01
const SurfaceShm = 0x02

const SurfaceHintResize = 0x10
const SurfaceHintRgb565 = 0x100

const WindowPreferredFormatNone = 0
const WindowPreferredFormatRgb565 = 1

const WindowBufferTypeEglWindow = 0
const WindowBufferTypeShm = 1

const CursorBottomLeft = 0
const CursorBottomRight = 1
const CursorBottom = 2
const CursorDragging = 3
const CursorLeftPtr = 4
const CursorLeft = 5
const CursorRight = 6
const CursorTopLeft = 7
const CursorTopRight = 8
const CursorTop = 9
const CursorIbeam = 10
const CursorHand1 = 11
const CursorWatch = 12
const CursorDndMove = 13
const CursorDndCopy = 14
const CursorDndForbidden = 15
const CursorBlank = 16

const ZwpRelativePointerManagerV1Version = 1
const ZwpPointerConstraintsV1Version = 1

type global struct {
	name    uint32
	iface   string
	version uint32
}

type Display struct {
	Display            *wl.Display
	registry           *wl.Registry
	compositor         *wl.Compositor
	subcompositor      *wl.Subcompositor
	shm                *wl.Shm
	dataDeviceManager  *wl.DataDeviceManager
	textCursorPosition *struct{}
	xdgShell           *zxdg.WmBase
	serial             uint32

	//display_fd        int32
	displayFdEvents uint32

	//display_task task
	//	pad4		uint64
	//	pad5		uint64
	//	pad6		uint64

	deferredList [2]uintptr
	//	pad7		uint64
	//	pad8		uint64

	running int32

	globalList []*global
	//	pad9		uint64
	//	pada		uint64
	windowList [2]*Window
	//	padb		uint64
	//	padc		uint64
	inputList []*Input
	//	padd		uint64
	//	pade		uint64
	outputList [2]*output
	//	padf		uint64
	//	padg		uint64

	theme       *theme
	cursorTheme *wlcursor.Theme
	cursors     *[lengthCursors]*wlcursor.Cursor

	xkbContext *struct{}

	/* A hack to get text extents for tooltips */
	dummySurface *cairo.Surface

	hasRgb565                int32
	dataDeviceManagerVersion uint32

	deferredListNew []runner

	//display_task_new os.Runner
	surface2window map[*wl.Surface]*Window

	globalHandler GlobalHandler

	userData interface{}

	seatHandler SeatHandler
}

type rectangle struct {
	x      int32
	y      int32
	width  int32
	height int32
}

type toysurface interface {
	prepare(dx int, dy int, width int32, height int32, flags uint32,
		bufferTransform uint32, bufferScale int32) cairo.Surface
	swap(bufferTransform uint32, bufferScale int32, serverAllocation *rectangle)
	acquire(ctx *struct{}) int
	release()
	destroy()
}

type surface struct {
	Window *Window

	surface_            *wl.Surface
	subsurface          *wl.Subsurface
	synchronized        int32
	synchronizedDefault int32
	toysurface          *toysurface
	Widget              *Widget
	redrawNeeded        int32

	frameCb  *wl.Callback
	lastTime uint32
	//	pad1	uint32

	allocation       rectangle
	serverAllocation rectangle

	inputRegion  *wl.Region
	opaqueRegion *wl.Region

	bufferType      int32
	bufferTransform int32
	bufferScale     int32

	cairoSurface cairo.Surface
}

func (s *surface) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	s.CallbackDone(ev.C, ev.CallbackData)
}

type Window struct {
	Display          *Display
	windowOutputList [2]uintptr

	title string

	savedAllocation   rectangle
	minAllocation     rectangle
	pendingAllocation rectangle
	lastGeometry      rectangle

	x, y int32

	redrawInhibited     int32
	redrawNeeded        int32
	redrawTaskScheduled int32

	//redraw_task task

	//	pad1	uint64
	//	pad2	uint64

	resizeNeeded int32
	custom       int32
	focused      int32

	resizing int32

	fullscreen int32
	maximized  int32

	preferredFormat int

	mainSurface *surface
	xdgSurface  *zxdg.Surface
	xdgToplevel *zxdg.Toplevel
	xdgPopup    *zxdg.Popup

	parent     *Window
	lastParent *Window

	/* struct surface::link, contains also mainSurface */
	subsurfaceList [2]*surface

	pointerLocked bool

	confined bool

	link [2]*Window

	Userdata WidgetHandler

	redrawRunner runner

	subsurfaceListNew []*surface

	keyboardHandler KeyboardHandler

	frame *windowFrame
}

func (Window *Window) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {
	Window.SurfaceConfigure(Window.xdgSurface, ev.Serial)
}

func (Window *Window) SurfaceConfigure(zxdgSurfaceV6 *zxdg.Surface, serial uint32) {

	_ = Window.xdgSurface.AckConfigure(serial)

	windowUninhibitRedraw(Window)

}

func (Window *Window) SetKeyboardHandler(handler KeyboardHandler) {

	Window.keyboardHandler = handler

}

func (Window *Window) HandleToplevelConfigure(ev zxdg.ToplevelConfigureEvent) {
	Window.ToplevelConfigure(Window.xdgToplevel, ev.Width, ev.Height, ev.States)
}

func (Window *Window) ToplevelConfigure(
	zxdgToplevelV6 *zxdg.Toplevel,
	width int32,
	height int32,
	states []int32,
) {

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
	} else if (Window.savedAllocation.width > 0) &&
		(Window.savedAllocation.height > 0) {
		Window.ScheduleResize(Window.savedAllocation.width, Window.savedAllocation.height)
	}

}

func (Window *Window) HandleToplevelClose(ev zxdg.ToplevelCloseEvent) {
	Window.ToplevelClose(Window.xdgToplevel)
}

func (Window *Window) ToplevelClose(zxdgToplevelV6 *zxdg.Toplevel) {

	Window.Display.Exit()
}

func SurfaceEnter(wlSurface *wl.Surface, wlOutput *wl.Output) {
}
func SurfaceLeave(wlSurface *wl.Surface, wlOutput *wl.Output) {
}

type Widget struct {
	Window     *Window
	surface    *surface
	tooltip    *struct{}
	childList  *widgetList
	allocation rectangle

	opaque        int32
	tooltipCount  int32
	defaultCursor int32

	/* If this is set to false then no cairo surface will be
	 * created before redrawing the surface. This is useful if the
	 * redraw handler is going to do completely custom rendering
	 * such as using EGL directly */
	useCairo int32

	Userdata WidgetHandler
}

type widgetList struct {
	l []*Widget
}

func (l *widgetList) Add(w *Widget) {
	l.l = append(l.l, w)
}

func (l *widgetList) Remove(w *Widget) {
	if len(l.l) > 0 {
		if (l.l)[0] == w {
			l.l = (l.l)[1:]
			return
		}
		if (l.l)[len(l.l)-1] == w {
			l.l = (l.l)[0 : len(l.l)-1]
			return
		}
	}

	for i, v := range l.l {
		if v == w {
			l.l = append((l.l)[0:i], (l.l)[i+1:]...)
		}
	}
}
func (l *widgetList) Insert(w *Widget) {
	w.childList = l
	l.Add(w)
}

type WidgetHandler interface {
	Resize(Widget *Widget, width int32, height int32)
	Redraw(Widget *Widget)
	Enter(Widget *Widget, Input *Input, x float32, y float32)
	Leave(Widget *Widget, Input *Input)
	Motion(Widget *Widget, Input *Input, time uint32, x float32, y float32) int
	Button(
		Widget *Widget,
		Input *Input,
		time uint32,
		button uint32,
		state wl.PointerButtonState,
		data WidgetHandler,
	)
	TouchUp(Widget *Widget, Input *Input, serial uint32, time uint32, id int32)
	TouchDown(
		Widget *Widget,
		Input *Input,
		serial uint32,
		time uint32,
		id int32,
		x float32,
		y float32,
	)
	TouchMotion(Widget *Widget, Input *Input, time uint32, id int32, x float32, y float32)
	TouchFrame(Widget *Widget, Input *Input)
	TouchCancel(Widget *Widget, width int32, height int32)
	Axis(Widget *Widget, Input *Input, time uint32, axis uint32, value wl.Fixed)
	AxisSource(Widget *Widget, Input *Input, source uint32)
	AxisStop(Widget *Widget, Input *Input, time uint32, axis uint32)
	AxisDiscrete(Widget *Widget, Input *Input, axis uint32, discrete int32)
	PointerFrame(Widget *Widget, Input *Input)
}

type xkbModMaskT uint32

type Input struct {
	Display            *Display
	seat               *wl.Seat
	pointer            *wl.Pointer
	keyboard           *wl.Keyboard
	touch              *wl.Touch
	touchPointList     [2]uintptr
	pointerFocus       *Window
	keyboardFocus      *Window
	touchFocus         int32
	currentCursor      int32
	cursorAnimStart    uint32
	cursorFrameCb      *wl.Callback
	cursorTimerStart   uint32
	cursorAnimCurrent  uint32
	cursorDelayFd      int32
	cursorTimerRunning bool
	//cursor_task          task
	pointerSurface     *wl.Surface
	modifiers          uint32
	pointerEnterSerial uint32
	cursorSerial       uint32
	sx                 float32
	sy                 float32

	focusWidget *Widget
	grab        *Widget
	grabButton  uint32

	dataDevice      *wl.DataDevice
	touchGrab       uint32
	touchGrabId     int32
	dragX           float32
	dragY           float32
	dragFocus       *Window
	dragEnterSerial uint32

	xkb struct {
		keymap *struct{}
		state  *struct{}

		controlMask xkbModMaskT
		altMask     xkbModMaskT
		shiftMask   xkbModMaskT
	}

	repeatRateSec   int32
	repeatRateNsec  int32
	repeatDelaySec  int32
	repeatDelayNsec int32

	//repeat_task     task
	repeatSym   uint32
	repeatKey   uint32
	repeatTime  uint32
	seatVersion int32
}

func (i *Input) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	i.CallbackDone(ev.C, ev.CallbackData)
}

type KeyboardHandler interface {
	Key(
		window *Window,
		input *Input,
		time uint32,
		key uint32,
		unicode uint32,
		state wl.KeyboardKeyState,
		data WidgetHandler,
	)
	Focus(window *Window, device *Input)
}

func inputRemoveKeyboardFocus(input *Input) {
	var window = input.keyboardFocus

	if window == nil {
		return
	}

	if window.keyboardHandler != nil {
		window.keyboardHandler.Focus(window, nil)
	}

	input.keyboardFocus = nil
}

type output struct {
	Display        *Display
	output         *wl.Output
	serverOutputId uint32
	allocation     rectangle
	link           [2]*output
	transform      int32
	scale          int32
	maker          string
	model          string
}

type shmPool struct {
	pool *wl.ShmPool
	size uintptr
	used uintptr
	data []byte
}

const CursorDefault = 100
const CursorUnset = 101

//line 509
func surfaceToBufferSize(
	bufferTransform uint32,
	bufferScale int32,
	width *int32,
	height *int32,
) {

	switch bufferTransform {
	case wl.OutputTransform90:
		fallthrough
	case wl.OutputTransform270:
		fallthrough
	case wl.OutputTransformFlipped90:
		fallthrough
	case wl.OutputTransformFlipped270:
		*width, *height = *height, *width
	}

	*width *= bufferScale
	*height *= bufferScale
}

//line 532
func bufferToSurfaceSize(
	bufferTransform uint32,
	bufferScale int32,
	width *int32,
	height *int32,
) {
	switch bufferTransform {
	case wl.OutputTransform90:
		fallthrough
	case wl.OutputTransform270:
		fallthrough
	case wl.OutputTransformFlipped90:
		fallthrough
	case wl.OutputTransformFlipped270:
		*width, *height = *height, *width

	}

	*width /= bufferScale
	*height /= bufferScale
}

func (i *Input) HandlePointerEnter(ev wl.PointerEnterEvent) {
	i.PointerEnter(nil, ev.Serial, ev.Surface, ev.SurfaceX, ev.SurfaceY)
}

func (i *Input) PointerEnter(
	wlPointer *wl.Pointer,
	serial uint32,
	surface *wl.Surface,
	sx float32,
	sy float32,
) {

	if nil == surface {
		/* enter event for a Window we've just destroyed */
		return
	}

	var Window = i.Display.surface2window[surface]

	if surface != Window.mainSurface.surface_ {
		//		DBG("Ignoring Input event from subsurface %p\n", surface);
		return
	}

	i.Display.serial = serial
	i.pointerEnterSerial = serial
	i.pointerFocus = Window

	i.sx = sx
	i.sy = sy

}

func (i *Input) HandlePointerLeave(ev wl.PointerLeaveEvent) {
	i.PointerLeave(nil, ev.Serial, ev.Surface)
}

func (i *Input) PointerLeave(wlPointer *wl.Pointer, serial uint32, wlSurface *wl.Surface) {

	i.Display.serial = serial
	inputRemovePointerFocus(i)

}

func (i *Input) HandlePointerMotion(ev wl.PointerMotionEvent) {
	i.PointerMotion(ev.P, ev.Time, ev.SurfaceX, ev.SurfaceY)
}

func (i *Input) PointerMotion(
	wlPointer *wl.Pointer,
	time uint32,
	surfaceX float32,
	surfaceY float32,
) {

	pointerHandleMotion(i, wlPointer, time, surfaceX, surfaceY)
}

func (i *Input) HandlePointerButton(ev wl.PointerButtonEvent) {
	i.PointerButton(ev.P, ev.Serial, ev.Time, ev.Button, ev.State)
}

func (i *Input) PointerButton(
	wlPointer *wl.Pointer,
	serial uint32,
	time uint32,
	button uint32,
	stateW uint32,
) {
	var widget *Widget
	var state = wl.PointerButtonState(stateW)

	i.Display.serial = serial
	if i.focusWidget != nil && i.grab == nil &&
		state == wl.PointerButtonStatePressed {
		inputGrab(i, i.focusWidget, button)
	}

	widget = i.grab
	if widget != nil && widget.Userdata != nil {
		widget.Userdata.Button(widget,
			i, time,
			button, state,
			i.grab.Userdata)
	}

	if i.grab != nil && i.grabButton == button &&
		state == wl.PointerButtonStateReleased {
		inputUngrab(i)
	}

}

func (i *Input) HandlePointerAxis(ev wl.PointerAxisEvent) {

}

func (*Input) PointerAxis(wlPointer *wl.Pointer, time uint32, axis uint32, value wl.Fixed) {
}

func (i *Input) HandlePointerFrame(ev wl.PointerFrameEvent) {

}

func (*Input) PointerFrame(wlPointer *wl.Pointer) {
}

func (i *Input) HandlePointerAxisSource(ev wl.PointerAxisSourceEvent) {

}

func (*Input) PointerAxisSource(wlPointer *wl.Pointer, axisSource uint32) {
}

func (i *Input) HandlePointerAxisStop(ev wl.PointerAxisStopEvent) {

}

func (*Input) PointerAxisStop(wlPointer *wl.Pointer, time uint32, axis uint32) {
}

func (i *Input) HandlePointerAxisDiscrete(ev wl.PointerAxisDiscreteEvent) {

}

func (*Input) PointerAxisDiscrete(wlPointer *wl.Pointer, axis uint32, discrete int32) {
}

type SeatHandler interface {
	Capabilities(i *Input, seat *wl.Seat, caps uint32)
	Name(i *Input, seat *wl.Seat, name string)
}

func (i *Input) HandleSeatCapabilities(ev wl.SeatCapabilitiesEvent) {
	i.SeatCapabilities(i.seat, ev.Capabilities)
}

func (i *Input) HandleSeatName(ev wl.SeatNameEvent) {
	i.SeatName(i.seat, ev.Name)
}

func (input *Input) SeatCapabilities(seat *wl.Seat, caps uint32) {

	if ((caps & wl.SeatCapabilityPointer) != 0) && (input.pointer == nil) {
		var err error
		input.pointer, err = seat.GetPointer()
		if err != nil {
			fmt.Println(err)
			return
		}
		wlclient.PointerSetUserData(input.pointer, input)
		wlclient.PointerAddListener(input.pointer, input)

	} else if ((caps & wl.SeatCapabilityPointer) == 0) && (nil != input.pointer) {
		if input.seatVersion >= wl.PointerReleaseSinceVersion {
			_ = input.pointer.Release()
		} else {
			wlclient.PointerDestroy(input.pointer)
		}
		input.pointer = nil
	}

	if ((caps & wl.SeatCapabilityKeyboard) != 0) && input.keyboard == nil {
		var err error
		input.keyboard, err = seat.GetKeyboard()
		if err != nil {
			fmt.Println(err)
			return
		}
		wlclient.KeyboardSetUserData(input.keyboard, input)
		wlclient.KeyboardAddListener(input.keyboard, input)
	} else if 0 == (caps&wl.SeatCapabilityKeyboard) && input.keyboard != nil {
		if input.seatVersion >= wl.KeyboardReleaseSinceVersion {
			_ = input.keyboard.Release()
		} else {
			wlclient.KeyboardDestroy(input.keyboard)
		}
		input.keyboard = nil
	}

	if ((caps & wl.SeatCapabilityTouch) != 0) && input.touch == nil {
		var err error
		input.touch, err = seat.GetTouch()
		if err != nil {
			fmt.Println(err)
			return
		}
		wlclient.TouchSetUserData(input.touch, input)
		wlclient.TouchAddListener(input.touch, input)
	} else if 0 == (caps&wl.SeatCapabilityTouch) && input.touch != nil {
		if input.seatVersion >= wl.TouchReleaseSinceVersion {
			input.touch.Release()
		} else {
			wlclient.TouchDestroy(input.touch)
		}
		input.touch = nil
	}

	if input.Display.seatHandler != nil {
		input.Display.seatHandler.Capabilities(input, seat, caps)
	}

}

func (input *Input) HandleKeyboardEnter(e wl.KeyboardEnterEvent) {
	println("ENTER")
}

func (input *Input) keyboard_handle_key(keyboard *wl.Keyboard,
	serial uint32, time uint32, key uint32,
	state_w uint32) {
	var window *Window = input.keyboardFocus
	var state = wl.KeyboardKeyState(state_w)

	input.Display.serial = serial
	var code = key + 8
	if window == nil || input.xkb.state == nil {
		return
	}

	/* We only use input grabs for pointer events for now, so just
	 * ignore key presses if a grab is active.  We expand the key
	 * event delivery mechanism to route events to widgets to
	 * properly handle key grabs.  In the meantime, this prevents
	 * key event delivery while a grab is active. */
	if input.grab != nil && input.grabButton == 0 {
		return
	}

	_ = state
	_ = code

	//TODO:
	//num_syms = xkb_state_key_get_syms(input.xkb.state, code, &syms);
	/*
		sym = XKB_KEY_NoSymbol;
		if (num_syms == 1) {
			sym = syms[0];
		}


		if (sym == XKB_KEY_F5 && input.modifiers == MOD_ALT_MASK) {
			if (state == WL_KEYBOARD_KEY_STATE_PRESSED) {
				window_set_maximized(window, !window.maximized);
			}
		} else if (sym == XKB_KEY_F11 &&
			   window.fullscreen_handler &&
			   state == WL_KEYBOARD_KEY_STATE_PRESSED) {
			window.fullscreen_handler(window, window.user_data);
		} else if (sym == XKB_KEY_F4 &&
			   input.modifiers == MOD_ALT_MASK &&
			   state == WL_KEYBOARD_KEY_STATE_PRESSED) {
			window_close(window);
		} else if (window.key_handler) {
			if (state == WL_KEYBOARD_KEY_STATE_PRESSED) {
				sym = process_key_press(sym, input);
			}

			(*window.key_handler)(window, input, time, key,
					       sym, state, window.user_data);
		}

		if (state == WL_KEYBOARD_KEY_STATE_RELEASED &&
		    key == input.repeat_key) {
			toytimer_disarm(&input.repeat_timer);
		} else if (state == WL_KEYBOARD_KEY_STATE_PRESSED &&
			   xkb_keymap_key_repeats(input.xkb.keymap, code)) {
			input.repeat_sym = sym;
			input.repeat_key = key;
			input.repeat_time = time;
			its.it_interval.tv_sec = input.repeat_rate_sec;
			its.it_interval.tv_nsec = input.repeat_rate_nsec;
			its.it_value.tv_sec = input.repeat_delay_sec;
			its.it_value.tv_nsec = input.repeat_delay_nsec;
			toytimer_arm(&input.repeat_timer, &its);
		}
	*/
}

func (input *Input) HandleKeyboardKey(e wl.KeyboardKeyEvent) {
	input.keyboard_handle_key(nil, e.Serial, e.Time, e.Key, e.State)
}

func (input *Input) HandleKeyboardKeymap(e wl.KeyboardKeymapEvent) {

}
func (input *Input) HandleKeyboardLeave(e wl.KeyboardLeaveEvent) {
	println("LEAVE")
}
func (input *Input) HandleKeyboardModifiers(e wl.KeyboardModifiersEvent) {

}
func (input *Input) HandleKeyboardRepeatInfo(e wl.KeyboardRepeatInfoEvent) {

}
func (input *Input) HandleTouchCancel(e wl.TouchCancelEvent) {

}

func (input *Input) HandleTouchDown(e wl.TouchDownEvent) {

}

func (input *Input) HandleTouchFrame(e wl.TouchFrameEvent) {

}
func (input *Input) HandleTouchMotion(e wl.TouchMotionEvent) {

}
func (input *Input) HandleTouchOrientation(e wl.TouchOrientationEvent) {

}
func (input *Input) HandleTouchShape(e wl.TouchShapeEvent) {

}
func (input *Input) HandleTouchUp(e wl.TouchUpEvent) {

}
func (i *Input) SeatName(wlSeat *wl.Seat, name string) {
	if i.Display.seatHandler != nil {
		i.Display.seatHandler.Name(i, wlSeat, name)
	}
}

// line 2697
func inputGrab(input *Input, widget *Widget, button uint32) {
	input.grab = widget
	input.grabButton = button

	inputSetFocusWidget(input, widget, input.sx, input.sy)
}

// line 2706
func inputUngrab(input *Input) {

	input.grab = nil
	if input.pointerFocus != nil {
		var widget = windowFindWidget(input.pointerFocus,
			int32(input.sx), int32(input.sy))
		inputSetFocusWidget(input, widget, input.sx, input.sy)
	}
}

type shmSurfaceData struct {
	buffer *wl.Buffer
	pool   *shmPool
}

//line 734
func shmSurfaceDataDestroy(data *shmSurfaceData) {
	data.buffer.Destroy()
	if data.pool != nil {
		shmPoolDestroy(data.pool)
	}
}

//line 744
func makeShmPool(Display *Display, size uintptr, data *[]byte) (pool *wl.ShmPool) {
	fd, err := os.CreateAnonymousFile(int64(size))
	if err != nil {
		println("creating a buffer file failed")
		println(err.Error())
		return nil
	}

	*data, err = os.Mmap(int(fd.Fd()), 0, int(size), os.ProtRead|os.ProtWrite, os.MapShared)
	if err != nil {
		println("mmap failed")
		fd.Close()
		return nil
	}

	pool, err = Display.shm.CreatePool(fd.Fd(), int32(size))
	if err != nil {
		println("create pool failed")
		fd.Close()
		return nil
	}

	fd.Close()

	return pool
}

//line 772
func shmPoolCreate(Display *Display, size uintptr) *shmPool {
	var pool = &shmPool{}

	pool.pool = makeShmPool(Display, size, &pool.data)

	if pool.pool == nil {
		return nil
	}

	pool.size = size
	pool.used = 0

	return pool
}

//line 792
func shmPoolAllocate(pool *shmPool, size uintptr, offset *int) (ret []byte) {

	if pool.used+size > pool.size {
		return nil
	}

	*offset = int(pool.used)
	ret = pool.data[pool.used:]
	pool.used += size
	pool.data = pool.data[0:pool.used]

	return ret
}

//line 804
/* destroy the pool. this does not unmap the memory though */
func shmPoolDestroy(pool *shmPool) {

	err := os.Munmap(pool.data)
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
func dataLengthForShmSurface(rect *rectangle) uintptr {
	var stride = int32(cairo.FormatStrideForWidth(cairo.FormatArgb32, int(rect.width)))
	return uintptr(int(stride * rect.height))
}

func shmPoolReset(pool *shmPool) {
	pool.used = 0
}

//line 829
func displayCreateShmSurfaceFromPool(Display *Display,
	rectangle *rectangle,
	flags uint32, pool *shmPool) (*cairo.Surface, *shmSurfaceData) {
	var data = &shmSurfaceData{}
	var format uint32
	var surface cairo.Surface
	var cairoFormat cairo.Format
	var stride, length int
	var offset int
	var map_ []byte
	var err error

	if (flags&uint32(SurfaceHintRgb565) != 0) && Display.hasRgb565 != 0 {
		cairoFormat = cairo.FormatRgb16565
	} else {
		cairoFormat = cairo.FormatArgb32
	}

	stride = cairo.FormatStrideForWidth(cairoFormat, int(rectangle.width))

	length = stride * int(rectangle.height)
	data.pool = nil

	map_ = shmPoolAllocate(pool, uintptr(length), &offset)

	if map_ == nil {
		return nil, nil
	}

	surface = cairo.ImageSurfaceCreateForData(map_,
		cairoFormat,
		int(rectangle.width),
		int(rectangle.height),
		stride)

	surface.SetUserData(func() {

		shmSurfaceDataDestroy(data)
	})

	if (flags&uint32(SurfaceHintRgb565) != 0) && Display.hasRgb565 != 0 {
		format = wl.ShmFormatRgb565
	} else {
		if flags&SurfaceOpaque != 0 {
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
func displayCreateShmSurface(Display *Display,
	rectangle *rectangle, flags uint32,
	alternatePool *shmPool,
	dataRet **shmSurfaceData) *cairo.Surface {
	var data *shmSurfaceData
	var pool *shmPool
	var surface *cairo.Surface

	if alternatePool != nil {
		shmPoolReset(alternatePool)

		surface, data = displayCreateShmSurfaceFromPool(
			Display,
			rectangle,
			flags,
			alternatePool,
		)

		if surface != nil {
			goto out
		}
	}

	pool = shmPoolCreate(Display, dataLengthForShmSurface(rectangle))

	if pool == nil {
		return nil
	}

	surface, data =
		displayCreateShmSurfaceFromPool(Display, rectangle, flags, pool)

	if surface == nil {
		shmPoolDestroy(pool)
		return nil
	}

	/* make sure we destroy the pool when the surface is destroyed */
	data.pool = pool

out:
	if dataRet != nil {
		*dataRet = data
	}

	return surface
}

type shmSurfaceLeaf struct {
	cairoSurface *cairo.Surface
	/* 'data' is automatically destroyed, when 'cairo_surface' is */
	data *shmSurfaceData

	resizePool *shmPool
	busy       int32
}

func shmSurfaceLeafRelease(leaf *shmSurfaceLeaf) {
	if leaf.cairoSurface != nil {
		(*leaf.cairoSurface).Destroy()
	}
	/* leaf.data already destroyed via cairo private */
}

const MaxLeaves = 3

//line 983
type shmSurface struct {
	Display *Display
	surface *wl.Surface
	flags   uint32
	dx      int32
	dy      int32

	leaf    [MaxLeaves]shmSurfaceLeaf
	current *shmSurfaceLeaf
}

func shmSurfaceBufferRelease(surface *shmSurface, buffer *wl.Buffer) error {
	var leaf *shmSurfaceLeaf
	var i int
	var freeFound int

	for i = 0; i < MaxLeaves; i++ {
		leaf = &surface.leaf[i]
		if leaf.data != nil && leaf.data.buffer == buffer {
			leaf.busy = 0
			break
		}
	}
	if i >= MaxLeaves {
		return errors.New("unknown buffer released")
	}

	/* Leave one free leaf with storage, release others */
	freeFound = 0
	for i = 0; i < MaxLeaves; i++ {
		leaf = &surface.leaf[i]

		if (leaf.cairoSurface == nil) || (leaf.busy != 0) {
			continue
		}

		if freeFound == 0 {
			freeFound = 1
		} else {
			shmSurfaceLeafRelease(leaf)

		}
	}
	return nil
}

func (s *shmSurface) HandleBufferRelease(ev wl.BufferReleaseEvent) {
	s.BufferRelease(ev.B)
}

func (s *shmSurface) BufferRelease(buf *wl.Buffer) {
	shmSurfaceBufferRelease(s, buf)

}

func (s *shmSurface) prepare(dx int, dy int, width int32, height int32, flags uint32,
	bufferTransform uint32, bufferScale int32) cairo.Surface {

	var resizeHint = (flags & SurfaceHintResize) != 0
	surface := s
	var rect rectangle
	var leaf *shmSurfaceLeaf
	var i int

	surface.dx = int32(dx)
	surface.dy = int32(dy)

	for i = 0; i < MaxLeaves; i++ {
		if surface.leaf[i].busy != 0 {
			continue
		}

		if leaf == nil || surface.leaf[i].cairoSurface != nil {
			leaf = &surface.leaf[i]
		}
	}

	if nil == leaf {
		panic("all buffers are held by the server.\n")

	}

	if !resizeHint && (leaf.resizePool != nil) {
		(*leaf.cairoSurface).Destroy()
		leaf.cairoSurface = nil
		shmPoolDestroy(leaf.resizePool)
		leaf.resizePool = nil
	}

	surfaceToBufferSize(bufferTransform, bufferScale, &width, &height)

	if (leaf.cairoSurface != nil) &&
		(int32((*leaf.cairoSurface).ImageSurfaceGetWidth()) == width) &&
		(int32((*leaf.cairoSurface).ImageSurfaceGetHeight()) == height) {
		goto out
	}

	if leaf.cairoSurface != nil {
		(*leaf.cairoSurface).Destroy()
	}

	rect.width = width
	rect.height = height

	leaf.cairoSurface = displayCreateShmSurface(
		surface.Display,
		&rect,
		surface.flags,
		leaf.resizePool,
		&leaf.data,
	)

	if leaf.cairoSurface == nil {
		return nil
	}

	wlclient.BufferAddListener(leaf.data.buffer, surface)

out:
	surface.current = leaf

	return (*leaf.cairoSurface).Reference()
}

//line 1146
func shmSurfaceSwap(surface *shmSurface, bufferTransform uint32,
	bufferScale int32, serverAllocation *rectangle) {
	var leaf = surface.current

	serverAllocation.width =
		int32((*leaf.cairoSurface).ImageSurfaceGetWidth())
	serverAllocation.height =
		int32((*leaf.cairoSurface).ImageSurfaceGetHeight())

	bufferToSurfaceSize(bufferTransform, bufferScale,
		&serverAllocation.width,
		&serverAllocation.height)

	_ = surface.surface.Attach(leaf.data.buffer,
		surface.dx, surface.dy)
	_ = surface.surface.Damage(0, 0,
		serverAllocation.width, serverAllocation.height)
	_ = surface.surface.Commit()

	leaf.busy = 1
	surface.current = nil
}

func (s *shmSurface) swap(
	bufferTransform uint32,
	bufferScale int32,
	serverAllocation *rectangle,
) {
	shmSurfaceSwap(s, bufferTransform, bufferScale, serverAllocation)

}

func (*shmSurface) acquire(ctx *struct{}) int {
	return -1
}

func (*shmSurface) release() {
}

func shmSurfaceDestroy(surface *shmSurface) {
	var i int

	for i = 0; i < MaxLeaves; i++ {
		shmSurfaceLeafRelease(&surface.leaf[i])
	}
}

func (s *shmSurface) destroy() {

	shmSurfaceDestroy(s)
}

//line 1199
func shmSurfaceCreate(Display *Display, wlSurface *wl.Surface,
	flags uint32, rectangle *rectangle) toysurface {
	var surface = &shmSurface{}

	surface.Display = Display
	surface.surface = wlSurface
	surface.flags = flags

	return surface
}

const lengthCursors = 16

//line 1343
func createCursors(Display *Display) (err error) {

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
	Display.cursorTheme = theme

	var wlCursors = [lengthCursors]*wlcursor.Cursor{}

	Display.cursors = &wlCursors

	for i := range Cursors {
		for j := range Cursors[i] {

			var str = Cursors[i][j]

			str = str[:len(str)-1]

			cursor, err = Display.cursorTheme.GetCursor(str)
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
func destroyCursors(Display *Display) {
	Display.cursorTheme.Destroy()
}

//line 1402
func surfaceFlush(surface *surface) {
	if surface.cairoSurface == nil {
		return
	}

	if surface.opaqueRegion != nil {
		_ = surface.surface_.SetOpaqueRegion(surface.opaqueRegion)
		surface.opaqueRegion.Destroy()
		surface.opaqueRegion = nil
	}

	if surface.inputRegion != nil {
		_ = surface.surface_.SetInputRegion(surface.inputRegion)
		surface.inputRegion.Destroy()
		surface.inputRegion = nil
	}

	(*surface.toysurface).swap(uint32(surface.bufferTransform), surface.bufferScale,
		&surface.serverAllocation)

	surface.cairoSurface.Destroy()
	surface.cairoSurface = nil
}

//line 1462
func surfaceCreateSurface(surface *surface, flags uint32) {
	var Display = surface.Window.Display
	var allocation = surface.allocation

	if surface.toysurface == nil {
		var toy = shmSurfaceCreate(Display, surface.surface_, flags, &allocation)

		surface.toysurface = &toy
	}

	surface.cairoSurface = (*surface.toysurface).prepare(
		0, 0,
		allocation.width, allocation.height, flags,
		uint32(surface.bufferTransform), surface.bufferScale)

}

//line 1488
func windowCreateMainSurface(Window *Window) {
	var surface = Window.mainSurface
	var flags uint32 = 0

	if Window.resizing != 0 {
		flags |= SurfaceHintResize
	}

	if Window.preferredFormat == WindowPreferredFormatRgb565 {
		flags |= SurfaceHintRgb565
	}

	surfaceCreateSurface(surface, flags)

}

//line 1552
func surfaceDestroy(surface *surface) {
	if surface.frameCb != nil {
		wlclient.CallbackDestroy(surface.frameCb)
	}

	if surface.inputRegion != nil {
		wlclient.RegionDestroy(surface.inputRegion)
	}

	if surface.opaqueRegion != nil {
		wlclient.RegionDestroy(surface.opaqueRegion)
	}

	if surface.subsurface != nil {
		wlclient.SubsurfaceDestroy(surface.subsurface)
	}

	surface.surface_.Destroy()

	if surface.toysurface != nil {
		(*surface.toysurface).destroy()
	}

}

//line 1577
func (Window *Window) Destroy() {

	if Window.xdgToplevel != nil {
		Window.xdgToplevel.Destroy()
	}
	if Window.xdgPopup != nil {
		Window.xdgPopup.Destroy()
	}
	if Window.xdgSurface != nil {
		Window.xdgSurface.Destroy()
	}

	surfaceDestroy(Window.mainSurface)

}

//line 1624
func widgetFindWidget(Widget *Widget, x int32, y int32) *Widget {

	if Widget.allocation.x <= x &&
		x < Widget.allocation.x+Widget.allocation.width &&
		Widget.allocation.y <= y &&
		y < Widget.allocation.y+Widget.allocation.height {
		return Widget
	}

	return nil
}

//line 1645
func windowFindWidget(Window *Window, x int32, y int32) *Widget {
	var surface *surface
	var Widget *Widget

	for _, surface = range Window.subsurfaceListNew {

		Widget = widgetFindWidget(surface.Widget, x, y)

		if Widget != nil {
			return Widget
		}
	}

	return nil
}

//line 1655
func widgetCreate(Window *Window, surface *surface, data WidgetHandler) *Widget {
	var w = new(Widget)
	w.Window = Window
	w.surface = surface
	w.Userdata = data
	w.allocation = surface.allocation
	w.childList = new(widgetList)
	w.opaque = 0
	w.tooltip = nil
	w.tooltipCount = 0
	w.defaultCursor = CursorLeftPtr
	w.useCairo = 1

	return w
}

//line 1675
func (Window *Window) AddWidget(data WidgetHandler) *Widget {
	var w = widgetCreate(Window, Window.mainSurface, data)

	Window.mainSurface.Widget = w

	return w
}

//line 1702
func (parent *Widget) AddWidget(data WidgetHandler) *Widget {
	widget := widgetCreate(parent.Window, parent.surface, data)

	parent.childList.Insert(widget)

	return widget
}

//line 1701
func (parent *Widget) Destroy() {

	var surface = parent.surface

	/* Destroy the sub-surface along with the root Widget */
	if (surface.Widget == parent) && (surface.subsurface != nil) {
		surfaceDestroy(parent.surface)
	}

}

func (d *Display) SetSeatHandler(h SeatHandler) {
	d.seatHandler = h
}

func (d *Display) HandleWmBasePing(ev zxdg.WmBasePingEvent) {
	d.ShellPing(d.xdgShell, ev.Serial)
}

func (d *Display) ShellPing(shell *zxdg.WmBase, serial uint32) {
	shell.Pong(serial)
}

func minU32(a, b uint32) uint32 {
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

	d.globalList = append(d.globalList, global)

	switch iface {

	case "wl_compositor":
		d.compositor = wlclient.RegistryBindCompositorInterface(d.registry, id, 1)

	case "wl_output":

		displayAddOutput(d, id)
		// TODO
	case "wl_seat":

		displayAddInput(d, id, int(version))

	case "wl_shm":
		d.shm = wlclient.RegistryBindShmInterface(d.registry, id, 1)
		wlclient.ShmAddListener(d.shm, d)
	case "wl_data_device_manager":
		d.dataDeviceManagerVersion = minU32(version, 3)

		wlclient.RegistryBindDataDeviceManagerInterface(d.registry, id,
			d.dataDeviceManagerVersion)

	//case "zxdg_shell_v6":
	case "xdg_wm_base":

		d.xdgShell = wlclient.RegistryBindWmBaseInterface(d.registry, id, 1)

		zxdg.WmBaseAddListener(d.xdgShell, d)

	case "text_cursor_position":
	case "wl_subcompositor":

	default:

	}
	if d.globalHandler != nil {
		d.globalHandler.HandleGlobal(d, id, iface, version, d.userData)
	}

}
func (d *Display) RegistryGlobalRemove(wlRegistry *wl.Registry, name uint32) {

}

type GlobalHandler interface {
	HandleGlobal(d *Display, id uint32, iface string, version uint32, data interface{})
}

func (d *Display) SetGlobalHandler(gh GlobalHandler) {
	d.globalHandler = gh
	if gh == nil {
		return
	}
	for _, v := range d.globalList {
		d.globalHandler.HandleGlobal(d, v.name, v.iface, v.version, d.userData)
	}
}

func (d *Display) HandleShmFormat(e wl.ShmFormatEvent) {
	d.ShmFormat(nil, e.Format)
}
func (d *Display) ShmFormat(wlShm *wl.Shm, format uint32) {
	print("SHM FORMAT: ")
	println(format)
}

//line 1733
func widgetSetSize(Widget *Widget, width int32, height int32) {
	Widget.allocation.width = width
	Widget.allocation.height = height
}

//line 1740
func widgetSetAllocation(Widget *Widget, x int32, y int32, width int32, height int32) {
	Widget.allocation.x = x
	Widget.allocation.y = y
	widgetSetSize(Widget, width, height)
}

// line 1763
func widgetGetCairoSurface(Widget *Widget) cairo.Surface {
	var surface = Widget.surface
	var Window = Widget.Window

	if Widget.useCairo == 0 {
		panic("assert")
	}

	if nil == surface.cairoSurface {
		if surface == Window.mainSurface {
			windowCreateMainSurface(Window)

		} else {
			surfaceCreateSurface(surface, 0)

		}
	}

	return surface.cairoSurface
}

// line 1887
func (parent *Widget) WidgetGetLastTime() uint32 {
	return parent.surface.lastTime
}

//line 2013
func (parent *Widget) WidgetScheduleRedraw() {
	parent.surface.redrawNeeded = 1
	windowScheduleRedrawTask(parent.Window)
}

//line 2036
func (Window *Window) WindowGetSurface() cairo.Surface {
	var cairoSurface = widgetGetCairoSurface(Window.mainSurface.Widget)

	return cairoSurface.Reference()
}

func (Window *Window) FrameCreate(data WidgetHandler) *Widget {
	var buttons uint32

	if Window.custom != 0 {
		buttons = FrameButtonNone
	} else {
		buttons = FrameButtonAll
	}

	var frame = new(windowFrame)
	frame.frame = frameCreate(Window.Display.theme, 0, 0, buttons, Window.title, nil)
	if frame.frame == nil {
		frame = nil
		return nil
	}

	frame.widget = Window.AddWidget(frame)
	frame.child = frame.widget.AddWidget(data)

	Window.frame = frame

	return frame.child
}

//line 2614
func inputSetFocusWidget(Input *Input, focus *Widget,
	x float32, y float32) {
	var old, Widget *Widget
	var cursor int

	if focus == Input.focusWidget {
		return
	}

	old = Input.focusWidget
	if old != nil {
		Input.focusWidget = nil
	}

	if focus != nil {
		Widget = focus
		if Input.grab != nil {
			Widget = Input.grab
		}
		Input.focusWidget = focus
		cursor = int(Widget.defaultCursor)

		inputSetPointerImage(Input, cursor)
	}
}

//line 2714
func cancelPointerImageUpdate(Input *Input) {

}

// line 2718
func inputRemovePointerFocus(input_ *Input) {
	var Window = input_.pointerFocus

	if nil == Window {
		return
	}

	inputSetFocusWidget(input_, nil, 0, 0)

	input_.pointerFocus = nil
	input_.currentCursor = CursorUnset

	cancelPointerImageUpdate(input_)
}

// line 2776
func pointerHandleMotion(data *Input, pointer *wl.Pointer,
	time uint32, sx float32, sy float32) {
	var Input = data
	_ = Input
	var Window = Input.pointerFocus
	var Widget *Widget
	var cursor int

	if Window == nil {
		return
	}

	Input.sx = sx
	Input.sy = sy

	// when making the Window smaller - e.g. after an unmaximise we might
	// * still have a pending motion event that the compositor has picked
	// * based on the old surface dimensions. However, if we have an active
	// * grab, we expect to see Input from outside the Window anyway.

	if nil == Input.grab && (sx < float32(Window.mainSurface.allocation.x) ||
		sy < float32(Window.mainSurface.allocation.y) ||
		sx > float32(Window.mainSurface.allocation.width) ||
		sy > float32(Window.mainSurface.allocation.height)) {
		return
	}

	if !(Input.grab != nil && Input.grabButton != 0) {
		Widget = windowFindWidget(Window, int32(sx), int32(sy))
		inputSetFocusWidget(Input, Widget, sx, sy)

	}

	if Input.grab != nil {
		Widget = Input.grab
	} else {
		Widget = Input.focusWidget
	}
	if Widget != nil {
		if Widget.Userdata != nil {
			cursor = Widget.Userdata.Motion(Input.focusWidget,
				Input, time, sx, sy)
		} else {
			cursor = int(Widget.defaultCursor)
		}
	} else {
		cursor = CursorLeftPtr
	}
	_ = cursor

	inputSetPointerImage(Input, cursor)
}

//line 3552
func inputGetSeat(Input *Input) *wl.Seat {
	return Input.seat
}

//line 3754
func inputSetPointerImageIndex(Input *Input, index int) {
	var buffer *wl.Buffer
	var cursor *wlcursor.Cursor
	var image wlcursor.Image

	if Input.pointer == nil {
		return
	}

	cursor = Input.Display.cursors[Input.currentCursor]
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

	_ = Input.pointerSurface.Attach(buffer, 0, 0)
	_ = Input.pointerSurface.Damage(0, 0,
		int32(image.GetWidth()), int32(image.GetHeight()))
	_ = Input.pointerSurface.Commit()
	wlcursor.PointerSetCursor(Input.pointer, Input.pointerEnterSerial, Input.pointerSurface,
		int32(image.GetHotspotX()), int32(image.GetHotspotY()))

}

//line 3789
func inputSetPointerSpecial(Input *Input) bool {
	if Input.currentCursor == CursorBlank {
		wlcursor.PointerSetCursor(Input.pointer,
			Input.pointerEnterSerial,
			nil, 0, 0)
		return true
	}

	if Input.currentCursor == CursorUnset {
		return true
	}

	return false
}

//line 3805
func schedulePointerImageUpdate(Input *Input,
	cursor *wlcursor.Cursor,
	duration uint32,
	forceFrame bool) {
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
	if !forceFrame && (duration > 100) {
		return
	}

	/* for short durations we'll just spin on frame callbacks for
	 * accurate timing - the way any kind of timing sensitive animation
	 * should really be done. */
	cb, err := Input.pointerSurface.Frame()
	if err != nil {
		fmt.Println(err)
		return
	}

	Input.cursorFrameCb = cb

	wlclient.CallbackAddListener(Input.cursorFrameCb, Input)

}

func (i *Input) CallbackDone(wlCallback *wl.Callback, callbackData uint32) {
	pointerSurfaceFrameCallback(i, wlCallback, callbackData)
}

//line 3842
func pointerSurfaceFrameCallback(Input *Input, callback *wl.Callback, time uint32) {
	var cursor *wlcursor.Cursor
	var i int
	var duration uint32
	var forceFrame = true

	cancelPointerImageUpdate(Input)

	if callback != nil {
		if callback != Input.cursorFrameCb {
			panic("assert")
		}
		wlclient.CallbackDestroy(callback)
		Input.cursorFrameCb = nil
		forceFrame = false
	}

	if Input.pointer == nil {
		return
	}

	if inputSetPointerSpecial(Input) {
		return
	}

	cursor = Input.Display.cursors[Input.currentCursor]
	if cursor == nil {
		return
	}

	/* FIXME We don't have the current time on the first call so we set
	 * the animation start to the time of the first frame callback. */
	if time == 0 {
		Input.cursorAnimStart = 0
	} else if Input.cursorAnimStart == 0 {
		Input.cursorAnimStart = time
	}

	Input.cursorAnimCurrent = time

	if time == 0 || Input.cursorAnimStart == 0 {
		duration = 0
		i = 0
	} else {
		frameDuration := cursor.FrameAndDuration(time - Input.cursorAnimStart)

		i, duration = frameDuration.FrameIndex, frameDuration.FrameDuration
	}

	if cursor.ImageCount() > 1 {
		schedulePointerImageUpdate(Input, cursor, duration,
			forceFrame)
	}

	inputSetPointerImageIndex(Input, i)
}

//line 3925
func inputSetPointerImage(Input *Input, pointer int) {
	var force bool

	if Input.pointer == nil {
		return
	}

	if Input.pointerEnterSerial > Input.cursorSerial {
		force = true
	}

	if !force && pointer == int(Input.currentCursor) {
		return
	}

	Input.currentCursor = int32(pointer)
	Input.cursorSerial = Input.pointerEnterSerial
	if Input.cursorFrameCb == nil {
		pointerSurfaceFrameCallback(Input, nil, 0)
	} else if force && (!inputSetPointerSpecial(Input)) {
		/* The current frame callback may be stuck if, for instance,
		 * the set cursor request was processed by the server after
		 * this client lost the focus. In this case the cursor surface
		 * might not be mapped and the frame callback wouldn't ever
		 * complete. Send a set_cursor and attach to try to map the
		 * cursor surface again so that the callback will finish */

		inputSetPointerImageIndex(Input, 0)
	}
}

// line 4104
func surfaceResize(surface *surface) {
	var Widget = surface.Widget

	if (surface.allocation.width != Widget.allocation.width) ||
		(surface.allocation.height != Widget.allocation.height) {
		windowScheduleRedraw(Widget.Window)

	}

	surface.allocation = Widget.allocation

}

//line 4144
func windowDoResize(Window *Window) {
	widgetSetAllocation(Window.mainSurface.Widget,
		Window.pendingAllocation.x,
		Window.pendingAllocation.y,
		Window.pendingAllocation.width,
		Window.pendingAllocation.height)

	surfaceResize(Window.mainSurface)

	if (Window.fullscreen != 0) && (Window.maximized != 0) {
		Window.savedAllocation = Window.pendingAllocation
	}
}

//line 4191
func idleResize(Window *Window) {
	Window.resizeNeeded = 0
	Window.redrawNeeded = 1

	windowDoResize(Window)
}

//line 4223
func (Window *Window) ScheduleResize(width int32, height int32) {
	/* We should probably get these numbers from the theme. */
	const minWidth = 200
	const minHeight = 200

	Window.pendingAllocation.x = 0
	Window.pendingAllocation.y = 0
	Window.pendingAllocation.width = width
	Window.pendingAllocation.height = height

	if Window.minAllocation.width == 0 {
		if width < minWidth {
			Window.minAllocation.width = minWidth
		} else {
			Window.minAllocation.width = width
		}
		if height < minHeight {
			Window.minAllocation.height = minHeight
		} else {
			Window.minAllocation.height = height
		}
	}

	if Window.pendingAllocation.width < Window.minAllocation.width {
		Window.pendingAllocation.width = Window.minAllocation.width
	}
	if Window.pendingAllocation.height < Window.minAllocation.height {
		Window.pendingAllocation.height = Window.minAllocation.height
	}

	Window.resizeNeeded = 1
	windowScheduleRedraw(Window)
}

//line 4254
func (parent *Widget) ScheduleResize(width int32, height int32) {
	parent.Window.ScheduleResize(width, height)
}

//line 4269
func windowInhibitRedraw(Window *Window) {
	Window.redrawInhibited = 1
	Window.redrawTaskScheduled = 0
}

// line 4284
func windowUninhibitRedraw(Window *Window) {
	Window.redrawInhibited = 0
	if (Window.redrawNeeded != 0) || (Window.resizeNeeded != 0) {
		windowScheduleRedrawTask(Window)
	}
}

//line 4521
func windowGetAllocation(Window *Window, allocation *rectangle) {
	*allocation = Window.mainSurface.allocation
}

//line 4445
func windowGetGeometry(Window *Window, geometry *rectangle) {
	if Window.fullscreen != 0 {
		windowGetAllocation(Window, geometry)
	}
}

//line 4458
func windowSyncGeometry(Window *Window) {
	var geometry rectangle

	if Window.xdgSurface == nil {
		return
	}

	windowGetGeometry(Window, &geometry)

	if geometry.x == Window.lastGeometry.x &&
		geometry.y == Window.lastGeometry.y &&
		geometry.width == Window.lastGeometry.width &&
		geometry.height == Window.lastGeometry.height {
		return
	}

	_ = Window.xdgSurface.SetWindowGeometry(
		geometry.x,
		geometry.y,
		geometry.width,
		geometry.height)
	Window.lastGeometry = geometry
}

// line 4480
func windowFlush(Window *Window) {

	if Window.redrawInhibited != 0 {
		panic("assert\n")
	}

	if Window.custom == 0 {
		if Window.xdgSurface != nil {
			windowSyncGeometry(Window)

		}

	}

	surfaceFlush(Window.mainSurface)

}

// line 4505
func widgetRedraw(Widget *Widget) {
	if Widget.Userdata != nil {
		Widget.Userdata.Redraw(Widget)
	}
}

//line 4517
func (s *surface) CallbackDone(callback *wl.Callback, time uint32) {
	wlclient.CallbackDestroy(callback)
	s.frameCb = nil

	s.lastTime = time

	if (s.redrawNeeded != 0) || (s.Window.redrawNeeded != 0) {

		windowScheduleRedrawTask(s.Window)
	}
}

//line 4545
func surfaceRedraw(surface *surface) int {

	if (surface.Window.redrawNeeded == 0) && (surface.redrawNeeded == 0) {
		return 0
	}

	// Whole-Window redraw forces a redraw even if the previous has
	// not yet hit the screen
	if nil != surface.frameCb {
		if surface.Window.redrawNeeded == 0 {
			return 0
		}

		wlclient.CallbackDestroy(surface.frameCb)
	}

	cb, err := surface.surface_.Frame()
	if err != nil {
		fmt.Println(err)
	} else {
		surface.frameCb = cb

		// add listener here
		wlclient.CallbackAddListener(surface.frameCb, surface)
	}

	surface.redrawNeeded = 0

	widgetRedraw(surface.Widget)

	return 0
}

// This is the alternative to idle_redraw
// line 4617
func (Window *Window) Run(events uint32) {

	Window.redrawTaskScheduled = 0

	if Window.resizeNeeded != 0 {
		if nil != Window.mainSurface.frameCb {
			return
		}

		idleResize(Window)

	}

	surfaceRedraw(Window.mainSurface)

	Window.redrawNeeded = 0
	windowFlush(Window)

}

//line 4619

func windowScheduleRedrawTask(Window *Window) {
	if Window.redrawInhibited != 0 {
		return
	}

	if Window.redrawTaskScheduled == 0 {

		Window.redrawRunner = Window
		displayDefer(Window.Display /*&Window.redraw_task,*/, Window)
		Window.redrawTaskScheduled = 1
	}
}

// line 4636
func windowScheduleRedraw(Window *Window) {
	windowScheduleRedrawTask(Window)
}

// line 4793
func (Window *Window) SetTitle(title string) {

	if Window.xdgToplevel != nil {
		_ = Window.xdgToplevel.SetTitle(title)
	}
}

// line 5178
func surfaceCreate(Window *Window) *surface {
	var Display = Window.Display
	var surface = &surface{}
	surface.Window = Window
	surf, err := Display.compositor.CreateSurface()
	if err != nil {
		panic(err.Error())
		return nil
	}
	surface.surface_ = surf

	surface.bufferScale = 1
	wlclient.SurfaceAddListener(surface.surface_, SurfaceEnter, SurfaceLeave)

	Window.subsurfaceListNew = append(Window.subsurfaceListNew, surface)

	return surface
}

// line 5219
func windowCreateInternal(Display *Display, custom int) *Window {

	var Window = &Window{}
	var surface_ *surface

	Window.Display = Display
	surface_ = surfaceCreate(Window)

	Window.mainSurface = surface_

	if (custom > 0) || (Display.xdgShell != nil) {
	} else {
		panic("assertion failed")
	}
	Window.custom = (int32)(custom)
	Window.preferredFormat = WindowPreferredFormatNone

	surface_.bufferType = WindowBufferTypeShm

	wlclient.SurfaceSetUserData(surface_.surface_, uintptr(0))
	Display.surface2window[surface_.surface_] = Window

	return Window
}

//line 5250
func Create(Display *Display) *Window {
	var Window = windowCreateInternal(Display, 0)

	if Window.Display.xdgShell != nil {
		surf, err := Window.Display.xdgShell.GetSurface(Window.mainSurface.surface_)
		if err != nil {
			fmt.Println(err)
			return nil
		} else {
			Window.xdgSurface = surf
		}

		Window.xdgSurface.AddListener(Window)

		tl, err := Window.xdgSurface.GetToplevel()
		if err != nil {
			fmt.Println(err)
			return nil
		} else {
			Window.xdgToplevel = tl
		}

		zxdg.ToplevelAddListener(Window.xdgToplevel, Window)

		windowInhibitRedraw(Window)

		_ = Window.mainSurface.surface_.Commit()
	}

	return Window
}

// line 5592
func (Window *Window) SetBufferType(t int32) {
	Window.mainSurface.bufferType = t
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

func (o *output) OutputGeometry(wlOutput *wl.Output, x int, y int, physicalWidth int,
	physicalHeight int, subpixel int, maker string, model string, transform int) {

	o.maker = maker
	o.model = model

}
func (o *output) OutputDone(wlOutput *wl.Output) {

}
func (o *output) OutputScale(wlOutput *wl.Output, factor int32) {

}

func (o *output) OutputMode(
	wlOutput *wl.Output,
	flags uint32,
	width int,
	height int,
	refresh int,
) {

}

// line 5771
func displayAddOutput(d *Display, id uint32) {

	var output = &output{}

	output.Display = d
	output.scale = 1
	output.output = wlclient.RegistryBindOutputInterface(d.registry, id, 2)

	output.serverOutputId = id

	wlclient.OutputAddListener(output.output, output)

}

//line 5925
func displayAddInput(d *Display, id uint32, displaySeatVersion int) {

	var input_ *Input
	var seatVersion = min(displaySeatVersion, 7)

	_ = seatVersion

	input_ = new(Input)

	input_.Display = d
	input_.seat = wlclient.RegistryBindSeatInterface(d.registry, id, uint32(seatVersion))
	input_.touchFocus = 0
	input_.pointerFocus = nil
	input_.keyboardFocus = nil
	input_.seatVersion = int32(seatVersion)

	d.inputList = append(d.inputList, input_)

	wlclient.SeatAddListener(input_.seat, input_)

	ps, err := d.compositor.CreateSurface()
	if err != nil {
		fmt.Println(err)
	} else {
		input_.pointerSurface = ps
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

	_ = createCursors(d)

	return d, nil
}

func (d *Display) BindUnstableInterface(name uint32, iface string, version uint32) wl.Proxy {
	return wlclient.RegistryBindUnstableInterface(d.registry, name, iface, version)
}

func (d *Display) SetUserData(data interface{}) {
	d.userData = data
}

//line 6387
func (d *Display) Destroy() {

	if d.dummySurface != nil {
		(*d.dummySurface).Destroy()
	}

	destroyCursors(d)

	if d.xdgShell != nil {
		d.xdgShell.Destroy()
	}

	if d.shm != nil {
		wlclient.ShmDestroy(d.shm)
	}

	wlclient.RegistryDestroy(d.registry)

	wlclient.DisplayDisconnect(d.Display)
}

//line 6478
func displayDefer(Display *Display /*task *task,*/, fun runner) {

	Display.deferredListNew = append(Display.deferredListNew, fun)
}

//line 6501
func DisplayRun(Display *Display) {

	Display.running = 1
	for {

		for len(Display.deferredListNew) > 0 {

			Display.deferredListNew[0].Run(0)

			Display.deferredListNew = Display.deferredListNew[1:]

		}

		if Display.running == 0 {
			break
		}

		if wlclient.DisplayRun(Display.Display) != nil {
			return
		}

	}
}

func (d *Display) Exit() {
	d.running = 0
}
