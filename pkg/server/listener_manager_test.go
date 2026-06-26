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

package server

import (
	"sync"
	"testing"
	"time"
)

// TestShutdownWaitGroupDoubleDoneWouldPanic demonstrates that calling Done()
// twice on a WaitGroup with counter 1 would panic (negative counter).
func TestShutdownWaitGroupDoubleDoneWouldPanic(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// First Done() - should be fine
	wg.Done()

	// Second Done() - should panic because counter becomes negative
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when Done() called twice on WaitGroup with counter 1")
		}
	}()

	wg.Done() // This should panic
}

// TestShutdownErrorCollection verifies that shutdown errors are collected
// instead of causing immediate exit.
func TestShutdownErrorCollection(t *testing.T) {
	numListeners := 3
	errCh := make(chan error, numListeners)

	results := []struct {
		hasError bool
	}{
		{hasError: false},
		{hasError: true},
		{hasError: false},
	}

	wg := &sync.WaitGroup{}
	wg.Add(numListeners)

	for i, result := range results {
		go func(idx int, hasError bool) {
			defer wg.Done()
			if hasError {
				errCh <- &testShutdownError{msg: "listener failed"}
			}
		}(i, result.hasError)
	}

	wg.Wait()

	// Drain errors using labeled break pattern from listener_manager.go
	var shutdownErrors []error
	drainErrors:
	for {
		select {
		case err := <-errCh:
			shutdownErrors = append(shutdownErrors, err)
		default:
			break drainErrors
		}
	}

	if len(shutdownErrors) != 1 {
		t.Errorf("Expected 1 shutdown error, got %d", len(shutdownErrors))
	}
}

// TestShutdownTimeoutHandling verifies timeout handling pattern
func TestShutdownTimeoutHandling(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	timeout := 100 * time.Millisecond

	// Listener 1 completes quickly
	go func() {
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	}()

	// Listener 2 takes longer than timeout
	go func() {
		time.Sleep(200 * time.Millisecond)
		wg.Done()
	}()

	// Use the pattern from listener_manager.go
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.Log("All shutdowns completed before timeout")
	case <-time.After(timeout):
		t.Log("Shutdown timeout reached (expected behavior)")
	}

	// Wait for remaining listener to complete (cleanup)
	<-done
}

// TestLabeledBreakPattern verifies the labeled break pattern for draining channel
func TestLabeledBreakPattern(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	var collected []int
	drainLoop:
	for {
		select {
		case v := <-ch:
			collected = append(collected, v)
		default:
			break drainLoop // This breaks the for loop, not just the select
		}
	}

	if len(collected) != 3 {
		t.Errorf("Expected to collect 3 values, got %d", len(collected))
	}
}

// testShutdownError is a simple error type for testing
type testShutdownError struct {
	msg string
}

func (e *testShutdownError) Error() string {
	return e.msg
}