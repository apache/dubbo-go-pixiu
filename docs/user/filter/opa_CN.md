# OPA 过滤器 (dgp.filter.http.opa)

[English](opa.md) | 中文

---

## 中文

### 概述
`dgp.filter.http.opa` 过滤器通过 Rego 策略将授权决策委托给 Open Policy Agent (OPA)。该过滤器评估每个 HTTP 请求并根据 Rego 策略决定是否允许或拒绝请求。

过滤器支持两种运行模式：
- **服务器模式（推荐）**：通过 HTTP API 调用独立的 OPA 服务器进行策略评估，支持集中式策略管理和热更新。
- **嵌入模式**：策略通过内联 Rego 模块加载，并使用 OPA 的内置查询引擎进行评估。

### 运行模式

#### 服务器模式（推荐）
- 通过 REST API 调用外部 OPA 服务器 (`server_url`)。
- 支持集中式策略管理和动态更新，无需重启 Pixiu。
- 适合生产环境和大规模部署。
- 配置简单，只需指定 OPA 服务器地址和决策路径。

#### 嵌入模式（向后兼容）
- 从配置项 `policy` 读取 **Rego 模块源码字符串**。
- 从配置项 `entrypoint` 读取 **Rego 查询字符串**。
- 策略在 Pixiu 进程内评估，更新策略需要重启。
- 适合简单场景或测试环境。

> **自动模式选择**：
> - 如果配置了 `server_url`，则使用**服务器模式**（优先）。
> - 如果未配置 `server_url` 但配置了 `policy`，则使用**嵌入模式**。

### 配置结构

#### 服务器模式配置（推荐）

```yaml
filters:
  - name: dgp.filter.httpconnectionmanager
    config:
      route_config:
        # ... 你的路由
      http_filters:
        - name: dgp.filter.http.opa
          config:
            # OPA 服务器地址
            server_url: "http://opa-server:8181"
            # 决策路径（OPA REST API 路径）
            decision_path: "/v1/data/http/authz/allow"
            # 请求超时时间（毫秒），默认 100ms
            timeout_ms: 100
            # 可选：Bearer Token 认证
            # bearer_token: "your-secret-token"
        # HTTP proxy 过滤器应该在 OPA 过滤器之后
        - name: dgp.filter.http.proxy
          config:
            # ... proxy config
```

#### 嵌入模式配置（向后兼容）

```yaml
filters:
  - name: dgp.filter.httpconnectionmanager
    config:
      route_config:
        # ... 你的路由
      http_filters:
        - name: dgp.filter.http.opa
          config:
            policy: |
              package http.authz

              default allow = false

              allow {
                input.method == "GET"
                input.path == "/status"
              }
            entrypoint: "data.http.authz.allow"
        # HTTP proxy 过滤器应该在 OPA 过滤器之后
        - name: dgp.filter.http.proxy
          config:
            # ... proxy config
```



#### 字段说明

**服务器模式字段：**

- **`server_url`**（字符串，服务器模式必填）
  - **含义：** OPA 服务器地址，如 `http://opa-server:8181` 或 `https://opa.example.com:8181`。
  - **数据类型：** `string`。
  - **示例：** `"http://localhost:8181"`。

- **`decision_path`**（字符串，服务器模式必填）
  - **含义：** OPA REST API 决策路径，格式为 `/v1/data/<package>/<rule>`。
  - **数据类型：** `string`。
  - **示例：** `"/v1/data/http/authz/allow"`。

- **`timeout_ms`**（整数，可选）
  - **含义：** 请求 OPA 服务器的超时时间（毫秒）。
  - **数据类型：** `int`。
  - **默认值：** `100`（100 毫秒）。

- **`bearer_token`**（字符串，可选）
  - **含义：** 用于 OPA 服务器认证的 Bearer Token。
  - **数据类型：** `string`。
  - **说明：** 如果 OPA 服务器需要认证，可配置此字段。

**嵌入模式字段（向后兼容）：**

