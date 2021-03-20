package wl

import (
	"bytes"
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

type Event struct {
	Pid    ProxyId
	Opcode uint32
	Data   []byte
	scms   []syscall.SocketControlMessage
	off    int
	oob    []byte
}

/*
	Okay, so if you pass a file descriptor across a
	UNIX domain socket, you may actually receive it
	on an earlier call to recvmsg. If we don't do
	anything about this we end up getting a file
	descriptor on the wrong wayland protocol message.

	See https://keithp.com/blogs/fd-passing/

	This is a hacky solution:
	- have a map from client fd to list of receive fds
	- when we receive a fd over the socket add it to
	  the clients list of fds
	- when we call event.FD() take the earliest FD
	  from the list

	UPDATE: Scratch the above. The actual solution is to
	store a lists of incoming fds in the Context
	and modify the generator to set FDs like:
		ev.Fd = p.Context().NextFD()
*/

func (c *Context) readEvent() (*Event, error) {
	buf := bytePool.Take(8)
	control := bytePool.Take(24)

	n, oobn, _, _, err := c.conn.ReadMsgUnix(buf[:], control)
	if err != nil {
		return nil, err
	}
	if n != 8 {
		return nil, fmt.Errorf("Unable to read message header.")
	}
	ev := new(Event)
	if oobn > 0 {
		if oobn > len(control) {
			return nil, fmt.Errorf("Unsufficient control msg buffer")
		}
		scms, err := syscall.ParseSocketControlMessage(control)
		if err != nil {
			return nil, fmt.Errorf("Control message parse error: %s", err)
		}
		ev.scms = scms
	}

	ev.Pid = ProxyId(order.Uint32(buf[0:4]))
	ev.Opcode = uint32(order.Uint16(buf[4:6]))
	size := uint32(order.Uint16(buf[6:8]))

	// subtract 8 bytes from header
	data := bytePool.Take(int(size) - 8)
	n, err = c.conn.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(size)-8 {
		return nil, fmt.Errorf("Invalid message size.")
	}
	ev.Data = data

	bytePool.Give(buf)
	bytePool.Give(control)

	return ev, nil
}

func ReadEventUnix(fd int) (Event, error) {
	buf := bytePool.Take(8)
	control := bytePool.Take(24)

	n, oobn, _, _, err := unix.Recvmsg(fd, buf, control, unix.MSG_DONTWAIT)
	if err != nil {
		return Event{}, err
	}
	if n != 8 {
		return Event{}, fmt.Errorf("Unable to read message header.")
	}

	// ev := new(Event)
	var ev Event

	if oobn > 0 {
		if oobn > len(control) {
			return Event{}, fmt.Errorf("Unsufficient control msg buffer")
		}
		scms, err := syscall.ParseSocketControlMessage(control)
		if err != nil {
			return Event{}, fmt.Errorf("Control message parse error: %s", err)
		}
		ev.scms = scms
	}

	ev.Pid = ProxyId(order.Uint32(buf[0:4]))
	// fmt.Println("id", ev.Pid)
	ev.Opcode = uint32(order.Uint16(buf[4:6]))
	size := uint32(order.Uint16(buf[6:8]))

	// subtract 8 bytes from header
	data := bytePool.Take(int(size) - 8)
	n, err = unix.Read(fd, data)
	// n, err = c.conn.Read(data)
	if err != nil {
		return Event{}, err
	}
	if n != int(size)-8 {
		return Event{}, fmt.Errorf("Invalid message size.")
	}
	ev.Data = data
	fmt.Println("data", data)
	bytePool.Give(buf)
	bytePool.Give(control)

	return ev, nil
}

func (c *Context) ReadEvent() (Event, error) {
	buf := bytePool.Take(8)
	control := bytePool.Take(24)

	var ev Event
	if c.fds == nil {
		c.fds = make([]uintptr, 0)
	}

	n, oobn, _, _, err := unix.Recvmsg(c.SockFD, buf, control, unix.MSG_DONTWAIT)
	if err != nil {
		return ev, err
	}
	if n != 8 {
		return ev, fmt.Errorf("Unable to read message header.")
	}

	// ev := new(Event)
	// var ev Event

	if oobn > 0 {
		if oobn > len(control) {
			return ev, fmt.Errorf("Unsufficient control msg buffer")
		}
		scms, err := syscall.ParseSocketControlMessage(control)
		if err != nil {
			return ev, fmt.Errorf("Control message parse error: %s", err)
		}
		ev.scms = scms
		fds, err := syscall.ParseUnixRights(&ev.scms[0])
		if err != nil {
			fmt.Print("Failed to extract fd")
		}
		for _, fd := range fds {
			c.AddFD(uintptr(fd))
		}
	}

	ev.Pid = ProxyId(order.Uint32(buf[0:4]))
	ev.Opcode = uint32(order.Uint16(buf[4:6]))
	size := uint32(order.Uint16(buf[6:8]))

	// subtract 8 bytes from header
	data := bytePool.Take(int(size) - 8)
	n, err = unix.Read(c.SockFD, data)
	// n, err = c.conn.Read(data)
	if err != nil {
		return ev, err
	}
	if n != int(size)-8 {
		return ev, fmt.Errorf("Invalid message size.")
	}
	ev.Data = data
	// fmt.Println(ev.Pid)
	// fmt.Println(ev.Opcode)
	// fmt.Println(data)
	bytePool.Give(buf)
	bytePool.Give(control)

	return ev, nil
}

func (ev *Event) FD() uintptr {
	if ev.scms == nil {
		return 0
	}
	fds, err := syscall.ParseUnixRights(&ev.scms[0])
	if err != nil {
		panic("Unable to parse unix rights")
	}
	//TODO is this required
	ev.scms = append(ev.scms, ev.scms[1:]...)
	return uintptr(fds[0])
}

func (ev *Event) Uint32() uint32 {
	buf := ev.next(4)
	if len(buf) != 4 {
		panic("Unable to read unsigned int")
	}
	return order.Uint32(buf)
}

func (ev *Event) Proxy(c *Context) Proxy {
	id := ev.Uint32()
	if id == 0 {
		return nil
	} else {
		return c.LookupProxy(ProxyId(id))
	}
}

func (ev *Event) String() string {
	l := int(ev.Uint32())
	buf := ev.next(l)
	if len(buf) != l {
		panic("Unable to read string")
	}
	ret := string(bytes.TrimRight(buf, "\x00"))
	//padding to 32 bit boundary
	if (l & 0x3) != 0 {
		ev.next(4 - (l & 0x3))
	}
	return ret
}

func (ev *Event) Int32() int32 {
	return int32(ev.Uint32())
}

func (ev *Event) Float32() float32 {
	return float32(FixedToFloat(ev.Int32()))
}

func (ev *Event) Array() []int32 {
	l := int(ev.Uint32())
	arr := make([]int32, l/4)
	for i := range arr {
		arr[i] = ev.Int32()
	}
	return arr
}

func (ev *Event) next(n int) []byte {
	ret := ev.Data[ev.off : ev.off+n]
	ev.off += n
	return ret
}
