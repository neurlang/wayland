package assets

import _ "embed"

//go:embed fonts/OpenSans-Light.ttf
var openSansLight []byte
//go:embed fonts/OpenSans-Regular.ttf
var openSansRegular []byte
//go:embed fonts/OpenSans-SemiBold.ttf
var openSansSemiBold []byte
//go:embed fonts/OpenSans-Bold.ttf
var openSansBold []byte
//go:embed fonts/OpenSans-ExtraBold.ttf
var openSansExtraBold []byte

func OpenSans(weight int) []byte {
	switch weight {
	case 300:
		// Light
		return openSansLight
	default:
		// Regular
		return openSansRegular
	case 600:
		// SemiBold
		return openSansSemiBold
	case 700:
		// Bold
		return openSansBold
	case 800:
		// ExtraBold
		return openSansExtraBold
	}
}
