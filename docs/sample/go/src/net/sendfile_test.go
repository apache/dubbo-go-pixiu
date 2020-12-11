// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !js

package net

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	newton       = "../testdata/Isaac.Newton-Opticks.txt"
	newtonLen    = 567198
	newtonSHA256 = "d4a9ac22462b35e7821a4f2706c211093da678620a8f9997989ee7cf8d507bbd"
)

func TestSendfile(t *testing.T) {
	ln, err := newLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	errc := make(chan error, 1)
	go func(ln Listener) {
		// Wait for a connection.
		conn, err := ln.Accept()
		if err != nil {
			errc <- err
			close(errc)
			return
		}

		go func() {
			defer close(errc)
			defer conn.Close()

			f, err := os.Open(newton)
			if err != nil {
				errc <- err
				return
			}
			defer f.Close()

			// Return file data using io.Copy, which should use
			// sendFile if available.
			sbytes, err := io.Copy(conn, f)
			if err != nil {
				errc <- err
				return
			}

			if sbytes != newtonLen {
				errc <- fmt.Errorf("sent %d bytes; expected %d", sbytes, newtonLen)
				return
			}
		}()
	}(ln)

	// Connect to listener to retrieve file and verify digest matches
	// expected.
	c, err := Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	h := sha256.New()
	rbytes, err := io.Copy(h, c)
	if err != nil {
		t.Error(err)
	}

	if rbytes != newtonLen {
		t.Errorf("received %d bytes; expected %d", rbytes, newtonLen)
	}

	if res := hex.EncodeToString(h.Sum(nil)); res != newtonSHA256 {
		t.Error("retrieved data hash did not match")
	}

	for err := range errc {
		t.Error(err)
	}
}

func TestSendfileParts(t *testing.T) {
	ln, err := newLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	errc := make(chan error, 1)
	go func(ln Listener) {
		// Wait for a connection.
		conn, err := ln.Accept()
		if err != nil {
			errc <- err
			close(errc)
			return
		}

		go func() {
			defer close(errc)
			defer conn.Close()

			f, err := os.Open(newton)
			if err != nil {
				errc <- err
				return
			}
			defer f.Close()

			for i := 0; i < 3; i++ {
				// Return file data using io.CopyN, which should use
				// sendFile if available.
				_, err = io.CopyN(conn, f, 3)
				if err != nil {
					errc <- err
					return
				}
			}
		}()
	}(ln)

	c, err := Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(c)

	if want, have := "Produced ", buf.String(); have != want {
		t.Errorf("unexpected server reply %q, want %q", have, want)
	}

	for err := range errc {
		t.Error(err)
	}
}

func TestSendfileSeeked(t *testing.T) {
	ln, err := newLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	const seekTo = 65 << 10
	const sendSize = 10 << 10

	errc := make(chan error, 1)
	go func(ln Listener) {
		// Wait for a connection.
		conn, err := ln.Accept()
		if err != nil {
			errc <- err
			close(errc)
			return
		}

		go func() {
			defer close(errc)
			defer conn.Close()

			f, err := os.Open(newton)
			if err != nil {
				errc <- err
				return
			}
			defer f.Close()
			if _, err := f.Seek(seekTo, os.SEEK_SET); err != nil {
				errc <- err
				return
			}

			_, err = io.CopyN(conn, f, sendSize)
			if err != nil {
				errc <- err
				return
			}
		}()
	}(ln)

	c, err := Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(c)

	if buf.Len() != sendSize {
		t.Errorf("Got %d bytes; want %d", buf.Len(), sendSize)
	}

	for err := range errc {
		t.Error(err)
	}
}

// Test that sendfile doesn't put a pipe into blocking mode.
func TestSendfilePipe(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		// These systems don't support deadlines on pipes.
		t.Skipf("skipping on %s", runtime.GOOS)
	}

	t.Parallel()

	ln, err := newLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	defer r.Close()

	copied := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Accept a connection and copy 1 byte from the read end of
		// the pipe to the connection. This will call into sendfile.
		defer wg.Done()
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		_, err = io.CopyN(conn, r, 1)
		if err != nil {
			t.Error(err)
			return
		}
		// Signal the main goroutine that we've copied the byte.
		close(copied)
	}()

	wg.Add(1)
	go func() {
		// Write 1 byte to the write end of the pipe.
		defer wg.Done()
		_, err := w.Write([]byte{'a'})
		if err != nil {
			t.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		// Connect to the server started two goroutines up and
		// discard any data that it writes.
		defer wg.Done()
		conn, err := Dial("tcp", ln.Addr().String())
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		io.Copy(ioutil.Discard, conn)
	}()

	// Wait for the byte to be copied, meaning that sendfile has
	// been called on the pipe.
	<-copied

	// Set a very short deadline on the read end of the pipe.
	if err := r.SetDeadline(time.Now().Add(time.Microsecond)); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		// Wait for much longer than the deadline and write a byte
		// to the pipe.
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		w.Write([]byte{'b'})
	}()

	// If this read does not time out, the pipe was incorrectly
	// put into blocking mode.
	_, err = r.Read(make([]byte, 1))
	if err == nil {
		t.Error("Read did not time out")
	} else if !os.IsTimeout(err) {
		t.Errorf("got error %v, expected a time out", err)
	}

	wg.Wait()
}
