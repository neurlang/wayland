package wl

import (
	"encoding/binary"
	"io"
	"os"
	"sync"
	"syscall"
	"testing"
)

// Client-allocated ids must reach the compositor in allocation order:
// libwayland's wl_map only accepts a new id that is the next free slot.
// Here many goroutines each allocate a callback (wl_display.sync) and send
// it, and the far end of a socketpair checks that the new ids it reads
// never skip.
func TestNewIdsReachTheWireInOrder(t *testing.T) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &Context{sockFD: fds[0], objects: make(map[ProxyId]Proxy)}
	display := NewDisplay(ctx) // id 1, never sent; the compositor owns it
	reader := os.NewFile(uintptr(fds[1]), "compositor")
	defer reader.Close()

	const workers, perWorker = 8, 200
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				if _, err := display.Sync(); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		syscall.Close(fds[0])
	}()

	want := ProxyId(2)
	hdr := make([]byte, 12) // object, size<<16|opcode, new_id
	for n := 0; n < workers*perWorker; n++ {
		if _, err := io.ReadFull(reader, hdr); err != nil {
			t.Fatalf("message %d: %v", n, err)
		}
		if got := ProxyId(binary.LittleEndian.Uint32(hdr[8:])); got != want {
			t.Fatalf("message %d: new id %d, want %d (ids reached the wire out of order)", n, got, want)
		}
		want++
	}
}

// A nil object argument goes out as id 0 instead of panicking with the
// send lock held.
func TestNilObjectArgument(t *testing.T) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fds[0])
	ctx := &Context{sockFD: fds[0], objects: make(map[ProxyId]Proxy)}
	reader := os.NewFile(uintptr(fds[1]), "compositor")
	defer reader.Close()

	dev := &DataDevice{}
	dev.SetContext(ctx)
	dev.SetId(5)
	if err := dev.SetSelection(nil, 42); err != nil {
		t.Fatal(err)
	}
	msg := make([]byte, 16) // object, size<<16|opcode, source, serial
	if _, err := io.ReadFull(reader, msg); err != nil {
		t.Fatal(err)
	}
	if src, serial := binary.LittleEndian.Uint32(msg[8:]), binary.LittleEndian.Uint32(msg[12:]); src != 0 || serial != 42 {
		t.Fatalf("set_selection(nil, 42) sent source %d serial %d", src, serial)
	}
}
