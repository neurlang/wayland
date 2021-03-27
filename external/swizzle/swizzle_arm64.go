// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code copied from "golang.org/x/exp/shiny/driver/internal/swizzle"

package swizzle

const useBGRA32 = true
const useBGRA16 = false
const useBGRA4 = false

func bgra32([]byte)
func bgra16([]byte) { return }
func bgra4([]byte) { return }
