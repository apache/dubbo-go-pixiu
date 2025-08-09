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

package noretry

import (
	"testing"
)

// TestNoRetryPolicy verifies the behavior of the NoRetryPolicy.
func TestNoRetryPolicy(t *testing.T) {
	// Test case 1: Basic lifecycle - one attempt, no retries.
	t.Run("should allow first attempt and block retries", func(t *testing.T) {
		// Create a new policy instance using the factory function.
		policy, err := newNoRetryPolicy(nil)
		if err != nil {
			t.Fatalf("newNoRetryPolicy should not return an error, but got: %v", err)
		}

		// 1. First call to Attempt() should return true to allow the initial request.
		// We pass a dummy error to simulate a failed attempt.
		if !policy.Attempt() {
			t.Error("First call to Attempt() should return true, but got false")
		}

		// 2. Second call to Attempt() should return false, as no retries are allowed.
		if policy.Attempt() {
			t.Error("Second call to Attempt() should return false, but got true")
		}

		// 3. Any later call should also return false.
		if policy.Attempt() {
			t.Error("Third call to Attempt() should return false, but got true")
		}
	})

	// Test case 2: Verify the Reset() method correctly resets the state.
	t.Run("should allow a new attempt after reset", func(t *testing.T) {
		policy, _ := newNoRetryPolicy(nil)

		// Simulate a full cycle: one attempt, one failed retry.
		policy.Attempt() // First attempt
		if policy.Attempt() {
			t.Fatal("Policy allowed a retry before being reset")
		}

		// 4. Reset the policy. This simulates starting a new request cycle (e.g., on a fallback endpoint).
		policy.Reset()

		// 5. After resetting, the first call to Attempt() should once again return true.
		if !policy.Attempt() {
			t.Error("Attempt() should return true after Reset(), but got false")
		}

		// 6. And the later call should return false again.
		if policy.Attempt() {
			t.Error("Attempt() should return false on the second try after Reset(), but got true")
		}
	})
}

// TestNoRetryPolicy_Factory ensures the factory function works as expected.
func TestNoRetryPolicy_Factory(t *testing.T) {
	t.Run("factory should create a valid instance", func(t *testing.T) {
		policy, err := newNoRetryPolicy(nil) // config map is not used, so pass nil

		if err != nil {
			t.Fatalf("newNoRetryPolicy failed: %v", err)
		}

		if policy == nil {
			t.Fatal("newNoRetryPolicy returned a nil policy")
		}

		// Ensure the returned type is the one we expect.
		_, ok := policy.(*NoRetryPolicy)
		if !ok {
			t.Fatalf("newNoRetryPolicy returned the wrong type: got %T", policy)
		}
	})
}
