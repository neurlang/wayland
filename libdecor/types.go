package libdecor

import (
	"unsafe"
)

type Libdecor struct {
	ptr unsafe.Pointer
}

type LibdecorFrame struct {
	ptr unsafe.Pointer
}

type LibdecorConfiguration struct {
	ptr unsafe.Pointer
}

type LibdecorState struct {
	ptr unsafe.Pointer
}

type LibdecorError int

type FrameInterface struct {
	Configure func(_, configuration, cube uintptr)
	Close     func(_, userData uintptr)
	Commit    func(_, userData uintptr)
}

type LibdecorFrameInterface [14]uintptr

const (
	LIBDECOR_ERROR_COMPOSITOR_INCOMPATIBLE LibdecorError = iota
	LIBDECOR_ERROR_INVALID_FRAME_CONFIGURATION
)

type LibdecorWindowState int

const (
	LIBDECOR_WINDOW_STATE_NONE         LibdecorWindowState = 0
	LIBDECOR_WINDOW_STATE_ACTIVE       LibdecorWindowState = 1 << 0
	LIBDECOR_WINDOW_STATE_MAXIMIZED    LibdecorWindowState = 1 << 1
	LIBDECOR_WINDOW_STATE_FULLSCREEN   LibdecorWindowState = 1 << 2
	LIBDECOR_WINDOW_STATE_TILED_LEFT   LibdecorWindowState = 1 << 3
	LIBDECOR_WINDOW_STATE_TILED_RIGHT  LibdecorWindowState = 1 << 4
	LIBDECOR_WINDOW_STATE_TILED_TOP    LibdecorWindowState = 1 << 5
	LIBDECOR_WINDOW_STATE_TILED_BOTTOM LibdecorWindowState = 1 << 6
	LIBDECOR_WINDOW_STATE_SUSPENDED    LibdecorWindowState = 1 << 7
	LIBDECOR_WINDOW_STATE_RESIZING     LibdecorWindowState = 1 << 8
)

// Other enums and constants are defined similarly...
