# KV Cache-Aware Routing

Current source anchors:

- `pkg/filter/ai/kvcache/config.go`
- `pkg/filter/ai/kvcache/filter.go`
- `pkg/filter/ai/kvcache/handlers.go`
- `docs/ai/kvcache.md`

`dgp.filter.ai.kvcache` is a Decode filter. It parses the request
body, records prompt hotness, optionally asks LMCache where a prompt is
already cached, writes `llm_preferred_endpoint_id` into request context,
and continues the chain. `dgp.filter.llm.proxy` then tries to route to
that endpoint id.

## Hard Contract

LMCache lookup `instance_id` must match Pixiu `clusters[].endpoints[].id`.

If LMCache returns `vllm-a`, then the LLM cluster must contain:

```yaml
endpoints:
  - id: "vllm-a"
    socket_address:
      address: "127.0.0.1"
      port: 8001
```

If no endpoint id matches, Pixiu falls back to normal load balancing.

## Filter Order

Use:

```yaml
http_filters:
  - name: dgp.filter.ai.kvcache
    config: ...
  - name: dgp.filter.llm.tokenizer
    config: ...
  - name: dgp.filter.llm.proxy
    config: ...
```

`kvcache` must run before proxy. `tokenizer` should also be before
proxy so request timing begins before the upstream call.

## Required Fields When Enabled

```yaml
enabled: true
vllm_endpoint: "http://127.0.0.1:8000"
lmcache_endpoint: "http://127.0.0.1:9000"
default_model: "Qwen2.5-3B-Instruct"
request_timeout: "2s"
lookup_routing_timeout: "50ms"
```

When `enabled: true`, `vllm_endpoint` and `lmcache_endpoint` are
required by current `Config.Validate()`.

## Ratios and Non-Negative Values

These must be in `[0,1]`:

- `cache_strategy.memory_threshold`
- `cache_strategy.load_threshold`

These must be non-negative:

- `token_cache.max_size`
- `cache_strategy.hot_content_threshold`
- `retry.max_attempts`
- `hot_max_records`
- `hot_max_keys`

## Best-Effort Behavior

KV cache management is best effort. Tokenize/lookup/pin/compress/evict
failures are logged with `[kvcache]` and should not block the main LLM
request path.
