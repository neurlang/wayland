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

const FORMAT_ARGB32 enum = 0
const FORMAT_RGB16_565 enum = 4

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

type simulated_surface struct {
	data   []byte
	width  int
	height int
	stride int
}

type simulated_surface_ref struct {
	surf *simulated_surface
}

func (s simulated_surface_ref) Destroy() {
	s.surf.Destroy()
	s.surf = nil
}

func (s simulated_surface_ref) Reference() Surface {
	return s.surf.Reference()
}

func (s *simulated_surface) Reference() Surface {
	return &simulated_surface_ref{surf: s}
}

func (s *simulated_surface) Destroy() {
}

func (s *simulated_surface_ref) ImageSurfaceGetData() []byte {
	return s.surf.ImageSurfaceGetData()
}

func (s *simulated_surface) ImageSurfaceGetData() []byte {
	return s.data
}

func (s *simulated_surface) ImageSurfaceGetWidth() int {
	return s.width
}

func (s *simulated_surface) ImageSurfaceGetHeight() int {
	return s.height
}

func (s *simulated_surface) ImageSurfaceGetStride() int {

	return s.stride
}

func (s *simulated_surface_ref) ImageSurfaceGetWidth() int {
	return s.surf.ImageSurfaceGetWidth()
}

func (s *simulated_surface_ref) ImageSurfaceGetHeight() int {
	return s.surf.ImageSurfaceGetHeight()
}

func (s *simulated_surface_ref) ImageSurfaceGetStride() int {
	return s.surf.ImageSurfaceGetStride()
}

func ImageSurfaceCreateForData(data []byte, cairo_format Format, width int, height int, stride int) Surface {
	return &simulated_surface{data: data, width: width, height: height, stride: stride}
}

func FormatStrideForWidth(cairo_format Format, width int) int {
	return width * 4
}

func (surface *simulated_surface) SetUserData(data func()) {
}
func (surface *simulated_surface_ref) SetUserData(data func()) {
	surface.surf.SetUserData(data)
}
