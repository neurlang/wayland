package mustard

import (
	gg "github.com/gogpu/gg"
	"image"
)

// CreateImageWidget - Creates and returns a new Image Widget
func CreateImageWidget(data image.Image) *ImageWidget {
	var widgets []Widget

	return &ImageWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			cursor: ArrowCursor,

			widgetType: imageWidget,

			backgroundColor: "#fff",
		},

		//path: path,
		img: data,
	}
}

// SetWidth - Sets the label width
func (label *ImageWidget) SetWidth(width float64) {
	label.box.width = width
	label.fixedWidth = true
	label.RequestReflow()
}

// SetHeight - Sets the label height
func (label *ImageWidget) SetHeight(height float64) {
	label.box.height = height
	label.fixedHeight = true
	label.RequestReflow()
}

func (im *ImageWidget) render(s Surface, time uint32) {
	context := makeContextFromCairo(s)

	top, left, _, _ := im.computedBox.GetCoords()
	context.DrawImage(gg.ImageBufFromImage(im.img), left+15, top+3)

}
