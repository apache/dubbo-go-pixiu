---
name: pixiu-llm-gateway
description: Creates and validates dubbo-go-pixiu LLM gateway conf.yaml. Use for LLM proxy/tokenizer/kvcache filters, llm_meta, vLLM, LMCache, retry/fallback, or Nacos LLM discovery. Do not use for MCP gateway config.
---

# pixiu-llm-gateway

## Purpose
Generate a complete or embeddable LLM gateway `conf.yaml` that configures Pixiu as a multi-provider HTTP proxy.

## When to use
- Use when:
  - Generating or fixing Pixiu LLM gateway `conf.yaml`, including LLM proxy/tokenizer/kvcache, `llm_meta`, retry/fallback, vLLM/LMCache, or Nacos LLM discovery.

- Do not use for:
  - MCP gateway, HTTP-to-Dubbo routes, or implementing new LLM-related filters.

## Inputs
- Choose upstream mode (required):
  - Required:
    - `upstream_mode`: `static` or `registry`

- `listener` (required):
  - Defaults:
    - `address.socket_address.address: 0.0.0.0`
    - `address.socket_address.port: 8888`

- `route_config` (required):
  - Defaults:
    - `routes[].match.prefix: /v1`
    - `routes[].route.cluster: llm`

- `dgp.filter.llm.proxy` (required):
  - Defaults:
    - `config.scheme: http` (use `https` when the upstream is an HTTPS endpoint)
    - `config.timeout: 60s`
    - `config.maxIdleConns: 100`
    - `config.maxIdleConnsPerHost: 100`
    - `config.maxConnsPerHost: 100`

- Static upstream (required when `upstream_mode: static`):
  - Required:
    - `endpoints[].ID`
    - `endpoints[].socket_address.address` or `endpoints[].socket_address.domains`
    - `endpoints[].socket_address.port`
  - Optional:
    - `endpoints[].llm_meta.provider`
    - `endpoints[].llm_meta.api_key`
    - `endpoints[].llm_meta.retry_policy.config`
  - Defaults:
    - `clusters[].name: llm`
    - `clusters[].lb_policy: RoundRobin`
    - `endpoints[].llm_meta.retry_policy.name: NoRetry`
    - `endpoints[].llm_meta.fallback: false`
    - `endpoints[].llm_meta.health_check_interval: 5000`

- Registry upstream (required when `upstream_mode: registry`):
  - Required:
    - `registries.nacos.address`
  - Optional:
    - `registries.nacos.group`
    - `registries.nacos.namespace`
    - `registries.nacos.username`
    - `registries.nacos.password`
  - Defaults:
    - `adapters[].id: llm-registry`
    - `adapters[].name: dgp.adapter.llmregistrycenter`
    - `registries.nacos.protocol: nacos`
    - `registries.nacos.timeout: 5s`
    - `registries.nacos.group: DEFAULT_GROUP`

- `dgp.filter.llm.tokenizer` (optional):
  - Defaults:
    - `config.log_to_console: false`

- `dgp.filter.ai.kvcache` (optional):
  - Required:
    - `config.enabled: true` (must be set explicitly when enabling kvcache)
    - `config.vllm_endpoint`
    - `config.lmcache_endpoint`
  - Optional:
    - `config.default_model`
    - `config.token_cache`
    - `config.cache_strategy`
  - Defaults:
    - `config.request_timeout: 2s`
    - `config.lookup_routing_timeout: 50ms`
    - `config.hot_window: 5m`
    - `config.hot_max_records: 300`
    - `config.max_idle_conns: 100`
    - `config.max_idle_conns_per_host: 100`
    - `config.max_conns_per_host: 100`
    - `config.retry.max_attempts: 3`
    - `config.retry.base_backoff: 100ms`
    - `config.retry.max_backoff: 2s`
    - `config.circuit_breaker.failure_threshold: 5`
    - `config.circuit_breaker.recovery_timeout: 10s`
    - `config.circuit_breaker.half_open_max_calls: 2`

## Workflow
1. Read current source before generating YAML:
   - `pkg/common/constant/key.go`.
   - `pkg/filter/llm/proxy/filter.go`.
   - `pkg/filter/llm/tokenizer/tokenizer.go`.
   - `pkg/filter/ai/kvcache/config.go` and `pkg/filter/ai/kvcache/handlers.go`.
   - `pkg/model/llm.go`, `pkg/model/cluster.go`, and `pkg/model/base.go`.
