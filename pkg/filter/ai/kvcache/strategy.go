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
)

type CacheStrategy struct {
	config        CacheStrategyConfig
	loadMonitor   *LoadMonitor
	lmcacheClient *LMCacheClient
	tokenManager  *TokenManager
}

type StrategyDecision struct {
	ShouldCompress bool
	ShouldPin      bool
	ShouldEvict    bool
	Reason         string
}

func NewCacheStrategy(cfg CacheStrategyConfig, client *LMCacheClient, tokenManager *TokenManager) *CacheStrategy {
	return &CacheStrategy{
		config:        cfg,
		loadMonitor:   NewLoadMonitor(),
		lmcacheClient: client,
		tokenManager:  tokenManager,
	}
}

func (cs *CacheStrategy) RecordRequest() {
	if cs == nil || cs.loadMonitor == nil {
		return
	}
	cs.loadMonitor.RecordRequest()
}

func (cs *CacheStrategy) MakeDecision(_ context.Context, cacheStatus *LookupResponse, model string, prompt string) *StrategyDecision {
	if cs == nil {
		return &StrategyDecision{}
	}
	decision := &StrategyDecision{}
	metrics := cs.loadMonitor.Snapshot()

	// load_threshold is validated as a ratio [0,1], so only ratio-based metrics
	// should participate in this decision.
	if cs.config.EnableCompression && cs.config.LoadThreshold > 0 &&
		metrics.CPUUsage >= cs.config.LoadThreshold {
		decision.ShouldCompress = true
		decision.Reason = "high_load"
	}
	if cs.config.EnableEviction && cs.config.MemoryThreshold > 0 && metrics.MemoryUsage >= cs.config.MemoryThreshold {
		decision.ShouldEvict = true
		decision.Reason = "memory_threshold"
	}
	if cs.config.EnablePinning && cs.tokenManager != nil &&
		cs.tokenManager.IsHot(model, prompt, cs.config.HotContentThreshold) {
		decision.ShouldPin = true
		decision.Reason = "hot_content"
	}
	return decision
}

func (cs *CacheStrategy) ExecuteDecision(ctx context.Context, decision *StrategyDecision, tokens []int) error {
	if cs == nil || decision == nil {
		return nil
	}
	if decision.ShouldCompress {
		_, err := cs.lmcacheClient.Compress(ctx, &CompressRequest{
			Tokens:     tokens,
			InstanceID: cs.config.CompressInstanceID,
			Location:   cs.config.CompressLocation,
			Method:     cs.config.CompressMethod,
		})
		if err != nil {
			return err
		}
	}
	if decision.ShouldPin {
		_, err := cs.lmcacheClient.Pin(ctx, &PinRequest{
			Tokens:     tokens,
			InstanceID: cs.config.PinInstanceID,
			Location:   cs.config.PinLocation,
		})
		if err != nil {
			return err
		}
	}
	if decision.ShouldEvict {
		_, err := cs.lmcacheClient.Evict(ctx, &EvictRequest{
			Tokens:     tokens,
			InstanceID: cs.config.EvictInstanceID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
