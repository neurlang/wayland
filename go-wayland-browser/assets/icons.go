package assets

import _ "embed"

//go:embed icons/arrowLeft.png
var arrowLeft []byte
func ArrowLeft() []byte {
	return arrowLeft
}

//go:embed icons/arrowRight.png
var arrowRight []byte
func ArrowRight() []byte {
	return arrowRight
}

//go:embed icons/downChevron.png
var downChevron []byte
func DownChevron() []byte {
	return downChevron
}

//go:embed icons/errorImage.png
var errorImage []byte
func ErrorImage() []byte {
	return errorImage
}

//go:embed icons/menu.png
var menu []byte
func Menu() []byte {
	return menu
}

//go:embed icons/reload.png
var reload []byte
func Reload() []byte {
	return reload
}

//go:embed icons/rightChevron.png
var rightChevron []byte
func RightChevron() []byte {
	return rightChevron
}

//go:embed icons/tools.png
var tools []byte
func Tools() []byte {
	return tools
}

