package wl

import (
	"encoding/binary"
	"errors"
	"io"
	"syscall"
	"testing"
)

type motionCounter int

func (c *motionCounter) HandlePointerMotion(PointerMotionEvent) { *c++ }

func pointerMotions(p *Pointer, n int) []byte {
	msg := make([]byte, 20)
	binary.NativeEndian.PutUint32(msg[0:], uint32(p.Id()))
	binary.NativeEndian.PutUint32(msg[4:], uint32(len(msg))<<16|2)
	out := make([]byte, 0, n*len(msg))
	for range n {
		out = append(out, msg...)
	}
	return out
}

func TestReadReadyWithNothingSentReadsNothing(t *testing.T) {
	ctx, _ := newSocketContext(t)
	read, err := ctx.ReadReady()
	if err != nil || read {
		t.Fatalf("ReadReady on a quiet socket = %v, %v; want false, nil", read, err)
	}
	if ctx.Buffered() {
		t.Fatal("a quiet socket buffered an event")
	}
}

func TestBufferedAndReadReadyDrainABurst(t *testing.T) {
	ctx, sock := newSocketContext(t)
	p := NewPointer(ctx)
	var got motionCounter
	p.AddMotionHandler(&got)
	const burst = 40
	if _, err := syscall.Write(sock, pointerMotions(p, burst)); err != nil {
		t.Fatal(err)
	}
	read, err := ctx.ReadReady()
	if err != nil || !read {
		t.Fatalf("ReadReady after a burst = %v, %v; want true, nil", read, err)
	}
	for ctx.Buffered() {
		if err := ctx.Run(); err != nil {
			t.Fatal(err)
		}
	}
	if got != burst {
		t.Fatalf("dispatched %d motions from the buffer, want %d", got, burst)
	}
	if read, err := ctx.ReadReady(); err != nil || read {
		t.Fatalf("ReadReady after the drain = %v, %v; want false, nil", read, err)
	}
}

func TestBufferedHoldsBackAPartialEvent(t *testing.T) {
	ctx, sock := newSocketContext(t)
	p := NewPointer(ctx)
	whole := pointerMotions(p, 1)
	if _, err := syscall.Write(sock, whole[:12]); err != nil {
		t.Fatal(err)
	}
	if read, err := ctx.ReadReady(); err != nil || !read {
		t.Fatalf("ReadReady = %v, %v; want true, nil", read, err)
	}
	if ctx.Buffered() {
		t.Fatal("a partial event reads as buffered")
	}
	if _, err := syscall.Write(sock, whole[12:]); err != nil {
		t.Fatal(err)
	}
	if _, err := ctx.ReadReady(); err != nil {
		t.Fatal(err)
	}
	if !ctx.Buffered() {
		t.Fatal("the completed event is not buffered")
	}
}

func TestReadReadyReportsAHangUp(t *testing.T) {
	ctx, sock := newSocketContext(t)
	if err := syscall.Shutdown(sock, syscall.SHUT_WR); err != nil {
		t.Fatal(err)
	}
	if _, err := ctx.ReadReady(); !errors.Is(err, io.EOF) {
		t.Fatalf("ReadReady after the peer hung up = %v, want io.EOF", err)
	}
}
