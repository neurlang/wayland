//go:build js

package window

import (
	"io"
	"syscall/js"
	"time"
	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/external/swizzle"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
	zxdg "github.com/neurlang/wayland/xdg"
	"github.com/neurlang/wayland/xkbcommon"
)

var (
	document = js.Global().Get("document")
	body     js.Value
	canvas   js.Value
	ctx      js.Value
	surface  *Widget
	windows  []*Window
	renderFn js.Func
)

type Display struct {
}

type Window struct {
	Display  *Display
	width    int32
	height   int32
	title    string
	handler  KeyboardHandler
	widgets  []*Widget
}

type Widget struct {
	window     *Window
	userdata   interface{}
	allocation Rectangle
	data       []byte
	width      int32
	height     int32
	stride     int32
	cleared    bool
}

type Input struct {
	lastKey string
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

type SeatHandler interface {
	Capabilities(i *Input, seat *wl.Seat, caps uint32)
	Name(i *Input, seat *wl.Seat, name string)
}

type GlobalHandler interface {
	HandleGlobal(d *Display, id uint32, iface string, version uint32, data interface{})
}

type FullscreenHandler interface {
	Fullscreen(*Window, WidgetHandler)
}

type CloseHandler interface {
	Close()
}

type ResizeHandler interface {
	MinimumSize() (int32, int32)
}

type DataHandler func(*Window, *Input, float32, float32, []string, *Window, WidgetHandler)

type Popuper interface {
	Render(cairo.Surface, uint32)
	Done()
	Configure() *Widget
}

type DataSource struct {
	CopyBuffer string
}

type Popup struct {
	Popup   *zxdg.Popup
	Display *Display
}

func (p *Popup) SetPopupHandler(ph Popuper) {}
func (p *Popup) BufferRelease(_ *wl.Buffer) {}
func (p *Popup) PopupGetSurface() cairo.Surface { return nil }
func (p *Popup) Destroy() {}

func DisplayCreate(args []string) (*Display, error) {
	initCanvas()
	return &Display{}, nil
}

func (d *Display) Destroy() {
}

func (d *Display) Exit() {
	js.Global().Get("window").Call("close")
}

func (d *Display) SetSeatHandler(_ SeatHandler) {
}

func (d *Display) CreateDataSource() (*DataSource, error) {
	return &DataSource{}, nil
}

func (d *Display) GetSerial() uint32 {
	return 0
}

func (d *Display) HandleRegistryGlobal(_ wl.RegistryGlobalEvent) {
}

func (d *Display) HandleRegistryGlobalRemove(_ wl.RegistryGlobalRemoveEvent) {
}

func (d *Display) HandleShmFormat(_ wl.ShmFormatEvent) {
}

func (d *Display) HandleWmBasePing(_ zxdg.WmBasePingEvent) {
}

func (d *Display) RegistryGlobal(_ *wl.Registry, _ uint32, _ string, _ uint32) {
}

func (d *Display) RegistryGlobalRemove(_ *wl.Registry, _ uint32) {
}

func (d *Display) SetGlobalHandler(_ GlobalHandler) {
}

func (d *Display) SetUserData(_ interface{}) {
}

func (d *Display) ShellPing(_ *zxdg.WmBase, _ uint32) {
}

func (d *Display) ShmFormat(_ *wl.Shm, _ uint32) {
}

func SurfaceEnter(_ *wl.Surface, _ *wl.Output) {
}

func SurfaceLeave(_ *wl.Surface, _ *wl.Output) {
}

// DisplayRun blocks with sleep for cooperative threading
func DisplayRun(d *Display) {
	startRenderLoop()
	for {
		time.Sleep(time.Hour)
	}
}

type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

const (
	CursorBottomLeft = 0
	CursorBottomRight = 1
	CursorBottom = 2
	CursorDragging = 3
	CursorLeftPtr = 4
	CursorLeft = 5
	CursorRight = 6
	CursorTopLeft = 7
	CursorTopRight = 8
	CursorTop = 9
	CursorIbeam = 10
	CursorHand1 = 11
	CursorWatch = 12
	CursorDndMove = 13
	CursorDndCopy = 14
	CursorDndForbidden = 15
	CursorBlank = 16
)

const BufferTypeShm = 1

func Create(d *Display) *Window {
	w := &Window{Display: d, width: 200, height: 200}
	w.widgets = make([]*Widget, 0)
	windows = append(windows, w)
	return w
}

func CreateUndecorated(d *Display) *Window {
	return Create(d)
}

func (w *Window) SetTitle(title string) {
	document.Set("title", title)
}

func (w *Window) SetBufferType(t int32) {
}

func (w *Window) SetKeyboardHandler(h KeyboardHandler) {
	w.handler = h
}

func (w *Window) SetFullscreenHandler(_ FullscreenHandler) {
}

func (w *Window) SetDecorationTheme(_ Theme) {
}

func (w *Window) SetFullscreen(_ bool) error {
	return nil
}

func (w *Window) SetMinimized() error {
	return nil
}

func (w *Window) ToggleMaximized() error {
	return nil
}

func (w *Window) SeMaximized(maximized bool) error {
	return nil
}

func (w *Window) SetMaximized(maximized bool) error {
	return nil
}

func (w *Window) UninhibitRedraw() {
}

func (w *Window) InhibitRedraw() {
}

func (w *Window) ScheduleResize(_ int32, _ int32) {
}

func (w *Window) AddPopupWidget(_ *Popup, _ WidgetHandler) *Widget {
	return nil
}


func (w *Window) CreatePopup(_ *wl.Seat, _, _, _, _, _ uint32) *Popup {
	return nil
}

func (w *Window) AddWidget(wh WidgetHandler) *Widget {
	widget := &Widget{window: w, userdata: wh}
	w.widgets = append(w.widgets, widget)
	return widget
}

func (w *Window) Destroy() {
	for i, win := range windows {
		if win == w {
			windows = append(windows[:i], windows[i+1:]...)
			break
		}
	}
}

func initCanvas() {
	console := js.Global().Get("console")
	
	if canvas.Truthy() {
		return
	}
	
	body = document.Get("body")
	
	existingCanvas := document.Call("getElementById", "canvas")
	if existingCanvas.Truthy() {
		canvas = existingCanvas
		
		ctx = canvas.Call("getContext", "2d")
		
		width := int(canvas.Get("width").Int())
		height := int(canvas.Get("height").Int())
		stride := width * 4
		surfaceData := make([]byte, stride*height)
		surface = &Widget{
			data:   surfaceData,
			width:  int32(width),
			height: int32(height),
			stride: int32(stride),
		}
		
		console.Call("log", "=== Using existing canvas ===")
	} else {
		width := 640
		height := 640
		
		canvas = document.Call("createElement", "canvas")
		
		canvas.Call("setAttribute", "width", js.ValueOf(width))
		canvas.Call("setAttribute", "height", js.ValueOf(height))
		
		body.Call("appendChild", canvas)
		
		ctx = canvas.Call("getContext", "2d")
		
		stride := width * 4
		surfaceData := make([]byte, stride*height)
		surface = &Widget{
			data:   surfaceData,
			width:  int32(width),
			height: int32(height),
			stride: int32(stride),
		}
		
		console.Call("log", "=== Created new canvas ===")
	}
	
	setupInputHandlers()
}

func setupInputHandlers() {
	if !canvas.Truthy() {
		return
	}
	
	canvas.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		jsButton := event.Get("button")
		var waylandButton uint32
		switch jsButton.Int() {
		case 0:
			waylandButton = 272
		case 1:
			waylandButton = 274
		case 2:
			waylandButton = 273
		}
		
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Button(*Widget, *Input, uint32, uint32, wl.PointerButtonState, WidgetHandler) }); ok {
				handler.Button(windows[0].widgets[0], &Input{}, uint32(time.Now().UnixNano()/1000000), waylandButton, wl.PointerButtonStatePressed, windows[0].widgets[0].userdata.(WidgetHandler))
			}
		}
		return nil
	}))
	
	canvas.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		jsButton := event.Get("button")
		var waylandButton uint32
		switch jsButton.Int() {
		case 0:
			waylandButton = 272
		case 1:
			waylandButton = 274
		case 2:
			waylandButton = 273
		}
		
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Button(*Widget, *Input, uint32, uint32, wl.PointerButtonState, WidgetHandler) }); ok {
				handler.Button(windows[0].widgets[0], &Input{}, uint32(time.Now().UnixNano()/1000000), waylandButton, wl.PointerButtonStateReleased, windows[0].widgets[0].userdata.(WidgetHandler))
			}
		}
		return nil
	}))
	
	canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		xVal := event.Get("offsetX")
		yVal := event.Get("offsetY")
		
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			cursor := -1
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Motion(*Widget, *Input, uint32, float32, float32) int }); ok {
				cursor = handler.Motion(windows[0].widgets[0], &Input{}, 0, float32(xVal.Float()), float32(yVal.Float()))
			}
			
			if cursor == CursorHand1 {
				style := canvas.Get("style")
				style.Set("cursor", "pointer")
			} else {
				style := canvas.Get("style")
				style.Set("cursor", "default")
			}
		}
		return nil
	}))
	
	js.Global().Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCodeVal := event.Get("keyCode")
		keyStr := event.Get("key").String()

		notUnicode := uint32(keyCodeVal.Int())
		switch keyStr {
		case "Backspace":
			notUnicode = xkbcommon.KeyBackspace
		case "Delete":
			notUnicode = xkbcommon.KeyDelete
		case "Enter":
			notUnicode = xkbcommon.KeyReturn
		case "ArrowLeft":
			notUnicode = xkbcommon.KeyLeft
		case "ArrowRight":
			notUnicode = xkbcommon.KeyRight
		case "ArrowUp":
			notUnicode = xkbcommon.KeyUp
		case "ArrowDown":
			notUnicode = xkbcommon.KeyDown
		case "Home":
			notUnicode = xkbcommon.KeyHome
		case "End":
			notUnicode = xkbcommon.KeyEnd
		default:
			runes := []rune(keyStr)
			if len(runes) > 0 {
				notUnicode = uint32(runes[0])
			}
		}

		if len(windows) > 0 && windows[0].handler != nil {
			input := &Input{lastKey: keyStr}
			windows[0].handler.Key(windows[0], input, 0, uint32(keyCodeVal.Int()), notUnicode, wl.KeyboardKeyStatePressed, windows[0].widgets[0].userdata.(WidgetHandler))
		}
		return nil
	}))

	js.Global().Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCodeVal := event.Get("keyCode")
		keyStr := event.Get("key").String()

		notUnicode := uint32(keyCodeVal.Int())
		switch keyStr {
		case "Backspace":
			notUnicode = xkbcommon.KeyBackspace
		case "Delete":
			notUnicode = xkbcommon.KeyDelete
		case "Enter":
			notUnicode = xkbcommon.KeyReturn
		case "ArrowLeft":
			notUnicode = xkbcommon.KeyLeft
		case "ArrowRight":
			notUnicode = xkbcommon.KeyRight
		case "ArrowUp":
			notUnicode = xkbcommon.KeyUp
		case "ArrowDown":
			notUnicode = xkbcommon.KeyDown
		case "Home":
			notUnicode = xkbcommon.KeyHome
		case "End":
			notUnicode = xkbcommon.KeyEnd
		default:
			runes := []rune(keyStr)
			if len(runes) > 0 {
				notUnicode = uint32(runes[0])
			}
		}

		if len(windows) > 0 && windows[0].handler != nil {
			input := &Input{lastKey: keyStr}
			windows[0].handler.Key(windows[0], input, 0, uint32(keyCodeVal.Int()), notUnicode, wl.KeyboardKeyStateReleased, windows[0].widgets[0].userdata.(WidgetHandler))
		}
		return nil
	}))
}

