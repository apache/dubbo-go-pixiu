---
name: pixiu-filter-author
description: |
  Create a new HTTP or Network filter for Apache dubbo-go-pixiu. Use whenever
  the user wants to add a custom filter, extend gateway per-request behavior,
  implement a new auth/logging/transformation/proxy filter, or mentions
  "pixiu filter", "dgp.filter", "http_filters", "filter chain", "custom filter",
  "extend pixiu", or is porting a filter from Envoy/Kong/APISIX to pixiu.
  Strongly prefer this skill over ad-hoc generation even when the user does
  not say the word "skill" — getting the SPI registration right is the single
  biggest source of "my filter does not run" bug reports.
allowed-tools: [Read, Grep, Glob, Edit, Write, Bash]
metadata:
  version: "0.1.2"
  domain: extension
  scope: implementation
  triggers: ["filter", "pixiu filter", "dgp.filter", "http_filters", "filter chain", "custom filter", "extend pixiu"]
  pixiu_min_version: "0.6.0"
  role: specialist
---

# pixiu-filter-author — Authoring HTTP / Network Filters

dubbo-go-pixiu's extension surface is Envoy-style: every filter is a Go
package that registers a `Plugin` in `init()`, and a central file
(`pkg/pluginregistry/registry.go`) blank-imports all such packages to
activate them. A single missing blank import is the most common reason a
newly written filter "does nothing" — so this skill treats registration as
a first-class concern, not an afterthought.

## When to Use

Use this skill when the user wants to:
- Add a new HTTP filter (auth, logging, header rewriting, rate limiting,
  request transformation, metric emission, etc.).
- Add a new Network filter (e.g. a custom connection manager).
- Port a filter from Envoy / Kong / APISIX onto pixiu.
- Fix an existing filter that "isn't running" — start at Step 5 (blank
  import) and Step 3 (phase).

**Do NOT use this skill for:**
- Editing the yaml that *uses* an existing filter → that is plain config
  work. Only include a minimal mounting snippet when you create a new
  filter in this skill.
- Mounting an already-registered filter such as CORS, JWT, OPA, MCP,
  LLM proxy, or Dubbo proxy. In that case do not create new filter code;
  do not touch `pkg/pluginregistry/registry.go`; answer with the
  smallest `http_filters` snippet and state that the filter already
  exists.
  If pressured to add a package anyway, explicitly answer:
  "this is already-registered; do not create new filter code; do not
  touch `pkg/pluginregistry/registry.go`."
- Adding a registry adapter (Nacos, ZK, Consul), a new protocol
  listener (WebSocket, MQTT), or a custom load-balancing algorithm —
  these are different SPI hubs and out of scope for this skill. Point
  the user to the relevant interface file
  (`pkg/common/extension/adapter/adapter.go`,
  `pkg/listener/listener.go`, or
  `pkg/cluster/loadbalancer/`) and let them work directly.

## Prerequisites

- pixiu ≥ 0.6.0. Verify with `grep dubbo-go-pixiu go.mod` in Step 0.
- The user should be inside a clone of `github.com/apache/dubbo-go-pixiu`
  (or a fork). If they are not, stop and ask — this skill edits repo files.

## Steps

### Step 0 — Verify Context (ALWAYS FIRST, no exceptions)

Before writing a single line of code, read the current shape of three
files. pixiu's SPI is stable but not frozen, and your instinct about
these interfaces can be one minor version stale:

1. `pkg/common/extension/filter/filter.go` — the four interfaces
   (`HttpFilterPlugin`, `HttpFilterFactory`, `HttpDecodeFilter`,
   `HttpEncodeFilter`), plus `NetworkFilterPlugin` / `NetworkFilter`.
   Read signatures, especially `PrepareFilterChain`.
2. `pkg/pluginregistry/registry.go` — the blank-import list you will
   extend in Step 5. Note the alphabetical order within groups.
