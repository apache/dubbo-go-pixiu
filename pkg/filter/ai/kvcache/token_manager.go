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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"github.com/go-resty/resty/v2"
)

const hotMapMinSweepInterval = 30 * time.Second

type TokenManager struct {
	httpClient     *resty.Client
	endpoint       string
	cache          sync.Map
	config         TokenCacheConfig
	circuitBreaker *CircuitBreaker

	cacheSize int64
	hitCount  int64
	missCount int64

	hotWindow  time.Duration
	hotMax     int
	hotMaxKeys int
	hotMu      sync.Mutex
	hotMap     map[string][]time.Time

	hotSweepInterval time.Duration
	hotLastSweep     time.Time
}

type TokenizeRequest struct {
	Model  string `json:"model,omitempty"`
	Prompt string `json:"prompt"`
}

type TokenizeResponse struct {
	Count  int   `json:"count"`
	Tokens []int `json:"tokens"`
	MaxLen int   `json:"max_model_len"`
}

type tokenCacheEntry struct {
	tokens    []int
	expiresAt time.Time
}

func NewTokenManager(endpoint string, httpClient *resty.Client, cfg TokenCacheConfig, cb *CircuitBreaker, hotWindow time.Duration, hotMax int, hotMaxKeys ...int) *TokenManager {
	maxKeys := 0
	if len(hotMaxKeys) > 0 {
		maxKeys = hotMaxKeys[0]
	}
	return &TokenManager{
		httpClient:       httpClient,
		endpoint:         endpoint,
		config:           cfg,
		circuitBreaker:   cb,
		hotWindow:        hotWindow,
		hotMax:           hotMax,
		hotMaxKeys:       maxKeys,
		hotMap:           make(map[string][]time.Time),
		hotSweepInterval: computeHotSweepInterval(hotWindow),
	}
}

func (tm *TokenManager) GetTokens(ctx context.Context, model string, prompt string, rawBody []byte) ([]int, error) {
	cacheKey := tm.cacheKey(model, prompt)
	if tm.config.Enabled {
		if tokens, ok := tm.loadCache(cacheKey); ok {
			atomic.AddInt64(&tm.hitCount, 1)
			return tokens, nil
		}
		atomic.AddInt64(&tm.missCount, 1)
	}

	tokens, err := tm.tokenize(ctx, model, prompt, rawBody)
	if err != nil {
		return nil, err
	}

	if tm.config.Enabled {
		tm.storeCache(cacheKey, tokens)
	}
	return tokens, nil
}

func (tm *TokenManager) GetCachedTokens(model string, prompt string) ([]int, bool) {
	if !tm.config.Enabled {
		return nil, false
	}
	cacheKey := tm.cacheKey(model, prompt)
	tokens, ok := tm.loadCache(cacheKey)
	if ok {
		atomic.AddInt64(&tm.hitCount, 1)
	} else {
		atomic.AddInt64(&tm.missCount, 1)
	}
	return tokens, ok
}

func (tm *TokenManager) InvalidateCache(model string, prompt string) {
	cacheKey := tm.cacheKey(model, prompt)
	tm.deleteCache(cacheKey)
}

func (tm *TokenManager) GetCacheStats() CacheStats {
	size := atomic.LoadInt64(&tm.cacheSize)
	hit := atomic.LoadInt64(&tm.hitCount)
	miss := atomic.LoadInt64(&tm.missCount)
	total := hit + miss
	var hitRate float64
	if total > 0 {
		hitRate = float64(hit) / float64(total)
	}
	return CacheStats{
		Size:      int(size),
		HitRate:   hitRate,
		HitCount:  hit,
		MissCount: miss,
	}
}

