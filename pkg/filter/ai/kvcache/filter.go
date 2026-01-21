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
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type Filter struct {
	cfg   *Config
	store *RedisStore
}

func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	if f.cfg == nil || !f.cfg.Enabled {
		return filter.Continue
	}
	incomingSessionID := strings.TrimSpace(hc.Request.Header.Get(HeaderSessionID))
	if f.cfg.DropClientHeaders != nil && *f.cfg.DropClientHeaders {
		stripIncomingHeaders(hc.Request.Header)
	}

	sessionID, generated := generateSessionID(incomingSessionID)
	if generated {
		hc.AddHeader(HeaderSessionID, sessionID)
	}

	modelName := getModels(hc.Request.URL, f.cfg.TargetModels)
	if !modelAllowed(modelName, f.cfg.TargetModels) {
		return filter.Continue
	}

	policy := defaultPolicy(f.cfg)
	ctx := getContext(hc)
	if sessionID != "" && f.store != nil {
		stored, err := f.store.GetPolicy(ctx, sessionID, policy)
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] kvcache redis get policy failed: %v", err)
		} else {
			policy = stored
		}
	}

	mode := resolveMode(policy)

	if f.cfg.CacheKey.IncludeSessionID {
		hc.Request.Header.Set(HeaderSessionID, sessionID)
	}
	if f.cfg.CacheKey.IncludeModel && modelName != "" {
		hc.Request.Header.Set(HeaderModel, modelName)
	}
	if mode != "" {
		hc.Request.Header.Set(HeaderCacheMode, mode)
	}

	storeContextValue(hc, ctxSessionIDKey, sessionID)
	storeContextValue(hc, ctxCacheModeKey, mode)
	storeContextValue(hc, ctxModelKey, modelName)

	return filter.Continue
}

func (f *Filter) Encode(hc *contexthttp.HttpContext) filter.FilterStatus {
	if f.cfg == nil || !f.cfg.Enabled {
		return filter.Continue
	}
	sessionID, _ := getContextValue(hc, ctxSessionIDKey).(string)
	if sessionID == "" || f.store == nil {
		return filter.Continue
	}
	mode, _ := getContextValue(hc, ctxCacheModeKey).(string)
	result := resolveCacheResult(readResponseHeader(hc, HeaderCacheHit))
	if result == "" {
		return filter.Continue
	}
	if err := f.store.UpdateStats(getContext(hc), sessionID, result, mode); err != nil {
		logger.Warnf("[dubbo-go-pixiu] kvcache redis update stats failed: %v", err)
	}
	return filter.Continue
}
