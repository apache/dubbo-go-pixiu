// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"cmd/compile/internal/gc"
	"cmd/compile/internal/ssa"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/obj/x86"
	"math"
)

// Generates code for v using 387 instructions.
func ssaGenValue387(s *gc.SSAGenState, v *ssa.Value) {
	// The SSA compiler pretends that it has an SSE backend.
	// If we don't have one of those, we need to translate
	// all the SSE ops to equivalent 387 ops. That's what this
	// function does.

	switch v.Op {
	case ssa.Op386MOVSSconst, ssa.Op386MOVSDconst:
		iv := uint64(v.AuxInt)
		if iv == 0x0000000000000000 { // +0.0
			s.Prog(x86.AFLDZ)
		} else if iv == 0x3ff0000000000000 { // +1.0
			s.Prog(x86.AFLD1)
		} else if iv == 0x8000000000000000 { // -0.0
			s.Prog(x86.AFLDZ)
			s.Prog(x86.AFCHS)
		} else if iv == 0xbff0000000000000 { // -1.0
			s.Prog(x86.AFLD1)
			s.Prog(x86.AFCHS)
		} else if iv == 0x400921fb54442d18 { // +pi
			s.Prog(x86.AFLDPI)
		} else if iv == 0xc00921fb54442d18 { // -pi
			s.Prog(x86.AFLDPI)
			s.Prog(x86.AFCHS)
		} else { // others
			p := s.Prog(loadPush(v.Type))
			p.From.Type = obj.TYPE_FCONST
			p.From.Val = math.Float64frombits(iv)
			p.To.Type = obj.TYPE_REG
			p.To.Reg = x86.REG_F0
		}
		popAndSave(s, v)

	case ssa.Op386MOVSSconst2, ssa.Op386MOVSDconst2:
		p := s.Prog(loadPush(v.Type))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = v.Args[0].Reg()
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)

	case ssa.Op386MOVSSload, ssa.Op386MOVSDload, ssa.Op386MOVSSloadidx1, ssa.Op386MOVSDloadidx1, ssa.Op386MOVSSloadidx4, ssa.Op386MOVSDloadidx8:
		p := s.Prog(loadPush(v.Type))
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = v.Args[0].Reg()
		gc.AddAux(&p.From, v)
		switch v.Op {
		case ssa.Op386MOVSSloadidx1, ssa.Op386MOVSDloadidx1:
			p.From.Scale = 1
			p.From.Index = v.Args[1].Reg()
			if p.From.Index == x86.REG_SP {
				p.From.Reg, p.From.Index = p.From.Index, p.From.Reg
			}
		case ssa.Op386MOVSSloadidx4:
			p.From.Scale = 4
			p.From.Index = v.Args[1].Reg()
		case ssa.Op386MOVSDloadidx8:
			p.From.Scale = 8
			p.From.Index = v.Args[1].Reg()
		}
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)

	case ssa.Op386MOVSSstore, ssa.Op386MOVSDstore:
		// Push to-be-stored value on top of stack.
		push(s, v.Args[1])

		// Pop and store value.
		var op obj.As
		switch v.Op {
		case ssa.Op386MOVSSstore:
			op = x86.AFMOVFP
		case ssa.Op386MOVSDstore:
			op = x86.AFMOVDP
		}
		p := s.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = v.Args[0].Reg()
		gc.AddAux(&p.To, v)

	case ssa.Op386MOVSSstoreidx1, ssa.Op386MOVSDstoreidx1, ssa.Op386MOVSSstoreidx4, ssa.Op386MOVSDstoreidx8:
		push(s, v.Args[2])
		var op obj.As
		switch v.Op {
		case ssa.Op386MOVSSstoreidx1, ssa.Op386MOVSSstoreidx4:
			op = x86.AFMOVFP
		case ssa.Op386MOVSDstoreidx1, ssa.Op386MOVSDstoreidx8:
			op = x86.AFMOVDP
		}
		p := s.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = v.Args[0].Reg()
		gc.AddAux(&p.To, v)
		switch v.Op {
		case ssa.Op386MOVSSstoreidx1, ssa.Op386MOVSDstoreidx1:
			p.To.Scale = 1
			p.To.Index = v.Args[1].Reg()
			if p.To.Index == x86.REG_SP {
				p.To.Reg, p.To.Index = p.To.Index, p.To.Reg
			}
		case ssa.Op386MOVSSstoreidx4:
			p.To.Scale = 4
			p.To.Index = v.Args[1].Reg()
		case ssa.Op386MOVSDstoreidx8:
			p.To.Scale = 8
			p.To.Index = v.Args[1].Reg()
		}

	case ssa.Op386ADDSS, ssa.Op386ADDSD, ssa.Op386SUBSS, ssa.Op386SUBSD,
		ssa.Op386MULSS, ssa.Op386MULSD, ssa.Op386DIVSS, ssa.Op386DIVSD:
		if v.Reg() != v.Args[0].Reg() {
			v.Fatalf("input[0] and output not in same register %s", v.LongString())
		}

		// Push arg1 on top of stack
		push(s, v.Args[1])

		// Set precision if needed.  64 bits is the default.
		switch v.Op {
		case ssa.Op386ADDSS, ssa.Op386SUBSS, ssa.Op386MULSS, ssa.Op386DIVSS:
			// Save AX so we can use it as scratch space.
			p := s.Prog(x86.AMOVL)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = x86.REG_AX
			s.AddrScratch(&p.To)
			// Install a 32-bit version of the control word.
			installControlWord(s, gc.ControlWord32, x86.REG_AX)
			// Restore AX.
			p = s.Prog(x86.AMOVL)
			s.AddrScratch(&p.From)
			p.To.Type = obj.TYPE_REG
			p.To.Reg = x86.REG_AX
		}

		var op obj.As
		switch v.Op {
		case ssa.Op386ADDSS, ssa.Op386ADDSD:
			op = x86.AFADDDP
		case ssa.Op386SUBSS, ssa.Op386SUBSD:
			op = x86.AFSUBDP
		case ssa.Op386MULSS, ssa.Op386MULSD:
			op = x86.AFMULDP
		case ssa.Op386DIVSS, ssa.Op386DIVSD:
			op = x86.AFDIVDP
		}
		p := s.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s.SSEto387[v.Reg()] + 1

		// Restore precision if needed.
		switch v.Op {
		case ssa.Op386ADDSS, ssa.Op386SUBSS, ssa.Op386MULSS, ssa.Op386DIVSS:
			restoreControlWord(s)
		}

	case ssa.Op386UCOMISS, ssa.Op386UCOMISD:
		push(s, v.Args[0])

		// Compare.
		p := s.Prog(x86.AFUCOMP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s.SSEto387[v.Args[1].Reg()] + 1

		// Save AX.
		p = s.Prog(x86.AMOVL)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_AX
		s.AddrScratch(&p.To)

		// Move status word into AX.
		p = s.Prog(x86.AFSTSW)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_AX

		// Then move the flags we need to the integer flags.
		s.Prog(x86.ASAHF)

		// Restore AX.
		p = s.Prog(x86.AMOVL)
		s.AddrScratch(&p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_AX

	case ssa.Op386SQRTSD:
		push(s, v.Args[0])
		s.Prog(x86.AFSQRT)
		popAndSave(s, v)

	case ssa.Op386FCHS:
		push(s, v.Args[0])
		s.Prog(x86.AFCHS)
		popAndSave(s, v)

	case ssa.Op386CVTSL2SS, ssa.Op386CVTSL2SD:
		p := s.Prog(x86.AMOVL)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = v.Args[0].Reg()
		s.AddrScratch(&p.To)
		p = s.Prog(x86.AFMOVL)
		s.AddrScratch(&p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)

	case ssa.Op386CVTTSD2SL, ssa.Op386CVTTSS2SL:
		push(s, v.Args[0])

		// Load control word which truncates (rounds towards zero).
		installControlWord(s, gc.ControlWord64trunc, v.Reg())

		// Now do the conversion.
		p := s.Prog(x86.AFMOVLP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		s.AddrScratch(&p.To)
		p = s.Prog(x86.AMOVL)
		s.AddrScratch(&p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = v.Reg()

		// Restore control word.
		restoreControlWord(s)

	case ssa.Op386CVTSS2SD:
		// float32 -> float64 is a nop
		push(s, v.Args[0])
		popAndSave(s, v)

	case ssa.Op386CVTSD2SS:
		// Round to nearest float32.
		push(s, v.Args[0])
		p := s.Prog(x86.AFMOVFP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		s.AddrScratch(&p.To)
		p = s.Prog(x86.AFMOVF)
		s.AddrScratch(&p.From)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		popAndSave(s, v)

	case ssa.OpLoadReg:
		if !v.Type.IsFloat() {
			ssaGenValue(s, v)
			return
		}
		// Load+push the value we need.
		p := s.Prog(loadPush(v.Type))
		gc.AddrAuto(&p.From, v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		// Move the value to its assigned register.
		popAndSave(s, v)

	case ssa.OpStoreReg:
		if !v.Type.IsFloat() {
			ssaGenValue(s, v)
			return
		}
		push(s, v.Args[0])
		var op obj.As
		switch v.Type.Size() {
		case 4:
			op = x86.AFMOVFP
		case 8:
			op = x86.AFMOVDP
		}
		p := s.Prog(op)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		gc.AddrAuto(&p.To, v)

	case ssa.OpCopy:
		if !v.Type.IsFloat() {
			ssaGenValue(s, v)
			return
		}
		push(s, v.Args[0])
		popAndSave(s, v)

	case ssa.Op386CALLstatic, ssa.Op386CALLclosure, ssa.Op386CALLinter:
		flush387(s) // Calls must empty the FP stack.
		fallthrough // then issue the call as normal
	default:
		ssaGenValue(s, v)
	}
}

// push pushes v onto the floating-point stack.  v must be in a register.
func push(s *gc.SSAGenState, v *ssa.Value) {
	p := s.Prog(x86.AFMOVD)
	p.From.Type = obj.TYPE_REG
	p.From.Reg = s.SSEto387[v.Reg()]
	p.To.Type = obj.TYPE_REG
	p.To.Reg = x86.REG_F0
}

// popAndSave pops a value off of the floating-point stack and stores
// it in the reigster assigned to v.
func popAndSave(s *gc.SSAGenState, v *ssa.Value) {
	r := v.Reg()
	if _, ok := s.SSEto387[r]; ok {
		// Pop value, write to correct register.
		p := s.Prog(x86.AFMOVDP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = s.SSEto387[v.Reg()] + 1
	} else {
		// Don't actually pop value. This 387 register is now the
		// new home for the not-yet-assigned-a-home SSE register.
		// Increase the register mapping of all other registers by one.
		for rSSE, r387 := range s.SSEto387 {
			s.SSEto387[rSSE] = r387 + 1
		}
		s.SSEto387[r] = x86.REG_F0
	}
}

// loadPush returns the opcode for load+push of the given type.
func loadPush(t *types.Type) obj.As {
	if t.Size() == 4 {
		return x86.AFMOVF
	}
	return x86.AFMOVD
}

// flush387 removes all entries from the 387 floating-point stack.
func flush387(s *gc.SSAGenState) {
	for k := range s.SSEto387 {
		p := s.Prog(x86.AFMOVDP)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_F0
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_F0
		delete(s.SSEto387, k)
	}
}

func ssaGenBlock387(s *gc.SSAGenState, b, next *ssa.Block) {
	// Empty the 387's FP stack before the block ends.
	flush387(s)

	ssaGenBlock(s, b, next)
}

// installControlWord saves the current floating-point control
// word and installs a new one loaded from cw.
// scratchReg must be an unused register.
// This call must be paired with restoreControlWord.
// Bytes 4-5 of the scratch space (s.AddrScratch) are used between
// this call and restoreControlWord.
func installControlWord(s *gc.SSAGenState, cw *obj.LSym, scratchReg int16) {
	// Save current control word.
	p := s.Prog(x86.AFSTCW)
	s.AddrScratch(&p.To)
	p.To.Offset += 4

	// Materialize address of new control word.
	// Note: this must be a seperate instruction to handle PIE correctly.
	// See issue 41503.
	p = s.Prog(x86.ALEAL)
	p.From.Type = obj.TYPE_MEM
	p.From.Name = obj.NAME_EXTERN
	p.From.Sym = cw
	p.To.Type = obj.TYPE_REG
	p.To.Reg = scratchReg

	// Load replacement control word.
	p = s.Prog(x86.AFLDCW)
	p.From.Type = obj.TYPE_MEM
	p.From.Reg = scratchReg
}
func restoreControlWord(s *gc.SSAGenState) {
	p := s.Prog(x86.AFLDCW)
	s.AddrScratch(&p.From)
	p.From.Offset += 4
}
