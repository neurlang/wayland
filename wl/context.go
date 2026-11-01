package wl

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	sys "github.com/neurlang/wayland/os"
	//"reflect"
)

// Context wraps the wayland connection together with the map of all Context objects (proxies)
type Context struct {
	mu        sync.RWMutex
	conn      *net.UnixConn
	sockFD    int
	currentId ProxyId
	objects   map[ProxyId]Proxy
	scms      []sys.SocketControlMessage

	// The compositor requires client-allocated ids to arrive in the order
	// they were allocated: libwayland's wl_map rejects a new id that is
	// not the next free slot ("not a valid new object id"). Register hands
	// out ids in order, but the request that carries an id to the wire is
	// written later, by whichever goroutine created the proxy, so two
	// goroutines can write their requests in the wrong order. sendMu
	// serialises writes, unsent holds the ids Register has handed out
	// that no request has carried yet, and sendCond lets a writer wait
	// until every smaller id has gone.
	sendMu   sync.Mutex
	sendCond *sync.Cond
	unsent   map[ProxyId]struct{}
}

// markUnsent records a freshly allocated client id as not yet on the wire.
func (ctx *Context) markUnsent(id ProxyId) {
	ctx.sendMu.Lock()
	if ctx.unsent == nil {
		ctx.unsent = make(map[ProxyId]struct{})
		ctx.sendCond = sync.NewCond(&ctx.sendMu)
	}
	ctx.unsent[id] = struct{}{}
	ctx.sendMu.Unlock()
}

// forgetUnsent drops an id from the unsent set and wakes any writer that
// was waiting for it. Called once the id is on the wire, or when the proxy
// is unregistered without ever being sent, so nobody waits for it forever.
func (ctx *Context) forgetUnsent(id ProxyId) {
	ctx.sendMu.Lock()
	ctx.forgetUnsentLocked(id)
	ctx.sendMu.Unlock()
}

func (ctx *Context) forgetUnsentLocked(id ProxyId) {
	if ctx.unsent == nil {
		return
	}
	if _, ok := ctx.unsent[id]; !ok {
		return
	}
	delete(ctx.unsent, id)
	ctx.sendCond.Broadcast()
}

// smallerUnsentLocked reports whether some id below id is still unsent.
// Callers hold sendMu.
func (ctx *Context) smallerUnsentLocked(id ProxyId) bool {
	for other := range ctx.unsent {
		if other < id {
			return true
		}
	}
	return false
}

func (ctx *Context) RegisterMapped(proxy Proxy, num uint32) {
	ctx.mu.Lock()
	proxy.SetId(ProxyId(num))
	proxy.SetContext(ctx)
	ctx.objects[ProxyId(num)] = proxy
	ctx.mu.Unlock()
}

// Register registers a proxy in the map of all Context objects (proxies)
func (ctx *Context) Register(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	// avoid colliding with server side mapped proxy id?
	startId := ctx.currentId
	for {
		ctx.currentId += 1
		if _, ok := ctx.objects[ctx.currentId]; !ok {
			break
		}
		// Prevent infinite loop if all IDs are exhausted
		if ctx.currentId == startId {
			panic("proxy ID space exhausted")
		}
	}
	proxy.SetId(ctx.currentId)
	proxy.SetContext(ctx)
	if c, ok := proxy.(*Registry); ok {
		SetUserData(c, &ctx)
	}
	ctx.objects[ctx.currentId] = proxy
	// The display is the one client-side proxy the compositor already
	// knows: its id (1) is fixed by the protocol and no request ever
	// carries it as a new_id. Tracking it as unsent would leave every
	// later new id -- starting with wl_display.get_registry -- waiting
	// for a write that never comes.
	if _, isDisplay := proxy.(*Display); !isDisplay {
		ctx.markUnsent(ctx.currentId)
	}
}

