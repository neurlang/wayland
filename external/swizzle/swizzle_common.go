// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code copied from "golang.org/x/exp/shiny/driver/internal/swizzle"

// Package swizzle provides functions for converting between RGBA pixel
// formats.
package swizzle

import "errors"

// ErrSliceNot32Bit is an error returned by BGRA in case the slice is not a multiple of 4
var ErrSliceNot32Bit = errors.New("input slice length is not a multiple of 4")

// BGRA converts a pixel buffer between Go's RGBA and other systems' BGRA byte
// orders.
//
// It returns the error ErrSliceNot32Bit if the input slice length is not a multiple of 4.
func BGRA(p []byte) error {
	if len(p)%4 != 0 {
		return ErrSliceNot32Bit
	}

	// Use asm code for 32-, 16- or 4-byte chunks, if supported.
	if useBGRA32 {
		n := len(p) &^ (32 - 1)
		bgra32(p[:n])
		p = p[n:]
	} else if useBGRA16 {
		n := len(p) &^ (16 - 1)
		bgra16(p[:n])
		p = p[n:]
	} else if useBGRA4 {
		bgra4(p)
		return nil
	}

	for i := 0; i < len(p); i += 4 {
		p[i+0], p[i+2] = p[i+2], p[i+0]
	}
	return nil
}
