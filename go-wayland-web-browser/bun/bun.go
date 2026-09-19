package bun

import (
	gg "github.com/gogpu/gg"
	hotdog "github.com/neurlang/wayland/go-wayland-web-browser/hotdog"
)

func RenderDocument(ctx *gg.Context, document *hotdog.Document, experimentalLayout bool) error {
	if !experimentalLayout {
		body, _ := document.DOM.FindChildByName("body")

		document.DOM.RenderBox.Width = float64(ctx.Width())
		document.DOM.RenderBox.Height = float64(ctx.Height())

		ctx.ClearWithColor(gg.RGBA{body.Style.BackgroundColor.R, body.Style.BackgroundColor.G, body.Style.BackgroundColor.B, 1})

		layoutDOM(ctx, body, 0)
	} else {
		html, err := document.DOM.FindChildByName("html")
		if err != nil {
			return err
		}

		renderTree := createRenderTree(html)
		renderTree.RenderBox.Width = float64(ctx.Width())
		renderTree.RenderBox.Height = float64(ctx.Height())

		layoutNode(ctx, renderTree)
		paintNode(ctx, renderTree)
		paintText(ctx, renderTree)

		renderTree.Print(0)
	}

	return nil
}

func getNodeContent(NodeDOM *hotdog.NodeDOM) string {
	return NodeDOM.Content
}

func getElementName(NodeDOM *hotdog.NodeDOM) string {
	return NodeDOM.Element
}

func getNodeChildren(NodeDOM *hotdog.NodeDOM) []*hotdog.NodeDOM {
	return NodeDOM.Children
}

func layoutDOM(ctx *gg.Context, node *hotdog.NodeDOM, childIdx int) {
	nodeChildren := getNodeChildren(node)

	node.RenderBox = &hotdog.RenderBox{}
	calculateNode(ctx, node, childIdx)

	for i := 0; i < len(nodeChildren); i++ {
		layoutDOM(ctx, nodeChildren[i], i)
		node.RenderBox.Height += nodeChildren[i].RenderBox.Height
	}

	paintNode(ctx, node)
}

func paintNode(ctx *gg.Context, node *hotdog.NodeDOM) {
	switch node.Style.Display {
	case "block":
		paintBlockElement(ctx, node)
	case "inline":
		paintInlineElement(ctx, node)
	case "list-item":
		paintListItemElement(ctx, node)
	}
}

func calculateNode(ctx *gg.Context, node *hotdog.NodeDOM, postion int) {
	switch node.Style.Display {
	case "block":
		calculateBlockLayout(ctx, node, postion)
	case "inline":
		calculateInlineLayout(ctx, node, postion)
	case "list-item":
		calculateListItemLayout(ctx, node, postion)
	}
}

func GetPageTitle(TreeDOM *hotdog.NodeDOM) string {
	const unnamed = "Unnamed"
	pageTitle := unnamed

	if getElementName(TreeDOM) == "title" {
		return getNodeContent(TreeDOM)
	}

	nodeChildren := getNodeChildren(TreeDOM)

	for i := 0; i < len(nodeChildren); i++ {

		nPageTitle := GetPageTitle(nodeChildren[i])

		if nPageTitle != unnamed {
			pageTitle = nPageTitle
			break
		}
	}

	return pageTitle
}
