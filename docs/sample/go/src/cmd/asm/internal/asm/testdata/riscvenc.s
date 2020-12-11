// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "../../../../../runtime/textflag.h"

TEXT asmtest(SB),DUPOK|NOSPLIT,$0
start:
	// Unprivileged ISA

	// 2.4: Integer Computational Instructions

	ADDI	$2047, X5, X6				// 1383f27f
	ADDI	$-2048, X5, X6				// 13830280
	ADDI	$2047, X5				// 9382f27f
	ADDI	$-2048, X5				// 93820280

	SLTI	$55, X5, X7				// 93a37203
	SLTIU	$55, X5, X7				// 93b37203

	ANDI	$1, X5, X6				// 13f31200
	ANDI	$1, X5					// 93f21200
	ORI	$1, X5, X6				// 13e31200
	ORI	$1, X5					// 93e21200
	XORI	$1, X5, X6				// 13c31200
	XORI	$1, X5					// 93c21200

	SLLI	$1, X5, X6				// 13931200
	SLLI	$1, X5					// 93921200
	SRLI	$1, X5, X6				// 13d31200
	SRLI	$1, X5					// 93d21200
	SRAI	$1, X5, X6				// 13d31240
	SRAI	$1, X5					// 93d21240

	ADD	X6, X5, X7				// b3836200
	ADD	X5, X6					// 33035300
	ADD	$2047, X5, X6				// 1383f27f
	ADD	$-2048, X5, X6				// 13830280
	ADD	$2047, X5				// 9382f27f
	ADD	$-2048, X5				// 93820280

	SLT	X6, X5, X7				// b3a36200
	SLT	$55, X5, X7				// 93a37203
	SLTU	X6, X5, X7				// b3b36200
	SLTU	$55, X5, X7				// 93b37203

	AND	X6, X5, X7				// b3f36200
	AND	X5, X6					// 33735300
	AND	$1, X5, X6				// 13f31200
	AND	$1, X5					// 93f21200
	OR	X6, X5, X7				// b3e36200
	OR	X5, X6					// 33635300
	OR	$1, X5, X6				// 13e31200
	OR	$1, X5					// 93e21200
	XOR	X6, X5, X7				// b3c36200
	XOR	X5, X6					// 33435300
	XOR	$1, X5, X6				// 13c31200
	XOR	$1, X5					// 93c21200

	AUIPC	$0, X10					// 17050000
	AUIPC	$0, X11					// 97050000
	AUIPC	$1, X10					// 17150000
	AUIPC	$-524288, X15				// 97070080
	AUIPC	$524287, X10				// 17f5ff7f

	LUI	$0, X15					// b7070000
	LUI	$167, X15				// b7770a00
	LUI	$-524288, X15				// b7070080
	LUI	$524287, X15				// b7f7ff7f

	SLL	X6, X5, X7				// b3936200
	SLL	X5, X6					// 33135300
	SLL	$1, X5, X6				// 13931200
	SLL	$1, X5					// 93921200
	SRL	X6, X5, X7				// b3d36200
	SRL	X5, X6					// 33535300
	SRL	$1, X5, X6				// 13d31200
	SRL	$1, X5					// 93d21200

	SUB	X6, X5, X7				// b3836240
	SUB	X5, X6					// 33035340

	SRA	X6, X5, X7				// b3d36240
	SRA	X5, X6					// 33535340
	SRA	$1, X5, X6				// 13d31240
	SRA	$1, X5					// 93d21240

	// 2.5: Control Transfer Instructions

	// These jumps and branches get printed as a jump or branch
	// to 2 because they transfer control to the second instruction
	// in the function (the first instruction being an invisible
	// stack pointer adjustment).
	JAL	X5, start	// JAL	X5, 2		// eff25ff0
	JALR	X6, (X5)				// 67830200
	JALR	X6, 4(X5)				// 67834200
	BEQ	X5, X6, start	// BEQ	X5, X6, 2	// e38c62ee
	BNE	X5, X6, start	// BNE	X5, X6, 2	// e39a62ee
	BLT	X5, X6, start	// BLT	X5, X6, 2	// e3c862ee
	BLTU	X5, X6, start	// BLTU	X5, X6, 2	// e3e662ee
	BGE	X5, X6, start	// BGE	X5, X6, 2	// e3d462ee
	BGEU	X5, X6, start	// BGEU	X5, X6, 2	// e3f262ee

	// 2.6: Load and Store Instructions
	LW	(X5), X6				// 03a30200
	LW	4(X5), X6				// 03a34200
	LWU	(X5), X6				// 03e30200
	LWU	4(X5), X6				// 03e34200
	LH	(X5), X6				// 03930200
	LH	4(X5), X6				// 03934200
	LHU	(X5), X6				// 03d30200
	LHU	4(X5), X6				// 03d34200
	LB	(X5), X6				// 03830200
	LB	4(X5), X6				// 03834200
	LBU	(X5), X6				// 03c30200
	LBU	4(X5), X6				// 03c34200

	SW	X5, (X6)				// 23205300
	SW	X5, 4(X6)				// 23225300
	SH	X5, (X6)				// 23105300
	SH	X5, 4(X6)				// 23125300
	SB	X5, (X6)				// 23005300
	SB	X5, 4(X6)				// 23025300

	// 5.2: Integer Computational Instructions (RV64I)
	ADDIW	$1, X5, X6				// 1b831200
	SLLIW	$1, X5, X6				// 1b931200
	SRLIW	$1, X5, X6				// 1bd31200
	SRAIW	$1, X5, X6				// 1bd31240
	ADDW	X5, X6, X7				// bb035300
	SLLW	X5, X6, X7				// bb135300
	SRLW	X5, X6, X7				// bb535300
	SUBW	X5, X6, X7				// bb035340
	SRAW	X5, X6, X7				// bb535340

	// 5.3: Load and Store Instructions (RV64I)
	LD	(X5), X6				// 03b30200
	LD	4(X5), X6				// 03b34200
	SD	X5, (X6)				// 23305300
	SD	X5, 4(X6)				// 23325300

	// 7.1: Multiplication Operations
	MUL	X5, X6, X7				// b3035302
	MULH	X5, X6, X7				// b3135302
	MULHU	X5, X6, X7				// b3335302
	MULHSU	X5, X6, X7				// b3235302
	MULW	X5, X6, X7				// bb035302
	DIV	X5, X6, X7				// b3435302
	DIVU	X5, X6, X7				// b3535302
	REM	X5, X6, X7				// b3635302
	REMU	X5, X6, X7				// b3735302
	DIVW	X5, X6, X7				// bb435302
	DIVUW	X5, X6, X7				// bb535302
	REMW	X5, X6, X7				// bb635302
	REMUW	X5, X6, X7				// bb735302

	// 10.1: Base Counters and Timers
	RDCYCLE		X5				// f32200c0
	RDTIME		X5				// f32210c0
	RDINSTRET	X5				// f32220c0

	// 11.5: Single-Precision Load and Store Instructions
	FLW	(X5), F0				// 07a00200
	FLW	4(X5), F0				// 07a04200
	FSW	F0, (X5)				// 27a00200
	FSW	F0, 4(X5)				// 27a20200

	// 11.6: Single-Precision Floating-Point Computational Instructions
	FADDS	F1, F0, F2				// 53011000
	FSUBS	F1, F0, F2				// 53011008
	FMULS	F1, F0, F2				// 53011010
	FDIVS	F1, F0, F2				// 53011018
	FMINS	F1, F0, F2				// 53011028
	FMAXS	F1, F0, F2				// 53111028
	FSQRTS	F0, F1					// d3000058

	// 11.7: Single-Precision Floating-Point Conversion and Move Instructions
	FCVTWS	F0, X5					// d31200c0
	FCVTLS	F0, X5					// d31220c0
	FCVTSW	X5, F0					// 538002d0
	FCVTSL	X5, F0					// 538022d0
	FCVTWUS	F0, X5					// d31210c0
	FCVTLUS	F0, X5					// d31230c0
	FCVTSWU	X5, F0					// 538012d0
	FCVTSLU	X5, F0					// 538032d0
	FSGNJS	F1, F0, F2				// 53011020
	FSGNJNS	F1, F0, F2				// 53111020
	FSGNJXS	F1, F0, F2				// 53211020
	FMVXS	F0, X5					// d30200e0
	FMVSX	X5, F0					// 538002f0
	FMVXW	F0, X5					// d30200e0
	FMVWX	X5, F0					// 538002f0

	// 11.8: Single-Precision Floating-Point Compare Instructions
	FEQS	F0, F1, X7				// d3a300a0
	FLTS	F0, F1, X7				// d39300a0
	FLES	F0, F1, X7				// d38300a0

	// 12.3: Double-Precision Load and Store Instructions
	FLD	(X5), F0				// 07b00200
	FLD	4(X5), F0				// 07b04200
	FSD	F0, (X5)				// 27b00200
	FSD	F0, 4(X5)				// 27b20200

	// 12.4: Double-Precision Floating-Point Computational Instructions
	FADDD	F1, F0, F2				// 53011002
	FSUBD	F1, F0, F2				// 5301100a
	FMULD	F1, F0, F2				// 53011012
	FDIVD	F1, F0, F2				// 5301101a
	FMIND	F1, F0, F2				// 5301102a
	FMAXD	F1, F0, F2				// 5311102a
	FSQRTD	F0, F1					// d300005a

	// 12.5: Double-Precision Floating-Point Conversion and Move Instructions
	FCVTWD	F0, X5					// d31200c2
	FCVTLD	F0, X5					// d31220c2
	FCVTDW	X5, F0					// 538002d2
	FCVTDL	X5, F0					// 538022d2
	FCVTWUD F0, X5					// d31210c2
	FCVTLUD F0, X5					// d31230c2
	FCVTDWU X5, F0					// 538012d2
	FCVTDLU X5, F0					// 538032d2
	FCVTSD	F0, F1					// d3001040
	FCVTDS	F0, F1					// d3000042
	FSGNJD	F1, F0, F2				// 53011022
	FSGNJND	F1, F0, F2				// 53111022
	FSGNJXD	F1, F0, F2				// 53211022
	FMVXD	F0, X5					// d30200e2
	FMVDX	X5, F0					// 538002f2

	// Privileged ISA

	// 3.2.1: Environment Call and Breakpoint
	ECALL						// 73000000
	SCALL						// 73000000
	EBREAK						// 73001000
	SBREAK						// 73001000

	// Arbitrary bytes (entered in little-endian mode)
	WORD	$0x12345678	// WORD $305419896	// 78563412
	WORD	$0x9abcdef0	// WORD $2596069104	// f0debc9a

	// MOV pseudo-instructions
	MOV	X5, X6					// 13830200
	MOV	$2047, X5				// 9b02f07f
	MOV	$-2048, X5				// 9b020080

	MOV	(X5), X6				// 03b30200
	MOV	4(X5), X6				// 03b34200
	MOVB	(X5), X6				// 03830200
	MOVB	4(X5), X6				// 03834200
	MOVH	(X5), X6				// 03930200
	MOVH	4(X5), X6				// 03934200
	MOVW	(X5), X6				// 03a30200
	MOVW	4(X5), X6				// 03a34200
	MOV	X5, (X6)				// 23305300
	MOV	X5, 4(X6)				// 23325300
	MOVB	X5, (X6)				// 23005300
	MOVB	X5, 4(X6)				// 23025300
	MOVH	X5, (X6)				// 23105300
	MOVH	X5, 4(X6)				// 23125300
	MOVW	X5, (X6)				// 23205300
	MOVW	X5, 4(X6)				// 23225300

	MOVF	4(X5), F0				// 07a04200
	MOVF	F0, 4(X5)				// 27a20200
	MOVF	F0, F1					// d3000020

	MOVD	4(X5), F0				// 07b04200
	MOVD	F0, 4(X5)				// 27b20200
	MOVD	F0, F1					// d3000022

	// These jumps can get printed as jumps to 2 because they go to the
	// second instruction in the function (the first instruction is an
	// invisible stack pointer adjustment).
	JMP	start		// JMP	2		// 6ff0dfcc
	JMP	(X5)					// 67800200
	JMP	4(X5)					// 67804200

	// JMP and CALL to symbol are encoded as:
	//	AUIPC $0, TMP
	//	JALR $0, TMP
	// with a R_RISCV_PCREL_ITYPE relocation - the linker resolves the
	// real address and updates the immediates for both instructions.
	CALL	asmtest(SB)				// 970f0000
	JMP	asmtest(SB)				// 970f0000

	SEQZ	X15, X15				// 93b71700
	SNEZ	X15, X15				// b337f000

	// F extension
	FNEGS	F0, F1					// d3100020

	// TODO(jsing): FNES gets encoded as FEQS+XORI - this should
	// be handled as a single *obj.Prog so that the full two
	// instruction encoding is tested here.
	FNES	F0, F1, X7				// d3a300a0

	// D extension
	FNEGD	F0, F1					// d3100022
	FEQD	F0, F1, X5				// d3a200a2
	FLTD	F0, F1, X5				// d39200a2
	FLED	F0, F1, X5				// d38200a2

	// TODO(jsing): FNED gets encoded as FEQD+XORI - this should
	// be handled as a single *obj.Prog so that the full two
	// instruction encoding is tested here.
	FNED	F0, F1, X5				// d3a200a2
