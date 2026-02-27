// Copyright 2021 Neurlang project
// SPDX-License-Identifier: MIT

package window

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	sys "github.com/neurlang/wayland/os"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
)

// Decoration constants
const (
	ShadowMargin   = 8 // graspable part of the border
	TitleHeight    = 24
	ButtonWidth    = 32
	SymDim         = 14
	ShadowBlurSize = 64
)

// Colors (RGBA)
var (
	ColTitle            = color.RGBA{0xEB, 0xEB, 0xEB, 0xFF} // Light gray titlebar
	ColTitleInact       = color.RGBA{0xF6, 0xF5, 0xF4, 0xFF} // Very light gray when inactive
	ColButtonMinHover   = color.RGBA{0xDA, 0xDA, 0xDA, 0xFF} // Slightly darker on hover
	ColButtonMaxHover   = color.RGBA{0xDA, 0xDA, 0xDA, 0xFF} // Slightly darker on hover
	ColButtonCloseHover = color.RGBA{0xE0, 0x1B, 0x24, 0xFF} // Red on hover (GNOME style)
	ColSym              = color.RGBA{0x2E, 0x34, 0x36, 0xFF} // Dark gray/black for symbols
	ColSymClose         = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // White symbol on red close button
	ColSymInact         = color.RGBA{0x9A, 0x99, 0x96, 0xFF} // Gray when inactive
)

// Font paths to try for titlebar text
var decorationFonts = []string{
	"/usr/share/fonts/dejavu-sans-fonts/DejaVuSans-Bold.ttf", // fedora
	"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",   // ubuntu
	"/usr/share/fonts/truetype/freefont/FreeSansBold.ttf",
	"/usr/share/fonts/liberation-sans/LiberationSans-Bold.ttf", // fedora
	"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
	"/usr/share/fonts/truetype/liberation2/LiberationSans-Bold.ttf",
}

// Component types
type componentType int

const (
	ComponentNone componentType = iota
	ComponentShadow
	ComponentTitle
	ComponentButtonMin
	ComponentButtonMax
	ComponentButtonClose
)

// DecorationSurface represents a decoration surface (shadow or titlebar)
type DecorationSurface struct {
	wlSurface    *wl.Surface
	wlSubsurface *wl.Subsurface
	buffer       *decorationBuffer
	x, y         int32
	width        int32
	height       int32
	scale        int32
	frameCb      *wl.Callback
}

type decorationBuffer struct {
	wlBuffer *wl.Buffer
	wlPool   *wl.ShmPool
	data     []byte
	width    int32
	height   int32
	stride   int32
	inUse    bool
}

// WindowDecoration manages client-side decorations
type WindowDecoration struct {
	window              *Window
	shadowSurf          *DecorationSurface
	titleSurf           *DecorationSurface
	shadowBlur          *image.RGBA
	active              bool
	hoverButton         componentType
	pointerX            float32
	pointerY            float32
	pointerSerial       uint32
	pendingShadowRedraw bool
	pendingTitleRedraw  bool
}

// HandleCallbackDone handles frame callbacks for decoration surfaces
func (d *WindowDecoration) HandleCallbackDone(ev wl.CallbackDoneEvent) {
	// Determine which surface this callback is for
	if d.shadowSurf != nil && d.shadowSurf.frameCb == ev.C {
		wlclient.CallbackDestroy(ev.C)
		d.shadowSurf.frameCb = nil

		// If another redraw is pending, schedule it
		if d.pendingShadowRedraw {
			d.pendingShadowRedraw = false
			d.commitShadow()
		}
	} else if d.titleSurf != nil && d.titleSurf.frameCb == ev.C {
		wlclient.CallbackDestroy(ev.C)
		d.titleSurf.frameCb = nil

		// If another redraw is pending, schedule it
		if d.pendingTitleRedraw {
			d.pendingTitleRedraw = false
			d.commitTitleBar()
		}
	}
}

// NewWindowDecoration creates a new decoration manager
func NewWindowDecoration(window *Window) *WindowDecoration {
	d := &WindowDecoration{
		window: window,
		active: false,
	}
	d.createShadowBlur()
	return d
}

