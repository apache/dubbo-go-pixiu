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
	"fmt"
	"strings"
	"time"
)

const (
	minRatio = 0.0
	maxRatio = 1.0

	defaultRequestTimeout       = 2 * time.Second
	defaultLookupRoutingTimeout = 50 * time.Millisecond
	defaultHotWindow            = 5 * time.Minute
	defaultHotMaxRecords        = 300
	defaultHotMaxKeys           = 0
	defaultRetryMaxAttempts     = 3
	defaultRetryBaseBackoff     = 100 * time.Millisecond
	defaultRetryMaxBackoff      = 2 * time.Second
	defaultCBFailureThreshold   = 5
	defaultCBRecoveryTimeout    = 10 * time.Second
	defaultCBHalfOpenMaxCalls   = 2
	defaultCompressMethod       = "zstd"
)

type Config struct {
	Enabled              bool                 `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	VLLMEndpoint         string               `yaml:"vllm_endpoint" json:"vllm_endpoint" mapstructure:"vllm_endpoint"`
	LMCacheEndpoint      string               `yaml:"lmcache_endpoint" json:"lmcache_endpoint" mapstructure:"lmcache_endpoint"`
	DefaultModel         string               `yaml:"default_model" json:"default_model" mapstructure:"default_model"`
	RequestTimeout       time.Duration        `yaml:"request_timeout" json:"request_timeout" mapstructure:"request_timeout"`
	LookupRoutingTimeout time.Duration        `yaml:"lookup_routing_timeout" json:"lookup_routing_timeout" mapstructure:"lookup_routing_timeout"`
	HotWindow            time.Duration        `yaml:"hot_window" json:"hot_window" mapstructure:"hot_window"`
	HotMaxRecords        int                  `yaml:"hot_max_records" json:"hot_max_records" mapstructure:"hot_max_records"`
	HotMaxKeys           int                  `yaml:"hot_max_keys" json:"hot_max_keys" mapstructure:"hot_max_keys"`
	MaxIdleConns         int                  `yaml:"max_idle_conns" json:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost  int                  `yaml:"max_idle_conns_per_host" json:"max_idle_conns_per_host" mapstructure:"max_idle_conns_per_host"`
	MaxConnsPerHost      int                  `yaml:"max_conns_per_host" json:"max_conns_per_host" mapstructure:"max_conns_per_host"`
	TokenCache           TokenCacheConfig     `yaml:"token_cache" json:"token_cache" mapstructure:"token_cache"`
	CacheStrategy        CacheStrategyConfig  `yaml:"cache_strategy" json:"cache_strategy" mapstructure:"cache_strategy"`
	CircuitBreaker       CircuitBreakerConfig `yaml:"circuit_breaker" json:"circuit_breaker" mapstructure:"circuit_breaker"`
	Retry                RetryConfig          `yaml:"retry" json:"retry" mapstructure:"retry"`
}

