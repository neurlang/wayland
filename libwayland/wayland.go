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

// Package libwayland implements a wayland protocol using wayland-client and
// wayland-cursor C libraries
package libwayland

/*
#cgo pkg-config: wayland-client wayland-cursor

#include <wayland-client.h>
#include <wayland-cursor.h>
#include <xdg-shell-client-protocol.h>
#include "fullscreen-shell-unstable-v1-client-protocol.h"

#include "bridge.h"

#include <errno.h>
*/
import "C"
import "unsafe"
import "sync"

// Interfaces

type Display C.struct_wl_display
type Registry C.struct_wl_registry
type Callback C.struct_wl_callback
type Compositor C.struct_wl_compositor
type ShmPool C.struct_wl_shm_pool
type Shm C.struct_wl_shm
type Buffer C.struct_wl_buffer
type DataOffer C.struct_wl_data_offer
type DataSource C.struct_wl_data_source
type DataDevice C.struct_wl_data_device
type DataDeviceManager C.struct_wl_data_device_manager
type Shell C.struct_wl_shell
type ShellSurface C.struct_wl_shell_surface
type Surface C.struct_wl_surface
type Seat C.struct_wl_seat
type Pointer C.struct_wl_pointer
type Keyboard C.struct_wl_keyboard
type Touch C.struct_wl_touch
type Output C.struct_wl_output
type Region C.struct_wl_region
type Subcompositor C.struct_wl_subcompositor
type Subsurface C.struct_wl_subsurface

// ZXDG interfaces

type XdgToplevel C.struct_xdg_toplevel
type XdgSurface C.struct_xdg_surface
type XdgWmBase C.struct_xdg_wm_base
type XdgPopup C.struct_xdg_popup

// ZWP interfaces

type RelativePointerV1 C.struct_zwp_relative_pointer_v1
type FullscreenShellV1 C.struct_zwp_fullscreen_shell_v1
type LockedPointerV1 C.struct_zwp_locked_pointer_v1
type ConfinedPointerV1 C.struct_zwp_confined_pointer_v1

// Permanent callback ///////////////////////////////////////////////////////////

type BufferListener interface {
	BufferRelease(buf *Buffer)
}

var buffer_listeners = []BufferListener{}

//export wlcallback_handle_buffer_release
func wlcallback_handle_buffer_release(data uintptr, buf *C.struct_wl_buffer) {
	buffer_listeners[data].BufferRelease((*Buffer)(buf))
}

func BufferAddListener(wl_buffer *Buffer, releaser BufferListener) {
	var i = uintptr(len(buffer_listeners))
	buffer_listeners = append(buffer_listeners, releaser)

	C._wl_buffer_add_listener((*C.struct_wl_buffer)(unsafe.Pointer(wl_buffer)), unsafe.Pointer(i))
}

//

type XdgToplevelListener interface {
	XdgToplevelConfigure(zxdg_toplevel *XdgToplevel, width int32, height int32, states []int32)
	XdgToplevelClose(zxdg_toplevel *XdgToplevel)
}

var callback_zxdg_toplevel_mutex sync.Mutex
var callback_zxdg_toplevel_listeners = []XdgToplevelListener{}

//export xdgcallback_handle_toplevel_configure
func xdgcallback_handle_toplevel_configure(data uintptr, cb *C.struct_xdg_toplevel,
	width int32, height int32, statptr unsafe.Pointer, statlen int) {
	callback_zxdg_toplevel_mutex.Lock()
	var configurer = callback_zxdg_toplevel_listeners[data]
	callback_zxdg_toplevel_mutex.Unlock()

	var stat []int32

	for i := 0; i < statlen; i += 4 {
		var statvalue = (*int32)((unsafe.Pointer)(uintptr(statptr) + uintptr(i)))
		stat = append(stat, *statvalue)
	}

	configurer.XdgToplevelConfigure((*XdgToplevel)(cb), width, height, stat)
}

//export xdgcallback_handle_toplevel_close
func xdgcallback_handle_toplevel_close(data uintptr, cb *C.struct_xdg_toplevel) {
	callback_zxdg_toplevel_mutex.Lock()
	var closer = callback_zxdg_toplevel_listeners[data]
	callback_zxdg_toplevel_mutex.Unlock()

	closer.XdgToplevelClose((*XdgToplevel)(cb))
}

func XdgToplevelAddListener(cb *XdgToplevel, confcloser XdgToplevelListener) {
	callback_zxdg_toplevel_mutex.Lock()
	var i = uintptr(len(callback_zxdg_toplevel_listeners))
	callback_zxdg_toplevel_listeners = append(callback_zxdg_toplevel_listeners, confcloser)
	callback_zxdg_toplevel_mutex.Unlock()

	C._xdg_toplevel_add_listener(unsafe.Pointer(cb), unsafe.Pointer(i))
}

