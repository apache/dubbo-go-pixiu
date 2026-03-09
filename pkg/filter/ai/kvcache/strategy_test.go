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
	"net/http"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/go-resty/resty/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheStrategyMakeDecision(t *testing.T) {
	tm := NewTokenManager("", resty.New(), TokenCacheConfig{}, nil, time.Minute, 10)
	tm.RecordHot("m1", "p1")
	tm.RecordHot("m1", "p1")

	cs := NewCacheStrategy(CacheStrategyConfig{
		EnableCompression:   true,
		EnableEviction:      true,
		EnablePinning:       true,
		LoadThreshold:       0,
		MemoryThreshold:     0.0000001,
		HotContentThreshold: 2,
	}, &LMCacheClient{}, tm)

	decision := cs.MakeDecision(context.Background(), nil, "m1", "p1")
	assert.True(t, decision.ShouldPin)
	assert.True(t, decision.ShouldEvict)
	assert.Equal(t, "hot_content", decision.Reason)
}

func TestCacheStrategyExecuteDecision(t *testing.T) {
	var (
		mu    sync.Mutex
		paths []string
	)
	restyClient := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		return newHTTPResponse(http.StatusOK, `{"event_id":"ok","num_tokens":2}`), nil
	})

	client := NewLMCacheClient("http://lmcache.local", restyClient, RetryConfig{
		MaxAttempts: 1,
		BaseBackoff: time.Millisecond,
		MaxBackoff:  time.Millisecond,
	}, nil)
	cs := NewCacheStrategy(CacheStrategyConfig{
		CompressInstanceID: "compress-i",
		CompressLocation:   "loc-a",
		CompressMethod:     "zstd",
		PinInstanceID:      "pin-i",
		PinLocation:        "loc-b",
		EvictInstanceID:    "evict-i",
	}, client, nil)

	decision := &StrategyDecision{
		ShouldCompress: true,
		ShouldPin:      true,
		ShouldEvict:    true,
	}
	err := cs.ExecuteDecision(context.Background(), decision, []int{1, 2})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"/compress", "/pin", "/evict"}, paths)
}

func TestCacheStrategyExecuteDecisionStopsOnFirstError(t *testing.T) {
	var (
		mu    sync.Mutex
		paths []string
	)
	restyClient := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		if r.URL.Path == "/compress" {
			return newHTTPResponse(http.StatusInternalServerError, "compress failed"), nil
		}
		return newHTTPResponse(http.StatusOK, `{"event_id":"ok","num_tokens":2}`), nil
	})

	client := NewLMCacheClient("http://lmcache.local", restyClient, RetryConfig{
		MaxAttempts: 1,
		BaseBackoff: time.Millisecond,
		MaxBackoff:  time.Millisecond,
	}, nil)
	cs := NewCacheStrategy(CacheStrategyConfig{
		CompressMethod: "zstd",
	}, client, nil)

	decision := &StrategyDecision{
		ShouldCompress: true,
		ShouldPin:      true,
		ShouldEvict:    true,
	}
	err := cs.ExecuteDecision(context.Background(), decision, []int{1, 2})
	assert.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"/compress"}, paths)
}
