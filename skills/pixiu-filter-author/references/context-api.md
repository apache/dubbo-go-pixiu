# `HttpContext` — Field Cheat Sheet

Source: `pkg/context/http/context.go`. Read that file when this reference
looks out of date; the struct evolves.

`HttpContext` is the per-request object your filter receives in both
`Decode` and `Encode`. It threads the original request, the upstream
response, the routed-client response, params, headers, and the filter
chain position.

## The fields filters actually use

| Field / method | Lives in | Typical use |
|---|---|---|
| `ctx.Request` | inbound `*http.Request` | read headers, method, path, raw body |
| `ctx.Writer` | `http.ResponseWriter` | write headers or body **before** handing control to the router; usually for Decode that terminates early (CORS preflight, rate-limit 429) |
| `ctx.GetHeader(name)` | — | case-insensitive header read on the inbound request |
| `ctx.AddHeader(k, v)` / `ctx.RemoveHeader(k)` | — | mutate the request headers forwarded upstream |
| `ctx.Params` | `map[string]any` | route params + filter-decoded values (JWT claims, etc.) |
| `ctx.SourceResp` | upstream's response (raw) | what the client/adapter returned; available in Encode |
| `ctx.TargetResp` | what pixiu will ship to client | what the client will see; Encode may mutate |
| `ctx.Err` | `error` | set by upstream failures; Encode can decide to mask or expose |
| `ctx.Route` | `*model.RouteAction` | matched route, cluster name, timeouts |
| `ctx.API` | `*router.API` | api-config route, only populated when the api filter is in play |
| `ctx.Next()` | — | explicitly advance; rarely needed — returning `filter.Continue` is the idiomatic way |
| `ctx.Abort()` | — | stop subsequent Decode filters; Encode still runs for filters already on-chain |
| `ctx.AddFinishCallbacks(f)` | — | run `f` after response is sent; use for cleanup / metric emission |
| `ctx.WriteErr(...)`, `ctx.WriteWithStatus(...)` | — | respond from within the filter (bypass upstream) |

## `SourceResp` vs `TargetResp` — the number-one Encode bug

- **`SourceResp`** is the raw response from the upstream (HTTP backend,
  Dubbo generic invoke result, gRPC message). Treat as read-only
  input.
- **`TargetResp`** is what pixiu will serialize and send to the client.
  This is the mutable one for an Encode filter.

A very common mistake is to mutate `SourceResp` in Encode and expect
the client to see the change. It does not. The chain may have already
copied data into `TargetResp`. Always check which field the filter you
are writing should touch.

## Headers: inbound vs outbound

- Inbound headers (what the client sent): `ctx.Request.Header` and the
  `ctx.GetHeader` shortcut.
- Outbound headers to the upstream: modify `ctx.Request.Header` in
  Decode — proxies forward the same object.
- Response headers to the client: `ctx.Writer.Header()`.

Do not assume the three maps share keys.

## Params — the cross-filter scratch pad

`ctx.Params` (a `map[string]any`) is how filters pass information to
each other. Examples:

- JWT filter stores claims under `"user"`.
- Rate-limit filter reads `"user"` to pick a quota bucket.
- OPA filter reads `"user"` to build its input document.

If you invent a new key, document it in `references/` of whichever
skill publishes it. Do not rely on keys from other teams without
checking.

## Lifecycle in one diagram

```
client → ctx = new HttpContext
       → Decode[0] → Decode[1] → ... → Decode[N]
                                         ↓
                                       upstream (Dubbo / HTTP / gRPC)
                                         ↓
       ← Encode[0] ← Encode[1] ← ... ← Encode[N]   (reverse order)
       ← write response  ← finish callbacks
```

Every Decode filter that returns `Continue` implicitly opts into having
its Encode half run. Conversely, if `ctx.Abort()` is called at Decode[i],
Encode[j] for j < i still runs (because those filters already entered
the chain), but j > i does not.
