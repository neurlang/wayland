package window

import (
	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/xdg"
	"github.com/neurlang/winc"
	"github.com/neurlang/winc/w32"
	"time"
	"sync"
	"syscall"
	"unsafe"
)

type Window struct {
	form    *winc.Form
	widgets map[*Widget]struct{}
	input   *Input

	parent_display *Display
	inhibited      bool
	maximized      bool
	down           [2]uint32
	up             [2]uint32

	Display *Display
}

func Create(d *Display) *Window {

	form := winc.NewForm(nil)
	form.Center()

	form.Show()


	w := &Window{form: form, parent_display: d}
	form.OnClose().Bind(func(arg *winc.Event) {
		d.count--
		for widget := range w.widgets {
			widget.destroyed = true
		}
		if d.count == 0 {
			winc.Exit()
		} else {
			form.Close()
		}
	})

	d.count++

	return w
}

func (w *Window) SetKeyboardHandler(t KeyboardHandler) {

	w.input = &Input{}

	allRedrawer := func() {
		for widget := range w.widgets {
			if !w.inhibited {
				widget.destroyed = false
				widget.ScheduleRedraw()
				//redrawer(widget, winc.NewCanvasFromHwnd(w.form.Handle()))
			}
		}
	}
	w.form.OnKeyDown().Bind(func(arg *winc.Event) {

		vKey := uint32(arg.Data.(*winc.KeyDownEventData).VKey)

		code := uint32(arg.Data.(*winc.KeyDownEventData).Code)
		if w.down != [2]uint32{vKey, code} {

			t.Key(w, w.input, uint32(time.Now().UnixNano()/1000000), vKey, code, wl.KeyboardKeyStatePressed, nil)
			allRedrawer()
		}
		w.down = [2]uint32{vKey, code}
		w.up = [2]uint32{0, 0}
	})

	w.form.OnKeyUp().Bind(func(arg *winc.Event) {

		vKey := uint32(arg.Data.(*winc.KeyUpEventData).VKey)

		code := uint32(arg.Data.(*winc.KeyUpEventData).Code)

		if w.up != [2]uint32{vKey, code} {
			t.Key(w, w.input, uint32(time.Now().UnixNano()/1000000), vKey, code, wl.KeyboardKeyStateReleased, nil)
			allRedrawer()
		}
		w.up = [2]uint32{vKey, code}
		w.down = [2]uint32{0, 0}
	})

	w.form.OnSetFocus().Bind(func(arg *winc.Event) {
		t.Focus(w, w.input)
	})

	w.form.OnKillFocus().Bind(func(arg *winc.Event) {

		t.Focus(w, nil)
	})

}

func (w Window) SetFullscreenHandler(t interface{}) {

}

func (w *Window) SetTitle(s string) {
	w.form.SetText(s)
}

func (w Window) SetBufferType(shm interface{}) {

}





var (
	gdi32                = syscall.NewLazyDLL("gdi32.dll")
	procStretchDIBits    = gdi32.NewProc("StretchDIBits")
)

const (
	BI_RGB       = 0
	DIB_RGB_COLORS = 0
	SRCCOPY      = 0x00CC0020
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors [1]uint32
}

var mut sync.RWMutex

func redrawer(widget *Widget, canvas *winc.Canvas) {
	buf, w, h, hash := widget.getBufferAndAllocAndHash()
	if hash == widget.drawnHash || w <= 0 || h <= 0 {
		return
	}

	// Get the device context handle from the window
	hwnd := widget.parent_window.form.Handle()
	hdc := w32.GetDC(hwnd)
	if hdc == 0 {
		return
	}
	defer w32.ReleaseDC(hwnd, hdc)

	// Setup bitmap info for DIB (Device Independent Bitmap)
	var bi BITMAPINFO
	bi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bi.BmiHeader))
	bi.BmiHeader.BiWidth = int32(w)
	bi.BmiHeader.BiHeight = -int32(h) // Negative height for top-down bitmap
	bi.BmiHeader.BiPlanes = 1
	bi.BmiHeader.BiBitCount = 32 // 32 bits per pixel (RGBA)
	bi.BmiHeader.BiCompression = BI_RGB
	bi.BmiHeader.BiSizeImage = 0

	// Call StretchDIBits to draw the entire bitmap in one call
	procStretchDIBits.Call(
		uintptr(hdc),
		0,                        // destination x
		0,                        // destination y
		uintptr(w),              // destination width
		uintptr(h),              // destination height
		0,                        // source x
		0,                        // source y
		uintptr(w),              // source width
		uintptr(h),              // source height
		uintptr(unsafe.Pointer(&buf[0])), // pointer to bitmap bits
		uintptr(unsafe.Pointer(&bi)),     // pointer to BITMAPINFO
		DIB_RGB_COLORS,          // color usage
		SRCCOPY,                 // raster operation
	)

	widget.setHashHashesRects(hash, nil, nil)
}

