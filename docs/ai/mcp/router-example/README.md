# MCP Intelligent Tool Router — Runnable Example

This directory demonstrates the MCP intelligent tool router added for
[issue #937](https://github.com/apache/dubbo-go-pixiu/issues/937).

## Files

- `router.yaml` — a complete Pixiu gateway config with the router enabled:
  4 tools, 1 always-on policy rule (deny `internal`/`admin` tags), 2 workflow
  bundles, and the admin debug endpoint turned on.
- `demo.sh` — a curl walkthrough: initialize → tools/list (trimmed) →
  tools/call (denied vs allowed) → inspect the session plan.

## What it shows

1. **`tools/list` trimming** — `internal_dump` (tagged `internal`/`admin`) is
   filtered out of the listing by the always-on `block-privileged` policy rule.
2. **`tools/call` enforcement** — calling `internal_dump` is denied because it
   is not in the session's plan, even though the client knows its name.
3. **Allowed call** — `search_kb` is in the plan and is forwarded to the backend.
4. **Admin inspection** — from loopback, `GET /__mcp/router/plan/{session_id}`
   returns the session's plan (selected tools, decision traces, version, mode),
   available because `audit.payload_logging: true`.

The example uses `fallback: bundle_default` so an empty selection falls back to
the `safe-minimal` bundle. Treat that as a discovery safety net, not an
authorization boundary; use `fallback: fail_closed` when a denied subject should
see no tools.

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
claims and never re-validates tokens. The always-on `block-privileged` rule
needs no claims, so the demo works standalone.