// createShadowBlur creates a pre-rendered blurred shadow tile
func (d *WindowDecoration) createShadowBlur() {
	const size = 128
	const boundary = 20 // Smaller boundary = larger shadow area

	d.shadowBlur = image.NewRGBA(image.Rect(0, 0, size, size))
	dc := gg.NewContextForRGBA(d.shadowBlur)

	// Draw fully opaque black - blur will make it semi-transparent at edges
	dc.SetColor(color.RGBA{0, 0, 0, 255})
	dc.DrawRectangle(float64(boundary), float64(boundary),
		float64(size-2*boundary), float64(size-2*boundary))
	dc.Fill()

	blurImage(d.shadowBlur, ShadowBlurSize)
}

// Show creates and displays the decoration surfaces
func (d *WindowDecoration) Show() error {
	if err := d.createShadowSurface(); err != nil {
		return err
	}
	if err := d.createTitleSurface(); err != nil {
		return err
	}
	d.drawShadow()
	d.drawTitleBar()
	return nil
}

// Hide hides the decoration surfaces
func (d *WindowDecoration) Hide() {
	if d.shadowSurf != nil {
		d.hideSurface(d.shadowSurf)
	}
	if d.titleSurf != nil {
		d.hideSurface(d.titleSurf)
	}
}

// Destroy cleans up decoration resources
func (d *WindowDecoration) Destroy() {
	if d.shadowSurf != nil {
		d.destroySurface(d.shadowSurf)
		d.shadowSurf = nil
	}
	if d.titleSurf != nil {
		d.destroySurface(d.titleSurf)
		d.titleSurf = nil
	}
	d.shadowBlur = nil
}

// createShadowSurface creates the shadow decoration surface
func (d *WindowDecoration) createShadowSurface() error {
	var contentWidth, contentHeight int32
	if d.window.mainSurface != nil && d.window.mainSurface.allocation.Width > 0 {
		contentWidth = d.window.mainSurface.allocation.Width
		contentHeight = d.window.mainSurface.allocation.Height
	} else {
		contentWidth = d.window.pendingAllocation.Width
		contentHeight = d.window.pendingAllocation.Height
	}

	fmt.Printf("DEBUG: Shadow surface content size: %dx%d\n", contentWidth, contentHeight)

	if contentWidth <= 0 || contentHeight <= 0 {
		return fmt.Errorf("invalid content dimensions: %dx%d", contentWidth, contentHeight)
	}
	const maxDim = 4096
	if contentWidth > maxDim || contentHeight > maxDim {
		return fmt.Errorf("content too large: %dx%d (max %d)", contentWidth, contentHeight, maxDim)
	}

	width := contentWidth + 2*ShadowMargin
	height := contentHeight + 2*ShadowMargin + TitleHeight
	x := int32(-ShadowMargin)
	y := int32(-(ShadowMargin + TitleHeight))

	fmt.Printf("DEBUG: Shadow surface total size: %dx%d at (%d,%d)\n", width, height, x, y)

	surf, err := d.createDecorationSurface(x, y, width, height)
	if err != nil {
		return err
	}
	d.shadowSurf = surf
	return nil
}

// createTitleSurface creates the title bar decoration surface
func (d *WindowDecoration) createTitleSurface() error {
	var contentWidth int32
	if d.window.mainSurface != nil && d.window.mainSurface.allocation.Width > 0 {
		contentWidth = d.window.mainSurface.allocation.Width
	} else {
		contentWidth = d.window.pendingAllocation.Width
	}
	if contentWidth <= 0 {
		return fmt.Errorf("invalid content width: %d", contentWidth)
	}

	surf, err := d.createDecorationSurface(0, -TitleHeight, contentWidth, TitleHeight)
	if err != nil {
		return err
	}
	d.titleSurf = surf
	return nil
}

// createDecorationSurface creates a subsurface for decoration
func (d *WindowDecoration) createDecorationSurface(x, y, width, height int32) (*DecorationSurface, error) {
	display := d.window.Display

	wlSurf, err := display.compositor.CreateSurface()
	if err != nil {
		return nil, fmt.Errorf("CreateSurface failed: %w", err)
	}

	parent := d.window.mainSurface.surface_
	wlSubsurf, err := display.subcompositor.GetSubsurface(wlSurf, parent)
	if err != nil {
		_ = wlSurf.Destroy()
		return nil, fmt.Errorf("CreateSurface failed: %w", err)
	}

	_ = wlSubsurf.SetPosition(x, y)
	_ = wlSubsurf.PlaceBelow(parent)

	// Register surface in the surface2window map for input handling
	display.surface2window[wlSurf] = d.window

	return &DecorationSurface{
		wlSurface:    wlSurf,
		wlSubsurface: wlSubsurf,
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		scale:        1,
	}, nil
}