//

type XdgSurfaceListener interface {
	XdgSurfaceConfigure(zxdg_surface *XdgSurface, serial uint32)
}

var callback_zxdg_surface_mutex sync.Mutex
var callback_zxdg_surface_listeners = []XdgSurfaceListener{}

//export xdgcallback_handle_surface_configure
func xdgcallback_handle_surface_configure(data uintptr, cb *C.struct_xdg_surface, serial uint32) {
	callback_zxdg_surface_mutex.Lock()
	var configurer = callback_zxdg_surface_listeners[data]
	callback_zxdg_surface_mutex.Unlock()

	_ = configurer

	configurer.XdgSurfaceConfigure((*XdgSurface)(cb), serial)
}

func (wl_callback *XdgSurface) AddListener(confer XdgSurfaceListener) {
	callback_zxdg_surface_mutex.Lock()
	var i = uintptr(len(callback_zxdg_surface_listeners))
	callback_zxdg_surface_listeners = append(callback_zxdg_surface_listeners, confer)
	callback_zxdg_surface_mutex.Unlock()

	C._xdg_surface_add_listener(unsafe.Pointer(wl_callback), unsafe.Pointer(i))
}

//

type XdgWmBaseListener interface {
	XdgWmBasePing(shell *XdgWmBase, serial uint32)
}

var callback_zxdg_shell_mutex sync.Mutex
var callback_zxdg_shell_listeners = []XdgWmBaseListener{}

//export xdgcallback_handle_shell_ping
func xdgcallback_handle_shell_ping(data uintptr, cb *C.struct_xdg_wm_base, serial uint32) {
	callback_zxdg_shell_mutex.Lock()
	var pinger = callback_zxdg_shell_listeners[data]
	callback_zxdg_shell_mutex.Unlock()

	pinger.XdgWmBasePing((*XdgWmBase)(cb), serial)
}

func XdgWmBaseAddListener(wl_callback *XdgWmBase, pinger XdgWmBaseListener) {
	callback_zxdg_shell_mutex.Lock()
	var i = uintptr(len(callback_zxdg_shell_listeners))
	callback_zxdg_shell_listeners = append(callback_zxdg_shell_listeners, pinger)
	callback_zxdg_shell_mutex.Unlock()

	C._xdg_wm_base_add_listener(unsafe.Pointer(wl_callback), unsafe.Pointer(i))
}

//

type RegistryListener interface {
	RegistryGlobal(wl_registry *Registry, name uint32, iface string, version uint32)
	RegistryGlobalRemove(wl_registry *Registry, name uint32)
}

var callback_wl_registry_mutex sync.Mutex
var callback_wl_registry_listeners = []RegistryListener{}

//export wlcallback_registry_global
func wlcallback_registry_global(data uintptr, registry *Registry, id uint32, iface *C.char, version uint32) {
	callback_wl_registry_mutex.Lock()
	var doner = callback_wl_registry_listeners[data]
	callback_wl_registry_mutex.Unlock()

	var str = C.GoString(iface)
	_ = str

	doner.RegistryGlobal((*Registry)(registry), id, str, version)
}

//export wlcallback_registry_global_remove
func wlcallback_registry_global_remove(data uintptr, registry *Registry, name uint32) {
	callback_wl_registry_mutex.Lock()
	var doner = callback_wl_registry_listeners[data]
	callback_wl_registry_mutex.Unlock()

	doner.RegistryGlobalRemove((*Registry)(registry), name)
}

func RegistryAddListener(wl_callback *Registry, globrem RegistryListener) {
	callback_wl_registry_mutex.Lock()
	var i uintptr = uintptr(len(callback_wl_registry_listeners))
	callback_wl_registry_listeners = append(callback_wl_registry_listeners, globrem)
	callback_wl_registry_mutex.Unlock()

	C._wl_registry_add_listener((*C.struct_wl_registry)(unsafe.Pointer(wl_callback)), unsafe.Pointer(i))
}

//

type ShmListener interface {
	ShmFormat(wl_shm *Shm, format uint32)
}

var callback_wl_shm_mutex sync.Mutex
var callback_wl_shm_listeners = []ShmListener{}

//export wlcallback_shm_format
func wlcallback_shm_format(data uintptr, cb *C.struct_wl_shm, format uint32) {
	callback_wl_shm_mutex.Lock()
	var formater = callback_wl_shm_listeners[data]
	callback_wl_shm_mutex.Unlock()

	formater.ShmFormat((*Shm)(cb), format)
}

func ShmAddListener(wl_callback *Shm, formater ShmListener) {
	callback_wl_shm_mutex.Lock()
	var i = uintptr(len(callback_wl_shm_listeners))
	callback_wl_shm_listeners = append(callback_wl_shm_listeners, formater)
	callback_wl_shm_mutex.Unlock()

	C._wl_shm_add_listener((*C.struct_wl_shm)(unsafe.Pointer(wl_callback)), unsafe.Pointer(i))
}

