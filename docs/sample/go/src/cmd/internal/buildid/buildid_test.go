// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildid

import (
	"bytes"
	"crypto/sha256"
	"internal/obscuretestdata"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const (
	expectedID = "abcdefghijklmnopqrstuvwxyz.1234567890123456789012345678901234567890123456789012345678901234"
	newID      = "bcdefghijklmnopqrstuvwxyza.2345678901234567890123456789012345678901234567890123456789012341"
)

func TestReadFile(t *testing.T) {
	f, err := ioutil.TempFile("", "buildid-test-")
	if err != nil {
		t.Fatal(err)
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	f.Close()

	// Use obscured files to prevent Apple’s notarization service from
	// mistaking them as candidates for notarization and rejecting the entire
	// toolchain.
	// See golang.org/issue/34986
	var files = []string{
		"p.a.base64",
		"a.elf.base64",
		"a.macho.base64",
		"a.pe.base64",
	}

	for _, name := range files {
		f, err := obscuretestdata.DecodeToTempFile("testdata/" + name)
		if err != nil {
			t.Errorf("obscuretestdata.DecodeToTempFile(testdata/%s): %v", name, err)
			continue
		}
		defer os.Remove(f)
		id, err := ReadFile(f)
		if id != expectedID || err != nil {
			t.Errorf("ReadFile(testdata/%s) = %q, %v, want %q, nil", f, id, err, expectedID)
		}
		old := readSize
		readSize = 2048
		id, err = ReadFile(f)
		readSize = old
		if id != expectedID || err != nil {
			t.Errorf("ReadFile(%s) [readSize=2k] = %q, %v, want %q, nil", f, id, err, expectedID)
		}

		data, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		m, _, err := FindAndHash(bytes.NewReader(data), expectedID, 1024)
		if err != nil {
			t.Errorf("FindAndHash(%s): %v", f, err)
			continue
		}
		if err := ioutil.WriteFile(tmp, data, 0666); err != nil {
			t.Error(err)
			continue
		}
		tf, err := os.OpenFile(tmp, os.O_WRONLY, 0)
		if err != nil {
			t.Error(err)
			continue
		}
		err = Rewrite(tf, m, newID)
		err2 := tf.Close()
		if err != nil {
			t.Errorf("Rewrite(%s): %v", f, err)
			continue
		}
		if err2 != nil {
			t.Fatal(err2)
		}

		id, err = ReadFile(tmp)
		if id != newID || err != nil {
			t.Errorf("ReadFile(%s after Rewrite) = %q, %v, want %q, nil", f, id, err, newID)
		}
	}
}

func TestFindAndHash(t *testing.T) {
	buf := make([]byte, 64)
	buf2 := make([]byte, 64)
	id := make([]byte, 8)
	zero := make([]byte, 8)
	for i := range id {
		id[i] = byte(i)
	}
	numError := 0
	errorf := func(msg string, args ...interface{}) {
		t.Errorf(msg, args...)
		if numError++; numError > 20 {
			t.Logf("stopping after too many errors")
			t.FailNow()
		}
	}
	for bufSize := len(id); bufSize <= len(buf); bufSize++ {
		for j := range buf {
			for k := 0; k < 2*len(id) && j+k < len(buf); k++ {
				for i := range buf {
					buf[i] = 1
				}
				copy(buf[j:], id)
				copy(buf[j+k:], id)
				var m []int64
				if j+len(id) <= j+k {
					m = append(m, int64(j))
				}
				if j+k+len(id) <= len(buf) {
					m = append(m, int64(j+k))
				}
				copy(buf2, buf)
				for _, p := range m {
					copy(buf2[p:], zero)
				}
				h := sha256.Sum256(buf2)

				matches, hash, err := FindAndHash(bytes.NewReader(buf), string(id), bufSize)
				if err != nil {
					errorf("bufSize=%d j=%d k=%d: findAndHash: %v", bufSize, j, k, err)
					continue
				}
				if !reflect.DeepEqual(matches, m) {
					errorf("bufSize=%d j=%d k=%d: findAndHash: matches=%v, want %v", bufSize, j, k, matches, m)
					continue
				}
				if hash != h {
					errorf("bufSize=%d j=%d k=%d: findAndHash: matches correct, but hash=%x, want %x", bufSize, j, k, hash, h)
				}
			}
		}
	}
}
