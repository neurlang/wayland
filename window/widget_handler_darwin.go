package window

import "github.com/neurlang/wayland/wl"

// WidgetHandler defines the interface for widget event handlers
type WidgetHandler interface {
	Resize(Widget *Widget, width int32, height int32, pwidth int32, pheight int32)
	Redraw(Widget *Widget)
	Enter(Widget *Widget, Input *Input, x float32, y float32)
	Leave(Widget *Widget, Input *Input)
	Motion(Widget *Widget, Input *Input, time uint32, x float32, y float32) int
	Button(
		Widget *Widget,
		Input *Input,
		time uint32,
		button uint32,
		state wl.PointerButtonState,
		data WidgetHandler,
	)
	TouchUp(Widget *Widget, Input *Input, serial uint32, time uint32, id int32)
	TouchDown(
		Widget *Widget,
		Input *Input,
		serial uint32,
		time uint32,
		id int32,
		x float32,
		y float32,
	)
	TouchMotion(Widget *Widget, Input *Input, time uint32, id int32, x float32, y float32)
	TouchFrame(Widget *Widget, Input *Input)
	TouchCancel(Widget *Widget, width int32, height int32)
	Axis(Widget *Widget, Input *Input, time uint32, axis uint32, value float32)
	AxisSource(Widget *Widget, Input *Input, source uint32)
	AxisStop(Widget *Widget, Input *Input, time uint32, axis uint32)
	AxisDiscrete(Widget *Widget, Input *Input, axis uint32, discrete int32)
	PointerFrame(Widget *Widget, Input *Input)
}

// KeyboardHandler defines the interface for keyboard event handlers
type KeyboardHandler interface {
	Key(Window *Window, Input *Input, time uint32, vKey uint32, code uint32, state wl.KeyboardKeyState, data WidgetHandler)
	Focus(Window *Window, Input *Input)
}

// FullscreenHandler defines the interface for fullscreen events
type FullscreenHandler interface {
	Fullscreen(Window *Window, data WidgetHandler)
}

// CloseHandler defines the interface for window close events
type CloseHandler interface {
	Close(Window *Window)
}

// DataHandler defines the interface for data transfer events
type DataHandler interface {
	DataOffer(Window *Window, offer interface{})
}

// ResizeHandler defines the interface for resize events
type ResizeHandler interface {
	Resize(Window *Window, width int32, height int32)
}
