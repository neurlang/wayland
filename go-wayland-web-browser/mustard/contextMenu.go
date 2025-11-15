package mustard

import (
	"image"

	gg "github.com/danfragoso/thdwb/gg"
	"github.com/goki/freetype/truetype"
	assets "github.com/neurlang/wayland/go-wayland-web-browser/assets"
	window "github.com/neurlang/wayland/windowtrace"
	wl "github.com/neurlang/wayland/wl"
)

func (window *Window) EnableContextMenus() {
	window.contextMenu = &contextMenu{
		overlay: &Overlay{
			ref: "contextMenu",
		},
	}
}

func (window *Window) AddContextMenuEntry(entryText string, action func()) {
	window.contextMenu.entries = append(
		window.contextMenu.entries,
		&menuEntry{
			entryText: entryText,
			action:    action,
		},
	)
}

func (window *Window) DestroyContextMenu() {
	window.RemoveOverlay(window.contextMenu.overlay)
	window.contextMenu.entries = nil
	window.contextMenu.selectedEntry = nil
	window.SetCursor(DefaultCursor)
}

func prepEntry(ctx *gg.Context, entry string, width float64) string {
	w, _ := ctx.MeasureString(entry)

	if w < width {
		return entry
	}

	for i := 0; i < len(entry); i++ {
		nW, _ := ctx.MeasureString(entry[:len(entry)-i] + "...")

		if nW <= width {
			return entry[:len(entry)-i] + "..."
		}
	}

	return entry
}

func (window *Window) PrepareContextMenu() {

	menuWidth := float64(200)
	menuHeight := float64(len(window.contextMenu.entries) * 20)

	overlay := extractOverlay(
		menuWidth, menuHeight,
		image.Point{
			int(window.cursorX + menuWidth/2),
			int(window.cursorY + menuHeight/2),
		})

	window.SetContextMenuOverlay(overlay)
}
func (window *Window) DrawContextMenu() {
	window.PrepareContextMenu()
	window.contextMenu.DrawContextMenu()

}
func (contextMenu *contextMenu) DrawContextMenu() {

	menuLeft := float64(0)
	menuTop := float64(0)

	menuWidth := float64(200)
	menuHeight := float64(len(contextMenu.entries) * 20)

	surf := contextMenu.overlay.popup.PopupGetSurface()
	ctx := makeContextFromCairo(surf)

	ctx.DrawRectangle(0, 0, menuWidth, menuHeight)
	ctx.SetHexColor("#eee")
	ctx.Fill()

	font, _ := truetype.Parse(assets.OpenSans(400))
	ctx.SetHexColor("#222")
	ctx.SetFont(font, 16)

	textLeft := 4.

	for idx, entry := range contextMenu.entries {
		top, left := 16+float64(idx*20), 0.

		entry.setCoords(menuTop+top-16, menuLeft+left, menuWidth, 20)
		ctx.DrawString(prepEntry(ctx, entry.entryText, menuWidth-textLeft), textLeft, top)
		ctx.Fill()
	}

	ctx.DrawRectangle(0, 0, menuWidth, menuHeight)
	ctx.SetHexColor("#ddd")
	ctx.Stroke()

}

func (window *Window) SetContextMenuOverlay(overlay *Overlay) {
	window.contextMenu.overlay = overlay
	window.AddOverlay(overlay)

}
func (cm *contextMenu) Configure() *window.Widget {
	overlay := cm.overlay
	overlay.widget = overlay.window.window.AddPopupWidget(overlay.popup, overlay)
	return overlay.widget
}
func (cm *contextMenu) Done() {
	cm.entries = nil
	cm.selectedEntry = nil
	cm.overlay.window.DestroyContextMenu()

}
func (cm *contextMenu) Render(s Surface, time uint32) {
	ctx := makeContextFromCairo(s)

	menuWidth := float64(s.ImageSurfaceGetWidth())
	textLeft := 4.
	ctx.SetHexColor("#eee")
	ctx.Clear()

	font, _ := truetype.Parse(assets.OpenSans(400))
	ctx.SetHexColor("#222")
	ctx.SetFont(font, 16)

	for idx, entry := range cm.entries {
		if cm.selectedEntry == entry {
			ctx.DrawRectangle(0, float64(idx*20), menuWidth, 20)
			ctx.SetHexColor("#ccc")
			ctx.Fill()
		}

		ctx.SetHexColor("#222")
		ctx.DrawString(prepEntry(ctx, entry.entryText, menuWidth-textLeft), textLeft, 16+float64(idx*20))
		ctx.Fill()
	}
}

