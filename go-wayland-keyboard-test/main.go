//go:build darwin
// +build darwin

package main

import (
	"fmt"

	"github.com/neurlang/wayland/window"
	"github.com/neurlang/wayland/wl"
	xkbcommon "github.com/neurlang/wayland/xkbcommon"
)

type KeyboardTestApp struct {
	window      *window.Window
	widget      *window.Widget
	lastKey     string
	lastKeyCode uint32
	modifiers   string
}

func (app *KeyboardTestApp) Resize(widget *window.Widget, width int32, height int32, pwidth int32, pheight int32) {
	widget.ScheduleRedraw()
}

func (app *KeyboardTestApp) Redraw(widget *window.Widget) {
	surface := widget.ImageSurfaceGetData()
	width := widget.ImageSurfaceGetWidth()
	height := widget.ImageSurfaceGetHeight()
	stride := widget.ImageSurfaceGetStride()

	// Clear to dark gray
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := y*stride + x*4
			surface[offset+0] = 40  // B
			surface[offset+1] = 40  // G
			surface[offset+2] = 40  // R
			surface[offset+3] = 255 // A
		}
	}

}

func (app *KeyboardTestApp) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {}
func (app *KeyboardTestApp) Leave(widget *window.Widget, input *window.Input)                       {}
func (app *KeyboardTestApp) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {
	return 0
}
func (app *KeyboardTestApp) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {
}
func (app *KeyboardTestApp) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (app *KeyboardTestApp) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (app *KeyboardTestApp) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {
}
func (app *KeyboardTestApp) TouchFrame(widget *window.Widget, input *window.Input)        {}
func (app *KeyboardTestApp) TouchCancel(widget *window.Widget, width int32, height int32) {}
func (app *KeyboardTestApp) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value float32) {
}
func (app *KeyboardTestApp) AxisSource(widget *window.Widget, input *window.Input, source uint32) {}
func (app *KeyboardTestApp) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (app *KeyboardTestApp) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
}
func (app *KeyboardTestApp) PointerFrame(widget *window.Widget, input *window.Input) {}

// Key handles keyboard events
func (app *KeyboardTestApp) Key(win *window.Window, input *window.Input, time uint32, vKey uint32, code uint32, state wl.KeyboardKeyState, data window.WidgetHandler) {
	if state != wl.KeyboardKeyStatePressed {
		return
	}

	// Get modifiers
	mods := input.GetModifiers()
	modStr := ""
	if mods&window.ModShiftMask != 0 {
		modStr += "Shift "
	}
	if mods&window.ModControlMask != 0 {
		modStr += "Control "
	}
	if mods&window.ModAltMask != 0 {
		modStr += "Alt "
	}
	app.modifiers = modStr

	// Map key code to name
	keyName := getKeyName(vKey)
	if code >= 32 && code < 127 {
		keyName = fmt.Sprintf("%s ('%c')", keyName, rune(code))
	}

	app.lastKey = keyName
	app.lastKeyCode = vKey

	fmt.Printf("Key pressed: %s (vKey=%d, unicode=%d) modifiers=%s\n", keyName, vKey, code, modStr)

	// Check for Cmd+Q to quit
	if vKey == xkbcommon.KeyQ && mods&window.ModControlMask != 0 {
		fmt.Println("Cmd+Q pressed, quitting...")
		app.window.Destroy()
	}

	app.widget.ScheduleRedraw()
}

// Focus handles keyboard focus events
func (app *KeyboardTestApp) Focus(win *window.Window, input *window.Input) {
	fmt.Println("Window gained keyboard focus")
}

