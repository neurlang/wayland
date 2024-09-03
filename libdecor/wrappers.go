package libdecor

import (
	"fmt"
	"runtime"
	"unsafe"
)

// LibdecorNew creates a new libdecor object.
func LibdecorNew(wldisplay uintptr, iface uintptr) (*Libdecor, error) {
	ret := libdecorNew(wldisplay, iface)
	if ret == 0 {
		return nil, fmt.Errorf("failed to create new libdecor: %d", ret)
	}
	libdecor := &Libdecor{ptr: unsafe.Pointer(ret)}
	runtime.SetFinalizer(libdecor, libdecorUnref)
	return libdecor, nil
}

// LibdecorUnref Releases a reference on a libdecor object, and possibly free it.
//
// Parameter the object.  If it is nil, this function does nothing.
func LibdecorUnref(l **Libdecor) {
	if l != nil {
		*l = nil
	}
}
func libdecorUnref(libdecor *Libdecor) {
	if libdecor == nil {
		return
	}
	libdecorUnrefP(uintptr(libdecor.ptr))
}

// Unref releases the libdecor frame object.
func LibdecorFrameUnref(f **LibdecorFrame) {
	if f != nil && *f != nil {
		*f = nil
	}
}

// SetTitle sets the title for the libdecor frame.
func (f *LibdecorFrame) SetTitle(title string) {
	if len(title) == 0 || title[len(title)-1] != 0 {
		title += "\x00"
	}
	libdecorFrameSetTitle(uintptr(f.ptr), title)
}

// Decorate decorates a surface and returns a new LibdecorFrame object.
func (l *Libdecor) Decorate(wlsurface uintptr, iface uintptr, data uintptr) (*LibdecorFrame, error) {
	frame := libdecorDecorate(uintptr(l.ptr), wlsurface, uintptr(iface), uintptr(data))
	if frame == 0 {
		return nil, fmt.Errorf("failed to decorate the surface")
	}
	libdecorFrame := &LibdecorFrame{ptr: unsafe.Pointer(frame)}
	runtime.SetFinalizer(libdecorFrame, libdecorFrameUnref)
	return libdecorFrame, nil
}

func (l *Libdecor) Dispatch(n int) int {
	return libdecorDispatch(uintptr(l.ptr), n)
}

// libdecorFrameUnref is the finalizer function for LibdecorFrame.
func libdecorFrameUnref(frame *LibdecorFrame) {
	if frame != nil {
		libdecorFrameUnrefP(uintptr(frame.ptr))
	}
}

// Map maps the libdecor frame.
func (f *LibdecorFrame) Map() {
	libdecorFrameMap(uintptr(f.ptr))
}

// SetAppID sets the application ID for the libdecor frame.
func (f *LibdecorFrame) SetAppID(appID string) {
	if len(appID) == 0 || appID[len(appID)-1] != 0 {
		appID += "\x00"
	}
	libdecorFrameSetAppId(uintptr(f.ptr), appID)
}

// Create a new libdecor state.
func LibdecorStateNew(width, height int) *LibdecorState {
	ptr := libdecorStateNew(width, height)
	if ptr == 0 {
		return nil
	}
	return &LibdecorState{ptr: unsafe.Pointer(ptr)}
}

// Free a libdecor state.
func LibdecorStateFree(s **LibdecorState) {
	if s != nil && *s != nil && (*s).ptr != nil {
		libdecorStateFree(uintptr((*s).ptr))
		*s = nil // Avoid double free
	}
}

// Commit a libdecor frame.
func (f *LibdecorFrame) Commit(state *LibdecorState, configuration uintptr) {
	if f != nil && f.ptr != nil {
		libdecorFrameCommit(uintptr(f.ptr), uintptr(unsafe.Pointer(state.ptr)), configuration)
	}
}

// Commit a libdecor frame.
func (f *LibdecorFrame) IsFloating() bool {
	if f != nil && f.ptr != nil {
		return libdecorFrameIsFloat(uintptr(f.ptr))
	}
	return false
}

func (f *LibdecorFrame) ConfigurationGetContentSize(conf uintptr, w *int, h *int) bool {
	if f != nil && f.ptr != nil {
		return libdecorConfGetContS(conf, uintptr(f.ptr), w, h)
	}
	return false
}
func (f *LibdecorFrame) SetVisibility(v bool) {
	if f != nil && f.ptr != nil {
		libdecorFrameSetVisib(uintptr(f.ptr), v)
	}
}

func CreateFrameInterface(f FrameInterface) (i *LibdecorFrameInterface) {
	i = &LibdecorFrameInterface{}
	i[0] = uintptr(NewCallback(f.Configure))
	i[1] = uintptr(NewCallback(f.Close))
	i[2] = uintptr(NewCallback(f.Commit))

	return
}
