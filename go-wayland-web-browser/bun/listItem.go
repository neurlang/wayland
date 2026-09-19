package bun

import (
	gg "github.com/gogpu/gg"
	text "github.com/gogpu/gg/text"
	hotdog "github.com/neurlang/wayland/go-wayland-web-browser/hotdog"

	"strings"
)

func paintListItemElement(ctx *gg.Context, node *hotdog.NodeDOM) {
	ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.DrawCircle(node.RenderBox.Left-15, node.RenderBox.Top+node.Style.FontSize/2, 3)
	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.SetFont(SansSerif.Face(node.Style.FontSize, text.WithVariations(
		text.NewFontVariation("wght", float32(node.Style.FontWeight)),
	)))
	ctx.DrawStringWrapped(node.Content, node.RenderBox.Left, node.RenderBox.Top+1, 0, 0, node.RenderBox.Width, 1, gg.AlignLeft)
	ctx.Fill()
}

func calculateListItemLayout(ctx *gg.Context, node *hotdog.NodeDOM, childIdx int) {
	if node.Style.Width == 0 {
		node.RenderBox.Width = node.Parent.RenderBox.Width - 30
	}

	if node.Style.Height == 0 && len(node.Content) > 0 {
		ctx.SetFont(SansSerif.Face(node.Style.FontSize, text.WithVariations(
			text.NewFontVariation("wght", float32(node.Style.FontWeight)),
		)))
		_, h := ctx.MeasureMultilineString(strings.Repeat("\n", len(ctx.WordWrap(node.Content, node.Style.Width))-1), 1)
		node.RenderBox.Height = h + 2
	}

	if childIdx > 0 {
		prev := node.Parent.Children[childIdx-1]
		node.RenderBox.Top = prev.RenderBox.Top + prev.RenderBox.Height
	} else {
		node.RenderBox.Top = node.Parent.RenderBox.Top
	}

	node.RenderBox.Left = 30
}
