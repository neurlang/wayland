// Copyright 2021 Neurlang project
// SPDX-License-Identifier: MIT

package window

import (
	"image"
	"image/color"
	"math"
	
	"github.com/fogleman/gg"
	sys "github.com/neurlang/wayland/os"
	"github.com/neurlang/wayland/wl"
)

// Decoration constants
const (
	ShadowMargin   = 24 // graspable part of the border
	TitleHeight    = 24
	ButtonWidth    = 32
	SymDim         = 14
	ShadowBlurSize = 64
)

// Colors (RGBA)
var (
	ColTitle       = color.RGBA{0x08, 0x07, 0x06, 0xFF}
	ColTitleInact  = color.RGBA{0x30, 0x30, 0x30, 0xFF}
	ColButtonMin   = color.RGBA{0xFF, 0xBB, 0x00, 0xFF}
	ColButtonMax   = color.RGBA{0x23, 0x88, 0x23, 0xFF}
	ColButtonClose = color.RGBA{0xFB, 0x65, 0x42, 0xFF}
	ColButtonInact = color.RGBA{0x40, 0x40, 0x40, 0xFF}
	ColSym         = color.RGBA{0xF4, 0xF4, 0xEF, 0xFF}
	ColSymAct      = color.RGBA{0x20, 0x32, 0x2A, 0xFF}
	ColSymInact    = color.RGBA{0x90, 0x90, 0x90, 0xFF}
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
	data     []byte
	width    int32
	height   int32
	stride   int32
	inUse    bool
}

// WindowDecoration manages client-side decorations
type WindowDecoration struct {
	window      *Window
	shadowSurf  *DecorationSurface
	titleSurf   *DecorationSurface
	shadowBlur  *image.RGBA // Pre-rendered shadow tile
	active      bool
	hoverButton componentType
}

// NewWindowDecoration creates a new decoration manager
func NewWindowDecoration(window *Window) *WindowDecoration {
	d := &WindowDecoration{
		window: window,
		active: false,
	}
	
	// Create pre-rendered shadow blur tile
	d.createShadowBlur()
	
	return d
}

// createShadowBlur creates a pre-rendered blurred shadow tile
func (d *WindowDecoration) createShadowBlur() {
	const size = 128
	const boundary = 32
	
	d.shadowBlur = image.NewRGBA(image.Rect(0, 0, size, size))
	dc := gg.NewContextForRGBA(d.shadowBlur)
	
	dc.SetColor(color.Black)
	dc.DrawRectangle(float64(boundary), float64(boundary), 
		float64(size-2*boundary), float64(size-2*boundary))
	dc.Fill()
	
	// Apply blur
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
	contentWidth := d.window.mainSurface.allocation.Width
	contentHeight := d.window.mainSurface.allocation.Height
	
	width := contentWidth + 2*ShadowMargin
	height := contentHeight + 2*ShadowMargin + TitleHeight
	x := int32(-ShadowMargin)
	y := int32(-(ShadowMargin + TitleHeight))
	
	surf, err := d.createDecorationSurface(x, y, width, height, false)
	if err != nil {
		return err
	}
	
	d.shadowSurf = surf
	return nil
}

// createTitleSurface creates the title bar decoration surface
func (d *WindowDecoration) createTitleSurface() error {
	contentWidth := d.window.mainSurface.allocation.Width
	
	width := contentWidth
	height := int32(TitleHeight)
	x := int32(0)
	y := int32(-TitleHeight)
	
	surf, err := d.createDecorationSurface(x, y, width, height, true)
	if err != nil {
		return err
	}
	
	d.titleSurf = surf
	return nil
}

