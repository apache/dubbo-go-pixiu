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
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

import (
	"github.com/go-resty/resty/v2"
)

type LMCacheClient struct {
	httpClient     *resty.Client
	baseURL        string
	retry          RetryConfig
	circuitBreaker *CircuitBreaker
}

type PinRequest struct {
	Tokens     []int  `json:"tokens"`
	InstanceID string `json:"instance_id"`
	Location   string `json:"location"`
}

type LookupRequest struct {
	Tokens []int `json:"tokens"`
}

type CompressRequest struct {
	Tokens     []int  `json:"tokens"`
	InstanceID string `json:"instance_id"`
	Location   string `json:"location"`
	Method     string `json:"method"`
}

type EvictRequest struct {
	Tokens     []int  `json:"tokens"`
	InstanceID string `json:"instance_id"`
}

type PinResponse struct {
	EventID   string `json:"event_id"`
	NumTokens int    `json:"num_tokens"`
}

type CompressResponse struct {
	EventID   string `json:"event_id"`
	NumTokens int    `json:"num_tokens"`
}

type EvictResponse struct {
	EventID   string `json:"event_id"`
	NumTokens int    `json:"num_tokens"`
}

func NewLMCacheClient(baseURL string, httpClient *resty.Client, retry RetryConfig, cb *CircuitBreaker) *LMCacheClient {
	return &LMCacheClient{
		httpClient:     httpClient,
		baseURL:        strings.TrimRight(baseURL, "/"),
		retry:          retry,
		circuitBreaker: cb,
	}
}

func (lc *LMCacheClient) Pin(ctx context.Context, req *PinRequest) (*PinResponse, error) {
	var resp PinResponse
	if err := lc.doRequestWithRetry(ctx, "/pin", req, &resp, "pin"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (lc *LMCacheClient) Lookup(ctx context.Context, req *LookupRequest) (*LookupResponse, error) {
	var resp LookupResponse
	if err := lc.doRequestWithRetry(ctx, "/lookup", req, &resp, "lookup"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (lc *LMCacheClient) Compress(ctx context.Context, req *CompressRequest) (*CompressResponse, error) {
	var resp CompressResponse
	if err := lc.doRequestWithRetry(ctx, "/compress", req, &resp, "compress"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (lc *LMCacheClient) Evict(ctx context.Context, req *EvictRequest) (*EvictResponse, error) {
	var resp EvictResponse
	if err := lc.doRequestWithRetry(ctx, "/evict", req, &resp, "evict"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (lc *LMCacheClient) doRequestWithRetry(ctx context.Context, path string, payload any, out any, op string) error {
	var lastErr error
	maxAttempts := lc.retry.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err := lc.execute(func() error {
			return lc.doRequest(ctx, path, payload, out)
		})
		if err == nil {
			return nil
		}
		if err == ErrCircuitBreakerOpen {
			return err
		}
		lastErr = err
		backoff := lc.backoffDuration(attempt)
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return lastErr
}

func (lc *LMCacheClient) doRequest(ctx context.Context, path string, payload any, out any) error {
	url := lc.baseURL + path
	resp, err := lc.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)
	if err != nil {
		return fmt.Errorf("[kvcache] call lmcache: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("[kvcache] lmcache status %d: %s", resp.StatusCode(), strings.TrimSpace(string(resp.Body())))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(resp.Body(), out); err != nil {
		return fmt.Errorf("[kvcache] decode lmcache response: %w", err)
	}
	return nil
}

func (lc *LMCacheClient) execute(operation func() error) error {
	if lc.circuitBreaker == nil {
		return operation()
	}
	return lc.circuitBreaker.Execute(operation)
}

func (lc *LMCacheClient) backoffDuration(attempt int) time.Duration {
	backoff := lc.retry.BaseBackoff
	for i := 0; i < attempt; i++ {
		backoff *= 2
		if backoff >= lc.retry.MaxBackoff {
			return lc.retry.MaxBackoff
		}
	}
	if backoff > lc.retry.MaxBackoff {
		return lc.retry.MaxBackoff
	}
	return backoff
}
