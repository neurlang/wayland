// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code copied from "golang.org/x/exp/shiny/driver/internal/swizzle"

// +build !amd64,!arm64

package swizzle

const (
	useBGRA32 = false
	useBGRA16 = false
	useBGRA4  = false
)

func bgra32(p []byte) { return }
func bgra16(p []byte) { return }
func bgra4(p []byte)  { return }
