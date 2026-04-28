# LLM Endpoint and Proxy Configuration

Current source anchors:

- `pkg/common/constant/key.go`: exact Kind strings.
- `pkg/filter/llm/proxy/filter.go`: proxy behavior.
- `pkg/model/llm.go`: `llm_meta` fields.
- `pkg/model/cluster.go`: endpoint layout.
- `pkg/model/base.go`: `socket_address.domains` and `GetAddress()`.

## Filter Kinds

Use these current Kind strings:

- `dgp.filter.llm.proxy`
- `dgp.filter.llm.tokenizer`
- `dgp.filter.ai.kvcache`

Do not invent `dgp.filter.http.llm.*` names for current Pixiu.

## Static Cluster Example

```yaml
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
    lb_policy: "lb"
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

`dgp.filter.llm.proxy` reuses the inbound request path and query. If
clients call `/v1/chat/completions`, the upstream sees
`/v1/chat/completions`.

Keep `socket_address.domains` host-only. Current proxy code puts
`SocketAddress.GetAddress()` into `url.URL.Host`, so values such as
`https://api.openai.com` or `api.openai.com/v1` are not safe static
endpoint addresses.

## Endpoint `llm_meta`

`llm_meta` belongs under each endpoint, not under the cluster:

| Field | Meaning |
|---|---|
| `provider` | Human/metadata label for logs and registry conversion |
| `api_key` | Injected into upstream `Authorization: Bearer ...` |
| `retry_policy.name` | `NoRetry`, `CountBased`, or `ExponentialBackoff` |
| `retry_policy.config` | Policy-specific config map |
| `fallback` | Whether to pick the next endpoint after failure |
| `health_check_interval` | Millisecond cooldown before retrying unhealthy endpoint |

## Retry Policies

```yaml
retry_policy:
  name: "CountBased"
  config:
    times: 2
```

```yaml
retry_policy:
  name: "ExponentialBackoff"
  config:
    times: 3
    initialInterval: "200ms"
    maxInterval: "5s"
    multiplier: 2.0
```

Use `NoRetry` or omit retry config for a single attempt.

## API Key Handling

The proxy replaces the upstream Authorization header when
`llm_meta.api_key` is non-empty. For production examples, prefer
placeholders such as `${OPENAI_API_KEY}` and document how the deployment
system expands them. Do not paste real keys into generated configs.
