package window

import "github.com/neurlang/winc"
import "github.com/neurlang/wayland/wl"
import "github.com/neurlang/wayland/xdg"

type Display struct {
	count int
	//redraw chan func()
}

func (d *Display) Destroy() {

}

func (d *Display) CreateDataSource() (*DataSource, error) {
	return &DataSource{}, nil
}

func (d *Display) GetSerial() uint32 {
	return 0
}

func (d *Display) Exit() {
	winc.Exit()
}

func DisplayCreate(args []string) (*Display, error) {
	return &Display{
	}, nil

}

func DisplayRun(d *Display) {
	winc.RunMainLoop()
}

func (d *Display) SetSeatHandler(_ interface{}) {
}



// HandleRegistryGlobal is a dummy method for Display.
func (d *Display) HandleRegistryGlobal(_ wl.RegistryGlobalEvent) {
    // Dummy implementation
    // Add functionality here if needed
}

// HandleRegistryGlobalRemove is a dummy method for Display.
func (d *Display) HandleRegistryGlobalRemove(_ wl.RegistryGlobalRemoveEvent) {
    // Dummy implementation
}

// HandleShmFormat is a dummy method for Display.
func (d *Display) HandleShmFormat(_ wl.ShmFormatEvent) {
    // Dummy implementation
}

// HandleWmBasePing is a dummy method for Display.
func (d *Display) HandleWmBasePing(_ xdg.WmBasePingEvent) {
    // Dummy implementation
}

// RegistryGlobal is a dummy method for Display.
func (d *Display) RegistryGlobal(_ *wl.Registry, _ uint32, _ string, _ uint32) {
    // Dummy implementation
}

// RegistryGlobalRemove is a dummy method for Display.
func (d *Display) RegistryGlobalRemove(_ *wl.Registry, _ uint32) {
    // Dummy implementation
}


type GlobalHandler interface {
	HandleGlobal(d *Display, id uint32, iface string, version uint32, data interface{})
}

// SetGlobalHandler is a dummy method for Display.
func (d *Display) SetGlobalHandler(_ GlobalHandler) {
    // Dummy implementation
}

// SetUserData is a dummy method for Display.
func (d *Display) SetUserData(_ interface{}) {
    // Dummy implementation
}

// ShellPing is a dummy method for Display.
func (d *Display) ShellPing(*xdg.WmBase, uint32) {
    // Dummy implementation
}

// ShmFormat is a dummy method for Display.
func (d *Display) ShmFormat(*wl.Shm, uint32) {
    // Dummy implementation
}
