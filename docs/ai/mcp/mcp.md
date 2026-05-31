## MCP (Model Context Protocol) Gateway Configuration

English | [中文](./mcp_CN.md)

This document explains how to configure the MCP (Model Context Protocol) filters within your gateway, enabling you to securely expose backend HTTP APIs as callable "tools" for AI Agents.

### Introduction

The Model Context Protocol (MCP) serves as an intelligent bridge between AI Agents and your existing backend services. It dynamically translates a simple, unified protocol into standard HTTP requests, allowing agents to interact with your APIs as if they were native functions or tools. This approach simplifies agent development and provides a centralized point for security, control, and observability.

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

### Intelligent Tool Routing (`router`) Configuration

When an MCP server exposes a large catalog of tools, sending all of them to an LLM on every `tools/list` causes tool overload, context bloat, and an uncontrolled exposure surface. The optional `router` block adds a governance layer that trims `tools/list` per session and enforces the trimmed set at `tools/call`.

**The router is disabled by default.** Omitting the `router` block — or setting `router.enabled: false` — preserves the exact pre-router behavior: every tool is listed and every call is forwarded.

Two principles shape the design:

1. **Discovery vs. execution separation** — tools are trimmed at `tools/list` *and* re-validated at `tools/call`, so a tool name learned out-of-band still cannot be invoked unless it is part of the session's plan.
2. **Deterministic before semantic** — policy rules and workflow bundles decide authorization. No black-box model gates access.

The selection pipeline runs in a fixed order; each stage is independently toggleable:

```
policy  ->  workflow  ->  progressive
```

```yaml
- name: "dgp.filter.mcp.mcpserver"
  config:
    server_info: { name: "Pixiu MCP Server", version: "1.0.0" }
    endpoint: "/mcp"
    router:
      enabled: true
      fallback: "bundle_default"      # bundle_default (default) | fail_closed
      default_bundle: "safe-minimal"  # bundle used when a selection is empty
      enforce_on_call: true           # default true: reject calls outside the plan
      stages:
        policy: true                  # default true
        workflow: true                # default true
        progressive: false            # default false
      policy:
        rules:
          - name: "tenant-isolation"
            when: { claim: "tenant", equals: "acme" }
            allow_tags: ["acme", "shared"]
            deny_tags: ["internal", "admin"]
          - name: "low-risk-for-anonymous"
            when: { missing_claim: "sub" }
            max_risk: "low"           # low | medium | high
      workflows:
        - name: "support-agent"       # workflow names must be non-empty and unique
          tools: ["search_kb", "create_ticket", "get_user"]
          when: { claim: "agent_role", equals: "support" }
        - name: "safe-minimal"        # name-addressable bundle (no `when`)
          tools: ["ping", "health_check"]
      progressive:
        initial_bundle: "safe-minimal"
        expand_after_calls: 1
      audit:
        sample_rate: 0.0              # 0 disables decision logs; (0,1] samples
        payload_logging: false        # opt-in; also enables the admin endpoint
```

#### Tool Metadata (`tools[].meta`)

Routing stages act on optional per-tool metadata. Tools without `meta` are treated as untagged, low-risk, and visible — so existing configs keep working unchanged. Unknown `meta.risk` or `policy.max_risk` values are configuration errors and fail fast on startup or dynamic update.

```yaml
tools:
  - name: "get_user"
    description: "Get user information by ID"
    cluster: "user-svc"
    request: { method: "GET", path: "/api/users/{id}" }
    meta:
      tags: ["user", "read"]
      capabilities: ["user.read"]
      risk: "low"                     # low | medium | high (default low)
      discovery_visibility: true      # false => hidden from tools/list, still authorized/callable when selected
```

#### Pipeline Stages

| Stage | What it does |
|-------|--------------|
| **policy** | Hard filter on JWT claims. A rule applies when its `when` clause matches; matching rules enforce `allow_tags` (keep only intersecting tags), `deny_tags` (drop any match), and `max_risk` (drop tools above the risk ceiling). |
| **workflow** | The first workflow whose `when` clause matches keeps only that bundle's tools. Bundles without a `when` clause are name-addressable only (used by fallback / progressive). |
| **progressive** | A fresh session sees only `initial_bundle`; after `expand_after_calls` successful completed tool calls, the full filtered set is revealed. Authorization checks and backend failures do not advance the counter. When this stage is enabled, `progressive.initial_bundle` is required and must reference a defined workflow bundle. |

