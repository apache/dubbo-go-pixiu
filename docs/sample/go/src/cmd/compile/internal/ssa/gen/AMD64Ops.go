// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "strings"

// Notes:
//  - Integer types live in the low portion of registers. Upper portions are junk.
//  - Boolean types use the low-order byte of a register. 0=false, 1=true.
//    Upper bytes are junk.
//  - Floating-point types live in the low natural slot of an sse2 register.
//    Unused portions are junk.
//  - We do not use AH,BH,CH,DH registers.
//  - When doing sub-register operations, we try to write the whole
//    destination register to avoid a partial-register write.
//  - Unused portions of AuxInt (or the Val portion of ValAndOff) are
//    filled by sign-extending the used portion.  Users of AuxInt which interpret
//    AuxInt as unsigned (e.g. shifts) must be careful.
//  - All SymOff opcodes require their offset to fit in an int32.

// Suffixes encode the bit width of various instructions.
// Q (quad word) = 64 bit
// L (long word) = 32 bit
// W (word)      = 16 bit
// B (byte)      = 8 bit

// copied from ../../amd64/reg.go
var regNamesAMD64 = []string{
	"AX",
	"CX",
	"DX",
	"BX",
	"SP",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
	"X0",
	"X1",
	"X2",
	"X3",
	"X4",
	"X5",
	"X6",
	"X7",
	"X8",
	"X9",
	"X10",
	"X11",
	"X12",
	"X13",
	"X14",
	"X15",

	// If you add registers, update asyncPreempt in runtime

	// pseudo-registers
	"SB",
}

