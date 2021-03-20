package wlclient

import wl "github.com/neurlang/wayland/wl"
import xdg "github.com/neurlang/wayland/xdg"
import "syscall"

func DisplayDispatch(d *wl.Display) error {
	panic("not implemented")
}
func PointerSetUserData(p *wl.Pointer, data interface{}) {
	return
}

func SurfaceSetUserData(p *wl.Surface, data interface{}) {
	p.UserData = data
}

type PointerListener interface {
	wl.PointerEnterHandler
	wl.PointerLeaveHandler
	wl.PointerMotionHandler
	wl.PointerButtonHandler
	wl.PointerAxisHandler
	wl.PointerFrameHandler
	wl.PointerAxisSourceHandler
	wl.PointerAxisStopHandler
	wl.PointerAxisDiscreteHandler
}

func PointerAddListener(p *wl.Pointer, h PointerListener) {
	p.AddEnterHandler(h)
	p.AddLeaveHandler(h)
	p.AddMotionHandler(h)
	p.AddButtonHandler(h)
	p.AddAxisHandler(h)
	p.AddFrameHandler(h)
	p.AddAxisSourceHandler(h)
	p.AddAxisStopHandler(h)
	p.AddAxisDiscreteHandler(h)

}
func PointerDestroy(p *wl.Pointer) {
	panic("not implemented")
}

func BufferDestroy(p *wl.Buffer) {
	panic("not implemented")
}
func ShmDestroy(p *wl.Shm) {
	panic("not implemented")
}
func ShmPoolDestroy(p *wl.ShmPool) {
	panic("not implemented")
}
func RegistryDestroy(p *wl.Registry) {
	panic("not implemented")
}
func BufferAddListener(b *wl.Buffer, data wl.BufferReleaseHandler) {
	b.AddReleaseHandler(data)
}
func CallbackAddListener(c *wl.Callback, data wl.CallbackDoneHandler) {
	c.AddDoneHandler(data)
}

type OutputListener interface {
	wl.OutputGeometryHandler
	wl.OutputModeHandler
	wl.OutputDoneHandler
	wl.OutputScaleHandler
}

func OutputAddListener(o *wl.Output, h OutputListener) {
	o.AddGeometryHandler(h)
	o.AddModeHandler(h)
	o.AddDoneHandler(h)
	o.AddScaleHandler(h)
	return
}

type SeatListener interface {
	wl.SeatCapabilitiesHandler
	wl.SeatNameHandler
}

func SeatAddListener(s *wl.Seat, data SeatListener) {
	s.AddCapabilitiesHandler(data)
	s.AddNameHandler(data)

	return
}

type RegistryListener interface {
	wl.RegistryGlobalHandler
	wl.RegistryGlobalRemoveHandler
}

func RegistryAddListener(r *wl.Registry, data RegistryListener) {
	r.AddGlobalHandler(data)
	r.AddGlobalRemoveHandler(data)
	return
}

type SurfaceEnterLeaveListener interface {
	wl.SurfaceEnterHandler
	wl.SurfaceLeaveHandler
}

type SurfaceEnterLeave struct {
	surface   *wl.Surface
	callbacks [2]func(*wl.Surface, *wl.Output)
}

func (el *SurfaceEnterLeave) HandleSurfaceEnter(en wl.SurfaceEnterEvent) {
	el.callbacks[0](el.surface, en.Output)
}

func (el *SurfaceEnterLeave) HandleSurfaceLeave(le wl.SurfaceLeaveEvent) {
	el.callbacks[1](el.surface, le.Output)
}

func SurfaceAddListener(s *wl.Surface, enter func(*wl.Surface, *wl.Output), leave func(*wl.Surface, *wl.Output)) {
	el := &SurfaceEnterLeave{surface: s, callbacks: [2]func(*wl.Surface, *wl.Output){enter, leave}}
	s.AddEnterHandler(el)
	s.AddLeaveHandler(el)
	return
}

func ShmAddListener(p *wl.Shm, data wl.ShmFormatHandler) {

	p.AddFormatHandler(data)

	return

}
func RegionDestroy(p *wl.Region) {
	return
}
func CallbackDestroy(p *wl.Callback) {
	return
}
func SubsurfaceDestroy(p *wl.Subsurface) {
	return
}

func RegistryBindCompositorInterface(r *wl.Registry, name uint32, version uint32) *wl.Compositor {
	c := wl.NewCompositor(r.Ctx)
	r.Bind(name, "wl_compositor", version, c)
	return c
}

func RegistryBindShmInterface(r *wl.Registry, name uint32, version uint32) *wl.Shm {
	s := wl.NewShm(r.Ctx)
	r.Bind(name, "wl_shm", version, s)
	return s
	//return (*Shm)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_shm_interface, (C.uint32_t)(version)))
}

func RegistryBindDataDeviceManagerInterface(r *wl.Registry, name uint32, version uint32) *wl.DataDeviceManager {
	d := wl.NewDataDeviceManager(r.Ctx)
	r.Bind(name, "wl_data_device_manager", version, d)
	return d
	//return (*DataDeviceManager)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_data_device_manager_interface, (C.uint32_t)(version)))
}

func RegistryBindOutputInterface(r *wl.Registry, name uint32, version uint32) *wl.Output {
	d := wl.NewOutput(r.Ctx)
	r.Bind(name, "wl_output", version, d)
	return d
	//return (*Output)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_output_interface, (C.uint32_t)(version)))
}

func RegistryBindSeatInterface(r *wl.Registry, name uint32, version uint32) *wl.Seat {
	d := wl.NewSeat(r.Ctx)
	r.Bind(name, "wl_seat", version, d)
	return d
	//return (*Seat)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_seat_interface, (C.uint32_t)(version)))
}

func RegistryBindShellInterface(r *wl.Registry, name uint32, version uint32) *xdg.Shell {
	d := xdg.NewShell(r.Ctx)
	r.Bind(name, "zxdg_shell_v6", version, d)
	return d
	//return (*Seat)(C.wl_registry_bind((*C.struct_wl_registry)(unsafe.Pointer(r)), (C.uint32_t)(name), &C.wl_seat_interface, (C.uint32_t)(version)))
}

func DisplayConnect(name []byte) (*wl.Display, error) {
	return wl.Connect(string(name))
}
func DisplayGetFd(d *wl.Display) int {
	return d.Fd
}
func DisplayGetRegistry(d *wl.Display) (*wl.Registry, error) {
	return d.GetRegistry()
}
func DisplayRun(d *wl.Display) error {
	return d.Context().Run()
}
func DisplayRoundtrip(d *wl.Display) error {
	cb, err := d.Sync()
	if err != nil {
		return err
	}
	d.Context().RunTill(cb)
	return err
}
func DisplayDisconnect(display *wl.Display) {
	panic("not implemented")
	//C.wl_display_disconnect((*C.struct_wl_display)(display))
}
func DisplayDispatchPending(d *wl.Display) int {
	//println("DisplayDispatchPending: no-op")
	return 0
}

var EAgain = syscall.EAGAIN

func DisplayFlush(d *wl.Display) error {
	//println("DisplayFlush: no-op")
	return nil
}