- **`policy`**（字符串，嵌入模式必填）
  - **含义：** **Rego 模块源码**（内联字符串）。通过 `rego.Module("policy.rego", policy)` 加载。
  - **数据类型：** `string`（建议使用 YAML 多行格式 `|`）。

- **`entrypoint`**（字符串，嵌入模式必填）
  - **含义：** 传给 `rego.Query(...)` 的 **Rego 查询字符串**，应为合法查询，如 `data.<package>.<rule>`（如 `data.http.authz.allow`）。
  - **数据类型：** `string`。

#### 判定约定

- 如果查询结果集合非空且首个表达式值为 **`true`**，则请求放行。
- 否则（空结果或值≠`true`），请求被拒绝。

### 策略输入（`input`）

过滤器将 HTTP 请求转换为以下键值对（与当前实现一致，策略编写时请进行空值检查）：

```
input.method       # 请求方法，字符串
input.path         # URL Path，字符串
input.headers      # map[string][]string
input.client_ip    # 字符串
input.query        # map[string][]string（URL 查询参数）
input.host         # 字符串
input.remote_addr  # 字符串
input.user_agent   # 字符串
input.route        # 路由条目对象（结构可能变化）
input.api          # API 对象（结构可能变化）
input.params       # 路由参数 map
```

**重要提示：HTTP Header 规范化**

Go 的 `net/http` 包会自动将 HTTP header 名称规范化为 **Title-Case** 格式（如 `X-Api-Key`、`Content-Type`）。在策略中访问 headers 时，请使用规范化后的键名：

- ✅ 正确：`input.headers["X-Api-Key"]` 或 `input.headers["Content-Type"]`
- ❌ 错误：`input.headers["x-api-key"]` 或 `input.headers["x_api_key"]`

**示例：**
```rego
# 正确的 header 访问方式
allow {
    input.headers["X-Api-Key"][0] == "secret"
    input.headers["Content-Type"][0] == "application/json"
}
```

### OPA 服务器部署（服务器模式）

#### Docker 部署

**1. 启动 OPA 服务器**

```bash
docker run -d \
  --name opa-server \
  -p 8181:8181 \
  openpolicyagent/opa:latest \
  run --server --addr :8181 --log-level info
```

**2. 上传策略到 OPA 服务器**

```bash
# 创建策略文件 policy.rego
cat > policy.rego <<EOF
package http.authz

default allow = false

allow {
    input.method == "GET"
    input.path == "/status"
}
EOF

# 上传策略
curl -X PUT http://localhost:8181/v1/policies/authz \
  -H "Content-Type: text/plain" \
  --data-binary @policy.rego
```

**3. 测试策略**

```bash
curl -X POST http://localhost:8181/v1/data/http/authz/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "method": "GET",
      "path": "/status"
    }
  }'
```

#### Docker Compose 部署（推荐）

```yaml
version: '3.8'

services:
  opa-server:
    image: openpolicyagent/opa:latest
    container_name: opa-server
    ports:
      - "8181:8181"
    command:
      - "run"
      - "--server"
      - "--addr=:8181"
      - "--log-level=info"
    volumes:
      - ./opa-policies:/policies:ro
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8181/health"]
      interval: 10s
      timeout: 3s
      retries: 3
    restart: unless-stopped

  pixiu:
    image: apache/dubbo-go-pixiu:latest
    depends_on:
      opa-server:
        condition: service_healthy
    volumes:
      - ./configs:/configs:ro
    ports:
      - "8888:8888"
```

#### 策略热更新

OPA 服务器支持在运行时更新策略，无需重启：

```bash
# 更新策略
curl -X PUT http://localhost:8181/v1/policies/authz \
  -H "Content-Type: text/plain" \
  --data-binary @new-policy.rego
```

### 最小可用示例

**服务器模式示例**

**1）仅允许 GET /status（使用 OPA 服务器）**

```yaml
- name: dgp.filter.http.opa
  config:
    server_url: "http://opa-server:8181"
    decision_path: "/v1/data/http/authz/allow"
    timeout_ms: 100
```

对应的 OPA 策略（上传到服务器）：