func init() {
	// Make map from reg names to reg integers.
	if len(regNamesAMD64) > 64 {
		panic("too many registers")
	}
	num := map[string]int{}
	for i, name := range regNamesAMD64 {
		num[name] = i
	}
	buildReg := func(s string) regMask {
		m := regMask(0)
		for _, r := range strings.Split(s, " ") {
			if n, ok := num[r]; ok {
				m |= regMask(1) << uint(n)
				continue
			}
			panic("register " + r + " not found")
		}
		return m
	}

	// Common individual register masks
	var (
		ax         = buildReg("AX")
		cx         = buildReg("CX")
		dx         = buildReg("DX")
		bx         = buildReg("BX")
		gp         = buildReg("AX CX DX BX BP SI DI R8 R9 R10 R11 R12 R13 R14 R15")
		fp         = buildReg("X0 X1 X2 X3 X4 X5 X6 X7 X8 X9 X10 X11 X12 X13 X14 X15")
		gpsp       = gp | buildReg("SP")
		gpspsb     = gpsp | buildReg("SB")
		callerSave = gp | fp
	)
	// Common slices of register masks
	var (
		gponly = []regMask{gp}
		fponly = []regMask{fp}
	)

	// Common regInfo
	var (
		gp01           = regInfo{inputs: nil, outputs: gponly}
		gp11           = regInfo{inputs: []regMask{gp}, outputs: gponly}
		gp11sp         = regInfo{inputs: []regMask{gpsp}, outputs: gponly}
		gp11sb         = regInfo{inputs: []regMask{gpspsb}, outputs: gponly}
		gp21           = regInfo{inputs: []regMask{gp, gp}, outputs: gponly}
		gp21sp         = regInfo{inputs: []regMask{gpsp, gp}, outputs: gponly}
		gp21sb         = regInfo{inputs: []regMask{gpspsb, gpsp}, outputs: gponly}
		gp21shift      = regInfo{inputs: []regMask{gp, cx}, outputs: []regMask{gp}}
		gp11div        = regInfo{inputs: []regMask{ax, gpsp &^ dx}, outputs: []regMask{ax, dx}}
		gp21hmul       = regInfo{inputs: []regMask{ax, gpsp}, outputs: []regMask{dx}, clobbers: ax}
		gp21flags      = regInfo{inputs: []regMask{gp, gp}, outputs: []regMask{gp, 0}}
		gp2flags1flags = regInfo{inputs: []regMask{gp, gp, 0}, outputs: []regMask{gp, 0}}

		gp2flags     = regInfo{inputs: []regMask{gpsp, gpsp}}
		gp1flags     = regInfo{inputs: []regMask{gpsp}}
		gp0flagsLoad = regInfo{inputs: []regMask{gpspsb, 0}}
		gp1flagsLoad = regInfo{inputs: []regMask{gpspsb, gpsp, 0}}
		flagsgp      = regInfo{inputs: nil, outputs: gponly}

		gp11flags      = regInfo{inputs: []regMask{gp}, outputs: []regMask{gp, 0}}
		gp1flags1flags = regInfo{inputs: []regMask{gp, 0}, outputs: []regMask{gp, 0}}

		readflags = regInfo{inputs: nil, outputs: gponly}
		flagsgpax = regInfo{inputs: nil, clobbers: ax, outputs: []regMask{gp &^ ax}}

		gpload    = regInfo{inputs: []regMask{gpspsb, 0}, outputs: gponly}
		gp21load  = regInfo{inputs: []regMask{gp, gpspsb, 0}, outputs: gponly}
		gploadidx = regInfo{inputs: []regMask{gpspsb, gpsp, 0}, outputs: gponly}
		gp21pax   = regInfo{inputs: []regMask{gp &^ ax, gp}, outputs: []regMask{gp &^ ax}, clobbers: ax}

		gpstore         = regInfo{inputs: []regMask{gpspsb, gpsp, 0}}
		gpstoreconst    = regInfo{inputs: []regMask{gpspsb, 0}}
		gpstoreidx      = regInfo{inputs: []regMask{gpspsb, gpsp, gpsp, 0}}
		gpstoreconstidx = regInfo{inputs: []regMask{gpspsb, gpsp, 0}}
		gpstorexchg     = regInfo{inputs: []regMask{gp, gpspsb, 0}, outputs: []regMask{gp}}
		cmpxchg         = regInfo{inputs: []regMask{gp, ax, gp, 0}, outputs: []regMask{gp, 0}, clobbers: ax}

		fp01     = regInfo{inputs: nil, outputs: fponly}
		fp21     = regInfo{inputs: []regMask{fp, fp}, outputs: fponly}
		fp31     = regInfo{inputs: []regMask{fp, fp, fp}, outputs: fponly}
		fp21load = regInfo{inputs: []regMask{fp, gpspsb, 0}, outputs: fponly}
		fpgp     = regInfo{inputs: fponly, outputs: gponly}
		gpfp     = regInfo{inputs: gponly, outputs: fponly}
		fp11     = regInfo{inputs: fponly, outputs: fponly}
		fp2flags = regInfo{inputs: []regMask{fp, fp}}

		fpload    = regInfo{inputs: []regMask{gpspsb, 0}, outputs: fponly}
		fploadidx = regInfo{inputs: []regMask{gpspsb, gpsp, 0}, outputs: fponly}

		fpstore    = regInfo{inputs: []regMask{gpspsb, fp, 0}}
		fpstoreidx = regInfo{inputs: []regMask{gpspsb, gpsp, fp, 0}}
	)

	var AMD64ops = []opData{
		// fp ops
		{name: "ADDSS", argLength: 2, reg: fp21, asm: "ADDSS", commutative: true, resultInArg0: true}, // fp32 add
		{name: "ADDSD", argLength: 2, reg: fp21, asm: "ADDSD", commutative: true, resultInArg0: true}, // fp64 add
		{name: "SUBSS", argLength: 2, reg: fp21, asm: "SUBSS", resultInArg0: true},                    // fp32 sub
		{name: "SUBSD", argLength: 2, reg: fp21, asm: "SUBSD", resultInArg0: true},                    // fp64 sub
		{name: "MULSS", argLength: 2, reg: fp21, asm: "MULSS", commutative: true, resultInArg0: true}, // fp32 mul
		{name: "MULSD", argLength: 2, reg: fp21, asm: "MULSD", commutative: true, resultInArg0: true}, // fp64 mul
		{name: "DIVSS", argLength: 2, reg: fp21, asm: "DIVSS", resultInArg0: true},                    // fp32 div
		{name: "DIVSD", argLength: 2, reg: fp21, asm: "DIVSD", resultInArg0: true},                    // fp64 div

		{name: "MOVSSload", argLength: 2, reg: fpload, asm: "MOVSS", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"}, // fp32 load
		{name: "MOVSDload", argLength: 2, reg: fpload, asm: "MOVSD", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"}, // fp64 load
		{name: "MOVSSconst", reg: fp01, asm: "MOVSS", aux: "Float32", rematerializeable: true},                               // fp32 constant
		{name: "MOVSDconst", reg: fp01, asm: "MOVSD", aux: "Float64", rematerializeable: true},                               // fp64 constant
		{name: "MOVSSloadidx1", argLength: 3, reg: fploadidx, asm: "MOVSS", scale: 1, aux: "SymOff", symEffect: "Read"},      // fp32 load indexed by i
		{name: "MOVSSloadidx4", argLength: 3, reg: fploadidx, asm: "MOVSS", scale: 4, aux: "SymOff", symEffect: "Read"},      // fp32 load indexed by 4*i
		{name: "MOVSDloadidx1", argLength: 3, reg: fploadidx, asm: "MOVSD", scale: 1, aux: "SymOff", symEffect: "Read"},      // fp64 load indexed by i
		{name: "MOVSDloadidx8", argLength: 3, reg: fploadidx, asm: "MOVSD", scale: 8, aux: "SymOff", symEffect: "Read"},      // fp64 load indexed by 8*i

		{name: "MOVSSstore", argLength: 3, reg: fpstore, asm: "MOVSS", aux: "SymOff", faultOnNilArg0: true, symEffect: "Write"}, // fp32 store
		{name: "MOVSDstore", argLength: 3, reg: fpstore, asm: "MOVSD", aux: "SymOff", faultOnNilArg0: true, symEffect: "Write"}, // fp64 store
		{name: "MOVSSstoreidx1", argLength: 4, reg: fpstoreidx, asm: "MOVSS", scale: 1, aux: "SymOff", symEffect: "Write"},      // fp32 indexed by i store
		{name: "MOVSSstoreidx4", argLength: 4, reg: fpstoreidx, asm: "MOVSS", scale: 4, aux: "SymOff", symEffect: "Write"},      // fp32 indexed by 4i store
		{name: "MOVSDstoreidx1", argLength: 4, reg: fpstoreidx, asm: "MOVSD", scale: 1, aux: "SymOff", symEffect: "Write"},      // fp64 indexed by i store
		{name: "MOVSDstoreidx8", argLength: 4, reg: fpstoreidx, asm: "MOVSD", scale: 8, aux: "SymOff", symEffect: "Write"},      // fp64 indexed by 8i store

		{name: "ADDSSload", argLength: 3, reg: fp21load, asm: "ADDSS", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp32 arg0 + tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "ADDSDload", argLength: 3, reg: fp21load, asm: "ADDSD", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp64 arg0 + tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "SUBSSload", argLength: 3, reg: fp21load, asm: "SUBSS", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp32 arg0 - tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "SUBSDload", argLength: 3, reg: fp21load, asm: "SUBSD", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp64 arg0 - tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "MULSSload", argLength: 3, reg: fp21load, asm: "MULSS", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp32 arg0 * tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "MULSDload", argLength: 3, reg: fp21load, asm: "MULSD", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp64 arg0 * tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "DIVSSload", argLength: 3, reg: fp21load, asm: "DIVSS", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp32 arg0 / tmp, tmp loaded from arg1+auxint+aux, arg2 = mem
		{name: "DIVSDload", argLength: 3, reg: fp21load, asm: "DIVSD", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, symEffect: "Read"}, // fp64 arg0 / tmp, tmp loaded from arg1+auxint+aux, arg2 = mem

		// binary ops
		{name: "ADDQ", argLength: 2, reg: gp21sp, asm: "ADDQ", commutative: true, clobberFlags: true},                                                                   // arg0 + arg1
		{name: "ADDL", argLength: 2, reg: gp21sp, asm: "ADDL", commutative: true, clobberFlags: true},                                                                   // arg0 + arg1
		{name: "ADDQconst", argLength: 1, reg: gp11sp, asm: "ADDQ", aux: "Int32", typ: "UInt64", clobberFlags: true},                                                    // arg0 + auxint
		{name: "ADDLconst", argLength: 1, reg: gp11sp, asm: "ADDL", aux: "Int32", clobberFlags: true},                                                                   // arg0 + auxint
		{name: "ADDQconstmodify", argLength: 2, reg: gpstoreconst, asm: "ADDQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // add ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "ADDLconstmodify", argLength: 2, reg: gpstoreconst, asm: "ADDL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // add ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem

		{name: "SUBQ", argLength: 2, reg: gp21, asm: "SUBQ", resultInArg0: true, clobberFlags: true},                    // arg0 - arg1
		{name: "SUBL", argLength: 2, reg: gp21, asm: "SUBL", resultInArg0: true, clobberFlags: true},                    // arg0 - arg1
		{name: "SUBQconst", argLength: 1, reg: gp11, asm: "SUBQ", aux: "Int32", resultInArg0: true, clobberFlags: true}, // arg0 - auxint
		{name: "SUBLconst", argLength: 1, reg: gp11, asm: "SUBL", aux: "Int32", resultInArg0: true, clobberFlags: true}, // arg0 - auxint

		{name: "MULQ", argLength: 2, reg: gp21, asm: "IMULQ", commutative: true, resultInArg0: true, clobberFlags: true}, // arg0 * arg1
		{name: "MULL", argLength: 2, reg: gp21, asm: "IMULL", commutative: true, resultInArg0: true, clobberFlags: true}, // arg0 * arg1
		{name: "MULQconst", argLength: 1, reg: gp11, asm: "IMUL3Q", aux: "Int32", clobberFlags: true},                    // arg0 * auxint
		{name: "MULLconst", argLength: 1, reg: gp11, asm: "IMUL3L", aux: "Int32", clobberFlags: true},                    // arg0 * auxint

		{name: "MULLU", argLength: 2, reg: regInfo{inputs: []regMask{ax, gpsp}, outputs: []regMask{ax, 0}, clobbers: dx}, typ: "(UInt32,Flags)", asm: "MULL", commutative: true, clobberFlags: true}, // Let x = arg0*arg1 (full 32x32->64  unsigned multiply). Returns uint32(x), and flags set to overflow if uint32(x) != x.
		{name: "MULQU", argLength: 2, reg: regInfo{inputs: []regMask{ax, gpsp}, outputs: []regMask{ax, 0}, clobbers: dx}, typ: "(UInt64,Flags)", asm: "MULQ", commutative: true, clobberFlags: true}, // Let x = arg0*arg1 (full 64x64->128 unsigned multiply). Returns uint64(x), and flags set to overflow if uint64(x) != x.

		{name: "HMULQ", argLength: 2, reg: gp21hmul, commutative: true, asm: "IMULQ", clobberFlags: true}, // (arg0 * arg1) >> width
		{name: "HMULL", argLength: 2, reg: gp21hmul, commutative: true, asm: "IMULL", clobberFlags: true}, // (arg0 * arg1) >> width
		{name: "HMULQU", argLength: 2, reg: gp21hmul, commutative: true, asm: "MULQ", clobberFlags: true}, // (arg0 * arg1) >> width
		{name: "HMULLU", argLength: 2, reg: gp21hmul, commutative: true, asm: "MULL", clobberFlags: true}, // (arg0 * arg1) >> width

		{name: "AVGQU", argLength: 2, reg: gp21, commutative: true, resultInArg0: true, clobberFlags: true}, // (arg0 + arg1) / 2 as unsigned, all 64 result bits

		// For DIVQ, DIVL and DIVW, AuxInt non-zero means that the divisor has been proved to be not -1.
		{name: "DIVQ", argLength: 2, reg: gp11div, typ: "(Int64,Int64)", asm: "IDIVQ", aux: "Bool", clobberFlags: true}, // [arg0 / arg1, arg0 % arg1]
		{name: "DIVL", argLength: 2, reg: gp11div, typ: "(Int32,Int32)", asm: "IDIVL", aux: "Bool", clobberFlags: true}, // [arg0 / arg1, arg0 % arg1]
		{name: "DIVW", argLength: 2, reg: gp11div, typ: "(Int16,Int16)", asm: "IDIVW", aux: "Bool", clobberFlags: true}, // [arg0 / arg1, arg0 % arg1]

		{name: "DIVQU", argLength: 2, reg: gp11div, typ: "(UInt64,UInt64)", asm: "DIVQ", clobberFlags: true}, // [arg0 / arg1, arg0 % arg1]
		{name: "DIVLU", argLength: 2, reg: gp11div, typ: "(UInt32,UInt32)", asm: "DIVL", clobberFlags: true}, // [arg0 / arg1, arg0 % arg1]
		{name: "DIVWU", argLength: 2, reg: gp11div, typ: "(UInt16,UInt16)", asm: "DIVW", clobberFlags: true}, // [arg0 / arg1, arg0 % arg1]

		{name: "NEGLflags", argLength: 1, reg: gp11flags, typ: "(UInt32,Flags)", asm: "NEGL", resultInArg0: true}, // -arg0, flags set for 0-arg0.
		// The following 4 add opcodes return the low 64 bits of the sum in the first result and
		// the carry (the 65th bit) in the carry flag.
		{name: "ADDQcarry", argLength: 2, reg: gp21flags, typ: "(UInt64,Flags)", asm: "ADDQ", commutative: true, resultInArg0: true}, // r = arg0+arg1
		{name: "ADCQ", argLength: 3, reg: gp2flags1flags, typ: "(UInt64,Flags)", asm: "ADCQ", commutative: true, resultInArg0: true}, // r = arg0+arg1+carry(arg2)
		{name: "ADDQconstcarry", argLength: 1, reg: gp11flags, typ: "(UInt64,Flags)", asm: "ADDQ", aux: "Int32", resultInArg0: true}, // r = arg0+auxint
		{name: "ADCQconst", argLength: 2, reg: gp1flags1flags, typ: "(UInt64,Flags)", asm: "ADCQ", aux: "Int32", resultInArg0: true}, // r = arg0+auxint+carry(arg1)

		// The following 4 add opcodes return the low 64 bits of the difference in the first result and
		// the borrow (if the result is negative) in the carry flag.
		{name: "SUBQborrow", argLength: 2, reg: gp21flags, typ: "(UInt64,Flags)", asm: "SUBQ", resultInArg0: true},                    // r = arg0-arg1
		{name: "SBBQ", argLength: 3, reg: gp2flags1flags, typ: "(UInt64,Flags)", asm: "SBBQ", resultInArg0: true},                     // r = arg0-(arg1+carry(arg2))
		{name: "SUBQconstborrow", argLength: 1, reg: gp11flags, typ: "(UInt64,Flags)", asm: "SUBQ", aux: "Int32", resultInArg0: true}, // r = arg0-auxint
		{name: "SBBQconst", argLength: 2, reg: gp1flags1flags, typ: "(UInt64,Flags)", asm: "SBBQ", aux: "Int32", resultInArg0: true},  // r = arg0-(auxint+carry(arg1))

		{name: "MULQU2", argLength: 2, reg: regInfo{inputs: []regMask{ax, gpsp}, outputs: []regMask{dx, ax}}, commutative: true, asm: "MULQ", clobberFlags: true}, // arg0 * arg1, returns (hi, lo)
		{name: "DIVQU2", argLength: 3, reg: regInfo{inputs: []regMask{dx, ax, gpsp}, outputs: []regMask{ax, dx}}, asm: "DIVQ", clobberFlags: true},                // arg0:arg1 / arg2 (128-bit divided by 64-bit), returns (q, r)

		{name: "ANDQ", argLength: 2, reg: gp21, asm: "ANDQ", commutative: true, resultInArg0: true, clobberFlags: true},                                                 // arg0 & arg1
		{name: "ANDL", argLength: 2, reg: gp21, asm: "ANDL", commutative: true, resultInArg0: true, clobberFlags: true},                                                 // arg0 & arg1
		{name: "ANDQconst", argLength: 1, reg: gp11, asm: "ANDQ", aux: "Int32", resultInArg0: true, clobberFlags: true},                                                 // arg0 & auxint
		{name: "ANDLconst", argLength: 1, reg: gp11, asm: "ANDL", aux: "Int32", resultInArg0: true, clobberFlags: true},                                                 // arg0 & auxint
		{name: "ANDQconstmodify", argLength: 2, reg: gpstoreconst, asm: "ANDQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // and ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "ANDLconstmodify", argLength: 2, reg: gpstoreconst, asm: "ANDL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // and ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem

		{name: "ORQ", argLength: 2, reg: gp21, asm: "ORQ", commutative: true, resultInArg0: true, clobberFlags: true},                                                 // arg0 | arg1
		{name: "ORL", argLength: 2, reg: gp21, asm: "ORL", commutative: true, resultInArg0: true, clobberFlags: true},                                                 // arg0 | arg1
		{name: "ORQconst", argLength: 1, reg: gp11, asm: "ORQ", aux: "Int32", resultInArg0: true, clobberFlags: true},                                                 // arg0 | auxint
		{name: "ORLconst", argLength: 1, reg: gp11, asm: "ORL", aux: "Int32", resultInArg0: true, clobberFlags: true},                                                 // arg0 | auxint
		{name: "ORQconstmodify", argLength: 2, reg: gpstoreconst, asm: "ORQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // or ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "ORLconstmodify", argLength: 2, reg: gpstoreconst, asm: "ORL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // or ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem

		{name: "XORQ", argLength: 2, reg: gp21, asm: "XORQ", commutative: true, resultInArg0: true, clobberFlags: true},                                                 // arg0 ^ arg1
		{name: "XORL", argLength: 2, reg: gp21, asm: "XORL", commutative: true, resultInArg0: true, clobberFlags: true},                                                 // arg0 ^ arg1
		{name: "XORQconst", argLength: 1, reg: gp11, asm: "XORQ", aux: "Int32", resultInArg0: true, clobberFlags: true},                                                 // arg0 ^ auxint
		{name: "XORLconst", argLength: 1, reg: gp11, asm: "XORL", aux: "Int32", resultInArg0: true, clobberFlags: true},                                                 // arg0 ^ auxint
		{name: "XORQconstmodify", argLength: 2, reg: gpstoreconst, asm: "XORQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // xor ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "XORLconstmodify", argLength: 2, reg: gpstoreconst, asm: "XORL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // xor ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem

		{name: "CMPQ", argLength: 2, reg: gp2flags, asm: "CMPQ", typ: "Flags"},                    // arg0 compare to arg1
		{name: "CMPL", argLength: 2, reg: gp2flags, asm: "CMPL", typ: "Flags"},                    // arg0 compare to arg1
		{name: "CMPW", argLength: 2, reg: gp2flags, asm: "CMPW", typ: "Flags"},                    // arg0 compare to arg1
		{name: "CMPB", argLength: 2, reg: gp2flags, asm: "CMPB", typ: "Flags"},                    // arg0 compare to arg1
		{name: "CMPQconst", argLength: 1, reg: gp1flags, asm: "CMPQ", typ: "Flags", aux: "Int32"}, // arg0 compare to auxint
		{name: "CMPLconst", argLength: 1, reg: gp1flags, asm: "CMPL", typ: "Flags", aux: "Int32"}, // arg0 compare to auxint
		{name: "CMPWconst", argLength: 1, reg: gp1flags, asm: "CMPW", typ: "Flags", aux: "Int16"}, // arg0 compare to auxint
		{name: "CMPBconst", argLength: 1, reg: gp1flags, asm: "CMPB", typ: "Flags", aux: "Int8"},  // arg0 compare to auxint

		// compare *(arg0+auxint+aux) to arg1 (in that order). arg2=mem.
		{name: "CMPQload", argLength: 3, reg: gp1flagsLoad, asm: "CMPQ", aux: "SymOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},
		{name: "CMPLload", argLength: 3, reg: gp1flagsLoad, asm: "CMPL", aux: "SymOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},
		{name: "CMPWload", argLength: 3, reg: gp1flagsLoad, asm: "CMPW", aux: "SymOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},
		{name: "CMPBload", argLength: 3, reg: gp1flagsLoad, asm: "CMPB", aux: "SymOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},

		// compare *(arg0+ValAndOff(AuxInt).Off()+aux) to ValAndOff(AuxInt).Val() (in that order). arg1=mem.
		{name: "CMPQconstload", argLength: 2, reg: gp0flagsLoad, asm: "CMPQ", aux: "SymValAndOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},
		{name: "CMPLconstload", argLength: 2, reg: gp0flagsLoad, asm: "CMPL", aux: "SymValAndOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},
		{name: "CMPWconstload", argLength: 2, reg: gp0flagsLoad, asm: "CMPW", aux: "SymValAndOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},
		{name: "CMPBconstload", argLength: 2, reg: gp0flagsLoad, asm: "CMPB", aux: "SymValAndOff", typ: "Flags", symEffect: "Read", faultOnNilArg0: true},

		{name: "UCOMISS", argLength: 2, reg: fp2flags, asm: "UCOMISS", typ: "Flags"}, // arg0 compare to arg1, f32
		{name: "UCOMISD", argLength: 2, reg: fp2flags, asm: "UCOMISD", typ: "Flags"}, // arg0 compare to arg1, f64

		{name: "BTL", argLength: 2, reg: gp2flags, asm: "BTL", typ: "Flags"},                                           // test whether bit arg0%32 in arg1 is set
		{name: "BTQ", argLength: 2, reg: gp2flags, asm: "BTQ", typ: "Flags"},                                           // test whether bit arg0%64 in arg1 is set
		{name: "BTCL", argLength: 2, reg: gp21, asm: "BTCL", resultInArg0: true, clobberFlags: true},                   // complement bit arg1%32 in arg0
		{name: "BTCQ", argLength: 2, reg: gp21, asm: "BTCQ", resultInArg0: true, clobberFlags: true},                   // complement bit arg1%64 in arg0
		{name: "BTRL", argLength: 2, reg: gp21, asm: "BTRL", resultInArg0: true, clobberFlags: true},                   // reset bit arg1%32 in arg0
		{name: "BTRQ", argLength: 2, reg: gp21, asm: "BTRQ", resultInArg0: true, clobberFlags: true},                   // reset bit arg1%64 in arg0
		{name: "BTSL", argLength: 2, reg: gp21, asm: "BTSL", resultInArg0: true, clobberFlags: true},                   // set bit arg1%32 in arg0
		{name: "BTSQ", argLength: 2, reg: gp21, asm: "BTSQ", resultInArg0: true, clobberFlags: true},                   // set bit arg1%64 in arg0
		{name: "BTLconst", argLength: 1, reg: gp1flags, asm: "BTL", typ: "Flags", aux: "Int8"},                         // test whether bit auxint in arg0 is set, 0 <= auxint < 32
		{name: "BTQconst", argLength: 1, reg: gp1flags, asm: "BTQ", typ: "Flags", aux: "Int8"},                         // test whether bit auxint in arg0 is set, 0 <= auxint < 64
		{name: "BTCLconst", argLength: 1, reg: gp11, asm: "BTCL", resultInArg0: true, clobberFlags: true, aux: "Int8"}, // complement bit auxint in arg0, 0 <= auxint < 32
		{name: "BTCQconst", argLength: 1, reg: gp11, asm: "BTCQ", resultInArg0: true, clobberFlags: true, aux: "Int8"}, // complement bit auxint in arg0, 0 <= auxint < 64
		{name: "BTRLconst", argLength: 1, reg: gp11, asm: "BTRL", resultInArg0: true, clobberFlags: true, aux: "Int8"}, // reset bit auxint in arg0, 0 <= auxint < 32
		{name: "BTRQconst", argLength: 1, reg: gp11, asm: "BTRQ", resultInArg0: true, clobberFlags: true, aux: "Int8"}, // reset bit auxint in arg0, 0 <= auxint < 64
		{name: "BTSLconst", argLength: 1, reg: gp11, asm: "BTSL", resultInArg0: true, clobberFlags: true, aux: "Int8"}, // set bit auxint in arg0, 0 <= auxint < 32
		{name: "BTSQconst", argLength: 1, reg: gp11, asm: "BTSQ", resultInArg0: true, clobberFlags: true, aux: "Int8"}, // set bit auxint in arg0, 0 <= auxint < 64

		// direct bit operation on memory operand
		{name: "BTCQmodify", argLength: 3, reg: gpstore, asm: "BTCQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},     // complement bit arg1 in 64-bit arg0+auxint+aux, arg2=mem
		{name: "BTCLmodify", argLength: 3, reg: gpstore, asm: "BTCL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},     // complement bit arg1 in 32-bit arg0+auxint+aux, arg2=mem
		{name: "BTSQmodify", argLength: 3, reg: gpstore, asm: "BTSQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},     // set bit arg1 in 64-bit arg0+auxint+aux, arg2=mem
		{name: "BTSLmodify", argLength: 3, reg: gpstore, asm: "BTSL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},     // set bit arg1 in 32-bit arg0+auxint+aux, arg2=mem
		{name: "BTRQmodify", argLength: 3, reg: gpstore, asm: "BTRQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},     // reset bit arg1 in 64-bit arg0+auxint+aux, arg2=mem
		{name: "BTRLmodify", argLength: 3, reg: gpstore, asm: "BTRL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},     // reset bit arg1 in 32-bit arg0+auxint+aux, arg2=mem
		{name: "BTCQconstmodify", argLength: 2, reg: gpstoreconst, asm: "BTCQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // complement bit ValAndOff(AuxInt).Val() in 64-bit arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "BTCLconstmodify", argLength: 2, reg: gpstoreconst, asm: "BTCL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // complement bit ValAndOff(AuxInt).Val() in 32-bit arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "BTSQconstmodify", argLength: 2, reg: gpstoreconst, asm: "BTSQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // set bit ValAndOff(AuxInt).Val() in 64-bit arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "BTSLconstmodify", argLength: 2, reg: gpstoreconst, asm: "BTSL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // set bit ValAndOff(AuxInt).Val() in 32-bit arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "BTRQconstmodify", argLength: 2, reg: gpstoreconst, asm: "BTRQ", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // reset bit ValAndOff(AuxInt).Val() in 64-bit arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem
		{name: "BTRLconstmodify", argLength: 2, reg: gpstoreconst, asm: "BTRL", aux: "SymValAndOff", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // reset bit ValAndOff(AuxInt).Val() in 32-bit arg0+ValAndOff(AuxInt).Off()+aux, arg1=mem

		{name: "TESTQ", argLength: 2, reg: gp2flags, commutative: true, asm: "TESTQ", typ: "Flags"}, // (arg0 & arg1) compare to 0
		{name: "TESTL", argLength: 2, reg: gp2flags, commutative: true, asm: "TESTL", typ: "Flags"}, // (arg0 & arg1) compare to 0
		{name: "TESTW", argLength: 2, reg: gp2flags, commutative: true, asm: "TESTW", typ: "Flags"}, // (arg0 & arg1) compare to 0
		{name: "TESTB", argLength: 2, reg: gp2flags, commutative: true, asm: "TESTB", typ: "Flags"}, // (arg0 & arg1) compare to 0
		{name: "TESTQconst", argLength: 1, reg: gp1flags, asm: "TESTQ", typ: "Flags", aux: "Int32"}, // (arg0 & auxint) compare to 0
		{name: "TESTLconst", argLength: 1, reg: gp1flags, asm: "TESTL", typ: "Flags", aux: "Int32"}, // (arg0 & auxint) compare to 0
		{name: "TESTWconst", argLength: 1, reg: gp1flags, asm: "TESTW", typ: "Flags", aux: "Int16"}, // (arg0 & auxint) compare to 0
		{name: "TESTBconst", argLength: 1, reg: gp1flags, asm: "TESTB", typ: "Flags", aux: "Int8"},  // (arg0 & auxint) compare to 0

		{name: "SHLQ", argLength: 2, reg: gp21shift, asm: "SHLQ", resultInArg0: true, clobberFlags: true},              // arg0 << arg1, shift amount is mod 64
		{name: "SHLL", argLength: 2, reg: gp21shift, asm: "SHLL", resultInArg0: true, clobberFlags: true},              // arg0 << arg1, shift amount is mod 32
		{name: "SHLQconst", argLength: 1, reg: gp11, asm: "SHLQ", aux: "Int8", resultInArg0: true, clobberFlags: true}, // arg0 << auxint, shift amount 0-63
		{name: "SHLLconst", argLength: 1, reg: gp11, asm: "SHLL", aux: "Int8", resultInArg0: true, clobberFlags: true}, // arg0 << auxint, shift amount 0-31
		// Note: x86 is weird, the 16 and 8 byte shifts still use all 5 bits of shift amount!

		{name: "SHRQ", argLength: 2, reg: gp21shift, asm: "SHRQ", resultInArg0: true, clobberFlags: true},              // unsigned arg0 >> arg1, shift amount is mod 64
		{name: "SHRL", argLength: 2, reg: gp21shift, asm: "SHRL", resultInArg0: true, clobberFlags: true},              // unsigned uint32(arg0) >> arg1, shift amount is mod 32
		{name: "SHRW", argLength: 2, reg: gp21shift, asm: "SHRW", resultInArg0: true, clobberFlags: true},              // unsigned uint16(arg0) >> arg1, shift amount is mod 32
		{name: "SHRB", argLength: 2, reg: gp21shift, asm: "SHRB", resultInArg0: true, clobberFlags: true},              // unsigned uint8(arg0) >> arg1, shift amount is mod 32
		{name: "SHRQconst", argLength: 1, reg: gp11, asm: "SHRQ", aux: "Int8", resultInArg0: true, clobberFlags: true}, // unsigned arg0 >> auxint, shift amount 0-63
		{name: "SHRLconst", argLength: 1, reg: gp11, asm: "SHRL", aux: "Int8", resultInArg0: true, clobberFlags: true}, // unsigned uint32(arg0) >> auxint, shift amount 0-31
		{name: "SHRWconst", argLength: 1, reg: gp11, asm: "SHRW", aux: "Int8", resultInArg0: true, clobberFlags: true}, // unsigned uint16(arg0) >> auxint, shift amount 0-15
		{name: "SHRBconst", argLength: 1, reg: gp11, asm: "SHRB", aux: "Int8", resultInArg0: true, clobberFlags: true}, // unsigned uint8(arg0) >> auxint, shift amount 0-7

		{name: "SARQ", argLength: 2, reg: gp21shift, asm: "SARQ", resultInArg0: true, clobberFlags: true},              // signed arg0 >> arg1, shift amount is mod 64
		{name: "SARL", argLength: 2, reg: gp21shift, asm: "SARL", resultInArg0: true, clobberFlags: true},              // signed int32(arg0) >> arg1, shift amount is mod 32
		{name: "SARW", argLength: 2, reg: gp21shift, asm: "SARW", resultInArg0: true, clobberFlags: true},              // signed int16(arg0) >> arg1, shift amount is mod 32
		{name: "SARB", argLength: 2, reg: gp21shift, asm: "SARB", resultInArg0: true, clobberFlags: true},              // signed int8(arg0) >> arg1, shift amount is mod 32
		{name: "SARQconst", argLength: 1, reg: gp11, asm: "SARQ", aux: "Int8", resultInArg0: true, clobberFlags: true}, // signed arg0 >> auxint, shift amount 0-63
		{name: "SARLconst", argLength: 1, reg: gp11, asm: "SARL", aux: "Int8", resultInArg0: true, clobberFlags: true}, // signed int32(arg0) >> auxint, shift amount 0-31
		{name: "SARWconst", argLength: 1, reg: gp11, asm: "SARW", aux: "Int8", resultInArg0: true, clobberFlags: true}, // signed int16(arg0) >> auxint, shift amount 0-15
		{name: "SARBconst", argLength: 1, reg: gp11, asm: "SARB", aux: "Int8", resultInArg0: true, clobberFlags: true}, // signed int8(arg0) >> auxint, shift amount 0-7

		{name: "ROLQ", argLength: 2, reg: gp21shift, asm: "ROLQ", resultInArg0: true, clobberFlags: true},              // arg0 rotate left arg1 bits.
		{name: "ROLL", argLength: 2, reg: gp21shift, asm: "ROLL", resultInArg0: true, clobberFlags: true},              // arg0 rotate left arg1 bits.
		{name: "ROLW", argLength: 2, reg: gp21shift, asm: "ROLW", resultInArg0: true, clobberFlags: true},              // arg0 rotate left arg1 bits.
		{name: "ROLB", argLength: 2, reg: gp21shift, asm: "ROLB", resultInArg0: true, clobberFlags: true},              // arg0 rotate left arg1 bits.
		{name: "RORQ", argLength: 2, reg: gp21shift, asm: "RORQ", resultInArg0: true, clobberFlags: true},              // arg0 rotate right arg1 bits.
		{name: "RORL", argLength: 2, reg: gp21shift, asm: "RORL", resultInArg0: true, clobberFlags: true},              // arg0 rotate right arg1 bits.
		{name: "RORW", argLength: 2, reg: gp21shift, asm: "RORW", resultInArg0: true, clobberFlags: true},              // arg0 rotate right arg1 bits.
		{name: "RORB", argLength: 2, reg: gp21shift, asm: "RORB", resultInArg0: true, clobberFlags: true},              // arg0 rotate right arg1 bits.
		{name: "ROLQconst", argLength: 1, reg: gp11, asm: "ROLQ", aux: "Int8", resultInArg0: true, clobberFlags: true}, // arg0 rotate left auxint, rotate amount 0-63
		{name: "ROLLconst", argLength: 1, reg: gp11, asm: "ROLL", aux: "Int8", resultInArg0: true, clobberFlags: true}, // arg0 rotate left auxint, rotate amount 0-31
		{name: "ROLWconst", argLength: 1, reg: gp11, asm: "ROLW", aux: "Int8", resultInArg0: true, clobberFlags: true}, // arg0 rotate left auxint, rotate amount 0-15
		{name: "ROLBconst", argLength: 1, reg: gp11, asm: "ROLB", aux: "Int8", resultInArg0: true, clobberFlags: true}, // arg0 rotate left auxint, rotate amount 0-7

		{name: "ADDLload", argLength: 3, reg: gp21load, asm: "ADDL", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 + tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "ADDQload", argLength: 3, reg: gp21load, asm: "ADDQ", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 + tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "SUBQload", argLength: 3, reg: gp21load, asm: "SUBQ", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 - tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "SUBLload", argLength: 3, reg: gp21load, asm: "SUBL", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 - tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "ANDLload", argLength: 3, reg: gp21load, asm: "ANDL", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 & tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "ANDQload", argLength: 3, reg: gp21load, asm: "ANDQ", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 & tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "ORQload", argLength: 3, reg: gp21load, asm: "ORQ", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"},   // arg0 | tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "ORLload", argLength: 3, reg: gp21load, asm: "ORL", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"},   // arg0 | tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "XORQload", argLength: 3, reg: gp21load, asm: "XORQ", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 ^ tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem
		{name: "XORLload", argLength: 3, reg: gp21load, asm: "XORL", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, symEffect: "Read"}, // arg0 ^ tmp, tmp loaded from  arg1+auxint+aux, arg2 = mem

		// direct binary-op on memory (read-modify-write)
		{name: "ADDQmodify", argLength: 3, reg: gpstore, asm: "ADDQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) += arg1, arg2=mem
		{name: "SUBQmodify", argLength: 3, reg: gpstore, asm: "SUBQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) -= arg1, arg2=mem
		{name: "ANDQmodify", argLength: 3, reg: gpstore, asm: "ANDQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) &= arg1, arg2=mem
		{name: "ORQmodify", argLength: 3, reg: gpstore, asm: "ORQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},   // *(arg0+auxint+aux) |= arg1, arg2=mem
		{name: "XORQmodify", argLength: 3, reg: gpstore, asm: "XORQ", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) ^= arg1, arg2=mem
		{name: "ADDLmodify", argLength: 3, reg: gpstore, asm: "ADDL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) += arg1, arg2=mem
		{name: "SUBLmodify", argLength: 3, reg: gpstore, asm: "SUBL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) -= arg1, arg2=mem
		{name: "ANDLmodify", argLength: 3, reg: gpstore, asm: "ANDL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) &= arg1, arg2=mem
		{name: "ORLmodify", argLength: 3, reg: gpstore, asm: "ORL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"},   // *(arg0+auxint+aux) |= arg1, arg2=mem
		{name: "XORLmodify", argLength: 3, reg: gpstore, asm: "XORL", aux: "SymOff", typ: "Mem", clobberFlags: true, faultOnNilArg0: true, symEffect: "Read,Write"}, // *(arg0+auxint+aux) ^= arg1, arg2=mem

		// unary ops
		{name: "NEGQ", argLength: 1, reg: gp11, asm: "NEGQ", resultInArg0: true, clobberFlags: true}, // -arg0
		{name: "NEGL", argLength: 1, reg: gp11, asm: "NEGL", resultInArg0: true, clobberFlags: true}, // -arg0

		{name: "NOTQ", argLength: 1, reg: gp11, asm: "NOTQ", resultInArg0: true, clobberFlags: true}, // ^arg0
		{name: "NOTL", argLength: 1, reg: gp11, asm: "NOTL", resultInArg0: true, clobberFlags: true}, // ^arg0

		// BS{F,R}Q returns a tuple [result, flags]
		// result is undefined if the input is zero.
		// flags are set to "equal" if the input is zero, "not equal" otherwise.
		// BS{F,R}L returns only the result.
		{name: "BSFQ", argLength: 1, reg: gp11flags, asm: "BSFQ", typ: "(UInt64,Flags)"},        // # of low-order zeroes in 64-bit arg
		{name: "BSFL", argLength: 1, reg: gp11, asm: "BSFL", typ: "UInt32", clobberFlags: true}, // # of low-order zeroes in 32-bit arg
		{name: "BSRQ", argLength: 1, reg: gp11flags, asm: "BSRQ", typ: "(UInt64,Flags)"},        // # of high-order zeroes in 64-bit arg
		{name: "BSRL", argLength: 1, reg: gp11, asm: "BSRL", typ: "UInt32", clobberFlags: true}, // # of high-order zeroes in 32-bit arg

		// CMOV instructions: 64, 32 and 16-bit sizes.
		// if arg2 encodes a true result, return arg1, else arg0
		{name: "CMOVQEQ", argLength: 3, reg: gp21, asm: "CMOVQEQ", resultInArg0: true},
		{name: "CMOVQNE", argLength: 3, reg: gp21, asm: "CMOVQNE", resultInArg0: true},
		{name: "CMOVQLT", argLength: 3, reg: gp21, asm: "CMOVQLT", resultInArg0: true},
		{name: "CMOVQGT", argLength: 3, reg: gp21, asm: "CMOVQGT", resultInArg0: true},
		{name: "CMOVQLE", argLength: 3, reg: gp21, asm: "CMOVQLE", resultInArg0: true},
		{name: "CMOVQGE", argLength: 3, reg: gp21, asm: "CMOVQGE", resultInArg0: true},
		{name: "CMOVQLS", argLength: 3, reg: gp21, asm: "CMOVQLS", resultInArg0: true},
		{name: "CMOVQHI", argLength: 3, reg: gp21, asm: "CMOVQHI", resultInArg0: true},
		{name: "CMOVQCC", argLength: 3, reg: gp21, asm: "CMOVQCC", resultInArg0: true},
		{name: "CMOVQCS", argLength: 3, reg: gp21, asm: "CMOVQCS", resultInArg0: true},

		{name: "CMOVLEQ", argLength: 3, reg: gp21, asm: "CMOVLEQ", resultInArg0: true},
		{name: "CMOVLNE", argLength: 3, reg: gp21, asm: "CMOVLNE", resultInArg0: true},
		{name: "CMOVLLT", argLength: 3, reg: gp21, asm: "CMOVLLT", resultInArg0: true},
		{name: "CMOVLGT", argLength: 3, reg: gp21, asm: "CMOVLGT", resultInArg0: true},
		{name: "CMOVLLE", argLength: 3, reg: gp21, asm: "CMOVLLE", resultInArg0: true},
		{name: "CMOVLGE", argLength: 3, reg: gp21, asm: "CMOVLGE", resultInArg0: true},
		{name: "CMOVLLS", argLength: 3, reg: gp21, asm: "CMOVLLS", resultInArg0: true},
		{name: "CMOVLHI", argLength: 3, reg: gp21, asm: "CMOVLHI", resultInArg0: true},
		{name: "CMOVLCC", argLength: 3, reg: gp21, asm: "CMOVLCC", resultInArg0: true},
		{name: "CMOVLCS", argLength: 3, reg: gp21, asm: "CMOVLCS", resultInArg0: true},

		{name: "CMOVWEQ", argLength: 3, reg: gp21, asm: "CMOVWEQ", resultInArg0: true},
		{name: "CMOVWNE", argLength: 3, reg: gp21, asm: "CMOVWNE", resultInArg0: true},
		{name: "CMOVWLT", argLength: 3, reg: gp21, asm: "CMOVWLT", resultInArg0: true},
		{name: "CMOVWGT", argLength: 3, reg: gp21, asm: "CMOVWGT", resultInArg0: true},
		{name: "CMOVWLE", argLength: 3, reg: gp21, asm: "CMOVWLE", resultInArg0: true},
		{name: "CMOVWGE", argLength: 3, reg: gp21, asm: "CMOVWGE", resultInArg0: true},
		{name: "CMOVWLS", argLength: 3, reg: gp21, asm: "CMOVWLS", resultInArg0: true},
		{name: "CMOVWHI", argLength: 3, reg: gp21, asm: "CMOVWHI", resultInArg0: true},
		{name: "CMOVWCC", argLength: 3, reg: gp21, asm: "CMOVWCC", resultInArg0: true},
		{name: "CMOVWCS", argLength: 3, reg: gp21, asm: "CMOVWCS", resultInArg0: true},

		// CMOV with floating point instructions. We need separate pseudo-op to handle
		// InvertFlags correctly, and to generate special code that handles NaN (unordered flag).
		// NOTE: the fact that CMOV*EQF here is marked to generate CMOV*NE is not a bug. See
		// code generation in amd64/ssa.go.
		{name: "CMOVQEQF", argLength: 3, reg: gp21pax, asm: "CMOVQNE", resultInArg0: true},
		{name: "CMOVQNEF", argLength: 3, reg: gp21, asm: "CMOVQNE", resultInArg0: true},
		{name: "CMOVQGTF", argLength: 3, reg: gp21, asm: "CMOVQHI", resultInArg0: true},
		{name: "CMOVQGEF", argLength: 3, reg: gp21, asm: "CMOVQCC", resultInArg0: true},
		{name: "CMOVLEQF", argLength: 3, reg: gp21pax, asm: "CMOVLNE", resultInArg0: true},
		{name: "CMOVLNEF", argLength: 3, reg: gp21, asm: "CMOVLNE", resultInArg0: true},
		{name: "CMOVLGTF", argLength: 3, reg: gp21, asm: "CMOVLHI", resultInArg0: true},
		{name: "CMOVLGEF", argLength: 3, reg: gp21, asm: "CMOVLCC", resultInArg0: true},
		{name: "CMOVWEQF", argLength: 3, reg: gp21pax, asm: "CMOVWNE", resultInArg0: true},
		{name: "CMOVWNEF", argLength: 3, reg: gp21, asm: "CMOVWNE", resultInArg0: true},
		{name: "CMOVWGTF", argLength: 3, reg: gp21, asm: "CMOVWHI", resultInArg0: true},
		{name: "CMOVWGEF", argLength: 3, reg: gp21, asm: "CMOVWCC", resultInArg0: true},

		{name: "BSWAPQ", argLength: 1, reg: gp11, asm: "BSWAPQ", resultInArg0: true, clobberFlags: true}, // arg0 swap bytes
		{name: "BSWAPL", argLength: 1, reg: gp11, asm: "BSWAPL", resultInArg0: true, clobberFlags: true}, // arg0 swap bytes

		// POPCNT instructions aren't guaranteed to be on the target platform (they are SSE4).
		// Any use must be preceded by a successful check of runtime.x86HasPOPCNT.
		{name: "POPCNTQ", argLength: 1, reg: gp11, asm: "POPCNTQ", clobberFlags: true}, // count number of set bits in arg0
		{name: "POPCNTL", argLength: 1, reg: gp11, asm: "POPCNTL", clobberFlags: true}, // count number of set bits in arg0

		{name: "SQRTSD", argLength: 1, reg: fp11, asm: "SQRTSD"}, // sqrt(arg0)

		// ROUNDSD instruction isn't guaranteed to be on the target platform (it is SSE4.1)
		// Any use must be preceded by a successful check of runtime.x86HasSSE41.
		{name: "ROUNDSD", argLength: 1, reg: fp11, aux: "Int8", asm: "ROUNDSD"}, // rounds arg0 depending on auxint, 1 means math.Floor, 2 Ceil, 3 Trunc

		// VFMADD231SD only exists on platforms with the FMA3 instruction set.
		// Any use must be preceded by a successful check of runtime.support_fma.
		{name: "VFMADD231SD", argLength: 3, reg: fp31, resultInArg0: true, asm: "VFMADD231SD"},

		{name: "SBBQcarrymask", argLength: 1, reg: flagsgp, asm: "SBBQ"}, // (int64)(-1) if carry is set, 0 if carry is clear.
		{name: "SBBLcarrymask", argLength: 1, reg: flagsgp, asm: "SBBL"}, // (int32)(-1) if carry is set, 0 if carry is clear.
		// Note: SBBW and SBBB are subsumed by SBBL

		{name: "SETEQ", argLength: 1, reg: readflags, asm: "SETEQ"}, // extract == condition from arg0
		{name: "SETNE", argLength: 1, reg: readflags, asm: "SETNE"}, // extract != condition from arg0
		{name: "SETL", argLength: 1, reg: readflags, asm: "SETLT"},  // extract signed < condition from arg0
		{name: "SETLE", argLength: 1, reg: readflags, asm: "SETLE"}, // extract signed <= condition from arg0
		{name: "SETG", argLength: 1, reg: readflags, asm: "SETGT"},  // extract signed > condition from arg0
		{name: "SETGE", argLength: 1, reg: readflags, asm: "SETGE"}, // extract signed >= condition from arg0
		{name: "SETB", argLength: 1, reg: readflags, asm: "SETCS"},  // extract unsigned < condition from arg0
		{name: "SETBE", argLength: 1, reg: readflags, asm: "SETLS"}, // extract unsigned <= condition from arg0
		{name: "SETA", argLength: 1, reg: readflags, asm: "SETHI"},  // extract unsigned > condition from arg0
		{name: "SETAE", argLength: 1, reg: readflags, asm: "SETCC"}, // extract unsigned >= condition from arg0
		{name: "SETO", argLength: 1, reg: readflags, asm: "SETOS"},  // extract if overflow flag is set from arg0
		// Variants that store result to memory
		{name: "SETEQstore", argLength: 3, reg: gpstoreconst, asm: "SETEQ", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // extract == condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETNEstore", argLength: 3, reg: gpstoreconst, asm: "SETNE", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // extract != condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETLstore", argLength: 3, reg: gpstoreconst, asm: "SETLT", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},  // extract signed < condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETLEstore", argLength: 3, reg: gpstoreconst, asm: "SETLE", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // extract signed <= condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETGstore", argLength: 3, reg: gpstoreconst, asm: "SETGT", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},  // extract signed > condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETGEstore", argLength: 3, reg: gpstoreconst, asm: "SETGE", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // extract signed >= condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETBstore", argLength: 3, reg: gpstoreconst, asm: "SETCS", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},  // extract unsigned < condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETBEstore", argLength: 3, reg: gpstoreconst, asm: "SETLS", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // extract unsigned <= condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETAstore", argLength: 3, reg: gpstoreconst, asm: "SETHI", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},  // extract unsigned > condition from arg1 to arg0+auxint+aux, arg2=mem
		{name: "SETAEstore", argLength: 3, reg: gpstoreconst, asm: "SETCC", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // extract unsigned >= condition from arg1 to arg0+auxint+aux, arg2=mem
		// Need different opcodes for floating point conditions because
		// any comparison involving a NaN is always FALSE and thus
		// the patterns for inverting conditions cannot be used.
		{name: "SETEQF", argLength: 1, reg: flagsgpax, asm: "SETEQ", clobberFlags: true}, // extract == condition from arg0
		{name: "SETNEF", argLength: 1, reg: flagsgpax, asm: "SETNE", clobberFlags: true}, // extract != condition from arg0
		{name: "SETORD", argLength: 1, reg: flagsgp, asm: "SETPC"},                       // extract "ordered" (No Nan present) condition from arg0
		{name: "SETNAN", argLength: 1, reg: flagsgp, asm: "SETPS"},                       // extract "unordered" (Nan present) condition from arg0

		{name: "SETGF", argLength: 1, reg: flagsgp, asm: "SETHI"},  // extract floating > condition from arg0
		{name: "SETGEF", argLength: 1, reg: flagsgp, asm: "SETCC"}, // extract floating >= condition from arg0

		{name: "MOVBQSX", argLength: 1, reg: gp11, asm: "MOVBQSX"}, // sign extend arg0 from int8 to int64
		{name: "MOVBQZX", argLength: 1, reg: gp11, asm: "MOVBLZX"}, // zero extend arg0 from int8 to int64
		{name: "MOVWQSX", argLength: 1, reg: gp11, asm: "MOVWQSX"}, // sign extend arg0 from int16 to int64
		{name: "MOVWQZX", argLength: 1, reg: gp11, asm: "MOVWLZX"}, // zero extend arg0 from int16 to int64
		{name: "MOVLQSX", argLength: 1, reg: gp11, asm: "MOVLQSX"}, // sign extend arg0 from int32 to int64
		{name: "MOVLQZX", argLength: 1, reg: gp11, asm: "MOVL"},    // zero extend arg0 from int32 to int64

		{name: "MOVLconst", reg: gp01, asm: "MOVL", typ: "UInt32", aux: "Int32", rematerializeable: true}, // 32 low bits of auxint
		{name: "MOVQconst", reg: gp01, asm: "MOVQ", typ: "UInt64", aux: "Int64", rematerializeable: true}, // auxint

		{name: "CVTTSD2SL", argLength: 1, reg: fpgp, asm: "CVTTSD2SL"}, // convert float64 to int32
		{name: "CVTTSD2SQ", argLength: 1, reg: fpgp, asm: "CVTTSD2SQ"}, // convert float64 to int64
		{name: "CVTTSS2SL", argLength: 1, reg: fpgp, asm: "CVTTSS2SL"}, // convert float32 to int32
		{name: "CVTTSS2SQ", argLength: 1, reg: fpgp, asm: "CVTTSS2SQ"}, // convert float32 to int64
		{name: "CVTSL2SS", argLength: 1, reg: gpfp, asm: "CVTSL2SS"},   // convert int32 to float32
		{name: "CVTSL2SD", argLength: 1, reg: gpfp, asm: "CVTSL2SD"},   // convert int32 to float64
		{name: "CVTSQ2SS", argLength: 1, reg: gpfp, asm: "CVTSQ2SS"},   // convert int64 to float32
		{name: "CVTSQ2SD", argLength: 1, reg: gpfp, asm: "CVTSQ2SD"},   // convert int64 to float64
		{name: "CVTSD2SS", argLength: 1, reg: fp11, asm: "CVTSD2SS"},   // convert float64 to float32
		{name: "CVTSS2SD", argLength: 1, reg: fp11, asm: "CVTSS2SD"},   // convert float32 to float64

		// Move values between int and float registers, with no conversion.
		// TODO: should we have generic versions of these?
		{name: "MOVQi2f", argLength: 1, reg: gpfp, typ: "Float64"}, // move 64 bits from int to float reg
		{name: "MOVQf2i", argLength: 1, reg: fpgp, typ: "UInt64"},  // move 64 bits from float to int reg
		{name: "MOVLi2f", argLength: 1, reg: gpfp, typ: "Float32"}, // move 32 bits from int to float reg
		{name: "MOVLf2i", argLength: 1, reg: fpgp, typ: "UInt32"},  // move 32 bits from float to int reg, zero extend

		{name: "PXOR", argLength: 2, reg: fp21, asm: "PXOR", commutative: true, resultInArg0: true}, // exclusive or, applied to X regs for float negation.

		{name: "LEAQ", argLength: 1, reg: gp11sb, asm: "LEAQ", aux: "SymOff", rematerializeable: true, symEffect: "Addr"},      // arg0 + auxint + offset encoded in aux
		{name: "LEAL", argLength: 1, reg: gp11sb, asm: "LEAL", aux: "SymOff", rematerializeable: true, symEffect: "Addr"},      // arg0 + auxint + offset encoded in aux
		{name: "LEAW", argLength: 1, reg: gp11sb, asm: "LEAW", aux: "SymOff", rematerializeable: true, symEffect: "Addr"},      // arg0 + auxint + offset encoded in aux
		{name: "LEAQ1", argLength: 2, reg: gp21sb, asm: "LEAQ", scale: 1, commutative: true, aux: "SymOff", symEffect: "Addr"}, // arg0 + arg1 + auxint + aux
		{name: "LEAL1", argLength: 2, reg: gp21sb, asm: "LEAL", scale: 1, commutative: true, aux: "SymOff", symEffect: "Addr"}, // arg0 + arg1 + auxint + aux
		{name: "LEAW1", argLength: 2, reg: gp21sb, asm: "LEAW", scale: 1, commutative: true, aux: "SymOff", symEffect: "Addr"}, // arg0 + arg1 + auxint + aux
		{name: "LEAQ2", argLength: 2, reg: gp21sb, asm: "LEAQ", scale: 2, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 2*arg1 + auxint + aux
		{name: "LEAL2", argLength: 2, reg: gp21sb, asm: "LEAL", scale: 2, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 2*arg1 + auxint + aux
		{name: "LEAW2", argLength: 2, reg: gp21sb, asm: "LEAW", scale: 2, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 2*arg1 + auxint + aux
		{name: "LEAQ4", argLength: 2, reg: gp21sb, asm: "LEAQ", scale: 4, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 4*arg1 + auxint + aux
		{name: "LEAL4", argLength: 2, reg: gp21sb, asm: "LEAL", scale: 4, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 4*arg1 + auxint + aux
		{name: "LEAW4", argLength: 2, reg: gp21sb, asm: "LEAW", scale: 4, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 4*arg1 + auxint + aux
		{name: "LEAQ8", argLength: 2, reg: gp21sb, asm: "LEAQ", scale: 8, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 8*arg1 + auxint + aux
		{name: "LEAL8", argLength: 2, reg: gp21sb, asm: "LEAL", scale: 8, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 8*arg1 + auxint + aux
		{name: "LEAW8", argLength: 2, reg: gp21sb, asm: "LEAW", scale: 8, aux: "SymOff", symEffect: "Addr"},                    // arg0 + 8*arg1 + auxint + aux
		// Note: LEAx{1,2,4,8} must not have OpSB as either argument.

		// auxint+aux == add auxint and the offset of the symbol in aux (if any) to the effective address
		{name: "MOVBload", argLength: 2, reg: gpload, asm: "MOVBLZX", aux: "SymOff", typ: "UInt8", faultOnNilArg0: true, symEffect: "Read"},  // load byte from arg0+auxint+aux. arg1=mem.  Zero extend.
		{name: "MOVBQSXload", argLength: 2, reg: gpload, asm: "MOVBQSX", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"},             // ditto, sign extend to int64
		{name: "MOVWload", argLength: 2, reg: gpload, asm: "MOVWLZX", aux: "SymOff", typ: "UInt16", faultOnNilArg0: true, symEffect: "Read"}, // load 2 bytes from arg0+auxint+aux. arg1=mem.  Zero extend.
		{name: "MOVWQSXload", argLength: 2, reg: gpload, asm: "MOVWQSX", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"},             // ditto, sign extend to int64
		{name: "MOVLload", argLength: 2, reg: gpload, asm: "MOVL", aux: "SymOff", typ: "UInt32", faultOnNilArg0: true, symEffect: "Read"},    // load 4 bytes from arg0+auxint+aux. arg1=mem.  Zero extend.
		{name: "MOVLQSXload", argLength: 2, reg: gpload, asm: "MOVLQSX", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"},             // ditto, sign extend to int64
		{name: "MOVQload", argLength: 2, reg: gpload, asm: "MOVQ", aux: "SymOff", typ: "UInt64", faultOnNilArg0: true, symEffect: "Read"},    // load 8 bytes from arg0+auxint+aux. arg1=mem
		{name: "MOVBstore", argLength: 3, reg: gpstore, asm: "MOVB", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},    // store byte in arg1 to arg0+auxint+aux. arg2=mem
		{name: "MOVWstore", argLength: 3, reg: gpstore, asm: "MOVW", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},    // store 2 bytes in arg1 to arg0+auxint+aux. arg2=mem
		{name: "MOVLstore", argLength: 3, reg: gpstore, asm: "MOVL", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},    // store 4 bytes in arg1 to arg0+auxint+aux. arg2=mem
		{name: "MOVQstore", argLength: 3, reg: gpstore, asm: "MOVQ", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},    // store 8 bytes in arg1 to arg0+auxint+aux. arg2=mem
		{name: "MOVOload", argLength: 2, reg: fpload, asm: "MOVUPS", aux: "SymOff", typ: "Int128", faultOnNilArg0: true, symEffect: "Read"},  // load 16 bytes from arg0+auxint+aux. arg1=mem
		{name: "MOVOstore", argLength: 3, reg: fpstore, asm: "MOVUPS", aux: "SymOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"},  // store 16 bytes in arg1 to arg0+auxint+aux. arg2=mem

		// indexed loads/stores
		{name: "MOVBloadidx1", argLength: 3, reg: gploadidx, commutative: true, asm: "MOVBLZX", scale: 1, aux: "SymOff", typ: "UInt8", symEffect: "Read"},  // load a byte from arg0+arg1+auxint+aux. arg2=mem
		{name: "MOVWloadidx1", argLength: 3, reg: gploadidx, commutative: true, asm: "MOVWLZX", scale: 1, aux: "SymOff", typ: "UInt16", symEffect: "Read"}, // load 2 bytes from arg0+arg1+auxint+aux. arg2=mem
		{name: "MOVWloadidx2", argLength: 3, reg: gploadidx, asm: "MOVWLZX", scale: 2, aux: "SymOff", typ: "UInt16", symEffect: "Read"},                    // load 2 bytes from arg0+2*arg1+auxint+aux. arg2=mem
		{name: "MOVLloadidx1", argLength: 3, reg: gploadidx, commutative: true, asm: "MOVL", scale: 1, aux: "SymOff", typ: "UInt32", symEffect: "Read"},    // load 4 bytes from arg0+arg1+auxint+aux. arg2=mem
		{name: "MOVLloadidx4", argLength: 3, reg: gploadidx, asm: "MOVL", scale: 4, aux: "SymOff", typ: "UInt32", symEffect: "Read"},                       // load 4 bytes from arg0+4*arg1+auxint+aux. arg2=mem
		{name: "MOVLloadidx8", argLength: 3, reg: gploadidx, asm: "MOVL", scale: 8, aux: "SymOff", typ: "UInt32", symEffect: "Read"},                       // load 4 bytes from arg0+8*arg1+auxint+aux. arg2=mem
		{name: "MOVQloadidx1", argLength: 3, reg: gploadidx, commutative: true, asm: "MOVQ", scale: 1, aux: "SymOff", typ: "UInt64", symEffect: "Read"},    // load 8 bytes from arg0+arg1+auxint+aux. arg2=mem
		{name: "MOVQloadidx8", argLength: 3, reg: gploadidx, asm: "MOVQ", scale: 8, aux: "SymOff", typ: "UInt64", symEffect: "Read"},                       // load 8 bytes from arg0+8*arg1+auxint+aux. arg2=mem
		// TODO: sign-extending indexed loads
		// TODO: mark the MOVXstoreidx1 ops as commutative.  Generates too many rewrite rules at the moment.
		{name: "MOVBstoreidx1", argLength: 4, reg: gpstoreidx, asm: "MOVB", scale: 1, aux: "SymOff", symEffect: "Write"}, // store byte in arg2 to arg0+arg1+auxint+aux. arg3=mem
		{name: "MOVWstoreidx1", argLength: 4, reg: gpstoreidx, asm: "MOVW", scale: 1, aux: "SymOff", symEffect: "Write"}, // store 2 bytes in arg2 to arg0+arg1+auxint+aux. arg3=mem
		{name: "MOVWstoreidx2", argLength: 4, reg: gpstoreidx, asm: "MOVW", scale: 2, aux: "SymOff", symEffect: "Write"}, // store 2 bytes in arg2 to arg0+2*arg1+auxint+aux. arg3=mem
		{name: "MOVLstoreidx1", argLength: 4, reg: gpstoreidx, asm: "MOVL", scale: 1, aux: "SymOff", symEffect: "Write"}, // store 4 bytes in arg2 to arg0+arg1+auxint+aux. arg3=mem
		{name: "MOVLstoreidx4", argLength: 4, reg: gpstoreidx, asm: "MOVL", scale: 4, aux: "SymOff", symEffect: "Write"}, // store 4 bytes in arg2 to arg0+4*arg1+auxint+aux. arg3=mem
		{name: "MOVLstoreidx8", argLength: 4, reg: gpstoreidx, asm: "MOVL", scale: 8, aux: "SymOff", symEffect: "Write"}, // store 4 bytes in arg2 to arg0+8*arg1+auxint+aux. arg3=mem
		{name: "MOVQstoreidx1", argLength: 4, reg: gpstoreidx, asm: "MOVQ", scale: 1, aux: "SymOff", symEffect: "Write"}, // store 8 bytes in arg2 to arg0+arg1+auxint+aux. arg3=mem
		{name: "MOVQstoreidx8", argLength: 4, reg: gpstoreidx, asm: "MOVQ", scale: 8, aux: "SymOff", symEffect: "Write"}, // store 8 bytes in arg2 to arg0+8*arg1+auxint+aux. arg3=mem
		// TODO: add size-mismatched indexed loads, like MOVBstoreidx4.

		// For storeconst ops, the AuxInt field encodes both
		// the value to store and an address offset of the store.
		// Cast AuxInt to a ValAndOff to extract Val and Off fields.
		{name: "MOVBstoreconst", argLength: 2, reg: gpstoreconst, asm: "MOVB", aux: "SymValAndOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // store low byte of ValAndOff(AuxInt).Val() to arg0+ValAndOff(AuxInt).Off()+aux.  arg1=mem
		{name: "MOVWstoreconst", argLength: 2, reg: gpstoreconst, asm: "MOVW", aux: "SymValAndOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // store low 2 bytes of ...
		{name: "MOVLstoreconst", argLength: 2, reg: gpstoreconst, asm: "MOVL", aux: "SymValAndOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // store low 4 bytes of ...
		{name: "MOVQstoreconst", argLength: 2, reg: gpstoreconst, asm: "MOVQ", aux: "SymValAndOff", typ: "Mem", faultOnNilArg0: true, symEffect: "Write"}, // store 8 bytes of ...

		{name: "MOVBstoreconstidx1", argLength: 3, reg: gpstoreconstidx, asm: "MOVB", scale: 1, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store low byte of ValAndOff(AuxInt).Val() to arg0+1*arg1+ValAndOff(AuxInt).Off()+aux.  arg2=mem
		{name: "MOVWstoreconstidx1", argLength: 3, reg: gpstoreconstidx, asm: "MOVW", scale: 1, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store low 2 bytes of ... arg1 ...
		{name: "MOVWstoreconstidx2", argLength: 3, reg: gpstoreconstidx, asm: "MOVW", scale: 2, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store low 2 bytes of ... 2*arg1 ...
		{name: "MOVLstoreconstidx1", argLength: 3, reg: gpstoreconstidx, asm: "MOVL", scale: 1, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store low 4 bytes of ... arg1 ...
		{name: "MOVLstoreconstidx4", argLength: 3, reg: gpstoreconstidx, asm: "MOVL", scale: 4, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store low 4 bytes of ... 4*arg1 ...
		{name: "MOVQstoreconstidx1", argLength: 3, reg: gpstoreconstidx, asm: "MOVQ", scale: 1, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store 8 bytes of ... arg1 ...
		{name: "MOVQstoreconstidx8", argLength: 3, reg: gpstoreconstidx, asm: "MOVQ", scale: 8, aux: "SymValAndOff", typ: "Mem", symEffect: "Write"}, // store 8 bytes of ... 8*arg1 ...

		// arg0 = pointer to start of memory to zero
		// arg1 = value to store (will always be zero)
		// arg2 = mem
		// auxint = # of bytes to zero
		// returns mem
		{
			name:      "DUFFZERO",
			aux:       "Int64",
			argLength: 3,
			reg: regInfo{
				inputs:   []regMask{buildReg("DI"), buildReg("X0")},
				clobbers: buildReg("DI"),
			},
			faultOnNilArg0: true,
		},
		{name: "MOVOconst", reg: regInfo{nil, 0, []regMask{fp}}, typ: "Int128", aux: "Int128", rematerializeable: true},

		// arg0 = address of memory to zero
		// arg1 = # of 8-byte words to zero
		// arg2 = value to store (will always be zero)
		// arg3 = mem
		// returns mem
		{
			name:      "REPSTOSQ",
			argLength: 4,
			reg: regInfo{
				inputs:   []regMask{buildReg("DI"), buildReg("CX"), buildReg("AX")},
				clobbers: buildReg("DI CX"),
			},
			faultOnNilArg0: true,
		},

		{name: "CALLstatic", argLength: 1, reg: regInfo{clobbers: callerSave}, aux: "SymOff", clobberFlags: true, call: true, symEffect: "None"},                          // call static function aux.(*obj.LSym).  arg0=mem, auxint=argsize, returns mem
		{name: "CALLclosure", argLength: 3, reg: regInfo{inputs: []regMask{gpsp, buildReg("DX"), 0}, clobbers: callerSave}, aux: "Int64", clobberFlags: true, call: true}, // call function via closure.  arg0=codeptr, arg1=closure, arg2=mem, auxint=argsize, returns mem
		{name: "CALLinter", argLength: 2, reg: regInfo{inputs: []regMask{gp}, clobbers: callerSave}, aux: "Int64", clobberFlags: true, call: true},                        // call fn by pointer.  arg0=codeptr, arg1=mem, auxint=argsize, returns mem

		// arg0 = destination pointer
		// arg1 = source pointer
		// arg2 = mem
		// auxint = offset from duffcopy symbol to call
		// returns memory
		{
			name:      "DUFFCOPY",
			aux:       "Int64",
			argLength: 3,
			reg: regInfo{
				inputs:   []regMask{buildReg("DI"), buildReg("SI")},
				clobbers: buildReg("DI SI X0"), // uses X0 as a temporary
			},
			clobberFlags:   true,
			faultOnNilArg0: true,
			faultOnNilArg1: true,
		},

		// arg0 = destination pointer
		// arg1 = source pointer
		// arg2 = # of 8-byte words to copy
		// arg3 = mem
		// returns memory
		{
			name:      "REPMOVSQ",
			argLength: 4,
			reg: regInfo{
				inputs:   []regMask{buildReg("DI"), buildReg("SI"), buildReg("CX")},
				clobbers: buildReg("DI SI CX"),
			},
			faultOnNilArg0: true,
			faultOnNilArg1: true,
		},

		// (InvertFlags (CMPQ a b)) == (CMPQ b a)
		// So if we want (SETL (CMPQ a b)) but we can't do that because a is a constant,
		// then we do (SETL (InvertFlags (CMPQ b a))) instead.
		// Rewrites will convert this to (SETG (CMPQ b a)).
		// InvertFlags is a pseudo-op which can't appear in assembly output.
		{name: "InvertFlags", argLength: 1}, // reverse direction of arg0

		// Pseudo-ops
		{name: "LoweredGetG", argLength: 1, reg: gp01}, // arg0=mem
		// Scheduler ensures LoweredGetClosurePtr occurs only in entry block,
		// and sorts it to the very beginning of the block to prevent other
		// use of DX (the closure pointer)
		{name: "LoweredGetClosurePtr", reg: regInfo{outputs: []regMask{buildReg("DX")}}, zeroWidth: true},
		// LoweredGetCallerPC evaluates to the PC to which its "caller" will return.
		// I.e., if f calls g "calls" getcallerpc,
		// the result should be the PC within f that g will return to.
		// See runtime/stubs.go for a more detailed discussion.
		{name: "LoweredGetCallerPC", reg: gp01, rematerializeable: true},
		// LoweredGetCallerSP returns the SP of the caller of the current function.
		{name: "LoweredGetCallerSP", reg: gp01, rematerializeable: true},
		//arg0=ptr,arg1=mem, returns void.  Faults if ptr is nil.
		{name: "LoweredNilCheck", argLength: 2, reg: regInfo{inputs: []regMask{gpsp}}, clobberFlags: true, nilCheck: true, faultOnNilArg0: true},
		// LoweredWB invokes runtime.gcWriteBarrier. arg0=destptr, arg1=srcptr, arg2=mem, aux=runtime.gcWriteBarrier
		// It saves all GP registers if necessary, but may clobber others.
		{name: "LoweredWB", argLength: 3, reg: regInfo{inputs: []regMask{buildReg("DI"), ax}, clobbers: callerSave &^ gp}, clobberFlags: true, aux: "Sym", symEffect: "None"},

		// There are three of these functions so that they can have three different register inputs.
		// When we check 0 <= c <= cap (A), then 0 <= b <= c (B), then 0 <= a <= b (C), we want the
		// default registers to match so we don't need to copy registers around unnecessarily.
		{name: "LoweredPanicBoundsA", argLength: 3, aux: "Int64", reg: regInfo{inputs: []regMask{dx, bx}}, typ: "Mem"}, // arg0=idx, arg1=len, arg2=mem, returns memory. AuxInt contains report code (see PanicBounds in generic.go).
		{name: "LoweredPanicBoundsB", argLength: 3, aux: "Int64", reg: regInfo{inputs: []regMask{cx, dx}}, typ: "Mem"}, // arg0=idx, arg1=len, arg2=mem, returns memory. AuxInt contains report code (see PanicBounds in generic.go).
		{name: "LoweredPanicBoundsC", argLength: 3, aux: "Int64", reg: regInfo{inputs: []regMask{ax, cx}}, typ: "Mem"}, // arg0=idx, arg1=len, arg2=mem, returns memory. AuxInt contains report code (see PanicBounds in generic.go).

		// Constant flag values. For any comparison, there are 5 possible
		// outcomes: the three from the signed total order (<,==,>) and the
		// three from the unsigned total order. The == cases overlap.
		// Note: there's a sixth "unordered" outcome for floating-point
		// comparisons, but we don't use such a beast yet.
		// These ops are for temporary use by rewrite rules. They
		// cannot appear in the generated assembly.
		{name: "FlagEQ"},     // equal
		{name: "FlagLT_ULT"}, // signed < and unsigned <
		{name: "FlagLT_UGT"}, // signed < and unsigned >
		{name: "FlagGT_UGT"}, // signed > and unsigned >
		{name: "FlagGT_ULT"}, // signed > and unsigned <

		// Atomic loads.  These are just normal loads but return <value,memory> tuples
		// so they can be properly ordered with other loads.
		// load from arg0+auxint+aux.  arg1=mem.
		{name: "MOVBatomicload", argLength: 2, reg: gpload, asm: "MOVB", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"},
		{name: "MOVLatomicload", argLength: 2, reg: gpload, asm: "MOVL", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"},
		{name: "MOVQatomicload", argLength: 2, reg: gpload, asm: "MOVQ", aux: "SymOff", faultOnNilArg0: true, symEffect: "Read"},

		// Atomic stores and exchanges.  Stores use XCHG to get the right memory ordering semantics.
		// store arg0 to arg1+auxint+aux, arg2=mem.
		// These ops return a tuple of <old contents of *(arg1+auxint+aux), memory>.
		// Note: arg0 and arg1 are backwards compared to MOVLstore (to facilitate resultInArg0)!
		{name: "XCHGB", argLength: 3, reg: gpstorexchg, asm: "XCHGB", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, hasSideEffects: true, symEffect: "RdWr"},
		{name: "XCHGL", argLength: 3, reg: gpstorexchg, asm: "XCHGL", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, hasSideEffects: true, symEffect: "RdWr"},
		{name: "XCHGQ", argLength: 3, reg: gpstorexchg, asm: "XCHGQ", aux: "SymOff", resultInArg0: true, faultOnNilArg1: true, hasSideEffects: true, symEffect: "RdWr"},

		// Atomic adds.
		// *(arg1+auxint+aux) += arg0.  arg2=mem.
		// Returns a tuple of <old contents of *(arg1+auxint+aux), memory>.
		// Note: arg0 and arg1 are backwards compared to MOVLstore (to facilitate resultInArg0)!
		{name: "XADDLlock", argLength: 3, reg: gpstorexchg, asm: "XADDL", typ: "(UInt32,Mem)", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, hasSideEffects: true, symEffect: "RdWr"},
		{name: "XADDQlock", argLength: 3, reg: gpstorexchg, asm: "XADDQ", typ: "(UInt64,Mem)", aux: "SymOff", resultInArg0: true, clobberFlags: true, faultOnNilArg1: true, hasSideEffects: true, symEffect: "RdWr"},
		{name: "AddTupleFirst32", argLength: 2}, // arg1=tuple <x,y>.  Returns <x+arg0,y>.
		{name: "AddTupleFirst64", argLength: 2}, // arg1=tuple <x,y>.  Returns <x+arg0,y>.

		// Compare and swap.
		// arg0 = pointer, arg1 = old value, arg2 = new value, arg3 = memory.
		// if *(arg0+auxint+aux) == arg1 {
		//   *(arg0+auxint+aux) = arg2
		//   return (true, memory)
		// } else {
		//   return (false, memory)
		// }
		// Note that these instructions also return the old value in AX, but we ignore it.
		// TODO: have these return flags instead of bool.  The current system generates:
		//    CMPXCHGQ ...
		//    SETEQ AX
		//    CMPB  AX, $0
		//    JNE ...
		// instead of just
		//    CMPXCHGQ ...
		//    JEQ ...
		// but we can't do that because memory-using ops can't generate flags yet
		// (flagalloc wants to move flag-generating instructions around).
		{name: "CMPXCHGLlock", argLength: 4, reg: cmpxchg, asm: "CMPXCHGL", aux: "SymOff", clobberFlags: true, faultOnNilArg0: true, hasSideEffects: true, symEffect: "RdWr"},
		{name: "CMPXCHGQlock", argLength: 4, reg: cmpxchg, asm: "CMPXCHGQ", aux: "SymOff", clobberFlags: true, faultOnNilArg0: true, hasSideEffects: true, symEffect: "RdWr"},

		// Atomic memory updates.
		{name: "ANDBlock", argLength: 3, reg: gpstore, asm: "ANDB", aux: "SymOff", clobberFlags: true, faultOnNilArg0: true, hasSideEffects: true, symEffect: "RdWr"}, // *(arg0+auxint+aux) &= arg1
		{name: "ORBlock", argLength: 3, reg: gpstore, asm: "ORB", aux: "SymOff", clobberFlags: true, faultOnNilArg0: true, hasSideEffects: true, symEffect: "RdWr"},   // *(arg0+auxint+aux) |= arg1
	}

	var AMD64blocks = []blockData{
		{name: "EQ", controls: 1},
		{name: "NE", controls: 1},
		{name: "LT", controls: 1},
		{name: "LE", controls: 1},
		{name: "GT", controls: 1},
		{name: "GE", controls: 1},
		{name: "OS", controls: 1},
		{name: "OC", controls: 1},
		{name: "ULT", controls: 1},
		{name: "ULE", controls: 1},
		{name: "UGT", controls: 1},
		{name: "UGE", controls: 1},
		{name: "EQF", controls: 1},
		{name: "NEF", controls: 1},
		{name: "ORD", controls: 1}, // FP, ordered comparison (parity zero)
		{name: "NAN", controls: 1}, // FP, unordered comparison (parity one)
	}

	archs = append(archs, arch{
		name:            "AMD64",
		pkg:             "cmd/internal/obj/x86",
		genfile:         "../../amd64/ssa.go",
		ops:             AMD64ops,
		blocks:          AMD64blocks,
		regnames:        regNamesAMD64,
		gpregmask:       gp,
		fpregmask:       fp,
		framepointerreg: int8(num["BP"]),
		linkreg:         -1, // not used
	})
}
