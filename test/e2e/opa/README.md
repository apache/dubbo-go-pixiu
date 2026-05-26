# PR2: Admin OPA → Gateway OPA Full-Link E2E

This directory holds the PR2 end-to-end test artifacts. PR2's scope is
**verifying the full chain**: admin REST configures a policy, the OPA Server
holds it, and the gateway OPA filter consumes that same policy when making
allow/deny decisions on live HTTP traffic.

## Why PR2 lives mostly as one Go test file

The CI-runnable test is `admin/initialize/e2e_opa_test.go`. It stays in the
`initialize` package so it can reuse PR1's JWT signer (`signToken`),
multipart helper (`putMultipart`), and HTTP driver (`doReq`) without
duplicating them. The directory you are reading provides the *runner* and
*operator-facing docs*; the test code itself is one package away in the
admin tree so `go test ./admin/...` keeps picking it up.

## What the suite verifies

The harness stands up:

1. A **smart in-process OPA mock** (`regoMockOPA`) — uses the real
   `github.com/open-policy-agent/opa/rego` library to compile and evaluate
   modules, so policy semantics in tests match what a real OPA daemon would
   do. No docker, no etcd, no real OPA binary required.
2. **The real admin Gin router** via `initialize.Routers()` with
   `adminconfig.Bootstrap.OPA.ServerURL` pointed at the mock.
3. **The real gateway OPA filter** from `pkg/filter/opa` (via the public
   `Plugin.CreateFilterFactory()` API) pointed at the same mock URL.

Each test publishes a policy through the admin REST PUT, then drives one or
more HTTP requests through the gateway filter and asserts the decision.

| Test | Scenario | What it proves |
|---|---|---|
| `TestE2E_AllowedThroughFullChain` | PUT "allow if GET" → GET request | Admin → OPA → gateway end-to-end allow returns `filter.Continue` with no local reply |
| `TestE2E_DeniedThroughFullChain` | Same policy, POST request | Deny returns `filter.Stop` + 403, short-circuits before upstream |
| `TestE2E_DefaultDenyForAllRequests` | `default allow := false` only | 5-method matrix all denied (no rule shape can sneak past) |
| `TestE2E_PolicyHotReload` | PUT v1, then PUT v2 (no restart) | Behaviour flips on the very next request — the headline OPA-server-mode value |
| `TestE2E_DeleteCausesMissingResultFailClosed` | PUT then DELETE | Gateway returns 502 (matches `test_opa.md` §6.6) |
| `TestE2E_HeaderBasedAllowDeny` | Policy on `input.headers["X-Role"]` | Title-cased header propagation works (admin/user/missing variants) |
| `TestE2E_GatewayTimeoutFailClosed` | 200ms decision delay, 50ms gateway timeout | Returns 504, elapsed time bounded under 180ms |
| `TestE2E_PolicyIDOverrideRoutesThroughGateway` | PUT with form-level `policy_id` override | The override flag from PR1 propagates all the way to a working gateway decision |

## Running

```bash
# Default — all PR2 cases, no -v
./test/e2e/opa/run.sh

# Verbose
VERBOSE=1 ./test/e2e/opa/run.sh

# Subset by name
./test/e2e/opa/run.sh -run AllowedThroughFullChain

# Or directly:
go test -count=1 -run TestE2E_ -v ./admin/initialize/
```

Expected output (verbose):

```
=== RUN   TestE2E_AllowedThroughFullChain
--- PASS: TestE2E_AllowedThroughFullChain (0.03s)
=== RUN   TestE2E_DeniedThroughFullChain
--- PASS: TestE2E_DeniedThroughFullChain (0.01s)
=== RUN   TestE2E_DefaultDenyForAllRequests
--- PASS: TestE2E_DefaultDenyForAllRequests (0.01s)
=== RUN   TestE2E_PolicyHotReload
--- PASS: TestE2E_PolicyHotReload (0.01s)
=== RUN   TestE2E_DeleteCausesMissingResultFailClosed
--- PASS: TestE2E_DeleteCausesMissingResultFailClosed (0.01s)
=== RUN   TestE2E_HeaderBasedAllowDeny
--- PASS: TestE2E_HeaderBasedAllowDeny (0.01s)
=== RUN   TestE2E_GatewayTimeoutFailClosed
--- PASS: TestE2E_GatewayTimeoutFailClosed (0.21s)
=== RUN   TestE2E_PolicyIDOverrideRoutesThroughGateway
--- PASS: TestE2E_PolicyIDOverrideRoutesThroughGateway (0.01s)
PASS
ok      github.com/apache/dubbo-go-pixiu/admin/initialize       0.342s
```

## Reuse — no duplicated implementation

The suite **only consumes** existing code:

| Component reused | Where |
|---|---|
| Admin Gin router + all middleware | `admin/initialize/router.go` (PR1) |
| Admin OPA controllers | `admin/controller/opa/opa.go` (PR1) |
| Admin business logic (PUT/GET/DELETE) | `admin/logic/logic.go` (PR1) |
| OPA bootstrap config (timeout, defaults) | `admin/config/opa.go` (PR1) |
| Gateway OPA filter (server mode) | `pkg/filter/opa/opa.go` (existing) |
| JWT auth middleware + SignKey | `admin/controller/auth/auth.go` (existing) |
| Rego compile + eval | `github.com/open-policy-agent/opa/rego` (already in `go.mod`) |

No production code is touched by PR2. If any of the existing components
regress, this suite is the highest-signal alarm — a single failed E2E case
points at the seam that broke.

## What is *not* covered here (manual / out of CI scope)

- Running the real `openpolicyagent/opa` docker image to validate Rego v0 vs
  v1 compatibility nuances — see `test_opa.md` §5.
- End-to-end network performance + concurrency (`ab` benchmarking) — see
  `test_opa.md` §6.8.
- Wiring through the full `cmd/pixiu` gateway binary's listener stack —
  the OPA filter is a self-contained `HttpFilterFactory`, and exercising
  it directly (as PR2 does) is faster and just as discriminating for the
  decision logic.

If those manual scenarios regress, `test_opa.md` already documents the exact
curl / docker recipes to reproduce them.

## Dependency on PR1

PR1 added:
- `admin/controller/opa/opa.go`
- `admin/logic/logic.go` (the `BizPut/Get/DeleteOPAPolicy` helpers)
- `admin/initialize/router.go` mounts (the three `/config/api/opa/policy` routes)
- `admin/config/opa.go` (Bootstrap fields)

PR2 imports those packages and exercises them through the real
`initialize.Routers()` engine — so this suite cannot meaningfully run until
PR1 is merged. After PR1 is in, PR2 is mergeable as a single self-contained
addition (one Go test file + this runner directory).
