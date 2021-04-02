// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package os

import "golang.org/x/sys/unix"
import "syscall"

// SocketControlMessage
type SocketControlMessage = syscall.SocketControlMessage
type Sockaddr = unix.Sockaddr

// ParseSocketControlMessage calls a system call to parse a Socket Control Message
func ParseSocketControlMessage(b []byte) ([]SocketControlMessage, error) {
	scms, err := syscall.ParseSocketControlMessage(b)
	return []SocketControlMessage(scms), err
}

// ParseUnixRights calls a system call to parse unix rights
func ParseUnixRights(m *SocketControlMessage) (fds []int, err error) {
	return syscall.ParseUnixRights(m)
}

func fallocate(fd int, mode uint32, off int64, size int64) error {
	return syscall.Fallocate(fd, mode, off, size)
}

// UnixRights calls a system call
func UnixRights(fd int) []byte {
	return syscall.UnixRights(fd)
}

// Sendmsg sends information on fd using a Sendmsg system call
func Sendmsg(fd int, msg []byte, oob []byte, sockaddr Sockaddr, z int) error {
	return unix.Sendmsg(fd, msg, oob, sockaddr, z)
}

// Mmap calls the system call to map memory on a fd
func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return syscall.Mmap(fd, offset, length, prot, flags)
}

// Mmap calls the system call to unmap memory
func Munmap(data []byte) error {
	return syscall.Munmap(data)
}

// Pages may be read
const ProtRead = syscall.PROT_READ

// Pages may be written
const ProtWrite = syscall.PROT_WRITE

// Share this mapping
const MapShared = syscall.MAP_SHARED
