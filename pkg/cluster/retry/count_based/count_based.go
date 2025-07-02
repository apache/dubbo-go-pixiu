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

package count_based

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	retry.RegisterRetryPolicy(model.RetryerCountBased, newCountBasedRetry)
}

type CountBasedRetry struct {
	MaxAttempts uint
	currentTry  uint
}

func (r *CountBasedRetry) Attempt(err error) bool {
	if r.currentTry < r.MaxAttempts {
		r.currentTry++
		return true
	}
	return false
}

func (r *CountBasedRetry) Reset() {
	r.currentTry = 0
}

func newCountBasedRetry(config map[string]any) (retry.Retryer, error) {
	timesValue, exists := config["times"]
	if !exists {
		return nil, fmt.Errorf("'times' field is missing in retry configuration")
	}

	timesUint, ok := timesValue.(int)
	if !ok {
		return nil, fmt.Errorf("invalid type for 'retry.count_based.times', expected int but got %T", timesValue)
	}

	// Total attempts = 1 initial try plus number of retries.
	return &CountBasedRetry{MaxAttempts: uint(timesUint) + 1}, nil
}
