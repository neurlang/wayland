package main

import "sort"
import "fmt"

type Canvas interface {
	PutRGB(ObjectPosition, [][3]byte, int, [3]byte, [3]byte, bool)
	GetTime() uint32
}

type ObjectPosition struct {
	X, Y int
}

func (o *ObjectPosition) Less(p ObjectPosition) bool {
	if o.Y < p.Y {
		return true
	}
	if o.Y == p.Y && o.X < p.X {
		return true
	}
	return false
}

type StringCell struct {
	Pos          ObjectPosition
	String       string
	CellWidth    int
	CellHeight   int
	Font         *Font
	BgRGB, FgRGB [3]byte
	Flip         bool
}

func (sc *StringCell) Render(c Canvas) {
	c.PutRGB(sc.Pos, sc.Font.GetRGBTexture(sc.String), sc.CellWidth, sc.BgRGB, sc.FgRGB, sc.Flip)
}

type StringGrid struct {
	Pos                   ObjectPosition
	LineNumbers           int
	LastColHint           int
	XCells                int
	YCells                int
	Content               []string
	CellWidth             int
	CellHeight            int
	Font                  *Font
	IbeamCursor           ObjectPosition
	IbeamCursorBlinkFix   uint32
	SelectionCursor       ObjectPosition
	Hover                 ObjectPosition
	Selecting, IsSelected bool
	ContentFgColor        map[[2]int][3]byte
	lineLen               []int
	FlipColor             bool
}

func (sg *StringGrid) Button(up bool) {
	if up {
		sg.IsSelected = sg.Selecting
		sg.Selecting = false
	} else {
		sg.Selecting = true
		sg.IsSelected = false
		sg.SelectionCursor = sg.Hover
		sg.IbeamCursor = sg.Hover
	}
}
func (sg *StringGrid) Motion(pos ObjectPosition) {

	pos.X -= sg.LineNumbers

	for pos.X > 0 && (sg.XCells*pos.Y+pos.X-1) < len(sg.Content) && sg.Content[sg.XCells*pos.Y+pos.X-1] == "" {
		pos.X--
	}

	if pos.X < 0 {
		pos.X = 0
	}

	sg.Hover = pos

	if sg.Selecting {
		sg.IbeamCursor = sg.Hover

	}
}

func (sg *StringGrid) FgColor(x, y int) [3]byte {
	for i := x; i >= 0 && i > x-16; i-- {
		if sg.ContentFgColor != nil {
			if c, ok := sg.ContentFgColor[[2]int{i, y}]; ok {
				return c
			}
		}
	}
	return [3]byte{255, 255, 255}
}

func (sg *StringGrid) Width() int {
	return sg.XCells * sg.CellWidth
}

func (sg *StringGrid) Height() int {
	return sg.YCells * sg.CellHeight
}

func (sg *StringGrid) IsSelection() bool {
	if !(sg.Selecting || sg.IsSelected) {
		return false
	}
	return true
}
func (sg *StringGrid) Selected(x, y int) bool {
	if !sg.IsSelection() {
		return false
	}
	var objs = [3]ObjectPosition{sg.SelectionCursor, sg.IbeamCursor, {x, y}}
	sort.Slice(objs[:], func(i, j int) bool {
		return objs[i].Less(objs[j])
	})
	return objs[1] == ObjectPosition{x, y} && objs[1] != objs[2]
}

func (sg *StringGrid) RowFocused(y int) bool {
	return sg.IbeamCursor.Y == y
}

func (sg *StringGrid) Render(c Canvas) {
	for y := 0; y < sg.YCells; y++ {
		var linenum = fmt.Sprintf("% "+fmt.Sprint(sg.LineNumbers-1)+"d   ", y+1)
		for x := 0; x < sg.LineNumbers; x++ {

			var bgcolor = [3]byte{0, 13, 26}
			var fgcolor = [3]byte{0, 101, 191}

			var cell = &StringCell{
				Pos: ObjectPosition{
					sg.Pos.X + x*sg.CellWidth,
					sg.Pos.Y + y*sg.CellHeight,
				},
				String:     string(linenum[x]),
				CellWidth:  sg.CellWidth,
				CellHeight: sg.CellHeight,
				Font:       sg.Font,
				BgRGB:      bgcolor,
				FgRGB:      fgcolor,
				Flip:       sg.FlipColor,
			}
			cell.Render(c)
		}

		for x := sg.LineNumbers; x < sg.XCells; x++ {

			xx := x - sg.LineNumbers

			var selected = sg.Selected(xx, y)
			var bgcolor = [3]byte{0, 27, 51}
			var fgcolor = sg.FgColor(xx, y)
			if selected {
				bgcolor = [3]byte{0, 136, 255}
				fgcolor = [3]byte{255, 255, 255}
			} else if sg.RowFocused(y) {
				if x > sg.LastColHint {
					bgcolor = [3]byte{12, 68, 117}
				} else {
					bgcolor = [3]byte{0, 59, 112}
				}
			} else if x > sg.LastColHint {
				bgcolor = [3]byte{12, 37, 60}
			}
			fgcolor = maxColor(fgcolor, bgcolor)

			var cell = &StringCell{
				Pos: ObjectPosition{
					sg.Pos.X + x*sg.CellWidth,
					sg.Pos.Y + y*sg.CellHeight,
				},
				String:     sg.Content[sg.XCells*y+xx],
				CellWidth:  sg.CellWidth,
				CellHeight: sg.CellHeight,
				Font:       sg.Font,
				BgRGB:      bgcolor,
				FgRGB:      fgcolor,
				Flip:       sg.FlipColor,
			}
			cell.Render(c)
		}
	}

	if (c.GetTime()-uint32(sg.IbeamCursorBlinkFix))&512 == 0 {
		var cursor = &IbeamCursor{
			Pos: ObjectPosition{
				sg.Pos.X + (sg.IbeamCursor.X+sg.LineNumbers)*sg.CellWidth,
				sg.Pos.Y + sg.IbeamCursor.Y*sg.CellHeight,
			},
			CellHeight: sg.CellHeight,
			RGB:        [3]byte{127, 127, 127},
		}
		if cursor.Pos.X < 0 {
			cursor.Pos.X = 0
		}
		if cursor.Pos.Y < 0 {
			cursor.Pos.Y = 0
		}
		cursor.Render(c)
	}
}

type VisualScrollBar struct {
	Pos     ObjectPosition
	XCells  int
	YCells  int
	Content string
}

const VisualScrollBarCellWidth = 2
const VisualScrollBarCellHeight = 3

func (vsb *VisualScrollBar) Width() int {
	return vsb.XCells * VisualScrollBarCellWidth
}

func (vsb *VisualScrollBar) Height() int {
	return vsb.YCells * VisualScrollBarCellHeight
}

type IbeamCursor struct {
	Pos        ObjectPosition
	CellHeight int
	RGB        [3]byte
}

func (ic *IbeamCursor) Render(c Canvas) {
	var buf = make([][3]byte, ic.CellHeight*2)
	for i := range buf {
		buf[i] = ic.RGB
	}
	c.PutRGB(ic.Pos, buf, 2, [3]byte{0, 0, 0}, [3]byte{255, 255, 255}, false)
}
