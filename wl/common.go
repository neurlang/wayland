// Package wl implements the stable Wayland protocol
package wl

// Proxy identifier that is sent to compositor over the wayland socket
type ProxyId uint32

// Dispatcher is anything that can process an Event
type Dispatcher interface {
	Dispatch(*Event)
}

// Proxy object
type Proxy interface {
	Context() *Context
	SetContext(c *Context)
	Id() ProxyId
	SetId(id ProxyId)
}

// Base Proxy is a struct that stores Context and ProxyId explicitly
type BaseProxy struct {
	id  ProxyId
	ctx *Context
}

// BaseProxy implements Id to get ProxyId
func (p *BaseProxy) Id() ProxyId {
	return p.id
}

// BaseProxy implements SetId to set ProxyId
func (p *BaseProxy) SetId(id ProxyId) {
	p.id = id
}

// BaseProxy implements Context to get Context
func (p *BaseProxy) Context() *Context {
	return p.ctx
}

// BaseProxy implements SetContext to set Context
func (p *BaseProxy) SetContext(c *Context) {
	p.ctx = c
}

// BaseProxy Unregister removes this BaseProxy from the map of all Context objects
func (p *BaseProxy) Unregister(s string) {
	if p.ctx != nil {
		// fmt.Println("Removing object", p.id, s)
		delete(p.ctx.objects, p.id)
	}
}
