/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package exponentialbackoff

import (
	"testing"
	"time"
)

// TestExponentialBackoffRetry_Factory verifies the creation of policies via the factory.
func TestExponentialBackoffRetry_Factory(t *testing.T) {
	t.Run("should create policy with valid config", func(t *testing.T) {
		config := map[string]any{
			"times":           uint(3),
			"initialInterval": "50ms",
			"maxInterval":     "500ms",
			"multiplier":      2.0,
		}

		policy, err := newExponentialBackoffRetry(config)
		if err != nil {
			t.Fatalf("newExponentialBackoffRetry failed with valid config: %v", err)
		}

		if policy == nil {
			t.Fatal("newExponentialBackoffRetry returned a nil policy")
		}

		p, ok := policy.(*ExponentialBackoffRetry)
		if !ok {
			t.Fatalf("newExponentialBackoffRetry returned wrong type: got %T", policy)
		}

		if p.MaxAttempts != 4 { // 3 retries + 1 initial
			t.Errorf("expected MaxAttempts to be 4, got %d", p.MaxAttempts)
		}
		if p.InitialInterval != 50*time.Millisecond {
			t.Errorf("expected InitialInterval to be 50ms, got %v", p.InitialInterval)
		}
	})

	t.Run("should fail with invalid config", func(t *testing.T) {
		testCases := []struct {
			name   string
			config map[string]any
		}{
			{
				name: "invalid initialInterval format",
				config: map[string]any{
					"initialInterval": "50someting",
					"maxInterval":     "500ms",
				},
			},
			{
				name: "invalid maxInterval format",
				config: map[string]any{
					"initialInterval": "50ms",
					"maxInterval":     "500",
				},
			},
			{
				name: "invalid type for times",
				config: map[string]any{
					"times": "three",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := newExponentialBackoffRetry(tc.config)
				if err == nil {
					t.Error("expected an error for invalid config but got nil")
				}
			})
		}
	})
}

// TestExponentialBackoffRetry_Attempts verifies the number of attempts.
func TestExponentialBackoffRetry_Attempts(t *testing.T) {
	config := map[string]any{
		"times":           uint(2),
		"initialInterval": "1ms", // Use very short intervals to not slow down test
		"maxInterval":     "5ms",
		"multiplier":      2.0,
	}
	policy, _ := newExponentialBackoffRetry(config)

	// Should allow 3 attempts (1 initial + 2 retries)
	for i := 0; i < 3; i++ {
		if !policy.Attempt() {
			t.Fatalf("attempt %d should have been allowed, but was blocked", i+1)
		}
	}

	// The 4th attempt should be blocked
	if policy.Attempt() {
		t.Fatal("4th attempt should have been blocked, but was allowed")
	}
}

// TestExponentialBackoffRetry_Timing verifies the backoff delay.
func TestExponentialBackoffRetry_Timing(t *testing.T) {
	initialInterval := 20 * time.Millisecond
	multiplier := 2.0
	maxJitter := 100 * time.Millisecond

	config := map[string]any{
		"times":           uint(3),
		"initialInterval": initialInterval.String(),
		"maxInterval":     "1s",
		"multiplier":      multiplier,
	}
	policy, _ := newExponentialBackoffRetry(config)

	// 1st attempt: no delay
	start := time.Now()
	if !policy.Attempt() {
		t.Fatal("first attempt failed")
	}
	elapsed := time.Since(start)
	if elapsed >= initialInterval {
		t.Errorf("first attempt should have no delay, but took %v", elapsed)
	}

	// 2nd attempt (1st retry): delay should be ~20ms
	expectedDelay1 := initialInterval
	start = time.Now()
	if !policy.Attempt() {
		t.Fatal("second attempt failed")
	}
	elapsed = time.Since(start)
	if elapsed < expectedDelay1 || elapsed > expectedDelay1+maxJitter {
		t.Errorf("expected ~%v delay for 2nd attempt, but got %v", expectedDelay1, elapsed)
	}

	// 3rd attempt (2nd retry): delay should be ~40ms
	expectedDelay2 := time.Duration(float64(initialInterval) * multiplier)
	start = time.Now()
	if !policy.Attempt() {
		t.Fatal("third attempt failed")
	}
	elapsed = time.Since(start)
	if elapsed < expectedDelay2 || elapsed > expectedDelay2+maxJitter {
		t.Errorf("expected ~%v delay for 3rd attempt, but got %v", expectedDelay2, elapsed)
	}
}

// TestExponentialBackoffRetry_Reset verifies the policy can be reused after reset.
func TestExponentialBackoffRetry_Reset(t *testing.T) {
	config := map[string]any{
		"times":           uint(1),
		"initialInterval": "1ms",
		"maxInterval":     "5ms",
	}
	policy, _ := newExponentialBackoffRetry(config)

	// Run through a full cycle
	if !policy.Attempt() {
		t.Fatal("first attempt in first cycle failed")
	}
	if !policy.Attempt() {
		t.Fatal("second attempt in first cycle failed")
	}
	if policy.Attempt() {
		t.Fatal("policy allowed too many attempts in first cycle")
	}

	// Reset the policy
	policy.Reset()

	// Run through a second cycle, it should behave identically
	if !policy.Attempt() {
		t.Fatal("first attempt in second cycle failed after reset")
	}
	if !policy.Attempt() {
		t.Fatal("second attempt in second cycle failed after reset")
	}
	if policy.Attempt() {
		t.Fatal("policy allowed too many attempts in second cycle after reset")
	}
}
