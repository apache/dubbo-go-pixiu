// Code generated from gen/PPC64.rules; DO NOT EDIT.
// generated with: cd gen; go run *.go

package ssa

import "math"
import "cmd/internal/objabi"
import "cmd/compile/internal/types"

func rewriteValuePPC64(v *Value) bool {
	switch v.Op {
	case OpAbs:
		return rewriteValuePPC64_OpAbs_0(v)
	case OpAdd16:
		return rewriteValuePPC64_OpAdd16_0(v)
	case OpAdd32:
		return rewriteValuePPC64_OpAdd32_0(v)
	case OpAdd32F:
		return rewriteValuePPC64_OpAdd32F_0(v)
	case OpAdd64:
		return rewriteValuePPC64_OpAdd64_0(v)
	case OpAdd64F:
		return rewriteValuePPC64_OpAdd64F_0(v)
	case OpAdd64carry:
		return rewriteValuePPC64_OpAdd64carry_0(v)
	case OpAdd8:
		return rewriteValuePPC64_OpAdd8_0(v)
	case OpAddPtr:
		return rewriteValuePPC64_OpAddPtr_0(v)
	case OpAddr:
		return rewriteValuePPC64_OpAddr_0(v)
	case OpAnd16:
		return rewriteValuePPC64_OpAnd16_0(v)
	case OpAnd32:
		return rewriteValuePPC64_OpAnd32_0(v)
	case OpAnd64:
		return rewriteValuePPC64_OpAnd64_0(v)
	case OpAnd8:
		return rewriteValuePPC64_OpAnd8_0(v)
	case OpAndB:
		return rewriteValuePPC64_OpAndB_0(v)
	case OpAtomicAdd32:
		return rewriteValuePPC64_OpAtomicAdd32_0(v)
	case OpAtomicAdd64:
		return rewriteValuePPC64_OpAtomicAdd64_0(v)
	case OpAtomicAnd8:
		return rewriteValuePPC64_OpAtomicAnd8_0(v)
	case OpAtomicCompareAndSwap32:
		return rewriteValuePPC64_OpAtomicCompareAndSwap32_0(v)
	case OpAtomicCompareAndSwap64:
		return rewriteValuePPC64_OpAtomicCompareAndSwap64_0(v)
	case OpAtomicCompareAndSwapRel32:
		return rewriteValuePPC64_OpAtomicCompareAndSwapRel32_0(v)
	case OpAtomicExchange32:
		return rewriteValuePPC64_OpAtomicExchange32_0(v)
	case OpAtomicExchange64:
		return rewriteValuePPC64_OpAtomicExchange64_0(v)
	case OpAtomicLoad32:
		return rewriteValuePPC64_OpAtomicLoad32_0(v)
	case OpAtomicLoad64:
		return rewriteValuePPC64_OpAtomicLoad64_0(v)
	case OpAtomicLoad8:
		return rewriteValuePPC64_OpAtomicLoad8_0(v)
	case OpAtomicLoadAcq32:
		return rewriteValuePPC64_OpAtomicLoadAcq32_0(v)
	case OpAtomicLoadPtr:
		return rewriteValuePPC64_OpAtomicLoadPtr_0(v)
	case OpAtomicOr8:
		return rewriteValuePPC64_OpAtomicOr8_0(v)
	case OpAtomicStore32:
		return rewriteValuePPC64_OpAtomicStore32_0(v)
	case OpAtomicStore64:
		return rewriteValuePPC64_OpAtomicStore64_0(v)
	case OpAtomicStore8:
		return rewriteValuePPC64_OpAtomicStore8_0(v)
	case OpAtomicStoreRel32:
		return rewriteValuePPC64_OpAtomicStoreRel32_0(v)
	case OpAvg64u:
		return rewriteValuePPC64_OpAvg64u_0(v)
	case OpBitLen32:
		return rewriteValuePPC64_OpBitLen32_0(v)
	case OpBitLen64:
		return rewriteValuePPC64_OpBitLen64_0(v)
	case OpCeil:
		return rewriteValuePPC64_OpCeil_0(v)
	case OpClosureCall:
		return rewriteValuePPC64_OpClosureCall_0(v)
	case OpCom16:
		return rewriteValuePPC64_OpCom16_0(v)
	case OpCom32:
		return rewriteValuePPC64_OpCom32_0(v)
	case OpCom64:
		return rewriteValuePPC64_OpCom64_0(v)
	case OpCom8:
		return rewriteValuePPC64_OpCom8_0(v)
	case OpCondSelect:
		return rewriteValuePPC64_OpCondSelect_0(v)
	case OpConst16:
		return rewriteValuePPC64_OpConst16_0(v)
	case OpConst32:
		return rewriteValuePPC64_OpConst32_0(v)
	case OpConst32F:
		return rewriteValuePPC64_OpConst32F_0(v)
	case OpConst64:
		return rewriteValuePPC64_OpConst64_0(v)
	case OpConst64F:
		return rewriteValuePPC64_OpConst64F_0(v)
	case OpConst8:
		return rewriteValuePPC64_OpConst8_0(v)
	case OpConstBool:
		return rewriteValuePPC64_OpConstBool_0(v)
	case OpConstNil:
		return rewriteValuePPC64_OpConstNil_0(v)
	case OpCopysign:
		return rewriteValuePPC64_OpCopysign_0(v)
	case OpCtz16:
		return rewriteValuePPC64_OpCtz16_0(v)
	case OpCtz32:
		return rewriteValuePPC64_OpCtz32_0(v)
	case OpCtz32NonZero:
		return rewriteValuePPC64_OpCtz32NonZero_0(v)
	case OpCtz64:
		return rewriteValuePPC64_OpCtz64_0(v)
	case OpCtz64NonZero:
		return rewriteValuePPC64_OpCtz64NonZero_0(v)
	case OpCtz8:
		return rewriteValuePPC64_OpCtz8_0(v)
	case OpCvt32Fto32:
		return rewriteValuePPC64_OpCvt32Fto32_0(v)
	case OpCvt32Fto64:
		return rewriteValuePPC64_OpCvt32Fto64_0(v)
	case OpCvt32Fto64F:
		return rewriteValuePPC64_OpCvt32Fto64F_0(v)
	case OpCvt32to32F:
		return rewriteValuePPC64_OpCvt32to32F_0(v)
	case OpCvt32to64F:
		return rewriteValuePPC64_OpCvt32to64F_0(v)
	case OpCvt64Fto32:
		return rewriteValuePPC64_OpCvt64Fto32_0(v)
	case OpCvt64Fto32F:
		return rewriteValuePPC64_OpCvt64Fto32F_0(v)
	case OpCvt64Fto64:
		return rewriteValuePPC64_OpCvt64Fto64_0(v)
	case OpCvt64to32F:
		return rewriteValuePPC64_OpCvt64to32F_0(v)
	case OpCvt64to64F:
		return rewriteValuePPC64_OpCvt64to64F_0(v)
	case OpDiv16:
		return rewriteValuePPC64_OpDiv16_0(v)
	case OpDiv16u:
		return rewriteValuePPC64_OpDiv16u_0(v)
	case OpDiv32:
		return rewriteValuePPC64_OpDiv32_0(v)
	case OpDiv32F:
		return rewriteValuePPC64_OpDiv32F_0(v)
	case OpDiv32u:
		return rewriteValuePPC64_OpDiv32u_0(v)
	case OpDiv64:
		return rewriteValuePPC64_OpDiv64_0(v)
	case OpDiv64F:
		return rewriteValuePPC64_OpDiv64F_0(v)
	case OpDiv64u:
		return rewriteValuePPC64_OpDiv64u_0(v)
	case OpDiv8:
		return rewriteValuePPC64_OpDiv8_0(v)
	case OpDiv8u:
		return rewriteValuePPC64_OpDiv8u_0(v)
	case OpEq16:
		return rewriteValuePPC64_OpEq16_0(v)
	case OpEq32:
		return rewriteValuePPC64_OpEq32_0(v)
	case OpEq32F:
		return rewriteValuePPC64_OpEq32F_0(v)
	case OpEq64:
		return rewriteValuePPC64_OpEq64_0(v)
	case OpEq64F:
		return rewriteValuePPC64_OpEq64F_0(v)
	case OpEq8:
		return rewriteValuePPC64_OpEq8_0(v)
	case OpEqB:
		return rewriteValuePPC64_OpEqB_0(v)
	case OpEqPtr:
		return rewriteValuePPC64_OpEqPtr_0(v)
	case OpFMA:
		return rewriteValuePPC64_OpFMA_0(v)
	case OpFloor:
		return rewriteValuePPC64_OpFloor_0(v)
	case OpGeq16:
		return rewriteValuePPC64_OpGeq16_0(v)
	case OpGeq16U:
		return rewriteValuePPC64_OpGeq16U_0(v)
	case OpGeq32:
		return rewriteValuePPC64_OpGeq32_0(v)
	case OpGeq32F:
		return rewriteValuePPC64_OpGeq32F_0(v)
	case OpGeq32U:
		return rewriteValuePPC64_OpGeq32U_0(v)
	case OpGeq64:
		return rewriteValuePPC64_OpGeq64_0(v)
	case OpGeq64F:
		return rewriteValuePPC64_OpGeq64F_0(v)
	case OpGeq64U:
		return rewriteValuePPC64_OpGeq64U_0(v)
	case OpGeq8:
		return rewriteValuePPC64_OpGeq8_0(v)
	case OpGeq8U:
		return rewriteValuePPC64_OpGeq8U_0(v)
	case OpGetCallerPC:
		return rewriteValuePPC64_OpGetCallerPC_0(v)
	case OpGetCallerSP:
		return rewriteValuePPC64_OpGetCallerSP_0(v)
	case OpGetClosurePtr:
		return rewriteValuePPC64_OpGetClosurePtr_0(v)
	case OpGreater16:
		return rewriteValuePPC64_OpGreater16_0(v)
	case OpGreater16U:
		return rewriteValuePPC64_OpGreater16U_0(v)
	case OpGreater32:
		return rewriteValuePPC64_OpGreater32_0(v)
	case OpGreater32F:
		return rewriteValuePPC64_OpGreater32F_0(v)
	case OpGreater32U:
		return rewriteValuePPC64_OpGreater32U_0(v)
	case OpGreater64:
		return rewriteValuePPC64_OpGreater64_0(v)
	case OpGreater64F:
		return rewriteValuePPC64_OpGreater64F_0(v)
	case OpGreater64U:
		return rewriteValuePPC64_OpGreater64U_0(v)
	case OpGreater8:
		return rewriteValuePPC64_OpGreater8_0(v)
	case OpGreater8U:
		return rewriteValuePPC64_OpGreater8U_0(v)
	case OpHmul32:
		return rewriteValuePPC64_OpHmul32_0(v)
	case OpHmul32u:
		return rewriteValuePPC64_OpHmul32u_0(v)
	case OpHmul64:
		return rewriteValuePPC64_OpHmul64_0(v)
	case OpHmul64u:
		return rewriteValuePPC64_OpHmul64u_0(v)
	case OpInterCall:
		return rewriteValuePPC64_OpInterCall_0(v)
	case OpIsInBounds:
		return rewriteValuePPC64_OpIsInBounds_0(v)
	case OpIsNonNil:
		return rewriteValuePPC64_OpIsNonNil_0(v)
	case OpIsSliceInBounds:
		return rewriteValuePPC64_OpIsSliceInBounds_0(v)
	case OpLeq16:
		return rewriteValuePPC64_OpLeq16_0(v)
	case OpLeq16U:
		return rewriteValuePPC64_OpLeq16U_0(v)
	case OpLeq32:
		return rewriteValuePPC64_OpLeq32_0(v)
	case OpLeq32F:
		return rewriteValuePPC64_OpLeq32F_0(v)
	case OpLeq32U:
		return rewriteValuePPC64_OpLeq32U_0(v)
	case OpLeq64:
		return rewriteValuePPC64_OpLeq64_0(v)
	case OpLeq64F:
		return rewriteValuePPC64_OpLeq64F_0(v)
	case OpLeq64U:
		return rewriteValuePPC64_OpLeq64U_0(v)
	case OpLeq8:
		return rewriteValuePPC64_OpLeq8_0(v)
	case OpLeq8U:
		return rewriteValuePPC64_OpLeq8U_0(v)
	case OpLess16:
		return rewriteValuePPC64_OpLess16_0(v)
	case OpLess16U:
		return rewriteValuePPC64_OpLess16U_0(v)
	case OpLess32:
		return rewriteValuePPC64_OpLess32_0(v)
	case OpLess32F:
		return rewriteValuePPC64_OpLess32F_0(v)
	case OpLess32U:
		return rewriteValuePPC64_OpLess32U_0(v)
	case OpLess64:
		return rewriteValuePPC64_OpLess64_0(v)
	case OpLess64F:
		return rewriteValuePPC64_OpLess64F_0(v)
	case OpLess64U:
		return rewriteValuePPC64_OpLess64U_0(v)
	case OpLess8:
		return rewriteValuePPC64_OpLess8_0(v)
	case OpLess8U:
		return rewriteValuePPC64_OpLess8U_0(v)
	case OpLoad:
		return rewriteValuePPC64_OpLoad_0(v)
	case OpLocalAddr:
		return rewriteValuePPC64_OpLocalAddr_0(v)
	case OpLsh16x16:
		return rewriteValuePPC64_OpLsh16x16_0(v)
	case OpLsh16x32:
		return rewriteValuePPC64_OpLsh16x32_0(v)
	case OpLsh16x64:
		return rewriteValuePPC64_OpLsh16x64_0(v)
	case OpLsh16x8:
		return rewriteValuePPC64_OpLsh16x8_0(v)
	case OpLsh32x16:
		return rewriteValuePPC64_OpLsh32x16_0(v)
	case OpLsh32x32:
		return rewriteValuePPC64_OpLsh32x32_0(v)
	case OpLsh32x64:
		return rewriteValuePPC64_OpLsh32x64_0(v)
	case OpLsh32x8:
		return rewriteValuePPC64_OpLsh32x8_0(v)
	case OpLsh64x16:
		return rewriteValuePPC64_OpLsh64x16_0(v)
	case OpLsh64x32:
		return rewriteValuePPC64_OpLsh64x32_0(v)
	case OpLsh64x64:
		return rewriteValuePPC64_OpLsh64x64_0(v)
	case OpLsh64x8:
		return rewriteValuePPC64_OpLsh64x8_0(v)
	case OpLsh8x16:
		return rewriteValuePPC64_OpLsh8x16_0(v)
	case OpLsh8x32:
		return rewriteValuePPC64_OpLsh8x32_0(v)
	case OpLsh8x64:
		return rewriteValuePPC64_OpLsh8x64_0(v)
	case OpLsh8x8:
		return rewriteValuePPC64_OpLsh8x8_0(v)
	case OpMod16:
		return rewriteValuePPC64_OpMod16_0(v)
	case OpMod16u:
		return rewriteValuePPC64_OpMod16u_0(v)
	case OpMod32:
		return rewriteValuePPC64_OpMod32_0(v)
	case OpMod32u:
		return rewriteValuePPC64_OpMod32u_0(v)
	case OpMod64:
		return rewriteValuePPC64_OpMod64_0(v)
	case OpMod64u:
		return rewriteValuePPC64_OpMod64u_0(v)
	case OpMod8:
		return rewriteValuePPC64_OpMod8_0(v)
	case OpMod8u:
		return rewriteValuePPC64_OpMod8u_0(v)
	case OpMove:
		return rewriteValuePPC64_OpMove_0(v) || rewriteValuePPC64_OpMove_10(v)
	case OpMul16:
		return rewriteValuePPC64_OpMul16_0(v)
	case OpMul32:
		return rewriteValuePPC64_OpMul32_0(v)
	case OpMul32F:
		return rewriteValuePPC64_OpMul32F_0(v)
	case OpMul64:
		return rewriteValuePPC64_OpMul64_0(v)
	case OpMul64F:
		return rewriteValuePPC64_OpMul64F_0(v)
	case OpMul64uhilo:
		return rewriteValuePPC64_OpMul64uhilo_0(v)
	case OpMul8:
		return rewriteValuePPC64_OpMul8_0(v)
	case OpNeg16:
		return rewriteValuePPC64_OpNeg16_0(v)
	case OpNeg32:
		return rewriteValuePPC64_OpNeg32_0(v)
	case OpNeg32F:
		return rewriteValuePPC64_OpNeg32F_0(v)
	case OpNeg64:
		return rewriteValuePPC64_OpNeg64_0(v)
	case OpNeg64F:
		return rewriteValuePPC64_OpNeg64F_0(v)
	case OpNeg8:
		return rewriteValuePPC64_OpNeg8_0(v)
	case OpNeq16:
		return rewriteValuePPC64_OpNeq16_0(v)
	case OpNeq32:
		return rewriteValuePPC64_OpNeq32_0(v)
	case OpNeq32F:
		return rewriteValuePPC64_OpNeq32F_0(v)
	case OpNeq64:
		return rewriteValuePPC64_OpNeq64_0(v)
	case OpNeq64F:
		return rewriteValuePPC64_OpNeq64F_0(v)
	case OpNeq8:
		return rewriteValuePPC64_OpNeq8_0(v)
	case OpNeqB:
		return rewriteValuePPC64_OpNeqB_0(v)
	case OpNeqPtr:
		return rewriteValuePPC64_OpNeqPtr_0(v)
	case OpNilCheck:
		return rewriteValuePPC64_OpNilCheck_0(v)
	case OpNot:
		return rewriteValuePPC64_OpNot_0(v)
	case OpOffPtr:
		return rewriteValuePPC64_OpOffPtr_0(v)
	case OpOr16:
		return rewriteValuePPC64_OpOr16_0(v)
	case OpOr32:
		return rewriteValuePPC64_OpOr32_0(v)
	case OpOr64:
		return rewriteValuePPC64_OpOr64_0(v)
	case OpOr8:
		return rewriteValuePPC64_OpOr8_0(v)
	case OpOrB:
		return rewriteValuePPC64_OpOrB_0(v)
	case OpPPC64ADD:
		return rewriteValuePPC64_OpPPC64ADD_0(v)
	case OpPPC64ADDconst:
		return rewriteValuePPC64_OpPPC64ADDconst_0(v)
	case OpPPC64AND:
		return rewriteValuePPC64_OpPPC64AND_0(v) || rewriteValuePPC64_OpPPC64AND_10(v)
	case OpPPC64ANDconst:
		return rewriteValuePPC64_OpPPC64ANDconst_0(v) || rewriteValuePPC64_OpPPC64ANDconst_10(v)
	case OpPPC64CMP:
		return rewriteValuePPC64_OpPPC64CMP_0(v)
	case OpPPC64CMPU:
		return rewriteValuePPC64_OpPPC64CMPU_0(v)
	case OpPPC64CMPUconst:
		return rewriteValuePPC64_OpPPC64CMPUconst_0(v)
	case OpPPC64CMPW:
		return rewriteValuePPC64_OpPPC64CMPW_0(v)
	case OpPPC64CMPWU:
		return rewriteValuePPC64_OpPPC64CMPWU_0(v)
	case OpPPC64CMPWUconst:
		return rewriteValuePPC64_OpPPC64CMPWUconst_0(v)
	case OpPPC64CMPWconst:
		return rewriteValuePPC64_OpPPC64CMPWconst_0(v)
	case OpPPC64CMPconst:
		return rewriteValuePPC64_OpPPC64CMPconst_0(v)
	case OpPPC64Equal:
		return rewriteValuePPC64_OpPPC64Equal_0(v)
	case OpPPC64FABS:
		return rewriteValuePPC64_OpPPC64FABS_0(v)
	case OpPPC64FADD:
		return rewriteValuePPC64_OpPPC64FADD_0(v)
	case OpPPC64FADDS:
		return rewriteValuePPC64_OpPPC64FADDS_0(v)
	case OpPPC64FCEIL:
		return rewriteValuePPC64_OpPPC64FCEIL_0(v)
	case OpPPC64FFLOOR:
		return rewriteValuePPC64_OpPPC64FFLOOR_0(v)
	case OpPPC64FGreaterEqual:
		return rewriteValuePPC64_OpPPC64FGreaterEqual_0(v)
	case OpPPC64FGreaterThan:
		return rewriteValuePPC64_OpPPC64FGreaterThan_0(v)
	case OpPPC64FLessEqual:
		return rewriteValuePPC64_OpPPC64FLessEqual_0(v)
	case OpPPC64FLessThan:
		return rewriteValuePPC64_OpPPC64FLessThan_0(v)
	case OpPPC64FMOVDload:
		return rewriteValuePPC64_OpPPC64FMOVDload_0(v)
	case OpPPC64FMOVDstore:
		return rewriteValuePPC64_OpPPC64FMOVDstore_0(v)
	case OpPPC64FMOVSload:
		return rewriteValuePPC64_OpPPC64FMOVSload_0(v)
	case OpPPC64FMOVSstore:
		return rewriteValuePPC64_OpPPC64FMOVSstore_0(v)
	case OpPPC64FNEG:
		return rewriteValuePPC64_OpPPC64FNEG_0(v)
	case OpPPC64FSQRT:
		return rewriteValuePPC64_OpPPC64FSQRT_0(v)
	case OpPPC64FSUB:
		return rewriteValuePPC64_OpPPC64FSUB_0(v)
	case OpPPC64FSUBS:
		return rewriteValuePPC64_OpPPC64FSUBS_0(v)
	case OpPPC64FTRUNC:
		return rewriteValuePPC64_OpPPC64FTRUNC_0(v)
	case OpPPC64GreaterEqual:
		return rewriteValuePPC64_OpPPC64GreaterEqual_0(v)
	case OpPPC64GreaterThan:
		return rewriteValuePPC64_OpPPC64GreaterThan_0(v)
	case OpPPC64ISEL:
		return rewriteValuePPC64_OpPPC64ISEL_0(v) || rewriteValuePPC64_OpPPC64ISEL_10(v) || rewriteValuePPC64_OpPPC64ISEL_20(v)
	case OpPPC64ISELB:
		return rewriteValuePPC64_OpPPC64ISELB_0(v) || rewriteValuePPC64_OpPPC64ISELB_10(v) || rewriteValuePPC64_OpPPC64ISELB_20(v)
	case OpPPC64LessEqual:
		return rewriteValuePPC64_OpPPC64LessEqual_0(v)
	case OpPPC64LessThan:
		return rewriteValuePPC64_OpPPC64LessThan_0(v)
	case OpPPC64MFVSRD:
		return rewriteValuePPC64_OpPPC64MFVSRD_0(v)
	case OpPPC64MOVBZload:
		return rewriteValuePPC64_OpPPC64MOVBZload_0(v)
	case OpPPC64MOVBZloadidx:
		return rewriteValuePPC64_OpPPC64MOVBZloadidx_0(v)
	case OpPPC64MOVBZreg:
		return rewriteValuePPC64_OpPPC64MOVBZreg_0(v) || rewriteValuePPC64_OpPPC64MOVBZreg_10(v)
	case OpPPC64MOVBreg:
		return rewriteValuePPC64_OpPPC64MOVBreg_0(v) || rewriteValuePPC64_OpPPC64MOVBreg_10(v)
	case OpPPC64MOVBstore:
		return rewriteValuePPC64_OpPPC64MOVBstore_0(v) || rewriteValuePPC64_OpPPC64MOVBstore_10(v) || rewriteValuePPC64_OpPPC64MOVBstore_20(v)
	case OpPPC64MOVBstoreidx:
		return rewriteValuePPC64_OpPPC64MOVBstoreidx_0(v) || rewriteValuePPC64_OpPPC64MOVBstoreidx_10(v)
	case OpPPC64MOVBstorezero:
		return rewriteValuePPC64_OpPPC64MOVBstorezero_0(v)
	case OpPPC64MOVDload:
		return rewriteValuePPC64_OpPPC64MOVDload_0(v)
	case OpPPC64MOVDloadidx:
		return rewriteValuePPC64_OpPPC64MOVDloadidx_0(v)
	case OpPPC64MOVDstore:
		return rewriteValuePPC64_OpPPC64MOVDstore_0(v)
	case OpPPC64MOVDstoreidx:
		return rewriteValuePPC64_OpPPC64MOVDstoreidx_0(v)
	case OpPPC64MOVDstorezero:
		return rewriteValuePPC64_OpPPC64MOVDstorezero_0(v)
	case OpPPC64MOVHBRstore:
		return rewriteValuePPC64_OpPPC64MOVHBRstore_0(v)
	case OpPPC64MOVHZload:
		return rewriteValuePPC64_OpPPC64MOVHZload_0(v)
	case OpPPC64MOVHZloadidx:
		return rewriteValuePPC64_OpPPC64MOVHZloadidx_0(v)
	case OpPPC64MOVHZreg:
		return rewriteValuePPC64_OpPPC64MOVHZreg_0(v) || rewriteValuePPC64_OpPPC64MOVHZreg_10(v)
	case OpPPC64MOVHload:
		return rewriteValuePPC64_OpPPC64MOVHload_0(v)
	case OpPPC64MOVHloadidx:
		return rewriteValuePPC64_OpPPC64MOVHloadidx_0(v)
	case OpPPC64MOVHreg:
		return rewriteValuePPC64_OpPPC64MOVHreg_0(v) || rewriteValuePPC64_OpPPC64MOVHreg_10(v)
	case OpPPC64MOVHstore:
		return rewriteValuePPC64_OpPPC64MOVHstore_0(v)
	case OpPPC64MOVHstoreidx:
		return rewriteValuePPC64_OpPPC64MOVHstoreidx_0(v)
	case OpPPC64MOVHstorezero:
		return rewriteValuePPC64_OpPPC64MOVHstorezero_0(v)
	case OpPPC64MOVWBRstore:
		return rewriteValuePPC64_OpPPC64MOVWBRstore_0(v)
	case OpPPC64MOVWZload:
		return rewriteValuePPC64_OpPPC64MOVWZload_0(v)
	case OpPPC64MOVWZloadidx:
		return rewriteValuePPC64_OpPPC64MOVWZloadidx_0(v)
	case OpPPC64MOVWZreg:
		return rewriteValuePPC64_OpPPC64MOVWZreg_0(v) || rewriteValuePPC64_OpPPC64MOVWZreg_10(v) || rewriteValuePPC64_OpPPC64MOVWZreg_20(v)
	case OpPPC64MOVWload:
		return rewriteValuePPC64_OpPPC64MOVWload_0(v)
	case OpPPC64MOVWloadidx:
		return rewriteValuePPC64_OpPPC64MOVWloadidx_0(v)
	case OpPPC64MOVWreg:
		return rewriteValuePPC64_OpPPC64MOVWreg_0(v) || rewriteValuePPC64_OpPPC64MOVWreg_10(v)
	case OpPPC64MOVWstore:
		return rewriteValuePPC64_OpPPC64MOVWstore_0(v)
	case OpPPC64MOVWstoreidx:
		return rewriteValuePPC64_OpPPC64MOVWstoreidx_0(v)
	case OpPPC64MOVWstorezero:
		return rewriteValuePPC64_OpPPC64MOVWstorezero_0(v)
	case OpPPC64MTVSRD:
		return rewriteValuePPC64_OpPPC64MTVSRD_0(v)
	case OpPPC64MaskIfNotCarry:
		return rewriteValuePPC64_OpPPC64MaskIfNotCarry_0(v)
	case OpPPC64NotEqual:
		return rewriteValuePPC64_OpPPC64NotEqual_0(v)
	case OpPPC64OR:
		return rewriteValuePPC64_OpPPC64OR_0(v) || rewriteValuePPC64_OpPPC64OR_10(v) || rewriteValuePPC64_OpPPC64OR_20(v) || rewriteValuePPC64_OpPPC64OR_30(v) || rewriteValuePPC64_OpPPC64OR_40(v) || rewriteValuePPC64_OpPPC64OR_50(v) || rewriteValuePPC64_OpPPC64OR_60(v) || rewriteValuePPC64_OpPPC64OR_70(v) || rewriteValuePPC64_OpPPC64OR_80(v) || rewriteValuePPC64_OpPPC64OR_90(v) || rewriteValuePPC64_OpPPC64OR_100(v) || rewriteValuePPC64_OpPPC64OR_110(v)
	case OpPPC64ORN:
		return rewriteValuePPC64_OpPPC64ORN_0(v)
	case OpPPC64ORconst:
		return rewriteValuePPC64_OpPPC64ORconst_0(v)
	case OpPPC64ROTL:
		return rewriteValuePPC64_OpPPC64ROTL_0(v)
	case OpPPC64ROTLW:
		return rewriteValuePPC64_OpPPC64ROTLW_0(v)
	case OpPPC64SUB:
		return rewriteValuePPC64_OpPPC64SUB_0(v)
	case OpPPC64XOR:
		return rewriteValuePPC64_OpPPC64XOR_0(v) || rewriteValuePPC64_OpPPC64XOR_10(v)
	case OpPPC64XORconst:
		return rewriteValuePPC64_OpPPC64XORconst_0(v)
	case OpPanicBounds:
		return rewriteValuePPC64_OpPanicBounds_0(v)
	case OpPopCount16:
		return rewriteValuePPC64_OpPopCount16_0(v)
	case OpPopCount32:
		return rewriteValuePPC64_OpPopCount32_0(v)
	case OpPopCount64:
		return rewriteValuePPC64_OpPopCount64_0(v)
	case OpPopCount8:
		return rewriteValuePPC64_OpPopCount8_0(v)
	case OpRotateLeft16:
		return rewriteValuePPC64_OpRotateLeft16_0(v)
	case OpRotateLeft32:
		return rewriteValuePPC64_OpRotateLeft32_0(v)
	case OpRotateLeft64:
		return rewriteValuePPC64_OpRotateLeft64_0(v)
	case OpRotateLeft8:
		return rewriteValuePPC64_OpRotateLeft8_0(v)
	case OpRound:
		return rewriteValuePPC64_OpRound_0(v)
	case OpRound32F:
		return rewriteValuePPC64_OpRound32F_0(v)
	case OpRound64F:
		return rewriteValuePPC64_OpRound64F_0(v)
	case OpRsh16Ux16:
		return rewriteValuePPC64_OpRsh16Ux16_0(v)
	case OpRsh16Ux32:
		return rewriteValuePPC64_OpRsh16Ux32_0(v)
	case OpRsh16Ux64:
		return rewriteValuePPC64_OpRsh16Ux64_0(v)
	case OpRsh16Ux8:
		return rewriteValuePPC64_OpRsh16Ux8_0(v)
	case OpRsh16x16:
		return rewriteValuePPC64_OpRsh16x16_0(v)
	case OpRsh16x32:
		return rewriteValuePPC64_OpRsh16x32_0(v)
	case OpRsh16x64:
		return rewriteValuePPC64_OpRsh16x64_0(v)
	case OpRsh16x8:
		return rewriteValuePPC64_OpRsh16x8_0(v)
	case OpRsh32Ux16:
		return rewriteValuePPC64_OpRsh32Ux16_0(v)
	case OpRsh32Ux32:
		return rewriteValuePPC64_OpRsh32Ux32_0(v)
	case OpRsh32Ux64:
		return rewriteValuePPC64_OpRsh32Ux64_0(v) || rewriteValuePPC64_OpRsh32Ux64_10(v)
	case OpRsh32Ux8:
		return rewriteValuePPC64_OpRsh32Ux8_0(v)
	case OpRsh32x16:
		return rewriteValuePPC64_OpRsh32x16_0(v)
	case OpRsh32x32:
		return rewriteValuePPC64_OpRsh32x32_0(v)
	case OpRsh32x64:
		return rewriteValuePPC64_OpRsh32x64_0(v) || rewriteValuePPC64_OpRsh32x64_10(v)
	case OpRsh32x8:
		return rewriteValuePPC64_OpRsh32x8_0(v)
	case OpRsh64Ux16:
		return rewriteValuePPC64_OpRsh64Ux16_0(v)
	case OpRsh64Ux32:
		return rewriteValuePPC64_OpRsh64Ux32_0(v)
	case OpRsh64Ux64:
		return rewriteValuePPC64_OpRsh64Ux64_0(v) || rewriteValuePPC64_OpRsh64Ux64_10(v)
	case OpRsh64Ux8:
		return rewriteValuePPC64_OpRsh64Ux8_0(v)
	case OpRsh64x16:
		return rewriteValuePPC64_OpRsh64x16_0(v)
	case OpRsh64x32:
		return rewriteValuePPC64_OpRsh64x32_0(v)
	case OpRsh64x64:
		return rewriteValuePPC64_OpRsh64x64_0(v) || rewriteValuePPC64_OpRsh64x64_10(v)
	case OpRsh64x8:
		return rewriteValuePPC64_OpRsh64x8_0(v)
	case OpRsh8Ux16:
		return rewriteValuePPC64_OpRsh8Ux16_0(v)
	case OpRsh8Ux32:
		return rewriteValuePPC64_OpRsh8Ux32_0(v)
	case OpRsh8Ux64:
		return rewriteValuePPC64_OpRsh8Ux64_0(v)
	case OpRsh8Ux8:
		return rewriteValuePPC64_OpRsh8Ux8_0(v)
	case OpRsh8x16:
		return rewriteValuePPC64_OpRsh8x16_0(v)
	case OpRsh8x32:
		return rewriteValuePPC64_OpRsh8x32_0(v)
	case OpRsh8x64:
		return rewriteValuePPC64_OpRsh8x64_0(v)
	case OpRsh8x8:
		return rewriteValuePPC64_OpRsh8x8_0(v)
	case OpSignExt16to32:
		return rewriteValuePPC64_OpSignExt16to32_0(v)
	case OpSignExt16to64:
		return rewriteValuePPC64_OpSignExt16to64_0(v)
	case OpSignExt32to64:
		return rewriteValuePPC64_OpSignExt32to64_0(v)
	case OpSignExt8to16:
		return rewriteValuePPC64_OpSignExt8to16_0(v)
	case OpSignExt8to32:
		return rewriteValuePPC64_OpSignExt8to32_0(v)
	case OpSignExt8to64:
		return rewriteValuePPC64_OpSignExt8to64_0(v)
	case OpSlicemask:
		return rewriteValuePPC64_OpSlicemask_0(v)
	case OpSqrt:
		return rewriteValuePPC64_OpSqrt_0(v)
	case OpStaticCall:
		return rewriteValuePPC64_OpStaticCall_0(v)
	case OpStore:
		return rewriteValuePPC64_OpStore_0(v)
	case OpSub16:
		return rewriteValuePPC64_OpSub16_0(v)
	case OpSub32:
		return rewriteValuePPC64_OpSub32_0(v)
	case OpSub32F:
		return rewriteValuePPC64_OpSub32F_0(v)
	case OpSub64:
		return rewriteValuePPC64_OpSub64_0(v)
	case OpSub64F:
		return rewriteValuePPC64_OpSub64F_0(v)
	case OpSub8:
		return rewriteValuePPC64_OpSub8_0(v)
	case OpSubPtr:
		return rewriteValuePPC64_OpSubPtr_0(v)
	case OpTrunc:
		return rewriteValuePPC64_OpTrunc_0(v)
	case OpTrunc16to8:
		return rewriteValuePPC64_OpTrunc16to8_0(v)
	case OpTrunc32to16:
		return rewriteValuePPC64_OpTrunc32to16_0(v)
	case OpTrunc32to8:
		return rewriteValuePPC64_OpTrunc32to8_0(v)
	case OpTrunc64to16:
		return rewriteValuePPC64_OpTrunc64to16_0(v)
	case OpTrunc64to32:
		return rewriteValuePPC64_OpTrunc64to32_0(v)
	case OpTrunc64to8:
		return rewriteValuePPC64_OpTrunc64to8_0(v)
	case OpWB:
		return rewriteValuePPC64_OpWB_0(v)
	case OpXor16:
		return rewriteValuePPC64_OpXor16_0(v)
	case OpXor32:
		return rewriteValuePPC64_OpXor32_0(v)
	case OpXor64:
		return rewriteValuePPC64_OpXor64_0(v)
	case OpXor8:
		return rewriteValuePPC64_OpXor8_0(v)
	case OpZero:
		return rewriteValuePPC64_OpZero_0(v) || rewriteValuePPC64_OpZero_10(v)
	case OpZeroExt16to32:
		return rewriteValuePPC64_OpZeroExt16to32_0(v)
	case OpZeroExt16to64:
		return rewriteValuePPC64_OpZeroExt16to64_0(v)
	case OpZeroExt32to64:
		return rewriteValuePPC64_OpZeroExt32to64_0(v)
	case OpZeroExt8to16:
		return rewriteValuePPC64_OpZeroExt8to16_0(v)
	case OpZeroExt8to32:
		return rewriteValuePPC64_OpZeroExt8to32_0(v)
	case OpZeroExt8to64:
		return rewriteValuePPC64_OpZeroExt8to64_0(v)
	}
	return false
}
func rewriteValuePPC64_OpAbs_0(v *Value) bool {
	// match: (Abs x)
	// result: (FABS x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FABS)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpAdd16_0(v *Value) bool {
	// match: (Add16 x y)
	// result: (ADD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAdd32_0(v *Value) bool {
	// match: (Add32 x y)
	// result: (ADD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAdd32F_0(v *Value) bool {
	// match: (Add32F x y)
	// result: (FADDS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FADDS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAdd64_0(v *Value) bool {
	// match: (Add64 x y)
	// result: (ADD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAdd64F_0(v *Value) bool {
	// match: (Add64F x y)
	// result: (FADD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAdd64carry_0(v *Value) bool {
	// match: (Add64carry x y c)
	// result: (LoweredAdd64Carry x y c)
	for {
		c := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpPPC64LoweredAdd64Carry)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(c)
		return true
	}
}
func rewriteValuePPC64_OpAdd8_0(v *Value) bool {
	// match: (Add8 x y)
	// result: (ADD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAddPtr_0(v *Value) bool {
	// match: (AddPtr x y)
	// result: (ADD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAddr_0(v *Value) bool {
	// match: (Addr {sym} base)
	// result: (MOVDaddr {sym} base)
	for {
		sym := v.Aux
		base := v.Args[0]
		v.reset(OpPPC64MOVDaddr)
		v.Aux = sym
		v.AddArg(base)
		return true
	}
}
func rewriteValuePPC64_OpAnd16_0(v *Value) bool {
	// match: (And16 x y)
	// result: (AND x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAnd32_0(v *Value) bool {
	// match: (And32 x y)
	// result: (AND x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAnd64_0(v *Value) bool {
	// match: (And64 x y)
	// result: (AND x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAnd8_0(v *Value) bool {
	// match: (And8 x y)
	// result: (AND x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAndB_0(v *Value) bool {
	// match: (AndB x y)
	// result: (AND x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpAtomicAdd32_0(v *Value) bool {
	// match: (AtomicAdd32 ptr val mem)
	// result: (LoweredAtomicAdd32 ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicAdd32)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicAdd64_0(v *Value) bool {
	// match: (AtomicAdd64 ptr val mem)
	// result: (LoweredAtomicAdd64 ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicAdd64)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicAnd8_0(v *Value) bool {
	// match: (AtomicAnd8 ptr val mem)
	// result: (LoweredAtomicAnd8 ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicAnd8)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicCompareAndSwap32_0(v *Value) bool {
	// match: (AtomicCompareAndSwap32 ptr old new_ mem)
	// result: (LoweredAtomicCas32 [1] ptr old new_ mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		old := v.Args[1]
		new_ := v.Args[2]
		v.reset(OpPPC64LoweredAtomicCas32)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(old)
		v.AddArg(new_)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicCompareAndSwap64_0(v *Value) bool {
	// match: (AtomicCompareAndSwap64 ptr old new_ mem)
	// result: (LoweredAtomicCas64 [1] ptr old new_ mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		old := v.Args[1]
		new_ := v.Args[2]
		v.reset(OpPPC64LoweredAtomicCas64)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(old)
		v.AddArg(new_)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicCompareAndSwapRel32_0(v *Value) bool {
	// match: (AtomicCompareAndSwapRel32 ptr old new_ mem)
	// result: (LoweredAtomicCas32 [0] ptr old new_ mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		old := v.Args[1]
		new_ := v.Args[2]
		v.reset(OpPPC64LoweredAtomicCas32)
		v.AuxInt = 0
		v.AddArg(ptr)
		v.AddArg(old)
		v.AddArg(new_)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicExchange32_0(v *Value) bool {
	// match: (AtomicExchange32 ptr val mem)
	// result: (LoweredAtomicExchange32 ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicExchange32)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicExchange64_0(v *Value) bool {
	// match: (AtomicExchange64 ptr val mem)
	// result: (LoweredAtomicExchange64 ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicExchange64)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicLoad32_0(v *Value) bool {
	// match: (AtomicLoad32 ptr mem)
	// result: (LoweredAtomicLoad32 [1] ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredAtomicLoad32)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicLoad64_0(v *Value) bool {
	// match: (AtomicLoad64 ptr mem)
	// result: (LoweredAtomicLoad64 [1] ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredAtomicLoad64)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicLoad8_0(v *Value) bool {
	// match: (AtomicLoad8 ptr mem)
	// result: (LoweredAtomicLoad8 [1] ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredAtomicLoad8)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicLoadAcq32_0(v *Value) bool {
	// match: (AtomicLoadAcq32 ptr mem)
	// result: (LoweredAtomicLoad32 [0] ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredAtomicLoad32)
		v.AuxInt = 0
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicLoadPtr_0(v *Value) bool {
	// match: (AtomicLoadPtr ptr mem)
	// result: (LoweredAtomicLoadPtr [1] ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredAtomicLoadPtr)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicOr8_0(v *Value) bool {
	// match: (AtomicOr8 ptr val mem)
	// result: (LoweredAtomicOr8 ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicOr8)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicStore32_0(v *Value) bool {
	// match: (AtomicStore32 ptr val mem)
	// result: (LoweredAtomicStore32 [1] ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicStore32)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicStore64_0(v *Value) bool {
	// match: (AtomicStore64 ptr val mem)
	// result: (LoweredAtomicStore64 [1] ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicStore64)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicStore8_0(v *Value) bool {
	// match: (AtomicStore8 ptr val mem)
	// result: (LoweredAtomicStore8 [1] ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicStore8)
		v.AuxInt = 1
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAtomicStoreRel32_0(v *Value) bool {
	// match: (AtomicStoreRel32 ptr val mem)
	// result: (LoweredAtomicStore32 [0] ptr val mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		v.reset(OpPPC64LoweredAtomicStore32)
		v.AuxInt = 0
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpAvg64u_0(v *Value) bool {
	b := v.Block
	// match: (Avg64u <t> x y)
	// result: (ADD (SRDconst <t> (SUB <t> x y) [1]) y)
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ADD)
		v0 := b.NewValue0(v.Pos, OpPPC64SRDconst, t)
		v0.AuxInt = 1
		v1 := b.NewValue0(v.Pos, OpPPC64SUB, t)
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpBitLen32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (BitLen32 x)
	// result: (SUB (MOVDconst [32]) (CNTLZW <typ.Int> x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 32
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64CNTLZW, typ.Int)
		v1.AddArg(x)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpBitLen64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (BitLen64 x)
	// result: (SUB (MOVDconst [64]) (CNTLZD <typ.Int> x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 64
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64CNTLZD, typ.Int)
		v1.AddArg(x)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpCeil_0(v *Value) bool {
	// match: (Ceil x)
	// result: (FCEIL x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FCEIL)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpClosureCall_0(v *Value) bool {
	// match: (ClosureCall [argwid] entry closure mem)
	// result: (CALLclosure [argwid] entry closure mem)
	for {
		argwid := v.AuxInt
		mem := v.Args[2]
		entry := v.Args[0]
		closure := v.Args[1]
		v.reset(OpPPC64CALLclosure)
		v.AuxInt = argwid
		v.AddArg(entry)
		v.AddArg(closure)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpCom16_0(v *Value) bool {
	// match: (Com16 x)
	// result: (NOR x x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NOR)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCom32_0(v *Value) bool {
	// match: (Com32 x)
	// result: (NOR x x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NOR)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCom64_0(v *Value) bool {
	// match: (Com64 x)
	// result: (NOR x x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NOR)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCom8_0(v *Value) bool {
	// match: (Com8 x)
	// result: (NOR x x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NOR)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCondSelect_0(v *Value) bool {
	b := v.Block
	// match: (CondSelect x y bool)
	// cond: flagArg(bool) != nil
	// result: (ISEL [2] x y bool)
	for {
		bool := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		if !(flagArg(bool) != nil) {
			break
		}
		v.reset(OpPPC64ISEL)
		v.AuxInt = 2
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(bool)
		return true
	}
	// match: (CondSelect x y bool)
	// cond: flagArg(bool) == nil
	// result: (ISEL [2] x y (CMPWconst [0] bool))
	for {
		bool := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		if !(flagArg(bool) == nil) {
			break
		}
		v.reset(OpPPC64ISEL)
		v.AuxInt = 2
		v.AddArg(x)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWconst, types.TypeFlags)
		v0.AuxInt = 0
		v0.AddArg(bool)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpConst16_0(v *Value) bool {
	// match: (Const16 [val])
	// result: (MOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValuePPC64_OpConst32_0(v *Value) bool {
	// match: (Const32 [val])
	// result: (MOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValuePPC64_OpConst32F_0(v *Value) bool {
	// match: (Const32F [val])
	// result: (FMOVSconst [val])
	for {
		val := v.AuxInt
		v.reset(OpPPC64FMOVSconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValuePPC64_OpConst64_0(v *Value) bool {
	// match: (Const64 [val])
	// result: (MOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValuePPC64_OpConst64F_0(v *Value) bool {
	// match: (Const64F [val])
	// result: (FMOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValuePPC64_OpConst8_0(v *Value) bool {
	// match: (Const8 [val])
	// result: (MOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValuePPC64_OpConstBool_0(v *Value) bool {
	// match: (ConstBool [b])
	// result: (MOVDconst [b])
	for {
		b := v.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = b
		return true
	}
}
func rewriteValuePPC64_OpConstNil_0(v *Value) bool {
	// match: (ConstNil)
	// result: (MOVDconst [0])
	for {
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
}
func rewriteValuePPC64_OpCopysign_0(v *Value) bool {
	// match: (Copysign x y)
	// result: (FCPSGN y x)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FCPSGN)
		v.AddArg(y)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCtz16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz16 x)
	// result: (POPCNTW (MOVHZreg (ANDN <typ.Int16> (ADDconst <typ.Int16> [-1] x) x)))
	for {
		x := v.Args[0]
		v.reset(OpPPC64POPCNTW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v1 := b.NewValue0(v.Pos, OpPPC64ANDN, typ.Int16)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconst, typ.Int16)
		v2.AuxInt = -1
		v2.AddArg(x)
		v1.AddArg(v2)
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCtz32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz32 x)
	// cond: objabi.GOPPC64<=8
	// result: (POPCNTW (MOVWZreg (ANDN <typ.Int> (ADDconst <typ.Int> [-1] x) x)))
	for {
		x := v.Args[0]
		if !(objabi.GOPPC64 <= 8) {
			break
		}
		v.reset(OpPPC64POPCNTW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWZreg, typ.Int64)
		v1 := b.NewValue0(v.Pos, OpPPC64ANDN, typ.Int)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconst, typ.Int)
		v2.AuxInt = -1
		v2.AddArg(x)
		v1.AddArg(v2)
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	// match: (Ctz32 x)
	// result: (CNTTZW (MOVWZreg x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64CNTTZW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCtz32NonZero_0(v *Value) bool {
	// match: (Ctz32NonZero x)
	// result: (Ctz32 x)
	for {
		x := v.Args[0]
		v.reset(OpCtz32)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCtz64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz64 x)
	// cond: objabi.GOPPC64<=8
	// result: (POPCNTD (ANDN <typ.Int64> (ADDconst <typ.Int64> [-1] x) x))
	for {
		x := v.Args[0]
		if !(objabi.GOPPC64 <= 8) {
			break
		}
		v.reset(OpPPC64POPCNTD)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDN, typ.Int64)
		v1 := b.NewValue0(v.Pos, OpPPC64ADDconst, typ.Int64)
		v1.AuxInt = -1
		v1.AddArg(x)
		v0.AddArg(v1)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Ctz64 x)
	// result: (CNTTZD x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64CNTTZD)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCtz64NonZero_0(v *Value) bool {
	// match: (Ctz64NonZero x)
	// result: (Ctz64 x)
	for {
		x := v.Args[0]
		v.reset(OpCtz64)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCtz8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz8 x)
	// result: (POPCNTB (MOVBZreg (ANDN <typ.UInt8> (ADDconst <typ.UInt8> [-1] x) x)))
	for {
		x := v.Args[0]
		v.reset(OpPPC64POPCNTB)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v1 := b.NewValue0(v.Pos, OpPPC64ANDN, typ.UInt8)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconst, typ.UInt8)
		v2.AuxInt = -1
		v2.AddArg(x)
		v1.AddArg(v2)
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt32Fto32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32Fto32 x)
	// result: (MFVSRD (FCTIWZ x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64MFVSRD)
		v0 := b.NewValue0(v.Pos, OpPPC64FCTIWZ, typ.Float64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt32Fto64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32Fto64 x)
	// result: (MFVSRD (FCTIDZ x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64MFVSRD)
		v0 := b.NewValue0(v.Pos, OpPPC64FCTIDZ, typ.Float64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt32Fto64F_0(v *Value) bool {
	// match: (Cvt32Fto64F x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCvt32to32F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32to32F x)
	// result: (FCFIDS (MTVSRD (SignExt32to64 x)))
	for {
		x := v.Args[0]
		v.reset(OpPPC64FCFIDS)
		v0 := b.NewValue0(v.Pos, OpPPC64MTVSRD, typ.Float64)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt32to64F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32to64F x)
	// result: (FCFID (MTVSRD (SignExt32to64 x)))
	for {
		x := v.Args[0]
		v.reset(OpPPC64FCFID)
		v0 := b.NewValue0(v.Pos, OpPPC64MTVSRD, typ.Float64)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt64Fto32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt64Fto32 x)
	// result: (MFVSRD (FCTIWZ x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64MFVSRD)
		v0 := b.NewValue0(v.Pos, OpPPC64FCTIWZ, typ.Float64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt64Fto32F_0(v *Value) bool {
	// match: (Cvt64Fto32F x)
	// result: (FRSP x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FRSP)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpCvt64Fto64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt64Fto64 x)
	// result: (MFVSRD (FCTIDZ x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64MFVSRD)
		v0 := b.NewValue0(v.Pos, OpPPC64FCTIDZ, typ.Float64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt64to32F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt64to32F x)
	// result: (FCFIDS (MTVSRD x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64FCFIDS)
		v0 := b.NewValue0(v.Pos, OpPPC64MTVSRD, typ.Float64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpCvt64to64F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt64to64F x)
	// result: (FCFID (MTVSRD x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64FCFID)
		v0 := b.NewValue0(v.Pos, OpPPC64MTVSRD, typ.Float64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpDiv16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div16 x y)
	// result: (DIVW (SignExt16to32 x) (SignExt16to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVW)
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpDiv16u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div16u x y)
	// result: (DIVWU (ZeroExt16to32 x) (ZeroExt16to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVWU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpDiv32_0(v *Value) bool {
	// match: (Div32 x y)
	// result: (DIVW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpDiv32F_0(v *Value) bool {
	// match: (Div32F x y)
	// result: (FDIVS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FDIVS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpDiv32u_0(v *Value) bool {
	// match: (Div32u x y)
	// result: (DIVWU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVWU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpDiv64_0(v *Value) bool {
	// match: (Div64 x y)
	// result: (DIVD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpDiv64F_0(v *Value) bool {
	// match: (Div64F x y)
	// result: (FDIV x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FDIV)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpDiv64u_0(v *Value) bool {
	// match: (Div64u x y)
	// result: (DIVDU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVDU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpDiv8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div8 x y)
	// result: (DIVW (SignExt8to32 x) (SignExt8to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpDiv8u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div8u x y)
	// result: (DIVWU (ZeroExt8to32 x) (ZeroExt8to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64DIVWU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpEq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Eq16 x y)
	// cond: isSigned(x.Type) && isSigned(y.Type)
	// result: (Equal (CMPW (SignExt16to32 x) (SignExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(isSigned(x.Type) && isSigned(y.Type)) {
			break
		}
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Eq16 x y)
	// result: (Equal (CMPW (ZeroExt16to32 x) (ZeroExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEq32_0(v *Value) bool {
	b := v.Block
	// match: (Eq32 x y)
	// result: (Equal (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEq32F_0(v *Value) bool {
	b := v.Block
	// match: (Eq32F x y)
	// result: (Equal (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEq64_0(v *Value) bool {
	b := v.Block
	// match: (Eq64 x y)
	// result: (Equal (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEq64F_0(v *Value) bool {
	b := v.Block
	// match: (Eq64F x y)
	// result: (Equal (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Eq8 x y)
	// cond: isSigned(x.Type) && isSigned(y.Type)
	// result: (Equal (CMPW (SignExt8to32 x) (SignExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(isSigned(x.Type) && isSigned(y.Type)) {
			break
		}
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Eq8 x y)
	// result: (Equal (CMPW (ZeroExt8to32 x) (ZeroExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEqB_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (EqB x y)
	// result: (ANDconst [1] (EQV x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = 1
		v0 := b.NewValue0(v.Pos, OpPPC64EQV, typ.Int64)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpEqPtr_0(v *Value) bool {
	b := v.Block
	// match: (EqPtr x y)
	// result: (Equal (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64Equal)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpFMA_0(v *Value) bool {
	// match: (FMA x y z)
	// result: (FMADD x y z)
	for {
		z := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpPPC64FMADD)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
}
func rewriteValuePPC64_OpFloor_0(v *Value) bool {
	// match: (Floor x)
	// result: (FFLOOR x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FFLOOR)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpGeq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq16 x y)
	// result: (GreaterEqual (CMPW (SignExt16to32 x) (SignExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq16U x y)
	// result: (GreaterEqual (CMPWU (ZeroExt16to32 x) (ZeroExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq32_0(v *Value) bool {
	b := v.Block
	// match: (Geq32 x y)
	// result: (GreaterEqual (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq32F_0(v *Value) bool {
	b := v.Block
	// match: (Geq32F x y)
	// result: (FGreaterEqual (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FGreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq32U_0(v *Value) bool {
	b := v.Block
	// match: (Geq32U x y)
	// result: (GreaterEqual (CMPWU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq64_0(v *Value) bool {
	b := v.Block
	// match: (Geq64 x y)
	// result: (GreaterEqual (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq64F_0(v *Value) bool {
	b := v.Block
	// match: (Geq64F x y)
	// result: (FGreaterEqual (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FGreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq64U_0(v *Value) bool {
	b := v.Block
	// match: (Geq64U x y)
	// result: (GreaterEqual (CMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq8 x y)
	// result: (GreaterEqual (CMPW (SignExt8to32 x) (SignExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGeq8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq8U x y)
	// result: (GreaterEqual (CMPWU (ZeroExt8to32 x) (ZeroExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGetCallerPC_0(v *Value) bool {
	// match: (GetCallerPC)
	// result: (LoweredGetCallerPC)
	for {
		v.reset(OpPPC64LoweredGetCallerPC)
		return true
	}
}
func rewriteValuePPC64_OpGetCallerSP_0(v *Value) bool {
	// match: (GetCallerSP)
	// result: (LoweredGetCallerSP)
	for {
		v.reset(OpPPC64LoweredGetCallerSP)
		return true
	}
}
func rewriteValuePPC64_OpGetClosurePtr_0(v *Value) bool {
	// match: (GetClosurePtr)
	// result: (LoweredGetClosurePtr)
	for {
		v.reset(OpPPC64LoweredGetClosurePtr)
		return true
	}
}
func rewriteValuePPC64_OpGreater16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater16 x y)
	// result: (GreaterThan (CMPW (SignExt16to32 x) (SignExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater16U x y)
	// result: (GreaterThan (CMPWU (ZeroExt16to32 x) (ZeroExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater32_0(v *Value) bool {
	b := v.Block
	// match: (Greater32 x y)
	// result: (GreaterThan (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater32F_0(v *Value) bool {
	b := v.Block
	// match: (Greater32F x y)
	// result: (FGreaterThan (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FGreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater32U_0(v *Value) bool {
	b := v.Block
	// match: (Greater32U x y)
	// result: (GreaterThan (CMPWU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater64_0(v *Value) bool {
	b := v.Block
	// match: (Greater64 x y)
	// result: (GreaterThan (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater64F_0(v *Value) bool {
	b := v.Block
	// match: (Greater64F x y)
	// result: (FGreaterThan (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FGreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater64U_0(v *Value) bool {
	b := v.Block
	// match: (Greater64U x y)
	// result: (GreaterThan (CMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater8 x y)
	// result: (GreaterThan (CMPW (SignExt8to32 x) (SignExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpGreater8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater8U x y)
	// result: (GreaterThan (CMPWU (ZeroExt8to32 x) (ZeroExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64GreaterThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpHmul32_0(v *Value) bool {
	// match: (Hmul32 x y)
	// result: (MULHW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULHW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpHmul32u_0(v *Value) bool {
	// match: (Hmul32u x y)
	// result: (MULHWU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULHWU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpHmul64_0(v *Value) bool {
	// match: (Hmul64 x y)
	// result: (MULHD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULHD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpHmul64u_0(v *Value) bool {
	// match: (Hmul64u x y)
	// result: (MULHDU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULHDU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpInterCall_0(v *Value) bool {
	// match: (InterCall [argwid] entry mem)
	// result: (CALLinter [argwid] entry mem)
	for {
		argwid := v.AuxInt
		mem := v.Args[1]
		entry := v.Args[0]
		v.reset(OpPPC64CALLinter)
		v.AuxInt = argwid
		v.AddArg(entry)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpIsInBounds_0(v *Value) bool {
	b := v.Block
	// match: (IsInBounds idx len)
	// result: (LessThan (CMPU idx len))
	for {
		len := v.Args[1]
		idx := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPU, types.TypeFlags)
		v0.AddArg(idx)
		v0.AddArg(len)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpIsNonNil_0(v *Value) bool {
	b := v.Block
	// match: (IsNonNil ptr)
	// result: (NotEqual (CMPconst [0] ptr))
	for {
		ptr := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPconst, types.TypeFlags)
		v0.AuxInt = 0
		v0.AddArg(ptr)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpIsSliceInBounds_0(v *Value) bool {
	b := v.Block
	// match: (IsSliceInBounds idx len)
	// result: (LessEqual (CMPU idx len))
	for {
		len := v.Args[1]
		idx := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPU, types.TypeFlags)
		v0.AddArg(idx)
		v0.AddArg(len)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq16 x y)
	// result: (LessEqual (CMPW (SignExt16to32 x) (SignExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq16U x y)
	// result: (LessEqual (CMPWU (ZeroExt16to32 x) (ZeroExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq32_0(v *Value) bool {
	b := v.Block
	// match: (Leq32 x y)
	// result: (LessEqual (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq32F_0(v *Value) bool {
	b := v.Block
	// match: (Leq32F x y)
	// result: (FLessEqual (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FLessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq32U_0(v *Value) bool {
	b := v.Block
	// match: (Leq32U x y)
	// result: (LessEqual (CMPWU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq64_0(v *Value) bool {
	b := v.Block
	// match: (Leq64 x y)
	// result: (LessEqual (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq64F_0(v *Value) bool {
	b := v.Block
	// match: (Leq64F x y)
	// result: (FLessEqual (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FLessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq64U_0(v *Value) bool {
	b := v.Block
	// match: (Leq64U x y)
	// result: (LessEqual (CMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq8 x y)
	// result: (LessEqual (CMPW (SignExt8to32 x) (SignExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLeq8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq8U x y)
	// result: (LessEqual (CMPWU (ZeroExt8to32 x) (ZeroExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less16 x y)
	// result: (LessThan (CMPW (SignExt16to32 x) (SignExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less16U x y)
	// result: (LessThan (CMPWU (ZeroExt16to32 x) (ZeroExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess32_0(v *Value) bool {
	b := v.Block
	// match: (Less32 x y)
	// result: (LessThan (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess32F_0(v *Value) bool {
	b := v.Block
	// match: (Less32F x y)
	// result: (FLessThan (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FLessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess32U_0(v *Value) bool {
	b := v.Block
	// match: (Less32U x y)
	// result: (LessThan (CMPWU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess64_0(v *Value) bool {
	b := v.Block
	// match: (Less64 x y)
	// result: (LessThan (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess64F_0(v *Value) bool {
	b := v.Block
	// match: (Less64F x y)
	// result: (FLessThan (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FLessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess64U_0(v *Value) bool {
	b := v.Block
	// match: (Less64U x y)
	// result: (LessThan (CMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less8 x y)
	// result: (LessThan (CMPW (SignExt8to32 x) (SignExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLess8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less8U x y)
	// result: (LessThan (CMPWU (ZeroExt8to32 x) (ZeroExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LessThan)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWU, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLoad_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Load <t> ptr mem)
	// cond: (is64BitInt(t) || isPtr(t))
	// result: (MOVDload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is64BitInt(t) || isPtr(t)) {
			break
		}
		v.reset(OpPPC64MOVDload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is32BitInt(t) && isSigned(t)
	// result: (MOVWload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is32BitInt(t) && isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVWload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is32BitInt(t) && !isSigned(t)
	// result: (MOVWZload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is32BitInt(t) && !isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVWZload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is16BitInt(t) && isSigned(t)
	// result: (MOVHload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is16BitInt(t) && isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVHload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is16BitInt(t) && !isSigned(t)
	// result: (MOVHZload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is16BitInt(t) && !isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVHZload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.IsBoolean()
	// result: (MOVBZload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.IsBoolean()) {
			break
		}
		v.reset(OpPPC64MOVBZload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is8BitInt(t) && isSigned(t)
	// result: (MOVBreg (MOVBZload ptr mem))
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is8BitInt(t) && isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVBreg)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZload, typ.UInt8)
		v0.AddArg(ptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is8BitInt(t) && !isSigned(t)
	// result: (MOVBZload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is8BitInt(t) && !isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVBZload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is32BitFloat(t)
	// result: (FMOVSload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is32BitFloat(t)) {
			break
		}
		v.reset(OpPPC64FMOVSload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is64BitFloat(t)
	// result: (FMOVDload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is64BitFloat(t)) {
			break
		}
		v.reset(OpPPC64FMOVDload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpLocalAddr_0(v *Value) bool {
	// match: (LocalAddr {sym} base _)
	// result: (MOVDaddr {sym} base)
	for {
		sym := v.Aux
		_ = v.Args[1]
		base := v.Args[0]
		v.reset(OpPPC64MOVDaddr)
		v.Aux = sym
		v.AddArg(base)
		return true
	}
}
func rewriteValuePPC64_OpLsh16x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh16x16 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -16
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh16x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x32 x (Const64 [c]))
	// cond: uint32(c) < 16
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 16) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh16x32 x (MOVDconst [c]))
	// cond: uint32(c) < 16
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 16) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh16x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh16x32 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -16
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh16x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x64 x (Const64 [c]))
	// cond: uint64(c) < 16
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 16) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh16x64 _ (Const64 [c]))
	// cond: uint64(c) >= 16
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 16) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh16x64 x (MOVDconst [c]))
	// cond: uint64(c) < 16
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 16) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh16x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh16x64 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -16
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh16x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh16x8 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -16
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh32x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh32x16 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh32x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x32 x (Const64 [c]))
	// cond: uint32(c) < 32
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 32) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh32x32 x (MOVDconst [c]))
	// cond: uint32(c) < 32
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 32) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh32x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh32x32 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh32x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x64 x (Const64 [c]))
	// cond: uint64(c) < 32
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 32) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh32x64 _ (Const64 [c]))
	// cond: uint64(c) >= 32
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 32) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh32x64 x (MOVDconst [c]))
	// cond: uint64(c) < 32
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 32) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh32x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh32x64 x (AND y (MOVDconst [31])))
	// result: (SLW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64MOVDconst || v_1_1.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Lsh32x64 x (AND (MOVDconst [31]) y))
	// result: (SLW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		y := v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Lsh32x64 x (ANDconst <typ.Int32> [31] y))
	// result: (SLW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64ANDconst || v_1.Type != typ.Int32 || v_1.AuxInt != 31 {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Lsh32x64 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh32x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh32x8 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh64x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SLD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh64x16 x y)
	// result: (SLD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh64x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x32 x (Const64 [c]))
	// cond: uint32(c) < 64
	// result: (SLDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 64) {
			break
		}
		v.reset(OpPPC64SLDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh64x32 x (MOVDconst [c]))
	// cond: uint32(c) < 64
	// result: (SLDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 64) {
			break
		}
		v.reset(OpPPC64SLDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh64x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SLD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh64x32 x y)
	// result: (SLD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh64x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x64 x (Const64 [c]))
	// cond: uint64(c) < 64
	// result: (SLDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpPPC64SLDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh64x64 _ (Const64 [c]))
	// cond: uint64(c) >= 64
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 64) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh64x64 x (MOVDconst [c]))
	// cond: uint64(c) < 64
	// result: (SLDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpPPC64SLDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh64x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SLD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh64x64 x (AND y (MOVDconst [63])))
	// result: (SLD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64MOVDconst || v_1_1.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Lsh64x64 x (AND (MOVDconst [63]) y))
	// result: (SLD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		y := v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Lsh64x64 x (ANDconst <typ.Int64> [63] y))
	// result: (SLD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64ANDconst || v_1.Type != typ.Int64 || v_1.AuxInt != 63 {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Lsh64x64 x y)
	// result: (SLD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh64x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SLD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh64x8 x y)
	// result: (SLD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh8x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh8x16 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -8
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh8x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x32 x (Const64 [c]))
	// cond: uint32(c) < 8
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 8) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh8x32 x (MOVDconst [c]))
	// cond: uint32(c) < 8
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 8) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh8x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh8x32 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -8
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh8x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x64 x (Const64 [c]))
	// cond: uint64(c) < 8
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 8) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh8x64 _ (Const64 [c]))
	// cond: uint64(c) >= 8
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 8) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh8x64 x (MOVDconst [c]))
	// cond: uint64(c) < 8
	// result: (SLWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 8) {
			break
		}
		v.reset(OpPPC64SLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh8x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh8x64 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -8
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpLsh8x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh8x8 x y)
	// result: (SLW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SLW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -8
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpMod16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod16 x y)
	// result: (Mod32 (SignExt16to32 x) (SignExt16to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpMod32)
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpMod16u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod16u x y)
	// result: (Mod32u (ZeroExt16to32 x) (ZeroExt16to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpMod32u)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpMod32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod32 x y)
	// result: (SUB x (MULLW y (DIVW x y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64MULLW, typ.Int32)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64DIVW, typ.Int32)
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpMod32u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod32u x y)
	// result: (SUB x (MULLW y (DIVWU x y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64MULLW, typ.Int32)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64DIVWU, typ.Int32)
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpMod64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod64 x y)
	// result: (SUB x (MULLD y (DIVD x y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64MULLD, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64DIVD, typ.Int64)
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpMod64u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod64u x y)
	// result: (SUB x (MULLD y (DIVDU x y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64MULLD, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64DIVDU, typ.Int64)
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpMod8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod8 x y)
	// result: (Mod32 (SignExt8to32 x) (SignExt8to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpMod32)
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpMod8u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod8u x y)
	// result: (Mod32u (ZeroExt8to32 x) (ZeroExt8to32 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpMod32u)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpMove_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Move [0] _ _ mem)
	// result: mem
	for {
		if v.AuxInt != 0 {
			break
		}
		mem := v.Args[2]
		v.reset(OpCopy)
		v.Type = mem.Type
		v.AddArg(mem)
		return true
	}
	// match: (Move [1] dst src mem)
	// result: (MOVBstore dst (MOVBZload src mem) mem)
	for {
		if v.AuxInt != 1 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVBstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZload, typ.UInt8)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [2] dst src mem)
	// result: (MOVHstore dst (MOVHZload src mem) mem)
	for {
		if v.AuxInt != 2 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVHstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZload, typ.UInt16)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [4] dst src mem)
	// result: (MOVWstore dst (MOVWZload src mem) mem)
	for {
		if v.AuxInt != 4 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVWstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWZload, typ.UInt32)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [8] {t} dst src mem)
	// cond: t.(*types.Type).Alignment()%4 == 0
	// result: (MOVDstore dst (MOVDload src mem) mem)
	for {
		if v.AuxInt != 8 {
			break
		}
		t := v.Aux
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(t.(*types.Type).Alignment()%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDload, typ.Int64)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [8] dst src mem)
	// result: (MOVWstore [4] dst (MOVWZload [4] src mem) (MOVWstore dst (MOVWZload src mem) mem))
	for {
		if v.AuxInt != 8 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWZload, typ.UInt32)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVWstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpPPC64MOVWZload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [3] dst src mem)
	// result: (MOVBstore [2] dst (MOVBZload [2] src mem) (MOVHstore dst (MOVHload src mem) mem))
	for {
		if v.AuxInt != 3 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = 2
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZload, typ.UInt8)
		v0.AuxInt = 2
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVHstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpPPC64MOVHload, typ.Int16)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [5] dst src mem)
	// result: (MOVBstore [4] dst (MOVBZload [4] src mem) (MOVWstore dst (MOVWZload src mem) mem))
	for {
		if v.AuxInt != 5 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZload, typ.UInt8)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVWstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpPPC64MOVWZload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [6] dst src mem)
	// result: (MOVHstore [4] dst (MOVHZload [4] src mem) (MOVWstore dst (MOVWZload src mem) mem))
	for {
		if v.AuxInt != 6 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZload, typ.UInt16)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVWstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpPPC64MOVWZload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [7] dst src mem)
	// result: (MOVBstore [6] dst (MOVBZload [6] src mem) (MOVHstore [4] dst (MOVHZload [4] src mem) (MOVWstore dst (MOVWZload src mem) mem)))
	for {
		if v.AuxInt != 7 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = 6
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZload, typ.UInt8)
		v0.AuxInt = 6
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVHstore, types.TypeMem)
		v1.AuxInt = 4
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpPPC64MOVHZload, typ.UInt16)
		v2.AuxInt = 4
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v3 := b.NewValue0(v.Pos, OpPPC64MOVWstore, types.TypeMem)
		v3.AddArg(dst)
		v4 := b.NewValue0(v.Pos, OpPPC64MOVWZload, typ.UInt32)
		v4.AddArg(src)
		v4.AddArg(mem)
		v3.AddArg(v4)
		v3.AddArg(mem)
		v1.AddArg(v3)
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValuePPC64_OpMove_10(v *Value) bool {
	// match: (Move [s] dst src mem)
	// cond: s > 8
	// result: (LoweredMove [s] dst src mem)
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s > 8) {
			break
		}
		v.reset(OpPPC64LoweredMove)
		v.AuxInt = s
		v.AddArg(dst)
		v.AddArg(src)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpMul16_0(v *Value) bool {
	// match: (Mul16 x y)
	// result: (MULLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpMul32_0(v *Value) bool {
	// match: (Mul32 x y)
	// result: (MULLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpMul32F_0(v *Value) bool {
	// match: (Mul32F x y)
	// result: (FMULS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FMULS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpMul64_0(v *Value) bool {
	// match: (Mul64 x y)
	// result: (MULLD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULLD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpMul64F_0(v *Value) bool {
	// match: (Mul64F x y)
	// result: (FMUL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FMUL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpMul64uhilo_0(v *Value) bool {
	// match: (Mul64uhilo x y)
	// result: (LoweredMuluhilo x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64LoweredMuluhilo)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpMul8_0(v *Value) bool {
	// match: (Mul8 x y)
	// result: (MULLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64MULLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpNeg16_0(v *Value) bool {
	// match: (Neg16 x)
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpNeg32_0(v *Value) bool {
	// match: (Neg32 x)
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpNeg32F_0(v *Value) bool {
	// match: (Neg32F x)
	// result: (FNEG x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FNEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpNeg64_0(v *Value) bool {
	// match: (Neg64 x)
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpNeg64F_0(v *Value) bool {
	// match: (Neg64F x)
	// result: (FNEG x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FNEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpNeg8_0(v *Value) bool {
	// match: (Neg8 x)
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpNeq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neq16 x y)
	// cond: isSigned(x.Type) && isSigned(y.Type)
	// result: (NotEqual (CMPW (SignExt16to32 x) (SignExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(isSigned(x.Type) && isSigned(y.Type)) {
			break
		}
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Neq16 x y)
	// result: (NotEqual (CMPW (ZeroExt16to32 x) (ZeroExt16to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNeq32_0(v *Value) bool {
	b := v.Block
	// match: (Neq32 x y)
	// result: (NotEqual (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNeq32F_0(v *Value) bool {
	b := v.Block
	// match: (Neq32F x y)
	// result: (NotEqual (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNeq64_0(v *Value) bool {
	b := v.Block
	// match: (Neq64 x y)
	// result: (NotEqual (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNeq64F_0(v *Value) bool {
	b := v.Block
	// match: (Neq64F x y)
	// result: (NotEqual (FCMPU x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64FCMPU, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNeq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neq8 x y)
	// cond: isSigned(x.Type) && isSigned(y.Type)
	// result: (NotEqual (CMPW (SignExt8to32 x) (SignExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(isSigned(x.Type) && isSigned(y.Type)) {
			break
		}
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Neq8 x y)
	// result: (NotEqual (CMPW (ZeroExt8to32 x) (ZeroExt8to32 y)))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPW, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNeqB_0(v *Value) bool {
	// match: (NeqB x y)
	// result: (XOR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpNeqPtr_0(v *Value) bool {
	b := v.Block
	// match: (NeqPtr x y)
	// result: (NotEqual (CMP x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64NotEqual)
		v0 := b.NewValue0(v.Pos, OpPPC64CMP, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpNilCheck_0(v *Value) bool {
	// match: (NilCheck ptr mem)
	// result: (LoweredNilCheck ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredNilCheck)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpNot_0(v *Value) bool {
	// match: (Not x)
	// result: (XORconst [1] x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64XORconst)
		v.AuxInt = 1
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpOffPtr_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (OffPtr [off] ptr)
	// result: (ADD (MOVDconst <typ.Int64> [off]) ptr)
	for {
		off := v.AuxInt
		ptr := v.Args[0]
		v.reset(OpPPC64ADD)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = off
		v.AddArg(v0)
		v.AddArg(ptr)
		return true
	}
}
func rewriteValuePPC64_OpOr16_0(v *Value) bool {
	// match: (Or16 x y)
	// result: (OR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpOr32_0(v *Value) bool {
	// match: (Or32 x y)
	// result: (OR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpOr64_0(v *Value) bool {
	// match: (Or64 x y)
	// result: (OR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpOr8_0(v *Value) bool {
	// match: (Or8 x y)
	// result: (OR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpOrB_0(v *Value) bool {
	// match: (OrB x y)
	// result: (OR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpPPC64ADD_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ADD (SLDconst x [c]) (SRDconst x [d]))
	// cond: d == 64-c
	// result: (ROTLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 64-c) {
			break
		}
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADD (SRDconst x [d]) (SLDconst x [c]))
	// cond: d == 64-c
	// result: (ROTLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLDconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 64-c) {
			break
		}
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADD (SLWconst x [c]) (SRWconst x [d]))
	// cond: d == 32-c
	// result: (ROTLWconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADD (SRWconst x [d]) (SLWconst x [c]))
	// cond: d == 32-c
	// result: (ROTLWconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLWconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADD (SLD x (ANDconst <typ.Int64> [63] y)) (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y))))
	// result: (ROTL x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLD {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64ANDconst || v_0_1.Type != typ.Int64 || v_0_1.AuxInt != 63 {
			break
		}
		y := v_0_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRD {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64SUB || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 64 {
			break
		}
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64ANDconst || v_1_1_1.Type != typ.UInt || v_1_1_1.AuxInt != 63 || y != v_1_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADD (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y))) (SLD x (ANDconst <typ.Int64> [63] y)))
	// result: (ROTL x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRD {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64SUB || v_0_1.Type != typ.UInt {
			break
		}
		_ = v_0_1.Args[1]
		v_0_1_0 := v_0_1.Args[0]
		if v_0_1_0.Op != OpPPC64MOVDconst || v_0_1_0.AuxInt != 64 {
			break
		}
		v_0_1_1 := v_0_1.Args[1]
		if v_0_1_1.Op != OpPPC64ANDconst || v_0_1_1.Type != typ.UInt || v_0_1_1.AuxInt != 63 {
			break
		}
		y := v_0_1_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLD {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.Int64 || v_1_1.AuxInt != 63 || y != v_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADD (SLW x (ANDconst <typ.Int32> [31] y)) (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y))))
	// result: (ROTLW x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLW {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64ANDconst || v_0_1.Type != typ.Int32 || v_0_1.AuxInt != 31 {
			break
		}
		y := v_0_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRW {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64SUB || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 32 {
			break
		}
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64ANDconst || v_1_1_1.Type != typ.UInt || v_1_1_1.AuxInt != 31 || y != v_1_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADD (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y))) (SLW x (ANDconst <typ.Int32> [31] y)))
	// result: (ROTLW x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRW {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64SUB || v_0_1.Type != typ.UInt {
			break
		}
		_ = v_0_1.Args[1]
		v_0_1_0 := v_0_1.Args[0]
		if v_0_1_0.Op != OpPPC64MOVDconst || v_0_1_0.AuxInt != 32 {
			break
		}
		v_0_1_1 := v_0_1.Args[1]
		if v_0_1_1.Op != OpPPC64ANDconst || v_0_1_1.Type != typ.UInt || v_0_1_1.AuxInt != 31 {
			break
		}
		y := v_0_1_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLW {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.Int32 || v_1_1.AuxInt != 31 || y != v_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADD x (MOVDconst [c]))
	// cond: is32Bit(c)
	// result: (ADDconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is32Bit(c)) {
			break
		}
		v.reset(OpPPC64ADDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADD (MOVDconst [c]) x)
	// cond: is32Bit(c)
	// result: (ADDconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(is32Bit(c)) {
			break
		}
		v.reset(OpPPC64ADDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ADDconst_0(v *Value) bool {
	// match: (ADDconst [c] (ADDconst [d] x))
	// cond: is32Bit(c+d)
	// result: (ADDconst [c+d] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(OpPPC64ADDconst)
		v.AuxInt = c + d
		v.AddArg(x)
		return true
	}
	// match: (ADDconst [0] x)
	// result: x
	for {
		if v.AuxInt != 0 {
			break
		}
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ADDconst [c] (MOVDaddr [d] {sym} x))
	// result: (MOVDaddr [c+d] {sym} x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDaddr {
			break
		}
		d := v_0.AuxInt
		sym := v_0.Aux
		x := v_0.Args[0]
		v.reset(OpPPC64MOVDaddr)
		v.AuxInt = c + d
		v.Aux = sym
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64AND_0(v *Value) bool {
	// match: (AND x (NOR y y))
	// result: (ANDN x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64NOR {
			break
		}
		y := v_1.Args[1]
		if y != v_1.Args[0] {
			break
		}
		v.reset(OpPPC64ANDN)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (AND (NOR y y) x)
	// result: (ANDN x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64NOR {
			break
		}
		y := v_0.Args[1]
		if y != v_0.Args[0] {
			break
		}
		v.reset(OpPPC64ANDN)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (AND (MOVDconst [c]) (MOVDconst [d]))
	// result: (MOVDconst [c&d])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		d := v_1.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c & d
		return true
	}
	// match: (AND (MOVDconst [d]) (MOVDconst [c]))
	// result: (MOVDconst [c&d])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		d := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c & d
		return true
	}
	// match: (AND x (MOVDconst [c]))
	// cond: isU16Bit(c)
	// result: (ANDconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(isU16Bit(c)) {
			break
		}
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (AND (MOVDconst [c]) x)
	// cond: isU16Bit(c)
	// result: (ANDconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(isU16Bit(c)) {
			break
		}
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (AND (MOVDconst [c]) y:(MOVWZreg _))
	// cond: c&0xFFFFFFFF == 0xFFFFFFFF
	// result: y
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		y := v.Args[1]
		if y.Op != OpPPC64MOVWZreg || !(c&0xFFFFFFFF == 0xFFFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (AND y:(MOVWZreg _) (MOVDconst [c]))
	// cond: c&0xFFFFFFFF == 0xFFFFFFFF
	// result: y
	for {
		_ = v.Args[1]
		y := v.Args[0]
		if y.Op != OpPPC64MOVWZreg {
			break
		}
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(c&0xFFFFFFFF == 0xFFFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (AND (MOVDconst [0xFFFFFFFF]) y:(MOVWreg x))
	// result: (MOVWZreg x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst || v_0.AuxInt != 0xFFFFFFFF {
			break
		}
		y := v.Args[1]
		if y.Op != OpPPC64MOVWreg {
			break
		}
		x := y.Args[0]
		v.reset(OpPPC64MOVWZreg)
		v.AddArg(x)
		return true
	}
	// match: (AND y:(MOVWreg x) (MOVDconst [0xFFFFFFFF]))
	// result: (MOVWZreg x)
	for {
		_ = v.Args[1]
		y := v.Args[0]
		if y.Op != OpPPC64MOVWreg {
			break
		}
		x := y.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst || v_1.AuxInt != 0xFFFFFFFF {
			break
		}
		v.reset(OpPPC64MOVWZreg)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64AND_10(v *Value) bool {
	// match: (AND (MOVDconst [c]) x:(MOVBZload _ _))
	// result: (ANDconst [c&0xFF] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		x := v.Args[1]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFF
		v.AddArg(x)
		return true
	}
	// match: (AND x:(MOVBZload _ _) (MOVDconst [c]))
	// result: (ANDconst [c&0xFF] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFF
		v.AddArg(x)
		return true
	}
	// match: (AND x:(MOVBZload _ _) (MOVDconst [c]))
	// result: (ANDconst [c&0xFF] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFF
		v.AddArg(x)
		return true
	}
	// match: (AND (MOVDconst [c]) x:(MOVBZload _ _))
	// result: (ANDconst [c&0xFF] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		x := v.Args[1]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFF
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ANDconst_0(v *Value) bool {
	// match: (ANDconst [c] (ANDconst [d] x))
	// result: (ANDconst [c&d] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ANDconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & d
		v.AddArg(x)
		return true
	}
	// match: (ANDconst [-1] x)
	// result: x
	for {
		if v.AuxInt != -1 {
			break
		}
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ANDconst [0] _)
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 0 {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ANDconst [c] y:(MOVBZreg _))
	// cond: c&0xFF == 0xFF
	// result: y
	for {
		c := v.AuxInt
		y := v.Args[0]
		if y.Op != OpPPC64MOVBZreg || !(c&0xFF == 0xFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ANDconst [0xFF] y:(MOVBreg _))
	// result: y
	for {
		if v.AuxInt != 0xFF {
			break
		}
		y := v.Args[0]
		if y.Op != OpPPC64MOVBreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ANDconst [c] y:(MOVHZreg _))
	// cond: c&0xFFFF == 0xFFFF
	// result: y
	for {
		c := v.AuxInt
		y := v.Args[0]
		if y.Op != OpPPC64MOVHZreg || !(c&0xFFFF == 0xFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ANDconst [0xFFFF] y:(MOVHreg _))
	// result: y
	for {
		if v.AuxInt != 0xFFFF {
			break
		}
		y := v.Args[0]
		if y.Op != OpPPC64MOVHreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ANDconst [c] (MOVBreg x))
	// result: (ANDconst [c&0xFF] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVBreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFF
		v.AddArg(x)
		return true
	}
	// match: (ANDconst [c] (MOVBZreg x))
	// result: (ANDconst [c&0xFF] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFF
		v.AddArg(x)
		return true
	}
	// match: (ANDconst [c] (MOVHreg x))
	// result: (ANDconst [c&0xFFFF] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVHreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFFFF
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ANDconst_10(v *Value) bool {
	// match: (ANDconst [c] (MOVHZreg x))
	// result: (ANDconst [c&0xFFFF] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFFFF
		v.AddArg(x)
		return true
	}
	// match: (ANDconst [c] (MOVWreg x))
	// result: (ANDconst [c&0xFFFFFFFF] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVWreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFFFFFFFF
		v.AddArg(x)
		return true
	}
	// match: (ANDconst [c] (MOVWZreg x))
	// result: (ANDconst [c&0xFFFFFFFF] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64ANDconst)
		v.AuxInt = c & 0xFFFFFFFF
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMP_0(v *Value) bool {
	b := v.Block
	// match: (CMP x (MOVDconst [c]))
	// cond: is16Bit(c)
	// result: (CMPconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64CMPconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (CMP (MOVDconst [c]) y)
	// cond: is16Bit(c)
	// result: (InvertFlags (CMPconst y [c]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64InvertFlags)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPconst, types.TypeFlags)
		v0.AuxInt = c
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPU_0(v *Value) bool {
	b := v.Block
	// match: (CMPU x (MOVDconst [c]))
	// cond: isU16Bit(c)
	// result: (CMPUconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(isU16Bit(c)) {
			break
		}
		v.reset(OpPPC64CMPUconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (CMPU (MOVDconst [c]) y)
	// cond: isU16Bit(c)
	// result: (InvertFlags (CMPUconst y [c]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(isU16Bit(c)) {
			break
		}
		v.reset(OpPPC64InvertFlags)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPUconst, types.TypeFlags)
		v0.AuxInt = c
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPUconst_0(v *Value) bool {
	// match: (CMPUconst (MOVDconst [x]) [y])
	// cond: x==y
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(x == y) {
			break
		}
		v.reset(OpPPC64FlagEQ)
		return true
	}
	// match: (CMPUconst (MOVDconst [x]) [y])
	// cond: uint64(x)<uint64(y)
	// result: (FlagLT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(uint64(x) < uint64(y)) {
			break
		}
		v.reset(OpPPC64FlagLT)
		return true
	}
	// match: (CMPUconst (MOVDconst [x]) [y])
	// cond: uint64(x)>uint64(y)
	// result: (FlagGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(uint64(x) > uint64(y)) {
			break
		}
		v.reset(OpPPC64FlagGT)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPW_0(v *Value) bool {
	b := v.Block
	// match: (CMPW x (MOVWreg y))
	// result: (CMPW x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWreg {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64CMPW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPW (MOVWreg x) y)
	// result: (CMPW x y)
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVWreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64CMPW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPW x (MOVDconst [c]))
	// cond: is16Bit(c)
	// result: (CMPWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64CMPWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (CMPW (MOVDconst [c]) y)
	// cond: is16Bit(c)
	// result: (InvertFlags (CMPWconst y [c]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64InvertFlags)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWconst, types.TypeFlags)
		v0.AuxInt = c
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPWU_0(v *Value) bool {
	b := v.Block
	// match: (CMPWU x (MOVWZreg y))
	// result: (CMPWU x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWZreg {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64CMPWU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPWU (MOVWZreg x) y)
	// result: (CMPWU x y)
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64CMPWU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPWU x (MOVDconst [c]))
	// cond: isU16Bit(c)
	// result: (CMPWUconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(isU16Bit(c)) {
			break
		}
		v.reset(OpPPC64CMPWUconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (CMPWU (MOVDconst [c]) y)
	// cond: isU16Bit(c)
	// result: (InvertFlags (CMPWUconst y [c]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(isU16Bit(c)) {
			break
		}
		v.reset(OpPPC64InvertFlags)
		v0 := b.NewValue0(v.Pos, OpPPC64CMPWUconst, types.TypeFlags)
		v0.AuxInt = c
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPWUconst_0(v *Value) bool {
	// match: (CMPWUconst (MOVDconst [x]) [y])
	// cond: int32(x)==int32(y)
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) == int32(y)) {
			break
		}
		v.reset(OpPPC64FlagEQ)
		return true
	}
	// match: (CMPWUconst (MOVDconst [x]) [y])
	// cond: uint32(x)<uint32(y)
	// result: (FlagLT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(uint32(x) < uint32(y)) {
			break
		}
		v.reset(OpPPC64FlagLT)
		return true
	}
	// match: (CMPWUconst (MOVDconst [x]) [y])
	// cond: uint32(x)>uint32(y)
	// result: (FlagGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(uint32(x) > uint32(y)) {
			break
		}
		v.reset(OpPPC64FlagGT)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPWconst_0(v *Value) bool {
	// match: (CMPWconst (MOVDconst [x]) [y])
	// cond: int32(x)==int32(y)
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) == int32(y)) {
			break
		}
		v.reset(OpPPC64FlagEQ)
		return true
	}
	// match: (CMPWconst (MOVDconst [x]) [y])
	// cond: int32(x)<int32(y)
	// result: (FlagLT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) < int32(y)) {
			break
		}
		v.reset(OpPPC64FlagLT)
		return true
	}
	// match: (CMPWconst (MOVDconst [x]) [y])
	// cond: int32(x)>int32(y)
	// result: (FlagGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) > int32(y)) {
			break
		}
		v.reset(OpPPC64FlagGT)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64CMPconst_0(v *Value) bool {
	// match: (CMPconst (MOVDconst [x]) [y])
	// cond: x==y
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(x == y) {
			break
		}
		v.reset(OpPPC64FlagEQ)
		return true
	}
	// match: (CMPconst (MOVDconst [x]) [y])
	// cond: x<y
	// result: (FlagLT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(x < y) {
			break
		}
		v.reset(OpPPC64FlagLT)
		return true
	}
	// match: (CMPconst (MOVDconst [x]) [y])
	// cond: x>y
	// result: (FlagGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		x := v_0.AuxInt
		if !(x > y) {
			break
		}
		v.reset(OpPPC64FlagGT)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64Equal_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Equal (FlagEQ))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (Equal (FlagLT))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Equal (FlagGT))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Equal (InvertFlags x))
	// result: (Equal x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64Equal)
		v.AddArg(x)
		return true
	}
	// match: (Equal cmp)
	// result: (ISELB [2] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 2
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64FABS_0(v *Value) bool {
	// match: (FABS (FMOVDconst [x]))
	// result: (FMOVDconst [auxFrom64F(math.Abs(auxTo64F(x)))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMOVDconst {
			break
		}
		x := v_0.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = auxFrom64F(math.Abs(auxTo64F(x)))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FADD_0(v *Value) bool {
	// match: (FADD (FMUL x y) z)
	// result: (FMADD x y z)
	for {
		z := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMUL {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(OpPPC64FMADD)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
	// match: (FADD z (FMUL x y))
	// result: (FMADD x y z)
	for {
		_ = v.Args[1]
		z := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FMUL {
			break
		}
		y := v_1.Args[1]
		x := v_1.Args[0]
		v.reset(OpPPC64FMADD)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FADDS_0(v *Value) bool {
	// match: (FADDS (FMULS x y) z)
	// result: (FMADDS x y z)
	for {
		z := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMULS {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(OpPPC64FMADDS)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
	// match: (FADDS z (FMULS x y))
	// result: (FMADDS x y z)
	for {
		_ = v.Args[1]
		z := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FMULS {
			break
		}
		y := v_1.Args[1]
		x := v_1.Args[0]
		v.reset(OpPPC64FMADDS)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FCEIL_0(v *Value) bool {
	// match: (FCEIL (FMOVDconst [x]))
	// result: (FMOVDconst [auxFrom64F(math.Ceil(auxTo64F(x)))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMOVDconst {
			break
		}
		x := v_0.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = auxFrom64F(math.Ceil(auxTo64F(x)))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FFLOOR_0(v *Value) bool {
	// match: (FFLOOR (FMOVDconst [x]))
	// result: (FMOVDconst [auxFrom64F(math.Floor(auxTo64F(x)))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMOVDconst {
			break
		}
		x := v_0.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = auxFrom64F(math.Floor(auxTo64F(x)))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FGreaterEqual_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (FGreaterEqual cmp)
	// result: (ISEL [2] (MOVDconst [1]) (ISELB [1] (MOVDconst [1]) cmp) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISEL)
		v.AuxInt = 2
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ISELB, typ.Int32)
		v1.AuxInt = 1
		v2 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v2.AuxInt = 1
		v1.AddArg(v2)
		v1.AddArg(cmp)
		v.AddArg(v1)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64FGreaterThan_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (FGreaterThan cmp)
	// result: (ISELB [1] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 1
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64FLessEqual_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (FLessEqual cmp)
	// result: (ISEL [2] (MOVDconst [1]) (ISELB [0] (MOVDconst [1]) cmp) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISEL)
		v.AuxInt = 2
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ISELB, typ.Int32)
		v1.AuxInt = 0
		v2 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v2.AuxInt = 1
		v1.AddArg(v2)
		v1.AddArg(cmp)
		v.AddArg(v1)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64FLessThan_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (FLessThan cmp)
	// result: (ISELB [0] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 0
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64FMOVDload_0(v *Value) bool {
	// match: (FMOVDload [off] {sym} ptr (MOVDstore [off] {sym} ptr x _))
	// result: (MTVSRD x)
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDstore || v_1.AuxInt != off || v_1.Aux != sym {
			break
		}
		_ = v_1.Args[2]
		if ptr != v_1.Args[0] {
			break
		}
		x := v_1.Args[1]
		v.reset(OpPPC64MTVSRD)
		v.AddArg(x)
		return true
	}
	// match: (FMOVDload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (FMOVDload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64FMOVDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (FMOVDload [off1] {sym} (ADDconst [off2] ptr) mem)
	// cond: is16Bit(off1+off2)
	// result: (FMOVDload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64FMOVDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FMOVDstore_0(v *Value) bool {
	// match: (FMOVDstore [off] {sym} ptr (MTVSRD x) mem)
	// result: (MOVDstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MTVSRD {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVDstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (FMOVDstore [off1] {sym} (ADDconst [off2] ptr) val mem)
	// cond: is16Bit(off1+off2)
	// result: (FMOVDstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64FMOVDstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (FMOVDstore [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) val mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (FMOVDstore [off1+off2] {mergeSym(sym1,sym2)} ptr val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		val := v.Args[1]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64FMOVDstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FMOVSload_0(v *Value) bool {
	// match: (FMOVSload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (FMOVSload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64FMOVSload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (FMOVSload [off1] {sym} (ADDconst [off2] ptr) mem)
	// cond: is16Bit(off1+off2)
	// result: (FMOVSload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64FMOVSload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FMOVSstore_0(v *Value) bool {
	// match: (FMOVSstore [off1] {sym} (ADDconst [off2] ptr) val mem)
	// cond: is16Bit(off1+off2)
	// result: (FMOVSstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64FMOVSstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (FMOVSstore [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) val mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (FMOVSstore [off1+off2] {mergeSym(sym1,sym2)} ptr val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		val := v.Args[1]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64FMOVSstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FNEG_0(v *Value) bool {
	// match: (FNEG (FABS x))
	// result: (FNABS x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FABS {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64FNABS)
		v.AddArg(x)
		return true
	}
	// match: (FNEG (FNABS x))
	// result: (FABS x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FNABS {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64FABS)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FSQRT_0(v *Value) bool {
	// match: (FSQRT (FMOVDconst [x]))
	// result: (FMOVDconst [auxFrom64F(math.Sqrt(auxTo64F(x)))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMOVDconst {
			break
		}
		x := v_0.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = auxFrom64F(math.Sqrt(auxTo64F(x)))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FSUB_0(v *Value) bool {
	// match: (FSUB (FMUL x y) z)
	// result: (FMSUB x y z)
	for {
		z := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMUL {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(OpPPC64FMSUB)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FSUBS_0(v *Value) bool {
	// match: (FSUBS (FMULS x y) z)
	// result: (FMSUBS x y z)
	for {
		z := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMULS {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(OpPPC64FMSUBS)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(z)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64FTRUNC_0(v *Value) bool {
	// match: (FTRUNC (FMOVDconst [x]))
	// result: (FMOVDconst [auxFrom64F(math.Trunc(auxTo64F(x)))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMOVDconst {
			break
		}
		x := v_0.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = auxFrom64F(math.Trunc(auxTo64F(x)))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64GreaterEqual_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (GreaterEqual (FlagEQ))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (GreaterEqual (FlagLT))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (GreaterEqual (FlagGT))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (GreaterEqual (InvertFlags x))
	// result: (LessEqual x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64LessEqual)
		v.AddArg(x)
		return true
	}
	// match: (GreaterEqual cmp)
	// result: (ISELB [4] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 4
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64GreaterThan_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (GreaterThan (FlagEQ))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (GreaterThan (FlagLT))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (GreaterThan (FlagGT))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (GreaterThan (InvertFlags x))
	// result: (LessThan x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64LessThan)
		v.AddArg(x)
		return true
	}
	// match: (GreaterThan cmp)
	// result: (ISELB [1] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 1
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64ISEL_0(v *Value) bool {
	// match: (ISEL [2] x _ (FlagEQ))
	// result: x
	for {
		if v.AuxInt != 2 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [2] _ y (FlagLT))
	// result: y
	for {
		if v.AuxInt != 2 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [2] _ y (FlagGT))
	// result: y
	for {
		if v.AuxInt != 2 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [6] _ y (FlagEQ))
	// result: y
	for {
		if v.AuxInt != 6 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [6] x _ (FlagLT))
	// result: x
	for {
		if v.AuxInt != 6 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [6] x _ (FlagGT))
	// result: x
	for {
		if v.AuxInt != 6 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [0] _ y (FlagEQ))
	// result: y
	for {
		if v.AuxInt != 0 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [0] _ y (FlagGT))
	// result: y
	for {
		if v.AuxInt != 0 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [0] x _ (FlagLT))
	// result: x
	for {
		if v.AuxInt != 0 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [5] _ x (FlagEQ))
	// result: x
	for {
		if v.AuxInt != 5 {
			break
		}
		_ = v.Args[2]
		x := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ISEL_10(v *Value) bool {
	// match: (ISEL [5] _ x (FlagLT))
	// result: x
	for {
		if v.AuxInt != 5 {
			break
		}
		_ = v.Args[2]
		x := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [5] y _ (FlagGT))
	// result: y
	for {
		if v.AuxInt != 5 {
			break
		}
		_ = v.Args[2]
		y := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [1] _ y (FlagEQ))
	// result: y
	for {
		if v.AuxInt != 1 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [1] _ y (FlagLT))
	// result: y
	for {
		if v.AuxInt != 1 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [1] x _ (FlagGT))
	// result: x
	for {
		if v.AuxInt != 1 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [4] x _ (FlagEQ))
	// result: x
	for {
		if v.AuxInt != 4 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [4] x _ (FlagGT))
	// result: x
	for {
		if v.AuxInt != 4 {
			break
		}
		_ = v.Args[2]
		x := v.Args[0]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ISEL [4] _ y (FlagLT))
	// result: y
	for {
		if v.AuxInt != 4 {
			break
		}
		_ = v.Args[2]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (ISEL [n] x y (InvertFlags bool))
	// cond: n%4 == 0
	// result: (ISEL [n+1] x y bool)
	for {
		n := v.AuxInt
		_ = v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64InvertFlags {
			break
		}
		bool := v_2.Args[0]
		if !(n%4 == 0) {
			break
		}
		v.reset(OpPPC64ISEL)
		v.AuxInt = n + 1
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(bool)
		return true
	}
	// match: (ISEL [n] x y (InvertFlags bool))
	// cond: n%4 == 1
	// result: (ISEL [n-1] x y bool)
	for {
		n := v.AuxInt
		_ = v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64InvertFlags {
			break
		}
		bool := v_2.Args[0]
		if !(n%4 == 1) {
			break
		}
		v.reset(OpPPC64ISEL)
		v.AuxInt = n - 1
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(bool)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ISEL_20(v *Value) bool {
	// match: (ISEL [n] x y (InvertFlags bool))
	// cond: n%4 == 2
	// result: (ISEL [n] x y bool)
	for {
		n := v.AuxInt
		_ = v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64InvertFlags {
			break
		}
		bool := v_2.Args[0]
		if !(n%4 == 2) {
			break
		}
		v.reset(OpPPC64ISEL)
		v.AuxInt = n
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(bool)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ISELB_0(v *Value) bool {
	// match: (ISELB [0] _ (FlagLT))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 0 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [0] _ (FlagGT))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 0 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [0] _ (FlagEQ))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 0 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [1] _ (FlagGT))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 1 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [1] _ (FlagLT))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 1 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [1] _ (FlagEQ))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 1 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [2] _ (FlagEQ))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 2 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [2] _ (FlagLT))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 2 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [2] _ (FlagGT))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 2 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [4] _ (FlagLT))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 4 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ISELB_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ISELB [4] _ (FlagGT))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 4 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [4] _ (FlagEQ))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 4 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [5] _ (FlagGT))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 5 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [5] _ (FlagLT))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 5 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [5] _ (FlagEQ))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 5 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [6] _ (FlagEQ))
	// result: (MOVDconst [0])
	for {
		if v.AuxInt != 6 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (ISELB [6] _ (FlagLT))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 6 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [6] _ (FlagGT))
	// result: (MOVDconst [1])
	for {
		if v.AuxInt != 6 {
			break
		}
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (ISELB [n] (MOVDconst [1]) (InvertFlags bool))
	// cond: n%4 == 0
	// result: (ISELB [n+1] (MOVDconst [1]) bool)
	for {
		n := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst || v_0.AuxInt != 1 {
			break
		}
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64InvertFlags {
			break
		}
		bool := v_1.Args[0]
		if !(n%4 == 0) {
			break
		}
		v.reset(OpPPC64ISELB)
		v.AuxInt = n + 1
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(bool)
		return true
	}
	// match: (ISELB [n] (MOVDconst [1]) (InvertFlags bool))
	// cond: n%4 == 1
	// result: (ISELB [n-1] (MOVDconst [1]) bool)
	for {
		n := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst || v_0.AuxInt != 1 {
			break
		}
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64InvertFlags {
			break
		}
		bool := v_1.Args[0]
		if !(n%4 == 1) {
			break
		}
		v.reset(OpPPC64ISELB)
		v.AuxInt = n - 1
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(bool)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ISELB_20(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ISELB [n] (MOVDconst [1]) (InvertFlags bool))
	// cond: n%4 == 2
	// result: (ISELB [n] (MOVDconst [1]) bool)
	for {
		n := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst || v_0.AuxInt != 1 {
			break
		}
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64InvertFlags {
			break
		}
		bool := v_1.Args[0]
		if !(n%4 == 2) {
			break
		}
		v.reset(OpPPC64ISELB)
		v.AuxInt = n
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(bool)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64LessEqual_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (LessEqual (FlagEQ))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (LessEqual (FlagLT))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (LessEqual (FlagGT))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (LessEqual (InvertFlags x))
	// result: (GreaterEqual x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64GreaterEqual)
		v.AddArg(x)
		return true
	}
	// match: (LessEqual cmp)
	// result: (ISELB [5] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 5
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64LessThan_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (LessThan (FlagEQ))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (LessThan (FlagLT))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (LessThan (FlagGT))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (LessThan (InvertFlags x))
	// result: (GreaterThan x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64GreaterThan)
		v.AddArg(x)
		return true
	}
	// match: (LessThan cmp)
	// result: (ISELB [0] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 0
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64MFVSRD_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MFVSRD (FMOVDconst [c]))
	// result: (MOVDconst [c])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FMOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c
		return true
	}
	// match: (MFVSRD x:(FMOVDload [off] {sym} ptr mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVDload [off] {sym} ptr mem)
	for {
		x := v.Args[0]
		if x.Op != OpPPC64FMOVDload {
			break
		}
		off := x.AuxInt
		sym := x.Aux
		mem := x.Args[1]
		ptr := x.Args[0]
		if !(x.Uses == 1 && clobber(x)) {
			break
		}
		b = x.Block
		v0 := b.NewValue0(x.Pos, OpPPC64MOVDload, typ.Int64)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBZload_0(v *Value) bool {
	// match: (MOVBZload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVBZload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVBZload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBZload [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVBZload [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVBZload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBZload [0] {sym} p:(ADD ptr idx) mem)
	// cond: sym == nil && p.Uses == 1
	// result: (MOVBZloadidx ptr idx mem)
	for {
		if v.AuxInt != 0 {
			break
		}
		sym := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		if !(sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVBZloadidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBZloadidx_0(v *Value) bool {
	// match: (MOVBZloadidx ptr (MOVDconst [c]) mem)
	// cond: is16Bit(c)
	// result: (MOVBZload [c] ptr mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVBZload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBZloadidx (MOVDconst [c]) ptr mem)
	// cond: is16Bit(c)
	// result: (MOVBZload [c] ptr mem)
	for {
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVBZload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBZreg_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVBZreg y:(ANDconst [c] _))
	// cond: uint64(c) <= 0xFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64ANDconst {
			break
		}
		c := y.AuxInt
		if !(uint64(c) <= 0xFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVBZreg (SRWconst [c] (MOVBZreg x)))
	// result: (SRWconst [c] (MOVBZreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVBZreg (SRWconst [c] x))
	// cond: sizeof(x.Type) == 8
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(sizeof(x.Type) == 8) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBZreg (SRDconst [c] x))
	// cond: c>=56
	// result: (SRDconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c >= 56) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBZreg (SRWconst [c] x))
	// cond: c>=24
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c >= 24) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBZreg y:(MOVBZreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVBZreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVBZreg (MOVBreg x))
	// result: (MOVBZreg x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVBreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
	// match: (MOVBZreg x:(MOVBZload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVBZreg x:(MOVBZloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVBZreg x:(Arg <t>))
	// cond: is8BitInt(t) && !isSigned(t)
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpArg {
			break
		}
		t := x.Type
		if !(is8BitInt(t) && !isSigned(t)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBZreg_10(v *Value) bool {
	// match: (MOVBZreg (MOVDconst [c]))
	// result: (MOVDconst [int64(uint8(c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = int64(uint8(c))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBreg_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVBreg y:(ANDconst [c] _))
	// cond: uint64(c) <= 0x7F
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64ANDconst {
			break
		}
		c := y.AuxInt
		if !(uint64(c) <= 0x7F) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVBreg (SRAWconst [c] (MOVBreg x)))
	// result: (SRAWconst [c] (MOVBreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVBreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVBreg (SRAWconst [c] x))
	// cond: sizeof(x.Type) == 8
	// result: (SRAWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(sizeof(x.Type) == 8) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBreg (SRDconst [c] x))
	// cond: c>56
	// result: (SRDconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c > 56) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBreg (SRDconst [c] x))
	// cond: c==56
	// result: (SRADconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c == 56) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBreg (SRWconst [c] x))
	// cond: c>24
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c > 24) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBreg (SRWconst [c] x))
	// cond: c==24
	// result: (SRAWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c == 24) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVBreg y:(MOVBreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVBreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVBreg (MOVBZreg x))
	// result: (MOVBreg x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
	// match: (MOVBreg x:(Arg <t>))
	// cond: is8BitInt(t) && isSigned(t)
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpArg {
			break
		}
		t := x.Type
		if !(is8BitInt(t) && isSigned(t)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBreg_10(v *Value) bool {
	// match: (MOVBreg (MOVDconst [c]))
	// result: (MOVDconst [int64(int8(c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = int64(int8(c))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBstore_0(v *Value) bool {
	// match: (MOVBstore [off1] {sym} (ADDconst [off2] x) val mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVBstore [off1+off2] {sym} x val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		val := v.Args[1]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) val mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVBstore [off1+off2] {mergeSym(sym1,sym2)} ptr val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		val := v.Args[1]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVDconst [0]) mem)
	// result: (MOVBstorezero [off] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst || v_1.AuxInt != 0 {
			break
		}
		v.reset(OpPPC64MOVBstorezero)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} p:(ADD ptr idx) val mem)
	// cond: off == 0 && sym == nil && p.Uses == 1
	// result: (MOVBstoreidx ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		val := v.Args[1]
		if !(off == 0 && sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVBstoreidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVBreg x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVBreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVBZreg x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVHreg x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVHreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVHZreg x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVWreg x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVWZreg x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBstore_10(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (MOVBstore [off] {sym} ptr (SRWconst (MOVHreg x) [c]) mem)
	// cond: c <= 8
	// result: (MOVBstore [off] {sym} ptr (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		c := v_1.AuxInt
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVHreg {
			break
		}
		x := v_1_0.Args[0]
		if !(c <= 8) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (SRWconst (MOVHZreg x) [c]) mem)
	// cond: c <= 8
	// result: (MOVBstore [off] {sym} ptr (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		c := v_1.AuxInt
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_1_0.Args[0]
		if !(c <= 8) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (SRWconst (MOVWreg x) [c]) mem)
	// cond: c <= 24
	// result: (MOVBstore [off] {sym} ptr (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		c := v_1.AuxInt
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVWreg {
			break
		}
		x := v_1_0.Args[0]
		if !(c <= 24) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (SRWconst (MOVWZreg x) [c]) mem)
	// cond: c <= 24
	// result: (MOVBstore [off] {sym} ptr (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		c := v_1.AuxInt
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_1_0.Args[0]
		if !(c <= 24) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i1] {s} p (SRWconst w [24]) x0:(MOVBstore [i0] {s} p (SRWconst w [16]) mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)
	// result: (MOVHstore [i0] {s} p (SRWconst <typ.UInt16> w [16]) mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst || v_1.AuxInt != 24 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] {
			break
		}
		x0_1 := x0.Args[1]
		if x0_1.Op != OpPPC64SRWconst || x0_1.AuxInt != 16 || w != x0_1.Args[0] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v0 := b.NewValue0(x0.Pos, OpPPC64SRWconst, typ.UInt16)
		v0.AuxInt = 16
		v0.AddArg(w)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i1] {s} p (SRDconst w [24]) x0:(MOVBstore [i0] {s} p (SRDconst w [16]) mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)
	// result: (MOVHstore [i0] {s} p (SRWconst <typ.UInt16> w [16]) mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst || v_1.AuxInt != 24 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] {
			break
		}
		x0_1 := x0.Args[1]
		if x0_1.Op != OpPPC64SRDconst || x0_1.AuxInt != 16 || w != x0_1.Args[0] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v0 := b.NewValue0(x0.Pos, OpPPC64SRWconst, typ.UInt16)
		v0.AuxInt = 16
		v0.AddArg(w)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i1] {s} p (SRWconst w [8]) x0:(MOVBstore [i0] {s} p w mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)
	// result: (MOVHstore [i0] {s} p w mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst || v_1.AuxInt != 8 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] || w != x0.Args[1] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i1] {s} p (SRDconst w [8]) x0:(MOVBstore [i0] {s} p w mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)
	// result: (MOVHstore [i0] {s} p w mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst || v_1.AuxInt != 8 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] || w != x0.Args[1] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i3] {s} p w x0:(MOVBstore [i2] {s} p (SRWconst w [8]) x1:(MOVBstore [i1] {s} p (SRWconst w [16]) x2:(MOVBstore [i0] {s} p (SRWconst w [24]) mem))))
	// cond: !config.BigEndian && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && clobber(x0) && clobber(x1) && clobber(x2)
	// result: (MOVWBRstore (MOVDaddr <typ.Uintptr> [i0] {s} p) w mem)
	for {
		i3 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		w := v.Args[1]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i2 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] {
			break
		}
		x0_1 := x0.Args[1]
		if x0_1.Op != OpPPC64SRWconst || x0_1.AuxInt != 8 || w != x0_1.Args[0] {
			break
		}
		x1 := x0.Args[2]
		if x1.Op != OpPPC64MOVBstore {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] {
			break
		}
		x1_1 := x1.Args[1]
		if x1_1.Op != OpPPC64SRWconst || x1_1.AuxInt != 16 || w != x1_1.Args[0] {
			break
		}
		x2 := x1.Args[2]
		if x2.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		mem := x2.Args[2]
		if p != x2.Args[0] {
			break
		}
		x2_1 := x2.Args[1]
		if x2_1.Op != OpPPC64SRWconst || x2_1.AuxInt != 24 || w != x2_1.Args[0] || !(!config.BigEndian && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && clobber(x0) && clobber(x1) && clobber(x2)) {
			break
		}
		v.reset(OpPPC64MOVWBRstore)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v.AddArg(v0)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i1] {s} p w x0:(MOVBstore [i0] {s} p (SRWconst w [8]) mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)
	// result: (MOVHBRstore (MOVDaddr <typ.Uintptr> [i0] {s} p) w mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		w := v.Args[1]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] {
			break
		}
		x0_1 := x0.Args[1]
		if x0_1.Op != OpPPC64SRWconst || x0_1.AuxInt != 8 || w != x0_1.Args[0] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+1 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVHBRstore)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v.AddArg(v0)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBstore_20(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (MOVBstore [i7] {s} p (SRDconst w [56]) x0:(MOVBstore [i6] {s} p (SRDconst w [48]) x1:(MOVBstore [i5] {s} p (SRDconst w [40]) x2:(MOVBstore [i4] {s} p (SRDconst w [32]) x3:(MOVWstore [i0] {s} p w mem)))))
	// cond: !config.BigEndian && i0%4 == 0 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3)
	// result: (MOVDstore [i0] {s} p w mem)
	for {
		i7 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst || v_1.AuxInt != 56 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i6 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] {
			break
		}
		x0_1 := x0.Args[1]
		if x0_1.Op != OpPPC64SRDconst || x0_1.AuxInt != 48 || w != x0_1.Args[0] {
			break
		}
		x1 := x0.Args[2]
		if x1.Op != OpPPC64MOVBstore {
			break
		}
		i5 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] {
			break
		}
		x1_1 := x1.Args[1]
		if x1_1.Op != OpPPC64SRDconst || x1_1.AuxInt != 40 || w != x1_1.Args[0] {
			break
		}
		x2 := x1.Args[2]
		if x2.Op != OpPPC64MOVBstore {
			break
		}
		i4 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] {
			break
		}
		x2_1 := x2.Args[1]
		if x2_1.Op != OpPPC64SRDconst || x2_1.AuxInt != 32 || w != x2_1.Args[0] {
			break
		}
		x3 := x2.Args[2]
		if x3.Op != OpPPC64MOVWstore {
			break
		}
		i0 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		mem := x3.Args[2]
		if p != x3.Args[0] || w != x3.Args[1] || !(!config.BigEndian && i0%4 == 0 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3)) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i7] {s} p w x0:(MOVBstore [i6] {s} p (SRDconst w [8]) x1:(MOVBstore [i5] {s} p (SRDconst w [16]) x2:(MOVBstore [i4] {s} p (SRDconst w [24]) x3:(MOVBstore [i3] {s} p (SRDconst w [32]) x4:(MOVBstore [i2] {s} p (SRDconst w [40]) x5:(MOVBstore [i1] {s} p (SRDconst w [48]) x6:(MOVBstore [i0] {s} p (SRDconst w [56]) mem))))))))
	// cond: !config.BigEndian && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6)
	// result: (MOVDBRstore (MOVDaddr <typ.Uintptr> [i0] {s} p) w mem)
	for {
		i7 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		w := v.Args[1]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVBstore {
			break
		}
		i6 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] {
			break
		}
		x0_1 := x0.Args[1]
		if x0_1.Op != OpPPC64SRDconst || x0_1.AuxInt != 8 || w != x0_1.Args[0] {
			break
		}
		x1 := x0.Args[2]
		if x1.Op != OpPPC64MOVBstore {
			break
		}
		i5 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] {
			break
		}
		x1_1 := x1.Args[1]
		if x1_1.Op != OpPPC64SRDconst || x1_1.AuxInt != 16 || w != x1_1.Args[0] {
			break
		}
		x2 := x1.Args[2]
		if x2.Op != OpPPC64MOVBstore {
			break
		}
		i4 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] {
			break
		}
		x2_1 := x2.Args[1]
		if x2_1.Op != OpPPC64SRDconst || x2_1.AuxInt != 24 || w != x2_1.Args[0] {
			break
		}
		x3 := x2.Args[2]
		if x3.Op != OpPPC64MOVBstore {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[2]
		if p != x3.Args[0] {
			break
		}
		x3_1 := x3.Args[1]
		if x3_1.Op != OpPPC64SRDconst || x3_1.AuxInt != 32 || w != x3_1.Args[0] {
			break
		}
		x4 := x3.Args[2]
		if x4.Op != OpPPC64MOVBstore {
			break
		}
		i2 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[2]
		if p != x4.Args[0] {
			break
		}
		x4_1 := x4.Args[1]
		if x4_1.Op != OpPPC64SRDconst || x4_1.AuxInt != 40 || w != x4_1.Args[0] {
			break
		}
		x5 := x4.Args[2]
		if x5.Op != OpPPC64MOVBstore {
			break
		}
		i1 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[2]
		if p != x5.Args[0] {
			break
		}
		x5_1 := x5.Args[1]
		if x5_1.Op != OpPPC64SRDconst || x5_1.AuxInt != 48 || w != x5_1.Args[0] {
			break
		}
		x6 := x5.Args[2]
		if x6.Op != OpPPC64MOVBstore {
			break
		}
		i0 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		mem := x6.Args[2]
		if p != x6.Args[0] {
			break
		}
		x6_1 := x6.Args[1]
		if x6_1.Op != OpPPC64SRDconst || x6_1.AuxInt != 56 || w != x6_1.Args[0] || !(!config.BigEndian && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6)) {
			break
		}
		v.reset(OpPPC64MOVDBRstore)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v.AddArg(v0)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBstoreidx_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVBstoreidx ptr (MOVDconst [c]) val mem)
	// cond: is16Bit(c)
	// result: (MOVBstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		val := v.Args[2]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx (MOVDconst [c]) ptr val mem)
	// cond: is16Bit(c)
	// result: (MOVBstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		val := v.Args[2]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (MOVBreg x) mem)
	// result: (MOVBstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVBreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (MOVBZreg x) mem)
	// result: (MOVBstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (MOVHreg x) mem)
	// result: (MOVBstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVHreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (MOVHZreg x) mem)
	// result: (MOVBstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (MOVWreg x) mem)
	// result: (MOVBstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVWreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (MOVWZreg x) mem)
	// result: (MOVBstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (SRWconst (MOVHreg x) [c]) mem)
	// cond: c <= 8
	// result: (MOVBstoreidx [off] {sym} ptr idx (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64SRWconst {
			break
		}
		c := v_2.AuxInt
		v_2_0 := v_2.Args[0]
		if v_2_0.Op != OpPPC64MOVHreg {
			break
		}
		x := v_2_0.Args[0]
		if !(c <= 8) {
			break
		}
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (SRWconst (MOVHZreg x) [c]) mem)
	// cond: c <= 8
	// result: (MOVBstoreidx [off] {sym} ptr idx (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64SRWconst {
			break
		}
		c := v_2.AuxInt
		v_2_0 := v_2.Args[0]
		if v_2_0.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_2_0.Args[0]
		if !(c <= 8) {
			break
		}
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBstoreidx_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVBstoreidx [off] {sym} ptr idx (SRWconst (MOVWreg x) [c]) mem)
	// cond: c <= 24
	// result: (MOVBstoreidx [off] {sym} ptr idx (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64SRWconst {
			break
		}
		c := v_2.AuxInt
		v_2_0 := v_2.Args[0]
		if v_2_0.Op != OpPPC64MOVWreg {
			break
		}
		x := v_2_0.Args[0]
		if !(c <= 24) {
			break
		}
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx [off] {sym} ptr idx (SRWconst (MOVWZreg x) [c]) mem)
	// cond: c <= 24
	// result: (MOVBstoreidx [off] {sym} ptr idx (SRWconst <typ.UInt32> x [c]) mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64SRWconst {
			break
		}
		c := v_2.AuxInt
		v_2_0 := v_2.Args[0]
		if v_2_0.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_2_0.Args[0]
		if !(c <= 24) {
			break
		}
		v.reset(OpPPC64MOVBstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v0 := b.NewValue0(v.Pos, OpPPC64SRWconst, typ.UInt32)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVBstorezero_0(v *Value) bool {
	// match: (MOVBstorezero [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVBstorezero [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVBstorezero)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstorezero [off1] {sym1} p:(MOVDaddr [off2] {sym2} x) mem)
	// cond: canMergeSym(sym1,sym2) && (x.Op != OpSB || p.Uses == 1)
	// result: (MOVBstorezero [off1+off2] {mergeSym(sym1,sym2)} x mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		x := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (x.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVBstorezero)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVDload_0(v *Value) bool {
	// match: (MOVDload [off] {sym} ptr (FMOVDstore [off] {sym} ptr x _))
	// result: (MFVSRD x)
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64FMOVDstore || v_1.AuxInt != off || v_1.Aux != sym {
			break
		}
		_ = v_1.Args[2]
		if ptr != v_1.Args[0] {
			break
		}
		x := v_1.Args[1]
		v.reset(OpPPC64MFVSRD)
		v.AddArg(x)
		return true
	}
	// match: (MOVDload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0
	// result: (MOVDload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDload [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2) && (off1+off2)%4 == 0
	// result: (MOVDload [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1+off2) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDload [0] {sym} p:(ADD ptr idx) mem)
	// cond: sym == nil && p.Uses == 1
	// result: (MOVDloadidx ptr idx mem)
	for {
		if v.AuxInt != 0 {
			break
		}
		sym := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		if !(sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVDloadidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVDloadidx_0(v *Value) bool {
	// match: (MOVDloadidx ptr (MOVDconst [c]) mem)
	// cond: is16Bit(c) && c%4 == 0
	// result: (MOVDload [c] ptr mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c) && c%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDloadidx (MOVDconst [c]) ptr mem)
	// cond: is16Bit(c) && c%4 == 0
	// result: (MOVDload [c] ptr mem)
	for {
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		if !(is16Bit(c) && c%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVDstore_0(v *Value) bool {
	// match: (MOVDstore [off] {sym} ptr (MFVSRD x) mem)
	// result: (FMOVDstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MFVSRD {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64FMOVDstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDstore [off1] {sym} (ADDconst [off2] x) val mem)
	// cond: is16Bit(off1+off2) && (off1+off2)%4 == 0
	// result: (MOVDstore [off1+off2] {sym} x val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		val := v.Args[1]
		if !(is16Bit(off1+off2) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDstore [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) val mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0
	// result: (MOVDstore [off1+off2] {mergeSym(sym1,sym2)} ptr val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		val := v.Args[1]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDstore [off] {sym} ptr (MOVDconst [0]) mem)
	// result: (MOVDstorezero [off] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst || v_1.AuxInt != 0 {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDstore [off] {sym} p:(ADD ptr idx) val mem)
	// cond: off == 0 && sym == nil && p.Uses == 1
	// result: (MOVDstoreidx ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		val := v.Args[1]
		if !(off == 0 && sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVDstoreidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVDstoreidx_0(v *Value) bool {
	// match: (MOVDstoreidx ptr (MOVDconst [c]) val mem)
	// cond: is16Bit(c) && c%4 == 0
	// result: (MOVDstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		val := v.Args[2]
		if !(is16Bit(c) && c%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDstoreidx (MOVDconst [c]) ptr val mem)
	// cond: is16Bit(c) && c%4 == 0
	// result: (MOVDstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		val := v.Args[2]
		if !(is16Bit(c) && c%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVDstorezero_0(v *Value) bool {
	// match: (MOVDstorezero [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2) && (off1+off2)%4 == 0
	// result: (MOVDstorezero [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1+off2) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVDstorezero [off1] {sym1} p:(MOVDaddr [off2] {sym2} x) mem)
	// cond: canMergeSym(sym1,sym2) && (x.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0
	// result: (MOVDstorezero [off1+off2] {mergeSym(sym1,sym2)} x mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		x := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (x.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHBRstore_0(v *Value) bool {
	// match: (MOVHBRstore {sym} ptr (MOVHreg x) mem)
	// result: (MOVHBRstore {sym} ptr x mem)
	for {
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVHreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHBRstore)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHBRstore {sym} ptr (MOVHZreg x) mem)
	// result: (MOVHBRstore {sym} ptr x mem)
	for {
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHBRstore)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHBRstore {sym} ptr (MOVWreg x) mem)
	// result: (MOVHBRstore {sym} ptr x mem)
	for {
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHBRstore)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHBRstore {sym} ptr (MOVWZreg x) mem)
	// result: (MOVHBRstore {sym} ptr x mem)
	for {
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHBRstore)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHZload_0(v *Value) bool {
	// match: (MOVHZload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVHZload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVHZload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHZload [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVHZload [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVHZload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHZload [0] {sym} p:(ADD ptr idx) mem)
	// cond: sym == nil && p.Uses == 1
	// result: (MOVHZloadidx ptr idx mem)
	for {
		if v.AuxInt != 0 {
			break
		}
		sym := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		if !(sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVHZloadidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHZloadidx_0(v *Value) bool {
	// match: (MOVHZloadidx ptr (MOVDconst [c]) mem)
	// cond: is16Bit(c)
	// result: (MOVHZload [c] ptr mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVHZload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHZloadidx (MOVDconst [c]) ptr mem)
	// cond: is16Bit(c)
	// result: (MOVHZload [c] ptr mem)
	for {
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVHZload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHZreg_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVHZreg y:(ANDconst [c] _))
	// cond: uint64(c) <= 0xFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64ANDconst {
			break
		}
		c := y.AuxInt
		if !(uint64(c) <= 0xFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVHZreg (SRWconst [c] (MOVBZreg x)))
	// result: (SRWconst [c] (MOVBZreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVHZreg (SRWconst [c] (MOVHZreg x)))
	// result: (SRWconst [c] (MOVHZreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVHZreg (SRWconst [c] x))
	// cond: sizeof(x.Type) <= 16
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(sizeof(x.Type) <= 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg (SRDconst [c] x))
	// cond: c>=48
	// result: (SRDconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c >= 48) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg (SRWconst [c] x))
	// cond: c>=16
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c >= 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg y:(MOVHZreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHZreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVHZreg y:(MOVBZreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVBZreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVHZreg y:(MOVHBRload _ _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHBRload {
			break
		}
		_ = y.Args[1]
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVHZreg y:(MOVHreg x))
	// result: (MOVHZreg x)
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHreg {
			break
		}
		x := y.Args[0]
		v.reset(OpPPC64MOVHZreg)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHZreg_10(v *Value) bool {
	// match: (MOVHZreg x:(MOVBZload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg x:(MOVBZloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg x:(MOVHZload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg x:(MOVHZloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHZloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg x:(Arg <t>))
	// cond: (is8BitInt(t) || is16BitInt(t)) && !isSigned(t)
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpArg {
			break
		}
		t := x.Type
		if !((is8BitInt(t) || is16BitInt(t)) && !isSigned(t)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHZreg (MOVDconst [c]))
	// result: (MOVDconst [int64(uint16(c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = int64(uint16(c))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHload_0(v *Value) bool {
	// match: (MOVHload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVHload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVHload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHload [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVHload [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVHload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHload [0] {sym} p:(ADD ptr idx) mem)
	// cond: sym == nil && p.Uses == 1
	// result: (MOVHloadidx ptr idx mem)
	for {
		if v.AuxInt != 0 {
			break
		}
		sym := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		if !(sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVHloadidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHloadidx_0(v *Value) bool {
	// match: (MOVHloadidx ptr (MOVDconst [c]) mem)
	// cond: is16Bit(c)
	// result: (MOVHload [c] ptr mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVHload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHloadidx (MOVDconst [c]) ptr mem)
	// cond: is16Bit(c)
	// result: (MOVHload [c] ptr mem)
	for {
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVHload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHreg_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVHreg y:(ANDconst [c] _))
	// cond: uint64(c) <= 0x7FFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64ANDconst {
			break
		}
		c := y.AuxInt
		if !(uint64(c) <= 0x7FFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVHreg (SRAWconst [c] (MOVBreg x)))
	// result: (SRAWconst [c] (MOVBreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVBreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVHreg (SRAWconst [c] (MOVHreg x)))
	// result: (SRAWconst [c] (MOVHreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVHreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVHreg (SRAWconst [c] x))
	// cond: sizeof(x.Type) <= 16
	// result: (SRAWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(sizeof(x.Type) <= 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg (SRDconst [c] x))
	// cond: c>48
	// result: (SRDconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c > 48) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg (SRDconst [c] x))
	// cond: c==48
	// result: (SRADconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c == 48) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg (SRWconst [c] x))
	// cond: c>16
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c > 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg (SRWconst [c] x))
	// cond: c==16
	// result: (SRAWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c == 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg y:(MOVHreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVHreg y:(MOVBreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVBreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHreg_10(v *Value) bool {
	// match: (MOVHreg y:(MOVHZreg x))
	// result: (MOVHreg x)
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHZreg {
			break
		}
		x := y.Args[0]
		v.reset(OpPPC64MOVHreg)
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg x:(MOVHload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg x:(MOVHloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg x:(Arg <t>))
	// cond: (is8BitInt(t) || is16BitInt(t)) && isSigned(t)
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpArg {
			break
		}
		t := x.Type
		if !((is8BitInt(t) || is16BitInt(t)) && isSigned(t)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVHreg (MOVDconst [c]))
	// result: (MOVDconst [int64(int16(c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = int64(int16(c))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHstore_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVHstore [off1] {sym} (ADDconst [off2] x) val mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVHstore [off1+off2] {sym} x val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		val := v.Args[1]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) val mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVHstore [off1+off2] {mergeSym(sym1,sym2)} ptr val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		val := v.Args[1]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off] {sym} ptr (MOVDconst [0]) mem)
	// result: (MOVHstorezero [off] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst || v_1.AuxInt != 0 {
			break
		}
		v.reset(OpPPC64MOVHstorezero)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off] {sym} p:(ADD ptr idx) val mem)
	// cond: off == 0 && sym == nil && p.Uses == 1
	// result: (MOVHstoreidx ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		val := v.Args[1]
		if !(off == 0 && sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVHstoreidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off] {sym} ptr (MOVHreg x) mem)
	// result: (MOVHstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVHreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off] {sym} ptr (MOVHZreg x) mem)
	// result: (MOVHstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off] {sym} ptr (MOVWreg x) mem)
	// result: (MOVHstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [off] {sym} ptr (MOVWZreg x) mem)
	// result: (MOVHstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [i1] {s} p (SRWconst w [16]) x0:(MOVHstore [i0] {s} p w mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+2 && clobber(x0)
	// result: (MOVWstore [i0] {s} p w mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst || v_1.AuxInt != 16 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVHstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] || w != x0.Args[1] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+2 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstore [i1] {s} p (SRDconst w [16]) x0:(MOVHstore [i0] {s} p w mem))
	// cond: !config.BigEndian && x0.Uses == 1 && i1 == i0+2 && clobber(x0)
	// result: (MOVWstore [i0] {s} p w mem)
	for {
		i1 := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst || v_1.AuxInt != 16 {
			break
		}
		w := v_1.Args[0]
		x0 := v.Args[2]
		if x0.Op != OpPPC64MOVHstore {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		mem := x0.Args[2]
		if p != x0.Args[0] || w != x0.Args[1] || !(!config.BigEndian && x0.Uses == 1 && i1 == i0+2 && clobber(x0)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = i0
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHstoreidx_0(v *Value) bool {
	// match: (MOVHstoreidx ptr (MOVDconst [c]) val mem)
	// cond: is16Bit(c)
	// result: (MOVHstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		val := v.Args[2]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstoreidx (MOVDconst [c]) ptr val mem)
	// cond: is16Bit(c)
	// result: (MOVHstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		val := v.Args[2]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstoreidx [off] {sym} ptr idx (MOVHreg x) mem)
	// result: (MOVHstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVHreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVHstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstoreidx [off] {sym} ptr idx (MOVHZreg x) mem)
	// result: (MOVHstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVHstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstoreidx [off] {sym} ptr idx (MOVWreg x) mem)
	// result: (MOVHstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVWreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVHstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstoreidx [off] {sym} ptr idx (MOVWZreg x) mem)
	// result: (MOVHstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVHstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVHstorezero_0(v *Value) bool {
	// match: (MOVHstorezero [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVHstorezero [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVHstorezero)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVHstorezero [off1] {sym1} p:(MOVDaddr [off2] {sym2} x) mem)
	// cond: canMergeSym(sym1,sym2) && (x.Op != OpSB || p.Uses == 1)
	// result: (MOVHstorezero [off1+off2] {mergeSym(sym1,sym2)} x mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		x := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (x.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVHstorezero)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWBRstore_0(v *Value) bool {
	// match: (MOVWBRstore {sym} ptr (MOVWreg x) mem)
	// result: (MOVWBRstore {sym} ptr x mem)
	for {
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVWBRstore)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWBRstore {sym} ptr (MOVWZreg x) mem)
	// result: (MOVWBRstore {sym} ptr x mem)
	for {
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVWBRstore)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWZload_0(v *Value) bool {
	// match: (MOVWZload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVWZload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVWZload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWZload [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVWZload [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVWZload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWZload [0] {sym} p:(ADD ptr idx) mem)
	// cond: sym == nil && p.Uses == 1
	// result: (MOVWZloadidx ptr idx mem)
	for {
		if v.AuxInt != 0 {
			break
		}
		sym := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		if !(sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVWZloadidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWZloadidx_0(v *Value) bool {
	// match: (MOVWZloadidx ptr (MOVDconst [c]) mem)
	// cond: is16Bit(c)
	// result: (MOVWZload [c] ptr mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVWZload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWZloadidx (MOVDconst [c]) ptr mem)
	// cond: is16Bit(c)
	// result: (MOVWZload [c] ptr mem)
	for {
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVWZload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWZreg_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVWZreg y:(ANDconst [c] _))
	// cond: uint64(c) <= 0xFFFFFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64ANDconst {
			break
		}
		c := y.AuxInt
		if !(uint64(c) <= 0xFFFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg y:(AND (MOVDconst [c]) _))
	// cond: uint64(c) <= 0xFFFFFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64AND {
			break
		}
		_ = y.Args[1]
		y_0 := y.Args[0]
		if y_0.Op != OpPPC64MOVDconst {
			break
		}
		c := y_0.AuxInt
		if !(uint64(c) <= 0xFFFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg y:(AND _ (MOVDconst [c])))
	// cond: uint64(c) <= 0xFFFFFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64AND {
			break
		}
		_ = y.Args[1]
		y_1 := y.Args[1]
		if y_1.Op != OpPPC64MOVDconst {
			break
		}
		c := y_1.AuxInt
		if !(uint64(c) <= 0xFFFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg (SRWconst [c] (MOVBZreg x)))
	// result: (SRWconst [c] (MOVBZreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVBZreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVWZreg (SRWconst [c] (MOVHZreg x)))
	// result: (SRWconst [c] (MOVHZreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVHZreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVWZreg (SRWconst [c] (MOVWZreg x)))
	// result: (SRWconst [c] (MOVWZreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVWZreg (SRWconst [c] x))
	// cond: sizeof(x.Type) <= 32
	// result: (SRWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(sizeof(x.Type) <= 32) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg (SRDconst [c] x))
	// cond: c>=32
	// result: (SRDconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c >= 32) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg y:(MOVWZreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVWZreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg y:(MOVHZreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHZreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWZreg_10(v *Value) bool {
	// match: (MOVWZreg y:(MOVBZreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVBZreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg y:(MOVHBRload _ _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHBRload {
			break
		}
		_ = y.Args[1]
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg y:(MOVWBRload _ _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVWBRload {
			break
		}
		_ = y.Args[1]
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWZreg y:(MOVWreg x))
	// result: (MOVWZreg x)
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVWreg {
			break
		}
		x := y.Args[0]
		v.reset(OpPPC64MOVWZreg)
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg x:(MOVBZload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg x:(MOVBZloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVBZloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg x:(MOVHZload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg x:(MOVHZloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHZloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg x:(MOVWZload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVWZload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg x:(MOVWZloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVWZloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWZreg_20(v *Value) bool {
	// match: (MOVWZreg x:(Arg <t>))
	// cond: (is8BitInt(t) || is16BitInt(t) || is32BitInt(t)) && !isSigned(t)
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpArg {
			break
		}
		t := x.Type
		if !((is8BitInt(t) || is16BitInt(t) || is32BitInt(t)) && !isSigned(t)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWZreg (MOVDconst [c]))
	// result: (MOVDconst [int64(uint32(c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = int64(uint32(c))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWload_0(v *Value) bool {
	// match: (MOVWload [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0
	// result: (MOVWload [off1+off2] {mergeSym(sym1,sym2)} ptr mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVWload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2) && (off1+off2)%4 == 0
	// result: (MOVWload [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1+off2) && (off1+off2)%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVWload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [0] {sym} p:(ADD ptr idx) mem)
	// cond: sym == nil && p.Uses == 1
	// result: (MOVWloadidx ptr idx mem)
	for {
		if v.AuxInt != 0 {
			break
		}
		sym := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		if !(sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVWloadidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWloadidx_0(v *Value) bool {
	// match: (MOVWloadidx ptr (MOVDconst [c]) mem)
	// cond: is16Bit(c) && c%4 == 0
	// result: (MOVWload [c] ptr mem)
	for {
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is16Bit(c) && c%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVWload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx (MOVDconst [c]) ptr mem)
	// cond: is16Bit(c) && c%4 == 0
	// result: (MOVWload [c] ptr mem)
	for {
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		if !(is16Bit(c) && c%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVWload)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWreg_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MOVWreg y:(ANDconst [c] _))
	// cond: uint64(c) <= 0xFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64ANDconst {
			break
		}
		c := y.AuxInt
		if !(uint64(c) <= 0xFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWreg y:(AND (MOVDconst [c]) _))
	// cond: uint64(c) <= 0x7FFFFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64AND {
			break
		}
		_ = y.Args[1]
		y_0 := y.Args[0]
		if y_0.Op != OpPPC64MOVDconst {
			break
		}
		c := y_0.AuxInt
		if !(uint64(c) <= 0x7FFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWreg y:(AND _ (MOVDconst [c])))
	// cond: uint64(c) <= 0x7FFFFFFF
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64AND {
			break
		}
		_ = y.Args[1]
		y_1 := y.Args[1]
		if y_1.Op != OpPPC64MOVDconst {
			break
		}
		c := y_1.AuxInt
		if !(uint64(c) <= 0x7FFFFFFF) {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWreg (SRAWconst [c] (MOVBreg x)))
	// result: (SRAWconst [c] (MOVBreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVBreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVWreg (SRAWconst [c] (MOVHreg x)))
	// result: (SRAWconst [c] (MOVHreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVHreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVWreg (SRAWconst [c] (MOVWreg x)))
	// result: (SRAWconst [c] (MOVWreg x))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64MOVWreg {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MOVWreg (SRAWconst [c] x))
	// cond: sizeof(x.Type) <= 32
	// result: (SRAWconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRAWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(sizeof(x.Type) <= 32) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg (SRDconst [c] x))
	// cond: c>32
	// result: (SRDconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c > 32) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg (SRDconst [c] x))
	// cond: c==32
	// result: (SRADconst [c] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c == 32) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg y:(MOVWreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVWreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWreg_10(v *Value) bool {
	// match: (MOVWreg y:(MOVHreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVHreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWreg y:(MOVBreg _))
	// result: y
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVBreg {
			break
		}
		v.reset(OpCopy)
		v.Type = y.Type
		v.AddArg(y)
		return true
	}
	// match: (MOVWreg y:(MOVWZreg x))
	// result: (MOVWreg x)
	for {
		y := v.Args[0]
		if y.Op != OpPPC64MOVWZreg {
			break
		}
		x := y.Args[0]
		v.reset(OpPPC64MOVWreg)
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg x:(MOVHload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg x:(MOVHloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVHloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg x:(MOVWload _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVWload {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg x:(MOVWloadidx _ _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVWloadidx {
			break
		}
		_ = x.Args[2]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg x:(Arg <t>))
	// cond: (is8BitInt(t) || is16BitInt(t) || is32BitInt(t)) && isSigned(t)
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpArg {
			break
		}
		t := x.Type
		if !((is8BitInt(t) || is16BitInt(t) || is32BitInt(t)) && isSigned(t)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVWreg (MOVDconst [c]))
	// result: (MOVDconst [int64(int32(c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = int64(int32(c))
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWstore_0(v *Value) bool {
	// match: (MOVWstore [off1] {sym} (ADDconst [off2] x) val mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVWstore [off1+off2] {sym} x val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		val := v.Args[1]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off1] {sym1} p:(MOVDaddr [off2] {sym2} ptr) val mem)
	// cond: canMergeSym(sym1,sym2) && (ptr.Op != OpSB || p.Uses == 1)
	// result: (MOVWstore [off1+off2] {mergeSym(sym1,sym2)} ptr val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		ptr := p.Args[0]
		val := v.Args[1]
		if !(canMergeSym(sym1, sym2) && (ptr.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} ptr (MOVDconst [0]) mem)
	// result: (MOVWstorezero [off] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst || v_1.AuxInt != 0 {
			break
		}
		v.reset(OpPPC64MOVWstorezero)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} p:(ADD ptr idx) val mem)
	// cond: off == 0 && sym == nil && p.Uses == 1
	// result: (MOVWstoreidx ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		p := v.Args[0]
		if p.Op != OpPPC64ADD {
			break
		}
		idx := p.Args[1]
		ptr := p.Args[0]
		val := v.Args[1]
		if !(off == 0 && sym == nil && p.Uses == 1) {
			break
		}
		v.reset(OpPPC64MOVWstoreidx)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} ptr (MOVWreg x) mem)
	// result: (MOVWstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} ptr (MOVWZreg x) mem)
	// result: (MOVWstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_1.Args[0]
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWstoreidx_0(v *Value) bool {
	// match: (MOVWstoreidx ptr (MOVDconst [c]) val mem)
	// cond: is16Bit(c)
	// result: (MOVWstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		val := v.Args[2]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx (MOVDconst [c]) ptr val mem)
	// cond: is16Bit(c)
	// result: (MOVWstore [c] ptr val mem)
	for {
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		ptr := v.Args[1]
		val := v.Args[2]
		if !(is16Bit(c)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AuxInt = c
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx [off] {sym} ptr idx (MOVWreg x) mem)
	// result: (MOVWstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVWreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVWstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx [off] {sym} ptr idx (MOVWZreg x) mem)
	// result: (MOVWstoreidx [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != OpPPC64MOVWZreg {
			break
		}
		x := v_2.Args[0]
		v.reset(OpPPC64MOVWstoreidx)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MOVWstorezero_0(v *Value) bool {
	// match: (MOVWstorezero [off1] {sym} (ADDconst [off2] x) mem)
	// cond: is16Bit(off1+off2)
	// result: (MOVWstorezero [off1+off2] {sym} x mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconst {
			break
		}
		off2 := v_0.AuxInt
		x := v_0.Args[0]
		if !(is16Bit(off1 + off2)) {
			break
		}
		v.reset(OpPPC64MOVWstorezero)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstorezero [off1] {sym1} p:(MOVDaddr [off2] {sym2} x) mem)
	// cond: canMergeSym(sym1,sym2) && (x.Op != OpSB || p.Uses == 1)
	// result: (MOVWstorezero [off1+off2] {mergeSym(sym1,sym2)} x mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		p := v.Args[0]
		if p.Op != OpPPC64MOVDaddr {
			break
		}
		off2 := p.AuxInt
		sym2 := p.Aux
		x := p.Args[0]
		if !(canMergeSym(sym1, sym2) && (x.Op != OpSB || p.Uses == 1)) {
			break
		}
		v.reset(OpPPC64MOVWstorezero)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MTVSRD_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (MTVSRD (MOVDconst [c]))
	// result: (FMOVDconst [c])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v.reset(OpPPC64FMOVDconst)
		v.AuxInt = c
		return true
	}
	// match: (MTVSRD x:(MOVDload [off] {sym} ptr mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (FMOVDload [off] {sym} ptr mem)
	for {
		x := v.Args[0]
		if x.Op != OpPPC64MOVDload {
			break
		}
		off := x.AuxInt
		sym := x.Aux
		mem := x.Args[1]
		ptr := x.Args[0]
		if !(x.Uses == 1 && clobber(x)) {
			break
		}
		b = x.Block
		v0 := b.NewValue0(x.Pos, OpPPC64FMOVDload, typ.Float64)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64MaskIfNotCarry_0(v *Value) bool {
	// match: (MaskIfNotCarry (ADDconstForCarry [c] (ANDconst [d] _)))
	// cond: c < 0 && d > 0 && c + d < 0
	// result: (MOVDconst [-1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ADDconstForCarry {
			break
		}
		c := v_0.AuxInt
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpPPC64ANDconst {
			break
		}
		d := v_0_0.AuxInt
		if !(c < 0 && d > 0 && c+d < 0) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = -1
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64NotEqual_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (NotEqual (FlagEQ))
	// result: (MOVDconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagEQ {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (NotEqual (FlagLT))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagLT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (NotEqual (FlagGT))
	// result: (MOVDconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64FlagGT {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 1
		return true
	}
	// match: (NotEqual (InvertFlags x))
	// result: (NotEqual x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(OpPPC64NotEqual)
		v.AddArg(x)
		return true
	}
	// match: (NotEqual cmp)
	// result: (ISELB [6] (MOVDconst [1]) cmp)
	for {
		cmp := v.Args[0]
		v.reset(OpPPC64ISELB)
		v.AuxInt = 6
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v0.AuxInt = 1
		v.AddArg(v0)
		v.AddArg(cmp)
		return true
	}
}
func rewriteValuePPC64_OpPPC64OR_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (OR (SLDconst x [c]) (SRDconst x [d]))
	// cond: d == 64-c
	// result: (ROTLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 64-c) {
			break
		}
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (OR (SRDconst x [d]) (SLDconst x [c]))
	// cond: d == 64-c
	// result: (ROTLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLDconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 64-c) {
			break
		}
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (OR (SLWconst x [c]) (SRWconst x [d]))
	// cond: d == 32-c
	// result: (ROTLWconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (OR (SRWconst x [d]) (SLWconst x [c]))
	// cond: d == 32-c
	// result: (ROTLWconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLWconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (OR (SLD x (ANDconst <typ.Int64> [63] y)) (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y))))
	// result: (ROTL x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLD {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64ANDconst || v_0_1.Type != typ.Int64 || v_0_1.AuxInt != 63 {
			break
		}
		y := v_0_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRD {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64SUB || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 64 {
			break
		}
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64ANDconst || v_1_1_1.Type != typ.UInt || v_1_1_1.AuxInt != 63 || y != v_1_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (OR (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y))) (SLD x (ANDconst <typ.Int64> [63] y)))
	// result: (ROTL x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRD {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64SUB || v_0_1.Type != typ.UInt {
			break
		}
		_ = v_0_1.Args[1]
		v_0_1_0 := v_0_1.Args[0]
		if v_0_1_0.Op != OpPPC64MOVDconst || v_0_1_0.AuxInt != 64 {
			break
		}
		v_0_1_1 := v_0_1.Args[1]
		if v_0_1_1.Op != OpPPC64ANDconst || v_0_1_1.Type != typ.UInt || v_0_1_1.AuxInt != 63 {
			break
		}
		y := v_0_1_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLD {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.Int64 || v_1_1.AuxInt != 63 || y != v_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (OR (SLW x (ANDconst <typ.Int32> [31] y)) (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y))))
	// result: (ROTLW x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLW {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64ANDconst || v_0_1.Type != typ.Int32 || v_0_1.AuxInt != 31 {
			break
		}
		y := v_0_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRW {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64SUB || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 32 {
			break
		}
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64ANDconst || v_1_1_1.Type != typ.UInt || v_1_1_1.AuxInt != 31 || y != v_1_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (OR (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y))) (SLW x (ANDconst <typ.Int32> [31] y)))
	// result: (ROTLW x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRW {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64SUB || v_0_1.Type != typ.UInt {
			break
		}
		_ = v_0_1.Args[1]
		v_0_1_0 := v_0_1.Args[0]
		if v_0_1_0.Op != OpPPC64MOVDconst || v_0_1_0.AuxInt != 32 {
			break
		}
		v_0_1_1 := v_0_1.Args[1]
		if v_0_1_1.Op != OpPPC64ANDconst || v_0_1_1.Type != typ.UInt || v_0_1_1.AuxInt != 31 {
			break
		}
		y := v_0_1_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLW {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.Int32 || v_1_1.AuxInt != 31 || y != v_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (OR (MOVDconst [c]) (MOVDconst [d]))
	// result: (MOVDconst [c|d])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		d := v_1.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c | d
		return true
	}
	// match: (OR (MOVDconst [d]) (MOVDconst [c]))
	// result: (MOVDconst [c|d])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		d := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c | d
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_10(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR x (MOVDconst [c]))
	// cond: isU32Bit(c)
	// result: (ORconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(isU32Bit(c)) {
			break
		}
		v.reset(OpPPC64ORconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (OR (MOVDconst [c]) x)
	// cond: isU32Bit(c)
	// result: (ORconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(isU32Bit(c)) {
			break
		}
		v.reset(OpPPC64ORconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i0] {s} p mem) o1:(SLWconst x1:(MOVBZload [i1] {s} p mem) [8]))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o1 := v.Args[1]
		if o1.Op != OpPPC64SLWconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVHZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o1:(SLWconst x1:(MOVBZload [i1] {s} p mem) [8]) x0:(MOVBZload [i0] {s} p mem))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o1 := v.Args[0]
		if o1.Op != OpPPC64SLWconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVHZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i0] {s} p mem) o1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [8]))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o1 := v.Args[1]
		if o1.Op != OpPPC64SLDconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVHZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [8]) x0:(MOVBZload [i0] {s} p mem))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o1 := v.Args[0]
		if o1.Op != OpPPC64SLDconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVHZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i1] {s} p mem) o1:(SLWconst x1:(MOVBZload [i0] {s} p mem) [8]))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o1 := v.Args[1]
		if o1.Op != OpPPC64SLWconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVHBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o1:(SLWconst x1:(MOVBZload [i0] {s} p mem) [8]) x0:(MOVBZload [i1] {s} p mem))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o1 := v.Args[0]
		if o1.Op != OpPPC64SLWconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVHBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i1] {s} p mem) o1:(SLDconst x1:(MOVBZload [i0] {s} p mem) [8]))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o1 := v.Args[1]
		if o1.Op != OpPPC64SLDconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVHBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o1:(SLDconst x1:(MOVBZload [i0] {s} p mem) [8]) x0:(MOVBZload [i1] {s} p mem))
	// cond: !config.BigEndian && i1 == i0+1 && x0.Uses ==1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)
	// result: @mergePoint(b,x0,x1) (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o1 := v.Args[0]
		if o1.Op != OpPPC64SLDconst || o1.AuxInt != 8 {
			break
		}
		x1 := o1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && o1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(o1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVHBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_20(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> s0:(SLWconst x0:(MOVBZload [i1] {s} p mem) [n1]) s1:(SLWconst x1:(MOVBZload [i0] {s} p mem) [n2]))
	// cond: !config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)
	// result: @mergePoint(b,x0,x1) (SLDconst <t> (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [n1])
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLWconst {
			break
		}
		n1 := s0.AuxInt
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLWconst {
			break
		}
		n2 := s1.AuxInt
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = n1
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVHBRload, t)
		v2 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s1:(SLWconst x1:(MOVBZload [i0] {s} p mem) [n2]) s0:(SLWconst x0:(MOVBZload [i1] {s} p mem) [n1]))
	// cond: !config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)
	// result: @mergePoint(b,x0,x1) (SLDconst <t> (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [n1])
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLWconst {
			break
		}
		n2 := s1.AuxInt
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLWconst {
			break
		}
		n1 := s0.AuxInt
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = n1
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVHBRload, t)
		v2 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i1] {s} p mem) [n1]) s1:(SLDconst x1:(MOVBZload [i0] {s} p mem) [n2]))
	// cond: !config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)
	// result: @mergePoint(b,x0,x1) (SLDconst <t> (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [n1])
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst {
			break
		}
		n1 := s0.AuxInt
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLDconst {
			break
		}
		n2 := s1.AuxInt
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = n1
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVHBRload, t)
		v2 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s1:(SLDconst x1:(MOVBZload [i0] {s} p mem) [n2]) s0:(SLDconst x0:(MOVBZload [i1] {s} p mem) [n1]))
	// cond: !config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)
	// result: @mergePoint(b,x0,x1) (SLDconst <t> (MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [n1])
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLDconst {
			break
		}
		n2 := s1.AuxInt
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst {
			break
		}
		n1 := s0.AuxInt
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && n1%8 == 0 && n2 == n1+8 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0) && clobber(s1)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = n1
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVHBRload, t)
		v2 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s1:(SLWconst x2:(MOVBZload [i3] {s} p mem) [24]) o0:(OR <t> s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [16]) x0:(MOVHZload [i0] {s} p mem)))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLWconst x2:(MOVBZload [i3] {s} p mem) [24]) o0:(OR <t> x0:(MOVHZload [i0] {s} p mem) s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [16])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [16]) x0:(MOVHZload [i0] {s} p mem)) s1:(SLWconst x2:(MOVBZload [i3] {s} p mem) [24]))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> x0:(MOVHZload [i0] {s} p mem) s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [16])) s1:(SLWconst x2:(MOVBZload [i3] {s} p mem) [24]))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLDconst x2:(MOVBZload [i3] {s} p mem) [24]) o0:(OR <t> s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [16]) x0:(MOVHZload [i0] {s} p mem)))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLDconst x2:(MOVBZload [i3] {s} p mem) [24]) o0:(OR <t> x0:(MOVHZload [i0] {s} p mem) s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [16])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_30(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> o0:(OR <t> s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [16]) x0:(MOVHZload [i0] {s} p mem)) s1:(SLDconst x2:(MOVBZload [i3] {s} p mem) [24]))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> x0:(MOVHZload [i0] {s} p mem) s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [16])) s1:(SLDconst x2:(MOVBZload [i3] {s} p mem) [24]))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses ==1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWZload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWZload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLWconst x2:(MOVBZload [i0] {s} p mem) [24]) o0:(OR <t> s0:(SLWconst x1:(MOVBZload [i1] {s} p mem) [16]) x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem)))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLWconst x2:(MOVBZload [i0] {s} p mem) [24]) o0:(OR <t> x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) s0:(SLWconst x1:(MOVBZload [i1] {s} p mem) [16])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLWconst x1:(MOVBZload [i1] {s} p mem) [16]) x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem)) s1:(SLWconst x2:(MOVBZload [i0] {s} p mem) [24]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) s0:(SLWconst x1:(MOVBZload [i1] {s} p mem) [16])) s1:(SLWconst x2:(MOVBZload [i0] {s} p mem) [24]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		mem := x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		s := x0_0.Aux
		p := x0_0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLDconst x2:(MOVBZload [i0] {s} p mem) [24]) o0:(OR <t> s0:(SLDconst x1:(MOVBZload [i1] {s} p mem) [16]) x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem)))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s1:(SLDconst x2:(MOVBZload [i0] {s} p mem) [24]) o0:(OR <t> x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) s0:(SLDconst x1:(MOVBZload [i1] {s} p mem) [16])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLDconst x1:(MOVBZload [i1] {s} p mem) [16]) x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem)) s1:(SLDconst x2:(MOVBZload [i0] {s} p mem) [24]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := o0.Args[1]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) s0:(SLDconst x1:(MOVBZload [i1] {s} p mem) [16])) s1:(SLDconst x2:(MOVBZload [i0] {s} p mem) [24]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		mem := x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		s := x0_0.Aux
		p := x0_0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_40(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> x0:(MOVBZload [i3] {s} p mem) o0:(OR <t> s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [8]) s1:(SLWconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		_ = x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		if x2_0.Aux != s || p != x2_0.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i3] {s} p mem) o0:(OR <t> s1:(SLWconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16]) s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [8])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		_ = x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		if x2_0.Aux != s || p != x2_0.Args[0] || mem != x2.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [8]) s1:(SLWconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16])) x0:(MOVBZload [i3] {s} p mem))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		_ = x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		if x2_0.Aux != s || p != x2_0.Args[0] || mem != x2.Args[1] {
			break
		}
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s1:(SLWconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16]) s0:(SLWconst x1:(MOVBZload [i2] {s} p mem) [8])) x0:(MOVBZload [i3] {s} p mem))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLWconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		mem := x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		s := x2_0.Aux
		p := x2_0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i3] {s} p mem) o0:(OR <t> s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [8]) s1:(SLDconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		_ = x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		if x2_0.Aux != s || p != x2_0.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x0:(MOVBZload [i3] {s} p mem) o0:(OR <t> s1:(SLDconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16]) s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [8])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		_ = x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		if x2_0.Aux != s || p != x2_0.Args[0] || mem != x2.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [8]) s1:(SLDconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16])) x0:(MOVBZload [i3] {s} p mem))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		_ = x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		if x2_0.Aux != s || p != x2_0.Args[0] || mem != x2.Args[1] {
			break
		}
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x2:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [16]) s0:(SLDconst x1:(MOVBZload [i2] {s} p mem) [8])) x0:(MOVBZload [i3] {s} p mem))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 16 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != OpPPC64MOVHBRload || x2.Type != t {
			break
		}
		mem := x2.Args[1]
		x2_0 := x2.Args[0]
		if x2_0.Op != OpPPC64MOVDaddr || x2_0.Type != typ.Uintptr {
			break
		}
		i0 := x2_0.AuxInt
		s := x2_0.Aux
		p := x2_0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		x0 := v.Args[1]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s2:(SLDconst x2:(MOVBZload [i3] {s} p mem) [32]) o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i2] {s} p mem) [40]) s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [48])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		s2 := v.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 32 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 40 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 48 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i0 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s2:(SLDconst x2:(MOVBZload [i3] {s} p mem) [32]) o0:(OR <t> s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [48]) s1:(SLDconst x1:(MOVBZload [i2] {s} p mem) [40])))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		s2 := v.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 32 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 48 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i0 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 40 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_50(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i2] {s} p mem) [40]) s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [48])) s2:(SLDconst x2:(MOVBZload [i3] {s} p mem) [32]))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 40 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 48 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i0 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s2 := v.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 32 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [48]) s1:(SLDconst x1:(MOVBZload [i2] {s} p mem) [40])) s2:(SLDconst x2:(MOVBZload [i3] {s} p mem) [32]))
	// cond: !config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 48 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		mem := x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i0 := x0_0.AuxInt
		s := x0_0.Aux
		p := x0_0.Args[0]
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 40 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s2 := v.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 32 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s2:(SLDconst x2:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) [32])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		s2 := v.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 56 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s2:(SLDconst x2:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) [32]) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		s2 := v.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 56 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) [32])) s2:(SLDconst x2:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		s0 := o0.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		_ = x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		if x0_0.Aux != s || p != x0_0.Args[0] || mem != x0.Args[1] {
			break
		}
		s2 := v.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 56 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> o0:(OR <t> s0:(SLDconst x0:(MOVHBRload <t> (MOVDaddr <typ.Uintptr> [i2] {s} p) mem) [32]) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])) s2:(SLDconst x2:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (SLDconst <t> (MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVHBRload || x0.Type != t {
			break
		}
		mem := x0.Args[1]
		x0_0 := x0.Args[0]
		if x0_0.Op != OpPPC64MOVDaddr || x0_0.Type != typ.Uintptr {
			break
		}
		i2 := x0_0.AuxInt
		s := x0_0.Aux
		p := x0_0.Args[0]
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s2 := v.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 56 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && o0.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, OpPPC64SLDconst, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = 32
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVWBRload, t)
		v2 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v2.AuxInt = i0
		v2.Aux = s
		v2.AddArg(p)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v0.AddArg(v1)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem)))))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32])))))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x4.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem)) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]))))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x5.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32])) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]))))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x5.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_60(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem))) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]))) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> o4:(OR <t> o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem)) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40])) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]) o5:(OR <t> o4:(OR <t> o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32])) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40])) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s6 := v.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem)))) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32])))) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem)) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]))) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48]) o4:(OR <t> o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32])) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]))) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s5 := o5.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem))) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		s := x5.Aux
		mem := x5.Args[1]
		p := x5.Args[0]
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40]) o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]))) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s4 := o4.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		s := x5.Aux
		mem := x5.Args[1]
		p := x5.Args[0]
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_70(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> o5:(OR <t> o4:(OR <t> o3:(OR <t> s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32]) x0:(MOVWZload {s} [i0] p mem)) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40])) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s3 := o3.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		s := x4.Aux
		mem := x4.Args[1]
		p := x4.Args[0]
		x0 := o3.Args[1]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> o3:(OR <t> x0:(MOVWZload {s} [i0] p mem) s3:(SLDconst x4:(MOVBZload [i4] {s} p mem) [32])) s4:(SLDconst x5:(MOVBZload [i5] {s} p mem) [40])) s5:(SLDconst x6:(MOVBZload [i6] {s} p mem) [48])) s6:(SLDconst x7:(MOVBZload [i7] {s} p mem) [56]))
	// cond: !config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses ==1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber (s6) && clobber(o3) && clobber(o4) && clobber(o5)
	// result: @mergePoint(b,x0,x4,x5,x6,x7) (MOVDload <t> {s} [i0] p mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		x0 := o3.Args[0]
		if x0.Op != OpPPC64MOVWZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s3 := o3.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x4 := s3.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s4 := o4.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 40 {
			break
		}
		x5 := s4.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s5 := o5.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 48 {
			break
		}
		x6 := s5.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		s6 := v.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 56 {
			break
		}
		x7 := s6.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i0%4 == 0 && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x0.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s3.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x0, x4, x5, x6, x7) != nil && clobber(x0) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(s3) && clobber(s4) && clobber(s5) && clobber(s6) && clobber(o3) && clobber(o4) && clobber(o5)) {
			break
		}
		b = mergePoint(b, x0, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem)))))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x4.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x4.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32])))))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x3.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x3.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem)) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]))))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32])) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]))))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x2.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x2.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem))) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]))) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> o1:(OR <t> o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem)) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40])) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]) o0:(OR <t> o1:(OR <t> o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32])) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40])) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		o0 := v.Args[1]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x1.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x1.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_80(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem)))) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32])))) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem)) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]))) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48]) o1:(OR <t> o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32])) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]))) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		s1 := o0.Args[0]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		o1 := o0.Args[1]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem))) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> o1:(OR <t> s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40]) o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]))) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		s2 := o1.Args[0]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o2 := o1.Args[1]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		if x3.Aux != s {
			break
		}
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> o1:(OR <t> o2:(OR <t> s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32]) x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem)) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40])) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		s3 := o2.Args[0]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		s := x3.Aux
		mem := x3.Args[1]
		p := x3.Args[0]
		x4 := o2.Args[1]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		_ = x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		if p != x4_0.Args[0] || mem != x4.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o0:(OR <t> o1:(OR <t> o2:(OR <t> x4:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i4] p) mem) s3:(SLDconst x3:(MOVBZload [i3] {s} p mem) [32])) s2:(SLDconst x2:(MOVBZload [i2] {s} p mem) [40])) s1:(SLDconst x1:(MOVBZload [i1] {s} p mem) [48])) s0:(SLDconst x0:(MOVBZload [i0] {s} p mem) [56]))
	// cond: !config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)
	// result: @mergePoint(b,x0,x1,x2,x3,x4) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != OpPPC64OR || o0.Type != t {
			break
		}
		_ = o0.Args[1]
		o1 := o0.Args[0]
		if o1.Op != OpPPC64OR || o1.Type != t {
			break
		}
		_ = o1.Args[1]
		o2 := o1.Args[0]
		if o2.Op != OpPPC64OR || o2.Type != t {
			break
		}
		_ = o2.Args[1]
		x4 := o2.Args[0]
		if x4.Op != OpPPC64MOVWBRload || x4.Type != t {
			break
		}
		mem := x4.Args[1]
		x4_0 := x4.Args[0]
		if x4_0.Op != OpPPC64MOVDaddr || x4_0.Type != typ.Uintptr {
			break
		}
		i4 := x4_0.AuxInt
		p := x4_0.Args[0]
		s3 := o2.Args[1]
		if s3.Op != OpPPC64SLDconst || s3.AuxInt != 32 {
			break
		}
		x3 := s3.Args[0]
		if x3.Op != OpPPC64MOVBZload {
			break
		}
		i3 := x3.AuxInt
		s := x3.Aux
		_ = x3.Args[1]
		if p != x3.Args[0] || mem != x3.Args[1] {
			break
		}
		s2 := o1.Args[1]
		if s2.Op != OpPPC64SLDconst || s2.AuxInt != 40 {
			break
		}
		x2 := s2.Args[0]
		if x2.Op != OpPPC64MOVBZload {
			break
		}
		i2 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] {
			break
		}
		s1 := o0.Args[1]
		if s1.Op != OpPPC64SLDconst || s1.AuxInt != 48 {
			break
		}
		x1 := s1.Args[0]
		if x1.Op != OpPPC64MOVBZload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] {
			break
		}
		s0 := v.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 56 {
			break
		}
		x0 := s0.Args[0]
		if x0.Op != OpPPC64MOVBZload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(!config.BigEndian && i1 == i0+1 && i2 == i0+2 && i3 == i0+3 && i4 == i0+4 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && x3.Uses == 1 && x4.Uses == 1 && o0.Uses == 1 && o1.Uses == 1 && o2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && s2.Uses == 1 && s3.Uses == 1 && mergePoint(b, x0, x1, x2, x3, x4) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(x3) && clobber(x4) && clobber(o0) && clobber(o1) && clobber(o2) && clobber(s0) && clobber(s1) && clobber(s2) && clobber(s3)) {
			break
		}
		b = mergePoint(b, x0, x1, x2, x3, x4)
		v0 := b.NewValue0(x0.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x0.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x3.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x3.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x4.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x4.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_90(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x5.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x5.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x5.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x5.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_100(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		s := x5.Aux
		mem := x5.Args[1]
		p := x5.Args[0]
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		s := x5.Aux
		mem := x5.Args[1]
		p := x5.Args[0]
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		s := x4.Aux
		mem := x4.Args[1]
		p := x4.Args[0]
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> o3:(OR <t> s0:(SLWconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLWconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		mem := x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		s := x3_0.Aux
		p := x3_0.Args[0]
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x3.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x3.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x4.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x4.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x5.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x5.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x5.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x5.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64OR_110(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> x7:(MOVBZload [i7] {s} p mem) o5:(OR <t> o4:(OR <t> o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		x7 := v.Args[0]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		s := x7.Aux
		mem := x7.Args[1]
		p := x7.Args[0]
		o5 := v.Args[1]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x6.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x6.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8]) o4:(OR <t> o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]))) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		s6 := o5.Args[0]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		s := x6.Aux
		mem := x6.Args[1]
		p := x6.Args[0]
		o4 := o5.Args[1]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		s := x5.Aux
		mem := x5.Args[1]
		p := x5.Args[0]
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16]) o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]))) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		s5 := o4.Args[0]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		s := x5.Aux
		mem := x5.Args[1]
		p := x5.Args[0]
		o3 := o4.Args[1]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> o3:(OR <t> s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24]) s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s4 := o3.Args[0]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		s := x4.Aux
		mem := x4.Args[1]
		p := x4.Args[0]
		s0 := o3.Args[1]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		_ = x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		if x3_0.Aux != s || p != x3_0.Args[0] || mem != x3.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	// match: (OR <t> o5:(OR <t> o4:(OR <t> o3:(OR <t> s0:(SLDconst x3:(MOVWBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem) [32]) s4:(SLDconst x4:(MOVBZload [i4] {s} p mem) [24])) s5:(SLDconst x5:(MOVBZload [i5] {s} p mem) [16])) s6:(SLDconst x6:(MOVBZload [i6] {s} p mem) [8])) x7:(MOVBZload [i7] {s} p mem))
	// cond: !config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)
	// result: @mergePoint(b,x3,x4,x5,x6,x7) (MOVDBRload <t> (MOVDaddr <typ.Uintptr> [i0] {s} p) mem)
	for {
		t := v.Type
		_ = v.Args[1]
		o5 := v.Args[0]
		if o5.Op != OpPPC64OR || o5.Type != t {
			break
		}
		_ = o5.Args[1]
		o4 := o5.Args[0]
		if o4.Op != OpPPC64OR || o4.Type != t {
			break
		}
		_ = o4.Args[1]
		o3 := o4.Args[0]
		if o3.Op != OpPPC64OR || o3.Type != t {
			break
		}
		_ = o3.Args[1]
		s0 := o3.Args[0]
		if s0.Op != OpPPC64SLDconst || s0.AuxInt != 32 {
			break
		}
		x3 := s0.Args[0]
		if x3.Op != OpPPC64MOVWBRload || x3.Type != t {
			break
		}
		mem := x3.Args[1]
		x3_0 := x3.Args[0]
		if x3_0.Op != OpPPC64MOVDaddr || x3_0.Type != typ.Uintptr {
			break
		}
		i0 := x3_0.AuxInt
		s := x3_0.Aux
		p := x3_0.Args[0]
		s4 := o3.Args[1]
		if s4.Op != OpPPC64SLDconst || s4.AuxInt != 24 {
			break
		}
		x4 := s4.Args[0]
		if x4.Op != OpPPC64MOVBZload {
			break
		}
		i4 := x4.AuxInt
		if x4.Aux != s {
			break
		}
		_ = x4.Args[1]
		if p != x4.Args[0] || mem != x4.Args[1] {
			break
		}
		s5 := o4.Args[1]
		if s5.Op != OpPPC64SLDconst || s5.AuxInt != 16 {
			break
		}
		x5 := s5.Args[0]
		if x5.Op != OpPPC64MOVBZload {
			break
		}
		i5 := x5.AuxInt
		if x5.Aux != s {
			break
		}
		_ = x5.Args[1]
		if p != x5.Args[0] || mem != x5.Args[1] {
			break
		}
		s6 := o5.Args[1]
		if s6.Op != OpPPC64SLDconst || s6.AuxInt != 8 {
			break
		}
		x6 := s6.Args[0]
		if x6.Op != OpPPC64MOVBZload {
			break
		}
		i6 := x6.AuxInt
		if x6.Aux != s {
			break
		}
		_ = x6.Args[1]
		if p != x6.Args[0] || mem != x6.Args[1] {
			break
		}
		x7 := v.Args[1]
		if x7.Op != OpPPC64MOVBZload {
			break
		}
		i7 := x7.AuxInt
		if x7.Aux != s {
			break
		}
		_ = x7.Args[1]
		if p != x7.Args[0] || mem != x7.Args[1] || !(!config.BigEndian && i4 == i0+4 && i5 == i0+5 && i6 == i0+6 && i7 == i0+7 && x3.Uses == 1 && x4.Uses == 1 && x5.Uses == 1 && x6.Uses == 1 && x7.Uses == 1 && o3.Uses == 1 && o4.Uses == 1 && o5.Uses == 1 && s0.Uses == 1 && s4.Uses == 1 && s5.Uses == 1 && s6.Uses == 1 && mergePoint(b, x3, x4, x5, x6, x7) != nil && clobber(x3) && clobber(x4) && clobber(x5) && clobber(x6) && clobber(x7) && clobber(o3) && clobber(o4) && clobber(o5) && clobber(s0) && clobber(s4) && clobber(s5) && clobber(s6)) {
			break
		}
		b = mergePoint(b, x3, x4, x5, x6, x7)
		v0 := b.NewValue0(x7.Pos, OpPPC64MOVDBRload, t)
		v.reset(OpCopy)
		v.AddArg(v0)
		v1 := b.NewValue0(x7.Pos, OpPPC64MOVDaddr, typ.Uintptr)
		v1.AuxInt = i0
		v1.Aux = s
		v1.AddArg(p)
		v0.AddArg(v1)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ORN_0(v *Value) bool {
	// match: (ORN x (MOVDconst [-1]))
	// result: x
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst || v_1.AuxInt != -1 {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ORconst_0(v *Value) bool {
	// match: (ORconst [c] (ORconst [d] x))
	// result: (ORconst [c|d] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64ORconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(OpPPC64ORconst)
		v.AuxInt = c | d
		v.AddArg(x)
		return true
	}
	// match: (ORconst [-1] _)
	// result: (MOVDconst [-1])
	for {
		if v.AuxInt != -1 {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = -1
		return true
	}
	// match: (ORconst [0] x)
	// result: x
	for {
		if v.AuxInt != 0 {
			break
		}
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ROTL_0(v *Value) bool {
	// match: (ROTL x (MOVDconst [c]))
	// result: (ROTLconst x [c&63])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c & 63
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64ROTLW_0(v *Value) bool {
	// match: (ROTLW x (MOVDconst [c]))
	// result: (ROTLWconst x [c&31])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64SUB_0(v *Value) bool {
	// match: (SUB x (MOVDconst [c]))
	// cond: is32Bit(-c)
	// result: (ADDconst [-c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(is32Bit(-c)) {
			break
		}
		v.reset(OpPPC64ADDconst)
		v.AuxInt = -c
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64XOR_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (XOR (SLDconst x [c]) (SRDconst x [d]))
	// cond: d == 64-c
	// result: (ROTLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLDconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRDconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 64-c) {
			break
		}
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XOR (SRDconst x [d]) (SLDconst x [c]))
	// cond: d == 64-c
	// result: (ROTLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRDconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLDconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 64-c) {
			break
		}
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XOR (SLWconst x [c]) (SRWconst x [d]))
	// cond: d == 32-c
	// result: (ROTLWconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLWconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRWconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XOR (SRWconst x [d]) (SLWconst x [c]))
	// cond: d == 32-c
	// result: (ROTLWconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLWconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XOR (SLD x (ANDconst <typ.Int64> [63] y)) (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y))))
	// result: (ROTL x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLD {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64ANDconst || v_0_1.Type != typ.Int64 || v_0_1.AuxInt != 63 {
			break
		}
		y := v_0_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRD {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64SUB || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 64 {
			break
		}
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64ANDconst || v_1_1_1.Type != typ.UInt || v_1_1_1.AuxInt != 63 || y != v_1_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (XOR (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y))) (SLD x (ANDconst <typ.Int64> [63] y)))
	// result: (ROTL x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRD {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64SUB || v_0_1.Type != typ.UInt {
			break
		}
		_ = v_0_1.Args[1]
		v_0_1_0 := v_0_1.Args[0]
		if v_0_1_0.Op != OpPPC64MOVDconst || v_0_1_0.AuxInt != 64 {
			break
		}
		v_0_1_1 := v_0_1.Args[1]
		if v_0_1_1.Op != OpPPC64ANDconst || v_0_1_1.Type != typ.UInt || v_0_1_1.AuxInt != 63 {
			break
		}
		y := v_0_1_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLD {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.Int64 || v_1_1.AuxInt != 63 || y != v_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (XOR (SLW x (ANDconst <typ.Int32> [31] y)) (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y))))
	// result: (ROTLW x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SLW {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64ANDconst || v_0_1.Type != typ.Int32 || v_0_1.AuxInt != 31 {
			break
		}
		y := v_0_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SRW {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64SUB || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 32 {
			break
		}
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64ANDconst || v_1_1_1.Type != typ.UInt || v_1_1_1.AuxInt != 31 || y != v_1_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (XOR (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y))) (SLW x (ANDconst <typ.Int32> [31] y)))
	// result: (ROTLW x y)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64SRW {
			break
		}
		_ = v_0.Args[1]
		x := v_0.Args[0]
		v_0_1 := v_0.Args[1]
		if v_0_1.Op != OpPPC64SUB || v_0_1.Type != typ.UInt {
			break
		}
		_ = v_0_1.Args[1]
		v_0_1_0 := v_0_1.Args[0]
		if v_0_1_0.Op != OpPPC64MOVDconst || v_0_1_0.AuxInt != 32 {
			break
		}
		v_0_1_1 := v_0_1.Args[1]
		if v_0_1_1.Op != OpPPC64ANDconst || v_0_1_1.Type != typ.UInt || v_0_1_1.AuxInt != 31 {
			break
		}
		y := v_0_1_1.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SLW {
			break
		}
		_ = v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.Int32 || v_1_1.AuxInt != 31 || y != v_1_1.Args[0] {
			break
		}
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (XOR (MOVDconst [c]) (MOVDconst [d]))
	// result: (MOVDconst [c^d])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		d := v_1.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c ^ d
		return true
	}
	// match: (XOR (MOVDconst [d]) (MOVDconst [c]))
	// result: (MOVDconst [c^d])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		d := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = c ^ d
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64XOR_10(v *Value) bool {
	// match: (XOR x (MOVDconst [c]))
	// cond: isU32Bit(c)
	// result: (XORconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(isU32Bit(c)) {
			break
		}
		v.reset(OpPPC64XORconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XOR (MOVDconst [c]) x)
	// cond: isU32Bit(c)
	// result: (XORconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64MOVDconst {
			break
		}
		c := v_0.AuxInt
		if !(isU32Bit(c)) {
			break
		}
		v.reset(OpPPC64XORconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPPC64XORconst_0(v *Value) bool {
	// match: (XORconst [c] (XORconst [d] x))
	// result: (XORconst [c^d] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpPPC64XORconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(OpPPC64XORconst)
		v.AuxInt = c ^ d
		v.AddArg(x)
		return true
	}
	// match: (XORconst [0] x)
	// result: x
	for {
		if v.AuxInt != 0 {
			break
		}
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPanicBounds_0(v *Value) bool {
	// match: (PanicBounds [kind] x y mem)
	// cond: boundsABI(kind) == 0
	// result: (LoweredPanicBoundsA [kind] x y mem)
	for {
		kind := v.AuxInt
		mem := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		if !(boundsABI(kind) == 0) {
			break
		}
		v.reset(OpPPC64LoweredPanicBoundsA)
		v.AuxInt = kind
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	// match: (PanicBounds [kind] x y mem)
	// cond: boundsABI(kind) == 1
	// result: (LoweredPanicBoundsB [kind] x y mem)
	for {
		kind := v.AuxInt
		mem := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		if !(boundsABI(kind) == 1) {
			break
		}
		v.reset(OpPPC64LoweredPanicBoundsB)
		v.AuxInt = kind
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	// match: (PanicBounds [kind] x y mem)
	// cond: boundsABI(kind) == 2
	// result: (LoweredPanicBoundsC [kind] x y mem)
	for {
		kind := v.AuxInt
		mem := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		if !(boundsABI(kind) == 2) {
			break
		}
		v.reset(OpPPC64LoweredPanicBoundsC)
		v.AuxInt = kind
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpPopCount16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (PopCount16 x)
	// result: (POPCNTW (MOVHZreg x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64POPCNTW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpPopCount32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (PopCount32 x)
	// result: (POPCNTW (MOVWZreg x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64POPCNTW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpPopCount64_0(v *Value) bool {
	// match: (PopCount64 x)
	// result: (POPCNTD x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64POPCNTD)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpPopCount8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (PopCount8 x)
	// result: (POPCNTB (MOVBZreg x))
	for {
		x := v.Args[0]
		v.reset(OpPPC64POPCNTB)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRotateLeft16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (RotateLeft16 <t> x (MOVDconst [c]))
	// result: (Or16 (Lsh16x64 <t> x (MOVDconst [c&15])) (Rsh16Ux64 <t> x (MOVDconst [-c&15])))
	for {
		t := v.Type
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpOr16)
		v0 := b.NewValue0(v.Pos, OpLsh16x64, t)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = c & 15
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpRsh16Ux64, t)
		v2.AddArg(x)
		v3 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v3.AuxInt = -c & 15
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
	return false
}
func rewriteValuePPC64_OpRotateLeft32_0(v *Value) bool {
	// match: (RotateLeft32 x (MOVDconst [c]))
	// result: (ROTLWconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64ROTLWconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	// match: (RotateLeft32 x y)
	// result: (ROTLW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ROTLW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpRotateLeft64_0(v *Value) bool {
	// match: (RotateLeft64 x (MOVDconst [c]))
	// result: (ROTLconst [c&63] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpPPC64ROTLconst)
		v.AuxInt = c & 63
		v.AddArg(x)
		return true
	}
	// match: (RotateLeft64 x y)
	// result: (ROTL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64ROTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpRotateLeft8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (RotateLeft8 <t> x (MOVDconst [c]))
	// result: (Or8 (Lsh8x64 <t> x (MOVDconst [c&7])) (Rsh8Ux64 <t> x (MOVDconst [-c&7])))
	for {
		t := v.Type
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		v.reset(OpOr8)
		v0 := b.NewValue0(v.Pos, OpLsh8x64, t)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = c & 7
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpRsh8Ux64, t)
		v2.AddArg(x)
		v3 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v3.AuxInt = -c & 7
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
	return false
}
func rewriteValuePPC64_OpRound_0(v *Value) bool {
	// match: (Round x)
	// result: (FROUND x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FROUND)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpRound32F_0(v *Value) bool {
	// match: (Round32F x)
	// result: (LoweredRound32F x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64LoweredRound32F)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpRound64F_0(v *Value) bool {
	// match: (Round64F x)
	// result: (LoweredRound64F x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64LoweredRound64F)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpRsh16Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVHZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16Ux16 x y)
	// result: (SRW (ZeroExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v4 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux32 x (Const64 [c]))
	// cond: uint32(c) < 16
	// result: (SRWconst (ZeroExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16Ux32 x (MOVDconst [c]))
	// cond: uint32(c) < 16
	// result: (SRWconst (ZeroExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16Ux32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVHZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16Ux32 x y)
	// result: (SRW (ZeroExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v4 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux64 x (Const64 [c]))
	// cond: uint64(c) < 16
	// result: (SRWconst (ZeroExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 16
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 16) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh16Ux64 x (MOVDconst [c]))
	// cond: uint64(c) < 16
	// result: (SRWconst (ZeroExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 16) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16Ux64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVHZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16Ux64 x y)
	// result: (SRW (ZeroExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVHZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16Ux8 x y)
	// result: (SRW (ZeroExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v4 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVHreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16x16 x y)
	// result: (SRAW (SignExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v4 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x32 x (Const64 [c]))
	// cond: uint32(c) < 16
	// result: (SRAWconst (SignExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16x32 x (MOVDconst [c]))
	// cond: uint32(c) < 16
	// result: (SRAWconst (SignExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVHreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16x32 x y)
	// result: (SRAW (SignExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v4 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x64 x (Const64 [c]))
	// cond: uint64(c) < 16
	// result: (SRAWconst (SignExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16x64 x (Const64 [c]))
	// cond: uint64(c) >= 16
	// result: (SRAWconst (SignExt16to32 x) [63])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = 63
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16x64 x (MOVDconst [c]))
	// cond: uint64(c) < 16
	// result: (SRAWconst (SignExt16to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 16) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVHreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16x64 x y)
	// result: (SRAW (SignExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh16x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVHreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh16x8 x y)
	// result: (SRAW (SignExt16to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-16] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt16to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -16
		v4 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh32Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32Ux16 x y)
	// result: (SRW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux32 x (Const64 [c]))
	// cond: uint32(c) < 32
	// result: (SRWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 32) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32Ux32 x (MOVDconst [c]))
	// cond: uint32(c) < 32
	// result: (SRWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 32) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32Ux32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32Ux32 x y)
	// result: (SRW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux64 x (Const64 [c]))
	// cond: uint64(c) < 32
	// result: (SRWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 32) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 32
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 32) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh32Ux64 x (MOVDconst [c]))
	// cond: uint64(c) < 32
	// result: (SRWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 32) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32Ux64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32Ux64 x (AND y (MOVDconst [31])))
	// result: (SRW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64MOVDconst || v_1_1.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32Ux64 x (AND (MOVDconst [31]) y))
	// result: (SRW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		y := v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32Ux64 x (ANDconst <typ.UInt> [31] y))
	// result: (SRW x (ANDconst <typ.UInt> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64ANDconst || v_1.Type != typ.UInt || v_1.AuxInt != 31 {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32Ux64 x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	// result: (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 32 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.UInt || v_1_1.AuxInt != 31 {
			break
		}
		y := v_1_1.Args[0]
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 31
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32Ux64 x (SUB <typ.UInt> (MOVDconst [32]) (AND <typ.UInt> y (MOVDconst [31]))))
	// result: (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 32 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		y := v_1_1.Args[0]
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64MOVDconst || v_1_1_1.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 31
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32Ux64 x (SUB <typ.UInt> (MOVDconst [32]) (AND <typ.UInt> (MOVDconst [31]) y)))
	// result: (SRW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 32 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		y := v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 31
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpRsh32Ux64_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux64 x y)
	// result: (SRW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32Ux8 x y)
	// result: (SRW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32x16 x y)
	// result: (SRAW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x32 x (Const64 [c]))
	// cond: uint32(c) < 32
	// result: (SRAWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 32) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x32 x (MOVDconst [c]))
	// cond: uint32(c) < 32
	// result: (SRAWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 32) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32x32 x y)
	// result: (SRAW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x64 x (Const64 [c]))
	// cond: uint64(c) < 32
	// result: (SRAWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 32) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x64 x (Const64 [c]))
	// cond: uint64(c) >= 32
	// result: (SRAWconst x [63])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 32) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = 63
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x64 x (MOVDconst [c]))
	// cond: uint64(c) < 32
	// result: (SRAWconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 32) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32x64 x (AND y (MOVDconst [31])))
	// result: (SRAW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64MOVDconst || v_1_1.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32x64 x (AND (MOVDconst [31]) y))
	// result: (SRAW x (ANDconst <typ.Int32> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		y := v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int32)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32x64 x (ANDconst <typ.UInt> [31] y))
	// result: (SRAW x (ANDconst <typ.UInt> [31] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64ANDconst || v_1.Type != typ.UInt || v_1.AuxInt != 31 {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v0.AuxInt = 31
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32x64 x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	// result: (SRAW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 32 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.UInt || v_1_1.AuxInt != 31 {
			break
		}
		y := v_1_1.Args[0]
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 31
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32x64 x (SUB <typ.UInt> (MOVDconst [32]) (AND <typ.UInt> y (MOVDconst [31]))))
	// result: (SRAW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 32 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		y := v_1_1.Args[0]
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64MOVDconst || v_1_1_1.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 31
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh32x64 x (SUB <typ.UInt> (MOVDconst [32]) (AND <typ.UInt> (MOVDconst [31]) y)))
	// result: (SRAW x (SUB <typ.UInt> (MOVDconst [32]) (ANDconst <typ.UInt> [31] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 32 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		y := v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 31 {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 31
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpRsh32x64_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x64 x y)
	// result: (SRAW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh32x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh32x8 x y)
	// result: (SRAW x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-32] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -32
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64Ux16 x y)
	// result: (SRD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux32 x (Const64 [c]))
	// cond: uint32(c) < 64
	// result: (SRDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 64) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64Ux32 x (MOVDconst [c]))
	// cond: uint32(c) < 64
	// result: (SRDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 64) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64Ux32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64Ux32 x y)
	// result: (SRD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux64 x (Const64 [c]))
	// cond: uint64(c) < 64
	// result: (SRDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 64
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 64) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh64Ux64 x (MOVDconst [c]))
	// cond: uint64(c) < 64
	// result: (SRDconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpPPC64SRDconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64Ux64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64Ux64 x (AND y (MOVDconst [63])))
	// result: (SRD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64MOVDconst || v_1_1.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64Ux64 x (AND (MOVDconst [63]) y))
	// result: (SRD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		y := v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64Ux64 x (ANDconst <typ.UInt> [63] y))
	// result: (SRD x (ANDconst <typ.UInt> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64ANDconst || v_1.Type != typ.UInt || v_1.AuxInt != 63 {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64Ux64 x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	// result: (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 64 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.UInt || v_1_1.AuxInt != 63 {
			break
		}
		y := v_1_1.Args[0]
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 64
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 63
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64Ux64 x (SUB <typ.UInt> (MOVDconst [64]) (AND <typ.UInt> y (MOVDconst [63]))))
	// result: (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 64 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		y := v_1_1.Args[0]
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64MOVDconst || v_1_1_1.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 64
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 63
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64Ux64 x (SUB <typ.UInt> (MOVDconst [64]) (AND <typ.UInt> (MOVDconst [63]) y)))
	// result: (SRD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 64 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		y := v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 64
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 63
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpRsh64Ux64_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux64 x y)
	// result: (SRD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64Ux8 x y)
	// result: (SRD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64x16 x y)
	// result: (SRAD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x32 x (Const64 [c]))
	// cond: uint32(c) < 64
	// result: (SRADconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 64) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64x32 x (MOVDconst [c]))
	// cond: uint32(c) < 64
	// result: (SRADconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 64) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64x32 x y)
	// result: (SRAD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x64 x (Const64 [c]))
	// cond: uint64(c) < 64
	// result: (SRADconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64x64 x (Const64 [c]))
	// cond: uint64(c) >= 64
	// result: (SRADconst x [63])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 64) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = 63
		v.AddArg(x)
		return true
	}
	// match: (Rsh64x64 x (MOVDconst [c]))
	// cond: uint64(c) < 64
	// result: (SRADconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpPPC64SRADconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh64x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64x64 x (AND y (MOVDconst [63])))
	// result: (SRAD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64MOVDconst || v_1_1.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x (AND (MOVDconst [63]) y))
	// result: (SRAD x (ANDconst <typ.Int64> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64AND {
			break
		}
		y := v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.Int64)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x (ANDconst <typ.UInt> [63] y))
	// result: (SRAD x (ANDconst <typ.UInt> [63] y))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64ANDconst || v_1.Type != typ.UInt || v_1.AuxInt != 63 {
			break
		}
		y := v_1.Args[0]
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v0.AuxInt = 63
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	// result: (SRAD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 64 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64ANDconst || v_1_1.Type != typ.UInt || v_1_1.AuxInt != 63 {
			break
		}
		y := v_1_1.Args[0]
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 64
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 63
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x (SUB <typ.UInt> (MOVDconst [64]) (AND <typ.UInt> y (MOVDconst [63]))))
	// result: (SRAD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 64 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		_ = v_1_1.Args[1]
		y := v_1_1.Args[0]
		v_1_1_1 := v_1_1.Args[1]
		if v_1_1_1.Op != OpPPC64MOVDconst || v_1_1_1.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 64
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 63
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x (SUB <typ.UInt> (MOVDconst [64]) (AND <typ.UInt> (MOVDconst [63]) y)))
	// result: (SRAD x (SUB <typ.UInt> (MOVDconst [64]) (ANDconst <typ.UInt> [63] y)))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64SUB || v_1.Type != typ.UInt {
			break
		}
		_ = v_1.Args[1]
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpPPC64MOVDconst || v_1_0.AuxInt != 64 {
			break
		}
		v_1_1 := v_1.Args[1]
		if v_1_1.Op != OpPPC64AND || v_1_1.Type != typ.UInt {
			break
		}
		y := v_1_1.Args[1]
		v_1_1_0 := v_1_1.Args[0]
		if v_1_1_0.Op != OpPPC64MOVDconst || v_1_1_0.AuxInt != 63 {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64SUB, typ.UInt)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDconst, typ.Int64)
		v1.AuxInt = 64
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpPPC64ANDconst, typ.UInt)
		v2.AuxInt = 63
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpRsh64x64_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x64 x y)
	// result: (SRAD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v2.AddArg(y)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh64x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64x8 x y)
	// result: (SRAD x (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-64] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAD)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v2 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v2.AuxInt = -64
		v3 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpRsh8Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVBZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8Ux16 x y)
	// result: (SRW (ZeroExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v4 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux32 x (Const64 [c]))
	// cond: uint32(c) < 8
	// result: (SRWconst (ZeroExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 8) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8Ux32 x (MOVDconst [c]))
	// cond: uint32(c) < 8
	// result: (SRWconst (ZeroExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 8) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8Ux32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVBZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8Ux32 x y)
	// result: (SRW (ZeroExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v4 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux64 x (Const64 [c]))
	// cond: uint64(c) < 8
	// result: (SRWconst (ZeroExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 8) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 8
	// result: (MOVDconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 8) {
			break
		}
		v.reset(OpPPC64MOVDconst)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh8Ux64 x (MOVDconst [c]))
	// cond: uint64(c) < 8
	// result: (SRWconst (ZeroExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 8) {
			break
		}
		v.reset(OpPPC64SRWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8Ux64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVBZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8Ux64 x y)
	// result: (SRW (ZeroExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRW (MOVBZreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBZreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8Ux8 x y)
	// result: (SRW (ZeroExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRW)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to32, typ.UInt32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v4 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x16 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVBreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8x16 x y)
	// result: (SRAW (SignExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt16to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v4 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x32 x (Const64 [c]))
	// cond: uint32(c) < 8
	// result: (SRAWconst (SignExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 8) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8x32 x (MOVDconst [c]))
	// cond: uint32(c) < 8
	// result: (SRAWconst (SignExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint32(c) < 8) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8x32 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVBreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8x32 x y)
	// result: (SRAW (SignExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt32to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v4 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x64 x (Const64 [c]))
	// cond: uint64(c) < 8
	// result: (SRAWconst (SignExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 8) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8x64 x (Const64 [c]))
	// cond: uint64(c) >= 8
	// result: (SRAWconst (SignExt8to32 x) [63])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpConst64 {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 8) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = 63
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8x64 x (MOVDconst [c]))
	// cond: uint64(c) < 8
	// result: (SRAWconst (SignExt8to32 x) [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpPPC64MOVDconst {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 8) {
			break
		}
		v.reset(OpPPC64SRAWconst)
		v.AuxInt = c
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8x64 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVBreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8x64 x y)
	// result: (SRAW (SignExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] y))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpRsh8x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x8 x y)
	// cond: shiftIsBounded(v)
	// result: (SRAW (MOVBreg x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVBreg, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
	// match: (Rsh8x8 x y)
	// result: (SRAW (SignExt8to32 x) (ORN y <typ.Int64> (MaskIfNotCarry (ADDconstForCarry [-8] (ZeroExt8to64 y)))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SRAW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to32, typ.Int32)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpPPC64ORN, typ.Int64)
		v1.AddArg(y)
		v2 := b.NewValue0(v.Pos, OpPPC64MaskIfNotCarry, typ.Int64)
		v3 := b.NewValue0(v.Pos, OpPPC64ADDconstForCarry, types.TypeFlags)
		v3.AuxInt = -8
		v4 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v4.AddArg(y)
		v3.AddArg(v4)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValuePPC64_OpSignExt16to32_0(v *Value) bool {
	// match: (SignExt16to32 x)
	// result: (MOVHreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpSignExt16to64_0(v *Value) bool {
	// match: (SignExt16to64 x)
	// result: (MOVHreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpSignExt32to64_0(v *Value) bool {
	// match: (SignExt32to64 x)
	// result: (MOVWreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVWreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpSignExt8to16_0(v *Value) bool {
	// match: (SignExt8to16 x)
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpSignExt8to32_0(v *Value) bool {
	// match: (SignExt8to32 x)
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpSignExt8to64_0(v *Value) bool {
	// match: (SignExt8to64 x)
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpSlicemask_0(v *Value) bool {
	b := v.Block
	// match: (Slicemask <t> x)
	// result: (SRADconst (NEG <t> x) [63])
	for {
		t := v.Type
		x := v.Args[0]
		v.reset(OpPPC64SRADconst)
		v.AuxInt = 63
		v0 := b.NewValue0(v.Pos, OpPPC64NEG, t)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuePPC64_OpSqrt_0(v *Value) bool {
	// match: (Sqrt x)
	// result: (FSQRT x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FSQRT)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpStaticCall_0(v *Value) bool {
	// match: (StaticCall [argwid] {target} mem)
	// result: (CALLstatic [argwid] {target} mem)
	for {
		argwid := v.AuxInt
		target := v.Aux
		mem := v.Args[0]
		v.reset(OpPPC64CALLstatic)
		v.AuxInt = argwid
		v.Aux = target
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpStore_0(v *Value) bool {
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 8 && is64BitFloat(val.Type)
	// result: (FMOVDstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 8 && is64BitFloat(val.Type)) {
			break
		}
		v.reset(OpPPC64FMOVDstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 8 && is32BitFloat(val.Type)
	// result: (FMOVDstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 8 && is32BitFloat(val.Type)) {
			break
		}
		v.reset(OpPPC64FMOVDstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 4 && is32BitFloat(val.Type)
	// result: (FMOVSstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 4 && is32BitFloat(val.Type)) {
			break
		}
		v.reset(OpPPC64FMOVSstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 8 && (is64BitInt(val.Type) || isPtr(val.Type))
	// result: (MOVDstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 8 && (is64BitInt(val.Type) || isPtr(val.Type))) {
			break
		}
		v.reset(OpPPC64MOVDstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 4 && is32BitInt(val.Type)
	// result: (MOVWstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 4 && is32BitInt(val.Type)) {
			break
		}
		v.reset(OpPPC64MOVWstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 2
	// result: (MOVHstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 2) {
			break
		}
		v.reset(OpPPC64MOVHstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 1
	// result: (MOVBstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 1) {
			break
		}
		v.reset(OpPPC64MOVBstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValuePPC64_OpSub16_0(v *Value) bool {
	// match: (Sub16 x y)
	// result: (SUB x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpSub32_0(v *Value) bool {
	// match: (Sub32 x y)
	// result: (SUB x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpSub32F_0(v *Value) bool {
	// match: (Sub32F x y)
	// result: (FSUBS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FSUBS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpSub64_0(v *Value) bool {
	// match: (Sub64 x y)
	// result: (SUB x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpSub64F_0(v *Value) bool {
	// match: (Sub64F x y)
	// result: (FSUB x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64FSUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpSub8_0(v *Value) bool {
	// match: (Sub8 x y)
	// result: (SUB x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpSubPtr_0(v *Value) bool {
	// match: (SubPtr x y)
	// result: (SUB x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpTrunc_0(v *Value) bool {
	// match: (Trunc x)
	// result: (FTRUNC x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64FTRUNC)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpTrunc16to8_0(v *Value) bool {
	// match: (Trunc16to8 <t> x)
	// cond: isSigned(t)
	// result: (MOVBreg x)
	for {
		t := v.Type
		x := v.Args[0]
		if !(isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
	// match: (Trunc16to8 x)
	// result: (MOVBZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpTrunc32to16_0(v *Value) bool {
	// match: (Trunc32to16 <t> x)
	// cond: isSigned(t)
	// result: (MOVHreg x)
	for {
		t := v.Type
		x := v.Args[0]
		if !(isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVHreg)
		v.AddArg(x)
		return true
	}
	// match: (Trunc32to16 x)
	// result: (MOVHZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVHZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpTrunc32to8_0(v *Value) bool {
	// match: (Trunc32to8 <t> x)
	// cond: isSigned(t)
	// result: (MOVBreg x)
	for {
		t := v.Type
		x := v.Args[0]
		if !(isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
	// match: (Trunc32to8 x)
	// result: (MOVBZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpTrunc64to16_0(v *Value) bool {
	// match: (Trunc64to16 <t> x)
	// cond: isSigned(t)
	// result: (MOVHreg x)
	for {
		t := v.Type
		x := v.Args[0]
		if !(isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVHreg)
		v.AddArg(x)
		return true
	}
	// match: (Trunc64to16 x)
	// result: (MOVHZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVHZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpTrunc64to32_0(v *Value) bool {
	// match: (Trunc64to32 <t> x)
	// cond: isSigned(t)
	// result: (MOVWreg x)
	for {
		t := v.Type
		x := v.Args[0]
		if !(isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVWreg)
		v.AddArg(x)
		return true
	}
	// match: (Trunc64to32 x)
	// result: (MOVWZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVWZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpTrunc64to8_0(v *Value) bool {
	// match: (Trunc64to8 <t> x)
	// cond: isSigned(t)
	// result: (MOVBreg x)
	for {
		t := v.Type
		x := v.Args[0]
		if !(isSigned(t)) {
			break
		}
		v.reset(OpPPC64MOVBreg)
		v.AddArg(x)
		return true
	}
	// match: (Trunc64to8 x)
	// result: (MOVBZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpWB_0(v *Value) bool {
	// match: (WB {fn} destptr srcptr mem)
	// result: (LoweredWB {fn} destptr srcptr mem)
	for {
		fn := v.Aux
		mem := v.Args[2]
		destptr := v.Args[0]
		srcptr := v.Args[1]
		v.reset(OpPPC64LoweredWB)
		v.Aux = fn
		v.AddArg(destptr)
		v.AddArg(srcptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpXor16_0(v *Value) bool {
	// match: (Xor16 x y)
	// result: (XOR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpXor32_0(v *Value) bool {
	// match: (Xor32 x y)
	// result: (XOR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpXor64_0(v *Value) bool {
	// match: (Xor64 x y)
	// result: (XOR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpXor8_0(v *Value) bool {
	// match: (Xor8 x y)
	// result: (XOR x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpPPC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValuePPC64_OpZero_0(v *Value) bool {
	b := v.Block
	// match: (Zero [0] _ mem)
	// result: mem
	for {
		if v.AuxInt != 0 {
			break
		}
		mem := v.Args[1]
		v.reset(OpCopy)
		v.Type = mem.Type
		v.AddArg(mem)
		return true
	}
	// match: (Zero [1] destptr mem)
	// result: (MOVBstorezero destptr mem)
	for {
		if v.AuxInt != 1 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVBstorezero)
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [2] destptr mem)
	// result: (MOVHstorezero destptr mem)
	for {
		if v.AuxInt != 2 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVHstorezero)
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [3] destptr mem)
	// result: (MOVBstorezero [2] destptr (MOVHstorezero destptr mem))
	for {
		if v.AuxInt != 3 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVBstorezero)
		v.AuxInt = 2
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHstorezero, types.TypeMem)
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [4] destptr mem)
	// result: (MOVWstorezero destptr mem)
	for {
		if v.AuxInt != 4 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVWstorezero)
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [5] destptr mem)
	// result: (MOVBstorezero [4] destptr (MOVWstorezero destptr mem))
	for {
		if v.AuxInt != 5 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVBstorezero)
		v.AuxInt = 4
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWstorezero, types.TypeMem)
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [6] destptr mem)
	// result: (MOVHstorezero [4] destptr (MOVWstorezero destptr mem))
	for {
		if v.AuxInt != 6 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVHstorezero)
		v.AuxInt = 4
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWstorezero, types.TypeMem)
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [7] destptr mem)
	// result: (MOVBstorezero [6] destptr (MOVHstorezero [4] destptr (MOVWstorezero destptr mem)))
	for {
		if v.AuxInt != 7 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVBstorezero)
		v.AuxInt = 6
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVHstorezero, types.TypeMem)
		v0.AuxInt = 4
		v0.AddArg(destptr)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVWstorezero, types.TypeMem)
		v1.AddArg(destptr)
		v1.AddArg(mem)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [8] {t} destptr mem)
	// cond: t.(*types.Type).Alignment()%4 == 0
	// result: (MOVDstorezero destptr mem)
	for {
		if v.AuxInt != 8 {
			break
		}
		t := v.Aux
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(t.(*types.Type).Alignment()%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [8] destptr mem)
	// result: (MOVWstorezero [4] destptr (MOVWstorezero [0] destptr mem))
	for {
		if v.AuxInt != 8 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpPPC64MOVWstorezero)
		v.AuxInt = 4
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVWstorezero, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValuePPC64_OpZero_10(v *Value) bool {
	b := v.Block
	// match: (Zero [12] {t} destptr mem)
	// cond: t.(*types.Type).Alignment()%4 == 0
	// result: (MOVWstorezero [8] destptr (MOVDstorezero [0] destptr mem))
	for {
		if v.AuxInt != 12 {
			break
		}
		t := v.Aux
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(t.(*types.Type).Alignment()%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVWstorezero)
		v.AuxInt = 8
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [16] {t} destptr mem)
	// cond: t.(*types.Type).Alignment()%4 == 0
	// result: (MOVDstorezero [8] destptr (MOVDstorezero [0] destptr mem))
	for {
		if v.AuxInt != 16 {
			break
		}
		t := v.Aux
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(t.(*types.Type).Alignment()%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AuxInt = 8
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [24] {t} destptr mem)
	// cond: t.(*types.Type).Alignment()%4 == 0
	// result: (MOVDstorezero [16] destptr (MOVDstorezero [8] destptr (MOVDstorezero [0] destptr mem)))
	for {
		if v.AuxInt != 24 {
			break
		}
		t := v.Aux
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(t.(*types.Type).Alignment()%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AuxInt = 16
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v0.AuxInt = 8
		v0.AddArg(destptr)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v1.AuxInt = 0
		v1.AddArg(destptr)
		v1.AddArg(mem)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [32] {t} destptr mem)
	// cond: t.(*types.Type).Alignment()%4 == 0
	// result: (MOVDstorezero [24] destptr (MOVDstorezero [16] destptr (MOVDstorezero [8] destptr (MOVDstorezero [0] destptr mem))))
	for {
		if v.AuxInt != 32 {
			break
		}
		t := v.Aux
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(t.(*types.Type).Alignment()%4 == 0) {
			break
		}
		v.reset(OpPPC64MOVDstorezero)
		v.AuxInt = 24
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v0.AuxInt = 16
		v0.AddArg(destptr)
		v1 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v1.AuxInt = 8
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpPPC64MOVDstorezero, types.TypeMem)
		v2.AuxInt = 0
		v2.AddArg(destptr)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [s] ptr mem)
	// result: (LoweredZero [s] ptr mem)
	for {
		s := v.AuxInt
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpPPC64LoweredZero)
		v.AuxInt = s
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValuePPC64_OpZeroExt16to32_0(v *Value) bool {
	// match: (ZeroExt16to32 x)
	// result: (MOVHZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVHZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpZeroExt16to64_0(v *Value) bool {
	// match: (ZeroExt16to64 x)
	// result: (MOVHZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVHZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpZeroExt32to64_0(v *Value) bool {
	// match: (ZeroExt32to64 x)
	// result: (MOVWZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVWZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpZeroExt8to16_0(v *Value) bool {
	// match: (ZeroExt8to16 x)
	// result: (MOVBZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpZeroExt8to32_0(v *Value) bool {
	// match: (ZeroExt8to32 x)
	// result: (MOVBZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValuePPC64_OpZeroExt8to64_0(v *Value) bool {
	// match: (ZeroExt8to64 x)
	// result: (MOVBZreg x)
	for {
		x := v.Args[0]
		v.reset(OpPPC64MOVBZreg)
		v.AddArg(x)
		return true
	}
}
func rewriteBlockPPC64(b *Block) bool {
	switch b.Kind {
	case BlockPPC64EQ:
		// match: (EQ (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (EQ (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (EQ (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (EQ (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (EQ (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (EQ (FlagLT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagLT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (EQ (FlagGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (EQ (InvertFlags cmp) yes no)
		// result: (EQ cmp yes no)
		for b.Controls[0].Op == OpPPC64InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(BlockPPC64EQ)
			b.AddControl(cmp)
			return true
		}
		// match: (EQ (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (EQ (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (EQ (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (EQ (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (EQ (CMPconst [0] z:(AND x y)) yes no)
		// cond: z.Uses == 1
		// result: (EQ (ANDCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64AND {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (EQ (CMPconst [0] z:(OR x y)) yes no)
		// cond: z.Uses == 1
		// result: (EQ (ORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64OR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (EQ (CMPconst [0] z:(XOR x y)) yes no)
		// cond: z.Uses == 1
		// result: (EQ (XORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64XOR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64EQ)
			v0 := b.NewValue0(v_0.Pos, OpPPC64XORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
	case BlockPPC64GE:
		// match: (GE (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (GE (FlagLT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagLT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GE (FlagGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (GE (InvertFlags cmp) yes no)
		// result: (LE cmp yes no)
		for b.Controls[0].Op == OpPPC64InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(BlockPPC64LE)
			b.AddControl(cmp)
			return true
		}
		// match: (GE (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (GE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64GE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (GE (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (GE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64GE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (GE (CMPconst [0] z:(AND x y)) yes no)
		// cond: z.Uses == 1
		// result: (GE (ANDCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64AND {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64GE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (GE (CMPconst [0] z:(OR x y)) yes no)
		// cond: z.Uses == 1
		// result: (GE (ORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64OR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64GE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (GE (CMPconst [0] z:(XOR x y)) yes no)
		// cond: z.Uses == 1
		// result: (GE (XORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64XOR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64GE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64XORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
	case BlockPPC64GT:
		// match: (GT (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GT (FlagLT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagLT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GT (FlagGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (GT (InvertFlags cmp) yes no)
		// result: (LT cmp yes no)
		for b.Controls[0].Op == OpPPC64InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(BlockPPC64LT)
			b.AddControl(cmp)
			return true
		}
		// match: (GT (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (GT (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64GT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (GT (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (GT (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64GT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (GT (CMPconst [0] z:(AND x y)) yes no)
		// cond: z.Uses == 1
		// result: (GT (ANDCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64AND {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64GT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (GT (CMPconst [0] z:(OR x y)) yes no)
		// cond: z.Uses == 1
		// result: (GT (ORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64OR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64GT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (GT (CMPconst [0] z:(XOR x y)) yes no)
		// cond: z.Uses == 1
		// result: (GT (XORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64XOR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64GT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64XORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
	case BlockIf:
		// match: (If (Equal cc) yes no)
		// result: (EQ cc yes no)
		for b.Controls[0].Op == OpPPC64Equal {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64EQ)
			b.AddControl(cc)
			return true
		}
		// match: (If (NotEqual cc) yes no)
		// result: (NE cc yes no)
		for b.Controls[0].Op == OpPPC64NotEqual {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64NE)
			b.AddControl(cc)
			return true
		}
		// match: (If (LessThan cc) yes no)
		// result: (LT cc yes no)
		for b.Controls[0].Op == OpPPC64LessThan {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64LT)
			b.AddControl(cc)
			return true
		}
		// match: (If (LessEqual cc) yes no)
		// result: (LE cc yes no)
		for b.Controls[0].Op == OpPPC64LessEqual {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64LE)
			b.AddControl(cc)
			return true
		}
		// match: (If (GreaterThan cc) yes no)
		// result: (GT cc yes no)
		for b.Controls[0].Op == OpPPC64GreaterThan {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64GT)
			b.AddControl(cc)
			return true
		}
		// match: (If (GreaterEqual cc) yes no)
		// result: (GE cc yes no)
		for b.Controls[0].Op == OpPPC64GreaterEqual {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64GE)
			b.AddControl(cc)
			return true
		}
		// match: (If (FLessThan cc) yes no)
		// result: (FLT cc yes no)
		for b.Controls[0].Op == OpPPC64FLessThan {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64FLT)
			b.AddControl(cc)
			return true
		}
		// match: (If (FLessEqual cc) yes no)
		// result: (FLE cc yes no)
		for b.Controls[0].Op == OpPPC64FLessEqual {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64FLE)
			b.AddControl(cc)
			return true
		}
		// match: (If (FGreaterThan cc) yes no)
		// result: (FGT cc yes no)
		for b.Controls[0].Op == OpPPC64FGreaterThan {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64FGT)
			b.AddControl(cc)
			return true
		}
		// match: (If (FGreaterEqual cc) yes no)
		// result: (FGE cc yes no)
		for b.Controls[0].Op == OpPPC64FGreaterEqual {
			v_0 := b.Controls[0]
			cc := v_0.Args[0]
			b.Reset(BlockPPC64FGE)
			b.AddControl(cc)
			return true
		}
		// match: (If cond yes no)
		// result: (NE (CMPWconst [0] cond) yes no)
		for {
			cond := b.Controls[0]
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(cond.Pos, OpPPC64CMPWconst, types.TypeFlags)
			v0.AuxInt = 0
			v0.AddArg(cond)
			b.AddControl(v0)
			return true
		}
	case BlockPPC64LE:
		// match: (LE (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LE (FlagLT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagLT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LE (FlagGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (LE (InvertFlags cmp) yes no)
		// result: (GE cmp yes no)
		for b.Controls[0].Op == OpPPC64InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(BlockPPC64GE)
			b.AddControl(cmp)
			return true
		}
		// match: (LE (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (LE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64LE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (LE (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (LE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64LE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (LE (CMPconst [0] z:(AND x y)) yes no)
		// cond: z.Uses == 1
		// result: (LE (ANDCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64AND {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64LE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (LE (CMPconst [0] z:(OR x y)) yes no)
		// cond: z.Uses == 1
		// result: (LE (ORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64OR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64LE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (LE (CMPconst [0] z:(XOR x y)) yes no)
		// cond: z.Uses == 1
		// result: (LE (XORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64XOR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64LE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64XORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
	case BlockPPC64LT:
		// match: (LT (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (LT (FlagLT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagLT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LT (FlagGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (LT (InvertFlags cmp) yes no)
		// result: (GT cmp yes no)
		for b.Controls[0].Op == OpPPC64InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(BlockPPC64GT)
			b.AddControl(cmp)
			return true
		}
		// match: (LT (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (LT (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64LT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (LT (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (LT (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64LT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (LT (CMPconst [0] z:(AND x y)) yes no)
		// cond: z.Uses == 1
		// result: (LT (ANDCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64AND {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64LT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (LT (CMPconst [0] z:(OR x y)) yes no)
		// cond: z.Uses == 1
		// result: (LT (ORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64OR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64LT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (LT (CMPconst [0] z:(XOR x y)) yes no)
		// cond: z.Uses == 1
		// result: (LT (XORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64XOR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64LT)
			v0 := b.NewValue0(v_0.Pos, OpPPC64XORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
	case BlockPPC64NE:
		// match: (NE (CMPWconst [0] (Equal cc)) yes no)
		// result: (EQ cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64Equal {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64EQ)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (NotEqual cc)) yes no)
		// result: (NE cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64NotEqual {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64NE)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (LessThan cc)) yes no)
		// result: (LT cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64LessThan {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64LT)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (LessEqual cc)) yes no)
		// result: (LE cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64LessEqual {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64LE)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (GreaterThan cc)) yes no)
		// result: (GT cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64GreaterThan {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64GT)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (GreaterEqual cc)) yes no)
		// result: (GE cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64GreaterEqual {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64GE)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (FLessThan cc)) yes no)
		// result: (FLT cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64FLessThan {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64FLT)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (FLessEqual cc)) yes no)
		// result: (FLE cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64FLessEqual {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64FLE)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (FGreaterThan cc)) yes no)
		// result: (FGT cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64FGreaterThan {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64FGT)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPWconst [0] (FGreaterEqual cc)) yes no)
		// result: (FGE cc yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64FGreaterEqual {
				break
			}
			cc := v_0_0.Args[0]
			b.Reset(BlockPPC64FGE)
			b.AddControl(cc)
			return true
		}
		// match: (NE (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (NE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (NE (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (NE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (NE (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == OpPPC64FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (NE (FlagLT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagLT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (NE (FlagGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == OpPPC64FlagGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (NE (InvertFlags cmp) yes no)
		// result: (NE cmp yes no)
		for b.Controls[0].Op == OpPPC64InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(BlockPPC64NE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (CMPconst [0] (ANDconst [c] x)) yes no)
		// result: (NE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (NE (CMPWconst [0] (ANDconst [c] x)) yes no)
		// result: (NE (ANDCCconst [c] x) yes no)
		for b.Controls[0].Op == OpPPC64CMPWconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != OpPPC64ANDconst {
				break
			}
			c := v_0_0.AuxInt
			x := v_0_0.Args[0]
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCCconst, types.TypeFlags)
			v0.AuxInt = c
			v0.AddArg(x)
			b.AddControl(v0)
			return true
		}
		// match: (NE (CMPconst [0] z:(AND x y)) yes no)
		// cond: z.Uses == 1
		// result: (NE (ANDCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64AND {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ANDCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (NE (CMPconst [0] z:(OR x y)) yes no)
		// cond: z.Uses == 1
		// result: (NE (ORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64OR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64ORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
		// match: (NE (CMPconst [0] z:(XOR x y)) yes no)
		// cond: z.Uses == 1
		// result: (NE (XORCC x y) yes no)
		for b.Controls[0].Op == OpPPC64CMPconst {
			v_0 := b.Controls[0]
			if v_0.AuxInt != 0 {
				break
			}
			z := v_0.Args[0]
			if z.Op != OpPPC64XOR {
				break
			}
			y := z.Args[1]
			x := z.Args[0]
			if !(z.Uses == 1) {
				break
			}
			b.Reset(BlockPPC64NE)
			v0 := b.NewValue0(v_0.Pos, OpPPC64XORCC, types.TypeFlags)
			v0.AddArg(x)
			v0.AddArg(y)
			b.AddControl(v0)
			return true
		}
	}
	return false
}
