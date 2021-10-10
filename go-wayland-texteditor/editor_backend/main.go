package main

import (
    "log"
    "net/http"
    "encoding/json"
    "io/ioutil"
)


var file = [][]string{
[]string{"H","e","l","l","o",},
[]string{"w","o","r","l","d",},
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
		if (w.X == len(row)) {
			file[w.Y] = append(row, file[w.Y+1]...)
			file = append(file[:w.Y+1], file[w.Y+2:]...)
			return &wr
		}
		row = append(row[:w.X], row[w.X+1:]...)
		file[w.Y] = row
	case "Backspace":
		if (w.X == 0) {
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
		wr.MoveX = -1
	default:
		if (w.X == len(row)) {
			row = append(row, w.Key)
			file[w.Y] = row
		} else if (w.Insert) {
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
	X, Y int
	Key string
	Insert bool
}

type ContentRequest struct {
	Xpos, Ypos, Width, Height int
	Write *WriteRequest
}

type ContentResponse struct {
	Content []string
	Write *WriteResponse
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
	if cr.Write != nil {
		resp.Write = handlerWrite(cr.Write)
	}
	
	

	
	for y := cr.Ypos; y < cr.Ypos+cr.Height; y++ {
	for x := cr.Xpos; x < cr.Xpos+cr.Width; x++ {
		if y >= len(file) || x >= len(file[y]) {
			resp.Content = append(resp.Content, " ")
		} else {
			resp.Content = append(resp.Content, file[y][x])
		}
	}
	}

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
