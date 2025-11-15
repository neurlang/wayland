package linux

import (
	"sync"

	"github.com/neurlang/wayland/wl"
)

// Type aliases for Wayland core types used by linux-dmabuf protocol
type BaseProxy = wl.BaseProxy
type Context = wl.Context
type Event = wl.Event
type Surface = wl.Surface

// Buffer is a wrapper around wl.Buffer for use in linux-dmabuf protocol
type Buffer struct {
	wl.BaseProxy
	mu                    sync.RWMutex
	privateBufferReleases map[wl.BufferReleaseHandler]struct{}
}

// initBuffer initializes the Buffer object's handler maps
func (ret *Buffer) initBuffer() {
	ret.privateBufferReleases = make(map[wl.BufferReleaseHandler]struct{})
}

// NewBuffer creates a new Buffer object
func NewBuffer(ctx *Context) *Buffer {
	ret := new(Buffer)
	ret.initBuffer()
	ctx.Register(ret)
	return ret
}

// Destroy destroys the buffer
func (p *Buffer) Destroy() error {
	return p.Context().SendRequest(p, 0)
}

// Dispatch dispatches events for Buffer
func (p *Buffer) Dispatch(event *Event) {
	switch event.Opcode {
	case 0:
		if len(p.privateBufferReleases) > 0 {
			ev := wl.BufferReleaseEvent{}
			p.mu.RLock()
			for h := range p.privateBufferReleases {
				h.HandleBufferRelease(ev)
			}
			p.mu.RUnlock()
		}
	}
}

// AddReleaseHandler adds a Release event handler
func (p *Buffer) AddReleaseHandler(h wl.BufferReleaseHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.privateBufferReleases[h] = struct{}{}
}

// RemoveReleaseHandler removes a Release event handler
func (p *Buffer) RemoveReleaseHandler(h wl.BufferReleaseHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.privateBufferReleases, h)
}
