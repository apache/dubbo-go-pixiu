---
name: pixiu-mcp-gateway
description: Creates and validates dubbo-go-pixiu MCP gateway conf.yaml. Use for MCP server filters, tools/list, tools/call, Mcp-Session-Id, MCP auth, Nacos MCP registry, or exposing backend HTTP APIs as MCP tools. Do not use for LLM model proxying.
---

# pixiu-mcp-gateway

## Purpose
Generate a complete or embeddable MCP gateway `conf.yaml` that exposes backend HTTP APIs as MCP tools.

## When to use
- Use when:
  - Exposing backend HTTP APIs as MCP tools through the Pixiu MCP gateway, in static or Nacos dynamic mode, or configuring resources/templates/prompts, OAuth/JWT protection, or Streamable HTTP/SSE.

- Do not use for:
  - LLM model proxying or implementing a new MCP-related filter.

## Inputs
- Choose tool source (required):
  - Required:
    - `tool_mode`: `static` or `registry`

- `listener` (required):
  - Defaults:
    - `address.socket_address.address: 0.0.0.0`
    - `address.socket_address.port: 8888`

- `route_config` (required):
  - Defaults:
    - `routes[].match.prefix: /mcp`
    - `routes[].route.cluster: mcp-route` (route cluster name)

- `dgp.filter.mcp.mcpserver` (required):
  - Defaults:
    - `config.endpoint: /mcp`
    - `config.server_info.name: Pixiu MCP Server`
    - `config.server_info.version: 1.0.0`
    - `config.server_info.description: MCP Server powered by Apache Dubbo-go-pixiu`
    - `config.server_info.instructions: Use the provided tools to interact with backend services.`

- Static tools (required when `tool_mode: static`):
  - Required:
    - `tools[].name`
    - `tools[].request.path`
    - Endpoint under `static_resources.clusters[]` for each backend cluster referenced by a tool:
      - `clusters[].endpoints[].socket_address.address`
      - `clusters[].endpoints[].socket_address.port`
  - Optional:
    - `tools[].description`
    - `tools[].args[]`
  - Defaults (`tools[].cluster` and `clusters[].name` reference the same name; changing one requires changing the other):
    - `tools[].cluster: mcp-backend`
    - `clusters[].name: mcp-backend`
    - `tools[].request.method: GET`
    - `tools[].request.timeout: 30s`
    - `tools[].args[].type: string`

- Registry tools (required when `tool_mode: registry`):
  - Required:
    - `registries.nacos.address`
  - Optional:
    - `registries.nacos.group`
    - `registries.nacos.namespace`
    - `registries.nacos.username`
    - `registries.nacos.password`
  - Defaults:
    - `adapters[].id: mcp-nacos`
    - `adapters[].name: dgp.adapter.mcpserver`
    - `registries.nacos.protocol: nacos`
    - `registries.nacos.timeout: 5s`
    - `registries.nacos.group: DEFAULT_GROUP`

- `resources[]` (optional):
  - Required:
    - `resources[].name`
    - `resources[].uri`
    - `resources[].source.type`
    - `resources[].source.content`
  - Optional:
    - `resources[].description`
    - `resources[].mime_type`

- `resource_templates[]` (optional):
  - Required:
    - `resource_templates[].name`
    - `resource_templates[].uri_template`
  - Optional:
    - `resource_templates[].title`
    - `resource_templates[].description`
    - `resource_templates[].mime_type`
    - `resource_templates[].parameters[]`
    - `resource_templates[].annotations`

- `prompts[]` (optional):
  - Required:
    - `prompts[].name`
    - `prompts[].messages[].role`
    - `prompts[].messages[].content`
  - Optional:
    - `prompts[].title`
    - `prompts[].description`
    - `prompts[].arguments[]`

- `dgp.filter.http.auth.mcp` (optional):
  - Required:
    - `resource_metadata.resource`
    - `resource_metadata.authorization_servers[]`
    - `providers[].name`
    - `providers[].issuer`
    - `providers[].jwks`
  - Defaults:
    - `resource_metadata.path: /.well-known/oauth-protected-resource`
    - `providers[].audience`: defaults to `resource_metadata.resource`
    - `rules[].cluster`: defaults to `route_config.routes[].route.cluster`

## Workflow
1. When field shape or behavior is uncertain, read current source before generating YAML:
   - `pkg/common/constant/key.go`.
   - `pkg/model/mcpserver.go`.
   - `pkg/filter/mcp/mcpserver/plugin.go`, `filter.go`, `handlers.go`, and `transport/`.
   - Read `pkg/filter/auth/mcp/config.go` and `filter.go` when auth is required.
   - Read `pkg/adapter/mcpserver/registrycenter.go` and `pkg/adapter/mcpserver/registry/nacos/` when Nacos is required.
