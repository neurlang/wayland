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

type Context struct {
	mu           sync.RWMutex
	conn         *net.UnixConn
	sockFD       int
	currentId    ProxyId
	objects      map[ProxyId]Proxy
	dispatchChan chan struct{}
	exitChan     chan struct{}
}

func (ctx *Context) Register(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.currentId += 1
	proxy.SetId(ctx.currentId)
	proxy.SetContext(ctx)
	ctx.objects[ctx.currentId] = proxy
}

func (ctx *Context) LookupProxy(id ProxyId) Proxy {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	proxy, ok := ctx.objects[id]
	if !ok {
		return nil
	}
	return proxy
}

func Connect(addr string) (ret *Display, err error) {
	runtime_dir := os.Getenv("XDG_RUNTIME_DIR")
	if runtime_dir == "" {
		return nil, errors.New("variable XDG_RUNTIME_DIR not set in the environment")
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
	c.dispatchChan = make(chan struct{})
	c.exitChan = make(chan struct{})
	c.conn, err = net.DialUnix("unix", nil, &net.UnixAddr{Name: addr, Net: "unix"})
	if err != nil {
		return nil, err
	}
	c.conn.SetReadDeadline(time.Time{})
	//dispatch events in separate gorutine
	//go c.Run()
	return NewDisplay(c), nil
}

var errFoundMyCallback = errors.New("run found my callback")

func (c *Context) RunTill(cb *Callback) (err error) {
	for {
		err = c.Run(cb)
		if err == errFoundMyCallback {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Context) Run(cb *Callback) error {
	// ctx := context.Background()

	ev, err := c.readEvent()
	if err != nil {
		if err == io.EOF {
			return errors.New("connection closed")
		}

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return errors.New("timeout error")
		}

		log.Fatal(err)
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
		} else {
			return errors.New("not dispatched")
		}
	} else {
		return errors.New("proxy NULL")
	}
	return nil
}
