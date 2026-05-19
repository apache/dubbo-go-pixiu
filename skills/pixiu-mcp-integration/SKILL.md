---
name: pixiu-mcp-integration
description: |
  Use when configuring dubbo-go-pixiu as an MCP gateway:
  `dgp.filter.mcp.mcpserver`, Streamable HTTP/SSE, `tools/list`,
  `tools/call`, OAuth/JWT `dgp.filter.http.auth.mcp`, Nacos
  `dgp.adapter.mcpserver`, or stdio MCP bridge guidance.
allowed-tools: [Read, Grep, Glob, Edit, Write, Bash]
metadata:
  version: "0.1.1"
  domain: ai-gateway
  scope: generate-and-validate
  triggers: ["MCP", "Model Context Protocol", "mcpserver", "tools/list", "tools/call", "Mcp-Session-Id", "dgp.filter.mcp.mcpserver", "dgp.filter.http.auth.mcp", "mcp gateway", "MCP tool"]
  pixiu_min_version: "0.6.0"
  experimental: true
  role: specialist
---

# pixiu-mcp-integration — Exposing HTTP APIs as MCP Tools

Pixiu's current MCP integration is an HTTP filter that speaks MCP
Streamable HTTP/SSE to clients and maps `tools/call` to backend HTTP
clusters. It can also expose resources, resource templates, prompts,
and optional OAuth/JWT protection.

Important boundary: current Pixiu does not directly spawn or manage
stdio MCP servers. For a stdio server, put an external bridge in front
of it, then configure Pixiu to call the bridge as an HTTP backend tool.

## When to Use

Use this skill when the user wants to:

- Expose existing backend HTTP APIs as MCP tools.
- Configure MCP endpoint metadata, tools, resources, templates, or
  prompts.
- Support Streamable HTTP / SSE MCP clients.
- Protect the MCP endpoint with OAuth 2.0 protected-resource metadata
  and JWT validation.
- Load MCP tool definitions dynamically from Nacos.
- Integrate stdio MCP servers by documenting the required bridge.

Do not use this skill for:

- LLM model proxying; use `pixiu-llm-gateway`.
- HTTP to Dubbo API mapping; use `pixiu-http-to-dubbo`.
- Writing a new filter implementation; use `pixiu-filter-author`.

## Step 0 — Verify Current Source First

Read current source before generating config:

1. `pkg/common/constant/key.go` for Kind strings.
2. `pkg/model/mcpserver.go` for `McpServerConfig`, tool args,
   resources, templates, and prompts.
3. `pkg/filter/mcp/mcpserver/plugin.go`, `filter.go`, `handlers.go`,
   and `transport/` for endpoint matching, Streamable HTTP/SSE, and
   JSON-RPC handling.
4. `pkg/filter/auth/mcp/config.go` and `filter.go` for optional auth.
5. `pkg/adapter/mcpserver/registrycenter.go` and
   `pkg/adapter/mcpserver/registry/nacos/` for dynamic tool discovery.
6. `docs/ai/mcp/mcp.md` if present.

If current source disagrees with this skill, trust the source and note
the difference.

## Step 1 — Gather the Nine Things

Do not generate final config until these are explicit:

1. MCP endpoint path, usually `/mcp`.
2. Server name, version, description, and instructions.
3. Static tools or Nacos-managed tools.
4. For every static tool: name, description, upstream cluster,
   request method/path/timeout, and argument schema.
5. Backend clusters for each tool.
6. Whether resources, resource templates, or prompts should be exposed.
7. Transport expectation: JSON POST only, SSE GET stream, or both.
8. Auth requirement: none, or `dgp.filter.http.auth.mcp` with issuer,
   JWKS URL, resource URI, and protected cluster.
9. If the user says stdio MCP: which external bridge will expose it as
   HTTP/SSE. Pixiu itself does not run the stdio process.

## Step 2 — Shape the Static MCP Server Config

Before generating a static tool config, read `pkg/model/mcpserver.go`
and the current MCP filter source. Keep the generated config compact but
include these essentials:

- `endpoint` must match the client URL path. Current filter checks
  `ctx.Request.URL.Path == cfg.Endpoint`.
- The `mcp-backend` route cluster is a routing anchor for the MCP
  endpoint. Each tool has its own `cluster` for the actual backend call.
- Declare every static tool cluster under `static_resources.clusters[]`.
- Put `dgp.filter.mcp.mcpserver` before `dgp.filter.http.httpproxy`.
- `tools/list`, `resources/list`, `prompts/list`, `initialize`,
  `ping`, and `notifications/initialized` are terminal methods handled
  by the MCP filter.
- `tools/call` continues through the chain to an HTTP proxy, then
  Encode converts the backend response into MCP tool-call output.
