// Package wlclient implements a wayland-client like api
package wlclient

import "github.com/neurlang/wayland/wl"
import "github.com/neurlang/wayland/xdg"
import "github.com/neurlang/wayland/unstable"

func DisplayDispatch(d *wl.Display) error {
	return d.Context().Run()
}
func PointerSetUserData(p *wl.Pointer, data interface{}) {
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
}
func ShmDestroy(p *wl.Shm) {
}
func RegistryDestroy(p *wl.Registry) {
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
}

type SeatListener interface {
	wl.SeatCapabilitiesHandler
	wl.SeatNameHandler
}

func SeatAddListener(s *wl.Seat, data SeatListener) {
	s.AddCapabilitiesHandler(data)
	s.AddNameHandler(data)
}

type RegistryListener interface {
	wl.RegistryGlobalHandler
	wl.RegistryGlobalRemoveHandler
}

func RegistryAddListener(r *wl.Registry, data RegistryListener) {
	r.AddGlobalHandler(data)
	r.AddGlobalRemoveHandler(data)
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

func SurfaceAddListener(
	s *wl.Surface,
	enter func(*wl.Surface, *wl.Output),
	leave func(*wl.Surface, *wl.Output),
) {
	el := &SurfaceEnterLeave{surface: s, callbacks: [2]func(*wl.Surface, *wl.Output){enter, leave}}
	s.AddEnterHandler(el)
	s.AddLeaveHandler(el)
}

func ShmAddListener(p *wl.Shm, data wl.ShmFormatHandler) {
	p.AddFormatHandler(data)
}
func RegionDestroy(p *wl.Region) {
}
func CallbackDestroy(p *wl.Callback) {
}
func SubsurfaceDestroy(p *wl.Subsurface) {
}

func RegistryBindCompositorInterface(r *wl.Registry, name uint32, version uint32) *wl.Compositor {
	c := wl.NewCompositor(r.Ctx)
	_ = r.Bind(name, "wl_compositor", version, c)
	return c
}

func RegistryBindShmInterface(r *wl.Registry, name uint32, version uint32) *wl.Shm {
	s := wl.NewShm(r.Ctx)
	_ = r.Bind(name, "wl_shm", version, s)
	return s
}

func RegistryBindDataDeviceManagerInterface(
	r *wl.Registry,
	name uint32,
	version uint32,
) *wl.DataDeviceManager {
	d := wl.NewDataDeviceManager(r.Ctx)
	_ = r.Bind(name, "wl_data_device_manager", version, d)
	return d
}

func RegistryBindOutputInterface(r *wl.Registry, name uint32, version uint32) *wl.Output {
	d := wl.NewOutput(r.Ctx)
	_ = r.Bind(name, "wl_output", version, d)
	return d
}

func RegistryBindSeatInterface(r *wl.Registry, name uint32, version uint32) *wl.Seat {
	d := wl.NewSeat(r.Ctx)
	_ = r.Bind(name, "wl_seat", version, d)
	return d
}

func RegistryBindWmBaseInterface(r *wl.Registry, name uint32, version uint32) *xdg.WmBase {
	d := xdg.NewShell(r.Ctx)
	_ = r.Bind(name, "xdg_wm_base", version, d)
	return d
}

func RegistryBindUnstableInterface(
	r *wl.Registry,
	name uint32,
	iface string,
	version uint32,
) wl.Proxy {
	function := unstable.GetNewFunc(iface)
	if function == nil {
		return nil
	}
	d := function(r.Ctx)
	_ = r.Bind(name, iface, version, d)
	return d
}

func DisplayConnect(name []byte) (*wl.Display, error) {
	return wl.Connect(string(name))
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
	_ = d.Context().RunTill(cb)
	return err
}
func DisplayDisconnect(display *wl.Display) {
	display.Context().Close()
}
