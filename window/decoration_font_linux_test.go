//go:build linux

package window

import (
	"testing"

	"github.com/gogpu/gg"
)

// A system with none of decorationFonts still gets a title: the built-in font.
func TestDecorationFallbackFaceMeasuresText(t *testing.T) {
	face := decorationFallbackFace(12)
	if face == nil {
		t.Fatal("no fallback face")
	}
	if w := face.Advance("f4 - Panels"); w <= 0 {
		t.Fatalf("the fallback face measured a title as %v wide", w)
	}
}

func TestLoadDecorationFontSetsAFace(t *testing.T) {
	dc := gg.NewContext(200, 24)
	loadDecorationFont(dc)
	if w, _ := dc.MeasureString("Title"); w <= 0 {
		t.Fatalf("no title font after loadDecorationFont: width %v", w)
	}
}
