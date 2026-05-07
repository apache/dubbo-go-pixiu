---
name: pixiu-llm-gateway
description: |
  Configure Apache dubbo-go-pixiu as an AI / LLM gateway. Use when the
  user asks to proxy OpenAI-compatible, Claude-compatible, Gemini, vLLM,
  DeepSeek, or self-hosted model endpoints; configure `dgp.filter.llm.proxy`,
  `dgp.filter.llm.tokenizer`, `dgp.filter.ai.kvcache`, LLM retry/fallback,
  endpoint `llm_meta`, token metrics, LMCache/vLLM cache-aware routing, or
  Nacos LLM registry discovery. This skill is for Pixiu `conf.yaml`, not
  `api_config.yaml`.
allowed-tools: [Read, Grep, Glob, Edit, Write, Bash]
metadata:
  version: "0.1.1"
  domain: ai-gateway
  scope: generate-and-validate
  triggers: ["LLM gateway", "AI gateway", "OpenAI proxy", "Claude proxy", "vLLM", "DeepSeek", "kvcache", "LMCache", "tokenizer", "token billing", "llm_meta", "dgp.filter.llm.proxy"]
  pixiu_min_version: "0.6.0"
  experimental: true
  role: specialist
---

# pixiu-llm-gateway — Configuring Pixiu as an LLM Gateway

Pixiu's LLM path is config-driven. The gateway routes normal HTTP
requests to LLM-compatible upstreams through `dgp.filter.llm.proxy`,
optionally records token metrics with `dgp.filter.llm.tokenizer`, and
optionally hints cache-local routing with `dgp.filter.ai.kvcache`.

This skill writes `conf.yaml` fragments. It does not write
`api_config.yaml`; HTTP-to-Dubbo API mapping is a different code path.

## When to Use

Use this skill when the user wants to:

- Expose an OpenAI-compatible chat/completions endpoint through Pixiu.
- Route between multiple LLM providers or self-hosted vLLM instances.
- Configure endpoint-level API keys, retry policy, fallback, or health
  cooldown through `llm_meta`.
- Enable token metrics / token accounting through
  `dgp.filter.llm.tokenizer`.
- Enable KV cache-aware routing through `dgp.filter.ai.kvcache`.
- Discover LLM endpoints dynamically from Nacos via
  `dgp.adapter.llmregistrycenter`.

Do not use this skill for:

- HTTP to Dubbo routes; use `pixiu-http-to-dubbo`.
- MCP tool exposure; use `pixiu-mcp-integration`.
- Authoring a new filter implementation; use `pixiu-filter-author`.

## Step 0 — Verify Current Source First

The LLM module is moving quickly. Before generating yaml, read the
current repo shape:

1. `pkg/common/constant/key.go` for the exact Kind strings.
2. `pkg/filter/llm/proxy/filter.go` for proxy config, endpoint
   selection, API key injection, retry, and fallback behavior.
3. `pkg/filter/llm/tokenizer/tokenizer.go` for token metric behavior.
4. `pkg/filter/ai/kvcache/config.go` and `handlers.go` for cache-aware
   routing and the `endpoint.id` contract.
5. `pkg/model/llm.go`, `pkg/model/cluster.go`, and `pkg/model/base.go`
   for `llm_meta`, endpoint, and `socket_address` fields.
6. `configs/ai_kvcache_config.yaml` and docs under `docs/ai/` if
   present.

If source disagrees with this skill, trust source and note the
difference in the final answer.

## Step 1 — Gather the Eight Things

Do not produce a final config until these are explicit:

1. HTTP listener address and route prefix/path, usually `/v1/`.
2. Upstream mode: static `clusters[]` or Nacos LLM registry.
3. Provider endpoints: host/domain, scheme (`http` or `https`), and
   path convention (`/v1/chat/completions`, `/v1/completions`, etc.).
4. API key handling: per-endpoint `llm_meta.api_key`, injected secret
   placeholder, or caller-supplied Authorization header.
5. Retry policy per endpoint: `NoRetry`, `CountBased`, or
   `ExponentialBackoff`.
6. Fallback behavior: whether endpoint failure should move to the next
   endpoint in the same cluster.
7. Token metrics: whether to add `dgp.filter.llm.tokenizer` and whether
   console logging is acceptable.
8. KV cache: whether vLLM `/tokenize` and LMCache controller endpoints
   exist, and whether LMCache `instance_id` values match Pixiu
   `endpoint.id`.

It is fine to propose safe defaults, but make assumptions explicit.

## Step 2 — Shape the Static Gateway Config

Minimal static gateway shape:

```yaml
static_resources:
  listeners:
    - name: "net/http"
      protocol_type: "HTTP"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8888
      filter_chains:
        filters:
          - name: dgp.filter.httpconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: "/v1/"
                    route:
                      cluster: "llm_cluster"
                      cluster_not_found_response_code: 503
              http_filters:
                - name: dgp.filter.llm.tokenizer
                  config:
                    log_to_console: false
                - name: dgp.filter.llm.proxy
                  config:
                    scheme: "https"
                    timeout: "60s"
  clusters:
    - name: "llm_cluster"
      lb_policy: "RoundRobin"
      endpoints:
        - id: "openai-primary"
          socket_address:
            domains:
              - "api.openai.com"
          llm_meta:
            provider: "openai"
            api_key: "${OPENAI_API_KEY}"
            fallback: true
            health_check_interval: 5000
            retry_policy:
              name: "ExponentialBackoff"
              config:
                times: 3
                initialInterval: "200ms"
                maxInterval: "5s"
                multiplier: 2.0
```