2. Guide and read user input:
   1. Read existing config and already provided user information first; do not ask again for information already present or stated.
   2. Ask one `Inputs` group at a time. Each time, output only that group's required fields, with a short explanation after each field.
   3. If the current group's required fields are incomplete, ask only for the missing fields and do not move to the next group.
   4. After the current group's required fields are complete, ask whether to fill that group's optional fields; if yes, list those optional fields with short explanations.
   5. After optional fields are skipped or completed, apply that group's defaults and output default-value information; defaults must not override existing config or user input.
3. Choose static or registry path by `tool_mode`: static writes tools into the MCP filter; registry uses the MCP server adapter.
4. Generate listener, route, and filters: route the MCP endpoint to the route cluster, and put the auth filter before the MCP server when needed.
5. Generate tool backend: write static tool backend clusters under `static_resources.clusters[]`; map args with `path`, `query`, or `body`.

## Output format
- Show the relevant YAML fragments.

## Validation
- Verify `dgp.filter.http.auth.mcp` `rules[].cluster` matches the MCP route cluster, such as `mcp-route`, not a tool backend cluster.
- Verify MCP filter-chain order: `dgp.filter.http.auth.mcp` -> `dgp.filter.mcp.mcpserver` -> `dgp.filter.http.httpproxy`; ignore missing filters within this order chain.
- Verify every static tool `cluster` has a same-name cluster declaration under `static_resources.clusters[]`.
- Verify `dgp.filter.mcp.mcpserver.config.endpoint` exactly matches the path in the client call URL.

## Examples
Static mode (`conf.yaml`):

```yaml
static_resources:
  listeners:
    - name: net/http
      protocol_type: HTTP
      address:
        socket_address:
          address: 0.0.0.0
          port: 8888
      filter_chains:
        filters:
          - name: dgp.filter.httpconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: /mcp
                    route:
                      cluster: mcp-route
              http_filters:
                - name: dgp.filter.http.auth.mcp
                  config:
                    resource_metadata:
                      path: /.well-known/oauth-protected-resource/mcp
                      resource: http://localhost:8888/mcp
                      authorization_servers:
                        - http://localhost:9000
                    providers:
                      - name: local
                        issuer: http://localhost:9000
                        jwks: http://localhost:9000/.well-known/jwks.json
                        audience: http://localhost:8888/mcp
                    rules:
                      - cluster: mcp-route
                - name: dgp.filter.mcp.mcpserver
                  config:
                    endpoint: /mcp
                    server_info:
                      name: demo-mcp
                      version: "1.0.0"
                      description: Demo MCP server for user APIs
                      instructions: Use the exposed tools to query and create users.
                    tools:
                      - name: search_users
                        description: Search users by status
                        cluster: users
                        request:
                          method: GET
                          path: /users
                          timeout: 5s
                        args:
                          - name: status
                            type: string
                            in: query
                            description: User status filter
                            required: false
                            default: active
                            enum:
                              - active
                              - disabled
                    resources:
                      - name: service-status
                        uri: pixiu://service/status
                        description: Current service status
                        mime_type: application/json
                        source:
                          type: inline
                          content: '{"status":"ok"}'
                    resource_templates:
                      - name: user-profile
                        uri_template: pixiu://users/{id}
                        title: User profile
                        description: User profile by ID
                        mime_type: application/json
                        parameters:
                          - name: id
                            type: string
                            description: User ID
                            required: true
                        annotations:
                          audience:
                            - assistant
                          priority: 0.7
                    prompts:
                      - name: summarize_user
                        title: Summarize user
                        description: Summarize a user profile
                        arguments:
                          - name: id
                            description: User ID
                            required: true
                        messages:
                          - role: user
                            content: "Summarize user {{id}}."
                - name: dgp.filter.http.httpproxy
  clusters:
    - name: mcp-route
      endpoints:
        - socket_address:
            address: 127.0.0.1
            port: 8080
    - name: users
      endpoints:
        - socket_address:
            address: 127.0.0.1
            port: 8080
```

Registry mode example (`conf.yaml`):

```yaml
static_resources:
  listeners:
    - name: net/http
      protocol_type: HTTP
      address:
        socket_address:
          address: 0.0.0.0
          port: 8888
      filter_chains:
        filters:
          - name: dgp.filter.httpconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: /mcp
                    route:
                      cluster: mcp-route
              http_filters:
                - name: dgp.filter.mcp.mcpserver
                  config:
                    endpoint: /mcp
                    server_info:
                      name: demo-mcp
                      version: "1.0.0"
                      description: Demo MCP server with Nacos-managed tools
                      instructions: Tool definitions are loaded from Nacos.
                - name: dgp.filter.http.httpproxy
  clusters:
    - name: mcp-route
      endpoints:
        - socket_address:
            address: 127.0.0.1
            port: 8080
  adapters:
    - id: mcp-nacos
      name: dgp.adapter.mcpserver
      config:
        registries:
          nacos:
            protocol: nacos
            address: "127.0.0.1:8848"
            timeout: "5s"
            group: DEFAULT_GROUP
```