func extractOverlay(width, height float64, position image.Point) *Overlay {
	return &Overlay{
		ref:    "contextMenu",
		active: true,

		top:  0,
		left: 0,

		width:  width,
		height: height,

		position: position,
	}
}

func (window *Window) SelectEntry(entry *menuEntry) {
	window.contextMenu.selectedEntry = entry
	//window.PrepareContextMenu()
	window.SetCursor(PointerCursor)
}

func (window *Window) DeselectEntries() {
	if window.contextMenu.selectedEntry != nil {
		window.contextMenu.selectedEntry = nil
		//window.PrepareContextMenu()
		window.SetCursor(DefaultCursor)
	}
}

func (entry *menuEntry) PointIntersects(x, y float64) bool {
	top, left, width, height := entry.getCoords()
	if x > left &&
		x < left+width &&
		y > top &&
		y < top+height {
		return true
	}

	return false
}

func (entry *menuEntry) getCoords() (float64, float64, float64, float64) {
	return entry.top, entry.left, entry.width, entry.height
}

func (entry *menuEntry) setCoords(top, left, width, height float64) {
	entry.top, entry.left = top, left
	entry.width, entry.height = width, height
}

func (overlay *Overlay) Axis(widget *window.Widget, input *window.Input, time uint32, axis uint32, value float32) {
}
func (overlay *Overlay) AxisSource(widget *window.Widget, input *window.Input, source uint32) {
}
func (overlay *Overlay) AxisStop(widget *window.Widget, input *window.Input, time uint32, axis uint32) {
}
func (overlay *Overlay) AxisDiscrete(widget *window.Widget, input *window.Input, axis uint32, discrete int32) {
	if axis == 0 {
		overlay.window.ProcessScroll(0, -float64(discrete))
	} else {
		overlay.window.ProcessScroll(-float64(discrete), 0)
	}
}
func (overlay *Overlay) Motion(widget *window.Widget, input *window.Input, time uint32, x float32, y float32) int {

	overlay.window.cursorX = float64(x)
	overlay.window.cursorY = float64(y)

	for _, f := range overlay.window.pointerPositionEventListeners {

		f(float64(x), float64(y))
	}

	overlay.window.ProcessPointerPosition()

	return overlay.window.cursor
}
func (overlay *Overlay) Button(widget *window.Widget, input *window.Input, time uint32, button uint32, state wl.PointerButtonState, data window.WidgetHandler) {

	if state == 1 {
		overlay.window.clickSerial = overlay.window.window.Display.GetSerial()
		return
	}

	overlay.window.ProcessPointerClick(int(button))
}

func (overlay *Overlay) PointerFrame(widget *window.Widget, input *window.Input) {
}
func (overlay *Overlay) Enter(widget *window.Widget, input *window.Input, x float32, y float32) {

	overlay.window.cursorX = float64(x)
	overlay.window.cursorY = float64(y)

	for _, f := range overlay.window.pointerPositionEventListeners {

		f(float64(x), float64(y))
	}

	overlay.window.ProcessPointerPosition()

}
func (overlay *Overlay) Leave(widget *window.Widget, input *window.Input) {
}
func (overlay *Overlay) Redraw(widget *window.Widget) {
}
func (overlay *Overlay) Resize(widget *window.Widget, width int32, height int32, pwidth int32, pheight int32) {
}

func (overlay *Overlay) TouchUp(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32) {
}
func (overlay *Overlay) TouchDown(widget *window.Widget, input *window.Input, serial uint32, time uint32, id int32, x float32, y float32) {
}
func (overlay *Overlay) TouchMotion(widget *window.Widget, input *window.Input, time uint32, id int32, x float32, y float32) {
}
func (overlay *Overlay) TouchFrame(widget *window.Widget, input *window.Input) {
}
func (overlay *Overlay) TouchCancel(widget *window.Widget, width int32, height int32) {
}
