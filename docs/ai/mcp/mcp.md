## MCP (Meta Channel Protocol) Gateway Configuration

English | [中文](./mcp_CN.md)

This document explains how to configure the MCP (Meta Channel Protocol) filters within your gateway, enabling you to securely expose backend HTTP APIs as callable "tools" for AI Agents.

### Introduction

The Meta Channel Protocol (MCP) serves as an intelligent bridge between AI Agents and your existing backend services. It dynamically translates a simple, unified protocol into standard HTTP requests, allowing agents to interact with your APIs as if they were native functions or tools. This approach simplifies agent development and provides a centralized point for security, control, and observability.

There are two primary filters for setting up an MCP endpoint:

1.  **`dgp.filter.mcp.mcpserver`**: The core MCP server filter that defines the server's identity and exposes backend APIs as tools.
2.  **`dgp.filter.http.auth.mcp`**: An optional but recommended security filter that protects the MCP endpoint using OAuth 2.0 and JWT-based authorization.

---

### MCP Server Filter (`dgp.filter.mcp.mcpserver`) Configuration

This filter is the heart of the MCP gateway. It's responsible for defining what tools are available and how they map to your backend HTTP services.

#### Server Information (`server_info`)

The `server_info` block provides metadata about your MCP server, which can be useful for discovery and diagnostics.

```yaml
server_info:
  name: "MCP OAuth Sample Server"
  version: "1.0.0"
  description: "MCP Server protected by OAuth for tools demonstration"
  instructions: "Use read/write tokens to interact with the mock server API via MCP"
```

-   **`name` (`string`)**: The display name of the MCP server.
-   **`version` (`string`)**: The version of the server.
-   **`description` (`string`)**: A brief description of the server's purpose.
-   **`instructions` (`string`)**: Instructions for clients or developers on how to interact with the exposed tools.

#### Tools Configuration (`tools`)

The `tools` section is an array where each item defines a single tool that the MCP gateway exposes. Each tool corresponds to a specific backend HTTP API endpoint.

```yaml
tools:
  - name: "get_user"
    description: "Get user information by ID with optional profile details"
    cluster: "mock-server"
    # ... request and args configuration ...
```

-   **`name` (`string`)**: The unique name of the tool. This is the identifier AI Agents will use to call it.
-   **`description` (`string`)**: A clear, concise description of what the tool does. This is crucial for LLMs to understand the tool's capabilities.
-   **`cluster` (`string`)**: The name of the upstream cluster that will handle requests for this tool.

##### Request Definition (`request`)

The `request` object specifies the details of the HTTP request that the gateway will make to the upstream cluster when the tool is called.

```yaml
request:
  method: "GET"
  path: "/api/users/{id}"
  timeout: "10s"
  headers:
    Content-Type: "application/json"
```

-   **`method` (`string`)**: The HTTP method (e.g., `GET`, `POST`, `PUT`, `DELETE`).
-   **`path` (`string`)**: The request path for the upstream service. You can use placeholders like `{arg_name}` for path parameters.
-   **`timeout` (`string`)**: The timeout for the upstream request (e.g., `5s`, `100ms`).
-   **`headers` (`object`)**: A key-value map of static HTTP headers to include in the upstream request.

##### Argument Definition (`args`)

The `args` array defines the parameters that a tool accepts. This schema allows the gateway to validate incoming arguments and correctly place them into the upstream HTTP request.

```yaml
args:
  - name: "id"
    type: "integer"
    in: "path"
    description: "User ID to retrieve"
    required: true
```

Each argument object contains the following fields:

| Field         | Type      | Description                                                                                                                              |
|---------------|-----------|------------------------------------------------------------------------------------------------------------------------------------------|
| `name`        | `string`  | The name of the argument.                                                                                                                |
| `type`        | `string`  | The data type of the argument (`string`, `integer`, `number`, `boolean`).                                                                |
| `in`          | `string`  | Specifies where the argument should be placed in the HTTP request: `path`, `query`, or `body`.                                             |
| `description` | `string`  | A detailed description of the argument's purpose, which helps LLMs use it correctly.                                                     |
| `required`    | `boolean` | Whether the argument is mandatory. Defaults to `false`.                                                                                  |
| `default`     | `any`     | A default value to use if the argument is not provided.                                                                                  |
| `enum`        | `array`   | An array of allowed values for the argument, providing a form of validation.                                                             |

---

### MCP Auth Filter (`dgp.filter.http.auth.mcp`) Configuration

