package main

import (
	wl "github.com/neurlang/wayland/wl"
	cursor "github.com/neurlang/wayland/wlcursor"
	xdg_shell "github.com/neurlang/wayland/xdg"

	"log"
)

const (
	pointerEventEnter        = 1 << 0
	pointerEventLeave        = 1 << 1
	pointerEventMotion       = 1 << 2
	pointerEventButton       = 1 << 3
	pointerEventAxis         = 1 << 4
	pointerEventAxisSource   = 1 << 5
	pointerEventAxisStop     = 1 << 6
	pointerEventAxisDiscrete = 1 << 7
)

// From linux/input-event-codes.h
const (
	BtnLeft   = 0x110
	BtnRight  = 0x111
	BtnMiddle = 0x112
)

type pointerEvent struct {
	eventMask          int
	surfaceX, surfaceY uint32
	button, state      uint32
	time               uint32
	serial             uint32
	axes               [2]struct {
		valid    bool
		value    int32
		discrete int32
	}
	axisSource uint32
}

func (app *appState) attachPointer() {
	pointer, err := app.seat.GetPointer()
	if err != nil {
		log.Fatal("unable to register pointer interface")
	}
	app.pointer = pointer
	pointer.AddEnterHandler(app)
	pointer.AddLeaveHandler(app)
	pointer.AddMotionHandler(app)
	pointer.AddButtonHandler(app)
	pointer.AddAxisHandler(app)
	pointer.AddAxisSourceHandler(app)
	pointer.AddAxisStopHandler(app)
	pointer.AddAxisDiscreteHandler(app)
	pointer.AddFrameHandler(app)

	log.Print("pointer interface registered")
}

func (app *appState) releasePointer() {
	app.pointer.RemoveEnterHandler(app)
	app.pointer.RemoveLeaveHandler(app)
	app.pointer.RemoveMotionHandler(app)
	app.pointer.RemoveButtonHandler(app)
	app.pointer.RemoveAxisHandler(app)
	app.pointer.RemoveAxisSourceHandler(app)
	app.pointer.RemoveAxisStopHandler(app)
	app.pointer.RemoveAxisDiscreteHandler(app)
	app.pointer.RemoveFrameHandler(app)

	if err := app.pointer.Release(); err != nil {
		log.Println("unable to release pointer interface")
	}
	app.pointer = nil

	log.Print("pointer interface released")
}

func (app *appState) HandlePointerEnter(e wl.PointerEnterEvent) {
	app.pointerEvent.eventMask |= pointerEventEnter
	app.pointerEvent.serial = e.Serial
	app.pointerEvent.surfaceX = uint32(e.SurfaceX)
	app.pointerEvent.surfaceY = uint32(e.SurfaceY)
}

func (app *appState) HandlePointerLeave(e wl.PointerLeaveEvent) {
	app.pointerEvent.eventMask |= pointerEventLeave
	app.pointerEvent.serial = e.Serial
}

func (app *appState) HandlePointerMotion(e wl.PointerMotionEvent) {
	app.pointerEvent.eventMask |= pointerEventMotion
	app.pointerEvent.time = e.Time
	app.pointerEvent.surfaceX = uint32(e.SurfaceX)
	app.pointerEvent.surfaceY = uint32(e.SurfaceY)
}

func (app *appState) HandlePointerButton(e wl.PointerButtonEvent) {
	app.pointerEvent.eventMask |= pointerEventButton
	app.pointerEvent.serial = e.Serial
	app.pointerEvent.time = e.Time
	app.pointerEvent.button = e.Button
	app.pointerEvent.state = e.State
}

func (app *appState) HandlePointerAxis(e wl.PointerAxisEvent) {
	app.pointerEvent.eventMask |= pointerEventAxis
	app.pointerEvent.time = e.Time
	app.pointerEvent.axes[e.Axis].valid = true
	app.pointerEvent.axes[e.Axis].value = int32(e.Value)
}

func (app *appState) HandlePointerAxisSource(e wl.PointerAxisSourceEvent) {
	app.pointerEvent.eventMask |= pointerEventAxis
	app.pointerEvent.axisSource = e.AxisSource
}

func (app *appState) HandlePointerAxisStop(e wl.PointerAxisStopEvent) {
	app.pointerEvent.eventMask |= pointerEventAxisStop
	app.pointerEvent.time = e.Time
	app.pointerEvent.axes[e.Axis].valid = true
}

func (app *appState) HandlePointerAxisDiscrete(e wl.PointerAxisDiscreteEvent) {
	app.pointerEvent.eventMask |= pointerEventAxisDiscrete
	app.pointerEvent.axes[e.Axis].valid = true
	app.pointerEvent.axes[e.Axis].discrete = e.Discrete
}

var axisName = map[int]string{
	wl.PointerAxisVerticalScroll:   "vertical",
	wl.PointerAxisHorizontalScroll: "horizontal",
}

var axisSource = map[uint32]string{
	wl.PointerAxisSourceWheel:      "wheel",
	wl.PointerAxisSourceFinger:     "finger",
	wl.PointerAxisSourceContinuous: "continuous",
	wl.PointerAxisSourceWheelTilt:  "wheel tilt",
}

