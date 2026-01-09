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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

import (
	"github.com/redis/go-redis/v9"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.AIKVCacheFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin 作为过滤器入口
	Plugin struct{}

	// FilterFactory 持有配置和共享实例
	FilterFactory struct {
		config       *Config
		cacheManager KVCacheManager
		keyGenerator *CacheKeyGenerator
	}

	// Filter 处理请求链路
	Filter struct {
		factory *FilterFactory
	}

	// Config 表示过滤器的 YAML/JSON 配置
	Config struct {
		Enabled      bool           `yaml:"enabled" json:"enabled"`
		Backend      string         `yaml:"backend" json:"backend"` // 目前仅支持 redis
		EngineType   string         `yaml:"engine_type" json:"engine_type"`
		RedisConfig  RedisConfig    `yaml:"redis" json:"redis"`
		CacheKey     CacheKeyConfig `yaml:"cache_key" json:"cache_key"`
		DefaultTTL   time.Duration  `yaml:"default_ttl" json:"default_ttl"`
		TargetModels []string       `yaml:"target_models" json:"target_models"`
	}

	// RedisConfig 保存 redis 客户端配置
	RedisConfig struct {
		ClusterMode bool     `yaml:"cluster_mode" json:"cluster_mode"`
		Addrs       []string `yaml:"addrs" json:"addrs"`
		Password    string   `yaml:"password" json:"password"`
		KeyPrefix   string   `yaml:"key_prefix" json:"key_prefix"`
		PoolSize    int      `yaml:"pool_size" json:"pool_size"`
		RedisURL    string   `yaml:"redis_url" json:"redis_url"` // 供 vLLM 构造 cache url
	}

	// CacheKeyConfig 控制生成 key 的组成
	CacheKeyConfig struct {
		IncludeModel     bool `yaml:"include_model" json:"include_model"`
		IncludeSessionID bool `yaml:"include_session_id" json:"include_session_id"`
		IncludePrompt    bool `yaml:"include_prompt" json:"include_prompt"`
	}

	// CacheKeyRequest 是生成 key 的输入
	CacheKeyRequest struct {
		SessionID   string
		ModelName   string
		Prompt      string
		Parameters  map[string]any
		RoundNumber int
	}

	// CacheKeyGenerator 生成可复现的缓存键
	CacheKeyGenerator struct {
		prefix string
		cfg    CacheKeyConfig
	}

	// KVCache 表示存储在 redis 的缓存实体
	KVCache struct {
		SessionID  string            `json:"session_id"`
		ModelName  string            `json:"model_name"`
		PromptHash string            `json:"prompt_hash"`
		Keys       []byte            `json:"keys"`
		Values     []byte            `json:"values"`
		TokenCount int               `json:"token_count"`
		CreatedAt  time.Time         `json:"created_at"`
		AccessedAt time.Time         `json:"accessed_at"`
		Metadata   map[string]string `json:"metadata"`
	}

	// KVCacheManager 抽象缓存读写接口
	KVCacheManager interface {
		Get(ctx context.Context, key string) (*KVCache, error)
		Set(ctx context.Context, key string, cache *KVCache, ttl time.Duration) error
		Close() error
	}

	redisCacheManager struct {
		client    redis.UniversalClient
		keyPrefix string
	}
)

// Kind 返回过滤器标识
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory 创建工厂实例
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		config: &Config{
			Enabled:    true,
			Backend:    "redis",
			DefaultTTL: time.Hour,
			CacheKey: CacheKeyConfig{
				IncludeModel:     true,
				IncludeSessionID: true,
				IncludePrompt:    true,
			},
		},
	}, nil
}

// Config 返回配置实例
func (factory *FilterFactory) Config() any {
	return factory.config
}

