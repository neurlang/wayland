package window

import (
	cairo "github.com/neurlang/wayland/cairoshim"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/xdg"
	"github.com/neurlang/winc"
	"github.com/neurlang/winc/w32"
	"time"
	"sync"
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

const patch_dim_x = 512
const patch_dim_y = 512

type rect struct {
    x0, y0, x1, y1 int  // Coordinates of the rectangle
    r, g, b        byte // Color of the rectangle
}



func reduceDrawRectsNew(oldRects []rect, buffer []byte, x0, y0, width, height int) (ret, oldok []rect) {
    done := make([][]bool, height)
    for i := range done {
        done[i] = make([]bool, width)
    }
outer:
    for _, oldRect := range oldRects {
    	//continue
        for y := oldRect.y0; y <= oldRect.y1 && y < height; y++ {
            for x := oldRect.x0; x <= oldRect.x1 && x < width; x++ {
         	 r, g, b := getColor(buffer, x, y, width)
         	 if r != oldRect.r || g != oldRect.g || b != oldRect.b {
         	 	continue outer
         	 }
            }
        }
        for y := oldRect.y0; y <= oldRect.y1 && y < height; y++ {
            for x := oldRect.x0; x <= oldRect.x1 && x < width; x++ {
		 done[y][x] = true
            }
        }
        oldok = append(oldok, oldRect)
    }
    
    
        for y := y0; y < y0 + patch_dim_y && y < height; y++ {
            for x := x0; x < x0 + patch_dim_x && x < width; x++ {
                if done[y][x] {continue;}
         	 r, g, b := getColor(buffer, x, y, width)
          	rect, _ := findMaxRectNew(buffer, done, x0, y0, x, y, r, g, b, width, height)
          	ret = append(ret, rect)
          	
          	
		// Step 2: Mark pixels of the largest rectangle as done
		for y := rect.y0; y <= rect.y1; y++ {
		    for x := rect.x0; x <= rect.x1; x++ {
		        done[y][x] = true
		    }
		}
            
        }}
        return

}

// IterateCoordinates iterates over coordinates and applies the callback function
func iterateCoordinates(n int, callback func(x, y int) bool) {

	for i := 1; i < 2*n; i++ {
		for row := 0; row <= i>>1; row++ {
		
			if !callback(row, i>>1) {
				return
			}
		}
		i++
		for col := 0; col < i>>1; col++ {

			if !callback(i>>1, col) {
				return
			}
		}
	}
}
// Function to find the largest rectangle or run for a given color starting from (x0, y0)
func findMaxRectNew(buffer []byte, done [][]bool, x0x, y0y, x0, y0 int, r, g, b byte, width, height int) (rect, int) {
    largestRect := rect{}
    largestArea := 0
    largestAffected := 0

    // Step 1: Find the longest horizontal run
    maxX := x0
    horizontalAffected := 0
    for maxX < x0x + patch_dim_x && maxX < width && (!done[y0][maxX] && (sameColor(buffer, maxX, y0, r, g, b, width))) {
        if !done[y0][maxX] {
            horizontalAffected++
        }
        maxX++
    }
    maxX--

    horizontalRect := rect{x0, y0, maxX, y0, r, g, b}

    // Update largest rect based on horizontal run
    if horizontalAffected > largestArea {
        largestRect = horizontalRect
        largestArea = horizontalAffected
        largestAffected = horizontalAffected
    }

    // Step 2: Find the longest vertical run
    maxY := y0
    verticalAffected := 0
    for maxY < y0y + patch_dim_y && maxY < height && (!done[maxY][x0] && (sameColor(buffer, x0, maxY, r, g, b, width))) {
        if !done[maxY][x0] {
            verticalAffected++
        }
        maxY++
    }
    maxY--

    verticalRect := rect{x0, y0, x0, maxY, r, g, b}

    // Update largest rect based on vertical run
    if verticalAffected > largestArea {
        largestRect = verticalRect
        largestArea = verticalAffected
        largestAffected = verticalAffected
    }

    // Step 3: Try larger rectangle (covering horizontal and vertical area)
    var left, right, up, down int = width-1, x0, height-1, y0
    accumAffected := 0
    accumAdding := 0

    iterateCoordinates(patch_dim_y, func(xd, yd int) bool {
        x := x0 + xd
        y := y0 + yd
        corner := xd == yd || xd == yd + 1

        if x >= width || y >= height {
            return false
        }
        if x >= x0x + patch_dim_x || y >= y0y + patch_dim_y {
            return false
        }
        
        if done[y][x] || !sameColor(buffer, x, y, r, g, b, width) {
            return false
        }
        

        // If the pixel is part of the rectangle, include it in affected count
        {
            if left > x {
	    		left = x
	    		}
	    if up > y {
	    		up = y
	    	}
            accumAdding++
        }
        

        if corner {
	    	accumAffected += accumAdding
	    	accumAdding = 0
            	// Update the bounds for the rectangle
            	right = x
            	down = y
        }

        return true
    })

    // Only return the larger rectangle if it expands beyond the initial rect
    if accumAffected > 0 && right > left && down > up {
        return rect{left, up, right, down, r, g, b}, accumAffected
    }

    return largestRect, largestAffected
}


// Helper function to compare pixel color at two locations
func sameColor(buffer []byte, x1, y1 int, r, g, b byte, width int) bool {
    r1, g1, b1 := getColor(buffer, x1, y1, width)
    return r == r1 && g == g1 && b == b1
}

// Function to get pixel color at (x, y)
func getColor(buffer []byte, x, y, width int) (byte, byte, byte) {
    index := (y*width + x) * BUFFER_BYTES // Assuming 3 bytes per pixel (RGB)
    return buffer[index], buffer[index+1], buffer[index+2]
}



//func timeTrack(start time.Time, name string) {
//    elapsed := time.Since(start)
//    println(name, "took", elapsed.String())
//}

var mut sync.RWMutex
func redrawer(widget *Widget, canvas *winc.Canvas) {
	//defer timeTrack(time.Now(), "redrawer")

	buf, w, h, hash := widget.getBufferAndAllocAndHash()
	if hash == widget.drawnHash || w <= 0 || h <= 0 {
		return
	}
	//println(widget.allocation_width * widget.allocation_height * BUFFER_BYTES, len(widget.buffer))
	wg := sync.WaitGroup{}
	var drawn = make(map[int]uint64)
	var oks = make(map[[2]int][]rect)

	for y := 0; y < h; y += patch_dim_y {
		end := y + patch_dim_y
		if end > h {
			end = h
		}
		hash := hashBuffer(buf, y, end, w)
		mut.RLock()
		if val, ok := widget.drawnHashes[y]; ok && val == hash {
			mut.RUnlock()
			continue
		}
		mut.RUnlock()
		drawn[y] = hash

	for x := 0; x < w; x += patch_dim_x {
		wg.Add(1)
		go func(x, y int) {
			mut.RLock()
			oldRects := widget.drawnRects[[2]int{x, y}]
			_ = oldRects
			//if len(oldRects) > 1024 {
			//	delete(widget.drawnRects, [2]int{x, y})
			//}
			mut.RUnlock()

			
			rects, oldOks := reduceDrawRectsNew(oldRects, buf, x, y, w, h)

			if len(rects) > 0 {

				var lastR, lastG, lastB byte = rects[0].r, rects[0].g, rects[0].b //first color

				brush := winc.NewSolidColorBrush(winc.RGB(lastB, lastG, lastR))
				//pen := winc.NewPen(1, 1, brush)
				
				for _, rectangle := range rects {
					
				
					if rectangle.r != lastR || rectangle.g != lastG || rectangle.b != lastB {
						brush.Dispose()
						//pen.Dispose()
						lastR = rectangle.r
						lastG = rectangle.g
						lastB = rectangle.b
						brush = winc.NewSolidColorBrush(winc.RGB(lastB, lastG, lastR))
						//pen = winc.NewPen(1, 1, brush)
					}

						// Draw a rectangle
					rect := winc.NewRect(
						rectangle.x0,
						rectangle.y0,
						rectangle.x1+1,
						rectangle.y1+1)
					mut.Lock()
					canvas.FillRect(rect, brush)
					mut.Unlock()
						
				}
				brush.Dispose()
			}

			for _, r := range rects {
				if r.x0 != r.x1 || r.y0 != r.y1 {
					oldOks = append(oldOks, r)
				}
			}
			mut.Lock()
			oks[[2]int{x, y}] = oldOks
			mut.Unlock()
			wg.Done()
		}(x, y)
	}}
	
	wg.Wait()
	
	widget.setHashHashesRects(hash, drawn, oks)

	
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
		widget.drawnRects = nil
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
