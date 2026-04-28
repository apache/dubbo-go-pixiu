# MCP Nacos Registry Integration

Current source anchors:

- `pkg/adapter/mcpserver/registrycenter.go`
- `pkg/adapter/mcpserver/registry/nacos/`
- `docs/ai/mcp/mcp.md`

Use `dgp.adapter.mcpserver` when MCP tool definitions are managed by
Nacos. The listener still needs `dgp.filter.mcp.mcpserver`; the adapter
updates the in-process tool registry.

## Adapter Config

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

Current adapter code only handles `protocol: nacos`; other registry
protocols are skipped.

## Dynamic Tool Effects

When Nacos config changes:

- The adapter applies `McpServerConfig` into the global MCP dynamic
  registry.
- For tools with `backend_url`, it parses host/port and registers an
  endpoint in the cluster named by `tool.cluster`.

## Static Plus Dynamic

It is valid to configure `server_info` statically and load tools
dynamically. Keep a minimal `dgp.filter.mcp.mcpserver` config in
`http_filters` even when all tools come from Nacos.

## Nacos Operational Notes

- Use `namespace` for environment isolation.
- Use `group` for service grouping.
- Ensure generated `backend_url` values are valid URLs with host and
  port; malformed URLs are logged and skipped by the adapter.
