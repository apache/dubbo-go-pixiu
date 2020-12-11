// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import (
	"cmd/compile/internal/types"
	"cmd/internal/src"
	"fmt"
	"math"
	"sort"
	"strings"
)

// A Value represents a value in the SSA representation of the program.
// The ID and Type fields must not be modified. The remainder may be modified
// if they preserve the value of the Value (e.g. changing a (mul 2 x) to an (add x x)).
type Value struct {
	// A unique identifier for the value. For performance we allocate these IDs
	// densely starting at 1.  There is no guarantee that there won't be occasional holes, though.
	ID ID

	// The operation that computes this value. See op.go.
	Op Op

	// The type of this value. Normally this will be a Go type, but there
	// are a few other pseudo-types, see ../types/type.go.
	Type *types.Type

	// Auxiliary info for this value. The type of this information depends on the opcode and type.
	// AuxInt is used for integer values, Aux is used for other values.
	// Floats are stored in AuxInt using math.Float64bits(f).
	// Unused portions of AuxInt are filled by sign-extending the used portion,
	// even if the represented value is unsigned.
	// Users of AuxInt which interpret AuxInt as unsigned (e.g. shifts) must be careful.
	// Use Value.AuxUnsigned to get the zero-extended value of AuxInt.
	AuxInt int64
	Aux    interface{}

	// Arguments of this value
	Args []*Value

	// Containing basic block
	Block *Block

	// Source position
	Pos src.XPos

	// Use count. Each appearance in Value.Args and Block.Controls counts once.
	Uses int32

	// wasm: Value stays on the WebAssembly stack. This value will not get a "register" (WebAssembly variable)
	// nor a slot on Go stack, and the generation of this value is delayed to its use time.
	OnWasmStack bool

	// Storage for the first three args
	argstorage [3]*Value
}

// Examples:
// Opcode          aux   args
//  OpAdd          nil      2
//  OpConst     string      0    string constant
//  OpConst      int64      0    int64 constant
//  OpAddcq      int64      1    amd64 op: v = arg[0] + constant

// short form print. Just v#.
func (v *Value) String() string {
	if v == nil {
		return "nil" // should never happen, but not panicking helps with debugging
	}
	return fmt.Sprintf("v%d", v.ID)
}

func (v *Value) AuxInt8() int8 {
	if opcodeTable[v.Op].auxType != auxInt8 {
		v.Fatalf("op %s doesn't have an int8 aux field", v.Op)
	}
	return int8(v.AuxInt)
}

func (v *Value) AuxInt16() int16 {
	if opcodeTable[v.Op].auxType != auxInt16 {
		v.Fatalf("op %s doesn't have an int16 aux field", v.Op)
	}
	return int16(v.AuxInt)
}

func (v *Value) AuxInt32() int32 {
	if opcodeTable[v.Op].auxType != auxInt32 {
		v.Fatalf("op %s doesn't have an int32 aux field", v.Op)
	}
	return int32(v.AuxInt)
}

// AuxUnsigned returns v.AuxInt as an unsigned value for OpConst*.
// v.AuxInt is always sign-extended to 64 bits, even if the
// represented value is unsigned. This undoes that sign extension.
func (v *Value) AuxUnsigned() uint64 {
	c := v.AuxInt
	switch v.Op {
	case OpConst64:
		return uint64(c)
	case OpConst32:
		return uint64(uint32(c))
	case OpConst16:
		return uint64(uint16(c))
	case OpConst8:
		return uint64(uint8(c))
	}
	v.Fatalf("op %s isn't OpConst*", v.Op)
	return 0
}

func (v *Value) AuxFloat() float64 {
	if opcodeTable[v.Op].auxType != auxFloat32 && opcodeTable[v.Op].auxType != auxFloat64 {
		v.Fatalf("op %s doesn't have a float aux field", v.Op)
	}
	return math.Float64frombits(uint64(v.AuxInt))
}
func (v *Value) AuxValAndOff() ValAndOff {
	if opcodeTable[v.Op].auxType != auxSymValAndOff {
		v.Fatalf("op %s doesn't have a ValAndOff aux field", v.Op)
	}
	return ValAndOff(v.AuxInt)
}