// Unregister unregisters a proxy in the map of all Context objects (proxies)
func (ctx *Context) Unregister(id ProxyId) {
	ctx.mu.Lock()
	if ctx.objects != nil {
		delete(ctx.objects, id)
	}
	ctx.mu.Unlock()
	ctx.forgetUnsent(id)
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

// ErrXdgRuntimeDirNotSet is returned by Connect when the operating system does not provide the required
// XDG_RUNTIME_DIR environment variable
var ErrXdgRuntimeDirNotSet = errors.New("variable XDG_RUNTIME_DIR not set in the environment")

// Connect connects to a Wayland compositor running on a specific wayland unix socket.
// An empty addr means WAYLAND_DISPLAY, else wayland-0. As in libwayland's
// wl_display_connect, an absolute addr is the socket's path as is, and a
// relative one names a socket in XDG_RUNTIME_DIR.
func Connect(addr string) (ret *Display, err error) {
	addr, err = socketPath(addr)
	if err != nil {
		return nil, err
	}
	c := new(Context)
	c.objects = make(map[ProxyId]Proxy)
	c.currentId = 0
	c.conn, err = net.DialUnix("unix", nil, &net.UnixAddr{Name: addr, Net: "unix"})
	if err != nil {
		return nil, err
	}
	err = c.conn.SetReadDeadline(time.Time{})
	if err != nil {
		c.conn.Close()
		return nil, err
	}
	//DON'T dispatch events in separate goroutine
	//go c.Run()
	return NewDisplay(c), nil
}

// socketPath resolves the socket Connect dials for addr.
func socketPath(addr string) (string, error) {
	if addr == "" {
		addr = os.Getenv("WAYLAND_DISPLAY")
	}
	if addr == "" {
		addr = "wayland-0"
	}
	if filepath.IsAbs(addr) {
		return addr, nil
	}
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return "", ErrXdgRuntimeDirNotSet
	}
	return filepath.Join(runtimeDir, addr), nil
}

var errFoundMyCallback = errors.New("run found my callback")

// RunTill (Context RunTill) runs until a specific callback or an error occurs, see Context Run
// for a description of a likely errors
func (ctx *Context) RunTill(cb *Callback) (err error) {
	for {
		err = ctx.run(cb)
		if err == errFoundMyCallback {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// ErrContextRunEventReadingError (Context Run event reading error), use InternalError to get the underlying cause
var ErrContextRunEventReadingError = errors.New("event reading error")

// ErrContextRunConnectionClosed (Context Run connection closed)
var ErrContextRunConnectionClosed = errors.New("connection closed")

// ErrContextRunTimeout (Context Run timeout error)
var ErrContextRunTimeout = errors.New("timeout error")

// ErrContextRunProtocolError (Context Run protocol error), use InternalError to get the underlying cause
var ErrContextRunProtocolError = errors.New("protocol error")

// ErrContextRunNotDispatched (Context Run not dispatched)
var ErrContextRunNotDispatched = errors.New("not dispatched")

// ErrContextRunProxyNil (Context Run proxy nil)
var ErrContextRunProxyNil = errors.New("proxy nil")

// Run (Context Run) reads and processes one event, a specific ErrContextRunXXX error
// may be returned in case of failure
func (ctx *Context) Run() error {
	return ctx.run(nil)
}

// ErrContextNil (Error context is nil) occurs if the thread closes context and then
// it wants to run, another thread probably cannot close it safely
var ErrContextNil = errors.New("context is nil")

// ErrContextConnNil (Error context conn is nil) occurs if the thread closes context and then
// it wants to run, another thread probably cannot close it safely
var ErrContextConnNil = errors.New("context conn is nil")

func (ctx *Context) run(cb *Callback) error {
	// ctx := context.Background()

	if ctx == nil {
		return ErrContextNil
	}

	ev, err := ctx.readEvent()
	if err != nil {
		if err == io.EOF {
			return ErrContextRunConnectionClosed
		}

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return ErrContextRunTimeout
		}

		return combinedError{ErrContextRunEventReadingError, err}
	}

	proxy := ctx.LookupProxy(ev.Pid)
	if proxy != nil {
		if dispatcher, ok := proxy.(Dispatcher); dispatcher != nil && ok {
			if foundCb, ok := dispatcher.(*Callback); ok {
				if foundCb == cb {
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

// Close (Context Close) closes Wayland connection
func (ctx *Context) Close() (err error) {
	if ctx == nil {
		return
	}
	ctx.mu.Lock()
	if ctx.conn != nil {
		err = ctx.conn.Close()
		ctx.conn = nil
	}
	ctx.sockFD = -1
	/*
		for i, v := range ctx.objects {
			print("close-time garbage: ")
			print(i)
			print(": ")
			println(reflect.TypeOf(v).String())
		}
	*/
	ctx.objects = nil
	ctx.scms = nil
	ctx.mu.Unlock()
	CleanContextData(ctx)
	return err
}
