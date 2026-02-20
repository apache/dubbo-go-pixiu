## AI KVCache Filter Configuration

English | [中文](./kvcache_CN.md)

This document explains how to configure and use the `dgp.filter.ai.kvcache` filter in Dubbo-go-Pixiu.

The filter integrates with vLLM (`/tokenize`) and LMCache controller APIs (`/lookup`, `/pin`, `/compress`, `/evict`) to:

- provide cache-aware routing hints
- trigger cache-management actions asynchronously
- keep the main request path non-blocking

---

### Architecture and Request Flow

`dgp.filter.ai.kvcache` is an HTTP decode filter. A typical request flow is:

1. Parse request body and extract `model` and `prompt` (or fallback from `messages`).
2. Record local hotness statistics (`model + prompt`) in the token manager.
3. Try cache-aware routing:
   - read token cache for prompt
   - call LMCache `/lookup`
   - set a preferred endpoint hint in context (`llm_preferred_endpoint_id`)
4. Start an async cache-management goroutine (best-effort):
   - call vLLM `/tokenize`
   - call LMCache `/lookup` if needed
   - execute strategy decisions (`compress` / `pin` / `evict`)
5. Continue the filter chain immediately (main request is not blocked by cache management).

---

### Routing Contract (Important)

Current cache-aware routing uses **instance id matching**:

- The kvcache filter writes `llm_preferred_endpoint_id` into request context.
- `dgp.filter.llm.proxy` reads this value and tries to select an endpoint by `endpoint.id`.

So for routing to work:

- `LMCache lookup instance_id` must equal `pixiu cluster endpoint.id`.

If no match exists, the request falls back to normal load-balancing behavior.

Note:

- Current implementation does **not** call LMCache `/query_worker_info`.
- This means kvcache routing is based on instance-id contract, not dynamic IP/port discovery.

---

### Configuration Example

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

### Key Config Fields

`enabled`
- Enable/disable kvcache filter.

`vllm_endpoint`
- Base URL for `/tokenize`.

`lmcache_endpoint`
- Base URL for LMCache controller APIs.

`lookup_routing_timeout`
- Short timeout for synchronous routing lookup in decode path.

`token_cache`
- Local in-memory token cache settings.
- Cache key is SHA-256 of `model + "\x00" + prompt`.

`cache_strategy.load_threshold`
- Ratio in `[0,1]`.
- Current code compares this threshold with measured CPU usage ratio.

`cache_strategy.memory_threshold`
- Ratio in `[0,1]`, used for eviction decisions.

`cache_strategy.hot_content_threshold`
- Minimum access count within `hot_window` to mark content as hot for pinning.

`retry`
- Retry settings for LMCache API calls (`lookup/pin/compress/evict`).

`circuit_breaker`
- Protection for tokenizer and LMCache client calls.

---

### Operational Notes

1. Best-effort behavior

- If tokenization or LMCache calls fail, main request path still continues.
- Failures are logged with `[kvcache]` prefix.

2. Context cancellation

- Cache-management work uses request-scoped context and timeout.
- Client cancel/timeout can stop ongoing background operations.

3. Compression trigger semantics

- `load_threshold` is treated as a ratio.
- Compression decision currently uses CPU usage ratio check.

4. Real engine vs mock

- Full performance benefit requires real vLLM + LMCache deployment.
- For smoke validation, you can run mock `/tokenize` and LMCache APIs to verify chain integration and routing hint behavior.

---

### Validation Checklist

1. `endpoint.id` matches LMCache `instance_id` for targeted routing.
2. `dgp.filter.ai.kvcache` is placed before `dgp.filter.llm.proxy`.
3. Logs show:
   - routing lookup success/failure
   - strategy action failures (if any)
4. Compare baseline vs enabled:
   - hit ratio
   - p95/p99 latency
   - upstream load distribution