// drawShadow renders the shadow decoration with frame callback synchronization
func (d *WindowDecoration) drawShadow() {
	if d.shadowSurf == nil {
		return
	}

	// If we're waiting for a frame callback, mark as pending
	if d.shadowSurf.frameCb != nil {
		d.pendingShadowRedraw = true
		return
	}

	d.renderShadow()
	d.commitShadow()
}

// renderShadow renders shadow to buffer without committing
func (d *WindowDecoration) renderShadow() {
	if d.shadowSurf == nil {
		return
	}
	surf := d.shadowSurf

	// Only recreate buffer if size changed
	needsNewBuffer := surf.buffer == nil ||
		surf.buffer.width != surf.width ||
		surf.buffer.height != surf.height

	if needsNewBuffer {
		if surf.buffer != nil {
			d.freeBuffer(surf.buffer)
		}
		buf, err := d.createBuffer(surf.width, surf.height, false)
		if err != nil {
			return
		}
		surf.buffer = buf
	}

	img := imageFromBuffer(surf.buffer)
	dc := gg.NewContextForRGBA(img)

	dc.SetColor(color.Transparent)
	dc.Clear()

	// Draw shadow as a gradient from edges inward
	shadowSize := float64(ShadowMargin)
	w := float64(surf.width)
	h := float64(surf.height)

	// Draw shadow on all four sides
	for i := 0; i < int(shadowSize); i++ {
		fi := float64(i)
		alpha := uint8(255 * (fi / shadowSize) * 0.5) // Fade from 0% to 50%, darkest near window
		shadowColor := color.RGBA{0, 0, 0, alpha}
		dc.SetColor(shadowColor)

		// Top shadow (excluding left and right corners at this level)
		dc.DrawRectangle(fi+1, fi, w-2*fi-2, 1)
		dc.Fill()

		// Bottom shadow (excluding left and right corners at this level)
		dc.DrawRectangle(fi+1, h-fi-1, w-2*fi-2, 1)
		dc.Fill()

		// Left shadow (full height)
		dc.DrawRectangle(fi, fi, 1, h-2*fi)
		dc.Fill()

		// Right shadow (full height)
		dc.DrawRectangle(w-fi-1, fi, 1, h-2*fi)
		dc.Fill()
	}

	// Mask out the content area (make it transparent)
	dc.SetColor(color.Transparent)
	dc.DrawRectangle(float64(ShadowMargin), float64(ShadowMargin+TitleHeight),
		float64(d.window.pendingAllocation.Width),
		float64(d.window.pendingAllocation.Height))
	dc.SetFillStyle(gg.NewSolidPattern(color.Transparent))
	dc.Fill()
}

// commitShadow commits the shadow surface with frame callback
func (d *WindowDecoration) commitShadow() {
	if d.shadowSurf == nil || d.shadowSurf.buffer == nil {
		return
	}
	surf := d.shadowSurf

	// Request frame callback before commit
	cb, err := surf.wlSurface.Frame()
	if err == nil {
		surf.frameCb = cb
		cb.AddDoneHandler(d)
	}

	_ = surf.wlSurface.Attach(surf.buffer.wlBuffer, 0, 0)
	_ = surf.wlSurface.Damage(0, 0, surf.width, surf.height)
	_ = surf.wlSurface.Commit()
	surf.buffer.inUse = true
}

// drawTitleBar renders the title bar with buttons with frame callback synchronization
func (d *WindowDecoration) drawTitleBar() {
	if d.titleSurf == nil {
		return
	}

	// If we're waiting for a frame callback, mark as pending
	if d.titleSurf.frameCb != nil {
		d.pendingTitleRedraw = true
		return
	}

	d.renderTitleBar()
	d.commitTitleBar()
}

