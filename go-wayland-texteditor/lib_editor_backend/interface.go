package lib_editor_backend

type Interface interface {
	Call(proc, payload string) string
}

type iface struct {}

func (i *iface) Call(proc, data string) string {
	if proc == "/content" {
		return handleContent(data)
	} else {
		return handleScrollbar(proc)
	}
}

func New() (*iface, error) {
	return &iface{}, nil
}


