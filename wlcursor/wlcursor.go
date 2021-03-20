package wlcursor

import wl "github.com/neurlang/wayland/wl"

type Theme struct {
}

type Cursor struct {
}

func (c *Cursor) GetCursorImage(index int) *Image {
	println("GetCursorImage: unimplemented")

	return nil
}
func (c *Cursor) FrameAndDuration(a uint32, b *uint32) int {
	panic("unimplemented")
}

func (c *Cursor) ImageCount() int {
	return 1
}

type Image struct {
}

func (i *Image) GetBuffer() *wl.Buffer {
	panic("unimplemented")
}
func (i *Image) GetWidth() int {
	panic("unimplemented")
}
func (i *Image) GetHeight() int {
	panic("unimplemented")
}

func (i *Image) GetHotspotX() int {
	return 0
}

func (i *Image) GetHotspotY() int {
	return 0
}

func LoadTheme(name []byte, size int, shm *wl.Shm) (t *Theme) {
	println("LoadTheme: unimplemented")
	return &Theme{}
}

func (t *Theme) GetCursor(name []byte) *Cursor {
	println("GetCursor: unimplemented")
	return &Cursor{}
}

func (t *Theme) Destroy() {
	println("Theme Destroy: unimplemented")
}

func PointerSetCursor(_ *wl.Pointer, _ uint32, _ *wl.Surface, _ int32, _ int32) {
	panic("unimplemented")
}