// long form print.  v# = opcode <type> [aux] args [: reg] (names)
func (v *Value) LongString() string {
	s := fmt.Sprintf("v%d = %s", v.ID, v.Op)
	s += " <" + v.Type.String() + ">"
	s += v.auxString()
	for _, a := range v.Args {
		s += fmt.Sprintf(" %v", a)
	}
	var r []Location
	if v.Block != nil {
		r = v.Block.Func.RegAlloc
	}
	if int(v.ID) < len(r) && r[v.ID] != nil {
		s += " : " + r[v.ID].String()
	}
	var names []string
	if v.Block != nil {
		for name, values := range v.Block.Func.NamedValues {
			for _, value := range values {
				if value == v {
					names = append(names, name.String())
					break // drop duplicates.
				}
			}
		}
	}
	if len(names) != 0 {
		sort.Strings(names) // Otherwise a source of variation in debugging output.
		s += " (" + strings.Join(names, ", ") + ")"
	}
	return s
}

func (v *Value) auxString() string {
	switch opcodeTable[v.Op].auxType {
	case auxBool:
		if v.AuxInt == 0 {
			return " [false]"
		} else {
			return " [true]"
		}
	case auxInt8:
		return fmt.Sprintf(" [%d]", v.AuxInt8())
	case auxInt16:
		return fmt.Sprintf(" [%d]", v.AuxInt16())
	case auxInt32:
		return fmt.Sprintf(" [%d]", v.AuxInt32())
	case auxInt64, auxInt128:
		return fmt.Sprintf(" [%d]", v.AuxInt)
	case auxFloat32, auxFloat64:
		return fmt.Sprintf(" [%g]", v.AuxFloat())
	case auxString:
		return fmt.Sprintf(" {%q}", v.Aux)
	case auxSym, auxTyp, auxArchSpecific:
		if v.Aux != nil {
			return fmt.Sprintf(" {%v}", v.Aux)
		}
	case auxSymOff, auxTypSize:
		s := ""
		if v.Aux != nil {
			s = fmt.Sprintf(" {%v}", v.Aux)
		}
		if v.AuxInt != 0 {
			s += fmt.Sprintf(" [%v]", v.AuxInt)
		}
		return s
	case auxSymValAndOff:
		s := ""
		if v.Aux != nil {
			s = fmt.Sprintf(" {%v}", v.Aux)
		}
		return s + fmt.Sprintf(" [%s]", v.AuxValAndOff())
	case auxCCop:
		return fmt.Sprintf(" {%s}", v.Aux.(Op))
	}
	return ""
}

// If/when midstack inlining is enabled (-l=4), the compiler gets both larger and slower.
// Not-inlining this method is a help (*Value.reset and *Block.NewValue0 are similar).
//go:noinline
func (v *Value) AddArg(w *Value) {
	if v.Args == nil {
		v.resetArgs() // use argstorage
	}
	v.Args = append(v.Args, w)
	w.Uses++
}
func (v *Value) AddArgs(a ...*Value) {
	if v.Args == nil {
		v.resetArgs() // use argstorage
	}
	v.Args = append(v.Args, a...)
	for _, x := range a {
		x.Uses++
	}
}
func (v *Value) SetArg(i int, w *Value) {
	v.Args[i].Uses--
	v.Args[i] = w
	w.Uses++
}
func (v *Value) RemoveArg(i int) {
	v.Args[i].Uses--
	copy(v.Args[i:], v.Args[i+1:])
	v.Args[len(v.Args)-1] = nil // aid GC
	v.Args = v.Args[:len(v.Args)-1]
}
func (v *Value) SetArgs1(a *Value) {
	v.resetArgs()
	v.AddArg(a)
}
func (v *Value) SetArgs2(a *Value, b *Value) {
	v.resetArgs()
	v.AddArg(a)
	v.AddArg(b)
}

func (v *Value) resetArgs() {
	for _, a := range v.Args {
		a.Uses--
	}
	v.argstorage[0] = nil
	v.argstorage[1] = nil
	v.argstorage[2] = nil
	v.Args = v.argstorage[:0]
}

func (v *Value) reset(op Op) {
	v.Op = op
	if op != OpCopy && notStmtBoundary(op) {
		// Special case for OpCopy because of how it is used in rewrite
		v.Pos = v.Pos.WithNotStmt()
	}
	v.resetArgs()
	v.AuxInt = 0
	v.Aux = nil
}

