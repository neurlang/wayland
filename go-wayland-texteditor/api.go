package main

import (
	"encoding/json"
)

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
	LineCount  int
	EndLineLen int
	LineLens   []int
	Xpos, Ypos int
	Copy       *CopyResponse
	Erase      *EraseResponse
	Write      *WriteResponse
	MultiClick *MultiClickResponse
	Paste      *struct{}
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


func load_content(creq ContentRequest) (*ContentResponse, error) {
	requestByte, err := json.Marshal(creq)
	if err != nil {
		return nil, err
	}
	resp := libInstance.Call("/content", string(requestByte))
	var cr ContentResponse
	err = json.Unmarshal([]byte(resp), &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}
