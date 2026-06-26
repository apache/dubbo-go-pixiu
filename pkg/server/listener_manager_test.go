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
	"errors"
	"sync"
	"testing"
	"time"
)

// TestShutdownListenersNoErrors tests shutdownListeners when all listeners
// shut down successfully without errors.
func TestShutdownListenersNoErrors(t *testing.T) {
	// Create fake listeners that shut down successfully
	shutdownFuncs := []ShutdownFunc{
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
	}

	timeout := 1 * time.Second
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(errs), errs)
	}
	if timedOut {
		t.Error("Expected no timeout, but timed out")
	}
}

// TestShutdownListenersWithErrors tests shutdownListeners when some listeners
// return errors during shutdown.
func TestShutdownListenersWithErrors(t *testing.T) {
	testErr := errors.New("shutdown failed")
	shutdownFuncs := []ShutdownFunc{
		func() error { return nil },
		func() error { return testErr },
		func() error { return nil },
	}

	timeout := 1 * time.Second
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	if len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout, but timed out")
	}
	if errs[0].Error() != testErr.Error() {
		t.Errorf("Expected error '%s', got '%s'", testErr.Error(), errs[0].Error())
	}
}

// TestShutdownListenersAllErrors tests shutdownListeners when all listeners
// return errors.
func TestShutdownListenersAllErrors(t *testing.T) {
	shutdownFuncs := []ShutdownFunc{
		func() error { return errors.New("error 1") },
		func() error { return errors.New("error 2") },
		func() error { return errors.New("error 3") },
	}

	timeout := 1 * time.Second
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	if len(errs) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout, but timed out")
	}
}

// TestShutdownListenersTimeout tests shutdownListeners when listeners take
// longer than the timeout.
func TestShutdownListenersTimeout(t *testing.T) {
	shutdownFuncs := []ShutdownFunc{
		func() error { time.Sleep(10 * time.Millisecond); return nil },
		func() error { time.Sleep(200 * time.Millisecond); return nil }, // exceeds timeout
		func() error { time.Sleep(10 * time.Millisecond); return nil },
	}

	timeout := 50 * time.Millisecond
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	// Should have timed out but still collected errors from fast listeners
	if !timedOut {
		t.Error("Expected timeout, but did not time out")
	}
	// No errors in this case (all listeners return nil), but timeout should be set
	if len(errs) != 0 {
		t.Errorf("Expected no errors (all returned nil), got %d", len(errs))
	}
}

// TestShutdownListenersRaceCondition verifies that errors are correctly collected
// even when listeners have slow error reporting after completing their work.
// This addresses the race condition Copilot identified in PR #993.
func TestShutdownListenersRaceCondition(t *testing.T) {
	// Simulate slow error send after shutdown completes
	shutdownFuncs := []ShutdownFunc{
		func() error {
			time.Sleep(10 * time.Millisecond)
			// Simulate slow error channel send
			return errors.New("delayed error 1")
		},
		func() error {
			time.Sleep(10 * time.Millisecond)
			return errors.New("delayed error 2")
		},
	}

	timeout := 1 * time.Second
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	// Both errors should be collected despite the delay
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got %d - race condition not fixed", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout, but timed out")
	}
}

// TestShutdownListenersEmpty tests shutdownListeners with no listeners.
func TestShutdownListenersEmpty(t *testing.T) {
	shutdownFuncs := []ShutdownFunc{}

	timeout := 1 * time.Second
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout, but timed out")
	}
}

// TestShutdownListenersTimeoutWithError tests that errors are collected
// even when timeout occurs.
func TestShutdownListenersTimeoutWithError(t *testing.T) {
	shutdownFuncs := []ShutdownFunc{
		func() error { return errors.New("fast error") },
		func() error { time.Sleep(200 * time.Millisecond); return nil }, // exceeds timeout
	}

	timeout := 50 * time.Millisecond
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout)

	if !timedOut {
		t.Error("Expected timeout, but did not time out")
	}
	// The fast error should still be collected
	if len(errs) != 1 {
		t.Errorf("Expected 1 error from fast listener, got %d", len(errs))
	}
}

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
