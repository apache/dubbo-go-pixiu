# LLM Gateway Troubleshooting

## 404 or route not matched

- Check `route_config.routes[].match.prefix` or `path`.
- Confirm the client path includes the upstream path you expect; proxy
  forwards the inbound path.
- Confirm `dgp.filter.httpconnectionmanager` owns the `http_filters`
  list.

## 503 cluster not found

- `route.cluster` must match `static_resources.clusters[].name`.
- With Nacos discovery, ensure the instance metadata has the same
  `cluster` value.

## Upstream 401

- `llm_meta.api_key` may be missing or wrong.
- If the caller should provide Authorization, leave endpoint
  `api_key` empty; otherwise proxy overwrites Authorization with
  `Bearer <api_key>`.
- Never paste production keys into generated examples.

## Fallback does not happen

- `llm_meta.fallback` must be `true` on the failing endpoint.
- There must be another endpoint in the same cluster.
- Retry policy errors can skip an endpoint before a request is made.

## Token metrics missing

- Ensure `dgp.filter.llm.tokenizer` is in `http_filters`.
- Put tokenizer before `dgp.filter.llm.proxy`.
- Some upstreams omit OpenAI-style `usage` data; streaming metrics rely
  on parseable SSE `data:` frames.

## KV cache route hint ignored

- LMCache returned an `instance_id` that does not equal any endpoint
  `id`.
- `dgp.filter.ai.kvcache` was placed after proxy.
- `lookup_routing_timeout` is too low for the LMCache controller.
- `enabled` is false or `vllm_endpoint` / `lmcache_endpoint` is empty.