```rego
package http.authz

default allow = false

allow {
    input.method == "GET"
    input.path == "/status"
}
```

**2）基于请求头校验（使用 OPA 服务器）**

```yaml
- name: dgp.filter.http.opa
  config:
    server_url: "http://opa-server:8181"
    decision_path: "/v1/data/http/authz/allow"
    timeout_ms: 100
    bearer_token: "your-secret-token"
```

对应的 OPA 策略：

```rego
package http.authz

default allow = false

allow {
    input.headers["X-Api-Key"][0] == "secret"
}
```

**嵌入模式示例**

**1）仅允许 GET /status（嵌入模式）**

```yaml
- name: dgp.filter.http.opa
  config:
    policy: |
      package http.authz
      default allow = false
      allow { input.method == "GET"; input.path == "/status" }
    entrypoint: "data.http.authz.allow"
```

**2）基于请求头校验（嵌入模式）**

```yaml
- name: dgp.filter.http.opa
  config:
    policy: |
      package http.authz
      default allow = false
      allow {
        input.headers["x-api-key"][0] == "secret"
      }
    entrypoint: "data.http.authz.allow"
```

### 限制与说明

**返回格式支持**

OPA 策略支持多种返回格式，过滤器会自动识别：

**格式 1：布尔值（简单场景）**
```rego
package http.authz
default allow = false
allow {
    input.method == "GET"
}
```

**格式 2：对象格式 - allow 字段（OPA 常见习惯，推荐）**
```rego
package http.authz

default decision = {"allow": false}

decision = {"allow": true, "reason": "admin user"} {
    input.headers["X-Role"][0] == "admin"
}
```

**格式 3：对象格式 - result 字段**
```rego
package http.authz

decision = {"result": true, "metadata": {"user": "alice"}} {
    input.headers["X-Api-Key"][0] == "secret"
}
```

过滤器会按以下优先级提取决策结果：
1. 直接布尔值 → 使用该值
2. 对象的 `allow` 字段 → 使用 `allow` 的布尔值
3. 对象的 `result` 字段 → 使用 `result` 的布尔值
4. 无法识别 → 默认拒绝（返回 403）

---

**服务器模式：**
- **网络延迟**：由于需要调用外部 OPA 服务器，会引入网络延迟（通常 5-50ms）。
- **OPA 服务器可用性**：需要确保 OPA 服务器高可用，建议部署多实例。
- **超时配置**：默认超时 100ms 可能在跨服务/云环境下偏紧，生产环境建议根据网络延迟调整（如 300-500ms）。

**嵌入模式：**
- **没有自定义拒绝响应**：过滤器不会将策略输出映射到 HTTP 状态码或响应体。
- **仅支持内联模块加载**：策略来自配置字符串，不读取外部文件。
- **策略更新需重启**：修改策略后需要重启 Pixiu 服务。

### 故障排查

**服务器模式：**
- **连接失败**：检查 `server_url` 配置是否正确，确保 OPA 服务器正在运行并可访问。
- **超时错误**：增加 `timeout_ms` 值，或检查 OPA 服务器性能。
- **认证失败**：如果 OPA 服务器需要认证，确保 `bearer_token` 配置正确。
- **决策路径错误**：确保 `decision_path` 与 OPA 服务器中的策略路径匹配。
- **意外拒绝**：使用 `curl` 直接测试 OPA 服务器的决策 API，验证策略逻辑。

**嵌入模式：**
- **意外拒绝**：检查查询是否正确（如 `data.http.authz.allow`），并确保策略在给定的 `input` 下返回 **`true`**。
- **策略编译错误**：在嵌入策略之前，先使用 `opa eval` 本地验证 Rego 语法。
- **空结果或类型不符**：请检查 `headers`/`query`，确保路径和方法匹配。

**通用排查：**
- **查看日志**：Pixiu 会记录 OPA 评估的详细日志，包括错误信息和决策结果。
- **测试策略**：先在 OPA 服务器或本地使用 `opa eval` 测试策略，确保逻辑正确。