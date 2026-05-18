//go:build js

package window

import (
	"syscall/js"
)

func startRenderWorker() {
	var frameFunc js.Func
	frameFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(windows) > 0 && len(windows[0].widgets) > 0 && windows[0].widgets[0].userdata != nil {
			if handler, ok := windows[0].widgets[0].userdata.(interface{ Redraw(*Widget) }); ok {
				handler.Redraw(windows[0].widgets[0])
			}
		}
		renderToCanvas()
		js.Global().Call("requestAnimationFrame", frameFunc)
		return nil
	})
	js.Global().Call("requestAnimationFrame", frameFunc)
}
