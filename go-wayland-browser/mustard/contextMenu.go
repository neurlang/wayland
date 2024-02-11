package mustard

import (
	assets "github.com/danfragoso/thdwb/assets"
	gg "github.com/danfragoso/thdwb/gg"
	"github.com/goki/freetype/truetype"
	"image"
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

	menuLeft := float64(0)
	menuTop := float64(0)

	//  menuWidth := float64(200)
	// menuHeight := float64(len(window.contextMenu.entries) * 20)

	overlay := extractOverlay(
		100, 100,
		image.Point{
			int(menuLeft),
			int(menuTop),
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

func extractOverlay(width, height float64, postion image.Point) *Overlay {
	return &Overlay{
		ref:    "contextMenu",
		active: true,

		top:  float64(postion.Y),
		left: float64(postion.X),

		width:  width,
		height: height,

		position: postion,
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
