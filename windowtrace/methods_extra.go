package windowtrace

import (
	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/window"
	"github.com/neurlang/wayland/wl"
	zxdg "github.com/neurlang/wayland/xdg"
)

func NewWindowDecoration(win *Window) *WindowDecoration {
	println("func NewWindowDecoration")
	nd := window.NewWindowDecoration(win.Window)
	return (*WindowDecoration)(nd)
}

func (w *Window) SetFullscreenHandler(handler FullscreenHandler) {
	println("func SetFullscreenHandler")
	(w.Window).SetFullscreenHandler(fullscreenHandler{handler})
}

func (w *Widget) GetAllocation() Rectangle {
	println("func GetAllocation")
	return ((*window.Widget)(w)).GetAllocation()
}

func (w *Widget) SetAllocation(a int32, b int32, c int32, d int32) {
	println("func SetAllocation")
	((*window.Widget)(w)).SetAllocation(a, b, c, d)
}

func (w *Widget) WidgetGetLastTime() uint32 {
	println("func WidgetGetLastTime")
	return ((*window.Widget)(w)).WidgetGetLastTime()
}

func (w *Widget) ScheduleRedraw() {
	println("func ScheduleRedraw")
	((*window.Widget)(w)).ScheduleRedraw()
}

func (w *Widget) SetUserDataWidgetHandler(wh WidgetHandler) {
	println("func SetUserDataWidgetHandler")
	((*window.Widget)(w)).SetUserDataWidgetHandler(widgetHandler{wh})
}

func (w *Widget) Destroy() {
	println("func Destroy")
	((*window.Widget)(w)).Destroy()
}

func (w *Widget) ImageSurfaceGetData() []byte {
	println("func ImageSurfaceGetData")
	return ((*window.Widget)(w)).ImageSurfaceGetData()
}

func (w *Widget) ImageSurfaceGetWidth() int {
	println("func ImageSurfaceGetWidth")
	return ((*window.Widget)(w)).ImageSurfaceGetWidth()
}

func (w *Widget) ImageSurfaceGetHeight() int {
	println("func ImageSurfaceGetHeight")
	return ((*window.Widget)(w)).ImageSurfaceGetHeight()
}

func (w *Widget) ImageSurfaceGetStride() int {
	println("func ImageSurfaceGetStride")
	return ((*window.Widget)(w)).ImageSurfaceGetStride()
}

func (w *Widget) Reference() cairo.Surface {
	println("func Reference")
	return ((*window.Widget)(w)).Reference()
}

func (w *Widget) SetDestructor(f func()) {
	println("func SetDestructor")
	((*window.Widget)(w)).SetDestructor(f)
}

func (w *Widget) SetUserData(f func()) {
	println("func SetUserData")
	((*window.Widget)(w)).SetUserData(f)
}

// Input methods
func (input *Input) PointerEnter(wlPointer *wl.Pointer, serial uint32, surface *wl.Surface, sx float32, sy float32) {
	println("func PointerEnter")
	((*window.Input)(input)).PointerEnter(wlPointer, serial, surface, sx, sy)
}

func (input *Input) PointerLeave(wlPointer *wl.Pointer, serial uint32, wlSurface *wl.Surface) {
	println("func PointerLeave")
	((*window.Input)(input)).PointerLeave(wlPointer, serial, wlSurface)
}

func (input *Input) PointerMotion(wlPointer *wl.Pointer, time uint32, surfaceX float32, surfaceY float32) {
	println("func PointerMotion")
	((*window.Input)(input)).PointerMotion(wlPointer, time, surfaceX, surfaceY)
}

func (input *Input) PointerButton(wlPointer *wl.Pointer, serial uint32, time uint32, button uint32, stateW uint32) {
	println("func PointerButton")
	((*window.Input)(input)).PointerButton(wlPointer, serial, time, button, stateW)
}

func (input *Input) PointerFrame(wlPointer *wl.Pointer) {
	println("func PointerFrame")
	((*window.Input)(input)).PointerFrame(wlPointer)
}

func (input *Input) PointerAxis(wlPointer *wl.Pointer, time uint32, axis uint32, value float32) {
	println("func PointerAxis")
	((*window.Input)(input)).PointerAxis(wlPointer, time, axis, value)
}

func (input *Input) PointerAxisSource(wlPointer *wl.Pointer, axisSource uint32) {
	println("func PointerAxisSource")
	((*window.Input)(input)).PointerAxisSource(wlPointer, axisSource)
}

func (input *Input) PointerAxisStop(wlPointer *wl.Pointer, time uint32, axis uint32) {
	println("func PointerAxisStop")
	((*window.Input)(input)).PointerAxisStop(wlPointer, time, axis)
}

func (input *Input) PointerAxisDiscrete(wlPointer *wl.Pointer, axis uint32, discrete int32) {
	println("func PointerAxisDiscrete")
	((*window.Input)(input)).PointerAxisDiscrete(wlPointer, axis, discrete)
}