2. Guide and read user input:
   1. Read existing config and already provided user information first; do not ask again for information already present or stated.
   2. Ask one `Inputs` group at a time. Each time, output only that group's required fields, with a short explanation after each field.
   3. If the current group's required fields are incomplete, ask only for the missing fields and do not move to the next group.
   4. After the current group's required fields are complete, ask whether to fill that group's optional fields; if yes, list those optional fields with short explanations.
   5. After optional fields are skipped or completed, apply that group's defaults and output default-value information; defaults must not override existing config or user input.
3. Choose static or registry path by `upstream_mode`: static uses `static_resources.clusters[]`; registry uses the LLM registry adapter.
4. Generate listener, route, and filters: use HCM to carry the LLM route and HTTP filters; add `dgp.filter.ai.kvcache`, tokenizer, and `dgp.filter.llm.proxy` as needed.
5. Generate LLM upstream: static endpoints go under `socket_address` and `llm_meta`; ensure the LLM cluster does not mix in ordinary HTTP endpoints.

## Output format
- Show the relevant YAML fragments.

## Validation
- Verify LLM filter-chain order: `dgp.filter.ai.kvcache` -> `dgp.filter.llm.tokenizer` -> `dgp.filter.llm.proxy`; ignore missing filters within this order chain.
- Verify `scheme` is under `dgp.filter.llm.proxy.config`, not under an endpoint; `socket_address.domains` contains host names only, such as `api.openai.com`, not full URLs or paths.
- When KV cache is enabled, verify the LMCache-side `instance_id` exactly matches the Pixiu endpoint `ID`.

## Examples
Complete LLM route example (`conf.yaml`):

```yaml
static_resources:
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
                      prefix: /v1
                    route:
                      cluster: llm
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
                    hot_max_keys: 1000
                    max_idle_conns: 100
                    max_idle_conns_per_host: 100
                    max_conns_per_host: 100
                    token_cache:
                      enabled: true
                      max_size: 1024
                      ttl: "10m"
                    cache_strategy:
                      enable_compression: true
                      enable_pinning: true
                      enable_eviction: true
                      memory_threshold: 0.85
                      hot_content_threshold: 10
                      load_threshold: 0.7
                      pin_instance_id: "vllm-instance-1"
                      pin_location: "LocalCPUBackend"
                      compress_instance_id: "vllm-instance-1"
                      compress_location: "LocalCPUBackend"
                      compress_method: "zstd"
                      evict_instance_id: "vllm-instance-1"
                    circuit_breaker:
                      failure_threshold: 5
                      recovery_timeout: "10s"
                      half_open_max_calls: 2
                    retry:
                      max_attempts: 3
                      base_backoff: "100ms"
                      max_backoff: "2s"
                - name: dgp.filter.llm.tokenizer
                  config:
                    log_to_console: false
                - name: dgp.filter.llm.proxy
                  config:
                    scheme: http
                    timeout: "60s"
                    maxIdleConns: 100
                    maxIdleConnsPerHost: 100
                    maxConnsPerHost: 100
  clusters:
    - name: llm
      lb_policy: RoundRobin
      endpoints:
        - ID: vllm-instance-1
          socket_address:
            address: 127.0.0.1
            port: 8000
          llm_meta:
            provider: vllm
            api_key: "<api-key>"
            fallback: false
            health_check_interval: 5000
            retry_policy:
              name: ExponentialBackoff
              config:
                times: 3
                initialInterval: "200ms"
                maxInterval: "5s"
                multiplier: 2.0
```

Nacos LLM registry mode example (`conf.yaml`):

```yaml
static_resources:
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
                      prefix: /v1
                    route:
                      cluster: llm
              http_filters:
                - name: dgp.filter.llm.proxy
                  config:
                    scheme: http
                    timeout: "60s"
  adapters:
    - id: llm-nacos
      name: dgp.adapter.llmregistrycenter
      config:
        registries:
          nacos:
            protocol: nacos
            address: "127.0.0.1:8848"
            timeout: "5s"
            group: DEFAULT_GROUP
            namespace: public
```
