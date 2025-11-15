package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const tabSize = 8

var fileMutex sync.Mutex
var fileTab = map[string][][]string{
	"": {{"H", "e", "l", "l", "o", "c", "r", "u", "e", "l"},
		{"w", "o", "r", "l", "d"}},
}

var fileColor [][5]int

func handlerCopy(tab string, p *CopyRequest) *CopyResponse {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file := fileTab[tab]

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
	cr.Buffer = [][]byte{{}}
	for y := p.Y0; y <= p.Y1; y++ {
		for x := 0; x < len(file[y]); x++ {
			if y == p.Y0 && x < p.X0 {
				continue
			}
			if y == p.Y1 && x >= p.X1 {
				break
			}
			cr.Buffer[len(cr.Buffer)-1] = append(cr.Buffer[len(cr.Buffer)-1], []byte(file[y][x])...)
		}
		if y != p.Y1 {
			cr.Buffer = append(cr.Buffer, []byte{})
		}
	}
	return &cr
}
func isCombinable(r rune) bool {
	return r > 128
}
func isCombiner(r rune) bool {
	// combining diacritical marks block
	if r >= 0x300 && r <= 0x36f {
		return true
	}
	switch r {
	case 0x900, 0x901, 0x902, 0x903, 0x93a, 0x93b, 0x93c, 0x93e, 0x93f, 0x940:
		return true // devanagari
	case 0x941, 0x942, 0x943, 0x944, 0x945, 0x946, 0x947, 0x948, 0x949, 0x94a:
		return true // devanagari
	case 0x94b, 0x94c, 0x94d, 0x94e, 0x94f, 0x955, 0x956, 0x957, 0x962, 0x963:
		return true // devanagari
	default:
		return false // not a combiner
	}
}
func handlerErase(tab string, e *EraseRequest) *EraseResponse {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file := fileTab[tab]

	if (e.Y0 > e.Y1) || ((e.Y0 == e.Y1) && (e.X0 > e.X1)) {
		e.X0, e.X1 = e.X1, e.X0
		e.Y0, e.Y1 = e.Y1, e.Y0
	}
	if e.X1 == 0 && e.Y1 >= len(file) && e.Y1 > 0 { // all file erase
		e.Y1 = len(file) - 1
		e.X1 = len(file[e.Y1])
	}
	if e.Y0 >= len(file) {
		return nil
	}
	if e.Y1 >= len(file) {
		return nil
	}
	if e.X0 > len(file[e.Y0]) {
		return nil
	}
	if e.X1 > len(file[e.Y1]) {
		return nil
	}

	if e.Y0 != e.Y1 {
		file[e.Y0] = file[e.Y0][:e.X0]
		file[e.Y1] = file[e.Y1][e.X1:]

		file[e.Y0] = append(file[e.Y0], file[e.Y1]...)

		file = append(file[:e.Y0+1], file[e.Y1+1:]...)

		fileTab[tab] = file
	} else {
		file[e.Y0] = append(file[e.Y0][:e.X0], file[e.Y0][e.X1:]...)
	}
	return &EraseResponse{Erased: true}
}
func handlerPaste(tab string, p *PasteRequest) *struct{} {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file := fileTab[tab]

	if p.Y >= len(file) {
		return nil
	}
	if p.X > len(file[p.Y]) {
		return nil
	}
	temp := p.Buffer
	if len(temp) > 0 {
		file = append(file[:p.Y+1], append(make([][]string, len(temp)-1), file[p.Y+1:]...)...)

		fileTab[tab] = file
	}
	var rrow []string

	for i, subarray := range temp {
		if i > 0 {
			p.X = 0
			p.Y++
		}
		if len(subarray) == 0 && i+1 == len(temp) {
			break
		}
		if p.Y >= len(file) {
			file = append(file, []string{})

			fileTab[tab] = file
		}
		var row = file[p.Y]
		if p.X > len(row) {
			p.X = 0
		}
		if rrow == nil {
			rrow = row[p.X:]
		}
		row = row[:p.X:p.X]

		var isCombinAble = false

		for _, c := range []rune(string(subarray)) { //nolint:gosimple,staticcheck // subarray is []byte, need runes for isCombiner()
			var char = string(c)
			if char == "\t" {
				for len(row)&(tabSize-1) != (tabSize - 1) {
					row = append(row, "")
				}
			}
			if (len(row) > 0) && isCombiner(c) && isCombinAble {
				row[len(row)-1] = row[len(row)-1] + char
			} else {
				isCombinAble = isCombinable(c)
				row = append(row, char)
			}
		}
		file[p.Y] = row
	}
	if p.Y >= len(file) {
		file = append(file, []string{})

		fileTab[tab] = file
	}
	file[p.Y] = append(file[p.Y], rrow...)
	return &struct{}{}
}

