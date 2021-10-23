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
	Buffer []byte
}

type ContentRequest struct {
	Xpos, Ypos, Width, Height int
	Write                     *WriteRequest
	Paste                     *PasteRequest
}
type ContentResponse struct {
	Content []string
	FgColor [][5]int
	Write   *WriteResponse
}

type WriteResponse struct {
	MoveX, MoveY int
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
	var cr ContentResponse
	//println(string(body))
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}
