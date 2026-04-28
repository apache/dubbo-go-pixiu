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
	"sync"
	"time"
)

var ErrCircuitBreakerOpen = errors.New("kvcache circuit breaker open")

type CircuitState int

const (
	CircuitClosed   CircuitState = 0
	CircuitOpen     CircuitState = 1
	CircuitHalfOpen CircuitState = 2
)

type CircuitBreaker struct {
	state         CircuitState
	failureCount  int
	lastFailTime  time.Time
	halfOpenCalls int
	config        CircuitBreakerConfig
	mutex         sync.Mutex
}

type CircuitBreakerConfig struct {
	FailureThreshold int           `yaml:"failure_threshold" json:"failure_threshold" mapstructure:"failure_threshold"`
	RecoveryTimeout  time.Duration `yaml:"recovery_timeout" json:"recovery_timeout" mapstructure:"recovery_timeout"`
	HalfOpenMaxCalls int           `yaml:"half_open_max_calls" json:"half_open_max_calls" mapstructure:"half_open_max_calls"`
}

func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{config: cfg, state: CircuitClosed}
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	if cb == nil {
		return operation()
	}
	if !cb.allow() {
		return ErrCircuitBreakerOpen
	}
	err := operation()
	cb.recordResult(err)
	return err
}

func (cb *CircuitBreaker) allow() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case CircuitOpen:
		if time.Since(cb.lastFailTime) >= cb.config.RecoveryTimeout {
			cb.state = CircuitHalfOpen
			cb.halfOpenCalls = 0
		} else {
			return false
		}
	case CircuitHalfOpen:
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			return false
		}
		cb.halfOpenCalls++
	}
	return true
}

func (cb *CircuitBreaker) recordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err == nil {
		cb.failureCount = 0
		if cb.state == CircuitHalfOpen || cb.state == CircuitOpen {
			cb.state = CircuitClosed
			cb.halfOpenCalls = 0
		}
		return
	}

	cb.failureCount++
	cb.lastFailTime = time.Now()
	if cb.failureCount >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
	}
}