`when` clauses (used by both policy rules and workflows) support: `claim` + `equals`, `claim` + `in: [...]`, `claim` + `regex`, `missing_claim`, or `claim` alone (presence check). The `sub` and `tenant` claims are promoted from the validated JWT for convenient matching. Claims come from the [MCP Auth Filter](#mcp-auth-filter-dgpfilterhttpauthmcp-configuration); without that filter in the chain, claim-based rules simply do not match. The router consumes already-validated claims and never re-validates tokens. Workflow names must be non-empty and unique because fallback and progressive disclosure address bundles by name.

#### Fallback

If the pipeline produces an empty selection (e.g. overly strict rules), the `fallback` strategy decides the outcome:

- **`bundle_default`** (default): expose the `default_bundle` workflow (intersected with the policy-allowed live catalog). `default_bundle` must be non-empty, the workflow list must exist, and the named workflow must contain at least one tool name, or the gateway fails to start.
- **`fail_closed`**: return no tools. There is intentionally no `fail_open` — a governance-layer failure must never *widen* the exposure surface.

> **Note — fallback does not distinguish "denied" from "over-filtered".** The pipeline treats *any* empty result as the fallback trigger, whether it came from overly strict rules or from a rule that *intentionally* denies every tool for a subject. With `fallback: bundle_default`, a subject you meant to fully lock out will therefore still see the `default_bundle` (intersected with candidates). If you want "deny means an empty tool set," use `fallback: fail_closed`; the `default_bundle` is only ever a safety net, not an authorization boundary.

#### Enforcement at `tools/call`

With `enforce_on_call: true` (default), a `tools/call` for a tool not in the session's plan is rejected with a tool-call error. The router re-validates the plan against the current claims, router config version, and live tool catalog before authorizing the call, so stale plans are recomputed instead of treated as long-lived credentials. A client that skips `tools/list` entirely has no plan and is therefore denied (reason `no_session_plan`). Set `enforce_on_call: false` to allow calls without plan enforcement.

Tools with `meta.discovery_visibility: false` are omitted from the plan's `visible_tool_names` / `tools/list` view but remain in the authorized `tool_names` set when selected by policy, workflow, or progressive stages. This supports hidden-but-callable tools for clients that already know the tool name while keeping discovery quieter.

#### Observability

When the router is enabled it publishes Prometheus metrics under the `pixiu_mcp_tool_router_*` namespace:

| Metric | Type | Labels |
|--------|------|--------|
| `select_total` | counter | `result` (ok/fallback/cached), `mode` |
| `selection_latency_ms` | histogram | `stage` |
| `candidates_count` / `selected_count` | histogram | — |
| `fallback_total` | counter | `reason` |
| `call_denied_total` | counter | `reason` (no_session_plan / not_in_plan / stale_plan_recompute_failed) |
| `plans_active` | gauge | — |

Decision logs are off by default. Set `audit.sample_rate` to a value in `(0,1]` to emit structured, PII-safe records (`event: mcp_router_decision`) carrying counts, mode, per-stage drop tallies, and metadata version. Bounded denied-tool samples are included only when `audit.payload_logging: true`.

#### Admin Debug Endpoint

When `audit.payload_logging: true`, the gateway serves `GET /__mcp/router/plan/{session_id}` to loopback clients, returning the session's current plan (selected tool names, decision traces, version, mode) as JSON. Non-loopback clients and configurations with payload logging off receive `404`, so its existence is not observable in production by default. Because plans reveal authorization state, only enable this behind access controls.

#### Multi-Instance Note

Session plans are stored in-process. Across multiple Pixiu instances, route the same `Mcp-Session-Id` to the same instance (sticky session, e.g. a load-balancer hash on the header) so a session sees a consistent plan. A shared/distributed plan store is a planned future enhancement.

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

---

### Using Nacos as MCP Server Registry

Pixiu supports dynamic discovery and management of MCP tool configurations through Nacos 3.0+. By using Nacos as a registry center, you can centrally manage MCP tool definitions and achieve dynamic configuration updates without restarting the gateway.

#### Adapter Configuration (`adapters`)

To enable Nacos integration, you need to add an `adapters` section to your configuration file. The adapter is responsible for connecting to the Nacos registry and subscribing to MCP service configurations.

```yaml
adapters:
  - id: "mcp-nacos-adapter"
    name: "dgp.adapter.mcpserver"
    config:
      registries:
        nacos:
          protocol: "nacos"
          address: "127.0.0.1:8848"
          timeout: "5s"
          username: "nacos"
          password: "nacos"
```

#### Adapter Configuration Fields

`id`

- **Type**: `string`
- **Description**: A unique identifier for the adapter. Used to identify this adapter instance in logs and monitoring.

`name`

- **Type**: `string`
- **Description**: The adapter type name. For MCP server Nacos integration, must be `dgp.adapter.mcpserver`.

`config`

- **Type**: `object`
- **Description**: The specific configuration for the adapter, including registry connection information.

##### Registry Configuration (`registries`)

`registries` is a key-value mapping where the key is the registry name (e.g., `nacos`) and the value is the configuration object for that registry.

`protocol`

- **Type**: `string`
- **Description**: The registry protocol type. Currently supports `nacos`.

`address`

- **Type**: `string`
- **Description**: The Nacos server address in the format `host:port`. For example, `127.0.0.1:8848`.

`timeout`

- **Type**: `string`
- **Description**: Connection timeout duration. Supported time units include `s` (seconds), `ms` (milliseconds), etc. For example, `5s` means 5 seconds.

`username`

- **Type**: `string`
- **Description**: Nacos authentication username. Required if the Nacos server has authentication enabled.

`password`

- **Type**: `string`
- **Description**: Nacos authentication password. Required if the Nacos server has authentication enabled.

`namespace` (optional)

- **Type**: `string`
- **Description**: Nacos namespace ID. Used for environment isolation. If not specified, the default namespace is used.

`group` (optional)

- **Type**: `string`
- **Description**: Nacos service group. Used for service group management. If not specified, the default group is used.

---

### Complete Nacos Integration Configuration Example

The following example demonstrates a complete configuration using Nacos as the MCP server configuration source. The gateway automatically retrieves tool definitions from Nacos and routes requests based on the configuration.

```yaml
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
                  # All MCP requests route to the protected cluster
                  - match:
                      prefix: "/"
                    route:
                      cluster: "mcp-protected"
                      cluster_not_found_response_code: 505
              http_filters:
                # MCP Server Filter
                - name: "dgp.filter.mcp.mcpserver"
                  config:
                    server_info:
                      name: "MCP Nacos Example Server"
                      version: "1.0.0"
                      description: "MCP Server with tools dynamically loaded from Nacos"
                      instructions: "Tool configurations for this server are centrally managed by Nacos"

                # Downstream HTTP Proxy
                - name: "dgp.filter.http.httpproxy"

  clusters:
    # Virtual cluster for routing rules
    - name: "mcp-protected"
      type: "STATIC"
      lb_policy: "ROUND_ROBIN"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 8081

# Nacos Adapter Configuration
adapters:
  - id: "mcp-nacos-adapter"
    name: "dgp.adapter.mcpserver"
    config:
      registries:
        nacos:
          protocol: "nacos"
          address: "127.0.0.1:8848"
          timeout: "5s"
          username: "nacos"
          password: "nacos"
          namespace: ""  # Optional: use default namespace
          group: "DEFAULT_GROUP"  # Optional: use default group
```

---

### Usage Steps

#### 1. Prepare Nacos Environment

Ensure you have Nacos 3.0 or higher installed and running. You can access the Nacos console at `http://<nacos-server-ip>:8848/nacos`.

#### 2. Configure MCP Service in Nacos

1. **Login to Nacos Console**
2. **Navigate to MCP Management**: Find and click "MCP Management" in the left sidebar
3. **Create MCP Server**:
   - Click "MCP List" → "Create MCP Server"
   - **Type**: Select `streamable`
   - **Tools**: Select "Import from OpenAPI" and upload your OpenAPI specification file

4. **Verify and Correct Configuration**:
   - After uploading successfully, Nacos will automatically parse the OpenAPI file and generate the tool list
   - **Important**: Check that all tool backend addresses are correct (Nacos 3.0 may have path parsing issues)
   - Ensure backend addresses are in the format `http://host:port`, not `http:/host:port`

5. **Publish Service**: After confirming all configurations are correct, click "Publish"

> **Note**: The current version of Pixiu only supports connecting to a single MCP Server instance.

#### 3. Start Pixiu Gateway

Start Pixiu using the configuration file containing the Nacos adapter configuration:

```bash
cd /path/to/dubbo-go-pixiu
go run cmd/pixiu/*.go gateway start -c /path/to/your/config.yaml
```

After starting, Pixiu will:

- Connect to the Nacos registry
- Subscribe to MCP service configurations
- Dynamically load tool definitions
- Automatically handle configuration updates (without restart)

#### 4. Verify Integration

You can verify that the Nacos integration is working correctly through:

1. **Check Logs**: Review Pixiu startup logs to confirm successful connection to Nacos
2. **Test Tool Calls**: Use an MCP client (such as MCP Inspector) to connect to `http://localhost:8888/mcp` and test tool invocations
3. **Dynamic Update Test**: Modify tool configurations in the Nacos console and verify that changes take effect automatically

---

### Best Practices

1. **Environment Isolation**: Use Nacos namespace features to isolate configurations for different environments (development, testing, production)
2. **Configuration Backup**: Regularly back up MCP configurations in Nacos to prevent accidental loss
3. **Monitoring and Alerting**: Configure Nacos connection status monitoring to detect connection issues promptly
4. **Canary Releases**: Leverage Nacos configuration management capabilities to implement canary releases of tool configurations

---

### Example Reference

Complete usage examples and configuration files can be found in the `mcp/nacos` directory of the [dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu-samples) project.
