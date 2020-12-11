// errorcheck

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

type T struct{}
type A = T
type B = T

func (T) m() {}
func (T) m() {} // ERROR "redeclared"
func (A) m() {} // ERROR "redeclared"
func (A) m() {} // ERROR "redeclared"
func (B) m() {} // ERROR "redeclared"
func (B) m() {} // ERROR "redeclared"

func (*T) m() {} // ERROR "redeclared"
func (*A) m() {} // ERROR "redeclared"
func (*B) m() {} // ERROR "redeclared"
