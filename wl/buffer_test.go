package wl

import "testing"

type releaseHandler struct{}

func (releaseHandler) HandleBufferRelease(BufferReleaseEvent) {}

func TestNewServerBufferTakesHandlers(t *testing.T) {
	b := NewServerBuffer()
	if b.Id() != 0 {
		t.Fatalf("id = %d, want unregistered", b.Id())
	}
	b.AddReleaseHandler(releaseHandler{})
}
