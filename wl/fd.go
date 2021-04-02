package wl

func (e *KeyboardKeymapEvent) Fd() (uintptr, error) {
	return e.fd, e.fdError
}
