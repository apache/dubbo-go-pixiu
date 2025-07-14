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
	"fmt"
	"math"
	"math/rand"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"

	"github.com/mitchellh/mapstructure"
)

func init() {
	retry.RegisterRetryPolicy(model.RetryerExponentialBackoff, newExponentialBackoffRetry)
}

type ExponentialBackoffRetry struct {
	MaxAttempts     uint
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	currentTry      uint
}

type ExponentialBackoffConfig struct {
	Times           uint    `mapstructure:"times" default:"3"`
	InitialInterval string  `mapstructure:"initialInterval" default:"100ms"`
	MaxInterval     string  `mapstructure:"maxInterval" default:"5s"`
	Multiplier      float64 `mapstructure:"multiplier" default:"2.0"`
}

func (e *ExponentialBackoffRetry) Attempt(err error) bool {
	if e.currentTry >= e.MaxAttempts {
		return false
	}

	// Don't wait before the first try
	if e.currentTry > 0 {
		backoff := float64(e.InitialInterval) * math.Pow(e.Multiplier, float64(e.currentTry-1))
		cappedBackoff := time.Duration(math.Min(backoff, float64(e.MaxInterval)))
		// Add jitter to prevent thundering herd
		jitter := time.Duration(rand.Intn(100)) * time.Millisecond // NOSONAR
		time.Sleep(cappedBackoff + jitter)
	}

	e.currentTry++
	return true
}

func (e *ExponentialBackoffRetry) Reset() {
	e.currentTry = 0
}

func newExponentialBackoffRetry(config map[string]any) (retry.Retryer, error) {
	var cfg ExponentialBackoffConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode exponential backoff config: %w", err)
	}

	initial, err := time.ParseDuration(cfg.InitialInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid initialInterval: %w", err)
	}
	duration, err := time.ParseDuration(cfg.MaxInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid maxInterval: %w", err)
	}

	return &ExponentialBackoffRetry{
		MaxAttempts:     cfg.Times + 1,
		InitialInterval: initial,
		MaxInterval:     duration,
		Multiplier:      cfg.Multiplier,
	}, nil
}