3. `pkg/filter/cors/cors.go` — the shortest idiomatic HTTP filter in the
   tree. A single file, all four types (`Plugin`, `FilterFactory`,
   `Filter`, `Config`) declared together. Note: `cors` ships **without**
   a `_test.go` — that matches roughly half of the existing filter
   packages. See Step 7 for the project's actual convention.

Then `ls pkg/filter/` and `ls pkg/filter/http/` to see the existing
directory conventions. Do not assume — layouts evolve (e.g. `cors` lives
at `pkg/filter/cors/`, while proxy filters live at
`pkg/filter/http/httpproxy/`, `pkg/filter/http/dubboproxy/`).

If an interface detail is unclear, stay in the source file instead of
relying on memory.

### Step 1 — Clarify the Filter Shape (STOP and ask)

Do not write code until the user has answered, in plain words:

1. **Phase**: Decode (request, before upstream), Encode (response, after
   upstream), or both?
2. **Kind**: the string that identifies the filter in yaml. It MUST start
   with `dgp.filter.http.` for HTTP filters or
   match Pixiu's current Network filter constants, usually
   `dgp.filter.network.`. The built-in HTTP connection manager is the
   special network-filter Kind `dgp.filter.httpconnectionmanager`.
   Suggest a name after checking `pkg/common/constant/key.go`; confirm
   the user is happy.
3. **Config fields**: what yaml keys does the user want? (`allow_origin`,
   `max_request_bytes`, etc.)
4. **Where to put the package**: the convention is
   `pkg/filter/<name>/` for ordinary filters and
   `pkg/filter/http/<name>/` for *proxy-style* filters (the ones that
   actually call upstream). If unsure, ask — getting this right once
   beats moving files later.

If the user says "just write something reasonable", pick a minimal
request-logging filter as the placeholder and be explicit about the
assumptions in a short summary before coding.

### Step 2 — Create the Package Skeleton

Create the directory and initial files. Two shapes are idiomatic in the
existing tree:

- **Single file** (cors, csrf, jwt, header, host, tracing): everything
  lives in `<name>.go`. This is the dominant pattern for non-proxy
  filters — about 20 of the existing filter packages follow it.
- **Split files** (network/dubboproxy, mcp/mcpserver, ai/kvcache, etc.):
  `plugin.go`, `filter.go`, `config.go`. Use this only when the package
  has more than ~200 LOC or several distinct types worth separating.

Test files (`*_test.go`) are **not** part of either default shape —
they are added on a per-filter basis when the logic warrants it. See
Step 7.

Create the initial files directly, using the closest in-tree filter as
the pattern. Keep the first version minimal: the four interface types,
`Kind`, `init()` registration, `Config`, and empty `Decode` / `Encode`
methods as appropriate. The user still owns every Config field.

### Step 3 — Implement the Four Interfaces

The four types for an HTTP filter, in order:

```go
type Plugin struct{}
func (p *Plugin) Kind() string { return Kind }
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
    return &FilterFactory{cfg: &Config{}}, nil
}

type FilterFactory struct{ cfg *Config }
func (f *FilterFactory) Config() any { return f.cfg }
func (f *FilterFactory) Apply() error { /* validate / defaults */ return nil }
func (f *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
    inst := &Filter{cfg: f.cfg.DeepCopy()} // or a shallow copy of the pieces you need
    chain.AppendDecodeFilters(inst)        // or AppendEncodeFilters, or both
    return nil
}

type Filter struct{ cfg *Config }
func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus { return filter.Continue }
// and/or:
func (f *Filter) Encode(ctx *http.HttpContext) filter.FilterStatus { return filter.Continue }
```

A few subtle rules the interface alone does not enforce:

- **Do not reuse `factory.cfg` directly in the Filter instance.** The
  factory's config pointer can be hot-reloaded at runtime. Copy the
  fields you need into the Filter when `PrepareFilterChain` runs. The
  CORS filter's `factory.cfg.DeepCopy()` pattern is the canonical
  reference.
