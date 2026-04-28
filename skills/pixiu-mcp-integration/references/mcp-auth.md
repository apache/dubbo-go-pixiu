# MCP Auth Filter

Current source anchors:

- `pkg/filter/auth/mcp/config.go`
- `pkg/filter/auth/mcp/filter.go`
- `pkg/filter/auth/mcp/internal/validator/`

Use `dgp.filter.http.auth.mcp` when the MCP endpoint must be protected
by OAuth 2.0 protected-resource metadata and JWT validation.

## Filter Order

```yaml
http_filters:
  - name: dgp.filter.http.auth.mcp
    config: ...
  - name: dgp.filter.mcp.mcpserver
    config: ...
  - name: dgp.filter.http.httpproxy
```

Auth must run before the MCP server filter.

## Config

```yaml
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
```

Required:

- `resource_metadata.resource`
- `resource_metadata.authorization_servers`
- at least one provider
- each provider `name`, `issuer`, and `jwks`

If provider `audience` is omitted, current code defaults it to
`resource_metadata.resource`.

## Rules Match Route Cluster

Rules are matched against the route entry's cluster, not the tool's
backend cluster. If route config sends `/mcp` to `mcp-backend`, then
auth should use:

```yaml
rules:
  - cluster: "mcp-backend"
```

## Runtime Behavior

- The metadata path returns the OAuth protected-resource metadata JSON.
- Missing/invalid bearer tokens return `401` with `WWW-Authenticate`.
- On success, the Authorization header is removed before forwarding to
  downstream services to avoid token leakage.
