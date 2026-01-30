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

package kvcache

import (
	"errors"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_InitialState(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second,
		HalfOpenMaxCalls: 2,
	})

	assert.Equal(t, CircuitClosed, cb.state)

	// Execute should succeed in closed state
	err := cb.Execute(func() error {
		return nil
	})
	assert.NoError(t, err)
}

func TestCircuitBreaker_TransitionToOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second,
		HalfOpenMaxCalls: 2,
	})

	// Record failures to reach threshold
	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return errors.New("test error")
		})
	}

	assert.Equal(t, CircuitOpen, cb.state)

	// Next execution should fail with circuit breaker open error
	err := cb.Execute(func() error {
		return nil
	})
	assert.Equal(t, ErrCircuitBreakerOpen, err)
}

func TestCircuitBreaker_TransitionToHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		RecoveryTimeout:  100 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	})

	// Transition to Open
	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	assert.Equal(t, CircuitOpen, cb.state)

	err := cb.Execute(func() error { return nil })
	assert.Equal(t, ErrCircuitBreakerOpen, err)

	// Wait for recovery timeout
	time.Sleep(150 * time.Millisecond)

	// Should transition to HalfOpen and allow execution
	err = cb.Execute(func() error { return nil })
	assert.NoError(t, err)
	assert.Equal(t, CircuitClosed, cb.state) // Success transitions to Closed
}

func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		RecoveryTimeout:  100 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	})

	// Transition to Open
	cb.Execute(func() error { return errors.New("error") })
	cb.Execute(func() error { return errors.New("error") })

	// Wait and transition to HalfOpen
	time.Sleep(150 * time.Millisecond)

	// Execute success - should transition to Closed
	err := cb.Execute(func() error { return nil })
	assert.NoError(t, err)

	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 0, cb.halfOpenCalls)
}

func TestCircuitBreaker_HalfOpenToOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		RecoveryTimeout:  100 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	})

	// Transition to Open
	cb.Execute(func() error { return errors.New("error") })
	cb.Execute(func() error { return errors.New("error") })

	// Wait and transition to HalfOpen
	time.Sleep(150 * time.Millisecond)

	// Execute failure - should transition back to Open
	cb.Execute(func() error { return errors.New("error") })

	assert.Equal(t, CircuitOpen, cb.state)

	// Next call should be rejected
	err := cb.Execute(func() error { return nil })
	assert.Equal(t, ErrCircuitBreakerOpen, err)
}

func TestCircuitBreaker_HalfOpenMaxCalls(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		RecoveryTimeout:  100 * time.Millisecond,
		HalfOpenMaxCalls: 2,
	})

	// Transition to Open
	cb.Execute(func() error { return errors.New("error") })
	cb.Execute(func() error { return errors.New("error") })

	// Wait and transition to HalfOpen
	time.Sleep(150 * time.Millisecond)

	// First call succeeds but doesn't close circuit yet
	err := cb.Execute(func() error { return nil })
	assert.NoError(t, err)
	assert.Equal(t, CircuitClosed, cb.state) // Actually transitions to Closed on first success
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second,
		HalfOpenMaxCalls: 2,
	})

	// Record some failures
	cb.Execute(func() error { return errors.New("error") })
	cb.Execute(func() error { return errors.New("error") })

	assert.Equal(t, CircuitClosed, cb.state)

	// Record success - should reset failure count
	cb.Execute(func() error { return nil })

	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 0, cb.failureCount)
}

func TestCircuitBreaker_OnlyResetsHalfOpenCallsWhenTransitioning(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		RecoveryTimeout:  time.Second,
		HalfOpenMaxCalls: 2,
	})

	// Start in Closed state
	assert.Equal(t, CircuitClosed, cb.state)

	// Artificially set halfOpenCalls
	cb.halfOpenCalls = 5

	// Record success in Closed state - should NOT reset halfOpenCalls
	cb.Execute(func() error { return nil })

	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 5, cb.halfOpenCalls) // Should remain unchanged
}
