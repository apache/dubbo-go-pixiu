// asmcheck

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import "math"

var sink64 [8]float64

func approx(x float64) {
	// s390x:"FIDBR\t[$]6"
	// arm64:"FRINTPD"
	// ppc64:"FRIP"
	// ppc64le:"FRIP"
	// wasm:"F64Ceil"
	sink64[0] = math.Ceil(x)

	// s390x:"FIDBR\t[$]7"
	// arm64:"FRINTMD"
	// ppc64:"FRIM"
	// ppc64le:"FRIM"
	// wasm:"F64Floor"
	sink64[1] = math.Floor(x)

	// s390x:"FIDBR\t[$]1"
	// arm64:"FRINTAD"
	// ppc64:"FRIN"
	// ppc64le:"FRIN"
	sink64[2] = math.Round(x)

	// s390x:"FIDBR\t[$]5"
	// arm64:"FRINTZD"
	// ppc64:"FRIZ"
	// ppc64le:"FRIZ"
	// wasm:"F64Trunc"
	sink64[3] = math.Trunc(x)

	// s390x:"FIDBR\t[$]4"
	// arm64:"FRINTND"
	// wasm:"F64Nearest"
	sink64[4] = math.RoundToEven(x)
}

func sqrt(x float64) float64 {
	// amd64:"SQRTSD"
	// 386/387:"FSQRT" 386/sse2:"SQRTSD"
	// arm64:"FSQRTD"
	// arm/7:"SQRTD"
	// mips/hardfloat:"SQRTD" mips/softfloat:-"SQRTD"
	// mips64/hardfloat:"SQRTD" mips64/softfloat:-"SQRTD"
	// wasm:"F64Sqrt"
	return math.Sqrt(x)
}

// Check that it's using integer registers
func abs(x, y float64) {
	// amd64:"BTRQ\t[$]63"
	// arm64:"FABSD\t"
	// s390x:"LPDFR\t",-"MOVD\t"     (no integer load/store)
	// ppc64:"FABS\t"
	// ppc64le:"FABS\t"
	// wasm:"F64Abs"
	// arm/6:"ABSD\t"
	sink64[0] = math.Abs(x)

	// amd64:"BTRQ\t[$]63","PXOR"    (TODO: this should be BTSQ)
	// s390x:"LNDFR\t",-"MOVD\t"     (no integer load/store)
	// ppc64:"FNABS\t"
	// ppc64le:"FNABS\t"
	sink64[1] = -math.Abs(y)
}

// Check that it's using integer registers
func abs32(x float32) float32 {
	// s390x:"LPDFR",-"LDEBR",-"LEDBR"     (no float64 conversion)
	return float32(math.Abs(float64(x)))
}

// Check that it's using integer registers
func copysign(a, b, c float64) {
	// amd64:"BTRQ\t[$]63","ANDQ","ORQ"
	// s390x:"CPSDR",-"MOVD"         (no integer load/store)
	// ppc64:"FCPSGN"
	// ppc64le:"FCPSGN"
	// wasm:"F64Copysign"
	sink64[0] = math.Copysign(a, b)

	// amd64:"BTSQ\t[$]63"
	// s390x:"LNDFR\t",-"MOVD\t"     (no integer load/store)
	// ppc64:"FCPSGN"
	// ppc64le:"FCPSGN"
	// arm64:"ORR", -"AND"
	sink64[1] = math.Copysign(c, -1)

	// Like math.Copysign(c, -1), but with integer operations. Useful
	// for platforms that have a copysign opcode to see if it's detected.
	// s390x:"LNDFR\t",-"MOVD\t"     (no integer load/store)
	sink64[2] = math.Float64frombits(math.Float64bits(a) | 1<<63)

	// amd64:"ANDQ","ORQ"
	// s390x:"CPSDR\t",-"MOVD\t"     (no integer load/store)
	// ppc64:"FCPSGN"
	// ppc64le:"FCPSGN"
	sink64[3] = math.Copysign(-1, c)
}

func fma(x, y, z float64) float64 {
	// amd64:"VFMADD231SD"
	// arm/6:"FMULAD"
	// arm64:"FMADDD"
	// s390x:"FMADD"
	// ppc64:"FMADD"
	// ppc64le:"FMADD"
	return math.FMA(x, y, z)
}

func fromFloat64(f64 float64) uint64 {
	// amd64:"MOVQ\tX.*, [^X].*"
	// arm64:"FMOVD\tF.*, R.*"
	// ppc64:"MFVSRD"
	// ppc64le:"MFVSRD"
	return math.Float64bits(f64+1) + 1
}

func fromFloat32(f32 float32) uint32 {
	// amd64:"MOVL\tX.*, [^X].*"
	// arm64:"FMOVS\tF.*, R.*"
	return math.Float32bits(f32+1) + 1
}

func toFloat64(u64 uint64) float64 {
	// amd64:"MOVQ\t[^X].*, X.*"
	// arm64:"FMOVD\tR.*, F.*"
	// ppc64:"MTVSRD"
	// ppc64le:"MTVSRD"
	return math.Float64frombits(u64+1) + 1
}

func toFloat32(u32 uint32) float32 {
	// amd64:"MOVL\t[^X].*, X.*"
	// arm64:"FMOVS\tR.*, F.*"
	return math.Float32frombits(u32+1) + 1
}

// Test that comparisons with constants converted to float
// are evaluated at compile-time

func constantCheck64() bool {
	// amd64:"MOVB\t[$]0",-"FCMP",-"MOVB\t[$]1"
	// s390x:"MOV(B|BZ|D)\t[$]0,",-"FCMPU",-"MOV(B|BZ|D)\t[$]1,"
	return 0.5 == float64(uint32(1)) || 1.5 > float64(uint64(1<<63)) || math.NaN() == math.NaN()
}

func constantCheck32() bool {
	// amd64:"MOVB\t[$]1",-"FCMP",-"MOVB\t[$]0"
	// s390x:"MOV(B|BZ|D)\t[$]1,",-"FCMPU",-"MOV(B|BZ|D)\t[$]0,"
	return float32(0.5) <= float32(int64(1)) && float32(1.5) >= float32(int32(-1<<31)) && float32(math.NaN()) != float32(math.NaN())
}

// Test that integer constants are converted to floating point constants
// at compile-time

func constantConvert32(x float32) float32 {
	// amd64:"MOVSS\t[$]f32.3f800000\\(SB\\)"
	// s390x:"FMOVS\t[$]f32.3f800000\\(SB\\)"
	// ppc64:"FMOVS\t[$]f32.3f800000\\(SB\\)"
	// ppc64le:"FMOVS\t[$]f32.3f800000\\(SB\\)"
	// arm64:"FMOVS\t[$]\\(1.0\\)"
	if x > math.Float32frombits(0x3f800000) {
		return -x
	}
	return x
}

func constantConvertInt32(x uint32) uint32 {
	// amd64:-"MOVSS"
	// s390x:-"FMOVS"
	// ppc64:-"FMOVS"
	// ppc64le:-"FMOVS"
	// arm64:-"FMOVS"
	if x > math.Float32bits(1) {
		return -x
	}
	return x
}