// renderTitleBar renders titlebar to buffer without committing
func (d *WindowDecoration) renderTitleBar() {
	if d.titleSurf == nil {
		return
	}
	surf := d.titleSurf

	// Only recreate buffer if size changed
	needsNewBuffer := surf.buffer == nil ||
		surf.buffer.width != surf.width ||
		surf.buffer.height != surf.height

	if needsNewBuffer {
		if surf.buffer != nil {
			d.freeBuffer(surf.buffer)
		}
		buf, err := d.createBuffer(surf.width, surf.height, true)
		if err != nil {
			return
		}
		surf.buffer = buf
	}

	img := imageFromBuffer(surf.buffer)
	dc := gg.NewContextForRGBA(img)

	colTitle := ColTitle
	if !d.active {
		colTitle = ColTitleInact
	}
	dc.SetColor(colTitle)
	dc.Clear()

	// Draw buttons first
	d.drawButton(dc, ComponentButtonMin, surf.width-3*ButtonWidth, 0)
	d.drawButton(dc, ComponentButtonMax, surf.width-2*ButtonWidth, 0)
	d.drawButton(dc, ComponentButtonClose, surf.width-ButtonWidth, 0)

	// Draw title text last so it's on top
	d.drawTitleText(dc, int(surf.width))
}

// commitTitleBar commits the titlebar surface with frame callback
func (d *WindowDecoration) commitTitleBar() {
	if d.titleSurf == nil || d.titleSurf.buffer == nil {
		return
	}
	surf := d.titleSurf

	// Request frame callback before commit
	cb, err := surf.wlSurface.Frame()
	if err == nil {
		surf.frameCb = cb
		cb.AddDoneHandler(d)
	}

	_ = surf.wlSurface.Attach(surf.buffer.wlBuffer, 0, 0)
	_ = surf.wlSurface.Damage(0, 0, surf.width, surf.height)
	_ = surf.wlSurface.Commit()
	surf.buffer.inUse = true
}

// drawTitleText renders the window title
func (d *WindowDecoration) drawTitleText(dc *gg.Context, titleWidth int) {
	if d.window == nil || d.window.title == "" {
		return
	}

	// Load font - try each path until one works
	fontLoaded := false
	for _, fontPath := range decorationFonts {
		if err := dc.LoadFontFace(fontPath, 12); err == nil {
			fontLoaded = true
			break
		}
	}

	if !fontLoaded {
		// Can't load any font
		return
	}

	// Calculate available space for title (excluding buttons on right)
	availableWidth := float64(titleWidth - 3*ButtonWidth - 40)
	if availableWidth < 50 {
		return
	}

	// Truncate title if too long
	title := d.window.title
	textWidth, _ := dc.MeasureString(title)

	for textWidth > availableWidth && len(title) > 0 {
		title = title[:len(title)-1]
		textWidth, _ = dc.MeasureString(title + "...")
	}
	if len(title) < len(d.window.title) && len(title) > 0 {
		title = title + "..."
	}

	// Position text
	textX := 20.0 + availableWidth/2.0
	textY := float64(TitleHeight) / 2.0

	// Set text color - black when active, gray when inactive
	if !d.active {
		dc.SetRGB255(154, 153, 150)
	} else {
		dc.SetRGB255(0, 0, 0)
	}

	// Draw the title text
	dc.DrawStringAnchored(title, textX, textY, 0.5, 0.5)
}

// drawButton renders a window button (min/max/close)
func (d *WindowDecoration) drawButton(dc *gg.Context, btnType componentType, x, y int32) {
	colTitle := ColTitle
	if !d.active {
		colTitle = ColTitleInact
	}

	var btnCol color.Color
	var symCol color.Color
	isHover := d.hoverButton == btnType

	// Determine button background color
	switch btnType {
	case ComponentButtonMin:
		if isHover && d.active {
			btnCol = ColButtonMinHover
		} else {
			btnCol = colTitle
		}
	case ComponentButtonMax:
		if isHover && d.active {
			btnCol = ColButtonMaxHover
		} else {
			btnCol = colTitle
		}
	case ComponentButtonClose:
		if isHover && d.active {
			btnCol = ColButtonCloseHover
		} else {
			btnCol = colTitle
		}
	default:
		btnCol = colTitle
	}

	// Draw button background
	dc.SetColor(btnCol)
	dc.DrawRectangle(float64(x), float64(y), ButtonWidth, TitleHeight)
	dc.Fill()

	// Determine symbol color
	if !d.active {
		symCol = ColSymInact
	} else if btnType == ComponentButtonClose && isHover {
		// White symbol on red close button when hovering
		symCol = ColSymClose
	} else {
		symCol = ColSym
	}

	dc.SetColor(symCol)
	dc.SetLineWidth(1.5) // Slightly thicker lines for better visibility

	symX := float64(x) + ButtonWidth/2.0 - SymDim/2.0 + 0.5
	symY := float64(y) + TitleHeight/2.0 - SymDim/2.0 + 0.5

	switch btnType {
	case ComponentButtonMin:
		// Horizontal line for minimize
		dc.DrawLine(symX, symY+SymDim-1, symX+SymDim-1, symY+SymDim-1)
		dc.Stroke()
	case ComponentButtonMax:
		if d.window.maximized {
			// Two overlapping squares for unmaximize
			const small = 10
			const offset = 2
			// Back square
			dc.DrawRectangle(symX+offset, symY+offset, small, small)
			dc.Stroke()
			// Front square
			dc.DrawRectangle(symX, symY, small, small)
			dc.Stroke()
		} else {
			// Single square for maximize
			dc.DrawRectangle(symX, symY, SymDim-1, SymDim-1)
			dc.Stroke()
		}
	case ComponentButtonClose:
		// X symbol for close
		dc.DrawLine(symX, symY, symX+SymDim-1, symY+SymDim-1)
		dc.DrawLine(symX+SymDim-1, symY, symX, symY+SymDim-1)
		dc.Stroke()
	}
}