- When explaining ordering, spell it out as
  `dgp.filter.http.auth.mcp` before `dgp.filter.mcp.mcpserver`, and
  `dgp.filter.mcp.mcpserver` before `dgp.filter.http.httpproxy`.
  Terminal methods (`initialize`, `tools/list`, `resources/list`,
  `prompts/list`, `ping`) stop at the MCP filter; `tools/call`
  continues to `httpproxy`.
- Current source declares `request.headers`, but `buildBackendRequest`
  does not apply them yet. Do not rely on static tool headers for
  runtime behavior; request bodies still get `Content-Type:
  application/json` automatically.

## Step 3 — Add OAuth/JWT Protection When Required

Place `dgp.filter.http.auth.mcp` before `dgp.filter.mcp.mcpserver`:

```yaml
http_filters:
  - name: dgp.filter.http.auth.mcp
    config:
      resource_metadata:
        path: "/.well-known/oauth-protected-resource/mcp"
        resource: "https://mcp.example.com/mcp"
        authorization_servers:
          - "https://auth.example.com"
      providers:
        - name: "main"
          issuer: "https://auth.example.com"
          jwks: "https://auth.example.com/.well-known/jwks.json"
          audience: "https://mcp.example.com/mcp"
      rules:
        - cluster: "mcp-backend"
  - name: dgp.filter.mcp.mcpserver
    config: ...
```

Auth rules match the route entry's cluster. If `rules[].cluster` does
not match the MCP route cluster, requests will not be protected.
On successful validation the current auth filter removes the
`Authorization` header before forwarding. Do not promise that the
caller token reaches the tool backend; if the backend needs credentials,
design an explicit downstream auth strategy instead of relying on the
validated bearer token being forwarded.
Safe alternatives are backend-side service auth, an external bridge or
proxy that injects credentials, mTLS or network policy, or a source
change that explicitly applies outbound headers.

Before generating auth config, read the current MCP auth filter source.

## Step 4 — Use Nacos Only for Dynamic Tool Definitions

Static tools are easier. Use `dgp.adapter.mcpserver` only when a Nacos
MCP registry supplies tool definitions.

```yaml
adapters:
  - id: "mcp-nacos-adapter"
    name: dgp.adapter.mcpserver
    config:
      registries:
        nacos:
          protocol: "nacos"
          address: "127.0.0.1:8848"
          timeout: "5s"
          username: "nacos"
          password: "nacos"
          namespace: ""
          group: "DEFAULT_GROUP"
```

The listener still needs `dgp.filter.mcp.mcpserver`; the adapter updates
the in-process tool registry and may register endpoints for tools with
`backend_url`.

Before generating Nacos instructions, read the current MCP registry
adapter source.

## Step 5 — Validate and Smoke Test

Before booting Pixiu, inspect `conf.yaml` directly. Check yaml syntax,
filter order, MCP server config shape, route targets, static tool
definitions, auth settings, and Nacos adapter basics.

Smoke-test Streamable HTTP:

```sh
curl -i -X POST http://localhost:8888/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"curl","version":"1.0.0"}}}'
```

For SSE, the client should call `GET /mcp` with
`Accept: text/event-stream` and then reuse the returned
`Mcp-Session-Id` on POST requests.

## Cross-Cutting Rules

### Always

- Use exact current Kinds: `dgp.filter.mcp.mcpserver`,
  `dgp.filter.http.auth.mcp`, and `dgp.adapter.mcpserver`.
- Put auth before MCP server when auth is enabled.
- Keep `endpoint` and route prefix/path aligned.
- Ensure every tool `cluster` exists in `static_resources.clusters[]`
  unless it is dynamically supplied by the adapter.
- Use arg `in` values `path`, `query`, or `body`.
- Use arg `type` values `string`, `integer`, `number`, or `boolean`.

### Never

- Claim Pixiu can directly run a stdio MCP server. It needs an external
  HTTP/SSE bridge for stdio-based servers.
- Put `tools` at top level; they belong under the MCP filter config.
- Put `dgp.filter.mcp.mcpserver` after `httpproxy`; terminal MCP
  methods must be handled before proxying.
- Forget `dgp.filter.http.httpproxy` when `tools/call` needs to reach a
  backend HTTP service.
- Configure auth `rules[].cluster` with the tool backend cluster when
  the MCP route itself uses a different route cluster.
- Rely on `request.headers` for static tool backend credentials in the
  current source. The model field exists, but `buildBackendRequest` does
  not apply it yet; use an external bridge/proxy, backend-side auth, or
  a source change before promising static header injection.

## Source Files To Read

- `pkg/model/mcpserver.go`
- `pkg/filter/mcp/mcpserver/`
- `pkg/filter/auth/mcp/`
- `pkg/adapter/mcpserver/`
- `docs/ai/mcp/` if present.
