// +build linux,!ppc64,!riscv64
// run

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test that a -B option is passed through when using both internal
// and external linking mode.

// TODO(jsing): Re-enable on riscv64 when it has support for external
// linking - see golang.org/issue/36739

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	test("internal")
	test("external")
}

func test(linkmode string) {
	out, err := exec.Command("go", "run", "-ldflags", "-B=0x12345678 -linkmode="+linkmode, filepath.Join("fixedbugs", "issue10607a.go")).CombinedOutput()
	if err != nil {
		fmt.Printf("BUG: linkmode=%s %v\n%s\n", linkmode, err, out)
		os.Exit(1)
	}
}
