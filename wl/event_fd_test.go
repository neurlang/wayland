package wl

import (
	"encoding/binary"
	"net"
	"os"
	"syscall"
	"testing"
)

// A compositor can pass fds on several events that the client reads
// before it decodes any of them (a keymap fd, then a clipboard pipe).
// Each queued control message must keep its own fd: the read buffer it
// was parsed from goes back to the pool and is reused by the next read.
func TestQueuedFdsSurviveTheNextRead(t *testing.T) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fds[1])
	f := os.NewFile(uintptr(fds[0]), "client")
	c, err := net.FileConn(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	conn := c.(*net.UnixConn)
	defer conn.Close()
	ctx := &Context{conn: conn, objects: make(map[ProxyId]Proxy)}
	kb := NewKeyboard(ctx)

	first, second := tempFile(t), tempFile(t)
	for _, file := range []*os.File{first, second} {
		msg := make([]byte, 16) // wl_keyboard.keymap: format, fd, size
		binary.LittleEndian.PutUint32(msg[0:], uint32(kb.Id()))
		binary.LittleEndian.PutUint32(msg[4:], 16<<16)
		if err := syscall.Sendmsg(fds[1], msg, syscall.UnixRights(int(file.Fd())), nil, 0); err != nil {
			t.Fatal(err)
		}
	}

	ev, err := ctx.readEvent()
	if err != nil {
		t.Fatal(err)
	}
	if err := ctx.claimFDs(ev, kb); err != nil {
		t.Fatal(err)
	}
	ev2, err := ctx.readEvent()
	if err != nil {
		t.Fatal(err)
	}
	if err := ctx.claimFDs(ev2, kb); err != nil {
		t.Fatal(err)
	}
	for i, c := range []struct {
		ev   *Event
		want *os.File
	}{{ev, first}, {ev2, second}} {
		fd, err := c.ev.FD()
		if err != nil {
			t.Fatalf("event %d: %v", i, err)
		}
		same := sameFile(t, int(fd), c.want)
		syscall.Close(int(fd))
		if !same {
			t.Fatalf("event %d: fd %d is not the file sent with it", i, fd)
		}
	}
}

func tempFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "fd")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

func sameFile(t *testing.T, fd int, want *os.File) bool {
	t.Helper()
	var got, exp syscall.Stat_t
	if err := syscall.Fstat(fd, &got); err != nil {
		t.Fatalf("fstat received fd %d: %v", fd, err)
	}
	if err := syscall.Fstat(int(want.Fd()), &exp); err != nil {
		t.Fatal(err)
	}
	return got.Dev == exp.Dev && got.Ino == exp.Ino
}
