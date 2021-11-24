package main

import "image"
import "image/color"
import "image/png"
import "bytes"

func bias_towards_larger(n int) int {
	n = 255 - (((n & 15) * ((n >> 4) & 15)) + ((n >> 8) & 15))
	return n
}

func reprocess_scrollbar_row(set func(x, y, r, g, b int), row []string, y int) {
	for x, a := range row {
		if (a != " ") && (a != "") && (a != "\t") {
			var hash = bias_towards_larger(int(hashstr(a)))
			set(x, y, hash, hash, hash)
		}
	}
}

func reprocess_scrollbar(file [][]string) (out []byte, err error) {

	width := 96
	height := len(file) * 2

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for y := 0; y < height; y++ {

		if y&1 == 1 {
			continue
		}

		for x := 0; x < width; x++ {

			reprocess_scrollbar_row(func(x, y, r, g, b int) {
				img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 0xff})
			}, file[y/2], y)

		}
	}

	var buffer bytes.Buffer

	err = png.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
