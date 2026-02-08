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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

import (
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const llmPreferredEndpointIDKey = "llm_preferred_endpoint_id"

func (f *Filter) manageCache(ctx context.Context, model string, prompt string, rawBody []byte, cacheStatus *LookupResponse, lookupDone bool) {
	if ctx.Err() != nil {
		return
	}
	tokens, err := f.tokenManager.GetTokens(ctx, model, prompt, rawBody)
	if err != nil {
		logger.Warnf("[kvcache] tokenize failed: %v", err)
		return
	}
	if ctx.Err() != nil {
		return
	}
	if !lookupDone || cacheStatus == nil {
		cacheStatus, err = f.lmcacheClient.Lookup(ctx, &LookupRequest{Tokens: tokens})
		if err != nil {
			logger.Warnf("[kvcache] lookup failed: %v", err)
			return
		}
	}
	decision := f.cacheStrategy.MakeDecision(ctx, cacheStatus, model, prompt)
	if ctx.Err() != nil {
		return
	}
	if err := f.cacheStrategy.ExecuteDecision(ctx, decision, tokens); err != nil {
		logger.Warnf("[kvcache] execute strategy failed: %v", err)
	}
}

func readRequestBody(req *http.Request) ([]byte, error) {
	if req == nil || req.Body == nil {
		return nil, nil
	}
	if req.GetBody != nil {
		reader, err := req.GetBody()
		if err == nil {
			defer reader.Close()
			return io.ReadAll(reader)
		}
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodyBytes)), nil
	}
	return bodyBytes, nil
}

func extractPromptAndModel(body []byte) (string, string, error) {
	if len(body) == 0 {
		return "", "", nil
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", err
	}
	model, _ := payload["model"].(string)
	prompt := coercePrompt(payload["prompt"])
	if prompt == "" {
		prompt = extractPromptFromMessages(payload["messages"])
	}
	return strings.TrimSpace(prompt), model, nil
}

func coercePrompt(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				parts = append(parts, str)
			}
		}
		return strings.Join(parts, "\n")
	default:
		return ""
	}
}

func extractPromptFromMessages(value any) string {
	msgs, ok := value.([]any)
	if !ok {
		return ""
	}
	parts := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			continue
		}
		if content, ok := msgMap["content"].(string); ok {
			parts = append(parts, content)
		}
	}
	return strings.Join(parts, "\n")
}

func selectPreferredInstanceID(resp *LookupResponse) string {
	if resp == nil || len(resp.LayoutInfo) == 0 {
		return ""
	}
	var (
		selected string
		maxCount int
	)
	for instanceID, layout := range resp.LayoutInfo {
		if layout.TokenCount > maxCount || selected == "" {
			selected = instanceID
			maxCount = layout.TokenCount
		}
	}
	return selected
}

func effectiveTimeout(hc *contexthttp.HttpContext, cfg *Config) time.Duration {
	if cfg == nil {
		return 0
	}
	timeout := cfg.RequestTimeout
	if hc != nil && hc.Timeout > 0 && (timeout <= 0 || hc.Timeout < timeout) {
		timeout = hc.Timeout
	}
	if timeout <= 0 {
		return 2 * time.Second
	}
	return timeout
}

func requestScopedContext(hc *contexthttp.HttpContext) context.Context {
	if hc != nil && hc.Request != nil {
		return hc.Request.Context()
	}
	if hc != nil && hc.Ctx != nil {
		return hc.Ctx
	}
	return context.Background()
}

func (f *Filter) tryRouteToCachedInstance(hc *contexthttp.HttpContext, model string, prompt string) (*LookupResponse, bool) {
	if f == nil || f.tokenManager == nil || f.lmcacheClient == nil {
		return nil, false
	}
	tokens, ok := f.tokenManager.GetCachedTokens(model, prompt)
	if !ok || len(tokens) == 0 {
		logger.Debugf("[kvcache] routing lookup skipped: token cache miss")
		return nil, false
	}
	timeout := effectiveTimeout(hc, f.cfg)
	if f.cfg != nil && f.cfg.LookupRoutingTimeout > 0 && f.cfg.LookupRoutingTimeout < timeout {
		timeout = f.cfg.LookupRoutingTimeout
	}
	ctx, cancel := context.WithTimeout(requestScopedContext(hc), timeout)
	defer cancel()
	cacheStatus, err := f.lmcacheClient.Lookup(ctx, &LookupRequest{Tokens: tokens})
	if err != nil {
		logger.Debugf("[kvcache] routing lookup failed: %v", err)
		return nil, false
	}
	instanceID := selectPreferredInstanceID(cacheStatus)
	if instanceID == "" {
		logger.Debugf("[kvcache] routing lookup returned empty instance")
		return cacheStatus, false
	}
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[llmPreferredEndpointIDKey] = instanceID
	logger.Debugf("[kvcache] routing preferred endpoint set: %s", instanceID)
	return cacheStatus, true
}