type SeatListener interface {
	SeatCapabilities(wl_seat *Seat, capabilities uint32)
	SeatName(wl_seat *Seat, name string)
}

var seat_mutex sync.Mutex
var seat_listeners = []SeatListener{}

//export wlcallback_handle_seat_capabilities
func wlcallback_handle_seat_capabilities(data uintptr, seat *C.struct_wl_seat, capabilities uint32) {
	seat_mutex.Lock()
	var doner = seat_listeners[data]
	seat_mutex.Unlock()

	doner.SeatCapabilities((*Seat)(seat), capabilities)
}

//export wlcallback_handle_seat_name
func wlcallback_handle_seat_name(data uintptr, seat *C.struct_wl_seat, name *C.char) {
	seat_mutex.Lock()
	var doner = seat_listeners[data]
	seat_mutex.Unlock()

	var str = C.GoString(name)
	_ = str

	doner.SeatName((*Seat)(seat), str)
}

func SeatAddListener(wl_seat *Seat, doner SeatListener) {
	seat_mutex.Lock()
	var i = uintptr(len(seat_listeners))
	seat_listeners = append(seat_listeners, doner)
	seat_mutex.Unlock()

	C._wl_seat_add_listener((*C.struct_wl_seat)(unsafe.Pointer(wl_seat)), unsafe.Pointer(i))
}

type OutputListener interface {
	OutputGeometry(wl_output *Output, x int, y int, physical_width int, physical_height int, subpixel int, make string, model string, transform int)
	OutputDone(wl_output *Output)
	OutputScale(wl_output *Output, scale int32)
	OutputMode(wl_output *Output, flags uint32, width int, height int, refresh int)
}

var output_listener_mutex sync.Mutex
var output_listener_listeners = []OutputListener{}

//export wlcallback_handle_output_geometry
func wlcallback_handle_output_geometry(data uintptr, ptr *C.struct_wl_output,
	x int, y int, physical_width int, physical_height int, subpixel int, make *C.char, model *C.char, transform int) {

	output_listener_mutex.Lock()
	var doner = output_listener_listeners[data]
	output_listener_mutex.Unlock()

	var strmake = C.GoString(make)

	var strmodel = C.GoString(model)

	doner.OutputGeometry((*Output)(ptr), x, y, physical_width, physical_height, subpixel, strmake, strmodel, transform)
}

//export wlcallback_handle_output_done
func wlcallback_handle_output_done(data uintptr, ptr *C.struct_wl_output) {

	output_listener_mutex.Lock()
	var doner = output_listener_listeners[data]
	output_listener_mutex.Unlock()

	doner.OutputDone((*Output)(ptr))
}

//export wlcallback_handle_output_scale
func wlcallback_handle_output_scale(data uintptr, ptr *C.struct_wl_output, scale int32) {

	output_listener_mutex.Lock()
	var doner = output_listener_listeners[data]
	output_listener_mutex.Unlock()

	doner.OutputScale((*Output)(ptr), scale)
}

//export wlcallback_handle_output_mode
func wlcallback_handle_output_mode(data uintptr, ptr *C.struct_wl_output, flags uint32, width int, height int, refresh int) {

	output_listener_mutex.Lock()
	var doner = output_listener_listeners[data]
	output_listener_mutex.Unlock()

	doner.OutputMode((*Output)(ptr), flags, width, height, refresh)
}

func OutputAddListener(wl_output *Output, doner OutputListener) {
	output_listener_mutex.Lock()
	var i = uintptr(len(output_listener_listeners))
	output_listener_listeners = append(output_listener_listeners, doner)
	output_listener_mutex.Unlock()

	C._wl_output_add_listener((*C.struct_wl_output)(unsafe.Pointer(wl_output)), unsafe.Pointer(i))
}

type PointerListener interface {
	PointerEnter(wl_pointer *Pointer, serial uint32, wl_surface *Surface, surface_x Fixed, surface_y Fixed)
	PointerLeave(wl_pointer *Pointer, serial uint32, wl_surface *Surface)
	PointerMotion(wl_pointer *Pointer, time uint32, surface_x Fixed, surface_y Fixed)
	PointerButton(wl_pointer *Pointer, serial uint32, time uint32, button uint32, state uint32)
	PointerAxis(wl_pointer *Pointer, time uint32, axis uint32, value Fixed)
	PointerFrame(wl_pointer *Pointer)
	PointerAxisSource(wl_pointer *Pointer, axis_source uint32)
	PointerAxisStop(wl_pointer *Pointer, time uint32, axis uint32)
	PointerAxisDiscrete(wl_pointer *Pointer, axis uint32, discrete int32)
}

