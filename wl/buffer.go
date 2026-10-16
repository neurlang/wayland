package wl

// NewServerBuffer returns an unregistered Buffer for a wl_buffer the
// compositor creates in an event, such as zwp_linux_buffer_params_v1.created.
func NewServerBuffer() *Buffer {
	ret := new(Buffer)
	ret.initBuffer()
	return ret
}
