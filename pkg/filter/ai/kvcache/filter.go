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

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"

	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

import "github.com/go-resty/resty/v2"

const (
	Kind = constant.AIKVCacheFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct{}

	FilterFactory struct {
		cfg        *Config
		httpClient *http.Client
		resty      *resty.Client
	}

	Filter struct {
		cfg           *Config
		tokenManager  *TokenManager
		lmcacheClient *LMCacheClient
		cacheStrategy *CacheStrategy
	}
)

func (p *Plugin) Kind() string { return Kind }

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any { return factory.cfg }

func (factory *FilterFactory) Apply() error {
	factory.cfg.ApplyDefaults()
	if err := factory.cfg.Validate(); err != nil {
		return err
	}
	cfg := factory.cfg
	factory.httpClient = &http.Client{
		Timeout: cfg.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
		},
	}
	factory.resty = resty.NewWithClient(factory.httpClient).
		SetTimeout(cfg.RequestTimeout)
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(_ *contexthttp.HttpContext, chain filter.FilterChain) error {
	cfgCopy := factory.cfg.DeepCopy()
	cbToken := NewCircuitBreaker(cfgCopy.CircuitBreaker)
	cbLMCache := NewCircuitBreaker(cfgCopy.CircuitBreaker)
	tokenManager := NewTokenManager(cfgCopy.VLLMEndpoint, factory.resty, cfgCopy.TokenCache, cbToken, cfgCopy.HotWindow, cfgCopy.HotMaxRecords)
	lmcacheClient := NewLMCacheClient(cfgCopy.LMCacheEndpoint, factory.resty, cfgCopy.Retry, cbLMCache)
	cacheStrategy := NewCacheStrategy(cfgCopy.CacheStrategy, lmcacheClient, tokenManager)

	f := &Filter{
		cfg:           cfgCopy,
		tokenManager:  tokenManager,
		lmcacheClient: lmcacheClient,
		cacheStrategy: cacheStrategy,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	if f.cfg == nil || !f.cfg.Enabled {
		return filter.Continue
	}
	if f.cacheStrategy != nil {
		f.cacheStrategy.RecordRequest()
	}
	body, err := readRequestBody(hc.Request)
	if err != nil {
		logger.Warnf("[KVCache] read request body failed: %v", err)
		return filter.Continue
	}
	prompt, model, err := extractPromptAndModel(body)
	if err != nil {
		logger.Warnf("[KVCache] parse request body failed: %v", err)
		return filter.Continue
	}
	if prompt == "" {
		return filter.Continue
	}
	if model == "" {
		model = f.cfg.DefaultModel
	}

	f.tokenManager.RecordHot(model, prompt)

	cacheStatus, routed := f.tryRouteToCachedInstance(hc, model, prompt)

	ctx, cancel := context.WithTimeout(hc.Ctx, effectiveTimeout(hc, f.cfg))
	go func() {
		defer cancel()
		f.manageCache(ctx, model, prompt, body, cacheStatus, routed)
	}()
	return filter.Continue
}

func (f *Filter) tryRouteToCachedInstance(hc *contexthttp.HttpContext, model string, prompt string) (*LookupResponse, bool) {
	if f == nil || f.tokenManager == nil || f.lmcacheClient == nil {
		return nil, false
	}
	tokens, ok := f.tokenManager.GetCachedTokens(model, prompt)
	if !ok || len(tokens) == 0 {
		logger.Debugf("[KVCache] routing lookup skipped: token cache miss")
		return nil, false
	}
	timeout := effectiveTimeout(hc, f.cfg)
	if f.cfg != nil && f.cfg.LookupRoutingTimeout > 0 && f.cfg.LookupRoutingTimeout < timeout {
		timeout = f.cfg.LookupRoutingTimeout
	}
	ctx, cancel := context.WithTimeout(hc.Ctx, timeout)
	defer cancel()
	cacheStatus, err := f.lmcacheClient.Lookup(ctx, &LookupRequest{Tokens: tokens})
	if err != nil {
		logger.Debugf("[KVCache] routing lookup failed: %v", err)
		return nil, false
	}
	instanceID := selectPreferredInstanceID(cacheStatus)
	if instanceID == "" {
		logger.Debugf("[KVCache] routing lookup returned empty instance")
		return cacheStatus, false
	}
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[constant.LLMPreferredEndpointID] = instanceID
	logger.Debugf("[KVCache] routing preferred endpoint set: %s", instanceID)
	return cacheStatus, true
}
