package main

import (
	"fmt"
	"time"

	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
)

type waylandMonitor struct {
	x, y, w, h, refresh int
	name                string
}

type waylandHandler struct {
	display        *wl.Display
	registry       *wl.Registry
	output         *wl.Output
	monitors       []waylandMonitor
	currentMonitor waylandMonitor
}

func (w *waylandHandler) HandleOutputGeometry(event wl.OutputGeometryEvent) {
	w.currentMonitor.x = int(event.X)
	w.currentMonitor.y = int(event.Y)
	w.currentMonitor.name = event.Model
}

func (w *waylandHandler) HandleOutputMode(event wl.OutputModeEvent) {
	w.currentMonitor.w = int(event.Width)
	w.currentMonitor.h = int(event.Height)
	w.currentMonitor.refresh = int(event.Refresh)
}

func (w *waylandHandler) HandleOutputDone(event wl.OutputDoneEvent) {
	w.monitors = append(w.monitors, w.currentMonitor)
	w.currentMonitor = waylandMonitor{}
}

func (w *waylandHandler) HandleOutputScale(event wl.OutputScaleEvent) {
	// This does not handle fractional scaling, presumably need xdg_output.logical_size for that
}

func (w *waylandHandler) HandleRegistryGlobal(event wl.RegistryGlobalEvent) {
	if event.Interface == "wl_output" {
		w.output = wlclient.RegistryBindOutputInterface(w.registry, event.Name, event.Version)
		wlclient.OutputAddListener(w.output, w)
	}
}

func (w *waylandHandler) HandleRegistryGlobalRemove(event wl.RegistryGlobalRemoveEvent) {
}

func enumMonitors() ([]waylandMonitor, error) {
	disp, err := wlclient.DisplayConnect(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer wlclient.DisplayDisconnect(disp)

	registry, err := disp.GetRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %v", err)
	}
	defer wlclient.RegistryDestroy(registry)

	handler := waylandHandler{display: disp, registry: registry}

	wlclient.RegistryAddListener(registry, &handler)

	for i := range 2 {
		err = wlclient.DisplayRoundtrip(disp)
		if err != nil {
			return nil, fmt.Errorf("failed to roundtrip %d: %v", i, err)
		}
	}

	return handler.monitors, nil
}

func main() {
	for range 2 {
		go func() {
			monitors, err := enumMonitors()
			if err != nil {
				panic(err)
			}

			fmt.Printf("GOROUTINE Monitors: %v\n", monitors)
		}()
	}
	time.Sleep(time.Second * 2)
}
