package main

import (
	"image"
	"os"

	syscall "golang.org/x/sys/unix"

	"github.com/neurlang/wayland/internal/swizzle"
	sys "github.com/neurlang/wayland/os"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
	"github.com/neurlang/wayland/wlcursor"
	"github.com/neurlang/wayland/xdg"
	"github.com/nfnt/resize"

	"log"
)

// Global app state
type appState struct {
	appID         string
	title         string
	pImage        *image.RGBA
	width, height int32
	frame         *image.RGBA
	exit          bool

	display    *wl.Display
	registry   *wl.Registry
	shm        *wl.Shm
	compositor *wl.Compositor
	wmBase     *xdg.WmBase
	seat       *wl.Seat

	surface     *wl.Surface
	xdgSurface  *xdg.Surface
	xdgTopLevel *xdg.Toplevel

	keyboard *wl.Keyboard
	pointer  *wl.Pointer

	pointerEvent  pointerEvent
	cursorTheme   *wlcursor.Theme
	cursors       map[string]*cursorData
	currentCursor string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s file.jpg", os.Args[0])
	}

	const (
		clampedWidth  = 1920
		clampedHeight = 1080
	)

	fileName := os.Args[1]

	// Read the image file to *image.RGBA
	pImage, err := rgbaImageFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// Create a proxy image for large images, makes resizing a little better
	if pImage.Rect.Dy() > pImage.Rect.Dx() && pImage.Rect.Dy() > clampedHeight {
		pImage = resize.Resize(0, clampedHeight, pImage, resize.Bilinear).(*image.RGBA)
		log.Print("creating proxy image, resizing by height clamped to", clampedHeight)
	} else if pImage.Rect.Dx() > pImage.Rect.Dy() && pImage.Rect.Dx() > clampedWidth {
		pImage = resize.Resize(clampedWidth, 0, pImage, resize.Bilinear).(*image.RGBA)
		log.Print("creating proxy image, resizing by width clamped to", clampedWidth)
	}

	// Resize again, for first frame
	frameImage := resize.Resize(0, 480, pImage, resize.Bilinear).(*image.RGBA)
	frameRect := frameImage.Bounds()

	app := &appState{
		// Set the title to `cat.jpg - imageview`
		title: fileName + " - imageviewer",
		appID: "imageviewer",
		// Keep proxy image in cache, for use in resizing
		pImage:  pImage,
		width:   int32(frameRect.Dx()),
		height:  int32(frameRect.Dy()),
		frame:   frameImage,
		cursors: make(map[string]*cursorData),
	}

	// Connect to wayland server
	display, err := wl.Connect("")
	if err != nil {
		log.Fatalf("unable to connect to wayland server: %v", err)
	}
	app.display = display

	display.AddErrorHandler(app)

	// Start other stuff in function for simplicity
	run(app)

	log.Println("closing")

	// Release the pointer if registered
	if app.pointer != nil {
		app.releasePointer()
	}

	// Release the keyboard if registered
	if app.keyboard != nil {
		app.releaseKeyboard()
	}

	// Release wl_seat handlers
	if app.seat != nil {
		app.seat.RemoveCapabilitiesHandler(app)
		app.seat.RemoveNameHandler(app)

		if err := app.seat.Release(); err != nil {
			log.Println("unable to destroy wl_seat:", err)
		}
		app.seat = nil
	}

	// Release xdg_wmbase
	if app.wmBase != nil {
		app.wmBase.RemovePingHandler(app)

		if err := app.wmBase.Destroy(); err != nil {
			log.Println("unable to destroy xdg_wm_base:", err)
		}
		app.wmBase = nil
	}

	for _, c := range app.cursors {
		if err := c.wlCursor.Destroy(); err != nil {
			log.Println("unable to destroy cursor", c.wlCursor.Name, ":", err)
		}

		if err := c.surface.Destroy(); err != nil {
			log.Println("unable to destroy wl_surface of cursor", c.wlCursor.Name, ":", err)
		}
	}

	if app.cursorTheme != nil {
		if err := app.cursorTheme.Destroy(); err != nil {
			log.Println("unable to destroy cursor theme:", err)
		}
	}

	// Close the wayland server connection
	app.Context().Close()
}

