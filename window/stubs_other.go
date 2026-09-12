//go:build windows || js
// +build windows js

package window

import (
	"github.com/neurlang/wayland/wl"
	zxdg "github.com/neurlang/wayland/xdg"
)

type ComponentType int
type DecorationSurface struct{}
type WindowDecoration struct{}

type AxisValue120Handler interface {
	AxisValue120(Widget *Widget, Input *Input, axis uint32, value120 int32)
}

const (
	ComponentNone ComponentType = iota
	ComponentShadow
	ComponentTitle
	ComponentButtonMin
	ComponentButtonMax
	ComponentButtonClose
)

const CursorDefault = 100
const CursorUnset = 101
const MaxLeaves = 3

const ShadowMargin = 8
const TitleHeight = 24
const ButtonWidth = 32
const SymDim = 14
const ShadowBlurSize = 64

const (
	TYPE_NONE byte = iota
	TYPE_TOPLEVEL
	TYPE_FULLSCREEN
	TYPE_MAXIMIZED
	TYPE_TRANSIENT
	TYPE_MENU
	TYPE_CUSTOM
)

func NewWindowDecoration(window *Window) *WindowDecoration { return nil }

func (d *WindowDecoration) HandleCallbackDone(ev wl.CallbackDoneEvent) {}
func (d *WindowDecoration) HandlePointerEnter(serial uint32, x, y float32) {}
func (d *WindowDecoration) HandlePointerLeave() {}
func (d *WindowDecoration) HandlePointerButton(serial uint32, button uint32, state wl.PointerButtonState) {}
func (d *WindowDecoration) HandlePointerMotion(x, y float32) {}
func (d *WindowDecoration) SetActive(active bool) {}
func (d *WindowDecoration) Show() error { return nil }
func (d *WindowDecoration) Hide() {}
func (d *WindowDecoration) Destroy() {}
func (d *WindowDecoration) UpdateSize() {}
func (d *WindowDecoration) UpdateSizeForResize(contentWidth, contentHeight int32) {}
func (d *WindowDecoration) SetHoverButton(btn ComponentType) {}

func (w *Window) ScheduleRedraw() {}
func (w *Window) ToplevelClose(t *zxdg.Toplevel) {}
func (w *Window) ToplevelConfigure(t *zxdg.Toplevel, width int32, height int32, states []int32) {}
func (w *Window) SurfaceConfigure(s *zxdg.Surface, serial uint32) {}
func (w *Window) Run(events uint32) {}
func (w *Window) SetCloseHandler(handler CloseHandler) {}
func (w *Window) SetDataHandler(win *Window, handler DataHandler) {}
func (w *Window) HandleToplevelClose(ev zxdg.ToplevelCloseEvent) {}
func (w *Window) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {}
func (w *Window) HandleToplevelConfigure(ev zxdg.ToplevelConfigureEvent) {}

func (p *Popup) HandlePopupConfigure(ev zxdg.PopupConfigureEvent) {}
func (p *Popup) HandlePopupPopupDone(ev zxdg.PopupPopupDoneEvent) {}
func (p *Popup) PopupConfigure(x, y, width, height int32) {}
func (p *Popup) PopupPopupDone() {}
func (p *Popup) SurfaceConfigure(serial uint32) {}
func (p *Popup) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {}
func (p *Popup) HandleCallbackDone(ev wl.CallbackDoneEvent) {}
func (p *Popup) HandleBufferRelease(ev wl.BufferReleaseEvent) {}

