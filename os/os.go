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

// Package os implements an operating system routines useful for graphics
package os

import "os"
import "unsafe"
import "syscall"
import "errors"
import "crypto/rand"
import "fmt"

const PROT_READ = syscall.PROT_READ
const PROT_WRITE = syscall.PROT_WRITE
const MAP_SHARED = syscall.MAP_SHARED

// Length carries length and a bitness of a slice
type Length int64

// Create a 8bit length
func Len8(n int) Length {
	return Length(n)*16 + 1
}

// Create a 32bit length
func Len32(n int) Length {
	return Length(n)*16 + 4
}

// What returns the bitness of a Length
func (l Length) What() byte {
	return byte(int64(l) & 15)
}

// Int returns the length as an integer
func (l Length) Int() int {
	return int(int64(l) / 16)
}

// Error returned by Mmap
var ErrUnsupportedBitness = errors.New("unsupported bitness")

// Mmap maps a 32bit or a 8bit slice
func Mmap(fd int, offset int64, length Length, prot int, flags int) (interface{}, error) {
	switch length.What() {
	case 1:
		return syscall.Mmap(fd, offset, length.Int(), prot, flags)
	case 4:
		return Mmap32(fd, offset, length.Int(), prot, flags)
	default:
		return nil, ErrUnsupportedBitness
	}
}

// Munmap unmaps a 32bit or a 8bit slice
func Munmap(buf interface{}) error {
	switch v := ((interface{})(buf)).(type) {
	case []byte:
		return syscall.Munmap(v)
	case []uint32:
		return Munmap32(v)
	}
	return nil
}

func mvetype(dst, src *interface{}) {
	*(*uintptr)(unsafe.Pointer(dst)) = *(*uintptr)(unsafe.Pointer(src))
}

// Mmap32: Like MMap but for uint32 array
func Mmap32(fd int, offset int64, length int, prot int, flags int) (ret []uint32, err error) {
	data, err := syscall.Mmap(fd, offset, length, prot, flags)
	if err != nil {
		return nil, err
	}
	var a, b interface{} = data[: len(data)/4 : cap(data)/4], ret
	mvetype(&a, &b)
	return a.([]uint32), nil
}

// Munmap32: Like MUnmap but for uint32 array
func Munmap32(arr []uint32) (err error) {
	var data []byte
	var a, b interface{} = data[: len(arr)*4 : cap(arr)*4], arr
	mvetype(&a, &b)
	return syscall.Munmap(a.([]byte))
}

// MkOsTemp: Golang version of the popular C library function call
// The string can contain the patern consistng of XXX that will be replaced
// with a high-entropy alphanumeric sequence, if you want more entropic string
// you can put more XXX (in multiples of 3 X) up to the recommended value of 27 X
// shorter sequence of XXX will make your MkOsTemp more prone to the failure
// the buffer tmpname will be overwritten by the high entropic buffer
// x1, x2, x3 are the three X characters we are replacing, it can be another
func MkOsTemp(tmpname []byte, flags int, x1 byte, x2 byte, x3 byte) (int, error) {
	const alphabet = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var randbuf [5 * 9]byte
	var randbuf_i byte

	for i := 3; i <= len(tmpname); i++ {
		if (tmpname[i-3] == x1) && (tmpname[i-2] == x2) && (tmpname[i-1] == x3) {

			if randbuf_i == 0 {

				rand.Read(randbuf[:])

				for o := 0; o < 9; o++ {
					n := (uint64(randbuf[5*o+0])) |
						(uint64(randbuf[5*o+1]) << 8) |
						(uint64(randbuf[5*o+2]) << 16) |
						(uint64(randbuf[5*o+3]) << 24) |
						(uint64(randbuf[5*o+4]) << 32)
					randbuf[3*o+0] = alphabet[n%61]
					n /= 61
					randbuf[3*o+1] = alphabet[n%61]
					n /= 61
					randbuf[3*o+2] = alphabet[n%61]
				}

			}

			for j := 0; j < 3; j++ {
				tmpname[i-3] = randbuf[randbuf_i]
				randbuf_i++
				i++
			}
			i--
			randbuf_i %= 27
		}
	}

	//println(string(tmpname))

	return syscall.Open(string(tmpname), (flags & ^syscall.O_ACCMODE)|os.O_RDWR|os.O_CREATE|os.O_EXCL, syscall.S_IRUSR|syscall.S_IWUSR)
}

// Creates Tmp file that will be cloexec. In case of the ErrUnlink error, the fd is valid.
func CreateTmpfileCloexec(tmpname string) (int, error) {

	var namebuf = []byte(tmpname)

	var fd, err = MkOsTemp(namebuf, syscall.O_CLOEXEC, 'X', 'X', 'X')
	if fd < 0 {
		return fd, fmt.Errorf("CreateTmpfileCloexec: fd=%d", fd)
	}
	if err != nil {
		return fd, fmt.Errorf("CreateTmpfileCloexec(%s): %w", namebuf, err)
	}

	if syscall.Unlink(string(namebuf)) != nil {
		return fd, ErrUnlink
	}

	return fd, nil
}

var ErrUnlink = errors.New("CreateTmpfileCloexec: unlink error")

// OsCreateAnonymousFile: in case of the ErrUnlink error, the fd is valid.
// The file just isn't anonymous and can't be deleted. You can either ignore the ErrUnlink
// error and proceed, but it is your responsibility to Close the fd.
// In case of other errors, the fd is not valid and does not need to be closed.
func CreateAnonymousFile(size int64) (fd int, err error) {

	const template = "/go-lang-shared-XXXXXXXXXXXXXXXXXXXXXXXXXXX"

	path := os.Getenv("XDG_RUNTIME_DIR")

	fd, err = CreateTmpfileCloexec(path + template)
	if err != nil && err != ErrUnlink {
		return fd, err
	}

	err2 := syscall.Fallocate(fd, 0, 0, size)
	if err2 != nil {
		syscall.Close(fd)
		return -1, err2
	}

	return fd, err
}

// Close the fd
func Close(fd int) error {
	return syscall.Close(fd)
}
