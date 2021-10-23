package main

import "fmt"

type Paste struct {
	buffer   []byte
	Textarea *textarea
}

func (p *Paste) Write(b []byte) (int, error) {
	p.buffer = append(p.buffer, b...)
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
			Buffer: p.buffer,
		}})
	if err != nil {
		fmt.Println(err)
		return err
	}

	p.Textarea.handleContent(content)
	return nil
}
