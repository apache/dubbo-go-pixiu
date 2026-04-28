# MCP Transport Notes

Current source anchors:

- `pkg/filter/mcp/mcpserver/filter.go`
- `pkg/filter/mcp/mcpserver/transport/`
- `pkg/common/constant/http.go`

Pixiu's MCP filter supports Streamable HTTP behavior:

- `POST /mcp` carries JSON-RPC requests.
- `GET /mcp` with `Accept: text/event-stream` establishes an SSE stream.
- `Mcp-Session-Id` is returned/used for session continuity.
- `Mcp-Protocol-Version` may be supplied by clients; server currently
  responds with `2025-06-18`.

## Initialize

```sh
curl -i -X POST http://localhost:8888/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"curl","version":"1.0.0"}}}'
```

Expect JSON response and `Mcp-Session-Id`.

## Initialized Notification

```sh
curl -i -X POST http://localhost:8888/mcp \
  -H 'Content-Type: application/json' \
  -H 'Mcp-Session-Id: <session-id>' \
  -d '{"jsonrpc":"2.0","method":"notifications/initialized"}'
```

Expect `202 Accepted`.

## List Tools

```sh
curl -s -X POST http://localhost:8888/mcp \
  -H 'Content-Type: application/json' \
  -H 'Mcp-Session-Id: <session-id>' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
```

## SSE Stream

```sh
curl -N http://localhost:8888/mcp \
  -H 'Accept: text/event-stream' \
  -H 'Mcp-Session-Id: <session-id>'
```

GET without `Accept: text/event-stream` should return 406.

## stdio MCP Servers

Pixiu does not spawn stdio MCP processes. To integrate a stdio server:

1. Run an external bridge that exposes the stdio server over HTTP/SSE.
2. Configure Pixiu MCP tools to call that bridge as a backend HTTP
   service, or expose the bridge directly behind `httpproxy`.
3. Document bridge lifecycle and health checks outside Pixiu.
