package xkbcommon

/*
#cgo pkg-config: xkbcommon
#cgo LDFLAGS: -ldl

#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"

// Context is an Opaque top level library context object.
type Context struct {
	cx *C.struct_xkb_context
}

// Keymap is an Opaque compiled keymap object.
type Keymap struct {
	km *C.struct_xkb_keymap
}

// State is an Opaque keyboard state object.
type State struct {
	st *C.struct_xkb_state
}

// ComposeState is an Opaque Compose state object.
type ComposeState struct {
	cs *C.struct_xkb_compose_state
}

// ComposeTable is an Opaque Compose table object.
type ComposeTable struct {
	ct *C.struct_xkb_compose_table
}