var pointer_listener_mutex sync.Mutex
var pointer_listener_listeners = []PointerListener{}

//export wlcallback_handle_pointer_enter
func wlcallback_handle_pointer_enter(data uintptr, ptr *C.struct_wl_pointer,
	serial uint32, sfc *C.struct_wl_surface, surface_x Fixed, surface_y Fixed) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerEnter((*Pointer)(ptr), serial, (*Surface)(sfc), surface_x, surface_y)
}

//export wlcallback_handle_pointer_leave
func wlcallback_handle_pointer_leave(data uintptr, ptr *C.struct_wl_pointer,
	serial uint32, sfc *C.struct_wl_surface) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerLeave((*Pointer)(ptr), serial, (*Surface)(sfc))
}

//export wlcallback_handle_pointer_motion
func wlcallback_handle_pointer_motion(data uintptr, ptr *C.struct_wl_pointer,
	time uint32, surface_x Fixed, surface_y Fixed) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerMotion((*Pointer)(ptr), time, surface_x, surface_y)
}

//export wlcallback_handle_pointer_button
func wlcallback_handle_pointer_button(data uintptr, ptr *C.struct_wl_pointer,
	serial uint32, time uint32, button uint32, state uint32) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerButton((*Pointer)(ptr), serial, time, button, state)
}

//export wlcallback_handle_pointer_axis
func wlcallback_handle_pointer_axis(data uintptr, ptr *C.struct_wl_pointer,
	time uint32, axis uint32, value Fixed) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerAxis((*Pointer)(ptr), time, axis, value)
}

//export wlcallback_handle_pointer_frame
func wlcallback_handle_pointer_frame(data uintptr, ptr *C.struct_wl_pointer) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerFrame((*Pointer)(ptr))
}

//export wlcallback_handle_pointer_axis_source
func wlcallback_handle_pointer_axis_source(data uintptr, ptr *C.struct_wl_pointer,
	axis_source uint32) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerAxisSource((*Pointer)(ptr), axis_source)
}

//export wlcallback_handle_pointer_axis_stop
func wlcallback_handle_pointer_axis_stop(data uintptr, ptr *C.struct_wl_pointer,
	time uint32, axis uint32) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerAxisStop((*Pointer)(ptr), time, axis)
}

//export wlcallback_handle_pointer_axis_discrete
func wlcallback_handle_pointer_axis_discrete(data uintptr, ptr *C.struct_wl_pointer,
	axis uint32, discrete int32) {

	pointer_listener_mutex.Lock()
	var doner = pointer_listener_listeners[data]
	pointer_listener_mutex.Unlock()

	doner.PointerAxisDiscrete((*Pointer)(ptr), axis, discrete)
}

func PointerAddListener(wl_pointer *Pointer, doner PointerListener) uintptr {
	pointer_listener_mutex.Lock()
	var i = uintptr(len(pointer_listener_listeners))
	pointer_listener_listeners = append(pointer_listener_listeners, doner)
	pointer_listener_mutex.Unlock()

	C._wl_pointer_add_listener((*C.struct_wl_pointer)(unsafe.Pointer(wl_pointer)), unsafe.Pointer(i))

	return i
}

// Temporary callback ////////////////////////////////////////////////////////////

//surface add listener

type SurfaceEnterListener func(ptr uintptr, wl_surface *Surface, wl_output *Output)
type SurfaceLeaveListener func(ptr uintptr, wl_surface *Surface, wl_output *Output)

type surfacesetuserdata struct {
	e SurfaceEnterListener
	l SurfaceLeaveListener
	p uintptr
}

var surface_mutex sync.Mutex
var surface_listeners_map map[*Surface]surfacesetuserdata

//export wlcallback_handle_surface_enter
func wlcallback_handle_surface_enter(data uintptr, sf *C.struct_wl_surface, out *C.struct_wl_output) {

	surface_mutex.Lock()
	var listener = surface_listeners_map[(*Surface)(sf)]
	surface_mutex.Unlock()

	listener.e(listener.p, (*Surface)(sf), (*Output)(out))
}

//export wlcallback_handle_surface_leave
func wlcallback_handle_surface_leave(data uintptr, sf *C.struct_wl_surface, out *C.struct_wl_output) {

	surface_mutex.Lock()
	var listener = surface_listeners_map[(*Surface)(sf)]
	surface_mutex.Unlock()

	listener.l(listener.p, (*Surface)(sf), (*Output)(out))
}

func SurfaceAddListener(wl_surface *Surface, se SurfaceEnterListener, sl SurfaceLeaveListener, sp uintptr) {
	surface_mutex.Lock()

	if surface_listeners_map == nil {
		surface_listeners_map = make(map[*Surface]surfacesetuserdata)
	}

	var put = surfacesetuserdata{e: se, l: sl, p: sp}

	surface_listeners_map[wl_surface] = put

	surface_mutex.Unlock()

	C._wl_surface_add_listener((*C.struct_wl_surface)(unsafe.Pointer(wl_surface)), unsafe.Pointer(wl_surface))
}

