// asmcheck

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import "unsafe"

// This file contains code generation tests related to the comparison
// operators.

// -------------- //
//    Equality    //
// -------------- //

// Check that compare to constant string use 2/4/8 byte compares

func CompareString1(s string) bool {
	// amd64:`CMPW\t\(.*\), [$]`
	// arm64:`MOVHU\t\(.*\), [R]`,`CMPW\t[$]`
	// ppc64le:`MOVHZ\t\(.*\), [R]`,`CMPW\t.*, [$]`
	// s390x:`MOVHBR\t\(.*\), [R]`,`CMPW\t.*, [$]`
	return s == "xx"
}

func CompareString2(s string) bool {
	// amd64:`CMPL\t\(.*\), [$]`
	// arm64:`MOVWU\t\(.*\), [R]`,`CMPW\t.*, [R]`
	// ppc64le:`MOVWZ\t\(.*\), [R]`,`CMPW\t.*, [R]`
	// s390x:`MOVWBR\t\(.*\), [R]`,`CMPW\t.*, [$]`
	return s == "xxxx"
}

func CompareString3(s string) bool {
	// amd64:`CMPQ\t\(.*\), [A-Z]`
	// arm64:-`CMPW\t`
	// ppc64:-`CMPW\t`
	// ppc64le:-`CMPW\t`
	// s390x:-`CMPW\t`
	return s == "xxxxxxxx"
}

// Check that arrays compare use 2/4/8 byte compares

func CompareArray1(a, b [2]byte) bool {
	// amd64:`CMPW\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	// arm64:-`MOVBU\t`
	// ppc64le:-`MOVBZ\t`
	// s390x:-`MOVBZ\t`
	return a == b
}

func CompareArray2(a, b [3]uint16) bool {
	// amd64:`CMPL\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	// amd64:`CMPW\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	return a == b
}

func CompareArray3(a, b [3]int16) bool {
	// amd64:`CMPL\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	// amd64:`CMPW\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	return a == b
}

func CompareArray4(a, b [12]int8) bool {
	// amd64:`CMPQ\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	// amd64:`CMPL\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	return a == b
}

func CompareArray5(a, b [15]byte) bool {
	// amd64:`CMPQ\t""[.+_a-z0-9]+\(SP\), [A-Z]`
	return a == b
}

// This was a TODO in mapaccess1_faststr
func CompareArray6(a, b unsafe.Pointer) bool {
	// amd64:`CMPL\t\(.*\), [A-Z]`
	// arm64:`MOVWU\t\(.*\), [R]`,`CMPW\t.*, [R]`
	// ppc64le:`MOVWZ\t\(.*\), [R]`,`CMPW\t.*, [R]`
	// s390x:`MOVWBR\t\(.*\), [R]`,`CMPW\t.*, [R]`
	return *((*[4]byte)(a)) != *((*[4]byte)(b))
}

// -------------- //
//    Ordering    //
// -------------- //

// Test that LEAQ/ADDQconst are folded into SETx ops

func CmpFold(x uint32) bool {
	// amd64:`SETHI\t.*\(SP\)`
	return x > 4
}

// Test that direct comparisons with memory are generated when
// possible

func CmpMem1(p int, q *int) bool {
	// amd64:`CMPQ\t\(.*\), [A-Z]`
	return p < *q
}

func CmpMem2(p *int, q int) bool {
	// amd64:`CMPQ\t\(.*\), [A-Z]`
	return *p < q
}

func CmpMem3(p *int) bool {
	// amd64:`CMPQ\t\(.*\), [$]7`
	return *p < 7
}

func CmpMem4(p *int) bool {
	// amd64:`CMPQ\t\(.*\), [$]7`
	return 7 < *p
}

func CmpMem5(p **int) {
	// amd64:`CMPL\truntime.writeBarrier\(SB\), [$]0`
	*p = nil
}

