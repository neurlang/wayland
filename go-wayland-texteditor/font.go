package main

import "os"
import "image/png"
import "image/color"
import "strings"
import "fmt"

type Font struct {
	cellx   int
	celly   int
	mapping map[string][][3]byte
}

func (f *Font) GetRGBTexture(code string) [][3]byte {

	var a, ok = f.mapping[code]
	if !ok {
		return nil
	}
	return a
}

func (f *Font) Load(name, descriptor string) error {
	file, err := os.Open(name)
	if err != nil {
		print("Font not found: ")
		println(name)
		return err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		print("Cannot decode png: ")
		println(name)
		return err
	}
	b := img.Bounds()

	var width = b.Max.X - b.Min.X
	var height = b.Max.Y - b.Min.Y

	var buffer = strings.Split(strings.ReplaceAll(descriptor, "\r\n", "\n"), "\n")
	var buf0 = strings.Split(buffer[0], "\t")

	var cellx = width / len(buf0)
	var celly = height / len(buffer)

	if f.mapping == nil {
		f.cellx = cellx
		f.celly = celly
	} else if f.cellx != cellx || f.celly != celly {
		return fmt.Errorf("only same cell sized fonts can be merged")
	}

	var mapping = make(map[string][2]int)
	var mapping2 = make(map[[2]int][][3]byte)

	for y, v := range buffer {
		var buf = strings.Split(v, "\t")
		for x, cell := range buf {

			mapping[cell] = [2]int{x, y}
		}
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		var iy = (y - b.Min.Y) / f.celly
		for x := b.Min.X; x < b.Max.X; x++ {
			var ix = (x - b.Min.X) / f.cellx
			var i = [2]int{ix, iy}

			var sli = mapping2[i]

			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

			sli = append(sli, [3]byte{c.R, c.G, c.B})

			mapping2[i] = sli
		}
	}
	if f.mapping == nil {
		f.mapping = make(map[string][][3]byte)
	}
	for k, v := range mapping {
		f.mapping[k] = mapping2[v]
	}
	return nil
}
