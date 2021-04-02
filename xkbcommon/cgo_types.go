package xkbcommon

/*
#cgo pkg-config: xkbcommon
#cgo LDFLAGS: -ldl

#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"

type Context struct {
	cx *C.struct_xkb_context
}
type Keymap struct {
	km *C.struct_xkb_keymap
}
type State struct {
	st *C.struct_xkb_state
}
type ComposeState struct {
	cs *C.struct_xkb_compose_state
}
type ComposeTable struct {
	ct *C.struct_xkb_compose_table
}
