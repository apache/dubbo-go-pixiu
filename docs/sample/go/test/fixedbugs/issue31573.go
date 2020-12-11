// errorcheck -0 -m

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

func f(...*int) {} // ERROR "can inline f$"

func g() {
	defer f()                   // ERROR "... argument does not escape$"
	defer f(new(int))           // ERROR "... argument does not escape$" "new\(int\) does not escape$"
	defer f(new(int), new(int)) // ERROR "... argument does not escape$" "new\(int\) does not escape$"

	defer f(nil...)
	defer f([]*int{}...)                   // ERROR "\[\]\*int literal does not escape$"
	defer f([]*int{new(int)}...)           // ERROR "\[\]\*int literal does not escape$" "new\(int\) does not escape$"
	defer f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int literal does not escape$" "new\(int\) does not escape$"

	go f()                   // ERROR "... argument escapes to heap$"
	go f(new(int))           // ERROR "... argument escapes to heap$" "new\(int\) escapes to heap$"
	go f(new(int), new(int)) // ERROR "... argument escapes to heap$" "new\(int\) escapes to heap$"

	go f(nil...)
	go f([]*int{}...)                   // ERROR "\[\]\*int literal escapes to heap$"
	go f([]*int{new(int)}...)           // ERROR "\[\]\*int literal escapes to heap$" "new\(int\) escapes to heap$"
	go f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int literal escapes to heap$" "new\(int\) escapes to heap$"

	for {
		defer f()                   // ERROR "... argument escapes to heap$"
		defer f(new(int))           // ERROR "... argument escapes to heap$" "new\(int\) escapes to heap$"
		defer f(new(int), new(int)) // ERROR "... argument escapes to heap$" "new\(int\) escapes to heap$"

		defer f(nil...)
		defer f([]*int{}...)                   // ERROR "\[\]\*int literal escapes to heap$"
		defer f([]*int{new(int)}...)           // ERROR "\[\]\*int literal escapes to heap$" "new\(int\) escapes to heap$"
		defer f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int literal escapes to heap$" "new\(int\) escapes to heap$"

		go f()                   // ERROR "... argument escapes to heap$"
		go f(new(int))           // ERROR "... argument escapes to heap$" "new\(int\) escapes to heap$"
		go f(new(int), new(int)) // ERROR "... argument escapes to heap$" "new\(int\) escapes to heap$"

		go f(nil...)
		go f([]*int{}...)                   // ERROR "\[\]\*int literal escapes to heap$"
		go f([]*int{new(int)}...)           // ERROR "\[\]\*int literal escapes to heap$" "new\(int\) escapes to heap$"
		go f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int literal escapes to heap$" "new\(int\) escapes to heap$"
	}
}
