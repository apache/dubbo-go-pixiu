// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import "sync"

// The line numbering of this file is important for TestTBHelper.

func notHelper(t *T, msg string) {
	t.Error(msg)
}

func helper(t *T, msg string) {
	t.Helper()
	t.Error(msg)
}

func notHelperCallingHelper(t *T, msg string) {
	helper(t, msg)
}

func helperCallingHelper(t *T, msg string) {
	t.Helper()
	helper(t, msg)
}

func testHelper(t *T) {
	// Check combinations of directly and indirectly
	// calling helper functions.
	notHelper(t, "0")
	helper(t, "1")
	notHelperCallingHelper(t, "2")
	helperCallingHelper(t, "3")

	// Check a function literal closing over t that uses Helper.
	fn := func(msg string) {
		t.Helper()
		t.Error(msg)
	}
	fn("4")

	t.Run("sub", func(t *T) {
		helper(t, "5")
		notHelperCallingHelper(t, "6")
		// Check that calling Helper from inside a subtest entry function
		// works as if it were in an ordinary function call.
		t.Helper()
		t.Error("7")
	})

	// Check that calling Helper from inside a top-level test function
	// has no effect.
	t.Helper()
	t.Error("8")
}

func parallelTestHelper(t *T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			notHelperCallingHelper(t, "parallel")
			wg.Done()
		}()
	}
	wg.Wait()
}
