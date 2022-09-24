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

	var comments, strings, is_golang, is_csharp bool

	is_golang = detect_golang(file)
	is_csharp = detect_csharp(file)

	// Set color for each pixel.
	for y := 0; y < height; y++ {

		var color_out [][5]int

		if y&1 == 1 {
			continue
		}

		if is_golang {
			color_out = append(color_out, reprocess_syntax_highlighting_row_golang(file[y/2], y/2, &comments, &strings)...)
		}
		if is_csharp {
			color_out = append(color_out, reprocess_syntax_highlighting_row_csharp(file[y/2], y/2, &comments, &strings)...)
		}

		var xx = 0
		var docolorxx = -1

		reprocess_scrollbar_row(func(x, y, r, g, b int) {

			for ; xx < len(color_out); xx++ {
				if color_out[xx][0] >= x {
					if color_out[xx][0] == x {
						docolorxx = xx
					}
					break
				} else {
					docolorxx = xx
				}
			}

			if (docolorxx >= 0) && (docolorxx < len(color_out)) {
				r *= color_out[docolorxx][2]
				g *= color_out[docolorxx][3]
				b *= color_out[docolorxx][4]
				r /= 255
				g /= 255
				b /= 255
			}

			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 0xff})
		}, file[y/2], y)

	}

	var buffer bytes.Buffer

	err = png.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
