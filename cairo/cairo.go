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

// Package cairo implements a cgo cairo api
package cairo

/*
#cgo pkg-config: cairo

#include <cairo.h>
#include "bridge.h"
*/
import "C"
import "unsafe"
import cairoshim "github.com/neurlang/wayland/cairoshim"

type Surface = cairoshim.Surface

const FORMAT_ARGB32 = C.CAIRO_FORMAT_ARGB32
const FORMAT_RGB16_565 = C.CAIRO_FORMAT_RGB16_565

type Format C.cairo_format_t
type UserDataKey C.struct__cairo_user_data_key

type simulated_surface struct {
	data   interface{}
	width  int
	height int
	stride int
	format Format
	fun    func()
}

type simulated_surface_ref struct {
	surf *simulated_surface
}

type cairo_surface C.cairo_surface_t

type cairo_surface_ref struct {
	surf *cairo_surface
}

func (s cairo_surface_ref) Destroy() {
	s.surf.Destroy()
	s.surf = nil
}

func (s simulated_surface_ref) Destroy() {
	s.surf.Destroy()
	s.surf = nil
}

func (s cairo_surface_ref) Reference() Surface {
	return s.surf.Reference()
}

func (s simulated_surface_ref) Reference() Surface {
	return s.surf.Reference()
}

func (s *simulated_surface) Reference() Surface {
	return &simulated_surface_ref{surf: s}
}

func (s *cairo_surface) Reference() Surface {
	return &cairo_surface_ref{surf: (*cairo_surface)(C.cairo_surface_reference((*C.cairo_surface_t)(s)))}
}

func (s *cairo_surface) Destroy() {
	C.cairo_surface_destroy((*C.struct__cairo_surface)(s))
}

func (s *simulated_surface) Destroy() {
}

func (s *cairo_surface) ImageSurfaceGetRawData() uintptr {
	return (uintptr)(unsafe.Pointer(C.cairo_image_surface_get_data((*C.struct__cairo_surface)(s))))
}

func (s *cairo_surface) ImageSurfaceGetData(bytes int) interface{} {

	var newCArray = (*uint32)(unsafe.Pointer(C.cairo_image_surface_get_data((*C.struct__cairo_surface)(s))))
	return (*[1 << 30]uint32)(unsafe.Pointer(newCArray))[:bytes:bytes]
}

func (s *cairo_surface_ref) ImageSurfaceGetData(bytes int) interface{} {
	return s.surf.ImageSurfaceGetData(bytes)
}

func (s *simulated_surface_ref) ImageSurfaceGetData(bytes int) interface{} {
	return s.surf.ImageSurfaceGetData(bytes)
}

func (s *cairo_surface) ImageSurfaceGetWidth() int {
	return int(C.cairo_image_surface_get_width((*C.struct__cairo_surface)(s)))
}

func (s *cairo_surface) ImageSurfaceGetHeight() int {
	return int(C.cairo_image_surface_get_height((*C.struct__cairo_surface)(s)))
}

func (s *cairo_surface) ImageSurfaceGetStride() int {
	return int(C.cairo_image_surface_get_stride((*C.struct__cairo_surface)(s)))
}

func (s *simulated_surface) ImageSurfaceGetRawData() uintptr {
	if real_data, ok := s.data.([]byte); ok {
		return uintptr(unsafe.Pointer(&real_data[0]))
	}
	return uintptr(0)
}

func (s *simulated_surface) ImageSurfaceGetData(bytes int) interface{} {
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

func (s *simulated_surface_ref) ImageSurfaceGetRawData() uintptr {
	return s.surf.ImageSurfaceGetRawData()
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

func (s *cairo_surface_ref) ImageSurfaceGetRawData() uintptr {
	return s.surf.ImageSurfaceGetRawData()
}

func (s *cairo_surface_ref) ImageSurfaceGetWidth() int {
	return s.surf.ImageSurfaceGetWidth()
}

func (s *cairo_surface_ref) ImageSurfaceGetHeight() int {
	return s.surf.ImageSurfaceGetHeight()
}

func (s *cairo_surface_ref) ImageSurfaceGetStride() int {
	return s.surf.ImageSurfaceGetStride()
}

func ImageSurfaceCreateForData(data interface{}, cairo_format Format, width int, height int, stride int) cairoshim.Surface {

	println(width)

	if real_data, ok := data.([]byte); ok {

		return (*cairo_surface)(C.cairo_image_surface_create_for_data((*C.uchar)(&real_data[0]),
			C.cairo_format_t(cairo_format),
			C.int(width),
			C.int(height),
			C.int(stride)))
	}

	return &simulated_surface{data: data, width: width, height: height, stride: stride}
}

func FormatStrideForWidth(cairo_format Format, width int) int {
	return int(C.cairo_format_stride_for_width(C.cairo_format_t(cairo_format), C.int(width)))
}

//export cairocallback_cairo_destroy_func
func cairocallback_cairo_destroy_func(ptr unsafe.Pointer) {
	var id = uintptr(ptr)

	// TODO mutex
	var fun = surface_user_datas[id]
	delete(surface_user_datas, id)
	// TODO mutex

	fun()
}

var surface_user_datas map[uintptr]func()

func (surface *cairo_surface) SetUserData(data func()) {

	// TODO mutex

	if surface_user_datas == nil {
		surface_user_datas = make(map[uintptr]func())
	}

	var id = uintptr(len(surface_user_datas))

	surface_user_datas[id] = data

	// TODO mutex

	C._cairo_surface_set_user_data((*C.struct__cairo_surface)(surface), unsafe.Pointer(id))

}

func (surface *cairo_surface_ref) SetUserData(data func()) {
}
func (surface *simulated_surface) SetUserData(data func()) {
}
func (surface *simulated_surface_ref) SetUserData(data func()) {
	surface.surf.SetUserData(data)
}