func getKeyName(keyCode uint32) string {
	switch keyCode {
	case xkbcommon.KEYa:
		return "A"
	case xkbcommon.KEYb:
		return "B"
	case xkbcommon.KEYc:
		return "C"
	case xkbcommon.KEYd:
		return "D"
	case xkbcommon.KEYe:
		return "E"
	case xkbcommon.KEYf:
		return "F"
	case xkbcommon.KEYg:
		return "G"
	case xkbcommon.KEYh:
		return "H"
	case xkbcommon.KEYi:
		return "I"
	case xkbcommon.KEYj:
		return "J"
	case xkbcommon.KEYk:
		return "K"
	case xkbcommon.KEYl:
		return "L"
	case xkbcommon.KEYm:
		return "M"
	case xkbcommon.KEYn:
		return "N"
	case xkbcommon.KEYo:
		return "O"
	case xkbcommon.KEYp:
		return "P"
	case xkbcommon.KEYq:
		return "Q"
	case xkbcommon.KEYr:
		return "R"
	case xkbcommon.KEYs:
		return "S"
	case xkbcommon.KEYt:
		return "T"
	case xkbcommon.KEYu:
		return "U"
	case xkbcommon.KEYv:
		return "V"
	case xkbcommon.KEYw:
		return "W"
	case xkbcommon.KEYx:
		return "X"
	case xkbcommon.KEYy:
		return "Y"
	case xkbcommon.KEYz:
		return "Z"
	case xkbcommon.Key0:
		return "0"
	case xkbcommon.Key1:
		return "1"
	case xkbcommon.Key2:
		return "2"
	case xkbcommon.Key3:
		return "3"
	case xkbcommon.Key4:
		return "4"
	case xkbcommon.Key5:
		return "5"
	case xkbcommon.Key6:
		return "6"
	case xkbcommon.Key7:
		return "7"
	case xkbcommon.Key8:
		return "8"
	case xkbcommon.Key9:
		return "9"
	case xkbcommon.KeySpace:
		return "Space"
	case xkbcommon.KeyReturn:
		return "Return"
	case xkbcommon.KeyBackspace:
		return "Backspace"
	case xkbcommon.KeyDelete:
		return "Delete"
	case xkbcommon.KeyTab:
		return "Tab"
	case xkbcommon.KeyEscape:
		return "Escape"
	case xkbcommon.KeyLeft:
		return "Left Arrow"
	case xkbcommon.KeyRight:
		return "Right Arrow"
	case xkbcommon.KeyUp:
		return "Up Arrow"
	case xkbcommon.KeyDown:
		return "Down Arrow"
	case xkbcommon.KeyHome:
		return "Home"
	case xkbcommon.KeyEnd:
		return "End"
	case xkbcommon.KeyPageUp:
		return "Page Up"
	case xkbcommon.KeyPageDown:
		return "Page Down"
	case xkbcommon.KeyF1:
		return "F1"
	case xkbcommon.KeyF2:
		return "F2"
	case xkbcommon.KeyF3:
		return "F3"
	case xkbcommon.KeyF4:
		return "F4"
	case xkbcommon.KeyF5:
		return "F5"
	case xkbcommon.KeyF6:
		return "F6"
	case xkbcommon.KeyF7:
		return "F7"
	case xkbcommon.KeyF8:
		return "F8"
	case xkbcommon.KeyF9:
		return "F9"
	case xkbcommon.KeyF10:
		return "F10"
	case xkbcommon.KeyF11:
		return "F11"
	case xkbcommon.KeyF12:
		return "F12"
	default:
		return fmt.Sprintf("Unknown (0x%X)", keyCode)
	}
}

func main() {
	display, err := window.DisplayCreate(nil)
	if err != nil {
		panic(err)
	}

	app := &KeyboardTestApp{}
	app.window = window.Create(display)
	app.window.SetTitle("Keyboard Test")

	// Set keyboard handler
	app.window.SetKeyboardHandler(app)

	// Add widget
	app.widget = app.window.AddWidget(app)

	// Set initial size
	app.window.ScheduleResize(600, 300)

	fmt.Println("Keyboard test started. Press any key to see events.")
	fmt.Println("Press Cmd+Q to quit.")

	window.DisplayRun(display)
}