//

type CallbackListener interface {
	CallbackDone(wl_callback *Callback, callback_data uint32)
}

var callback_frame_mutex sync.Mutex
var callback_frame_listeners = []CallbackListener{}

//export wlcallback_handle_frame
func wlcallback_handle_frame(data uintptr, cb *C.struct_wl_callback, time uint32) {
	callback_frame_mutex.Lock()
	var doner = callback_frame_listeners[data]
	callback_frame_listeners[data] = nil
	callback_frame_mutex.Unlock()

	doner.CallbackDone((*Callback)(cb), time)
}

func CallbackAddListener(wl_callback *Callback, doner CallbackListener) {
	callback_frame_mutex.Lock()

	var i uintptr
	for i = 0; i < uintptr(len(callback_frame_listeners)); i++ {
		if callback_frame_listeners[i] == nil {
			goto found
		}
	}
	callback_frame_listeners = append(callback_frame_listeners, nil)

found:

	callback_frame_listeners[i] = doner

	callback_frame_mutex.Unlock()

	C._wl_callback_add_listener((*C.struct_wl_callback)(unsafe.Pointer(wl_callback)), unsafe.Pointer(i))
}

// Specials //////////////////////////////////////////////////////////////////////////////

func (t *XdgToplevel) SetTitle(strtitle string) {

	var title = []byte(strtitle)

	// add nul byte
	if (len(title) == 0) || (title[len(title)-1] != 0) {
		title = append(title, 0)
	}

	C._xdg_toplevel_set_title(unsafe.Pointer(t), unsafe.Pointer(&title[0]), C.int(len(title)))
}

///////////////////////////////////////////////////////////////////////////////////////////

func SurfaceDamage(s *Surface, a, b, c, d int32) {
	C.wl_surface_damage((*C.struct_wl_surface)(unsafe.Pointer(s)), (C.int32_t)(a), (C.int32_t)(b),
		(C.int32_t)(c), (C.int32_t)(d))
}

func (c *Compositor) CreateSurface() (*Surface, error) {
	return (*Surface)(C.wl_compositor_create_surface((*C.struct_wl_compositor)(unsafe.Pointer(c)))), nil
}
func (c *Subcompositor) GetSubsurface(surface, parent *Surface) (*Subsurface, error) {
	return (*Subsurface)(C.wl_subcompositor_get_subsurface((*C.struct_wl_subcompositor)(unsafe.Pointer(c)),
		(*C.struct_wl_surface)(unsafe.Pointer(surface)),
		(*C.struct_wl_surface)(unsafe.Pointer(parent)))), nil
}

func (sub *Subsurface) SetDesync() {
	C.wl_subsurface_set_desync((*C.struct_wl_subsurface)(unsafe.Pointer(sub)))
}

func (sub *Subsurface) SetSync() {
	C.wl_subsurface_set_sync((*C.struct_wl_subsurface)(unsafe.Pointer(sub)))
}

func (sub *Subsurface) SetPosition(x, y int) {
	C.wl_subsurface_set_position((*C.struct_wl_subsurface)(unsafe.Pointer(sub)), (C.int32_t)(x), (C.int32_t)(y))
}

func DisplayDispatch(d *Display) int {
	return (int)(C.wl_display_dispatch((*C.struct_wl_display)(unsafe.Pointer(d))))
}

func DisplayRoundtrip(d *Display) int {
	return (int)(C.wl_display_roundtrip((*C.struct_wl_display)(unsafe.Pointer(d))))
}

func DisplayConnect(name []byte) (*Display, error) {
	if len(name) == 0 {
		return (*Display)(C.wl_display_connect((*C.char)(nil))), nil
	} else if 0 == name[len(name)-1] {
		return (*Display)(C.wl_display_connect((*C.char)(unsafe.Pointer(&name[0])))), nil
	}
	return nil, nil
}

func DisplayGetRegistry(d *Display) (*Registry, error) {
	return (*Registry)(C.wl_display_get_registry((*C.struct_wl_display)(unsafe.Pointer(d)))), nil
}

func (s *Surface) Commit() {
	C.wl_surface_commit((*C.struct_wl_surface)(unsafe.Pointer(s)))
}
func SurfaceFrame(s *Surface) *Callback {
	return (*Callback)(C.wl_surface_frame((*C.struct_wl_surface)(unsafe.Pointer(s))))
}

func CallbackDestroy(c *Callback) {
	C.wl_callback_destroy((*C.struct_wl_callback)(unsafe.Pointer(c)))
}