func (tm *TokenManager) tokenize(ctx context.Context, model string, prompt string, rawBody []byte) ([]int, error) {
	var tokens []int
	err := tm.execute(ctx, func() error {
		body, err := tm.buildTokenizeBody(model, prompt, rawBody)
		if err != nil {
			return err
		}
		resp, err := tm.doTokenizeRequest(ctx, body)
		if err != nil {
			return err
		}
		tokens = resp.Tokens
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (tm *TokenManager) buildTokenizeBody(model string, prompt string, rawBody []byte) (any, error) {
	if len(rawBody) > 0 {
		return rawBody, nil
	}
	return TokenizeRequest{Model: model, Prompt: prompt}, nil
}

func (tm *TokenManager) doTokenizeRequest(ctx context.Context, body any) (*TokenizeResponse, error) {
	tokenizeURL := strings.TrimRight(tm.endpoint, "/") + "/tokenize"
	resp, err := tm.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(tokenizeURL)
	if err != nil {
		return nil, fmt.Errorf("[kvcache] call tokenize: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("[kvcache] tokenize status %d: %s", resp.StatusCode(), strings.TrimSpace(string(resp.Body())))
	}
	var tokenResp TokenizeResponse
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return nil, fmt.Errorf("[kvcache] decode tokenize response: %w", err)
	}
	return &tokenResp, nil
}

func (tm *TokenManager) RecordHot(model string, prompt string) {
	if tm == nil || tm.hotWindow <= 0 || model == "" || prompt == "" {
		return
	}
	now := time.Now()
	key := tm.cacheKey(model, prompt)
	tm.hotMu.Lock()
	defer tm.hotMu.Unlock()
	tm.maybeSweepHotMapLocked(now)
	entries := tm.hotMap[key]
	entries = append(entries, now)
	entries = trimHotWindow(entries, now, tm.hotWindow)
	if tm.hotMax > 0 && len(entries) > tm.hotMax {
		entries = entries[len(entries)-tm.hotMax:]
	}
	if len(entries) == 0 {
		delete(tm.hotMap, key)
		return
	}
	tm.hotMap[key] = entries
	tm.enforceHotMapLimitLocked(now)
}

func (tm *TokenManager) IsHot(model string, prompt string, threshold int) bool {
	if tm == nil || tm.hotWindow <= 0 || threshold <= 0 || model == "" || prompt == "" {
		return false
	}
	now := time.Now()
	key := tm.cacheKey(model, prompt)
	tm.hotMu.Lock()
	defer tm.hotMu.Unlock()
	entries := tm.hotMap[key]
	if len(entries) == 0 {
		return false
	}
	entries = trimHotWindow(entries, now, tm.hotWindow)
	if tm.hotMax > 0 && len(entries) > tm.hotMax {
		entries = entries[len(entries)-tm.hotMax:]
	}
	if len(entries) == 0 {
		delete(tm.hotMap, key)
		return false
	}
	tm.hotMap[key] = entries
	return len(entries) >= threshold
}

func computeHotSweepInterval(hotWindow time.Duration) time.Duration {
	if hotWindow <= 0 {
		return 0
	}
	interval := hotWindow / 2
	if interval <= 0 {
		interval = hotWindow
	}
	if hotWindow > hotMapMinSweepInterval && interval < hotMapMinSweepInterval {
		interval = hotMapMinSweepInterval
	}
	return interval
}

func (tm *TokenManager) maybeSweepHotMapLocked(now time.Time) {
	if tm.hotSweepInterval <= 0 {
		return
	}
	if !tm.hotLastSweep.IsZero() && now.Sub(tm.hotLastSweep) < tm.hotSweepInterval {
		return
	}
	for key, entries := range tm.hotMap {
		entries = trimHotWindow(entries, now, tm.hotWindow)
		if tm.hotMax > 0 && len(entries) > tm.hotMax {
			entries = entries[len(entries)-tm.hotMax:]
		}
		if len(entries) == 0 {
			delete(tm.hotMap, key)
			continue
		}
		tm.hotMap[key] = entries
	}
	tm.hotLastSweep = now
}

func (tm *TokenManager) enforceHotMapLimitLocked(now time.Time) {
	if tm.hotMaxKeys <= 0 || len(tm.hotMap) <= tm.hotMaxKeys {
		return
	}

	// First, aggressively drop expired keys when we are already over the key cap.
	for key, entries := range tm.hotMap {
		entries = trimHotWindow(entries, now, tm.hotWindow)
		if tm.hotMax > 0 && len(entries) > tm.hotMax {
			entries = entries[len(entries)-tm.hotMax:]
		}
		if len(entries) == 0 {
			delete(tm.hotMap, key)
			continue
		}
		tm.hotMap[key] = entries
	}
	if len(tm.hotMap) <= tm.hotMaxKeys {
		return
	}

	type hotKeyLastSeen struct {
		key      string
		lastSeen time.Time
	}
	candidates := make([]hotKeyLastSeen, 0, len(tm.hotMap))
	for key, entries := range tm.hotMap {
		if len(entries) == 0 {
			delete(tm.hotMap, key)
			continue
		}
		candidates = append(candidates, hotKeyLastSeen{
			key:      key,
			lastSeen: entries[len(entries)-1],
		})
	}
	if len(candidates) <= tm.hotMaxKeys {
		return
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].lastSeen.Before(candidates[j].lastSeen)
	})
	excess := len(candidates) - tm.hotMaxKeys
	for i := 0; i < excess; i++ {
		delete(tm.hotMap, candidates[i].key)
	}
}

func trimHotWindow(entries []time.Time, now time.Time, window time.Duration) []time.Time {
	if window <= 0 || len(entries) == 0 {
		return entries
	}
	cutoff := now.Add(-window)
	idx := 0
	for idx < len(entries) && entries[idx].Before(cutoff) {
		idx++
	}
	if idx == 0 {
		return entries
	}
	return append([]time.Time(nil), entries[idx:]...)
}

func (tm *TokenManager) execute(ctx context.Context, operation func() error) error {
	if tm.circuitBreaker == nil {
		return operation()
	}
	return tm.circuitBreaker.Execute(operation)
}

func (tm *TokenManager) cacheKey(model string, prompt string) string {
	sum := sha256.Sum256([]byte(model + "\x00" + prompt))
	return hex.EncodeToString(sum[:])
}

func (tm *TokenManager) deleteCache(key string) {
	if _, loaded := tm.cache.LoadAndDelete(key); loaded {
		atomic.AddInt64(&tm.cacheSize, -1)
	}
}

func (tm *TokenManager) loadCache(key string) ([]int, bool) {
	entryAny, ok := tm.cache.Load(key)
	if !ok {
		return nil, false
	}
	entry, ok := entryAny.(*tokenCacheEntry)
	if !ok {
		tm.deleteCache(key)
		return nil, false
	}
	if tm.config.TTL > 0 && time.Now().After(entry.expiresAt) {
		tm.deleteCache(key)
		return nil, false
	}
	return entry.tokens, true
}

func (tm *TokenManager) storeCache(key string, tokens []int) {
	if tm.config.MaxSize > 0 {
		for atomic.LoadInt64(&tm.cacheSize) >= int64(tm.config.MaxSize) {
			if !tm.evictOne() {
				break
			}
		}
	}
	entry := &tokenCacheEntry{
		tokens:    tokens,
		expiresAt: time.Now().Add(tm.config.TTL),
	}
	if _, loaded := tm.cache.LoadOrStore(key, entry); !loaded {
		atomic.AddInt64(&tm.cacheSize, 1)
	}
}

func (tm *TokenManager) evictOne() bool {
	evicted := false
	tm.cache.Range(func(key, value any) bool {
		tm.deleteCache(key.(string))
		evicted = true
		return false
	})
	return evicted
}
