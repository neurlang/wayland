package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const tabSize = 8

var file = [][]string{
	{"H", "e", "l", "l", "o", "c", "r", "u", "e", "l"},
	{"w", "o", "r", "l", "d"},
}

func handlerCopy(p *CopyRequest) *CopyResponse {
	var cr CopyResponse
	if (p.Y0 > p.Y1) || ((p.Y0 == p.Y1) && (p.X0 > p.X1)) {
		p.X0, p.X1 = p.X1, p.X0
		p.Y0, p.Y1 = p.Y1, p.Y0
	}
	if p.Y0 >= len(file) {
		return &cr
	}
	if p.X0 > len(file[p.Y0]) {
		return &cr
	}
	if p.Y1 >= len(file) {
		return &cr
	}
	if p.X1 > len(file[p.Y1]) {
		return &cr
	}

	for y := p.Y0; y <= p.Y1; y++ {
		for x := 0; x < len(file[y]); x++ {
			if y == p.Y0 && x < p.X0 {
				continue
			}
			if y == p.Y1 && x >= p.X1 {
				break
			}
			cr.Buffer = append(cr.Buffer, []byte(file[y][x])...)
		}
		cr.Buffer = append(cr.Buffer, '\r', '\n')
	}
	return &cr
}

func handlerPaste(p *PasteRequest) {
	var array = string(p.Buffer)
	if p.Y >= len(file) {
		return
	}
	if p.X > len(file[p.Y]) {
		return
	}

	temp := strings.Split(strings.ReplaceAll(array, "\r\n", "\n"), "\n")
	if len(temp) > 0 {
		file = append(file[:p.Y+1], append(make([][]string, len(temp)-1), file[p.Y+1:]...)...)
	}
	var rrow []string

	for i, subarray := range temp {
		if len(subarray) == 0 && i+1 == len(temp) {
			break
		}
		array := []rune(subarray)
		if p.Y >= len(file) {
			file = append(file, []string{})
		}
		var row = file[p.Y]
		if p.X > len(row) {
			p.X = 0
		}
		if rrow == nil {
			rrow = row[p.X:]
		}
		row = row[:p.X:p.X]

		for _, c := range array {
			var char = string(c)
			if char == "\t" {
				for len(row)&(tabSize-1) != (tabSize - 1) {
					row = append(row, "")
				}
			}
			row = append(row, char)
		}
		file[p.Y] = row

		p.X = 0
		p.Y++
	}
	if p.Y >= len(file) {
		file = append(file, []string{})
	}
	file[p.Y] = append(file[p.Y], rrow...)
}

func handlerWrite(w *WriteRequest) *WriteResponse {
	var wr WriteResponse
	if w.Y >= len(file) {
		return &wr
	}
	var row = file[w.Y]
	if w.X > len(row) {
		return &wr
	}
	switch w.Key {
	case "Enter":
		file = append(file[:w.Y+1], file[w.Y:]...)
		file[w.Y] = append([]string{}, file[w.Y+1][:w.X]...)
		file[w.Y+1] = file[w.Y+1][w.X:]
		wr.MoveX = -w.X
		wr.MoveY = 1
	case "Delete":
		if w.X == len(row) {
			file[w.Y] = append(row, file[w.Y+1]...)
			file = append(file[:w.Y+1], file[w.Y+2:]...)
			return &wr
		}
		for row[w.X] == "" {
			row = append(row[:w.X], row[w.X+1:]...)
		}
		row = append(row[:w.X], row[w.X+1:]...)
		file[w.Y] = row
	case "Backspace":
	again:
		if w.X == 0 {
			if w.Y != 0 {
				wr.MoveX = len(file[w.Y-1])
				wr.MoveY = -1
				file[w.Y-1] = append(file[w.Y-1], row...)
				file = append(file[:w.Y], file[w.Y+1:]...)
			}

			return &wr
		}
		row = append(row[:w.X-1], row[w.X:]...)
		file[w.Y] = row
		wr.MoveX--
		w.X--
		if w.X >= 1 && row[w.X-1] == "" {
			goto again
		}
	case "\t":
		for {
			if w.X == len(row) {
				row = append(row, w.Key)
				file[w.Y] = row
			} else if w.Insert {
				row = append(row[:w.X+1], row[w.X:]...)
				file[w.Y] = row
			}
			if (w.X & (tabSize - 1)) == (tabSize - 1) {
				row[w.X] = w.Key
			} else {
				row[w.X] = ""
			}
			println(w.Key)
			wr.MoveX++
			w.X++
			if (w.X & (tabSize - 1)) == 0 {
				break
			}
		}
	default:
		if w.X == len(row) {
			row = append(row, w.Key)
			file[w.Y] = row
		} else if w.Insert && row[w.X] != "" {
			row = append(row[:w.X+1], row[w.X:]...)
			file[w.Y] = row
		}
		row[w.X] = w.Key
		println(w.Key)
		wr.MoveX = 1
	}
	return &wr
}

/////////////////

type WriteRequest struct {
	X, Y   int
	Key    string
	Insert bool
}

type PasteRequest struct {
	X, Y   int
	Buffer []byte
}

type ContentRequest struct {
	Xpos, Ypos, Width, Height int
	Copy                      *CopyRequest
	Write                     *WriteRequest
	Paste                     *PasteRequest
}

type ContentResponse struct {
	Content []string
	FgColor [][5]int
	Copy    *CopyResponse
	Write   *WriteResponse
}

type CopyRequest struct {
	X0, Y0, X1, Y1 int
}
type CopyResponse struct {
	Buffer []byte
}
type WriteResponse struct {
	MoveX, MoveY int
}

func handlerContent(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	println(string(body))

	var cr ContentRequest
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return
	}
	var resp ContentResponse
	if cr.Copy != nil {
		resp.Copy = handlerCopy(cr.Copy)
	}
	if cr.Write != nil {
		resp.Write = handlerWrite(cr.Write)
	}
	if cr.Paste != nil {
		handlerPaste(cr.Paste)
	}

	for y := cr.Ypos; y < cr.Ypos+cr.Height; y++ {
		for x := cr.Xpos; x < cr.Xpos+cr.Width; x++ {
			if y >= len(file) || x >= len(file[y]) {
				resp.Content = append(resp.Content, "")
			} else {
				resp.Content = append(resp.Content, file[y][x])
			}
		}
	}

	resp.FgColor = reprocess_syntax_highlighting_golang(file)

	bytes, err := json.Marshal(resp)
	if err != nil {
		return
	}

	w.Write(bytes)
}

func main() {
	http.HandleFunc("/content", handlerContent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