// createDecorationSurface creates a subsurface for decoration
func (d *WindowDecoration) createDecorationSurface(x, y, width, height int32, opaque bool) (*DecorationSurface, error) {
	display := d.window.Display
	
	// Create wayland surface
	wlSurf, err := display.compositor.CreateSurface()
	if err != nil {
		return nil, err
	}
	
	// Create subsurface
	parent := d.window.mainSurface.surface_
	wlSubsurf, err := display.subcompositor.GetSubsurface(wlSurf, parent)
	if err != nil {
		_ = wlSurf.Destroy()
		return nil, err
	}
	
	// Position subsurface
	_ = wlSubsurf.SetPosition(x, y)
	_ = wlSubsurf.PlaceBelow(parent)
	
	surf := &DecorationSurface{
		wlSurface:    wlSurf,
		wlSubsurface: wlSubsurf,
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		scale:        1,
	}
	
	return surf, nil
}

// drawShadow renders the shadow decoration
func (d *WindowDecoration) drawShadow() {
	if d.shadowSurf == nil {
		return
	}
	
	surf := d.shadowSurf
	
	// Create or reuse buffer
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
	
	// Create image from buffer
	img := imageFromBuffer(surf.buffer)
	dc := gg.NewContextForRGBA(img)
	
	// Clear with transparent
	dc.SetColor(color.Transparent)
	dc.Clear()
	
	// Render shadow
	renderShadow(dc, d.shadowBlur, 
		-ShadowMargin/2, -ShadowMargin/2,
		int(surf.width)+ShadowMargin, int(surf.height)+ShadowMargin,
		ShadowBlurSize, ShadowBlurSize)
	
	// Mask out the content area (make it transparent)
	dc.SetColor(color.Transparent)
	dc.DrawRectangle(float64(ShadowMargin), float64(ShadowMargin+TitleHeight),
		float64(d.window.mainSurface.allocation.Width),
		float64(d.window.mainSurface.allocation.Height))
	dc.SetFillStyle(gg.NewSolidPattern(color.Transparent))
	dc.Fill()
	
	// Copy back to buffer
	copyImageToBuffer(img, surf.buffer)
	
	// Attach and commit
	_ = surf.wlSurface.Attach(surf.buffer.wlBuffer, 0, 0)
	_ = surf.wlSurface.DamageBuffer(0, 0, surf.width, surf.height)
	_ = surf.wlSurface.Commit()
	surf.buffer.inUse = true
}

// drawTitleBar renders the title bar with buttons
func (d *WindowDecoration) drawTitleBar() {
	if d.titleSurf == nil {
		return
	}
	
	surf := d.titleSurf
	
	// Create or reuse buffer
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
	
	// Create image from buffer
	img := imageFromBuffer(surf.buffer)
	dc := gg.NewContextForRGBA(img)
	
	// Background
	colTitle := ColTitle
	if !d.active {
		colTitle = ColTitleInact
	}
	dc.SetColor(colTitle)
	dc.Clear()
	
	// Draw title text
	d.drawTitleText(dc, int(surf.width))
	
	// Draw buttons
	d.drawButton(dc, ComponentButtonMin, surf.width-3*ButtonWidth, 0)
	d.drawButton(dc, ComponentButtonMax, surf.width-2*ButtonWidth, 0)
	d.drawButton(dc, ComponentButtonClose, surf.width-ButtonWidth, 0)
	
	// Copy back to buffer
	copyImageToBuffer(img, surf.buffer)
	
	// Attach and commit
	_ = surf.wlSurface.Attach(surf.buffer.wlBuffer, 0, 0)
	_ = surf.wlSurface.DamageBuffer(0, 0, surf.width, surf.height)
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
	
	// Load font and measure text
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 12); err != nil {
		// Fallback if font not found
		dc.SetColor(colText)
		return
	}
	
	textWidth, textHeight := dc.MeasureString(d.window.title)
	
	// Center text, but keep away from buttons
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
	// Button background
	colTitle := ColTitle
	if !d.active {
		colTitle = ColTitleInact
	}
	
	var btnCol color.Color
	isHover := d.hoverButton == btnType
	
	switch btnType {
	case ComponentButtonMin:
		if isHover && d.active {
			btnCol = ColButtonMin
		} else {
			btnCol = colTitle
		}
	case ComponentButtonMax:
		if isHover && d.active {
			btnCol = ColButtonMax
		} else {
			btnCol = colTitle
		}
	case ComponentButtonClose:
		if isHover && d.active {
			btnCol = ColButtonClose
		} else {
			btnCol = colTitle
		}
	default:
		btnCol = colTitle
	}
	
	dc.SetColor(btnCol)
	dc.DrawRectangle(float64(x), float64(y), ButtonWidth, TitleHeight)
	dc.Fill()
	
	// Button symbol
	symCol := ColSym
	if !d.active {
		symCol = ColSymInact
	} else if isHover {
		symCol = ColSymAct
	}
	
	dc.SetColor(symCol)
	dc.SetLineWidth(1)
	
	// Symbol position
	symX := float64(x) + ButtonWidth/2.0 - SymDim/2.0 + 0.5
	symY := float64(y) + TitleHeight/2.0 - SymDim/2.0 + 0.5
	
	switch btnType {
	case ComponentButtonMin:
		// Minimize: horizontal line at bottom
		dc.DrawLine(symX, symY+SymDim-1, symX+SymDim-1, symY+SymDim-1)
		dc.Stroke()
		
	case ComponentButtonMax:
		// Maximize: rectangle
		if d.window.maximized {
			// Two overlapping rectangles for "restore"
			const small = 12
			dc.DrawRectangle(symX, symY+SymDim-small, small-1, small-1)
			dc.Stroke()
			dc.MoveTo(symX+SymDim-small, symY+SymDim-small)
			dc.LineTo(symX+SymDim-small, symY)
			dc.LineTo(symX+SymDim-1, symY)
			dc.LineTo(symX+SymDim-1, symY+small-1)
			dc.LineTo(symX+small-1, symY+small-1)
			dc.Stroke()
		} else {
			dc.DrawRectangle(symX, symY, SymDim-1, SymDim-1)
			dc.Stroke()
		}
		
	case ComponentButtonClose:
		// Close: X
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
		_ = surf.wlSurface.Destroy()
		surf.wlSurface = nil
	}
}

