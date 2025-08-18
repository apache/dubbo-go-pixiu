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

package retry

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	// RetryPolicy defines the interface for retry logic.
	RetryPolicy interface {
		// Attempt checks if a retry should be performed and potentially waits.
		// It returns true if the request should be attempted, false otherwise.
		Attempt() bool
		// Reset re-initializes the policy's state.
		Reset()
	}

	// RetryPolicyFactory creates an instance of a RetryPolicy from a config map.
	RetryPolicyFactory func(config map[string]any) (RetryPolicy, error)
)

// retryPolicyRegistry holds all available retry policy implementations.
var retryPolicyRegistry = make(map[model.RetryType]RetryPolicyFactory)

// RegisterRetryPolicy makes a retry policy implementation available by name.
// This function is called from init() functions in files that define a policy.
func RegisterRetryPolicy(name model.RetryType, factory RetryPolicyFactory) {
	name = name.ToLower()
	if _, exists := retryPolicyRegistry[name]; exists {
		logger.Warnf("retry policy type '%s' is being overwritten", name)
	}
	retryPolicyRegistry[name] = factory
}

// GetRetryPolicy dynamically creates a RetryPolicy based on endpoint metadata.
func GetRetryPolicy(endpoint *model.Endpoint) (RetryPolicy, error) {
	retryPolicy := endpoint.LLMMeta.RetryPolicy
	factory, exists := retryPolicyRegistry[retryPolicy.Name.ToLower()]
	if !exists {
		return nil, fmt.Errorf("unknown retry policy type '%s' specified for endpoint %s", retryPolicy.Name, endpoint.ID)
	}

	return factory(retryPolicy.Config)
}
