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

import "time"

const (
	HeaderSessionID   = "X-LMC-Session-ID"
	HeaderCacheMode   = "X-LMC-Cache-Mode"
	HeaderModel       = "X-LMC-Model"
	HeaderCacheHit    = "X-LMC-Cache-Hit"
	HeaderCacheDetail = "X-LMC-Cache-Detail"
	HeaderPolicy      = "X-LMC-Policy"
)

const (
	ModeReuse   = "reuse"
	ModeRefresh = "refresh"
	ModeFlush   = "flush"
)

const (
	CacheHit   = "HIT"
	CacheMiss  = "MISS"
	CacheError = "ERROR"
)

const (
	ctxSessionIDKey = "kvcache_session_id"
	ctxCacheModeKey = "kvcache_cache_mode"
	ctxModelKey     = "kvcache_model"
)

type Config struct {
	Enabled           bool           `yaml:"enabled" json:"enabled"`
	Backend           string         `yaml:"backend" json:"backend"`
	EngineType        string         `yaml:"engine_type" json:"engine_type"`
	Redis             RedisConfig    `yaml:"redis" json:"redis"`
	CacheKey          CacheKeyConfig `yaml:"cache_key" json:"cache_key"`
	DefaultTTL        time.Duration  `yaml:"default_ttl" json:"default_ttl"`
	TargetModels      []string       `yaml:"target_models" json:"target_models"`
	DefaultUseCache   *bool          `yaml:"default_use_cache" json:"default_use_cache"`
	DefaultMode       string         `yaml:"default_mode" json:"default_mode"`
	DropClientHeaders *bool          `yaml:"drop_client_headers" json:"drop_client_headers"`
}

type RedisConfig struct {
	ClusterMode  bool          `yaml:"cluster_mode" json:"cluster_mode"`
	Addrs        []string      `yaml:"addrs" json:"addrs"`
	Password     string        `yaml:"password" json:"password"`
	DB           int           `yaml:"db" json:"db"`
	KeyPrefix    string        `yaml:"key_prefix" json:"key_prefix"`
	PoolSize     int           `yaml:"pool_size" json:"pool_size"`
	RedisURL     string        `yaml:"redis_url" json:"redis_url"`
	DialTimeout  time.Duration `yaml:"dial_timeout" json:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
}

type CacheKeyConfig struct {
	IncludeModel     bool `yaml:"include_model" json:"include_model"`
	IncludeSessionID bool `yaml:"include_session_id" json:"include_session_id"`
	IncludePrompt    bool `yaml:"include_prompt" json:"include_prompt"`
}

type Policy struct {
	UseCache bool
	Mode     string
	TTL      time.Duration
}

type Stats struct {
	Hits     int64
	Misses   int64
	Errors   int64
	LastMode string
}
