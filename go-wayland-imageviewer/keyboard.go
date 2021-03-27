package main

import (
	"github.com/neurlang/wayland/wl"
	"log"
)

func (app *appState) attachKeyboard() {
	keyboard, err := app.seat.GetKeyboard()
	if err != nil {
		log.Fatal("unable to register keyboard interface")
	}
	app.keyboard = keyboard

	keyboard.AddKeyHandler(app)

	log.Print("keyboard interface registered")
}

func (app *appState) releaseKeyboard() {
	app.keyboard.RemoveKeyHandler(app)

	if err := app.keyboard.Release(); err != nil {
		log.Println("unable to release keyboard interface")
	}
	app.keyboard = nil

	log.Print("keyboard interface released")
}

func (app *appState) HandleKeyboardKey(e wl.KeyboardKeyEvent) {
	// close on "q"
	if e.Key == 16 {
		app.exit = true
	}
}
