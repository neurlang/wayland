package main

import "sync"
import "sort"
import "fmt"
import "unicode"

const rowHeight = 2
const whiteRectThick = 2
const scrollingOffset = 100

type Canvas interface {
	PutRGB(ObjectPosition, [][3]byte, int, [3]byte, [3]byte, bool)
	GetTime() uint32
}

type ObjectPosition struct {
	X, Y int
}

func (o *ObjectPosition) Less(p *ObjectPosition) bool {
	if o.Y < p.Y {
		return true
	}
	if o.Y == p.Y && o.X < p.X {
		return true
	}
	return false
}
func (o *ObjectPosition) Lesser(p *ObjectPosition) *ObjectPosition {
	if o.Less(p) {
		return o
	}
	return p
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

func (sg *StringGrid) IbeamCursorAbsolute() *ObjectPosition {

	if sg.IbeamCursorAbs.X != sg.FilePosition.X+sg.IbeamCursor.X {
		panic("x mismatch")
	}
	if sg.IbeamCursorAbs.Y != sg.FilePosition.Y+sg.IbeamCursor.Y {
		panic("y mismatch")
	}

	return &sg.IbeamCursorAbs
}
func (sg *StringGrid) SelectionCursorAbsolute() *ObjectPosition {
	if sg.SelectionCursorAbs.X != sg.FilePosition.X+sg.SelectionCursor.X {
		panic("x mismatch")
	}
	if sg.SelectionCursorAbs.Y != sg.FilePosition.Y+sg.SelectionCursor.Y {
		panic("y mismatch")
	}

	return &sg.SelectionCursorAbs
}

type StringGrid struct {
	Pos                   ObjectPosition
	LineCount             int
	EndLineLen            int
	LineNumbers           int
	LastColHint           int
	XCells                int
	YCells                int
	Content               []string
	CellWidth             int
	CellHeight            int
	Font                  *Font
	FilePosition          ObjectPosition
	IbeamCursor           ObjectPosition
	IbeamCursorAbs        ObjectPosition
	IbeamCursorBlinkFix   uint32
	SelectionCursor       ObjectPosition
	SelectionCursorAbs    ObjectPosition
	Hover                 ObjectPosition
	HoverOld              ObjectPosition
	Selecting, IsSelected bool
	ContentFgColor        map[[2]int][3]byte
	lineLen               []int
	BgColor               [3]byte
	FgColor               [3]byte
	FlipColor             bool
	LineLens              []int
	wasDoubleClick        bool
	MatchingBrace         ObjectPosition
	control               bool
}

func (sg *StringGrid) Control(ctrl bool) {
	sg.control = ctrl
}
func (sg *StringGrid) LineLen(y int) int {
	if y >= len(sg.LineLens) {
		return 0
	}
	return sg.LineLens[y]
}

func (sg *StringGrid) DoLineNumbers() {
	var maxLn = sg.YCells + sg.FilePosition.Y

	if maxLn > sg.LineCount {
		maxLn = sg.LineCount
	}

	println("DoLineNumbers", maxLn)

	for sg.LineNumbers = 2; maxLn > 0; sg.LineNumbers++ {
		maxLn /= 10
	}
}

func (sg *StringGrid) IsHoverButton() bool {
	return sg.Hover.X >= 0 && sg.Hover.Y >= 0
}

func (sg *StringGrid) IsHover(x, y float32, w, h int32) bool {
	var pos = sg.Pos
	if pos.X < 0 {
		pos.X += int(w)
	}
	if pos.Y < 0 {
		pos.Y += int(h)
	}
	if x < float32(pos.X)+float32(sg.CellWidth*sg.LineNumbers) {
		return false
	}
	if y < float32(pos.Y) {
		return false
	}
	if x-float32(pos.X) > float32(sg.CellWidth*sg.XCells) {
		return false
	}
	if y-float32(pos.Y) > float32(sg.CellHeight*sg.YCells) {
		return false
	}
	return true
}
func (sg *StringGrid) isTripleClick() bool {
	return !sg.Selecting && sg.IsSelected && sg.SelectionCursor != sg.Hover && sg.IbeamCursor == sg.Hover

}

func (sg *StringGrid) isDoubleClick() bool {
	return !sg.Selecting && sg.IsSelected && sg.SelectionCursor == sg.Hover && sg.IbeamCursor == sg.Hover
}

func (sg *StringGrid) WasDoubleClick() bool {
	return sg.wasDoubleClick
}

func (sg *StringGrid) Button(up bool) (changedSelection bool) {
	if up {
		changedSelection = false

		sg.IsSelected = sg.Selecting
		sg.Selecting = false

	} else if sg.isDoubleClick() && !sg.wasDoubleClick {
		changedSelection = true

		sg.wasDoubleClick = true
		sg.Selecting = true
		sg.IsSelected = false

		sg.ReMotion()

		var add = 0
		for ; isLetterOrDigit(sg.GetContentRune(sg.Hover.X+add, sg.Hover.Y)); add++ {
			if sg.Hover.Y < len(sg.LineLens) && sg.LineLens[sg.Hover.Y] <= sg.Hover.X+add {
				break
			}
		}

		sg.SelectionCursor = sg.Hover
		sg.IbeamCursor = sg.Hover

		sg.SelectionCursor.X += add

		sg.SelectionCursorAbs.X = sg.Hover.X + sg.FilePosition.X + add
		sg.SelectionCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y
		sg.IbeamCursorAbs.X = sg.Hover.X + sg.FilePosition.X
		sg.IbeamCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y

	} else if sg.isTripleClick() && sg.wasDoubleClick {
		changedSelection = true

		sg.wasDoubleClick = false
		sg.Selecting = true
		sg.IsSelected = false

		sg.ReMotion()

		sg.Hover.X = 0

		var add = 0
		for ; true; add++ {
			if sg.Hover.Y < len(sg.LineLens) && sg.LineLens[sg.Hover.Y] <= sg.Hover.X+add {
				break
			}
		}

		sg.SelectionCursor = sg.Hover
		sg.IbeamCursor = sg.Hover

		sg.SelectionCursor.X += add

		sg.SelectionCursorAbs.X = sg.Hover.X + sg.FilePosition.X + add
		sg.SelectionCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y
		sg.IbeamCursorAbs.X = sg.Hover.X + sg.FilePosition.X
		sg.IbeamCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y

	} else {
		changedSelection = false

		sg.Selecting = true
		sg.IsSelected = false
		sg.wasDoubleClick = false

		sg.ReMotion()

		sg.SelectionCursor = sg.Hover
		sg.IbeamCursor = sg.Hover
		sg.SelectionCursorAbs.X = sg.Hover.X + sg.FilePosition.X
		sg.SelectionCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y
		sg.IbeamCursorAbs.X = sg.Hover.X + sg.FilePosition.X
		sg.IbeamCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y
	}
	return
}
func (sg *StringGrid) ReMotion() {
	sg.Motion(sg.HoverOld)
}

func (sg *StringGrid) lookupBraceWhenMotion(pos ObjectPosition) {
	var color1 = sg.GetFgColor(pos.X, pos.Y)
	switch sg.GetContentRune(pos.X, pos.Y) {
	case '{', '[', '(':
		for ; pos.Y < len(sg.LineLens); pos.Y++ {
			for ; pos.X < sg.LineLens[pos.Y]; pos.X++ {
				switch sg.GetContentRune(pos.X, pos.Y) {
				case '}', ']', ')':
					if color1 == sg.GetFgColor(pos.X, pos.Y) {
						sg.MatchingBrace.X = pos.X
						sg.MatchingBrace.Y = pos.Y
						return
					}
				}
			}
			if pos.X == sg.LineLens[pos.Y] {
				pos.X = 0
			}
		}
	case '}', ']', ')':
		for ; pos.Y >= 0; pos.Y-- {
			if pos.X < 0 {
				pos.X = sg.LineLens[pos.Y]
			}
			for ; pos.X >= 0; pos.X-- {
				switch sg.GetContentRune(pos.X, pos.Y) {
				case '{', '[', '(':
					if color1 == sg.GetFgColor(pos.X, pos.Y) {
						sg.MatchingBrace.X = pos.X
						sg.MatchingBrace.Y = pos.Y
						return
					}
				}
			}
		}
	}

}

func isLetterOrDigit(l rune) bool {
	return unicode.IsLetter(l) || unicode.IsDigit(l)
}

func (sg *StringGrid) Motion(pos ObjectPosition) {

	sg.HoverOld = pos

	pos.X -= sg.LineNumbers
	if pos.Y < 0 {
		pos.Y = 0
	}
	if pos.X < 0 {
		pos.X = 0
	}
	if pos.Y >= sg.LineCount-sg.FilePosition.Y {
		pos.Y = sg.LineCount - 1 - sg.FilePosition.Y
		if pos.Y < 0 {
			pos.Y = 0
		}
	}
	if pos.X > 0 && pos.Y >= 0 && pos.Y < len(sg.LineLens) && sg.LineLens[pos.Y] < pos.X {
		pos.X = sg.LineLens[pos.Y]
	}
	for (pos.X > 0) && (len(sg.GetContent(pos.X-1, pos.Y)) == 0) {
		if sg.GetContent(pos.X, pos.Y) != "\t" {

			pos.X--

		} else {
			pos.X++
			break
		}
	}

	if sg.Selecting && sg.wasDoubleClick {
		if (pos.X > 0) && (pos.Y < len(sg.LineLens)) && (sg.LineLens[pos.Y] == pos.X) {
			pos.X--
		}

		for (pos.X > 0) && (isLetterOrDigit(sg.GetContentRune(pos.X, pos.Y))) {
			pos.X--
		}
		if (pos.X > 0) && (!isLetterOrDigit(sg.GetContentRune(pos.X, pos.Y))) {
			pos.X++
		}
		if pos.Y < len(sg.LineLens) && sg.LineLens[pos.Y] < pos.X {
			if !isLetterOrDigit(sg.GetContentRune(pos.X, pos.Y)) {
				pos.X++
			}
		}
	}

	sg.Hover = pos

	sg.MatchingBrace = pos

	sg.lookupBraceWhenMotion(pos)
	pos.X--
	sg.lookupBraceWhenMotion(pos)

	if sg.Selecting {
		sg.IbeamCursor = sg.Hover
		sg.IbeamCursorAbs.X = sg.Hover.X + sg.FilePosition.X
		sg.IbeamCursorAbs.Y = sg.Hover.Y + sg.FilePosition.Y
	}
}

func (sg *StringGrid) GetFgColor(x, y int) [3]byte {
	for i := x; i >= 0 && i > x-17; i-- {
		if sg.ContentFgColor != nil {
			if c, ok := sg.ContentFgColor[[2]int{i, y}]; ok {
				return c
			}
		}
	}
	return sg.FgColor
}

func (sg *StringGrid) Width() int {
	return sg.XCells * sg.CellWidth
}

func (sg *StringGrid) Height() int {
	return sg.YCells * sg.CellHeight
}

func (sg *StringGrid) IsSelectionStrict() bool {
	return sg.SelectionCursor != sg.IbeamCursor
}
func (sg *StringGrid) IsSelection() bool {
	if !(sg.Selecting || sg.IsSelected) {
		return false
	}
	return true
}
func (sg *StringGrid) Highlighted(x, y int) bool {
	if sg.Selecting && !sg.control &&
		((x == sg.Hover.X || x == sg.Hover.X-1) && y == sg.Hover.Y ||
			x == sg.MatchingBrace.X && y == sg.MatchingBrace.Y) {
		switch sg.GetContentRune(x, y) {
		case '}', ')', ']', '{', '(', '[':
			return true
		}
		return false
	}
	if sg.control {
		switch sg.GetContentRune(x, y) {
		case '}', ')', ']':
			var fg = sg.GetFgColor(x, y)
			return fg == sg.GetFgColor(sg.Hover.X, sg.Hover.Y) ||
				fg == sg.GetFgColor(sg.Hover.X-1, sg.Hover.Y)
		case '{', '(', '[':
			var fg = sg.GetFgColor(x, y)
			return fg == sg.GetFgColor(sg.Hover.X, sg.Hover.Y) ||
				fg == sg.GetFgColor(sg.Hover.X-1, sg.Hover.Y)
		}
	}
	return false
}
func (sg *StringGrid) Selected(x, y int) bool {
	if !sg.IsSelection() {
		return false
	}
	var objs = [3]ObjectPosition{sg.SelectionCursor, sg.IbeamCursor, {x, y}}
	sort.Slice(objs[:], func(i, j int) bool {
		return objs[i].Less(&objs[j])
	})
	return objs[1] == ObjectPosition{x, y} && objs[1] != objs[2]
}

func (sg *StringGrid) RowFocused(y int) bool {
	return sg.IbeamCursorAbs.Y == y+sg.FilePosition.Y
}
func (sg *StringGrid) GetContentRune(x, y int) rune {
	arr := []rune(sg.GetContent(x, y))
	if len(arr) == 0 {
		return 0
	}
	return arr[0]
}
func (sg *StringGrid) GetContent(x, y int) string {
	var pos = sg.XCells*y + x
	if len(sg.Content) <= pos {
		return ""
	}
	if pos < 0 {
		return ""
	}
	return sg.Content[pos]
}

func (sg *StringGrid) Render(c Canvas) {
	for y := 0; y < sg.YCells; y++ {
		var linenum = fmt.Sprintf("% "+fmt.Sprint(sg.LineNumbers-1)+"d   ", y+sg.FilePosition.Y+1)
		if y+sg.FilePosition.Y >= sg.LineCount {
			linenum = "                      "
		}
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
			var highlighted = sg.Highlighted(xx, y)
			var bgcolor = [3]byte{0, 27, 51}
			var fgcolor = sg.GetFgColor(xx, y)
			if selected {
				bgcolor = [3]byte{0, 136, 255}
				fgcolor = sg.FgColor
			} else if highlighted {
				bgcolor = [3]byte{0, 255, 128}
				fgcolor = sg.FgColor
			} else if sg.RowFocused(y) {
				if x > sg.LastColHint {
					bgcolor = [3]byte{12, 68, 117}
				} else {
					bgcolor = sg.BgColor
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
				String:     sg.GetContent(xx, y),
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
		if cursor.Pos.X >= 0 && cursor.Pos.Y >= 0 {
			cursor.Render(c)
		}
	}
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

type Scrollbar struct {
	Pos     ObjectPosition
	Width   int
	Height  int
	mut     sync.RWMutex
	RGB     [][3]byte
	RGBok   [][3]byte
	BgRGB   [3]byte
	FgRGB   [3]byte
	Flip    bool
	syncing bool
	Hover   ObjectPosition
	hovered bool

	// copy of StringGrid position
	FilePosition ObjectPosition
	XCells       int
	YCells       int
}

func (s *Scrollbar) Scroll() {

	var goY = s.Hover.Y
	goY /= rowHeight

	var oldY = s.FilePosition.Y

	if oldY > scrollingOffset {
		goY += oldY - scrollingOffset
	}

	println("scroll", goY, oldY)
	s.FilePosition = ObjectPosition{X: 0, Y: goY}
}
func (s *Scrollbar) IsHover(x, y float32, w, h int32) bool {
	var pos = s.Pos
	if pos.X < 0 {
		pos.X += int(w)
	}
	if pos.Y < 0 {
		pos.Y += int(h)
	}

	var goX = int(x) - pos.X
	var goY = int(y) - pos.Y
	s.mut.RLock()
	sw := int(s.Width)
	sh := int(s.Height)
	s.mut.RUnlock()

	s.Hover.X = goX
	s.Hover.Y = goY

	if (goX >= 0) && (goY >= 0) && (goX <= sw) && (goY <= sh) {

		s.Hover.Y -= s.YCells
		if s.Hover.Y < 0 {
			s.Hover.Y = 0
		}

		s.hovered = true
		return true
	}
	s.hovered = false
	return false
}
func (s *Scrollbar) IsHoverButton() bool {
	return s.hovered
}
func (s *Scrollbar) SyncTo(g *StringGrid) {
	g.FilePosition = s.FilePosition
	sg := g
	sg.SelectionCursor.X = sg.SelectionCursorAbs.X - sg.FilePosition.X
	sg.SelectionCursor.Y = sg.SelectionCursorAbs.Y - sg.FilePosition.Y
	sg.IbeamCursor.X = sg.IbeamCursorAbs.X - sg.FilePosition.X
	sg.IbeamCursor.Y = sg.IbeamCursorAbs.Y - sg.FilePosition.Y
}

func (s *Scrollbar) SyncWith(g *StringGrid) bool {
	s.FilePosition = g.FilePosition
	s.XCells = g.XCells
	s.YCells = g.YCells
	h := g.YCells * g.CellHeight
	var changed bool
	s.mut.Lock()
	changed = s.Height != h
	s.Height = h
	s.mut.Unlock()
	return changed
}

func ScrollbarSync(sb *Scrollbar, p []patchScrollbar, heightLines int) {
	sb.mut.Lock()
	sb.Height = heightLines * rowHeight
	if sb.syncing {
		sb.mut.Unlock()
		return
	}
	sb.syncing = true
	sb.mut.Unlock()

	go sb.Sync(p)
}
func (sb *Scrollbar) RenderRectangle(c Canvas, y int, bgRGB, fgRGB [3]byte) {
	//white rectangle:
	var white [][3]byte
	const width = 96
	if width > sb.YCells*whiteRectThick {
		if len(white) != width*whiteRectThick {
			white = make([][3]byte, width*whiteRectThick)
		}
	} else {
		if len(white) != sb.YCells*whiteRectThick*2 {
			white = make([][3]byte, sb.YCells*whiteRectThick*2)
		}
	}

	var lu, lb, ru ObjectPosition

	lu = sb.Pos
	lu.Y += y

	lb = lu
	lb.Y += sb.YCells * whiteRectThick

	ru = lu
	ru.X += width - whiteRectThick

	c.PutRGB(lu, white, width, bgRGB, fgRGB, true)
	c.PutRGB(lb, white, width, bgRGB, fgRGB, true)
	c.PutRGB(lu, white[0:sb.YCells*whiteRectThick*2], whiteRectThick, bgRGB, fgRGB, true)
	c.PutRGB(ru, white[0:sb.YCells*whiteRectThick*2], whiteRectThick, bgRGB, fgRGB, true)
}

func (sb *Scrollbar) Render(c Canvas) {

	movedPosition := sb.FilePosition.Y
	skippedPostion := 0
	if movedPosition > scrollingOffset {
		skippedPostion = movedPosition - scrollingOffset
		movedPosition = scrollingOffset
	}

	sb.mut.RLock()
	var renderbuf = sb.RGBok
	pos := rowHeight * sb.Width * skippedPostion
	if pos > len(renderbuf) {
		pos = len(renderbuf)
	}
	renderbuf = renderbuf[pos:]
	length := sb.Width * sb.Height
	sb.mut.RUnlock()
	if len(renderbuf) > length {
		renderbuf = renderbuf[:length]
	}
	c.PutRGB(sb.Pos, renderbuf, sb.Width, sb.BgRGB, sb.FgRGB, sb.Flip)

	if sb.hovered {
		sb.RenderRectangle(c, sb.Hover.Y, [3]byte{}, [3]byte{192, 192, 192})
	}
	sb.RenderRectangle(c, movedPosition*rowHeight, [3]byte{}, [3]byte{255, 255, 255})
}

type patchScrollbar struct {
	FileName string
	Pos      ObjectPosition
}

func (sb *Scrollbar) Patch(patch patchScrollbar, data [][3]byte) {
	sb.mut.RLock()
	w := sb.Width
	h := sb.Height
	sb.mut.RUnlock()
	if len(sb.RGB) < w*h {
		sb.RGB = append(sb.RGB, make([][3]byte, (w*h)-len(sb.RGB))...)
	}

	var i = 0
	for y := patch.Pos.Y * w; y < len(sb.RGB); y += w {
		var j = 0

		for x := patch.Pos.X; x < w; x++ {
			if i+j >= len(data) {
				break
			}
			sb.RGB[x+y] = data[i+j]
			j++
		}

		i += w
		if i >= len(data) {
			break
		}
	}
}

func (sb *Scrollbar) Sync(p []patchScrollbar) {

	var syncing = false

	for !syncing {

		for _, patch := range p {
			var data, err = downloadScrollbarPatch(patch.FileName)
			if err != nil {
				println(err.Error())
				continue
			}

			sb.Patch(patch, data)
		}
		var buff = make([][3]byte, sb.Width*sb.Height)
		copy(buff, sb.RGB)
		sb.mut.Lock()
		sb.RGBok = buff
		sb.mut.Unlock()

		//time.Sleep(time.Second)

		sb.mut.Lock()
		syncing = sb.syncing
		sb.syncing = false
		sb.mut.Unlock()

	}
}
