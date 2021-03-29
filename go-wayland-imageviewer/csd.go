package main

import (
	"github.com/fogleman/gg"
	"image"
)

type Decoration struct {
	Title string
}

type App interface {
	Pix() *[]uint8
	Stride() int
	IncreaseHeight(by int)
	Width() int
	Height() int
	FrameImage() image.Image
	Focused() bool
	SetImage(img *image.RGBA)
}

const Border = 8
const Titlebar = 20
const bm1 = Border - 1
const hb = Border / 2
const db = Border * 2

var decorationFonts = []string{
	"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
	"/usr/share/fonts/truetype/freefont/FreeSansBold.ttf",
	"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
	"/usr/share/fonts/truetype/liberation2/LiberationSans-Bold.ttf",
}

func (d *Decoration) clientSideDecoration(a App, just_border bool) {

	if !just_border {
		pixp := a.Pix()
		(*pixp) = append(make([]uint8, a.Stride()*(Border+Titlebar)), (*pixp)...)
		(*pixp) = append((*pixp), make([]uint8, a.Stride()*(Border))...)

		a.IncreaseHeight(2*Border + Titlebar)
	}
	dc := gg.NewContextForImage(a.FrameImage())

	w := float64(a.Width())
	h := float64(a.Height())

	dc.DrawRoundedRectangle(0, 0, w, h, hb)
	dc.SetFillRuleEvenOdd()
	dc.DrawRectangle(Border, Border+Titlebar, w-2*Border, h-2*Border-Titlebar)
	dc.SetRGB255(0, 0, 0)
	dc.Fill()

	dc.DrawRoundedRectangle(1, 1, w-2, h-2, bm1/2)
	dc.SetFillRuleEvenOdd()
	dc.DrawRectangle(bm1, bm1+Titlebar, w-2*bm1, h-2*bm1-Titlebar)

	if a.Focused() {
		dc.SetRGB255(192, 192, 192)
	} else {
		dc.SetRGB255(255, 255, 255)
	}
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	for _, v := range decorationFonts {
		if err := dc.LoadFontFace(v, 12); err == nil {
			break
		}
	}
	dc.DrawString(d.Title, Border, Border+Titlebar/2)

	a.SetImage(dc.Image().(*image.RGBA))
}
