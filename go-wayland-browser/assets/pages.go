package assets

import _ "embed"

func DefaultPage() []byte {
	return defaultPage
}

func HomePage() []byte {
	return homePage
}

//go:embed pages/defaultPage.html
var defaultPage []byte


/*
var DefaultPage = func() []byte {
	return defaultPage
}
*/
//go:embed pages/homePage.html
var homePage []byte


/*
var HomePage = func() []byte {
	return homePage
}
*/
