package mustard

import (
	"github.com/gogpu/gg"
	cairo "github.com/neurlang/wayland/cairoshim"
	assets "github.com/neurlang/wayland/go-wayland-web-browser/assets"
	bun "github.com/neurlang/wayland/go-wayland-web-browser/bun"
)

// CreateTreeWidget - Creates and returns a new Tree Widget
func CreateTreeWidget() *TreeWidget {
	var widgets []Widget
	font := bun.SansSerif

	//openIcon, _ := gg.LoadAsset(assets.DownChevron())
	//closeIcon, _ := gg.LoadAsset(assets.RightChevron())

	return &TreeWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			widgetType: treeWidget,

			cursor: ArrowCursor,

			backgroundColor: "#fff",

			font: font,
		},

		openIcon:  assets.RightChevron(),
		closeIcon: assets.Logo(),
		fontSize:  20,
		fontColor: "#ff00ff00",
	}
}

func (tree *TreeWidget) Click() {

	x, y := tree.window.GetCursorPosition()
	node := getIntersectedNode(tree.nodes, x, y)

	if node != nil {
		fireNodeEvents(tree, node, x, y, tree.fontSize)
	}

}

func (tree *TreeWidget) SetSelectCallback(selectCallback func(*TreeWidgetNode)) {
	tree.selectCallback = selectCallback
}

func desselectNodes(nodes []*TreeWidgetNode) {
	for _, node := range nodes {
		node.isSelected = false
		desselectNodes(node.Children)
	}
}

func (tree *TreeWidget) SelectNode(node *TreeWidgetNode) {
	for _, treeNode := range tree.nodes {
		treeNode.isSelected = false
		desselectNodes(treeNode.Children)
	}

	node.isSelected = true
}

func (tree *TreeWidget) SelectNodeByValue(value string) {
	if tree != nil {
		selectNodeByValue(tree.nodes, value)
	}
}

func selectNodeByValue(nodes []*TreeWidgetNode, value string) {
	for _, childNode := range nodes {
		selectNodeByValue(childNode.Children, value)

		if childNode.Value == value {
			childNode.isSelected = true
		} else {
			childNode.isSelected = false
		}
	}
}

func selectNode(nodes []*TreeWidgetNode, node *TreeWidgetNode) { //nolint:unused // Reserved for future use
	for _, childNode := range nodes {
		selectNode(childNode.Children, node)

		if childNode == node {
			node.isSelected = true
		} else {
			node.isSelected = false
		}
	}
}

func fireNodeEvents(tree *TreeWidget, node *TreeWidgetNode, x, y, nodeHeight float64) {
	t, l, _, _ := node.box.GetCoords()
	if y > t && y < t+nodeHeight {
		tree.SelectNode(node)
		tree.selectCallback(node)

		if x < l+25 {
			node.Toggle()
		}
	}
}

func getIntersectedNode(nodes []*TreeWidgetNode, x, y float64) *TreeWidgetNode {
	var intersectedNode *TreeWidgetNode
	for _, node := range nodes {
		if x > float64(node.box.left) &&
			x < float64(node.box.left+node.box.width) &&
			y > float64(node.box.top) &&
			y < float64(node.box.top+node.box.height) {

			if !node.isOpen {
				return node
			}

			intersectedNode = node
			childIntersectedNode := getIntersectedNode(node.Children, x, y)
			if childIntersectedNode != nil {
				intersectedNode = childIntersectedNode
			}
		}
	}

	return intersectedNode
}

// SetWidth - Sets the tree width
func (tree *TreeWidget) SetWidth(width float64) {
	tree.box.width = width
	tree.fixedWidth = true
	tree.RequestReflow()
}

// SetHeight - Sets the tree height
func (tree *TreeWidget) SetHeight(height float64) {
	tree.box.height = height
	tree.fixedHeight = true
	tree.RequestReflow()
}

// SetFontSize - Sets the tree font size
func (tree *TreeWidget) SetFontSize(fontSize float64) {
	tree.fontSize = fontSize
	tree.needsRepaint = true
}

// SetFontColor - Sets the tree font color
func (tree *TreeWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		tree.fontColor = fontColor
		tree.needsRepaint = true
	}
}

// SetBackgroundColor - Sets the tree background color
func (tree *TreeWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		tree.backgroundColor = backgroundColor
		tree.needsRepaint = true
	}
}

func (tree *TreeWidget) render(s cairo.Surface, time uint32) {
	pm := gg.NewPixmapFromBuffer(s.ImageSurfaceGetData(), s.ImageSurfaceGetWidth(), s.ImageSurfaceGetHeight())
	context := gg.NewContext(s.ImageSurfaceGetWidth(), s.ImageSurfaceGetHeight(), gg.WithPixmap(pm))
	top, left, width, height := tree.computedBox.GetCoords()

	context.SetHexColor(tree.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	for _, node := range tree.nodes {
		flowNode(context, node, tree, 0)
		drawNode(context, node, tree, 0)
	}

}

func flowNode(context *gg.Context, node *TreeWidgetNode, tree *TreeWidget, level int) {
	node.box.left = tree.fontSize * float64(level)
	node.box.height = tree.fontSize + 4
	node.box.width = tree.computedBox.width

	prevSibling := node.PreviousSibling()
	if node.Parent == nil {
		if prevSibling != nil {
			node.box.top = prevSibling.box.top + prevSibling.box.height
		} else {
			node.box.top = 0
		}
	} else {
		if prevSibling != nil {
			node.box.top = prevSibling.box.top + prevSibling.box.height
		} else {
			node.box.top = node.Parent.box.top + node.Parent.box.height
		}
	}

	if node.isOpen {
		for _, childNode := range node.Children {
			flowNode(context, childNode, tree, level+1)
			node.box.height += childNode.box.height
		}
	}

}

func drawNode(context *gg.Context, node *TreeWidgetNode, tree *TreeWidget, level int) {
	top, left, width, height := node.box.GetCoords()

	if node.isSelected {
		context.SetHexColor("#7db1ff32")
		context.DrawRectangle(float64(tree.computedBox.left), float64(top), float64(width), float64(height))
		context.Fill()

	} else {
		context.SetHexColor(tree.backgroundColor)
		context.DrawRectangle(float64(tree.computedBox.left), float64(top), float64(width), float64(height))
		context.Fill()
	}

	context.SetHexColor(tree.fontColor)
	context.SetFont(tree.font.Face(tree.fontSize))
	context.SetColor(gg.Black)
	context.DrawString(node.Key, float64(left)+20, float64(top))
	context.Fill()

	if len(node.Children) > 0 {
		if node.isOpen {
			context.DrawImage(gg.ImageBufFromImage(tree.openIcon), left+4, top+1)

			for _, childNode := range node.Children {
				drawNode(context, childNode, tree, level+1)
			}
		} else {
			context.Push()
			context.Rotate(40)
			context.DrawImage(gg.ImageBufFromImage(tree.closeIcon), left+4, top+1)
			context.Pop()
		}
	}

}
