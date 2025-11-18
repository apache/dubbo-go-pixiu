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
	"fmt"
	"strconv"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	retryTimesKey     = "times"
	defaultRetryTimes = 2
)

func init() {
	retry.RegisterRetryPolicy(model.RetryerCountBased, newCountBasedRetry)
}

type CountBasedRetry struct {
	MaxAttempts uint
	retryTimes  uint
}

func (r *CountBasedRetry) Attempt() bool {
	if r.retryTimes < r.MaxAttempts {
		r.retryTimes++
		return true
	}
	return false
}

func (r *CountBasedRetry) Reset() {
	r.retryTimes = 0
}

func newCountBasedRetry(config map[string]any) (retry.RetryPolicy, error) {
	timesValue, exists := config[retryTimesKey]
	if !exists {
		timesValue = defaultRetryTimes
	}

	var times int

	switch v := timesValue.(type) {
	case int:
		times = v
	case float64:
		if v != float64(int(v)) {
			return nil, fmt.Errorf("invalid float value for '%s', must be a whole number, but got %f", retryTimesKey, v)
		}
		times = int(v)
	case string:
		parsedTimes, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("could not parse string value for '%s': %w", retryTimesKey, err)
		}
		times = parsedTimes
	default:
		return nil, fmt.Errorf("invalid type for '%s', expected a number or a numeric string, but got %T", retryTimesKey, timesValue)
	}

	// must be non-negative
	if times < 0 {
		return nil, fmt.Errorf("value for '%s' cannot be negative, but got %d", retryTimesKey, times)
	}

	// Total attempts = 1 initial try plus number of retries.
	return &CountBasedRetry{MaxAttempts: uint(times) + 1}, nil
}
