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

package model

import (
	"strings"
)

// RetryPolicy holds the raw configuration for a policy from the endpoint metadata.
type RetryPolicy struct {
	Name   RetryType      `mapstructure:"name" default:"NoRetry"`
	Config map[string]any `mapstructure:"config"`
}

// RetryType the retry policy enum
type RetryType string

const (
	RetryerNoRetry            RetryType = "NoRetry"
	RetryerExponentialBackoff RetryType = "ExponentialBackoff"
	RetryerCountBased         RetryType = "CountBased"
)

var RetryTypeValue = map[string]RetryType{
	"NoRetry":            RetryerNoRetry,
	"ExponentialBackoff": RetryerExponentialBackoff,
	"CountBased":         RetryerCountBased,
}

func (r RetryType) ToLower() RetryType {
	return RetryType(strings.ToLower(string(r)))
}
