// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

func gcstats(name string, n int, t time.Duration) {
	st := new(runtime.MemStats)
	runtime.ReadMemStats(st)
	nprocs := runtime.GOMAXPROCS(-1)
	cpus := ""
	if nprocs != 1 {
		cpus = fmt.Sprintf("-%d", nprocs)
	}
	fmt.Printf("garbage.%sMem%s Alloc=%d/%d Heap=%d NextGC=%d Mallocs=%d\n", name, cpus, st.Alloc, st.TotalAlloc, st.Sys, st.NextGC, st.Mallocs)
	fmt.Printf("garbage.%s%s %d %d ns/op\n", name, cpus, n, t.Nanoseconds()/int64(n))
	fmt.Printf("garbage.%sLastPause%s 1 %d ns/op\n", name, cpus, st.PauseNs[(st.NumGC-1)%uint32(len(st.PauseNs))])
	fmt.Printf("garbage.%sPause%s %d %d ns/op\n", name, cpus, st.NumGC, int64(st.PauseTotalNs)/int64(st.NumGC))
	nn := int(st.NumGC)
	if nn >= len(st.PauseNs) {
		nn = len(st.PauseNs)
	}
	t1, t2, t3, t4, t5 := tukey5(st.PauseNs[0:nn])
	fmt.Printf("garbage.%sPause5%s: %d %d %d %d %d\n", name, cpus, t1, t2, t3, t4, t5)

	//	fmt.Printf("garbage.%sScan: %v\n", name, st.ScanDist)
}

type T []uint64

func (t T) Len() int           { return len(t) }
func (t T) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t T) Less(i, j int) bool { return t[i] < t[j] }

func tukey5(raw []uint64) (lo, q1, q2, q3, hi uint64) {
	x := make(T, len(raw))
	copy(x, raw)
	sort.Sort(T(x))
	lo = x[0]
	q1 = x[len(x)/4]
	q2 = x[len(x)/2]
	q3 = x[len(x)*3/4]
	hi = x[len(x)-1]
	return
}
