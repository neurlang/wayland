package libdecor

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

func NewCallback(fn interface{}) uintptr {
	return purego.NewCallback(fn)
}

// getSystemLibrary returns the name of the system library based on OS.
func getSystemLibrary() string {
	switch runtime.GOOS {
	case "linux":
		return "libdecor-0.so"
	default:
		panic(fmt.Sprintf("GOOS=%s is not supported", runtime.GOOS))
	}
}

// getSystemLibraryDotVersion returns the library version suffix if needed.
func getSystemLibraryDotVersion() string {
	switch runtime.GOOS {
	case "linux":
		return ".0" // Adjust if specific version is required.
	default:
		panic(fmt.Sprintf("GOOS=%s is not supported", runtime.GOOS))
	}
}

// Function types for libdecor functions.
var (
	libdecorNew           func(uintptr, uintptr) uintptr
	libdecorUnrefP        func(uintptr)
	libdecorFrameUnrefP   func(uintptr)
	libdecorFrameSetTitle func(uintptr, string)
	libdecorDecorate      func(uintptr, uintptr, uintptr, uintptr) uintptr
	libdecorFrameMap      func(uintptr)
	libdecorFrameSetAppId func(uintptr, string)
	libdecorDispatch      func(uintptr, int) int
	libdecorStateNew      func(int, int) uintptr
	libdecorFrameCommit   func(uintptr, uintptr, uintptr)
	libdecorStateFree     func(uintptr)
	libdecorFrameIsFloat  func(uintptr) bool
	libdecorConfGetContS  func(uintptr, uintptr, *int, *int) bool
	libdecorFrameSetVisib func(uintptr, bool)

	available bool
)

func Available() bool {
	return available
}

func init() {
	// Load the library
	libdecorLib, err := purego.Dlopen(getSystemLibrary(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		libdecorLib, err = purego.Dlopen(getSystemLibrary()+getSystemLibraryDotVersion(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			// Library not available, leave available = false
			return
		}
	}

	// Register library functions - if any fail, library is not fully available
	defer func() {
		if r := recover(); r != nil {
			// RegisterLibFunc panicked, library not fully available
			available = false
		}
	}()

	purego.RegisterLibFunc(&libdecorNew, libdecorLib, "libdecor_new")
	purego.RegisterLibFunc(&libdecorUnrefP, libdecorLib, "libdecor_unref")
	purego.RegisterLibFunc(&libdecorFrameUnrefP, libdecorLib, "libdecor_frame_unref")
	purego.RegisterLibFunc(&libdecorFrameSetTitle, libdecorLib, "libdecor_frame_set_title")
	purego.RegisterLibFunc(&libdecorDecorate, libdecorLib, "libdecor_decorate")
	purego.RegisterLibFunc(&libdecorFrameMap, libdecorLib, "libdecor_frame_map")
	purego.RegisterLibFunc(&libdecorFrameSetAppId, libdecorLib, "libdecor_frame_set_app_id")
	purego.RegisterLibFunc(&libdecorDispatch, libdecorLib, "libdecor_dispatch")
	purego.RegisterLibFunc(&libdecorStateNew, libdecorLib, "libdecor_state_new")
	purego.RegisterLibFunc(&libdecorFrameCommit, libdecorLib, "libdecor_frame_commit")
	purego.RegisterLibFunc(&libdecorStateFree, libdecorLib, "libdecor_state_free")
	purego.RegisterLibFunc(&libdecorFrameIsFloat, libdecorLib, "libdecor_frame_is_floating")
	purego.RegisterLibFunc(&libdecorConfGetContS, libdecorLib, "libdecor_configuration_get_content_size")
	purego.RegisterLibFunc(&libdecorFrameSetVisib, libdecorLib, "libdecor_frame_set_visibility")

	available = true
}
