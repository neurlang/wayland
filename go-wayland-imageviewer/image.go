package main

import (
	"bufio"
	"image"
	"image/draw"
	"os"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func rgbaImageFromFile(filePath string) (*image.RGBA, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	img, _, err := image.Decode(br)
	if err != nil {
		return nil, err
	}

	rgbaImage, ok := img.(*image.RGBA)
	if !ok {
		// Convert to RGBA if not already RGBA
		rect := img.Bounds()
		rgbaImage = image.NewRGBA(rect)
		draw.Draw(rgbaImage, rect, img, rect.Min, draw.Src)
	}

	return rgbaImage, nil
}