// createBuffer creates a shared memory buffer
func (d *WindowDecoration) createBuffer(width, height int32, opaque bool) (*decorationBuffer, error) {
	display := d.window.Display
	
	stride := int32(width * 4) // ARGB32
	size := stride * height
	
	// Create anonymous file
	fd, err := sys.CreateAnonymousFile(int64(size))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	
	// Map memory
	data, err := sys.Mmap(int(fd.Fd()), 0, int(size), sys.ProtRead|sys.ProtWrite, sys.MapShared)
	if err != nil {
		return nil, err
	}
	
	// Create shm pool
	pool, err := display.shm.CreatePool(fd.Fd(), size)
	if err != nil {
		_ = sys.Munmap(data)
		return nil, err
	}
	
	// Create buffer
	format := uint32(wl.ShmFormatArgb8888)
	if opaque {
		format = uint32(wl.ShmFormatXrgb8888)
	}
	
	wlBuf, err := pool.CreateBuffer(0, width, height, stride, format)
	if err != nil {
		_ = pool.Destroy()
		_ = sys.Munmap(data)
		return nil, err
	}
	
	_ = pool.Destroy()
	
	buf := &decorationBuffer{
		wlBuffer: wlBuf,
		data:     data,
		width:    width,
		height:   height,
		stride:   stride,
		inUse:    false,
	}
	
	return buf, nil
}

// freeBuffer frees a buffer
func (d *WindowDecoration) freeBuffer(buf *decorationBuffer) {
	if buf.wlBuffer != nil {
		_ = buf.wlBuffer.Destroy()
		buf.wlBuffer = nil
	}
	if buf.data != nil {
		_ = sys.Munmap(buf.data)
		buf.data = nil
	}
}

// imageFromBuffer creates an RGBA image from a buffer
func imageFromBuffer(buf *decorationBuffer) *image.RGBA {
	return &image.RGBA{
		Pix:    buf.data,
		Stride: int(buf.stride),
		Rect:   image.Rect(0, 0, int(buf.width), int(buf.height)),
	}
}

