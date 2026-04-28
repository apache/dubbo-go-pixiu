# MCP Server Filter Configuration

Current source anchors:

- `pkg/model/mcpserver.go`
- `pkg/filter/mcp/mcpserver/plugin.go`
- `pkg/filter/mcp/mcpserver/filter.go`
- `pkg/filter/mcp/mcpserver/handlers.go`

## Kind and Endpoint

Use current Kind:

```yaml
- name: dgp.filter.mcp.mcpserver
  config:
    endpoint: "/mcp"
```

Current filter matches MCP requests by exact path:
`ctx.Request.URL.Path == cfg.Endpoint`.

## Server Info

```yaml
server_info:
  name: "Pixiu MCP Gateway"
  version: "1.0.0"
  description: "Expose backend APIs as MCP tools"
  instructions: "Use tools/list then tools/call."
```

## Tools

Each tool maps an MCP `tools/call` to a backend HTTP request.

```yaml
tools:
  - name: "get_user"
    description: "Get user information by ID"
    cluster: "user-service"
    request:
      method: "GET"
      path: "/api/users/{id}"
      timeout: "10s"
    args:
      - name: "id"
        type: "integer"
        in: "path"
        description: "User ID"
        required: true
```

Tool fields:

| Field | Meaning |
|---|---|
| `name` | MCP tool name used by clients |
| `description` | LLM-readable tool behavior |
| `cluster` | Backend HTTP cluster to call |
| `backend_url` | Optional dynamic-registry URL used by adapter flows |
| `request.method` | `GET`, `POST`, `PUT`, `DELETE`, etc. |
| `request.path` | Backend path, may include `{name}` placeholders |
| `request.timeout` | Upstream timeout string |
| `request.headers` | Present in the model, but current `buildBackendRequest` does not apply it at runtime |

Current source note: `RequestConfig.Headers` exists in the model, but the
MCP tool-call path does not copy those headers into the backend request
yet. For POST/PUT tools with body args, Pixiu sets
`Content-Type: application/json` automatically.

Arg fields:

| Field | Values |
|---|---|
| `type` | `string`, `integer`, `number`, `boolean` |
| `in` | `path`, `query`, `body` |

Path placeholders such as `/api/users/{id}` should have a matching arg
with `name: id` and `in: path`.

## Resources

```yaml
resources:
  - name: "service-doc"
    uri: "docs://service"
    description: "Backend service documentation"
    mime_type: "text/markdown"
    source:
      type: "inline"
      content: "# Service Docs"
```

Supported source shapes in the model are `file`, `url`, `inline`, and
`template`.

## Resource Templates

```yaml
resource_templates:
  - name: "user-doc"
    uri_template: "docs://users/{id}"
    title: "User docs"
    description: "Documentation for a user"
    mime_type: "text/plain"
    parameters:
      - name: "id"
        type: "string"
        required: true
```

## Prompts

```yaml
prompts:
  - name: "summarize_user"
    title: "Summarize user"
    description: "Summarize user data"
    arguments:
      - name: "user_id"
        required: true
    messages:
      - role: "user"
        content: "Summarize user {{user_id}}"
```

Prompt message roles are `user` or `assistant`.
