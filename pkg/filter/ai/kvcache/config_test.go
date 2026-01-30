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
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateDisabled(t *testing.T) {
	cfg := &Config{
		Enabled: false,
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_ValidateMissingVLLMEndpoint(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		LMCacheEndpoint: "http://lmcache:8080",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vllm_endpoint is required")
}

func TestConfig_ValidateMissingLMCacheEndpoint(t *testing.T) {
	cfg := &Config{
		Enabled:      true,
		VLLMEndpoint: "http://vllm:8000",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lmcache_endpoint is required")
}

func TestConfig_ValidateNegativeMaxSize(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		TokenCache: TokenCacheConfig{
			MaxSize: -1,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max_size must be >= 0")
}

func TestConfig_ValidateMemoryThresholdTooLow(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold: -0.1,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "memory_threshold must be between 0 and 1")
}

func TestConfig_ValidateMemoryThresholdTooHigh(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold: 1.5,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "memory_threshold must be between 0 and 1")
}

func TestConfig_ValidateLoadThresholdInvalid(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold: 0.8,
			LoadThreshold:   -0.5,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "load_threshold must be between 0 and 1")
}

func TestConfig_ValidateHotContentThresholdNegative(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold:      0.8,
			LoadThreshold:        0.7,
			HotContentThreshold:  -10,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hot_content_threshold must be >= 0")
}

func TestConfig_ValidateRetryMaxAttemptsNegative(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold:     0.8,
			LoadThreshold:       0.7,
			HotContentThreshold: 100,
		},
		Retry: RetryConfig{
			MaxAttempts: -1,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max_attempts must be >= 0")
}

func TestConfig_ValidateRetryBackoffNegative(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold:     0.8,
			LoadThreshold:       0.7,
			HotContentThreshold: 100,
		},
		Retry: RetryConfig{
			MaxAttempts: 3,
			BaseBackoff: -100 * time.Millisecond,
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backoff durations must be >= 0")
}

func TestConfig_ValidateSuccess(t *testing.T) {
	cfg := &Config{
		Enabled:         true,
		VLLMEndpoint:    "http://vllm:8000",
		LMCacheEndpoint: "http://lmcache:8080",
		TokenCache: TokenCacheConfig{
			MaxSize: 1000,
		},
		CacheStrategy: CacheStrategyConfig{
			MemoryThreshold:     0.8,
			LoadThreshold:       0.7,
			HotContentThreshold: 100,
		},
		Retry: RetryConfig{
			MaxAttempts: 3,
			BaseBackoff: 100 * time.Millisecond,
			MaxBackoff:  2 * time.Second,
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_ApplyDefaults(t *testing.T) {
	cfg := &Config{}

	cfg.ApplyDefaults()

	assert.Equal(t, 2*time.Second, cfg.RequestTimeout)
	assert.Equal(t, 50*time.Millisecond, cfg.LookupRoutingTimeout)
	assert.Equal(t, 3, cfg.Retry.MaxAttempts)
	assert.Equal(t, 100*time.Millisecond, cfg.Retry.BaseBackoff)
	assert.Equal(t, 2*time.Second, cfg.Retry.MaxBackoff)
}

func TestConfig_ApplyDefaultsDoesNotOverride(t *testing.T) {
	cfg := &Config{
		RequestTimeout:       5 * time.Second,
		LookupRoutingTimeout: 100 * time.Millisecond,
		Retry: RetryConfig{
			MaxAttempts: 5,
			BaseBackoff: 200 * time.Millisecond,
			MaxBackoff:  5 * time.Second,
		},
	}

	cfg.ApplyDefaults()

	assert.Equal(t, 5*time.Second, cfg.RequestTimeout)
	assert.Equal(t, 100*time.Millisecond, cfg.LookupRoutingTimeout)
	assert.Equal(t, 5, cfg.Retry.MaxAttempts)
	assert.Equal(t, 200*time.Millisecond, cfg.Retry.BaseBackoff)
	assert.Equal(t, 5*time.Second, cfg.Retry.MaxBackoff)
}
