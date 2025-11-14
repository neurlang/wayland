package main

import "fmt"
import window "github.com/neurlang/wayland/windowtrace"
import "os"
import "bufio"

func (textarea *textarea) selectAll() {

	var endX = textarea.StringGrid.EndLineLen
	var endY = textarea.StringGrid.LineCount - 1
	if endY < 0 {
		endY = 0
	}

	textarea.StringGrid.SelectionCursor = ObjectPosition{-textarea.StringGrid.FilePosition.X, -textarea.StringGrid.FilePosition.Y}
	textarea.StringGrid.SelectionCursorAbs = ObjectPosition{0, 0}
	textarea.StringGrid.IbeamCursor = ObjectPosition{
		endX - textarea.StringGrid.FilePosition.X,
		endY - textarea.StringGrid.FilePosition.Y,
	}
	textarea.StringGrid.IbeamCursorAbs = ObjectPosition{endX, endY}
	textarea.StringGrid.Selecting = false
	textarea.StringGrid.IsSelected = true

}

func (textarea *textarea) copyOperation(input *window.Input, isX bool) {
	var erase = &EraseRequest{
		X0: textarea.StringGrid.IbeamCursorAbsolute().X,     /*+ textarea.StringGrid.FilePosition.X*/
		Y0: textarea.StringGrid.IbeamCursorAbsolute().Y,     /*+ textarea.StringGrid.FilePosition.Y*/
		X1: textarea.StringGrid.SelectionCursorAbsolute().X, /*+ textarea.StringGrid.FilePosition.X*/
		Y1: textarea.StringGrid.SelectionCursorAbsolute().Y, /*+ textarea.StringGrid.FilePosition.Y*/
	}
	if !isX {
		erase = nil
	}

	content, err := load_content(ContentRequest{
		Xpos:   textarea.StringGrid.FilePosition.X,
		Ypos:   textarea.StringGrid.FilePosition.Y,
		Width:  textarea.StringGrid.XCells,
		Height: textarea.StringGrid.YCells,
		Erase:  erase,
		Copy: &CopyRequest{
			X0: textarea.StringGrid.IbeamCursorAbsolute().X,     /*+ textarea.StringGrid.FilePosition.X*/
			Y0: textarea.StringGrid.IbeamCursorAbsolute().Y,     /*+ textarea.StringGrid.FilePosition.Y*/
			X1: textarea.StringGrid.SelectionCursorAbsolute().X, /*+ textarea.StringGrid.FilePosition.X*/
			Y1: textarea.StringGrid.SelectionCursorAbsolute().Y, /*+ textarea.StringGrid.FilePosition.Y*/
		}})
	if err != nil {
		fmt.Println(err)
		return
	}

	textarea.handleContent(content, false)

	if textarea.src != nil {

		textarea.src.RemoveListener(textarea)

		//textarea.src.Destroy()
		//textarea.src.Unregister()
	}

	src, err := textarea.display.CreateDataSource()
	if err != nil {
		fmt.Println(err)
		return
	}
	textarea.src = src

	textarea.src.CopyBuffer = ""
	for i, buf := range content.Copy.Buffer {
		if i+1 == len(content.Copy.Buffer) {
			textarea.src.CopyBuffer += string(buf)
		} else {
			textarea.src.CopyBuffer += string(buf) + "\n"
		}
	}

	textarea.src.Offer("UTF8_STRING")
	textarea.src.Offer("text/plain;charset=utf-8")
	textarea.src.Offer("text/plain;charset=UTF-8")
	textarea.src.AddListener(textarea)

	input.DeviceSetSelection(textarea.src, textarea.display.GetSerial())
}

func (textarea *textarea) saveToFileOperation(filename string) {
	var endX = textarea.StringGrid.EndLineLen
	var endY = textarea.StringGrid.LineCount - 1
	if endY < 0 {
		endY = 0
	}

	content, err := load_content(ContentRequest{
		Xpos:   textarea.StringGrid.FilePosition.X,
		Ypos:   textarea.StringGrid.FilePosition.Y,
		Width:  textarea.StringGrid.XCells,
		Height: textarea.StringGrid.YCells,
		Erase:  nil,
		Copy: &CopyRequest{
			X0: 0,    /*+ textarea.StringGrid.FilePosition.X*/
			Y0: 0,    /*+ textarea.StringGrid.FilePosition.Y*/
			X1: endX, /*+ textarea.StringGrid.FilePosition.X*/
			Y1: endY, /*+ textarea.StringGrid.FilePosition.Y*/
		}})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open or create the file for writing
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(file)

	// Write each byte array to the file followed by a newline
	for i, line := range content.Copy.Buffer {

		if i > 0 {
			_, err := writer.WriteString("\n")
			if err != nil {
				fmt.Println("Error writing newline to file:", err)
				return
			}
		}
		_, err := writer.Write(line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

	}

	// Flush the buffer
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer:", err)
		return
	}
}
