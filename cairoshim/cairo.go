// Copyright 2021 Neurlang project

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Package cairo implements a simple dummy cairo api
package cairo

import "runtime"

type enum = int64

// FormatArgb32 pixel is a 32-bit quantity, with alpha in the upper 8 bits, then red, then green, then blue. The 32-bit quantities are stored native-endian. Pre-multiplied alpha is used. (That is, 50% transparent red is 0x80800000, not 0x80ff0000.)
const FormatArgb32 enum = 0

// FormatRgb16565 pixel is a 16-bit quantity, with red in the upper 5 bits, then green in the next 6, then blue in the lowest 5 bits
const FormatRgb16565 enum = 4

// Surface is a Cairo Surface
type Surface interface {
	Reference() Surface
	Destroy()
	SetUserData(data func())
	SetDestructor(destructor func())
	ImageSurfaceGetData() []byte
	ImageSurfaceGetWidth() int
	ImageSurfaceGetHeight() int
	ImageSurfaceGetStride() int
}

// Format represents an enum
type Format = enum

type simulatedSurface struct {
	data   []byte
	width  int
	height int
	stride int

	destructor func()
}

type simulatedSurfaceRef struct {
	surf *simulatedSurface
}

// Destroy destroys a simulatedSurfaceRef
func (s simulatedSurfaceRef) Destroy() {
	s.surf.Destroy()
	s.surf = nil
}

// References a simulatedSurfaceRef
func (s simulatedSurfaceRef) Reference() Surface {
	return s.surf.Reference()
}

// References a simulatedSurface
func (s *simulatedSurface) Reference() Surface {
	return &simulatedSurfaceRef{surf: s}
}

// Destroy destroys a simulatedSurface
func (s *simulatedSurface) Destroy() {
}

// ImageSurfaceGetData gets the raw image surface data
func (s *simulatedSurfaceRef) ImageSurfaceGetData() []byte {
	return s.surf.ImageSurfaceGetData()
}

// ImageSurfaceGetData gets the raw image surface data
func (s *simulatedSurface) ImageSurfaceGetData() []byte {
	return s.data
}

// ImageSurfaceGetWidth gets width
func (s *simulatedSurface) ImageSurfaceGetWidth() int {
	return s.width
}

// ImageSurfaceGetHeight gets height
func (s *simulatedSurface) ImageSurfaceGetHeight() int {
	return s.height
}

// ImageSurfaceGetHeight gets stride
func (s *simulatedSurface) ImageSurfaceGetStride() int {

	return s.stride
}

// ImageSurfaceGetWidth gets width
func (s *simulatedSurfaceRef) ImageSurfaceGetWidth() int {
	return s.surf.ImageSurfaceGetWidth()
}

// ImageSurfaceGetHeight gets height
func (s *simulatedSurfaceRef) ImageSurfaceGetHeight() int {
	return s.surf.ImageSurfaceGetHeight()
}

// ImageSurfaceGetHeight gets height
func (s *simulatedSurfaceRef) ImageSurfaceGetStride() int {
	return s.surf.ImageSurfaceGetStride()
}

// ImageSurfaceCreateForData creates a simulatedSurface
func ImageSurfaceCreateForData(
	data []byte,
	cairoFormat Format,
	width int,
	height int,
	stride int,
) Surface {
	return &simulatedSurface{data: data, width: width, height: height, stride: stride}
}

// FormatStrideForWidth gets the stride for a concrete image
func FormatStrideForWidth(cairoFormat Format, width int) int {
	return width * 4
}

// SetUserData sets destructor, unused
func (s *simulatedSurface) SetDestructor(data func()) {
	runtime.SetFinalizer(s, func(interface{}) {
		data()
	})
}

// SetUserData sets destructor, unused
func (s *simulatedSurfaceRef) SetDestructor(data func()) {
	s.surf.SetDestructor(data)
}

// SetUserData sets User Data, unused
func (s *simulatedSurface) SetUserData(data func()) {
}

// SetUserData sets User Data, unused
func (s *simulatedSurfaceRef) SetUserData(data func()) {
}
