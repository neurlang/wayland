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

// Control interface for an epoll file descriptor
func Epollctl(epfd uintptr, op uintptr, fd uintptr, event uintptr, fdp *int32, r *Runner) {
	var ptr = (*unsafe.Pointer)(unsafe.Pointer(&fdp))
	*ptr = unsafe.Pointer(DoFlagRunner(r))

	syscall.RawSyscall6(syscall.SYS_EPOLL_CTL, epfd, op, fd, event, 0, 0)
}

// Type Runner must be Unflagged before Run is called, otherwise it will crash
// You can check if Runner is Flagged using the IsFlagRunner call
type Runner interface {
	Run(events uint32)
}

func DoFlagRunner(r *Runner) *Runner {
	return (*Runner)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) | 1))
}
func UnFlagRunner(r *Runner) *Runner {
	return (*Runner)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) & uintptr(0xFFFFFFFFFFFFFFFE)))
}
func IsFlagRunner(r *Runner) bool {
	return (uintptr(unsafe.Pointer(r)) & 1) == 1
}

const PROT_READ = syscall.PROT_READ
const PROT_WRITE = syscall.PROT_WRITE
const MAP_SHARED = syscall.MAP_SHARED

func mvetype(dst, src *interface{}) {
	*(*uintptr)(unsafe.Pointer(dst)) = *(*uintptr)(unsafe.Pointer(src))
}

// Length carries length and a bitness of a slice
type Length int64

// 8bit length
const LEN8 byte = 1

// 32bit length
const LEN32 byte = 4

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
var ErrUnsupportedBitness = errors.New("Unsupported Bitness")

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
// the string *tmpname will be overwritten by the high entropic string
func MkOsTemp(tmpname *string, flags int) (int, error) {

	var randbuf [27]byte
	rand.Read(randbuf[:])

	var randbuf_i byte

	const alphabet = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	const XXX = "XXX"

	for i := range randbuf {
		randbuf[i] = alphabet[randbuf[i]%byte(len(alphabet))]
	}

	for i := len(XXX); i <= len(*tmpname); i++ {
		if ((*tmpname)[i-len(XXX)] == 'X') && ((*tmpname)[i-len(XXX):i] == XXX) {
			*tmpname = (*tmpname)[0:i-len(XXX)] + string(randbuf[randbuf_i:randbuf_i+byte(len(XXX))]) + (*tmpname)[i:]
			randbuf_i += byte(len(XXX))
			if randbuf_i == byte(len(randbuf)) {
				break
			}
			i += len(XXX) - 1
		}
	}

	//println(*tmpname)

	return syscall.Open((*tmpname), (flags & ^syscall.O_ACCMODE)|os.O_RDWR|os.O_CREATE|os.O_EXCL, syscall.S_IRUSR|syscall.S_IWUSR)
}

// Creates Tmp file that will be cloexec. In case of the ErrUnlink error, the fd is valid.
func CreateTmpfileCloexec(tmpname string) (int, error) {
	var fd, err = MkOsTemp(&tmpname, syscall.O_CLOEXEC)
	if fd < 0 {
		return fd, fmt.Errorf("CreateTmpfileCloexec: fd=%d", fd)
	}
	if err != nil {
		return fd, fmt.Errorf("CreateTmpfileCloexec(%s): %w", tmpname, err)
	}

	if syscall.Unlink(tmpname) != nil {
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

func os_fd_set_cloexec(fd int) int {

	if fd == -1 {
		return -1
	}

	flags, err := fcntl(fd, syscall.F_GETFD, 0)
	if (flags == -1) || (err != nil) {
		return -1
	}

	flags, err = fcntl(fd, syscall.F_SETFD, flags|syscall.FD_CLOEXEC)
	if (flags == -1) || (err != nil) {
		return -1
	}

	return 0
}

func set_cloexec_or_close(fd int) int {
	if os_fd_set_cloexec(fd) != 0 {
		Close(fd)
		return -1
	}
	return fd
}

// Error returned by SocketpairCloexec
var ErrWrongFds = errors.New("SocketpairCloexec: wrong fds")

// Error returned by SocketpairCloexec
var ErrBadFds = errors.New("SocketpairCloexec: bad fds")

// can return ErrWrongFds, ErrBadFds, or other syscall error
func SocketpairCloexec(domain int, typ int, protocol int) (fdsv [2]int, err error) {
	fdsv, err = syscall.Socketpair(domain, typ|syscall.SOCK_CLOEXEC, protocol)
	if (err == nil) || (err != syscall.EINVAL) {
		return fdsv, err
	}

	fdsv, err = syscall.Socketpair(domain, typ, protocol)
	if err != nil {
		return fdsv, err
	}

	fdsv[0] = set_cloexec_or_close(fdsv[0])
	fdsv[1] = set_cloexec_or_close(fdsv[1])

	if fdsv[0] != -1 && fdsv[1] != -1 {
		return fdsv, ErrWrongFds
	}

	Close(fdsv[0])
	Close(fdsv[1])
	return fdsv, ErrBadFds
}

// Error returned by EpollCreateCloexec
var ErrWrongFd = errors.New("EpollCreateCloexec: wrong fd")

// can return ErrWrongFd or other syscall error
func EpollCreateCloexec() (fd int, err error) {

	fd, err = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if fd >= 0 {
		return fd, err
	}
	if err != syscall.EINVAL {
		return -1, err
	}

	fd, err = syscall.EpollCreate(1)
	fd = set_cloexec_or_close(fd)
	if fd == -1 {
		return fd, ErrWrongFd
	}
	return fd, nil
}