func SurfaceAttach(s *Surface, b *Buffer, x int32, y int32) {
	C.wl_surface_attach((*C.struct_wl_surface)(unsafe.Pointer(s)), (*C.struct_wl_buffer)(unsafe.Pointer(b)), C.int32_t(x), C.int32_t(y))
}

func (s *XdgSurface) GetToplevel() (t *XdgToplevel, err error) {
	return (*XdgToplevel)(C.xdg_surface_get_toplevel((*C.struct_xdg_surface)(unsafe.Pointer(s)))), nil
}

func (sh *XdgWmBase) GetSurface(wls *Surface) (s *XdgSurface, err error) {
	return (*XdgSurface)(C.xdg_wm_base_get_xdg_surface((*C.struct_xdg_wm_base)(unsafe.Pointer(sh)), (*C.struct_wl_surface)(unsafe.Pointer(wls)))), nil
}

func XdgWmBasePong(sh *XdgWmBase, serial uint32) {
	C.xdg_wm_base_pong((*C.struct_xdg_wm_base)(unsafe.Pointer(sh)), (C.uint32_t)(serial))
}

func XdgSurfaceAckConfigure(sf *XdgSurface, serial uint32) {
	C.xdg_surface_ack_configure((*C.struct_xdg_surface)(unsafe.Pointer(sf)), (C.uint32_t)(serial))
}

func RegistryBindCompositorInterface(r *Registry, name uint32, version uint32) *Compositor {
	return (*Compositor)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_compositor_interface, (C.uint32_t)(version)))
}
func RegistryBindSubcompositorInterface(r *Registry, name uint32, version uint32) *Subcompositor {
	return (*Subcompositor)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_subcompositor_interface, (C.uint32_t)(version)))
}

func RegistryBindShmInterface(r *Registry, name uint32, version uint32) *Shm {
	return (*Shm)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_shm_interface, (C.uint32_t)(version)))
}

//

func RegistryBindXdgWmBaseInterface(r *Registry, name uint32, version uint32) *XdgWmBase {
	return (*XdgWmBase)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.xdg_wm_base_interface, (C.uint32_t)(version)))
}

func RegistryBindFullscreenShellV1Interface(r *Registry, name uint32, version uint32) *FullscreenShellV1 {
	return (*FullscreenShellV1)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.zwp_fullscreen_shell_v1_interface, (C.uint32_t)(version)))
}

func RegistryBindSeatInterface(r *Registry, name uint32, version uint32) *Seat {
	return (*Seat)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_seat_interface, (C.uint32_t)(version)))
}

func ShmCreatePool(shm *Shm, fd int32, size int32) *ShmPool {
	return (*ShmPool)(C.wl_shm_create_pool((*C.struct_wl_shm)(unsafe.Pointer(shm)),
		C.int32_t(fd), C.int32_t(size)))
}

func ShmPoolDestroy(p *ShmPool) {
	C.wl_shm_pool_destroy((*C.struct_wl_shm_pool)(unsafe.Pointer(p)))
}

func ShmPoolCreateBuffer(p *ShmPool, offset int32, width int32, height int32, stride int32, format uint32) *Buffer {
	return (*Buffer)(C.wl_shm_pool_create_buffer((*C.struct_wl_shm_pool)(unsafe.Pointer(p)), C.int32_t(offset),
		C.int32_t(width), C.int32_t(height),
		C.int32_t(stride), C.uint32_t(format)))
}

const SHM_FORMAT_XRGB8888 = C.WL_SHM_FORMAT_XRGB8888
const SHM_FORMAT_RGB565 = C.WL_SHM_FORMAT_RGB565
const SHM_FORMAT_ARGB8888 = C.WL_SHM_FORMAT_ARGB8888

// Fixed is a fraction type
type Fixed int32

type PointerButtonState int

type RelativePointerManagerV1 C.struct_zwp_relative_pointer_manager_v1
type PointerConstraintsV1 C.struct_zwp_pointer_constraints_v1

func DisplayGetFd(d *Display) int {
	return (int)(C.wl_display_get_fd((*C.struct_wl_display)(unsafe.Pointer(d))))
}

func DisplayPrepareRead(d *Display) int {
	return (int)(C.wl_display_prepare_read((*C.struct_wl_display)(unsafe.Pointer(d))))
}

const ErrAgain = int(C.EAGAIN)

func DisplayFlush(d *Display) (n int, err error) {
	cn, errno := C.wl_display_flush((*C.struct_wl_display)(unsafe.Pointer(d)))
	n = int(cn)
	if n == -1 {
		err = errno
	}

	return
}
func DisplayCancelRead(d *Display) {
	C.wl_display_cancel_read((*C.struct_wl_display)(d))
}
func DisplayReadEvents(d *Display) {
	C.wl_display_read_events((*C.struct_wl_display)(d))
}
func RegistryBindDataDeviceManager(r *Registry, name uint32, version uint32) *DataDeviceManager {
	return (*DataDeviceManager)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_data_device_manager_interface, (C.uint32_t)(version)))
}

