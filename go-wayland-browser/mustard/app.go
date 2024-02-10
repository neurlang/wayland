package mustard
import (
	"github.com/neurlang/wayland/window"
)
func CreateNewApp(name string) *App {
	display, err := window.DisplayCreate(nil)
	if err != nil {
		println(err.Error())
		return nil
	}
    return &App{display: display}
}

func (app *App) Run(callback func()) {
  window.DisplayRun(app.display)
}


func (app *App) AddWindow(w *Window) {
    // implementation here
	w.window = window.Create(app.display)
	w.window.SetTitle(w.title)
	w.window.SetBufferType(window.BufferTypeShm)
	w.rootFrame.widget = w.window.AddWidget(w.rootFrame)
	w.window.SetKeyboardHandler(w.rootFrame)
	w.rootFrame.widget.ScheduleResize(1024, 768)
}

func (app *App) DestroyWindow(window *Window) {
    // implementation here
}

func setWidgetWindow(widget Widget, window *Window) {
	widget.SetWindow(window)

	for _, childWidget := range widget.Widgets() {
		setWidgetWindow(childWidget, window)
	}
}