func handlerWrite(tab string, w *WriteRequest) *WriteResponse {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file := fileTab[tab]

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

		fileTab[tab] = file
	case "Delete":
		if w.X == len(row) {
			if len(file) <= w.Y+1 {
				return &wr
			}
			file[w.Y] = append(row, file[w.Y+1]...)
			file = append(file[:w.Y+1], file[w.Y+2:]...)

			fileTab[tab] = file

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

				fileTab[tab] = file
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
	Buffer [][]byte
}

type ContentRequest struct {
	Xpos, Ypos, Width, Height int
	Tab                       string
	Copy                      *CopyRequest
	Erase                     *EraseRequest
	Write                     *WriteRequest
	Paste                     *PasteRequest
	MultiClick                *MultiClickRequest
}

type ContentResponse struct {
	Content    [][]string
	FgColor    [][5]int
	LineLens   []int
	LineCount  int
	EndLineLen int
	Xpos, Ypos int
	Copy       *CopyResponse
	Erase      *EraseResponse
	Write      *WriteResponse
	Paste      *struct{}
	MultiClick *MultiClickResponse
}
type EraseRequest struct {
	X0, Y0, X1, Y1 int
}
type CopyRequest struct {
	X0, Y0, X1, Y1 int
}
type MultiClickRequest struct {
	Double bool
}
type MultiClickResponse struct {
}
type CopyResponse struct {
	Buffer [][]byte
}
type WriteResponse struct {
	MoveX, MoveY int
}
type EraseResponse struct {
	Erased bool
}

func handlerContent(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
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
	var recolor bool
	if cr.Copy != nil {
		resp.Copy = handlerCopy(cr.Tab, cr.Copy)
	}
	if cr.Erase != nil {
		resp.Erase = handlerErase(cr.Tab, cr.Erase)
		if resp.Erase != nil {
			recolor = true
		}
	}
	if cr.Write != nil {
		resp.Write = handlerWrite(cr.Tab, cr.Write)
		if resp.Write != nil {
			recolor = true
		}
	}
	if cr.Paste != nil {
		resp.Paste = handlerPaste(cr.Tab, cr.Paste)
		if resp.Paste != nil {
			recolor = true
		}
	}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file := fileTab[cr.Tab]

	if recolor {
		var highlights = []func([][]string) [][5]int{
			reprocess_syntax_highlighting_golang,
			reprocess_syntax_highlighting_csharp,
		}
		for _, highlight := range highlights {
			fileColor = highlight(file)
			recolor = fileColor == nil
			if !recolor {
				break
			}
		}
	}

	var minI = 0
	var maxI = len(fileColor)

	for i := range fileColor {
		if fileColor[i][1] < cr.Ypos {
			minI = i
		} else if fileColor[i][1] >= cr.Ypos+cr.Height {
			maxI = i
			break
		}
	}
	resp.FgColor = fileColor[minI:maxI]

	for y := cr.Ypos; y < len(file) && y < cr.Ypos+cr.Height; y++ {
		resp.LineLens = append(resp.LineLens, len(file[y]))
		resp.Content = append(resp.Content, nil)
		for x := cr.Xpos; x < cr.Xpos+cr.Width; x++ {
			if !(y >= len(file) || x >= len(file[y])) {
				resp.Content[len(resp.Content)-1] = append(resp.Content[len(resp.Content)-1], file[y][x])
			}
		}
	}

	resp.LineCount = len(file)
	if len(file) > 0 {
		resp.EndLineLen = len(file[len(file)-1])
	}
	resp.Xpos = cr.Xpos
	resp.Ypos = cr.Ypos

	bytes, err := json.Marshal(resp)
	if err != nil {
		return
	}

	if _, err := w.Write(bytes); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
func handlerScrollbar(w http.ResponseWriter, req *http.Request) {

	path := req.URL.Path[strings.LastIndex(req.URL.Path, "/"):]
	path = strings.TrimSuffix(path, ".png")
	path = strings.TrimPrefix(path, "/")
	if path == "live" {
		path = ""
	}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file := fileTab[path]

	body, err := reprocess_scrollbar(file)
	if err != nil {
		return
	}
	if _, err := w.Write(body); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
func main() {
	http.HandleFunc("/content", handlerContent)
	http.HandleFunc("/scrollbar/", handlerScrollbar)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
