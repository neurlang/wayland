package main

import (
	"github.com/neurlang/wayland/wl"
	cursor "github.com/neurlang/wayland/wlcursor"
	xdgshell "github.com/neurlang/wayland/xdg"

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
	app.pointer.Unregister()
	app.pointer = nil

	log.Print("pointer interface released")
}

func (app *appState) HandlePointerEnter(e wl.PointerEnterEvent) {
	app.pointerEvent.eventMask &= ^pointerEventLeave
	app.pointerEvent.eventMask |= pointerEventEnter
	app.pointerEvent.serial = e.Serial
	app.pointerEvent.surfaceX = uint32(e.SurfaceX)
	app.pointerEvent.surfaceY = uint32(e.SurfaceY)
}

func (app *appState) HandlePointerLeave(e wl.PointerLeaveEvent) {
	app.pointerEvent.eventMask &= ^pointerEventEnter
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
	xdgshell.ToplevelResizeEdgeTop:         cursor.TopSide,
	xdgshell.ToplevelResizeEdgeTopLeft:     cursor.TopLeftCorner,
	xdgshell.ToplevelResizeEdgeTopRight:    cursor.TopRightCorner,
	xdgshell.ToplevelResizeEdgeBottom:      cursor.BottomSide,
	xdgshell.ToplevelResizeEdgeBottomLeft:  cursor.BottomLeftCorner,
	xdgshell.ToplevelResizeEdgeBottomRight: cursor.BottomRightCorner,
	xdgshell.ToplevelResizeEdgeLeft:        cursor.LeftSide,
	xdgshell.ToplevelResizeEdgeRight:       cursor.RightSide,
	xdgshell.ToplevelResizeEdgeNone:        cursor.LeftPtr,
}

func (app *appState) pointerFrameMotionEvent(e pointerEvent) {
	log.Printf("motion %v, %v", e.surfaceX, e.surfaceY)

	edge := componentEdge(uint32(app.width), uint32(app.height), e.surfaceX, e.surfaceY, Border)
	cursorName, ok := cursorMap[edge]
	if ok {
		app.trySetCursor(e.serial, cursorName)
	} else {
		println("cursor not in map")
	}

}
func (app *appState) pointerFrameAxisEvent(e pointerEvent) {
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

func (app *appState) pointerFrameButtonEvent() {
	e := &app.pointerEvent
	if wl.PointerButtonState(e.state) == wl.PointerButtonStateReleased {
		log.Printf("button %d released", e.button)
		
		if app.decoration != nil {
			app.decoration.LeftActive, app.decoration.RightActive = app.decoration.activeLeftRight(app, float64(e.surfaceX), float64(e.surfaceY))
			if app.decoration.RightActive == 1 {
				app.exit = true
			}
			
			if app.decoration.RightActive == 2 {
				if app.decoration.Maximized {
				
					if nil == app.xdgTopLevel.UnsetMaximized() {
						app.decoration.Maximized = false
					}
				} else {
					if nil == app.xdgTopLevel.SetMaximized() {
						app.decoration.Maximized = true
					}
				}
			}
			
			app.decoration.LeftActive, app.decoration.RightActive = 0, 0
			
			app.redecorate()
		}
		
	} else {
		log.Printf("button %d pressed", e.button)

		switch e.button {
		case BtnLeft:
		
			if app.decoration != nil {
				app.decoration.LeftActive, app.decoration.RightActive = app.decoration.activeLeftRight(app, float64(e.surfaceX), float64(e.surfaceY))
				if app.decoration.LeftActive != 0 || app.decoration.RightActive != 0 {
					app.redecorate()
					break
				}
			}
		
			edge := componentEdge(uint32(app.width), uint32(app.height), e.surfaceX, e.surfaceY, 8)
			if edge != xdgshell.ToplevelResizeEdgeNone {
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
func (app *appState) HandlePointerFrame(_ wl.PointerFrameEvent) {
	e := app.pointerEvent

	if (e.eventMask & pointerEventEnter) != 0 {
		app.pointerEvent.eventMask &= ^pointerEventEnter
		log.Printf("entered %v, %v", e.surfaceX, e.surfaceY)
		//app.trySetCursor(e.serial, cursor.LeftPtr)
	}

	if (e.eventMask & pointerEventLeave) != 0 {
		app.pointerEvent.eventMask &= ^pointerEventLeave
		log.Print("leave")
	}
	if (e.eventMask & pointerEventMotion) != 0 {
		app.pointerEvent.eventMask &= ^pointerEventMotion
		app.pointerFrameMotionEvent(e)
	}
	if (e.eventMask & pointerEventButton) != 0 {
		app.pointerEvent.eventMask &= ^pointerEventButton
		app.pointerFrameButtonEvent()
	}

	const axisEvents = pointerEventAxis | pointerEventAxisSource | pointerEventAxisStop | pointerEventAxisDiscrete

	if (e.eventMask & axisEvents) != 0 {
		app.pointerEvent.eventMask &= ^axisEvents
		app.pointerFrameAxisEvent(e)
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
			return xdgshell.ToplevelResizeEdgeTopLeft
		} else if right {
			log.Print("ToplevelResizeEdgeTopRight")
			return xdgshell.ToplevelResizeEdgeTopRight
		} else {
			log.Print("ToplevelResizeEdgeTop")
			return xdgshell.ToplevelResizeEdgeTop
		}
	} else if bottom {
		if left {
			log.Print("ToplevelResizeEdgeBottomLeft")
			return xdgshell.ToplevelResizeEdgeBottomLeft
		} else if right {
			log.Print("ToplevelResizeEdgeBottomRight")
			return xdgshell.ToplevelResizeEdgeBottomRight
		} else {
			log.Print("ToplevelResizeEdgeBottom")
			return xdgshell.ToplevelResizeEdgeBottom
		}
	} else if left {
		log.Print("ToplevelResizeEdgeLeft")
		return xdgshell.ToplevelResizeEdgeLeft
	} else if right {
		log.Print("ToplevelResizeEdgeRight")
		return xdgshell.ToplevelResizeEdgeRight
	} else {
		log.Print("ToplevelResizeEdgeNone")
		return xdgshell.ToplevelResizeEdgeNone
	}
}

type cursorData struct {
	wlCursor *cursor.Cursor

	surface *wl.Surface
}
func (app *appState) trySetCursor(serial uint32, cursorName string) {
	//if cursorName != app.currentCursor {
		print("SERIAL: ")
		println(serial)
		app.setCursor(serial, cursorName)
	//}
}
func (app *appState) setCursor(serial uint32, cursorName string) {
	c, ok := app.cursors[cursorName]
	if !ok {
		log.Printf("unable to get %v cursor", cursorName)
		return
	}

	image := c.wlCursor.Images[0]
	if err := app.pointer.SetCursor(
		serial, c.surface,
		int32(image.GetHotspotX()), int32(image.GetHotspotY()),
	); err != nil {
		log.Print("unable to set cursor")
	}
	c.surface.Attach(image.GetBuffer(), 0,0)
	c.surface.Damage(0,0,int32(image.GetWidth()),int32(image.GetHeight()))
	c.surface.Commit()

	app.currentCursor = cursorName
}
