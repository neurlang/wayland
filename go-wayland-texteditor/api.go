package main

import "net/http"
import "io/ioutil"
import "encoding/json"
import "bytes"

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
	Copy                      *CopyRequest
	Erase                     *EraseRequest
	Write                     *WriteRequest
	Paste                     *PasteRequest
}
type ContentResponse struct {
	Content    [][]string
	FgColor    [][5]int
	LineCount  int
	EndLineLen int
	LineLens   []int
	Xpos, Ypos int
	Copy       *CopyResponse
	Erase      *EraseResponse
	Write      *WriteResponse
	Paste      *struct{}
}
type EraseRequest struct {
	X0, Y0, X1, Y1 int
}
type CopyRequest struct {
	X0, Y0, X1, Y1 int
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

func (req *ContentRequest) Reader() *bytes.Reader {
	requestByte, _ := json.Marshal(req)
	requestReader := bytes.NewReader(requestByte)
	return requestReader
}

func load_content(creq ContentRequest) (*ContentResponse, error) {
	resp, err := http.Post("http://127.0.0.1:8080/content", "application/json", creq.Reader())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cr ContentResponse
	//println(string(body))
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}