func (input *Input) PointerEnter(wlPointer *wl.Pointer, serial uint32, surface *wl.Surface, sx float32, sy float32) {}
func (input *Input) PointerLeave(wlPointer *wl.Pointer, serial uint32, wlSurface *wl.Surface) {}
func (input *Input) PointerMotion(wlPointer *wl.Pointer, time uint32, surfaceX float32, surfaceY float32) {}
func (input *Input) PointerButton(wlPointer *wl.Pointer, serial uint32, time uint32, button uint32, stateW uint32) {}
func (input *Input) PointerFrame(wlPointer *wl.Pointer) {}
func (input *Input) PointerAxis(wlPointer *wl.Pointer, time uint32, axis uint32, value float32) {}
func (input *Input) PointerAxisSource(wlPointer *wl.Pointer, axisSource uint32) {}
func (input *Input) PointerAxisStop(wlPointer *wl.Pointer, time uint32, axis uint32) {}
func (input *Input) PointerAxisDiscrete(wlPointer *wl.Pointer, axis uint32, discrete int32) {}
func (input *Input) CallbackDone(wlCallback *wl.Callback, callbackData uint32) {}
func (input *Input) HandleCallbackDone(ev wl.CallbackDoneEvent) {}
func (input *Input) HandleKeyboardKey(e wl.KeyboardKeyEvent) {}
func (input *Input) HandleKeyboardKeymap(e wl.KeyboardKeymapEvent) {}
func (input *Input) HandleKeyboardLeave(e wl.KeyboardLeaveEvent) {}
func (input *Input) HandleKeyboardModifiers(e wl.KeyboardModifiersEvent) {}
func (input *Input) HandleKeyboardRepeatInfo(e wl.KeyboardRepeatInfoEvent) {}
func (input *Input) HandleKeyboardEnter(e wl.KeyboardEnterEvent) {}
func (input *Input) HandleTouchCancel(e wl.TouchCancelEvent) {}
func (input *Input) HandleTouchDown(e wl.TouchDownEvent) {}
func (input *Input) HandleTouchFrame(e wl.TouchFrameEvent) {}
func (input *Input) HandleTouchMotion(e wl.TouchMotionEvent) {}
func (input *Input) HandleTouchOrientation(e wl.TouchOrientationEvent) {}
func (input *Input) HandleTouchShape(e wl.TouchShapeEvent) {}
func (input *Input) HandleTouchUp(e wl.TouchUpEvent) {}
func (input *Input) SeatName(wlSeat *wl.Seat, name string) {}
func (input *Input) SeatCapabilities(seat *wl.Seat, caps uint32) {}
func (input *Input) HandleSeatName(ev wl.SeatNameEvent) {}
func (input *Input) HandleSeatCapabilities(ev wl.SeatCapabilitiesEvent) {}
func (input *Input) HandlePointerEnter(ev wl.PointerEnterEvent) {}
func (input *Input) HandlePointerLeave(ev wl.PointerLeaveEvent) {}
func (input *Input) HandlePointerMotion(ev wl.PointerMotionEvent) {}
func (input *Input) HandlePointerButton(ev wl.PointerButtonEvent) {}
func (input *Input) HandlePointerFrame(ev wl.PointerFrameEvent) {}
func (input *Input) HandlePointerAxis(ev wl.PointerAxisEvent) {}
func (input *Input) HandlePointerAxisSource(ev wl.PointerAxisSourceEvent) {}
func (input *Input) HandlePointerAxisStop(ev wl.PointerAxisStopEvent) {}
func (input *Input) HandlePointerAxisDiscrete(ev wl.PointerAxisDiscreteEvent) {}
func (input *Input) HandlePointerAxisValue120(ev wl.PointerAxisValue120Event) {}
func (input *Input) PointerAxisValue120(wlPointer *wl.Pointer, axis uint32, value120 int32) {}
func (input *Input) HandleDataDeviceEnter(ev wl.DataDeviceEnterEvent) {}
func (input *Input) HandleDataDeviceLeave(ev wl.DataDeviceLeaveEvent) {}
func (input *Input) HandleDataDeviceMotion(ev wl.DataDeviceMotionEvent) {}
func (input *Input) HandleDataDeviceDrop(ev wl.DataDeviceDropEvent) {}
func (input *Input) HandleDataDeviceSelection(ev wl.DataDeviceSelectionEvent) {}
func (input *Input) HandleDataDeviceDataOffer(ev wl.DataDeviceDataOfferEvent) {}
func (input *Input) Destroy() {}
