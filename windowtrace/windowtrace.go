package windowtrace

import (
	"io"

	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/window"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
	zxdg "github.com/neurlang/wayland/xdg"
)

type Display window.Display
type Window struct {
	*window.Window
	Display *Display
}

type DataHandler window.DataHandler
type CloseHandler interface {
	Close()
}

const SurfaceOpaque = window.SurfaceOpaque
const SurfaceShm = window.SurfaceShm

const SurfaceHintResize = window.SurfaceHintResize
const SurfaceHintRgb565 = window.SurfaceHintRgb565

const PreferredFormatNone = window.PreferredFormatNone
const PreferredFormatRgb565 = window.PreferredFormatRgb565

const BufferTypeEglWindow = window.BufferTypeEglWindow
const BufferTypeShm = window.BufferTypeShm
const (
	CursorBottomLeft   = window.CursorBottomLeft
	CursorBottomRight  = window.CursorBottomRight
	CursorBottom       = window.CursorBottom
	CursorDragging     = window.CursorDragging
	CursorLeftPtr      = window.CursorLeftPtr
	CursorLeft         = window.CursorLeft
	CursorRight        = window.CursorRight
	CursorTopLeft      = window.CursorTopLeft
	CursorTopRight     = window.CursorTopRight
	CursorTop          = window.CursorTop
	CursorIbeam        = window.CursorIbeam
	CursorHand1        = window.CursorHand1
	CursorWatch        = window.CursorWatch
	CursorDndMove      = window.CursorDndMove
	CursorDndCopy      = window.CursorDndCopy
	CursorDndForbidden = window.CursorDndForbidden
	CursorBlank        = window.CursorBlank
)

const ZwpRelativePointerManagerV1Version = window.ZwpRelativePointerManagerV1Version
const ZwpPointerConstraintsV1Version = window.ZwpPointerConstraintsV1Version



func DisplayCreate(argv []string) (d *Display, e error) {
	println("func DisplayCreate")
	disp, err := window.DisplayCreate(argv)
	return (*Display)(disp), err
}

func DisplayRun(d *Display) {
	println("func DisplayRun")
	window.DisplayRun((*window.Display)(d))
}

func Create(d *Display) *Window {
	println("func Create")
	w := &Window{
		Window:  (window.Create((*window.Display)(d))),
		Display: d,
	}
	return w
}

func CreateUndecorated(d *Display) *Window {
	println("func CreateUndecorated")
	w := &Window{
		Window:  (window.CreateUndecorated((*window.Display)(d))),
		Display: d,
	}
	return w
}

func SurfaceEnter(wlSurface *wl.Surface, wlOutput *wl.Output) {
	println("func SurfaceEnter")
	window.SurfaceEnter(wlSurface, wlOutput)
}

func SurfaceLeave(wlSurface *wl.Surface, wlOutput *wl.Output) {
	println("func SurfaceLeave")
	window.SurfaceLeave(wlSurface, wlOutput)
}

type DataSource window.DataSource

func (d *Display) SetSeatHandler(h SeatHandler) {
	println("func SetSeatHandler")
	((*window.Display)(d)).SetSeatHandler(seatHandler{h})
}

type GlobalHandler interface {
	HandleGlobal(d *Display, id uint32, iface string, version uint32, data interface{})
}

type globalHandler struct {
	GlobalHandler
}

func (g globalHandler) HandleGlobal(d *window.Display, id uint32, iface string, version uint32, data interface{}) {
	println("func HandleGlobal")
	g.GlobalHandler.HandleGlobal((*Display)(d), id, iface, version, data)
}

type Input window.Input

type keyboardHandler struct {
	KeyboardHandler
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
	Focus(window *Window, input *Input)
}

func (k keyboardHandler) Key(
	window *window.Window,
	input *window.Input,
	time uint32,
	key uint32,
	notUnicode uint32,
	state wl.KeyboardKeyState,
	data window.WidgetHandler,
) {
	println("func Key")
	var wh2 WidgetHandler
	wh, ok := data.(widgetHandler)
	if ok {
		wh2 = wh.WidgetHandler
	}
	k.KeyboardHandler.Key(
		&Window{Window: window, Display: (*Display)(window.Display)},
		(*Input)(input),
		time,
		key,
		notUnicode,
		state,
		wh2,
	)
}

func (k keyboardHandler) Focus(window *window.Window, input *window.Input) {
	println("func Focus")
	k.KeyboardHandler.Focus(&Window{Window: window, Display: (*Display)(window.Display)}, (*Input)(input))
}

type Popup struct {
	Popup   *zxdg.Popup
	Display *Display
	nested  *window.Popup
}