// copyImageToBuffer copies image data back to buffer (no-op since they share memory)
func copyImageToBuffer(img *image.RGBA, buf *decorationBuffer) {
	// No-op: img.Pix and buf.data point to the same memory
}

// blurImage applies a Gaussian blur to an RGBA image
func blurImage(img *image.RGBA, margin int) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Create destination buffer
	dst := image.NewRGBA(bounds)
	
	// Gaussian kernel
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
				R: uint8(r / a),
				G: uint8(g / a),
				B: uint8(b / a),
				A: uint8(alpha / a),
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
				R: uint8(r / a),
				G: uint8(g / a),
				B: uint8(b / a),
				A: uint8(alpha / a),
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
		
		// Adjust if shadows are larger than surface
		if height < 2*shadowHeight {
			shadowHeight = (height + (1 - fy)) / 2
		}
		if width < 2*shadowWidth {
			shadowWidth = (width + (1 - fx)) / 2
		}
		
		// Calculate source and destination rectangles
		srcX := fx * (128 - shadowWidth)
		srcY := fy * (128 - shadowHeight)
		dstX := x + fx*(width-shadowWidth)
		dstY := y + fy*(height-shadowHeight)
		
		// Draw corner
		dc.Push()
		dc.DrawRectangle(float64(dstX), float64(dstY), float64(shadowWidth), float64(shadowHeight))
		dc.Clip()
		
		// Extract and draw the tile section
		tilePart := shadowTile.SubImage(image.Rect(srcX, srcY, srcX+shadowWidth, srcY+shadowHeight))
		dc.DrawImage(tilePart, dstX, dstY)
		dc.Pop()
	}
	
	// Top and bottom stretches
	shadowWidth := width - 2*margin
	shadowHeight := topMargin
	if height < 2*shadowHeight {
		shadowHeight = height / 2
	}
	
	if shadowWidth > 0 && shadowHeight > 0 {
		// Top stretch
		dc.Push()
		dc.DrawRectangle(float64(x+margin), float64(y), float64(shadowWidth), float64(shadowHeight))
		dc.Clip()
		
		// Stretch middle part of shadow tile
		tilePart := shadowTile.SubImage(image.Rect(60, 0, 68, shadowHeight))
		for sx := x + margin; sx < x+margin+shadowWidth; sx += 8 {
			dc.DrawImage(tilePart, sx, y)
		}
		dc.Pop()
		
		// Bottom stretch
		dc.Push()
		dc.DrawRectangle(float64(x+margin), float64(y+height-margin), float64(shadowWidth), float64(margin))
		dc.Clip()
		
		tilePart = shadowTile.SubImage(image.Rect(60, 128-margin, 68, 128))
		for sx := x + margin; sx < x+margin+shadowWidth; sx += 8 {
			dc.DrawImage(tilePart, sx, y+height-margin)
		}
		dc.Pop()
	}
	
	// Left and right stretches
	shadowWidth = margin
	if width < 2*shadowWidth {
		shadowWidth = width / 2
	}
	shadowHeight = height - margin - topMargin
	
	if shadowHeight > 0 && shadowWidth > 0 {
		// Left stretch
		dc.Push()
		dc.DrawRectangle(float64(x), float64(y+topMargin), float64(shadowWidth), float64(shadowHeight))
		dc.Clip()
		
		tilePart := shadowTile.SubImage(image.Rect(0, 60, shadowWidth, 68))
		for sy := y + topMargin; sy < y+topMargin+shadowHeight; sy += 8 {
			dc.DrawImage(tilePart, x, sy)
		}
		dc.Pop()
		
		// Right stretch
		dc.Push()
		dc.DrawRectangle(float64(x+width-shadowWidth), float64(y+topMargin), float64(shadowWidth), float64(shadowHeight))
		dc.Clip()
		
		tilePart = shadowTile.SubImage(image.Rect(128-shadowWidth, 60, 128, 68))
		for sy := y + topMargin; sy < y+topMargin+shadowHeight; sy += 8 {
			dc.DrawImage(tilePart, x+width-shadowWidth, sy)
		}
		dc.Pop()
	}
}
