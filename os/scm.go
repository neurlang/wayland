// +build !aix,!darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package os

import "errors"

// SocketControlMessage
type SocketControlMessage struct{}
type Sockaddr struct{}

var ErrUnsupportedOS = errors.New("unsupported os")

// ParseSocketControlMessage calls a system call to parse a Socket Control Message
func ParseSocketControlMessage([]byte) (scms []SocketControlMessage, err error) {
	return nil, ErrUnsupportedOS
}

// ParseUnixRights calls a system call to parse unix rights
func ParseUnixRights(*SocketControlMessage) (fds []int, err error) {
	return nil, ErrUnsupportedOS
}

func fallocate(fd int, mode uint32, off int64, size int64) error {
	return ErrUnsupportedOS
}

// UnixRights calls a system call
func UnixRights(int) []byte {
	return nil
}

// Sendmsg sends information on fd using a Sendmsg system call
func Sendmsg(fd int, msg []byte, oob []byte, sockaddr Sockaddr, z int) error {
	return ErrUnsupportedOS
}

// Mmap calls the system call to map memory on a fd
func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return nil, ErrUnsupportedOS
}

// Mmap calls the system call to unmap memory
func Munmap(data []byte) error {
	return ErrUnsupportedOS
}

// Pages may be read
const ProtRead = 0x1

// Pages may be written
const ProtWrite = 0x2

// Share this mapping
const MapShared = 0x01
