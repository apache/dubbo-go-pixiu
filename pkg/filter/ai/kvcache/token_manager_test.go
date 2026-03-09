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
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenManagerGetTokensUsesCache(t *testing.T) {
	var callCount int64
	client := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/tokenize", r.URL.Path)
		atomic.AddInt64(&callCount, 1)
		return newHTTPResponse(http.StatusOK, `{"count":2,"tokens":[11,22],"max_model_len":4096}`), nil
	})

	tm := NewTokenManager("http://tokenizer.local", client, TokenCacheConfig{
		Enabled: true,
		MaxSize: 10,
		TTL:     time.Minute,
	}, nil, time.Minute, 10)

	tokens1, err := tm.GetTokens(context.Background(), "m1", "hello", nil)
	require.NoError(t, err)
	assert.Equal(t, []int{11, 22}, tokens1)

	tokens2, err := tm.GetTokens(context.Background(), "m1", "hello", nil)
	require.NoError(t, err)
	assert.Equal(t, []int{11, 22}, tokens2)
	assert.Equal(t, int64(1), atomic.LoadInt64(&callCount))

	stats := tm.GetCacheStats()
	assert.Equal(t, 1, stats.Size)
	assert.Equal(t, int64(1), stats.HitCount)
	assert.Equal(t, int64(1), stats.MissCount)
	assert.Equal(t, 0.5, stats.HitRate)
}

func TestTokenManagerRawBodyAndErrorHandling(t *testing.T) {
	t.Run("raw body is passed through", func(t *testing.T) {
		var gotBody []byte
		client := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
			var err error
			gotBody, err = io.ReadAll(r.Body)
			require.NoError(t, err)
			return newHTTPResponse(http.StatusOK, `{"count":1,"tokens":[7],"max_model_len":4096}`), nil
		})

		tm := NewTokenManager("http://tokenizer.local", client, TokenCacheConfig{Enabled: false}, nil, time.Minute, 10)
		rawBody := []byte(`{"model":"raw-model","prompt":"raw-prompt"}`)
		tokens, err := tm.GetTokens(context.Background(), "ignored", "ignored", rawBody)
		require.NoError(t, err)
		assert.Equal(t, []int{7}, tokens)
		assert.JSONEq(t, string(rawBody), string(gotBody))
	})

	t.Run("status code error does not populate cache", func(t *testing.T) {
		client := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
			return newHTTPResponse(http.StatusInternalServerError, "boom"), nil
		})

		tm := NewTokenManager("http://tokenizer.local", client, TokenCacheConfig{
			Enabled: true,
			MaxSize: 10,
			TTL:     time.Minute,
		}, nil, time.Minute, 10)

		tokens, err := tm.GetTokens(context.Background(), "m1", "p1", nil)
		assert.Nil(t, tokens)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tokenize status 500")
		assert.Equal(t, 0, tm.GetCacheStats().Size)
	})
}