func (p *Popup) SetPopupHandler(ph Popuper) {
	println("func SetPopupHandler")
	p.nested.SetPopupHandler(popuper{ph})
}

func (parent *Widget) AddWidget(data WidgetHandler) *Widget {
	println("func AddWidget")
	return (*Widget)((*window.Widget)(parent).AddWidget(widgetHandler{data}))
}

type Popuper interface {
	Render(cairo.Surface, uint32)
	Done()
	Configure() *Widget
}

type popuper struct {
	Popuper
}

func (p popuper) Render(s cairo.Surface, n uint32) {
	println("func Render")
	p.Popuper.Render(s, n)
}
func (p popuper) Done() {
	println("func Done")
	p.Popuper.Done()
}
func (p popuper) Configure() *window.Widget {
	println("func Configure")
	return (*window.Widget)(p.Popuper.Configure())
}

type Rectangle = window.Rectangle
type ResizeHandler interface {
	MinimumSize() (int32, int32)
}
type seatHandler struct {
	SeatHandler
}

type SeatHandler interface {
	Capabilities(i *Input, seat *wl.Seat, caps uint32)
	Name(i *Input, seat *wl.Seat, name string)
}

func (s seatHandler) Capabilities(i *window.Input, seat *wl.Seat, caps uint32) {
	println("func Capabilities")
	s.SeatHandler.Capabilities((*Input)(i), seat, caps)
}

func (s seatHandler) Name(i *window.Input, seat *wl.Seat, name string) {
	println("func Name")
	s.SeatHandler.Name((*Input)(i), seat, name)
}

type Widget window.Widget

type widgetHandler struct {
	WidgetHandler
}

type WidgetHandler interface {
	Resize(widget *Widget, width int32, height int32, pwidth int32, pheight int32)
	Redraw(widget *Widget)
	Enter(widget *Widget, input *Input, x float32, y float32)
	Leave(widget *Widget, input *Input)
	Motion(widget *Widget, input *Input, time uint32, x float32, y float32) int
	Button(
		widget *Widget,
		input *Input,
		time uint32,
		button uint32,
		state wl.PointerButtonState,
		data WidgetHandler,
	)
	TouchUp(widget *Widget, input *Input, serial uint32, time uint32, id int32)
	TouchDown(widget *Widget, input *Input, serial uint32, time uint32, id int32, x float32, y float32)
	TouchMotion(widget *Widget, input *Input, time uint32, id int32, x float32, y float32)
	TouchFrame(widget *Widget, input *Input)
	TouchCancel(widget *Widget, width int32, height int32)
	Axis(widget *Widget, input *Input, time uint32, axis uint32, value float32)
	AxisSource(widget *Widget, input *Input, source uint32)
	AxisStop(widget *Widget, input *Input, time uint32, axis uint32)
	AxisDiscrete(widget *Widget, input *Input, axis uint32, discrete int32)
	PointerFrame(widget *Widget, input *Input)
}

type AxisValue120Handler interface {
	AxisValue120(widget *Widget, input *Input, axis uint32, value120 int32)
}

func (w widgetHandler) Resize(widget *window.Widget, width int32, height int32, pwidth int32, pheight int32) {
	println("func Resize")
	w.WidgetHandler.Resize((*Widget)(widget), width, height, pwidth, pheight)
}

func (w widgetHandler) Redraw(widget *window.Widget) {
	println("func Redraw")
	w.WidgetHandler.Redraw((*Widget)(widget))
}

func (w widgetHandler) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
	println("func Enter")
	w.WidgetHandler.Enter((*Widget)(widget), (*Input)(input), x, y)
}

func (w widgetHandler) Leave(widget *window.Widget, input *window.Input) {
	println("func Leave")
	w.WidgetHandler.Leave((*Widget)(widget), (*Input)(input))
}

func (w widgetHandler) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {
	println("func Motion")
	return w.WidgetHandler.Motion((*Widget)(widget), (*Input)(input), time, x, y)
}

func (w widgetHandler) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
	println("func Button")
	w.WidgetHandler.Button((*Widget)(widget), (*Input)(input), time, button, state, w.WidgetHandler)
}

func (w widgetHandler) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
	println("func TouchUp")
	w.WidgetHandler.TouchUp((*Widget)(widget), (*Input)(input), serial, time, id)
}

func (w widgetHandler) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
	println("func TouchDown")
	w.WidgetHandler.TouchDown((*Widget)(widget), (*Input)(input), serial, time, id, x, y)
}

