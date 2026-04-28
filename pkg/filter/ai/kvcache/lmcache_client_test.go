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
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLMCacheClientRetrySuccess(t *testing.T) {
	var attempts int64
	restyClient := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/lookup", r.URL.Path)
		attempt := atomic.AddInt64(&attempts, 1)
		if attempt < 3 {
			return newHTTPResponse(http.StatusBadGateway, "temporary failure"), nil
		}
		return newHTTPResponse(http.StatusOK, `{"event_id":"evt-1","layout_info":{"node-a":{"0":"mem0","1":6}}}`), nil
	})

	client := NewLMCacheClient("http://lmcache.local", restyClient, RetryConfig{
		MaxAttempts: 3,
		BaseBackoff: time.Millisecond,
		MaxBackoff:  2 * time.Millisecond,
	}, nil)

	resp, err := client.Lookup(context.Background(), &LookupRequest{Tokens: []int{1, 2}})
	require.NoError(t, err)
	assert.Equal(t, int64(3), atomic.LoadInt64(&attempts))
	assert.Equal(t, "evt-1", resp.EventID)
	assert.Equal(t, 6, resp.LayoutInfo["node-a"].TokenCount)
}

func TestLMCacheClientContextCancelDuringBackoff(t *testing.T) {
	restyClient := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		return newHTTPResponse(http.StatusBadGateway, "retry"), nil
	})

	client := NewLMCacheClient("http://lmcache.local", restyClient, RetryConfig{
		MaxAttempts: 3,
		BaseBackoff: 100 * time.Millisecond,
		MaxBackoff:  100 * time.Millisecond,
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := client.Pin(ctx, &PinRequest{Tokens: []int{1}})
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestLMCacheClientCircuitBreakerOpen(t *testing.T) {
	var called int64
	restyClient := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt64(&called, 1)
		return newHTTPResponse(http.StatusOK, `{"event_id":"evt-2","num_tokens":1}`), nil
	})

	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  time.Minute,
		HalfOpenMaxCalls: 1,
	})
	cb.state = CircuitOpen
	cb.lastFailTime = time.Now()

	client := NewLMCacheClient("http://lmcache.local", restyClient, RetryConfig{
		MaxAttempts: 2,
		BaseBackoff: time.Millisecond,
		MaxBackoff:  time.Millisecond,
	}, cb)
	_, err := client.Evict(context.Background(), &EvictRequest{Tokens: []int{1}})
	assert.ErrorIs(t, err, ErrCircuitBreakerOpen)
	assert.Equal(t, int64(0), atomic.LoadInt64(&called))
}
