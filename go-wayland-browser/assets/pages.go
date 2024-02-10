package assets

import _ "embed"

//go:embed pages/defaultPage.html
var defaultPage []byte
func DefaultPage() []byte {
	return defaultPage
}

//go:embed pages/homePage.html
var homePage []byte
func HomePage() []byte {
	return homePage
}