func (w widgetHandler) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {
	println("func TouchMotion")
	w.WidgetHandler.TouchMotion((*Widget)(widget), (*Input)(input), time, id, x, y)
}

func (w widgetHandler) TouchFrame(widget *window.Widget, input *window.Input) {
	println("func TouchFrame")
	w.WidgetHandler.TouchFrame((*Widget)(widget), (*Input)(input))
}

func (w widgetHandler) TouchCancel(widget *window.Widget, width int32, height int32) {
	println("func TouchCancel")
	w.WidgetHandler.TouchCancel((*Widget)(widget), width, height)
}

func (w widgetHandler) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value float32) {
	println("func Axis")
	w.WidgetHandler.Axis((*Widget)(widget), (*Input)(input), time, axis, value)
}

func (w widgetHandler) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
	println("func AxisSource")
	w.WidgetHandler.AxisSource((*Widget)(widget), (*Input)(input), source)
}

func (w widgetHandler) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
	println("func AxisStop")
	w.WidgetHandler.AxisStop((*Widget)(widget), (*Input)(input), time, axis)
}

func (w widgetHandler) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
	println("func AxisDiscrete")
	w.WidgetHandler.AxisDiscrete((*Widget)(widget), (*Input)(input), axis, discrete)
}

func (w widgetHandler) PointerFrame(widget *window.Widget, input *window.Input) {
	println("func PointerFrame")
	w.WidgetHandler.PointerFrame((*Widget)(widget), (*Input)(input))
}

func (w *Window) SetTitle(title string) {
	println("func SetTitle")
	(w.Window).SetTitle(title)
}

func (w *Window) SetBufferType(t int32) {
	println("func SetBufferType")
	(w.Window).SetBufferType(t)
}

func (w *Window) AddWidget(data WidgetHandler) *Widget {
	println("func AddWidget")
	return (*Widget)((w.Window).AddWidget(widgetHandler{data}))
}

func (w *Window) SetKeyboardHandler(handler KeyboardHandler) {
	println("func SetKeyboardHandler")
	(w.Window).SetKeyboardHandler(keyboardHandler{handler})
}

func (w *Window) ScheduleResize(width int32, height int32) {
	println("func ScheduleResize")
	(w.Window).ScheduleResize(width, height)
}

func (w *Window) Destroy() {
	println("func Destroy")
	(w.Window).Destroy()
}

func (w *Window) AddPopupWidget(p *Popup, data WidgetHandler) *Widget {
	println("func AddPopupWidget")
	return (*Widget)((w.Window).AddPopupWidget((*window.Popup)(p.nested), widgetHandler{data}))
}

func (parent *Widget) ScheduleResize(width int32, height int32) {
	println("func ScheduleResize")
	((*window.Widget)(parent)).ScheduleResize(width, height)
}

func (w *Window) SetFullscreen(fullscreen bool) error {
	println("func SetFullscreen")
	return (w.Window).SetFullscreen(fullscreen)
}

func (w *Window) SetMinimized() error {
	println("func SetMinimized")
	return (w.Window).SetMinimized()
}

func (w *Window) SeMaximized(maximized bool) error {
	println("func SeMaximized")
	return (w.Window).SeMaximized(maximized)
}

func (w *Window) ToggleMaximized() error {
	println("func ToggleMaximized")
	return (w.Window).ToggleMaximized()
}

func (w *Window) UninhibitRedraw() {
	println("func UninhibitRedraw")
	(w.Window).UninhibitRedraw()
}

func (w *Window) InhibitRedraw() {
	println("func InhibitRedraw")
	(w.Window).InhibitRedraw()
}

func (w *Window) ScheduleRedraw() {
	println("func ScheduleRedraw")
	(w.Window).ScheduleRedraw()
}

func (w *Window) ToplevelClose(zxdgToplevelV6 *zxdg.Toplevel) {
	println("func ToplevelClose")
	(w.Window).ToplevelClose(zxdgToplevelV6)
}

func (w *Window) ToplevelConfigure(zxdgToplevelV6 *zxdg.Toplevel, width int32, height int32, states []int32) {
	println("func ToplevelConfigure")
	(w.Window).ToplevelConfigure(zxdgToplevelV6, width, height, states)
}

func (w *Window) SurfaceConfigure(zxdgSurfaceV6 *zxdg.Surface, serial uint32) {
	println("func SurfaceConfigure")
	(w.Window).SurfaceConfigure(zxdgSurfaceV6, serial)
}

func (w *Window) Run(events uint32) {
	println("func Run")
	(w.Window).Run(events)
}

func (w *Window) SetCloseHandler(handler CloseHandler) {
	println("func SetCloseHandler")
	(w.Window).SetCloseHandler(handler)
}

