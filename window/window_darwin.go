// +build darwin

package window

import (
	"sync"

	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/xdg"
)

type Display struct {
	count   int
	windows []*Window
	mu      sync.RWMutex
}

type Window struct {
	darwinHandle   *darwinWindowHandle
	widgets        map[*Widget]struct{}
	input          *Input
	parent_display *Display
	inhibited      bool
	maximized      bool
	fullscreen     bool
	width          int32
	height         int32
	
	Display *Display
	Popup   *xdg.Popup
}

type Popup struct {
	Display *Display
	Popup   *xdg.Popup
}

// DisplayCreate creates a new Display instance for macOS
func DisplayCreate(args []string) (*Display, error) {
	return &Display{
		windows: make([]*Window, 0),
	}, nil
}

// Destroy cleans up the Display
func (d *Display) Destroy() {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	for _, w := range d.windows {
		if w.darwinHandle != nil {
			darwin_destroyWindow(w.darwinHandle)
		}
	}
	d.windows = nil
}

// DisplayRun starts the macOS event loop
func DisplayRun(d *Display) {
	darwin_runMainLoop()
}

// Exit stops the event loop
func (d *Display) Exit() {
	darwin_stopMainLoop()
}

// Create creates a new Window
func Create(d *Display) *Window {
	w := &Window{
		parent_display: d,
		Display:        d,
		widgets:        make(map[*Widget]struct{}),
		input:          &Input{},
		width:          800,
		height:         600,
		darwinHandle:   nil, // Don't create window yet - wait for first resize
	}
	
	d.mu.Lock()
	d.windows = append(d.windows, w)
	d.count++
	d.mu.Unlock()
	
	return w
}

// Destroy closes the window
func (w *Window) Destroy() {
	if w.darwinHandle != nil {
		darwin_destroyWindow(w.darwinHandle)
		w.darwinHandle = nil
	}
	
	w.parent_display.mu.Lock()
	w.parent_display.count--
	if w.parent_display.count == 0 {
		w.parent_display.Exit()
	}
	w.parent_display.mu.Unlock()
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) {
	if w.darwinHandle != nil {
		darwin_setTitle(w.darwinHandle, title)
	}
}

// SetFullscreen toggles fullscreen mode
func (w *Window) SetFullscreen(fullscreen bool) error {
	if w.darwinHandle != nil {
		darwin_setFullscreen(w.darwinHandle, fullscreen)
		w.fullscreen = fullscreen
	}
	return nil
}

// SeMaximized sets the maximized state
func (w *Window) SeMaximized(maximized bool) error {
	// Note: maximized functionality not yet implemented in CGO layer
	w.maximized = maximized
	return nil
}

// SetMaximized sets the maximized state (calls SeMaximized for compatibility)
func (w *Window) SetMaximized(maximized bool) error {
	return w.SeMaximized(maximized)
}

// ToggleMaximized toggles the maximized state
func (w *Window) ToggleMaximized() error {
	w.maximized = !w.maximized
	// Note: maximized functionality not yet implemented in CGO layer
	return nil
}

// SetKeyboardHandler sets the keyboard event handler
func (w *Window) SetKeyboardHandler(handler KeyboardHandler) {
	w.input.keyboardHandler = handler
}

// SetFullscreenHandler sets the fullscreen handler
func (w *Window) SetFullscreenHandler(handler interface{}) {
	// Placeholder for fullscreen handler
}

// SetBufferType is a placeholder for buffer type setting
func (w *Window) SetBufferType(bufferType int32) {
	// Not needed for macOS
}

// ScheduleResize resizes the window
func (w *Window) ScheduleResize(width int32, height int32) {
	if width < 32 {
		width = 32
	}
	if height < 32 {
		height = 32
	}
	
	w.width = width
	w.height = height
	
	// Create window on first resize call with the correct size
	if w.darwinHandle == nil {
		w.darwinHandle = darwin_createWindow(w.width, w.height, "Window", w)
		// Start display link after window is created
		darwin_startDisplayLink(w.darwinHandle)
	} else {
		// Resize existing window
		darwin_resizeWindow(w.darwinHandle, width, height)
	}
	
	// Update all widgets
	for widget := range w.widgets {
		widget.SetAllocation(0, 0, width, height)
		widget.drawnHash = 0
		widget.drawnHashes = nil
		if widget.handler != nil {
			widget.handler.Resize(widget, width, height, width, height)
		}
	}
	
	// Request redraw after resize
	darwin_requestRedraw(w.darwinHandle)
}