Notes:

- `dgp.filter.llm.proxy` forwards the original request path and query
  to the selected endpoint. If the client calls `/v1/chat/completions`,
  the upstream receives that same path.
- `socket_address.domains[0]` should be a host, not a full URL and not
  a path-bearing value. Current `llm.proxy` puts it into `url.URL.Host`;
  use `api.openai.com`, not `https://api.openai.com` or
  `api.openai.com/v1`.
- `dgp.filter.llm.proxy.config.scheme` is filter-level. Every endpoint
  reached by the same HTTP filter instance uses that one scheme, so do
  not mix `https` providers and `http` local vLLM endpoints in the same
  route/filter configuration. Split them by route/listener/cluster
  arrangement or normalize the upstream scheme.
- When correcting adversarial configs, explicitly say: `scheme` is
  filter-level, `socket_address.domains` is host-only, and full URLs
  such as `https://api.openai.com/v1` do not belong in `domains`.
- `lb_policy` must be one of `Rand`, `RoundRobin`, `RingHashing`,
  `MaglevHashing`, or `WeightRandom`; reject `lb` and check
  `pkg/model/cluster.go` before inventing a policy.
- Endpoint API keys are injected as `Authorization: Bearer <api_key>`.
  Do not hard-code production secrets; use placeholders or the user's
  secret management convention.
- `dgp.filter.llm.tokenizer` should be before `dgp.filter.llm.proxy` so
  Decode starts timing before the upstream call. Its Encode phase reads
  proxy attempt data and token usage from the response.
- Tokenizer metrics read OpenAI-style `usage` fields for unary responses
  and SSE `data:` frames for streaming responses. If upstream omits token
  usage, report partial metrics instead of inventing counts.

## Step 3 — Add KV Cache Only When the Contract Exists

KV cache-aware routing is optional. It only works when LMCache returns
an `instance_id` that matches a Pixiu endpoint `id`.

```yaml
http_filters:
  - name: dgp.filter.ai.kvcache
    config:
      enabled: true
      vllm_endpoint: "http://127.0.0.1:8000"
      lmcache_endpoint: "http://127.0.0.1:9000"
      default_model: "Qwen2.5-3B-Instruct"
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
  - name: dgp.filter.llm.tokenizer
    config:
      log_to_console: false
  - name: dgp.filter.llm.proxy
    config:
      scheme: "http"
      timeout: "60s"
```

Before generating a kvcache config, read the current
`pkg/filter/ai/kvcache/` source and any checked-in AI config examples.

## Step 4 — Use Nacos Discovery Only for Dynamic Endpoint Discovery

Static clusters are simpler. Use `dgp.adapter.llmregistrycenter` only
when LLM services register endpoint metadata into Nacos.

```yaml
adapters:
  - id: "llm-nacos"
    name: dgp.adapter.llmregistrycenter
    config:
      registries:
        nacos:
          protocol: nacos
          address: "127.0.0.1:8848"
          timeout: "5s"
          group: "test_llm_registry_group"
          namespace: "public"
```

The Nacos instance metadata must include `cluster`, `id`, and LLM
metadata keys such as `llm-meta.api_key`,
`llm-meta.retry_policy.name`, and `llm-meta.fallback`. Read the current
LLM registry adapter source before generating Nacos instructions.

## Step 5 — Validate

Before booting Pixiu, inspect `conf.yaml` directly. Check yaml syntax,
filter order, required LLM cluster fields, `llm_meta` placement,
kvcache instance-id references, and Nacos adapter basics.

## Cross-Cutting Rules

### Always

- Generate `conf.yaml` fragments, not `api_config.yaml`.
- Use exact current Kind strings: `dgp.filter.llm.proxy`,
  `dgp.filter.llm.tokenizer`, `dgp.filter.ai.kvcache`, and
  `dgp.adapter.llmregistrycenter`.
- Use singular `llm-meta.api_key` in Nacos metadata. Older docs may say
  `llm-meta.api_keys`, but current source reads only the singular key.
- Put `dgp.filter.ai.kvcache` and `dgp.filter.llm.tokenizer` before
  `dgp.filter.llm.proxy` when enabled.
- Keep `llm_meta` under each `clusters[].endpoints[]` entry.
- Use endpoint `id` values that can be matched by kvcache / LMCache
  instance IDs.
- Use duration strings with units: `"60s"`, `"200ms"`, `"5m"`.

### Never

- Put LLM provider credentials in `api_config.yaml`.
- Put `llm_meta` under the cluster instead of endpoints.
- Put `scheme` under individual endpoints or put full URLs in
  `socket_address.domains`. The proxy has one filter-level scheme and
  endpoint domains must be host-only.
- Use Nacos `llm-meta.api_keys` plural for current source.
- Use old names such as `dgp.filter.http.llm.proxy` unless Step 0 proves
  the target branch changed the Kind constants.
- Enable kvcache without `vllm_endpoint`, `lmcache_endpoint`, and an
  endpoint-id matching story.
- Promise cost billing from config alone. Current tokenizer records
  token metrics; pricing/currency math belongs outside the current
  filter unless the target branch adds it.

## Source Files To Read

- `pkg/filter/llm/proxy/`
- `pkg/filter/llm/tokenizer/`
- `pkg/filter/ai/kvcache/`
- `pkg/adapter/llmregistry/`
- `pkg/model/llm.go`, `pkg/model/cluster.go`, and current AI config
  examples.