func run(app *appState) {
	// Get global interfaces registry
	registry, err := app.display.GetRegistry()
	if err != nil {
		log.Fatalf("unable to get global registry object: %v", err)
	}
	app.registry = registry

	// Add global interfaces registrar handler
	registry.AddGlobalHandler(app)
	// Wait for interfaces to register
	wlclient.DisplayRoundtrip(app.display)

	log.Print("all interfaces registered")

	// Create a wl_surface for toplevel window
	surface, err := app.compositor.CreateSurface()
	if err != nil {
		log.Fatalf("unable to create compositor surface: %v", err)
	}
	app.surface = surface
	log.Print("created new wl_surface")

	// attach wl_surface to xdg_wmbase to get toplevel
	// handle
	xdgSurface, err := app.wmBase.GetSurface(surface)
	if err != nil {
		log.Fatalf("unable to get xdg_surface: %v", err)
	}
	app.xdgSurface = xdgSurface
	log.Print("got xdg_surface")

	// Add xdg_surface configure handler `app.HandleSurfaceConfigure`
	xdgSurface.AddConfigureHandler(app)
	log.Print("added configure handler")

	// Get toplevel
	xdgTopLevel, err := xdgSurface.GetToplevel()
	if err != nil {
		log.Fatalf("unable to get xdg_toplevel: %v", err)
	}
	app.xdgTopLevel = xdgTopLevel
	log.Print("got xdg_toplevel")

	// Add xdg_toplevel configure handler for window resizing
	xdgTopLevel.AddConfigureHandler(app)
	// Add xdg_toplevel close handler
	xdgTopLevel.AddCloseHandler(app)

	// Set title
	if err2 := xdgTopLevel.SetTitle(app.title); err2 != nil {
		log.Fatalf("unable to set toplevel title: %v", err2)
	}
	// Set appID
	if err2 := xdgTopLevel.SetAppID(app.appID); err2 != nil {
		log.Fatalf("unable to set toplevel appID: %v", err2)
	}
	// Commit the state changes (title & appID) to the server
	if err2 := app.surface.Commit(); err2 != nil {
		log.Fatalf("unable to commit surface state: %v", err2)
	}

	// Preload required cursors
	app.loadCursors()

	// Start the dispatch loop
	for !app.exit {
		err := app.Context().Run()
		if err != nil {
			log.Fatalf("error when running: %v", err)
		}
	}

}

func (app *appState) Context() *wl.Context {
	return app.display.Context()
}

func (app *appState) HandleRegistryGlobal(e wl.RegistryGlobalEvent) {
	log.Printf("discovered an interface: %q\n", e.Interface)

	switch e.Interface {
	case "wl_shm":
		shm := wl.NewShm(app.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, shm)
		if err != nil {
			log.Fatalf("unable to bind wl_shm interface: %v", err)
		}
		app.shm = shm
	case "wl_compositor":
		compositor := wl.NewCompositor(app.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, compositor)
		if err != nil {
			log.Fatalf("unable to bind wl_compositor interface: %v", err)
		}
		app.compositor = compositor
	case "xdg_wm_base":
		wmBase := xdg.NewWmBase(app.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, wmBase)
		if err != nil {
			log.Fatalf("unable to bind xdg_wm_base interface: %v", err)
		}
		app.wmBase = wmBase
		// Add xdg_wmbase ping handler `app.HandleWmBasePing`
		wmBase.AddPingHandler(app)
	case "wl_seat":
		seat := wl.NewSeat(app.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, seat)
		if err != nil {
			log.Fatalf("unable to bind wl_seat interface: %v", err)
		}
		app.seat = seat
		// Add Keyboard & Pointer handlers
		seat.AddCapabilitiesHandler(app)
		seat.AddNameHandler(app)
	}
}

func (app *appState) HandleSurfaceConfigure(e xdg.SurfaceConfigureEvent) {
	// Send ack to xdg_surface that we have a frame.
	if err := app.xdgSurface.AckConfigure(e.Serial); err != nil {
		log.Fatal("unable to ack xdg surface configure")
	}

	// Draw frame
	buffer := app.drawFrame()

	// Attach new frame to the surface
	if err := app.surface.Attach(buffer, 0, 0); err != nil {
		log.Fatalf("unable to attach buffer to surface: %v", err)
	}
	// Commit the surface state
	if err := app.surface.Commit(); err != nil {
		log.Fatalf("unable to commit surface state: %v", err)
	}
}

func (app *appState) HandleToplevelConfigure(e xdg.ToplevelConfigureEvent) {
	width := e.Width
	height := e.Height

	if width == 0 || height == 0 {
		// Compositor is deferring to us
		return
	}

	if width == app.width && height == app.height {
		// No need to resize
		return
	}

	// Resize the proxy image to new frame size
	// and set it to frame image
	log.Print("resizing frame")
	app.frame = resize.Resize(uint(width), uint(height), app.pImage, resize.Bilinear).(*image.RGBA)
	log.Print("done resizing frame")

	// Update app size
	app.width = width
	app.height = height
}