- **Decode vs Encode is about the phase, not the direction.**  Mutating
  `ctx.TargetResp` (the outbound body) from a Decode filter is a
  category error — read `pkg/context/http/` before touching response
  state.
- **Response-body mutation is an Encode-only checklist.** For redaction,
  watermarking, compression, or similar response transforms, inspect
  existing response-transforming filters and their tests.
  Operate on `ctx.TargetResp` only after confirming the concrete
  response type. Handle `*client.UnaryResponse` deliberately, pass
  streaming responses through unless the user explicitly asked for
  streaming support, gate by `Content-Type` when the behavior is textual,
  and write table-driven tests for empty body, unsupported content type,
  configured/default values, and invalid config.

For Network filters, the interface surface is larger
(`ServeHTTP`, `OnData`, `OnTripleData`, etc.); reuse an existing one
(`pkg/filter/network/httpconnectionmanager/`) as a template.

### Step 4 — Register in `init()`

```go
func init() {
    filter.RegisterHttpFilter(&Plugin{})
}
```

For Network filters the call is `filter.RegisterNetworkFilterPlugin(&Plugin{})`.

### Step 5 — Add the Blank Import (THIS is where filters die silently)

Edit `pkg/pluginregistry/registry.go` and insert, in alphabetical order
within the correct grouping:

```go
_ "github.com/apache/dubbo-go-pixiu/pkg/filter/<name>"
```

(Or `.../pkg/filter/http/<name>` if the package lives there.)

If this line is missing, the filter compiles fine, `go test` passes, and
pixiu boots — but the filter never registers, and yaml using its Kind
fails with `no filter found for name ...`. This is *the* pixiu rite of
passage.

In answers, spell this out: package-level tests can pass while the
gateway binary never imports the package, so the plugin `init()` never
registers. The missing blank import is still required.

Manually check the new package against `pkg/pluginregistry/registry.go`:
if the package calls `filter.RegisterHttpFilter` or
`filter.RegisterNetworkFilterPlugin`, it must be blank-imported there.
Review the existing grouping and alphabetical order before inserting the
new line.

If the user asked for a "full workflow", do not stop after writing the
filter package. The final answer must include either the applied
`registry.go` edit or a precise patch hunk the user can apply.

### Step 6 — Wire Up Config

Add your filter to `configs/conf.yaml` (or whatever bootstrap the user
is running), under the `http_filters` list of the
`dgp.filter.httpconnectionmanager` network filter. Order matters: CORS /
request auth / tracing normally run before proxy-style filters, while
response-shaping filters must be placed where their Encode phase will
see the intended response.

```yaml
- name: dgp.filter.http.<name>
  config:
    <field>: <value>
```

Do NOT put the filter under `static_resources.listeners[].filter_chains[].filters[]`
directly — that's the Network filter level. HTTP filters always live
inside `dgp.filter.httpconnectionmanager`'s `http_filters`.

### Step 7 — Tests (optional; follow the project's actual convention)

Pixiu does **not** have a "every filter must have tests" policy. Decide
from the behavior you add, then mirror a nearby filter with comparable
complexity:

- **Skip the test** when your filter is a thin config-driven header
  rewriter, an auth check that delegates to a library, a one-line proxy
  switch, or otherwise mostly setup — matches cors / csrf / jwt /
  httpproxy.
- **Add a table-driven test** when your filter has branching logic,
  state transitions, or response-body transformation — matches
  sentinel/ratelimit / opa / accesslog.

If you choose to test, copy the style of nearby table-driven tests in
`pkg/filter/`. If you choose not to test, do not generate an empty stub
`*_test.go` — empty test files are noise and do not match repo style.

When in doubt, ask the user. Either answer is consistent with the
project; pick deliberately.

### Step 8 — Verify End-to-End

