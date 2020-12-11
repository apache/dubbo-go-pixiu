// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"internal/testenv"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func loadSyms(t *testing.T) map[string]string {
	cmd := exec.Command(testenv.GoToolPath(t), "tool", "nm", os.Args[0])
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go tool nm %v: %v\n%s", os.Args[0], err, string(out))
	}
	syms := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		if len(f) < 3 {
			continue
		}
		syms[f[2]] = f[0]
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading symbols: %v", err)
	}
	return syms
}

func runAddr2Line(t *testing.T, exepath, addr string) (funcname, path, lineno string) {
	cmd := exec.Command(exepath, os.Args[0])
	cmd.Stdin = strings.NewReader(addr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go tool addr2line %v: %v\n%s", os.Args[0], err, string(out))
	}
	f := strings.Split(string(out), "\n")
	if len(f) < 3 && f[2] == "" {
		t.Fatal("addr2line output must have 2 lines")
	}
	funcname = f[0]
	pathAndLineNo := f[1]
	f = strings.Split(pathAndLineNo, ":")
	if runtime.GOOS == "windows" && len(f) == 3 {
		// Reattach drive letter.
		f = []string{f[0] + ":" + f[1], f[2]}
	}
	if len(f) != 2 {
		t.Fatalf("no line number found in %q", pathAndLineNo)
	}
	return funcname, f[0], f[1]
}

const symName = "cmd/addr2line.TestAddr2Line"

func testAddr2Line(t *testing.T, exepath, addr string) {
	funcName, srcPath, srcLineNo := runAddr2Line(t, exepath, addr)
	if symName != funcName {
		t.Fatalf("expected function name %v; got %v", symName, funcName)
	}
	fi1, err := os.Stat("addr2line_test.go")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	fi2, err := os.Stat(srcPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if !os.SameFile(fi1, fi2) {
		t.Fatalf("addr2line_test.go and %s are not same file", srcPath)
	}
	if srcLineNo != "89" {
		t.Fatalf("line number = %v; want 89", srcLineNo)
	}
}

// This is line 88. The test depends on that.
func TestAddr2Line(t *testing.T) {
	testenv.MustHaveGoBuild(t)

	tmpDir, err := ioutil.TempDir("", "TestAddr2Line")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// Build copy of test binary with debug symbols,
	// since the one running now may not have them.
	exepath := filepath.Join(tmpDir, "testaddr2line_test.exe")
	out, err := exec.Command(testenv.GoToolPath(t), "test", "-c", "-o", exepath, "cmd/addr2line").CombinedOutput()
	if err != nil {
		t.Fatalf("go test -c -o %v cmd/addr2line: %v\n%s", exepath, err, string(out))
	}
	os.Args[0] = exepath

	syms := loadSyms(t)

	exepath = filepath.Join(tmpDir, "testaddr2line.exe")
	out, err = exec.Command(testenv.GoToolPath(t), "build", "-o", exepath, "cmd/addr2line").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmd/addr2line: %v\n%s", exepath, err, string(out))
	}

	testAddr2Line(t, exepath, syms[symName])
	testAddr2Line(t, exepath, "0x"+syms[symName])
}