func (w *Window) AddWidget(t WidgetHandler) (widget *Widget) {
	widget = &Widget{
		parent_window: w,
		handler:       t,
	}
	if w.widgets == nil {
		w.widgets = make(map[*Widget]struct{})
	}
	w.widgets[widget] = struct{}{}

	//canvasRedrawer := func(canvas *winc.Canvas) {
		//t.Redraw(widget)
		//redrawer(widget, canvas)
		//widget.ScheduleRedraw()
	//}

	//allRedrawer := func() {
	//	if !w.inhibited {
			//t.Redraw(widget)
			//redrawer(widget, winc.NewCanvasFromHwnd(w.form.Handle()))
			//widget.ScheduleRedraw()
		//}
	//}
	w.form.OnPaint().Bind(func(arg *winc.Event) {

		if !w.inhibited && !widget.destroyed {
			//t.Redraw(widget)
			//canvasRedrawer(arg.Data.(*winc.PaintEventData).Canvas)
			//widget.ScheduleRedraw()
		}

	})

	w.form.OnSize().Bind(func(arg *winc.Event) {
		widget.destroyed = false

		xy := arg.Data.(*winc.SizeEventData)

		for widget := range w.widgets {
			widget.SetAllocation(0, 0, int32(xy.X), int32(xy.Y))
		}
		widget.drawnHash = 0
		widget.drawnHashes = nil
		t.Resize(widget, int32(xy.X), int32(xy.Y), int32(xy.X), int32(xy.Y))
		widget.ScheduleRedraw()
		//allRedrawer()
		//widget.ScheduleRedraw()
	})

	w.form.OnMouseMove().Bind(func(arg *winc.Event) {
		xy := arg.Data.(*winc.MouseEventData)
		t.Motion(widget, w.input, uint32(time.Now().UnixNano()/1000000), float32(xy.X), float32(xy.Y))
		//allRedrawer()
	})
	w.form.OnMouseHover().Bind(func(arg *winc.Event) {
		xy := arg.Data.(*winc.MouseEventData)

		t.Motion(widget, w.input, uint32(time.Now().UnixNano()/1000000), float32(xy.X), float32(xy.Y))
		//allRedrawer()
	})

	w.form.OnLBDown().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 272, wl.PointerButtonStatePressed, t)
		//allRedrawer()
	})
	w.form.OnLBUp().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 272, wl.PointerButtonStateReleased, t)
		//allRedrawer()
	})
	w.form.OnRBDown().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 273, wl.PointerButtonStatePressed, t)
		//allRedrawer()
	})
	w.form.OnRBUp().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 273, wl.PointerButtonStateReleased, t)
		//allRedrawer()
	})
	
	return
}

func (w *Window) Destroy() {
	winc.Exit()
}

func (w *Window) SetFullscreen(fullscreen bool) {
	w.form.Fullscreen()
}

func (w *Window) WindowGetSurface() cairo.Surface {
	for widget := range w.widgets {
		return widget
	}
	return nil
}

func (w *Window) UninhibitRedraw() {
	w.inhibited = false

}

func (w *Window) SetMinimized() {
	w.form.Minimise()
}

func (w *Window) ToggleMaximized() {
	if w.maximized {
		w.form.Restore()
	} else {
		w.form.Maximise()
	}
	w.maximized = !w.maximized

}