func (w *Window) WindowGetSurface() cairo.Surface {
	if surface == nil {
		initCanvas()
	}
	return surface
}

func (w *Widget) Reference() cairo.Surface {
	return w
}

func (w *Widget) SetUserData(data func()) {
}

func (w *Widget) SetDestructor(destructor func()) {
}

func (w *Widget) ImageSurfaceGetData() []byte {
	return w.data
}

func (w *Widget) ImageSurfaceGetWidth() int {
	return int(w.width)
}

func (w *Widget) ImageSurfaceGetHeight() int {
	return int(w.height)
}

func (w *Widget) ImageSurfaceGetStride() int {
	return int(w.stride)
}

func renderToCanvas() {
	if surface == nil || len(surface.data) == 0 || !ctx.Truthy() {
		return
	}

	imgData := js.Global().Get("ImageData").New(js.ValueOf(surface.width), js.ValueOf(surface.height))
	var buf = make([]byte, len(surface.data), len(surface.data))
	copy(buf, surface.data)
	swizzle.BGRA(buf)
	js.CopyBytesToJS(imgData.Get("data"), buf)
	ctx.Call("putImageData", imgData, 0, 0)
}

func (w *Widget) SetUserDataWidgetHandler(wh WidgetHandler) {
	w.userdata = wh
}

