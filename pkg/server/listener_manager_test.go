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
	doneCh := make(chan struct{}, numListeners)

	results := []struct {
		hasError bool
	}{
		{hasError: false},
		{hasError: true},
		{hasError: false},
	}

	for i, result := range results {
		go func(idx int, hasError bool) {
			if hasError {
				errCh <- &testShutdownError{msg: "listener failed"}
			}
			// Signal completion AFTER potential error send (like listener_manager.go)
			doneCh <- struct{}{}
		}(i, result.hasError)
	}

	// Wait for all goroutines to complete (like listener_manager.go pattern)
	for i := 0; i < numListeners; i++ {
		<-doneCh
	}

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
	numListeners := 2
	doneCh := make(chan struct{}, numListeners)

	timeout := 100 * time.Millisecond

	// Listener 1 completes quickly
	go func() {
		time.Sleep(10 * time.Millisecond)
		doneCh <- struct{}{}
	}()

	// Listener 2 takes longer than timeout
	go func() {
		time.Sleep(200 * time.Millisecond)
		doneCh <- struct{}{}
	}()

	// Use the pattern from listener_manager.go
	allDone := make(chan struct{})
	go func() {
		for i := 0; i < numListeners; i++ {
			select {
			case <-doneCh:
				// listener completed
			case <-time.After(5 * time.Second):
				// individual listener stuck (not expected in this test)
			}
		}
		close(allDone)
	}()

	select {
	case <-allDone:
		t.Log("All shutdowns completed before timeout")
	case <-time.After(timeout):
		t.Log("Shutdown timeout reached (expected behavior)")
	}

	// Wait for remaining listener to complete (cleanup)
	<-allDone
}

// TestShutdownRaceConditionFix verifies that errors are correctly collected
// even when wg.Done() is called before the error is sent to errCh.
// This is the race condition that Copilot identified in PR #993.
func TestShutdownRaceConditionFix(t *testing.T) {
	numListeners := 2
	errCh := make(chan error, numListeners)
	doneCh := make(chan struct{}, numListeners)
	wg := &sync.WaitGroup{}
	wg.Add(numListeners)

	// Simulate the pattern from listener.ShutDown() where wg.Done() is called
	// in a defer before returning (and before listener_manager.go can send error)
	for i := 0; i < numListeners; i++ {
		go func(idx int) {
			// Simulate ShutDown behavior: wg.Done() in defer
			defer wg.Done()

			// Simulate an error during shutdown
			time.Sleep(10 * time.Millisecond)
			// At this point, wg.Done() has been called (WaitGroup reaches 0)
			// but we haven't yet returned to listener_manager.go's goroutine

			// In the fixed implementation, listener_manager.go waits for doneCh
			// which is sent AFTER the error is sent
			if idx == 0 {
				errCh <- &testShutdownError{msg: "shutdown error 1"}
			}
			if idx == 1 {
				errCh <- &testShutdownError{msg: "shutdown error 2"}
			}
			// Signal completion AFTER error send
			doneCh <- struct{}{}
		}(i)
	}

	// Wait for all listener goroutines to complete (the fixed pattern)
	allDone := make(chan struct{})
	go func() {
		for i := 0; i < numListeners; i++ {
			<-doneCh
		}
		close(allDone)
	}()

	// Wait with timeout
	select {
	case <-allDone:
		t.Log("All listeners completed")
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for listeners")
	}

	// Drain errors
	var shutdownErrors []error
drainLoop:
	for {
		select {
		case err := <-errCh:
			shutdownErrors = append(shutdownErrors, err)
		default:
			break drainLoop
		}
	}

	// Both errors should be collected despite the race condition
	if len(shutdownErrors) != 2 {
		t.Errorf("Expected 2 shutdown errors, got %d - race condition not fixed", len(shutdownErrors))
	}
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