type TokenCacheConfig struct {
	MaxSize int           `yaml:"max_size" json:"max_size" mapstructure:"max_size"`
	TTL     time.Duration `yaml:"ttl" json:"ttl" mapstructure:"ttl"`
	Enabled bool          `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
}

type CacheStrategyConfig struct {
	EnableCompression   bool    `yaml:"enable_compression" json:"enable_compression" mapstructure:"enable_compression"`
	EnablePinning       bool    `yaml:"enable_pinning" json:"enable_pinning" mapstructure:"enable_pinning"`
	EnableEviction      bool    `yaml:"enable_eviction" json:"enable_eviction" mapstructure:"enable_eviction"`
	MemoryThreshold     float64 `yaml:"memory_threshold" json:"memory_threshold" mapstructure:"memory_threshold"`
	HotContentThreshold int     `yaml:"hot_content_threshold" json:"hot_content_threshold" mapstructure:"hot_content_threshold"`
	LoadThreshold       float64 `yaml:"load_threshold" json:"load_threshold" mapstructure:"load_threshold"`
	PinInstanceID       string  `yaml:"pin_instance_id" json:"pin_instance_id" mapstructure:"pin_instance_id"`
	PinLocation         string  `yaml:"pin_location" json:"pin_location" mapstructure:"pin_location"`
	CompressInstanceID  string  `yaml:"compress_instance_id" json:"compress_instance_id" mapstructure:"compress_instance_id"`
	CompressLocation    string  `yaml:"compress_location" json:"compress_location" mapstructure:"compress_location"`
	CompressMethod      string  `yaml:"compress_method" json:"compress_method" mapstructure:"compress_method"`
	EvictInstanceID     string  `yaml:"evict_instance_id" json:"evict_instance_id" mapstructure:"evict_instance_id"`
}

type RetryConfig struct {
	MaxAttempts int           `yaml:"max_attempts" json:"max_attempts" mapstructure:"max_attempts"`
	BaseBackoff time.Duration `yaml:"base_backoff" json:"base_backoff" mapstructure:"base_backoff"`
	MaxBackoff  time.Duration `yaml:"max_backoff" json:"max_backoff" mapstructure:"max_backoff"`
}

func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if strings.TrimSpace(c.VLLMEndpoint) == "" {
		return fmt.Errorf("[kvcache] vllm_endpoint is required when enabled")
	}
	if strings.TrimSpace(c.LMCacheEndpoint) == "" {
		return fmt.Errorf("[kvcache] lmcache_endpoint is required when enabled")
	}
	if c.TokenCache.MaxSize < 0 {
		return fmt.Errorf("[kvcache] token_cache.max_size must be >= 0")
	}
	if c.CacheStrategy.MemoryThreshold < minRatio || c.CacheStrategy.MemoryThreshold > maxRatio {
		return fmt.Errorf("[kvcache] cache_strategy.memory_threshold must be between 0 and 1")
	}
	if c.CacheStrategy.LoadThreshold < minRatio || c.CacheStrategy.LoadThreshold > maxRatio {
		return fmt.Errorf("[kvcache] cache_strategy.load_threshold must be between 0 and 1")
	}
	if c.CacheStrategy.HotContentThreshold < 0 {
		return fmt.Errorf("[kvcache] cache_strategy.hot_content_threshold must be >= 0")
	}
	if c.Retry.MaxAttempts < 0 {
		return fmt.Errorf("[kvcache] retry.max_attempts must be >= 0")
	}
	if c.Retry.BaseBackoff < 0 || c.Retry.MaxBackoff < 0 {
		return fmt.Errorf("[kvcache] retry backoff durations must be >= 0")
	}
	if c.HotWindow < 0 {
		return fmt.Errorf("[kvcache] hot_window must be >= 0")
	}
	if c.HotMaxRecords < 0 {
		return fmt.Errorf("[kvcache] hot_max_records must be >= 0")
	}
	if c.HotMaxKeys < 0 {
		return fmt.Errorf("[kvcache] hot_max_keys must be >= 0")
	}
	return nil
}

func (c *Config) ApplyDefaults() {
	if c.RequestTimeout <= 0 {
		c.RequestTimeout = defaultRequestTimeout
	}
	if c.LookupRoutingTimeout <= 0 {
		c.LookupRoutingTimeout = defaultLookupRoutingTimeout
	}
	if c.HotWindow <= 0 {
		c.HotWindow = defaultHotWindow
	}
	if c.HotMaxRecords <= 0 {
		c.HotMaxRecords = defaultHotMaxRecords
	}
	if c.HotMaxKeys < 0 {
		c.HotMaxKeys = defaultHotMaxKeys
	}
	if c.Retry.MaxAttempts <= 0 {
		c.Retry.MaxAttempts = defaultRetryMaxAttempts
	}
	if c.Retry.BaseBackoff <= 0 {
		c.Retry.BaseBackoff = defaultRetryBaseBackoff
	}
	if c.Retry.MaxBackoff <= 0 {
		c.Retry.MaxBackoff = defaultRetryMaxBackoff
	}
	if c.CircuitBreaker.FailureThreshold <= 0 {
		c.CircuitBreaker.FailureThreshold = defaultCBFailureThreshold
	}
	if c.CircuitBreaker.RecoveryTimeout <= 0 {
		c.CircuitBreaker.RecoveryTimeout = defaultCBRecoveryTimeout
	}
	if c.CircuitBreaker.HalfOpenMaxCalls <= 0 {
		c.CircuitBreaker.HalfOpenMaxCalls = defaultCBHalfOpenMaxCalls
	}
	if c.CacheStrategy.CompressMethod == "" {
		c.CacheStrategy.CompressMethod = defaultCompressMethod
	}
}

func (c *Config) DeepCopy() *Config {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
}