// SetActive updates the active state and redraws
func (d *WindowDecoration) SetActive(active bool) {
	if d.active != active {
		d.active = active
		d.drawTitleBar()
	}
}

// SetHoverButton updates the hover state and redraws
func (d *WindowDecoration) SetHoverButton(btn componentType) {
	if d.hoverButton != btn {
		d.hoverButton = btn
		// Only redraw titlebar if we're actually on the titlebar
		// Don't redraw when hovering over shadow/border areas
		if d.titleSurf != nil {
			d.drawTitleBar()
		}
	}
}

// Redraw redraws all decoration surfaces
func (d *WindowDecoration) Redraw() {
	d.drawShadow()
	d.drawTitleBar()
}

// UpdateSize updates decoration surfaces when window is resized
func (d *WindowDecoration) UpdateSize() {
	if d.window == nil || d.window.mainSurface == nil {
		return
	}

	var contentWidth, contentHeight int32
	if d.window.mainSurface.allocation.Width > 0 {
		contentWidth = d.window.mainSurface.allocation.Width
		contentHeight = d.window.mainSurface.allocation.Height
	} else {
		contentWidth = d.window.pendingAllocation.Width
		contentHeight = d.window.pendingAllocation.Height
	}

	if contentWidth <= 0 || contentHeight <= 0 {
		return
	}

	// Update shadow surface size
	if d.shadowSurf != nil {
		newWidth := contentWidth + 2*ShadowMargin
		newHeight := contentHeight + 2*ShadowMargin + TitleHeight

		// Only update if size actually changed
		if d.shadowSurf.width != newWidth || d.shadowSurf.height != newHeight {
			d.shadowSurf.width = newWidth
			d.shadowSurf.height = newHeight
			// Buffer will be recreated on next draw
			d.drawShadow()
		}
	}

	// Update titlebar surface size
	if d.titleSurf != nil {
		newWidth := contentWidth

		// Only update if size actually changed
		if d.titleSurf.width != newWidth {
			d.titleSurf.width = newWidth
			// Buffer will be recreated on next draw
			d.drawTitleBar()
		}
	}
}

// UpdateSizeForResize updates decoration surfaces during interactive resize
// This is called on every configure event during resize for smooth tracking
func (d *WindowDecoration) UpdateSizeForResize(contentWidth, contentHeight int32) {
	if contentWidth <= 0 || contentHeight <= 0 {
		return
	}

	// Update shadow surface size
	if d.shadowSurf != nil {
		newWidth := contentWidth + 2*ShadowMargin
		newHeight := contentHeight + 2*ShadowMargin + TitleHeight

		// Only update if size actually changed
		if d.shadowSurf.width != newWidth || d.shadowSurf.height != newHeight {
			d.shadowSurf.width = newWidth
			d.shadowSurf.height = newHeight
			// Render and commit with frame callback
			d.renderShadow()
			d.commitShadow()
		}
	}

	// Update titlebar surface size
	if d.titleSurf != nil {
		newWidth := contentWidth

		// Only update if size actually changed
		if d.titleSurf.width != newWidth {
			d.titleSurf.width = newWidth
			// Render and commit with frame callback
			d.renderTitleBar()
			d.commitTitleBar()
		}
	}
}

