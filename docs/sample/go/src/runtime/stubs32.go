// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386 arm mips mipsle

package runtime

import "unsafe"

// Declarations for runtime services implemented in C or assembly that
// are only present on 32 bit systems.

func call16(typ, fn, arg unsafe.Pointer, n, retoffset uint32)
