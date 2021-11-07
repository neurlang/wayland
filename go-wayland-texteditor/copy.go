package main

import "io"
import "os"
import "fmt"
import "strings"

type Copy struct {
	file     *os.File
	Textarea *textarea
}

func (p *Copy) Receive(fd uintptr, name string) error {
	p.file = os.NewFile(fd, name)

	p.Textarea.mutex.Lock()
	defer p.Textarea.mutex.Unlock()

	if p.Textarea.StringGrid.IbeamCursor == p.Textarea.StringGrid.SelectionCursor {
		fmt.Println(io.Copy(p, strings.NewReader(p.Textarea.srcClipboard)))
		fmt.Println(p.Close())
		return nil
	}

	content, err := load_content(ContentRequest{
		Width:  p.Textarea.StringGrid.XCells,
		Height: p.Textarea.StringGrid.YCells,
		Copy: &CopyRequest{
			X0: p.Textarea.StringGrid.IbeamCursor.X,
			Y0: p.Textarea.StringGrid.IbeamCursor.Y,
			X1: p.Textarea.StringGrid.SelectionCursor.X,
			Y1: p.Textarea.StringGrid.SelectionCursor.Y,
		}})
	if err != nil {
		fmt.Println(err)
		return err
	}

	p.Textarea.handleContent(content)

	p.Textarea.srcClipboard = string(content.Copy.Buffer)

	fmt.Println(io.Copy(p, strings.NewReader(p.Textarea.srcClipboard)))
	fmt.Println(p.Close())

	return nil
}
func (p *Copy) Write(buf []byte) (int, error) {

	println("COPYING:", len(buf))

	return p.file.Write(buf)
}

func (p *Copy) Close() error {
	return p.file.Close()
}
