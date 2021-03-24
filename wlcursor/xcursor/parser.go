package xcursor

import (
	"bytes"
	"encoding/binary"

	"github.com/neurlang/wayland/internal/swizzle"
)

type Image struct {
	PixRGBA  []uint8
	PixBGRA  []uint8
	Height   uint32
	HotspotX uint32
	HotspotY uint32
	Delay    uint32
	Size     uint32
	Width    uint32
}

type toc struct {
	toctype uint32
	subtype uint32
	pos     uint32
}

func parseHeader(buf *bytes.Buffer) uint32 {
	buf.Next(4) // skip "Xcur"

	buf.Next(4)
	buf.Next(4)
	nToc := binary.LittleEndian.Uint32(buf.Next(4))

	return nToc
}

func parseToc(buf *bytes.Buffer) toc {
	tocType := binary.LittleEndian.Uint32(buf.Next(4))
	subType := binary.LittleEndian.Uint32(buf.Next(4))
	pos := binary.LittleEndian.Uint32(buf.Next(4))

	return toc{
		toctype: tocType,
		subtype: subType,
		pos:     pos,
	}
}

func parseImg(b []byte) Image {
	buf := bytes.NewBuffer(b)
	buf.Next(8) // skip header (header size, type)
	size := binary.LittleEndian.Uint32(buf.Next(4))
	buf.Next(4) // skip image version
	width := binary.LittleEndian.Uint32(buf.Next(4))
	height := binary.LittleEndian.Uint32(buf.Next(4))
	hotspotX := binary.LittleEndian.Uint32(buf.Next(4))
	hotspotY := binary.LittleEndian.Uint32(buf.Next(4))
	delay := binary.LittleEndian.Uint32(buf.Next(4))

	imageLength := 4 * width * height
	pixRGBA := make([]uint8, imageLength)
	_, _ = buf.Read(pixRGBA)

	pixBGRA := make([]uint8, imageLength)
	copy(pixBGRA, pixRGBA)
	swizzle.BGRA(pixBGRA)

	return Image{
		Size:     size,
		Width:    width,
		Height:   height,
		HotspotX: hotspotX,
		HotspotY: hotspotY,
		Delay:    delay,
		PixRGBA:  pixRGBA,
		PixBGRA:  pixBGRA,
	}
}

func ParseXcursor(content []byte) []Image {
	buf := bytes.NewBuffer(content)
	ntoc := parseHeader(buf)
	imgs := make([]Image, ntoc)

	for i := uint32(0); i < ntoc; i++ {
		toc := parseToc(buf)

		if toc.toctype == 0xfffd_0002 {
			index := toc.pos
			img := parseImg(content[index:])
			imgs[i] = img
		}
	}

	return imgs
}