func (app *appState) loadCursors() {
	// Load default cursor theme
	theme, err := wlcursor.LoadTheme(24, app.shm)
	if err != nil {
		log.Fatalf("unable to load cursor theme: %v", err)
	}
	app.cursorTheme = theme

	// Create
	for _, name := range []string{
		wlcursor.BottomLeftCorner,
		wlcursor.BottomRightCorner,
		wlcursor.BottomSide,
		wlcursor.LeftPtr,
		wlcursor.LeftSide,
		wlcursor.RightSide,
		wlcursor.TopLeftCorner,
		wlcursor.TopRightCorner,
		wlcursor.TopSide,
	} {
		// Get wl_cursor
		c, err := theme.GetCursor(name)
		if err != nil {
			log.Fatalf("unable to get %v cursor: %v", name, err)
		}

		// Create a wl_surface for cursor
		surface, err := app.compositor.CreateSurface()
		if err != nil {
			log.Fatalf("unable to create compositor surface: %v", err)
		}
		log.Print("created new wl_surface for cursor: ", c.Name)

		// For now get the first image (there are multiple images because of animated cursors)
		// will figure out cursor animation afterwards
		image := c.Images[0]

		// Attach the first image to wl_surface
		if err := surface.Attach(image.GetBuffer(), 0, 0); err != nil {
			log.Fatalf("unable to attach cursor image buffer to cursor suface: %v", err)
		}
		// Commit the surface state changes
		if err2 := surface.Commit(); err2 != nil {
			log.Fatalf("unable to commit surface state: %v", err2)
		}

		// Store the surface for later (immediate) use
		app.cursors[name] = &cursorData{
			wlCursor: c,
			surface:  surface,
		}
	}
}

func (app *appState) drawFrame() *wl.Buffer {
	log.Print("drawing frame")

	stride := app.width * 4
	size := stride * app.height

	file, err := sys.CreateAnonymousFile(int64(size))
	if err != nil {
		log.Fatalf("unable to create a temporary file: %v", err)
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("unable to create mapping: %v", err)
	}

	pool, err := app.shm.CreatePool(file.Fd(), size)
	if err != nil {
		log.Fatalf("unable to create shm pool: %v", err)
	}

	buf, err := pool.CreateBuffer(0, app.width, app.height, stride, wl.ShmFormatArgb8888)
	if err != nil {
		log.Fatalf("unable to create wlclient.Buffer from shm pool: %v", err)
	}
	if err := pool.Destroy(); err != nil {
		log.Printf("unable to destroy shm pool: %v", err)
	}
	if err := file.Close(); err != nil {
		log.Printf("unable to close file: %v", err)
	}

	// Convert RGBA to BGRA
	imgData := make([]byte, len(app.frame.Pix))
	copy(imgData, app.frame.Pix)
	swizzle.BGRA(imgData)

	// Copy image to buffer
	copy(data, imgData)

	if err := syscall.Munmap(data); err != nil {
		log.Printf("unable to delete mapping: %v", err)
	}
	buf.AddReleaseHandler(bufferReleaser{buf: buf})

	log.Print("drawing frame complete")
	return buf
}

type bufferReleaser struct {
	buf *wl.Buffer
}

func (b bufferReleaser) HandleBufferRelease(e wl.BufferReleaseEvent) {
	if err := b.buf.Destroy(); err != nil {
		log.Printf("unable to destroy buffer: %v", err)
	}
}

func (app *appState) HandleSeatCapabilities(e wl.SeatCapabilitiesEvent) {
	havePointer := (e.Capabilities * wl.SeatCapabilityPointer) != 0

	if havePointer && app.pointer == nil {
		app.attachPointer()
	} else if !havePointer && app.pointer != nil {
		app.releasePointer()
	}

	haveKeyboard := (e.Capabilities * wl.SeatCapabilityKeyboard) != 0

	if haveKeyboard && app.keyboard == nil {
		app.attachKeyboard()
	} else if !haveKeyboard && app.keyboard != nil {
		app.releaseKeyboard()
	}
}

func (*appState) HandleSeatName(e wl.SeatNameEvent) {
	log.Printf("seat name: %v", e.Name)
}

// HandleDisplayError handles wlclient.Display errors
func (*appState) HandleDisplayError(e wl.DisplayErrorEvent) {
	// Just log.Fatal for now
	log.Fatalf("display error event: %v", e)
}

// HandleWmBasePing handles xdg ping by doing a Pong request
func (app *appState) HandleWmBasePing(e xdg.WmBasePingEvent) {
	log.Printf("xdg_wmbase ping: serial=%v", e.Serial)
	app.wmBase.Pong(e.Serial)
	log.Print("xdg_wmbase pong sent")
}

func (app *appState) HandleToplevelClose(_ xdg.ToplevelCloseEvent) {
	app.exit = true
}
