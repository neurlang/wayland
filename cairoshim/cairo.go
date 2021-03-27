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

type enum = int64

const FormatArgb32 enum = 0
const FormatRgb16565 enum = 4

type Surface interface {
	Reference() Surface
	Destroy()
	SetUserData(data func())
	ImageSurfaceGetData() []byte
	ImageSurfaceGetWidth() int
	ImageSurfaceGetHeight() int
	ImageSurfaceGetStride() int
}
type Format = enum

type simulatedSurface struct {
	data   []byte
	width  int
	height int
	stride int
}

type simulatedSurfaceRef struct {
	surf *simulatedSurface
}

func (s simulatedSurfaceRef) Destroy() {
	s.surf.Destroy()
	s.surf = nil
}

func (s simulatedSurfaceRef) Reference() Surface {
	return s.surf.Reference()
}

func (s *simulatedSurface) Reference() Surface {
	return &simulatedSurfaceRef{surf: s}
}

func (s *simulatedSurface) Destroy() {
}

func (s *simulatedSurfaceRef) ImageSurfaceGetData() []byte {
	return s.surf.ImageSurfaceGetData()
}

func (s *simulatedSurface) ImageSurfaceGetData() []byte {
	return s.data
}

func (s *simulatedSurface) ImageSurfaceGetWidth() int {
	return s.width
}

func (s *simulatedSurface) ImageSurfaceGetHeight() int {
	return s.height
}

func (s *simulatedSurface) ImageSurfaceGetStride() int {

	return s.stride
}

func (s *simulatedSurfaceRef) ImageSurfaceGetWidth() int {
	return s.surf.ImageSurfaceGetWidth()
}

func (s *simulatedSurfaceRef) ImageSurfaceGetHeight() int {
	return s.surf.ImageSurfaceGetHeight()
}

func (s *simulatedSurfaceRef) ImageSurfaceGetStride() int {
	return s.surf.ImageSurfaceGetStride()
}

func ImageSurfaceCreateForData(
	data []byte,
	cairoFormat Format,
	width int,
	height int,
	stride int,
) Surface {
	return &simulatedSurface{data: data, width: width, height: height, stride: stride}
}

func FormatStrideForWidth(cairoFormat Format, width int) int {
	return width * 4
}

func (s *simulatedSurface) SetUserData(data func()) {
}
func (s *simulatedSurfaceRef) SetUserData(data func()) {
	s.surf.SetUserData(data)
}