func (input *Input) GetModifiers() ModType {
	println("func GetModifiers")
	return (ModType)(((*window.Input)(input)).GetModifiers())
}

func (input *Input) GetRune(sym *uint32, v uint32) (r rune) {
	println("func GetRune")
	return ((*window.Input)(input)).GetRune(sym, v)
}

func (input *Input) GetUtf8() []byte {
	println("func GetUtf8")
	return ((*window.Input)(input)).GetUtf8()
}

func (input *Input) Destroy() {
	println("func Destroy")
	((*window.Input)(input)).Destroy()
}

func (input *Input) CallbackDone(wlCallback *wl.Callback, callbackData uint32) {
	println("func CallbackDone")
	((*window.Input)(input)).CallbackDone(wlCallback, callbackData)
}

func (input *Input) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	println("func HandleCallbackDone")
	((*window.Input)(input)).HandleCallbackDone(ev)
}

func (input *Input) HandleKeyboardKey(e wl.KeyboardKeyEvent) {
	println("func HandleKeyboardKey")
	((*window.Input)(input)).HandleKeyboardKey(e)
}

func (input *Input) HandleKeyboardKeymap(e wl.KeyboardKeymapEvent) {
	println("func HandleKeyboardKeymap")
	((*window.Input)(input)).HandleKeyboardKeymap(e)
}

func (input *Input) HandleKeyboardLeave(e wl.KeyboardLeaveEvent) {
	println("func HandleKeyboardLeave")
	((*window.Input)(input)).HandleKeyboardLeave(e)
}

func (input *Input) HandleKeyboardModifiers(e wl.KeyboardModifiersEvent) {
	println("func HandleKeyboardModifiers")
	((*window.Input)(input)).HandleKeyboardModifiers(e)
}

func (input *Input) HandleKeyboardRepeatInfo(e wl.KeyboardRepeatInfoEvent) {
	println("func HandleKeyboardRepeatInfo")
	((*window.Input)(input)).HandleKeyboardRepeatInfo(e)
}

func (input *Input) HandleKeyboardEnter(e wl.KeyboardEnterEvent) {
	println("func HandleKeyboardEnter")
	((*window.Input)(input)).HandleKeyboardEnter(e)
}

func (input *Input) HandleTouchCancel(e wl.TouchCancelEvent) {
	println("func HandleTouchCancel")
	((*window.Input)(input)).HandleTouchCancel(e)
}

func (input *Input) HandleTouchDown(e wl.TouchDownEvent) {
	println("func HandleTouchDown")
	((*window.Input)(input)).HandleTouchDown(e)
}

func (input *Input) HandleTouchFrame(e wl.TouchFrameEvent) {
	println("func HandleTouchFrame")
	((*window.Input)(input)).HandleTouchFrame(e)
}

func (input *Input) HandleTouchMotion(e wl.TouchMotionEvent) {
	println("func HandleTouchMotion")
	((*window.Input)(input)).HandleTouchMotion(e)
}

func (input *Input) HandleTouchOrientation(e wl.TouchOrientationEvent) {
	println("func HandleTouchOrientation")
	((*window.Input)(input)).HandleTouchOrientation(e)
}

func (input *Input) HandleTouchShape(e wl.TouchShapeEvent) {
	println("func HandleTouchShape")
	((*window.Input)(input)).HandleTouchShape(e)
}

func (input *Input) HandleTouchUp(e wl.TouchUpEvent) {
	println("func HandleTouchUp")
	((*window.Input)(input)).HandleTouchUp(e)
}

func (input *Input) SeatName(wlSeat *wl.Seat, name string) {
	println("func SeatName")
	((*window.Input)(input)).SeatName(wlSeat, name)
}

func (input *Input) SeatCapabilities(seat *wl.Seat, caps uint32) {
	println("func SeatCapabilities")
	((*window.Input)(input)).SeatCapabilities(seat, caps)
}

func (input *Input) HandleSeatName(ev wl.SeatNameEvent) {
	println("func HandleSeatName")
	((*window.Input)(input)).HandleSeatName(ev)
}

func (input *Input) HandleSeatCapabilities(ev wl.SeatCapabilitiesEvent) {
	println("func HandleSeatCapabilities")
	((*window.Input)(input)).HandleSeatCapabilities(ev)
}

func (input *Input) HandlePointerEnter(ev wl.PointerEnterEvent) {
	println("func HandlePointerEnter")
	((*window.Input)(input)).HandlePointerEnter(ev)
}

func (input *Input) HandlePointerLeave(ev wl.PointerLeaveEvent) {
	println("func HandlePointerLeave")
	((*window.Input)(input)).HandlePointerLeave(ev)
}

func (input *Input) HandlePointerMotion(ev wl.PointerMotionEvent) {
	println("func HandlePointerMotion")
	((*window.Input)(input)).HandlePointerMotion(ev)
}

