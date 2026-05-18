//go:build js

package window

import (
	"syscall/js"
	"time"
	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/wl"
)

var (
	document = js.Global().Get("document")
	body     js.Value
	canvas   js.Value
	ctx      js.Value
	surface  *Surface
	windows  []*Window
	renderFn js.Func
)

type Display struct {
}

type Window struct {
	display  *Display
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
}

type Input struct {
}

type WidgetHandler interface{}

func DisplayCreate(args []string) (*Display, error) {
	initCanvas()
	return &Display{}, nil
}

func (d *Display) Destroy() {
}

func (d *Display) Exit() {
	js.Global().Get("window").Call("close")
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
	w := &Window{display: d, width: 200, height: 200}
	w.widgets = make([]*Widget, 0)
	windows = append(windows, w)
	return w
}

func (w *Window) SetTitle(title string) {
	document.Set("title", title)
}

func (w *Window) SetBufferType(t int32) {
}

func (w *Window) SetKeyboardHandler(h KeyboardHandler) {
	w.handler = h
}

func (w *Window) AddWidget(wh interface{}) *Widget {
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
		surface = &Surface{
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
		surface = &Surface{
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
		buttonVal := event.Get("button")
		
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Button(*Widget, *Input, uint32, uint32, wl.PointerButtonState, WidgetHandler) }); ok {
				handler.Button(windows[0].widgets[0], &Input{}, 0, uint32(buttonVal.Int()), wl.PointerButtonStatePressed, windows[0].widgets[0].userdata.(WidgetHandler))
			}
		}
		return nil
	}))
	
	canvas.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		buttonVal := event.Get("button")
		
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Button(*Widget, *Input, uint32, uint32, wl.PointerButtonState, WidgetHandler) }); ok {
				handler.Button(windows[0].widgets[0], &Input{}, 0, uint32(buttonVal.Int()), wl.PointerButtonStateReleased, windows[0].widgets[0].userdata.(WidgetHandler))
			}
		}
		return nil
	}))
	
	canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		xVal := event.Get("offsetX")
		yVal := event.Get("offsetY")
		
		if len(windows) > 0 && windows[0].widgets[0].userdata != nil {
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
		
		if len(windows) > 0 && windows[0].handler != nil {
			windows[0].handler.Key(windows[0], &Input{}, 0, uint32(keyCodeVal.Int()), uint32(keyCodeVal.Int()), wl.KeyboardKeyStatePressed, windows[0].widgets[0].userdata.(WidgetHandler))
		}
		return nil
	}))
	
	js.Global().Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCodeVal := event.Get("keyCode")
		
		if len(windows) > 0 && windows[0].handler != nil {
			windows[0].handler.Key(windows[0], &Input{}, 0, uint32(keyCodeVal.Int()), uint32(keyCodeVal.Int()), wl.KeyboardKeyStateReleased, windows[0].widgets[0].userdata.(WidgetHandler))
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

type Surface struct {
	data    []byte
	width   int32
	height  int32
	stride  int32
	cleared bool
}

func (s *Surface) Reference() cairo.Surface {
	return s
}

func (s *Surface) Destroy() {
}

func (s *Surface) SetUserData(data func()) {
}

func (s *Surface) SetDestructor(destructor func()) {
}

func (s *Surface) ImageSurfaceGetData() []byte {
	return s.data
}

func (s *Surface) ImageSurfaceGetWidth() int {
	return int(s.width)
}

func (s *Surface) ImageSurfaceGetHeight() int {
	return int(s.height)
}

func (s *Surface) ImageSurfaceGetStride() int {
	return int(s.stride)
}

func renderToCanvas() {
	if surface == nil || len(surface.data) == 0 || !ctx.Truthy() {
		return
	}

	imgData := js.Global().Get("ImageData").New(js.ValueOf(surface.width), js.ValueOf(surface.height))
	js.CopyBytesToJS(imgData.Get("data"), surface.data)
	ctx.Call("putImageData", imgData, 0, 0)
}

func (w *Widget) SetUserDataWidgetHandler(wh interface{}) {
	w.userdata = wh
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
