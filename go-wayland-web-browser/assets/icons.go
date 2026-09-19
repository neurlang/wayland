package assets

import "bytes"
import . "image"
import "image/png"
import _ "embed"

func decode(r []byte) (img Image) {
	img, _ = png.Decode(bytes.NewReader(r))
	return
}

//go:embed icons/arrowLeft.png
var arrowLeft []byte

func ArrowLeft() Image {
	return decode(arrowLeft)
}

//go:embed icons/arrowRight.png
var arrowRight []byte

func ArrowRight() Image {
	return decode(arrowRight)
}

//go:embed icons/downChevron.png
var downChevron []byte

func DownChevron() Image {
	return decode(downChevron)
}

//go:embed icons/errorImage.png
var errorImage []byte

func ErrorImage() Image {
	return decode(errorImage)
}

//go:embed icons/menu.png
var menu []byte

func Menu() Image {
	return decode(menu)
}

//go:embed icons/reload.png
var reload []byte

func Reload() Image {
	return decode(reload)
}

//go:embed icons/rightChevron.png
var rightChevron []byte

func RightChevron() Image {
	return decode(rightChevron)
}

//go:embed icons/tools.png
var tools []byte

func Tools() Image {
	return decode(tools)
}

//go:embed icons/gopher.png
var logo []byte

func Logo() Image {
	return decode(logo)
}
