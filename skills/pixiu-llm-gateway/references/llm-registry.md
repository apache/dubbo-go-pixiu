# LLM Registry Discovery

Current source anchors:

- `pkg/adapter/llmregistry/registrycenter.go`
- `pkg/adapter/llmregistry/registry/nacos/listener.go`
- `docs/ai/registry.md`

Use `dgp.adapter.llmregistrycenter` only when LLM services publish
endpoint metadata into Nacos. Static `clusters[]` are simpler for one
or two stable upstreams.

## Adapter Config

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

## Required Nacos Instance Metadata

Each LLM instance should publish metadata such as:

```json
{
  "cluster": "deepseek_cluster",
  "id": "deepseek-primary",
  "name": "DeepSeek primary",
  "address": "api.deepseek.com",
  "llm-meta.fallback": "true",
  "llm-meta.api_key": "key-placeholder",
  "llm-meta.retry_policy.name": "ExponentialBackoff",
  "llm-meta.retry_policy.config": "{\"times\":3,\"initialInterval\":\"200ms\",\"maxInterval\":\"5s\",\"multiplier\":2.0}"
}
```

Important details:

- `cluster` groups endpoints into a Pixiu cluster.
- `id` must be unique within the cluster.
- `address` can override Nacos IP/port when the gateway should call a
  public provider host.
- Metadata keys use hyphenated `llm-meta.*`, not yaml `llm_meta`.
- Current docs mention both `llm-meta.api_key` and older typo-like
  `llm-meta.api_keys`; current source reads `llm-meta.api_key`.

## Dynamic vs Static Config

Even with registry discovery, the listener route and LLM filters still
live in `conf.yaml`. Nacos supplies endpoints, not the entire HTTP
listener.
