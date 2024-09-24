package mustard

import (
	window "github.com/neurlang/wayland/windowtrace"
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
	w.window = window.Create(app.display)
	w.window.SetTitle(w.title)
	w.window.SetBufferType(window.BufferTypeShm)
	w.rootFrame.widget = w.window.AddWidget(w.rootFrame)
	w.window.SetKeyboardHandler(w.rootFrame)
	w.window.Display.SetSeatHandler(w)

	// hack, we reuse these variables for the initial window size
	w.rootFrame.widget.ScheduleResize(int32(w.cursorX), int32(w.cursorY))
	w.cursorX, w.cursorY = 0, 0
}

func (app *App) DestroyWindow(window *Window) {
	window.window.Destroy()
	window.window = nil
}

func setWidgetWindow(widget Widget, window *Window) {
	widget.SetWindow(window)

	for _, childWidget := range widget.Widgets() {
		setWidgetWindow(childWidget, window)
	}
}
