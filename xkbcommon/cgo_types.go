package xkbcommon

/*
#cgo pkg-config: xkbcommon
#cgo LDFLAGS: -ldl

#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"

type Context = C.struct_xkb_context
type Keymap = C.struct_xkb_keymap
type State = C.struct_xkb_state
type ComposeState = C.struct_xkb_compose_state
type ComposeTable = C.struct_xkb_compose_table
