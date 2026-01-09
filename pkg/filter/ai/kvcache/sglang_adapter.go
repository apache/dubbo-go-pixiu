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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// SGLangAdapter 适配 SGLang API 的 KV Cache 管理
type SGLangAdapter struct {
	cacheManager KVCacheManager
	keyGenerator *CacheKeyGenerator
}

// ChatCompletionRequest SGLang/OpenAI Chat Completion 请求
type ChatCompletionRequest struct {
	Model             string                 `json:"model"`
	Messages          []Message              `json:"messages"`
	Temperature       float64                `json:"temperature,omitempty"`
	MaxTokens         int                    `json:"max_tokens,omitempty"`
	Stream            bool                   `json:"stream,omitempty"`
	SessionID         string                 `json:"session_id,omitempty"`
	ReturnKVCache     bool                   `json:"return_kv_cache,omitempty"`
	KVCache           *KVCachePayload        `json:"kv_cache,omitempty"`
	EnablePrefixCache bool                   `json:"enable_prefix_cache,omitempty"`
	ExtraParams       map[string]interface{} `json:"-"`
}

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// KVCachePayload KV Cache 数据载荷
type KVCachePayload struct {
	CacheKey      string            `json:"cache_key"`
	CacheData     string            `json:"cache_data,omitempty"`     // Base64 编码的 KV Cache
	CacheURL      string            `json:"cache_url,omitempty"`      // Redis URL（vLLM 使用）
	CacheMetadata map[string]string `json:"cache_metadata,omitempty"` // 元数据
}

// ChatCompletionResponse SGLang/OpenAI Chat Completion 响应
type ChatCompletionResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Choices []Choice        `json:"choices"`
	Usage   Usage           `json:"usage"`
	KVCache *KVCachePayload `json:"kv_cache,omitempty"`
}

// Choice 响应选项
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage Token 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewSGLangAdapter 创建 SGLang 适配器
func NewSGLangAdapter(cacheManager KVCacheManager, keyGen *CacheKeyGenerator) *SGLangAdapter {
	return &SGLangAdapter{
		cacheManager: cacheManager,
		keyGenerator: keyGen,
	}
}

// InjectKVCache 在 Decode 阶段注入 KV Cache
func (a *SGLangAdapter) InjectKVCache(ctx *contexthttp.HttpContext) (string, error) {
	// 1. 读取并解析请求体
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}
	ctx.Request.Body.Close()

	var req ChatCompletionRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return "", fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 2. 提取会话信息
	sessionID := req.SessionID
	if sessionID == "" {
		// 如果没有 session_id，尝试从 Header 获取
		sessionID = ctx.Request.Header.Get("X-Session-ID")
	}
	if sessionID == "" {
		// 没有会话 ID，跳过缓存
		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return "", nil
	}

	// 3. 生成缓存键
	cacheKey := a.keyGenerator.Generate(&CacheKeyRequest{
		SessionID:   sessionID,
		ModelName:   req.Model,
		Prompt:      a.extractPrompt(&req),
		Parameters:  a.extractParams(&req),
		RoundNumber: a.extractRoundNumber(ctx),
	})

	// 4. 从 Redis 查询 KV Cache
	reqCtx := ctx.Request.Context()
	cache, err := a.cacheManager.Get(reqCtx, cacheKey)
	if err == nil && cache != nil {
		// 缓存命中！注入到请求
		logger.Infof("[KVCache] Cache HIT for session %s, key: %s", sessionID, cacheKey)

		encodedKV, encodeErr := encodeKVPair(cache.Keys, cache.Values)
		if encodeErr != nil {
			return "", fmt.Errorf("failed to encode kv cache: %w", encodeErr)
		}

		req.KVCache = &KVCachePayload{
			CacheKey:  cacheKey,
			CacheData: encodedKV,
		}

		// 标记启用 KV Cache
		req.ReturnKVCache = true

		// 重新序列化请求
		modifiedBody, err := json.Marshal(req)
		if err != nil {
			return "", fmt.Errorf("failed to marshal modified request: %w", err)
		}

		// 替换请求体
		ctx.Request.Body = io.NopCloser(bytes.NewReader(modifiedBody))
		ctx.Request.ContentLength = int64(len(modifiedBody))

		recordCacheHit(ctx, req.Model)
		return cacheKey, nil
	}

	// 缓存未命中，但需要请求 SGLang 返回 KV Cache
	logger.Infof("[KVCache] Cache MISS for session %s, key: %s", sessionID, cacheKey)
	req.ReturnKVCache = true

	modifiedBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal modified request: %w", err)
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(modifiedBody))
	ctx.Request.ContentLength = int64(len(modifiedBody))

	recordCacheMiss(ctx, req.Model)
	return cacheKey, nil
}

