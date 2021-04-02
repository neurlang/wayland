package window

import cairo "github.com/neurlang/wayland/cairoshim"
import "github.com/neurlang/wayland/wl"

const FrameStatusNone = 0
const FrameStatusRepaint = 0x1
const FrameStatusMinimize = 0x2
const FrameStatusMaximize = 0x4
const FrameStatusClose = 0x8
const FrameStatusMenu = 0x10
const FrameStatusResize = 0x20
const FrameStatusMove = 0x40
const FrameStatusAll = FrameStatusRepaint | FrameStatusMinimize | FrameStatusMaximize |
	FrameStatusClose | FrameStatusMenu | FrameStatusResize |
	FrameStatusMove

const FrameFlagActive = 0x1
const FrameFlagMaximized = 0x2

const FrameButtonNone = 0
const FrameButtonClose = 0x1
const FrameButtonMaximize = 0x2
const FrameButtonMinimize = 0x4
const FrameButtonAll = FrameButtonClose | FrameButtonMaximize | FrameButtonMinimize

type windowFrame struct {
	frame  *frame
	widget *Widget
	child  *Widget
}

func (*windowFrame) Resize(Widget *Widget, width int32, height int32, pwidth int32, pheight int32) {}
func (*windowFrame) Redraw(Widget *Widget)                                                         {}
func (*windowFrame) Enter(Widget *Widget, Input *Input, x float32, y float32)                      {}
func (*windowFrame) Leave(Widget *Widget, Input *Input)                                            {}
func (*windowFrame) Motion(Widget *Widget, Input *Input, time uint32, x float32, y float32) int {
	return CursorWatch
}

func (*windowFrame) Button(
	Widget *Widget,
	Input *Input,
	time uint32,
	button uint32,
	state wl.PointerButtonState,
	data WidgetHandler,
) {
}
func (*windowFrame) TouchUp(Widget *Widget, Input *Input, serial uint32, time uint32, id int32) {}

func (*windowFrame) TouchDown(
	Widget *Widget,
	Input *Input,
	serial uint32,
	time uint32,
	id int32,
	x float32,
	y float32,
) {
}

func (*windowFrame) TouchMotion(
	Widget *Widget,
	Input *Input,
	time uint32,
	id int32,
	x float32,
	y float32,
) {
}
func (*windowFrame) TouchFrame(Widget *Widget, Input *Input)                                     {}
func (*windowFrame) TouchCancel(Widget *Widget, width int32, height int32)                       {}
func (*windowFrame) Axis(Widget *Widget, Input *Input, time uint32, axis uint32, value wl.Fixed) {}
func (*windowFrame) AxisSource(Widget *Widget, Input *Input, source uint32)                      {}
func (*windowFrame) AxisStop(Widget *Widget, Input *Input, time uint32, axis uint32)             {}
func (*windowFrame) AxisDiscrete(Widget *Widget, Input *Input, axis uint32, discrete int32)      {}
func (*windowFrame) PointerFrame(Widget *Widget, Input *Input)                                   {}

type frame struct {
	width, height int32
	title         string
	flags         uint32
	theme         *theme

	interior struct {
		x, y          int32
		width, height int32
	}
	shadowMargin  int
	opaqueMargin  int
	geometryDirty int

	status uint32
}

type theme struct {
}

func frameCreate(
	t *theme,
	width, height int32,
	buttons uint32,
	title string,
	icon cairo.Surface,
) *frame {
	return new(frame)
}
