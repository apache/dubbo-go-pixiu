// Code generated from gen/Wasm.rules; DO NOT EDIT.
// generated with: cd gen; go run *.go

package ssa

import "cmd/internal/objabi"
import "cmd/compile/internal/types"

func rewriteValueWasm(v *Value) bool {
	switch v.Op {
	case OpAbs:
		return rewriteValueWasm_OpAbs_0(v)
	case OpAdd16:
		return rewriteValueWasm_OpAdd16_0(v)
	case OpAdd32:
		return rewriteValueWasm_OpAdd32_0(v)
	case OpAdd32F:
		return rewriteValueWasm_OpAdd32F_0(v)
	case OpAdd64:
		return rewriteValueWasm_OpAdd64_0(v)
	case OpAdd64F:
		return rewriteValueWasm_OpAdd64F_0(v)
	case OpAdd8:
		return rewriteValueWasm_OpAdd8_0(v)
	case OpAddPtr:
		return rewriteValueWasm_OpAddPtr_0(v)
	case OpAddr:
		return rewriteValueWasm_OpAddr_0(v)
	case OpAnd16:
		return rewriteValueWasm_OpAnd16_0(v)
	case OpAnd32:
		return rewriteValueWasm_OpAnd32_0(v)
	case OpAnd64:
		return rewriteValueWasm_OpAnd64_0(v)
	case OpAnd8:
		return rewriteValueWasm_OpAnd8_0(v)
	case OpAndB:
		return rewriteValueWasm_OpAndB_0(v)
	case OpBitLen64:
		return rewriteValueWasm_OpBitLen64_0(v)
	case OpCeil:
		return rewriteValueWasm_OpCeil_0(v)
	case OpClosureCall:
		return rewriteValueWasm_OpClosureCall_0(v)
	case OpCom16:
		return rewriteValueWasm_OpCom16_0(v)
	case OpCom32:
		return rewriteValueWasm_OpCom32_0(v)
	case OpCom64:
		return rewriteValueWasm_OpCom64_0(v)
	case OpCom8:
		return rewriteValueWasm_OpCom8_0(v)
	case OpCondSelect:
		return rewriteValueWasm_OpCondSelect_0(v)
	case OpConst16:
		return rewriteValueWasm_OpConst16_0(v)
	case OpConst32:
		return rewriteValueWasm_OpConst32_0(v)
	case OpConst32F:
		return rewriteValueWasm_OpConst32F_0(v)
	case OpConst64:
		return rewriteValueWasm_OpConst64_0(v)
	case OpConst64F:
		return rewriteValueWasm_OpConst64F_0(v)
	case OpConst8:
		return rewriteValueWasm_OpConst8_0(v)
	case OpConstBool:
		return rewriteValueWasm_OpConstBool_0(v)
	case OpConstNil:
		return rewriteValueWasm_OpConstNil_0(v)
	case OpConvert:
		return rewriteValueWasm_OpConvert_0(v)
	case OpCopysign:
		return rewriteValueWasm_OpCopysign_0(v)
	case OpCtz16:
		return rewriteValueWasm_OpCtz16_0(v)
	case OpCtz16NonZero:
		return rewriteValueWasm_OpCtz16NonZero_0(v)
	case OpCtz32:
		return rewriteValueWasm_OpCtz32_0(v)
	case OpCtz32NonZero:
		return rewriteValueWasm_OpCtz32NonZero_0(v)
	case OpCtz64:
		return rewriteValueWasm_OpCtz64_0(v)
	case OpCtz64NonZero:
		return rewriteValueWasm_OpCtz64NonZero_0(v)
	case OpCtz8:
		return rewriteValueWasm_OpCtz8_0(v)
	case OpCtz8NonZero:
		return rewriteValueWasm_OpCtz8NonZero_0(v)
	case OpCvt32Fto32:
		return rewriteValueWasm_OpCvt32Fto32_0(v)
	case OpCvt32Fto32U:
		return rewriteValueWasm_OpCvt32Fto32U_0(v)
	case OpCvt32Fto64:
		return rewriteValueWasm_OpCvt32Fto64_0(v)
	case OpCvt32Fto64F:
		return rewriteValueWasm_OpCvt32Fto64F_0(v)
	case OpCvt32Fto64U:
		return rewriteValueWasm_OpCvt32Fto64U_0(v)
	case OpCvt32Uto32F:
		return rewriteValueWasm_OpCvt32Uto32F_0(v)
	case OpCvt32Uto64F:
		return rewriteValueWasm_OpCvt32Uto64F_0(v)
	case OpCvt32to32F:
		return rewriteValueWasm_OpCvt32to32F_0(v)
	case OpCvt32to64F:
		return rewriteValueWasm_OpCvt32to64F_0(v)
	case OpCvt64Fto32:
		return rewriteValueWasm_OpCvt64Fto32_0(v)
	case OpCvt64Fto32F:
		return rewriteValueWasm_OpCvt64Fto32F_0(v)
	case OpCvt64Fto32U:
		return rewriteValueWasm_OpCvt64Fto32U_0(v)
	case OpCvt64Fto64:
		return rewriteValueWasm_OpCvt64Fto64_0(v)
	case OpCvt64Fto64U:
		return rewriteValueWasm_OpCvt64Fto64U_0(v)
	case OpCvt64Uto32F:
		return rewriteValueWasm_OpCvt64Uto32F_0(v)
	case OpCvt64Uto64F:
		return rewriteValueWasm_OpCvt64Uto64F_0(v)
	case OpCvt64to32F:
		return rewriteValueWasm_OpCvt64to32F_0(v)
	case OpCvt64to64F:
		return rewriteValueWasm_OpCvt64to64F_0(v)
	case OpDiv16:
		return rewriteValueWasm_OpDiv16_0(v)
	case OpDiv16u:
		return rewriteValueWasm_OpDiv16u_0(v)
	case OpDiv32:
		return rewriteValueWasm_OpDiv32_0(v)
	case OpDiv32F:
		return rewriteValueWasm_OpDiv32F_0(v)
	case OpDiv32u:
		return rewriteValueWasm_OpDiv32u_0(v)
	case OpDiv64:
		return rewriteValueWasm_OpDiv64_0(v)
	case OpDiv64F:
		return rewriteValueWasm_OpDiv64F_0(v)
	case OpDiv64u:
		return rewriteValueWasm_OpDiv64u_0(v)
	case OpDiv8:
		return rewriteValueWasm_OpDiv8_0(v)
	case OpDiv8u:
		return rewriteValueWasm_OpDiv8u_0(v)
	case OpEq16:
		return rewriteValueWasm_OpEq16_0(v)
	case OpEq32:
		return rewriteValueWasm_OpEq32_0(v)
	case OpEq32F:
		return rewriteValueWasm_OpEq32F_0(v)
	case OpEq64:
		return rewriteValueWasm_OpEq64_0(v)
	case OpEq64F:
		return rewriteValueWasm_OpEq64F_0(v)
	case OpEq8:
		return rewriteValueWasm_OpEq8_0(v)
	case OpEqB:
		return rewriteValueWasm_OpEqB_0(v)
	case OpEqPtr:
		return rewriteValueWasm_OpEqPtr_0(v)
	case OpFloor:
		return rewriteValueWasm_OpFloor_0(v)
	case OpGeq16:
		return rewriteValueWasm_OpGeq16_0(v)
	case OpGeq16U:
		return rewriteValueWasm_OpGeq16U_0(v)
	case OpGeq32:
		return rewriteValueWasm_OpGeq32_0(v)
	case OpGeq32F:
		return rewriteValueWasm_OpGeq32F_0(v)
	case OpGeq32U:
		return rewriteValueWasm_OpGeq32U_0(v)
	case OpGeq64:
		return rewriteValueWasm_OpGeq64_0(v)
	case OpGeq64F:
		return rewriteValueWasm_OpGeq64F_0(v)
	case OpGeq64U:
		return rewriteValueWasm_OpGeq64U_0(v)
	case OpGeq8:
		return rewriteValueWasm_OpGeq8_0(v)
	case OpGeq8U:
		return rewriteValueWasm_OpGeq8U_0(v)
	case OpGetCallerPC:
		return rewriteValueWasm_OpGetCallerPC_0(v)
	case OpGetCallerSP:
		return rewriteValueWasm_OpGetCallerSP_0(v)
	case OpGetClosurePtr:
		return rewriteValueWasm_OpGetClosurePtr_0(v)
	case OpGreater16:
		return rewriteValueWasm_OpGreater16_0(v)
	case OpGreater16U:
		return rewriteValueWasm_OpGreater16U_0(v)
	case OpGreater32:
		return rewriteValueWasm_OpGreater32_0(v)
	case OpGreater32F:
		return rewriteValueWasm_OpGreater32F_0(v)
	case OpGreater32U:
		return rewriteValueWasm_OpGreater32U_0(v)
	case OpGreater64:
		return rewriteValueWasm_OpGreater64_0(v)
	case OpGreater64F:
		return rewriteValueWasm_OpGreater64F_0(v)
	case OpGreater64U:
		return rewriteValueWasm_OpGreater64U_0(v)
	case OpGreater8:
		return rewriteValueWasm_OpGreater8_0(v)
	case OpGreater8U:
		return rewriteValueWasm_OpGreater8U_0(v)
	case OpInterCall:
		return rewriteValueWasm_OpInterCall_0(v)
	case OpIsInBounds:
		return rewriteValueWasm_OpIsInBounds_0(v)
	case OpIsNonNil:
		return rewriteValueWasm_OpIsNonNil_0(v)
	case OpIsSliceInBounds:
		return rewriteValueWasm_OpIsSliceInBounds_0(v)
	case OpLeq16:
		return rewriteValueWasm_OpLeq16_0(v)
	case OpLeq16U:
		return rewriteValueWasm_OpLeq16U_0(v)
	case OpLeq32:
		return rewriteValueWasm_OpLeq32_0(v)
	case OpLeq32F:
		return rewriteValueWasm_OpLeq32F_0(v)
	case OpLeq32U:
		return rewriteValueWasm_OpLeq32U_0(v)
	case OpLeq64:
		return rewriteValueWasm_OpLeq64_0(v)
	case OpLeq64F:
		return rewriteValueWasm_OpLeq64F_0(v)
	case OpLeq64U:
		return rewriteValueWasm_OpLeq64U_0(v)
	case OpLeq8:
		return rewriteValueWasm_OpLeq8_0(v)
	case OpLeq8U:
		return rewriteValueWasm_OpLeq8U_0(v)
	case OpLess16:
		return rewriteValueWasm_OpLess16_0(v)
	case OpLess16U:
		return rewriteValueWasm_OpLess16U_0(v)
	case OpLess32:
		return rewriteValueWasm_OpLess32_0(v)
	case OpLess32F:
		return rewriteValueWasm_OpLess32F_0(v)
	case OpLess32U:
		return rewriteValueWasm_OpLess32U_0(v)
	case OpLess64:
		return rewriteValueWasm_OpLess64_0(v)
	case OpLess64F:
		return rewriteValueWasm_OpLess64F_0(v)
	case OpLess64U:
		return rewriteValueWasm_OpLess64U_0(v)
	case OpLess8:
		return rewriteValueWasm_OpLess8_0(v)
	case OpLess8U:
		return rewriteValueWasm_OpLess8U_0(v)
	case OpLoad:
		return rewriteValueWasm_OpLoad_0(v)
	case OpLocalAddr:
		return rewriteValueWasm_OpLocalAddr_0(v)
	case OpLsh16x16:
		return rewriteValueWasm_OpLsh16x16_0(v)
	case OpLsh16x32:
		return rewriteValueWasm_OpLsh16x32_0(v)
	case OpLsh16x64:
		return rewriteValueWasm_OpLsh16x64_0(v)
	case OpLsh16x8:
		return rewriteValueWasm_OpLsh16x8_0(v)
	case OpLsh32x16:
		return rewriteValueWasm_OpLsh32x16_0(v)
	case OpLsh32x32:
		return rewriteValueWasm_OpLsh32x32_0(v)
	case OpLsh32x64:
		return rewriteValueWasm_OpLsh32x64_0(v)
	case OpLsh32x8:
		return rewriteValueWasm_OpLsh32x8_0(v)
	case OpLsh64x16:
		return rewriteValueWasm_OpLsh64x16_0(v)
	case OpLsh64x32:
		return rewriteValueWasm_OpLsh64x32_0(v)
	case OpLsh64x64:
		return rewriteValueWasm_OpLsh64x64_0(v)
	case OpLsh64x8:
		return rewriteValueWasm_OpLsh64x8_0(v)
	case OpLsh8x16:
		return rewriteValueWasm_OpLsh8x16_0(v)
	case OpLsh8x32:
		return rewriteValueWasm_OpLsh8x32_0(v)
	case OpLsh8x64:
		return rewriteValueWasm_OpLsh8x64_0(v)
	case OpLsh8x8:
		return rewriteValueWasm_OpLsh8x8_0(v)
	case OpMod16:
		return rewriteValueWasm_OpMod16_0(v)
	case OpMod16u:
		return rewriteValueWasm_OpMod16u_0(v)
	case OpMod32:
		return rewriteValueWasm_OpMod32_0(v)
	case OpMod32u:
		return rewriteValueWasm_OpMod32u_0(v)
	case OpMod64:
		return rewriteValueWasm_OpMod64_0(v)
	case OpMod64u:
		return rewriteValueWasm_OpMod64u_0(v)
	case OpMod8:
		return rewriteValueWasm_OpMod8_0(v)
	case OpMod8u:
		return rewriteValueWasm_OpMod8u_0(v)
	case OpMove:
		return rewriteValueWasm_OpMove_0(v) || rewriteValueWasm_OpMove_10(v)
	case OpMul16:
		return rewriteValueWasm_OpMul16_0(v)
	case OpMul32:
		return rewriteValueWasm_OpMul32_0(v)
	case OpMul32F:
		return rewriteValueWasm_OpMul32F_0(v)
	case OpMul64:
		return rewriteValueWasm_OpMul64_0(v)
	case OpMul64F:
		return rewriteValueWasm_OpMul64F_0(v)
	case OpMul8:
		return rewriteValueWasm_OpMul8_0(v)
	case OpNeg16:
		return rewriteValueWasm_OpNeg16_0(v)
	case OpNeg32:
		return rewriteValueWasm_OpNeg32_0(v)
	case OpNeg32F:
		return rewriteValueWasm_OpNeg32F_0(v)
	case OpNeg64:
		return rewriteValueWasm_OpNeg64_0(v)
	case OpNeg64F:
		return rewriteValueWasm_OpNeg64F_0(v)
	case OpNeg8:
		return rewriteValueWasm_OpNeg8_0(v)
	case OpNeq16:
		return rewriteValueWasm_OpNeq16_0(v)
	case OpNeq32:
		return rewriteValueWasm_OpNeq32_0(v)
	case OpNeq32F:
		return rewriteValueWasm_OpNeq32F_0(v)
	case OpNeq64:
		return rewriteValueWasm_OpNeq64_0(v)
	case OpNeq64F:
		return rewriteValueWasm_OpNeq64F_0(v)
	case OpNeq8:
		return rewriteValueWasm_OpNeq8_0(v)
	case OpNeqB:
		return rewriteValueWasm_OpNeqB_0(v)
	case OpNeqPtr:
		return rewriteValueWasm_OpNeqPtr_0(v)
	case OpNilCheck:
		return rewriteValueWasm_OpNilCheck_0(v)
	case OpNot:
		return rewriteValueWasm_OpNot_0(v)
	case OpOffPtr:
		return rewriteValueWasm_OpOffPtr_0(v)
	case OpOr16:
		return rewriteValueWasm_OpOr16_0(v)
	case OpOr32:
		return rewriteValueWasm_OpOr32_0(v)
	case OpOr64:
		return rewriteValueWasm_OpOr64_0(v)
	case OpOr8:
		return rewriteValueWasm_OpOr8_0(v)
	case OpOrB:
		return rewriteValueWasm_OpOrB_0(v)
	case OpPopCount16:
		return rewriteValueWasm_OpPopCount16_0(v)
	case OpPopCount32:
		return rewriteValueWasm_OpPopCount32_0(v)
	case OpPopCount64:
		return rewriteValueWasm_OpPopCount64_0(v)
	case OpPopCount8:
		return rewriteValueWasm_OpPopCount8_0(v)
	case OpRotateLeft16:
		return rewriteValueWasm_OpRotateLeft16_0(v)
	case OpRotateLeft32:
		return rewriteValueWasm_OpRotateLeft32_0(v)
	case OpRotateLeft64:
		return rewriteValueWasm_OpRotateLeft64_0(v)
	case OpRotateLeft8:
		return rewriteValueWasm_OpRotateLeft8_0(v)
	case OpRound32F:
		return rewriteValueWasm_OpRound32F_0(v)
	case OpRound64F:
		return rewriteValueWasm_OpRound64F_0(v)
	case OpRoundToEven:
		return rewriteValueWasm_OpRoundToEven_0(v)
	case OpRsh16Ux16:
		return rewriteValueWasm_OpRsh16Ux16_0(v)
	case OpRsh16Ux32:
		return rewriteValueWasm_OpRsh16Ux32_0(v)
	case OpRsh16Ux64:
		return rewriteValueWasm_OpRsh16Ux64_0(v)
	case OpRsh16Ux8:
		return rewriteValueWasm_OpRsh16Ux8_0(v)
	case OpRsh16x16:
		return rewriteValueWasm_OpRsh16x16_0(v)
	case OpRsh16x32:
		return rewriteValueWasm_OpRsh16x32_0(v)
	case OpRsh16x64:
		return rewriteValueWasm_OpRsh16x64_0(v)
	case OpRsh16x8:
		return rewriteValueWasm_OpRsh16x8_0(v)
	case OpRsh32Ux16:
		return rewriteValueWasm_OpRsh32Ux16_0(v)
	case OpRsh32Ux32:
		return rewriteValueWasm_OpRsh32Ux32_0(v)
	case OpRsh32Ux64:
		return rewriteValueWasm_OpRsh32Ux64_0(v)
	case OpRsh32Ux8:
		return rewriteValueWasm_OpRsh32Ux8_0(v)
	case OpRsh32x16:
		return rewriteValueWasm_OpRsh32x16_0(v)
	case OpRsh32x32:
		return rewriteValueWasm_OpRsh32x32_0(v)
	case OpRsh32x64:
		return rewriteValueWasm_OpRsh32x64_0(v)
	case OpRsh32x8:
		return rewriteValueWasm_OpRsh32x8_0(v)
	case OpRsh64Ux16:
		return rewriteValueWasm_OpRsh64Ux16_0(v)
	case OpRsh64Ux32:
		return rewriteValueWasm_OpRsh64Ux32_0(v)
	case OpRsh64Ux64:
		return rewriteValueWasm_OpRsh64Ux64_0(v)
	case OpRsh64Ux8:
		return rewriteValueWasm_OpRsh64Ux8_0(v)
	case OpRsh64x16:
		return rewriteValueWasm_OpRsh64x16_0(v)
	case OpRsh64x32:
		return rewriteValueWasm_OpRsh64x32_0(v)
	case OpRsh64x64:
		return rewriteValueWasm_OpRsh64x64_0(v)
	case OpRsh64x8:
		return rewriteValueWasm_OpRsh64x8_0(v)
	case OpRsh8Ux16:
		return rewriteValueWasm_OpRsh8Ux16_0(v)
	case OpRsh8Ux32:
		return rewriteValueWasm_OpRsh8Ux32_0(v)
	case OpRsh8Ux64:
		return rewriteValueWasm_OpRsh8Ux64_0(v)
	case OpRsh8Ux8:
		return rewriteValueWasm_OpRsh8Ux8_0(v)
	case OpRsh8x16:
		return rewriteValueWasm_OpRsh8x16_0(v)
	case OpRsh8x32:
		return rewriteValueWasm_OpRsh8x32_0(v)
	case OpRsh8x64:
		return rewriteValueWasm_OpRsh8x64_0(v)
	case OpRsh8x8:
		return rewriteValueWasm_OpRsh8x8_0(v)
	case OpSignExt16to32:
		return rewriteValueWasm_OpSignExt16to32_0(v)
	case OpSignExt16to64:
		return rewriteValueWasm_OpSignExt16to64_0(v)
	case OpSignExt32to64:
		return rewriteValueWasm_OpSignExt32to64_0(v)
	case OpSignExt8to16:
		return rewriteValueWasm_OpSignExt8to16_0(v)
	case OpSignExt8to32:
		return rewriteValueWasm_OpSignExt8to32_0(v)
	case OpSignExt8to64:
		return rewriteValueWasm_OpSignExt8to64_0(v)
	case OpSlicemask:
		return rewriteValueWasm_OpSlicemask_0(v)
	case OpSqrt:
		return rewriteValueWasm_OpSqrt_0(v)
	case OpStaticCall:
		return rewriteValueWasm_OpStaticCall_0(v)
	case OpStore:
		return rewriteValueWasm_OpStore_0(v)
	case OpSub16:
		return rewriteValueWasm_OpSub16_0(v)
	case OpSub32:
		return rewriteValueWasm_OpSub32_0(v)
	case OpSub32F:
		return rewriteValueWasm_OpSub32F_0(v)
	case OpSub64:
		return rewriteValueWasm_OpSub64_0(v)
	case OpSub64F:
		return rewriteValueWasm_OpSub64F_0(v)
	case OpSub8:
		return rewriteValueWasm_OpSub8_0(v)
	case OpSubPtr:
		return rewriteValueWasm_OpSubPtr_0(v)
	case OpTrunc:
		return rewriteValueWasm_OpTrunc_0(v)
	case OpTrunc16to8:
		return rewriteValueWasm_OpTrunc16to8_0(v)
	case OpTrunc32to16:
		return rewriteValueWasm_OpTrunc32to16_0(v)
	case OpTrunc32to8:
		return rewriteValueWasm_OpTrunc32to8_0(v)
	case OpTrunc64to16:
		return rewriteValueWasm_OpTrunc64to16_0(v)
	case OpTrunc64to32:
		return rewriteValueWasm_OpTrunc64to32_0(v)
	case OpTrunc64to8:
		return rewriteValueWasm_OpTrunc64to8_0(v)
	case OpWB:
		return rewriteValueWasm_OpWB_0(v)
	case OpWasmF64Add:
		return rewriteValueWasm_OpWasmF64Add_0(v)
	case OpWasmF64Mul:
		return rewriteValueWasm_OpWasmF64Mul_0(v)
	case OpWasmI64Add:
		return rewriteValueWasm_OpWasmI64Add_0(v)
	case OpWasmI64AddConst:
		return rewriteValueWasm_OpWasmI64AddConst_0(v)
	case OpWasmI64And:
		return rewriteValueWasm_OpWasmI64And_0(v)
	case OpWasmI64Eq:
		return rewriteValueWasm_OpWasmI64Eq_0(v)
	case OpWasmI64Eqz:
		return rewriteValueWasm_OpWasmI64Eqz_0(v)
	case OpWasmI64Load:
		return rewriteValueWasm_OpWasmI64Load_0(v)
	case OpWasmI64Load16S:
		return rewriteValueWasm_OpWasmI64Load16S_0(v)
	case OpWasmI64Load16U:
		return rewriteValueWasm_OpWasmI64Load16U_0(v)
	case OpWasmI64Load32S:
		return rewriteValueWasm_OpWasmI64Load32S_0(v)
	case OpWasmI64Load32U:
		return rewriteValueWasm_OpWasmI64Load32U_0(v)
	case OpWasmI64Load8S:
		return rewriteValueWasm_OpWasmI64Load8S_0(v)
	case OpWasmI64Load8U:
		return rewriteValueWasm_OpWasmI64Load8U_0(v)
	case OpWasmI64Mul:
		return rewriteValueWasm_OpWasmI64Mul_0(v)
	case OpWasmI64Ne:
		return rewriteValueWasm_OpWasmI64Ne_0(v)
	case OpWasmI64Or:
		return rewriteValueWasm_OpWasmI64Or_0(v)
	case OpWasmI64Shl:
		return rewriteValueWasm_OpWasmI64Shl_0(v)
	case OpWasmI64ShrS:
		return rewriteValueWasm_OpWasmI64ShrS_0(v)
	case OpWasmI64ShrU:
		return rewriteValueWasm_OpWasmI64ShrU_0(v)
	case OpWasmI64Store:
		return rewriteValueWasm_OpWasmI64Store_0(v)
	case OpWasmI64Store16:
		return rewriteValueWasm_OpWasmI64Store16_0(v)
	case OpWasmI64Store32:
		return rewriteValueWasm_OpWasmI64Store32_0(v)
	case OpWasmI64Store8:
		return rewriteValueWasm_OpWasmI64Store8_0(v)
	case OpWasmI64Xor:
		return rewriteValueWasm_OpWasmI64Xor_0(v)
	case OpXor16:
		return rewriteValueWasm_OpXor16_0(v)
	case OpXor32:
		return rewriteValueWasm_OpXor32_0(v)
	case OpXor64:
		return rewriteValueWasm_OpXor64_0(v)
	case OpXor8:
		return rewriteValueWasm_OpXor8_0(v)
	case OpZero:
		return rewriteValueWasm_OpZero_0(v) || rewriteValueWasm_OpZero_10(v)
	case OpZeroExt16to32:
		return rewriteValueWasm_OpZeroExt16to32_0(v)
	case OpZeroExt16to64:
		return rewriteValueWasm_OpZeroExt16to64_0(v)
	case OpZeroExt32to64:
		return rewriteValueWasm_OpZeroExt32to64_0(v)
	case OpZeroExt8to16:
		return rewriteValueWasm_OpZeroExt8to16_0(v)
	case OpZeroExt8to32:
		return rewriteValueWasm_OpZeroExt8to32_0(v)
	case OpZeroExt8to64:
		return rewriteValueWasm_OpZeroExt8to64_0(v)
	}
	return false
}
func rewriteValueWasm_OpAbs_0(v *Value) bool {
	// match: (Abs x)
	// result: (F64Abs x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Abs)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpAdd16_0(v *Value) bool {
	// match: (Add16 x y)
	// result: (I64Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAdd32_0(v *Value) bool {
	// match: (Add32 x y)
	// result: (I64Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAdd32F_0(v *Value) bool {
	// match: (Add32F x y)
	// result: (F32Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAdd64_0(v *Value) bool {
	// match: (Add64 x y)
	// result: (I64Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAdd64F_0(v *Value) bool {
	// match: (Add64F x y)
	// result: (F64Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAdd8_0(v *Value) bool {
	// match: (Add8 x y)
	// result: (I64Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAddPtr_0(v *Value) bool {
	// match: (AddPtr x y)
	// result: (I64Add x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Add)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAddr_0(v *Value) bool {
	// match: (Addr {sym} base)
	// result: (LoweredAddr {sym} base)
	for {
		sym := v.Aux
		base := v.Args[0]
		v.reset(OpWasmLoweredAddr)
		v.Aux = sym
		v.AddArg(base)
		return true
	}
}
func rewriteValueWasm_OpAnd16_0(v *Value) bool {
	// match: (And16 x y)
	// result: (I64And x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAnd32_0(v *Value) bool {
	// match: (And32 x y)
	// result: (I64And x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAnd64_0(v *Value) bool {
	// match: (And64 x y)
	// result: (I64And x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAnd8_0(v *Value) bool {
	// match: (And8 x y)
	// result: (I64And x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpAndB_0(v *Value) bool {
	// match: (AndB x y)
	// result: (I64And x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpBitLen64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (BitLen64 x)
	// result: (I64Sub (I64Const [64]) (I64Clz x))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 64
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Clz, typ.Int64)
		v1.AddArg(x)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpCeil_0(v *Value) bool {
	// match: (Ceil x)
	// result: (F64Ceil x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Ceil)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpClosureCall_0(v *Value) bool {
	// match: (ClosureCall [argwid] entry closure mem)
	// result: (LoweredClosureCall [argwid] entry closure mem)
	for {
		argwid := v.AuxInt
		mem := v.Args[2]
		entry := v.Args[0]
		closure := v.Args[1]
		v.reset(OpWasmLoweredClosureCall)
		v.AuxInt = argwid
		v.AddArg(entry)
		v.AddArg(closure)
		v.AddArg(mem)
		return true
	}
}
func rewriteValueWasm_OpCom16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Com16 x)
	// result: (I64Xor x (I64Const [-1]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = -1
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCom32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Com32 x)
	// result: (I64Xor x (I64Const [-1]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = -1
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCom64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Com64 x)
	// result: (I64Xor x (I64Const [-1]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = -1
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCom8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Com8 x)
	// result: (I64Xor x (I64Const [-1]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = -1
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCondSelect_0(v *Value) bool {
	// match: (CondSelect <t> x y cond)
	// result: (Select <t> x y cond)
	for {
		t := v.Type
		cond := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpWasmSelect)
		v.Type = t
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(cond)
		return true
	}
}
func rewriteValueWasm_OpConst16_0(v *Value) bool {
	// match: (Const16 [val])
	// result: (I64Const [val])
	for {
		val := v.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = val
		return true
	}
}
func rewriteValueWasm_OpConst32_0(v *Value) bool {
	// match: (Const32 [val])
	// result: (I64Const [val])
	for {
		val := v.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = val
		return true
	}
}
func rewriteValueWasm_OpConst32F_0(v *Value) bool {
	// match: (Const32F [val])
	// result: (F32Const [val])
	for {
		val := v.AuxInt
		v.reset(OpWasmF32Const)
		v.AuxInt = val
		return true
	}
}
func rewriteValueWasm_OpConst64_0(v *Value) bool {
	// match: (Const64 [val])
	// result: (I64Const [val])
	for {
		val := v.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = val
		return true
	}
}
func rewriteValueWasm_OpConst64F_0(v *Value) bool {
	// match: (Const64F [val])
	// result: (F64Const [val])
	for {
		val := v.AuxInt
		v.reset(OpWasmF64Const)
		v.AuxInt = val
		return true
	}
}
func rewriteValueWasm_OpConst8_0(v *Value) bool {
	// match: (Const8 [val])
	// result: (I64Const [val])
	for {
		val := v.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = val
		return true
	}
}
func rewriteValueWasm_OpConstBool_0(v *Value) bool {
	// match: (ConstBool [b])
	// result: (I64Const [b])
	for {
		b := v.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = b
		return true
	}
}
func rewriteValueWasm_OpConstNil_0(v *Value) bool {
	// match: (ConstNil)
	// result: (I64Const [0])
	for {
		v.reset(OpWasmI64Const)
		v.AuxInt = 0
		return true
	}
}
func rewriteValueWasm_OpConvert_0(v *Value) bool {
	// match: (Convert <t> x mem)
	// result: (LoweredConvert <t> x mem)
	for {
		t := v.Type
		mem := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmLoweredConvert)
		v.Type = t
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
}
func rewriteValueWasm_OpCopysign_0(v *Value) bool {
	// match: (Copysign x y)
	// result: (F64Copysign x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Copysign)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpCtz16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz16 x)
	// result: (I64Ctz (I64Or x (I64Const [0x10000])))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v0 := b.NewValue0(v.Pos, OpWasmI64Or, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 0x10000
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCtz16NonZero_0(v *Value) bool {
	// match: (Ctz16NonZero x)
	// result: (I64Ctz x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCtz32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz32 x)
	// result: (I64Ctz (I64Or x (I64Const [0x100000000])))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v0 := b.NewValue0(v.Pos, OpWasmI64Or, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 0x100000000
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCtz32NonZero_0(v *Value) bool {
	// match: (Ctz32NonZero x)
	// result: (I64Ctz x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCtz64_0(v *Value) bool {
	// match: (Ctz64 x)
	// result: (I64Ctz x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCtz64NonZero_0(v *Value) bool {
	// match: (Ctz64NonZero x)
	// result: (I64Ctz x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCtz8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz8 x)
	// result: (I64Ctz (I64Or x (I64Const [0x100])))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v0 := b.NewValue0(v.Pos, OpWasmI64Or, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 0x100
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCtz8NonZero_0(v *Value) bool {
	// match: (Ctz8NonZero x)
	// result: (I64Ctz x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Ctz)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt32Fto32_0(v *Value) bool {
	// match: (Cvt32Fto32 x)
	// result: (I64TruncSatF32S x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF32S)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt32Fto32U_0(v *Value) bool {
	// match: (Cvt32Fto32U x)
	// result: (I64TruncSatF32U x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF32U)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt32Fto64_0(v *Value) bool {
	// match: (Cvt32Fto64 x)
	// result: (I64TruncSatF32S x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF32S)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt32Fto64F_0(v *Value) bool {
	// match: (Cvt32Fto64F x)
	// result: (F64PromoteF32 x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64PromoteF32)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt32Fto64U_0(v *Value) bool {
	// match: (Cvt32Fto64U x)
	// result: (I64TruncSatF32U x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF32U)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt32Uto32F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32Uto32F x)
	// result: (F32ConvertI64U (ZeroExt32to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmF32ConvertI64U)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCvt32Uto64F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32Uto64F x)
	// result: (F64ConvertI64U (ZeroExt32to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmF64ConvertI64U)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCvt32to32F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32to32F x)
	// result: (F32ConvertI64S (SignExt32to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmF32ConvertI64S)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCvt32to64F_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Cvt32to64F x)
	// result: (F64ConvertI64S (SignExt32to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmF64ConvertI64S)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpCvt64Fto32_0(v *Value) bool {
	// match: (Cvt64Fto32 x)
	// result: (I64TruncSatF64S x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF64S)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64Fto32F_0(v *Value) bool {
	// match: (Cvt64Fto32F x)
	// result: (F32DemoteF64 x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF32DemoteF64)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64Fto32U_0(v *Value) bool {
	// match: (Cvt64Fto32U x)
	// result: (I64TruncSatF64U x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF64U)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64Fto64_0(v *Value) bool {
	// match: (Cvt64Fto64 x)
	// result: (I64TruncSatF64S x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF64S)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64Fto64U_0(v *Value) bool {
	// match: (Cvt64Fto64U x)
	// result: (I64TruncSatF64U x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64TruncSatF64U)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64Uto32F_0(v *Value) bool {
	// match: (Cvt64Uto32F x)
	// result: (F32ConvertI64U x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF32ConvertI64U)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64Uto64F_0(v *Value) bool {
	// match: (Cvt64Uto64F x)
	// result: (F64ConvertI64U x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64ConvertI64U)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64to32F_0(v *Value) bool {
	// match: (Cvt64to32F x)
	// result: (F32ConvertI64S x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF32ConvertI64S)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpCvt64to64F_0(v *Value) bool {
	// match: (Cvt64to64F x)
	// result: (F64ConvertI64S x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64ConvertI64S)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpDiv16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div16 x y)
	// result: (I64DivS (SignExt16to64 x) (SignExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivS)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpDiv16u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div16u x y)
	// result: (I64DivU (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpDiv32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div32 x y)
	// result: (I64DivS (SignExt32to64 x) (SignExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivS)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpDiv32F_0(v *Value) bool {
	// match: (Div32F x y)
	// result: (F32Div x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Div)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpDiv32u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div32u x y)
	// result: (I64DivU (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivU)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpDiv64_0(v *Value) bool {
	// match: (Div64 x y)
	// result: (I64DivS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpDiv64F_0(v *Value) bool {
	// match: (Div64F x y)
	// result: (F64Div x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Div)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpDiv64u_0(v *Value) bool {
	// match: (Div64u x y)
	// result: (I64DivU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpDiv8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div8 x y)
	// result: (I64DivS (SignExt8to64 x) (SignExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivS)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpDiv8u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div8u x y)
	// result: (I64DivU (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64DivU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpEq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Eq16 x y)
	// result: (I64Eq (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Eq)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpEq32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Eq32 x y)
	// result: (I64Eq (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Eq)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpEq32F_0(v *Value) bool {
	// match: (Eq32F x y)
	// result: (F32Eq x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Eq)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpEq64_0(v *Value) bool {
	// match: (Eq64 x y)
	// result: (I64Eq x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Eq)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpEq64F_0(v *Value) bool {
	// match: (Eq64F x y)
	// result: (F64Eq x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Eq)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpEq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Eq8 x y)
	// result: (I64Eq (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Eq)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpEqB_0(v *Value) bool {
	// match: (EqB x y)
	// result: (I64Eq x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Eq)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpEqPtr_0(v *Value) bool {
	// match: (EqPtr x y)
	// result: (I64Eq x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Eq)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpFloor_0(v *Value) bool {
	// match: (Floor x)
	// result: (F64Floor x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Floor)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpGeq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq16 x y)
	// result: (I64GeS (SignExt16to64 x) (SignExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeS)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGeq16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq16U x y)
	// result: (I64GeU (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGeq32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq32 x y)
	// result: (I64GeS (SignExt32to64 x) (SignExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeS)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGeq32F_0(v *Value) bool {
	// match: (Geq32F x y)
	// result: (F32Ge x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Ge)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGeq32U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq32U x y)
	// result: (I64GeU (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeU)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGeq64_0(v *Value) bool {
	// match: (Geq64 x y)
	// result: (I64GeS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGeq64F_0(v *Value) bool {
	// match: (Geq64F x y)
	// result: (F64Ge x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Ge)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGeq64U_0(v *Value) bool {
	// match: (Geq64U x y)
	// result: (I64GeU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGeq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq8 x y)
	// result: (I64GeS (SignExt8to64 x) (SignExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeS)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGeq8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Geq8U x y)
	// result: (I64GeU (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GeU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGetCallerPC_0(v *Value) bool {
	// match: (GetCallerPC)
	// result: (LoweredGetCallerPC)
	for {
		v.reset(OpWasmLoweredGetCallerPC)
		return true
	}
}
func rewriteValueWasm_OpGetCallerSP_0(v *Value) bool {
	// match: (GetCallerSP)
	// result: (LoweredGetCallerSP)
	for {
		v.reset(OpWasmLoweredGetCallerSP)
		return true
	}
}
func rewriteValueWasm_OpGetClosurePtr_0(v *Value) bool {
	// match: (GetClosurePtr)
	// result: (LoweredGetClosurePtr)
	for {
		v.reset(OpWasmLoweredGetClosurePtr)
		return true
	}
}
func rewriteValueWasm_OpGreater16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater16 x y)
	// result: (I64GtS (SignExt16to64 x) (SignExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtS)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGreater16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater16U x y)
	// result: (I64GtU (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGreater32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater32 x y)
	// result: (I64GtS (SignExt32to64 x) (SignExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtS)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGreater32F_0(v *Value) bool {
	// match: (Greater32F x y)
	// result: (F32Gt x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Gt)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGreater32U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater32U x y)
	// result: (I64GtU (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtU)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGreater64_0(v *Value) bool {
	// match: (Greater64 x y)
	// result: (I64GtS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGreater64F_0(v *Value) bool {
	// match: (Greater64F x y)
	// result: (F64Gt x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Gt)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGreater64U_0(v *Value) bool {
	// match: (Greater64U x y)
	// result: (I64GtU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpGreater8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater8 x y)
	// result: (I64GtS (SignExt8to64 x) (SignExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtS)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpGreater8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Greater8U x y)
	// result: (I64GtU (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64GtU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpInterCall_0(v *Value) bool {
	// match: (InterCall [argwid] entry mem)
	// result: (LoweredInterCall [argwid] entry mem)
	for {
		argwid := v.AuxInt
		mem := v.Args[1]
		entry := v.Args[0]
		v.reset(OpWasmLoweredInterCall)
		v.AuxInt = argwid
		v.AddArg(entry)
		v.AddArg(mem)
		return true
	}
}
func rewriteValueWasm_OpIsInBounds_0(v *Value) bool {
	// match: (IsInBounds idx len)
	// result: (I64LtU idx len)
	for {
		len := v.Args[1]
		idx := v.Args[0]
		v.reset(OpWasmI64LtU)
		v.AddArg(idx)
		v.AddArg(len)
		return true
	}
}
func rewriteValueWasm_OpIsNonNil_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (IsNonNil p)
	// result: (I64Eqz (I64Eqz p))
	for {
		p := v.Args[0]
		v.reset(OpWasmI64Eqz)
		v0 := b.NewValue0(v.Pos, OpWasmI64Eqz, typ.Bool)
		v0.AddArg(p)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpIsSliceInBounds_0(v *Value) bool {
	// match: (IsSliceInBounds idx len)
	// result: (I64LeU idx len)
	for {
		len := v.Args[1]
		idx := v.Args[0]
		v.reset(OpWasmI64LeU)
		v.AddArg(idx)
		v.AddArg(len)
		return true
	}
}
func rewriteValueWasm_OpLeq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq16 x y)
	// result: (I64LeS (SignExt16to64 x) (SignExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeS)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLeq16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq16U x y)
	// result: (I64LeU (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLeq32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq32 x y)
	// result: (I64LeS (SignExt32to64 x) (SignExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeS)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLeq32F_0(v *Value) bool {
	// match: (Leq32F x y)
	// result: (F32Le x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Le)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLeq32U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq32U x y)
	// result: (I64LeU (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeU)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLeq64_0(v *Value) bool {
	// match: (Leq64 x y)
	// result: (I64LeS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLeq64F_0(v *Value) bool {
	// match: (Leq64F x y)
	// result: (F64Le x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Le)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLeq64U_0(v *Value) bool {
	// match: (Leq64U x y)
	// result: (I64LeU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLeq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq8 x y)
	// result: (I64LeS (SignExt8to64 x) (SignExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeS)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLeq8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Leq8U x y)
	// result: (I64LeU (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LeU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLess16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less16 x y)
	// result: (I64LtS (SignExt16to64 x) (SignExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtS)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLess16U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less16U x y)
	// result: (I64LtU (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLess32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less32 x y)
	// result: (I64LtS (SignExt32to64 x) (SignExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtS)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLess32F_0(v *Value) bool {
	// match: (Less32F x y)
	// result: (F32Lt x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Lt)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLess32U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less32U x y)
	// result: (I64LtU (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtU)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLess64_0(v *Value) bool {
	// match: (Less64 x y)
	// result: (I64LtS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLess64F_0(v *Value) bool {
	// match: (Less64F x y)
	// result: (F64Lt x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Lt)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLess64U_0(v *Value) bool {
	// match: (Less64U x y)
	// result: (I64LtU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLess8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less8 x y)
	// result: (I64LtS (SignExt8to64 x) (SignExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtS)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLess8U_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Less8U x y)
	// result: (I64LtU (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64LtU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpLoad_0(v *Value) bool {
	// match: (Load <t> ptr mem)
	// cond: is32BitFloat(t)
	// result: (F32Load ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is32BitFloat(t)) {
			break
		}
		v.reset(OpWasmF32Load)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is64BitFloat(t)
	// result: (F64Load ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is64BitFloat(t)) {
			break
		}
		v.reset(OpWasmF64Load)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 8
	// result: (I64Load ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 8) {
			break
		}
		v.reset(OpWasmI64Load)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 4 && !t.IsSigned()
	// result: (I64Load32U ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 4 && !t.IsSigned()) {
			break
		}
		v.reset(OpWasmI64Load32U)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 4 && t.IsSigned()
	// result: (I64Load32S ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 4 && t.IsSigned()) {
			break
		}
		v.reset(OpWasmI64Load32S)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 2 && !t.IsSigned()
	// result: (I64Load16U ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 2 && !t.IsSigned()) {
			break
		}
		v.reset(OpWasmI64Load16U)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 2 && t.IsSigned()
	// result: (I64Load16S ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 2 && t.IsSigned()) {
			break
		}
		v.reset(OpWasmI64Load16S)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 1 && !t.IsSigned()
	// result: (I64Load8U ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 1 && !t.IsSigned()) {
			break
		}
		v.reset(OpWasmI64Load8U)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: t.Size() == 1 && t.IsSigned()
	// result: (I64Load8S ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.Size() == 1 && t.IsSigned()) {
			break
		}
		v.reset(OpWasmI64Load8S)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpLocalAddr_0(v *Value) bool {
	// match: (LocalAddr {sym} base _)
	// result: (LoweredAddr {sym} base)
	for {
		sym := v.Aux
		_ = v.Args[1]
		base := v.Args[0]
		v.reset(OpWasmLoweredAddr)
		v.Aux = sym
		v.AddArg(base)
		return true
	}
}
func rewriteValueWasm_OpLsh16x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x16 x y)
	// result: (Lsh64x64 x (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh16x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x32 x y)
	// result: (Lsh64x64 x (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh16x64_0(v *Value) bool {
	// match: (Lsh16x64 x y)
	// result: (Lsh64x64 x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLsh16x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh16x8 x y)
	// result: (Lsh64x64 x (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh32x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x16 x y)
	// result: (Lsh64x64 x (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh32x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x32 x y)
	// result: (Lsh64x64 x (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh32x64_0(v *Value) bool {
	// match: (Lsh32x64 x y)
	// result: (Lsh64x64 x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLsh32x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh32x8 x y)
	// result: (Lsh64x64 x (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh64x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x16 x y)
	// result: (Lsh64x64 x (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh64x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x32 x y)
	// result: (Lsh64x64 x (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh64x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x64 x y)
	// cond: shiftIsBounded(v)
	// result: (I64Shl x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpWasmI64Shl)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Lsh64x64 x (I64Const [c]))
	// cond: uint64(c) < 64
	// result: (I64Shl x (I64Const [c]))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpWasmI64Shl)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = c
		v.AddArg(v0)
		return true
	}
	// match: (Lsh64x64 x (I64Const [c]))
	// cond: uint64(c) >= 64
	// result: (I64Const [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 64) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh64x64 x y)
	// result: (Select (I64Shl x y) (I64Const [0]) (I64LtU y (I64Const [64])))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmSelect)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 0
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpWasmI64LtU, typ.Bool)
		v2.AddArg(y)
		v3 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v3.AuxInt = 64
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpLsh64x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh64x8 x y)
	// result: (Lsh64x64 x (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh8x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x16 x y)
	// result: (Lsh64x64 x (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh8x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x32 x y)
	// result: (Lsh64x64 x (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpLsh8x64_0(v *Value) bool {
	// match: (Lsh8x64 x y)
	// result: (Lsh64x64 x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpLsh8x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Lsh8x8 x y)
	// result: (Lsh64x64 x (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpLsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpMod16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod16 x y)
	// result: (I64RemS (SignExt16to64 x) (SignExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemS)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpMod16u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod16u x y)
	// result: (I64RemU (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemU)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpMod32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod32 x y)
	// result: (I64RemS (SignExt32to64 x) (SignExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemS)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpMod32u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod32u x y)
	// result: (I64RemU (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemU)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpMod64_0(v *Value) bool {
	// match: (Mod64 x y)
	// result: (I64RemS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMod64u_0(v *Value) bool {
	// match: (Mod64u x y)
	// result: (I64RemU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMod8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod8 x y)
	// result: (I64RemS (SignExt8to64 x) (SignExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemS)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpMod8u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod8u x y)
	// result: (I64RemU (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64RemU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpMove_0(v *Value) bool {
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
	// result: (I64Store8 dst (I64Load8U src mem) mem)
	for {
		if v.AuxInt != 1 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store8)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load8U, typ.UInt8)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [2] dst src mem)
	// result: (I64Store16 dst (I64Load16U src mem) mem)
	for {
		if v.AuxInt != 2 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store16)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load16U, typ.UInt16)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [4] dst src mem)
	// result: (I64Store32 dst (I64Load32U src mem) mem)
	for {
		if v.AuxInt != 4 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store32)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load32U, typ.UInt32)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [8] dst src mem)
	// result: (I64Store dst (I64Load src mem) mem)
	for {
		if v.AuxInt != 8 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [16] dst src mem)
	// result: (I64Store [8] dst (I64Load [8] src mem) (I64Store dst (I64Load src mem) mem))
	for {
		if v.AuxInt != 16 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store)
		v.AuxInt = 8
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v0.AuxInt = 8
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [3] dst src mem)
	// result: (I64Store8 [2] dst (I64Load8U [2] src mem) (I64Store16 dst (I64Load16U src mem) mem))
	for {
		if v.AuxInt != 3 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store8)
		v.AuxInt = 2
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load8U, typ.UInt8)
		v0.AuxInt = 2
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store16, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpWasmI64Load16U, typ.UInt16)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [5] dst src mem)
	// result: (I64Store8 [4] dst (I64Load8U [4] src mem) (I64Store32 dst (I64Load32U src mem) mem))
	for {
		if v.AuxInt != 5 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store8)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load8U, typ.UInt8)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store32, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpWasmI64Load32U, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [6] dst src mem)
	// result: (I64Store16 [4] dst (I64Load16U [4] src mem) (I64Store32 dst (I64Load32U src mem) mem))
	for {
		if v.AuxInt != 6 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store16)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load16U, typ.UInt16)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store32, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpWasmI64Load32U, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [7] dst src mem)
	// result: (I64Store32 [3] dst (I64Load32U [3] src mem) (I64Store32 dst (I64Load32U src mem) mem))
	for {
		if v.AuxInt != 7 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(OpWasmI64Store32)
		v.AuxInt = 3
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load32U, typ.UInt32)
		v0.AuxInt = 3
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store32, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpWasmI64Load32U, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValueWasm_OpMove_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Move [s] dst src mem)
	// cond: s > 8 && s < 16
	// result: (I64Store [s-8] dst (I64Load [s-8] src mem) (I64Store dst (I64Load src mem) mem))
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s > 8 && s < 16) {
			break
		}
		v.reset(OpWasmI64Store)
		v.AuxInt = s - 8
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v0.AuxInt = s - 8
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [s] dst src mem)
	// cond: s > 16 && s%16 != 0 && s%16 <= 8
	// result: (Move [s-s%16] (OffPtr <dst.Type> dst [s%16]) (OffPtr <src.Type> src [s%16]) (I64Store dst (I64Load src mem) mem))
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s > 16 && s%16 != 0 && s%16 <= 8) {
			break
		}
		v.reset(OpMove)
		v.AuxInt = s - s%16
		v0 := b.NewValue0(v.Pos, OpOffPtr, dst.Type)
		v0.AuxInt = s % 16
		v0.AddArg(dst)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpOffPtr, src.Type)
		v1.AuxInt = s % 16
		v1.AddArg(src)
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v2.AddArg(dst)
		v3 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v3.AddArg(src)
		v3.AddArg(mem)
		v2.AddArg(v3)
		v2.AddArg(mem)
		v.AddArg(v2)
		return true
	}
	// match: (Move [s] dst src mem)
	// cond: s > 16 && s%16 != 0 && s%16 > 8
	// result: (Move [s-s%16] (OffPtr <dst.Type> dst [s%16]) (OffPtr <src.Type> src [s%16]) (I64Store [8] dst (I64Load [8] src mem) (I64Store dst (I64Load src mem) mem)))
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s > 16 && s%16 != 0 && s%16 > 8) {
			break
		}
		v.reset(OpMove)
		v.AuxInt = s - s%16
		v0 := b.NewValue0(v.Pos, OpOffPtr, dst.Type)
		v0.AuxInt = s % 16
		v0.AddArg(dst)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpOffPtr, src.Type)
		v1.AuxInt = s % 16
		v1.AddArg(src)
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v2.AuxInt = 8
		v2.AddArg(dst)
		v3 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v3.AuxInt = 8
		v3.AddArg(src)
		v3.AddArg(mem)
		v2.AddArg(v3)
		v4 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v4.AddArg(dst)
		v5 := b.NewValue0(v.Pos, OpWasmI64Load, typ.UInt64)
		v5.AddArg(src)
		v5.AddArg(mem)
		v4.AddArg(v5)
		v4.AddArg(mem)
		v2.AddArg(v4)
		v.AddArg(v2)
		return true
	}
	// match: (Move [s] dst src mem)
	// cond: s%8 == 0
	// result: (LoweredMove [s/8] dst src mem)
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s%8 == 0) {
			break
		}
		v.reset(OpWasmLoweredMove)
		v.AuxInt = s / 8
		v.AddArg(dst)
		v.AddArg(src)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpMul16_0(v *Value) bool {
	// match: (Mul16 x y)
	// result: (I64Mul x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Mul)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMul32_0(v *Value) bool {
	// match: (Mul32 x y)
	// result: (I64Mul x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Mul)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMul32F_0(v *Value) bool {
	// match: (Mul32F x y)
	// result: (F32Mul x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Mul)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMul64_0(v *Value) bool {
	// match: (Mul64 x y)
	// result: (I64Mul x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Mul)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMul64F_0(v *Value) bool {
	// match: (Mul64F x y)
	// result: (F64Mul x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Mul)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpMul8_0(v *Value) bool {
	// match: (Mul8 x y)
	// result: (I64Mul x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Mul)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpNeg16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neg16 x)
	// result: (I64Sub (I64Const [0]) x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpNeg32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neg32 x)
	// result: (I64Sub (I64Const [0]) x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpNeg32F_0(v *Value) bool {
	// match: (Neg32F x)
	// result: (F32Neg x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF32Neg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpNeg64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neg64 x)
	// result: (I64Sub (I64Const [0]) x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpNeg64F_0(v *Value) bool {
	// match: (Neg64F x)
	// result: (F64Neg x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Neg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpNeg8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neg8 x)
	// result: (I64Sub (I64Const [0]) x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpNeq16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neq16 x y)
	// result: (I64Ne (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Ne)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpNeq32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neq32 x y)
	// result: (I64Ne (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Ne)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpNeq32F_0(v *Value) bool {
	// match: (Neq32F x y)
	// result: (F32Ne x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Ne)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpNeq64_0(v *Value) bool {
	// match: (Neq64 x y)
	// result: (I64Ne x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Ne)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpNeq64F_0(v *Value) bool {
	// match: (Neq64F x y)
	// result: (F64Ne x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Ne)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpNeq8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Neq8 x y)
	// result: (I64Ne (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Ne)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpNeqB_0(v *Value) bool {
	// match: (NeqB x y)
	// result: (I64Ne x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Ne)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpNeqPtr_0(v *Value) bool {
	// match: (NeqPtr x y)
	// result: (I64Ne x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Ne)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpNilCheck_0(v *Value) bool {
	// match: (NilCheck ptr mem)
	// result: (LoweredNilCheck ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(OpWasmLoweredNilCheck)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValueWasm_OpNot_0(v *Value) bool {
	// match: (Not x)
	// result: (I64Eqz x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Eqz)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpOffPtr_0(v *Value) bool {
	// match: (OffPtr [off] ptr)
	// result: (I64AddConst [off] ptr)
	for {
		off := v.AuxInt
		ptr := v.Args[0]
		v.reset(OpWasmI64AddConst)
		v.AuxInt = off
		v.AddArg(ptr)
		return true
	}
}
func rewriteValueWasm_OpOr16_0(v *Value) bool {
	// match: (Or16 x y)
	// result: (I64Or x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Or)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpOr32_0(v *Value) bool {
	// match: (Or32 x y)
	// result: (I64Or x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Or)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpOr64_0(v *Value) bool {
	// match: (Or64 x y)
	// result: (I64Or x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Or)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpOr8_0(v *Value) bool {
	// match: (Or8 x y)
	// result: (I64Or x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Or)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpOrB_0(v *Value) bool {
	// match: (OrB x y)
	// result: (I64Or x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Or)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpPopCount16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (PopCount16 x)
	// result: (I64Popcnt (ZeroExt16to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Popcnt)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpPopCount32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (PopCount32 x)
	// result: (I64Popcnt (ZeroExt32to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Popcnt)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpPopCount64_0(v *Value) bool {
	// match: (PopCount64 x)
	// result: (I64Popcnt x)
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Popcnt)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpPopCount8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (PopCount8 x)
	// result: (I64Popcnt (ZeroExt8to64 x))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64Popcnt)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRotateLeft16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (RotateLeft16 <t> x (I64Const [c]))
	// result: (Or16 (Lsh16x64 <t> x (I64Const [c&15])) (Rsh16Ux64 <t> x (I64Const [-c&15])))
	for {
		t := v.Type
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		v.reset(OpOr16)
		v0 := b.NewValue0(v.Pos, OpLsh16x64, t)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = c & 15
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpRsh16Ux64, t)
		v2.AddArg(x)
		v3 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v3.AuxInt = -c & 15
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
	return false
}
func rewriteValueWasm_OpRotateLeft32_0(v *Value) bool {
	// match: (RotateLeft32 x y)
	// result: (I32Rotl x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI32Rotl)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRotateLeft64_0(v *Value) bool {
	// match: (RotateLeft64 x y)
	// result: (I64Rotl x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Rotl)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRotateLeft8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (RotateLeft8 <t> x (I64Const [c]))
	// result: (Or8 (Lsh8x64 <t> x (I64Const [c&7])) (Rsh8Ux64 <t> x (I64Const [-c&7])))
	for {
		t := v.Type
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		v.reset(OpOr8)
		v0 := b.NewValue0(v.Pos, OpLsh8x64, t)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = c & 7
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpRsh8Ux64, t)
		v2.AddArg(x)
		v3 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v3.AuxInt = -c & 7
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
	return false
}
func rewriteValueWasm_OpRound32F_0(v *Value) bool {
	// match: (Round32F x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpRound64F_0(v *Value) bool {
	// match: (Round64F x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpRoundToEven_0(v *Value) bool {
	// match: (RoundToEven x)
	// result: (F64Nearest x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Nearest)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpRsh16Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux16 x y)
	// result: (Rsh64Ux64 (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh16Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux32 x y)
	// result: (Rsh64Ux64 (ZeroExt16to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh16Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux64 x y)
	// result: (Rsh64Ux64 (ZeroExt16to64 x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRsh16Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16Ux8 x y)
	// result: (Rsh64Ux64 (ZeroExt16to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh16x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x16 x y)
	// result: (Rsh64x64 (SignExt16to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh16x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x32 x y)
	// result: (Rsh64x64 (SignExt16to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh16x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x64 x y)
	// result: (Rsh64x64 (SignExt16to64 x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRsh16x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh16x8 x y)
	// result: (Rsh64x64 (SignExt16to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt16to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh32Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux16 x y)
	// result: (Rsh64Ux64 (ZeroExt32to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh32Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux32 x y)
	// result: (Rsh64Ux64 (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh32Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux64 x y)
	// result: (Rsh64Ux64 (ZeroExt32to64 x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRsh32Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32Ux8 x y)
	// result: (Rsh64Ux64 (ZeroExt32to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh32x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x16 x y)
	// result: (Rsh64x64 (SignExt32to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh32x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x32 x y)
	// result: (Rsh64x64 (SignExt32to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh32x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x64 x y)
	// result: (Rsh64x64 (SignExt32to64 x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRsh32x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh32x8 x y)
	// result: (Rsh64x64 (SignExt32to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt32to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh64Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux16 x y)
	// result: (Rsh64Ux64 x (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh64Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux32 x y)
	// result: (Rsh64Ux64 x (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh64Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux64 x y)
	// cond: shiftIsBounded(v)
	// result: (I64ShrU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpWasmI64ShrU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64Ux64 x (I64Const [c]))
	// cond: uint64(c) < 64
	// result: (I64ShrU x (I64Const [c]))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpWasmI64ShrU)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = c
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64Ux64 x (I64Const [c]))
	// cond: uint64(c) >= 64
	// result: (I64Const [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 64) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh64Ux64 x y)
	// result: (Select (I64ShrU x y) (I64Const [0]) (I64LtU y (I64Const [64])))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmSelect)
		v0 := b.NewValue0(v.Pos, OpWasmI64ShrU, typ.Int64)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 0
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpWasmI64LtU, typ.Bool)
		v2.AddArg(y)
		v3 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v3.AuxInt = 64
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpRsh64Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64Ux8 x y)
	// result: (Rsh64Ux64 x (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh64x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x16 x y)
	// result: (Rsh64x64 x (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh64x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x32 x y)
	// result: (Rsh64x64 x (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh64x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x64 x y)
	// cond: shiftIsBounded(v)
	// result: (I64ShrS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		if !(shiftIsBounded(v)) {
			break
		}
		v.reset(OpWasmI64ShrS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (Rsh64x64 x (I64Const [c]))
	// cond: uint64(c) < 64
	// result: (I64ShrS x (I64Const [c]))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) < 64) {
			break
		}
		v.reset(OpWasmI64ShrS)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = c
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x (I64Const [c]))
	// cond: uint64(c) >= 64
	// result: (I64ShrS x (I64Const [63]))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		c := v_1.AuxInt
		if !(uint64(c) >= 64) {
			break
		}
		v.reset(OpWasmI64ShrS)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 63
		v.AddArg(v0)
		return true
	}
	// match: (Rsh64x64 x y)
	// result: (I64ShrS x (Select <typ.Int64> y (I64Const [63]) (I64LtU y (I64Const [64]))))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmSelect, typ.Int64)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 63
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpWasmI64LtU, typ.Bool)
		v2.AddArg(y)
		v3 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v3.AuxInt = 64
		v2.AddArg(v3)
		v0.AddArg(v2)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh64x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh64x8 x y)
	// result: (Rsh64x64 x (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpRsh8Ux16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux16 x y)
	// result: (Rsh64Ux64 (ZeroExt8to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh8Ux32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux32 x y)
	// result: (Rsh64Ux64 (ZeroExt8to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh8Ux64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux64 x y)
	// result: (Rsh64Ux64 (ZeroExt8to64 x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRsh8Ux8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8Ux8 x y)
	// result: (Rsh64Ux64 (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64Ux64)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh8x16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x16 x y)
	// result: (Rsh64x64 (SignExt8to64 x) (ZeroExt16to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt16to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh8x32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x32 x y)
	// result: (Rsh64x64 (SignExt8to64 x) (ZeroExt32to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt32to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpRsh8x64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x64 x y)
	// result: (Rsh64x64 (SignExt8to64 x) y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpRsh8x8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Rsh8x8 x y)
	// result: (Rsh64x64 (SignExt8to64 x) (ZeroExt8to64 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpRsh64x64)
		v0 := b.NewValue0(v.Pos, OpSignExt8to64, typ.Int64)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to64, typ.UInt64)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueWasm_OpSignExt16to32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (SignExt16to32 x:(I64Load16S _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load16S {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SignExt16to32 x)
	// cond: objabi.GOWASM.SignExt
	// result: (I64Extend16S x)
	for {
		x := v.Args[0]
		if !(objabi.GOWASM.SignExt) {
			break
		}
		v.reset(OpWasmI64Extend16S)
		v.AddArg(x)
		return true
	}
	// match: (SignExt16to32 x)
	// result: (I64ShrS (I64Shl x (I64Const [48])) (I64Const [48]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 48
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 48
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSignExt16to64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (SignExt16to64 x:(I64Load16S _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load16S {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SignExt16to64 x)
	// cond: objabi.GOWASM.SignExt
	// result: (I64Extend16S x)
	for {
		x := v.Args[0]
		if !(objabi.GOWASM.SignExt) {
			break
		}
		v.reset(OpWasmI64Extend16S)
		v.AddArg(x)
		return true
	}
	// match: (SignExt16to64 x)
	// result: (I64ShrS (I64Shl x (I64Const [48])) (I64Const [48]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 48
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 48
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSignExt32to64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (SignExt32to64 x:(I64Load32S _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load32S {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SignExt32to64 x)
	// cond: objabi.GOWASM.SignExt
	// result: (I64Extend32S x)
	for {
		x := v.Args[0]
		if !(objabi.GOWASM.SignExt) {
			break
		}
		v.reset(OpWasmI64Extend32S)
		v.AddArg(x)
		return true
	}
	// match: (SignExt32to64 x)
	// result: (I64ShrS (I64Shl x (I64Const [32])) (I64Const [32]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 32
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 32
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSignExt8to16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (SignExt8to16 x:(I64Load8S _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load8S {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SignExt8to16 x)
	// cond: objabi.GOWASM.SignExt
	// result: (I64Extend8S x)
	for {
		x := v.Args[0]
		if !(objabi.GOWASM.SignExt) {
			break
		}
		v.reset(OpWasmI64Extend8S)
		v.AddArg(x)
		return true
	}
	// match: (SignExt8to16 x)
	// result: (I64ShrS (I64Shl x (I64Const [56])) (I64Const [56]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 56
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 56
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSignExt8to32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (SignExt8to32 x:(I64Load8S _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load8S {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SignExt8to32 x)
	// cond: objabi.GOWASM.SignExt
	// result: (I64Extend8S x)
	for {
		x := v.Args[0]
		if !(objabi.GOWASM.SignExt) {
			break
		}
		v.reset(OpWasmI64Extend8S)
		v.AddArg(x)
		return true
	}
	// match: (SignExt8to32 x)
	// result: (I64ShrS (I64Shl x (I64Const [56])) (I64Const [56]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 56
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 56
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSignExt8to64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (SignExt8to64 x:(I64Load8S _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load8S {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SignExt8to64 x)
	// cond: objabi.GOWASM.SignExt
	// result: (I64Extend8S x)
	for {
		x := v.Args[0]
		if !(objabi.GOWASM.SignExt) {
			break
		}
		v.reset(OpWasmI64Extend8S)
		v.AddArg(x)
		return true
	}
	// match: (SignExt8to64 x)
	// result: (I64ShrS (I64Shl x (I64Const [56])) (I64Const [56]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Shl, typ.Int64)
		v0.AddArg(x)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 56
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 56
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSlicemask_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Slicemask x)
	// result: (I64ShrS (I64Sub (I64Const [0]) x) (I64Const [63]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64ShrS)
		v0 := b.NewValue0(v.Pos, OpWasmI64Sub, typ.Int64)
		v1 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v1.AuxInt = 0
		v0.AddArg(v1)
		v0.AddArg(x)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 63
		v.AddArg(v2)
		return true
	}
}
func rewriteValueWasm_OpSqrt_0(v *Value) bool {
	// match: (Sqrt x)
	// result: (F64Sqrt x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Sqrt)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpStaticCall_0(v *Value) bool {
	// match: (StaticCall [argwid] {target} mem)
	// result: (LoweredStaticCall [argwid] {target} mem)
	for {
		argwid := v.AuxInt
		target := v.Aux
		mem := v.Args[0]
		v.reset(OpWasmLoweredStaticCall)
		v.AuxInt = argwid
		v.Aux = target
		v.AddArg(mem)
		return true
	}
}
func rewriteValueWasm_OpStore_0(v *Value) bool {
	// match: (Store {t} ptr val mem)
	// cond: is64BitFloat(t.(*types.Type))
	// result: (F64Store ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(is64BitFloat(t.(*types.Type))) {
			break
		}
		v.reset(OpWasmF64Store)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: is32BitFloat(t.(*types.Type))
	// result: (F32Store ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(is32BitFloat(t.(*types.Type))) {
			break
		}
		v.reset(OpWasmF32Store)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 8
	// result: (I64Store ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 8) {
			break
		}
		v.reset(OpWasmI64Store)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 4
	// result: (I64Store32 ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 4) {
			break
		}
		v.reset(OpWasmI64Store32)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 2
	// result: (I64Store16 ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 2) {
			break
		}
		v.reset(OpWasmI64Store16)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 1
	// result: (I64Store8 ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 1) {
			break
		}
		v.reset(OpWasmI64Store8)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpSub16_0(v *Value) bool {
	// match: (Sub16 x y)
	// result: (I64Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpSub32_0(v *Value) bool {
	// match: (Sub32 x y)
	// result: (I64Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpSub32F_0(v *Value) bool {
	// match: (Sub32F x y)
	// result: (F32Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF32Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpSub64_0(v *Value) bool {
	// match: (Sub64 x y)
	// result: (I64Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpSub64F_0(v *Value) bool {
	// match: (Sub64F x y)
	// result: (F64Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmF64Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpSub8_0(v *Value) bool {
	// match: (Sub8 x y)
	// result: (I64Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpSubPtr_0(v *Value) bool {
	// match: (SubPtr x y)
	// result: (I64Sub x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Sub)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpTrunc_0(v *Value) bool {
	// match: (Trunc x)
	// result: (F64Trunc x)
	for {
		x := v.Args[0]
		v.reset(OpWasmF64Trunc)
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpTrunc16to8_0(v *Value) bool {
	// match: (Trunc16to8 x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpTrunc32to16_0(v *Value) bool {
	// match: (Trunc32to16 x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpTrunc32to8_0(v *Value) bool {
	// match: (Trunc32to8 x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpTrunc64to16_0(v *Value) bool {
	// match: (Trunc64to16 x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpTrunc64to32_0(v *Value) bool {
	// match: (Trunc64to32 x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpTrunc64to8_0(v *Value) bool {
	// match: (Trunc64to8 x)
	// result: x
	for {
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
}
func rewriteValueWasm_OpWB_0(v *Value) bool {
	// match: (WB {fn} destptr srcptr mem)
	// result: (LoweredWB {fn} destptr srcptr mem)
	for {
		fn := v.Aux
		mem := v.Args[2]
		destptr := v.Args[0]
		srcptr := v.Args[1]
		v.reset(OpWasmLoweredWB)
		v.Aux = fn
		v.AddArg(destptr)
		v.AddArg(srcptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValueWasm_OpWasmF64Add_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (F64Add (F64Const [x]) (F64Const [y]))
	// result: (F64Const [auxFrom64F(auxTo64F(x) + auxTo64F(y))])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmF64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmF64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmF64Const)
		v.AuxInt = auxFrom64F(auxTo64F(x) + auxTo64F(y))
		return true
	}
	// match: (F64Add (F64Const [x]) y)
	// result: (F64Add y (F64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmF64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmF64Add)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmF64Const, typ.Float64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmF64Mul_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (F64Mul (F64Const [x]) (F64Const [y]))
	// result: (F64Const [auxFrom64F(auxTo64F(x) * auxTo64F(y))])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmF64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmF64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmF64Const)
		v.AuxInt = auxFrom64F(auxTo64F(x) * auxTo64F(y))
		return true
	}
	// match: (F64Mul (F64Const [x]) y)
	// result: (F64Mul y (F64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmF64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmF64Mul)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmF64Const, typ.Float64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Add_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64Add (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x + y])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x + y
		return true
	}
	// match: (I64Add (I64Const [x]) y)
	// result: (I64Add y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64Add)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	// match: (I64Add x (I64Const [y]))
	// result: (I64AddConst [y] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64AddConst)
		v.AuxInt = y
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64AddConst_0(v *Value) bool {
	// match: (I64AddConst [0] x)
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
	// match: (I64AddConst [off] (LoweredAddr {sym} [off2] base))
	// cond: isU32Bit(off+off2)
	// result: (LoweredAddr {sym} [off+off2] base)
	for {
		off := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != OpWasmLoweredAddr {
			break
		}
		off2 := v_0.AuxInt
		sym := v_0.Aux
		base := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmLoweredAddr)
		v.AuxInt = off + off2
		v.Aux = sym
		v.AddArg(base)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64And_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64And (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x & y])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x & y
		return true
	}
	// match: (I64And (I64Const [x]) y)
	// result: (I64And y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64And)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Eq_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64Eq (I64Const [x]) (I64Const [y]))
	// cond: x == y
	// result: (I64Const [1])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		if !(x == y) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = 1
		return true
	}
	// match: (I64Eq (I64Const [x]) (I64Const [y]))
	// cond: x != y
	// result: (I64Const [0])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		if !(x != y) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = 0
		return true
	}
	// match: (I64Eq (I64Const [x]) y)
	// result: (I64Eq y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64Eq)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	// match: (I64Eq x (I64Const [0]))
	// result: (I64Eqz x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const || v_1.AuxInt != 0 {
			break
		}
		v.reset(OpWasmI64Eqz)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Eqz_0(v *Value) bool {
	// match: (I64Eqz (I64Eqz (I64Eqz x)))
	// result: (I64Eqz x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Eqz {
			break
		}
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpWasmI64Eqz {
			break
		}
		x := v_0_0.Args[0]
		v.reset(OpWasmI64Eqz)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (I64Load [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (I64Load [off] (LoweredAddr {sym} [off2] (SB)) _)
	// cond: symIsRO(sym) && isU32Bit(off+off2)
	// result: (I64Const [int64(read64(sym, off+off2, config.BigEndian))])
	for {
		off := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmLoweredAddr {
			break
		}
		off2 := v_0.AuxInt
		sym := v_0.Aux
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpSB || !(symIsRO(sym) && isU32Bit(off+off2)) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = int64(read64(sym, off+off2, config.BigEndian))
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load16S_0(v *Value) bool {
	// match: (I64Load16S [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load16S [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load16S)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load16U_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (I64Load16U [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load16U [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load16U)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (I64Load16U [off] (LoweredAddr {sym} [off2] (SB)) _)
	// cond: symIsRO(sym) && isU32Bit(off+off2)
	// result: (I64Const [int64(read16(sym, off+off2, config.BigEndian))])
	for {
		off := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmLoweredAddr {
			break
		}
		off2 := v_0.AuxInt
		sym := v_0.Aux
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpSB || !(symIsRO(sym) && isU32Bit(off+off2)) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = int64(read16(sym, off+off2, config.BigEndian))
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load32S_0(v *Value) bool {
	// match: (I64Load32S [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load32S [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load32S)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load32U_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (I64Load32U [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load32U [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load32U)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (I64Load32U [off] (LoweredAddr {sym} [off2] (SB)) _)
	// cond: symIsRO(sym) && isU32Bit(off+off2)
	// result: (I64Const [int64(read32(sym, off+off2, config.BigEndian))])
	for {
		off := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmLoweredAddr {
			break
		}
		off2 := v_0.AuxInt
		sym := v_0.Aux
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpSB || !(symIsRO(sym) && isU32Bit(off+off2)) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = int64(read32(sym, off+off2, config.BigEndian))
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load8S_0(v *Value) bool {
	// match: (I64Load8S [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load8S [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load8S)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Load8U_0(v *Value) bool {
	// match: (I64Load8U [off] (I64AddConst [off2] ptr) mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Load8U [off+off2] ptr mem)
	for {
		off := v.AuxInt
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Load8U)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (I64Load8U [off] (LoweredAddr {sym} [off2] (SB)) _)
	// cond: symIsRO(sym) && isU32Bit(off+off2)
	// result: (I64Const [int64(read8(sym, off+off2))])
	for {
		off := v.AuxInt
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmLoweredAddr {
			break
		}
		off2 := v_0.AuxInt
		sym := v_0.Aux
		v_0_0 := v_0.Args[0]
		if v_0_0.Op != OpSB || !(symIsRO(sym) && isU32Bit(off+off2)) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = int64(read8(sym, off+off2))
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Mul_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64Mul (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x * y])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x * y
		return true
	}
	// match: (I64Mul (I64Const [x]) y)
	// result: (I64Mul y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64Mul)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Ne_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64Ne (I64Const [x]) (I64Const [y]))
	// cond: x == y
	// result: (I64Const [0])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		if !(x == y) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = 0
		return true
	}
	// match: (I64Ne (I64Const [x]) (I64Const [y]))
	// cond: x != y
	// result: (I64Const [1])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		if !(x != y) {
			break
		}
		v.reset(OpWasmI64Const)
		v.AuxInt = 1
		return true
	}
	// match: (I64Ne (I64Const [x]) y)
	// result: (I64Ne y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64Ne)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	// match: (I64Ne x (I64Const [0]))
	// result: (I64Eqz (I64Eqz x))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const || v_1.AuxInt != 0 {
			break
		}
		v.reset(OpWasmI64Eqz)
		v0 := b.NewValue0(v.Pos, OpWasmI64Eqz, typ.Bool)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Or_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64Or (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x | y])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x | y
		return true
	}
	// match: (I64Or (I64Const [x]) y)
	// result: (I64Or y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64Or)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Shl_0(v *Value) bool {
	// match: (I64Shl (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x << uint64(y)])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x << uint64(y)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64ShrS_0(v *Value) bool {
	// match: (I64ShrS (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x >> uint64(y)])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x >> uint64(y)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64ShrU_0(v *Value) bool {
	// match: (I64ShrU (I64Const [x]) (I64Const [y]))
	// result: (I64Const [int64(uint64(x) >> uint64(y))])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = int64(uint64(x) >> uint64(y))
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Store_0(v *Value) bool {
	// match: (I64Store [off] (I64AddConst [off2] ptr) val mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Store [off+off2] ptr val mem)
	for {
		off := v.AuxInt
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Store)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Store16_0(v *Value) bool {
	// match: (I64Store16 [off] (I64AddConst [off2] ptr) val mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Store16 [off+off2] ptr val mem)
	for {
		off := v.AuxInt
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Store16)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Store32_0(v *Value) bool {
	// match: (I64Store32 [off] (I64AddConst [off2] ptr) val mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Store32 [off+off2] ptr val mem)
	for {
		off := v.AuxInt
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Store32)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Store8_0(v *Value) bool {
	// match: (I64Store8 [off] (I64AddConst [off2] ptr) val mem)
	// cond: isU32Bit(off+off2)
	// result: (I64Store8 [off+off2] ptr val mem)
	for {
		off := v.AuxInt
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64AddConst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(isU32Bit(off + off2)) {
			break
		}
		v.reset(OpWasmI64Store8)
		v.AuxInt = off + off2
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpWasmI64Xor_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (I64Xor (I64Const [x]) (I64Const [y]))
	// result: (I64Const [x ^ y])
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v_1 := v.Args[1]
		if v_1.Op != OpWasmI64Const {
			break
		}
		y := v_1.AuxInt
		v.reset(OpWasmI64Const)
		v.AuxInt = x ^ y
		return true
	}
	// match: (I64Xor (I64Const [x]) y)
	// result: (I64Xor y (I64Const [x]))
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpWasmI64Const {
			break
		}
		x := v_0.AuxInt
		v.reset(OpWasmI64Xor)
		v.AddArg(y)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = x
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValueWasm_OpXor16_0(v *Value) bool {
	// match: (Xor16 x y)
	// result: (I64Xor x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpXor32_0(v *Value) bool {
	// match: (Xor32 x y)
	// result: (I64Xor x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpXor64_0(v *Value) bool {
	// match: (Xor64 x y)
	// result: (I64Xor x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpXor8_0(v *Value) bool {
	// match: (Xor8 x y)
	// result: (I64Xor x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(OpWasmI64Xor)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueWasm_OpZero_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
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
	// result: (I64Store8 destptr (I64Const [0]) mem)
	for {
		if v.AuxInt != 1 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store8)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [2] destptr mem)
	// result: (I64Store16 destptr (I64Const [0]) mem)
	for {
		if v.AuxInt != 2 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store16)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [4] destptr mem)
	// result: (I64Store32 destptr (I64Const [0]) mem)
	for {
		if v.AuxInt != 4 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store32)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [8] destptr mem)
	// result: (I64Store destptr (I64Const [0]) mem)
	for {
		if v.AuxInt != 8 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [3] destptr mem)
	// result: (I64Store8 [2] destptr (I64Const [0]) (I64Store16 destptr (I64Const [0]) mem))
	for {
		if v.AuxInt != 3 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store8)
		v.AuxInt = 2
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store16, types.TypeMem)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [5] destptr mem)
	// result: (I64Store8 [4] destptr (I64Const [0]) (I64Store32 destptr (I64Const [0]) mem))
	for {
		if v.AuxInt != 5 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store8)
		v.AuxInt = 4
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store32, types.TypeMem)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [6] destptr mem)
	// result: (I64Store16 [4] destptr (I64Const [0]) (I64Store32 destptr (I64Const [0]) mem))
	for {
		if v.AuxInt != 6 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store16)
		v.AuxInt = 4
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store32, types.TypeMem)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [7] destptr mem)
	// result: (I64Store32 [3] destptr (I64Const [0]) (I64Store32 destptr (I64Const [0]) mem))
	for {
		if v.AuxInt != 7 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store32)
		v.AuxInt = 3
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store32, types.TypeMem)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [s] destptr mem)
	// cond: s%8 != 0 && s > 8
	// result: (Zero [s-s%8] (OffPtr <destptr.Type> destptr [s%8]) (I64Store destptr (I64Const [0]) mem))
	for {
		s := v.AuxInt
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(s%8 != 0 && s > 8) {
			break
		}
		v.reset(OpZero)
		v.AuxInt = s - s%8
		v0 := b.NewValue0(v.Pos, OpOffPtr, destptr.Type)
		v0.AuxInt = s % 8
		v0.AddArg(destptr)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValueWasm_OpZero_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Zero [16] destptr mem)
	// result: (I64Store [8] destptr (I64Const [0]) (I64Store destptr (I64Const [0]) mem))
	for {
		if v.AuxInt != 16 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store)
		v.AuxInt = 8
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [24] destptr mem)
	// result: (I64Store [16] destptr (I64Const [0]) (I64Store [8] destptr (I64Const [0]) (I64Store destptr (I64Const [0]) mem)))
	for {
		if v.AuxInt != 24 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store)
		v.AuxInt = 16
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v1.AuxInt = 8
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v3 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v3.AddArg(destptr)
		v4 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v4.AuxInt = 0
		v3.AddArg(v4)
		v3.AddArg(mem)
		v1.AddArg(v3)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [32] destptr mem)
	// result: (I64Store [24] destptr (I64Const [0]) (I64Store [16] destptr (I64Const [0]) (I64Store [8] destptr (I64Const [0]) (I64Store destptr (I64Const [0]) mem))))
	for {
		if v.AuxInt != 32 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(OpWasmI64Store)
		v.AuxInt = 24
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v1.AuxInt = 16
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v2.AuxInt = 0
		v1.AddArg(v2)
		v3 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v3.AuxInt = 8
		v3.AddArg(destptr)
		v4 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v4.AuxInt = 0
		v3.AddArg(v4)
		v5 := b.NewValue0(v.Pos, OpWasmI64Store, types.TypeMem)
		v5.AddArg(destptr)
		v6 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v6.AuxInt = 0
		v5.AddArg(v6)
		v5.AddArg(mem)
		v3.AddArg(v5)
		v1.AddArg(v3)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [s] destptr mem)
	// cond: s%8 == 0 && s > 32
	// result: (LoweredZero [s/8] destptr mem)
	for {
		s := v.AuxInt
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(s%8 == 0 && s > 32) {
			break
		}
		v.reset(OpWasmLoweredZero)
		v.AuxInt = s / 8
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValueWasm_OpZeroExt16to32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ZeroExt16to32 x:(I64Load16U _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load16U {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ZeroExt16to32 x)
	// result: (I64And x (I64Const [0xffff]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0xffff
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpZeroExt16to64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ZeroExt16to64 x:(I64Load16U _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load16U {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ZeroExt16to64 x)
	// result: (I64And x (I64Const [0xffff]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0xffff
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpZeroExt32to64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ZeroExt32to64 x:(I64Load32U _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load32U {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ZeroExt32to64 x)
	// result: (I64And x (I64Const [0xffffffff]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0xffffffff
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpZeroExt8to16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ZeroExt8to16 x:(I64Load8U _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load8U {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ZeroExt8to16 x)
	// result: (I64And x (I64Const [0xff]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0xff
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpZeroExt8to32_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ZeroExt8to32 x:(I64Load8U _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load8U {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ZeroExt8to32 x)
	// result: (I64And x (I64Const [0xff]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0xff
		v.AddArg(v0)
		return true
	}
}
func rewriteValueWasm_OpZeroExt8to64_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ZeroExt8to64 x:(I64Load8U _ _))
	// result: x
	for {
		x := v.Args[0]
		if x.Op != OpWasmI64Load8U {
			break
		}
		_ = x.Args[1]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ZeroExt8to64 x)
	// result: (I64And x (I64Const [0xff]))
	for {
		x := v.Args[0]
		v.reset(OpWasmI64And)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, OpWasmI64Const, typ.Int64)
		v0.AuxInt = 0xff
		v.AddArg(v0)
		return true
	}
}
func rewriteBlockWasm(b *Block) bool {
	switch b.Kind {
	}
	return false
}