- `go build ./...` from repo root must succeed.
- If you wrote tests: `go test ./pkg/filter/<name>/...` must pass.
- If the user has the gateway binary, `./dubbo-go-pixiu gateway start -c
  configs/conf.yaml` should log the filter's Kind among the registered
  filters at boot.
- A `curl` against the listener port should exhibit the new behavior.

For full-workflow answers, emit a compact delivery checklist:

- Filter package path and the `Kind` string.
- `pkg/pluginregistry/registry.go` blank-import patch or applied edit.
- Minimal `http_filters` yaml snippet mounting the new filter.
- Test decision: either real test files and commands, or a short reason
  why no `_test.go` is consistent with nearby pixiu filters.
- Verification commands run or, if not runnable, the exact commands the
  user should run.

## Cross-Cutting Rules

### Always

- Name the `Kind` constant with the correct prefix. HTTP filters use
  `dgp.filter.http.*`; Network filters must match the current constants
  in `pkg/common/constant/key.go` such as `dgp.filter.network.*` or the
  special `dgp.filter.httpconnectionmanager`. Put it in a `const` at
  the top of the package.
- Use the project's logger: `import ".../pkg/logger"` and
  `logger.Infof(...)`. Never `fmt.Println`, never stdlib `log`.
- Deep-copy config into the per-request `Filter` instance — the factory
  config is live and can change under you.
- Keep imports in three alphabetically ordered groups (stdlib, third
  party, internal) — this is a project-wide convention and the CI
  linter enforces it.
- Match the surrounding package style on testing: look at neighbors of
  comparable complexity (cors / jwt for simple, opa / sentinel for
  branching) and follow their pattern.
- In full-workflow mode, always produce registry and config artifacts.
  A filter package alone is incomplete even when it compiles.

### Never

- Expose internal error details in HTTP responses. Return generic
  client-facing errors and log full diagnostics server-side only.
- Skip the blank import. This is the #1 filter bug. If something does
  not seem to register, re-check this before anything else.
- Accept `dgp.filter.networkfilter.*`. That prefix is invalid in current
  Pixiu; HTTP filters use `dgp.filter.http.*`, and Network filter Kinds
  must be checked against `pkg/common/constant/key.go`.
- Mutate `ctx.TargetResp` from a Decode filter. Use Encode.
- Rewrite response streams unless the user explicitly requested
  streaming support and you have designed backpressure/lifetime handling.
- Launch a goroutine inside Decode/Encode without a parent context and
  timeout — it will outlive the request.
- Return `nil, error` from `CreateFilterFactory` for a non-fatal config
  problem — instead return a zero-valued factory and log a warning; this
  keeps pixiu bootable.
- Generate an empty stub `*_test.go` "to be filled in later". Either
  write a real test or omit the file.

## Common Pitfalls

1. **"My filter doesn't run" = missing blank import.** 90% of the time.
   Check `pkg/pluginregistry/registry.go` first.
2. **Phase confusion.** Auth → Decode; response shaping → Encode.
   Filters are invoked Decode forward, Encode reversed — draw it out
   once and the rest follows.
3. **`ctx.Abort` vs `ctx.Next`.** `Abort` stops subsequent Decode
   filters, but Encode filters of already-entered filters still run.
4. **Per-request state in the factory.** The factory is shared across
   all requests on a listener; per-request state belongs in the Filter
   struct created by `PrepareFilterChain`.
5. **Dubbo-aware filters running too early.** Anything that needs the
   request's Dubbo parameters parsed must run *after*
   `dgp.filter.httpconnectionmanager` has done its work — effectively,
   it must be inside its `http_filters` list.
6. **Yaml tags omitted on Config struct.** Without
   `yaml:"field_name" mapstructure:"field_name"`, your field will not
   bind. Follow CORS's tag pattern exactly.

## Source Files To Read

- `pkg/common/extension/filter/filter.go`
- `pkg/pluginregistry/registry.go`
- `pkg/filter/cors/cors.go`
- Nearby filters under `pkg/filter/` or `pkg/filter/http/` with similar
  complexity.
