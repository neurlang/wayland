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

	if !p.Textarea.StringGrid.IsSelection() {
		//fmt.Println(io.Copy(p, strings.NewReader(p.Textarea.srcClipboard)))
		fmt.Println(p.Close())
		return nil
	}

	var clipbrd = string(p.Textarea.src.CopyBuffer)

	fmt.Println(io.Copy(p, strings.NewReader(clipbrd)))
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
