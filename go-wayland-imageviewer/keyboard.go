package main

import (
	"github.com/neurlang/wayland/wl"
	"log"
)

const (
	keyboardEventEnter = 1 << 8
	keyboardEventLeave = 1 << 9
)

func (app *appState) UnFocused() bool {
	var ret1 = (app.pointerEvent.eventMask & keyboardEventLeave) != 0
	return ret1 && !app.pointerEvent.moveWindow
}

func (app *appState) redecorate(configuring bool) {
	if app.decoration != nil {

		app.frame.Rect.Min.X = 0
		app.frame.Rect.Min.Y = 0
		app.frame.Rect.Max.X = int(app.width)
		app.frame.Rect.Max.Y = int(app.height)

		app.decoration.clientSideDecoration(app, true, configuring)

		app.width = int32(app.frame.Rect.Max.X)
		app.height = int32(app.frame.Rect.Max.Y)

		// Draw frame
		buffer := app.drawFrame()

		// Attach new frame to the surface
		if err := app.surface.Attach(buffer, 0, 0); err != nil {
			log.Fatalf("unable to attach buffer to surface: %v", err)
		}

		// Damage the surface
		if err := app.surface.DamageBuffer(0, 0, app.width, app.height); err != nil {
			log.Fatalf("unable to damage buffer: %v", err)
		}

		// Commit the surface state
		if err := app.surface.Commit(); err != nil {
			log.Fatalf("unable to commit surface state: %v", err)
		}
	}
}

func (app *appState) HandleKeyboardEnter(wl.KeyboardEnterEvent) {
	app.pointerEvent.eventMask &= ^keyboardEventLeave
	app.pointerEvent.eventMask |= keyboardEventEnter

	app.redecorate(true)

}

func (app *appState) HandleKeyboardLeave(wl.KeyboardLeaveEvent) {
	app.pointerEvent.eventMask &= ^keyboardEventEnter
	app.pointerEvent.eventMask |= keyboardEventLeave

	app.redecorate(false)

	app.pointerEvent.moveWindow = false

}
func (app *appState) attachKeyboard() {
	keyboard, err := app.seat.GetKeyboard()
	if err != nil {
		log.Fatal("unable to register keyboard interface")
	}
	app.keyboard = keyboard

	keyboard.AddKeyHandler(app)
	keyboard.AddEnterHandler(app)
	keyboard.AddLeaveHandler(app)

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