// copyInto makes a new value identical to v and adds it to the end of b.
// unlike copyIntoWithXPos this does not check for v.Pos being a statement.
func (v *Value) copyInto(b *Block) *Value {
	c := b.NewValue0(v.Pos.WithNotStmt(), v.Op, v.Type) // Lose the position, this causes line number churn otherwise.
	c.Aux = v.Aux
	c.AuxInt = v.AuxInt
	c.AddArgs(v.Args...)
	for _, a := range v.Args {
		if a.Type.IsMemory() {
			v.Fatalf("can't move a value with a memory arg %s", v.LongString())
		}
	}
	return c
}

// copyIntoWithXPos makes a new value identical to v and adds it to the end of b.
// The supplied position is used as the position of the new value.
// Because this is used for rematerialization, check for case that (rematerialized)
// input to value with position 'pos' carried a statement mark, and that the supplied
// position (of the instruction using the rematerialized value) is not marked, and
// preserve that mark if its line matches the supplied position.
func (v *Value) copyIntoWithXPos(b *Block, pos src.XPos) *Value {
	if v.Pos.IsStmt() == src.PosIsStmt && pos.IsStmt() != src.PosIsStmt && v.Pos.SameFileAndLine(pos) {
		pos = pos.WithIsStmt()
	}
	c := b.NewValue0(pos, v.Op, v.Type)
	c.Aux = v.Aux
	c.AuxInt = v.AuxInt
	c.AddArgs(v.Args...)
	for _, a := range v.Args {
		if a.Type.IsMemory() {
			v.Fatalf("can't move a value with a memory arg %s", v.LongString())
		}
	}
	return c
}

func (v *Value) Logf(msg string, args ...interface{}) { v.Block.Logf(msg, args...) }
func (v *Value) Log() bool                            { return v.Block.Log() }
func (v *Value) Fatalf(msg string, args ...interface{}) {
	v.Block.Func.fe.Fatalf(v.Pos, msg, args...)
}

// isGenericIntConst reports whether v is a generic integer constant.
func (v *Value) isGenericIntConst() bool {
	return v != nil && (v.Op == OpConst64 || v.Op == OpConst32 || v.Op == OpConst16 || v.Op == OpConst8)
}

// Reg returns the register assigned to v, in cmd/internal/obj/$ARCH numbering.
func (v *Value) Reg() int16 {
	reg := v.Block.Func.RegAlloc[v.ID]
	if reg == nil {
		v.Fatalf("nil register for value: %s\n%s\n", v.LongString(), v.Block.Func)
	}
	return reg.(*Register).objNum
}

// Reg0 returns the register assigned to the first output of v, in cmd/internal/obj/$ARCH numbering.
func (v *Value) Reg0() int16 {
	reg := v.Block.Func.RegAlloc[v.ID].(LocPair)[0]
	if reg == nil {
		v.Fatalf("nil first register for value: %s\n%s\n", v.LongString(), v.Block.Func)
	}
	return reg.(*Register).objNum
}

// Reg1 returns the register assigned to the second output of v, in cmd/internal/obj/$ARCH numbering.
func (v *Value) Reg1() int16 {
	reg := v.Block.Func.RegAlloc[v.ID].(LocPair)[1]
	if reg == nil {
		v.Fatalf("nil second register for value: %s\n%s\n", v.LongString(), v.Block.Func)
	}
	return reg.(*Register).objNum
}

func (v *Value) RegName() string {
	reg := v.Block.Func.RegAlloc[v.ID]
	if reg == nil {
		v.Fatalf("nil register for value: %s\n%s\n", v.LongString(), v.Block.Func)
	}
	return reg.(*Register).name
}

// MemoryArg returns the memory argument for the Value.
// The returned value, if non-nil, will be memory-typed (or a tuple with a memory-typed second part).
// Otherwise, nil is returned.
func (v *Value) MemoryArg() *Value {
	if v.Op == OpPhi {
		v.Fatalf("MemoryArg on Phi")
	}
	na := len(v.Args)
	if na == 0 {
		return nil
	}
	if m := v.Args[na-1]; m.Type.IsMemory() {
		return m
	}
	return nil
}

// LackingPos indicates whether v is a value that is unlikely to have a correct
// position assigned to it.  Ignoring such values leads to more user-friendly positions
// assigned to nearby values and the blocks containing them.
func (v *Value) LackingPos() bool {
	// The exact definition of LackingPos is somewhat heuristically defined and may change
	// in the future, for example if some of these operations are generated more carefully
	// with respect to their source position.
	return v.Op == OpVarDef || v.Op == OpVarKill || v.Op == OpVarLive || v.Op == OpPhi ||
		(v.Op == OpFwdRef || v.Op == OpCopy) && v.Type == types.TypeMem
}
