package window

import cairo "github.com/neurlang/wayland/cairoshim"
import "github.com/neurlang/winc"
import "sync"
import "time"
import "github.com/spaolacci/murmur3"
type Widget struct {
	userdata interface{}
	//canvas                     *winc.Canvas
	buffer []byte
	swapbuffer []byte
	drawnHash  uint64
	drawnHashes map[int]uint64
	drawnRects  map[[2]int][]rect

	allocation_x, allocation_y int

	parent_window *Window
	handler       WidgetHandler

	allocation_width  int
	allocation_height int
	destroyed         bool
	scheduled         bool
	draw_mut sync.Mutex
}

const BUFFER_BYTES = 4

func (w *Widget) ImageSurfaceGetData() []byte {
	if len(w.buffer) == 0 {

		w.buffer = make([]byte, BUFFER_BYTES*w.allocation_width*w.allocation_height, BUFFER_BYTES*w.allocation_width*w.allocation_height)
	}
	return w.buffer
}

func (w *Widget) ImageSurfaceGetWidth() int {
	return w.allocation_width
}

func (w *Widget) ImageSurfaceGetHeight() int {
	return w.allocation_height
}

func hashBuffer(buf []byte, y, end, w int) uint64 {
	hash := murmur3.Sum64WithSeed(
		buf[
			BUFFER_BYTES*y*w:
			BUFFER_BYTES*end*w], uint32(y))
	return hash
}

func (w *Widget) getBufferAndAllocAndHash() ([]byte, int, int, uint64) {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()
	w.swapbuffer = make([]byte, len(w.buffer), len(w.buffer))
	copy(w.swapbuffer, w.buffer)
	return w.swapbuffer, w.allocation_width, w.allocation_height, murmur3.Sum64WithSeed(w.swapbuffer, 0)
}

func (w *Widget) ImageSurfaceGetStride() int {
	return w.allocation_width * BUFFER_BYTES
}

func (w *Widget) Reference() cairo.Surface {
	return w
}

func (w *Widget) SetDestructor(f func()) {
}
func (w *Widget) SetUserData(f func()) {
}
func (w *Widget) SetUserDataWidgetHandler(wh WidgetHandler) {
	w.userdata = wh
}


func (w *Widget) ScheduleResize(width int32, height int32) {
	println("ScheduleResize", width, height)

	var bx = (w.parent_window.form.Width() - w.parent_window.form.ClientWidth())
	var by = w.parent_window.form.Height() - w.parent_window.form.ClientHeight()

	// simple impl
	w.parent_window.form.SetSize(int(width)+bx, int(height)+by)
	w.draw_mut.Lock()
	w.drawnHash = 0
	w.drawnHashes = nil
	w.drawnRects = nil
	w.draw_mut.Unlock()
	w.handler.Resize(w, int32(width), int32(height), int32(width), int32(height))
}

func (w *Widget) Destroy() {
	w.destroyed = true
}

func (w *Widget) SetAllocation(x int32, y int32, pwidth int32, pheight int32) {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()
	
	if pwidth <= 0 || pheight <= 0 {
		w.buffer = nil
		w.swapbuffer = nil
		w.drawnHash = 0
		w.drawnHashes = nil
		w.drawnRects = nil
		w.allocation_x = int(x)
		w.allocation_y = int(y)
		w.allocation_width = 0
		w.allocation_height = 0
	}
	
	w.allocation_x = int(x)
	w.allocation_y = int(y)
	w.allocation_width = int(pwidth)
	w.allocation_height = int(pheight)
	w.buffer = make([]byte, BUFFER_BYTES*w.allocation_width*w.allocation_height, BUFFER_BYTES*w.allocation_width*w.allocation_height)
	w.swapbuffer = nil
	w.drawnHash = 0
	w.drawnHashes = nil
	w.drawnRects = nil
}

func (w *Widget) setHashHashesRects(hash uint64, hashes map[int]uint64, rects map[[2]int][]rect) {
	w.draw_mut.Lock()
	defer w.draw_mut.Unlock()
	if w.swapbuffer == nil {
		return
	}
	w.drawnHash = hash
	w.drawnHashes = hashes
	w.drawnRects = rects
	w.swapbuffer = nil
}

func (w *Widget) WidgetGetLastTime() uint32 {
	return 0
}

func (w *Widget) ScheduleRedraw() {
	go func() {
	//println("ScheduleRedraw")
	w.draw_mut.Lock()
	is_sch := w.scheduled
	w.draw_mut.Unlock()
	if !is_sch {
		w.draw_mut.Lock()
		w.scheduled = true
		
		w.handler.Redraw(w)
		w.draw_mut.Unlock()
		//
		w.parent_window.form.Invalidate(false)
		redrawer(w, winc.NewCanvasFromHwnd(w.parent_window.form.Handle()))
		w.parent_window.form.Invalidate(false)
		w.draw_mut.Lock()
		w.scheduled = false
		w.draw_mut.Unlock()
		time.Sleep(8*time.Millisecond)
		w.ScheduleRedraw()
	}
	}()
}

type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}
// FIXME: unimplemented
func (w *Widget) GetAllocation() Rectangle {
	// TODO: implement this
	return Rectangle{}
}

// FIXME: unimplemented
func (w *Widget) AddWidget(_ WidgetHandler) *Widget {
	// TODO: implement this
	return w
}
