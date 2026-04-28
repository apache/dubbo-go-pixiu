# MCP Gateway Troubleshooting

## 404 or route not matched

- Check the route prefix/path covers the MCP endpoint.
- Current MCP filter only handles requests where URL path exactly equals
  `config.endpoint`, usually `/mcp`.

## `tools/list` is empty

- Static `tools` may be missing under the MCP filter config.
- If using Nacos, check the adapter is configured and logs show dynamic
  config applied.
- Confirm `dgp.filter.mcp.mcpserver` appears before proxy filters.

## `tools/call` does not reach backend

- Every tool `cluster` must exist in `static_resources.clusters[]` or be
  dynamically supplied by the adapter.
- `dgp.filter.http.httpproxy` must appear after the MCP server filter.
- Path placeholders like `{id}` should have matching `args` entries with
  `in: path`.

## SSE GET returns 406

- Client must send `Accept: text/event-stream`.
- Use `Mcp-Session-Id` from initialize if resuming a session.

## Auth does not protect the endpoint

- `dgp.filter.http.auth.mcp` must be before `dgp.filter.mcp.mcpserver`.
- `rules[].cluster` must match the route cluster for the MCP endpoint.
- `resource_metadata.resource` and provider `audience` must match the
  token audience expected by the issuer.

## stdio MCP server does not work

Pixiu does not launch stdio MCP servers. Run a bridge outside Pixiu and
expose it over HTTP/SSE, then configure Pixiu routes/tools to call that
bridge.