func (input *Input) HandlePointerButton(ev wl.PointerButtonEvent) {
	println("func HandlePointerButton")
	((*window.Input)(input)).HandlePointerButton(ev)
}

func (input *Input) HandlePointerFrame(ev wl.PointerFrameEvent) {
	println("func HandlePointerFrame")
	((*window.Input)(input)).HandlePointerFrame(ev)
}

func (input *Input) HandlePointerAxis(ev wl.PointerAxisEvent) {
	println("func HandlePointerAxis")
	((*window.Input)(input)).HandlePointerAxis(ev)
}

func (input *Input) HandlePointerAxisSource(ev wl.PointerAxisSourceEvent) {
	println("func HandlePointerAxisSource")
	((*window.Input)(input)).HandlePointerAxisSource(ev)
}

func (input *Input) HandlePointerAxisStop(ev wl.PointerAxisStopEvent) {
	println("func HandlePointerAxisStop")
	((*window.Input)(input)).HandlePointerAxisStop(ev)
}

func (input *Input) HandlePointerAxisDiscrete(ev wl.PointerAxisDiscreteEvent) {
	println("func HandlePointerAxisDiscrete")
	((*window.Input)(input)).HandlePointerAxisDiscrete(ev)
}

func (input *Input) HandlePointerAxisValue120(ev wl.PointerAxisValue120Event) {
	println("func HandlePointerAxisValue120")
	((*window.Input)(input)).HandlePointerAxisValue120(ev)
}

func (input *Input) PointerAxisValue120(wlPointer *wl.Pointer, axis uint32, value120 int32) {
	println("func PointerAxisValue120")
	((*window.Input)(input)).PointerAxisValue120(wlPointer, axis, value120)
}

// Popup methods
func (p *Popup) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {
	println("func HandleSurfaceConfigure")
	(p.nested).HandleSurfaceConfigure(ev)
}

func (p *Popup) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	println("func HandleCallbackDone")
	(p.nested).HandleCallbackDone(ev)
}

func (p *Popup) HandleBufferRelease(ev wl.BufferReleaseEvent) {
	println("func HandleBufferRelease")
	(p.nested).HandleBufferRelease(ev)
}

// Window event handlers
func (w *Window) HandleToplevelClose(ev zxdg.ToplevelCloseEvent) {
	println("func HandleToplevelClose")
	(w.Window).HandleToplevelClose(ev)
}

func (w *Window) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {
	println("func HandleSurfaceConfigure")
	(w.Window).HandleSurfaceConfigure(ev)
}

func (w *Window) HandleToplevelConfigure(ev zxdg.ToplevelConfigureEvent) {
	println("func HandleToplevelConfigure")
	(w.Window).HandleToplevelConfigure(ev)
}

// WindowDecoration methods
func (d *WindowDecoration) SetHoverButton(btn ComponentType) {
	println("func SetHoverButton")
	((*window.WindowDecoration)(d)).SetHoverButton(window.ComponentType(btn))
}

func (d *WindowDecoration) HandlePointerEnter(serial uint32, x, y float32) {
	println("func HandlePointerEnter")
	((*window.WindowDecoration)(d)).HandlePointerEnter(serial, x, y)
}

func (d *WindowDecoration) HandlePointerLeave() {
	println("func HandlePointerLeave")
	((*window.WindowDecoration)(d)).HandlePointerLeave()
}

func (d *WindowDecoration) HandlePointerButton(serial uint32, button uint32, state wl.PointerButtonState) {
	println("func HandlePointerButton")
	((*window.WindowDecoration)(d)).HandlePointerButton(serial, button, state)
}

func (d *WindowDecoration) HandlePointerMotion(x, y float32) {
	println("func HandlePointerMotion")
	((*window.WindowDecoration)(d)).HandlePointerMotion(x, y)
}

func (d *WindowDecoration) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	println("func HandleCallbackDone")
	((*window.WindowDecoration)(d)).HandleCallbackDone(ev)
}

func (d *WindowDecoration) SetActive(active bool) {
	println("func SetActive")
	((*window.WindowDecoration)(d)).SetActive(active)
}

func (d *WindowDecoration) Show() error {
	println("func Show")
	return ((*window.WindowDecoration)(d)).Show()
}

func (d *WindowDecoration) Hide() {
	println("func Hide")
	((*window.WindowDecoration)(d)).Hide()
}

func (d *WindowDecoration) Destroy() {
	println("func Destroy")
	((*window.WindowDecoration)(d)).Destroy()
}

func (d *WindowDecoration) UpdateSize() {
	println("func UpdateSize")
	((*window.WindowDecoration)(d)).UpdateSize()
}

func (d *WindowDecoration) UpdateSizeForResize(contentWidth, contentHeight int32) {
	println("func UpdateSizeForResize")
	((*window.WindowDecoration)(d)).UpdateSizeForResize(contentWidth, contentHeight)
}