func (w *Window) SetDataHandler(win *Window, handler DataHandler) {
	println("func SetDataHandler")
	(w.Window).SetDataHandler(win.Window, window.DataHandler(handler))
}

func (w *Window) SetMaximized(maximized bool) error {
	println("func SetMaximized")
	return (w.Window).SetMaximized(maximized)
}

func (w *Window) SetDecorationTheme(theme Theme) {
	println("func SetDecorationTheme")
	(w.Window).SetDecorationTheme(window.Theme(theme))
}

func (w *Window) WindowGetSurface() cairo.Surface {
	println("func WindowGetSurface")
	return ((*window.Window)(w.Window)).WindowGetSurface()
}

func (w *Window) CreatePopup(seat *wl.Seat, clickSerial, width, height, x, y uint32) *Popup {
	println("func CreatePopup")
	p := ((*window.Window)(w.Window)).CreatePopup(seat, clickSerial, width, height, x, y)
	return &Popup{
		nested:  p,
		Popup:   p.Popup,
		Display: (*Display)(p.Display),
	}
}

func (ds *DataSource) AddListener(l wlclient.DataSourceListener) {
	println("func AddListener")
	((*window.DataSource)(ds)).AddListener(l)
}

func (ds *DataSource) Offer(str string) {
	println("func Offer")
	((*window.DataSource)(ds)).Offer(str)
}

func (ds *DataSource) RemoveListener(l wlclient.DataSourceListener) {
	println("func RemoveListener")
	((*window.DataSource)(ds)).RemoveListener(l)
}

func (d *Display) CreateDataSource() (*DataSource, error) {
	println("func CreateDataSource")
	dataSource, err := ((*window.Display)(d)).CreateDataSource()
	return (*DataSource)(dataSource), err
}

func (d *Display) Destroy() {
	println("func Destroy")
	((*window.Display)(d)).Destroy()
}

func (d *Display) Exit() {
	println("func Exit")
	((*window.Display)(d)).Exit()
}

func (d *Display) GetSerial() uint32 {
	println("func GetSerial")
	return ((*window.Display)(d)).GetSerial()
}

func (d *Display) HandleRegistryGlobal(e wl.RegistryGlobalEvent) {
	println("func HandleRegistryGlobal")
	((*window.Display)(d)).HandleRegistryGlobal(e)
}

func (d *Display) HandleRegistryGlobalRemove(e wl.RegistryGlobalRemoveEvent) {
	println("func HandleRegistryGlobalRemove")
	((*window.Display)(d)).HandleRegistryGlobalRemove(e)
}

func (d *Display) HandleShmFormat(e wl.ShmFormatEvent) {
	println("func HandleShmFormat")
	((*window.Display)(d)).HandleShmFormat(e)
}

func (d *Display) HandleWmBasePing(ev zxdg.WmBasePingEvent) {
	println("func HandleWmBasePing")
	((*window.Display)(d)).HandleWmBasePing(ev)
}

func (d *Display) RegistryGlobal(registry *wl.Registry, id uint32, iface string, version uint32) {
	println("func RegistryGlobal")
	((*window.Display)(d)).RegistryGlobal(registry, id, iface, version)
}

func (d *Display) RegistryGlobalRemove(wlRegistry *wl.Registry, name uint32) {
	println("func RegistryGlobalRemove")
	((*window.Display)(d)).RegistryGlobalRemove(wlRegistry, name)
}

func (d *Display) SetGlobalHandler(gh GlobalHandler) {
	println("func SetGlobalHandler")
	((*window.Display)(d)).SetGlobalHandler(globalHandler{gh})
}

func (d *Display) SetUserData(data interface{}) {
	println("func SetUserData")
	((*window.Display)(d)).SetUserData(data)
}

func (d *Display) ShellPing(shell *zxdg.WmBase, serial uint32) {
	println("func ShellPing")
	((*window.Display)(d)).ShellPing(shell, serial)
}

func (d *Display) ShmFormat(wlShm *wl.Shm, format uint32) {
	println("func ShmFormat")
	((*window.Display)(d)).ShmFormat(wlShm, format)
}

func (w *Popup) BufferRelease(buffer *wl.Buffer) {
	println("func BufferRelease")
	((*window.Popup)(w.nested)).BufferRelease(buffer)
}

func (p *Popup) PopupGetSurface() cairo.Surface {
	println("func PopupGetSurface")
	return ((*window.Popup)(p.nested)).PopupGetSurface()
}

