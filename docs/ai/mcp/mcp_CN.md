## MCP (Model Context Protocol) 网关配置

[English](./mcp.md) | 中文

本文档解释了如何在您的网关中配置 MCP (Model Context Protocol) 过滤器，使您能够安全地将后端 HTTP API 暴露为可供 AI Agent 调用的"工具"。

### 简介

MCP (Model Context Protocol) 是一个智能桥梁，连接 AI Agent 与您现有的后端服务。它将一个简单、统一的协议动态转换为标准的 HTTP 请求，允许 Agent 与您的 API 进行交互，就好像它们是本地函数或工具一样。这种方法简化了 Agent 的开发，并为安全性、控制和可观察性提供了一个中心化的管理点。

设置 MCP 端点主要涉及两个过滤器：

1.  **`dgp.filter.mcp.mcpserver`**: 核心的 MCP 服务器过滤器，用于定义服务器的身份并将后端 API 暴露为工具。
2.  **`dgp.filter.http.auth.mcp`**: 一个可选但推荐使用的安全过滤器，它使用 OAuth 2.0 和基于 JWT 的授权来保护 MCP 端点。

---

### MCP 服务器过滤器 (`dgp.filter.mcp.mcpserver`) 配置

此过滤器是 MCP 网关的核心。它负责定义哪些工具可用，以及它们如何映射到您的后端 HTTP 服务。

#### 服务器信息 (`server_info`)

`server_info` 块提供了有关您的 MCP 服务器的元数据，这对于服务发现和诊断非常有用。

```yaml
server_info:
  name: "MCP OAuth 示例服务器"
  version: "1.0.0"
  description: "一个用于工具演示的受 OAuth 保护的 MCP 服务器"
  instructions: "使用读/写令牌通过 MCP 与模拟服务器 API 进行交互"
```

-   **`name` (`string`)**: MCP 服务器的显示名称。
-   **`version` (`string`)**: 服务器的版本。
-   **`description` (`string`)**: 服务器用途的简要描述。
-   **`instructions` (`string`)**: 为客户端或开发者提供的关于如何与暴露的工具进行交互的说明。

#### 工具配置 (`tools`)

`tools` 部分是一个数组，其中每个项目定义了 MCP 网关暴露的单个工具。每个工具对应一个特定的后端 HTTP API 端点。

```yaml
tools:
  - name: "get_user"
    description: "通过 ID 获取用户信息，可选择是否包含个人资料详情"
    cluster: "mock-server"
    # ... 请求和参数配置 ...
```

-   **`name` (`string`)**: 工具的唯一名称。这是 AI Agent 调用它时将使用的标识符。
-   **`description` (`string`)**: 对工具功能的清晰、简洁的描述。这对于 LLM 理解工具的能力至关重要。
-   **`cluster` (`string`)**: 将处理此工具请求的上游集群的名称。

##### 请求定义 (`request`)

`request` 对象指定了当工具被调用时，网关将向上游集群发出的 HTTP 请求的详细信息。

```yaml
request:
  method: "GET"
  path: "/api/users/{id}"
  timeout: "10s"
  headers:
    Content-Type: "application/json"
```

-   **`method` (`string`)**: HTTP 方法 (例如, `GET`, `POST`, `PUT`, `DELETE`)。
-   **`path` (`string`)**: 上游服务的请求路径。您可以使用像 `{arg_name}` 这样的占位符来表示路径参数。
-   **`timeout` (`string`)**: 上游请求的超时时间 (例如, `5s`, `100ms`)。
-   **`headers` (`object`)**: 一个键值对映射，定义了要包含在上游请求中的静态 HTTP 头。

##### 参数定义 (`args`)

`args` 数组定义了工具接受的参数。此模式允许网关验证传入的参数，并将它们正确地放入上游 HTTP 请求中。

```yaml
args:
  - name: "id"
    type: "integer"
    in: "path"
    description: "要检索的用户 ID"
    required: true
```

每个参数对象包含以下字段：

| 字段          | 类型      | 描述                                                                                                                              |
|---------------|-----------|-----------------------------------------------------------------------------------------------------------------------------------|
| `name`        | `string`  | 参数的名称。                                                                                                                      |
| `type`        | `string`  | 参数的数据类型 (`string`, `integer`, `number`, `boolean`)。                                                                     |
| `in`          | `string`  | 指定参数在 HTTP 请求中的位置：`path`、`query` 或 `body`。                                                                          |
| `description` | `string`  | 参数用途的详细描述，帮助 LLM 正确使用它。                                                                                         |
| `required`    | `boolean` | 参数是否为必需。默认为 `false`。                                                                                                  |
| `default`     | `any`     | 如果未提供参数，则使用的默认值。                                                                                                  |
| `enum`        | `array`   | 允许值的数组，提供一种验证形式。                                                                                                  |

---

### MCP 认证过滤器 (`dgp.filter.http.auth.mcp`) 配置

