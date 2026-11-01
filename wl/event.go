package wl

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/neurlang/wayland/os"
	"github.com/yalue/native_endian"
)

// Event is the Wayland event (e.g. a response) from the compositor
type Event struct {
	Pid    ProxyId
	Opcode uint32
	Data   []byte
	off    int
	err    error
	fds    []int
}

// ErrReadHeader (Error unable to read message header) is returned when it is not possible to read enough bytes from the unix socket,
// use Unwrap() to get the underlying cause and GetExternal to get this cause
var ErrReadHeader = errors.New("unable to read message header")

// ErrControlMsgBuffer (Error insufficient control msg buffer) is returned when the kernel discarded control data, and the fds in it, for lack of buffer space
var ErrControlMsgBuffer = errors.New("insufficient control msg buffer")

// ErrControlMsgParseError (Error control message parse error) is returned when the unix socket control message cannot be parsed, use
// Unwrap() to get the underlying cause and GetExternal to get this cause
var ErrControlMsgParseError = errors.New("control message parse error")

// ErrInvalidMsgSize (Error invalid message size) is returned when a message header declares a size smaller than the header
var ErrInvalidMsgSize = errors.New("invalid message size")

// ErrReadPayload (Error cannot read message) is returned when the payload message cannot be read, use Unwrap() to get the
// underlying cause and GetExternal to get this cause
var ErrReadPayload = errors.New("cannot read message")

const (
	libwaylandMaxFDsOut = 28
	inBufferSize        = 4096
	headerSize          = 8
)

var controlLen = os.CmsgSpace(libwaylandMaxFDsOut * 4)

func (ctx *Context) readEvent() (*Event, error) {
	if ctx.conn == nil {
		return nil, ErrContextConnNil
	}
	header, err := ctx.fill(headerSize)
	if err != nil {
		return nil, combinedError{ErrReadHeader, err}
	}
	size := int(native_endian.NativeEndian().Uint16(header[6:8]))
	if size < headerSize {
		return nil, ErrInvalidMsgSize
	}
	msg, err := ctx.fill(size)
	if err != nil {
		return nil, combinedError{ErrReadPayload, err}
	}
	ctx.inStart += size
	return &Event{
		Pid:    ProxyId(native_endian.NativeEndian().Uint32(msg[0:4])),
		Opcode: uint32(native_endian.NativeEndian().Uint16(msg[4:6])),
		Data:   msg[headerSize:],
	}, nil
}

func (ctx *Context) fill(n int) ([]byte, error) {
	for len(ctx.in)-ctx.inStart < n {
		if ctx.inStart > 0 {
			ctx.in = ctx.in[:copy(ctx.in, ctx.in[ctx.inStart:])]
			ctx.inStart = 0
		}
		if cap(ctx.in) < n || cap(ctx.in) == len(ctx.in) {
			grown := make([]byte, len(ctx.in), max(n, inBufferSize))
			copy(grown, ctx.in)
			ctx.in = grown
		}
		m, err := ctx.read(ctx.in[len(ctx.in):cap(ctx.in)])
		if m > 0 {
			ctx.in = ctx.in[:len(ctx.in)+m]
		}
		if err != nil {
			return nil, err
		}
	}
	return ctx.in[ctx.inStart : ctx.inStart+n], nil
}

func (ctx *Context) read(b []byte) (int, error) {
	if ctx.control == nil {
		ctx.control = make([]byte, controlLen)
	}
	n, oobn, flags, _, err := ctx.conn.ReadMsgUnix(b, ctx.control)
	if oobn > 0 {
		if qerr := ctx.queueFDs(ctx.control[:oobn]); qerr != nil {
			return n, qerr
		}
	}
	if err != nil {
		return n, err
	}
	if flags&os.MsgCtrunc != 0 {
		return n, ErrControlMsgBuffer
	}
	return n, nil
}

