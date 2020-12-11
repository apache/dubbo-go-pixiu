// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build mips64 mips64le

#include "textflag.h"

// bool cas(uint32 *ptr, uint32 old, uint32 new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT ·Cas(SB), NOSPLIT, $0-17
	MOVV	ptr+0(FP), R1
	MOVW	old+8(FP), R2
	MOVW	new+12(FP), R5
	SYNC
cas_again:
	MOVV	R5, R3
	LL	(R1), R4
	BNE	R2, R4, cas_fail
	SC	R3, (R1)
	BEQ	R3, cas_again
	MOVV	$1, R1
	MOVB	R1, ret+16(FP)
	SYNC
	RET
cas_fail:
	MOVV	$0, R1
	JMP	-4(PC)

// bool	cas64(uint64 *ptr, uint64 old, uint64 new)
// Atomically:
//	if(*val == *old){
//		*val = new;
//		return 1;
//	} else {
//		return 0;
//	}
TEXT ·Cas64(SB), NOSPLIT, $0-25
	MOVV	ptr+0(FP), R1
	MOVV	old+8(FP), R2
	MOVV	new+16(FP), R5
	SYNC
cas64_again:
	MOVV	R5, R3
	LLV	(R1), R4
	BNE	R2, R4, cas64_fail
	SCV	R3, (R1)
	BEQ	R3, cas64_again
	MOVV	$1, R1
	MOVB	R1, ret+24(FP)
	SYNC
	RET
cas64_fail:
	MOVV	$0, R1
	JMP	-4(PC)

TEXT ·Casuintptr(SB), NOSPLIT, $0-25
	JMP	·Cas64(SB)

TEXT ·CasRel(SB), NOSPLIT, $0-17
	JMP	·Cas(SB)

TEXT ·Loaduintptr(SB),  NOSPLIT|NOFRAME, $0-16
	JMP	·Load64(SB)

TEXT ·Loaduint(SB), NOSPLIT|NOFRAME, $0-16
	JMP	·Load64(SB)

TEXT ·Storeuintptr(SB), NOSPLIT, $0-16
	JMP	·Store64(SB)

TEXT ·Xadduintptr(SB), NOSPLIT, $0-24
	JMP	·Xadd64(SB)

TEXT ·Loadint64(SB), NOSPLIT, $0-16
	JMP	·Load64(SB)

TEXT ·Xaddint64(SB), NOSPLIT, $0-24
	JMP	·Xadd64(SB)

// bool casp(void **val, void *old, void *new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT ·Casp1(SB), NOSPLIT, $0-25
	JMP runtime∕internal∕atomic·Cas64(SB)

// uint32 xadd(uint32 volatile *ptr, int32 delta)
// Atomically:
//	*val += delta;
//	return *val;
TEXT ·Xadd(SB), NOSPLIT, $0-20
	MOVV	ptr+0(FP), R2
	MOVW	delta+8(FP), R3
	SYNC
	LL	(R2), R1
	ADDU	R1, R3, R4
	MOVV	R4, R1
	SC	R4, (R2)
	BEQ	R4, -4(PC)
	MOVW	R1, ret+16(FP)
	SYNC
	RET

TEXT ·Xadd64(SB), NOSPLIT, $0-24
	MOVV	ptr+0(FP), R2
	MOVV	delta+8(FP), R3
	SYNC
	LLV	(R2), R1
	ADDVU	R1, R3, R4
	MOVV	R4, R1
	SCV	R4, (R2)
	BEQ	R4, -4(PC)
	MOVV	R1, ret+16(FP)
	SYNC
	RET

TEXT ·Xchg(SB), NOSPLIT, $0-20
	MOVV	ptr+0(FP), R2
	MOVW	new+8(FP), R5

	SYNC
	MOVV	R5, R3
	LL	(R2), R1
	SC	R3, (R2)
	BEQ	R3, -3(PC)
	MOVW	R1, ret+16(FP)
	SYNC
	RET

TEXT ·Xchg64(SB), NOSPLIT, $0-24
	MOVV	ptr+0(FP), R2
	MOVV	new+8(FP), R5

	SYNC
	MOVV	R5, R3
	LLV	(R2), R1
	SCV	R3, (R2)
	BEQ	R3, -3(PC)
	MOVV	R1, ret+16(FP)
	SYNC
	RET

TEXT ·Xchguintptr(SB), NOSPLIT, $0-24
	JMP	·Xchg64(SB)

TEXT ·StorepNoWB(SB), NOSPLIT, $0-16
	JMP	·Store64(SB)

TEXT ·StoreRel(SB), NOSPLIT, $0-12
	JMP	·Store(SB)

TEXT ·Store(SB), NOSPLIT, $0-12
	MOVV	ptr+0(FP), R1
	MOVW	val+8(FP), R2
	SYNC
	MOVW	R2, 0(R1)
	SYNC
	RET

TEXT ·Store8(SB), NOSPLIT, $0-9
	MOVV	ptr+0(FP), R1
	MOVB	val+8(FP), R2
	SYNC
	MOVB	R2, 0(R1)
	SYNC
	RET

TEXT ·Store64(SB), NOSPLIT, $0-16
	MOVV	ptr+0(FP), R1
	MOVV	val+8(FP), R2
	SYNC
	MOVV	R2, 0(R1)
	SYNC
	RET

// void	Or8(byte volatile*, byte);
TEXT ·Or8(SB), NOSPLIT, $0-9
	MOVV	ptr+0(FP), R1
	MOVBU	val+8(FP), R2
	// Align ptr down to 4 bytes so we can use 32-bit load/store.
	MOVV	$~3, R3
	AND	R1, R3
	// Compute val shift.
#ifdef GOARCH_mips64
	// Big endian.  ptr = ptr ^ 3
	XOR	$3, R1
#endif
	// R4 = ((ptr & 3) * 8)
	AND	$3, R1, R4
	SLLV	$3, R4
	// Shift val for aligned ptr. R2 = val << R4
	SLLV	R4, R2

	SYNC
	LL	(R3), R4
	OR	R2, R4
	SC	R4, (R3)
	BEQ	R4, -4(PC)
	SYNC
	RET

// void	And8(byte volatile*, byte);
TEXT ·And8(SB), NOSPLIT, $0-9
	MOVV	ptr+0(FP), R1
	MOVBU	val+8(FP), R2
	// Align ptr down to 4 bytes so we can use 32-bit load/store.
	MOVV	$~3, R3
	AND	R1, R3
	// Compute val shift.
#ifdef GOARCH_mips64
	// Big endian.  ptr = ptr ^ 3
	XOR	$3, R1
#endif
	// R4 = ((ptr & 3) * 8)
	AND	$3, R1, R4
	SLLV	$3, R4
	// Shift val for aligned ptr. R2 = val << R4 | ^(0xFF << R4)
	MOVV	$0xFF, R5
	SLLV	R4, R2
	SLLV	R4, R5
	NOR	R0, R5
	OR	R5, R2

	SYNC
	LL	(R3), R4
	AND	R2, R4
	SC	R4, (R3)
	BEQ	R4, -4(PC)
	SYNC
	RET
