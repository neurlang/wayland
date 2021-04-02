package xkbcommon

func KeymapNewFromString(*struct{}, []byte, int, int) *Keymap {
	return new(Keymap)
}

func ComposeTableNewFromLocale(*struct{}, string, int) *ComposeTable {
	return new(ComposeTable)
}
func StateNew(*Keymap) *State {
	return new(State)
}
func ComposeStateNew(*ComposeTable, int) *ComposeState {
	return new(ComposeState)
}
func KeymapUnref(*Keymap) {

}
func ComposeTableUnref(*ComposeTable) {

}
func ComposeStateUnref(*ComposeState) {

}
func StateUnref(*State) {

}

func KeymapModGetIndex(*Keymap, string) uint {
	return 0
}

func StateKeyGetSyms(state *State, code uint32) []uint32 {
	return []uint32{code}
}

func KeymapKeyRepeats(keymap *Keymap, code uint32) bool {
	return true
}

func ComposeStateGetStatus(*ComposeState) uint32 {
	return 0
}
func ComposeStateGetOneSym(*ComposeState) uint32 {
	return 0
}
func ComposeStateFeed(*ComposeState, uint32) uint32 {
	return 0
}
