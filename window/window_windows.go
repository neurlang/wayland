package window

import (
	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/winc"
	"github.com/neurlang/winc/w32"
	"time"
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
	form.OnClose().Bind(func(arg *winc.Event) {
		winc.Exit()
	})

	return &Window{form: form, parent_display: d}
}

func (w *Window) SetKeyboardHandler(t KeyboardHandler) {

	w.input = &Input{}

	allRedrawer := func() {
		for widget := range w.widgets {
			if !w.inhibited {
				widget.destroyed = false
				widget.ScheduleRedraw()
				redrawer(widget, winc.NewCanvasFromHwnd(w.form.Handle()))
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
func redrawer(widget *Widget, canvas *winc.Canvas) {
outer:
	for i := 0; i < len(widget.buffer); i += BUFFER_BYTES {
		now := i / BUFFER_BYTES
		var r = widget.buffer[i+2]
		var g = widget.buffer[i+1]
		var b = widget.buffer[i]

		for j := i + BUFFER_BYTES; j <= len(widget.buffer); j += BUFFER_BYTES {
			now2 := (j - BUFFER_BYTES) / BUFFER_BYTES
			next2 := j / BUFFER_BYTES
			var same bool
			if j < len(widget.buffer) {
				var rNext = widget.buffer[j+2]
				var gNext = widget.buffer[j+1]
				var bNext = widget.buffer[j]
				same = r == rNext && g == gNext && b == bNext
			} else {
				same = false
			}
			if !same || (now/widget.allocation_width != next2/widget.allocation_width) {
				var cleared bool
				for clearer := now + 1; clearer < now2; clearer++ {
					if _, ok := widget.drawn[clearer]; ok {
						delete(widget.drawn, clearer)
						cleared = true
					}
					if _, ok := widget.drawn2[clearer]; ok {
						delete(widget.drawn2, clearer)
						cleared = true
					}

				}

				if cleared || (widget.drawn[now] != [4]byte{r, g, b, 1} || widget.drawn2[now2] != [4]byte{r, g, b, 1}) {

					widget.drawn[now] = [4]byte{r, g, b, 1}
					widget.drawn2[now2] = [4]byte{r, g, b, 1}

					// Draw a rectangle
					rect := winc.NewRect(
						now%widget.allocation_width+0,
						now/widget.allocation_width+0,
						now2%widget.allocation_width+1,
						now/widget.allocation_width+1)

					brush := winc.NewSolidColorBrush(winc.RGB(r, g, b))
					canvas.FillRect(rect, brush)
					brush.Dispose()
					i = j - BUFFER_BYTES
					continue outer
				} else {
					i = j - BUFFER_BYTES
					continue outer
				}
			}
		}
	}
}

func (w *Window) AddWidget(t WidgetHandler) (widget *Widget) {
	widget = &Widget{
		drawn:         make(map[int][4]byte),
		drawn2:        make(map[int][4]byte),
		parent_window: w,
		handler:       t,
	}
	if w.widgets == nil {
		w.widgets = make(map[*Widget]struct{})
	}
	w.widgets[widget] = struct{}{}

	canvasRedrawer := func(canvas *winc.Canvas) {
		redrawer(widget, canvas)
	}

	allRedrawer := func() {
		if !w.inhibited {
			t.Redraw(widget)
			redrawer(widget, winc.NewCanvasFromHwnd(w.form.Handle()))
		}
	}
	w.form.OnPaint().Bind(func(arg *winc.Event) {

		if !w.inhibited && !widget.destroyed {
			t.Redraw(widget)
			canvasRedrawer(arg.Data.(*winc.PaintEventData).Canvas)
		}

	})

	w.form.OnSize().Bind(func(arg *winc.Event) {
		widget.destroyed = false

		xy := arg.Data.(*winc.SizeEventData)

		for widget := range w.widgets {
			widget.SetAllocation(0, 0, int32(xy.X), int32(xy.Y))
		}

		t.Resize(widget, int32(xy.X), int32(xy.Y), int32(xy.X), int32(xy.Y))
	})

	w.form.OnMouseMove().Bind(func(arg *winc.Event) {
		xy := arg.Data.(*winc.MouseEventData)
		t.Motion(widget, w.input, uint32(time.Now().UnixNano()/1000000), float32(xy.X), float32(xy.Y))
		allRedrawer()
	})
	w.form.OnMouseHover().Bind(func(arg *winc.Event) {
		xy := arg.Data.(*winc.MouseEventData)

		t.Motion(widget, w.input, uint32(time.Now().UnixNano()/1000000), float32(xy.X), float32(xy.Y))
		allRedrawer()
	})

	w.form.OnLBDown().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 272, wl.PointerButtonStatePressed, t)
		allRedrawer()
	})
	w.form.OnLBUp().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 272, wl.PointerButtonStateReleased, t)
		allRedrawer()
	})
	w.form.OnRBDown().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 273, wl.PointerButtonStatePressed, t)
		allRedrawer()
	})
	w.form.OnRBUp().Bind(func(arg *winc.Event) {
		t.Button(widget, w.input, uint32(time.Now().UnixNano()/1000000), 273, wl.PointerButtonStateReleased, t)
		allRedrawer()
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

	for widget := range w.widgets {
		w.parent_display.mustResize = append(w.parent_display.mustResize, &mustResize{
			widget,
			width, height,
		})
	}

}

func (w *Window) AddPopupWidget(p *Popup, handler WidgetHandler) *Widget {
	p.widget.handler = handler

	canvasRedrawer := func(canvas *winc.Canvas) {
		redrawer(&p.widget, canvas)
	}
	allRedrawer := func() {
		if !w.inhibited {
			if p.form != nil {
				if p.Popuper != nil {
					p.Popuper.Render(&p.widget, 0)
				}
				redrawer(&p.widget, winc.NewCanvasFromHwnd(p.form.Handle()))
			}
		}
	}
	p.form.OnPaint().Bind(func(arg *winc.Event) {

		if !p.inhibited && !p.widget.destroyed {
			p.Popuper.Render(&p.widget, 0)
			canvasRedrawer(arg.Data.(*winc.PaintEventData).Canvas)
		}

	})
	p.form.OnSize().Bind(func(arg *winc.Event) {
		p.widget.destroyed = false

		xy := arg.Data.(*winc.SizeEventData)

		p.widget.SetAllocation(0, 0, int32(xy.X), int32(xy.Y))

		p.Popuper.Configure()
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

	popup.widget.SetAllocation(int(x), int(y), int32(width), int32(height))
	return
}

type Popuper interface {
	Render(cairo.Surface, uint32)
	Done()
	Configure() *Widget
}

type Popup struct {
	form *winc.Form

	widget Widget

	Popuper Popuper

	inhibited, configured bool
}

func (p *Popup) Destroy() {
	form := p.form
	if form != nil {
		p.form = nil
		form.Close()
	}
	p.widget.destroyed = true
	p.Popuper = nil
	p.inhibited = true
	p.configured = true
}

func (p *Popup) PopupGetSurface() cairo.Surface {

	if !p.configured {
		p.Popuper.Configure()
		p.configured = true
	}

	return &p.widget
}