// hideSurface hides a decoration surface
func (d *WindowDecoration) hideSurface(surf *DecorationSurface) {
	if surf.wlSurface != nil {
		_ = surf.wlSurface.Attach(nil, 0, 0)
		_ = surf.wlSurface.Commit()
	}
}

// destroySurface destroys a decoration surface
func (d *WindowDecoration) destroySurface(surf *DecorationSurface) {
	if surf.frameCb != nil {
		wlclient.CallbackDestroy(surf.frameCb)
		surf.frameCb = nil
	}
	if surf.buffer != nil {
		d.freeBuffer(surf.buffer)
		surf.buffer = nil
	}
	if surf.wlSubsurface != nil {
		_ = surf.wlSubsurface.Destroy()
		surf.wlSubsurface = nil
	}
	if surf.wlSurface != nil {
		// Unregister from surface2window map
		delete(d.window.Display.surface2window, surf.wlSurface)
		_ = surf.wlSurface.Destroy()
		surf.wlSurface = nil
	}
}

// createBuffer creates a shared memory buffer
func (d *WindowDecoration) createBuffer(width, height int32, opaque bool) (*decorationBuffer, error) {
	display := d.window.Display

	stride := width * 4
	size := stride * height

	fd, err := sys.CreateAnonymousFile(int64(size))
	if err != nil {
		return nil, fmt.Errorf("CreateAnonymousFile failed: %w", err)
	}

	data, err := sys.Mmap(int(fd.Fd()), 0, int(size), sys.ProtRead|sys.ProtWrite, sys.MapShared)
	if err != nil {
		fd.Close()
		return nil, fmt.Errorf("Mmap failed: %w", err)
	}

	pool, err := display.shm.CreatePool(fd.Fd(), size)
	if err != nil {
		_ = sys.Munmap(data)
		fd.Close()
		return nil, fmt.Errorf("CreatePool failed: %w", err)
	}
	fd.Close()

	format := uint32(wl.ShmFormatArgb8888)
	if opaque {
		format = uint32(wl.ShmFormatXrgb8888)
	}

	wlBuf, err := pool.CreateBuffer(0, width, height, stride, format)
	if err != nil {
		_ = pool.Destroy()
		_ = sys.Munmap(data)
		return nil, fmt.Errorf("CreateBuffer failed: %w", err)
	}

	return &decorationBuffer{
		wlBuffer: wlBuf,
		wlPool:   pool,
		data:     data,
		width:    width,
		height:   height,
		stride:   stride,
	}, nil
}

// freeBuffer frees a buffer
func (d *WindowDecoration) freeBuffer(buf *decorationBuffer) {
	if buf.wlBuffer != nil {
		_ = buf.wlBuffer.Destroy()
	}
	if buf.wlPool != nil {
		_ = buf.wlPool.Destroy()
	}
	if buf.data != nil {
		_ = sys.Munmap(buf.data)
	}
}

// imageFromBuffer creates an RGBA image backed by the buffer's shared memory
func imageFromBuffer(buf *decorationBuffer) *image.RGBA {
	return &image.RGBA{
		Pix:    buf.data,
		Stride: int(buf.stride),
		Rect:   image.Rect(0, 0, int(buf.width), int(buf.height)),
	}
}

// blurImage applies a Gaussian blur to an RGBA image
func blurImage(img *image.RGBA, margin int) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dst := image.NewRGBA(bounds)

	const kernelSize = 71
	kernel := make([]uint32, kernelSize)
	half := kernelSize / 2
	var a uint32

	for i := 0; i < kernelSize; i++ {
		f := float64(i - half)
		kernel[i] = uint32(math.Exp(-f*f/float64(kernelSize)) * 10000)
		a += kernel[i]
	}

	// Horizontal pass
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if margin < x && x < width-margin {
				dst.Set(x, y, img.At(x, y))
				continue
			}
			var r, g, b, alpha uint32
			for k := 0; k < kernelSize; k++ {
				xk := x - half + k
				if xk < 0 || xk >= width {
					continue
				}
				c := img.RGBAAt(xk, y)
				r += uint32(c.R) * kernel[k]
				g += uint32(c.G) * kernel[k]
				b += uint32(c.B) * kernel[k]
				alpha += uint32(c.A) * kernel[k]
			}
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8(r / a), G: uint8(g / a),
				B: uint8(b / a), A: uint8(alpha / a),
			})
		}
	}

	// Vertical pass
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if margin <= y && y < height-margin {
				img.Set(x, y, dst.At(x, y))
				continue
			}
			var r, g, b, alpha uint32
			for k := 0; k < kernelSize; k++ {
				yk := y - half + k
				if yk < 0 || yk >= height {
					continue
				}
				c := dst.RGBAAt(x, yk)
				r += uint32(c.R) * kernel[k]
				g += uint32(c.G) * kernel[k]
				b += uint32(c.B) * kernel[k]
				alpha += uint32(c.A) * kernel[k]
			}
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(r / a), G: uint8(g / a),
				B: uint8(b / a), A: uint8(alpha / a),
			})
		}
	}
}

