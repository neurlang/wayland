package swizzle

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBGRAConvertsBuffersOfAnyPixelLength(t *testing.T) {
	for _, pixelCount := range []int{0, 1, 8, 256, 576} {
		t.Run(fmt.Sprintf("%d_pixels", pixelCount), func(t *testing.T) {
			buf := make([]byte, pixelCount*4)
			for i := 0; i < len(buf); i += 4 {
				buf[i+0] = byte(i / 4)
				buf[i+1] = 0x11
				buf[i+2] = 0x22
				buf[i+3] = 0x33
			}
			want := append([]byte(nil), buf...)
			for i := 0; i < len(want); i += 4 {
				want[i], want[i+2] = want[i+2], want[i]
			}

			if err := BGRA(buf); err != nil {
				t.Fatalf("BGRA returned an error: %v", err)
			}
			if !bytes.Equal(buf, want) {
				t.Fatal("BGRA produced incorrect pixel order")
			}
		})
	}
}

func TestBGRARejectsPartialPixel(t *testing.T) {
	if err := BGRA(make([]byte, 3)); err != ErrSliceNot32Bit {
		t.Fatalf("BGRA error = %v, want %v", err, ErrSliceNot32Bit)
	}
}