func CmpMem6(a []int) int {
	// 386:`CMPL\s8\([A-Z]+\),`
	// amd64:`CMPQ\s16\([A-Z]+\),`
	if a[1] > a[2] {
		return 1
	} else {
		return 2
	}
}

// Check tbz/tbnz are generated when comparing against zero on arm64

func CmpZero1(a int32, ptr *int) {
	if a < 0 { // arm64:"TBZ"
		*ptr = 0
	}
}

func CmpZero2(a int64, ptr *int) {
	if a < 0 { // arm64:"TBZ"
		*ptr = 0
	}
}

func CmpZero3(a int32, ptr *int) {
	if a >= 0 { // arm64:"TBNZ"
		*ptr = 0
	}
}

func CmpZero4(a int64, ptr *int) {
	if a >= 0 { // arm64:"TBNZ"
		*ptr = 0
	}
}

func CmpToZero(a, b, d int32, e, f int64) int32 {
	// arm:`TST`,-`AND`
	// arm64:`TSTW`,-`AND`
	// 386:`TESTL`,-`ANDL`
	// amd64:`TESTL`,-`ANDL`
	c0 := a&b < 0
	// arm:`CMN`,-`ADD`
	// arm64:`CMNW`,-`ADD`
	c1 := a+b < 0
	// arm:`TEQ`,-`XOR`
	c2 := a^b < 0
	// arm64:`TST`,-`AND`
	// amd64:`TESTQ`,-`ANDQ`
	c3 := e&f < 0
	// arm64:`CMN`,-`ADD`
	c4 := e+f < 0
	// not optimized to single CMNW/CMN due to further use of b+d
	// arm64:`ADD`,-`CMNW`
	// arm:`ADD`,-`CMN`
	c5 := b+d == 0
	// not optimized to single TSTW/TST due to further use of a&d
	// arm64:`AND`,-`TSTW`
	// arm:`AND`,-`TST`
	// 386:`ANDL`
	c6 := a&d >= 0
	// arm64:`TST\sR[0-9]+<<3,\sR[0-9]+`
	c7 := e&(f<<3) < 0
	// arm64:`CMN\sR[0-9]+<<3,\sR[0-9]+`
	c8 := e+(f<<3) < 0
	if c0 {
		return 1
	} else if c1 {
		return 2
	} else if c2 {
		return 3
	} else if c3 {
		return 4
	} else if c4 {
		return 5
	} else if c5 {
		return b + d
	} else if c6 {
		return a & d
	} else if c7 {
		return 7
	} else if c8 {
		return 8
	} else {
		return 0
	}
}

func CmpLogicalToZero(a, b, c uint32, d, e uint64) uint64 {

	// ppc64:"ANDCC",-"CMPW"
	// ppc64le:"ANDCC",-"CMPW"
	// wasm:"I64Eqz",-"I32Eqz",-"I64ExtendI32U",-"I32WrapI64"
	if a&63 == 0 {
		return 1
	}

	// ppc64:"ANDCC",-"CMP"
	// ppc64le:"ANDCC",-"CMP"
	// wasm:"I64Eqz",-"I32Eqz",-"I64ExtendI32U",-"I32WrapI64"
	if d&255 == 0 {
		return 1
	}

	// ppc64:"ANDCC",-"CMP"
	// ppc64le:"ANDCC",-"CMP"
	// wasm:"I64Eqz",-"I32Eqz",-"I64ExtendI32U",-"I32WrapI64"
	if d&e == 0 {
		return 1
	}
	// ppc64:"ORCC",-"CMP"
	// ppc64le:"ORCC",-"CMP"
	// wasm:"I64Eqz",-"I32Eqz",-"I64ExtendI32U",-"I32WrapI64"
	if d|e == 0 {
		return 1
	}

	// ppc64:"XORCC",-"CMP"
	// ppc64le:"XORCC",-"CMP"
	// wasm:"I64Eqz","I32Eqz",-"I64ExtendI32U",-"I32WrapI64"
	if e^d == 0 {
		return 1
	}
	return 0
}
