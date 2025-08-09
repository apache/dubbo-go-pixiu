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

package retry

import (
	"fmt"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// --- Test Helper Implementations ---
// dummyRetryPolicy is a mock implementation of the RetryPolicy interface for testing.
type dummyRetryPolicy struct {
	configValue string
}

func (d *dummyRetryPolicy) Attempt() bool { return false }
func (d *dummyRetryPolicy) Reset()        {}

// newDummyRetryer is a mock factory function for creating dummyRetryPolicy instances.
func newDummyRetryer(config map[string]any) (RetryPolicy, error) {
	val, ok := config["key"].(string)
	if !ok {
		return nil, fmt.Errorf("config missing 'key' or key is not a string")
	}
	return &dummyRetryPolicy{configValue: val}, nil
}

// anotherDummyRetryPolicy is a different implementation to test overwriting.
type anotherDummyRetryPolicy struct{}

func (d *anotherDummyRetryPolicy) Attempt() bool { return false }
func (d *anotherDummyRetryPolicy) Reset()        {}

func newAnotherDummyRetryer(config map[string]any) (RetryPolicy, error) {
	return &anotherDummyRetryPolicy{}, nil
}

// cleanupRegistry resets the global registry between tests to ensure they are isolated.
func cleanupRegistry() {
	retryPolicyRegistry = make(map[model.RetryType]RetryPolicyFactory)
}

// --- Test Cases ---

func TestRetryPolicyRegistry(t *testing.T) {

	t.Run("successful registration and retrieval", func(t *testing.T) {
		cleanupRegistry()
		policyName := model.RetryType("dummy")
		RegisterRetryPolicy(policyName, newDummyRetryer)

		endpoint := &model.Endpoint{
			ID: "test-endpoint-1",
			LLMMeta: &model.LLMMeta{
				RetryPolicy: model.RetryPolicy{
					Name: policyName,
					Config: map[string]any{
						"key": "test-value",
					},
				},
			},
		}

		policy, err := GetRetryPolicy(endpoint)
		if err != nil {
			t.Fatalf("GetRetryPolicy() returned an unexpected error: %v", err)
		}
		if policy == nil {
			t.Fatal("GetRetryPolicy() returned a nil policy")
		}

		p, ok := policy.(*dummyRetryPolicy)
		if !ok {
			t.Fatalf("Expected policy of type *dummyRetryPolicy, but got %T", policy)
		}
		if p.configValue != "test-value" {
			t.Errorf("Expected config value 'test-value', but got '%s'", p.configValue)
		}
	})

	t.Run("policy not found", func(t *testing.T) {
		cleanupRegistry()
		endpoint := &model.Endpoint{
			ID: "test-endpoint-2",
			LLMMeta: &model.LLMMeta{
				RetryPolicy: model.RetryPolicy{
					Name: "nonexistent",
				},
			},
		}

		_, err := GetRetryPolicy(endpoint)
		if err == nil {
			t.Fatal("GetRetryPolicy() should have returned an error for a nonexistent policy, but it was nil")
		}
	})

	t.Run("case-insensitive registration and retrieval", func(t *testing.T) {
		cleanupRegistry()
		// Register with a mixed-case name
		RegisterRetryPolicy("MyDummyPolicy", newDummyRetryer)

		// Retrieve with a lowercase name
		endpoint := &model.Endpoint{
			ID: "test-endpoint-3",
			LLMMeta: &model.LLMMeta{
				RetryPolicy: model.RetryPolicy{
					Name: "mydummypolicy",
					Config: map[string]any{
						"key": "value",
					},
				},
			},
		}

		policy, err := GetRetryPolicy(endpoint)
		if err != nil {
			t.Fatalf("GetRetryPolicy() failed with case-insensitive name: %v", err)
		}
		if _, ok := policy.(*dummyRetryPolicy); !ok {
			t.Fatalf("Expected policy of type *dummyRetryPolicy, but got %T", policy)
		}
	})

	t.Run("policy overwrite", func(t *testing.T) {
		cleanupRegistry()
		policyName := model.RetryType("overwrite_test")

		// Register the first implementation
		RegisterRetryPolicy(policyName, newDummyRetryer)
		// Register the second implementation with the same name
		RegisterRetryPolicy(policyName, newAnotherDummyRetryer)

		endpoint := &model.Endpoint{
			ID: "test-endpoint-4",
			LLMMeta: &model.LLMMeta{
				RetryPolicy: model.RetryPolicy{
					Name: policyName,
				},
			},
		}

		policy, err := GetRetryPolicy(endpoint)
		if err != nil {
			t.Fatalf("GetRetryPolicy() failed after overwrite: %v", err)
		}

		// Check that the retrieved policy is from the *second* factory
		if _, ok := policy.(*anotherDummyRetryPolicy); !ok {
			t.Fatalf("Expected policy to be of the overwritten type *anotherDummyRetryPolicy, but got %T", policy)
		}
	})
}
