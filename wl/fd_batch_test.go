package wl

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

type sendRecorder []DataSourceSendEvent

func (r *sendRecorder) HandleDataSourceSend(ev DataSourceSendEvent) {
	*r = append(*r, ev)
}

func newSocketContext(t *testing.T) (*Context, int) {
	t.Helper()
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	return socketContext(t, fds)
}

func socketContext(t testing.TB, fds [2]int) (*Context, int) {
	t.Helper()
	t.Cleanup(func() { syscall.Close(fds[1]) })
	f := os.NewFile(uintptr(fds[0]), "client")
	conn, err := net.FileConn(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	ctx := &Context{conn: conn.(*net.UnixConn), sockFD: -1, objects: make(map[ProxyId]Proxy)}
	t.Cleanup(func() { ctx.Close() })
	return ctx, fds[1]
}

func dataSourceSend(id ProxyId, mime string) []byte {
	l := len(mime) + 1
	padded := (l + 3) &^ 3
	msg := make([]byte, 12+padded)
	binary.NativeEndian.PutUint32(msg[0:], uint32(id))
	binary.NativeEndian.PutUint32(msg[4:], uint32(len(msg))<<16|1)
	binary.NativeEndian.PutUint32(msg[8:], uint32(l))
	copy(msg[12:], mime)
	return msg
}

func compositorFlush(t *testing.T, sock int, msgs [][]byte, fds []int) {
	t.Helper()
	var data []byte
	for _, m := range msgs {
		data = append(data, m...)
	}
	if err := syscall.Sendmsg(sock, data, syscall.UnixRights(fds...), nil, 0); err != nil {
		t.Fatal(err)
	}
}

func pipes(t *testing.T, n int) (readers []*os.File, writers []int) {
	t.Helper()
	for range n {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { r.Close(); w.Close() })
		readers = append(readers, r)
		writers = append(writers, int(w.Fd()))
	}
	return readers, writers
}

func TestEventsOfOneFlushEachGetTheirFD(t *testing.T) {
	ctx, sock := newSocketContext(t)
	src := NewDataSource(ctx)
	var got sendRecorder
	src.AddSendHandler(&got)

	const events = 3
	readers, writers := pipes(t, events)
	msgs := make([][]byte, events)
	for i := range msgs {
		msgs[i] = dataSourceSend(src.Id(), "text/plain")
	}
	compositorFlush(t, sock, msgs, writers)

	for range events {
		if err := ctx.Run(); err != nil {
			t.Fatal(err)
		}
	}
	if len(got) != events {
		t.Fatalf("dispatched %d send events, want %d", len(got), events)
	}
	for i, ev := range got {
		if ev.FdError != nil {
			t.Fatalf("event %d: %v", i, ev.FdError)
		}
		w := os.NewFile(ev.Fd, "received")
		if _, err := w.Write([]byte{byte(i)}); err != nil {
			t.Fatalf("event %d: write to its fd: %v", i, err)
		}
		w.Close()
		b := make([]byte, 1)
		if _, err := io.ReadFull(readers[i], b); err != nil {
			t.Fatalf("event %d: read its pipe: %v", i, err)
		}
		if b[0] != byte(i) {
			t.Fatalf("event %d got the fd of pipe %d", i, b[0])
		}
	}
}

func TestCloseClosesFDsNoEventTook(t *testing.T) {
	ctx, sock := newSocketContext(t)
	src := NewDataSource(ctx)
	var got sendRecorder
	src.AddSendHandler(&got)

	_, writers := pipes(t, 2)
	compositorFlush(t, sock, [][]byte{dataSourceSend(src.Id(), "text/plain")}, writers)
	if err := ctx.Run(); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].FdError != nil {
		t.Fatalf("send events %+v, want one with an fd", got)
	}
	defer syscall.Close(int(got[0].Fd))
	unused := ctx.fds[0]

	if err := ctx.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := fcntlGetFD(unused); !errors.Is(err, syscall.EBADF) {
		t.Fatalf("fd %d no event took is still open after Close (fcntl: %v)", unused, err)
	}
}

func openFDs(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		t.Fatal(err)
	}
	return len(entries)
}

func assertOwnFD(t *testing.T, ev DataSourceSendEvent, own *os.File) {
	t.Helper()
	if ev.FdError != nil {
		t.Fatal(ev.FdError)
	}
	w := os.NewFile(ev.Fd, "received")
	_, err := w.Write([]byte{1})
	w.Close()
	if err != nil {
		t.Fatalf("write to the event's fd: %v", err)
	}
	if err := own.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(own, make([]byte, 1)); err != nil {
		t.Fatalf("the event got another event's fd: its own pipe read %v", err)
	}
}

func TestUnhandledFDEventLeavesLaterFDsAligned(t *testing.T) {
	for _, tc := range []struct {
		name    string
		destroy bool
		runErr  error
	}{
		{name: "no handler"},
		{name: "destroyed object", destroy: true, runErr: ErrContextRunProxyNil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, sock := newSocketContext(t)
			other := NewDataSource(ctx)
			src := NewDataSource(ctx)
			var got sendRecorder
			src.AddSendHandler(&got)
			if tc.destroy {
				other.Unregister()
			}
			readers, writers := pipes(t, 2)
			before := openFDs(t)
			compositorFlush(t, sock, [][]byte{dataSourceSend(other.Id(), "text/plain")}, writers[:1])
			if err := ctx.Run(); !errors.Is(err, tc.runErr) {
				t.Fatalf("run of the unhandled event: %v, want %v", err, tc.runErr)
			}
			if after := openFDs(t); after != before {
				t.Fatalf("%d fds open after the unhandled event, %d before", after, before)
			}
			compositorFlush(t, sock, [][]byte{dataSourceSend(src.Id(), "text/plain")}, writers[1:])
			if err := ctx.Run(); err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("dispatched %d send events, want 1", len(got))
			}
			assertOwnFD(t, got[0], readers[1])
		})
	}
}

func TestDeleteIdReleasesADestroyedObject(t *testing.T) {
	ctx, sock := newSocketContext(t)
	ctx.newDisplay()
	src := NewDataSource(ctx)
	src.Unregister()
	if _, ok := ctx.zombies[src.Id()]; !ok {
		t.Fatal("a destroyed data source is not kept until its id is deleted")
	}
	msg := make([]byte, 12)
	binary.NativeEndian.PutUint32(msg[0:], 1)
	binary.NativeEndian.PutUint32(msg[4:], uint32(len(msg))<<16|1)
	binary.NativeEndian.PutUint32(msg[8:], uint32(src.Id()))
	if _, err := syscall.Write(sock, msg); err != nil {
		t.Fatal(err)
	}
	if err := ctx.Run(); err != nil {
		t.Fatal(err)
	}
	if _, ok := ctx.zombies[src.Id()]; ok {
		t.Fatal("wl_display.delete_id left the destroyed data source behind")
	}
}

func fcntlGetFD(fd int) (uintptr, error) {
	r, _, errno := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), syscall.F_GETFD, 0)
	if errno != 0 {
		return r, errno
	}
	return r, nil
}
