package window

import "github.com/tadvi/winc"

type Display struct {
	mustResize []*mustResize
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
	return &Display{}, nil

}

func DisplayRun(d *Display) {

	// hit scheduled resize
	for _, resize := range d.mustResize {

		resize.w.SetAllocation(0, 0, resize.width, resize.height)

		resize.w.handler.Resize(resize.w, resize.width, resize.height, resize.width, resize.height)
	}
	d.mustResize = nil

	winc.RunMainLoop()
}
