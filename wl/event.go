package wl

import (
	"bytes"
	"fmt"
	"syscall"
)

type Event struct {
	Pid    ProxyId
	Opcode uint32
	Data   []byte
	scms   []syscall.SocketControlMessage
	off    int
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
		return nil, fmt.Errorf("unable to read message header")
	}
	ev := new(Event)
	if oobn > 0 {
		if oobn > len(control) {
			return nil, fmt.Errorf("unsufficient control msg buffer")
		}
		scms, err := syscall.ParseSocketControlMessage(control)
		if err != nil {
			return nil, fmt.Errorf("control message parse error: %s", err)
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
		return nil, fmt.Errorf("invalid message size")
	}
	ev.Data = data

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
		panic("unable to parse unix rights")
	}
	//TODO is this required
	ev.scms = append(ev.scms, ev.scms[1:]...)
	return uintptr(fds[0])
}

func (ev *Event) Uint32() uint32 {
	buf := ev.next(4)
	if len(buf) != 4 {
		panic("unable to read unsigned int")
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
		panic("unable to read string")
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
