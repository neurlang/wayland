// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix

// FcntlInt performs a fcntl syscall on fd with the provided command and argument.
func FcntlInt(fd uintptr, cmd, arg int) (int, error) {
	return fcntl(int(fd), cmd, arg)
}