func (p *Popup) Destroy() {
	println("func Destroy")
	((*window.Popup)(p.nested)).Destroy()
}

func (p *Popup) HandlePopupConfigure(ev zxdg.PopupConfigureEvent) {
	println("func HandlePopupConfigure")
	(p.nested).HandlePopupConfigure(ev)
}

func (p *Popup) HandlePopupPopupDone(ev zxdg.PopupPopupDoneEvent) {
	println("func HandlePopupPopupDone")
	(p.nested).HandlePopupPopupDone(ev)
}

func (p *Popup) PopupConfigure(x, y, width, height int32) {
	println("func PopupConfigure")
	(p.nested).PopupConfigure(x, y, width, height)
}

func (p *Popup) PopupPopupDone() {
	println("func PopupPopupDone")
	(p.nested).PopupPopupDone()
}

func (p *Popup) SurfaceConfigure(serial uint32) {
	println("func SurfaceConfigure")
	(p.nested).SurfaceConfigure(serial)
}

func (input *Input) DeviceSetSelection(ds *DataSource, num uint32) {
	println("func DeviceSetSelection")
	((*window.Input)(input)).DeviceSetSelection((*window.DataSource)(ds), num)
}

func (input *Input) ReceiveSelectionData(str string, val io.WriteCloser) error {
	println("func ReceiveSelectionData")
	return ((*window.Input)(input)).ReceiveSelectionData(str, val)
}

func (input *Input) HandleDataDeviceEnter(ev wl.DataDeviceEnterEvent) {
	println("func HandleDataDeviceEnter")
	((*window.Input)(input)).HandleDataDeviceEnter(ev)
}

func (input *Input) HandleDataDeviceLeave(ev wl.DataDeviceLeaveEvent) {
	println("func HandleDataDeviceLeave")
	((*window.Input)(input)).HandleDataDeviceLeave(ev)
}

func (input *Input) HandleDataDeviceMotion(ev wl.DataDeviceMotionEvent) {
	println("func HandleDataDeviceMotion")
	((*window.Input)(input)).HandleDataDeviceMotion(ev)
}

func (input *Input) HandleDataDeviceDrop(ev wl.DataDeviceDropEvent) {
	println("func HandleDataDeviceDrop")
	((*window.Input)(input)).HandleDataDeviceDrop(ev)
}

func (input *Input) HandleDataDeviceSelection(ev wl.DataDeviceSelectionEvent) {
	println("func HandleDataDeviceSelection")
	((*window.Input)(input)).HandleDataDeviceSelection(ev)
}

func (input *Input) HandleDataDeviceDataOffer(ev wl.DataDeviceDataOfferEvent) {
	println("func HandleDataDeviceDataOffer")
	((*window.Input)(input)).HandleDataDeviceDataOffer(ev)
}


type FullscreenHandler interface {
	Fullscreen(*Window, WidgetHandler)
}

type fullscreenHandler struct {
	FullscreenHandler
}

func (f fullscreenHandler) Fullscreen(w *window.Window, h window.WidgetHandler) {
	wh, ok := h.(widgetHandler)
	var whan WidgetHandler
	if ok {
		whan = wh.WidgetHandler
	}
	println("func Fullscreen")
	f.FullscreenHandler.Fullscreen(&Window{Window: w, Display: (*Display)(w.Display)}, whan)
}

type Theme window.Theme

const ThemeLight = Theme(window.ThemeLight)
const ThemeDark = Theme(window.ThemeDark)

type WindowDecoration window.WindowDecoration
type ComponentType = window.ComponentType
type DecorationSurface window.DecorationSurface

const (
	ComponentNone      = window.ComponentNone
	ComponentShadow    = window.ComponentShadow
	ComponentTitle     = window.ComponentTitle
	ComponentButtonMin = window.ComponentButtonMin
	ComponentButtonMax = window.ComponentButtonMax
	ComponentButtonClose = window.ComponentButtonClose
)


const (
	CursorDefault = window.CursorDefault
	CursorUnset   = window.CursorUnset
	MaxLeaves     = window.MaxLeaves
)

const (
	ShadowMargin  = window.ShadowMargin
	TitleHeight   = window.TitleHeight
	ButtonWidth   = window.ButtonWidth
	SymDim        = window.SymDim
	ShadowBlurSize = window.ShadowBlurSize
)

const (
	CursorCrosshair        = window.CursorCrosshair
	CursorClosedHand       = window.CursorClosedHand
	CursorDisappearingItem = window.CursorDisappearingItem
	CursorResizeLeftRight  = window.CursorResizeLeftRight
	CursorResizeUpDown     = window.CursorResizeUpDown
)