// Apply 在配置加载后初始化共享组件
func (factory *FilterFactory) Apply() error {
	if !factory.config.Enabled {
		return nil
	}

	if factory.config.Backend != "" && factory.config.Backend != "redis" {
		return fmt.Errorf("[KVCache] unsupported backend %s", factory.config.Backend)
	}

	// 初始化 redis manager
	managerCfg := factory.config.RedisConfig
	// 为避免双重前缀，Redis 层不再追加前缀
	managerCfg.KeyPrefix = ""
	manager, err := newRedisCacheManager(managerCfg)
	if err != nil {
		return err
	}
	factory.cacheManager = manager

	prefix := "pixiu:kvcache"
	if factory.config.RedisConfig.KeyPrefix != "" {
		prefix = factory.config.RedisConfig.KeyPrefix
	}
	factory.keyGenerator = &CacheKeyGenerator{
		prefix: prefix,
		cfg:    factory.config.CacheKey,
	}

	if factory.config.DefaultTTL <= 0 {
		factory.config.DefaultTTL = time.Hour
	}

	return nil
}

// PrepareFilterChain 将过滤器挂载到链路
func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{factory: factory}
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)
	return nil
}

// Decode 阶段尝试把缓存注入上游请求
func (f *Filter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	if ctx.Params == nil {
		ctx.Params = make(map[string]any)
	}

	if !f.factory.config.Enabled {
		return filter.Continue
	}

	modelName := extractModelName(ctx)
	if !f.factory.isTargetModel(modelName) {
		return filter.Continue
	}

	var cacheKey string
	var err error

	engineType := f.detectEngineType(ctx)
	switch engineType {
	case "sglang":
		adapter := NewSGLangAdapter(f.factory.cacheManager, f.factory.keyGenerator)
		cacheKey, err = adapter.InjectKVCache(ctx)
	case "vllm":
		adapter := NewVLLMAdapter(f.factory.cacheManager, f.factory.keyGenerator, f.factory.config.RedisConfig.getRedisURL())
		cacheKey, err = adapter.InjectKVCache(ctx)
	default:
		adapter := NewSGLangAdapter(f.factory.cacheManager, f.factory.keyGenerator)
		cacheKey, err = adapter.InjectKVCache(ctx)
	}

	if err != nil {
		logger.Errorf("[KVCache] Failed to inject KV Cache: %v", err)
	}

	ctx.Params["kv_cache_key"] = cacheKey
	ctx.Params["engine_type"] = engineType

	return filter.Continue
}

// Encode 阶段从上游响应提取并保存缓存
func (f *Filter) Encode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	cacheKeyObj, ok := ctx.Params["kv_cache_key"]
	if !ok {
		return filter.Continue
	}
	cacheKey, _ := cacheKeyObj.(string)

	engineType, _ := ctx.Params["engine_type"].(string)

	var err error
	switch engineType {
	case "sglang":
		adapter := NewSGLangAdapter(f.factory.cacheManager, f.factory.keyGenerator)
		err = adapter.ExtractAndSaveKVCache(ctx, cacheKey)
	case "vllm":
		adapter := NewVLLMAdapter(f.factory.cacheManager, f.factory.keyGenerator, f.factory.config.RedisConfig.getRedisURL())
		err = adapter.ExtractAndSaveKVCache(ctx, cacheKey)
	default:
		adapter := NewSGLangAdapter(f.factory.cacheManager, f.factory.keyGenerator)
		err = adapter.ExtractAndSaveKVCache(ctx, cacheKey)
	}

	if err != nil {
		logger.Errorf("[KVCache] Failed to extract KV Cache: %v", err)
	}

	return filter.Continue
}

// detectEngineType 根据配置/头/cluster 推断推理引擎类型
func (f *Filter) detectEngineType(ctx *contexthttp.HttpContext) string {
	if engineType := f.factory.config.EngineType; engineType != "" {
		return engineType
	}

	if engine := ctx.Request.Header.Get("X-Engine-Type"); engine != "" {
		return engine
	}

	if rEntry := ctx.GetRouteEntry(); rEntry != nil {
		if strings.Contains(rEntry.Cluster, "sglang") {
			return "sglang"
		}
		if strings.Contains(rEntry.Cluster, "vllm") {
			return "vllm"
		}
	}

	// 默认 sglang
	return "sglang"
}

// extractModelName 在不破坏 body 的前提下获取模型名
func extractModelName(ctx *contexthttp.HttpContext) string {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return ctx.Request.Header.Get("X-Model")
	}
	_ = ctx.Request.Body.Close()

	var req struct {
		Model string `json:"model"`
	}
	_ = json.Unmarshal(bodyBytes, &req)

	restoreBody(ctx, bodyBytes)

	if req.Model != "" {
		return req.Model
	}
	return ctx.Request.Header.Get("X-Model")
}

