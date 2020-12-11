// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import "testing"

// TestNlzNto tests nlz/nto of the same number which is used in some of
// the rewrite rules.
func TestNlzNto(t *testing.T) {
	// construct the bit pattern 000...111, nlz(x) + nto(0) = 64
	var x int64
	for i := int64(0); i < 64; i++ {
		if got := nto(x); got != i {
			t.Errorf("expected nto(0x%X) = %d, got %d", x, i, got)
		}
		if got := nlz(x); got != 64-i {
			t.Errorf("expected nlz(0x%X) = %d, got %d", x, 64-i, got)
		}
		x = (x << 1) | 1
	}

	x = 0
	// construct the bit pattern 000...111, with bit 33 set as well.
	for i := int64(0); i < 64; i++ {
		tx := x | (1 << 32)
		// nto should be the number of bits we've shifted on, with an extra bit
		// at iter 32
		ntoExp := i
		if ntoExp == 32 {
			ntoExp = 33
		}
		if got := nto(tx); got != ntoExp {
			t.Errorf("expected nto(0x%X) = %d, got %d", tx, ntoExp, got)
		}

		// sinec bit 33 is set, nlz can be no greater than 31
		nlzExp := 64 - i
		if nlzExp > 31 {
			nlzExp = 31
		}
		if got := nlz(tx); got != nlzExp {
			t.Errorf("expected nlz(0x%X) = %d, got %d", tx, nlzExp, got)
		}
		x = (x << 1) | 1
	}

}

func TestNlz(t *testing.T) {
	var nlzTests = []struct {
		v   int64
		exp int64
	}{{0x00, 64},
		{0x01, 63},
		{0x0F, 60},
		{0xFF, 56},
		{0xffffFFFF, 32},
		{-0x01, 0}}

	for _, tc := range nlzTests {
		if got := nlz(tc.v); got != tc.exp {
			t.Errorf("expected nlz(0x%X) = %d, got %d", tc.v, tc.exp, got)
		}
	}
}

func TestNto(t *testing.T) {
	var ntoTests = []struct {
		v   int64
		exp int64
	}{{0x00, 0},
		{0x01, 1},
		{0x0F, 4},
		{0xFF, 8},
		{0xffffFFFF, 32},
		{-0x01, 64}}

	for _, tc := range ntoTests {
		if got := nto(tc.v); got != tc.exp {
			t.Errorf("expected nto(0x%X) = %d, got %d", tc.v, tc.exp, got)
		}
	}
}

func TestLog2(t *testing.T) {
	var log2Tests = []struct {
		v   int64
		exp int64
	}{{0, -1}, // nlz expects log2(0) == -1
		{1, 0},
		{2, 1},
		{4, 2},
		{7, 2},
		{8, 3},
		{9, 3},
		{1024, 10}}

	for _, tc := range log2Tests {
		if got := log2(tc.v); got != tc.exp {
			t.Errorf("expected log2(%d) = %d, got %d", tc.v, tc.exp, got)
		}
	}
}

// We generate memmove for copy(x[1:], x[:]), however we may change it to OpMove,
// because size is known. Check that OpMove is alias-safe, or we did call memmove.
func TestMove(t *testing.T) {
	x := [...]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}
	copy(x[1:], x[:])
	for i := 1; i < len(x); i++ {
		if int(x[i]) != i {
			t.Errorf("Memmove got converted to OpMove in alias-unsafe way. Got %d insted of %d in position %d", int(x[i]), i, i+1)
		}
	}
}

func TestMoveSmall(t *testing.T) {
	x := [...]byte{1, 2, 3, 4, 5, 6, 7}
	copy(x[1:], x[:])
	for i := 1; i < len(x); i++ {
		if int(x[i]) != i {
			t.Errorf("Memmove got converted to OpMove in alias-unsafe way. Got %d instead of %d in position %d", int(x[i]), i, i+1)
		}
	}
}
