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

package countbased

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

// TestCountBasedRetry_Factory verifies the creation of policies via the factory.
func TestCountBasedRetry_Factory(t *testing.T) {
	t.Run("should create policy with valid config", func(t *testing.T) {
		config := map[string]any{
			"times": 2, // Corresponds to 2 retries
		}
		policy, err := newCountBasedRetry(config)
		if err != nil {
			t.Fatalf("newCountBasedRetry failed with a valid config: %v", err)
		}

		p, ok := policy.(*CountBasedRetry)
		if !ok {
			t.Fatalf("newCountBasedRetry returned wrong type: got %T", policy)
		}

		// 1 initial try + 2 retries = 3 attempts
		if p.MaxAttempts != 3 {
			t.Errorf("expected MaxAttempts to be 3, but got %d", p.MaxAttempts)
		}
	})

	t.Run("should fail when times key is missing", func(t *testing.T) {
		config := map[string]any{
			"other_key": "some_value",
		}
		policy, _ := newCountBasedRetry(config)
		assert.Equal(t, policy.(*CountBasedRetry).MaxAttempts, uint(3), "expected 'times' key to default to 3")
	})

	t.Run("should fail when times has invalid type", func(t *testing.T) {
		config := map[string]any{
			"times": "two",
		}
		_, err := newCountBasedRetry(config)
		if err == nil {
			t.Error("expected an error for invalid type of 'times', but got nil")
		}
	})
}

// TestCountBasedRetry_Attempts verifies the number of attempts.
func TestCountBasedRetry_Attempts(t *testing.T) {
	config := map[string]any{"times": 2}
	policy, _ := newCountBasedRetry(config)

	// Should allow exactly 3 attempts (1 initial + 2 retries)
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

// TestCountBasedRetry_Reset verifies the policy can be reused after reset.
func TestCountBasedRetry_Reset(t *testing.T) {
	config := map[string]any{"times": 1} // 2 total attempts
	policy, _ := newCountBasedRetry(config)

	// --- First Cycle ---
	// Run through a full cycle of attempts
	if !policy.Attempt() {
		t.Fatal("first attempt in first cycle failed")
	}
	if !policy.Attempt() {
		t.Fatal("second attempt in first cycle failed")
	}
	// Verify it's exhausted
	if policy.Attempt() {
		t.Fatal("policy allowed too many attempts in first cycle")
	}

	// --- Reset ---
	policy.Reset()

	// --- Second Cycle ---
	// It should behave identically to the first cycle
	if !policy.Attempt() {
		t.Fatal("first attempt in second cycle failed after reset")
	}
	if !policy.Attempt() {
		t.Fatal("second attempt in second cycle failed after reset")
	}
	// Verify it's exhausted again
	if policy.Attempt() {
		t.Fatal("policy allowed too many attempts in second cycle after reset")
	}
}
