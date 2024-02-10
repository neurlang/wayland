package mustard

import (
	"github.com/neurlang/wayland/window"
	gg "github.com/danfragoso/thdwb/gg"
	"image"
)

// Define Constants
func SetGLFWHints() {
    // implementation here
}

// Implement type App
type App struct {
    display *window.Display
}

type Window struct {
	window *window.Window
	contextMenu *contextMenu
	rootFrame *Frame
	title string
	registeredTrees []*TreeWidget
	registeredButtons []*ButtonWidget
	registeredInputs  []*InputWidget

	cursorX float64
	cursorY float64
	cursor int

	pointerPositionEventListeners []func(float64, float64)
	scrollEventListeners          []func(int)
	clickEventListeners           []func(MustardKey)

	hasActiveOverlay bool

	selectedWidget Widget
	activeInput *InputWidget
}






func (window *Window) SetTitle(str string) {
    // implementation here
	window.title = str 
}

// RemoveStaticOverlay removes static overlay from the window
func (w *Window) RemoveStaticOverlay(str string) {
    // Dummy implementation
    println("Removing static overlay from the window", str)
}
// RemoveStaticOverlay removes static overlay from the window
func (w *Window) AddStaticOverlay(str string) {
    // Dummy implementation
    println("Adding static overlay from the window", str)
}


func (window *Window) RegisterButton(button *ButtonWidget, callback func()) {
	button.onClick = callback
	window.registeredButtons = append(window.registeredButtons, button)
}

func (window *Window) RegisterInput(input *InputWidget) {
	window.registeredInputs = append(window.registeredInputs, input)
}

func (window *Window) AttachPointerPositionEventListener(callback func(pointerX, pointerY float64)) {
	window.pointerPositionEventListeners = append(window.pointerPositionEventListeners, callback)
}

func (window *Window) AttachScrollEventListener(callback func(direction int)) {
	window.scrollEventListeners = append(window.scrollEventListeners, callback)
}

func (window *Window) AttachClickEventListener(callback func(MustardKey)) {
	window.clickEventListeners = append(window.clickEventListeners, callback)
}

func (window *Window) GetCursorPosition() (float64, float64) {
	return window.cursorX, window.cursorY
}

func (window *Window) RegisterTree(tree *TreeWidget) {
	window.registeredTrees = append(window.registeredTrees, tree)
}

func (window *Window) SetRootFrame(f *Frame) {
	f.window = window
	window.rootFrame = f
}
func (window *Window) Show() {
}

// SetCursor sets the cursor of the window
func (w *Window) SetCursor(cursorType cursorType) {
	w.cursor = int(cursorType)
}

// CreateNewWindow creates a new window
func CreateNewWindow(string, int, int, bool) *Window {
	w := &Window{} // Returning a dummy window
	//w.addEvents()

	return w
}


// CreateStaticOverlay creates a new static overlay
func CreateStaticOverlay(string, *gg.Context, image.Point) string {
    return ""
}
// AddOverlay adds an overlay to the window
func (w *Window) AddOverlay(content *Overlay) {

}

// RemoveOverlay removes the overlay from the window
func (w *Window) RemoveOverlay(*Overlay) {

}

func drawRootFrame(window *Window, w, h float64) {
	window.rootFrame.computedBox.SetCoords(0, 0, w, h)
}
func (window *Window) RequestReflowWith(w, h float64) {
	drawRootFrame(window, w, h)
}
func (window *Window) RequestReflow() {
	rect := window.rootFrame.widget.GetAllocation()
	drawRootFrame(window, float64(rect.Width), float64(rect.Height))
}


func (window *Window) addEvents() {


/*
	window.glw.SetSizeCallback(func(w *glfw.Window, width, height int) {
		xscale, yscale := float32(1), float32(1)
		if window.hiDPI {
			xscale, yscale = w.GetContentScale()
		}

		swidth := int(float32(width) / xscale)
		sheight := int(float32(height) / yscale)

		window.width, window.height = swidth, sheight
		window.RecreateContext()
		//window.RecreateOverlayContext()
		window.needsReflow = true
	})

	window.glw.SetCharCallback(func(w *glfw.Window, char rune) {
	})

	window.glw.SetCloseCallback(func(w *glfw.Window) {
		os.Exit(0)
	})

	window.glw.SetKeyCallback(func(w *glfw.Window, key glfw.Key, sc int, action glfw.Action, mods glfw.ModifierKey) {
		
	})

	window.glw.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Release {
			window.ProcessPointerClick(button)
		}
	})

	window.glw.SetScrollCallback(func(w *glfw.Window, x, y float64) {
		window.ProcessScroll(x, y)
	})
*/
}

