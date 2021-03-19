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

import wl "github.com/neurlang/wayland/wayland"
import zxdg "github.com/neurlang/wayland/wayland"
import zwp "github.com/neurlang/wayland/wayland"
import ivi "github.com/neurlang/wayland/wayland"
import "fmt"
import "github.com/neurlang/wayland/os"

type display struct {
	display         *wl.Display
	registry        *wl.Registry
	compositor      *wl.Compositor
	shell           *zxdg.ShellV6
	fshell          *zwp.FullscreenShellV1
	shm             *wl.Shm
	has_xrgb        bool
	ivi_application *ivi.IviApplication
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
	xdg_surface        *zxdg.SurfaceV6
	xdg_toplevel       *zxdg.ToplevelV6
	ivi_surface        *ivi.IviSurface
	buffers            [2]buffer
	prev_buffer_id     int
	callback           *wl.Callback
	wait_for_configure bool
}

func (mybuf *buffer) BufferRelease(buf *wl.Buffer) {
	mybuf.busy = false
}

func create_shm_buffer(w *window, buffer *buffer, width int, height int, format uint32) error {

	stride := width * 4
	var size = stride * height

	fd, err := os.CreateAnonymousFile(int64(size))
	if err != nil {
		return err
	}

	data, err := os.Mmap32(fd, 0,
		size, os.PROT_READ|os.PROT_WRITE, os.MAP_SHARED)
	if err != nil {
		os.Close(fd)
		return err
	}

	pool := wl.ShmCreatePool((w.display.shm), int32(fd), int32(size))
	buffer.buffer = wl.ShmPoolCreateBuffer(pool, 0, int32(width), int32(height), int32(stride), uint32(format))

	// add buffer releaser here
	wl.BufferAddListener((buffer.buffer), buffer)
	wl.ShmPoolDestroy(pool)
	os.Close(fd)

	buffer.shm_data = data

	return nil
}

func (window *window) SurfaceV6Configure(sf *zxdg.SurfaceV6, serial uint32) {

	zxdg.SurfaceV6AckConfigure(sf, serial)

	if window.wait_for_configure {
		redraw(window, nil, 0)
		window.wait_for_configure = false
	}

}

func (window *window) ToplevelV6Configure(zxdg_toplevel_v6 *zxdg.ToplevelV6,
	width int32, height int32, states []int32) {

}

func (window *window) ToplevelV6Close(tl *zxdg.ToplevelV6) {
}

func (d *display) ShmFormat(shm *wl.Shm, format uint32) {

	if format == wl.SHM_FORMAT_XRGB8888 {
		d.has_xrgb = true
	}

}

func (d *display) ShellV6Ping(shell *zxdg.ShellV6, serial uint32) {
	zxdg.ShellV6Pong(shell, serial)
}

func (disp *display) RegistryGlobal(reg *wl.Registry, goid uint32, face string,
	version uint32) {

	goFace := face

	switch goFace {
	case "wl_compositor":
		disp.compositor = wl.RegistryBindCompositorInterface(disp.registry, goid, 1)

	case "zxdg_shell_v6":
		disp.shell = zxdg.RegistryBindShellV6Interface(disp.registry, goid, 1)

		zxdg.ShellV6AddListener((disp.shell), disp)

	case "zwp_fullscreen_shell_v1":
		disp.fshell = zwp.RegistryBindFullscreenShellV1Interface(disp.registry, goid, 1)
	case "wl_shm":
		disp.shm = wl.RegistryBindShmInterface(disp.registry, goid, 1)
		wl.ShmAddListener((disp.shm), disp)

	case "ivi_application":
		disp.ivi_application = ivi.RegistryBindIviApplicationInterface(disp.registry, goid, 1)

	default:
		fmt.Println("Other register global", goFace)
	}
}

func (*display) RegistryGlobalRemove(r *wl.Registry, u uint32) {

}

func create_display() *display {
	disp := &display{
		has_xrgb: false,
	}
	disp.display = wl.DisplayConnect(nil)
	fmt.Printf("Type of wl_display value: %T\n", disp.display)

	if disp.display == nil {
		panic("Could not connect to Wayland.")
	}

	disp.registry = (wl.DisplayGetRegistry(disp.display))
	wl.RegistryAddListener((disp.registry), disp)
	wl.DisplayRoundtrip(disp.display)
	if disp.shm == nil {
		panic("No wl_shm global\n")
	}

	wl.DisplayRoundtrip(disp.display)

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

	win.surface = wl.CompositorCreateSurface((disp.compositor))

	if disp.shell != nil {
		win.xdg_surface =
			zxdg.ShellV6GetXdgSurface((disp.shell),
				(win.surface))
		if win.xdg_surface == nil {
			panic("xdg surface is nil")
		}

		zxdg.SurfaceV6AddListener((win.xdg_surface), win)

		win.xdg_toplevel = zxdg.SurfaceV6GetToplevel((win.xdg_surface))

		if win.xdg_toplevel == nil {
			panic("xdg surface is nil")
		}

		zxdg.ToplevelV6AddListener((win.xdg_toplevel), win)

		zxdg.ToplevelV6SetTitle((win.xdg_toplevel), []byte("Title"))

		wl.SurfaceCommit((win.surface))
		win.wait_for_configure = true

	} else if disp.fshell != nil {

		print("Fshell\n")

	} else if disp.ivi_application != nil {

		print("Ivi shell\n")

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
			wl.SHM_FORMAT_XRGB8888)
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

func (win *window) CallbackDone(callback *wl.Callback, time uint32) {
	redraw(win, callback, time)
}

func redraw(win *window, callback *wl.Callback, time uint32) {

	buffer := window_next_buffer(win)

	win.paintPixels(uint32(time), buffer)
	wl.SurfaceAttach(win.surface, (buffer.buffer), 0, 0)
	wl.SurfaceDamage(win.surface, 0, 0,
		int32(win.width), int32(win.height))

	if callback != nil {
		wl.CallbackDestroy(callback)
	}

	win.callback = wl.SurfaceFrame(win.surface)
	wl.CallbackAddListener((win.callback), win)
	wl.SurfaceCommit(win.surface)
	buffer.busy = true
}

func main() {

	display := create_display()
	window := create_window(display, 250, 250)
	if window == nil {
		return
	}

	// Initialise damage to full surface, so the padding gets painted
	wl.SurfaceDamage(window.surface, int32(0), int32(0), int32(window.width), int32(window.height))

	if !window.wait_for_configure {
		redraw(window, nil, 0)
	}

	for wl.DisplayDispatch(display.display) != -1 {
	}

	print("go-wayland-simpleshm exiting\n")

}
