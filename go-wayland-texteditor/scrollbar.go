package main

import (
	"image/color"
	"image/png"
	"strings"
)

const scrollTestFilename = "live.png"

func downloadScrollbarPatch(filename string) ([][3]byte, error) {
	resp := libInstance.Call("/scrollbar/"+filename, "")
	reader := strings.NewReader(resp)
	img, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}
	var buf [][3]byte
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			buf = append(buf, [3]byte{c.R, c.G, c.B})
		}
	}
	println("length=", len(buf))
	return buf, nil
}