This filter adds a layer of security to your MCP endpoint, ensuring that only authenticated and authorized clients can invoke the tools. It validates JWTs provided by clients against a configured identity provider.

#### Resource Metadata (`resource_metadata`)

This section defines the resource being protected and points to the authorization server(s) that can grant access to it.

```yaml
resource_metadata:
  path: "/.well-known/oauth-protected-resource/mcp"
  resource: "http://localhost:8888/mcp"
  authorization_servers:
    - "http://localhost:9000"
```

-   **`path` (`string`)**: The path where the resource metadata is exposed, following standards for discoverability.
-   **`resource` (`string`)**: The identifier for the resource being protected (typically the URL of the MCP endpoint itself).
-   **`authorization_servers` (`array` of `string`)**: A list of trusted authorization server URLs.

#### Providers (`providers`)

This section defines the trusted JWT issuers (identity providers). The gateway will use this information to validate the signature of incoming JWTs.

```yaml
providers:
  - name: "local"
    issuer: "http://localhost:9000"
    jwks: "http://localhost:9000/.well-known/jwks.json"
```

-   **`name` (`string`)**: A unique name for this provider configuration.
-   **`issuer` (`string`)**: The `iss` (issuer) claim expected in the JWT. This must match the issuer's identifier.
-   **`jwks` (`string`)**: The URL of the JSON Web Key Set (JWKS) endpoint, where the public keys for verifying JWT signatures are published.

#### Rules (`rules`)

Rules connect the authentication policy to specific upstream clusters.

```yaml
rules:
  - cluster: "mcp-protected"
```

-   **`cluster` (`string`)**: The name of the cluster to which this authentication and authorization policy applies.

---

### Complete Configuration Example

This example demonstrates a complete setup for an MCP gateway. It includes an MCP server with several tools and is protected by the MCP auth filter.

```yaml
# Full pixiu_mcp_auth_test.yaml example
static_resources:
  listeners:
    - name: "net/http"
      protocol_type: "HTTP"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8888
      filter_chains:
        filters:
          - name: "dgp.filter.httpconnectionmanager"
            config:
              route_config:
                routes:
                  # All traffic is routed to a virtual cluster.
                  - match:
                      prefix: "/"
                    route:
                      cluster: "mcp-protected"
              http_filters:
                # (Optional) MCP Authorization Filter to protect the endpoint
                - name: "dgp.filter.http.auth.mcp"
                  config:
                    resource_metadata:
                      path: "/.well-known/oauth-protected-resource/mcp"
                      resource: "http://localhost:8888/mcp"
                      authorization_servers:
                        - "http://localhost:9000"
                    providers:
                      - name: "local"
                        issuer: "http://localhost:9000"
                        jwks: "http://localhost:9000/.well-known/jwks.json"
                    rules:
                      - cluster: "mcp-protected"

                # Core MCP Server Filter
                - name: "dgp.filter.mcp.mcpserver"
                  config:
                    server_info:
                      name: "MCP OAuth Sample Server"
                      version: "1.0.0"
                      description: "MCP Server protected by OAuth for tools demonstration"
                      instructions: "Use appropriate tokens to interact with the mock server API via MCP"
                    
                    tools:
                      # Tool 1: Get a user by ID
                      - name: "get_user"
                        description: "Get user information by ID"
                        cluster: "mock-server"
                        request:
                          method: "GET"
                          path: "/api/users/{id}"
                          timeout: "10s"
                        args:
                          - name: "id"
                            type: "integer"
                            in: "path"
                            description: "User ID to retrieve"
                            required: true

                      # Tool 2: Create a new user
                      - name: "create_user"
                        description: "Create a new user account"
                        cluster: "mock-server"
                        request:
                          method: "POST"
                          path: "/api/users"
                          timeout: "10s"
                          headers:
                            Content-Type: "application/json"
                        args:
                          - name: "name"
                            type: "string"
                            in: "body"
                            description: "User's full name"
                            required: true
                          - name: "email"
                            type: "string"
                            in: "body"
                            description: "User's email address"
                            required: true
                
                # Standard HTTP Proxy filter for downstream requests
                - name: "dgp.filter.http.httpproxy"

  clusters:
    # Upstream backend service
    - name: "mock-server"
      type: "STATIC"
      lb_policy: "ROUND_ROBIN"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 8081

    # Virtual cluster used for routing rules
    - name: "mcp-protected"
      type: "STATIC"
      lb_policy: "ROUND_ROBIN"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 8081
```
