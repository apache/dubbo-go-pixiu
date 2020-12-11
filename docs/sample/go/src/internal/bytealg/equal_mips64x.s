// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build mips64 mips64le

#include "go_asm.h"
#include "textflag.h"

#define	REGCTXT	R22

// memequal(a, b unsafe.Pointer, size uintptr) bool
TEXT runtime·memequal(SB),NOSPLIT|NOFRAME,$0-25
	MOVV	a+0(FP), R1
	MOVV	b+8(FP), R2
	BEQ	R1, R2, eq
	MOVV	size+16(FP), R3
	ADDV	R1, R3, R4
loop:
	BNE	R1, R4, test
	MOVV	$1, R1
	MOVB	R1, ret+24(FP)
	RET
test:
	MOVBU	(R1), R6
	ADDV	$1, R1
	MOVBU	(R2), R7
	ADDV	$1, R2
	BEQ	R6, R7, loop

	MOVB	R0, ret+24(FP)
	RET
eq:
	MOVV	$1, R1
	MOVB	R1, ret+24(FP)
	RET

// memequal_varlen(a, b unsafe.Pointer) bool
TEXT runtime·memequal_varlen(SB),NOSPLIT,$40-17
	MOVV	a+0(FP), R1
	MOVV	b+8(FP), R2
	BEQ	R1, R2, eq
	MOVV	8(REGCTXT), R3    // compiler stores size at offset 8 in the closure
	MOVV	R1, 8(R29)
	MOVV	R2, 16(R29)
	MOVV	R3, 24(R29)
	JAL	runtime·memequal(SB)
	MOVBU	32(R29), R1
	MOVB	R1, ret+16(FP)
	RET
eq:
	MOVV	$1, R1
	MOVB	R1, ret+16(FP)
	RET
