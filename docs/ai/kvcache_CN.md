## AI KVCache 过滤器配置

[English](./kvcache.md) | 中文

本文档说明如何在 Dubbo-go-Pixiu 中配置和使用 `dgp.filter.ai.kvcache` 过滤器。

该过滤器通过对接 vLLM（`/tokenize`）与 LMCache controller API（`/lookup`、`/pin`、`/compress`、`/evict`），实现：

- cache 感知的路由提示
- 异步缓存管理动作触发
- 主请求链路非阻塞

---

### 架构与请求链路

`dgp.filter.ai.kvcache` 是一个 HTTP Decode 过滤器。典型处理流程如下：

1. 读取请求体，提取 `model` 与 `prompt`（必要时从 `messages` 回退提取）。
2. 在 TokenManager 中记录本地热点统计（`model + prompt`）。
3. 尝试 cache 感知路由：
   - 从本地 token cache 获取 token
   - 调用 LMCache `/lookup`
   - 把优选端点提示写入上下文（`llm_preferred_endpoint_id`）
4. 启动异步缓存管理协程（best-effort）：
   - 调用 vLLM `/tokenize`
   - 必要时调用 LMCache `/lookup`
   - 按策略执行 `compress` / `pin` / `evict`
5. 立即放行后续过滤器链（主请求不被缓存管理阻塞）。

---

### 路由约定（重点）

当前 cache 路由依赖 **instance id 对齐**：

- kvcache 过滤器在上下文写入 `llm_preferred_endpoint_id`
- `dgp.filter.llm.proxy` 读取该值并按 `endpoint.id` 选目标实例

因此要生效必须满足：

- `LMCache lookup 返回的 instance_id` 与 `pixiu cluster endpoint.id` 一致

如果不一致，请求会自动回退到正常负载均衡。

说明：

- 当前实现 **不会** 调 LMCache `/query_worker_info`
- 也就是路由依据是 instance-id 合约，不是动态查询 IP/Port

---

### 配置示例

```yaml
listeners:
  - name: net/http
    protocol_type: HTTP
    address:
      socket_address:
        address: 0.0.0.0
        port: 8888
    filter_chains:
      filters:
        - name: dgp.filter.httpconnectionmanager
          config:
            route_config:
              routes:
                - match:
                    prefix: /
                  route:
                    cluster: vllm_cluster
            http_filters:
              - name: dgp.filter.ai.kvcache
                config:
                  enabled: true
                  vllm_endpoint: "http://127.0.0.1:8000"
                  lmcache_endpoint: "http://127.0.0.1:9000"
                  default_model: "demo"
                  request_timeout: "2s"
                  lookup_routing_timeout: "50ms"
                  hot_window: "5m"
                  hot_max_records: 300
                  token_cache:
                    enabled: true
                    max_size: 1024
                    ttl: "10m"
                  cache_strategy:
                    enable_compression: true
                    enable_pinning: true
                    enable_eviction: true
                    load_threshold: 0.7
                    memory_threshold: 0.85
                    hot_content_threshold: 10
                    pin_instance_id: "vllm-instance-1"
                    pin_location: "LocalCPUBackend"
                    compress_instance_id: "vllm-instance-1"
                    compress_location: "LocalCPUBackend"
                    compress_method: "zstd"
                    evict_instance_id: "vllm-instance-1"
                  retry:
                    max_attempts: 3
                    base_backoff: "100ms"
                    max_backoff: "2s"
                  circuit_breaker:
                    failure_threshold: 5
                    recovery_timeout: "10s"
                    half_open_max_calls: 2
              - name: dgp.filter.llm.proxy
```

---

### 关键配置项说明

`enabled`
- 开关，是否启用 kvcache 过滤器。

`vllm_endpoint`
- `/tokenize` 的上游地址。

`lmcache_endpoint`
- LMCache controller API 的上游地址。

`lookup_routing_timeout`
- Decode 路径里同步 lookup 的短超时时间。

`token_cache`
- 本地内存 token 缓存配置。
- 当前 cache key 为 `model + "\x00" + prompt` 的 SHA-256。

`cache_strategy.load_threshold`
- 比例值，范围 `[0,1]`。
- 当前代码里用于和 CPU 使用率比例比较，决定是否触发压缩。

`cache_strategy.memory_threshold`
- 比例值，范围 `[0,1]`，用于驱逐决策。

`cache_strategy.hot_content_threshold`
- 在 `hot_window` 内达到该访问次数后，判定为热点内容，用于 pin。

`retry`
- LMCache API（`lookup/pin/compress/evict`）调用的重试参数。

`circuit_breaker`
- tokenizer 与 LMCache 调用的熔断保护参数。

---

### 运行与行为说明

1. Best-effort 行为

- tokenize 或 LMCache 调用失败时，主请求链路仍会继续。
- 错误日志统一使用 `[kvcache]` 前缀。

2. 请求取消传播

- 缓存管理使用 request-scoped context + timeout。
- 客户端取消/超时可以及时中断后台缓存任务。

3. 压缩触发语义

- `load_threshold` 按比例处理。
- 当前压缩决策依据 CPU 使用率比例。

4. 真实引擎与 mock

- 真正收益评估需要真实 vLLM + LMCache 实例部署。
- 联调阶段可先用 mock `/tokenize` 与 mock LMCache API 验证链路接入和路由提示是否正确。

---

### 验证清单

1. `endpoint.id` 与 LMCache `instance_id` 对齐。
2. `dgp.filter.ai.kvcache` 放在 `dgp.filter.llm.proxy` 之前。
3. 日志中可观测到：
   - routing lookup 成功/失败
   - strategy 执行失败（若有）
4. 对比开启前后：
   - 命中率
   - p95/p99 延迟
   - 上游实例流量分布
