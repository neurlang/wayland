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
	width          int32
	height         int32
	inhibited      bool
	maximized      bool
	fullscreen     bool
	decorated      bool  // Whether window has decorations (title bar, etc.)
	
	popupList      [][5]uintptr

	Display *Display
}

type Popup struct {
	Popup   *xdg.Popup
	popuper      Popuper
	Display      *Display
	popupWindow  *Window
	parentWindow *Window
	x            int32
	y            int32
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

// Create creates a new Window with decorations
func Create(d *Display) *Window {
	w := &Window{
		parent_display: d,
		Display:        d,
		widgets:        make(map[*Widget]struct{}),
		input:          &Input{},
		width:          800,
		height:         600,
		decorated:      true,  // Window has decorations
		darwinHandle:   nil, // Don't create window yet - wait for first resize
	}
	
	d.mu.Lock()
	d.windows = append(d.windows, w)
	d.count++
	d.mu.Unlock()
	
	return w
}

// CreateUndecorated creates a new Window without decorations (borderless)
func CreateUndecorated(d *Display) *Window {
	w := &Window{
		parent_display: d,
		Display:        d,
		widgets:        make(map[*Widget]struct{}),
		input:          &Input{},
		width:          800,
		height:         600,
		decorated:      false,  // Window has no decorations
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
	
	// Only decrement count and potentially exit if this is NOT a popup window
	// Popup windows shouldn't affect the main window count
	if w.parent_display != nil {
		w.parent_display.mu.Lock()
		
		// Check if this window is in the windows list (main windows only)
		isMainWindow := false
		for i, win := range w.parent_display.windows {
			if win == w {
				// Remove from windows list
				w.parent_display.windows = append(w.parent_display.windows[:i], w.parent_display.windows[i+1:]...)
				isMainWindow = true
				break
			}
		}
		
		// Only decrement count and exit for main windows
		if isMainWindow {
			w.parent_display.count--
			if w.parent_display.count == 0 {
				w.parent_display.Exit()
			}
		}
		
		w.parent_display.mu.Unlock()
	}
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
		w.darwinHandle = darwin_createWindow(w.width, w.height, "Window", w.decorated, w)
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

// CreatePopup creates a popup window for macOS
func (w *Window) CreatePopup(seat *wl.Seat, clickSerial, width, height, x, y uint32) *Popup {
	println("[DEBUG] CreatePopup called")
	println("[DEBUG]   width:", width, "height:", height)
	println("[DEBUG]   x:", x, "y:", y)
	println("[DEBUG]   parent window:", w)
	
	// Create a new undecorated window for the popup
	popupWindow := CreateUndecorated(w.Display)
	popupWindow.width = int32(width)
	popupWindow.height = int32(height)
	
	// Store parent window reference
	popupWindow.parent_display = w.Display
	
	println("[DEBUG] Created popup window:", popupWindow)
	
	return &Popup{
		Display:      w.Display,
		Popup:        nil,
		popupWindow:  popupWindow,
		parentWindow: w,
		x:            int32(x),
		y:            int32(y),
	}
}

// AddPopupWidget adds a popup widget to the popup window
func (w *Window) AddPopupWidget(p *Popup, handler WidgetHandler) *Widget {
	println("[DEBUG] AddPopupWidget called")
	
	println("[DEBUG] Popup window size:", p.popupWindow.width, "x", p.popupWindow.height)
	
	// IMPORTANT: Create the popup window BEFORE adding widget
	// This ensures the window has a valid darwinHandle and Cairo surface
	if p.popupWindow.darwinHandle == nil {
		println("[DEBUG] Creating popup window with ScheduleResize")
		p.popupWindow.ScheduleResize(p.popupWindow.width, p.popupWindow.height)
	}
	
	// Add widget to the popup window
	widget := p.popupWindow.AddWidget(handler)
	println("[DEBUG] Widget added to popup window")
	
	// Position the popup relative to parent window AFTER window is created
	if p.parentWindow != nil && p.parentWindow.darwinHandle != nil && p.popupWindow.darwinHandle != nil {
		if p.popupWindow == p.parentWindow {
			panic("can't be same")
		}
		println("[DEBUG] Positioning popup at offset:", p.x, p.y)
		darwin_positionPopup(p.popupWindow.darwinHandle, p.parentWindow.darwinHandle, p.x, p.y)
	} else {
		println("[DEBUG] ERROR: Cannot position popup - missing handles")
		if p.parentWindow == nil {
			println("[DEBUG]   parentWindow is nil")
		} else if p.parentWindow.darwinHandle == nil {
			println("[DEBUG]   parentWindow.darwinHandle is nil")
		}
		if p.popupWindow.darwinHandle == nil {
			println("[DEBUG]   popupWindow.darwinHandle is nil")
		}
	}
	
	println("[DEBUG] Popup widget added successfully")


	return widget
}

// Popup methods (stubs for macOS)

type Popuper interface {
	Render(cairo.Surface, uint32)
	Done()
	Configure() *Widget
}

func (p *Popup) SetPopupHandler(handler Popuper) {
	p.popuper = handler

}

func (p *Popup) BufferRelease(buffer *wl.Buffer) {
	// Not implemented for macOS
}

func (p *Popup) PopupGetSurface() (ret cairo.Surface) {
	println("[DEBUG] PopupGetSurface called")
	if p.popupWindow != nil {
		ret = p.popupWindow.WindowGetSurface()
	}
	if ret == nil {
		if p.popuper != nil {
			p.popuper.Configure()
		}
		if p.popupWindow != nil {
			ret = p.popupWindow.WindowGetSurface()
		}
	}
	return ret
}

func (p *Popup) Destroy() {
	println("[DEBUG] Popup Destroy called")
	goPopupGone(p.parentWindow, p.popupWindow)
	if p.popupWindow != nil {
		p.popupWindow.Destroy()
		p.popupWindow = nil
	}
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


