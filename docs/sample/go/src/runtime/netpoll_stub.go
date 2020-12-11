// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package runtime

import "runtime/internal/atomic"

var netpollInited uint32
var netpollWaiters uint32

var netpollStubLock mutex
var netpollNote note

// netpollBroken, protected by netpollBrokenLock, avoids a double notewakeup.
var netpollBrokenLock mutex
var netpollBroken bool

func netpollGenericInit() {
	atomic.Store(&netpollInited, 1)
}

func netpollBreak() {
	lock(&netpollBrokenLock)
	broken := netpollBroken
	netpollBroken = true
	if !broken {
		notewakeup(&netpollNote)
	}
	unlock(&netpollBrokenLock)
}

// Polls for ready network connections.
// Returns list of goroutines that become runnable.
func netpoll(delay int64) gList {
	// Implementation for platforms that do not support
	// integrated network poller.
	if delay != 0 {
		// This lock ensures that only one goroutine tries to use
		// the note. It should normally be completely uncontended.
		lock(&netpollStubLock)

		lock(&netpollBrokenLock)
		noteclear(&netpollNote)
		netpollBroken = false
		unlock(&netpollBrokenLock)

		notetsleep(&netpollNote, delay)
		unlock(&netpollStubLock)
	}
	return gList{}
}

func netpollinited() bool {
	return atomic.Load(&netpollInited) != 0
}
