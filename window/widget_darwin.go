package window

import (
	"sync"
	"time"

	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/spaolacci/murmur3"
)

// Widget represents a drawable area within a window on macOS
type Widget struct {
	userdata interface{}
	buffer   []byte
	swapbuffer []byte
	drawnHash   uint64
	drawnHashes map[int]uint64

	allocation_x      int
	allocation_y      int
	allocation_width  int
	allocation_height int

	parent_window *Window
	handler       WidgetHandler

	destroyed bool
	scheduled bool
	draw_mut  sync.Mutex
}

const BUFFER_BYTES = 4

// ImageSurfaceGetData returns the pixel buffer for drawing
func (w *Widget) ImageSurfaceGetData() []byte {
	if len(w.buffer) == 0 {
		w.buffer = make([]byte, BUFFER_BYTES*w.allocation_width*w.allocation_height)
	}
	return w.buffer
}

// ImageSurfaceGetWidth returns the width of the surface
func (w *Widget) ImageSurfaceGetWidth() int {
	return w.allocation_width
}

// ImageSurfaceGetHeight returns the height of the surface
func (w *Widget) ImageSurfaceGetHeight() int {
	return w.allocation_height
}

// ImageSurfaceGetStride returns the stride (bytes per row)
func (w *Widget) ImageSurfaceGetStride() int {
	return w.allocation_width * BUFFER_BYTES
}

// hashBuffer computes a hash for a portion of the buffer
func hashBuffer(buf []byte, y, end, w int) uint64 {
	hash := murmur3.Sum64WithSeed(
		buf[BUFFER_BYTES*y*w:BUFFER_BYTES*end*w], uint32(y))
	return hash
}

// getBufferAndAllocAndHash returns buffer data and its hash for rendering
func (w *Widget) getBufferAndAllocAndHash() ([]byte, int, int, uint64) {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()
	w.swapbuffer = make([]byte, len(w.buffer))
	copy(w.swapbuffer, w.buffer)
	return w.swapbuffer, w.allocation_width, w.allocation_height, murmur3.Sum64WithSeed(w.swapbuffer, 0)
}

// Reference returns the widget as a cairo surface
func (w *Widget) Reference() cairo.Surface {
	return w
}

// SetDestructor sets a destructor callback (placeholder)
func (w *Widget) SetDestructor(f func()) {
	// Not needed for macOS
}

// SetUserData sets user data callback (placeholder)
func (w *Widget) SetUserData(f func()) {
	// Not needed for macOS
}

// SetUserDataWidgetHandler sets the widget handler
func (w *Widget) SetUserDataWidgetHandler(wh WidgetHandler) {
	w.userdata = wh
}

// ScheduleResize requests a window resize
func (w *Widget) ScheduleResize(width int32, height int32) {
	if w.parent_window == nil {
		return
	}

	// On macOS, when widget is resized, also resize the parent window
	w.parent_window.ScheduleResize(width, height)

	w.draw_mut.Lock()
	w.drawnHash = 0
	w.drawnHashes = nil
	w.draw_mut.Unlock()

	// No delay needed - window is created with correct size
	if w.handler != nil {
		w.handler.Resize(w, width, height, width, height)
	}
}

// Destroy marks the widget as destroyed
func (w *Widget) Destroy() {
	w.destroyed = true
}

// SetAllocation sets the widget's position and size
func (w *Widget) SetAllocation(x int32, y int32, pwidth int32, pheight int32) {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()

	if pwidth <= 0 || pheight <= 0 {
		w.buffer = nil
		w.swapbuffer = nil
		w.drawnHash = 0
		w.drawnHashes = nil
		w.allocation_x = int(x)
		w.allocation_y = int(y)
		w.allocation_width = 0
		w.allocation_height = 0
		return
	}

	w.allocation_x = int(x)
	w.allocation_y = int(y)
	w.allocation_width = int(pwidth)
	w.allocation_height = int(pheight)
	w.buffer = make([]byte, BUFFER_BYTES*w.allocation_width*w.allocation_height)
	w.swapbuffer = nil
	w.drawnHash = 0
	w.drawnHashes = nil
}

// setHashHashesRects updates the drawn hash after rendering
func (w *Widget) setHashHashesRects(hash uint64, hashes map[int]uint64, rects interface{}) {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()
	if w.swapbuffer == nil {
		return
	}
	w.drawnHash = hash
	w.drawnHashes = hashes
	w.swapbuffer = nil
}

// WidgetGetLastTime returns the last event timestamp
func (w *Widget) WidgetGetLastTime() uint32 {
	return 0
}

// ScheduleRedraw schedules a redraw of the widget
func (w *Widget) ScheduleRedraw() {
	go func() {
		w.draw_mut.Lock()
		is_sch := w.scheduled
		w.draw_mut.Unlock()
		
		if !is_sch {
			w.draw_mut.Lock()
			w.scheduled = true
			
			// Call the handler's redraw method
			if w.handler != nil {
				w.handler.Redraw(w)
			}
			w.draw_mut.Unlock()

			// Trigger window redraw
			if w.parent_window != nil {
				w.parent_window.Redraw()
			}

			w.draw_mut.Lock()
			w.scheduled = false
			w.draw_mut.Unlock()

			// Small delay before allowing next redraw
			time.Sleep(8 * time.Millisecond)
			w.ScheduleRedraw()
		}
	}()
}

// Rectangle represents a rectangular area
type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// GetAllocation returns the widget's current allocation
func (w *Widget) GetAllocation() Rectangle {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()
	return Rectangle{
		X:      int32(w.allocation_x),
		Y:      int32(w.allocation_y),
		Width:  int32(w.allocation_width),
		Height: int32(w.allocation_height),
	}
}

// AddWidget adds a child widget (not implemented for macOS)
func (w *Widget) AddWidget(_ WidgetHandler) *Widget {
	// TODO: implement child widgets if needed
	return w
}

// Input represents input state for keyboard and mouse
type Input struct {
	keyboardHandler KeyboardHandler
}

// Input methods for compatibility

func (input *Input) GetModifiers() ModType {
	// Not implemented for macOS
	return 0
}

func (input *Input) GetRune(sym *uint32, v uint32) rune {
	// Not implemented for macOS
	return 0
}

func (input *Input) GetUtf8() []byte {
	// Not implemented for macOS
	return nil
}

func (input *Input) DeviceSetSelection(ds *DataSource, num uint32) {
	// Not implemented for macOS
}

func (input *Input) ReceiveSelectionData(str string, val interface{}) error {
	// Not implemented for macOS
	return nil
}

// DataSource represents a data source for clipboard/drag-drop
type DataSource struct {
	CopyBuffer string
}

// DataSource methods for compatibility

func (ds *DataSource) AddListener(l interface{}) {
	// Not implemented for macOS
}

func (ds *DataSource) Offer(str string) {
	// Not implemented for macOS
}

func (ds *DataSource) RemoveListener(l interface{}) {
	// Not implemented for macOS
}