func restoreBody(ctx *contexthttp.HttpContext, body []byte) {
	reader := bytes.NewReader(body)
	ctx.Request.Body = io.NopCloser(reader)
	ctx.Request.ContentLength = int64(len(body))
	if ctx.Request.GetBody == nil {
		ctx.Request.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
	}
}

// isTargetModel 检查模型是否在目标列表内（空列表代表全部通过）
func (factory *FilterFactory) isTargetModel(model string) bool {
	if len(factory.config.TargetModels) == 0 || model == "" {
		return true
	}
	for _, m := range factory.config.TargetModels {
		if m == model {
			return true
		}
	}
	return false
}

// recordCacheHit 记录命中（可在此扩展指标）
func recordCacheHit(ctx *contexthttp.HttpContext, model string) {
	logger.Infof("[KVCache] cache hit, model=%s", model)
}

// recordCacheMiss 记录未命中（可在此扩展指标）
func recordCacheMiss(ctx *contexthttp.HttpContext, model string) {
	logger.Infof("[KVCache] cache miss, model=%s", model)
}

// Generate 根据配置拼装缓存键
func (g *CacheKeyGenerator) Generate(req *CacheKeyRequest) string {
	parts := []string{g.prefix}
	if g.cfg.IncludeModel && req.ModelName != "" {
		parts = append(parts, "model", req.ModelName)
	}
	if g.cfg.IncludeSessionID && req.SessionID != "" {
		parts = append(parts, "sid", req.SessionID)
	}
	if req.RoundNumber > 0 {
		parts = append(parts, "round", fmt.Sprintf("%d", req.RoundNumber))
	}
	if g.cfg.IncludePrompt && req.Prompt != "" {
		sum := sha256.Sum256([]byte(req.Prompt))
		parts = append(parts, "p", hex.EncodeToString(sum[:8]))
	}
	return strings.Join(parts, ":")
}

func newRedisCacheManager(cfg RedisConfig) (KVCacheManager, error) {
	if len(cfg.Addrs) == 0 {
		return nil, errors.New("[KVCache] redis addrs required")
	}

	options := &redis.UniversalOptions{
		Addrs:      cfg.Addrs,
		Password:   cfg.Password,
		PoolSize:   cfg.PoolSize,
		MasterName: "",
	}
	client := redis.NewUniversalClient(options)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("[KVCache] failed to ping redis: %w", err)
	}

	return &redisCacheManager{
		client:    client,
		keyPrefix: cfg.KeyPrefix,
	}, nil
}

func (m *redisCacheManager) prefixed(key string) string {
	if m.keyPrefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", strings.TrimSuffix(m.keyPrefix, ":"), key)
}

func (m *redisCacheManager) Get(ctx context.Context, key string) (*KVCache, error) {
	raw, err := m.client.Get(ctx, m.prefixed(key)).Bytes()
	if err != nil {
		return nil, err
	}
	cache := &KVCache{}
	if err := json.Unmarshal(raw, cache); err != nil {
		return nil, err
	}
	cache.AccessedAt = time.Now()
	return cache, nil
}

func (m *redisCacheManager) Set(ctx context.Context, key string, cache *KVCache, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = time.Hour
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return m.client.Set(ctx, m.prefixed(key), data, ttl).Err()
}

func (m *redisCacheManager) Close() error {
	return m.client.Close()
}

// getRedisURL 返回 vLLM 适配器需要的 redis URL
func (r RedisConfig) getRedisURL() string {
	if r.RedisURL != "" {
		return r.RedisURL
	}
	if len(r.Addrs) > 0 {
		if r.Password != "" {
			return fmt.Sprintf("redis://:%s@%s", r.Password, r.Addrs[0])
		}
		return fmt.Sprintf("redis://%s", r.Addrs[0])
	}
	return ""
}

// maskRedisURL 将 redis URL 中的密码脱敏，防止日志泄漏
func maskRedisURL(raw string) string {
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.User != nil {
		username := u.User.Username()
		if _, hasPwd := u.User.Password(); hasPwd {
			if username != "" {
				u.User = url.UserPassword(username, "***")
			} else {
				u.User = url.UserPassword("", "***")
			}
		}
	}
	return u.String()
}
