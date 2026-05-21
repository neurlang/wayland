package main

import (
	"io"
	"log"
	"net/http"
	lib "github.com/neurlang/wayland/go-wayland-texteditor/lib_editor_backend"
)

func init() {
	var lib, err = lib.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/content", func (w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		println(string(body))
		resp := lib.Call("/content", string(body))
		if _, err := w.Write([]byte(resp)); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})
	http.HandleFunc("/scrollbar/", func (w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		println(string(body))
		resp := lib.Call(r.URL.Path, string(body))
		if _, err := w.Write([]byte(resp)); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})
	go http.ListenAndServe(":8080", nil)
}
