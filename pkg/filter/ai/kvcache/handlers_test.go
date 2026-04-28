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
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/go-resty/resty/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
)

func TestExtractPromptAndModel(t *testing.T) {
	testCases := []struct {
		name       string
		body       string
		wantPrompt string
		wantModel  string
		wantErr    bool
	}{
		{
			name:       "prompt field",
			body:       `{"model":"m1","prompt":"hello"}`,
			wantPrompt: "hello",
			wantModel:  "m1",
		},
		{
			name:       "prompt array",
			body:       `{"model":"m2","prompt":["a","b"]}`,
			wantPrompt: "a\nb",
			wantModel:  "m2",
		},
		{
			name:       "messages fallback",
			body:       `{"model":"m3","messages":[{"role":"user","content":"hi"},{"role":"assistant","content":"there"}]}`,
			wantPrompt: "hi\nthere",
			wantModel:  "m3",
		},
		{
			name:    "invalid json",
			body:    `{"model":`,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompt, model, err := extractPromptAndModel([]byte(tc.body))
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantPrompt, prompt)
			assert.Equal(t, tc.wantModel, model)
		})
	}
}

func TestTryRouteToCachedInstance(t *testing.T) {
	var lookupCalls int64
	restyClient := newRestyClientWithRoundTripper(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt64(&lookupCalls, 1)
		assert.Equal(t, "/lookup", r.URL.Path)
		return newHTTPResponse(http.StatusOK, `{"event_id":"evt","layout_info":{"node-a":{"0":"x","1":2},"node-b":{"0":"y","1":8}}}`), nil
	})

	tm := NewTokenManager("", resty.New(), TokenCacheConfig{
		Enabled: true,
		MaxSize: 10,
		TTL:     time.Minute,
	}, nil, time.Minute, 10)
	key := tm.cacheKey("m1", "prompt-1")
	tm.storeCache(key, []int{1, 2, 3})

	lmcacheClient := NewLMCacheClient("http://lmcache.local", restyClient, RetryConfig{
		MaxAttempts: 1,
		BaseBackoff: time.Millisecond,
		MaxBackoff:  time.Millisecond,
	}, nil)
	f := &Filter{
		cfg: &Config{
			RequestTimeout:       time.Second,
			LookupRoutingTimeout: 100 * time.Millisecond,
		},
		tokenManager:  tm,
		lmcacheClient: lmcacheClient,
	}

	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	require.NoError(t, err)
	hc := mock.GetMockHTTPContext(req)

	cacheStatus, routed := f.tryRouteToCachedInstance(hc, "m1", "prompt-1")
	require.True(t, routed)
	require.NotNil(t, cacheStatus)
	assert.Equal(t, "node-b", hc.Params[llmPreferredEndpointIDKey])
	assert.Equal(t, int64(1), atomic.LoadInt64(&lookupCalls))

	cacheStatus, routed = f.tryRouteToCachedInstance(hc, "m1", "prompt-2")
	assert.False(t, routed)
	assert.Nil(t, cacheStatus)
	assert.Equal(t, int64(1), atomic.LoadInt64(&lookupCalls))
}