// ExtractAndSaveKVCache 在 Encode 阶段提取并保存 KV Cache
func (a *SGLangAdapter) ExtractAndSaveKVCache(ctx *contexthttp.HttpContext, cacheKey string) error {
	if cacheKey == "" {
		return nil
	}

	// 1. 读取响应体
	resp, ok := ctx.SourceResp.(*http.Response)
	if !ok || resp == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	// 2. 解析响应
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		// 不是 JSON 格式，可能是流式响应，跳过
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return nil
	}

	// 3. 提取 KV Cache
	if chatResp.KVCache != nil && chatResp.KVCache.CacheData != "" {
		keys, values, decodeErr := decodeKVPair(chatResp.KVCache.CacheData)
		if decodeErr != nil {
			logger.Errorf("[KVCache] Failed to decode kv cache data for key %s: %v", cacheKey, decodeErr)
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			return nil
		}

		// 构建 KV Cache 对象
		kvCache := &KVCache{
			SessionID:  extractSessionID(ctx),
			ModelName:  chatResp.Model,
			PromptHash: cacheKey,
			Keys:       keys,
			Values:     values,
			TokenCount: chatResp.Usage.PromptTokens,
			CreatedAt:  time.Now(),
			AccessedAt: time.Now(),
			Metadata:   chatResp.KVCache.CacheMetadata,
		}

		// 4. 异步保存到 Redis
		go func() {
			saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := a.cacheManager.Set(saveCtx, cacheKey, kvCache, 3600*time.Second); err != nil {
				logger.Errorf("[KVCache] Failed to save cache for key %s: %v", cacheKey, err)
			} else {
				logger.Infof("[KVCache] Cache saved for key: %s, size: %d bytes",
					cacheKey, len(kvCache.Keys)+len(kvCache.Values))
			}
		}()
	}

	// 5. 恢复响应体（供后续 Filter 使用）
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return nil
}

// extractPrompt 提取 Prompt 文本
func (a *SGLangAdapter) extractPrompt(req *ChatCompletionRequest) string {
	if len(req.Messages) == 0 {
		return ""
	}
	// 使用最后一条用户消息作为 Prompt
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			return req.Messages[i].Content
		}
	}
	return req.Messages[len(req.Messages)-1].Content
}

// extractParams 提取推理参数
func (a *SGLangAdapter) extractParams(req *ChatCompletionRequest) map[string]any {
	params := make(map[string]any)
	if req.Temperature > 0 {
		params["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		params["max_tokens"] = req.MaxTokens
	}
	return params
}

// extractRoundNumber 提取对话轮次
func (a *SGLangAdapter) extractRoundNumber(ctx *contexthttp.HttpContext) int {
	// 从上下文或 Header 中提取
	if round := ctx.Request.Header.Get("X-Round-Number"); round != "" {
		var num int
		fmt.Sscanf(round, "%d", &num)
		return num
	}
	return 1
}

func extractSessionID(ctx *contexthttp.HttpContext) string {
	return ctx.Request.Header.Get("X-Session-ID")
}

// encodeKVPair 将键值对编码为 base64(JSON) 字符串，避免长度不一致导致的数据损坏
func encodeKVPair(keys, values []byte) (string, error) {
	payload := struct {
		K []byte `json:"k"`
		V []byte `json:"v"`
	}{
		K: keys,
		V: values,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

// decodeKVPair 解码 encodeKVPair 生成的字符串
func decodeKVPair(encoded string) ([]byte, []byte, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, nil, err
	}
	var payload struct {
		K []byte `json:"k"`
		V []byte `json:"v"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, nil, err
	}
	return payload.K, payload.V, nil
}
