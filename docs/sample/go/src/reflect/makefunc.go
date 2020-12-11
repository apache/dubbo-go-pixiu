// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// MakeFunc implementation.

package reflect

import (
	"unsafe"
)

// makeFuncImpl is the closure value implementing the function
// returned by MakeFunc.
// The first three words of this type must be kept in sync with
// methodValue and runtime.reflectMethodValue.
// Any changes should be reflected in all three.
type makeFuncImpl struct {
	code   uintptr
	stack  *bitVector // ptrmap for both args and results
	argLen uintptr    // just args
	ftyp   *funcType
	fn     func([]Value) []Value
}

// MakeFunc returns a new function of the given Type
// that wraps the function fn. When called, that new function
// does the following:
//
//	- converts its arguments to a slice of Values.
//	- runs results := fn(args).
//	- returns the results as a slice of Values, one per formal result.
//
// The implementation fn can assume that the argument Value slice
// has the number and type of arguments given by typ.
// If typ describes a variadic function, the final Value is itself
// a slice representing the variadic arguments, as in the
// body of a variadic function. The result Value slice returned by fn
// must have the number and type of results given by typ.
//
// The Value.Call method allows the caller to invoke a typed function
// in terms of Values; in contrast, MakeFunc allows the caller to implement
// a typed function in terms of Values.
//
// The Examples section of the documentation includes an illustration
// of how to use MakeFunc to build a swap function for different types.
//
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {
	if typ.Kind() != Func {
		panic("reflect: call of MakeFunc with non-Func type")
	}

	t := typ.common()
	ftyp := (*funcType)(unsafe.Pointer(t))

	// Indirect Go func value (dummy) to obtain
	// actual code address. (A Go func value is a pointer
	// to a C function pointer. https://golang.org/s/go11func.)
	dummy := makeFuncStub
	code := **(**uintptr)(unsafe.Pointer(&dummy))

	// makeFuncImpl contains a stack map for use by the runtime
	_, argLen, _, stack, _ := funcLayout(ftyp, nil)

	impl := &makeFuncImpl{code: code, stack: stack, argLen: argLen, ftyp: ftyp, fn: fn}

	return Value{t, unsafe.Pointer(impl), flag(Func)}
}

// makeFuncStub is an assembly function that is the code half of
// the function returned from MakeFunc. It expects a *callReflectFunc
// as its context register, and its job is to invoke callReflect(ctxt, frame)
// where ctxt is the context register and frame is a pointer to the first
// word in the passed-in argument frame.
func makeFuncStub()

// The first 3 words of this type must be kept in sync with
// makeFuncImpl and runtime.reflectMethodValue.
// Any changes should be reflected in all three.
type methodValue struct {
	fn     uintptr
	stack  *bitVector // ptrmap for both args and results
	argLen uintptr    // just args
	method int
	rcvr   Value
}

// makeMethodValue converts v from the rcvr+method index representation
// of a method value to an actual method func value, which is
// basically the receiver value with a special bit set, into a true
// func value - a value holding an actual func. The output is
// semantically equivalent to the input as far as the user of package
// reflect can tell, but the true func representation can be handled
// by code like Convert and Interface and Assign.
func makeMethodValue(op string, v Value) Value {
	if v.flag&flagMethod == 0 {
		panic("reflect: internal error: invalid use of makeMethodValue")
	}

	// Ignoring the flagMethod bit, v describes the receiver, not the method type.
	fl := v.flag & (flagRO | flagAddr | flagIndir)
	fl |= flag(v.typ.Kind())
	rcvr := Value{v.typ, v.ptr, fl}

	// v.Type returns the actual type of the method value.
	ftyp := (*funcType)(unsafe.Pointer(v.Type().(*rtype)))

	// Indirect Go func value (dummy) to obtain
	// actual code address. (A Go func value is a pointer
	// to a C function pointer. https://golang.org/s/go11func.)
	dummy := methodValueCall
	code := **(**uintptr)(unsafe.Pointer(&dummy))

	// methodValue contains a stack map for use by the runtime
	_, argLen, _, stack, _ := funcLayout(ftyp, nil)

	fv := &methodValue{
		fn:     code,
		stack:  stack,
		argLen: argLen,
		method: int(v.flag) >> flagMethodShift,
		rcvr:   rcvr,
	}

	// Cause panic if method is not appropriate.
	// The panic would still happen during the call if we omit this,
	// but we want Interface() and other operations to fail early.
	methodReceiver(op, fv.rcvr, fv.method)

	return Value{&ftyp.rtype, unsafe.Pointer(fv), v.flag&flagRO | flag(Func)}
}

// methodValueCall is an assembly function that is the code half of
// the function returned from makeMethodValue. It expects a *methodValue
// as its context register, and its job is to invoke callMethod(ctxt, frame)
// where ctxt is the context register and frame is a pointer to the first
// word in the passed-in argument frame.
func methodValueCall()
