# The Four (Plus Two) HTTP Filter Interfaces

Source of truth: `pkg/common/extension/filter/filter.go`. Read that file
first; this reference exists to annotate the shape and the "why".

## Interfaces, in implementation order

### 1. `HttpFilterPlugin`

```go
type HttpFilterPlugin interface {
    Kind() string                                   // unique identifier, e.g. "dgp.filter.http.cors"
    CreateFilterFactory() (HttpFilterFactory, error)
}
```

Every package has exactly one `Plugin` struct that implements this. It
is what gets registered in `init()` via `filter.RegisterHttpFilter`.
`Kind()` returns a package-level `const Kind = "..."` — define it once
at the top of the file and reuse everywhere.

### 2. `HttpFilterFactory`

```go
type HttpFilterFactory interface {
    Config() any
    Apply() error
    PrepareFilterChain(ctx *http.HttpContext, chain FilterChain) error
}
```

Factories are **listener-scoped singletons** (one per listener into
which the filter is wired). `Config()` returns a pointer so that
pixiu's config manager can deserialize yaml into it. `Apply()` runs
**after** deserialization — this is where you validate fields and set
defaults. `PrepareFilterChain` is called **per request**: construct a
fresh `Filter` instance, append it to the chain.

Two correctness rules that bit many contributors:

- `Config()` must return a **pointer**, not the struct by value. Yaml
  binding silently no-ops on a non-pointer.
- Do not share `factory.cfg` pointer with the per-request `Filter`.
  Copy out the fields you need so a hot-reload of the factory config
  does not race with in-flight requests.

### 3. `HttpDecodeFilter`

```go
type HttpDecodeFilter interface {
    Decode(ctx *http.HttpContext) FilterStatus
}
```

Invoked in **declaration order** as the request travels inbound. Return
`filter.Continue` to move on or `filter.Stop` to short-circuit. Common
decode work: authentication, rate limiting, path rewriting, enriching
`ctx.Params`, emitting request-size metrics.

### 4. `HttpEncodeFilter`

```go
type HttpEncodeFilter interface {
    Encode(ctx *http.HttpContext) FilterStatus
}
```

Invoked in **reverse declaration order** on the response path. Common
encode work: response header rewriting, body transformation, masking
PII, emitting response-size metrics.

A single filter package may implement either, neither (rare), or both.
If both, `PrepareFilterChain` calls both `chain.AppendDecodeFilters`
and `chain.AppendEncodeFilters`.

### 5. `NetworkFilterPlugin` (L4)

```go
type NetworkFilterPlugin interface {
    Kind() string
    CreateFilter(config any) (NetworkFilter, error)
    Config() any
}
```

Network filters sit one layer down — they own the raw connection and
dispatch per-protocol. `dgp.filter.httpconnectionmanager` is a network
filter; its `http_filters` list is where HTTP filters above attach.
Unless you are adding a wholly new L4 behavior, you do not need to
write a Network filter.

### 6. `NetworkFilter`

Large surface: `ServeHTTP`, `OnDecode`, `OnEncode`, `OnData`,
`OnTripleData`, `OnUnaryRPC`, `OnStreamRPC`, `Close`. Reuse an existing
implementation as a template; see
`pkg/filter/network/httpconnectionmanager/` and
`pkg/filter/network/grpcconnectionmanager/`.

## Reference filter walkthrough — CORS

`pkg/filter/cors/cors.go` is the shortest idiomatic HTTP filter. ~150
lines, single file, Decode-only. Pattern to copy:

```go
const Kind = constant.HTTPCorsFilter // "dgp.filter.http.cors"

func init() { filter.RegisterHttpFilter(&Plugin{}) }

type (
    Plugin        struct{}
    FilterFactory struct{ cfg *Config }
    Filter        struct{ cfg *Config }
    Config        struct { /* yaml-tagged fields */ }
)

func (p *Plugin) Kind() string                                   { return Kind }
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) { return &FilterFactory{cfg: &Config{}}, nil }

func (f *FilterFactory) Config() any { return f.cfg }
func (f *FilterFactory) Apply() error { return nil }
func (f *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
    inst := &Filter{cfg: f.cfg.DeepCopy()}
    chain.AppendDecodeFilters(inst)
    return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus { ... }
```

Add `Encode` and `chain.AppendEncodeFilters(inst)` when you need both
phases.

## Where to put `Kind` constants

Two patterns in the tree:

- **Shared registry** (`pkg/common/constant/`): most built-in filters
  put their kind names there so other packages can import them
  cleanly. Prefer this when the filter is being merged upstream.
- **Package-local const**: fine for private / experimental filters.
  Declare `const Kind = "dgp.filter.http.<name>"` in your file.

Be consistent within a PR — do not mix.
