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

// Go Wayland SimpleShm demo
package main

import "github.com/neurlang/wayland/wl"
import zxdg "github.com/neurlang/wayland/xdg"
import "github.com/neurlang/wayland/wlclient"
import "log"
import "github.com/neurlang/wayland/os"

type display struct {
	display    *wl.Display
	registry   *wl.Registry
	compositor *wl.Compositor
	shell      *zxdg.WmBase
	shm        *wl.Shm
	hasXrgb    bool
}
type buffer struct {
	buffer  *wl.Buffer
	shmData []byte
	busy    bool
}

type window struct {
	display          *display
	width, height    int
	surface          *wl.Surface
	xdgSurface       *zxdg.Surface
	xdgToplevel      *zxdg.Toplevel
	buffers          [2]buffer
	callback         *wl.Callback
	waitForConfigure bool
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (mybuf *buffer) HandleBufferRelease(wl.BufferReleaseEvent) {
	mybuf.busy = false
}

func createShmBuffer(w *window, buffer *buffer, width int, height int, format uint32) error {

	stride := width * 4
	var size = stride * height

	fd, err := os.CreateAnonymousFile(int64(size))
	if err != nil {
		return err
	}

	defer fd.Close()

	data, err := os.Mmap(int(fd.Fd()), 0,
		size, os.ProtRead|os.ProtWrite, os.MapShared)
	if err != nil {
		return err
	}

	pool, err := w.display.shm.CreatePool(fd.Fd(), int32(size))
	if err != nil {
		return err
	}
	buf, err := pool.CreateBuffer(0, int32(width), int32(height), int32(stride), format)
	if err != nil {
		return err
	}

	buffer.buffer = buf
	// add buffer releaser here
	wlclient.BufferAddListener(buffer.buffer, buffer)
	pool.Destroy()

	buffer.shmData = data

	return nil
}
func (window *window) HandleSurfaceConfigure(ev zxdg.SurfaceConfigureEvent) {
	window.SurfaceConfigure(window.xdgSurface, ev.Serial)
}
func (window *window) SurfaceConfigure(sf *zxdg.Surface, serial uint32) {

	handle(sf.AckConfigure(serial))

	if window.waitForConfigure {
		redraw(window, nil, 0)
		window.waitForConfigure = false
	}

}

func (window *window) HandleToplevelConfigure(ev zxdg.ToplevelConfigureEvent) {

}

func (window *window) HandleToplevelClose(ev zxdg.ToplevelCloseEvent) {
}

func (d *display) HandleShmFormat(ev wl.ShmFormatEvent) {

	if ev.Format == wl.ShmFormatXrgb8888 {
		d.hasXrgb = true
	}

}

func (d *display) HandleWmBasePing(ev zxdg.WmBasePingEvent) {
	d.shell.Pong(ev.Serial)
}

func (d *display) HandleRegistryGlobal(ev wl.RegistryGlobalEvent) {
	d.RegistryGlobal(d.registry, ev.Name, ev.Interface, ev.Version)
}

func (d *display) RegistryGlobal(reg *wl.Registry, goid uint32, face string,
	version uint32) {

	goFace := face

	switch goFace {
	case "wl_compositor":
		d.compositor = wlclient.RegistryBindCompositorInterface(d.registry, goid, 1)

	case "xdg_wm_base":
		d.shell = wlclient.RegistryBindWmBaseInterface(d.registry, goid, 1)

		zxdg.WmBaseAddListener(d.shell, d)

	case "zwp_fullscreen_shell_v1":

	case "wl_shm":
		d.shm = wlclient.RegistryBindShmInterface(d.registry, goid, 1)
		wlclient.ShmAddListener(d.shm, d)

	default:
		log.Println("Other register global", goFace)
	}
}

func (*display) HandleRegistryGlobalRemove(ev wl.RegistryGlobalRemoveEvent) {

}

func createDisplay() *display {
	disp := &display{
		hasXrgb: false,
	}
	d, err := wlclient.DisplayConnect(nil)
	if err != nil {
		println("Could not connect to Wayland.")
		handle(err)
	}

	disp.display = d

	reg, err := disp.display.GetRegistry()
	if err != nil {
		println("Could not get Registry.")
		handle(err)
	}

	disp.registry = reg

	wlclient.RegistryAddListener(disp.registry, disp)
	handle(wlclient.DisplayRoundtrip(disp.display))

	if disp.shm == nil {
		log.Fatal("No wl_shm global\n")
	}

	handle(wlclient.DisplayRoundtrip(disp.display))

	if !disp.hasXrgb {
		log.Fatal("WL_SHM_FORMAT_XRGB32 not available\n")
	}

	return disp
}

func createWindow(disp *display, width, height int) *window {

	win := &window{
		callback: nil,
		display:  disp,
		width:    width,
		height:   height,
	}

	surf, err := disp.compositor.CreateSurface()
	if err != nil {
		println("cannot create surface")
		handle(err)
	}
	win.surface = surf

	if disp.shell != nil {
		xdgSurf, err := disp.shell.GetSurface(win.surface)
		if err != nil {
			println("cannot get xdg surface")
			handle(err)
		}
		win.xdgSurface = xdgSurf

		win.xdgSurface.AddListener(win)

		xdgToplevel, err := win.xdgSurface.GetToplevel()
		if err != nil {
			println("xdg toplevel is wrong")
			handle(err)
		}

		win.xdgToplevel = xdgToplevel

		zxdg.ToplevelAddListener(win.xdgToplevel, win)

		handle(win.xdgToplevel.SetTitle("Title"))

		handle(win.surface.Commit())
		win.waitForConfigure = true

	} else {
		log.Fatal("Unknown shell\n")
	}

	return win
}
func windowNextBuffer(window *window) *buffer {
	var buffer *buffer

	if !window.buffers[0].busy {
		buffer = &window.buffers[0]
	} else if !window.buffers[1].busy {
		buffer = &window.buffers[1]
	} else {
		log.Println("All Buffers Busy")
		return nil
	}

	if buffer.buffer == nil {
		err := createShmBuffer(window, buffer,
			window.width, window.height,
			wl.ShmFormatXrgb8888)
		if err != nil {
			log.Println(err)
			return nil
		}

		/* paint the padding */
		for i := range buffer.shmData {
			buffer.shmData[i] = 0xff
		}
	}

	return buffer
}

func (window *window) paintPixels(time uint32, buf *buffer) {
	var or int
	halfw, halfh := window.width/2, window.height/2
	if halfw < halfh {
		or = halfw - 8
	} else {
		or = halfh - 8
	}
	ir := or - 32
	or *= or
	ir *= ir

	var iter = 0

	for y := 0; y < window.height; y++ {
		y2 := (y - halfh) * (y - halfh)
		for x := 0; x < window.width; x++ {
			var v int
			r2 := (x-halfw)*(x-halfw) + y2
			if r2 < ir {
				v = (r2/32 + int(time)/64) * 0x0080401
			} else if r2 < or {
				v = (y + int(time)/32) * 0x0080401
			} else {
				v = (x + int(time)/16) * 0x0080401
			}
			a := 0
			r := v >> 16
			g := v >> 8
			b := v
			buf.shmData[iter] = byte(b)
			iter++
			buf.shmData[iter] = byte(g)
			iter++
			buf.shmData[iter] = byte(r)
			iter++
			buf.shmData[iter] = byte(a)
			iter++
		}
	}
}
func (window *window) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	window.CallbackDone(ev.C, ev.CallbackData)
}
func (window *window) CallbackDone(callback *wl.Callback, time uint32) {
	redraw(window, callback, time)
}

func redraw(win *window, callback *wl.Callback, time uint32) {
	if win.callback != nil {
		wlclient.CallbackDestroy(win.callback)
	}
	buffer := windowNextBuffer(win)

	win.paintPixels(time, buffer)
	handle(win.surface.Attach(buffer.buffer, 0, 0))
	handle(win.surface.Damage(0, 0,
		int32(win.width), int32(win.height)))

	if callback != nil {
		wlclient.CallbackDestroy(callback)
	}

	cb, err := win.surface.Frame()
	handle(err)
	win.callback = cb

	wlclient.CallbackAddListener(win.callback, win)
	handle(win.surface.Commit())
	buffer.busy = true
}

func main() {

	display := createDisplay()
	window := createWindow(display, 250, 250)
	if window == nil {
		return
	}

	// Initialise damage to full surface, so the padding gets painted
	handle(window.surface.Damage(0, 0, int32(window.width), int32(window.height)))

	if !window.waitForConfigure {
		redraw(window, nil, 0)
	}

	for wlclient.DisplayDispatch(display.display) == nil {
	}

	print("go-wayland-simpleshm exiting\n")

}
