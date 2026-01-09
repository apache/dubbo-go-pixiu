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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// VLLMAdapter 适配 vLLM API 的 KV Cache 管理
type VLLMAdapter struct {
	cacheManager KVCacheManager
	keyGenerator *CacheKeyGenerator
	redisURL     string // Redis 集群地址，用于 vLLM 直接访问
	safeRedisURL string // 脱敏后的 redisURL，避免日志泄漏
}

// VLLMKVCacheConfig vLLM KV Cache 配置
type VLLMKVCacheConfig struct {
	EnableOffload  bool   `json:"enable_offload"`
	OffloadBackend string `json:"offload_backend"` // "redis", "s3", "file"
	CacheKey       string `json:"cache_key"`
	RedisURL       string `json:"redis_url,omitempty"`
	S3Bucket       string `json:"s3_bucket,omitempty"`
	FilePath       string `json:"file_path,omitempty"`
}

// VLLMExternalKVCache 外部 KV Cache 引用
type VLLMExternalKVCache struct {
	CacheKey string `json:"cache_key"`
	CacheURL string `json:"cache_url"` // 完整的 Redis/S3/File URL
}

// VLLMChatCompletionRequest vLLM 扩展的请求
type VLLMChatCompletionRequest struct {
	ChatCompletionRequest                      // 嵌入标准请求
	KVCacheConfig         *VLLMKVCacheConfig   `json:"kv_cache_config,omitempty"`
	ExternalKVCache       *VLLMExternalKVCache `json:"external_kv_cache,omitempty"`
}

// VLLMKVCacheInfo vLLM KV Cache 信息
type VLLMKVCacheInfo struct {
	BlockTable  []int `json:"block_table"`
	NumBlocks   int   `json:"num_blocks"`
	BlockSize   int   `json:"block_size"`
	TotalTokens int   `json:"total_tokens"`
}

// VLLMChatCompletionResponse vLLM 扩展的响应
type VLLMChatCompletionResponse struct {
	ChatCompletionResponse                  // 嵌入标准响应
	KVCacheInfo            *VLLMKVCacheInfo `json:"kv_cache_info,omitempty"`
}

// NewVLLMAdapter 创建 vLLM 适配器
func NewVLLMAdapter(cacheManager KVCacheManager, keyGen *CacheKeyGenerator, redisURL string) *VLLMAdapter {
	return &VLLMAdapter{
		cacheManager: cacheManager,
		keyGenerator: keyGen,
		redisURL:     redisURL,
		safeRedisURL: maskRedisURL(redisURL),
	}
}

// InjectKVCache 在 Decode 阶段注入 KV Cache（vLLM 方式）
func (a *VLLMAdapter) InjectKVCache(ctx *contexthttp.HttpContext) (string, error) {
	// 1. 读取并解析请求体
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}
	ctx.Request.Body.Close()

	var req VLLMChatCompletionRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return "", fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// 2. 提取会话信息
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = ctx.Request.Header.Get("X-Session-ID")
	}
	if sessionID == "" {
		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return "", nil
	}

	// 3. 生成缓存键
	cacheKey := a.keyGenerator.Generate(&CacheKeyRequest{
		SessionID:   sessionID,
		ModelName:   req.Model,
		Prompt:      extractPromptFromRequest(&req.ChatCompletionRequest),
		RoundNumber: extractRoundNumberFromContext(ctx),
	})

	// 4. 检查 Redis 中是否有缓存
	reqCtx := ctx.Request.Context()
	cache, err := a.cacheManager.Get(reqCtx, cacheKey)

	if err == nil && cache != nil {
		// 缓存命中！vLLM 支持直接引用 Redis URL
		logger.Infof("[KVCache] Cache HIT for session %s, key: %s (vLLM)", sessionID, cacheKey)

		// 构建 Redis URL
		cacheURL := fmt.Sprintf("%s/kvcache/%s", a.redisURL, cacheKey)

		req.ExternalKVCache = &VLLMExternalKVCache{
			CacheKey: cacheKey,
			CacheURL: cacheURL,
		}

		recordCacheHit(ctx, req.Model)
	} else {
		// 缓存未命中，配置 vLLM 将 KV Cache 卸载到 Redis
		logger.Infof("[KVCache] Cache MISS for session %s, key: %s (vLLM)", sessionID, cacheKey)

		req.KVCacheConfig = &VLLMKVCacheConfig{
			EnableOffload:  true,
			OffloadBackend: "redis",
			CacheKey:       cacheKey,
			RedisURL:       a.redisURL,
		}

		recordCacheMiss(ctx, req.Model)
	}

	// 5. 重新序列化请求
	modifiedBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal modified request: %w", err)
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(modifiedBody))
	ctx.Request.ContentLength = int64(len(modifiedBody))

	return cacheKey, nil
}

// ExtractAndSaveKVCache 在 Encode 阶段提取并保存 KV Cache（vLLM 方式）
func (a *VLLMAdapter) ExtractAndSaveKVCache(ctx *contexthttp.HttpContext, cacheKey string) error {
	if cacheKey == "" {
		return nil
	}

	// vLLM 已经直接将 KV Cache 写入 Redis，Pixiu 只需记录元数据
	resp, ok := ctx.SourceResp.(*http.Response)
	if !ok || resp == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	var vllmResp VLLMChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &vllmResp); err != nil {
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return nil
	}

	// 记录 KV Cache 元数据到 Pixiu
	if vllmResp.KVCacheInfo != nil {
		logger.Infof("[KVCache] vLLM saved cache for key: %s, blocks: %d, tokens: %d",
			cacheKey, vllmResp.KVCacheInfo.NumBlocks, vllmResp.KVCacheInfo.TotalTokens)

		// 可选：将元数据保存到数据库或监控系统
		go a.recordKVCacheMetadata(cacheKey, vllmResp.KVCacheInfo)
	}

	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return nil
}

func (a *VLLMAdapter) recordKVCacheMetadata(cacheKey string, info *VLLMKVCacheInfo) {
	// 实现元数据记录逻辑（可选）
	// 例如：保存到 Prometheus、数据库等
}

func extractPromptFromRequest(req *ChatCompletionRequest) string {
	if len(req.Messages) == 0 {
		return ""
	}
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			return req.Messages[i].Content
		}
	}
	return req.Messages[len(req.Messages)-1].Content
}

func extractRoundNumberFromContext(ctx *contexthttp.HttpContext) int {
	if round := ctx.Request.Header.Get("X-Round-Number"); round != "" {
		var num int
		if _, err := fmt.Sscanf(round, "%d", &num); err == nil && num > 0 {
			return num
		}
	}
	return 1
}
