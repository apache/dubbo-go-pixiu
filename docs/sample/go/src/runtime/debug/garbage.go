// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"runtime"
	"sort"
	"time"
)

// GCStats collect information about recent garbage collections.
type GCStats struct {
	LastGC         time.Time       // time of last collection
	NumGC          int64           // number of garbage collections
	PauseTotal     time.Duration   // total pause for all collections
	Pause          []time.Duration // pause history, most recent first
	PauseEnd       []time.Time     // pause end times history, most recent first
	PauseQuantiles []time.Duration
}

// ReadGCStats reads statistics about garbage collection into stats.
// The number of entries in the pause history is system-dependent;
// stats.Pause slice will be reused if large enough, reallocated otherwise.
// ReadGCStats may use the full capacity of the stats.Pause slice.
// If stats.PauseQuantiles is non-empty, ReadGCStats fills it with quantiles
// summarizing the distribution of pause time. For example, if
// len(stats.PauseQuantiles) is 5, it will be filled with the minimum,
// 25%, 50%, 75%, and maximum pause times.
func ReadGCStats(stats *GCStats) {
	// Create a buffer with space for at least two copies of the
	// pause history tracked by the runtime. One will be returned
	// to the caller and the other will be used as transfer buffer
	// for end times history and as a temporary buffer for
	// computing quantiles.
	const maxPause = len(((*runtime.MemStats)(nil)).PauseNs)
	if cap(stats.Pause) < 2*maxPause+3 {
		stats.Pause = make([]time.Duration, 2*maxPause+3)
	}

	// readGCStats fills in the pause and end times histories (up to
	// maxPause entries) and then three more: Unix ns time of last GC,
	// number of GC, and total pause time in nanoseconds. Here we
	// depend on the fact that time.Duration's native unit is
	// nanoseconds, so the pauses and the total pause time do not need
	// any conversion.
	readGCStats(&stats.Pause)
	n := len(stats.Pause) - 3
	stats.LastGC = time.Unix(0, int64(stats.Pause[n]))
	stats.NumGC = int64(stats.Pause[n+1])
	stats.PauseTotal = stats.Pause[n+2]
	n /= 2 // buffer holds pauses and end times
	stats.Pause = stats.Pause[:n]

	if cap(stats.PauseEnd) < maxPause {
		stats.PauseEnd = make([]time.Time, 0, maxPause)
	}
	stats.PauseEnd = stats.PauseEnd[:0]
	for _, ns := range stats.Pause[n : n+n] {
		stats.PauseEnd = append(stats.PauseEnd, time.Unix(0, int64(ns)))
	}

	if len(stats.PauseQuantiles) > 0 {
		if n == 0 {
			for i := range stats.PauseQuantiles {
				stats.PauseQuantiles[i] = 0
			}
		} else {
			// There's room for a second copy of the data in stats.Pause.
			// See the allocation at the top of the function.
			sorted := stats.Pause[n : n+n]
			copy(sorted, stats.Pause)
			sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
			nq := len(stats.PauseQuantiles) - 1
			for i := 0; i < nq; i++ {
				stats.PauseQuantiles[i] = sorted[len(sorted)*i/nq]
			}
			stats.PauseQuantiles[nq] = sorted[len(sorted)-1]
		}
	}
}

// SetGCPercent sets the garbage collection target percentage:
// a collection is triggered when the ratio of freshly allocated data
// to live data remaining after the previous collection reaches this percentage.
// SetGCPercent returns the previous setting.
// The initial setting is the value of the GOGC environment variable
// at startup, or 100 if the variable is not set.
// A negative percentage disables garbage collection.
func SetGCPercent(percent int) int {
	return int(setGCPercent(int32(percent)))
}

// FreeOSMemory forces a garbage collection followed by an
// attempt to return as much memory to the operating system
// as possible. (Even if this is not called, the runtime gradually
// returns memory to the operating system in a background task.)
func FreeOSMemory() {
	freeOSMemory()
}

// SetMaxStack sets the maximum amount of memory that
// can be used by a single goroutine stack.
// If any goroutine exceeds this limit while growing its stack,
// the program crashes.
// SetMaxStack returns the previous setting.
// The initial setting is 1 GB on 64-bit systems, 250 MB on 32-bit systems.
//
// SetMaxStack is useful mainly for limiting the damage done by
// goroutines that enter an infinite recursion. It only limits future
// stack growth.
func SetMaxStack(bytes int) int {
	return setMaxStack(bytes)
}

// SetMaxThreads sets the maximum number of operating system
// threads that the Go program can use. If it attempts to use more than
// this many, the program crashes.
// SetMaxThreads returns the previous setting.
// The initial setting is 10,000 threads.
//
// The limit controls the number of operating system threads, not the number
// of goroutines. A Go program creates a new thread only when a goroutine
// is ready to run but all the existing threads are blocked in system calls, cgo calls,
// or are locked to other goroutines due to use of runtime.LockOSThread.
//
// SetMaxThreads is useful mainly for limiting the damage done by
// programs that create an unbounded number of threads. The idea is
// to take down the program before it takes down the operating system.
func SetMaxThreads(threads int) int {
	return setMaxThreads(threads)
}

// SetPanicOnFault controls the runtime's behavior when a program faults
// at an unexpected (non-nil) address. Such faults are typically caused by
// bugs such as runtime memory corruption, so the default response is to crash
// the program. Programs working with memory-mapped files or unsafe
// manipulation of memory may cause faults at non-nil addresses in less
// dramatic situations; SetPanicOnFault allows such programs to request
// that the runtime trigger only a panic, not a crash.
// SetPanicOnFault applies only to the current goroutine.
// It returns the previous setting.
func SetPanicOnFault(enabled bool) bool {
	return setPanicOnFault(enabled)
}

// WriteHeapDump writes a description of the heap and the objects in
// it to the given file descriptor.
//
// WriteHeapDump suspends the execution of all goroutines until the heap
// dump is completely written.  Thus, the file descriptor must not be
// connected to a pipe or socket whose other end is in the same Go
// process; instead, use a temporary file or network socket.
//
// The heap dump format is defined at https://golang.org/s/go15heapdump.
func WriteHeapDump(fd uintptr)

// SetTraceback sets the amount of detail printed by the runtime in
// the traceback it prints before exiting due to an unrecovered panic
// or an internal runtime error.
// The level argument takes the same values as the GOTRACEBACK
// environment variable. For example, SetTraceback("all") ensure
// that the program prints all goroutines when it crashes.
// See the package runtime documentation for details.
// If SetTraceback is called with a level lower than that of the
// environment variable, the call is ignored.
func SetTraceback(level string)
