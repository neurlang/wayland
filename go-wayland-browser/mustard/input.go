package mustard

import (
	assets "github.com/neurlang/wayland/go-wayland-browser/assets"

	"github.com/goki/freetype/truetype"
)

// CreateInputWidget - Creates and returns a new Input Widget
func CreateInputWidget() *InputWidget {
	var widgets []Widget
	font, _ := truetype.Parse(assets.OpenSans(400))

	return &InputWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			widgetType: inputWidget,

			cursor: IbeamCursor,

			backgroundColor: "#fff",

			font: font,
		},

		fontSize:  20,
		fontColor: "#000",
	}
}

// SetWidth - Sets the input width
func (input *InputWidget) SetWidth(width float64) {
	input.box.width = width
	input.fixedWidth = true
	input.RequestReflow()
}

// SetHeight - Sets the input height
func (input *InputWidget) SetHeight(height float64) {
	input.box.height = height
	input.fixedHeight = true
	input.RequestReflow()
}

// SetFontSize - Sets the input font size
func (input *InputWidget) SetFontSize(fontSize float64) {
	input.fontSize = fontSize
	input.needsRepaint = true
}

func (input *InputWidget) SetReturnCallback(returnCallback func()) {
	input.returnCallback = returnCallback
}

// SetFontColor - Sets the input font color
func (input *InputWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		input.fontColor = fontColor
		input.needsRepaint = true
	}
}

// SetFontColor - Sets the input font color
func (input *InputWidget) SetValue(value string) {
	input.value = value
	input.needsRepaint = true
}

// SetFontColor - Sets the input font color
func (input *InputWidget) GetValue() string {
	return input.value
}

func (input *InputWidget) GetCursorPos() int {
	return input.cursorPosition
}

// SetBackgroundColor - Sets the input background color
func (input *InputWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		input.backgroundColor = backgroundColor
		input.needsRepaint = true
	}
}

func (input *InputWidget) render(s Surface, time uint32) {
	context := makeContextFromCairo(s)

	input.padding = 4
	top, left, width, height := input.computedBox.GetCoords()
	totalPadding := input.padding * 2

	if input.selected {
		context.SetHexColor("#e4e4e4")
		context.SetHexColor("#e4e4e4")
		//context.Clear()
	} else {
		context.SetHexColor("#efefef")
		context.SetHexColor("#efefef")
		//context.Clear()
	}

	if input.active {
		context.SetHexColor("#fff")
		context.SetHexColor("#fff")
		//context.Clear()
	} else {
		input.cursorPosition = 0
	}

	tx, ty := float64(left+totalPadding/2), float64(top+totalPadding/2)

	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	context.SetHexColor("#2f2f2f")

	context.SetFont(input.font, input.fontSize)
	w, h := context.MeasureString(input.value)

	cursorP := width - totalPadding*2
	cP, _ := context.MeasureString(input.value[len(input.value)+input.cursorPosition:])
	cursorP = cursorP - cP

	if cursorP > 0 {
		input.cursorFloat = true
	} else {
		input.cursorFloat = false
	}

	valueBigggerThanInput := w > float64(width)-input.fontSize
	if valueBigggerThanInput && input.active {
		if cursorP > 0 {
			context.DrawStringAnchored(input.value, tx+float64(width)-input.fontSize, ty+float64(height+totalPadding/2)/2, 1, 0)
		} else {
			context.DrawStringAnchored(input.value, tx+cP, ty+float64(height+totalPadding/2)/2, 1, 0)
		}
	} else {
		context.DrawString(input.value, tx, ty+float64(height+totalPadding/2)/2)
	}

	context.Fill()

	if input.active {
		context.SetHexColor("#000")

		if valueBigggerThanInput {
			if cursorP > 0 {
				context.DrawRectangle(tx+cursorP, ty+h/4, 1.3, float64(input.fontSize))
			} else {
				context.DrawRectangle(tx, ty+h/4, 1.3, float64(input.fontSize))
			}

		} else {
			cursorDefaultPosition, _ := context.MeasureString(input.value[:len(input.value)+input.cursorPosition])
			context.DrawRectangle(tx+cursorDefaultPosition, ty+h/4, 1.3, float64(input.fontSize))
		}

		context.Fill()
	}

	context.SetHexColor("#000")
	context.SetLineWidth(.4)

	context.DrawRectangle(
		left+1,
		top+1,
		width-2,
		height-2,
	)

	context.SetLineJoinRound()
	context.Stroke()
}
