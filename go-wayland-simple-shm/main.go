// Copyright 2021 Neurlang project

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

//Go Wayland SimpleShm demo
package main

import wl "github.com/neurlang/wayland/wl"
import zxdg "github.com/neurlang/wayland/xdg"
import wlclient "github.com/neurlang/wayland/wlclient"
import "fmt"
import "github.com/neurlang/wayland/os"

type display struct {
	display    *wl.Display
	registry   *wl.Registry
	compositor *wl.Compositor
	shell      *zxdg.Shell
	shm        *wl.Shm
	has_xrgb   bool
}
type buffer struct {
	buffer   *wl.Buffer
	shm_data []uint32
	busy     bool
}

type window struct {
	display            *display
	width, height      int
	surface            *wl.Surface
	xdg_surface        *zxdg.Surface
	xdg_toplevel       *zxdg.Toplevel
	buffers            [2]buffer
	callback           *wl.Callback
	wait_for_configure bool
}

func (mybuf *buffer) HandleBufferRelease(wl.BufferReleaseEvent) {
	mybuf.busy = false
}

func create_shm_buffer(w *window, buffer *buffer, width int, height int, format uint32) error {

	stride := width * 4
	var size = stride * height

	fd, err := os.CreateAnonymousFile(int64(size))
	if err != nil {
		return err
	}

	defer os.Close(fd)

	data, err := os.Mmap32(fd, 0,
		size, os.PROT_READ|os.PROT_WRITE, os.MAP_SHARED)
	if err != nil {
		return err
	}

	pool, err := w.display.shm.CreatePool(uintptr(fd), int32(size))
	if err != nil {
		return err
	}
	buf, err := pool.CreateBuffer(0, int32(width), int32(height), int32(stride), uint32(format))
	if err != nil {
		return err
	}

	buffer.buffer = buf
	// add buffer releaser here
	wlclient.BufferAddListener((buffer.buffer), buffer)
	pool.Destroy()

	buffer.shm_data = data

	return nil
}
func (window *window) HandleXdgSurfaceConfigure(ev zxdg.XdgSurfaceConfigureEvent) {
	window.SurfaceConfigure(window.xdg_surface, ev.Serial)
}
func (window *window) SurfaceConfigure(sf *zxdg.Surface, serial uint32) {

	sf.AckConfigure(serial)

	if window.wait_for_configure {
		redraw(window, nil, 0)
		window.wait_for_configure = false
	}

}

func (window *window) HandleToplevelConfigure(ev zxdg.ToplevelConfigureEvent) {

}

func (window *window) HandleToplevelClose(ev zxdg.ToplevelCloseEvent) {
}

func (d *display) HandleShmFormat(ev wl.ShmFormatEvent) {

	if ev.Format == wl.ShmFormatXrgb8888 {
		d.has_xrgb = true
	}

}

func (d *display) HandleXdgWmBasePing(ev zxdg.XdgWmBasePingEvent) {
	d.shell.Pong(ev.Serial)
}

func (disp *display) HandleRegistryGlobal(ev wl.RegistryGlobalEvent) {
	disp.RegistryGlobal(disp.registry, ev.Name, ev.Interface, ev.Version)
}

func (disp *display) RegistryGlobal(reg *wl.Registry, goid uint32, face string,
	version uint32) {

	goFace := face

	switch goFace {
	case "wl_compositor":
		disp.compositor = wlclient.RegistryBindCompositorInterface(disp.registry, goid, 1)

	case "zxdg_shell_v6":
		disp.shell = wlclient.RegistryBindShellInterface(disp.registry, goid, 1)

		zxdg.WmBaseAddListener((disp.shell), disp)

	case "zwp_fullscreen_shell_v1":

	case "wl_shm":
		disp.shm = wlclient.RegistryBindShmInterface(disp.registry, goid, 1)
		wlclient.ShmAddListener((disp.shm), disp)

	default:
		fmt.Println("Other register global", goFace)
	}
}

func (*display) HandleRegistryGlobalRemove(ev wl.RegistryGlobalRemoveEvent) {

}

