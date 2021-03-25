package window

import cairo "github.com/neurlang/wayland/cairoshim"
import wl "github.com/neurlang/wayland/wl"

const FRAME_STATUS_NONE = 0
const FRAME_STATUS_REPAINT = 0x1
const FRAME_STATUS_MINIMIZE = 0x2
const FRAME_STATUS_MAXIMIZE = 0x4
const FRAME_STATUS_CLOSE = 0x8
const FRAME_STATUS_MENU = 0x10
const FRAME_STATUS_RESIZE = 0x20
const FRAME_STATUS_MOVE = 0x40
const FRAME_STATUS_ALL = FRAME_STATUS_REPAINT | FRAME_STATUS_MINIMIZE | FRAME_STATUS_MAXIMIZE |
	FRAME_STATUS_CLOSE | FRAME_STATUS_MENU | FRAME_STATUS_RESIZE |
	FRAME_STATUS_MOVE

const FRAME_FLAG_ACTIVE = 0x1
const FRAME_FLAG_MAXIMIZED = 0x2

const FRAME_BUTTON_NONE = 0
const FRAME_BUTTON_CLOSE = 0x1
const FRAME_BUTTON_MAXIMIZE = 0x2
const FRAME_BUTTON_MINIMIZE = 0x4
const FRAME_BUTTON_ALL = FRAME_BUTTON_CLOSE | FRAME_BUTTON_MAXIMIZE | FRAME_BUTTON_MINIMIZE

type window_frame struct {
	frame  *frame
	widget *Widget
	child  *Widget
}

func (*window_frame) Resize(Widget *Widget, width int32, height int32)         {}
func (*window_frame) Redraw(Widget *Widget)                                    {}
func (*window_frame) Enter(Widget *Widget, Input *Input, x float32, y float32) {}
func (*window_frame) Leave(Widget *Widget, Input *Input)                       {}
func (*window_frame) Motion(Widget *Widget, Input *Input, time uint32, x float32, y float32) int {
	return CURSOR_WATCH
}
func (*window_frame) Button(Widget *Widget, Input *Input, time uint32, button uint32, state wl.PointerButtonState, data WidgetHandler) {
}
func (*window_frame) TouchUp(Widget *Widget, Input *Input, serial uint32, time uint32, id int32) {}
func (*window_frame) TouchDown(Widget *Widget, Input *Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (*window_frame) TouchMotion(Widget *Widget, Input *Input, time uint32, id int32, x float32, y float32) {
}
func (*window_frame) TouchFrame(Widget *Widget, Input *Input)                                     {}
func (*window_frame) TouchCancel(Widget *Widget, width int32, height int32)                       {}
func (*window_frame) Axis(Widget *Widget, Input *Input, time uint32, axis uint32, value wl.Fixed) {}
func (*window_frame) AxisSource(Widget *Widget, Input *Input, source uint32)                      {}
func (*window_frame) AxisStop(Widget *Widget, Input *Input, time uint32, axis uint32)             {}
func (*window_frame) AxisDiscrete(Widget *Widget, Input *Input, axis uint32, discrete int32)      {}
func (*window_frame) PointerFrame(Widget *Widget, Input *Input)                                   {}

type frame struct {
	width, height int32
	title         string
	flags         uint32
	theme         *theme

	interior struct {
		x, y          int32
		width, height int32
	}
	shadow_margin  int
	opaque_margin  int
	geometry_dirty int

	status uint32
}

type theme struct {
}

func frame_create(t *theme, width, height int32, buttons uint32, title string, icon cairo.Surface) *frame {
	return new(frame)
}
