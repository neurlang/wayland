package main

import "net/http"
import "image/png"
import "image/color"

const scrollTestFilename = "live.png"

func downloadScrollbarPatch(filename string) ([][3]byte, error) {
	resp, err := http.Get("http://127.0.0.1:8080/scrollbar/" + filename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, err := png.Decode(resp.Body)
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
