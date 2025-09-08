# OPA 过滤器 (dgp.filter.http.opa)

[English](opa.md) | 中文

---

## 中文

### 概述
`dgp.filter.http.opa` 过滤器通过 Rego 策略将授权决策委托给 Open Policy Agent (OPA)。该过滤器评估每个 HTTP 请求并根据 Rego 策略决定是否允许或拒绝请求。策略通过内联 Rego 模块加载，并使用 OPA 的内置查询引擎进行评估。

### 实际行为
- 从配置项 `policy` 读取 **Rego 模块源码字符串**。
- 从配置项 `entrypoint` 读取 **Rego 查询字符串**。
- 每次请求构造 `input` 对象（见下），并评估该查询。
- 如果查询结果为 `true`，则放行请求；否则拒绝请求。

> 目前过滤器**不支持**：外部文件或 URI 加载、自定义拒绝状态码或返回自定义错误体等。

### 配置结构
将过滤器添加到 HTTP 连接管理器的 `http_filters` 列表中：

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
        # HTTP proxy 过滤器应该在OPA 过滤器之后
        - name: dgp.filter.http.proxy
          config:
          	# ... proxy config
```



#### 字段说明

- **`policy`**（字符串，必填）
  - **含义：** **Rego 模块源码**（内联字符串）。通过 `rego.Module("policy.rego", policy)` 加载。
  - **数据类型：** `string`（建议使用 YAML 多行格式 `|`）。
  - **说明：** 当前版本不支持外部文件路径或 bundle URI。
- **`entrypoint`**（字符串，必填）
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

### 最小可用示例

**1）仅允许 GET /status**

```yaml
- name: dgp.filter.http.opa
  config:
    policy: |
      package http.authz
      default allow = false
      allow { input.method == "GET"; input.path == "/status" }
    entrypoint: "data.http.authz.allow"
```

**2）基于请求头校验**

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

- **返回类型必须是布尔值**：只有 `true` 会被视为放行；对象（如 `{allow: true}`）不会被特殊处理。
- **没有自定义拒绝响应**：过滤器不会将策略输出映射到 HTTP 状态码或响应体。
- **仅支持内联模块加载**：策略来自配置字符串，不读取外部文件。

### 故障排查

- **意外拒绝**：检查查询是否正确（如 `data.http.authz.allow`），并确保策略在给定的 `input` 下返回 **`true`**。
- **策略编译错误**：在嵌入策略之前，先使用 `opa eval` 本地验证 Rego 语法。
- **空结果或类型不符**：请检查 `headers`/`query`，确保路径和方法匹配。