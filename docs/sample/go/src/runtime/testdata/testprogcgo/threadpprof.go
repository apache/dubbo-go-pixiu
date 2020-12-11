// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!windows

package main

// Run a slow C function saving a CPU profile.

/*
#include <stdint.h>
#include <time.h>
#include <pthread.h>

int threadSalt1;
int threadSalt2;

void cpuHogThread() {
	int foo = threadSalt1;
	int i;

	for (i = 0; i < 100000; i++) {
		if (foo > 0) {
			foo *= foo;
		} else {
			foo *= foo + 1;
		}
	}
	threadSalt2 = foo;
}

void cpuHogThread2() {
}

static int cpuHogThreadCount;

struct cgoTracebackArg {
	uintptr_t  context;
	uintptr_t  sigContext;
	uintptr_t* buf;
	uintptr_t  max;
};

// pprofCgoThreadTraceback is passed to runtime.SetCgoTraceback.
// For testing purposes it pretends that all CPU hits in C code are in cpuHog.
void pprofCgoThreadTraceback(void* parg) {
	struct cgoTracebackArg* arg = (struct cgoTracebackArg*)(parg);
	arg->buf[0] = (uintptr_t)(cpuHogThread) + 0x10;
	arg->buf[1] = (uintptr_t)(cpuHogThread2) + 0x4;
	arg->buf[2] = 0;
	__sync_add_and_fetch(&cpuHogThreadCount, 1);
}

// getCPUHogThreadCount fetches the number of times we've seen cpuHogThread
// in the traceback.
int getCPUHogThreadCount() {
	return __sync_add_and_fetch(&cpuHogThreadCount, 0);
}

static void* cpuHogDriver(void* arg __attribute__ ((unused))) {
	while (1) {
		cpuHogThread();
	}
	return 0;
}

void runCPUHogThread(void) {
	pthread_t tid;
	pthread_create(&tid, 0, cpuHogDriver, 0);
}
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
	"unsafe"
)

func init() {
	register("CgoPprofThread", CgoPprofThread)
	register("CgoPprofThreadNoTraceback", CgoPprofThreadNoTraceback)
}

func CgoPprofThread() {
	runtime.SetCgoTraceback(0, unsafe.Pointer(C.pprofCgoThreadTraceback), nil, nil)
	pprofThread()
}

func CgoPprofThreadNoTraceback() {
	pprofThread()
}

func pprofThread() {
	f, err := ioutil.TempFile("", "prof")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	C.runCPUHogThread()

	t0 := time.Now()
	for C.getCPUHogThreadCount() < 2 && time.Since(t0) < time.Second {
		time.Sleep(100 * time.Millisecond)
	}

	pprof.StopCPUProfile()

	name := f.Name()
	if err := f.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(name)
}
