// Code generated from gen/386.rules; DO NOT EDIT.
// generated with: cd gen; go run *.go

package ssa

import "math"
import "cmd/compile/internal/types"

func rewriteValue386(v *Value) bool {
	switch v.Op {
	case Op386ADCL:
		return rewriteValue386_Op386ADCL_0(v)
	case Op386ADDL:
		return rewriteValue386_Op386ADDL_0(v) || rewriteValue386_Op386ADDL_10(v) || rewriteValue386_Op386ADDL_20(v)
	case Op386ADDLcarry:
		return rewriteValue386_Op386ADDLcarry_0(v)
	case Op386ADDLconst:
		return rewriteValue386_Op386ADDLconst_0(v)
	case Op386ADDLconstmodify:
		return rewriteValue386_Op386ADDLconstmodify_0(v)
	case Op386ADDLconstmodifyidx4:
		return rewriteValue386_Op386ADDLconstmodifyidx4_0(v)
	case Op386ADDLload:
		return rewriteValue386_Op386ADDLload_0(v)
	case Op386ADDLloadidx4:
		return rewriteValue386_Op386ADDLloadidx4_0(v)
	case Op386ADDLmodify:
		return rewriteValue386_Op386ADDLmodify_0(v)
	case Op386ADDLmodifyidx4:
		return rewriteValue386_Op386ADDLmodifyidx4_0(v)
	case Op386ADDSD:
		return rewriteValue386_Op386ADDSD_0(v)
	case Op386ADDSDload:
		return rewriteValue386_Op386ADDSDload_0(v)
	case Op386ADDSS:
		return rewriteValue386_Op386ADDSS_0(v)
	case Op386ADDSSload:
		return rewriteValue386_Op386ADDSSload_0(v)
	case Op386ANDL:
		return rewriteValue386_Op386ANDL_0(v)
	case Op386ANDLconst:
		return rewriteValue386_Op386ANDLconst_0(v)
	case Op386ANDLconstmodify:
		return rewriteValue386_Op386ANDLconstmodify_0(v)
	case Op386ANDLconstmodifyidx4:
		return rewriteValue386_Op386ANDLconstmodifyidx4_0(v)
	case Op386ANDLload:
		return rewriteValue386_Op386ANDLload_0(v)
	case Op386ANDLloadidx4:
		return rewriteValue386_Op386ANDLloadidx4_0(v)
	case Op386ANDLmodify:
		return rewriteValue386_Op386ANDLmodify_0(v)
	case Op386ANDLmodifyidx4:
		return rewriteValue386_Op386ANDLmodifyidx4_0(v)
	case Op386CMPB:
		return rewriteValue386_Op386CMPB_0(v)
	case Op386CMPBconst:
		return rewriteValue386_Op386CMPBconst_0(v)
	case Op386CMPBload:
		return rewriteValue386_Op386CMPBload_0(v)
	case Op386CMPL:
		return rewriteValue386_Op386CMPL_0(v)
	case Op386CMPLconst:
		return rewriteValue386_Op386CMPLconst_0(v) || rewriteValue386_Op386CMPLconst_10(v)
	case Op386CMPLload:
		return rewriteValue386_Op386CMPLload_0(v)
	case Op386CMPW:
		return rewriteValue386_Op386CMPW_0(v)
	case Op386CMPWconst:
		return rewriteValue386_Op386CMPWconst_0(v)
	case Op386CMPWload:
		return rewriteValue386_Op386CMPWload_0(v)
	case Op386DIVSD:
		return rewriteValue386_Op386DIVSD_0(v)
	case Op386DIVSDload:
		return rewriteValue386_Op386DIVSDload_0(v)
	case Op386DIVSS:
		return rewriteValue386_Op386DIVSS_0(v)
	case Op386DIVSSload:
		return rewriteValue386_Op386DIVSSload_0(v)
	case Op386LEAL:
		return rewriteValue386_Op386LEAL_0(v)
	case Op386LEAL1:
		return rewriteValue386_Op386LEAL1_0(v)
	case Op386LEAL2:
		return rewriteValue386_Op386LEAL2_0(v)
	case Op386LEAL4:
		return rewriteValue386_Op386LEAL4_0(v)
	case Op386LEAL8:
		return rewriteValue386_Op386LEAL8_0(v)
	case Op386MOVBLSX:
		return rewriteValue386_Op386MOVBLSX_0(v)
	case Op386MOVBLSXload:
		return rewriteValue386_Op386MOVBLSXload_0(v)
	case Op386MOVBLZX:
		return rewriteValue386_Op386MOVBLZX_0(v)
	case Op386MOVBload:
		return rewriteValue386_Op386MOVBload_0(v)
	case Op386MOVBloadidx1:
		return rewriteValue386_Op386MOVBloadidx1_0(v)
	case Op386MOVBstore:
		return rewriteValue386_Op386MOVBstore_0(v) || rewriteValue386_Op386MOVBstore_10(v)
	case Op386MOVBstoreconst:
		return rewriteValue386_Op386MOVBstoreconst_0(v)
	case Op386MOVBstoreconstidx1:
		return rewriteValue386_Op386MOVBstoreconstidx1_0(v)
	case Op386MOVBstoreidx1:
		return rewriteValue386_Op386MOVBstoreidx1_0(v) || rewriteValue386_Op386MOVBstoreidx1_10(v) || rewriteValue386_Op386MOVBstoreidx1_20(v)
	case Op386MOVLload:
		return rewriteValue386_Op386MOVLload_0(v)
	case Op386MOVLloadidx1:
		return rewriteValue386_Op386MOVLloadidx1_0(v)
	case Op386MOVLloadidx4:
		return rewriteValue386_Op386MOVLloadidx4_0(v)
	case Op386MOVLstore:
		return rewriteValue386_Op386MOVLstore_0(v) || rewriteValue386_Op386MOVLstore_10(v) || rewriteValue386_Op386MOVLstore_20(v)
	case Op386MOVLstoreconst:
		return rewriteValue386_Op386MOVLstoreconst_0(v)
	case Op386MOVLstoreconstidx1:
		return rewriteValue386_Op386MOVLstoreconstidx1_0(v)
	case Op386MOVLstoreconstidx4:
		return rewriteValue386_Op386MOVLstoreconstidx4_0(v)
	case Op386MOVLstoreidx1:
		return rewriteValue386_Op386MOVLstoreidx1_0(v)
	case Op386MOVLstoreidx4:
		return rewriteValue386_Op386MOVLstoreidx4_0(v) || rewriteValue386_Op386MOVLstoreidx4_10(v)
	case Op386MOVSDconst:
		return rewriteValue386_Op386MOVSDconst_0(v)
	case Op386MOVSDload:
		return rewriteValue386_Op386MOVSDload_0(v)
	case Op386MOVSDloadidx1:
		return rewriteValue386_Op386MOVSDloadidx1_0(v)
	case Op386MOVSDloadidx8:
		return rewriteValue386_Op386MOVSDloadidx8_0(v)
	case Op386MOVSDstore:
		return rewriteValue386_Op386MOVSDstore_0(v)
	case Op386MOVSDstoreidx1:
		return rewriteValue386_Op386MOVSDstoreidx1_0(v)
	case Op386MOVSDstoreidx8:
		return rewriteValue386_Op386MOVSDstoreidx8_0(v)
	case Op386MOVSSconst:
		return rewriteValue386_Op386MOVSSconst_0(v)
	case Op386MOVSSload:
		return rewriteValue386_Op386MOVSSload_0(v)
	case Op386MOVSSloadidx1:
		return rewriteValue386_Op386MOVSSloadidx1_0(v)
	case Op386MOVSSloadidx4:
		return rewriteValue386_Op386MOVSSloadidx4_0(v)
	case Op386MOVSSstore:
		return rewriteValue386_Op386MOVSSstore_0(v)
	case Op386MOVSSstoreidx1:
		return rewriteValue386_Op386MOVSSstoreidx1_0(v)
	case Op386MOVSSstoreidx4:
		return rewriteValue386_Op386MOVSSstoreidx4_0(v)
	case Op386MOVWLSX:
		return rewriteValue386_Op386MOVWLSX_0(v)
	case Op386MOVWLSXload:
		return rewriteValue386_Op386MOVWLSXload_0(v)
	case Op386MOVWLZX:
		return rewriteValue386_Op386MOVWLZX_0(v)
	case Op386MOVWload:
		return rewriteValue386_Op386MOVWload_0(v)
	case Op386MOVWloadidx1:
		return rewriteValue386_Op386MOVWloadidx1_0(v)
	case Op386MOVWloadidx2:
		return rewriteValue386_Op386MOVWloadidx2_0(v)
	case Op386MOVWstore:
		return rewriteValue386_Op386MOVWstore_0(v)
	case Op386MOVWstoreconst:
		return rewriteValue386_Op386MOVWstoreconst_0(v)
	case Op386MOVWstoreconstidx1:
		return rewriteValue386_Op386MOVWstoreconstidx1_0(v)
	case Op386MOVWstoreconstidx2:
		return rewriteValue386_Op386MOVWstoreconstidx2_0(v)
	case Op386MOVWstoreidx1:
		return rewriteValue386_Op386MOVWstoreidx1_0(v) || rewriteValue386_Op386MOVWstoreidx1_10(v)
	case Op386MOVWstoreidx2:
		return rewriteValue386_Op386MOVWstoreidx2_0(v)
	case Op386MULL:
		return rewriteValue386_Op386MULL_0(v)
	case Op386MULLconst:
		return rewriteValue386_Op386MULLconst_0(v) || rewriteValue386_Op386MULLconst_10(v) || rewriteValue386_Op386MULLconst_20(v) || rewriteValue386_Op386MULLconst_30(v)
	case Op386MULLload:
		return rewriteValue386_Op386MULLload_0(v)
	case Op386MULLloadidx4:
		return rewriteValue386_Op386MULLloadidx4_0(v)
	case Op386MULSD:
		return rewriteValue386_Op386MULSD_0(v)
	case Op386MULSDload:
		return rewriteValue386_Op386MULSDload_0(v)
	case Op386MULSS:
		return rewriteValue386_Op386MULSS_0(v)
	case Op386MULSSload:
		return rewriteValue386_Op386MULSSload_0(v)
	case Op386NEGL:
		return rewriteValue386_Op386NEGL_0(v)
	case Op386NOTL:
		return rewriteValue386_Op386NOTL_0(v)
	case Op386ORL:
		return rewriteValue386_Op386ORL_0(v) || rewriteValue386_Op386ORL_10(v) || rewriteValue386_Op386ORL_20(v) || rewriteValue386_Op386ORL_30(v) || rewriteValue386_Op386ORL_40(v) || rewriteValue386_Op386ORL_50(v)
	case Op386ORLconst:
		return rewriteValue386_Op386ORLconst_0(v)
	case Op386ORLconstmodify:
		return rewriteValue386_Op386ORLconstmodify_0(v)
	case Op386ORLconstmodifyidx4:
		return rewriteValue386_Op386ORLconstmodifyidx4_0(v)
	case Op386ORLload:
		return rewriteValue386_Op386ORLload_0(v)
	case Op386ORLloadidx4:
		return rewriteValue386_Op386ORLloadidx4_0(v)
	case Op386ORLmodify:
		return rewriteValue386_Op386ORLmodify_0(v)
	case Op386ORLmodifyidx4:
		return rewriteValue386_Op386ORLmodifyidx4_0(v)
	case Op386ROLBconst:
		return rewriteValue386_Op386ROLBconst_0(v)
	case Op386ROLLconst:
		return rewriteValue386_Op386ROLLconst_0(v)
	case Op386ROLWconst:
		return rewriteValue386_Op386ROLWconst_0(v)
	case Op386SARB:
		return rewriteValue386_Op386SARB_0(v)
	case Op386SARBconst:
		return rewriteValue386_Op386SARBconst_0(v)
	case Op386SARL:
		return rewriteValue386_Op386SARL_0(v)
	case Op386SARLconst:
		return rewriteValue386_Op386SARLconst_0(v)
	case Op386SARW:
		return rewriteValue386_Op386SARW_0(v)
	case Op386SARWconst:
		return rewriteValue386_Op386SARWconst_0(v)
	case Op386SBBL:
		return rewriteValue386_Op386SBBL_0(v)
	case Op386SBBLcarrymask:
		return rewriteValue386_Op386SBBLcarrymask_0(v)
	case Op386SETA:
		return rewriteValue386_Op386SETA_0(v)
	case Op386SETAE:
		return rewriteValue386_Op386SETAE_0(v)
	case Op386SETB:
		return rewriteValue386_Op386SETB_0(v)
	case Op386SETBE:
		return rewriteValue386_Op386SETBE_0(v)
	case Op386SETEQ:
		return rewriteValue386_Op386SETEQ_0(v)
	case Op386SETG:
		return rewriteValue386_Op386SETG_0(v)
	case Op386SETGE:
		return rewriteValue386_Op386SETGE_0(v)
	case Op386SETL:
		return rewriteValue386_Op386SETL_0(v)
	case Op386SETLE:
		return rewriteValue386_Op386SETLE_0(v)
	case Op386SETNE:
		return rewriteValue386_Op386SETNE_0(v)
	case Op386SHLL:
		return rewriteValue386_Op386SHLL_0(v)
	case Op386SHLLconst:
		return rewriteValue386_Op386SHLLconst_0(v)
	case Op386SHRB:
		return rewriteValue386_Op386SHRB_0(v)
	case Op386SHRBconst:
		return rewriteValue386_Op386SHRBconst_0(v)
	case Op386SHRL:
		return rewriteValue386_Op386SHRL_0(v)
	case Op386SHRLconst:
		return rewriteValue386_Op386SHRLconst_0(v)
	case Op386SHRW:
		return rewriteValue386_Op386SHRW_0(v)
	case Op386SHRWconst:
		return rewriteValue386_Op386SHRWconst_0(v)
	case Op386SUBL:
		return rewriteValue386_Op386SUBL_0(v)
	case Op386SUBLcarry:
		return rewriteValue386_Op386SUBLcarry_0(v)
	case Op386SUBLconst:
		return rewriteValue386_Op386SUBLconst_0(v)
	case Op386SUBLload:
		return rewriteValue386_Op386SUBLload_0(v)
	case Op386SUBLloadidx4:
		return rewriteValue386_Op386SUBLloadidx4_0(v)
	case Op386SUBLmodify:
		return rewriteValue386_Op386SUBLmodify_0(v)
	case Op386SUBLmodifyidx4:
		return rewriteValue386_Op386SUBLmodifyidx4_0(v)
	case Op386SUBSD:
		return rewriteValue386_Op386SUBSD_0(v)
	case Op386SUBSDload:
		return rewriteValue386_Op386SUBSDload_0(v)
	case Op386SUBSS:
		return rewriteValue386_Op386SUBSS_0(v)
	case Op386SUBSSload:
		return rewriteValue386_Op386SUBSSload_0(v)
	case Op386XORL:
		return rewriteValue386_Op386XORL_0(v) || rewriteValue386_Op386XORL_10(v)
	case Op386XORLconst:
		return rewriteValue386_Op386XORLconst_0(v)
	case Op386XORLconstmodify:
		return rewriteValue386_Op386XORLconstmodify_0(v)
	case Op386XORLconstmodifyidx4:
		return rewriteValue386_Op386XORLconstmodifyidx4_0(v)
	case Op386XORLload:
		return rewriteValue386_Op386XORLload_0(v)
	case Op386XORLloadidx4:
		return rewriteValue386_Op386XORLloadidx4_0(v)
	case Op386XORLmodify:
		return rewriteValue386_Op386XORLmodify_0(v)
	case Op386XORLmodifyidx4:
		return rewriteValue386_Op386XORLmodifyidx4_0(v)
	case OpAdd16:
		return rewriteValue386_OpAdd16_0(v)
	case OpAdd32:
		return rewriteValue386_OpAdd32_0(v)
	case OpAdd32F:
		return rewriteValue386_OpAdd32F_0(v)
	case OpAdd32carry:
		return rewriteValue386_OpAdd32carry_0(v)
	case OpAdd32withcarry:
		return rewriteValue386_OpAdd32withcarry_0(v)
	case OpAdd64F:
		return rewriteValue386_OpAdd64F_0(v)
	case OpAdd8:
		return rewriteValue386_OpAdd8_0(v)
	case OpAddPtr:
		return rewriteValue386_OpAddPtr_0(v)
	case OpAddr:
		return rewriteValue386_OpAddr_0(v)
	case OpAnd16:
		return rewriteValue386_OpAnd16_0(v)
	case OpAnd32:
		return rewriteValue386_OpAnd32_0(v)
	case OpAnd8:
		return rewriteValue386_OpAnd8_0(v)
	case OpAndB:
		return rewriteValue386_OpAndB_0(v)
	case OpAvg32u:
		return rewriteValue386_OpAvg32u_0(v)
	case OpBswap32:
		return rewriteValue386_OpBswap32_0(v)
	case OpClosureCall:
		return rewriteValue386_OpClosureCall_0(v)
	case OpCom16:
		return rewriteValue386_OpCom16_0(v)
	case OpCom32:
		return rewriteValue386_OpCom32_0(v)
	case OpCom8:
		return rewriteValue386_OpCom8_0(v)
	case OpConst16:
		return rewriteValue386_OpConst16_0(v)
	case OpConst32:
		return rewriteValue386_OpConst32_0(v)
	case OpConst32F:
		return rewriteValue386_OpConst32F_0(v)
	case OpConst64F:
		return rewriteValue386_OpConst64F_0(v)
	case OpConst8:
		return rewriteValue386_OpConst8_0(v)
	case OpConstBool:
		return rewriteValue386_OpConstBool_0(v)
	case OpConstNil:
		return rewriteValue386_OpConstNil_0(v)
	case OpCtz16:
		return rewriteValue386_OpCtz16_0(v)
	case OpCtz16NonZero:
		return rewriteValue386_OpCtz16NonZero_0(v)
	case OpCvt32Fto32:
		return rewriteValue386_OpCvt32Fto32_0(v)
	case OpCvt32Fto64F:
		return rewriteValue386_OpCvt32Fto64F_0(v)
	case OpCvt32to32F:
		return rewriteValue386_OpCvt32to32F_0(v)
	case OpCvt32to64F:
		return rewriteValue386_OpCvt32to64F_0(v)
	case OpCvt64Fto32:
		return rewriteValue386_OpCvt64Fto32_0(v)
	case OpCvt64Fto32F:
		return rewriteValue386_OpCvt64Fto32F_0(v)
	case OpDiv16:
		return rewriteValue386_OpDiv16_0(v)
	case OpDiv16u:
		return rewriteValue386_OpDiv16u_0(v)
	case OpDiv32:
		return rewriteValue386_OpDiv32_0(v)
	case OpDiv32F:
		return rewriteValue386_OpDiv32F_0(v)
	case OpDiv32u:
		return rewriteValue386_OpDiv32u_0(v)
	case OpDiv64F:
		return rewriteValue386_OpDiv64F_0(v)
	case OpDiv8:
		return rewriteValue386_OpDiv8_0(v)
	case OpDiv8u:
		return rewriteValue386_OpDiv8u_0(v)
	case OpEq16:
		return rewriteValue386_OpEq16_0(v)
	case OpEq32:
		return rewriteValue386_OpEq32_0(v)
	case OpEq32F:
		return rewriteValue386_OpEq32F_0(v)
	case OpEq64F:
		return rewriteValue386_OpEq64F_0(v)
	case OpEq8:
		return rewriteValue386_OpEq8_0(v)
	case OpEqB:
		return rewriteValue386_OpEqB_0(v)
	case OpEqPtr:
		return rewriteValue386_OpEqPtr_0(v)
	case OpGeq16:
		return rewriteValue386_OpGeq16_0(v)
	case OpGeq16U:
		return rewriteValue386_OpGeq16U_0(v)
	case OpGeq32:
		return rewriteValue386_OpGeq32_0(v)
	case OpGeq32F:
		return rewriteValue386_OpGeq32F_0(v)
	case OpGeq32U:
		return rewriteValue386_OpGeq32U_0(v)
	case OpGeq64F:
		return rewriteValue386_OpGeq64F_0(v)
	case OpGeq8:
		return rewriteValue386_OpGeq8_0(v)
	case OpGeq8U:
		return rewriteValue386_OpGeq8U_0(v)
	case OpGetCallerPC:
		return rewriteValue386_OpGetCallerPC_0(v)
	case OpGetCallerSP:
		return rewriteValue386_OpGetCallerSP_0(v)
	case OpGetClosurePtr:
		return rewriteValue386_OpGetClosurePtr_0(v)
	case OpGetG:
		return rewriteValue386_OpGetG_0(v)
	case OpGreater16:
		return rewriteValue386_OpGreater16_0(v)
	case OpGreater16U:
		return rewriteValue386_OpGreater16U_0(v)
	case OpGreater32:
		return rewriteValue386_OpGreater32_0(v)
	case OpGreater32F:
		return rewriteValue386_OpGreater32F_0(v)
	case OpGreater32U:
		return rewriteValue386_OpGreater32U_0(v)
	case OpGreater64F:
		return rewriteValue386_OpGreater64F_0(v)
	case OpGreater8:
		return rewriteValue386_OpGreater8_0(v)
	case OpGreater8U:
		return rewriteValue386_OpGreater8U_0(v)
	case OpHmul32:
		return rewriteValue386_OpHmul32_0(v)
	case OpHmul32u:
		return rewriteValue386_OpHmul32u_0(v)
	case OpInterCall:
		return rewriteValue386_OpInterCall_0(v)
	case OpIsInBounds:
		return rewriteValue386_OpIsInBounds_0(v)
	case OpIsNonNil:
		return rewriteValue386_OpIsNonNil_0(v)
	case OpIsSliceInBounds:
		return rewriteValue386_OpIsSliceInBounds_0(v)
	case OpLeq16:
		return rewriteValue386_OpLeq16_0(v)
	case OpLeq16U:
		return rewriteValue386_OpLeq16U_0(v)
	case OpLeq32:
		return rewriteValue386_OpLeq32_0(v)
	case OpLeq32F:
		return rewriteValue386_OpLeq32F_0(v)
	case OpLeq32U:
		return rewriteValue386_OpLeq32U_0(v)
	case OpLeq64F:
		return rewriteValue386_OpLeq64F_0(v)
	case OpLeq8:
		return rewriteValue386_OpLeq8_0(v)
	case OpLeq8U:
		return rewriteValue386_OpLeq8U_0(v)
	case OpLess16:
		return rewriteValue386_OpLess16_0(v)
	case OpLess16U:
		return rewriteValue386_OpLess16U_0(v)
	case OpLess32:
		return rewriteValue386_OpLess32_0(v)
	case OpLess32F:
		return rewriteValue386_OpLess32F_0(v)
	case OpLess32U:
		return rewriteValue386_OpLess32U_0(v)
	case OpLess64F:
		return rewriteValue386_OpLess64F_0(v)
	case OpLess8:
		return rewriteValue386_OpLess8_0(v)
	case OpLess8U:
		return rewriteValue386_OpLess8U_0(v)
	case OpLoad:
		return rewriteValue386_OpLoad_0(v)
	case OpLocalAddr:
		return rewriteValue386_OpLocalAddr_0(v)
	case OpLsh16x16:
		return rewriteValue386_OpLsh16x16_0(v)
	case OpLsh16x32:
		return rewriteValue386_OpLsh16x32_0(v)
	case OpLsh16x64:
		return rewriteValue386_OpLsh16x64_0(v)
	case OpLsh16x8:
		return rewriteValue386_OpLsh16x8_0(v)
	case OpLsh32x16:
		return rewriteValue386_OpLsh32x16_0(v)
	case OpLsh32x32:
		return rewriteValue386_OpLsh32x32_0(v)
	case OpLsh32x64:
		return rewriteValue386_OpLsh32x64_0(v)
	case OpLsh32x8:
		return rewriteValue386_OpLsh32x8_0(v)
	case OpLsh8x16:
		return rewriteValue386_OpLsh8x16_0(v)
	case OpLsh8x32:
		return rewriteValue386_OpLsh8x32_0(v)
	case OpLsh8x64:
		return rewriteValue386_OpLsh8x64_0(v)
	case OpLsh8x8:
		return rewriteValue386_OpLsh8x8_0(v)
	case OpMod16:
		return rewriteValue386_OpMod16_0(v)
	case OpMod16u:
		return rewriteValue386_OpMod16u_0(v)
	case OpMod32:
		return rewriteValue386_OpMod32_0(v)
	case OpMod32u:
		return rewriteValue386_OpMod32u_0(v)
	case OpMod8:
		return rewriteValue386_OpMod8_0(v)
	case OpMod8u:
		return rewriteValue386_OpMod8u_0(v)
	case OpMove:
		return rewriteValue386_OpMove_0(v) || rewriteValue386_OpMove_10(v)
	case OpMul16:
		return rewriteValue386_OpMul16_0(v)
	case OpMul32:
		return rewriteValue386_OpMul32_0(v)
	case OpMul32F:
		return rewriteValue386_OpMul32F_0(v)
	case OpMul32uhilo:
		return rewriteValue386_OpMul32uhilo_0(v)
	case OpMul64F:
		return rewriteValue386_OpMul64F_0(v)
	case OpMul8:
		return rewriteValue386_OpMul8_0(v)
	case OpNeg16:
		return rewriteValue386_OpNeg16_0(v)
	case OpNeg32:
		return rewriteValue386_OpNeg32_0(v)
	case OpNeg32F:
		return rewriteValue386_OpNeg32F_0(v)
	case OpNeg64F:
		return rewriteValue386_OpNeg64F_0(v)
	case OpNeg8:
		return rewriteValue386_OpNeg8_0(v)
	case OpNeq16:
		return rewriteValue386_OpNeq16_0(v)
	case OpNeq32:
		return rewriteValue386_OpNeq32_0(v)
	case OpNeq32F:
		return rewriteValue386_OpNeq32F_0(v)
	case OpNeq64F:
		return rewriteValue386_OpNeq64F_0(v)
	case OpNeq8:
		return rewriteValue386_OpNeq8_0(v)
	case OpNeqB:
		return rewriteValue386_OpNeqB_0(v)
	case OpNeqPtr:
		return rewriteValue386_OpNeqPtr_0(v)
	case OpNilCheck:
		return rewriteValue386_OpNilCheck_0(v)
	case OpNot:
		return rewriteValue386_OpNot_0(v)
	case OpOffPtr:
		return rewriteValue386_OpOffPtr_0(v)
	case OpOr16:
		return rewriteValue386_OpOr16_0(v)
	case OpOr32:
		return rewriteValue386_OpOr32_0(v)
	case OpOr8:
		return rewriteValue386_OpOr8_0(v)
	case OpOrB:
		return rewriteValue386_OpOrB_0(v)
	case OpPanicBounds:
		return rewriteValue386_OpPanicBounds_0(v)
	case OpPanicExtend:
		return rewriteValue386_OpPanicExtend_0(v)
	case OpRotateLeft16:
		return rewriteValue386_OpRotateLeft16_0(v)
	case OpRotateLeft32:
		return rewriteValue386_OpRotateLeft32_0(v)
	case OpRotateLeft8:
		return rewriteValue386_OpRotateLeft8_0(v)
	case OpRound32F:
		return rewriteValue386_OpRound32F_0(v)
	case OpRound64F:
		return rewriteValue386_OpRound64F_0(v)
	case OpRsh16Ux16:
		return rewriteValue386_OpRsh16Ux16_0(v)
	case OpRsh16Ux32:
		return rewriteValue386_OpRsh16Ux32_0(v)
	case OpRsh16Ux64:
		return rewriteValue386_OpRsh16Ux64_0(v)
	case OpRsh16Ux8:
		return rewriteValue386_OpRsh16Ux8_0(v)
	case OpRsh16x16:
		return rewriteValue386_OpRsh16x16_0(v)
	case OpRsh16x32:
		return rewriteValue386_OpRsh16x32_0(v)
	case OpRsh16x64:
		return rewriteValue386_OpRsh16x64_0(v)
	case OpRsh16x8:
		return rewriteValue386_OpRsh16x8_0(v)
	case OpRsh32Ux16:
		return rewriteValue386_OpRsh32Ux16_0(v)
	case OpRsh32Ux32:
		return rewriteValue386_OpRsh32Ux32_0(v)
	case OpRsh32Ux64:
		return rewriteValue386_OpRsh32Ux64_0(v)
	case OpRsh32Ux8:
		return rewriteValue386_OpRsh32Ux8_0(v)
	case OpRsh32x16:
		return rewriteValue386_OpRsh32x16_0(v)
	case OpRsh32x32:
		return rewriteValue386_OpRsh32x32_0(v)
	case OpRsh32x64:
		return rewriteValue386_OpRsh32x64_0(v)
	case OpRsh32x8:
		return rewriteValue386_OpRsh32x8_0(v)
	case OpRsh8Ux16:
		return rewriteValue386_OpRsh8Ux16_0(v)
	case OpRsh8Ux32:
		return rewriteValue386_OpRsh8Ux32_0(v)
	case OpRsh8Ux64:
		return rewriteValue386_OpRsh8Ux64_0(v)
	case OpRsh8Ux8:
		return rewriteValue386_OpRsh8Ux8_0(v)
	case OpRsh8x16:
		return rewriteValue386_OpRsh8x16_0(v)
	case OpRsh8x32:
		return rewriteValue386_OpRsh8x32_0(v)
	case OpRsh8x64:
		return rewriteValue386_OpRsh8x64_0(v)
	case OpRsh8x8:
		return rewriteValue386_OpRsh8x8_0(v)
	case OpSelect0:
		return rewriteValue386_OpSelect0_0(v)
	case OpSelect1:
		return rewriteValue386_OpSelect1_0(v)
	case OpSignExt16to32:
		return rewriteValue386_OpSignExt16to32_0(v)
	case OpSignExt8to16:
		return rewriteValue386_OpSignExt8to16_0(v)
	case OpSignExt8to32:
		return rewriteValue386_OpSignExt8to32_0(v)
	case OpSignmask:
		return rewriteValue386_OpSignmask_0(v)
	case OpSlicemask:
		return rewriteValue386_OpSlicemask_0(v)
	case OpSqrt:
		return rewriteValue386_OpSqrt_0(v)
	case OpStaticCall:
		return rewriteValue386_OpStaticCall_0(v)
	case OpStore:
		return rewriteValue386_OpStore_0(v)
	case OpSub16:
		return rewriteValue386_OpSub16_0(v)
	case OpSub32:
		return rewriteValue386_OpSub32_0(v)
	case OpSub32F:
		return rewriteValue386_OpSub32F_0(v)
	case OpSub32carry:
		return rewriteValue386_OpSub32carry_0(v)
	case OpSub32withcarry:
		return rewriteValue386_OpSub32withcarry_0(v)
	case OpSub64F:
		return rewriteValue386_OpSub64F_0(v)
	case OpSub8:
		return rewriteValue386_OpSub8_0(v)
	case OpSubPtr:
		return rewriteValue386_OpSubPtr_0(v)
	case OpTrunc16to8:
		return rewriteValue386_OpTrunc16to8_0(v)
	case OpTrunc32to16:
		return rewriteValue386_OpTrunc32to16_0(v)
	case OpTrunc32to8:
		return rewriteValue386_OpTrunc32to8_0(v)
	case OpWB:
		return rewriteValue386_OpWB_0(v)
	case OpXor16:
		return rewriteValue386_OpXor16_0(v)
	case OpXor32:
		return rewriteValue386_OpXor32_0(v)
	case OpXor8:
		return rewriteValue386_OpXor8_0(v)
	case OpZero:
		return rewriteValue386_OpZero_0(v) || rewriteValue386_OpZero_10(v)
	case OpZeroExt16to32:
		return rewriteValue386_OpZeroExt16to32_0(v)
	case OpZeroExt8to16:
		return rewriteValue386_OpZeroExt8to16_0(v)
	case OpZeroExt8to32:
		return rewriteValue386_OpZeroExt8to32_0(v)
	case OpZeromask:
		return rewriteValue386_OpZeromask_0(v)
	}
	return false
}
func rewriteValue386_Op386ADCL_0(v *Value) bool {
	// match: (ADCL x (MOVLconst [c]) f)
	// result: (ADCLconst [c] x f)
	for {
		f := v.Args[2]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ADCLconst)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(f)
		return true
	}
	// match: (ADCL (MOVLconst [c]) x f)
	// result: (ADCLconst [c] x f)
	for {
		f := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		x := v.Args[1]
		v.reset(Op386ADCLconst)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(f)
		return true
	}
	// match: (ADCL (MOVLconst [c]) x f)
	// result: (ADCLconst [c] x f)
	for {
		f := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		x := v.Args[1]
		v.reset(Op386ADCLconst)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(f)
		return true
	}
	// match: (ADCL x (MOVLconst [c]) f)
	// result: (ADCLconst [c] x f)
	for {
		f := v.Args[2]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ADCLconst)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(f)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDL_0(v *Value) bool {
	// match: (ADDL x (MOVLconst [c]))
	// result: (ADDLconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ADDLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL (MOVLconst [c]) x)
	// result: (ADDLconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386ADDLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL (SHLLconst [c] x) (SHRLconst [d] x))
	// cond: d == 32-c
	// result: (ROLLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(Op386ROLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL (SHRLconst [d] x) (SHLLconst [c] x))
	// cond: d == 32-c
	// result: (ROLLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(Op386ROLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL <t> (SHLLconst x [c]) (SHRWconst x [d]))
	// cond: c < 16 && d == 16-c && t.Size() == 2
	// result: (ROLWconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRWconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 16 && d == 16-c && t.Size() == 2) {
			break
		}
		v.reset(Op386ROLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL <t> (SHRWconst x [d]) (SHLLconst x [c]))
	// cond: c < 16 && d == 16-c && t.Size() == 2
	// result: (ROLWconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 16 && d == 16-c && t.Size() == 2) {
			break
		}
		v.reset(Op386ROLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL <t> (SHLLconst x [c]) (SHRBconst x [d]))
	// cond: c < 8 && d == 8-c && t.Size() == 1
	// result: (ROLBconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRBconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 8 && d == 8-c && t.Size() == 1) {
			break
		}
		v.reset(Op386ROLBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL <t> (SHRBconst x [d]) (SHLLconst x [c]))
	// cond: c < 8 && d == 8-c && t.Size() == 1
	// result: (ROLBconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRBconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 8 && d == 8-c && t.Size() == 1) {
			break
		}
		v.reset(Op386ROLBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDL x (SHLLconst [3] y))
	// result: (LEAL8 x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 3 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL8)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL (SHLLconst [3] y) x)
	// result: (LEAL8 x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 3 {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386LEAL8)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDL_10(v *Value) bool {
	// match: (ADDL x (SHLLconst [2] y))
	// result: (LEAL4 x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 2 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL4)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL (SHLLconst [2] y) x)
	// result: (LEAL4 x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 2 {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386LEAL4)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL x (SHLLconst [1] y))
	// result: (LEAL2 x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL (SHLLconst [1] y) x)
	// result: (LEAL2 x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 1 {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL x (ADDL y y))
	// result: (LEAL2 x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDL {
			break
		}
		y := v_1.Args[1]
		if y != v_1.Args[0] {
			break
		}
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL (ADDL y y) x)
	// result: (LEAL2 x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		y := v_0.Args[1]
		if y != v_0.Args[0] {
			break
		}
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL x (ADDL x y))
	// result: (LEAL2 y x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDL {
			break
		}
		y := v_1.Args[1]
		if x != v_1.Args[0] {
			break
		}
		v.reset(Op386LEAL2)
		v.AddArg(y)
		v.AddArg(x)
		return true
	}
	// match: (ADDL x (ADDL y x))
	// result: (LEAL2 y x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDL {
			break
		}
		_ = v_1.Args[1]
		y := v_1.Args[0]
		if x != v_1.Args[1] {
			break
		}
		v.reset(Op386LEAL2)
		v.AddArg(y)
		v.AddArg(x)
		return true
	}
	// match: (ADDL (ADDL x y) x)
	// result: (LEAL2 y x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		y := v_0.Args[1]
		if x != v_0.Args[0] {
			break
		}
		v.reset(Op386LEAL2)
		v.AddArg(y)
		v.AddArg(x)
		return true
	}
	// match: (ADDL (ADDL y x) x)
	// result: (LEAL2 y x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		_ = v_0.Args[1]
		y := v_0.Args[0]
		if x != v_0.Args[1] {
			break
		}
		v.reset(Op386LEAL2)
		v.AddArg(y)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDL_20(v *Value) bool {
	// match: (ADDL (ADDLconst [c] x) y)
	// result: (LEAL1 [c] x y)
	for {
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386LEAL1)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL y (ADDLconst [c] x))
	// result: (LEAL1 [c] x y)
	for {
		_ = v.Args[1]
		y := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		c := v_1.AuxInt
		x := v_1.Args[0]
		v.reset(Op386LEAL1)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL x (LEAL [c] {s} y))
	// cond: x.Op != OpSB && y.Op != OpSB
	// result: (LEAL1 [c] {s} x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		c := v_1.AuxInt
		s := v_1.Aux
		y := v_1.Args[0]
		if !(x.Op != OpSB && y.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL (LEAL [c] {s} y) x)
	// cond: x.Op != OpSB && y.Op != OpSB
	// result: (LEAL1 [c] {s} x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		c := v_0.AuxInt
		s := v_0.Aux
		y := v_0.Args[0]
		if !(x.Op != OpSB && y.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL x l:(MOVLload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ADDLload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ADDL l:(MOVLload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ADDLload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ADDL x l:(MOVLloadidx4 [off] {sym} ptr idx mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ADDLloadidx4 x [off] {sym} ptr idx mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ADDL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ADDLloadidx4 x [off] {sym} ptr idx mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ADDL x (NEGL y))
	// result: (SUBL x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386NEGL {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386SUBL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDL (NEGL y) x)
	// result: (SUBL x y)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386NEGL {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386SUBL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLcarry_0(v *Value) bool {
	// match: (ADDLcarry x (MOVLconst [c]))
	// result: (ADDLconstcarry [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ADDLconstcarry)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ADDLcarry (MOVLconst [c]) x)
	// result: (ADDLconstcarry [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386ADDLconstcarry)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLconst_0(v *Value) bool {
	// match: (ADDLconst [c] (ADDL x y))
	// result: (LEAL1 [c] x y)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(Op386LEAL1)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDLconst [c] (LEAL [d] {s} x))
	// cond: is32Bit(c+d)
	// result: (LEAL [c+d] {s} x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		d := v_0.AuxInt
		s := v_0.Aux
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(Op386LEAL)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		return true
	}
	// match: (ADDLconst [c] (LEAL1 [d] {s} x y))
	// cond: is32Bit(c+d)
	// result: (LEAL1 [c+d] {s} x y)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		d := v_0.AuxInt
		s := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDLconst [c] (LEAL2 [d] {s} x y))
	// cond: is32Bit(c+d)
	// result: (LEAL2 [c+d] {s} x y)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL2 {
			break
		}
		d := v_0.AuxInt
		s := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(Op386LEAL2)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDLconst [c] (LEAL4 [d] {s} x y))
	// cond: is32Bit(c+d)
	// result: (LEAL4 [c+d] {s} x y)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		d := v_0.AuxInt
		s := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(Op386LEAL4)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDLconst [c] (LEAL8 [d] {s} x y))
	// cond: is32Bit(c+d)
	// result: (LEAL8 [c+d] {s} x y)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL8 {
			break
		}
		d := v_0.AuxInt
		s := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(Op386LEAL8)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (ADDLconst [c] x)
	// cond: int32(c)==0
	// result: x
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(int32(c) == 0) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ADDLconst [c] (MOVLconst [d]))
	// result: (MOVLconst [int64(int32(c+d))])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = int64(int32(c + d))
		return true
	}
	// match: (ADDLconst [c] (ADDLconst [d] x))
	// result: (ADDLconst [int64(int32(c+d))] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ADDLconst)
		v.AuxInt = int64(int32(c + d))
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLconstmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDLconstmodify [valoff1] {sym} (ADDLconst [off2] base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (ADDLconstmodify [ValAndOff(valoff1).add(off2)] {sym} base mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386ADDLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLconstmodify [valoff1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDLconstmodify [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLconstmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDLconstmodifyidx4 [valoff1] {sym} (ADDLconst [off2] base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (ADDLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386ADDLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLconstmodifyidx4 [valoff1] {sym} base (ADDLconst [off2] idx) mem)
	// cond: ValAndOff(valoff1).canAdd(off2*4)
	// result: (ADDLconstmodifyidx4 [ValAndOff(valoff1).add(off2*4)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2 * 4)) {
			break
		}
		v.reset(Op386ADDLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2 * 4)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLconstmodifyidx4 [valoff1] {sym1} (LEAL [off2] {sym2} base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDLload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (ADDLload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ADDLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDLload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLload [off1] {sym1} val (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (ADDLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL4 {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		idx := v_1.Args[1]
		ptr := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386ADDLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLloadidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDLloadidx4 [off1] {sym} val (ADDLconst [off2] base) idx mem)
	// cond: is32Bit(off1+off2)
	// result: (ADDLloadidx4 [off1+off2] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ADDLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLloadidx4 [off1] {sym} val base (ADDLconst [off2] idx) mem)
	// cond: is32Bit(off1+off2*4)
	// result: (ADDLloadidx4 [off1+off2*4] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		base := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386ADDLconst {
			break
		}
		off2 := v_2.AuxInt
		idx := v_2.Args[0]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386ADDLloadidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLloadidx4 [off1] {sym1} val (LEAL [off2] {sym2} base) idx mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val base idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDLmodify [off1] {sym} (ADDLconst [off2] base) val mem)
	// cond: is32Bit(off1+off2)
	// result: (ADDLmodify [off1+off2] {sym} base val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ADDLmodify)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLmodify [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDLmodify [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDLmodify)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDLmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDLmodifyidx4 [off1] {sym} (ADDLconst [off2] base) idx val mem)
	// cond: is32Bit(off1+off2)
	// result: (ADDLmodifyidx4 [off1+off2] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ADDLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLmodifyidx4 [off1] {sym} base (ADDLconst [off2] idx) val mem)
	// cond: is32Bit(off1+off2*4)
	// result: (ADDLmodifyidx4 [off1+off2*4] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386ADDLmodifyidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLmodifyidx4 [off1] {sym1} (LEAL [off2] {sym2} base) idx val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDLmodifyidx4 [off1+off2] {mergeSym(sym1,sym2)} base idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ADDLmodifyidx4 [off] {sym} ptr idx (MOVLconst [c]) mem)
	// cond: validValAndOff(c,off)
	// result: (ADDLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386MOVLconst {
			break
		}
		c := v_2.AuxInt
		if !(validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ADDLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDSD_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDSD x l:(MOVSDload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (ADDSDload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSDload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386ADDSDload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ADDSD l:(MOVSDload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (ADDSDload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVSDload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386ADDSDload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDSDload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDSDload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (ADDSDload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ADDSDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ADDSDload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDSDload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDSDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDSS_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDSS x l:(MOVSSload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (ADDSSload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSSload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386ADDSSload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ADDSS l:(MOVSSload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (ADDSSload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVSSload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386ADDSSload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ADDSSload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ADDSSload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (ADDSSload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ADDSSload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ADDSSload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ADDSSload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ADDSSload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDL_0(v *Value) bool {
	// match: (ANDL x (MOVLconst [c]))
	// result: (ANDLconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ANDLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ANDL (MOVLconst [c]) x)
	// result: (ANDLconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386ANDLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ANDL x l:(MOVLload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ANDLload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ANDL l:(MOVLload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ANDLload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ANDL x l:(MOVLloadidx4 [off] {sym} ptr idx mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ANDLloadidx4 x [off] {sym} ptr idx mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ANDL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ANDLloadidx4 x [off] {sym} ptr idx mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ANDL x x)
	// result: x
	for {
		x := v.Args[1]
		if x != v.Args[0] {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLconst_0(v *Value) bool {
	// match: (ANDLconst [c] (ANDLconst [d] x))
	// result: (ANDLconst [c & d] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ANDLconst)
		v.AuxInt = c & d
		v.AddArg(x)
		return true
	}
	// match: (ANDLconst [c] _)
	// cond: int32(c)==0
	// result: (MOVLconst [0])
	for {
		c := v.AuxInt
		if !(int32(c) == 0) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (ANDLconst [c] x)
	// cond: int32(c)==-1
	// result: x
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(int32(c) == -1) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ANDLconst [c] (MOVLconst [d]))
	// result: (MOVLconst [c&d])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = c & d
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLconstmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ANDLconstmodify [valoff1] {sym} (ADDLconst [off2] base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (ANDLconstmodify [ValAndOff(valoff1).add(off2)] {sym} base mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386ANDLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLconstmodify [valoff1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ANDLconstmodify [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ANDLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLconstmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ANDLconstmodifyidx4 [valoff1] {sym} (ADDLconst [off2] base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (ANDLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386ANDLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLconstmodifyidx4 [valoff1] {sym} base (ADDLconst [off2] idx) mem)
	// cond: ValAndOff(valoff1).canAdd(off2*4)
	// result: (ANDLconstmodifyidx4 [ValAndOff(valoff1).add(off2*4)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2 * 4)) {
			break
		}
		v.reset(Op386ANDLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2 * 4)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLconstmodifyidx4 [valoff1] {sym1} (LEAL [off2] {sym2} base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ANDLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ANDLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ANDLload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (ANDLload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ANDLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ANDLload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ANDLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLload [off1] {sym1} val (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (ANDLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL4 {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		idx := v_1.Args[1]
		ptr := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386ANDLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLloadidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ANDLloadidx4 [off1] {sym} val (ADDLconst [off2] base) idx mem)
	// cond: is32Bit(off1+off2)
	// result: (ANDLloadidx4 [off1+off2] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ANDLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLloadidx4 [off1] {sym} val base (ADDLconst [off2] idx) mem)
	// cond: is32Bit(off1+off2*4)
	// result: (ANDLloadidx4 [off1+off2*4] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		base := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386ADDLconst {
			break
		}
		off2 := v_2.AuxInt
		idx := v_2.Args[0]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386ANDLloadidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLloadidx4 [off1] {sym1} val (LEAL [off2] {sym2} base) idx mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ANDLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val base idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ANDLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ANDLmodify [off1] {sym} (ADDLconst [off2] base) val mem)
	// cond: is32Bit(off1+off2)
	// result: (ANDLmodify [off1+off2] {sym} base val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ANDLmodify)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLmodify [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ANDLmodify [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ANDLmodify)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ANDLmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ANDLmodifyidx4 [off1] {sym} (ADDLconst [off2] base) idx val mem)
	// cond: is32Bit(off1+off2)
	// result: (ANDLmodifyidx4 [off1+off2] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ANDLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLmodifyidx4 [off1] {sym} base (ADDLconst [off2] idx) val mem)
	// cond: is32Bit(off1+off2*4)
	// result: (ANDLmodifyidx4 [off1+off2*4] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386ANDLmodifyidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLmodifyidx4 [off1] {sym1} (LEAL [off2] {sym2} base) idx val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ANDLmodifyidx4 [off1+off2] {mergeSym(sym1,sym2)} base idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ANDLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ANDLmodifyidx4 [off] {sym} ptr idx (MOVLconst [c]) mem)
	// cond: validValAndOff(c,off)
	// result: (ANDLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386MOVLconst {
			break
		}
		c := v_2.AuxInt
		if !(validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ANDLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPB_0(v *Value) bool {
	b := v.Block
	// match: (CMPB x (MOVLconst [c]))
	// result: (CMPBconst x [int64(int8(c))])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386CMPBconst)
		v.AuxInt = int64(int8(c))
		v.AddArg(x)
		return true
	}
	// match: (CMPB (MOVLconst [c]) x)
	// result: (InvertFlags (CMPBconst x [int64(int8(c))]))
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386InvertFlags)
		v0 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v0.AuxInt = int64(int8(c))
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (CMPB l:(MOVBload {sym} [off] ptr mem) x)
	// cond: canMergeLoad(v, l) && clobber(l)
	// result: (CMPBload {sym} [off] ptr x mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVBload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoad(v, l) && clobber(l)) {
			break
		}
		v.reset(Op386CMPBload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (CMPB x l:(MOVBload {sym} [off] ptr mem))
	// cond: canMergeLoad(v, l) && clobber(l)
	// result: (InvertFlags (CMPBload {sym} [off] ptr x mem))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVBload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoad(v, l) && clobber(l)) {
			break
		}
		v.reset(Op386InvertFlags)
		v0 := b.NewValue0(l.Pos, Op386CMPBload, types.TypeFlags)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(x)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPBconst_0(v *Value) bool {
	b := v.Block
	// match: (CMPBconst (MOVLconst [x]) [y])
	// cond: int8(x)==int8(y)
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int8(x) == int8(y)) {
			break
		}
		v.reset(Op386FlagEQ)
		return true
	}
	// match: (CMPBconst (MOVLconst [x]) [y])
	// cond: int8(x)<int8(y) && uint8(x)<uint8(y)
	// result: (FlagLT_ULT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int8(x) < int8(y) && uint8(x) < uint8(y)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPBconst (MOVLconst [x]) [y])
	// cond: int8(x)<int8(y) && uint8(x)>uint8(y)
	// result: (FlagLT_UGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int8(x) < int8(y) && uint8(x) > uint8(y)) {
			break
		}
		v.reset(Op386FlagLT_UGT)
		return true
	}
	// match: (CMPBconst (MOVLconst [x]) [y])
	// cond: int8(x)>int8(y) && uint8(x)<uint8(y)
	// result: (FlagGT_ULT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int8(x) > int8(y) && uint8(x) < uint8(y)) {
			break
		}
		v.reset(Op386FlagGT_ULT)
		return true
	}
	// match: (CMPBconst (MOVLconst [x]) [y])
	// cond: int8(x)>int8(y) && uint8(x)>uint8(y)
	// result: (FlagGT_UGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int8(x) > int8(y) && uint8(x) > uint8(y)) {
			break
		}
		v.reset(Op386FlagGT_UGT)
		return true
	}
	// match: (CMPBconst (ANDLconst _ [m]) [n])
	// cond: 0 <= int8(m) && int8(m) < int8(n)
	// result: (FlagLT_ULT)
	for {
		n := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		m := v_0.AuxInt
		if !(0 <= int8(m) && int8(m) < int8(n)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPBconst l:(ANDL x y) [0])
	// cond: l.Uses==1
	// result: (TESTB x y)
	for {
		if v.AuxInt != 0 {
			break
		}
		l := v.Args[0]
		if l.Op != Op386ANDL {
			break
		}
		y := l.Args[1]
		x := l.Args[0]
		if !(l.Uses == 1) {
			break
		}
		v.reset(Op386TESTB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPBconst l:(ANDLconst [c] x) [0])
	// cond: l.Uses==1
	// result: (TESTBconst [int64(int8(c))] x)
	for {
		if v.AuxInt != 0 {
			break
		}
		l := v.Args[0]
		if l.Op != Op386ANDLconst {
			break
		}
		c := l.AuxInt
		x := l.Args[0]
		if !(l.Uses == 1) {
			break
		}
		v.reset(Op386TESTBconst)
		v.AuxInt = int64(int8(c))
		v.AddArg(x)
		return true
	}
	// match: (CMPBconst x [0])
	// result: (TESTB x x)
	for {
		if v.AuxInt != 0 {
			break
		}
		x := v.Args[0]
		v.reset(Op386TESTB)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
	// match: (CMPBconst l:(MOVBload {sym} [off] ptr mem) [c])
	// cond: l.Uses == 1 && validValAndOff(c, off) && clobber(l)
	// result: @l.Block (CMPBconstload {sym} [makeValAndOff(c,off)] ptr mem)
	for {
		c := v.AuxInt
		l := v.Args[0]
		if l.Op != Op386MOVBload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(l.Uses == 1 && validValAndOff(c, off) && clobber(l)) {
			break
		}
		b = l.Block
		v0 := b.NewValue0(l.Pos, Op386CMPBconstload, types.TypeFlags)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = makeValAndOff(c, off)
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPBload_0(v *Value) bool {
	// match: (CMPBload {sym} [off] ptr (MOVLconst [c]) mem)
	// cond: validValAndOff(int64(int8(c)),off)
	// result: (CMPBconstload {sym} [makeValAndOff(int64(int8(c)),off)] ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(validValAndOff(int64(int8(c)), off)) {
			break
		}
		v.reset(Op386CMPBconstload)
		v.AuxInt = makeValAndOff(int64(int8(c)), off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPL_0(v *Value) bool {
	b := v.Block
	// match: (CMPL x (MOVLconst [c]))
	// result: (CMPLconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386CMPLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (CMPL (MOVLconst [c]) x)
	// result: (InvertFlags (CMPLconst x [c]))
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386InvertFlags)
		v0 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (CMPL l:(MOVLload {sym} [off] ptr mem) x)
	// cond: canMergeLoad(v, l) && clobber(l)
	// result: (CMPLload {sym} [off] ptr x mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoad(v, l) && clobber(l)) {
			break
		}
		v.reset(Op386CMPLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (CMPL x l:(MOVLload {sym} [off] ptr mem))
	// cond: canMergeLoad(v, l) && clobber(l)
	// result: (InvertFlags (CMPLload {sym} [off] ptr x mem))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoad(v, l) && clobber(l)) {
			break
		}
		v.reset(Op386InvertFlags)
		v0 := b.NewValue0(l.Pos, Op386CMPLload, types.TypeFlags)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(x)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPLconst_0(v *Value) bool {
	// match: (CMPLconst (MOVLconst [x]) [y])
	// cond: int32(x)==int32(y)
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) == int32(y)) {
			break
		}
		v.reset(Op386FlagEQ)
		return true
	}
	// match: (CMPLconst (MOVLconst [x]) [y])
	// cond: int32(x)<int32(y) && uint32(x)<uint32(y)
	// result: (FlagLT_ULT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) < int32(y) && uint32(x) < uint32(y)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPLconst (MOVLconst [x]) [y])
	// cond: int32(x)<int32(y) && uint32(x)>uint32(y)
	// result: (FlagLT_UGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) < int32(y) && uint32(x) > uint32(y)) {
			break
		}
		v.reset(Op386FlagLT_UGT)
		return true
	}
	// match: (CMPLconst (MOVLconst [x]) [y])
	// cond: int32(x)>int32(y) && uint32(x)<uint32(y)
	// result: (FlagGT_ULT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) > int32(y) && uint32(x) < uint32(y)) {
			break
		}
		v.reset(Op386FlagGT_ULT)
		return true
	}
	// match: (CMPLconst (MOVLconst [x]) [y])
	// cond: int32(x)>int32(y) && uint32(x)>uint32(y)
	// result: (FlagGT_UGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int32(x) > int32(y) && uint32(x) > uint32(y)) {
			break
		}
		v.reset(Op386FlagGT_UGT)
		return true
	}
	// match: (CMPLconst (SHRLconst _ [c]) [n])
	// cond: 0 <= n && 0 < c && c <= 32 && (1<<uint64(32-c)) <= uint64(n)
	// result: (FlagLT_ULT)
	for {
		n := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRLconst {
			break
		}
		c := v_0.AuxInt
		if !(0 <= n && 0 < c && c <= 32 && (1<<uint64(32-c)) <= uint64(n)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPLconst (ANDLconst _ [m]) [n])
	// cond: 0 <= int32(m) && int32(m) < int32(n)
	// result: (FlagLT_ULT)
	for {
		n := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		m := v_0.AuxInt
		if !(0 <= int32(m) && int32(m) < int32(n)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPLconst l:(ANDL x y) [0])
	// cond: l.Uses==1
	// result: (TESTL x y)
	for {
		if v.AuxInt != 0 {
			break
		}
		l := v.Args[0]
		if l.Op != Op386ANDL {
			break
		}
		y := l.Args[1]
		x := l.Args[0]
		if !(l.Uses == 1) {
			break
		}
		v.reset(Op386TESTL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPLconst l:(ANDLconst [c] x) [0])
	// cond: l.Uses==1
	// result: (TESTLconst [c] x)
	for {
		if v.AuxInt != 0 {
			break
		}
		l := v.Args[0]
		if l.Op != Op386ANDLconst {
			break
		}
		c := l.AuxInt
		x := l.Args[0]
		if !(l.Uses == 1) {
			break
		}
		v.reset(Op386TESTLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (CMPLconst x [0])
	// result: (TESTL x x)
	for {
		if v.AuxInt != 0 {
			break
		}
		x := v.Args[0]
		v.reset(Op386TESTL)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPLconst_10(v *Value) bool {
	b := v.Block
	// match: (CMPLconst l:(MOVLload {sym} [off] ptr mem) [c])
	// cond: l.Uses == 1 && validValAndOff(c, off) && clobber(l)
	// result: @l.Block (CMPLconstload {sym} [makeValAndOff(c,off)] ptr mem)
	for {
		c := v.AuxInt
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(l.Uses == 1 && validValAndOff(c, off) && clobber(l)) {
			break
		}
		b = l.Block
		v0 := b.NewValue0(l.Pos, Op386CMPLconstload, types.TypeFlags)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = makeValAndOff(c, off)
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPLload_0(v *Value) bool {
	// match: (CMPLload {sym} [off] ptr (MOVLconst [c]) mem)
	// cond: validValAndOff(int64(int32(c)),off)
	// result: (CMPLconstload {sym} [makeValAndOff(int64(int32(c)),off)] ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(validValAndOff(int64(int32(c)), off)) {
			break
		}
		v.reset(Op386CMPLconstload)
		v.AuxInt = makeValAndOff(int64(int32(c)), off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPW_0(v *Value) bool {
	b := v.Block
	// match: (CMPW x (MOVLconst [c]))
	// result: (CMPWconst x [int64(int16(c))])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386CMPWconst)
		v.AuxInt = int64(int16(c))
		v.AddArg(x)
		return true
	}
	// match: (CMPW (MOVLconst [c]) x)
	// result: (InvertFlags (CMPWconst x [int64(int16(c))]))
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386InvertFlags)
		v0 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v0.AuxInt = int64(int16(c))
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (CMPW l:(MOVWload {sym} [off] ptr mem) x)
	// cond: canMergeLoad(v, l) && clobber(l)
	// result: (CMPWload {sym} [off] ptr x mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVWload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoad(v, l) && clobber(l)) {
			break
		}
		v.reset(Op386CMPWload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (CMPW x l:(MOVWload {sym} [off] ptr mem))
	// cond: canMergeLoad(v, l) && clobber(l)
	// result: (InvertFlags (CMPWload {sym} [off] ptr x mem))
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVWload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoad(v, l) && clobber(l)) {
			break
		}
		v.reset(Op386InvertFlags)
		v0 := b.NewValue0(l.Pos, Op386CMPWload, types.TypeFlags)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(x)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPWconst_0(v *Value) bool {
	b := v.Block
	// match: (CMPWconst (MOVLconst [x]) [y])
	// cond: int16(x)==int16(y)
	// result: (FlagEQ)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int16(x) == int16(y)) {
			break
		}
		v.reset(Op386FlagEQ)
		return true
	}
	// match: (CMPWconst (MOVLconst [x]) [y])
	// cond: int16(x)<int16(y) && uint16(x)<uint16(y)
	// result: (FlagLT_ULT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int16(x) < int16(y) && uint16(x) < uint16(y)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPWconst (MOVLconst [x]) [y])
	// cond: int16(x)<int16(y) && uint16(x)>uint16(y)
	// result: (FlagLT_UGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int16(x) < int16(y) && uint16(x) > uint16(y)) {
			break
		}
		v.reset(Op386FlagLT_UGT)
		return true
	}
	// match: (CMPWconst (MOVLconst [x]) [y])
	// cond: int16(x)>int16(y) && uint16(x)<uint16(y)
	// result: (FlagGT_ULT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int16(x) > int16(y) && uint16(x) < uint16(y)) {
			break
		}
		v.reset(Op386FlagGT_ULT)
		return true
	}
	// match: (CMPWconst (MOVLconst [x]) [y])
	// cond: int16(x)>int16(y) && uint16(x)>uint16(y)
	// result: (FlagGT_UGT)
	for {
		y := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		x := v_0.AuxInt
		if !(int16(x) > int16(y) && uint16(x) > uint16(y)) {
			break
		}
		v.reset(Op386FlagGT_UGT)
		return true
	}
	// match: (CMPWconst (ANDLconst _ [m]) [n])
	// cond: 0 <= int16(m) && int16(m) < int16(n)
	// result: (FlagLT_ULT)
	for {
		n := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		m := v_0.AuxInt
		if !(0 <= int16(m) && int16(m) < int16(n)) {
			break
		}
		v.reset(Op386FlagLT_ULT)
		return true
	}
	// match: (CMPWconst l:(ANDL x y) [0])
	// cond: l.Uses==1
	// result: (TESTW x y)
	for {
		if v.AuxInt != 0 {
			break
		}
		l := v.Args[0]
		if l.Op != Op386ANDL {
			break
		}
		y := l.Args[1]
		x := l.Args[0]
		if !(l.Uses == 1) {
			break
		}
		v.reset(Op386TESTW)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (CMPWconst l:(ANDLconst [c] x) [0])
	// cond: l.Uses==1
	// result: (TESTWconst [int64(int16(c))] x)
	for {
		if v.AuxInt != 0 {
			break
		}
		l := v.Args[0]
		if l.Op != Op386ANDLconst {
			break
		}
		c := l.AuxInt
		x := l.Args[0]
		if !(l.Uses == 1) {
			break
		}
		v.reset(Op386TESTWconst)
		v.AuxInt = int64(int16(c))
		v.AddArg(x)
		return true
	}
	// match: (CMPWconst x [0])
	// result: (TESTW x x)
	for {
		if v.AuxInt != 0 {
			break
		}
		x := v.Args[0]
		v.reset(Op386TESTW)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
	// match: (CMPWconst l:(MOVWload {sym} [off] ptr mem) [c])
	// cond: l.Uses == 1 && validValAndOff(c, off) && clobber(l)
	// result: @l.Block (CMPWconstload {sym} [makeValAndOff(c,off)] ptr mem)
	for {
		c := v.AuxInt
		l := v.Args[0]
		if l.Op != Op386MOVWload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(l.Uses == 1 && validValAndOff(c, off) && clobber(l)) {
			break
		}
		b = l.Block
		v0 := b.NewValue0(l.Pos, Op386CMPWconstload, types.TypeFlags)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = makeValAndOff(c, off)
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386CMPWload_0(v *Value) bool {
	// match: (CMPWload {sym} [off] ptr (MOVLconst [c]) mem)
	// cond: validValAndOff(int64(int16(c)),off)
	// result: (CMPWconstload {sym} [makeValAndOff(int64(int16(c)),off)] ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(validValAndOff(int64(int16(c)), off)) {
			break
		}
		v.reset(Op386CMPWconstload)
		v.AuxInt = makeValAndOff(int64(int16(c)), off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386DIVSD_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (DIVSD x l:(MOVSDload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (DIVSDload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSDload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386DIVSDload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386DIVSDload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (DIVSDload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (DIVSDload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386DIVSDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (DIVSDload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (DIVSDload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386DIVSDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386DIVSS_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (DIVSS x l:(MOVSSload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (DIVSSload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSSload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386DIVSSload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386DIVSSload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (DIVSSload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (DIVSSload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386DIVSSload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (DIVSSload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (DIVSSload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386DIVSSload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386LEAL_0(v *Value) bool {
	// match: (LEAL [c] {s} (ADDLconst [d] x))
	// cond: is32Bit(c+d)
	// result: (LEAL [c+d] {s} x)
	for {
		c := v.AuxInt
		s := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		if !(is32Bit(c + d)) {
			break
		}
		v.reset(Op386LEAL)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		return true
	}
	// match: (LEAL [c] {s} (ADDL x y))
	// cond: x.Op != OpSB && y.Op != OpSB
	// result: (LEAL1 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(x.Op != OpSB && y.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL [off1] {sym1} (LEAL [off2] {sym2} x))
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (LEAL [off1+off2] {mergeSym(sym1,sym2)} x)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386LEAL)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		return true
	}
	// match: (LEAL [off1] {sym1} (LEAL1 [off2] {sym2} x y))
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (LEAL1 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL [off1] {sym1} (LEAL2 [off2] {sym2} x y))
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (LEAL2 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL2 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386LEAL2)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL [off1] {sym1} (LEAL4 [off2] {sym2} x y))
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (LEAL4 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386LEAL4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL [off1] {sym1} (LEAL8 [off2] {sym2} x y))
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (LEAL8 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL8 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		y := v_0.Args[1]
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386LEAL8)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386LEAL1_0(v *Value) bool {
	// match: (LEAL1 [c] {s} (ADDLconst [d] x) y)
	// cond: is32Bit(c+d) && x.Op != OpSB
	// result: (LEAL1 [c+d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		if !(is32Bit(c+d) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} y (ADDLconst [d] x))
	// cond: is32Bit(c+d) && x.Op != OpSB
	// result: (LEAL1 [c+d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		y := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		x := v_1.Args[0]
		if !(is32Bit(c+d) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} x (SHLLconst [1] y))
	// result: (LEAL2 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL2)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} (SHLLconst [1] y) x)
	// result: (LEAL2 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 1 {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386LEAL2)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} x (SHLLconst [2] y))
	// result: (LEAL4 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 2 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL4)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} (SHLLconst [2] y) x)
	// result: (LEAL4 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 2 {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386LEAL4)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} x (SHLLconst [3] y))
	// result: (LEAL8 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 3 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL8)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [c] {s} (SHLLconst [3] y) x)
	// result: (LEAL8 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 3 {
			break
		}
		y := v_0.Args[0]
		v.reset(Op386LEAL8)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [off1] {sym1} (LEAL [off2] {sym2} x) y)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB
	// result: (LEAL1 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL1 [off1] {sym1} y (LEAL [off2] {sym2} x))
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB
	// result: (LEAL1 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		_ = v.Args[1]
		y := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		x := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386LEAL2_0(v *Value) bool {
	// match: (LEAL2 [c] {s} (ADDLconst [d] x) y)
	// cond: is32Bit(c+d) && x.Op != OpSB
	// result: (LEAL2 [c+d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		if !(is32Bit(c+d) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL2)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL2 [c] {s} x (ADDLconst [d] y))
	// cond: is32Bit(c+2*d) && y.Op != OpSB
	// result: (LEAL2 [c+2*d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		y := v_1.Args[0]
		if !(is32Bit(c+2*d) && y.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL2)
		v.AuxInt = c + 2*d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL2 [c] {s} x (SHLLconst [1] y))
	// result: (LEAL4 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL4)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL2 [c] {s} x (SHLLconst [2] y))
	// result: (LEAL8 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 2 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL8)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL2 [off1] {sym1} (LEAL [off2] {sym2} x) y)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB
	// result: (LEAL2 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL2)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386LEAL4_0(v *Value) bool {
	// match: (LEAL4 [c] {s} (ADDLconst [d] x) y)
	// cond: is32Bit(c+d) && x.Op != OpSB
	// result: (LEAL4 [c+d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		if !(is32Bit(c+d) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL4)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL4 [c] {s} x (ADDLconst [d] y))
	// cond: is32Bit(c+4*d) && y.Op != OpSB
	// result: (LEAL4 [c+4*d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		y := v_1.Args[0]
		if !(is32Bit(c+4*d) && y.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL4)
		v.AuxInt = c + 4*d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL4 [c] {s} x (SHLLconst [1] y))
	// result: (LEAL8 [c] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386LEAL8)
		v.AuxInt = c
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL4 [off1] {sym1} (LEAL [off2] {sym2} x) y)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB
	// result: (LEAL4 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386LEAL8_0(v *Value) bool {
	// match: (LEAL8 [c] {s} (ADDLconst [d] x) y)
	// cond: is32Bit(c+d) && x.Op != OpSB
	// result: (LEAL8 [c+d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		if !(is32Bit(c+d) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL8)
		v.AuxInt = c + d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL8 [c] {s} x (ADDLconst [d] y))
	// cond: is32Bit(c+8*d) && y.Op != OpSB
	// result: (LEAL8 [c+8*d] {s} x y)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		y := v_1.Args[0]
		if !(is32Bit(c+8*d) && y.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL8)
		v.AuxInt = c + 8*d
		v.Aux = s
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	// match: (LEAL8 [off1] {sym1} (LEAL [off2] {sym2} x) y)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB
	// result: (LEAL8 [off1+off2] {mergeSym(sym1,sym2)} x y)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		y := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		x := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && x.Op != OpSB) {
			break
		}
		v.reset(Op386LEAL8)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBLSX_0(v *Value) bool {
	b := v.Block
	// match: (MOVBLSX x:(MOVBload [off] {sym} ptr mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVBLSXload <v.Type> [off] {sym} ptr mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVBload {
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
		v0 := b.NewValue0(x.Pos, Op386MOVBLSXload, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVBLSX (ANDLconst [c] x))
	// cond: c & 0x80 == 0
	// result: (ANDLconst [c & 0x7f] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c&0x80 == 0) {
			break
		}
		v.reset(Op386ANDLconst)
		v.AuxInt = c & 0x7f
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBLSXload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVBLSXload [off] {sym} ptr (MOVBstore [off2] {sym2} ptr2 x _))
	// cond: sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)
	// result: (MOVBLSX x)
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVBstore {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		_ = v_1.Args[2]
		ptr2 := v_1.Args[0]
		x := v_1.Args[1]
		if !(sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)) {
			break
		}
		v.reset(Op386MOVBLSX)
		v.AddArg(x)
		return true
	}
	// match: (MOVBLSXload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVBLSXload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVBLSXload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBLZX_0(v *Value) bool {
	b := v.Block
	// match: (MOVBLZX x:(MOVBload [off] {sym} ptr mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVBload <v.Type> [off] {sym} ptr mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVBload {
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
		v0 := b.NewValue0(x.Pos, Op386MOVBload, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVBLZX x:(MOVBloadidx1 [off] {sym} ptr idx mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVBloadidx1 <v.Type> [off] {sym} ptr idx mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVBloadidx1 {
			break
		}
		off := x.AuxInt
		sym := x.Aux
		mem := x.Args[2]
		ptr := x.Args[0]
		idx := x.Args[1]
		if !(x.Uses == 1 && clobber(x)) {
			break
		}
		b = x.Block
		v0 := b.NewValue0(v.Pos, Op386MOVBloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVBLZX (ANDLconst [c] x))
	// result: (ANDLconst [c & 0xff] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ANDLconst)
		v.AuxInt = c & 0xff
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVBload [off] {sym} ptr (MOVBstore [off2] {sym2} ptr2 x _))
	// cond: sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)
	// result: (MOVBLZX x)
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVBstore {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		_ = v_1.Args[2]
		ptr2 := v_1.Args[0]
		x := v_1.Args[1]
		if !(sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)) {
			break
		}
		v.reset(Op386MOVBLZX)
		v.AddArg(x)
		return true
	}
	// match: (MOVBload [off1] {sym} (ADDLconst [off2] ptr) mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVBload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVBload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVBload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVBload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBload [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVBloadidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVBloadidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBload [off] {sym} (ADDL ptr idx) mem)
	// cond: ptr.Op != OpSB
	// result: (MOVBloadidx1 [off] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVBloadidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBload [off] {sym} (SB) _)
	// cond: symIsRO(sym)
	// result: (MOVLconst [int64(read8(sym, off))])
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpSB || !(symIsRO(sym)) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = int64(read8(sym, off))
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBloadidx1_0(v *Value) bool {
	// match: (MOVBloadidx1 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVBloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVBloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBloadidx1 [c] {sym} idx (ADDLconst [d] ptr) mem)
	// result: (MOVBloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		idx := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		ptr := v_1.Args[0]
		v.reset(Op386MOVBloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBloadidx1 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVBloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVBloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBloadidx1 [c] {sym} (ADDLconst [d] idx) ptr mem)
	// result: (MOVBloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		idx := v_0.Args[0]
		ptr := v.Args[1]
		v.reset(Op386MOVBloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstore_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVBstore [off] {sym} ptr (MOVBLSX x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVBLSX {
			break
		}
		x := v_1.Args[0]
		v.reset(Op386MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVBLZX x) mem)
	// result: (MOVBstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVBLZX {
			break
		}
		x := v_1.Args[0]
		v.reset(Op386MOVBstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off1] {sym} (ADDLconst [off2] ptr) val mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVBstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVBstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} ptr (MOVLconst [c]) mem)
	// cond: validOff(off)
	// result: (MOVBstoreconst [makeValAndOff(int64(int8(c)),off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(validOff(off)) {
			break
		}
		v.reset(Op386MOVBstoreconst)
		v.AuxInt = makeValAndOff(int64(int8(c)), off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVBstore [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVBstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVBstoreidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVBstoreidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [off] {sym} (ADDL ptr idx) val mem)
	// cond: ptr.Op != OpSB
	// result: (MOVBstoreidx1 [off] {sym} ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVBstoreidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i] {s} p (SHRWconst [8] w) x:(MOVBstore [i-1] {s} p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstore [i-1] {s} p w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRWconst || v_1.AuxInt != 8 {
			break
		}
		w := v_1.Args[0]
		x := v.Args[2]
		if x.Op != Op386MOVBstore || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] || w != x.Args[1] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i] {s} p (SHRLconst [8] w) x:(MOVBstore [i-1] {s} p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstore [i-1] {s} p w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst || v_1.AuxInt != 8 {
			break
		}
		w := v_1.Args[0]
		x := v.Args[2]
		if x.Op != Op386MOVBstore || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] || w != x.Args[1] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i] {s} p w x:(MOVBstore {s} [i+1] p (SHRWconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstore [i] {s} p w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		w := v.Args[1]
		x := v.Args[2]
		if x.Op != Op386MOVBstore || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] {
			break
		}
		x_1 := x.Args[1]
		if x_1.Op != Op386SHRWconst || x_1.AuxInt != 8 || w != x_1.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstore_10(v *Value) bool {
	// match: (MOVBstore [i] {s} p w x:(MOVBstore {s} [i+1] p (SHRLconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstore [i] {s} p w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		w := v.Args[1]
		x := v.Args[2]
		if x.Op != Op386MOVBstore || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] {
			break
		}
		x_1 := x.Args[1]
		if x_1.Op != Op386SHRLconst || x_1.AuxInt != 8 || w != x_1.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstore [i] {s} p (SHRLconst [j] w) x:(MOVBstore [i-1] {s} p w0:(SHRLconst [j-8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstore [i-1] {s} p w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst {
			break
		}
		j := v_1.AuxInt
		w := v_1.Args[0]
		x := v.Args[2]
		if x.Op != Op386MOVBstore || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] {
			break
		}
		w0 := x.Args[1]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-8 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstoreconst_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVBstoreconst [sc] {s} (ADDLconst [off] ptr) mem)
	// cond: ValAndOff(sc).canAdd(off)
	// result: (MOVBstoreconst [ValAndOff(sc).add(off)] {s} ptr mem)
	for {
		sc := v.AuxInt
		s := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(ValAndOff(sc).canAdd(off)) {
			break
		}
		v.reset(Op386MOVBstoreconst)
		v.AuxInt = ValAndOff(sc).add(off)
		v.Aux = s
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconst [sc] {sym1} (LEAL [off] {sym2} ptr) mem)
	// cond: canMergeSym(sym1, sym2) && ValAndOff(sc).canAdd(off) && (ptr.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVBstoreconst [ValAndOff(sc).add(off)] {mergeSym(sym1, sym2)} ptr mem)
	for {
		sc := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2) && ValAndOff(sc).canAdd(off) && (ptr.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVBstoreconst)
		v.AuxInt = ValAndOff(sc).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconst [x] {sym1} (LEAL1 [off] {sym2} ptr idx) mem)
	// cond: canMergeSym(sym1, sym2)
	// result: (MOVBstoreconstidx1 [ValAndOff(x).add(off)] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		x := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVBstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconst [x] {sym} (ADDL ptr idx) mem)
	// result: (MOVBstoreconstidx1 [x] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		v.reset(Op386MOVBstoreconstidx1)
		v.AuxInt = x
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconst [c] {s} p x:(MOVBstoreconst [a] {s} p mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 1 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVWstoreconst [makeValAndOff(ValAndOff(a).Val()&0xff | ValAndOff(c).Val()<<8, ValAndOff(a).Off())] {s} p mem)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		p := v.Args[0]
		x := v.Args[1]
		if x.Op != Op386MOVBstoreconst {
			break
		}
		a := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[1]
		if p != x.Args[0] || !(x.Uses == 1 && ValAndOff(a).Off()+1 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xff|ValAndOff(c).Val()<<8, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconst [a] {s} p x:(MOVBstoreconst [c] {s} p mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 1 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVWstoreconst [makeValAndOff(ValAndOff(a).Val()&0xff | ValAndOff(c).Val()<<8, ValAndOff(a).Off())] {s} p mem)
	for {
		a := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		p := v.Args[0]
		x := v.Args[1]
		if x.Op != Op386MOVBstoreconst {
			break
		}
		c := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[1]
		if p != x.Args[0] || !(x.Uses == 1 && ValAndOff(a).Off()+1 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xff|ValAndOff(c).Val()<<8, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstoreconstidx1_0(v *Value) bool {
	// match: (MOVBstoreconstidx1 [x] {sym} (ADDLconst [c] ptr) idx mem)
	// result: (MOVBstoreconstidx1 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		c := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVBstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconstidx1 [x] {sym} ptr (ADDLconst [c] idx) mem)
	// result: (MOVBstoreconstidx1 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		c := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVBstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreconstidx1 [c] {s} p i x:(MOVBstoreconstidx1 [a] {s} p i mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 1 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVWstoreconstidx1 [makeValAndOff(ValAndOff(a).Val()&0xff | ValAndOff(c).Val()<<8, ValAndOff(a).Off())] {s} p i mem)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		i := v.Args[1]
		x := v.Args[2]
		if x.Op != Op386MOVBstoreconstidx1 {
			break
		}
		a := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] || i != x.Args[1] || !(x.Uses == 1 && ValAndOff(a).Off()+1 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreconstidx1)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xff|ValAndOff(c).Val()<<8, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v.AddArg(i)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstoreidx1_0(v *Value) bool {
	// match: (MOVBstoreidx1 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVBstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVBstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [c] {sym} idx (ADDLconst [d] ptr) val mem)
	// result: (MOVBstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		idx := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		ptr := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVBstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVBstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVBstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [c] {sym} (ADDLconst [d] idx) ptr val mem)
	// result: (MOVBstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		idx := v_0.Args[0]
		ptr := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVBstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx (SHRLconst [8] w) x:(MOVBstoreidx1 [i-1] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx (SHRLconst [8] w) x:(MOVBstoreidx1 [i-1] {s} idx p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p (SHRLconst [8] w) x:(MOVBstoreidx1 [i-1] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p (SHRLconst [8] w) x:(MOVBstoreidx1 [i-1] {s} idx p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx (SHRWconst [8] w) x:(MOVBstoreidx1 [i-1] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRWconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx (SHRWconst [8] w) x:(MOVBstoreidx1 [i-1] {s} idx p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRWconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstoreidx1_10(v *Value) bool {
	// match: (MOVBstoreidx1 [i] {s} idx p (SHRWconst [8] w) x:(MOVBstoreidx1 [i-1] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRWconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p (SHRWconst [8] w) x:(MOVBstoreidx1 [i-1] {s} idx p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRWconst || v_2.AuxInt != 8 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx w x:(MOVBstoreidx1 [i+1] {s} p idx (SHRLconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRLconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx w x:(MOVBstoreidx1 [i+1] {s} idx p (SHRLconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRLconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p w x:(MOVBstoreidx1 [i+1] {s} p idx (SHRLconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRLconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p w x:(MOVBstoreidx1 [i+1] {s} idx p (SHRLconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRLconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx w x:(MOVBstoreidx1 [i+1] {s} p idx (SHRWconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRWconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx w x:(MOVBstoreidx1 [i+1] {s} idx p (SHRWconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRWconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p w x:(MOVBstoreidx1 [i+1] {s} p idx (SHRWconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRWconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p w x:(MOVBstoreidx1 [i+1] {s} idx p (SHRWconst [8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		w := v.Args[2]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i+1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		x_2 := x.Args[2]
		if x_2.Op != Op386SHRWconst || x_2.AuxInt != 8 || w != x_2.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVBstoreidx1_20(v *Value) bool {
	// match: (MOVBstoreidx1 [i] {s} p idx (SHRLconst [j] w) x:(MOVBstoreidx1 [i-1] {s} p idx w0:(SHRLconst [j-8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-8 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} p idx (SHRLconst [j] w) x:(MOVBstoreidx1 [i-1] {s} idx p w0:(SHRLconst [j-8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-8 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p (SHRLconst [j] w) x:(MOVBstoreidx1 [i-1] {s} p idx w0:(SHRLconst [j-8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-8 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVBstoreidx1 [i] {s} idx p (SHRLconst [j] w) x:(MOVBstoreidx1 [i-1] {s} idx p w0:(SHRLconst [j-8] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVWstoreidx1 [i-1] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVBstoreidx1 || x.AuxInt != i-1 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-8 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = i - 1
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVLload [off] {sym} ptr (MOVLstore [off2] {sym2} ptr2 x _))
	// cond: sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)
	// result: x
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLstore {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		_ = v_1.Args[2]
		ptr2 := v_1.Args[0]
		x := v_1.Args[1]
		if !(sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MOVLload [off1] {sym} (ADDLconst [off2] ptr) mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVLload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVLload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLload [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVLloadidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVLloadidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLload [off1] {sym1} (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLload [off] {sym} (ADDL ptr idx) mem)
	// cond: ptr.Op != OpSB
	// result: (MOVLloadidx1 [off] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVLloadidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLload [off] {sym} (SB) _)
	// cond: symIsRO(sym)
	// result: (MOVLconst [int64(int32(read32(sym, off, config.BigEndian)))])
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpSB || !(symIsRO(sym)) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = int64(int32(read32(sym, off, config.BigEndian)))
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLloadidx1_0(v *Value) bool {
	// match: (MOVLloadidx1 [c] {sym} ptr (SHLLconst [2] idx) mem)
	// result: (MOVLloadidx4 [c] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 2 {
			break
		}
		idx := v_1.Args[0]
		v.reset(Op386MOVLloadidx4)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLloadidx1 [c] {sym} (SHLLconst [2] idx) ptr mem)
	// result: (MOVLloadidx4 [c] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 2 {
			break
		}
		idx := v_0.Args[0]
		ptr := v.Args[1]
		v.reset(Op386MOVLloadidx4)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLloadidx1 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVLloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVLloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLloadidx1 [c] {sym} idx (ADDLconst [d] ptr) mem)
	// result: (MOVLloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		idx := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		ptr := v_1.Args[0]
		v.reset(Op386MOVLloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLloadidx1 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVLloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVLloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLloadidx1 [c] {sym} (ADDLconst [d] idx) ptr mem)
	// result: (MOVLloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		idx := v_0.Args[0]
		ptr := v.Args[1]
		v.reset(Op386MOVLloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLloadidx4_0(v *Value) bool {
	// match: (MOVLloadidx4 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVLloadidx4 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVLloadidx4)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLloadidx4 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVLloadidx4 [int64(int32(c+4*d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVLloadidx4)
		v.AuxInt = int64(int32(c + 4*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstore_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVLstore [off1] {sym} (ADDLconst [off2] ptr) val mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVLstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVLstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore [off] {sym} ptr (MOVLconst [c]) mem)
	// cond: validOff(off)
	// result: (MOVLstoreconst [makeValAndOff(int64(int32(c)),off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(validOff(off)) {
			break
		}
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(int64(int32(c)), off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVLstore [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVLstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVLstoreidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore [off1] {sym1} (LEAL4 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVLstoreidx4 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVLstoreidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore [off] {sym} (ADDL ptr idx) val mem)
	// cond: ptr.Op != OpSB
	// result: (MOVLstoreidx1 [off] {sym} ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ADDLload x [off] {sym} ptr mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (ADDLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ADDLload || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[2]
		x := y.Args[0]
		if ptr != y.Args[1] || mem != y.Args[2] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386ADDLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ANDLload x [off] {sym} ptr mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (ANDLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ANDLload || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[2]
		x := y.Args[0]
		if ptr != y.Args[1] || mem != y.Args[2] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386ANDLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ORLload x [off] {sym} ptr mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (ORLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ORLload || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[2]
		x := y.Args[0]
		if ptr != y.Args[1] || mem != y.Args[2] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386ORLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(XORLload x [off] {sym} ptr mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (XORLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386XORLload || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[2]
		x := y.Args[0]
		if ptr != y.Args[1] || mem != y.Args[2] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386XORLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstore_10(v *Value) bool {
	// match: (MOVLstore {sym} [off] ptr y:(ADDL l:(MOVLload [off] {sym} ptr mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ADDLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ADDL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ADDL x l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ADDLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ADDL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(SUBL l:(MOVLload [off] {sym} ptr mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (SUBLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386SUBL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386SUBLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ANDL l:(MOVLload [off] {sym} ptr mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ANDLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ANDL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ANDL x l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ANDLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ANDL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ORL l:(MOVLload [off] {sym} ptr mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ORLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ORL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ORLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ORL x l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ORLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ORL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ORLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(XORL l:(MOVLload [off] {sym} ptr mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (XORLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386XORL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386XORLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(XORL x l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (XORLmodify [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386XORL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386XORLmodify)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ADDLconst [c] l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (ADDLconstmodify [makeValAndOff(c,off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ADDLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ADDLconstmodify)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstore_20(v *Value) bool {
	// match: (MOVLstore {sym} [off] ptr y:(ANDLconst [c] l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (ANDLconstmodify [makeValAndOff(c,off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ANDLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ANDLconstmodify)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(ORLconst [c] l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (ORLconstmodify [makeValAndOff(c,off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386ORLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ORLconstmodify)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstore {sym} [off] ptr y:(XORLconst [c] l:(MOVLload [off] {sym} ptr mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (XORLconstmodify [makeValAndOff(c,off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		y := v.Args[1]
		if y.Op != Op386XORLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLload || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[1]
		if ptr != l.Args[0] || mem != l.Args[1] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386XORLconstmodify)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstoreconst_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVLstoreconst [sc] {s} (ADDLconst [off] ptr) mem)
	// cond: ValAndOff(sc).canAdd(off)
	// result: (MOVLstoreconst [ValAndOff(sc).add(off)] {s} ptr mem)
	for {
		sc := v.AuxInt
		s := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(ValAndOff(sc).canAdd(off)) {
			break
		}
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = ValAndOff(sc).add(off)
		v.Aux = s
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconst [sc] {sym1} (LEAL [off] {sym2} ptr) mem)
	// cond: canMergeSym(sym1, sym2) && ValAndOff(sc).canAdd(off) && (ptr.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVLstoreconst [ValAndOff(sc).add(off)] {mergeSym(sym1, sym2)} ptr mem)
	for {
		sc := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2) && ValAndOff(sc).canAdd(off) && (ptr.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = ValAndOff(sc).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconst [x] {sym1} (LEAL1 [off] {sym2} ptr idx) mem)
	// cond: canMergeSym(sym1, sym2)
	// result: (MOVLstoreconstidx1 [ValAndOff(x).add(off)] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		x := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVLstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconst [x] {sym1} (LEAL4 [off] {sym2} ptr idx) mem)
	// cond: canMergeSym(sym1, sym2)
	// result: (MOVLstoreconstidx4 [ValAndOff(x).add(off)] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		x := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVLstoreconstidx4)
		v.AuxInt = ValAndOff(x).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconst [x] {sym} (ADDL ptr idx) mem)
	// result: (MOVLstoreconstidx1 [x] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		v.reset(Op386MOVLstoreconstidx1)
		v.AuxInt = x
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstoreconstidx1_0(v *Value) bool {
	// match: (MOVLstoreconstidx1 [c] {sym} ptr (SHLLconst [2] idx) mem)
	// result: (MOVLstoreconstidx4 [c] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 2 {
			break
		}
		idx := v_1.Args[0]
		v.reset(Op386MOVLstoreconstidx4)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconstidx1 [x] {sym} (ADDLconst [c] ptr) idx mem)
	// result: (MOVLstoreconstidx1 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		c := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVLstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconstidx1 [x] {sym} ptr (ADDLconst [c] idx) mem)
	// result: (MOVLstoreconstidx1 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		c := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVLstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstoreconstidx4_0(v *Value) bool {
	// match: (MOVLstoreconstidx4 [x] {sym} (ADDLconst [c] ptr) idx mem)
	// result: (MOVLstoreconstidx4 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		c := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVLstoreconstidx4)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreconstidx4 [x] {sym} ptr (ADDLconst [c] idx) mem)
	// result: (MOVLstoreconstidx4 [ValAndOff(x).add(4*c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		c := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVLstoreconstidx4)
		v.AuxInt = ValAndOff(x).add(4 * c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstoreidx1_0(v *Value) bool {
	// match: (MOVLstoreidx1 [c] {sym} ptr (SHLLconst [2] idx) val mem)
	// result: (MOVLstoreidx4 [c] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 2 {
			break
		}
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx4)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx1 [c] {sym} (SHLLconst [2] idx) ptr val mem)
	// result: (MOVLstoreidx4 [c] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 2 {
			break
		}
		idx := v_0.Args[0]
		ptr := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx4)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx1 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVLstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx1 [c] {sym} idx (ADDLconst [d] ptr) val mem)
	// result: (MOVLstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		idx := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		ptr := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx1 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVLstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx1 [c] {sym} (ADDLconst [d] idx) ptr val mem)
	// result: (MOVLstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		idx := v_0.Args[0]
		ptr := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVLstoreidx4_0(v *Value) bool {
	// match: (MOVLstoreidx4 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVLstoreidx4 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx4)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVLstoreidx4 [int64(int32(c+4*d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVLstoreidx4)
		v.AuxInt = int64(int32(c + 4*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ADDLloadidx4 x [off] {sym} ptr idx mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (ADDLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ADDLloadidx4 || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[3]
		x := y.Args[0]
		if ptr != y.Args[1] || idx != y.Args[2] || mem != y.Args[3] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386ADDLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ANDLloadidx4 x [off] {sym} ptr idx mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (ANDLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ANDLloadidx4 || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[3]
		x := y.Args[0]
		if ptr != y.Args[1] || idx != y.Args[2] || mem != y.Args[3] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386ANDLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ORLloadidx4 x [off] {sym} ptr idx mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (ORLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ORLloadidx4 || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[3]
		x := y.Args[0]
		if ptr != y.Args[1] || idx != y.Args[2] || mem != y.Args[3] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386ORLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(XORLloadidx4 x [off] {sym} ptr idx mem) mem)
	// cond: y.Uses==1 && clobber(y)
	// result: (XORLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386XORLloadidx4 || y.AuxInt != off || y.Aux != sym {
			break
		}
		_ = y.Args[3]
		x := y.Args[0]
		if ptr != y.Args[1] || idx != y.Args[2] || mem != y.Args[3] || !(y.Uses == 1 && clobber(y)) {
			break
		}
		v.reset(Op386XORLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ADDL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ADDLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ADDL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ADDL x l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ADDLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ADDL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ADDLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(SUBL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (SUBLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386SUBL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386SUBLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ANDL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ANDLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ANDL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLmodifyidx4)
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
func rewriteValue386_Op386MOVLstoreidx4_10(v *Value) bool {
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ANDL x l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ANDLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ANDL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ANDLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ORL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ORLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ORL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ORLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ORL x l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (ORLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ORL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386ORLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(XORL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (XORLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386XORL {
			break
		}
		x := y.Args[1]
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386XORLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(XORL x l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l)
	// result: (XORLmodifyidx4 [off] {sym} ptr idx x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386XORL {
			break
		}
		_ = y.Args[1]
		x := y.Args[0]
		l := y.Args[1]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l)) {
			break
		}
		v.reset(Op386XORLmodifyidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ADDLconst [c] l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (ADDLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ADDLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ADDLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ANDLconst [c] l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (ANDLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ANDLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ANDLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(ORLconst [c] l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (ORLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386ORLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ORLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVLstoreidx4 {sym} [off] ptr idx y:(XORLconst [c] l:(MOVLloadidx4 [off] {sym} ptr idx mem)) mem)
	// cond: y.Uses==1 && l.Uses==1 && clobber(y) && clobber(l) && validValAndOff(c,off)
	// result: (XORLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		y := v.Args[2]
		if y.Op != Op386XORLconst {
			break
		}
		c := y.AuxInt
		l := y.Args[0]
		if l.Op != Op386MOVLloadidx4 || l.AuxInt != off || l.Aux != sym {
			break
		}
		_ = l.Args[2]
		if ptr != l.Args[0] || idx != l.Args[1] || mem != l.Args[2] || !(y.Uses == 1 && l.Uses == 1 && clobber(y) && clobber(l) && validValAndOff(c, off)) {
			break
		}
		v.reset(Op386XORLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDconst_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (MOVSDconst [c])
	// cond: config.ctxt.Flag_shared
	// result: (MOVSDconst2 (MOVSDconst1 [c]))
	for {
		c := v.AuxInt
		if !(config.ctxt.Flag_shared) {
			break
		}
		v.reset(Op386MOVSDconst2)
		v0 := b.NewValue0(v.Pos, Op386MOVSDconst1, typ.UInt32)
		v0.AuxInt = c
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVSDload [off1] {sym} (ADDLconst [off2] ptr) mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVSDload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVSDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVSDload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVSDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDload [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSDloadidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSDloadidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDload [off1] {sym1} (LEAL8 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSDloadidx8 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL8 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSDloadidx8)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDload [off] {sym} (ADDL ptr idx) mem)
	// cond: ptr.Op != OpSB
	// result: (MOVSDloadidx1 [off] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVSDloadidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDloadidx1_0(v *Value) bool {
	// match: (MOVSDloadidx1 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVSDloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVSDloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDloadidx1 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVSDloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVSDloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDloadidx8_0(v *Value) bool {
	// match: (MOVSDloadidx8 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVSDloadidx8 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVSDloadidx8)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDloadidx8 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVSDloadidx8 [int64(int32(c+8*d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVSDloadidx8)
		v.AuxInt = int64(int32(c + 8*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDstore_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVSDstore [off1] {sym} (ADDLconst [off2] ptr) val mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVSDstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVSDstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDstore [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVSDstore [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVSDstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDstore [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSDstoreidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSDstoreidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDstore [off1] {sym1} (LEAL8 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSDstoreidx8 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL8 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSDstoreidx8)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDstore [off] {sym} (ADDL ptr idx) val mem)
	// cond: ptr.Op != OpSB
	// result: (MOVSDstoreidx1 [off] {sym} ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVSDstoreidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDstoreidx1_0(v *Value) bool {
	// match: (MOVSDstoreidx1 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVSDstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVSDstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDstoreidx1 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVSDstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVSDstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSDstoreidx8_0(v *Value) bool {
	// match: (MOVSDstoreidx8 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVSDstoreidx8 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVSDstoreidx8)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSDstoreidx8 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVSDstoreidx8 [int64(int32(c+8*d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVSDstoreidx8)
		v.AuxInt = int64(int32(c + 8*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSconst_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (MOVSSconst [c])
	// cond: config.ctxt.Flag_shared
	// result: (MOVSSconst2 (MOVSSconst1 [c]))
	for {
		c := v.AuxInt
		if !(config.ctxt.Flag_shared) {
			break
		}
		v.reset(Op386MOVSSconst2)
		v0 := b.NewValue0(v.Pos, Op386MOVSSconst1, typ.UInt32)
		v0.AuxInt = c
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVSSload [off1] {sym} (ADDLconst [off2] ptr) mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVSSload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVSSload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVSSload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVSSload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSload [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSSloadidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSSloadidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSload [off1] {sym1} (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSSloadidx4 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSSloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSload [off] {sym} (ADDL ptr idx) mem)
	// cond: ptr.Op != OpSB
	// result: (MOVSSloadidx1 [off] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVSSloadidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSloadidx1_0(v *Value) bool {
	// match: (MOVSSloadidx1 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVSSloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVSSloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSloadidx1 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVSSloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVSSloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSloadidx4_0(v *Value) bool {
	// match: (MOVSSloadidx4 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVSSloadidx4 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVSSloadidx4)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSloadidx4 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVSSloadidx4 [int64(int32(c+4*d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVSSloadidx4)
		v.AuxInt = int64(int32(c + 4*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSstore_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVSSstore [off1] {sym} (ADDLconst [off2] ptr) val mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVSSstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVSSstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSstore [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVSSstore [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVSSstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSstore [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSSstoreidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSSstoreidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSstore [off1] {sym1} (LEAL4 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVSSstoreidx4 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL4 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVSSstoreidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSstore [off] {sym} (ADDL ptr idx) val mem)
	// cond: ptr.Op != OpSB
	// result: (MOVSSstoreidx1 [off] {sym} ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVSSstoreidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSstoreidx1_0(v *Value) bool {
	// match: (MOVSSstoreidx1 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVSSstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVSSstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSstoreidx1 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVSSstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVSSstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVSSstoreidx4_0(v *Value) bool {
	// match: (MOVSSstoreidx4 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVSSstoreidx4 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVSSstoreidx4)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVSSstoreidx4 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVSSstoreidx4 [int64(int32(c+4*d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVSSstoreidx4)
		v.AuxInt = int64(int32(c + 4*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWLSX_0(v *Value) bool {
	b := v.Block
	// match: (MOVWLSX x:(MOVWload [off] {sym} ptr mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVWLSXload <v.Type> [off] {sym} ptr mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVWload {
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
		v0 := b.NewValue0(x.Pos, Op386MOVWLSXload, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVWLSX (ANDLconst [c] x))
	// cond: c & 0x8000 == 0
	// result: (ANDLconst [c & 0x7fff] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		if !(c&0x8000 == 0) {
			break
		}
		v.reset(Op386ANDLconst)
		v.AuxInt = c & 0x7fff
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWLSXload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVWLSXload [off] {sym} ptr (MOVWstore [off2] {sym2} ptr2 x _))
	// cond: sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)
	// result: (MOVWLSX x)
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVWstore {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		_ = v_1.Args[2]
		ptr2 := v_1.Args[0]
		x := v_1.Args[1]
		if !(sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)) {
			break
		}
		v.reset(Op386MOVWLSX)
		v.AddArg(x)
		return true
	}
	// match: (MOVWLSXload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVWLSXload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVWLSXload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWLZX_0(v *Value) bool {
	b := v.Block
	// match: (MOVWLZX x:(MOVWload [off] {sym} ptr mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVWload <v.Type> [off] {sym} ptr mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVWload {
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
		v0 := b.NewValue0(x.Pos, Op386MOVWload, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVWLZX x:(MOVWloadidx1 [off] {sym} ptr idx mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVWloadidx1 <v.Type> [off] {sym} ptr idx mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVWloadidx1 {
			break
		}
		off := x.AuxInt
		sym := x.Aux
		mem := x.Args[2]
		ptr := x.Args[0]
		idx := x.Args[1]
		if !(x.Uses == 1 && clobber(x)) {
			break
		}
		b = x.Block
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVWLZX x:(MOVWloadidx2 [off] {sym} ptr idx mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: @x.Block (MOVWloadidx2 <v.Type> [off] {sym} ptr idx mem)
	for {
		x := v.Args[0]
		if x.Op != Op386MOVWloadidx2 {
			break
		}
		off := x.AuxInt
		sym := x.Aux
		mem := x.Args[2]
		ptr := x.Args[0]
		idx := x.Args[1]
		if !(x.Uses == 1 && clobber(x)) {
			break
		}
		b = x.Block
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx2, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = off
		v0.Aux = sym
		v0.AddArg(ptr)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (MOVWLZX (ANDLconst [c] x))
	// result: (ANDLconst [c & 0xffff] x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386ANDLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ANDLconst)
		v.AuxInt = c & 0xffff
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVWload [off] {sym} ptr (MOVWstore [off2] {sym2} ptr2 x _))
	// cond: sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)
	// result: (MOVWLZX x)
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVWstore {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		_ = v_1.Args[2]
		ptr2 := v_1.Args[0]
		x := v_1.Args[1]
		if !(sym == sym2 && off == off2 && isSamePtr(ptr, ptr2)) {
			break
		}
		v.reset(Op386MOVWLZX)
		v.AddArg(x)
		return true
	}
	// match: (MOVWload [off1] {sym} (ADDLconst [off2] ptr) mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVWload [off1+off2] {sym} ptr mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVWload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [off1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVWload [off1+off2] {mergeSym(sym1,sym2)} base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVWload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVWloadidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVWloadidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [off1] {sym1} (LEAL2 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVWloadidx2 [off1+off2] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL2 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVWloadidx2)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [off] {sym} (ADDL ptr idx) mem)
	// cond: ptr.Op != OpSB
	// result: (MOVWloadidx1 [off] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVWloadidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWload [off] {sym} (SB) _)
	// cond: symIsRO(sym)
	// result: (MOVLconst [int64(read16(sym, off, config.BigEndian))])
	for {
		off := v.AuxInt
		sym := v.Aux
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != OpSB || !(symIsRO(sym)) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = int64(read16(sym, off, config.BigEndian))
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWloadidx1_0(v *Value) bool {
	// match: (MOVWloadidx1 [c] {sym} ptr (SHLLconst [1] idx) mem)
	// result: (MOVWloadidx2 [c] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		idx := v_1.Args[0]
		v.reset(Op386MOVWloadidx2)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx1 [c] {sym} (SHLLconst [1] idx) ptr mem)
	// result: (MOVWloadidx2 [c] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 1 {
			break
		}
		idx := v_0.Args[0]
		ptr := v.Args[1]
		v.reset(Op386MOVWloadidx2)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx1 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVWloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVWloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx1 [c] {sym} idx (ADDLconst [d] ptr) mem)
	// result: (MOVWloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		idx := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		ptr := v_1.Args[0]
		v.reset(Op386MOVWloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx1 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVWloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVWloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx1 [c] {sym} (ADDLconst [d] idx) ptr mem)
	// result: (MOVWloadidx1 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		idx := v_0.Args[0]
		ptr := v.Args[1]
		v.reset(Op386MOVWloadidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWloadidx2_0(v *Value) bool {
	// match: (MOVWloadidx2 [c] {sym} (ADDLconst [d] ptr) idx mem)
	// result: (MOVWloadidx2 [int64(int32(c+d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVWloadidx2)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWloadidx2 [c] {sym} ptr (ADDLconst [d] idx) mem)
	// result: (MOVWloadidx2 [int64(int32(c+2*d))] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVWloadidx2)
		v.AuxInt = int64(int32(c + 2*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstore_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVWstore [off] {sym} ptr (MOVWLSX x) mem)
	// result: (MOVWstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVWLSX {
			break
		}
		x := v_1.Args[0]
		v.reset(Op386MOVWstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} ptr (MOVWLZX x) mem)
	// result: (MOVWstore [off] {sym} ptr x mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVWLZX {
			break
		}
		x := v_1.Args[0]
		v.reset(Op386MOVWstore)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(x)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off1] {sym} (ADDLconst [off2] ptr) val mem)
	// cond: is32Bit(off1+off2)
	// result: (MOVWstore [off1+off2] {sym} ptr val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} ptr (MOVLconst [c]) mem)
	// cond: validOff(off)
	// result: (MOVWstoreconst [makeValAndOff(int64(int16(c)),off)] {sym} ptr mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(validOff(off)) {
			break
		}
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = makeValAndOff(int64(int16(c)), off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVWstore [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVWstore)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off1] {sym1} (LEAL1 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVWstoreidx1 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off1] {sym1} (LEAL2 [off2] {sym2} ptr idx) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MOVWstoreidx2 [off1+off2] {mergeSym(sym1,sym2)} ptr idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL2 {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVWstoreidx2)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [off] {sym} (ADDL ptr idx) val mem)
	// cond: ptr.Op != OpSB
	// result: (MOVWstoreidx1 [off] {sym} ptr idx val mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		val := v.Args[1]
		if !(ptr.Op != OpSB) {
			break
		}
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [i] {s} p (SHRLconst [16] w) x:(MOVWstore [i-2] {s} p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstore [i-2] {s} p w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst || v_1.AuxInt != 16 {
			break
		}
		w := v_1.Args[0]
		x := v.Args[2]
		if x.Op != Op386MOVWstore || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] || w != x.Args[1] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstore)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstore [i] {s} p (SHRLconst [j] w) x:(MOVWstore [i-2] {s} p w0:(SHRLconst [j-16] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstore [i-2] {s} p w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst {
			break
		}
		j := v_1.AuxInt
		w := v_1.Args[0]
		x := v.Args[2]
		if x.Op != Op386MOVWstore || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] {
			break
		}
		w0 := x.Args[1]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-16 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstore)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstoreconst_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MOVWstoreconst [sc] {s} (ADDLconst [off] ptr) mem)
	// cond: ValAndOff(sc).canAdd(off)
	// result: (MOVWstoreconst [ValAndOff(sc).add(off)] {s} ptr mem)
	for {
		sc := v.AuxInt
		s := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off := v_0.AuxInt
		ptr := v_0.Args[0]
		if !(ValAndOff(sc).canAdd(off)) {
			break
		}
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = ValAndOff(sc).add(off)
		v.Aux = s
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconst [sc] {sym1} (LEAL [off] {sym2} ptr) mem)
	// cond: canMergeSym(sym1, sym2) && ValAndOff(sc).canAdd(off) && (ptr.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MOVWstoreconst [ValAndOff(sc).add(off)] {mergeSym(sym1, sym2)} ptr mem)
	for {
		sc := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2) && ValAndOff(sc).canAdd(off) && (ptr.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = ValAndOff(sc).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconst [x] {sym1} (LEAL1 [off] {sym2} ptr idx) mem)
	// cond: canMergeSym(sym1, sym2)
	// result: (MOVWstoreconstidx1 [ValAndOff(x).add(off)] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		x := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL1 {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVWstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconst [x] {sym1} (LEAL2 [off] {sym2} ptr idx) mem)
	// cond: canMergeSym(sym1, sym2)
	// result: (MOVWstoreconstidx2 [ValAndOff(x).add(off)] {mergeSym(sym1,sym2)} ptr idx mem)
	for {
		x := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL2 {
			break
		}
		off := v_0.AuxInt
		sym2 := v_0.Aux
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		if !(canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MOVWstoreconstidx2)
		v.AuxInt = ValAndOff(x).add(off)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconst [x] {sym} (ADDL ptr idx) mem)
	// result: (MOVWstoreconstidx1 [x] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDL {
			break
		}
		idx := v_0.Args[1]
		ptr := v_0.Args[0]
		v.reset(Op386MOVWstoreconstidx1)
		v.AuxInt = x
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconst [c] {s} p x:(MOVWstoreconst [a] {s} p mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 2 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVLstoreconst [makeValAndOff(ValAndOff(a).Val()&0xffff | ValAndOff(c).Val()<<16, ValAndOff(a).Off())] {s} p mem)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		p := v.Args[0]
		x := v.Args[1]
		if x.Op != Op386MOVWstoreconst {
			break
		}
		a := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[1]
		if p != x.Args[0] || !(x.Uses == 1 && ValAndOff(a).Off()+2 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xffff|ValAndOff(c).Val()<<16, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconst [a] {s} p x:(MOVWstoreconst [c] {s} p mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 2 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVLstoreconst [makeValAndOff(ValAndOff(a).Val()&0xffff | ValAndOff(c).Val()<<16, ValAndOff(a).Off())] {s} p mem)
	for {
		a := v.AuxInt
		s := v.Aux
		_ = v.Args[1]
		p := v.Args[0]
		x := v.Args[1]
		if x.Op != Op386MOVWstoreconst {
			break
		}
		c := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[1]
		if p != x.Args[0] || !(x.Uses == 1 && ValAndOff(a).Off()+2 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xffff|ValAndOff(c).Val()<<16, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstoreconstidx1_0(v *Value) bool {
	// match: (MOVWstoreconstidx1 [c] {sym} ptr (SHLLconst [1] idx) mem)
	// result: (MOVWstoreconstidx2 [c] {sym} ptr idx mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		idx := v_1.Args[0]
		v.reset(Op386MOVWstoreconstidx2)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconstidx1 [x] {sym} (ADDLconst [c] ptr) idx mem)
	// result: (MOVWstoreconstidx1 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		c := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVWstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconstidx1 [x] {sym} ptr (ADDLconst [c] idx) mem)
	// result: (MOVWstoreconstidx1 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		c := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVWstoreconstidx1)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconstidx1 [c] {s} p i x:(MOVWstoreconstidx1 [a] {s} p i mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 2 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVLstoreconstidx1 [makeValAndOff(ValAndOff(a).Val()&0xffff | ValAndOff(c).Val()<<16, ValAndOff(a).Off())] {s} p i mem)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		i := v.Args[1]
		x := v.Args[2]
		if x.Op != Op386MOVWstoreconstidx1 {
			break
		}
		a := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] || i != x.Args[1] || !(x.Uses == 1 && ValAndOff(a).Off()+2 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreconstidx1)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xffff|ValAndOff(c).Val()<<16, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v.AddArg(i)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstoreconstidx2_0(v *Value) bool {
	b := v.Block
	// match: (MOVWstoreconstidx2 [x] {sym} (ADDLconst [c] ptr) idx mem)
	// result: (MOVWstoreconstidx2 [ValAndOff(x).add(c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		c := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		v.reset(Op386MOVWstoreconstidx2)
		v.AuxInt = ValAndOff(x).add(c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconstidx2 [x] {sym} ptr (ADDLconst [c] idx) mem)
	// result: (MOVWstoreconstidx2 [ValAndOff(x).add(2*c)] {sym} ptr idx mem)
	for {
		x := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		c := v_1.AuxInt
		idx := v_1.Args[0]
		v.reset(Op386MOVWstoreconstidx2)
		v.AuxInt = ValAndOff(x).add(2 * c)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreconstidx2 [c] {s} p i x:(MOVWstoreconstidx2 [a] {s} p i mem))
	// cond: x.Uses == 1 && ValAndOff(a).Off() + 2 == ValAndOff(c).Off() && clobber(x)
	// result: (MOVLstoreconstidx1 [makeValAndOff(ValAndOff(a).Val()&0xffff | ValAndOff(c).Val()<<16, ValAndOff(a).Off())] {s} p (SHLLconst <i.Type> [1] i) mem)
	for {
		c := v.AuxInt
		s := v.Aux
		_ = v.Args[2]
		p := v.Args[0]
		i := v.Args[1]
		x := v.Args[2]
		if x.Op != Op386MOVWstoreconstidx2 {
			break
		}
		a := x.AuxInt
		if x.Aux != s {
			break
		}
		mem := x.Args[2]
		if p != x.Args[0] || i != x.Args[1] || !(x.Uses == 1 && ValAndOff(a).Off()+2 == ValAndOff(c).Off() && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreconstidx1)
		v.AuxInt = makeValAndOff(ValAndOff(a).Val()&0xffff|ValAndOff(c).Val()<<16, ValAndOff(a).Off())
		v.Aux = s
		v.AddArg(p)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, i.Type)
		v0.AuxInt = 1
		v0.AddArg(i)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstoreidx1_0(v *Value) bool {
	// match: (MOVWstoreidx1 [c] {sym} ptr (SHLLconst [1] idx) val mem)
	// result: (MOVWstoreidx2 [c] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst || v_1.AuxInt != 1 {
			break
		}
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx2)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [c] {sym} (SHLLconst [1] idx) ptr val mem)
	// result: (MOVWstoreidx2 [c] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst || v_0.AuxInt != 1 {
			break
		}
		idx := v_0.Args[0]
		ptr := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx2)
		v.AuxInt = c
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVWstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [c] {sym} idx (ADDLconst [d] ptr) val mem)
	// result: (MOVWstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		idx := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		ptr := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVWstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [c] {sym} (ADDLconst [d] idx) ptr val mem)
	// result: (MOVWstoreidx1 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		idx := v_0.Args[0]
		ptr := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx1)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} p idx (SHRLconst [16] w) x:(MOVWstoreidx1 [i-2] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 16 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} p idx (SHRLconst [16] w) x:(MOVWstoreidx1 [i-2] {s} idx p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 16 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} idx p (SHRLconst [16] w) x:(MOVWstoreidx1 [i-2] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 16 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} idx p (SHRLconst [16] w) x:(MOVWstoreidx1 [i-2] {s} idx p w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 16 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstoreidx1_10(v *Value) bool {
	// match: (MOVWstoreidx1 [i] {s} p idx (SHRLconst [j] w) x:(MOVWstoreidx1 [i-2] {s} p idx w0:(SHRLconst [j-16] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-16 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} p idx (SHRLconst [j] w) x:(MOVWstoreidx1 [i-2] {s} idx p w0:(SHRLconst [j-16] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-16 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} idx p (SHRLconst [j] w) x:(MOVWstoreidx1 [i-2] {s} p idx w0:(SHRLconst [j-16] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-16 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx1 [i] {s} idx p (SHRLconst [j] w) x:(MOVWstoreidx1 [i-2] {s} idx p w0:(SHRLconst [j-16] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p idx w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		idx := v.Args[0]
		p := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx1 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if idx != x.Args[0] || p != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-16 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v.AddArg(idx)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MOVWstoreidx2_0(v *Value) bool {
	b := v.Block
	// match: (MOVWstoreidx2 [c] {sym} (ADDLconst [d] ptr) idx val mem)
	// result: (MOVWstoreidx2 [int64(int32(c+d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		d := v_0.AuxInt
		ptr := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx2)
		v.AuxInt = int64(int32(c + d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx2 [c] {sym} ptr (ADDLconst [d] idx) val mem)
	// result: (MOVWstoreidx2 [int64(int32(c+2*d))] {sym} ptr idx val mem)
	for {
		c := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		d := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		v.reset(Op386MOVWstoreidx2)
		v.AuxInt = int64(int32(c + 2*d))
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx2 [i] {s} p idx (SHRLconst [16] w) x:(MOVWstoreidx2 [i-2] {s} p idx w mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p (SHLLconst <idx.Type> [1] idx) w mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst || v_2.AuxInt != 16 {
			break
		}
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx2 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] || w != x.Args[2] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, idx.Type)
		v0.AuxInt = 1
		v0.AddArg(idx)
		v.AddArg(v0)
		v.AddArg(w)
		v.AddArg(mem)
		return true
	}
	// match: (MOVWstoreidx2 [i] {s} p idx (SHRLconst [j] w) x:(MOVWstoreidx2 [i-2] {s} p idx w0:(SHRLconst [j-16] w) mem))
	// cond: x.Uses == 1 && clobber(x)
	// result: (MOVLstoreidx1 [i-2] {s} p (SHLLconst <idx.Type> [1] idx) w0 mem)
	for {
		i := v.AuxInt
		s := v.Aux
		_ = v.Args[3]
		p := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386SHRLconst {
			break
		}
		j := v_2.AuxInt
		w := v_2.Args[0]
		x := v.Args[3]
		if x.Op != Op386MOVWstoreidx2 || x.AuxInt != i-2 || x.Aux != s {
			break
		}
		mem := x.Args[3]
		if p != x.Args[0] || idx != x.Args[1] {
			break
		}
		w0 := x.Args[2]
		if w0.Op != Op386SHRLconst || w0.AuxInt != j-16 || w != w0.Args[0] || !(x.Uses == 1 && clobber(x)) {
			break
		}
		v.reset(Op386MOVLstoreidx1)
		v.AuxInt = i - 2
		v.Aux = s
		v.AddArg(p)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, idx.Type)
		v0.AuxInt = 1
		v0.AddArg(idx)
		v.AddArg(v0)
		v.AddArg(w0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULL_0(v *Value) bool {
	// match: (MULL x (MOVLconst [c]))
	// result: (MULLconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386MULLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MULL (MOVLconst [c]) x)
	// result: (MULLconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386MULLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (MULL x l:(MOVLload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (MULLload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386MULLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MULL l:(MOVLload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (MULLload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386MULLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MULL x l:(MOVLloadidx4 [off] {sym} ptr idx mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (MULLloadidx4 x [off] {sym} ptr idx mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386MULLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MULL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (MULLloadidx4 x [off] {sym} ptr idx mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386MULLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULLconst_0(v *Value) bool {
	b := v.Block
	// match: (MULLconst [c] (MULLconst [d] x))
	// result: (MULLconst [int64(int32(c * d))] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MULLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386MULLconst)
		v.AuxInt = int64(int32(c * d))
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [-9] x)
	// result: (NEGL (LEAL8 <v.Type> x x))
	for {
		if v.AuxInt != -9 {
			break
		}
		x := v.Args[0]
		v.reset(Op386NEGL)
		v0 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [-5] x)
	// result: (NEGL (LEAL4 <v.Type> x x))
	for {
		if v.AuxInt != -5 {
			break
		}
		x := v.Args[0]
		v.reset(Op386NEGL)
		v0 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [-3] x)
	// result: (NEGL (LEAL2 <v.Type> x x))
	for {
		if v.AuxInt != -3 {
			break
		}
		x := v.Args[0]
		v.reset(Op386NEGL)
		v0 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [-1] x)
	// result: (NEGL x)
	for {
		if v.AuxInt != -1 {
			break
		}
		x := v.Args[0]
		v.reset(Op386NEGL)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [0] _)
	// result: (MOVLconst [0])
	for {
		if v.AuxInt != 0 {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (MULLconst [1] x)
	// result: x
	for {
		if v.AuxInt != 1 {
			break
		}
		x := v.Args[0]
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [3] x)
	// result: (LEAL2 x x)
	for {
		if v.AuxInt != 3 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [5] x)
	// result: (LEAL4 x x)
	for {
		if v.AuxInt != 5 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL4)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [7] x)
	// result: (LEAL2 x (LEAL2 <v.Type> x x))
	for {
		if v.AuxInt != 7 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386MULLconst_10(v *Value) bool {
	b := v.Block
	// match: (MULLconst [9] x)
	// result: (LEAL8 x x)
	for {
		if v.AuxInt != 9 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v.AddArg(x)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [11] x)
	// result: (LEAL2 x (LEAL4 <v.Type> x x))
	for {
		if v.AuxInt != 11 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [13] x)
	// result: (LEAL4 x (LEAL2 <v.Type> x x))
	for {
		if v.AuxInt != 13 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL4)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [19] x)
	// result: (LEAL2 x (LEAL8 <v.Type> x x))
	for {
		if v.AuxInt != 19 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL2)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [21] x)
	// result: (LEAL4 x (LEAL4 <v.Type> x x))
	for {
		if v.AuxInt != 21 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL4)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [25] x)
	// result: (LEAL8 x (LEAL2 <v.Type> x x))
	for {
		if v.AuxInt != 25 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [27] x)
	// result: (LEAL8 (LEAL2 <v.Type> x x) (LEAL2 <v.Type> x x))
	for {
		if v.AuxInt != 27 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v0 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v1.AddArg(x)
		v1.AddArg(x)
		v.AddArg(v1)
		return true
	}
	// match: (MULLconst [37] x)
	// result: (LEAL4 x (LEAL8 <v.Type> x x))
	for {
		if v.AuxInt != 37 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL4)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [41] x)
	// result: (LEAL8 x (LEAL4 <v.Type> x x))
	for {
		if v.AuxInt != 41 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [45] x)
	// result: (LEAL8 (LEAL4 <v.Type> x x) (LEAL4 <v.Type> x x))
	for {
		if v.AuxInt != 45 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v0 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v1.AddArg(x)
		v1.AddArg(x)
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValue386_Op386MULLconst_20(v *Value) bool {
	b := v.Block
	// match: (MULLconst [73] x)
	// result: (LEAL8 x (LEAL8 <v.Type> x x))
	for {
		if v.AuxInt != 73 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [81] x)
	// result: (LEAL8 (LEAL8 <v.Type> x x) (LEAL8 <v.Type> x x))
	for {
		if v.AuxInt != 81 {
			break
		}
		x := v.Args[0]
		v.reset(Op386LEAL8)
		v0 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v1.AddArg(x)
		v1.AddArg(x)
		v.AddArg(v1)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: isPowerOfTwo(c+1) && c >= 15
	// result: (SUBL (SHLLconst <v.Type> [log2(c+1)] x) x)
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(isPowerOfTwo(c+1) && c >= 15) {
			break
		}
		v.reset(Op386SUBL)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, v.Type)
		v0.AuxInt = log2(c + 1)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: isPowerOfTwo(c-1) && c >= 17
	// result: (LEAL1 (SHLLconst <v.Type> [log2(c-1)] x) x)
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(isPowerOfTwo(c-1) && c >= 17) {
			break
		}
		v.reset(Op386LEAL1)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, v.Type)
		v0.AuxInt = log2(c - 1)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: isPowerOfTwo(c-2) && c >= 34
	// result: (LEAL2 (SHLLconst <v.Type> [log2(c-2)] x) x)
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(isPowerOfTwo(c-2) && c >= 34) {
			break
		}
		v.reset(Op386LEAL2)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, v.Type)
		v0.AuxInt = log2(c - 2)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: isPowerOfTwo(c-4) && c >= 68
	// result: (LEAL4 (SHLLconst <v.Type> [log2(c-4)] x) x)
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(isPowerOfTwo(c-4) && c >= 68) {
			break
		}
		v.reset(Op386LEAL4)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, v.Type)
		v0.AuxInt = log2(c - 4)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: isPowerOfTwo(c-8) && c >= 136
	// result: (LEAL8 (SHLLconst <v.Type> [log2(c-8)] x) x)
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(isPowerOfTwo(c-8) && c >= 136) {
			break
		}
		v.reset(Op386LEAL8)
		v0 := b.NewValue0(v.Pos, Op386SHLLconst, v.Type)
		v0.AuxInt = log2(c - 8)
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: c%3 == 0 && isPowerOfTwo(c/3)
	// result: (SHLLconst [log2(c/3)] (LEAL2 <v.Type> x x))
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(c%3 == 0 && isPowerOfTwo(c/3)) {
			break
		}
		v.reset(Op386SHLLconst)
		v.AuxInt = log2(c / 3)
		v0 := b.NewValue0(v.Pos, Op386LEAL2, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: c%5 == 0 && isPowerOfTwo(c/5)
	// result: (SHLLconst [log2(c/5)] (LEAL4 <v.Type> x x))
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(c%5 == 0 && isPowerOfTwo(c/5)) {
			break
		}
		v.reset(Op386SHLLconst)
		v.AuxInt = log2(c / 5)
		v0 := b.NewValue0(v.Pos, Op386LEAL4, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (MULLconst [c] x)
	// cond: c%9 == 0 && isPowerOfTwo(c/9)
	// result: (SHLLconst [log2(c/9)] (LEAL8 <v.Type> x x))
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(c%9 == 0 && isPowerOfTwo(c/9)) {
			break
		}
		v.reset(Op386SHLLconst)
		v.AuxInt = log2(c / 9)
		v0 := b.NewValue0(v.Pos, Op386LEAL8, v.Type)
		v0.AddArg(x)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_Op386MULLconst_30(v *Value) bool {
	// match: (MULLconst [c] (MOVLconst [d]))
	// result: (MOVLconst [int64(int32(c*d))])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = int64(int32(c * d))
		return true
	}
	return false
}
func rewriteValue386_Op386MULLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MULLload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (MULLload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MULLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MULLload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MULLload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MULLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MULLload [off1] {sym1} val (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (MULLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL4 {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		idx := v_1.Args[1]
		ptr := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386MULLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULLloadidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MULLloadidx4 [off1] {sym} val (ADDLconst [off2] base) idx mem)
	// cond: is32Bit(off1+off2)
	// result: (MULLloadidx4 [off1+off2] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MULLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MULLloadidx4 [off1] {sym} val base (ADDLconst [off2] idx) mem)
	// cond: is32Bit(off1+off2*4)
	// result: (MULLloadidx4 [off1+off2*4] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		base := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386ADDLconst {
			break
		}
		off2 := v_2.AuxInt
		idx := v_2.Args[0]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386MULLloadidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (MULLloadidx4 [off1] {sym1} val (LEAL [off2] {sym2} base) idx mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MULLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val base idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MULLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULSD_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MULSD x l:(MOVSDload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (MULSDload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSDload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386MULSDload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MULSD l:(MOVSDload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (MULSDload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVSDload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386MULSDload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULSDload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MULSDload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (MULSDload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MULSDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MULSDload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MULSDload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MULSDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULSS_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MULSS x l:(MOVSSload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (MULSSload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSSload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386MULSSload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (MULSS l:(MOVSSload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (MULSSload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVSSload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386MULSSload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386MULSSload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (MULSSload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (MULSSload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386MULSSload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (MULSSload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (MULSSload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386MULSSload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386NEGL_0(v *Value) bool {
	// match: (NEGL (MOVLconst [c]))
	// result: (MOVLconst [int64(int32(-c))])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = int64(int32(-c))
		return true
	}
	return false
}
func rewriteValue386_Op386NOTL_0(v *Value) bool {
	// match: (NOTL (MOVLconst [c]))
	// result: (MOVLconst [^c])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = ^c
		return true
	}
	return false
}
func rewriteValue386_Op386ORL_0(v *Value) bool {
	// match: (ORL x (MOVLconst [c]))
	// result: (ORLconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ORLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL (MOVLconst [c]) x)
	// result: (ORLconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386ORLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL (SHLLconst [c] x) (SHRLconst [d] x))
	// cond: d == 32-c
	// result: (ROLLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(Op386ROLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL (SHRLconst [d] x) (SHLLconst [c] x))
	// cond: d == 32-c
	// result: (ROLLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(Op386ROLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL <t> (SHLLconst x [c]) (SHRWconst x [d]))
	// cond: c < 16 && d == 16-c && t.Size() == 2
	// result: (ROLWconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRWconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 16 && d == 16-c && t.Size() == 2) {
			break
		}
		v.reset(Op386ROLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL <t> (SHRWconst x [d]) (SHLLconst x [c]))
	// cond: c < 16 && d == 16-c && t.Size() == 2
	// result: (ROLWconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 16 && d == 16-c && t.Size() == 2) {
			break
		}
		v.reset(Op386ROLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL <t> (SHLLconst x [c]) (SHRBconst x [d]))
	// cond: c < 8 && d == 8-c && t.Size() == 1
	// result: (ROLBconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRBconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 8 && d == 8-c && t.Size() == 1) {
			break
		}
		v.reset(Op386ROLBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL <t> (SHRBconst x [d]) (SHLLconst x [c]))
	// cond: c < 8 && d == 8-c && t.Size() == 1
	// result: (ROLBconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRBconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 8 && d == 8-c && t.Size() == 1) {
			break
		}
		v.reset(Op386ROLBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (ORL x l:(MOVLload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ORLload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ORLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (ORL l:(MOVLload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ORLload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ORLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORL_10(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (ORL x l:(MOVLloadidx4 [off] {sym} ptr idx mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ORLloadidx4 x [off] {sym} ptr idx mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ORLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ORL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (ORLloadidx4 x [off] {sym} ptr idx mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386ORLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ORL x x)
	// result: x
	for {
		x := v.Args[1]
		if x != v.Args[0] {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ORL x0:(MOVBload [i0] {s} p mem) s0:(SHLLconst [8] x1:(MOVBload [i1] {s} p mem)))
	// cond: i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWload [i0] {s} p mem)
	for {
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != Op386MOVBload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s0 := v.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBload {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x1.Pos, Op386MOVWload, typ.UInt16)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s0:(SHLLconst [8] x1:(MOVBload [i1] {s} p mem)) x0:(MOVBload [i0] {s} p mem))
	// cond: i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWload [i0] {s} p mem)
	for {
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBload {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := v.Args[1]
		if x0.Op != Op386MOVBload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(x0.Pos, Op386MOVWload, typ.UInt16)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWload [i0] {s} p mem) s0:(SHLLconst [16] x1:(MOVBload [i2] {s} p mem))) s1:(SHLLconst [24] x2:(MOVBload [i3] {s} p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLload [i0] {s} p mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWload {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[1]
		p := x0.Args[0]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBload {
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
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, Op386MOVLload, typ.UInt32)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBload [i2] {s} p mem)) x0:(MOVWload [i0] {s} p mem)) s1:(SHLLconst [24] x2:(MOVBload [i3] {s} p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLload [i0] {s} p mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBload {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[1]
		p := x1.Args[0]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWload {
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
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBload {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[1]
		if p != x2.Args[0] || mem != x2.Args[1] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x2.Pos, Op386MOVLload, typ.UInt32)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBload [i3] {s} p mem)) o0:(ORL x0:(MOVWload [i0] {s} p mem) s0:(SHLLconst [16] x1:(MOVBload [i2] {s} p mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLload [i0] {s} p mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWload {
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
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBload {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[1]
		if p != x1.Args[0] || mem != x1.Args[1] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x1.Pos, Op386MOVLload, typ.UInt32)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBload [i3] {s} p mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBload [i2] {s} p mem)) x0:(MOVWload [i0] {s} p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLload [i0] {s} p mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBload {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[1]
		p := x2.Args[0]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBload {
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
		if x0.Op != Op386MOVWload {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[1]
		if p != x0.Args[0] || mem != x0.Args[1] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(x0.Pos, Op386MOVLload, typ.UInt32)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL x0:(MOVBloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} p idx mem)))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		p := x0.Args[0]
		idx := x0.Args[1]
		s0 := v.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORL_20(v *Value) bool {
	b := v.Block
	// match: (ORL x0:(MOVBloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} p idx mem)))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		idx := x0.Args[0]
		p := x0.Args[1]
		s0 := v.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL x0:(MOVBloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} idx p mem)))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		p := x0.Args[0]
		idx := x0.Args[1]
		s0 := v.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL x0:(MOVBloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} idx p mem)))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		x0 := v.Args[0]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		idx := x0.Args[0]
		p := x0.Args[1]
		s0 := v.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} p idx mem)) x0:(MOVBloadidx1 [i0] {s} p idx mem))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		p := x1.Args[0]
		idx := x1.Args[1]
		x0 := v.Args[1]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} idx p mem)) x0:(MOVBloadidx1 [i0] {s} p idx mem))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		idx := x1.Args[0]
		p := x1.Args[1]
		x0 := v.Args[1]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} p idx mem)) x0:(MOVBloadidx1 [i0] {s} idx p mem))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		p := x1.Args[0]
		idx := x1.Args[1]
		x0 := v.Args[1]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s0:(SHLLconst [8] x1:(MOVBloadidx1 [i1] {s} idx p mem)) x0:(MOVBloadidx1 [i0] {s} idx p mem))
	// cond: i1==i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b,x0,x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)
	// result: @mergePoint(b,x0,x1) (MOVWloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s0 := v.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 8 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i1 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		idx := x1.Args[0]
		p := x1.Args[1]
		x0 := v.Args[1]
		if x0.Op != Op386MOVBloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] || !(i1 == i0+1 && x0.Uses == 1 && x1.Uses == 1 && s0.Uses == 1 && mergePoint(b, x0, x1) != nil && clobber(x0) && clobber(x1) && clobber(s0)) {
			break
		}
		b = mergePoint(b, x0, x1)
		v0 := b.NewValue0(v.Pos, Op386MOVWloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		p := x0.Args[0]
		idx := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		idx := x0.Args[0]
		p := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		p := x0.Args[0]
		idx := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORL_30(v *Value) bool {
	b := v.Block
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		idx := x0.Args[0]
		p := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		p := x1.Args[0]
		idx := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		idx := x1.Args[0]
		p := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		p := x1.Args[0]
		idx := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		idx := x1.Args[0]
		p := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if p != x2.Args[0] || idx != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		p := x0.Args[0]
		idx := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		idx := x0.Args[0]
		p := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		p := x0.Args[0]
		idx := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		s := x0.Aux
		mem := x0.Args[2]
		idx := x0.Args[0]
		p := x0.Args[1]
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		p := x1.Args[0]
		idx := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORL_40(v *Value) bool {
	b := v.Block
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		idx := x1.Args[0]
		p := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		p := x1.Args[0]
		idx := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)) s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		o0 := v.Args[0]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		s := x1.Aux
		mem := x1.Args[2]
		idx := x1.Args[0]
		p := x1.Args[1]
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s1 := v.Args[1]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		if x2.Aux != s {
			break
		}
		_ = x2.Args[2]
		if idx != x2.Args[0] || p != x2.Args[1] || mem != x2.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} p idx mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORL_50(v *Value) bool {
	b := v.Block
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL x0:(MOVWloadidx1 [i0] {s} idx p mem) s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem))))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		x0 := o0.Args[0]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] {
			break
		}
		s0 := o0.Args[1]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} p idx mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if p != x0.Args[0] || idx != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} p idx mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if p != x1.Args[0] || idx != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} p idx mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		p := x2.Args[0]
		idx := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	// match: (ORL s1:(SHLLconst [24] x2:(MOVBloadidx1 [i3] {s} idx p mem)) o0:(ORL s0:(SHLLconst [16] x1:(MOVBloadidx1 [i2] {s} idx p mem)) x0:(MOVWloadidx1 [i0] {s} idx p mem)))
	// cond: i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b,x0,x1,x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)
	// result: @mergePoint(b,x0,x1,x2) (MOVLloadidx1 <v.Type> [i0] {s} p idx mem)
	for {
		_ = v.Args[1]
		s1 := v.Args[0]
		if s1.Op != Op386SHLLconst || s1.AuxInt != 24 {
			break
		}
		x2 := s1.Args[0]
		if x2.Op != Op386MOVBloadidx1 {
			break
		}
		i3 := x2.AuxInt
		s := x2.Aux
		mem := x2.Args[2]
		idx := x2.Args[0]
		p := x2.Args[1]
		o0 := v.Args[1]
		if o0.Op != Op386ORL {
			break
		}
		_ = o0.Args[1]
		s0 := o0.Args[0]
		if s0.Op != Op386SHLLconst || s0.AuxInt != 16 {
			break
		}
		x1 := s0.Args[0]
		if x1.Op != Op386MOVBloadidx1 {
			break
		}
		i2 := x1.AuxInt
		if x1.Aux != s {
			break
		}
		_ = x1.Args[2]
		if idx != x1.Args[0] || p != x1.Args[1] || mem != x1.Args[2] {
			break
		}
		x0 := o0.Args[1]
		if x0.Op != Op386MOVWloadidx1 {
			break
		}
		i0 := x0.AuxInt
		if x0.Aux != s {
			break
		}
		_ = x0.Args[2]
		if idx != x0.Args[0] || p != x0.Args[1] || mem != x0.Args[2] || !(i2 == i0+2 && i3 == i0+3 && x0.Uses == 1 && x1.Uses == 1 && x2.Uses == 1 && s0.Uses == 1 && s1.Uses == 1 && o0.Uses == 1 && mergePoint(b, x0, x1, x2) != nil && clobber(x0) && clobber(x1) && clobber(x2) && clobber(s0) && clobber(s1) && clobber(o0)) {
			break
		}
		b = mergePoint(b, x0, x1, x2)
		v0 := b.NewValue0(v.Pos, Op386MOVLloadidx1, v.Type)
		v.reset(OpCopy)
		v.AddArg(v0)
		v0.AuxInt = i0
		v0.Aux = s
		v0.AddArg(p)
		v0.AddArg(idx)
		v0.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORLconst_0(v *Value) bool {
	// match: (ORLconst [c] x)
	// cond: int32(c)==0
	// result: x
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(int32(c) == 0) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (ORLconst [c] _)
	// cond: int32(c)==-1
	// result: (MOVLconst [-1])
	for {
		c := v.AuxInt
		if !(int32(c) == -1) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = -1
		return true
	}
	// match: (ORLconst [c] (MOVLconst [d]))
	// result: (MOVLconst [c|d])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = c | d
		return true
	}
	return false
}
func rewriteValue386_Op386ORLconstmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ORLconstmodify [valoff1] {sym} (ADDLconst [off2] base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (ORLconstmodify [ValAndOff(valoff1).add(off2)] {sym} base mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386ORLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ORLconstmodify [valoff1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ORLconstmodify [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ORLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORLconstmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ORLconstmodifyidx4 [valoff1] {sym} (ADDLconst [off2] base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (ORLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386ORLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ORLconstmodifyidx4 [valoff1] {sym} base (ADDLconst [off2] idx) mem)
	// cond: ValAndOff(valoff1).canAdd(off2*4)
	// result: (ORLconstmodifyidx4 [ValAndOff(valoff1).add(off2*4)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2 * 4)) {
			break
		}
		v.reset(Op386ORLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2 * 4)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ORLconstmodifyidx4 [valoff1] {sym1} (LEAL [off2] {sym2} base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ORLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ORLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ORLload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (ORLload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ORLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ORLload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ORLload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ORLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (ORLload [off1] {sym1} val (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (ORLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL4 {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		idx := v_1.Args[1]
		ptr := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386ORLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORLloadidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ORLloadidx4 [off1] {sym} val (ADDLconst [off2] base) idx mem)
	// cond: is32Bit(off1+off2)
	// result: (ORLloadidx4 [off1+off2] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ORLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ORLloadidx4 [off1] {sym} val base (ADDLconst [off2] idx) mem)
	// cond: is32Bit(off1+off2*4)
	// result: (ORLloadidx4 [off1+off2*4] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		base := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386ADDLconst {
			break
		}
		off2 := v_2.AuxInt
		idx := v_2.Args[0]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386ORLloadidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (ORLloadidx4 [off1] {sym1} val (LEAL [off2] {sym2} base) idx mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ORLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val base idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ORLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORLmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ORLmodify [off1] {sym} (ADDLconst [off2] base) val mem)
	// cond: is32Bit(off1+off2)
	// result: (ORLmodify [off1+off2] {sym} base val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ORLmodify)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ORLmodify [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ORLmodify [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ORLmodify)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ORLmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (ORLmodifyidx4 [off1] {sym} (ADDLconst [off2] base) idx val mem)
	// cond: is32Bit(off1+off2)
	// result: (ORLmodifyidx4 [off1+off2] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386ORLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ORLmodifyidx4 [off1] {sym} base (ADDLconst [off2] idx) val mem)
	// cond: is32Bit(off1+off2*4)
	// result: (ORLmodifyidx4 [off1+off2*4] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386ORLmodifyidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ORLmodifyidx4 [off1] {sym1} (LEAL [off2] {sym2} base) idx val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (ORLmodifyidx4 [off1+off2] {mergeSym(sym1,sym2)} base idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386ORLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (ORLmodifyidx4 [off] {sym} ptr idx (MOVLconst [c]) mem)
	// cond: validValAndOff(c,off)
	// result: (ORLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386MOVLconst {
			break
		}
		c := v_2.AuxInt
		if !(validValAndOff(c, off)) {
			break
		}
		v.reset(Op386ORLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386ROLBconst_0(v *Value) bool {
	// match: (ROLBconst [c] (ROLBconst [d] x))
	// result: (ROLBconst [(c+d)& 7] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ROLBconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ROLBconst)
		v.AuxInt = (c + d) & 7
		v.AddArg(x)
		return true
	}
	// match: (ROLBconst [0] x)
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
func rewriteValue386_Op386ROLLconst_0(v *Value) bool {
	// match: (ROLLconst [c] (ROLLconst [d] x))
	// result: (ROLLconst [(c+d)&31] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ROLLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ROLLconst)
		v.AuxInt = (c + d) & 31
		v.AddArg(x)
		return true
	}
	// match: (ROLLconst [0] x)
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
func rewriteValue386_Op386ROLWconst_0(v *Value) bool {
	// match: (ROLWconst [c] (ROLWconst [d] x))
	// result: (ROLWconst [(c+d)&15] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386ROLWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386ROLWconst)
		v.AuxInt = (c + d) & 15
		v.AddArg(x)
		return true
	}
	// match: (ROLWconst [0] x)
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
func rewriteValue386_Op386SARB_0(v *Value) bool {
	// match: (SARB x (MOVLconst [c]))
	// result: (SARBconst [min(c&31,7)] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SARBconst)
		v.AuxInt = min(c&31, 7)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386SARBconst_0(v *Value) bool {
	// match: (SARBconst x [0])
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
	// match: (SARBconst [c] (MOVLconst [d]))
	// result: (MOVLconst [d>>uint64(c)])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = d >> uint64(c)
		return true
	}
	return false
}
func rewriteValue386_Op386SARL_0(v *Value) bool {
	// match: (SARL x (MOVLconst [c]))
	// result: (SARLconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SARLconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	// match: (SARL x (ANDLconst [31] y))
	// result: (SARL x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ANDLconst || v_1.AuxInt != 31 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386SARL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386SARLconst_0(v *Value) bool {
	// match: (SARLconst x [0])
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
	// match: (SARLconst [c] (MOVLconst [d]))
	// result: (MOVLconst [d>>uint64(c)])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = d >> uint64(c)
		return true
	}
	return false
}
func rewriteValue386_Op386SARW_0(v *Value) bool {
	// match: (SARW x (MOVLconst [c]))
	// result: (SARWconst [min(c&31,15)] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SARWconst)
		v.AuxInt = min(c&31, 15)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386SARWconst_0(v *Value) bool {
	// match: (SARWconst x [0])
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
	// match: (SARWconst [c] (MOVLconst [d]))
	// result: (MOVLconst [d>>uint64(c)])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = d >> uint64(c)
		return true
	}
	return false
}
func rewriteValue386_Op386SBBL_0(v *Value) bool {
	// match: (SBBL x (MOVLconst [c]) f)
	// result: (SBBLconst [c] x f)
	for {
		f := v.Args[2]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SBBLconst)
		v.AuxInt = c
		v.AddArg(x)
		v.AddArg(f)
		return true
	}
	return false
}
func rewriteValue386_Op386SBBLcarrymask_0(v *Value) bool {
	// match: (SBBLcarrymask (FlagEQ))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SBBLcarrymask (FlagLT_ULT))
	// result: (MOVLconst [-1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = -1
		return true
	}
	// match: (SBBLcarrymask (FlagLT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SBBLcarrymask (FlagGT_ULT))
	// result: (MOVLconst [-1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = -1
		return true
	}
	// match: (SBBLcarrymask (FlagGT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SETA_0(v *Value) bool {
	// match: (SETA (InvertFlags x))
	// result: (SETB x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETB)
		v.AddArg(x)
		return true
	}
	// match: (SETA (FlagEQ))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETA (FlagLT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETA (FlagLT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETA (FlagGT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETA (FlagGT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	return false
}
func rewriteValue386_Op386SETAE_0(v *Value) bool {
	// match: (SETAE (InvertFlags x))
	// result: (SETBE x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETBE)
		v.AddArg(x)
		return true
	}
	// match: (SETAE (FlagEQ))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETAE (FlagLT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETAE (FlagLT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETAE (FlagGT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETAE (FlagGT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	return false
}
func rewriteValue386_Op386SETB_0(v *Value) bool {
	// match: (SETB (InvertFlags x))
	// result: (SETA x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETA)
		v.AddArg(x)
		return true
	}
	// match: (SETB (FlagEQ))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETB (FlagLT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETB (FlagLT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETB (FlagGT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETB (FlagGT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SETBE_0(v *Value) bool {
	// match: (SETBE (InvertFlags x))
	// result: (SETAE x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETAE)
		v.AddArg(x)
		return true
	}
	// match: (SETBE (FlagEQ))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETBE (FlagLT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETBE (FlagLT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETBE (FlagGT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETBE (FlagGT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SETEQ_0(v *Value) bool {
	// match: (SETEQ (InvertFlags x))
	// result: (SETEQ x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETEQ)
		v.AddArg(x)
		return true
	}
	// match: (SETEQ (FlagEQ))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETEQ (FlagLT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETEQ (FlagLT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETEQ (FlagGT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETEQ (FlagGT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SETG_0(v *Value) bool {
	// match: (SETG (InvertFlags x))
	// result: (SETL x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETL)
		v.AddArg(x)
		return true
	}
	// match: (SETG (FlagEQ))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETG (FlagLT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETG (FlagLT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETG (FlagGT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETG (FlagGT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	return false
}
func rewriteValue386_Op386SETGE_0(v *Value) bool {
	// match: (SETGE (InvertFlags x))
	// result: (SETLE x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETLE)
		v.AddArg(x)
		return true
	}
	// match: (SETGE (FlagEQ))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETGE (FlagLT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETGE (FlagLT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETGE (FlagGT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETGE (FlagGT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	return false
}
func rewriteValue386_Op386SETL_0(v *Value) bool {
	// match: (SETL (InvertFlags x))
	// result: (SETG x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETG)
		v.AddArg(x)
		return true
	}
	// match: (SETL (FlagEQ))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETL (FlagLT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETL (FlagLT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETL (FlagGT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETL (FlagGT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SETLE_0(v *Value) bool {
	// match: (SETLE (InvertFlags x))
	// result: (SETGE x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETGE)
		v.AddArg(x)
		return true
	}
	// match: (SETLE (FlagEQ))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETLE (FlagLT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETLE (FlagLT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETLE (FlagGT_ULT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETLE (FlagGT_UGT))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SETNE_0(v *Value) bool {
	// match: (SETNE (InvertFlags x))
	// result: (SETNE x)
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386InvertFlags {
			break
		}
		x := v_0.Args[0]
		v.reset(Op386SETNE)
		v.AddArg(x)
		return true
	}
	// match: (SETNE (FlagEQ))
	// result: (MOVLconst [0])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagEQ {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	// match: (SETNE (FlagLT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETNE (FlagLT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagLT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETNE (FlagGT_ULT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_ULT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	// match: (SETNE (FlagGT_UGT))
	// result: (MOVLconst [1])
	for {
		v_0 := v.Args[0]
		if v_0.Op != Op386FlagGT_UGT {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 1
		return true
	}
	return false
}
func rewriteValue386_Op386SHLL_0(v *Value) bool {
	// match: (SHLL x (MOVLconst [c]))
	// result: (SHLLconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SHLLconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	// match: (SHLL x (ANDLconst [31] y))
	// result: (SHLL x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ANDLconst || v_1.AuxInt != 31 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386SHLL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386SHLLconst_0(v *Value) bool {
	// match: (SHLLconst x [0])
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
func rewriteValue386_Op386SHRB_0(v *Value) bool {
	// match: (SHRB x (MOVLconst [c]))
	// cond: c&31 < 8
	// result: (SHRBconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(c&31 < 8) {
			break
		}
		v.reset(Op386SHRBconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	// match: (SHRB _ (MOVLconst [c]))
	// cond: c&31 >= 8
	// result: (MOVLconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(c&31 >= 8) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SHRBconst_0(v *Value) bool {
	// match: (SHRBconst x [0])
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
func rewriteValue386_Op386SHRL_0(v *Value) bool {
	// match: (SHRL x (MOVLconst [c]))
	// result: (SHRLconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SHRLconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	// match: (SHRL x (ANDLconst [31] y))
	// result: (SHRL x y)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ANDLconst || v_1.AuxInt != 31 {
			break
		}
		y := v_1.Args[0]
		v.reset(Op386SHRL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
	return false
}
func rewriteValue386_Op386SHRLconst_0(v *Value) bool {
	// match: (SHRLconst x [0])
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
func rewriteValue386_Op386SHRW_0(v *Value) bool {
	// match: (SHRW x (MOVLconst [c]))
	// cond: c&31 < 16
	// result: (SHRWconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(c&31 < 16) {
			break
		}
		v.reset(Op386SHRWconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	// match: (SHRW _ (MOVLconst [c]))
	// cond: c&31 >= 16
	// result: (MOVLconst [0])
	for {
		_ = v.Args[1]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		if !(c&31 >= 16) {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SHRWconst_0(v *Value) bool {
	// match: (SHRWconst x [0])
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
func rewriteValue386_Op386SUBL_0(v *Value) bool {
	b := v.Block
	// match: (SUBL x (MOVLconst [c]))
	// result: (SUBLconst x [c])
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SUBLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (SUBL (MOVLconst [c]) x)
	// result: (NEGL (SUBLconst <v.Type> x [c]))
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386NEGL)
		v0 := b.NewValue0(v.Pos, Op386SUBLconst, v.Type)
		v0.AuxInt = c
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (SUBL x l:(MOVLload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (SUBLload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386SUBLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (SUBL x l:(MOVLloadidx4 [off] {sym} ptr idx mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (SUBLloadidx4 x [off] {sym} ptr idx mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386SUBLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (SUBL x x)
	// result: (MOVLconst [0])
	for {
		x := v.Args[1]
		if x != v.Args[0] {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386SUBLcarry_0(v *Value) bool {
	// match: (SUBLcarry x (MOVLconst [c]))
	// result: (SUBLconstcarry [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386SUBLconstcarry)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBLconst_0(v *Value) bool {
	// match: (SUBLconst [c] x)
	// cond: int32(c) == 0
	// result: x
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(int32(c) == 0) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (SUBLconst [c] x)
	// result: (ADDLconst [int64(int32(-c))] x)
	for {
		c := v.AuxInt
		x := v.Args[0]
		v.reset(Op386ADDLconst)
		v.AuxInt = int64(int32(-c))
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_Op386SUBLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBLload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (SUBLload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386SUBLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (SUBLload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386SUBLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLload [off1] {sym1} val (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (SUBLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL4 {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		idx := v_1.Args[1]
		ptr := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386SUBLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBLloadidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBLloadidx4 [off1] {sym} val (ADDLconst [off2] base) idx mem)
	// cond: is32Bit(off1+off2)
	// result: (SUBLloadidx4 [off1+off2] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386SUBLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLloadidx4 [off1] {sym} val base (ADDLconst [off2] idx) mem)
	// cond: is32Bit(off1+off2*4)
	// result: (SUBLloadidx4 [off1+off2*4] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		base := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386ADDLconst {
			break
		}
		off2 := v_2.AuxInt
		idx := v_2.Args[0]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386SUBLloadidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLloadidx4 [off1] {sym1} val (LEAL [off2] {sym2} base) idx mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (SUBLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val base idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386SUBLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBLmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBLmodify [off1] {sym} (ADDLconst [off2] base) val mem)
	// cond: is32Bit(off1+off2)
	// result: (SUBLmodify [off1+off2] {sym} base val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386SUBLmodify)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLmodify [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (SUBLmodify [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386SUBLmodify)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBLmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBLmodifyidx4 [off1] {sym} (ADDLconst [off2] base) idx val mem)
	// cond: is32Bit(off1+off2)
	// result: (SUBLmodifyidx4 [off1+off2] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386SUBLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLmodifyidx4 [off1] {sym} base (ADDLconst [off2] idx) val mem)
	// cond: is32Bit(off1+off2*4)
	// result: (SUBLmodifyidx4 [off1+off2*4] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386SUBLmodifyidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLmodifyidx4 [off1] {sym1} (LEAL [off2] {sym2} base) idx val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (SUBLmodifyidx4 [off1+off2] {mergeSym(sym1,sym2)} base idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386SUBLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (SUBLmodifyidx4 [off] {sym} ptr idx (MOVLconst [c]) mem)
	// cond: validValAndOff(-c,off)
	// result: (ADDLconstmodifyidx4 [makeValAndOff(-c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386MOVLconst {
			break
		}
		c := v_2.AuxInt
		if !(validValAndOff(-c, off)) {
			break
		}
		v.reset(Op386ADDLconstmodifyidx4)
		v.AuxInt = makeValAndOff(-c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBSD_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBSD x l:(MOVSDload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (SUBSDload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSDload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386SUBSDload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBSDload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBSDload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (SUBSDload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386SUBSDload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (SUBSDload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (SUBSDload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386SUBSDload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBSS_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBSS x l:(MOVSSload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)
	// result: (SUBSSload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVSSload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && !config.use387 && clobber(l)) {
			break
		}
		v.reset(Op386SUBSSload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386SUBSSload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (SUBSSload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (SUBSSload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386SUBSSload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (SUBSSload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (SUBSSload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386SUBSSload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORL_0(v *Value) bool {
	// match: (XORL x (MOVLconst [c]))
	// result: (XORLconst [c] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386XORLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL (MOVLconst [c]) x)
	// result: (XORLconst [c] x)
	for {
		x := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		c := v_0.AuxInt
		v.reset(Op386XORLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL (SHLLconst [c] x) (SHRLconst [d] x))
	// cond: d == 32-c
	// result: (ROLLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRLconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(Op386ROLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL (SHRLconst [d] x) (SHLLconst [c] x))
	// cond: d == 32-c
	// result: (ROLLconst [c] x)
	for {
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(d == 32-c) {
			break
		}
		v.reset(Op386ROLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL <t> (SHLLconst x [c]) (SHRWconst x [d]))
	// cond: c < 16 && d == 16-c && t.Size() == 2
	// result: (ROLWconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRWconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 16 && d == 16-c && t.Size() == 2) {
			break
		}
		v.reset(Op386ROLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL <t> (SHRWconst x [d]) (SHLLconst x [c]))
	// cond: c < 16 && d == 16-c && t.Size() == 2
	// result: (ROLWconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRWconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 16 && d == 16-c && t.Size() == 2) {
			break
		}
		v.reset(Op386ROLWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL <t> (SHLLconst x [c]) (SHRBconst x [d]))
	// cond: c < 8 && d == 8-c && t.Size() == 1
	// result: (ROLBconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHLLconst {
			break
		}
		c := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHRBconst {
			break
		}
		d := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 8 && d == 8-c && t.Size() == 1) {
			break
		}
		v.reset(Op386ROLBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL <t> (SHRBconst x [d]) (SHLLconst x [c]))
	// cond: c < 8 && d == 8-c && t.Size() == 1
	// result: (ROLBconst x [c])
	for {
		t := v.Type
		_ = v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386SHRBconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386SHLLconst {
			break
		}
		c := v_1.AuxInt
		if x != v_1.Args[0] || !(c < 8 && d == 8-c && t.Size() == 1) {
			break
		}
		v.reset(Op386ROLBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (XORL x l:(MOVLload [off] {sym} ptr mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (XORLload x [off] {sym} ptr mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386XORLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (XORL l:(MOVLload [off] {sym} ptr mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (XORLload x [off] {sym} ptr mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLload {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[1]
		ptr := l.Args[0]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386XORLload)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORL_10(v *Value) bool {
	// match: (XORL x l:(MOVLloadidx4 [off] {sym} ptr idx mem))
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (XORLloadidx4 x [off] {sym} ptr idx mem)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		l := v.Args[1]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386XORLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (XORL l:(MOVLloadidx4 [off] {sym} ptr idx mem) x)
	// cond: canMergeLoadClobber(v, l, x) && clobber(l)
	// result: (XORLloadidx4 x [off] {sym} ptr idx mem)
	for {
		x := v.Args[1]
		l := v.Args[0]
		if l.Op != Op386MOVLloadidx4 {
			break
		}
		off := l.AuxInt
		sym := l.Aux
		mem := l.Args[2]
		ptr := l.Args[0]
		idx := l.Args[1]
		if !(canMergeLoadClobber(v, l, x) && clobber(l)) {
			break
		}
		v.reset(Op386XORLloadidx4)
		v.AuxInt = off
		v.Aux = sym
		v.AddArg(x)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (XORL x x)
	// result: (MOVLconst [0])
	for {
		x := v.Args[1]
		if x != v.Args[0] {
			break
		}
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_Op386XORLconst_0(v *Value) bool {
	// match: (XORLconst [c] (XORLconst [d] x))
	// result: (XORLconst [c ^ d] x)
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386XORLconst {
			break
		}
		d := v_0.AuxInt
		x := v_0.Args[0]
		v.reset(Op386XORLconst)
		v.AuxInt = c ^ d
		v.AddArg(x)
		return true
	}
	// match: (XORLconst [c] x)
	// cond: int32(c)==0
	// result: x
	for {
		c := v.AuxInt
		x := v.Args[0]
		if !(int32(c) == 0) {
			break
		}
		v.reset(OpCopy)
		v.Type = x.Type
		v.AddArg(x)
		return true
	}
	// match: (XORLconst [c] (MOVLconst [d]))
	// result: (MOVLconst [c^d])
	for {
		c := v.AuxInt
		v_0 := v.Args[0]
		if v_0.Op != Op386MOVLconst {
			break
		}
		d := v_0.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = c ^ d
		return true
	}
	return false
}
func rewriteValue386_Op386XORLconstmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (XORLconstmodify [valoff1] {sym} (ADDLconst [off2] base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (XORLconstmodify [ValAndOff(valoff1).add(off2)] {sym} base mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386XORLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (XORLconstmodify [valoff1] {sym1} (LEAL [off2] {sym2} base) mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (XORLconstmodify [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[1]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386XORLconstmodify)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORLconstmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (XORLconstmodifyidx4 [valoff1] {sym} (ADDLconst [off2] base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2)
	// result: (XORLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2)) {
			break
		}
		v.reset(Op386XORLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (XORLconstmodifyidx4 [valoff1] {sym} base (ADDLconst [off2] idx) mem)
	// cond: ValAndOff(valoff1).canAdd(off2*4)
	// result: (XORLconstmodifyidx4 [ValAndOff(valoff1).add(off2*4)] {sym} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		if !(ValAndOff(valoff1).canAdd(off2 * 4)) {
			break
		}
		v.reset(Op386XORLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2 * 4)
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (XORLconstmodifyidx4 [valoff1] {sym1} (LEAL [off2] {sym2} base) idx mem)
	// cond: ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (XORLconstmodifyidx4 [ValAndOff(valoff1).add(off2)] {mergeSym(sym1,sym2)} base idx mem)
	for {
		valoff1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		if !(ValAndOff(valoff1).canAdd(off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386XORLconstmodifyidx4)
		v.AuxInt = ValAndOff(valoff1).add(off2)
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORLload_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (XORLload [off1] {sym} val (ADDLconst [off2] base) mem)
	// cond: is32Bit(off1+off2)
	// result: (XORLload [off1+off2] {sym} val base mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386XORLload)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (XORLload [off1] {sym1} val (LEAL [off2] {sym2} base) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (XORLload [off1+off2] {mergeSym(sym1,sym2)} val base mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386XORLload)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(mem)
		return true
	}
	// match: (XORLload [off1] {sym1} val (LEAL4 [off2] {sym2} ptr idx) mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2)
	// result: (XORLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val ptr idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL4 {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		idx := v_1.Args[1]
		ptr := v_1.Args[0]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2)) {
			break
		}
		v.reset(Op386XORLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORLloadidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (XORLloadidx4 [off1] {sym} val (ADDLconst [off2] base) idx mem)
	// cond: is32Bit(off1+off2)
	// result: (XORLloadidx4 [off1+off2] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386XORLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (XORLloadidx4 [off1] {sym} val base (ADDLconst [off2] idx) mem)
	// cond: is32Bit(off1+off2*4)
	// result: (XORLloadidx4 [off1+off2*4] {sym} val base idx mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		base := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386ADDLconst {
			break
		}
		off2 := v_2.AuxInt
		idx := v_2.Args[0]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386XORLloadidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	// match: (XORLloadidx4 [off1] {sym1} val (LEAL [off2] {sym2} base) idx mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (XORLloadidx4 [off1+off2] {mergeSym(sym1,sym2)} val base idx mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		val := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386LEAL {
			break
		}
		off2 := v_1.AuxInt
		sym2 := v_1.Aux
		base := v_1.Args[0]
		idx := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386XORLloadidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(val)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORLmodify_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (XORLmodify [off1] {sym} (ADDLconst [off2] base) val mem)
	// cond: is32Bit(off1+off2)
	// result: (XORLmodify [off1+off2] {sym} base val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386XORLmodify)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (XORLmodify [off1] {sym1} (LEAL [off2] {sym2} base) val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (XORLmodify [off1+off2] {mergeSym(sym1,sym2)} base val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[2]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		val := v.Args[1]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386XORLmodify)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_Op386XORLmodifyidx4_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	// match: (XORLmodifyidx4 [off1] {sym} (ADDLconst [off2] base) idx val mem)
	// cond: is32Bit(off1+off2)
	// result: (XORLmodifyidx4 [off1+off2] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386ADDLconst {
			break
		}
		off2 := v_0.AuxInt
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1 + off2)) {
			break
		}
		v.reset(Op386XORLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (XORLmodifyidx4 [off1] {sym} base (ADDLconst [off2] idx) val mem)
	// cond: is32Bit(off1+off2*4)
	// result: (XORLmodifyidx4 [off1+off2*4] {sym} base idx val mem)
	for {
		off1 := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		base := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386ADDLconst {
			break
		}
		off2 := v_1.AuxInt
		idx := v_1.Args[0]
		val := v.Args[2]
		if !(is32Bit(off1 + off2*4)) {
			break
		}
		v.reset(Op386XORLmodifyidx4)
		v.AuxInt = off1 + off2*4
		v.Aux = sym
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (XORLmodifyidx4 [off1] {sym1} (LEAL [off2] {sym2} base) idx val mem)
	// cond: is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)
	// result: (XORLmodifyidx4 [off1+off2] {mergeSym(sym1,sym2)} base idx val mem)
	for {
		off1 := v.AuxInt
		sym1 := v.Aux
		mem := v.Args[3]
		v_0 := v.Args[0]
		if v_0.Op != Op386LEAL {
			break
		}
		off2 := v_0.AuxInt
		sym2 := v_0.Aux
		base := v_0.Args[0]
		idx := v.Args[1]
		val := v.Args[2]
		if !(is32Bit(off1+off2) && canMergeSym(sym1, sym2) && (base.Op != OpSB || !config.ctxt.Flag_shared)) {
			break
		}
		v.reset(Op386XORLmodifyidx4)
		v.AuxInt = off1 + off2
		v.Aux = mergeSym(sym1, sym2)
		v.AddArg(base)
		v.AddArg(idx)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (XORLmodifyidx4 [off] {sym} ptr idx (MOVLconst [c]) mem)
	// cond: validValAndOff(c,off)
	// result: (XORLconstmodifyidx4 [makeValAndOff(c,off)] {sym} ptr idx mem)
	for {
		off := v.AuxInt
		sym := v.Aux
		mem := v.Args[3]
		ptr := v.Args[0]
		idx := v.Args[1]
		v_2 := v.Args[2]
		if v_2.Op != Op386MOVLconst {
			break
		}
		c := v_2.AuxInt
		if !(validValAndOff(c, off)) {
			break
		}
		v.reset(Op386XORLconstmodifyidx4)
		v.AuxInt = makeValAndOff(c, off)
		v.Aux = sym
		v.AddArg(ptr)
		v.AddArg(idx)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpAdd16_0(v *Value) bool {
	// match: (Add16 x y)
	// result: (ADDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAdd32_0(v *Value) bool {
	// match: (Add32 x y)
	// result: (ADDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAdd32F_0(v *Value) bool {
	// match: (Add32F x y)
	// result: (ADDSS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDSS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAdd32carry_0(v *Value) bool {
	// match: (Add32carry x y)
	// result: (ADDLcarry x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDLcarry)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAdd32withcarry_0(v *Value) bool {
	// match: (Add32withcarry x y c)
	// result: (ADCL x y c)
	for {
		c := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v.reset(Op386ADCL)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(c)
		return true
	}
}
func rewriteValue386_OpAdd64F_0(v *Value) bool {
	// match: (Add64F x y)
	// result: (ADDSD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDSD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAdd8_0(v *Value) bool {
	// match: (Add8 x y)
	// result: (ADDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAddPtr_0(v *Value) bool {
	// match: (AddPtr x y)
	// result: (ADDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ADDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAddr_0(v *Value) bool {
	// match: (Addr {sym} base)
	// result: (LEAL {sym} base)
	for {
		sym := v.Aux
		base := v.Args[0]
		v.reset(Op386LEAL)
		v.Aux = sym
		v.AddArg(base)
		return true
	}
}
func rewriteValue386_OpAnd16_0(v *Value) bool {
	// match: (And16 x y)
	// result: (ANDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAnd32_0(v *Value) bool {
	// match: (And32 x y)
	// result: (ANDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAnd8_0(v *Value) bool {
	// match: (And8 x y)
	// result: (ANDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAndB_0(v *Value) bool {
	// match: (AndB x y)
	// result: (ANDL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpAvg32u_0(v *Value) bool {
	// match: (Avg32u x y)
	// result: (AVGLU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386AVGLU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpBswap32_0(v *Value) bool {
	// match: (Bswap32 x)
	// result: (BSWAPL x)
	for {
		x := v.Args[0]
		v.reset(Op386BSWAPL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpClosureCall_0(v *Value) bool {
	// match: (ClosureCall [argwid] entry closure mem)
	// result: (CALLclosure [argwid] entry closure mem)
	for {
		argwid := v.AuxInt
		mem := v.Args[2]
		entry := v.Args[0]
		closure := v.Args[1]
		v.reset(Op386CALLclosure)
		v.AuxInt = argwid
		v.AddArg(entry)
		v.AddArg(closure)
		v.AddArg(mem)
		return true
	}
}
func rewriteValue386_OpCom16_0(v *Value) bool {
	// match: (Com16 x)
	// result: (NOTL x)
	for {
		x := v.Args[0]
		v.reset(Op386NOTL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCom32_0(v *Value) bool {
	// match: (Com32 x)
	// result: (NOTL x)
	for {
		x := v.Args[0]
		v.reset(Op386NOTL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCom8_0(v *Value) bool {
	// match: (Com8 x)
	// result: (NOTL x)
	for {
		x := v.Args[0]
		v.reset(Op386NOTL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpConst16_0(v *Value) bool {
	// match: (Const16 [val])
	// result: (MOVLconst [val])
	for {
		val := v.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValue386_OpConst32_0(v *Value) bool {
	// match: (Const32 [val])
	// result: (MOVLconst [val])
	for {
		val := v.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValue386_OpConst32F_0(v *Value) bool {
	// match: (Const32F [val])
	// result: (MOVSSconst [val])
	for {
		val := v.AuxInt
		v.reset(Op386MOVSSconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValue386_OpConst64F_0(v *Value) bool {
	// match: (Const64F [val])
	// result: (MOVSDconst [val])
	for {
		val := v.AuxInt
		v.reset(Op386MOVSDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValue386_OpConst8_0(v *Value) bool {
	// match: (Const8 [val])
	// result: (MOVLconst [val])
	for {
		val := v.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValue386_OpConstBool_0(v *Value) bool {
	// match: (ConstBool [b])
	// result: (MOVLconst [b])
	for {
		b := v.AuxInt
		v.reset(Op386MOVLconst)
		v.AuxInt = b
		return true
	}
}
func rewriteValue386_OpConstNil_0(v *Value) bool {
	// match: (ConstNil)
	// result: (MOVLconst [0])
	for {
		v.reset(Op386MOVLconst)
		v.AuxInt = 0
		return true
	}
}
func rewriteValue386_OpCtz16_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Ctz16 x)
	// result: (BSFL (ORLconst <typ.UInt32> [0x10000] x))
	for {
		x := v.Args[0]
		v.reset(Op386BSFL)
		v0 := b.NewValue0(v.Pos, Op386ORLconst, typ.UInt32)
		v0.AuxInt = 0x10000
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpCtz16NonZero_0(v *Value) bool {
	// match: (Ctz16NonZero x)
	// result: (BSFL x)
	for {
		x := v.Args[0]
		v.reset(Op386BSFL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCvt32Fto32_0(v *Value) bool {
	// match: (Cvt32Fto32 x)
	// result: (CVTTSS2SL x)
	for {
		x := v.Args[0]
		v.reset(Op386CVTTSS2SL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCvt32Fto64F_0(v *Value) bool {
	// match: (Cvt32Fto64F x)
	// result: (CVTSS2SD x)
	for {
		x := v.Args[0]
		v.reset(Op386CVTSS2SD)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCvt32to32F_0(v *Value) bool {
	// match: (Cvt32to32F x)
	// result: (CVTSL2SS x)
	for {
		x := v.Args[0]
		v.reset(Op386CVTSL2SS)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCvt32to64F_0(v *Value) bool {
	// match: (Cvt32to64F x)
	// result: (CVTSL2SD x)
	for {
		x := v.Args[0]
		v.reset(Op386CVTSL2SD)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCvt64Fto32_0(v *Value) bool {
	// match: (Cvt64Fto32 x)
	// result: (CVTTSD2SL x)
	for {
		x := v.Args[0]
		v.reset(Op386CVTTSD2SL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpCvt64Fto32F_0(v *Value) bool {
	// match: (Cvt64Fto32F x)
	// result: (CVTSD2SS x)
	for {
		x := v.Args[0]
		v.reset(Op386CVTSD2SS)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpDiv16_0(v *Value) bool {
	// match: (Div16 [a] x y)
	// result: (DIVW [a] x y)
	for {
		a := v.AuxInt
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVW)
		v.AuxInt = a
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpDiv16u_0(v *Value) bool {
	// match: (Div16u x y)
	// result: (DIVWU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVWU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpDiv32_0(v *Value) bool {
	// match: (Div32 [a] x y)
	// result: (DIVL [a] x y)
	for {
		a := v.AuxInt
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVL)
		v.AuxInt = a
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpDiv32F_0(v *Value) bool {
	// match: (Div32F x y)
	// result: (DIVSS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVSS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpDiv32u_0(v *Value) bool {
	// match: (Div32u x y)
	// result: (DIVLU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVLU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpDiv64F_0(v *Value) bool {
	// match: (Div64F x y)
	// result: (DIVSD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVSD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpDiv8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div8 x y)
	// result: (DIVW (SignExt8to16 x) (SignExt8to16 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to16, typ.Int16)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to16, typ.Int16)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpDiv8u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Div8u x y)
	// result: (DIVWU (ZeroExt8to16 x) (ZeroExt8to16 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386DIVWU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to16, typ.UInt16)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to16, typ.UInt16)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpEq16_0(v *Value) bool {
	b := v.Block
	// match: (Eq16 x y)
	// result: (SETEQ (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQ)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpEq32_0(v *Value) bool {
	b := v.Block
	// match: (Eq32 x y)
	// result: (SETEQ (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQ)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpEq32F_0(v *Value) bool {
	b := v.Block
	// match: (Eq32F x y)
	// result: (SETEQF (UCOMISS x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISS, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpEq64F_0(v *Value) bool {
	b := v.Block
	// match: (Eq64F x y)
	// result: (SETEQF (UCOMISD x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISD, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpEq8_0(v *Value) bool {
	b := v.Block
	// match: (Eq8 x y)
	// result: (SETEQ (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQ)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpEqB_0(v *Value) bool {
	b := v.Block
	// match: (EqB x y)
	// result: (SETEQ (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQ)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpEqPtr_0(v *Value) bool {
	b := v.Block
	// match: (EqPtr x y)
	// result: (SETEQ (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETEQ)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq16_0(v *Value) bool {
	b := v.Block
	// match: (Geq16 x y)
	// result: (SETGE (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGE)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq16U_0(v *Value) bool {
	b := v.Block
	// match: (Geq16U x y)
	// result: (SETAE (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETAE)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq32_0(v *Value) bool {
	b := v.Block
	// match: (Geq32 x y)
	// result: (SETGE (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq32F_0(v *Value) bool {
	b := v.Block
	// match: (Geq32F x y)
	// result: (SETGEF (UCOMISS x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGEF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISS, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq32U_0(v *Value) bool {
	b := v.Block
	// match: (Geq32U x y)
	// result: (SETAE (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETAE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq64F_0(v *Value) bool {
	b := v.Block
	// match: (Geq64F x y)
	// result: (SETGEF (UCOMISD x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGEF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISD, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq8_0(v *Value) bool {
	b := v.Block
	// match: (Geq8 x y)
	// result: (SETGE (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGE)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGeq8U_0(v *Value) bool {
	b := v.Block
	// match: (Geq8U x y)
	// result: (SETAE (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETAE)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGetCallerPC_0(v *Value) bool {
	// match: (GetCallerPC)
	// result: (LoweredGetCallerPC)
	for {
		v.reset(Op386LoweredGetCallerPC)
		return true
	}
}
func rewriteValue386_OpGetCallerSP_0(v *Value) bool {
	// match: (GetCallerSP)
	// result: (LoweredGetCallerSP)
	for {
		v.reset(Op386LoweredGetCallerSP)
		return true
	}
}
func rewriteValue386_OpGetClosurePtr_0(v *Value) bool {
	// match: (GetClosurePtr)
	// result: (LoweredGetClosurePtr)
	for {
		v.reset(Op386LoweredGetClosurePtr)
		return true
	}
}
func rewriteValue386_OpGetG_0(v *Value) bool {
	// match: (GetG mem)
	// result: (LoweredGetG mem)
	for {
		mem := v.Args[0]
		v.reset(Op386LoweredGetG)
		v.AddArg(mem)
		return true
	}
}
func rewriteValue386_OpGreater16_0(v *Value) bool {
	b := v.Block
	// match: (Greater16 x y)
	// result: (SETG (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETG)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater16U_0(v *Value) bool {
	b := v.Block
	// match: (Greater16U x y)
	// result: (SETA (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETA)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater32_0(v *Value) bool {
	b := v.Block
	// match: (Greater32 x y)
	// result: (SETG (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETG)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater32F_0(v *Value) bool {
	b := v.Block
	// match: (Greater32F x y)
	// result: (SETGF (UCOMISS x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISS, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater32U_0(v *Value) bool {
	b := v.Block
	// match: (Greater32U x y)
	// result: (SETA (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETA)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater64F_0(v *Value) bool {
	b := v.Block
	// match: (Greater64F x y)
	// result: (SETGF (UCOMISD x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISD, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater8_0(v *Value) bool {
	b := v.Block
	// match: (Greater8 x y)
	// result: (SETG (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETG)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpGreater8U_0(v *Value) bool {
	b := v.Block
	// match: (Greater8U x y)
	// result: (SETA (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETA)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpHmul32_0(v *Value) bool {
	// match: (Hmul32 x y)
	// result: (HMULL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386HMULL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpHmul32u_0(v *Value) bool {
	// match: (Hmul32u x y)
	// result: (HMULLU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386HMULLU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpInterCall_0(v *Value) bool {
	// match: (InterCall [argwid] entry mem)
	// result: (CALLinter [argwid] entry mem)
	for {
		argwid := v.AuxInt
		mem := v.Args[1]
		entry := v.Args[0]
		v.reset(Op386CALLinter)
		v.AuxInt = argwid
		v.AddArg(entry)
		v.AddArg(mem)
		return true
	}
}
func rewriteValue386_OpIsInBounds_0(v *Value) bool {
	b := v.Block
	// match: (IsInBounds idx len)
	// result: (SETB (CMPL idx len))
	for {
		len := v.Args[1]
		idx := v.Args[0]
		v.reset(Op386SETB)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(idx)
		v0.AddArg(len)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpIsNonNil_0(v *Value) bool {
	b := v.Block
	// match: (IsNonNil p)
	// result: (SETNE (TESTL p p))
	for {
		p := v.Args[0]
		v.reset(Op386SETNE)
		v0 := b.NewValue0(v.Pos, Op386TESTL, types.TypeFlags)
		v0.AddArg(p)
		v0.AddArg(p)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpIsSliceInBounds_0(v *Value) bool {
	b := v.Block
	// match: (IsSliceInBounds idx len)
	// result: (SETBE (CMPL idx len))
	for {
		len := v.Args[1]
		idx := v.Args[0]
		v.reset(Op386SETBE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(idx)
		v0.AddArg(len)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq16_0(v *Value) bool {
	b := v.Block
	// match: (Leq16 x y)
	// result: (SETLE (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETLE)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq16U_0(v *Value) bool {
	b := v.Block
	// match: (Leq16U x y)
	// result: (SETBE (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETBE)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq32_0(v *Value) bool {
	b := v.Block
	// match: (Leq32 x y)
	// result: (SETLE (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETLE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq32F_0(v *Value) bool {
	b := v.Block
	// match: (Leq32F x y)
	// result: (SETGEF (UCOMISS y x))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGEF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISS, types.TypeFlags)
		v0.AddArg(y)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq32U_0(v *Value) bool {
	b := v.Block
	// match: (Leq32U x y)
	// result: (SETBE (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETBE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq64F_0(v *Value) bool {
	b := v.Block
	// match: (Leq64F x y)
	// result: (SETGEF (UCOMISD y x))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGEF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISD, types.TypeFlags)
		v0.AddArg(y)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq8_0(v *Value) bool {
	b := v.Block
	// match: (Leq8 x y)
	// result: (SETLE (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETLE)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLeq8U_0(v *Value) bool {
	b := v.Block
	// match: (Leq8U x y)
	// result: (SETBE (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETBE)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess16_0(v *Value) bool {
	b := v.Block
	// match: (Less16 x y)
	// result: (SETL (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETL)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess16U_0(v *Value) bool {
	b := v.Block
	// match: (Less16U x y)
	// result: (SETB (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETB)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess32_0(v *Value) bool {
	b := v.Block
	// match: (Less32 x y)
	// result: (SETL (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETL)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess32F_0(v *Value) bool {
	b := v.Block
	// match: (Less32F x y)
	// result: (SETGF (UCOMISS y x))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISS, types.TypeFlags)
		v0.AddArg(y)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess32U_0(v *Value) bool {
	b := v.Block
	// match: (Less32U x y)
	// result: (SETB (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETB)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess64F_0(v *Value) bool {
	b := v.Block
	// match: (Less64F x y)
	// result: (SETGF (UCOMISD y x))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETGF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISD, types.TypeFlags)
		v0.AddArg(y)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess8_0(v *Value) bool {
	b := v.Block
	// match: (Less8 x y)
	// result: (SETL (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETL)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLess8U_0(v *Value) bool {
	b := v.Block
	// match: (Less8U x y)
	// result: (SETB (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETB)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpLoad_0(v *Value) bool {
	// match: (Load <t> ptr mem)
	// cond: (is32BitInt(t) || isPtr(t))
	// result: (MOVLload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is32BitInt(t) || isPtr(t)) {
			break
		}
		v.reset(Op386MOVLload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is16BitInt(t)
	// result: (MOVWload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is16BitInt(t)) {
			break
		}
		v.reset(Op386MOVWload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: (t.IsBoolean() || is8BitInt(t))
	// result: (MOVBload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(t.IsBoolean() || is8BitInt(t)) {
			break
		}
		v.reset(Op386MOVBload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is32BitFloat(t)
	// result: (MOVSSload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is32BitFloat(t)) {
			break
		}
		v.reset(Op386MOVSSload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is64BitFloat(t)
	// result: (MOVSDload ptr mem)
	for {
		t := v.Type
		mem := v.Args[1]
		ptr := v.Args[0]
		if !(is64BitFloat(t)) {
			break
		}
		v.reset(Op386MOVSDload)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpLocalAddr_0(v *Value) bool {
	// match: (LocalAddr {sym} base _)
	// result: (LEAL {sym} base)
	for {
		sym := v.Aux
		_ = v.Args[1]
		base := v.Args[0]
		v.reset(Op386LEAL)
		v.Aux = sym
		v.AddArg(base)
		return true
	}
}
func rewriteValue386_OpLsh16x16_0(v *Value) bool {
	b := v.Block
	// match: (Lsh16x16 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPWconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh16x32_0(v *Value) bool {
	b := v.Block
	// match: (Lsh16x32 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPLconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh16x64_0(v *Value) bool {
	// match: (Lsh16x64 x (Const64 [c]))
	// cond: uint64(c) < 16
	// result: (SHLLconst x [c])
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
		v.reset(Op386SHLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh16x64 _ (Const64 [c]))
	// cond: uint64(c) >= 16
	// result: (Const16 [0])
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
		v.reset(OpConst16)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_OpLsh16x8_0(v *Value) bool {
	b := v.Block
	// match: (Lsh16x8 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPBconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh32x16_0(v *Value) bool {
	b := v.Block
	// match: (Lsh32x16 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPWconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh32x32_0(v *Value) bool {
	b := v.Block
	// match: (Lsh32x32 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPLconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh32x64_0(v *Value) bool {
	// match: (Lsh32x64 x (Const64 [c]))
	// cond: uint64(c) < 32
	// result: (SHLLconst x [c])
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
		v.reset(Op386SHLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh32x64 _ (Const64 [c]))
	// cond: uint64(c) >= 32
	// result: (Const32 [0])
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
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_OpLsh32x8_0(v *Value) bool {
	b := v.Block
	// match: (Lsh32x8 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPBconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh8x16_0(v *Value) bool {
	b := v.Block
	// match: (Lsh8x16 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPWconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh8x32_0(v *Value) bool {
	b := v.Block
	// match: (Lsh8x32 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPLconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpLsh8x64_0(v *Value) bool {
	// match: (Lsh8x64 x (Const64 [c]))
	// cond: uint64(c) < 8
	// result: (SHLLconst x [c])
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
		v.reset(Op386SHLLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Lsh8x64 _ (Const64 [c]))
	// cond: uint64(c) >= 8
	// result: (Const8 [0])
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
		v.reset(OpConst8)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_OpLsh8x8_0(v *Value) bool {
	b := v.Block
	// match: (Lsh8x8 <t> x y)
	// result: (ANDL (SHLL <t> x y) (SBBLcarrymask <t> (CMPBconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHLL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpMod16_0(v *Value) bool {
	// match: (Mod16 [a] x y)
	// result: (MODW [a] x y)
	for {
		a := v.AuxInt
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MODW)
		v.AuxInt = a
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMod16u_0(v *Value) bool {
	// match: (Mod16u x y)
	// result: (MODWU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MODWU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMod32_0(v *Value) bool {
	// match: (Mod32 [a] x y)
	// result: (MODL [a] x y)
	for {
		a := v.AuxInt
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MODL)
		v.AuxInt = a
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMod32u_0(v *Value) bool {
	// match: (Mod32u x y)
	// result: (MODLU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MODLU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMod8_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod8 x y)
	// result: (MODW (SignExt8to16 x) (SignExt8to16 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MODW)
		v0 := b.NewValue0(v.Pos, OpSignExt8to16, typ.Int16)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpSignExt8to16, typ.Int16)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpMod8u_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Mod8u x y)
	// result: (MODWU (ZeroExt8to16 x) (ZeroExt8to16 y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MODWU)
		v0 := b.NewValue0(v.Pos, OpZeroExt8to16, typ.UInt16)
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpZeroExt8to16, typ.UInt16)
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpMove_0(v *Value) bool {
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
	// result: (MOVBstore dst (MOVBload src mem) mem)
	for {
		if v.AuxInt != 1 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVBstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVBload, typ.UInt8)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [2] dst src mem)
	// result: (MOVWstore dst (MOVWload src mem) mem)
	for {
		if v.AuxInt != 2 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVWstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVWload, typ.UInt16)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [4] dst src mem)
	// result: (MOVLstore dst (MOVLload src mem) mem)
	for {
		if v.AuxInt != 4 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVLstore)
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Move [3] dst src mem)
	// result: (MOVBstore [2] dst (MOVBload [2] src mem) (MOVWstore dst (MOVWload src mem) mem))
	for {
		if v.AuxInt != 3 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVBstore)
		v.AuxInt = 2
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVBload, typ.UInt8)
		v0.AuxInt = 2
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVWstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, Op386MOVWload, typ.UInt16)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [5] dst src mem)
	// result: (MOVBstore [4] dst (MOVBload [4] src mem) (MOVLstore dst (MOVLload src mem) mem))
	for {
		if v.AuxInt != 5 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVBstore)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVBload, typ.UInt8)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVLstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [6] dst src mem)
	// result: (MOVWstore [4] dst (MOVWload [4] src mem) (MOVLstore dst (MOVLload src mem) mem))
	for {
		if v.AuxInt != 6 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVWstore)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVWload, typ.UInt16)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVLstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [7] dst src mem)
	// result: (MOVLstore [3] dst (MOVLload [3] src mem) (MOVLstore dst (MOVLload src mem) mem))
	for {
		if v.AuxInt != 7 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVLstore)
		v.AuxInt = 3
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v0.AuxInt = 3
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVLstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [8] dst src mem)
	// result: (MOVLstore [4] dst (MOVLload [4] src mem) (MOVLstore dst (MOVLload src mem) mem))
	for {
		if v.AuxInt != 8 {
			break
		}
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		v.reset(Op386MOVLstore)
		v.AuxInt = 4
		v.AddArg(dst)
		v0 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v0.AuxInt = 4
		v0.AddArg(src)
		v0.AddArg(mem)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVLstore, types.TypeMem)
		v1.AddArg(dst)
		v2 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v2.AddArg(src)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Move [s] dst src mem)
	// cond: s > 8 && s%4 != 0
	// result: (Move [s-s%4] (ADDLconst <dst.Type> dst [s%4]) (ADDLconst <src.Type> src [s%4]) (MOVLstore dst (MOVLload src mem) mem))
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s > 8 && s%4 != 0) {
			break
		}
		v.reset(OpMove)
		v.AuxInt = s - s%4
		v0 := b.NewValue0(v.Pos, Op386ADDLconst, dst.Type)
		v0.AuxInt = s % 4
		v0.AddArg(dst)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386ADDLconst, src.Type)
		v1.AuxInt = s % 4
		v1.AddArg(src)
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, Op386MOVLstore, types.TypeMem)
		v2.AddArg(dst)
		v3 := b.NewValue0(v.Pos, Op386MOVLload, typ.UInt32)
		v3.AddArg(src)
		v3.AddArg(mem)
		v2.AddArg(v3)
		v2.AddArg(mem)
		v.AddArg(v2)
		return true
	}
	return false
}
func rewriteValue386_OpMove_10(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (Move [s] dst src mem)
	// cond: s > 8 && s <= 4*128 && s%4 == 0 && !config.noDuffDevice
	// result: (DUFFCOPY [10*(128-s/4)] dst src mem)
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !(s > 8 && s <= 4*128 && s%4 == 0 && !config.noDuffDevice) {
			break
		}
		v.reset(Op386DUFFCOPY)
		v.AuxInt = 10 * (128 - s/4)
		v.AddArg(dst)
		v.AddArg(src)
		v.AddArg(mem)
		return true
	}
	// match: (Move [s] dst src mem)
	// cond: (s > 4*128 || config.noDuffDevice) && s%4 == 0
	// result: (REPMOVSL dst src (MOVLconst [s/4]) mem)
	for {
		s := v.AuxInt
		mem := v.Args[2]
		dst := v.Args[0]
		src := v.Args[1]
		if !((s > 4*128 || config.noDuffDevice) && s%4 == 0) {
			break
		}
		v.reset(Op386REPMOVSL)
		v.AddArg(dst)
		v.AddArg(src)
		v0 := b.NewValue0(v.Pos, Op386MOVLconst, typ.UInt32)
		v0.AuxInt = s / 4
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpMul16_0(v *Value) bool {
	// match: (Mul16 x y)
	// result: (MULL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MULL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMul32_0(v *Value) bool {
	// match: (Mul32 x y)
	// result: (MULL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MULL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMul32F_0(v *Value) bool {
	// match: (Mul32F x y)
	// result: (MULSS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MULSS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMul32uhilo_0(v *Value) bool {
	// match: (Mul32uhilo x y)
	// result: (MULLQU x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MULLQU)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMul64F_0(v *Value) bool {
	// match: (Mul64F x y)
	// result: (MULSD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MULSD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpMul8_0(v *Value) bool {
	// match: (Mul8 x y)
	// result: (MULL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386MULL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpNeg16_0(v *Value) bool {
	// match: (Neg16 x)
	// result: (NEGL x)
	for {
		x := v.Args[0]
		v.reset(Op386NEGL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpNeg32_0(v *Value) bool {
	// match: (Neg32 x)
	// result: (NEGL x)
	for {
		x := v.Args[0]
		v.reset(Op386NEGL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpNeg32F_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (Neg32F x)
	// cond: !config.use387
	// result: (PXOR x (MOVSSconst <typ.Float32> [auxFrom32F(float32(math.Copysign(0, -1)))]))
	for {
		x := v.Args[0]
		if !(!config.use387) {
			break
		}
		v.reset(Op386PXOR)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386MOVSSconst, typ.Float32)
		v0.AuxInt = auxFrom32F(float32(math.Copysign(0, -1)))
		v.AddArg(v0)
		return true
	}
	// match: (Neg32F x)
	// cond: config.use387
	// result: (FCHS x)
	for {
		x := v.Args[0]
		if !(config.use387) {
			break
		}
		v.reset(Op386FCHS)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpNeg64F_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (Neg64F x)
	// cond: !config.use387
	// result: (PXOR x (MOVSDconst <typ.Float64> [auxFrom64F(math.Copysign(0, -1))]))
	for {
		x := v.Args[0]
		if !(!config.use387) {
			break
		}
		v.reset(Op386PXOR)
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386MOVSDconst, typ.Float64)
		v0.AuxInt = auxFrom64F(math.Copysign(0, -1))
		v.AddArg(v0)
		return true
	}
	// match: (Neg64F x)
	// cond: config.use387
	// result: (FCHS x)
	for {
		x := v.Args[0]
		if !(config.use387) {
			break
		}
		v.reset(Op386FCHS)
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpNeg8_0(v *Value) bool {
	// match: (Neg8 x)
	// result: (NEGL x)
	for {
		x := v.Args[0]
		v.reset(Op386NEGL)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpNeq16_0(v *Value) bool {
	b := v.Block
	// match: (Neq16 x y)
	// result: (SETNE (CMPW x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNE)
		v0 := b.NewValue0(v.Pos, Op386CMPW, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNeq32_0(v *Value) bool {
	b := v.Block
	// match: (Neq32 x y)
	// result: (SETNE (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNeq32F_0(v *Value) bool {
	b := v.Block
	// match: (Neq32F x y)
	// result: (SETNEF (UCOMISS x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNEF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISS, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNeq64F_0(v *Value) bool {
	b := v.Block
	// match: (Neq64F x y)
	// result: (SETNEF (UCOMISD x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNEF)
		v0 := b.NewValue0(v.Pos, Op386UCOMISD, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNeq8_0(v *Value) bool {
	b := v.Block
	// match: (Neq8 x y)
	// result: (SETNE (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNE)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNeqB_0(v *Value) bool {
	b := v.Block
	// match: (NeqB x y)
	// result: (SETNE (CMPB x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNE)
		v0 := b.NewValue0(v.Pos, Op386CMPB, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNeqPtr_0(v *Value) bool {
	b := v.Block
	// match: (NeqPtr x y)
	// result: (SETNE (CMPL x y))
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SETNE)
		v0 := b.NewValue0(v.Pos, Op386CMPL, types.TypeFlags)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpNilCheck_0(v *Value) bool {
	// match: (NilCheck ptr mem)
	// result: (LoweredNilCheck ptr mem)
	for {
		mem := v.Args[1]
		ptr := v.Args[0]
		v.reset(Op386LoweredNilCheck)
		v.AddArg(ptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValue386_OpNot_0(v *Value) bool {
	// match: (Not x)
	// result: (XORLconst [1] x)
	for {
		x := v.Args[0]
		v.reset(Op386XORLconst)
		v.AuxInt = 1
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpOffPtr_0(v *Value) bool {
	// match: (OffPtr [off] ptr)
	// result: (ADDLconst [off] ptr)
	for {
		off := v.AuxInt
		ptr := v.Args[0]
		v.reset(Op386ADDLconst)
		v.AuxInt = off
		v.AddArg(ptr)
		return true
	}
}
func rewriteValue386_OpOr16_0(v *Value) bool {
	// match: (Or16 x y)
	// result: (ORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpOr32_0(v *Value) bool {
	// match: (Or32 x y)
	// result: (ORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpOr8_0(v *Value) bool {
	// match: (Or8 x y)
	// result: (ORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpOrB_0(v *Value) bool {
	// match: (OrB x y)
	// result: (ORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpPanicBounds_0(v *Value) bool {
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
		v.reset(Op386LoweredPanicBoundsA)
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
		v.reset(Op386LoweredPanicBoundsB)
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
		v.reset(Op386LoweredPanicBoundsC)
		v.AuxInt = kind
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpPanicExtend_0(v *Value) bool {
	// match: (PanicExtend [kind] hi lo y mem)
	// cond: boundsABI(kind) == 0
	// result: (LoweredPanicExtendA [kind] hi lo y mem)
	for {
		kind := v.AuxInt
		mem := v.Args[3]
		hi := v.Args[0]
		lo := v.Args[1]
		y := v.Args[2]
		if !(boundsABI(kind) == 0) {
			break
		}
		v.reset(Op386LoweredPanicExtendA)
		v.AuxInt = kind
		v.AddArg(hi)
		v.AddArg(lo)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	// match: (PanicExtend [kind] hi lo y mem)
	// cond: boundsABI(kind) == 1
	// result: (LoweredPanicExtendB [kind] hi lo y mem)
	for {
		kind := v.AuxInt
		mem := v.Args[3]
		hi := v.Args[0]
		lo := v.Args[1]
		y := v.Args[2]
		if !(boundsABI(kind) == 1) {
			break
		}
		v.reset(Op386LoweredPanicExtendB)
		v.AuxInt = kind
		v.AddArg(hi)
		v.AddArg(lo)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	// match: (PanicExtend [kind] hi lo y mem)
	// cond: boundsABI(kind) == 2
	// result: (LoweredPanicExtendC [kind] hi lo y mem)
	for {
		kind := v.AuxInt
		mem := v.Args[3]
		hi := v.Args[0]
		lo := v.Args[1]
		y := v.Args[2]
		if !(boundsABI(kind) == 2) {
			break
		}
		v.reset(Op386LoweredPanicExtendC)
		v.AuxInt = kind
		v.AddArg(hi)
		v.AddArg(lo)
		v.AddArg(y)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpRotateLeft16_0(v *Value) bool {
	// match: (RotateLeft16 x (MOVLconst [c]))
	// result: (ROLWconst [c&15] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ROLWconst)
		v.AuxInt = c & 15
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpRotateLeft32_0(v *Value) bool {
	// match: (RotateLeft32 x (MOVLconst [c]))
	// result: (ROLLconst [c&31] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ROLLconst)
		v.AuxInt = c & 31
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpRotateLeft8_0(v *Value) bool {
	// match: (RotateLeft8 x (MOVLconst [c]))
	// result: (ROLBconst [c&7] x)
	for {
		_ = v.Args[1]
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != Op386MOVLconst {
			break
		}
		c := v_1.AuxInt
		v.reset(Op386ROLBconst)
		v.AuxInt = c & 7
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpRound32F_0(v *Value) bool {
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
func rewriteValue386_OpRound64F_0(v *Value) bool {
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
func rewriteValue386_OpRsh16Ux16_0(v *Value) bool {
	b := v.Block
	// match: (Rsh16Ux16 <t> x y)
	// result: (ANDL (SHRW <t> x y) (SBBLcarrymask <t> (CMPWconst y [16])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRW, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v2.AuxInt = 16
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh16Ux32_0(v *Value) bool {
	b := v.Block
	// match: (Rsh16Ux32 <t> x y)
	// result: (ANDL (SHRW <t> x y) (SBBLcarrymask <t> (CMPLconst y [16])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRW, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v2.AuxInt = 16
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh16Ux64_0(v *Value) bool {
	// match: (Rsh16Ux64 x (Const64 [c]))
	// cond: uint64(c) < 16
	// result: (SHRWconst x [c])
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
		v.reset(Op386SHRWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh16Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 16
	// result: (Const16 [0])
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
		v.reset(OpConst16)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_OpRsh16Ux8_0(v *Value) bool {
	b := v.Block
	// match: (Rsh16Ux8 <t> x y)
	// result: (ANDL (SHRW <t> x y) (SBBLcarrymask <t> (CMPBconst y [16])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRW, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v2.AuxInt = 16
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh16x16_0(v *Value) bool {
	b := v.Block
	// match: (Rsh16x16 <t> x y)
	// result: (SARW <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPWconst y [16])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARW)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v3.AuxInt = 16
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh16x32_0(v *Value) bool {
	b := v.Block
	// match: (Rsh16x32 <t> x y)
	// result: (SARW <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPLconst y [16])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARW)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v3.AuxInt = 16
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh16x64_0(v *Value) bool {
	// match: (Rsh16x64 x (Const64 [c]))
	// cond: uint64(c) < 16
	// result: (SARWconst x [c])
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
		v.reset(Op386SARWconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh16x64 x (Const64 [c]))
	// cond: uint64(c) >= 16
	// result: (SARWconst x [15])
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
		v.reset(Op386SARWconst)
		v.AuxInt = 15
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpRsh16x8_0(v *Value) bool {
	b := v.Block
	// match: (Rsh16x8 <t> x y)
	// result: (SARW <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPBconst y [16])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARW)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v3.AuxInt = 16
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh32Ux16_0(v *Value) bool {
	b := v.Block
	// match: (Rsh32Ux16 <t> x y)
	// result: (ANDL (SHRL <t> x y) (SBBLcarrymask <t> (CMPWconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh32Ux32_0(v *Value) bool {
	b := v.Block
	// match: (Rsh32Ux32 <t> x y)
	// result: (ANDL (SHRL <t> x y) (SBBLcarrymask <t> (CMPLconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh32Ux64_0(v *Value) bool {
	// match: (Rsh32Ux64 x (Const64 [c]))
	// cond: uint64(c) < 32
	// result: (SHRLconst x [c])
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
		v.reset(Op386SHRLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 32
	// result: (Const32 [0])
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
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_OpRsh32Ux8_0(v *Value) bool {
	b := v.Block
	// match: (Rsh32Ux8 <t> x y)
	// result: (ANDL (SHRL <t> x y) (SBBLcarrymask <t> (CMPBconst y [32])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRL, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v2.AuxInt = 32
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh32x16_0(v *Value) bool {
	b := v.Block
	// match: (Rsh32x16 <t> x y)
	// result: (SARL <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPWconst y [32])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARL)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v3.AuxInt = 32
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh32x32_0(v *Value) bool {
	b := v.Block
	// match: (Rsh32x32 <t> x y)
	// result: (SARL <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPLconst y [32])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARL)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v3.AuxInt = 32
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh32x64_0(v *Value) bool {
	// match: (Rsh32x64 x (Const64 [c]))
	// cond: uint64(c) < 32
	// result: (SARLconst x [c])
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
		v.reset(Op386SARLconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x64 x (Const64 [c]))
	// cond: uint64(c) >= 32
	// result: (SARLconst x [31])
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
		v.reset(Op386SARLconst)
		v.AuxInt = 31
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpRsh32x8_0(v *Value) bool {
	b := v.Block
	// match: (Rsh32x8 <t> x y)
	// result: (SARL <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPBconst y [32])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARL)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v3.AuxInt = 32
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh8Ux16_0(v *Value) bool {
	b := v.Block
	// match: (Rsh8Ux16 <t> x y)
	// result: (ANDL (SHRB <t> x y) (SBBLcarrymask <t> (CMPWconst y [8])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRB, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v2.AuxInt = 8
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh8Ux32_0(v *Value) bool {
	b := v.Block
	// match: (Rsh8Ux32 <t> x y)
	// result: (ANDL (SHRB <t> x y) (SBBLcarrymask <t> (CMPLconst y [8])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRB, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v2.AuxInt = 8
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh8Ux64_0(v *Value) bool {
	// match: (Rsh8Ux64 x (Const64 [c]))
	// cond: uint64(c) < 8
	// result: (SHRBconst x [c])
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
		v.reset(Op386SHRBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh8Ux64 _ (Const64 [c]))
	// cond: uint64(c) >= 8
	// result: (Const8 [0])
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
		v.reset(OpConst8)
		v.AuxInt = 0
		return true
	}
	return false
}
func rewriteValue386_OpRsh8Ux8_0(v *Value) bool {
	b := v.Block
	// match: (Rsh8Ux8 <t> x y)
	// result: (ANDL (SHRB <t> x y) (SBBLcarrymask <t> (CMPBconst y [8])))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386ANDL)
		v0 := b.NewValue0(v.Pos, Op386SHRB, t)
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v2 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v2.AuxInt = 8
		v2.AddArg(y)
		v1.AddArg(v2)
		v.AddArg(v1)
		return true
	}
}
func rewriteValue386_OpRsh8x16_0(v *Value) bool {
	b := v.Block
	// match: (Rsh8x16 <t> x y)
	// result: (SARB <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPWconst y [8])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARB)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPWconst, types.TypeFlags)
		v3.AuxInt = 8
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh8x32_0(v *Value) bool {
	b := v.Block
	// match: (Rsh8x32 <t> x y)
	// result: (SARB <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPLconst y [8])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARB)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v3.AuxInt = 8
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpRsh8x64_0(v *Value) bool {
	// match: (Rsh8x64 x (Const64 [c]))
	// cond: uint64(c) < 8
	// result: (SARBconst x [c])
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
		v.reset(Op386SARBconst)
		v.AuxInt = c
		v.AddArg(x)
		return true
	}
	// match: (Rsh8x64 x (Const64 [c]))
	// cond: uint64(c) >= 8
	// result: (SARBconst x [7])
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
		v.reset(Op386SARBconst)
		v.AuxInt = 7
		v.AddArg(x)
		return true
	}
	return false
}
func rewriteValue386_OpRsh8x8_0(v *Value) bool {
	b := v.Block
	// match: (Rsh8x8 <t> x y)
	// result: (SARB <t> x (ORL <y.Type> y (NOTL <y.Type> (SBBLcarrymask <y.Type> (CMPBconst y [8])))))
	for {
		t := v.Type
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SARB)
		v.Type = t
		v.AddArg(x)
		v0 := b.NewValue0(v.Pos, Op386ORL, y.Type)
		v0.AddArg(y)
		v1 := b.NewValue0(v.Pos, Op386NOTL, y.Type)
		v2 := b.NewValue0(v.Pos, Op386SBBLcarrymask, y.Type)
		v3 := b.NewValue0(v.Pos, Op386CMPBconst, types.TypeFlags)
		v3.AuxInt = 8
		v3.AddArg(y)
		v2.AddArg(v3)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpSelect0_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Select0 (Mul32uover x y))
	// result: (Select0 <typ.UInt32> (MULLU x y))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpMul32uover {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(OpSelect0)
		v.Type = typ.UInt32
		v0 := b.NewValue0(v.Pos, Op386MULLU, types.NewTuple(typ.UInt32, types.TypeFlags))
		v0.AddArg(x)
		v0.AddArg(y)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_OpSelect1_0(v *Value) bool {
	b := v.Block
	typ := &b.Func.Config.Types
	// match: (Select1 (Mul32uover x y))
	// result: (SETO (Select1 <types.TypeFlags> (MULLU x y)))
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpMul32uover {
			break
		}
		y := v_0.Args[1]
		x := v_0.Args[0]
		v.reset(Op386SETO)
		v0 := b.NewValue0(v.Pos, OpSelect1, types.TypeFlags)
		v1 := b.NewValue0(v.Pos, Op386MULLU, types.NewTuple(typ.UInt32, types.TypeFlags))
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_OpSignExt16to32_0(v *Value) bool {
	// match: (SignExt16to32 x)
	// result: (MOVWLSX x)
	for {
		x := v.Args[0]
		v.reset(Op386MOVWLSX)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpSignExt8to16_0(v *Value) bool {
	// match: (SignExt8to16 x)
	// result: (MOVBLSX x)
	for {
		x := v.Args[0]
		v.reset(Op386MOVBLSX)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpSignExt8to32_0(v *Value) bool {
	// match: (SignExt8to32 x)
	// result: (MOVBLSX x)
	for {
		x := v.Args[0]
		v.reset(Op386MOVBLSX)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpSignmask_0(v *Value) bool {
	// match: (Signmask x)
	// result: (SARLconst x [31])
	for {
		x := v.Args[0]
		v.reset(Op386SARLconst)
		v.AuxInt = 31
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpSlicemask_0(v *Value) bool {
	b := v.Block
	// match: (Slicemask <t> x)
	// result: (SARLconst (NEGL <t> x) [31])
	for {
		t := v.Type
		x := v.Args[0]
		v.reset(Op386SARLconst)
		v.AuxInt = 31
		v0 := b.NewValue0(v.Pos, Op386NEGL, t)
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValue386_OpSqrt_0(v *Value) bool {
	// match: (Sqrt x)
	// result: (SQRTSD x)
	for {
		x := v.Args[0]
		v.reset(Op386SQRTSD)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpStaticCall_0(v *Value) bool {
	// match: (StaticCall [argwid] {target} mem)
	// result: (CALLstatic [argwid] {target} mem)
	for {
		argwid := v.AuxInt
		target := v.Aux
		mem := v.Args[0]
		v.reset(Op386CALLstatic)
		v.AuxInt = argwid
		v.Aux = target
		v.AddArg(mem)
		return true
	}
}
func rewriteValue386_OpStore_0(v *Value) bool {
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 8 && is64BitFloat(val.Type)
	// result: (MOVSDstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 8 && is64BitFloat(val.Type)) {
			break
		}
		v.reset(Op386MOVSDstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 4 && is32BitFloat(val.Type)
	// result: (MOVSSstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 4 && is32BitFloat(val.Type)) {
			break
		}
		v.reset(Op386MOVSSstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 4
	// result: (MOVLstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 4) {
			break
		}
		v.reset(Op386MOVLstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	// match: (Store {t} ptr val mem)
	// cond: t.(*types.Type).Size() == 2
	// result: (MOVWstore ptr val mem)
	for {
		t := v.Aux
		mem := v.Args[2]
		ptr := v.Args[0]
		val := v.Args[1]
		if !(t.(*types.Type).Size() == 2) {
			break
		}
		v.reset(Op386MOVWstore)
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
		v.reset(Op386MOVBstore)
		v.AddArg(ptr)
		v.AddArg(val)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpSub16_0(v *Value) bool {
	// match: (Sub16 x y)
	// result: (SUBL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpSub32_0(v *Value) bool {
	// match: (Sub32 x y)
	// result: (SUBL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpSub32F_0(v *Value) bool {
	// match: (Sub32F x y)
	// result: (SUBSS x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBSS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpSub32carry_0(v *Value) bool {
	// match: (Sub32carry x y)
	// result: (SUBLcarry x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBLcarry)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpSub32withcarry_0(v *Value) bool {
	// match: (Sub32withcarry x y c)
	// result: (SBBL x y c)
	for {
		c := v.Args[2]
		x := v.Args[0]
		y := v.Args[1]
		v.reset(Op386SBBL)
		v.AddArg(x)
		v.AddArg(y)
		v.AddArg(c)
		return true
	}
}
func rewriteValue386_OpSub64F_0(v *Value) bool {
	// match: (Sub64F x y)
	// result: (SUBSD x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBSD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpSub8_0(v *Value) bool {
	// match: (Sub8 x y)
	// result: (SUBL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpSubPtr_0(v *Value) bool {
	// match: (SubPtr x y)
	// result: (SUBL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386SUBL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpTrunc16to8_0(v *Value) bool {
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
func rewriteValue386_OpTrunc32to16_0(v *Value) bool {
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
func rewriteValue386_OpTrunc32to8_0(v *Value) bool {
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
func rewriteValue386_OpWB_0(v *Value) bool {
	// match: (WB {fn} destptr srcptr mem)
	// result: (LoweredWB {fn} destptr srcptr mem)
	for {
		fn := v.Aux
		mem := v.Args[2]
		destptr := v.Args[0]
		srcptr := v.Args[1]
		v.reset(Op386LoweredWB)
		v.Aux = fn
		v.AddArg(destptr)
		v.AddArg(srcptr)
		v.AddArg(mem)
		return true
	}
}
func rewriteValue386_OpXor16_0(v *Value) bool {
	// match: (Xor16 x y)
	// result: (XORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386XORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpXor32_0(v *Value) bool {
	// match: (Xor32 x y)
	// result: (XORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386XORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpXor8_0(v *Value) bool {
	// match: (Xor8 x y)
	// result: (XORL x y)
	for {
		y := v.Args[1]
		x := v.Args[0]
		v.reset(Op386XORL)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValue386_OpZero_0(v *Value) bool {
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
	// result: (MOVBstoreconst [0] destptr mem)
	for {
		if v.AuxInt != 1 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVBstoreconst)
		v.AuxInt = 0
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [2] destptr mem)
	// result: (MOVWstoreconst [0] destptr mem)
	for {
		if v.AuxInt != 2 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = 0
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [4] destptr mem)
	// result: (MOVLstoreconst [0] destptr mem)
	for {
		if v.AuxInt != 4 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = 0
		v.AddArg(destptr)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [3] destptr mem)
	// result: (MOVBstoreconst [makeValAndOff(0,2)] destptr (MOVWstoreconst [0] destptr mem))
	for {
		if v.AuxInt != 3 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVBstoreconst)
		v.AuxInt = makeValAndOff(0, 2)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVWstoreconst, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [5] destptr mem)
	// result: (MOVBstoreconst [makeValAndOff(0,4)] destptr (MOVLstoreconst [0] destptr mem))
	for {
		if v.AuxInt != 5 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVBstoreconst)
		v.AuxInt = makeValAndOff(0, 4)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [6] destptr mem)
	// result: (MOVWstoreconst [makeValAndOff(0,4)] destptr (MOVLstoreconst [0] destptr mem))
	for {
		if v.AuxInt != 6 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVWstoreconst)
		v.AuxInt = makeValAndOff(0, 4)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [7] destptr mem)
	// result: (MOVLstoreconst [makeValAndOff(0,3)] destptr (MOVLstoreconst [0] destptr mem))
	for {
		if v.AuxInt != 7 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(0, 3)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [s] destptr mem)
	// cond: s%4 != 0 && s > 4
	// result: (Zero [s-s%4] (ADDLconst destptr [s%4]) (MOVLstoreconst [0] destptr mem))
	for {
		s := v.AuxInt
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(s%4 != 0 && s > 4) {
			break
		}
		v.reset(OpZero)
		v.AuxInt = s - s%4
		v0 := b.NewValue0(v.Pos, Op386ADDLconst, typ.UInt32)
		v0.AuxInt = s % 4
		v0.AddArg(destptr)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v1.AuxInt = 0
		v1.AddArg(destptr)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	// match: (Zero [8] destptr mem)
	// result: (MOVLstoreconst [makeValAndOff(0,4)] destptr (MOVLstoreconst [0] destptr mem))
	for {
		if v.AuxInt != 8 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(0, 4)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v0.AuxInt = 0
		v0.AddArg(destptr)
		v0.AddArg(mem)
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteValue386_OpZero_10(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	typ := &b.Func.Config.Types
	// match: (Zero [12] destptr mem)
	// result: (MOVLstoreconst [makeValAndOff(0,8)] destptr (MOVLstoreconst [makeValAndOff(0,4)] destptr (MOVLstoreconst [0] destptr mem)))
	for {
		if v.AuxInt != 12 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(0, 8)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v0.AuxInt = makeValAndOff(0, 4)
		v0.AddArg(destptr)
		v1 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v1.AuxInt = 0
		v1.AddArg(destptr)
		v1.AddArg(mem)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [16] destptr mem)
	// result: (MOVLstoreconst [makeValAndOff(0,12)] destptr (MOVLstoreconst [makeValAndOff(0,8)] destptr (MOVLstoreconst [makeValAndOff(0,4)] destptr (MOVLstoreconst [0] destptr mem))))
	for {
		if v.AuxInt != 16 {
			break
		}
		mem := v.Args[1]
		destptr := v.Args[0]
		v.reset(Op386MOVLstoreconst)
		v.AuxInt = makeValAndOff(0, 12)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v0.AuxInt = makeValAndOff(0, 8)
		v0.AddArg(destptr)
		v1 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v1.AuxInt = makeValAndOff(0, 4)
		v1.AddArg(destptr)
		v2 := b.NewValue0(v.Pos, Op386MOVLstoreconst, types.TypeMem)
		v2.AuxInt = 0
		v2.AddArg(destptr)
		v2.AddArg(mem)
		v1.AddArg(v2)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
	// match: (Zero [s] destptr mem)
	// cond: s > 16 && s <= 4*128 && s%4 == 0 && !config.noDuffDevice
	// result: (DUFFZERO [1*(128-s/4)] destptr (MOVLconst [0]) mem)
	for {
		s := v.AuxInt
		mem := v.Args[1]
		destptr := v.Args[0]
		if !(s > 16 && s <= 4*128 && s%4 == 0 && !config.noDuffDevice) {
			break
		}
		v.reset(Op386DUFFZERO)
		v.AuxInt = 1 * (128 - s/4)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLconst, typ.UInt32)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(mem)
		return true
	}
	// match: (Zero [s] destptr mem)
	// cond: (s > 4*128 || (config.noDuffDevice && s > 16)) && s%4 == 0
	// result: (REPSTOSL destptr (MOVLconst [s/4]) (MOVLconst [0]) mem)
	for {
		s := v.AuxInt
		mem := v.Args[1]
		destptr := v.Args[0]
		if !((s > 4*128 || (config.noDuffDevice && s > 16)) && s%4 == 0) {
			break
		}
		v.reset(Op386REPSTOSL)
		v.AddArg(destptr)
		v0 := b.NewValue0(v.Pos, Op386MOVLconst, typ.UInt32)
		v0.AuxInt = s / 4
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, Op386MOVLconst, typ.UInt32)
		v1.AuxInt = 0
		v.AddArg(v1)
		v.AddArg(mem)
		return true
	}
	return false
}
func rewriteValue386_OpZeroExt16to32_0(v *Value) bool {
	// match: (ZeroExt16to32 x)
	// result: (MOVWLZX x)
	for {
		x := v.Args[0]
		v.reset(Op386MOVWLZX)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpZeroExt8to16_0(v *Value) bool {
	// match: (ZeroExt8to16 x)
	// result: (MOVBLZX x)
	for {
		x := v.Args[0]
		v.reset(Op386MOVBLZX)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpZeroExt8to32_0(v *Value) bool {
	// match: (ZeroExt8to32 x)
	// result: (MOVBLZX x)
	for {
		x := v.Args[0]
		v.reset(Op386MOVBLZX)
		v.AddArg(x)
		return true
	}
}
func rewriteValue386_OpZeromask_0(v *Value) bool {
	b := v.Block
	// match: (Zeromask <t> x)
	// result: (XORLconst [-1] (SBBLcarrymask <t> (CMPLconst x [1])))
	for {
		t := v.Type
		x := v.Args[0]
		v.reset(Op386XORLconst)
		v.AuxInt = -1
		v0 := b.NewValue0(v.Pos, Op386SBBLcarrymask, t)
		v1 := b.NewValue0(v.Pos, Op386CMPLconst, types.TypeFlags)
		v1.AuxInt = 1
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteBlock386(b *Block) bool {
	switch b.Kind {
	case Block386EQ:
		// match: (EQ (InvertFlags cmp) yes no)
		// result: (EQ cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386EQ)
			b.AddControl(cmp)
			return true
		}
		// match: (EQ (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (EQ (FlagLT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (EQ (FlagLT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (EQ (FlagGT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (EQ (FlagGT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
	case Block386GE:
		// match: (GE (InvertFlags cmp) yes no)
		// result: (LE cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386LE)
			b.AddControl(cmp)
			return true
		}
		// match: (GE (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (GE (FlagLT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GE (FlagLT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GE (FlagGT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (GE (FlagGT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			return true
		}
	case Block386GT:
		// match: (GT (InvertFlags cmp) yes no)
		// result: (LT cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386LT)
			b.AddControl(cmp)
			return true
		}
		// match: (GT (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GT (FlagLT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GT (FlagLT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (GT (FlagGT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (GT (FlagGT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			return true
		}
	case BlockIf:
		// match: (If (SETL cmp) yes no)
		// result: (LT cmp yes no)
		for b.Controls[0].Op == Op386SETL {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386LT)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETLE cmp) yes no)
		// result: (LE cmp yes no)
		for b.Controls[0].Op == Op386SETLE {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386LE)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETG cmp) yes no)
		// result: (GT cmp yes no)
		for b.Controls[0].Op == Op386SETG {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386GT)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETGE cmp) yes no)
		// result: (GE cmp yes no)
		for b.Controls[0].Op == Op386SETGE {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386GE)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETEQ cmp) yes no)
		// result: (EQ cmp yes no)
		for b.Controls[0].Op == Op386SETEQ {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386EQ)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETNE cmp) yes no)
		// result: (NE cmp yes no)
		for b.Controls[0].Op == Op386SETNE {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386NE)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETB cmp) yes no)
		// result: (ULT cmp yes no)
		for b.Controls[0].Op == Op386SETB {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386ULT)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETBE cmp) yes no)
		// result: (ULE cmp yes no)
		for b.Controls[0].Op == Op386SETBE {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386ULE)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETA cmp) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386SETA {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETAE cmp) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386SETAE {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETO cmp) yes no)
		// result: (OS cmp yes no)
		for b.Controls[0].Op == Op386SETO {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386OS)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETGF cmp) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386SETGF {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETGEF cmp) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386SETGEF {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETEQF cmp) yes no)
		// result: (EQF cmp yes no)
		for b.Controls[0].Op == Op386SETEQF {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386EQF)
			b.AddControl(cmp)
			return true
		}
		// match: (If (SETNEF cmp) yes no)
		// result: (NEF cmp yes no)
		for b.Controls[0].Op == Op386SETNEF {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386NEF)
			b.AddControl(cmp)
			return true
		}
		// match: (If cond yes no)
		// result: (NE (TESTB cond cond) yes no)
		for {
			cond := b.Controls[0]
			b.Reset(Block386NE)
			v0 := b.NewValue0(cond.Pos, Op386TESTB, types.TypeFlags)
			v0.AddArg(cond)
			v0.AddArg(cond)
			b.AddControl(v0)
			return true
		}
	case Block386LE:
		// match: (LE (InvertFlags cmp) yes no)
		// result: (GE cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386GE)
			b.AddControl(cmp)
			return true
		}
		// match: (LE (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LE (FlagLT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LE (FlagLT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LE (FlagGT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (LE (FlagGT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
	case Block386LT:
		// match: (LT (InvertFlags cmp) yes no)
		// result: (GT cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386GT)
			b.AddControl(cmp)
			return true
		}
		// match: (LT (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (LT (FlagLT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LT (FlagLT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (LT (FlagGT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (LT (FlagGT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
	case Block386NE:
		// match: (NE (TESTB (SETL cmp) (SETL cmp)) yes no)
		// result: (LT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETL {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETL || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386LT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETL cmp) (SETL cmp)) yes no)
		// result: (LT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETL {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETL || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386LT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETLE cmp) (SETLE cmp)) yes no)
		// result: (LE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETLE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETLE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386LE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETLE cmp) (SETLE cmp)) yes no)
		// result: (LE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETLE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETLE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386LE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETG cmp) (SETG cmp)) yes no)
		// result: (GT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETG {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETG || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386GT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETG cmp) (SETG cmp)) yes no)
		// result: (GT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETG {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETG || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386GT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETGE cmp) (SETGE cmp)) yes no)
		// result: (GE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETGE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETGE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386GE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETGE cmp) (SETGE cmp)) yes no)
		// result: (GE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETGE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETGE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386GE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETEQ cmp) (SETEQ cmp)) yes no)
		// result: (EQ cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETEQ {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETEQ || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386EQ)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETEQ cmp) (SETEQ cmp)) yes no)
		// result: (EQ cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETEQ {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETEQ || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386EQ)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETNE cmp) (SETNE cmp)) yes no)
		// result: (NE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETNE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETNE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386NE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETNE cmp) (SETNE cmp)) yes no)
		// result: (NE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETNE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETNE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386NE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETB cmp) (SETB cmp)) yes no)
		// result: (ULT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETB {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETB || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386ULT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETB cmp) (SETB cmp)) yes no)
		// result: (ULT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETB {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETB || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386ULT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETBE cmp) (SETBE cmp)) yes no)
		// result: (ULE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETBE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETBE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386ULE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETBE cmp) (SETBE cmp)) yes no)
		// result: (ULE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETBE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETBE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386ULE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETA cmp) (SETA cmp)) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETA {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETA || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETA cmp) (SETA cmp)) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETA {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETA || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETAE cmp) (SETAE cmp)) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETAE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETAE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETAE cmp) (SETAE cmp)) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETAE {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETAE || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETO cmp) (SETO cmp)) yes no)
		// result: (OS cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETO {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETO || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386OS)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETO cmp) (SETO cmp)) yes no)
		// result: (OS cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETO {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETO || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386OS)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETGF cmp) (SETGF cmp)) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETGF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETGF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETGF cmp) (SETGF cmp)) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETGF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETGF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETGEF cmp) (SETGEF cmp)) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETGEF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETGEF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETGEF cmp) (SETGEF cmp)) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETGEF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETGEF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETEQF cmp) (SETEQF cmp)) yes no)
		// result: (EQF cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETEQF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETEQF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386EQF)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETEQF cmp) (SETEQF cmp)) yes no)
		// result: (EQF cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETEQF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETEQF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386EQF)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETNEF cmp) (SETNEF cmp)) yes no)
		// result: (NEF cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETNEF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETNEF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386NEF)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (TESTB (SETNEF cmp) (SETNEF cmp)) yes no)
		// result: (NEF cmp yes no)
		for b.Controls[0].Op == Op386TESTB {
			v_0 := b.Controls[0]
			_ = v_0.Args[1]
			v_0_0 := v_0.Args[0]
			if v_0_0.Op != Op386SETNEF {
				break
			}
			cmp := v_0_0.Args[0]
			v_0_1 := v_0.Args[1]
			if v_0_1.Op != Op386SETNEF || cmp != v_0_1.Args[0] {
				break
			}
			b.Reset(Block386NEF)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (InvertFlags cmp) yes no)
		// result: (NE cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386NE)
			b.AddControl(cmp)
			return true
		}
		// match: (NE (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (NE (FlagLT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (NE (FlagLT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (NE (FlagGT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (NE (FlagGT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			return true
		}
	case Block386UGE:
		// match: (UGE (InvertFlags cmp) yes no)
		// result: (ULE cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386ULE)
			b.AddControl(cmp)
			return true
		}
		// match: (UGE (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (UGE (FlagLT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (UGE (FlagLT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (UGE (FlagGT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (UGE (FlagGT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			return true
		}
	case Block386UGT:
		// match: (UGT (InvertFlags cmp) yes no)
		// result: (ULT cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386ULT)
			b.AddControl(cmp)
			return true
		}
		// match: (UGT (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (UGT (FlagLT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (UGT (FlagLT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (UGT (FlagGT_ULT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (UGT (FlagGT_UGT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			return true
		}
	case Block386ULE:
		// match: (ULE (InvertFlags cmp) yes no)
		// result: (UGE cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386UGE)
			b.AddControl(cmp)
			return true
		}
		// match: (ULE (FlagEQ) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			return true
		}
		// match: (ULE (FlagLT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (ULE (FlagLT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (ULE (FlagGT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (ULE (FlagGT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
	case Block386ULT:
		// match: (ULT (InvertFlags cmp) yes no)
		// result: (UGT cmp yes no)
		for b.Controls[0].Op == Op386InvertFlags {
			v_0 := b.Controls[0]
			cmp := v_0.Args[0]
			b.Reset(Block386UGT)
			b.AddControl(cmp)
			return true
		}
		// match: (ULT (FlagEQ) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagEQ {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (ULT (FlagLT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagLT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (ULT (FlagLT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagLT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
		// match: (ULT (FlagGT_ULT) yes no)
		// result: (First yes no)
		for b.Controls[0].Op == Op386FlagGT_ULT {
			b.Reset(BlockFirst)
			return true
		}
		// match: (ULT (FlagGT_UGT) yes no)
		// result: (First no yes)
		for b.Controls[0].Op == Op386FlagGT_UGT {
			b.Reset(BlockFirst)
			b.swapSuccessors()
			return true
		}
	}
	return false
}
