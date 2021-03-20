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
	id  ProxyId
	ctx *Context
}

func (p *BaseProxy) Id() ProxyId {
	return p.id
}

func (p *BaseProxy) SetId(id ProxyId) {
	p.id = id
}

func (p *BaseProxy) Context() *Context {
	return p.ctx
}

func (p *BaseProxy) SetContext(c *Context) {
	p.ctx = c
}

func (p *BaseProxy) Unregister(s string) {
	if p.ctx != nil {
		// fmt.Println("Removing object", p.id, s)
		delete(p.ctx.objects, p.id)
	}
}