// AddWidget adds a widget to the window
func (w *Window) AddWidget(handler WidgetHandler) *Widget {
	widget := &Widget{
		parent_window: w,
		handler:       handler,
	}
	
	if w.widgets == nil {
		w.widgets = make(map[*Widget]struct{})
	}
	
	w.widgets[widget] = struct{}{}
	
	// Initial size
	widget.SetAllocation(0, 0, w.width, w.height)
	
	return widget
}

// WindowGetSurface returns the cairo surface for the window
func (w *Window) WindowGetSurface() cairo.Surface {
	for widget := range w.widgets {
		// Return a reference, not the widget itself
		// This allows the caller to Destroy() the surface without destroying the widget
		return widget.Reference()
	}
	return nil
}

// InhibitRedraw prevents redrawing
func (w *Window) InhibitRedraw() {
	w.inhibited = true
}

// UninhibitRedraw allows redrawing
func (w *Window) UninhibitRedraw() {
	w.inhibited = false
}

// SetMinimized minimizes the window
func (w *Window) SetMinimized() error {
	// Not implemented for macOS
	return nil
}

// ScheduleRedraw schedules a redraw for all widgets
func (w *Window) ScheduleRedraw() {
	if w.darwinHandle != nil {
		darwin_requestRedraw(w.darwinHandle)
	}
}

// Redraw performs the actual drawing
func (w *Window) Redraw() {
	if w.inhibited {
		return
	}
	
	for widget := range w.widgets {
		if !widget.destroyed {
			// Call handler's redraw method
			if widget.handler != nil {
				widget.handler.Redraw(widget)
			}
			
			// Get buffer and check if content changed
			buf, width, height, hash := widget.getBufferAndAllocAndHash()
			if hash != widget.drawnHash && width > 0 && height > 0 {
				darwin_drawBitmap(w.darwinHandle, buf, int32(width), int32(height))
				widget.setHashHashesRects(hash, nil, nil)
			}
		}
	}
}

// CreatePopup creates a popup window (not implemented for macOS)
func (w *Window) CreatePopup(seat *wl.Seat, clickSerial, width, height, x, y uint32) *Popup {
	// Popup windows not implemented for macOS
	return &Popup{
		Display: w.Display,
		Popup:   nil,
	}
}

// AddPopupWidget adds a popup widget (not implemented for macOS)
func (w *Window) AddPopupWidget(p *Popup, handler WidgetHandler) *Widget {
	// Popup widgets not implemented for macOS
	return w.AddWidget(handler)
}

// Popup methods (stubs for macOS)

func (p *Popup) SetPopupHandler(handler interface{}) {
	// Not implemented for macOS
}

func (p *Popup) BufferRelease(buffer *wl.Buffer) {
	// Not implemented for macOS
}

func (p *Popup) PopupGetSurface() cairo.Surface {
	// Not implemented for macOS
	return nil
}

func (p *Popup) Destroy() {
	// Not implemented for macOS
}

// Package-level functions for compatibility

func SurfaceEnter(wlSurface *wl.Surface, wlOutput *wl.Output) {
	// Not implemented for macOS
}

func SurfaceLeave(wlSurface *wl.Surface, wlOutput *wl.Output) {
	// Not implemented for macOS
}

// Dummy methods for compatibility

func (d *Display) CreateDataSource() (*DataSource, error) {
	return &DataSource{}, nil
}

func (d *Display) GetSerial() uint32 {
	return 0
}

func (d *Display) SetSeatHandler(_ interface{}) {}

func (d *Display) HandleRegistryGlobal(_ wl.RegistryGlobalEvent) {}

func (d *Display) HandleRegistryGlobalRemove(_ wl.RegistryGlobalRemoveEvent) {}

func (d *Display) HandleShmFormat(_ wl.ShmFormatEvent) {}

func (d *Display) HandleWmBasePing(_ xdg.WmBasePingEvent) {}

func (d *Display) RegistryGlobal(_ *wl.Registry, _ uint32, _ string, _ uint32) {}

func (d *Display) RegistryGlobalRemove(_ *wl.Registry, _ uint32) {}

type GlobalHandler interface {
	HandleGlobal(d *Display, id uint32, iface string, version uint32, data interface{})
}

func (d *Display) SetGlobalHandler(_ GlobalHandler) {}

func (d *Display) SetUserData(_ interface{}) {}

func (d *Display) ShellPing(*xdg.WmBase, uint32) {}

func (d *Display) ShmFormat(*wl.Shm, uint32) {}


