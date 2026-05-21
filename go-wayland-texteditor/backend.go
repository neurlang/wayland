package main

import (
	lib "github.com/neurlang/wayland/go-wayland-texteditor/lib_editor_backend"
)

var libInstance = getLibInstance()

func getLibInstance() lib.Interface {
	var l, err = lib.New()
	if err != nil {
		panic(err)
	}
	return l
}