func (w *Widget) GetAllocation() Rectangle {
	return w.allocation
}

func (w *Widget) SetAllocation(a int32, b int32, c int32, d int32) {
	w.allocation = Rectangle{a, b, c, d}
}

func (w *Widget) AddWidget(_ WidgetHandler) *Widget {
	return w
}

func (i *Input) GetModifiers() ModType {
	return 0
}

func (i *Input) GetRune(sym *uint32, key uint32) rune {
	switch *sym {
	case xkbcommon.KeyBackspace, xkbcommon.KeyDelete, xkbcommon.KeyReturn, xkbcommon.KeyKpEnter,
		xkbcommon.KeyLeft, xkbcommon.KeyRight, xkbcommon.KeyUp, xkbcommon.KeyDown,
		xkbcommon.KeyHome, xkbcommon.KeyEnd:
		return 0
	}
	return rune(*sym)
}

func (i *Input) GetUtf8() []byte {
	return []byte(i.lastKey)
}

func (i *Input) ReceiveSelectionData(_ string, _ io.WriteCloser) error {
	return nil
}

func (i *Input) DeviceSetSelection(_ *DataSource, _ uint32) {
}

func (s *DataSource) RemoveListener(_ wlclient.DataSourceListener) {
}

func (s *DataSource) Offer(_ string) {
}

func (s *DataSource) AddListener(_ wlclient.DataSourceListener) {
}

func (w *Widget) ScheduleResize(width int32, height int32) {
	w.allocation.Width = width
	w.allocation.Height = height
	if w.window != nil && w.userdata != nil {
		if handler, ok := w.userdata.(interface{ Resize(*Widget, int32, int32, int32, int32) }); ok {
			handler.Resize(w, width, height, width, height)
		}
	}
}

func (w *Widget) ScheduleRedraw() {
	renderToCanvas()
}

func (w *Widget) WidgetGetLastTime() uint32 {
	return 0
}

func (w *Widget) Destroy() {
}

type KeyboardHandler interface {
	Key(window *Window, input *Input, time uint32, key uint32, notUnicode uint32, state wl.KeyboardKeyState, data WidgetHandler)
	Focus(window *Window, device *Input)
}

// startRenderLoop sets up the rendering loop
func startRenderLoop() {
	renderFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Redraw(*Widget) }); ok {
				handler.Redraw(windows[0].widgets[0])
			}
		}
		renderToCanvas()
		js.Global().Call("requestAnimationFrame", renderFn)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFn)
}
