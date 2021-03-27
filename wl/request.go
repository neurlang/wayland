package wl

import (
	"errors"
	"github.com/neurlang/wayland/os"
	"net"

	"github.com/yalue/native_endian"
)

// Request is the request message from your program to the Wayland compositor
type Request struct {
	pid    ProxyId
	Opcode uint32
	data   []byte
	oob    []byte
}

var ErrContextSendRequestUnix = errors.New("unable to send request using unix")
var ErrContextSendRequestConn = errors.New("unable to send request using conn")

var ErrContextSendRequestUnixLength = errors.New("unable to send request using unix, WriteMsgUnix length check failed")

// Context SendRequest sends a specific request with arguments to the compositor
func (ctx *Context) SendRequest(proxy Proxy, opcode uint32, args ...interface{}) (err error) {
	req := Request{
		pid:    proxy.Id(),
		Opcode: opcode,
	}

	for _, arg := range args {
		err1 := req.Write(arg)
		if err1 != nil {
			err = err1
		}
	}

	if err != nil {
		return err
	}

	if ctx.conn != nil {
		return writeRequest(ctx.conn, req)
	} else {
		return writeRequestUnix(ctx.sockFD, req)
	}
}

// Request Write writes a specific request argument to the compositor
func (r *Request) Write(arg interface{}) error {
	switch t := arg.(type) {
	case Proxy:
		r.PutProxy(t)
	case uint32:
		r.PutUint32(t)
	case int32:
		r.PutInt32(t)
	case float32:
		r.PutFloat32(t)
	case string:
		r.PutString(t)
	case []int32:
		r.PutArray(t)
	case uintptr:
		r.PutFd(t)
	default:
		return errors.New("invalid Wayland request parameter type")
	}
	return nil
}

// Request PutUint32 writes an uint32 argument to the compositor
func (r *Request) PutUint32(u uint32) {
	buf := bytePool.Take(4)
	native_endian.NativeEndian().PutUint32(buf, u)
	r.data = append(r.data, buf...)
}

// Request PutProxy writes a proxy argument to the compositor
func (r *Request) PutProxy(p Proxy) {
	r.PutUint32(uint32(p.Id()))
}

// Request PutInt32 writes an int32 argument to the compositor
func (r *Request) PutInt32(i int32) {
	r.PutUint32(uint32(i))
}

// Request PutFloat32 writes a float32 argument to the compositor
func (r *Request) PutFloat32(f float32) {
	fx := FloatToFixed(float64(f))
	r.PutUint32(uint32(fx))
}

// Request PutString writes a string argument to the compositor
func (r *Request) PutString(s string) {
	tail := 4 - (len(s) & 0x3)
	r.PutUint32(uint32(len(s) + tail))
	r.data = append(r.data, []byte(s)...)
	// if padding required
	if tail > 0 {
		padding := make([]byte, tail)
		r.data = append(r.data, padding...)
	}
}

// Request PutArray writes an array argument to the compositor
func (r *Request) PutArray(a []int32) {
	r.PutUint32(uint32(len(a)))
	for _, e := range a {
		r.PutUint32(uint32(e))
	}
}

// Request PutFd writes a file descriptor argument to the compositor
func (r *Request) PutFd(fd uintptr) {
	rights := os.UnixRights(int(fd))
	r.oob = append(r.oob, rights...)
}

func writeRequest(conn *net.UnixConn, r Request) error {
	var header []byte
	// calculate message total size
	size := uint32(len(r.data) + 8)
	buf := make([]byte, 4)
	native_endian.NativeEndian().PutUint32(buf, uint32(r.pid))
	header = append(header, buf...)
	native_endian.NativeEndian().PutUint32(buf, size<<16|r.Opcode&0x0000ffff)
	header = append(header, buf...)

	d, c, err := conn.WriteMsgUnix(append(header, r.data...), r.oob, nil)
	if err != nil {
		return combinedError{ErrContextSendRequestConn, err}
	}
	if c != len(r.oob) || d != (len(header)+len(r.data)) {
		return ErrContextSendRequestUnixLength
	}
	bytePool.Give(r.data)

	return nil
}

func writeRequestUnix(fd int, r Request) error {
	var header []byte
	// calculate message total size
	size := uint32(len(r.data) + 8)
	buf := make([]byte, 4)
	native_endian.NativeEndian().PutUint32(buf, uint32(r.pid))
	header = append(header, buf...)
	native_endian.NativeEndian().PutUint32(buf, size<<16|r.Opcode&0x0000ffff)
	header = append(header, buf...)

	// unix.
	var addr os.Sockaddr
	err := os.Sendmsg(fd, append(header, r.data...), r.oob, addr, 0)
	if err != nil {
		return combinedError{ErrContextSendRequestUnix, err}
	}
	bytePool.Give(r.data)

	return nil
}