func create_display() *display {
	disp := &display{
		has_xrgb: false,
	}
	d, err := wlclient.DisplayConnect(nil)
	if err != nil {
		panic("Could not connect to Wayland.")
	}

	disp.display = d

	reg, err := disp.display.GetRegistry()
	if err != nil {
		panic("Could not get Registry.")
	}

	disp.registry = reg

	wlclient.RegistryAddListener((disp.registry), disp)
	wlclient.DisplayRoundtrip(disp.display)
	if disp.shm == nil {
		panic("No wl_shm global\n")
	}

	wlclient.DisplayRoundtrip(disp.display)

	if !disp.has_xrgb {
		panic("WL_SHM_FORMAT_XRGB32 not available\n")
	}

	return disp
}

func create_window(disp *display, width, height int) *window {

	win := &window{
		callback: nil,
		display:  disp,
		width:    width,
		height:   height,
	}

	surf, err := disp.compositor.CreateSurface()
	if err != nil {
		panic("cannot create surface")
	}
	win.surface = surf

	if disp.shell != nil {
		xdg_surf, err := disp.shell.GetXdgSurface(win.surface)
		if err != nil {
			panic("cannot get xdg surface")
		}
		win.xdg_surface = xdg_surf

		win.xdg_surface.AddListener(win)

		xdg_toplevel, err := win.xdg_surface.GetToplevel()
		if err != nil {
			panic("xdg toplevel is wrong")
		}

		win.xdg_toplevel = xdg_toplevel

		zxdg.ToplevelAddListener((win.xdg_toplevel), win)

		win.xdg_toplevel.SetTitle("Title")

		win.surface.Commit()
		win.wait_for_configure = true

	} else {
		print("Unknown shell\n")
		panic("")
	}

	return win
}
func window_next_buffer(window *window) *buffer {
	var buffer *buffer

	if !window.buffers[0].busy {
		buffer = &window.buffers[0]
	} else if !window.buffers[1].busy {
		buffer = &window.buffers[1]
	} else {
		print("All Buffers Busy\n")
		return nil
	}

	if buffer.buffer == nil {
		err := create_shm_buffer(window, buffer,
			window.width, window.height,
			wl.ShmFormatXrgb8888)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		/* paint the padding */
		for i := range buffer.shm_data {
			buffer.shm_data[i] = uint32(0xffffffff)
		}
	}

	return buffer
}

func (w *window) paintPixels(time uint32, buf *buffer) {
	var or int
	halfw, halfh := w.width/2, w.height/2
	if halfw < halfh {
		or = halfw - 8
	} else {
		or = halfh - 8
	}
	ir := or - 32
	or *= or
	ir *= ir

	var iter = 0

	for y := 0; y < w.height; y++ {
		y2 := (y - halfh) * (y - halfh)
		for x := 0; x < w.width; x++ {
			var v int
			r2 := (x-halfw)*(x-halfw) + y2
			if r2 < ir {
				v = (r2/32 + int(time)/64) * 0x0080401
			} else if r2 < or {
				v = (y + int(time)/32) * 0x0080401
			} else {
				v = (x + int(time)/16) * 0x0080401
			}
			v &= 0x00ffffff

			buf.shm_data[iter] = uint32(v)
			iter = iter + 1
		}
	}
}
func (win *window) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	win.CallbackDone(ev.C, ev.CallbackData)
}
func (win *window) CallbackDone(callback *wl.Callback, time uint32) {
	redraw(win, callback, time)
}

func redraw(win *window, callback *wl.Callback, time uint32) {

	buffer := window_next_buffer(win)

	win.paintPixels(uint32(time), buffer)
	win.surface.Attach(buffer.buffer, 0, 0)
	win.surface.Damage(0, 0,
		int32(win.width), int32(win.height))

	if callback != nil {
		wlclient.CallbackDestroy(callback)
	}

	cb, err := win.surface.Frame()
	if err != nil {
		panic("cannot get frame callback")
	}
	win.callback = cb

	wlclient.CallbackAddListener(win.callback, win)
	win.surface.Commit()
	buffer.busy = true
}

func main() {

	display := create_display()
	window := create_window(display, 250, 250)
	if window == nil {
		return
	}

	// Initialise damage to full surface, so the padding gets painted
	window.surface.Damage((0), (0), int32(window.width), int32(window.height))

	if !window.wait_for_configure {
		redraw(window, nil, 0)
	}

	for wlclient.DisplayDispatch(display.display) == nil {
	}

	print("go-wayland-simpleshm exiting\n")

}
