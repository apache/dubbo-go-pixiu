# Tokenizer Metrics

Current source anchor: `pkg/filter/llm/tokenizer/tokenizer.go`.

`dgp.filter.llm.tokenizer` records token and upstream metrics. It adds
both Decode and Encode filters: Decode captures the start time; Encode
reads the response and proxy attempt metadata.

## Config

```yaml
- name: dgp.filter.llm.tokenizer
  config:
    log_to_console: false
```

`log_to_console` is optional. Avoid enabling it in production examples
unless the user explicitly wants verbose output.

## Placement

Place tokenizer before `dgp.filter.llm.proxy`:

```yaml
http_filters:
  - name: dgp.filter.llm.tokenizer
  - name: dgp.filter.llm.proxy
```

This lets Decode start timing before the proxy performs the upstream
request. Encode later inspects unary and streaming responses.

## Metrics Emitted

Current metric names include:

- `pixiu_llm_prompt_tokens_total`
- `pixiu_llm_completion_tokens_total`
- `pixiu_llm_total_tokens_total`
- `pixiu_llm_upstream_requests_total`
- `pixiu_llm_upstream_requests_success_total`
- `pixiu_llm_upstream_requests_failure_total`
- `pixiu_llm_total_duration_microseconds_sum_total`
- `pixiu_llm_time_to_last_token_milliseconds_sum_total`
- `pixiu_llm_streaming_requests_total`

The tokenizer reads OpenAI-style `usage` fields for unary responses and
SSE `data:` frames for streaming responses. If an upstream omits token
usage, metrics may be partial.

## Not Cost Billing by Itself

The tokenizer measures tokens. It does not currently apply provider
pricing, currency conversion, tenant billing accounts, or quota
enforcement by itself. Do not promise cost billing unless the target
branch has additional billing code.