此过滤器为您的 MCP 端点增加了一个安全层，确保只有经过身份验证和授权的客户端才能调用工具。它根据配置的身份提供者验证客户端提供的 JWT。

#### 资源元数据 (`resource_metadata`)

此部分定义了受保护的资源，并指向可以授予其访问权限的授权服务器。

```yaml
resource_metadata:
  path: "/.well-known/oauth-protected-resource/mcp"
  resource: "http://localhost:8888/mcp"
  authorization_servers:
    - "http://localhost:9000"
```

-   **`path` (`string`)**: 资源元数据的暴露路径，遵循可发现性标准。
-   **`resource` (`string`)**: 受保护资源的标识符 (通常是 MCP 端点本身的 URL)。
-   **`authorization_servers` (`array` of `string`)**: 受信任的授权服务器 URL 列表。

#### 提供者 (`providers`)

此部分定义了受信任的 JWT 签发者 (身份提供者)。网关将使用此信息来验证传入 JWT 的签名。

```yaml
providers:
  - name: "local"
    issuer: "http://localhost:9000"
    jwks: "http://localhost:9000/.well-known/jwks.json"
```

-   **`name` (`string`)**: 此提供者配置的唯一名称。
-   **`issuer` (`string`)**: JWT 中预期的 `iss` (签发者) 声明。这必须与签发者的标识符匹配。
-   **`jwks` (`string`)**: JSON Web Key Set (JWKS) 端点的 URL，用于发布验证 JWT 签名的公钥。

#### 规则 (`rules`)

规则将认证策略连接到特定的上游集群。

```yaml
rules:
  - cluster: "mcp-protected"
```

-   **`cluster` (`string`)**: 此认证和授权策略适用的集群名称。

---

### 完整配置示例

此示例演示了 MCP 网关的完整设置。它包括一个带有多个工具的 MCP 服务器，并由 MCP 认证过滤器保护。

```yaml
# 完整的 pixiu_mcp_auth_test.yaml 示例
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
                  # 所有流量都路由到一个虚拟集群
                  - match:
                      prefix: "/"
                    route:
                      cluster: "mcp-protected"
              http_filters:
                # (可选) MCP 授权过滤器，用于保护端点
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

                # 核心 MCP 服务器过滤器
                - name: "dgp.filter.mcp.mcpserver"
                  config:
                    server_info:
                      name: "MCP OAuth 示例服务器"
                      version: "1.0.0"
                      description: "一个用于工具演示的受 OAuth 保护的 MCP 服务器"
                      instructions: "使用适当的令牌通过 MCP 与模拟服务器 API 进行交互"
                    
                    tools:
                      # 工具 1: 通过 ID 获取用户
                      - name: "get_user"
                        description: "通过 ID 获取用户信息"
                        cluster: "mock-server"
                        request:
                          method: "GET"
                          path: "/api/users/{id}"
                          timeout: "10s"
                        args:
                          - name: "id"
                            type: "integer"
                            in: "path"
                            description: "要检索的用户 ID"
                            required: true

                      # 工具 2: 创建新用户
                      - name: "create_user"
                        description: "创建一个新用户帐户"
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
                            description: "用户的全名"
                            required: true
                          - name: "email"
                            type: "string"
                            in: "body"
                            description: "用户的电子邮件地址"
                            required: true
                
                # 标准的下游 HTTP 代理过滤器
                - name: "dgp.filter.http.httpproxy"

  clusters:
    # 上游后端服务
    - name: "mock-server"
      type: "STATIC"
      lb_policy: "ROUND_ROBIN"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 8081

    # 用于路由规则的虚拟集群
    - name: "mcp-protected"
      type: "STATIC"
      lb_policy: "ROUND_ROBIN"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 8081
```

---

### 使用 Nacos 作为 MCP 服务器注册中心

Pixiu 支持通过 Nacos 3.0+ 动态发现和管理 MCP 工具配置。通过使用 Nacos 作为注册中心，您可以集中管理 MCP 工具定义，并实现动态配置更新，而无需重启网关。

#### Adapter 配置 (`adapters`)

要启用 Nacos 集成，您需要在配置文件中添加 `adapters` 部分。适配器负责连接到 Nacos 注册中心并订阅 MCP 服务配置。

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

#### Adapter 配置字段说明

`id`

- **类型**: `string`
- **描述**: 适配器的唯一标识符。用于在日志和监控中识别此适配器实例。

`name`

- **类型**: `string`
- **描述**: 适配器类型名称。对于 MCP 服务器的 Nacos 集成，必须使用 `dgp.adapter.mcpserver`。

`config`

- **类型**: `object`
- **描述**: 适配器的具体配置，包含注册中心连接信息。

##### 注册中心配置 (`registries`)

`registries` 是一个键值对映射，其中键是注册中心的名称（如 `nacos`），值是该注册中心的配置对象。

