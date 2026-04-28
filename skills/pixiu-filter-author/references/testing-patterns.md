# Testing a pixiu HTTP Filter (Optional)

Tests in `pkg/filter/` are not uniform. As of pixiu 0.6, roughly half
of the filter packages ship without `_test.go` — see SKILL.md Step 7
for the exact split. This reference exists for the case **when you have
decided** to write a test, so the new test fits the patterns the
existing ones use.

If your filter is a thin config-shaped wrapper (cors / csrf / jwt
style), this file is not for you — skip the test entirely.

## When a test is worth writing

- The filter has branching logic on request shape (rate-limit decisions,
  OPA policy results, conditional rewrites).
- The filter mutates response bodies based on upstream output.
- The filter has Apply()-time validation that should reject bad
  configs.
- You are debugging a specific issue and want a regression test
  pinning the fix.

If none of these apply, omit the test — the existing tree treats the
boot-time integration check (filter loads, request flows) as sufficient
for trivial filters.

## The standard table-driven skeleton

Use this layout when you do write a test. It mirrors the patterns used
in `pkg/filter/sentinel/ratelimit/`, `pkg/filter/opa/`, and
`pkg/filter/accesslog/`.

```go
package mynewfilter

import (
    "net/http"
    "net/http/httptest"
    "testing"

    pixiuHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
    "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
    "github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
    tests := []struct {
        name       string
        reqHeaders map[string]string
        cfg        *Config
        wantStatus filter.FilterStatus
        wantHeader map[string]string
    }{
        {
            name:       "no-op when feature disabled",
            reqHeaders: map[string]string{},
            cfg:        &Config{Enabled: false},
            wantStatus: filter.Continue,
        },
        {
            name:       "adds tracing header when enabled",
            reqHeaders: map[string]string{},
            cfg:        &Config{Enabled: true},
            wantStatus: filter.Continue,
            wantHeader: map[string]string{"X-Trace": "yes"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            ctx := newTestCtx(tc.reqHeaders)
            f := &Filter{cfg: tc.cfg}

            got := f.Decode(ctx)
            assert.Equal(t, tc.wantStatus, got)
            for k, v := range tc.wantHeader {
                assert.Equal(t, v, ctx.Request.Header.Get(k))
            }
        })
    }
}

// newTestCtx builds a minimal HttpContext that is safe to exercise
// for Decode-phase tests. For Encode-phase tests, also populate
// ctx.SourceResp / ctx.TargetResp.
func newTestCtx(headers map[string]string) *pixiuHttp.HttpContext {
    req := httptest.NewRequest(http.MethodGet, "/probe", nil)
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    rec := httptest.NewRecorder()

    ctx := &pixiuHttp.HttpContext{
        Request: req,
        Writer:  rec,
        Params:  map[string]any{},
    }
    return ctx
}
```

`testify/assert` is the project default — pixiu's existing tests use
it consistently.

## Coverage expectations

Pixiu's CI does not gate on a coverage number. The expectation in
review is "tests cover the branches the code introduced" — not "every
file has a test". Aim accordingly.

## What to avoid

- **Do not test `PrepareFilterChain` by itself.** It only appends to
  the chain; the important behavior is in `Decode` / `Encode`.
- **Do not test the blank import from Go test code.** Boot-time
  smoke tests catch that; unit tests cannot.
- **Do not pull in `pkg/server` or a real listener** to unit test a
  filter — build the context by hand. If you find yourself wanting a
  full server, you are writing an e2e test, which belongs elsewhere
  (`integration/` if the user is targeting the upstream repo).
- **Do not ship an empty test file as a placeholder.** Either include
  real tests or omit the file.

## Helpers worth stealing

- `httptest.NewRequest` + `httptest.NewRecorder` are enough for most
  Decode tests.
- For response-body assertions in Encode, put the bytes directly into
  `ctx.TargetResp` and call `f.Encode(ctx)`.
- For filters that read `ctx.Params["user"]`, inject the expected
  payload into the map before invoking — do not build a JWT.

## Summary

If you are not sure whether to test: look at the closest existing
filter (similar complexity, similar shape) and copy its decision. The
project's actual convention is "match the neighbors", not "always test"
or "never test".
