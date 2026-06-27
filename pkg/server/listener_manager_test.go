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
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// mockListenerService is a mock implementation of listener.ListenerService for testing.
type mockListenerService struct {
	startErr    error
	closeErr    error
	shutdownErr error
	refreshErr  error
	shutdownWg  *sync.WaitGroup // Will be set when ShutDown is called
}

func (m *mockListenerService) Start() error {
	return m.startErr
}

func (m *mockListenerService) Close() error {
	return m.closeErr
}

func (m *mockListenerService) ShutDown(wg any) error {
	// Cast the WaitGroup and call Done() to simulate real behavior
	if w, ok := wg.(*sync.WaitGroup); ok && w != nil {
		w.Done()
	}
	return m.shutdownErr
}

func (m *mockListenerService) Refresh(_ model.Listener) error {
	return m.refreshErr
}

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, defaultPerListenerTimeout)

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

// TestShutdownListenersStuckListener tests that shutdownListeners handles
// individual listeners that get stuck (exceed the per-listener timeout).
func TestShutdownListenersStuckListener(t *testing.T) {
	// Use a short per-listener timeout for testing
	perListenerTimeout := 100 * time.Millisecond
	// Create a listener that will get stuck
	shutdownFuncs := []ShutdownFunc{
		func() error { return nil }, // fast listener
		func() error {
			// This listener blocks longer than perListenerTimeout
			time.Sleep(200 * time.Millisecond)
			return nil
		},
	}

	// Use a timeout longer than perListenerTimeout so we can observe the stuck behavior
	timeout := 300 * time.Millisecond
	errs, timedOut := shutdownListeners(shutdownFuncs, timeout, perListenerTimeout)

	// Should not timeout because overall timeout is longer than listener sleep
	if timedOut {
		t.Error("Did not expect overall timeout")
	}
	// No errors because listeners return nil
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d", len(errs))
	}
}

// TestHandleShutdownSignalNoErrors tests handleShutdownSignal when all listeners
// shut down successfully.
func TestHandleShutdownSignalNoErrors(t *testing.T) {
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			"listener1": {ListenerService: &mockListenerService{shutdownErr: nil}},
			"listener2": {ListenerService: &mockListenerService{shutdownErr: nil}},
		},
		shutdownWG: &sync.WaitGroup{},
	}

	errs, timedOut := lm.handleShutdownSignal(os.Interrupt, 1*time.Second)

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout")
	}
}

// TestHandleShutdownSignalWithErrors tests handleShutdownSignal when some listeners
// fail to shut down.
func TestHandleShutdownSignalWithErrors(t *testing.T) {
	testErr := errors.New("shutdown failed")
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			"listener1": {ListenerService: &mockListenerService{shutdownErr: nil}},
			"listener2": {ListenerService: &mockListenerService{shutdownErr: testErr}},
		},
		shutdownWG: &sync.WaitGroup{},
	}

	errs, timedOut := lm.handleShutdownSignal(os.Interrupt, 1*time.Second)

	if len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout")
	}
	if errs[0].Error() != testErr.Error() {
		t.Errorf("Expected error '%s', got '%s'", testErr.Error(), errs[0].Error())
	}
}

// TestHandleShutdownSignalNoListeners tests handleShutdownSignal with no listeners.
func TestHandleShutdownSignalNoListeners(t *testing.T) {
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{},
		shutdownWG:            &sync.WaitGroup{},
	}

	errs, timedOut := lm.handleShutdownSignal(os.Interrupt, 1*time.Second)

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout")
	}
}

// TestHandleShutdownSignalTimeout tests handleShutdownSignal when shutdown times out.
func TestHandleShutdownSignalTimeout(t *testing.T) {
	// Create a slow listener that takes longer than the timeout
	slowListener := &mockSlowListenerService{duration: 100 * time.Millisecond}
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			"listener1": {ListenerService: slowListener},
		},
		shutdownWG: &sync.WaitGroup{},
	}

	// Very short timeout to trigger timeout condition (less than slowListener duration)
	errs, timedOut := lm.handleShutdownSignal(os.Interrupt, 10*time.Millisecond)

	// Should timeout because listener takes longer than the timeout
	if !timedOut {
		t.Error("Expected timeout")
	}
	if len(errs) != 0 {
		t.Errorf("Expected no errors (listener returns nil), got %d", len(errs))
	}
}

// TestHandleShutdownSignalDumpHeap tests handleShutdownSignal with a dump heap signal.
// Note: We skip actually writing heap dump in tests, but verify the code path is reached.
func TestHandleShutdownSignalDumpHeap(t *testing.T) {
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			"listener1": {ListenerService: &mockListenerService{}},
		},
		shutdownWG: &sync.WaitGroup{},
	}

	// SIGQUIT is typically in DumpHeapShutdownSignals
	errs, timedOut := lm.handleShutdownSignal(os.Signal(syscall.SIGQUIT), 1*time.Second)

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d", len(errs))
	}
	if timedOut {
		t.Error("Expected no timeout")
	}
}

// mockSlowListenerService is a mock that takes time to shutdown for testing timeout scenarios.
type mockSlowListenerService struct {
	duration time.Duration
}

func (m *mockSlowListenerService) Start() error { return nil }
func (m *mockSlowListenerService) Close() error { return nil }
func (m *mockSlowListenerService) ShutDown(wg any) error {
	time.Sleep(m.duration)
	if w, ok := wg.(*sync.WaitGroup); ok && w != nil {
		w.Done()
	}
	return nil
}
func (m *mockSlowListenerService) Refresh(_ model.Listener) error { return nil }
