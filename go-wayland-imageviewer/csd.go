package main

import (
	"github.com/fogleman/gg"
	"image"
	"unicode/utf8"
)

type DecorationButton struct {
	Text string
	Width int
	Pear bool
}

type Decoration struct {
	Title string
	Titlebar int
	LeftActive int
	LeftButtons []DecorationButton
	RightActive int
	RightButtons []DecorationButton
	Maximized bool
}

type DecoratedRgbaWindow interface {
	Focused() bool
	GetImage() *image.RGBA
	SetImage(img *image.RGBA)
}

const Border = 8
const bm1 = Border - 1
const hb = Border / 2
const db = Border * 2

var decorationFonts = []string{
	"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
	"/usr/share/fonts/truetype/freefont/FreeSansBold.ttf",
	"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
	"/usr/share/fonts/truetype/liberation2/LiberationSans-Bold.ttf",
}

func cutString(dc *gg.Context, w float64, out string) string {
	for len(out) > 0 {
	
		tw, _ := dc.MeasureString(out)
		if tw < w {
			return out
		}
		_, size := utf8.DecodeLastRuneInString(out)
		out = out[:len(out)-size]
	}
	return ""
}


func (d *Decoration) activeLeftRight(a DecoratedRgbaWindow, x, y float64) (l int, r int) {
	l = 0
	r = 0
	
	img := a.GetImage()
	w := float64(img.Rect.Max.X - img.Rect.Min.X)
	
	var x_rd float64
	for i, v := range d.RightButtons {
	
		lpos := x_rd+w-Border-1-float64(v.Width)
	
		if lpos < Border {
			break
		}
		pear := 0.
		if (!v.Pear) {
			pear = 1.
		}
		
		if x >= lpos && y >= Border/(2-pear)+2 && x < lpos+float64(v.Width) && y < Border/(2-pear)+1+float64(d.Titlebar)-2-pear*Border {
			r = i+1
			break
		}
		_ = i
		_ = v
		x_rd -= float64(v.Width) + Border
	}
	
	var x_d float64
	for i, v := range d.LeftButtons {
		lpos := x_d+Border+1
	
		if lpos + float64(v.Width) + Border >= w+x_rd {
			break
		}
	
		pear := 0.
		if (!v.Pear) {
			pear = 1.
		}
		
		if x >= lpos && y >= Border/(2-pear)+2 && x < lpos + float64(v.Width) && y < Border/(2-pear)+1 + float64(d.Titlebar)-2-pear*Border {
			l = i+1
			break
		}
		
		_ = i
		_ = v
		x_d += float64(v.Width) + Border
	}
	return
}

func (d *Decoration) clientSideDecoration(a DecoratedRgbaWindow, just_border bool) {

	img := a.GetImage()

	if !just_border {
		pixp := &img.Pix
		(*pixp) = append(make([]uint8, img.Stride*(Border+d.Titlebar)), (*pixp)...)
		(*pixp) = append((*pixp), make([]uint8, img.Stride*(Border))...)

		img.Rect.Max.Y += 2*Border + d.Titlebar
	}
	dc := gg.NewContextForImage(img)

	w := float64(img.Rect.Max.X - img.Rect.Min.X)
	h := float64(img.Rect.Max.Y - img.Rect.Min.Y)

	if d.Maximized {
		dc.DrawRectangle(1, 1, w-2, h-2)
	} else {
		dc.DrawRoundedRectangle(1, 1, w-2, h-2, hb)
	}
	dc.SetFillRuleEvenOdd()
	dc.DrawRectangle(Border, Border+float64(d.Titlebar), w-2*Border, h-2*Border-float64(d.Titlebar))
	
	dc.SetRGB255(0, 0, 0)
	dc.SetLineWidth(2.)
	dc.StrokePreserve()
	
	if a.Focused() {
		dc.SetRGB255(192, 192, 192)
	} else {
		dc.SetRGB255(255, 255, 255)
	}
	dc.Fill()
	dc.ClearPath()
	dc.SetFillRuleWinding()
	
	for _, v := range decorationFonts {
		if err := dc.LoadFontFace(v, 12); err == nil {
			break
		}
	}
	

	
	var x_rd float64
	for i, v := range d.RightButtons {
	
		lpos := x_rd+w-Border-1-float64(v.Width)
	
		if lpos < Border {
			break
		}
		pear := 0.
		if (!v.Pear) {
			pear = 1.
		}
		dc.DrawRoundedRectangle(lpos, Border/(2-pear)+2, float64(v.Width), float64(d.Titlebar)-2-pear*Border, hb)
		dc.SetRGB255(0, 0, 0)
		dc.SetLineWidth(2.)
		dc.StrokePreserve()
		if ((i+1) == d.RightActive) != (!v.Pear) {
			dc.SetRGB255(192, 192, 192)
		} else {
			dc.SetRGB255(255, 255, 255)
		}
		dc.Fill()
		dc.ClearPath()
		
		icon := 0.125
		if v.Text != "_" {
			icon = 0.25
		}
		
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(cutString(dc, float64(v.Width), v.Text), lpos+0.5*float64(v.Width), Border + float64(d.Titlebar)*icon, 0.5, 0.5)
		_ = i
		_ = v
		x_rd -= float64(v.Width) + Border
	}
	
	var x_d float64
	for i, v := range d.LeftButtons {
		lpos := x_d+Border+1
	
		if lpos + float64(v.Width) + Border >= w+x_rd {
			break
		}
	
		pear := 0.
		if (!v.Pear) {
			pear = 1.
		}
		dc.DrawRoundedRectangle(lpos, Border/(2-pear)+2, float64(v.Width), float64(d.Titlebar)-2-pear*Border, hb)
		dc.SetRGB255(0, 0, 0)
		dc.SetLineWidth(2.)
		dc.StrokePreserve()
		if ((i+1) == d.LeftActive) != (!v.Pear) {
			dc.SetRGB255(192, 192, 192)
		} else {
			dc.SetRGB255(255, 255, 255)
		}
		dc.Fill()
		dc.ClearPath()
		
		icon := 0.125
		if v.Text != "_" {
			icon = 0.25
		}
		
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(cutString(dc, float64(v.Width), v.Text), lpos+0.5*float64(v.Width), Border + float64(d.Titlebar)*icon, 0.5, 0.5)
		_ = i
		_ = v
		x_d += float64(v.Width) + Border
	}

	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(cutString(dc, (w-x_d+x_rd)-2*Border, d.Title), (w-x_d+x_rd)/2+x_d, Border + float64(d.Titlebar)/3, 0.5, 0.5)

	a.SetImage(dc.Image().(*image.RGBA))
}
