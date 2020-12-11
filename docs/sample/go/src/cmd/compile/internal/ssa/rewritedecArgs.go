// Code generated from gen/decArgs.rules; DO NOT EDIT.
// generated with: cd gen; go run *.go

package ssa

func rewriteValuedecArgs(v *Value) bool {
	switch v.Op {
	case OpArg:
		return rewriteValuedecArgs_OpArg_0(v) || rewriteValuedecArgs_OpArg_10(v)
	}
	return false
}
func rewriteValuedecArgs_OpArg_0(v *Value) bool {
	b := v.Block
	config := b.Func.Config
	fe := b.Func.fe
	typ := &b.Func.Config.Types
	// match: (Arg {n} [off])
	// cond: v.Type.IsString()
	// result: (StringMake (Arg <typ.BytePtr> {n} [off]) (Arg <typ.Int> {n} [off+config.PtrSize]))
	for {
		off := v.AuxInt
		n := v.Aux
		if !(v.Type.IsString()) {
			break
		}
		v.reset(OpStringMake)
		v0 := b.NewValue0(v.Pos, OpArg, typ.BytePtr)
		v0.AuxInt = off
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, typ.Int)
		v1.AuxInt = off + config.PtrSize
		v1.Aux = n
		v.AddArg(v1)
		return true
	}
	// match: (Arg {n} [off])
	// cond: v.Type.IsSlice()
	// result: (SliceMake (Arg <v.Type.Elem().PtrTo()> {n} [off]) (Arg <typ.Int> {n} [off+config.PtrSize]) (Arg <typ.Int> {n} [off+2*config.PtrSize]))
	for {
		off := v.AuxInt
		n := v.Aux
		if !(v.Type.IsSlice()) {
			break
		}
		v.reset(OpSliceMake)
		v0 := b.NewValue0(v.Pos, OpArg, v.Type.Elem().PtrTo())
		v0.AuxInt = off
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, typ.Int)
		v1.AuxInt = off + config.PtrSize
		v1.Aux = n
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpArg, typ.Int)
		v2.AuxInt = off + 2*config.PtrSize
		v2.Aux = n
		v.AddArg(v2)
		return true
	}
	// match: (Arg {n} [off])
	// cond: v.Type.IsInterface()
	// result: (IMake (Arg <typ.Uintptr> {n} [off]) (Arg <typ.BytePtr> {n} [off+config.PtrSize]))
	for {
		off := v.AuxInt
		n := v.Aux
		if !(v.Type.IsInterface()) {
			break
		}
		v.reset(OpIMake)
		v0 := b.NewValue0(v.Pos, OpArg, typ.Uintptr)
		v0.AuxInt = off
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, typ.BytePtr)
		v1.AuxInt = off + config.PtrSize
		v1.Aux = n
		v.AddArg(v1)
		return true
	}
	// match: (Arg {n} [off])
	// cond: v.Type.IsComplex() && v.Type.Size() == 16
	// result: (ComplexMake (Arg <typ.Float64> {n} [off]) (Arg <typ.Float64> {n} [off+8]))
	for {
		off := v.AuxInt
		n := v.Aux
		if !(v.Type.IsComplex() && v.Type.Size() == 16) {
			break
		}
		v.reset(OpComplexMake)
		v0 := b.NewValue0(v.Pos, OpArg, typ.Float64)
		v0.AuxInt = off
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, typ.Float64)
		v1.AuxInt = off + 8
		v1.Aux = n
		v.AddArg(v1)
		return true
	}
	// match: (Arg {n} [off])
	// cond: v.Type.IsComplex() && v.Type.Size() == 8
	// result: (ComplexMake (Arg <typ.Float32> {n} [off]) (Arg <typ.Float32> {n} [off+4]))
	for {
		off := v.AuxInt
		n := v.Aux
		if !(v.Type.IsComplex() && v.Type.Size() == 8) {
			break
		}
		v.reset(OpComplexMake)
		v0 := b.NewValue0(v.Pos, OpArg, typ.Float32)
		v0.AuxInt = off
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, typ.Float32)
		v1.AuxInt = off + 4
		v1.Aux = n
		v.AddArg(v1)
		return true
	}
	// match: (Arg <t>)
	// cond: t.IsStruct() && t.NumFields() == 0 && fe.CanSSA(t)
	// result: (StructMake0)
	for {
		t := v.Type
		if !(t.IsStruct() && t.NumFields() == 0 && fe.CanSSA(t)) {
			break
		}
		v.reset(OpStructMake0)
		return true
	}
	// match: (Arg <t> {n} [off])
	// cond: t.IsStruct() && t.NumFields() == 1 && fe.CanSSA(t)
	// result: (StructMake1 (Arg <t.FieldType(0)> {n} [off+t.FieldOff(0)]))
	for {
		t := v.Type
		off := v.AuxInt
		n := v.Aux
		if !(t.IsStruct() && t.NumFields() == 1 && fe.CanSSA(t)) {
			break
		}
		v.reset(OpStructMake1)
		v0 := b.NewValue0(v.Pos, OpArg, t.FieldType(0))
		v0.AuxInt = off + t.FieldOff(0)
		v0.Aux = n
		v.AddArg(v0)
		return true
	}
	// match: (Arg <t> {n} [off])
	// cond: t.IsStruct() && t.NumFields() == 2 && fe.CanSSA(t)
	// result: (StructMake2 (Arg <t.FieldType(0)> {n} [off+t.FieldOff(0)]) (Arg <t.FieldType(1)> {n} [off+t.FieldOff(1)]))
	for {
		t := v.Type
		off := v.AuxInt
		n := v.Aux
		if !(t.IsStruct() && t.NumFields() == 2 && fe.CanSSA(t)) {
			break
		}
		v.reset(OpStructMake2)
		v0 := b.NewValue0(v.Pos, OpArg, t.FieldType(0))
		v0.AuxInt = off + t.FieldOff(0)
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, t.FieldType(1))
		v1.AuxInt = off + t.FieldOff(1)
		v1.Aux = n
		v.AddArg(v1)
		return true
	}
	// match: (Arg <t> {n} [off])
	// cond: t.IsStruct() && t.NumFields() == 3 && fe.CanSSA(t)
	// result: (StructMake3 (Arg <t.FieldType(0)> {n} [off+t.FieldOff(0)]) (Arg <t.FieldType(1)> {n} [off+t.FieldOff(1)]) (Arg <t.FieldType(2)> {n} [off+t.FieldOff(2)]))
	for {
		t := v.Type
		off := v.AuxInt
		n := v.Aux
		if !(t.IsStruct() && t.NumFields() == 3 && fe.CanSSA(t)) {
			break
		}
		v.reset(OpStructMake3)
		v0 := b.NewValue0(v.Pos, OpArg, t.FieldType(0))
		v0.AuxInt = off + t.FieldOff(0)
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, t.FieldType(1))
		v1.AuxInt = off + t.FieldOff(1)
		v1.Aux = n
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpArg, t.FieldType(2))
		v2.AuxInt = off + t.FieldOff(2)
		v2.Aux = n
		v.AddArg(v2)
		return true
	}
	// match: (Arg <t> {n} [off])
	// cond: t.IsStruct() && t.NumFields() == 4 && fe.CanSSA(t)
	// result: (StructMake4 (Arg <t.FieldType(0)> {n} [off+t.FieldOff(0)]) (Arg <t.FieldType(1)> {n} [off+t.FieldOff(1)]) (Arg <t.FieldType(2)> {n} [off+t.FieldOff(2)]) (Arg <t.FieldType(3)> {n} [off+t.FieldOff(3)]))
	for {
		t := v.Type
		off := v.AuxInt
		n := v.Aux
		if !(t.IsStruct() && t.NumFields() == 4 && fe.CanSSA(t)) {
			break
		}
		v.reset(OpStructMake4)
		v0 := b.NewValue0(v.Pos, OpArg, t.FieldType(0))
		v0.AuxInt = off + t.FieldOff(0)
		v0.Aux = n
		v.AddArg(v0)
		v1 := b.NewValue0(v.Pos, OpArg, t.FieldType(1))
		v1.AuxInt = off + t.FieldOff(1)
		v1.Aux = n
		v.AddArg(v1)
		v2 := b.NewValue0(v.Pos, OpArg, t.FieldType(2))
		v2.AuxInt = off + t.FieldOff(2)
		v2.Aux = n
		v.AddArg(v2)
		v3 := b.NewValue0(v.Pos, OpArg, t.FieldType(3))
		v3.AuxInt = off + t.FieldOff(3)
		v3.Aux = n
		v.AddArg(v3)
		return true
	}
	return false
}
func rewriteValuedecArgs_OpArg_10(v *Value) bool {
	b := v.Block
	fe := b.Func.fe
	// match: (Arg <t>)
	// cond: t.IsArray() && t.NumElem() == 0
	// result: (ArrayMake0)
	for {
		t := v.Type
		if !(t.IsArray() && t.NumElem() == 0) {
			break
		}
		v.reset(OpArrayMake0)
		return true
	}
	// match: (Arg <t> {n} [off])
	// cond: t.IsArray() && t.NumElem() == 1 && fe.CanSSA(t)
	// result: (ArrayMake1 (Arg <t.Elem()> {n} [off]))
	for {
		t := v.Type
		off := v.AuxInt
		n := v.Aux
		if !(t.IsArray() && t.NumElem() == 1 && fe.CanSSA(t)) {
			break
		}
		v.reset(OpArrayMake1)
		v0 := b.NewValue0(v.Pos, OpArg, t.Elem())
		v0.AuxInt = off
		v0.Aux = n
		v.AddArg(v0)
		return true
	}
	return false
}
func rewriteBlockdecArgs(b *Block) bool {
	switch b.Kind {
	}
	return false
}
