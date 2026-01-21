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
	"context"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/redis/go-redis/v9"
)

const defaultKeyPrefix = "pixiu:kvcache"

// RedisStore provides minimal policy and stats persistence.
type RedisStore struct {
	client     redis.UniversalClient
	keyPrefix  string
	defaultTTL time.Duration
}

func NewRedisStore(cfg RedisConfig, defaultTTL time.Duration) *RedisStore {
	if len(cfg.Addrs) == 0 {
		return nil
	}
	keyPrefix := strings.TrimSpace(cfg.KeyPrefix)
	if keyPrefix == "" {
		keyPrefix = defaultKeyPrefix
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        cfg.Addrs,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
	return &RedisStore{
		client:     client,
		keyPrefix:  keyPrefix,
		defaultTTL: defaultTTL,
	}
}

func (s *RedisStore) policyKey(sessionID string) string {
	return s.keyPrefix + ":session:" + sessionID + ":policy"
}

func (s *RedisStore) statsKey(sessionID string) string {
	return s.keyPrefix + ":session:" + sessionID + ":stats"
}

func (s *RedisStore) GetPolicy(ctx context.Context, sessionID string, fallback Policy) (Policy, error) {
	if s == nil {
		return fallback, nil
	}
	res, err := s.client.HGetAll(ctx, s.policyKey(sessionID)).Result()
	if err != nil {
		return fallback, err
	}
	if len(res) == 0 {
		_ = s.SetPolicy(ctx, sessionID, fallback)
		return fallback, nil
	}
	policy := fallback
	if v, ok := res["use_cache"]; ok {
		policy.UseCache = strings.ToLower(v) == "true"
	}
	if v, ok := res["mode"]; ok && v != "" {
		policy.Mode = v
	}
	if v, ok := res["ttl_sec"]; ok {
		if ttlSec, err := strconv.ParseInt(v, 10, 64); err == nil {
			policy.TTL = time.Duration(ttlSec) * time.Second
		}
	}
	return policy, nil
}

func (s *RedisStore) SetPolicy(ctx context.Context, sessionID string, policy Policy) error {
	if s == nil {
		return nil
	}
	key := s.policyKey(sessionID)
	fields := map[string]any{
		"use_cache": strconv.FormatBool(policy.UseCache),
		"mode":      policy.Mode,
	}
	if policy.TTL > 0 {
		fields["ttl_sec"] = strconv.FormatInt(int64(policy.TTL.Seconds()), 10)
	}
	if err := s.client.HSet(ctx, key, fields).Err(); err != nil {
		return err
	}
	if policy.TTL > 0 {
		return s.client.Expire(ctx, key, policy.TTL).Err()
	}
	if s.defaultTTL > 0 {
		return s.client.Expire(ctx, key, s.defaultTTL).Err()
	}
	return nil
}

func (s *RedisStore) UpdateStats(ctx context.Context, sessionID string, result string, mode string) error {
	if s == nil {
		return nil
	}
	key := s.statsKey(sessionID)
	pipe := s.client.Pipeline()
	switch result {
	case CacheHit:
		pipe.HIncrBy(ctx, key, "hits", 1)
	case CacheMiss:
		pipe.HIncrBy(ctx, key, "misses", 1)
	case CacheError:
		pipe.HIncrBy(ctx, key, "errors", 1)
	}
	if mode != "" {
		pipe.HSet(ctx, key, "last_mode", mode)
	}
	pipe.HSet(ctx, key, "last_update_at", strconv.FormatInt(time.Now().Unix(), 10))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	if s.defaultTTL > 0 {
		return s.client.Expire(ctx, key, s.defaultTTL).Err()
	}
	return nil
}
