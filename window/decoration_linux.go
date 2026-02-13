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
)

// Decoration constants
const (
	ShadowMargin   = 8  // graspable part of the border
	TitleHeight    = 24
	ButtonWidth    = 32
	SymDim         = 14
	ShadowBlurSize = 64
)

// Colors (RGBA)
var (
	ColTitle          = color.RGBA{0xEB, 0xEB, 0xEB, 0xFF} // Light gray titlebar
	ColTitleInact     = color.RGBA{0xF6, 0xF5, 0xF4, 0xFF} // Very light gray when inactive
	ColButtonMinHover = color.RGBA{0xDA, 0xDA, 0xDA, 0xFF} // Slightly darker on hover
	ColButtonMaxHover = color.RGBA{0xDA, 0xDA, 0xDA, 0xFF} // Slightly darker on hover
	ColButtonCloseHover = color.RGBA{0xE0, 0x1B, 0x24, 0xFF} // Red on hover (GNOME style)
	ColSym            = color.RGBA{0x2E, 0x34, 0x36, 0xFF} // Dark gray/black for symbols
	ColSymClose       = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // White symbol on red close button
	ColSymInact       = color.RGBA{0x9A, 0x99, 0x96, 0xFF} // Gray when inactive
)

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
	window       *Window
	shadowSurf   *DecorationSurface
	titleSurf    *DecorationSurface
	shadowBlur   *image.RGBA
	active       bool
	hoverButton  componentType
	pointerX     float32
	pointerY     float32
	pointerSerial uint32
	isDragging   bool
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

// drawShadow renders the shadow decoration
func (d *WindowDecoration) drawShadow() {
	if d.shadowSurf == nil {
		return
	}
	surf := d.shadowSurf

	if surf.buffer == nil || surf.buffer.width != surf.width || surf.buffer.height != surf.height {
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
		alpha := uint8(255 * (fi/shadowSize) * 0.5) // Fade from 0% to 50%, darkest near window
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

	_ = surf.wlSurface.Attach(surf.buffer.wlBuffer, 0, 0)
	_ = surf.wlSurface.Damage(0, 0, surf.width, surf.height)
	_ = surf.wlSurface.Commit()
	surf.buffer.inUse = true
}

// drawTitleBar renders the title bar with buttons
func (d *WindowDecoration) drawTitleBar() {
	if d.titleSurf == nil {
		return
	}
	surf := d.titleSurf

	if surf.buffer == nil || surf.buffer.width != surf.width || surf.buffer.height != surf.height {
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

	d.drawTitleText(dc, int(surf.width))
	d.drawButton(dc, ComponentButtonMin, surf.width-3*ButtonWidth, 0)
	d.drawButton(dc, ComponentButtonMax, surf.width-2*ButtonWidth, 0)
	d.drawButton(dc, ComponentButtonClose, surf.width-ButtonWidth, 0)

	_ = surf.wlSurface.Attach(surf.buffer.wlBuffer, 0, 0)
	_ = surf.wlSurface.Damage(0, 0, surf.width, surf.height)
	_ = surf.wlSurface.Commit()
	surf.buffer.inUse = true
}

// drawTitleText renders the window title
func (d *WindowDecoration) drawTitleText(dc *gg.Context, titleWidth int) {
	if d.window.title == "" {
		return
	}

	colText := ColSym
	if !d.active {
		colText = ColSymInact
	}

	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 12); err != nil {
		return
	}

	textWidth, textHeight := dc.MeasureString(d.window.title)

	textX := float64(titleWidth)/2.0 - textWidth/2.0
	maxTextX := float64(titleWidth - 3*ButtonWidth - 10)
	if textX+textWidth > maxTextX {
		textX = maxTextX - textWidth
	}
	if textX < ButtonWidth {
		textX = ButtonWidth
	}
	textY := float64(TitleHeight)/2.0 + textHeight/2.0

	dc.SetColor(colText)
	dc.DrawString(d.window.title, textX, textY)
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
		d.drawTitleBar()
	}
}

// Redraw redraws all decoration surfaces
func (d *WindowDecoration) Redraw() {
	d.drawShadow()
	d.drawTitleBar()
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

// renderShadow renders a shadow using a pre-blurred tile
func renderShadow(dc *gg.Context, shadowTile *image.RGBA,
	x, y, width, height, margin, topMargin int) {

	// Render four corners
	for i := 0; i < 4; i++ {
		fx := i & 1
		fy := i >> 1

		shadowWidth := margin
		shadowHeight := topMargin
		if fy != 0 {
			shadowHeight = margin
		}
		if height < 2*shadowHeight {
			shadowHeight = (height + (1 - fy)) / 2
		}
		if width < 2*shadowWidth {
			shadowWidth = (width + (1 - fx)) / 2
		}

		srcX := fx * (128 - shadowWidth)
		srcY := fy * (128 - shadowHeight)
		dstX := x + fx*(width-shadowWidth)
		dstY := y + fy*(height-shadowHeight)

		dc.Push()
		dc.DrawRectangle(float64(dstX), float64(dstY), float64(shadowWidth), float64(shadowHeight))
		dc.Clip()
		tilePart := shadowTile.SubImage(image.Rect(srcX, srcY, srcX+shadowWidth, srcY+shadowHeight))
		dc.DrawImage(tilePart, dstX, dstY)
		dc.Pop()
	}

	// Top and bottom stretches
	sw := width - 2*margin
	sh := topMargin
	if height < 2*sh {
		sh = height / 2
	}
	if sw > 0 && sh > 0 {
		dc.Push()
		dc.DrawRectangle(float64(x+margin), float64(y), float64(sw), float64(sh))
		dc.Clip()
		tilePart := shadowTile.SubImage(image.Rect(60, 0, 68, sh))
		for sx := x + margin; sx < x+margin+sw; sx += 8 {
			dc.DrawImage(tilePart, sx, y)
		}
		dc.Pop()

		dc.Push()
		dc.DrawRectangle(float64(x+margin), float64(y+height-margin), float64(sw), float64(margin))
		dc.Clip()
		tilePart = shadowTile.SubImage(image.Rect(60, 128-margin, 68, 128))
		for sx := x + margin; sx < x+margin+sw; sx += 8 {
			dc.DrawImage(tilePart, sx, y+height-margin)
		}
		dc.Pop()
	}

	// Left and right stretches
	sw = margin
	if width < 2*sw {
		sw = width / 2
	}
	sh = height - margin - topMargin
	if sh > 0 && sw > 0 {
		dc.Push()
		dc.DrawRectangle(float64(x), float64(y+topMargin), float64(sw), float64(sh))
		dc.Clip()
		tilePart := shadowTile.SubImage(image.Rect(0, 60, sw, 68))
		for sy := y + topMargin; sy < y+topMargin+sh; sy += 8 {
			dc.DrawImage(tilePart, x, sy)
		}
		dc.Pop()

		dc.Push()
		dc.DrawRectangle(float64(x+width-sw), float64(y+topMargin), float64(sw), float64(sh))
		dc.Clip()
		tilePart = shadowTile.SubImage(image.Rect(128-sw, 60, 128, 68))
		for sy := y + topMargin; sy < y+topMargin+sh; sy += 8 {
			dc.DrawImage(tilePart, x+width-sw, sy)
		}
		dc.Pop()
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