func (w *Window) InhibitRedraw() {
	w.inhibited = true
}

func (w *Window) ScheduleResize(width int32, height int32) {

}

func (w *Window) AddPopupWidget(p *Popup, handler WidgetHandler) *Widget {
	p.widget.handler = handler

	canvasRedrawer := func(canvas *winc.Canvas) {
		p.widget.ScheduleRedraw()
	}
	allRedrawer := func() {
		if !w.inhibited {
			if p.form != nil {
				if p.popuper != nil {
					p.popuper.Render(&p.widget, 0)
				}
				p.widget.ScheduleRedraw()
			}
		}
	}
	p.form.OnPaint().Bind(func(arg *winc.Event) {

		if !p.inhibited && !p.widget.destroyed {
			p.popuper.Render(&p.widget, 0)
			canvasRedrawer(arg.Data.(*winc.PaintEventData).Canvas)
		}

	})
	p.form.OnSize().Bind(func(arg *winc.Event) {
		p.widget.destroyed = false

		xy := arg.Data.(*winc.SizeEventData)

		p.widget.SetAllocation(0, 0, int32(xy.X), int32(xy.Y))

		p.popuper.Configure()
	})

	hover := func(arg *winc.Event) {
		if p.widget.handler == nil {
			return
		}
		xy := arg.Data.(*winc.MouseEventData)

		p.widget.handler.Motion(&p.widget, w.input, uint32(time.Now().UnixNano()/1000000), float32(xy.X), float32(xy.Y))
		allRedrawer()
	}

	p.form.OnMouseMove().Bind(hover)
	p.form.OnMouseHover().Bind(hover)

	p.form.OnLBDown().Bind(func(arg *winc.Event) {
		p.widget.handler.Button(&p.widget, w.input, uint32(time.Now().UnixNano()/1000000), 272, wl.PointerButtonStatePressed, p.widget.handler)
		allRedrawer()
	})
	p.form.OnLBUp().Bind(func(arg *winc.Event) {
		p.widget.handler.Button(&p.widget, w.input, uint32(time.Now().UnixNano()/1000000), 272, wl.PointerButtonStateReleased, p.widget.handler)
		allRedrawer()
	})

	p.form.OnClose().Bind(func(arg *winc.Event) {
		p.form.OnMouseMove().Bind(nil)
		p.form.OnMouseHover().Bind(nil)
		p.form.OnLBDown().Bind(nil)
		p.form.OnLBUp().Bind(nil)
		p.Destroy()
	})
	p.form.Show()

	return &p.widget
}
func (w *Window) CreatePopup(_ *wl.Seat, _, width, height, x, y uint32) (popup *Popup) {
	form := winc.NewCustomForm(w.form, 0, w32.WS_POPUP)

	//var bx = (w.form.Width() - w.form.ClientWidth())
	//var by = w.form.Height() - w.form.ClientHeight()

	form.SetSize(int(width), int(height))
	form.SetPos(int(x), int(y)+int(height))

	popup = &Popup{
		form: form,
	}

	popup.widget.SetAllocation(int32(x), int32(y), int32(width), int32(height))
	return
}

func SurfaceEnter(wlSurface *wl.Surface, wlOutput *wl.Output) {
}
func SurfaceLeave(wlSurface *wl.Surface, wlOutput *wl.Output) {
}

type Popuper interface {
	Render(cairo.Surface, uint32)
	Done()
	Configure() *Widget
}

type Popup struct {
	Popup *xdg.Popup
	
	Display *Display

	form *winc.Form

	widget Widget

	popuper Popuper

	inhibited, configured bool
}

func (p *Popup) Destroy() {
	form := p.form
	if form != nil {
		p.form = nil
		form.Close()
	}
	p.widget.destroyed = true
	p.popuper = nil
	p.inhibited = true
	p.configured = true
}

func (p *Popup) SetPopupHandler(ph Popuper) {
	p.popuper = ph
}

func (p *Popup) PopupGetSurface() cairo.Surface {

	if !p.configured {
		p.popuper.Configure()
		p.configured = true
	}

	return &p.widget
}

func (p *Popup) BufferRelease(*wl.Buffer) {
}
