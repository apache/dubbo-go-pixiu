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
)

import (
	"github.com/stretchr/testify/assert"
)

func TestTokenManager_CacheKey(t *testing.T) {
	tm := &TokenManager{}

	key1 := tm.cacheKey("model1", "prompt1")
	key2 := tm.cacheKey("model1", "prompt1")
	key3 := tm.cacheKey("model2", "prompt1")
	key4 := tm.cacheKey("model1", "prompt2")

	// Same model and prompt should produce same key
	assert.Equal(t, key1, key2)

	// Different model should produce different key
	assert.NotEqual(t, key1, key3)

	// Different prompt should produce different key
	assert.NotEqual(t, key1, key4)
}

func TestTokenManager_StoreAndLoadCache(t *testing.T) {
	tm := &TokenManager{}

	tokens := []int{1, 2, 3, 4, 5}
	cacheKey := "test:key"

	// Store tokens
	tm.storeCache(cacheKey, tokens)

	// Load tokens
	loaded, ok := tm.loadCache(cacheKey)
	assert.True(t, ok)
	assert.Equal(t, tokens, loaded)
}

func TestTokenManager_LoadCacheMiss(t *testing.T) {
	tm := &TokenManager{}

	// Try to load non-existent key
	loaded, ok := tm.loadCache("nonexistent:key")
	assert.False(t, ok)
	assert.Nil(t, loaded)
}

func TestTokenManager_GetCachedTokensDisabled(t *testing.T) {
	tm := &TokenManager{
		config: TokenCacheConfig{
			Enabled: false,
		},
	}

	tokens, ok := tm.GetCachedTokens("model", "prompt")
	assert.False(t, ok)
	assert.Nil(t, tokens)
}

func TestTokenManager_GetCachedTokensHit(t *testing.T) {
	tm := &TokenManager{
		config: TokenCacheConfig{
			Enabled: true,
		},
	}

	// Store tokens first
	expectedTokens := []int{10, 20, 30}
	cacheKey := tm.cacheKey("model1", "test prompt")
	tm.storeCache(cacheKey, expectedTokens)

	// Get cached tokens
	tokens, ok := tm.GetCachedTokens("model1", "test prompt")
	assert.True(t, ok)
	assert.Equal(t, expectedTokens, tokens)
	assert.Equal(t, int64(1), tm.hitCount)
}

func TestTokenManager_GetCachedTokensMiss(t *testing.T) {
	tm := &TokenManager{
		config: TokenCacheConfig{
			Enabled: true,
		},
	}

	// Try to get non-existent tokens
	tokens, ok := tm.GetCachedTokens("model1", "nonexistent prompt")
	assert.False(t, ok)
	assert.Nil(t, tokens)
	assert.Equal(t, int64(1), tm.missCount)
}

func TestTokenManager_InvalidateCache(t *testing.T) {
	tm := &TokenManager{
		config: TokenCacheConfig{
			Enabled: true,
		},
	}

	// Store tokens
	tokens := []int{1, 2, 3}
	cacheKey := tm.cacheKey("model1", "prompt1")
	tm.storeCache(cacheKey, tokens)

	// Verify stored
	loaded, ok := tm.loadCache(cacheKey)
	assert.True(t, ok)
	assert.Equal(t, tokens, loaded)

	// Invalidate
	tm.InvalidateCache("model1", "prompt1")

	// Verify removed
	loaded, ok = tm.loadCache(cacheKey)
	assert.False(t, ok)
	assert.Nil(t, loaded)
}

func TestTokenManager_GetCacheStats(t *testing.T) {
	tm := &TokenManager{
		config: TokenCacheConfig{
			Enabled: true,
		},
	}

	// Simulate some hits and misses
	cacheKey := tm.cacheKey("model", "prompt")
	tm.storeCache(cacheKey, []int{1, 2, 3})

	// Hit
	tm.GetCachedTokens("model", "prompt")
	// Miss
	tm.GetCachedTokens("model", "other")
	// Hit
	tm.GetCachedTokens("model", "prompt")

	stats := tm.GetCacheStats()
	assert.Equal(t, int64(2), stats.HitCount)
	assert.Equal(t, int64(1), stats.MissCount)
	assert.InDelta(t, 0.666, stats.HitRate, 0.01)
}

func TestTokenManager_GetCacheStatsNoRequests(t *testing.T) {
	tm := &TokenManager{
		config: TokenCacheConfig{
			Enabled: true,
		},
	}

	stats := tm.GetCacheStats()
	assert.Equal(t, int64(0), stats.HitCount)
	assert.Equal(t, int64(0), stats.MissCount)
	assert.Equal(t, 0.0, stats.HitRate)
}
