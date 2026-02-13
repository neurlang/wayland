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

import "os"
import "io"

import sys "github.com/neurlang/wayland/os"
import xkb "github.com/neurlang/wayland/xkbcommon"
import "errors"

import "fmt"

type runner interface {
	Run(uint32)
}

const BufferTypeShm = 1

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

type global struct {
	name    uint32
	iface   string
	version uint32
}

type Display struct {
	Display               *wl.Display
	registry              *wl.Registry
	compositor            *wl.Compositor
	subcompositor         *wl.Subcompositor         //nolint:unused // Reserved for future use
	shm                   *wl.Shm
	dataDeviceManager     *wl.DataDeviceManager
	dataDeviceVersion     int                       //nolint:unused // Reserved for future use
	textCursorPosition    *struct{}                 //nolint:unused // Reserved for future use
	xdgShell              *zxdg.WmBase
	serial                uint32

	//display_fd        int32
	displayFdEvents       uint32                    //nolint:unused // Reserved for future use

	//display_task task
	//	pad4		uint64
	//	pad5		uint64
	//	pad6		uint64

	deferredList          [2]uintptr                //nolint:unused // Reserved for future use
	//	pad7		uint64
	//	pad8		uint64

	running               bool

	globalList            []*global
	//	pad9		uint64
	//	pada		uint64
	windowList            [2]*Window                //nolint:unused // Reserved for future use
	//	padb		uint64
	//	padc		uint64
	inputList             []*Input
	//	padd		uint64
	//	pade		uint64
	outputList            [2]*output                //nolint:unused // Reserved for future use
	//	padf		uint64
	//	padg		uint64

	theme                 *theme
	cursorTheme           *wlcursor.Theme
	cursors     *[lengthCursors]*wlcursor.Cursor

	xkbContext *xkb.Context

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

type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

type toysurface interface {
	prepare(dx int, dy int, width int32, height int32, flags uint32,
		bufferTransform uint32, bufferScale int32) cairo.Surface
	swap(bufferTransform uint32, bufferScale int32, serverAllocation *Rectangle)
	acquire(ctx *struct{}) int
	release()
	destroy()
}

type surface struct {
	Window *Window

	surface_            *wl.Surface
	subsurface          *wl.Subsurface
	synchronized        int32 //nolint:unused // Reserved for future use
	synchronizedDefault int32 //nolint:unused // Reserved for future use
	toysurface          *toysurface
	Widget              *Widget
	redrawNeeded        int32

	frameCb  *wl.Callback
	lastTime uint32
	//	pad1	uint32

	allocation       Rectangle
	serverAllocation Rectangle

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

type DataHandler func(*Window, *Input, float32, float32, []string, *Window, WidgetHandler)

type FullscreenHandler interface {
	Fullscreen(*Window, WidgetHandler)
}
type CloseHandler interface {
	Close()
}

const (
	TYPE_NONE byte = iota
	TYPE_TOPLEVEL
	TYPE_FULLSCREEN
	TYPE_MAXIMIZED
	TYPE_TRANSIENT
	TYPE_MENU
	TYPE_CUSTOM
)

type ResizeHandler interface {
	MinimumSize() (int32, int32)
}

type Popuper interface {
	Render(cairo.Surface, uint32)
	Done()
	Configure() *Widget
}

//line 1675
func (Window *Window) AddPopupWidget(p *Popup, data WidgetHandler) *Widget {
	var w = widgetCreate(Window, p.mainSurface, data)
	return w
}

type Popup struct {
	Popup    *zxdg.Popup
	xdgsurf  *zxdg.Surface
	Display  *Display
	rect     Rectangle
	surf     *wl.Surface
	callback *wl.Callback //nolint:unused // Reserved for future use

	popuper Popuper

	mainSurface *surface
}

func (w *Popup) Destroy() {
	_ = w.Popup.Destroy()
	_ = w.xdgsurf.Destroy()
	_ = w.surf.Destroy()

	for _, input := range w.Display.inputList {
		input.grab = nil
	}
}

func (w *Popup) HandleCallbackDone(ev wl.CallbackDoneEvent) {

	w.Destroy()
}

func (w *Popup) HandleBufferRelease(ev wl.BufferReleaseEvent) {
	w.BufferRelease(ev.B)
}

func (w *Popup) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {

	w.SurfaceConfigure(ev.Serial)
}
func (w *Popup) HandlePopupConfigure(ev zxdg.PopupConfigureEvent) {

	w.PopupConfigure(ev.X, ev.Y, ev.Width, ev.Height)
}
func (w *Popup) HandlePopupPopupDone(ev zxdg.PopupPopupDoneEvent) {

	w.PopupPopupDone()
}
func (w *Popup) BufferRelease(buffer *wl.Buffer) {

}
func (p *Popup) SetPopupHandler(ph Popuper) {
	p.popuper = ph
}

func (w *Popup) PopupConfigure(x, y, width, height int32) {
	//wind := (*Window2)(unsafe.Pointer(w.window))
	//xdgsurf := (*xdg.Surface)(unsafe.Pointer(wind.xdgSurface))

	println(x, y, width, height)

	w.rect.X = 0
	w.rect.Y = 0
	w.rect.Width = width
	w.rect.Height = height

	w.mainSurface = &surface{}
	w.mainSurface.allocation = Rectangle{x, y, width, height}

	widget := w.popuper.Configure()

	for _, input := range w.Display.inputList {
		input.grab = widget
	}

}
func (w *Popup) PopupPopupDone() {
	//wind := (*Window2)(unsafe.Pointer(w.window))
	//xdgsurf := (*xdg.Surface)(unsafe.Pointer(wind.xdgSurface))

	println("popup done")
	//w.Popup.Destroy()
	//w.xdgsurf.Destroy()
	//w.surf.Destroy()
	//surfaceDestroy(w.mainSurface)

	w.popuper.Done()

	for _, input := range w.Display.inputList {
		input.grab = nil
	}

}
func (w *Popup) SurfaceConfigure(serial uint32) {
	//wind := (*Window2)(unsafe.Pointer(w.window))
	//xdgsurf := (*xdg.Surface)(unsafe.Pointer(wind.xdgSurface))

	println("surface configure", serial)
	_ = w.xdgsurf.AckConfigure(serial)

	w.mainSurface = &surface{}

	w.mainSurface.surface_ = w.surf
	w.mainSurface.bufferScale = 1

	w.mainSurface.allocation = w.rect

	w.popupCreateSurface(w.mainSurface, 0)

	w.popuper.Render(w.mainSurface.cairoSurface, 0)

	surfaceFlush(w.mainSurface)
}

//line 1462
func (w *Popup) popupCreateSurface(surface *surface, flags uint32) {
	var Display = w.Display
	var allocation = surface.allocation

	if surface.toysurface == nil {
		var toy = shmSurfaceCreate(Display, surface.surface_, flags, &allocation)

		surface.toysurface = &toy
	}

	surface.cairoSurface = (*surface.toysurface).prepare(
		0, 0,
		allocation.Width, allocation.Height, flags,
		uint32(surface.bufferTransform), surface.bufferScale)

}

//line 2036
func (p *Popup) PopupGetSurface() cairo.Surface {
	if p == nil {
		return nil
	}
	var mainSurface = p.mainSurface
	if mainSurface == nil {
		return nil
	}
	var cairoSurface = mainSurface.cairoSurface
	if cairoSurface == nil {
		return nil
	}
	return cairoSurface.Reference()
}

func (w *Window) CreatePopup(seat *wl.Seat, clickSerial, width, height, x, y uint32) *Popup {
	disp := w.Display
	shell := disp.xdgShell

	p := &Popup{}

	surface := surfaceCreate(w)

	_ = shell
	_ = surface

	// Create a positioner
	positioner, err := shell.CreatePositioner()
	if err != nil {
		panic(err.Error())
	}

	_ = positioner.SetOffset(int32(x), int32(y))
	_ = positioner.SetSize(int32(width), int32(height))
	_ = positioner.SetAnchor(zxdg.PositionerAnchorTopLeft)

	// Set up the xdg_surface
	xdgSurface, err := shell.GetSurface(surface.surface_)
	if err != nil {
		panic(err.Error())

	}
	_ = xdgSurface
	xdgSurface.AddListener(p)

	p.xdgsurf = (*zxdg.Surface)((xdgSurface))

	// Get the popup
	popup, err := p.xdgsurf.GetPopup(w.xdgSurface, positioner)
	if err != nil {
		panic(err.Error())

	}

	_ = positioner.Destroy()

	popup.AddConfigureHandler(p)
	popup.AddPopupDoneHandler(p)

	_ = popup.Grab(seat, clickSerial)

	_ = surface.surface_.Commit()

	p.Popup = popup
	p.Display = disp
	p.surf = surface.surface_

	disp.surface2window[surface.surface_] = w

	return p
}

type Window struct {
	Display          *Display
	windowOutputList [2]uintptr //nolint:unused // Reserved for future use

	title string

	typ byte

	savedAllocation   Rectangle
	minAllocation     Rectangle
	pendingAllocation Rectangle
	lastGeometry      Rectangle

	x, y int32 //nolint:unused // Reserved for future use

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

	fullscreen bool
	maximized  bool

	preferredFormat int

	mainSurface *surface
	xdgSurface  *zxdg.Surface
	xdgToplevel *zxdg.Toplevel
	xdgPopup    *zxdg.Popup

	parent     *Window     //nolint:unused // Reserved for future use
	lastParent *Window     //nolint:unused // Reserved for future use

	/* struct surface::link, contains also mainSurface */
	subsurfaceList [2]*surface //nolint:unused // Reserved for future use

	pointerLocked bool //nolint:unused // Reserved for future use

	confined bool //nolint:unused // Reserved for future use

	link [2]*Window //nolint:unused // Reserved for future use

	Userdata WidgetHandler

	redrawRunner runner

	subsurfaceListNew []*surface

	keyboardHandler KeyboardHandler

	frame *windowFrame

	decoration *WindowDecoration
	decorationsRequested bool

	fullscreenHandler FullscreenHandler
	closeHandler      CloseHandler
	dataHandler       DataHandler

	fullscreenMethod uint32 //nolint:unused // Reserved for future use

	resizor ResizeHandler
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
func (Window *Window) SetFullscreenHandler(handler FullscreenHandler) {

	Window.fullscreenHandler = handler

}
func (Window *Window) SetCloseHandler(handler CloseHandler) {

	Window.closeHandler = handler

}

func (Window *Window) SetDataHandler(window *Window, handler DataHandler) {
	window.dataHandler = handler
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
	Window.maximized = false
	Window.fullscreen = false
	Window.resizing = 0
	Window.focused = 0

	for i := range states {
		switch states[i] {
		case zxdg.ToplevelStateMaximized:
			Window.maximized = true
		case zxdg.ToplevelStateFullscreen:
			Window.fullscreen = true
		case zxdg.ToplevelStateResizing:
			Window.resizing = 1
		case zxdg.ToplevelStateActivated:
			Window.focused = 1
		default:
			/* Unknown state */
		}
	}

	// better have this to be sure
	if (width < 0) || (height < 0) {
		return
	}

	if (width > 0) && (height > 0) {
		/* The width / height params are for Window geometry,
		 * but window_schedule_resize takes allocation. Add
		 * on the shadow margin to get the difference. */
		var margin int32 = 0

		Window.ScheduleResize(width+margin*2, height+margin*2)
	} else if (Window.savedAllocation.Width > 0) &&
		(Window.savedAllocation.Height > 0) {
		Window.ScheduleResize(Window.savedAllocation.Width, Window.savedAllocation.Height)
	}
	
	// Update decorations based on window state
	// Only after we have a valid size
	if Window.decoration != nil && Window.mainSurface != nil && 
		Window.mainSurface.allocation.Width > 0 && 
		Window.mainSurface.allocation.Height > 0 {
		
		// Update active state
		Window.decoration.SetActive(Window.focused == 1)
		
		// Show/hide decorations based on state
		if Window.fullscreen {
			// Hide decorations in fullscreen
			Window.decoration.Hide()
		} else if Window.decoration.shadowSurf == nil {
			// First time showing - create surfaces
			if err := Window.decoration.Show(); err != nil {
				// If decoration fails, disable it
				Window.decoration.Destroy()
				Window.decoration = nil
			}
		} else {
			// Already shown, just redraw
			Window.decoration.Redraw()
		}
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
	allocation Rectangle

	opaque        int32
	tooltipCount  int32
	defaultCursor int32

	/* If this is set to false then no cairo surface will be
	 * created before redrawing the surface. This is useful if the
	 * redraw handler is going to do completely custom rendering
	 * such as using EGL directly */
	useCairo int32

	userdata WidgetHandler
}

func (w *Widget) SetUserDataWidgetHandler(wh WidgetHandler) {
	w.userdata = wh
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
	Resize(Widget *Widget, width int32, height int32, pwidth int32, pheight int32)
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
	Axis(Widget *Widget, Input *Input, time uint32, axis uint32, value float32)
	AxisSource(Widget *Widget, Input *Input, source uint32)
	AxisStop(Widget *Widget, Input *Input, time uint32, axis uint32)
	AxisDiscrete(Widget *Widget, Input *Input, axis uint32, discrete int32)
	PointerFrame(Widget *Widget, Input *Input)
}

type Input struct {
	Display            *Display
	seat               *wl.Seat
	pointer            *wl.Pointer
	keyboard           *wl.Keyboard
	touch              *wl.Touch
	touchPointList     [2]uintptr //nolint:unused // Reserved for future use
	pointerFocus       *Window
	keyboardFocus      *Window
	touchFocus         int32
	currentCursor      int32
	cursorAnimStart    uint32
	cursorFrameCb      *wl.Callback
	cursorTimerStart   uint32      //nolint:unused // Reserved for future use
	cursorAnimCurrent  uint32
	cursorDelayFd      int32       //nolint:unused // Reserved for future use
	cursorTimerRunning bool        //nolint:unused // Reserved for future use
	//cursor_task          task
	pointerSurface     *wl.Surface
	modifiers          ModType
	pointerEnterSerial uint32
	cursorSerial       uint32
	sx                 float32
	sy                 float32

	focusWidget *Widget
	grab        *Widget
	grabButton  uint32

	dataDevice      *wl.DataDevice
	touchGrab       uint32  //nolint:unused // Reserved for future use
	touchGrabId     int32   //nolint:unused // Reserved for future use
	dragX           float32 //nolint:unused // Reserved for future use
	dragY           float32 //nolint:unused // Reserved for future use
	dragFocus       *Window //nolint:unused // Reserved for future use
	dragEnterSerial uint32  //nolint:unused // Reserved for future use

	xkb struct {
		keymap       *xkb.Keymap
		state        *xkb.State
		composeTable *xkb.ComposeTable
		composeState *xkb.ComposeState

		controlMask uint32
		altMask     uint32
		shiftMask   uint32
	}

	repeatRateSec   int32 //nolint:unused // Reserved for future use
	repeatRateNsec  int32 //nolint:unused // Reserved for future use
	repeatDelaySec  int32 //nolint:unused // Reserved for future use
	repeatDelayNsec int32 //nolint:unused // Reserved for future use

	//repeat_task     task
	repeatSym      uint32 //nolint:unused // Reserved for future use
	repeatKey      uint32 //nolint:unused // Reserved for future use
	repeatTime     uint32 //nolint:unused // Reserved for future use
	seatVersion    int32
	selectionOffer *dataOffer
	dragOffer      *dataOffer
	offerData      map[*wl.DataOffer]*dataOffer
}

func (input *Input) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	input.CallbackDone(ev.C, ev.CallbackData)
}

type dataOffer struct {
	input *Input
	offer *wl.DataOffer

	types []string
}

func (do *dataOffer) Destroy() {
	wlclient.DataOfferDestroy(do.offer)
	do.types = nil
	do.input = nil
}

func (do *dataOffer) HandleDataOfferOffer(ev wl.DataOfferOfferEvent) {

	println("HandleDataOfferOffer", ev.MimeType)

	do.types = append(do.types, ev.MimeType)
}

func (do *dataOffer) HandleDataOfferSourceActions(ev wl.DataOfferSourceActionsEvent) {
	println("HandleDataOfferSourceActions")

}

func (do *dataOffer) HandleDataOfferAction(ev wl.DataOfferActionEvent) {
	println("HandleDataOfferAction")
}

func (input *Input) HandleDataDeviceDataOffer(ev wl.DataDeviceDataOfferEvent) {

	println("HandleDataDeviceDataOffer")
	var offer dataOffer

	offer.input = input
	offer.offer = ev.Id

	wlclient.DataOfferAddListener(offer.offer, &offer)

	if input.offerData == nil {
		input.offerData = make(map[*wl.DataOffer]*dataOffer)
	}

	input.offerData[ev.Id] = &offer

	println("HandleDataDeviceDataOffer")
}

func (input *Input) HandleDataDeviceEnter(ev wl.DataDeviceEnterEvent) {
	println("HandleDataDeviceEnter")

	surface := ev.Surface

	if surface == nil {
		return
	}

	window, _ := wl.GetUserData[Window](surface)

	var typesData []string

	if ev.Id != nil {
		input.dragOffer = input.offerData[ev.Id]

		typesData = input.dragOffer.types

	}

	if window.dataHandler != nil {
		window.dataHandler(window, input, ev.X, ev.Y, typesData, window, window.Userdata)
	}
}

func (input *Input) HandleDataDeviceLeave(ev wl.DataDeviceLeaveEvent) {
	println("HandleDataDeviceLeave")

	if input.dragOffer != nil {
		_ = input.dragOffer.offer.Destroy()
		input.dragOffer.offer.Unregister()
	}

}

func (input *Input) HandleDataDeviceMotion(ev wl.DataDeviceMotionEvent) {
	println("HandleDataDeviceMotion")

}

func (input *Input) HandleDataDeviceDrop(ev wl.DataDeviceDropEvent) {
	println("HandleDataDeviceDrop")
}

func (input *Input) HandleDataDeviceSelection(ev wl.DataDeviceSelectionEvent) {

	println("HandleDataDeviceSelection")

	//println("HandleDataDeviceSelection", input.selectionOffer.offer, ev.Offer)

	if input.selectionOffer != nil {
		if input.selectionOffer.offer != nil {
			println("deleting")
			_ = input.selectionOffer.offer.Destroy()
			input.selectionOffer.offer.Unregister()
			input.selectionOffer.offer = nil
		}

	}

	if ev.Id != nil {

		var another = input.offerData[ev.Id]

		if another == input.selectionOffer {
			input.selectionOffer = nil
		} else {

			input.selectionOffer = another
		}
	} else {

		input.selectionOffer = nil
	}
	//println("HandleDataDeviceSelection", input.selectionOffer.offer, ev.Offer)
}

func (input *Input) DeviceSetSelection(src *DataSource, serial uint32) {
	if input.dataDevice != nil {
		_ = input.dataDevice.SetSelection(src.src, serial)
	}
}

func (input *Input) Destroy() {
	inputRemoveKeyboardFocus(input)
	inputRemovePointerFocus(input)

	if input.dragOffer != nil {
		input.dragOffer.Destroy()
		input.dragOffer = nil
	}
	if input.selectionOffer != nil {
		input.selectionOffer.Destroy()
		input.selectionOffer = nil
	}

	if input.dataDevice != nil {
		if input.Display.dataDeviceManagerVersion >= 2 {
			_ = input.dataDevice.Release()
		} else {
			wlclient.DataDeviceDestroy(input.dataDevice)
		}
		input.dataDevice = nil
	}

	if input.seatVersion >= wl.PointerReleaseSinceVersion {
		if input.touch != nil {
			_ = input.touch.Release()
		}
		if input.pointer != nil {
			_ = input.pointer.Release()
		}
		if input.keyboard != nil {
			_ = input.keyboard.Release()
		}
	} else {
		if input.touch != nil {
			wlclient.TouchDestroy(input.touch)
		}
		if input.pointer != nil {
			wlclient.PointerDestroy(input.pointer)
		}
		if input.keyboard != nil {
			wlclient.KeyboardDestroy(input.keyboard)
		}
	}
	input.touch = nil
	input.pointer = nil
	input.keyboard = nil
	_ = input.pointerSurface.Destroy()

	wlclient.SeatDestroy(input.seat)

}

type KeyboardHandler interface {
	Key(
		window *Window,
		input *Input,
		time uint32,
		key uint32,
		notUnicode uint32,
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
	allocation     Rectangle //nolint:unused // Reserved for future use
	link           [2]*output //nolint:unused // Reserved for future use
	transform      int32 //nolint:unused // Reserved for future use
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

func (input *Input) HandlePointerEnter(ev wl.PointerEnterEvent) {
	input.PointerEnter(nil, ev.Serial, ev.Surface, ev.SurfaceX, ev.SurfaceY)
}

func (input *Input) PointerEnter(
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

	var Window = input.Display.surface2window[surface]

	if Window == nil {
		return
	}

	// Check if this is a decoration surface
	if Window.decoration != nil {
		if Window.decoration.titleSurf != nil && surface == Window.decoration.titleSurf.wlSurface {
			input.Display.serial = serial
			input.pointerEnterSerial = serial
			input.pointerFocus = Window
			input.sx = sx
			input.sy = sy
			Window.decoration.HandlePointerEnter(serial, sx, sy)
			return
		}
		if Window.decoration.shadowSurf != nil && surface == Window.decoration.shadowSurf.wlSurface {
			// Shadow surface - ignore input
			return
		}
	}

	if surface != Window.mainSurface.surface_ {
		//		DBG("Ignoring Input event from subsurface %p\n", surface);
		return
	}

	input.Display.serial = serial
	input.pointerEnterSerial = serial
	input.pointerFocus = Window

	input.sx = sx
	input.sy = sy

}

func (input *Input) HandlePointerLeave(ev wl.PointerLeaveEvent) {
	input.PointerLeave(nil, ev.Serial, ev.Surface)
}

func (input *Input) PointerLeave(wlPointer *wl.Pointer, serial uint32, wlSurface *wl.Surface) {

	input.Display.serial = serial
	
	// Check if leaving a decoration surface
	if input.pointerFocus != nil && input.pointerFocus.decoration != nil {
		if input.pointerFocus.decoration.titleSurf != nil && 
			wlSurface == input.pointerFocus.decoration.titleSurf.wlSurface {
			input.pointerFocus.decoration.HandlePointerLeave()
		}
	}
	
	inputRemovePointerFocus(input)

}

func (input *Input) HandlePointerMotion(ev wl.PointerMotionEvent) {
	input.PointerMotion(ev.P, ev.Time, ev.SurfaceX, ev.SurfaceY)
}

func (input *Input) PointerMotion(
	wlPointer *wl.Pointer,
	time uint32,
	surfaceX float32,
	surfaceY float32,
) {

	pointerHandleMotion(input, wlPointer, time, surfaceX, surfaceY)
}

func (input *Input) HandlePointerButton(ev wl.PointerButtonEvent) {
	input.PointerButton(ev.P, ev.Serial, ev.Time, ev.Button, ev.State)
}

func (input *Input) PointerButton(
	wlPointer *wl.Pointer,
	serial uint32,
	time uint32,
	button uint32,
	stateW uint32,
) {
	var widget *Widget
	var state = wl.PointerButtonState(stateW)

	input.Display.serial = serial
	
	// Check if button event is on a decoration surface
	if input.pointerFocus != nil && input.pointerFocus.decoration != nil {
		if input.pointerFocus.decoration.hoverButton != ComponentNone {
			input.pointerFocus.decoration.HandlePointerButton(serial, button, state)
			return
		}
	}
	
	if input.focusWidget != nil && input.grab == nil &&
		state == wl.PointerButtonStatePressed {
		inputGrab(input, input.focusWidget, button)
	}

	widget = input.grab
	if widget != nil && widget.userdata != nil {
		widget.userdata.Button(widget,
			input, time,
			button, state,
			input.grab.userdata)
	}

	if input.grab != nil && input.grabButton == button &&
		state == wl.PointerButtonStateReleased {
		inputUngrab(input)
	}

}

func (input *Input) HandlePointerFrame(ev wl.PointerFrameEvent) {
}

func (input *Input) PointerFrame(wlPointer *wl.Pointer) {
}

func (input *Input) HandlePointerAxis(ev wl.PointerAxisEvent) {
	input.PointerAxis(nil, ev.Time, ev.Axis, ev.Value)
}

func (input *Input) PointerAxis(wlPointer *wl.Pointer, time uint32, axis uint32, value float32) {
	var Window = input.pointerFocus
	if Window == nil {
		return
	}

	var Widget *Widget
	if input.grab != nil {
		Widget = input.grab
	} else {
		Widget = input.focusWidget
	}
	if Widget == nil {
		return
	}
	if Widget.userdata != nil {
		Widget.userdata.Axis(Widget, input, time, axis, value)
	} else if Window.Userdata != nil {
		Window.Userdata.Axis(Widget, input, time, axis, value)
	}
}

func (input *Input) HandlePointerAxisSource(ev wl.PointerAxisSourceEvent) {
	input.PointerAxisSource(nil, ev.AxisSource)
}

func (input *Input) PointerAxisSource(wlPointer *wl.Pointer, axisSource uint32) {
	var Window = input.pointerFocus
	if Window == nil {
		return
	}
	var Widget *Widget
	if input.grab != nil {
		Widget = input.grab
	} else {
		Widget = input.focusWidget
	}
	if Widget == nil {
		return
	}
	if Widget.userdata != nil {
		Widget.userdata.AxisSource(Widget, input, axisSource)
	} else if Window.Userdata != nil {
		Window.Userdata.AxisSource(Widget, input, axisSource)
	}
}

func (input *Input) HandlePointerAxisStop(ev wl.PointerAxisStopEvent) {
	input.PointerAxisStop(nil, ev.Time, ev.Axis)
}

func (input *Input) PointerAxisStop(wlPointer *wl.Pointer, time uint32, axis uint32) {
	var Window = input.pointerFocus
	if Window == nil {
		return
	}
	var Widget *Widget
	if input.grab != nil {
		Widget = input.grab
	} else {
		Widget = input.focusWidget
	}
	if Widget == nil {
		return
	}
	if Widget.userdata != nil {
		Widget.userdata.AxisStop(Widget, input, time, axis)
	} else if Window.Userdata != nil {
		Window.Userdata.AxisStop(Widget, input, time, axis)
	}
}

func (input *Input) HandlePointerAxisDiscrete(ev wl.PointerAxisDiscreteEvent) {
	input.PointerAxisDiscrete(nil, ev.Axis, ev.Discrete)
}

func (input *Input) PointerAxisDiscrete(wlPointer *wl.Pointer, axis uint32, discrete int32) {
	var Window = input.pointerFocus
	if Window == nil {
		return
	}
	var Widget *Widget
	if input.grab != nil {
		Widget = input.grab
	} else {
		Widget = input.focusWidget
	}

	if Widget == nil {
		return
	}

	if Widget.userdata != nil {
		Widget.userdata.AxisDiscrete(Widget, input, axis, discrete)
	} else if Window.Userdata != nil {
		Window.Userdata.AxisDiscrete(Widget, input, axis, discrete)
	}
}

type SeatHandler interface {
	Capabilities(i *Input, seat *wl.Seat, caps uint32)
	Name(i *Input, seat *wl.Seat, name string)
}

func (input *Input) HandleSeatCapabilities(ev wl.SeatCapabilitiesEvent) {
	input.SeatCapabilities(input.seat, ev.Capabilities)
}

func (input *Input) HandleSeatName(ev wl.SeatNameEvent) {
	input.SeatName(input.seat, ev.Name)
}
func (input *Input) seatCapabilitiesPointer(seat *wl.Seat, caps uint32) {

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
}
func (input *Input) seatCapabilitiesKeyboard(seat *wl.Seat, caps uint32) {
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
}
func (input *Input) seatCapabilitiesTouch(seat *wl.Seat, caps uint32) {

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
			_ = input.touch.Release()
		} else {
			wlclient.TouchDestroy(input.touch)
		}
		input.touch = nil
	}
}

func (input *Input) SeatCapabilities(seat *wl.Seat, caps uint32) {

	input.seatCapabilitiesPointer(seat, caps)
	input.seatCapabilitiesKeyboard(seat, caps)
	input.seatCapabilitiesTouch(seat, caps)

	if input.Display.seatHandler != nil {
		input.Display.seatHandler.Capabilities(input, seat, caps)
	}

}

func (input *Input) HandleKeyboardEnter(e wl.KeyboardEnterEvent) {

	println("kbEnter")

	var window *Window
	var surface = e.Surface
	var serial = e.Serial

	if surface == nil {
		/* enter event for a window we've just destroyed */
		return
	}

	input.Display.serial = serial

	input.keyboardFocus, _ = wl.GetUserData[Window](surface)

	window = input.keyboardFocus
	if window != nil && window.keyboardHandler != nil {
		window.keyboardHandler.Focus(window, input)
	}

}

/* Translate symbols appropriately if a compose sequence is being entered */
func processKeyPress(sym uint32, input *Input) uint32 {
	if input.xkb.composeState == nil {
		return sym
	}
	if sym == xkb.KeyNoSymbol {
		return sym
	}
	if input.xkb.composeState.Feed(sym) != xkb.ComposeFeedAccepted {
		return sym
	}
	switch input.xkb.composeState.GetStatus() {
	case xkb.ComposeComposing:
		return xkb.KeyNoSymbol
	case xkb.ComposeComposed:
		return input.xkb.composeState.GetOneSym()
	case xkb.ComposeCancelled:
		return xkb.KeyNoSymbol
	case xkb.ComposeNothing:
		return sym
	default:
		return sym
	}
}

func (input *Input) GetModifiers() ModType {
	return input.modifiers
}

// This gets the UTF32 rune from the key sym ("notUnicode")
func (input *Input) GetRune(sym *uint32, _ uint32) (r rune) {
	r = rune(xkb.KeysymToUtf32(*sym))
	return
}

// This gets the compose UTF8 string from input
func (input *Input) GetUtf8() []byte {
	return input.xkb.composeState.GetUtf8()
}

func (input *Input) keyboardHandleKeyInternal(keyboard *wl.Keyboard,
	window *Window, sym uint32, state wl.KeyboardKeyState,
	time uint32, key uint32) {
	if sym == xkb.KeyF5 && input.modifiers == ModAltMask {
		if state == wl.KeyboardKeyStatePressed {
			_ = windowSetMaximized(window, !window.maximized)
		}
	} else if sym == xkb.KeyF11 &&
		window.fullscreenHandler != nil &&
		state == wl.KeyboardKeyStatePressed {
		window.fullscreenHandler.Fullscreen(window, window.Userdata)
	} else if sym == xkb.KeyF4 &&
		input.modifiers == ModAltMask &&
		state == wl.KeyboardKeyStatePressed {
		windowClose(window)
	} else if window.keyboardHandler != nil {
		if state == wl.KeyboardKeyStatePressed {
			sym = processKeyPress(sym, input)
		}

		window.keyboardHandler.Key(window, input, time, key,
			sym, state, window.Userdata)
	}
}
func (input *Input) keyboardHandleKey(keyboard *wl.Keyboard,
	serial uint32, time uint32, key uint32,
	stateW uint32) {
	var window = input.keyboardFocus
	var state = wl.KeyboardKeyState(stateW)

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

	var sym, _ = input.xkb.state.KeyGetSyms(code)

	input.keyboardHandleKeyInternal(keyboard, window, sym, state, time, key)

	// FIXME: Disarm repeat timer when state == wl.KeyboardKeyStateReleased && key == input.repeatKey
	// FIXME: Arm repeat timer when state == wl.KeyboardKeyStatePressed && input.xkb.keymap.KeyRepeats(code)
	_ = state
	_ = key
	_ = code
}

func (input *Input) HandleKeyboardKey(e wl.KeyboardKeyEvent) {
	input.keyboardHandleKey(nil, e.Serial, e.Time, e.Key, e.State)
}

func (input *Input) HandleKeyboardKeymap(e wl.KeyboardKeymapEvent) {

	var fd, err = e.Fd, e.FdError
	if err != nil {
		println(err.Error())
		return
	}

	var format = e.Format
	var size = e.Size

	var keymap *xkb.Keymap
	var state *xkb.State
	var composeTable *xkb.ComposeTable
	var composeState *xkb.ComposeState

	var locale string

	if input == nil {
		sys.Close(int(fd))
		return
	}

	if format != wl.KeyboardKeymapFormatXkbV1 {
		sys.Close(int(fd))
		return
	}

	mapStr, err := sys.Mmap(int(fd), 0, int(size), sys.ProtRead, sys.MapPrivate)
	if err != nil {
		sys.Close(int(fd))
		return
	}

	/* Set up XKB keymap */
	keymap = input.Display.xkbContext.KeymapNewFromString(mapStr,
		xkb.KeymapFormatTextV1,
		0)
	_ = sys.Munmap(mapStr)
	sys.Close(int(fd))

	if keymap == nil {
		println("failed to compile keymap")
		return
	}

	/* Set up XKB state */
	state = keymap.StateNew()
	if state == nil {
		println("failed to create XKB state")
		xkb.KeymapUnref(keymap)
		return
	}

	/* Look up the preferred locale, falling back to "C" as default */
	locale = os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LC_CTYPE")
		if locale == "" {
			locale = os.Getenv("LANG")
			if locale == "" {
				locale = "C"
			}
		}
	}

	/* Set up XKB compose table */
	composeTable =
		input.Display.xkbContext.ComposeTableNewFromLocale(locale,
			xkb.ComposeCompileNoFlags)
	if composeTable == nil {
		print("locale ")
		print(locale)
		println(": could not create XKB compose table for locale. Disabiling compose.")
	} else {
		/* Set up XKB compose state */
		composeState = xkb.ComposeStateNew(composeTable,
			xkb.ComposeStateNoFlags)
		if composeState == nil {
			println("could not create XKB compose state. Disabiling compose.")
			xkb.ComposeTableUnref(composeTable)

		} else {
			xkb.ComposeStateUnref(input.xkb.composeState)
			xkb.ComposeTableUnref(input.xkb.composeTable)
			input.xkb.composeState = composeState
			input.xkb.composeTable = composeTable
		}
	}

	xkb.KeymapUnref(input.xkb.keymap)
	xkb.StateUnref(input.xkb.state)
	input.xkb.keymap = keymap
	input.xkb.state = state

	input.xkb.controlMask =
		1 << input.xkb.keymap.ModGetIndex(xkb.ModNameCtrl)
	input.xkb.altMask =
		1 << input.xkb.keymap.ModGetIndex(xkb.ModNameAlt)
	input.xkb.shiftMask =
		1 << input.xkb.keymap.ModGetIndex(xkb.ModNameShift)

}
func (input *Input) HandleKeyboardLeave(e wl.KeyboardLeaveEvent) {
	var serial = e.Serial
	input.Display.serial = serial
	inputRemoveKeyboardFocus(input)

}
func (input *Input) HandleKeyboardModifiers(e wl.KeyboardModifiersEvent) {
	var mask uint32

	/* If we're not using a keymap, then we don't handle PC-style modifiers */
	if input.xkb.keymap == nil {
		return
	}

	input.xkb.state.UpdateMask(e.ModsDepressed, e.ModsLatched,
		e.ModsLocked, 0, 0, e.Group)

	mask = input.xkb.state.SerializeMods(xkb.StateModsDepressed | xkb.StateModsLatched)
	input.modifiers = 0
	if (mask & input.xkb.controlMask) != 0 {
		input.modifiers |= ModControlMask
	}
	if (mask & input.xkb.altMask) != 0 {
		input.modifiers |= ModAltMask
	}
	if (mask & input.xkb.shiftMask) != 0 {
		input.modifiers |= ModShiftMask
	}

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
func (input *Input) SeatName(wlSeat *wl.Seat, name string) {
	if input.Display.seatHandler != nil {
		input.Display.seatHandler.Name(input, wlSeat, name)
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
	_ = data.buffer.Destroy()
	if data.pool != nil {
		shmPoolDestroy(data.pool)
	}
}

//line 744
func makeShmPool(Display *Display, size uintptr, data *[]byte) (pool *wl.ShmPool) {
	fd, err := sys.CreateAnonymousFile(int64(size))
	if err != nil {
		println("creating a buffer file failed")
		println(size)
		println(err.Error())
		return nil
	}

	*data, err = sys.Mmap(int(fd.Fd()), 0, int(size), sys.ProtRead|sys.ProtWrite, sys.MapShared)
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

	err := sys.Munmap(pool.data)
	if err != nil {
		println(err)
	}
	if pool.pool != nil {
		_ = pool.pool.Destroy()
	}
	pool.data = nil
	pool.pool = nil
	pool.size = 0
	pool.used = 0
}

//line 820
func dataLengthForShmSurface(rect *Rectangle) uintptr {
	var stride = int32(cairo.FormatStrideForWidth(cairo.FormatArgb32, int(rect.Width)))
	return uintptr(int(stride * rect.Height))
}

func shmPoolReset(pool *shmPool) {
	pool.used = 0
}

//line 829
func displayCreateShmSurfaceFromPool(Display *Display,
	rectangle *Rectangle,
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

	stride = cairo.FormatStrideForWidth(cairoFormat, int(rectangle.Width))

	length = stride * int(rectangle.Height)
	data.pool = nil

	map_ = shmPoolAllocate(pool, uintptr(length), &offset)

	if map_ == nil {
		return nil, nil
	}

	surface = cairo.ImageSurfaceCreateForData(map_,
		cairoFormat,
		int(rectangle.Width),
		int(rectangle.Height),
		stride)

	surface.SetDestructor(func() {

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
		rectangle.Width,
		rectangle.Height,
		int32(stride), format)
	if err != nil {
		return nil, nil
	}

	return &surface, data
}

//line 886
func displayCreateShmSurface(Display *Display,
	rectangle *Rectangle, flags uint32,
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
	_ = shmSurfaceBufferRelease(s, buf)

}

func (s *shmSurface) prepare(dx int, dy int, width int32, height int32, flags uint32,
	bufferTransform uint32, bufferScale int32) cairo.Surface {

	var resizeHint = (flags & SurfaceHintResize) != 0
	surface := s
	var rect Rectangle
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

	rect.Width = width
	rect.Height = height

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
	bufferScale int32, serverAllocation *Rectangle) {
	var leaf = surface.current

	serverAllocation.Width =
		int32((*leaf.cairoSurface).ImageSurfaceGetWidth())
	serverAllocation.Height =
		int32((*leaf.cairoSurface).ImageSurfaceGetHeight())

	bufferToSurfaceSize(bufferTransform, bufferScale,
		&serverAllocation.Width,
		&serverAllocation.Height)

	_ = surface.surface.Attach(leaf.data.buffer,
		surface.dx, surface.dy)
	_ = surface.surface.Damage(0, 0,
		serverAllocation.Width, serverAllocation.Height)
	_ = surface.surface.Commit()

	leaf.busy = 1
	surface.current = nil
}

func (s *shmSurface) swap(
	bufferTransform uint32,
	bufferScale int32,
	serverAllocation *Rectangle,
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
	flags uint32, rectangle *Rectangle) toysurface {
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
		println(err.Error())
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
	_ = Display.cursorTheme.Destroy()
	Display.cursorTheme = nil
}

//line 1402
func surfaceFlush(surface *surface) {
	if surface.cairoSurface == nil {
		return
	}

	if surface.opaqueRegion != nil {
		_ = surface.surface_.SetOpaqueRegion(surface.opaqueRegion)
		_ = surface.opaqueRegion.Destroy()
		surface.opaqueRegion = nil
	}

	if surface.inputRegion != nil {
		_ = surface.surface_.SetInputRegion(surface.inputRegion)
		_ = surface.inputRegion.Destroy()
		surface.inputRegion = nil
	}

	(*surface.toysurface).swap(uint32(surface.bufferTransform), surface.bufferScale,
		&surface.serverAllocation)

	surface.cairoSurface.Destroy()
	surface.cairoSurface = nil
}

func windowClose(window *Window) {
	if window.closeHandler != nil {
		window.closeHandler.Close()
	} else {
		window.Display.Exit()
	}
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
		allocation.Width, allocation.Height, flags,
		uint32(surface.bufferTransform), surface.bufferScale)

}

//line 1488
func windowCreateMainSurface(Window *Window) {
	var surface = Window.mainSurface
	var flags uint32 = 0

	if Window.resizing != 0 {
		flags |= SurfaceHintResize
	}

	if Window.preferredFormat == PreferredFormatRgb565 {
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

	_ = surface.surface_.Destroy()

	if surface.toysurface != nil {
		(*surface.toysurface).destroy()
	}

}

//line 1577
func (Window *Window) Destroy() {
	
	// Clean up decorations first
	if Window.decoration != nil {
		Window.decoration.Destroy()
		Window.decoration = nil
	}

	if Window.xdgToplevel != nil {
		_ = Window.xdgToplevel.Destroy()
	}
	if Window.xdgPopup != nil {
		_ = Window.xdgPopup.Destroy()
	}
	if Window.xdgSurface != nil {
		_ = Window.xdgSurface.Destroy()
	}

	surfaceDestroy(Window.mainSurface)

}

//line 1624
func widgetFindWidget(Widget *Widget, x int32, y int32) *Widget {

	if Widget.allocation.X <= x &&
		x < Widget.allocation.X+Widget.allocation.Width &&
		Widget.allocation.Y <= y &&
		y < Widget.allocation.Y+Widget.allocation.Height {
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
	w.userdata = data
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

func (widget *Widget) GetAllocation() Rectangle {
	return widget.allocation
}

func (d *Display) SetSeatHandler(h SeatHandler) {
	d.seatHandler = h
}

func (d *Display) HandleWmBasePing(ev zxdg.WmBasePingEvent) {
	d.ShellPing(d.xdgShell, ev.Serial)
}

func (d *Display) ShellPing(shell *zxdg.WmBase, serial uint32) {
	_ = shell.Pong(serial)
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
		d.compositor = wlclient.RegistryBindCompositorInterface(d.registry, id, 3)

	case "wl_output":

		displayAddOutput(d, id)
		// TODO
	case "wl_seat":

		displayAddInput(d, id, int(version))

	case "wl_shm":
		d.shm = wlclient.RegistryBindShmInterface(d.registry, id, 1)
		wlclient.ShmAddListener(d.shm, d)
	case "wl_data_device_manager":
		displayAddDataDevice(d, id, version)

	//case "zxdg_shell_v6":
	case "xdg_wm_base":

		d.xdgShell = wlclient.RegistryBindWmBaseInterface(d.registry, id, 1)

		zxdg.WmBaseAddListener(d.xdgShell, d)

	case "text_cursor_position":
	case "wl_subcompositor":
		d.subcompositor = wlclient.RegistryBindSubcompositorInterface(d.registry, id, 1)

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
	Widget.allocation.Width = width
	Widget.allocation.Height = height
}

//line 1740
func (Widget *Widget) SetAllocation(x int32, y int32, width int32, height int32) {
	Widget.allocation.X = x
	Widget.allocation.Y = y
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
func (parent *Widget) ScheduleRedraw() {
	parent.surface.redrawNeeded = 1
	windowScheduleRedrawTask(parent.Window)
}

//line 2036
func (Window *Window) WindowGetSurface() cairo.Surface {
	var cairoSurface = widgetGetCairoSurface(Window.mainSurface.Widget)
	if cairoSurface == nil {
		return nil
	}
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

	frame.widget.userdata = data

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
		goto end
	}

	Input.sx = sx
	Input.sy = sy
	
	// Check if pointer is over a decoration surface
	if Window.decoration != nil && Window.decoration.titleSurf != nil {
		// If we're tracking decoration input, route motion there
		if Window.decoration.hoverButton != ComponentNone || Window.decoration.isDragging {
			Window.decoration.HandlePointerMotion(sx, sy)
			return
		}
	}

	// when making the Window smaller - e.g. after an unmaximise we might
	// * still have a pending motion event that the compositor has picked
	// * based on the old surface dimensions. However, if we have an active
	// * grab, we expect to see Input from outside the Window anyway.

	if nil == Input.grab && (sx < float32(Window.mainSurface.allocation.X) ||
		sy < float32(Window.mainSurface.allocation.Y) ||
		sx > float32(Window.mainSurface.allocation.Width) ||
		sy > float32(Window.mainSurface.allocation.Height)) {
		return
	}

	if !(Input.grab != nil && Input.grabButton != 0) {
		Widget = windowFindWidget(Window, int32(sx), int32(sy))
		inputSetFocusWidget(Input, Widget, sx, sy)

	}
end:

	if Input.grab != nil {
		Widget = Input.grab
	} else {
		Widget = Input.focusWidget
	}
	if Widget != nil {
		if Widget.userdata != nil {
			cursor = Widget.userdata.Motion(Input.focusWidget,
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

func (window *Window) SetFullscreen(fullscreen bool) error {
	if window.Display.xdgShell == nil {
		return errors.New("no xdg shell")
	}

	if (window.typ == TYPE_FULLSCREEN) == fullscreen {
		return errors.New("bad typ")
	}

	if fullscreen {
		window.typ = TYPE_FULLSCREEN
		return window.xdgToplevel.SetFullscreen(nil)
	} else {
		window.typ = TYPE_TOPLEVEL
		_ = window.xdgToplevel.UnsetFullscreen()
		window.ScheduleResize(window.savedAllocation.Width,
			window.savedAllocation.Height)
	}
	return nil
}

//line 3754
func inputSetPointerImageIndex(Input *Input, index int) {
	if Input.pointer == nil {
		print("input has no pointer\n")
		return
	}

	var cursor = Input.Display.cursors[Input.currentCursor]
	if cursor == nil {
		print("current cursor index out of range\n")
		return
	}

	var image = cursor.GetCursorImage(index)
	if image == nil {
		print("cursor index out of range\n")
		return
	}

	var buffer = image.GetBuffer()
	if buffer == nil {
		print("cursor buffer is nil\n")
		return
	}

	_ = Input.pointerSurface.Attach(buffer, 0, 0)
	_ = Input.pointerSurface.Damage(0, 0,
		int32(image.GetWidth()), int32(image.GetHeight()))
	_ = Input.pointerSurface.Commit()
	_ = wlcursor.PointerSetCursor(Input.pointer, Input.pointerEnterSerial, Input.pointerSurface,
		int32(image.GetHotspotX()), int32(image.GetHotspotY()))

}

//line 3789
func inputSetPointerSpecial(Input *Input) bool {
	if Input.currentCursor == CursorBlank {
		_ = wlcursor.PointerSetCursor(Input.pointer,
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

func (input *Input) CallbackDone(wlCallback *wl.Callback, callbackData uint32) {
	pointerSurfaceFrameCallback(input, wlCallback, callbackData)
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
		print("input has no pointer\n")
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

//line 4020
func (of *dataOffer) ReceiveData(mimeType string, function io.WriteCloser) error {
	var f1, f2, err = os.Pipe()
	if err != nil {
		return err
	}

	_ = of.offer.Receive(mimeType, f2.Fd())
	f2.Close()

	go func(f *os.File) {
		for {
			n, err := io.Copy(function, f)
			if n == 0 {
				function.Close()
				f.Close()
				return
			}
			if err != nil {
				function.Close()
				f.Close()
				return
			}
		}
	}(f1)

	return nil
}

//line 4062
func (input *Input) ReceiveSelectionData(mimeType string, function io.WriteCloser) error {

	if input.selectionOffer == nil {
		return errors.New("no offer")
	}
	if function == nil {
		return errors.New("nil function")
	}
	var found bool
	for _, p := range input.selectionOffer.types {
		if mimeType == p {
			found = true
			break
		}
	}
	if !found {
		return errors.New("not found")
	}
	return input.selectionOffer.ReceiveData(mimeType, function)
}

// line 4104
func surfaceResize(surface *surface) {
	var Widget = surface.Widget

	if Widget.userdata != nil {
		Widget.userdata.Resize(Widget,
			Widget.allocation.Width,
			Widget.allocation.Height,
			Widget.Window.pendingAllocation.Width,
			Widget.Window.pendingAllocation.Height)
	}

	if (surface.allocation.Width != Widget.allocation.Width) ||
		(surface.allocation.Height != Widget.allocation.Height) {
		Widget.Window.ScheduleRedraw()

	}

	surface.allocation = Widget.allocation

}

//line 4144
func windowDoResize(Window *Window) {

	Window.mainSurface.Widget.SetAllocation(
		Window.pendingAllocation.X,
		Window.pendingAllocation.Y,
		Window.pendingAllocation.Width,
		Window.pendingAllocation.Height)

	surfaceResize(Window.mainSurface)
	
	// Create or update decorations
	if Window.Display.subcompositor != nil && !Window.fullscreen &&
		Window.pendingAllocation.Width > 0 && Window.pendingAllocation.Height > 0 {

		if Window.decorationsRequested && Window.decoration == nil {
			Window.decorationsRequested = false
			Window.decoration = NewWindowDecoration(Window)
			if err := Window.decoration.Show(); err != nil {
				Window.decoration = nil
			} else {
				Window.decoration.SetActive(Window.focused == 1)
			}
		} else if Window.decoration != nil && Window.decoration.shadowSurf != nil {
			Window.decoration.Destroy()
			Window.decoration = NewWindowDecoration(Window)
			if err := Window.decoration.Show(); err != nil {
				Window.decoration = nil
			} else {
				Window.decoration.SetActive(Window.focused == 1)
			}
		}
	}
	if (!Window.fullscreen) && (!Window.maximized) {
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
	// this is an explicit change from upstream wayland/weston
	var minWidth int32 = 32
	var minHeight int32 = 32

	if Window.resizor != nil {
		minWidth, minHeight = Window.resizor.MinimumSize()
	}

	Window.pendingAllocation.X = 0
	Window.pendingAllocation.Y = 0
	Window.pendingAllocation.Width = width
	Window.pendingAllocation.Height = height

	// this is an explicit change from upstream wayland/weston
	if Window.minAllocation.Width < minWidth {
		Window.minAllocation.Width = minWidth
	}
	// this is an explicit change from upstream wayland/weston
	if Window.minAllocation.Width < minHeight {
		Window.minAllocation.Height = minHeight
	}

	if Window.pendingAllocation.Width < Window.minAllocation.Width {
		Window.pendingAllocation.Width = Window.minAllocation.Width
	}
	if Window.pendingAllocation.Height < Window.minAllocation.Height {
		Window.pendingAllocation.Height = Window.minAllocation.Height
	}

	Window.resizeNeeded = 1
	Window.ScheduleRedraw()
}

//line 4254
func (parent *Widget) ScheduleResize(width int32, height int32) {
	parent.Window.ScheduleResize(width, height)
}

//line 4269
func (Window *Window) InhibitRedraw() {
	Window.redrawInhibited = 1
	Window.redrawTaskScheduled = 0
}

func (Window *Window) UninhibitRedraw() {
	windowUninhibitRedraw(Window)
	windowScheduleRedrawTask(Window)
}

// line 4284
func windowUninhibitRedraw(Window *Window) {
	Window.redrawInhibited = 0
	if (Window.redrawNeeded != 0) || (Window.resizeNeeded != 0) {
		windowScheduleRedrawTask(Window)
	}
}

//line 4521
func windowGetAllocation(Window *Window, allocation *Rectangle) {
	*allocation = Window.mainSurface.allocation
}

//line 4445
func windowGetGeometry(Window *Window, geometry *Rectangle) {
	windowGetAllocation(Window, geometry)
}

//line 4458
func windowSyncGeometry(Window *Window) {
	var geometry Rectangle

	if Window.xdgSurface == nil {
		return
	}

	windowGetGeometry(Window, &geometry)

	if geometry.X == Window.lastGeometry.X &&
		geometry.Y == Window.lastGeometry.Y &&
		geometry.Width == Window.lastGeometry.Width &&
		geometry.Height == Window.lastGeometry.Height {
		return
	}
	_ = Window.xdgSurface.SetWindowGeometry(
		geometry.X,
		geometry.Y,
		geometry.Width,
		geometry.Height)
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
	if Widget.userdata != nil {
		Widget.userdata.Redraw(Widget)
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
func (Window *Window) ScheduleRedraw() {
	windowScheduleRedrawTask(Window)
}

func (Window *Window) ToggleMaximized() error {
	// extra feature: un-fullscreen using toggle maximized button if fullscreen
	if (Window.typ == TYPE_FULLSCREEN) && Window.fullscreen {
		_ = Window.SetFullscreen(false)
		return nil
	}
	return windowSetMaximized(Window, !Window.maximized)
}
func (Window *Window) SetMaximized(maximized bool) error {
	return windowSetMaximized(Window, maximized)
}

func (Window *Window) SetMinimized() error {
	return windowSetMinimized(Window)
}

func windowSetMaximized(window *Window, maximized bool) error {
	if window.xdgToplevel == nil {
		return errors.New("no_toplevel")
	}

	if window.maximized == maximized {
		return errors.New("already_set")
	}

	if maximized {
		window.savedAllocation = window.mainSurface.allocation
		return window.xdgToplevel.SetMaximized()
	} else {
		return window.xdgToplevel.UnsetMaximized()
	}
}

func windowSetMinimized(window *Window) error {
	if window.xdgToplevel == nil {
		return errors.New("no_toplevel")
	}

	return window.xdgToplevel.SetMinimized()
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
	Window.preferredFormat = PreferredFormatNone

	surface_.bufferType = BufferTypeShm

	wlclient.SurfaceSetUserData(surface_.surface_, Window)
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

		Window.InhibitRedraw()

		_ = Window.mainSurface.surface_.Commit()
		
		// Decorations will be created after first configure/resize
		// when we know the window dimensions
		Window.decorationsRequested = true
	}

	return Window
}

// CreateUndecorated creates a window without client-side decorations
func CreateUndecorated(Display *Display) *Window {
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

		Window.InhibitRedraw()

		_ = Window.mainSurface.surface_.Commit()
		
		// Explicitly disable decorations
		Window.decorationsRequested = false
	}

	return Window
}

// line 5592
func (Window *Window) SetBufferType(t int32) {
	Window.mainSurface.bufferType = t
}

// enableDecorations creates and shows window decorations (private)
func (Window *Window) enableDecorations() error {
	if Window.Display.subcompositor == nil {
		return fmt.Errorf("subcompositor not available")
	}
	if Window.decoration != nil {
		return nil
	}

	var width, height int32
	if Window.mainSurface != nil && Window.mainSurface.allocation.Width > 0 {
		width = Window.mainSurface.allocation.Width
		height = Window.mainSurface.allocation.Height
	} else if Window.pendingAllocation.Width > 0 {
		width = Window.pendingAllocation.Width
		height = Window.pendingAllocation.Height
	} else {
		return fmt.Errorf("window not yet configured")
	}
	_ = width
	_ = height

	Window.decoration = NewWindowDecoration(Window)
	if err := Window.decoration.Show(); err != nil {
		Window.decoration = nil
		return err
	}
	Window.decoration.SetActive(Window.focused == 1)
	return nil
}

// disableDecorations removes window decorations (private)
func (Window *Window) disableDecorations() {
	if Window.decoration != nil {
		Window.decoration.Destroy()
		Window.decoration = nil
	}
}

func minInt(a, b int) int {
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
	var seatVersion = minInt(displaySeatVersion, 7)

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

	if d.dataDeviceManager != nil {
		dev, err := d.dataDeviceManager.GetDataDevice(input_.seat)
		if err != nil {
			fmt.Println(err)
		} else if dev != nil {
			wlclient.DataDeviceAddListener(dev, input_)

			input_.dataDevice = dev
		}

	}

	ps, err := d.compositor.CreateSurface()
	if err != nil {
		fmt.Println(err)
	} else {
		input_.pointerSurface = ps
	}
}

func displayAddDataDevice(d *Display, id uint32, ddmVersion uint32) {
	d.dataDeviceManagerVersion = minU32(ddmVersion, 3)

	d.dataDeviceManager = wlclient.RegistryBindDataDeviceManagerInterface(d.registry, id,
		d.dataDeviceManagerVersion)

	for _, input := range d.inputList {
		if input.dataDevice == nil {

			dev, err := d.dataDeviceManager.GetDataDevice(input.seat)
			if err != nil {
				fmt.Println(err)
			} else if dev != nil {
				wlclient.DataDeviceAddListener(dev, input)

				input.dataDevice = dev
			}

		}
	}
}

// line 6237
func DisplayCreate(argv []string) (d *Display, e error) {

	d = &Display{}

	d.Display, e = wlclient.DisplayConnect(nil)
	if e != nil {
		return nil, fmt.Errorf("failed to connect to Wayland Display: %w", e)
	}

	d.xkbContext = xkb.ContextNew(xkb.ContextNoFlags)
	if d.xkbContext == nil {
		return nil, fmt.Errorf("failed to create XKB context: %w", e)
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
		_ = d.xdgShell.Destroy()
	}

	if d.shm != nil {
		wlclient.ShmDestroy(d.shm)
	}

	if d.dataDeviceManager != nil {
		wlclient.DataDeviceManagerDestroy(d.dataDeviceManager)
	}

	wlclient.RegistryDestroy(d.registry)

	wlclient.DisplayDisconnect(d.Display)
}

func (d *Display) GetSerial() uint32 {
	return d.serial
}

//line 6425
func (d *Display) CreateDataSource() (*DataSource, error) {
	if d.dataDeviceManager == nil {
		return nil, errors.New("Device manager does not exist")
	}
	ds, err := d.dataDeviceManager.CreateDataSource()

	return &DataSource{ds, ""}, err
}

//line 6478
func displayDefer(Display *Display /*task *task,*/, fun runner) {

	Display.deferredListNew = append(Display.deferredListNew, fun)
}

//line 6501
func DisplayRun(Display *Display) {

	Display.running = true
	for {

		for len(Display.deferredListNew) > 0 {

			Display.deferredListNew[0].Run(0)

			Display.deferredListNew = Display.deferredListNew[1:]

		}

		if !Display.running {
			break
		}

		if err := wlclient.DisplayRun(Display.Display); err != nil {
			fmt.Println(err)
			return
		}

	}
}

func (d *Display) Exit() {
	d.running = false
}
