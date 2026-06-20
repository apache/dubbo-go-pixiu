# MCP Intelligent Tool Router — Runnable Example

This directory demonstrates the MCP intelligent tool router.

## Files

- `router.yaml` — a complete Pixiu gateway config with the `router` block present,
  which enables tool governance: 4 tools, 1 unconditional policy rule (deny `internal`/`admin` tags), 2 workflow
  bundles, sampled decision logs, and decision detail logging for denied-tool samples.
- `demo.sh` — a curl walkthrough: initialize → tools/list (trimmed) →
  tools/call (denied vs allowed).

## What it shows

1. **`tools/list` trimming** — `internal_dump` (tagged `internal`/`admin`) is
   filtered out of the listing by the unconditional `block-privileged` policy rule.
2. **`tools/call` enforcement** — calling `internal_dump` is denied because it
   is not in the session's plan, even though the client knows its name.
3. **Allowed call** — `search_kb` is in the plan and is forwarded to the backend.
4. **Session-bound behavior** — `initialize` creates a fresh session; later
   `tools/list` and `tools/call` reuse that session instead of accepting a
   caller-supplied session on initialize.
5. **Notification semantics** — the server advertises
   `tools.listChanged=true`; clients do not declare this capability. Progressive
   or dynamic catalog changes are reported with `notifications/tools/list_changed`
   when the visible tool set changes.

The router default is `fail_closed`. This example explicitly sets
`fallback: bundle_default`, so no workflow match or an internal router error
falls back to the `safe-minimal` bundle after hard policy is applied. Explicit
policy denial still returns no tools.

## Running

```bash
# 1. Start a backend on :8081 (any HTTP server returning 200 works for the demo).
#    For example: python3 -m http.server 8081

# 2. Start Pixiu with this config.
dubbo-go-pixiu gateway start -c docs/ai/mcp/router-example/router.yaml

# 3. Run the walkthrough.
bash docs/ai/mcp/router-example/demo.sh
```

## Claim-based rules

`router.yaml` also includes a `tenant`-based isolation rule and an
`agent_role`-based workflow. These only take effect when the
[MCP auth filter](../mcp.md#mcp-auth-filter-dgpfilterhttpauthmcp-configuration)
is in the chain and populates JWT claims — the router consumes already-validated
claims and never re-validates tokens. The unconditional `block-privileged` rule
needs no claims, so the demo works standalone.

Dynamic Nacos updates for this PR support tool catalog updates only. Router-only
dynamic updates are rejected and should be applied by rebuilding/reloading the
filter configuration.