`protocol`

- **类型**: `string`
- **描述**: 注册中心协议类型。目前支持 `nacos`。

`address`

- **类型**: `string`
- **描述**: Nacos 服务器地址，格式为 `host:port`。例如 `127.0.0.1:8848`。

`timeout`

- **类型**: `string`
- **描述**: 连接超时时间。支持的时间单位包括 `s`（秒）、`ms`（毫秒）等。例如 `5s` 表示 5 秒。

`username`

- **类型**: `string`
- **描述**: Nacos 认证用户名。如果 Nacos 服务器启用了认证，此字段为必填。

`password`

- **类型**: `string`
- **描述**: Nacos 认证密码。如果 Nacos 服务器启用了认证，此字段为必填。

`namespace` (可选)

- **类型**: `string`
- **描述**: Nacos 命名空间 ID。用于环境隔离。如果不指定，则使用默认命名空间。

`group` (可选)

- **类型**: `string`
- **描述**: Nacos 服务分组。用于服务分组管理。如果不指定，则使用默认分组。

---

### 完整的 Nacos 集成配置示例

以下示例展示了一个使用 Nacos 作为 MCP 服务器配置源的完整配置。网关会自动从 Nacos 获取工具定义，并根据配置动态路由请求。

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
                  # 所有 MCP 请求路由到受保护的集群
                  - match:
                      prefix: "/"
                    route:
                      cluster: "mcp-protected"
                      cluster_not_found_response_code: 505
              http_filters:
                # MCP 服务器过滤器
                - name: "dgp.filter.mcp.mcpserver"
                  config:
                    server_info:
                      name: "MCP Nacos 示例服务器"
                      version: "1.0.0"
                      description: "从 Nacos 动态加载工具的 MCP 服务器"
                      instructions: "此服务器的工具配置由 Nacos 集中管理"

                # 下游 HTTP 代理
                - name: "dgp.filter.http.httpproxy"

  clusters:
    # 虚拟集群，用于路由规则
    - name: "mcp-protected"
      type: "STATIC"
      lb_policy: "ROUND_ROBIN"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 8081

# Nacos 适配器配置
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
          namespace: ""  # 可选：使用默认命名空间
          group: "DEFAULT_GROUP"  # 可选：使用默认分组
```

---

### 使用步骤

#### 1. 准备 Nacos 环境

确保您已安装并启动 Nacos 3.0 或更高版本。您可以通过访问 `http://<nacos-server-ip>:8848/nacos` 来访问 Nacos 控制台。

#### 2. 在 Nacos 中配置 MCP 服务

1. **登录 Nacos 控制台**
2. **进入 MCP 管理**：在左侧菜单栏找到并点击 "MCP管理"
3. **创建 MCP Server**：
   - 点击 "MCP列表" → "创建MCP Server"
   - **类型**：选择 `streamable`
   - **工具(Tools)**：选择 "从OpenAPI导入"，然后上传您的 OpenAPI 规范文件

4. **验证并修正配置**：
   - 上传成功后，Nacos 会自动解析 OpenAPI 文件并生成工具列表
   - **重要**：检查所有工具的后端地址是否正确（Nacos 3.0 版本可能存在路径解析问题）
   - 确保后端地址格式为 `http://host:port`，而不是 `http:/host:port`

5. **发布服务**：确认所有配置无误后，点击 "发布"

> **注意**：当前版本的 Pixiu 仅支持连接到单个 MCP Server 实例。

#### 3. 启动 Pixiu 网关

使用包含 Nacos 适配器配置的配置文件启动 Pixiu：

```bash
cd /path/to/dubbo-go-pixiu
go run cmd/pixiu/*.go gateway start -c /path/to/your/config.yaml
```

启动后，Pixiu 会：

- 连接到 Nacos 注册中心
- 订阅 MCP 服务配置
- 动态加载工具定义
- 自动处理配置更新（无需重启）

#### 4. 验证集成

您可以通过以下方式验证 Nacos 集成是否正常工作：

1. **检查日志**：查看 Pixiu 启动日志，确认已成功连接到 Nacos
2. **测试工具调用**：使用 MCP 客户端（如 MCP Inspector）连接到 `http://localhost:8888/mcp` 并测试工具调用
3. **动态更新测试**：在 Nacos 控制台中修改工具配置，验证更改是否自动生效

---

### 最佳实践

1. **环境隔离**：使用 Nacos 的命名空间功能来隔离不同环境（开发、测试、生产）的配置
2. **配置备份**：定期备份 Nacos 中的 MCP 配置，防止意外丢失
3. **监控与告警**：配置 Nacos 连接状态监控，及时发现连接问题
4. **灰度发布**：利用 Nacos 的配置管理能力，实现工具配置的灰度发布

---

### 示例参考

完整的使用示例和配置文件可以在 [dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu-samples) 项目的 `mcp/nacos` 目录中找到。
