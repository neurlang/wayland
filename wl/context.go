package wl

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

func init() {
	log.SetFlags(0)
}

type Context struct {
	mu           sync.RWMutex
	conn         *net.UnixConn
	SockFD       int
	currentId    ProxyId
	objects      map[ProxyId]Proxy
	dispatchChan chan struct{}
	exitChan     chan struct{}
	fds          []uintptr
}

func (ctx *Context) Register(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.currentId += 1
	proxy.SetId(ctx.currentId)
	proxy.SetContext(ctx)
	ctx.objects[ctx.currentId] = proxy
}

func (ctx *Context) RegisterId(proxy Proxy, id int) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.currentId += 1
	proxy.SetId(ProxyId(id))
	proxy.SetContext(ctx)
	ctx.objects[ProxyId(id)] = proxy
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

func (ctx *Context) unregister(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	delete(ctx.objects, proxy.Id())
}

func (ctx *Context) Unregister(proxy Proxy) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	delete(ctx.objects, proxy.Id())
}

func (ctx *Context) Objects() map[ProxyId]Proxy {
	return ctx.objects
}

func (c *Context) Close() {
	c.conn.Close()
	c.exitChan <- struct{}{}
	close(c.dispatchChan)

}

func (c *Context) Dispatch() chan<- struct{} {
	return c.dispatchChan
}

func (c *Context) AddFD(fd uintptr) {
	c.fds = append(c.fds, fd)
}

func (c *Context) NextFD() uintptr {
	if len(c.fds) > 0 {
		fd := c.fds[0]
		c.fds = c.fds[1:]
		return fd
	}
	return 0
}
func (ret *Display) GetFd() (int) {
	return ret.Fd
}


func Connect(addr string) (ret *Display, err error) {
	runtime_dir := os.Getenv("XDG_RUNTIME_DIR")
	if runtime_dir == "" {
		return nil, errors.New("XDG_RUNTIME_DIR not set in the environment.")
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

func NewClientConnect(fd int) *Display {
	c := new(Context)
	c.objects = make(map[ProxyId]Proxy)
	c.currentId = 0
	c.dispatchChan = make(chan struct{})
	c.exitChan = make(chan struct{})
	c.SockFD = fd
	return NewDisplay(c)
}

func Listen(addr string) (ret *Display, err error) {
	c := new(Context)
	c.objects = make(map[ProxyId]Proxy)

	runtime_dir := os.Getenv("XDG_RUNTIME_DIR")
	if runtime_dir == "" {
		return nil, errors.New("XDG_RUNTIME_DIR not set in the environment.")
	}
	if addr == "" {
		addr = os.Getenv("WAYLAND_DISPLAY")
	}
	if addr == "" {
		addr = "wayland-0"
	}
	addr = runtime_dir + "/" + addr

	sockFD, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Println(err)
		// runtime.Goexit()
		return NewDisplay(c), err
	}
	var sockAddr unix.SockaddrUnix
	sockAddr.Name = addr
	unix.Unlink(addr)
	err = unix.Bind(sockFD, &sockAddr)
	if err != nil {
		log.Printf("Couldn't bind %s\n", sockAddr.Name)
	}
	err = unix.Listen(sockFD, 64)
	if err != nil {
		log.Println(err)
		return NewDisplay(c), err
	}

	c.currentId = 0
	c.SockFD = sockFD
	return NewDisplay(c), nil
}

func ListenFD(addr string) (ret int, err error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return -1, errors.New("XDG_RUNTIME_DIR not set in the environment")
	}
	if addr == "" {
		addr = os.Getenv("WAYLAND_DISPLAY")
	}
	if addr == "" {
		addr = "wayland-0"
	}
	addr = runtimeDir + "/" + addr

	sockFD, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	var sockAddr unix.SockaddrUnix
	sockAddr.Name = addr
	unix.Unlink(addr)
	err = unix.Bind(sockFD, &sockAddr)
	if err != nil {
		log.Printf("Couldn't bind %s\n", sockAddr.Name)
	}
	err = unix.Listen(sockFD, 64)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return sockFD, nil
}

func  (c *Context) RunTill(cb *Callback) {
loop:
	for {
		ev, err := c.readEvent()
		if err != nil {
			if err == io.EOF {
				// connection closed
				break loop

			}

			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				log.Print("Timeout Error")
				continue
			}

			log.Fatal(err)
		}
		proxy := c.LookupProxy(ev.Pid)
		if proxy != nil {
			if dispatcher, ok := proxy.(Dispatcher); ok {
			
				if found_cb, ok := dispatcher.(*Callback); ok {
					if found_cb == cb {
						bytePool.Give(ev.Data)
						return
					}
				}
				dispatcher.Dispatch(ev)
				bytePool.Give(ev.Data)
			} else {
				log.Print("Not dispatched")
			}
		} else {
			log.Print("Proxy NULL")
		}
	
	}
}

func (c *Context) Run() {
	// ctx := context.Background()

loop:
	for {
		select {
		case <-c.dispatchChan:
			ev, err := c.readEvent()
			if err != nil {
				if err == io.EOF {
					log.Print("connection closed")
					// connection closed
					break loop

				}

				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					log.Print("Timeout Error")
					continue
				}

				log.Fatal(err)
			}

			proxy := c.LookupProxy(ev.Pid)
			if proxy != nil {
				if dispatcher, ok := proxy.(Dispatcher); ok {
					dispatcher.Dispatch(ev)
					bytePool.Give(ev.Data)
				} else {
					log.Print("Not dispatched")
				}
			} else {
				log.Print("Proxy NULL")
			}

		case <-c.exitChan:
			break loop
		}
	}
}
