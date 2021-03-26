package wl

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func init() {
	log.SetFlags(0)
}

// Context wraps the wayland connection together with the map of all Context objects (proxies)
type Context struct {
	mu        sync.RWMutex
	conn      *net.UnixConn
	sockFD    int
	currentId ProxyId
	objects   map[ProxyId]Proxy
}

// Register registers a proxy in the map of all Context objects (proxies)
func (ctx *Context) Register(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.currentId += 1
	proxy.SetId(ctx.currentId)
	proxy.SetContext(ctx)
	ctx.objects[ctx.currentId] = proxy
}

// LookupProxy looks up a specific proxy by it's Id in the map of all Context objects (proxies)
func (ctx *Context) LookupProxy(id ProxyId) Proxy {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	proxy, ok := ctx.objects[id]
	if !ok {
		return nil
	}
	return proxy
}

// This error is returned by Connect when the operating system does not provide the required XDG_RUNTIME_DIR environment variable
var ErrXdgRuntimeDirNotSet = errors.New("variable XDG_RUNTIME_DIR not set in the environment")

// Connect connects to a Wayland compositor running on a specific wayland unix socket
func Connect(addr string) (ret *Display, err error) {
	runtime_dir := os.Getenv("XDG_RUNTIME_DIR")
	if runtime_dir == "" {
		return nil, ErrXdgRuntimeDirNotSet
	}
	if addr == "" {
		addr = os.Getenv("WAYLAND_DISPLAY")
	}
	if addr == "" {
		addr = "wayland-0"
	}
	addr = runtime_dir + "/" + addr
	c := new(Context)
	c.objects = make(map[ProxyId]Proxy)
	c.currentId = 0
	c.conn, err = net.DialUnix("unix", nil, &net.UnixAddr{Name: addr, Net: "unix"})
	if err != nil {
		return nil, err
	}
	c.conn.SetReadDeadline(time.Time{})
	//DON'T dispatch events in separate gorutine
	//go c.Run()
	return NewDisplay(c), nil
}

var errFoundMyCallback = errors.New("run found my callback")

// Context RunTill runs until a specific callback or an error occurs, see Context Run
// for a description of a likely errors
func (c *Context) RunTill(cb *Callback) (err error) {
	for {
		err = c.run(cb)
		if err == errFoundMyCallback {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Context Run event reading error, use InternalError to get the underlying cause
var ErrContextRunEventReadingError = errors.New("event reading error")

// Context Run connection closed
var ErrContextRunConnectionClosed = errors.New("connection closed")

// Context Run timeout error
var ErrContextRunTimeout = errors.New("timeout error")

// Context Run protocol error, use InternalError to get the underlying cause
var ErrContextRunProtocolError = errors.New("protocol error")

// Context Run not dispatched
var ErrContextRunNotDispatched = errors.New("not dispatched")

// Context Run proxy nil
var ErrContextRunProxyNil = errors.New("proxy nil")

// Context Run reads and processes one event, a specific ErrContextRunXXX error
// may be returned in case of failure
func (c *Context) Run() error {
	return c.run(nil)
}

func (c *Context) run(cb *Callback) error {
	// ctx := context.Background()

	ev, err := c.readEvent()
	if err != nil {
		if err == io.EOF {
			return ErrContextRunConnectionClosed
		}

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return ErrContextRunTimeout
		}

		return combinedError{ErrContextRunEventReadingError, err}
	}

	proxy := c.LookupProxy(ev.Pid)
	if proxy != nil {
		if dispatcher, ok := proxy.(Dispatcher); ok {
			if found_cb, ok := dispatcher.(*Callback); dispatcher != nil && ok {
				if found_cb == cb {
					bytePool.Give(ev.Data)
					return errFoundMyCallback
				}
			}
			dispatcher.Dispatch(ev)
			bytePool.Give(ev.Data)

			if ev.err != nil {
				return combinedError{ErrContextRunProtocolError, ev.err}
			}
		} else {
			return ErrContextRunNotDispatched
		}
	} else {
		return ErrContextRunProxyNil
	}
	return nil
}
