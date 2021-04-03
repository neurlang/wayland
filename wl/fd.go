package wl

// Fd (KeyboardKeymapEvent Fd) is used internally to retrieve the file descriptor from an event
func (e *KeyboardKeymapEvent) Fd() (uintptr, error) {
	return e.fd, e.fdError
}