func RegistryBindOutputInterface(r *Registry, name uint32, version uint32) *Output {
	return (*Output)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_output_interface, (C.uint32_t)(version)))
}

const TOPLEVEL_STATE_MAXIMIZED = 1

/**
 * the surface is fullscreen
 */
const TOPLEVEL_STATE_FULLSCREEN = 2

/**
 * the surface is being resized
 */
const TOPLEVEL_STATE_RESIZING = 3

/**
 * the surface is now activated
 */
const TOPLEVEL_STATE_ACTIVATED = 4

const SEAT_CAPABILITY_POINTER = 1
const SEAT_CAPABILITY_KEYBOARD = 2
const SEAT_CAPABILITY_TOUCH = 4

func SeatGetPointer(seat *Seat) *Pointer {
	return (*Pointer)(C.wl_seat_get_pointer((*C.struct_wl_seat)(unsafe.Pointer(seat))))
}

func PointerRelease(p *Pointer) {
	// unused because of a strange bug

	C.wl_pointer_release((*C.struct_wl_pointer)(unsafe.Pointer(p)))
}

func PointerDestroy(p *Pointer) {
	C.wl_pointer_destroy((*C.struct_wl_pointer)(unsafe.Pointer(p)))
}

const POINTER_ENTER_SINCE_VERSION = 1
const POINTER_LEAVE_SINCE_VERSION = 1
const POINTER_MOTION_SINCE_VERSION = 1
const POINTER_BUTTON_SINCE_VERSION = 1
const POINTER_AXIS_SINCE_VERSION = 1
const POINTER_FRAME_SINCE_VERSION = 5
const POINTER_AXIS_SOURCE_SINCE_VERSION = 5
const POINTER_AXIS_STOP_SINCE_VERSION = 5
const POINTER_AXIS_DISCRETE_SINCE_VERSION = 5
const POINTER_SET_CURSOR_SINCE_VERSION = 1
const POINTER_RELEASE_SINCE_VERSION = 3

func FixedToDouble(num Fixed) float64 {
	return float64(num) * (1. / 256.)
}

func FixedToFloat(num Fixed) float32 {
	return float32(num) * (1. / 256.)
}

func SurfaceSetUserData(s *Surface, userdata uintptr) {
	C.wl_surface_set_user_data((*C.struct_wl_surface)(unsafe.Pointer(s)), unsafe.Pointer(userdata))
}

func PointerSetUserData(s *Pointer, userdata uintptr) {
	C.wl_pointer_set_user_data((*C.struct_wl_pointer)(unsafe.Pointer(s)), unsafe.Pointer(userdata))
}

func DisplayDispatchPending(d *Display) int {
	return (int)(C.wl_display_dispatch_pending((*C.struct_wl_display)(unsafe.Pointer(d))))
}

func XdgSurfaceSetWindowGeometry(surf *XdgSurface, x int32, y int32, width int32, height int32) {
	C.xdg_surface_set_window_geometry((*C.struct_xdg_surface)(unsafe.Pointer(surf)),
		C.int32_t(x), C.int32_t(y), C.int32_t(width), C.int32_t(height))
}

func CursorFrameAndDuration(cursor *Cursor, time uint32, duration *uint32) int {
	return int(C.wl_cursor_frame_and_duration(
		(*C.struct_wl_cursor)(unsafe.Pointer(cursor)),
		C.uint32_t(time),
		(*C.uint32_t)(duration)))
}

func CursorImageCount(cursor *Cursor) int {
	return int(cursor.image_count)
}

func CursorGetCursorImage(cursor *Cursor, index int) *CursorImage {
	var c = (*C.struct_wl_cursor)(unsafe.Pointer(cursor))

	if index >= (int)(c.image_count) {
		return nil
	}

	var i = uintptr(unsafe.Pointer(*c.images)) + unsafe.Sizeof(c)*uintptr(index)

	return (*CursorImage)(unsafe.Pointer(i))
}

func CursorImageGetWidth(img *CursorImage) int {
	var c = (*C.struct_wl_cursor_image)(unsafe.Pointer(img))
	return int(c.width)
}

func CursorImageGetHeight(img *CursorImage) int {
	var c = (*C.struct_wl_cursor_image)(unsafe.Pointer(img))
	return int(c.height)
}

func CursorImageGetHotspotX(img *CursorImage) int {
	var c = (*C.struct_wl_cursor_image)(unsafe.Pointer(img))
	return int(c.hotspot_x)
}

func CursorImageGetHotspotY(img *CursorImage) int {
	var c = (*C.struct_wl_cursor_image)(unsafe.Pointer(img))
	return int(c.hotspot_y)
}

