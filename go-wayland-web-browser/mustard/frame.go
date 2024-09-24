package mustard

import 	window "github.com/neurlang/wayland/windowtrace"
import wl "github.com/neurlang/wayland/wl"
import "github.com/neurlang/wayland/external/swizzle"

// CreateFrame - Creates and returns a new Frame
func CreateFrame(orientation FrameOrientation) *Frame {
	var widgets []Widget

	return &Frame{
		baseWidget: baseWidget{
			widgetType: frameWidget,

			needsRepaint: true,
			widgets:      widgets,

			backgroundColor: "#fff",
		},

		orientation: orientation,
	}
}

// SetBackgroundColor - Sets the frame background color
func (frame *Frame) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		frame.backgroundColor = backgroundColor
		frame.needsRepaint = true
	}
}

// SetWidth - Sets the frame width
func (frame *Frame) SetWidth(width float64) {
	frame.box.width = width
	frame.fixedWidth = true
	frame.RequestReflow()
}

// SetHeight - Sets the frame height
func (frame *Frame) SetHeight(height float64) {
	frame.box.height = height
	frame.fixedHeight = true
	frame.RequestReflow()
}

// SetHeight - Sets the frame height
func (frame *Frame) GetHeight() float64 {
	return frame.box.height
}

// kb

func (frame *Frame) Key(
	win *window.Window,
	input *window.Input,
	time uint32,
	key uint32,
	notUnicode uint32,
	state wl.KeyboardKeyState,
	data window.WidgetHandler,
) {
	println(key)
	const ModControlMask = window.ModControlMask
	var window = frame.window
	var action = state
	_ = action

	var entered = input.GetRune(&notUnicode, key)
	const Repeat = 2
	const Release = 0
	const Press = 1

	switch key {
	case 14:
		if action == Repeat || action == Release {
			if window.activeInput != nil && len(window.activeInput.value) > 0 {
				if window.activeInput.cursorPosition == 0 {
					window.activeInput.value = window.activeInput.value[:len(window.activeInput.value)-1]
				} else {
					inputVal, cursorPos := window.activeInput.value, window.activeInput.cursorPosition

					if cursorPos+len(inputVal) > 0 {
						window.activeInput.value = inputVal[:len(inputVal)+cursorPos-1] + inputVal[len(inputVal)+cursorPos:]
					}
				}
				window.activeInput.needsRepaint = true
			}
		}
		return
	case 1:
		if action == Release {
			window.DestroyContextMenu()

			if window.activeInput != nil {
				window.activeInput.active = false
				window.activeInput.selected = false
				window.activeInput.needsRepaint = true
				window.activeInput = nil
			}
		}

		break
	case 103:
		if action == Release || action == Repeat {
			window.ProcessArrowKeys("up")
		}
		break
	case 108:
		if action == Release || action == Repeat {
			window.ProcessArrowKeys("down")
		}
		break
	case 105:
		if action == Release || action == Repeat {
			window.ProcessArrowKeys("left")
		}
		return
	case 106:
		if action == Release || action == Repeat {
			window.ProcessArrowKeys("right")
		}
		return
	case 47:
		if action == Release && input.GetModifiers()&ModControlMask != 0 {
			if window.activeInput != nil {
				if window.activeInput.cursorPosition == 0 {
					window.activeInput.value = window.activeInput.value + "GetClipboardString()"
				} else {
					inputVal, cursorPos := window.activeInput.value, window.activeInput.cursorPosition
					window.activeInput.value = inputVal[:len(inputVal)+cursorPos] + "GetClipboardString()" + inputVal[len(inputVal)+cursorPos:]
				}
				window.activeInput.needsRepaint = true
			}
		}
		break
	case 28:
		if action == Release {
			window.ProcessReturnKey()
		}
		return
	case 54:
		return
	case 42:
		return
	}

	if action == Release {
		return
	}
	if window.activeInput != nil {
		inputVal, cursorPos := window.activeInput.value, window.activeInput.cursorPosition

		window.activeInput.value = inputVal[:len(inputVal)+cursorPos] + string(entered) + inputVal[len(inputVal)+cursorPos:]
		window.activeInput.needsRepaint = true
	}

}
func (frame *Frame) Focus(window *window.Window, device *window.Input) {

}

//end kb

func (frame *Frame) PointerFrame(widget *window.Widget, input *window.Input) {
}
func (frame *Frame) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {
	frame.window.cursorX = float64(x)
	frame.window.cursorY = float64(y)

	for _, f := range frame.window.pointerPositionEventListeners {

		f(float64(x), float64(y))
	}

	frame.window.ProcessPointerPosition()

}
func (frame *Frame) Leave(widget *window.Widget, input *window.Input) {
	println("leave")
}
func (frame *Frame) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	frame.window.cursorX = float64(x)
	frame.window.cursorY = float64(y)

	for _, f := range frame.window.pointerPositionEventListeners {

		f(float64(x), float64(y))
	}

	frame.window.ProcessPointerPosition()

	return frame.window.cursor
}
func (frame *Frame) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {

	if state == 1 {
		frame.window.clickSerial = frame.window.window.Display.GetSerial()
		return
	}

	frame.window.ProcessPointerClick(int(button))
}
func (frame *Frame) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value float32) {
}
func (frame *Frame) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (frame *Frame) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
	println("axis stop", axis)
}
func (frame *Frame) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
	if axis == 0 {
		frame.window.ProcessScroll(0, -float64(discrete))
	} else {
		frame.window.ProcessScroll(-float64(discrete), 0)
	}
}

func (frame *Frame) render(s Surface, time uint32) {
	context := makeContextFromCairo(s)

	top, left, width, height := frame.computedBox.GetCoords()

	context.SetHexColor(frame.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	childrenLen := len(frame.widgets)
	if childrenLen > 0 {
		childrenWidgets := getCoreWidgets(frame.widgets)
		childrenLayout := calculateChildrenWidgetsLayout(childrenWidgets, top, left, width, height, frame.orientation)

		for idx, child := range frame.Widgets() {
			child.ComputedBox().SetCoords(childrenLayout[idx].box.GetCoords())
			child.render(s, time)
		}
	}
}

func (frame *Frame) Redraw(widget *window.Widget) {

	var time = (uint32)(widget.WidgetGetLastTime())

	var surface = frame.window.window.WindowGetSurface()

	if surface != nil {

		frame.render(surface, time)
		swizzle.BGRA(surface.ImageSurfaceGetData())
		surface.Destroy()
	}

	widget.ScheduleRedraw()

}
func (frame *Frame) Resize(widget *window.Widget, width int32, height int32, pwidth int32, pheight int32) {
	frame.RequestReflowWith(float64(width), float64(height))
}

func (frame *Frame) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (frame *Frame) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
	println(x, y)
}
func (frame *Frame) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {
	println(x, y)
}
func (frame *Frame) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (frame *Frame) TouchCancel(widget *window.Widget, width int32, height int32) {
}