func (ctx *Context) queueFDs(control []byte) error {
	scms, err := os.ParseSocketControlMessage(control)
	if err != nil {
		return combinedError{ErrControlMsgParseError, err}
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	var closeErr error
	for i := range scms {
		fds, err := os.ParseUnixRights(&scms[i])
		if err != nil {
			return combinedError{ErrUnableToParseUnixRights, err}
		}
		if ctx.conn != nil {
			ctx.fds = append(ctx.fds, fds...)
			continue
		}
		for _, fd := range fds {
			closeErr = errors.Join(closeErr, os.Close(fd))
		}
	}
	return closeErr
}

type fdEvents interface {
	EventFDs(opcode uint32) int
}

func (ctx *Context) eventFDCount(ev *Event, proxy Proxy) int {
	if proxy != nil {
		if c, ok := proxy.(fdEvents); ok {
			return c.EventFDs(ev.Opcode)
		}
		return 0
	}
	ctx.mu.RLock()
	c, ok := ctx.zombies[ev.Pid]
	ctx.mu.RUnlock()
	if !ok {
		return 0
	}
	return c.EventFDs(ev.Opcode)
}

func (ctx *Context) claimFDs(ev *Event, proxy Proxy) error {
	n := ctx.eventFDCount(ev, proxy)
	if n == 0 {
		return nil
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if len(ctx.fds) < n {
		return fmt.Errorf("event %d of object %d carries %d fds, %d arrived: %w", ev.Opcode, ev.Pid, n, len(ctx.fds), ErrNoControlMsgs)
	}
	ev.fds = ctx.fds[:n:n]
	ctx.fds = ctx.fds[n:]
	return nil
}

func (ev *Event) closeUnclaimedFDs() error {
	var err error
	for _, fd := range ev.fds {
		err = errors.Join(err, os.Close(fd))
	}
	ev.fds = nil
	return err
}

// ErrNoControlMsgs (Error no socket control messages) is returned when an
// event carries more fd arguments than fds arrived for it
var ErrNoControlMsgs = errors.New("no socket control messages")

// ErrUnableToParseUnixRights (Error unable to parse unix rights)
var ErrUnableToParseUnixRights = errors.New("unable to parse unix rights")

// FD (Event FD) takes this event's next file descriptor, which the caller
// then owns, and an optional error
func (ev *Event) FD() (uintptr, error) {
	if ev.err != nil {
		return 0, ev.err
	}
	if len(ev.fds) == 0 {
		return 0, ErrNoControlMsgs
	}
	fd := ev.fds[0]
	ev.fds = ev.fds[1:]
	return uintptr(fd), nil
}

// ErrUnableToParseUint32 (Error unable to read unsigned int) is returned when the buffer is too short to contain a specific unsigned int
var ErrUnableToParseUint32 = errors.New("unable to read unsigned int")

// Uint32 (Event Uint32) decodes an Uint32 from the Event
func (ev *Event) Uint32() uint32 {
	buf := ev.next(4)
	if len(buf) != 4 {
		ev.err = ErrUnableToParseUint32
		return 0
	}
	return native_endian.NativeEndian().Uint32(buf)
}

// SafeCast casts to any type, it doesn't panic if proxy is nil
func SafeCast[T any](p Proxy) (out T) {
	out, _ = p.(T)
	return
}

// Proxy (Event Proxy) decodes Proxy by it's Id from the Event
func (ev *Event) Proxy(c *Context) Proxy {
	id := ev.Uint32()
	if id != 0 {
		return c.LookupProxy(ProxyId(id))
	}
	return nil
}

// ErrUnableToParseString (Error unable to parse string) is returned when the buffer is too short to contain a specific string
var ErrUnableToParseString = errors.New("unable to parse string")

// String (Event String) decodes a string from the Event
func (ev *Event) String() string {
	l := int(ev.Uint32())
	buf := ev.next(l)
	if len(buf) != l {
		ev.err = ErrUnableToParseString
		return ""
	}
	ret := string(bytes.TrimRight(buf, "\x00"))
	//padding to 32 bit boundary
	if (l & 0x3) != 0 {
		ev.next(4 - (l & 0x3))
	}
	return ret
}

// Int32 (Event Int32) decodes an Int32 from the Event
func (ev *Event) Int32() int32 {
	return int32(ev.Uint32())
}

// Float32 (Event Float32) decodes a Float32 from the Event
func (ev *Event) Float32() float32 {
	return float32(FixedToFloat(ev.Int32()))
}

// Array (Event Array) decodes an Array from the Event
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

func (ev *Event) NewId(i Proxy, c *Context) Proxy {
	c.RegisterMapped(i, ev.Uint32())
	return i
}