func CursorThemeGetCursor(theme *CursorTheme, str []byte) *Cursor {
	return (*Cursor)(C.wl_cursor_theme_get_cursor((*C.struct_wl_cursor_theme)(theme), (*C.char)(unsafe.Pointer(&str[0]))))
}

func CursorThemeLoad(name []byte, bar int, shm *Shm) *CursorTheme {
	if name == nil {
		return (*CursorTheme)(C.wl_cursor_theme_load((*C.char)(unsafe.Pointer(uintptr(0))), C.int(bar), (*C.struct_wl_shm)(shm)))
	} else {
		return (*CursorTheme)(C.wl_cursor_theme_load((*C.char)(unsafe.Pointer(&name[0])), C.int(bar), (*C.struct_wl_shm)(shm)))
	}
}

func CursorThemeDestroy(theme *CursorTheme) {
	C.wl_cursor_theme_destroy((*C.struct_wl_cursor_theme)(theme))
}

type Cursor C.struct_wl_cursor
type CursorImage C.struct_wl_cursor_image
type CursorTheme C.struct_wl_cursor_theme

func XdgPopupDestroy(p *XdgPopup) {
	C.xdg_popup_destroy((*C.struct_xdg_popup)(p))
}
func XdgSurfaceDestroy(s *XdgSurface) {
	C.xdg_surface_destroy((*C.struct_xdg_surface)(s))
}
func XdgToplevelDestroy(t *XdgToplevel) {
	C.xdg_toplevel_destroy((*C.struct_xdg_toplevel)(t))
}
func XdgWmBaseDestroy(s *XdgWmBase) {
	C.xdg_wm_base_destroy((*C.struct_xdg_wm_base)(s))
}

/**
 * no transform
 */
const OUTPUT_TRANSFORM_NORMAL = 0

/**
 * 90 degrees counter-clockwise
 */
const OUTPUT_TRANSFORM_90 = 1

/**
 * 180 degrees counter-clockwise
 */
const OUTPUT_TRANSFORM_180 = 2

/**
 * 270 degrees counter-clockwise
 */
const OUTPUT_TRANSFORM_270 = 3

/**
 * 180 degree flip around a vertical axis
 */
const OUTPUT_TRANSFORM_FLIPPED = 4

/**
 * flip and rotate 90 degrees counter-clockwise
 */
const OUTPUT_TRANSFORM_FLIPPED_90 = 5

/**
 * flip and rotate 180 degrees counter-clockwise
 */
const OUTPUT_TRANSFORM_FLIPPED_180 = 6

/**
 * flip and rotate 270 degrees counter-clockwise
 */
const OUTPUT_TRANSFORM_FLIPPED_270 = 7

func BufferDestroy(buf *Buffer) {
	C.wl_buffer_destroy((*C.struct_wl_buffer)(buf))
}

func CursorImageGetBuffer(image *CursorImage) *Buffer {
	return (*Buffer)(unsafe.Pointer(C.wl_cursor_image_get_buffer((*C.struct_wl_cursor_image)(unsafe.Pointer(image)))))
}

func DisplayDisconnect(display *Display) {
	C.wl_display_disconnect((*C.struct_wl_display)(display))
}

func SurfaceSetOpaqueRegion(s *Surface, r *Region) {
	C.wl_surface_set_opaque_region((*C.struct_wl_surface)(unsafe.Pointer(s)),
		(*C.struct_wl_region)(unsafe.Pointer(r)))
}
func SurfaceSetInputRegion(s *Surface, r *Region) {
	C.wl_surface_set_input_region((*C.struct_wl_surface)(unsafe.Pointer(s)),
		(*C.struct_wl_region)(unsafe.Pointer(r)))
}
func SurfaceDestroy(s *Surface) {
	C.wl_surface_destroy((*C.struct_wl_surface)(unsafe.Pointer(s)))
}

func RegionDestroy(r *Region) {
	C.wl_region_destroy((*C.struct_wl_region)(unsafe.Pointer(r)))
}

func RegistryDestroy(r *Registry) {
	C.wl_registry_destroy((*C.struct_wl_registry)(unsafe.Pointer(r)))
}
func SubsurfaceDestroy(r *Subsurface) {
	C.wl_subsurface_destroy((*C.struct_wl_subsurface)(unsafe.Pointer(r)))
}

func ShmDestroy(shm *Shm) {
	C.wl_shm_destroy((*C.struct_wl_shm)(shm))
}

func PointerSetCursor(pointer *Pointer, s uint32, surface *Surface, x int32, y int32) {
	C.wl_pointer_set_cursor((*C.struct_wl_pointer)(unsafe.Pointer(pointer)), C.uint32_t(s),
		(*C.struct_wl_surface)(unsafe.Pointer(surface)),
		C.int32_t((x)), C.int32_t((y)))
}
