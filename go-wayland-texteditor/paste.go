package main

import "fmt"

type Paste struct {
	linesBuffer [][]byte
	buffer      []byte
	Textarea    *textarea
}

func (p *Paste) Write(b []byte) (int, error) {

	println("WRITING:", len(b))

	p.buffer = append(p.buffer, b...)

	var wasR bool
	var consume int
	for i, ch := range p.buffer {
		if wasR && ch == '\n' {
			consume++
			wasR = false
			continue
		}
		wasR = false
		if ch == '\n' || ch == '\r' {
			p.linesBuffer = append(p.linesBuffer, p.buffer[consume:i])
			consume = i + 1
			wasR = ch == '\r'
		}
	}
	p.buffer = p.buffer[consume:]

	return len(b), nil
}

func (p *Paste) Close() error {

	p.Textarea.mutex.Lock()
	defer p.Textarea.mutex.Unlock()

	var paste = &PasteRequest{
		X:      p.Textarea.StringGrid.IbeamCursorAbsolute().X,
		Y:      p.Textarea.StringGrid.IbeamCursorAbsolute().Y,
		Buffer: append(p.linesBuffer, p.buffer),
	}

	var erase = &EraseRequest{
		X0: p.Textarea.StringGrid.IbeamCursorAbsolute().X,     /*+ textarea.StringGrid.FilePosition.X*/
		Y0: p.Textarea.StringGrid.IbeamCursorAbsolute().Y,     /*+ textarea.StringGrid.FilePosition.Y*/
		X1: p.Textarea.StringGrid.SelectionCursorAbsolute().X, /*+ textarea.StringGrid.FilePosition.X*/
		Y1: p.Textarea.StringGrid.SelectionCursorAbsolute().Y, /*+ textarea.StringGrid.FilePosition.Y*/
	}
	if !(p.Textarea.StringGrid.IsSelection() && p.Textarea.StringGrid.IsSelectionStrict()) {
		erase = nil
	} else {
		var pasteErase = &PasteRequest{
			X:      (&p.Textarea.StringGrid).IbeamCursorAbsolute().Lesser(p.Textarea.StringGrid.SelectionCursorAbsolute()).X, /*+ textarea.StringGrid.FilePosition.X*/
			Y:      (&p.Textarea.StringGrid).IbeamCursorAbsolute().Lesser(p.Textarea.StringGrid.SelectionCursorAbsolute()).Y, /*+ textarea.StringGrid.FilePosition.Y*/
			Buffer: append(p.linesBuffer, p.buffer),
		}
		paste = pasteErase
	}
	content, err := load_content(ContentRequest{
		Xpos:   p.Textarea.StringGrid.FilePosition.X,
		Ypos:   p.Textarea.StringGrid.FilePosition.Y,
		Width:  p.Textarea.StringGrid.XCells,
		Height: p.Textarea.StringGrid.YCells,
		Erase:  erase,
		Paste:  paste,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	p.Textarea.handleContent(content)
	return nil
}