// HandlePointerEnter handles pointer entering the titlebar
func (d *WindowDecoration) HandlePointerEnter(serial uint32, x, y float32) {
	d.pointerSerial = serial
	d.pointerX = x
	d.pointerY = y
	d.updateHoverButton()
}

// HandlePointerLeave handles pointer leaving the titlebar
func (d *WindowDecoration) HandlePointerLeave() {
	if d.hoverButton != ComponentNone {
		d.SetHoverButton(ComponentNone)
	}
}

// HandlePointerMotion handles pointer motion over the titlebar
func (d *WindowDecoration) HandlePointerMotion(x, y float32) {
	d.pointerX = x
	d.pointerY = y
	d.updateHoverButton()
}

// HandlePointerButton handles pointer button clicks on the titlebar
func (d *WindowDecoration) HandlePointerButton(serial uint32, button uint32, state wl.PointerButtonState) {
	d.pointerSerial = serial

	if state == wl.PointerButtonStatePressed && button == 272 { // Left click (BTN_LEFT)
		switch d.hoverButton {
		case ComponentButtonClose:
			d.handleClose()
		case ComponentButtonMax:
			d.handleMaximize()
		case ComponentButtonMin:
			d.handleMinimize()
		case ComponentTitle:
			d.handleDragStart(serial)
		}
	}
}

// updateHoverButton determines which button the pointer is over
func (d *WindowDecoration) updateHoverButton() {
	if d.titleSurf == nil {
		return
	}

	// Pointer coordinates are relative to the titlebar surface
	x := d.pointerX
	y := d.pointerY

	// Check if pointer is within titlebar bounds
	if y < 0 || y >= float32(TitleHeight) || x < 0 || x >= float32(d.titleSurf.width) {
		d.SetHoverButton(ComponentNone)
		return
	}

	// Check buttons (from right to left)
	closeX := float32(d.titleSurf.width - ButtonWidth)
	maxX := float32(d.titleSurf.width - 2*ButtonWidth)
	minX := float32(d.titleSurf.width - 3*ButtonWidth)

	if x >= closeX {
		d.SetHoverButton(ComponentButtonClose)
	} else if x >= maxX {
		d.SetHoverButton(ComponentButtonMax)
	} else if x >= minX {
		d.SetHoverButton(ComponentButtonMin)
	} else {
		d.SetHoverButton(ComponentTitle)
	}
}

// handleDragStart initiates window dragging
func (d *WindowDecoration) handleDragStart(serial uint32) {
	if d.window.xdgToplevel == nil {
		return
	}

	// Get the seat from the display's input list
	if len(d.window.Display.inputList) == 0 {
		return
	}
	input := d.window.Display.inputList[0]
	if input.seat == nil {
		return
	}

	// Start interactive move
	_ = d.window.xdgToplevel.Move(input.seat, serial)
}

// handleClose closes the window
func (d *WindowDecoration) handleClose() {
	if d.window.closeHandler != nil {
		d.window.closeHandler.Close()
	} else {
		// Default behavior: exit the display (same as compositor close event)
		d.window.Display.Exit()
	}
}

// handleMaximize toggles window maximization
func (d *WindowDecoration) handleMaximize() {
	if d.window.xdgToplevel == nil {
		return
	}

	if d.window.maximized {
		_ = d.window.xdgToplevel.UnsetMaximized()
	} else {
		_ = d.window.xdgToplevel.SetMaximized()
	}
}

// handleMinimize minimizes the window
func (d *WindowDecoration) handleMinimize() {
	if d.window.xdgToplevel == nil {
		return
	}

	_ = d.window.xdgToplevel.SetMinimized()
}
