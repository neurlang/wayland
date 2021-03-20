package wl

type ProxyId uint32

type Dispatcher interface {
	Dispatch(*Event)
	// Dispatch(*Event)
}

type Proxy interface {
	Context() *Context
	SetContext(c *Context)
	Id() ProxyId
	SetId(id ProxyId)
}

type BaseProxy struct {
	id        ProxyId
	version   uint32
	ctx       *Context
	container interface{}
}

func (p *BaseProxy) Id() ProxyId {
	return p.id
}

func (p *BaseProxy) SetId(id ProxyId) {
	p.id = id
}

func (p *BaseProxy) Version() uint32 {
	return p.version
}

func (p *BaseProxy) SetVersion(version uint32) {
	p.version = version
}

func (p *BaseProxy) Context() *Context {
	return p.ctx
}

func (p *BaseProxy) SetContext(c *Context) {
	p.ctx = c
}

func (p *BaseProxy) Container() interface{} {
	return p.container
}

func (p *BaseProxy) SetContainer(c interface{}) {
	p.container = c
}

func (p *BaseProxy) Unregister(s string) {
	if p.ctx != nil {
		// fmt.Println("Removing object", p.id, s)
		delete(p.ctx.objects, p.id)
	}
}

type Handler interface {
	Handle(ev interface{})
}

type eventHandler struct {
	f func(interface{})
}

func HandlerFunc(f func(interface{})) Handler {
	return &eventHandler{f}
}

func (h *eventHandler) Handle(ev interface{}) {
	h.f(ev)
}
