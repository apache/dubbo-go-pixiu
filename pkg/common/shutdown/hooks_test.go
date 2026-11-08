/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shutdown

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"
)

func TestHooksRunOnceInReverseOrder(t *testing.T) {
	var h hooks
	var order []int
	h.register(func() { order = append(order, 1) })
	h.register(func() { order = append(order, 2) })
	h.run()
	h.run()
	require.Equal(t, []int{2, 1}, order)
}

func TestConcurrentHooksWaitForCleanup(t *testing.T) {
	var h hooks
	var calls atomic.Int32
	entered := make(chan struct{})
	release := make(chan struct{})
	h.register(func() { calls.Add(1); close(entered); <-release })
	var workers sync.WaitGroup
	workers.Add(2)
	for i := 0; i < 2; i++ {
		go func() { defer workers.Done(); h.run() }()
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("cleanup did not start")
	}
	close(release)
	done := make(chan struct{})
	go func() { workers.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cleanup did not finish")
	}
	require.Equal(t, int32(1), calls.Load())
}

func TestProcessHooks(t *testing.T) {
	calls := 0
	RegisterHook(func() { calls++ })
	RunHooks()
	RunHooks()
	require.Equal(t, 1, calls)
}
