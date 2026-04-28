# `pkg/pluginregistry/registry.go` — The Single Boot Activator

This file exists because Go's `init()` only runs for packages that are
**imported somewhere in the build graph**. Filters register themselves
in `init()`, so unless something imports their package, the `init()`
never runs and the registration never happens. `pluginregistry` is the
one file that imports every filter/adapter/listener/lb package for side
effects — a single central blank-import list — and `cmd/pixiu/` imports
`pluginregistry`.

## Shape of the file

```go
package pluginregistry

import (
    _ "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry"
    _ "github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry"
    _ "github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver"
    // ...
    _ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/maglev"
    // ...
    _ "github.com/apache/dubbo-go-pixiu/pkg/filter/accesslog"
    _ "github.com/apache/dubbo-go-pixiu/pkg/filter/cors"
    // ...
    _ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/dubboproxy"
    _ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/httpproxy"
    // ...
    _ "github.com/apache/dubbo-go-pixiu/pkg/listener/http"
    // ...
)
```

All imports are in one block. Entries are **alphabetical within each
logical group** — adapters, lb, retry, filter (top-level),
filter/http/, filter/network/, filter/auth/, listener. When adding,
keep the alphabetical order *within* the group your package belongs to.

## How to add an import

For a new filter at `pkg/filter/mynewfilter/`:

```go
_ "github.com/apache/dubbo-go-pixiu/pkg/filter/mynewfilter"
```

Placement: in the top-level `pkg/filter/` group, alphabetically between
`_ "github.com/apache/dubbo-go-pixiu/pkg/filter/metric"` and
`_ "github.com/apache/dubbo-go-pixiu/pkg/filter/opa"` (or wherever
`mynewfilter` sorts).

For a new proxy-style filter at `pkg/filter/http/mynewproxy/`:

```go
_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/mynewproxy"
```

Placement: in the `pkg/filter/http/` group, alphabetically.

## What happens if you forget

- `go build ./...` — passes.
- `go test ./pkg/filter/mynewfilter/...` — passes.
- `./dubbo-go-pixiu gateway start -c configs/conf.yaml` — **fails at
  config load** with `no filter found for name dgp.filter.http.mynew`.

That boot-time error is the earliest signal. There is no compile-time
check because the blank import is literally an import for side effects
only.

## Automating the check

`scripts/gen-registry-import.sh` (in this skill) walks `pkg/filter/` and
`pkg/filter/http/`, grep's each package for `filter.RegisterHttpFilter`
or `filter.RegisterNetworkFilterPlugin`, and compares the resulting set
against the blank-import list. It prints the missing lines. It does NOT
edit the registry file — apply the diff by hand (this is a high-blast-
radius file; a human should read the insertion point).

Run it as the last step of every filter PR.

## Why not a magic auto-registration?

Two options have been proposed over the years:

1. `go generate` a registry file from package tags.
2. Reflect over package paths at runtime.

Both were rejected because `pluginregistry` is also the list of
**what gets built into the binary** — some deployments of pixiu ship a
cut-down registry for supply-chain or binary-size reasons. Keeping the
list explicit is a feature, not a bug.