var cursorMap = map[uint32]string{
	xdg_shell.ToplevelResizeEdgeTop:         cursor.TopSide,
	xdg_shell.ToplevelResizeEdgeTopLeft:     cursor.TopLeftCorner,
	xdg_shell.ToplevelResizeEdgeTopRight:    cursor.TopRightCorner,
	xdg_shell.ToplevelResizeEdgeBottom:      cursor.BottomSide,
	xdg_shell.ToplevelResizeEdgeBottomLeft:  cursor.BottomLeftCorner,
	xdg_shell.ToplevelResizeEdgeBottomRight: cursor.BottomRightCorner,
	xdg_shell.ToplevelResizeEdgeLeft:        cursor.LeftSide,
	xdg_shell.ToplevelResizeEdgeRight:       cursor.RightSide,
	xdg_shell.ToplevelResizeEdgeNone:        cursor.LeftPtr,
}

func (app *appState) HandlePointerFrame(_ wl.PointerFrameEvent) {
	e := app.pointerEvent

	if (e.eventMask & pointerEventEnter) != 0 {
		log.Printf("entered %v, %v", e.surfaceX, e.surfaceY)

		app.setCursor(e.serial, cursor.LeftPtr)
	}

	if (e.eventMask & pointerEventLeave) != 0 {
		log.Print("leave")
	}
	if (e.eventMask & pointerEventMotion) != 0 {
		log.Printf("motion %v, %v", e.surfaceX, e.surfaceY)

		edge := componentEdge(uint32(app.width), uint32(app.height), e.surfaceX, e.surfaceY, 8)
		cursorName, ok := cursorMap[edge]
		if ok && cursorName != app.currentCursor {
			app.setCursor(e.serial, cursorName)
		}
	}
	if (e.eventMask & pointerEventButton) != 0 {
		if e.state == wl.PointerButtonStateReleased {
			log.Printf("button %d released", e.button)
		} else {
			log.Printf("button %d pressed", e.button)

			switch e.button {
			case BtnLeft:
				edge := componentEdge(uint32(app.width), uint32(app.height), e.surfaceX, e.surfaceY, 8)
				if edge != xdg_shell.ToplevelResizeEdgeNone {
					if err := app.xdgTopLevel.Resize(app.seat, e.serial, edge); err != nil {
						log.Println("unable to start resize")
					}
				} else {
					if err := app.xdgTopLevel.Move(app.seat, e.serial); err != nil {
						log.Println("unable to start move")
					}
				}
			case BtnRight:
				if err := app.xdgTopLevel.ShowWindowMenu(app.seat, e.serial, int32(e.surfaceX), int32(e.surfaceY)); err != nil {
					log.Println("unable to show window menu")
				}
			}
		}
	}

	const axisEvents = pointerEventAxis | pointerEventAxisSource | pointerEventAxisStop | pointerEventAxisDiscrete

	if (e.eventMask & axisEvents) != 0 {
		for i := 0; i < 2; i++ {
			if !e.axes[i].valid {
				continue
			}
			log.Printf("%s axis ", axisName[i])
			if (e.eventMask & pointerEventAxis) != 0 {
				log.Printf("value %v", e.axes[i].value)
			}
			if (e.eventMask & pointerEventAxisDiscrete) != 0 {
				log.Printf("discrete %d ", e.axes[i].discrete)
			}
			if (e.eventMask & pointerEventAxisSource) != 0 {
				log.Printf("via %s", axisSource[e.axisSource])
			}
			if (e.eventMask & pointerEventAxisStop) != 0 {
				log.Printf("(stopped)")
			}
		}
	}
}

func componentEdge(width, height, pointerX, pointerY, margin uint32) uint32 {
	top := pointerY < margin
	bottom := pointerY > (height - margin)
	left := pointerX < margin
	right := pointerX > (width - margin)

	if top {
		if left {
			log.Print("ToplevelResizeEdgeTopLeft")
			return xdg_shell.ToplevelResizeEdgeTopLeft
		} else if right {
			log.Print("ToplevelResizeEdgeTopRight")
			return xdg_shell.ToplevelResizeEdgeTopRight
		} else {
			log.Print("ToplevelResizeEdgeTop")
			return xdg_shell.ToplevelResizeEdgeTop
		}
	} else if bottom {
		if left {
			log.Print("ToplevelResizeEdgeBottomLeft")
			return xdg_shell.ToplevelResizeEdgeBottomLeft
		} else if right {
			log.Print("ToplevelResizeEdgeBottomRight")
			return xdg_shell.ToplevelResizeEdgeBottomRight
		} else {
			log.Print("ToplevelResizeEdgeBottom")
			return xdg_shell.ToplevelResizeEdgeBottom
		}
	} else if left {
		log.Print("ToplevelResizeEdgeLeft")
		return xdg_shell.ToplevelResizeEdgeLeft
	} else if right {
		log.Print("ToplevelResizeEdgeRight")
		return xdg_shell.ToplevelResizeEdgeRight
	} else {
		log.Print("ToplevelResizeEdgeNone")
		return xdg_shell.ToplevelResizeEdgeNone
	}
}

type cursorData struct {
	wlCursor *cursor.Cursor

	surface *wl.Surface
}

func (app *appState) setCursor(serial uint32, cursorName string) {
	c, ok := app.cursors[cursorName]
	if !ok {
		log.Print("unable to get %v cursor", cursorName)
		return
	}

	image := c.wlCursor.Images[0]
	if err := app.pointer.SetCursor(
		serial, c.surface,
		int32(image.GetHotspotX()), int32(image.GetHotspotY()),
	); err != nil {
		log.Print("unable to set cursor")
	}

	app.currentCursor = cursorName
}
