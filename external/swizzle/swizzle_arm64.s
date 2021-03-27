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

// +build !codeanalysis

#include "textflag.h"

// func bgra32(p []byte)
TEXT ·bgra32(SB),NOSPLIT,$0-24
 WORD $0x910df003
 //add  x3, x0, #0x37c
 
 WORD $0xd1041004
 //sub  x4, x0, #0x104
 
l2:
 WORD $0x91020060
 //ADDW R0, R3, $0x80
 
 WORD $0xd503201f
 NOP
 
l3:
 WORD $0x39400002
 //ldrb w2, [x0]
 
 WORD $0x39400801
 //ldrb w1, [x0, #2]
 
 WORD $0x39000802
 //strb w2, [x0, #2]
 
 WORD $0x381fc401
 //strb w1, [x0], #-4
 
 WORD $0xeb03001f
 //cmp  x0, x3
 
 //54ffff61
 BNE l3
 //b.ne 4006c0
 
 WORD $0xd1020003
 //sub  x3, x0, #0x80
 
 WORD $0xeb04007f
 //cmp  x3, x4
 
 //54fffec1
 BNE l2
 //b.ne 4006b8
 
 //d65f03c0
 RET

