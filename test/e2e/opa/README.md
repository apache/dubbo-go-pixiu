# Admin OPA → Gateway OPA Full-Link E2E

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
| `TestE2E_PolicyIDOverrideRoutesThroughGateway` | PUT with form-level `policy_id` override 

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


