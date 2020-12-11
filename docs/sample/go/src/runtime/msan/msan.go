// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build msan,linux
// +build amd64 arm64

package msan

/*
#cgo CFLAGS: -fsanitize=memory
#cgo LDFLAGS: -fsanitize=memory

#include <stdint.h>
#include <sanitizer/msan_interface.h>

void __msan_read_go(void *addr, uintptr_t sz) {
	__msan_check_mem_is_initialized(addr, sz);
}

void __msan_write_go(void *addr, uintptr_t sz) {
	__msan_unpoison(addr, sz);
}

void __msan_malloc_go(void *addr, uintptr_t sz) {
	__msan_unpoison(addr, sz);
}

void __msan_free_go(void *addr, uintptr_t sz) {
	__msan_poison(addr, sz);
}
*/
import "C"
