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

	content, err := load_content(ContentRequest{
		Width:  p.Textarea.StringGrid.XCells,
		Height: p.Textarea.StringGrid.YCells,
		Paste: &PasteRequest{
			X:      p.Textarea.StringGrid.IbeamCursor.X,
			Y:      p.Textarea.StringGrid.IbeamCursor.Y,
			Buffer: append(p.linesBuffer, p.buffer),
		}})
	if err != nil {
		fmt.Println(err)
		return err
	}

	p.Textarea.handleContent(content)
	return nil
}
