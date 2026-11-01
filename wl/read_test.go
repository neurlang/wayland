package wl

import (
	"encoding/binary"
	"errors"
	"strings"
	"syscall"
	"testing"
)

func BenchmarkRunPointerMotion(b *testing.B) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		b.Fatal(err)
	}
	ctx, sock := socketContext(b, fds)
	p := NewPointer(ctx)
	const batch = 64
	msg := make([]byte, 20)
	binary.NativeEndian.PutUint32(msg[0:], uint32(p.Id()))
	binary.NativeEndian.PutUint32(msg[4:], uint32(len(msg))<<16|2)
	flush := make([]byte, 0, batch*len(msg))
	for range batch {
		flush = append(flush, msg...)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i += batch {
		if _, err := syscall.Write(sock, flush); err != nil {
			b.Fatal(err)
		}
		for range batch {
			if err := ctx.Run(); err != nil {
				b.Fatal(err)
			}
		}
	}
}

type targetRecorder []string

func (r *targetRecorder) HandleDataSourceTarget(ev DataSourceTargetEvent) {
	*r = append(*r, ev.MimeType)
}

func TestEventLargerThanTheReadBuffer(t *testing.T) {
	ctx, sock := newSocketContext(t)
	src := NewDataSource(ctx)
	var got targetRecorder
	src.AddTargetHandler(&got)
	mime := strings.Repeat("m", 3*inBufferSize)
	msg := dataSourceSend(src.Id(), mime)
	binary.NativeEndian.PutUint32(msg[4:], uint32(len(msg))<<16)
	short := dataSourceSend(src.Id(), "text/plain")
	binary.NativeEndian.PutUint32(short[4:], uint32(len(short))<<16)
	written := make(chan error, 1)
	go func() {
		_, err := syscall.Write(sock, append(msg, short...))
		written <- err
	}()
	for range 2 {
		if err := ctx.Run(); err != nil {
			t.Fatal(err)
		}
	}
	if err := <-written; err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != mime || got[1] != "text/plain" {
		t.Fatalf("target events %d long, want the %d-byte mime type then text/plain", len(got), len(mime))
	}
}

func TestSizeBelowTheHeaderIsAnError(t *testing.T) {
	ctx, sock := newSocketContext(t)
	msg := make([]byte, 8)
	binary.NativeEndian.PutUint32(msg[0:], 1)
	binary.NativeEndian.PutUint32(msg[4:], 4<<16)
	if _, err := syscall.Write(sock, msg); err != nil {
		t.Fatal(err)
	}
	if err := ctx.Run(); !errors.Is(err, ErrInvalidMsgSize) {
		t.Fatalf("run of a 4-byte message: %v, want %v", err, ErrInvalidMsgSize)
	}
}
