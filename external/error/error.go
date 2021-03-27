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

// Package error implements an external error
package error

import "errors"

type externalError interface {
	External() error
}

// External extracts the external (outermost) error cause from within the
// error returned, or nil in case the error is not external.
func GetExternal(e error) error {
	if err, ok := e.(externalError); ok {
		return err.External()
	}
	return nil
}

// Is can be used to check if the external or the internal error
// of a specific error is a specific target. It does a deep traversal.
func Is(specificError, target error) (result bool) {
	for !result && specificError != nil {
		result = errors.Is(specificError, target)
		specificError = GetExternal(specificError)
	}
	return result
}
