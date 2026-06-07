# OpenAPI 请求校验过滤器

[English](openapi.md) | 中文

---

## 概述

Pixiu 可以在 `dgp.filter.http.openapi` 中加载本地 OpenAPI 3.0/3.1 文件，并在请求转发到上游之前校验命中的请求。对于
OpenAPI 3.2 文档，当前仅覆盖这个 filter 已接入的标准 HTTP operation。

官方参考：

- [libopenapi](https://github.com/pb33f/libopenapi)
- [libopenapi validation](https://pb33f.io/libopenapi/validation/)
- [libopenapi validator](https://github.com/pb33f/libopenapi-validator)

如果校验失败，Pixiu 会直接返回本地 `400 Bad Request`，并停止后续过滤链。

## 当前 filter 已接入

- 本地 OpenAPI 3.0/3.1 文件加载，以及使用标准 HTTP operation 的 OpenAPI 3.2 文档
- 按 OpenAPI path 和 method 匹配请求，包括 `/users/{id}` 这类模板路径
- path、query、header 参数校验
- JSON request body 校验
- 由 SDK 执行的 OpenAPI schema 约束，包括 `required`、`type`、`enum`、`minimum`、`maximum`、`minLength`、`maxLength`

当前这个 filter 使用 `libopenapi` 解析 OpenAPI 文档，并使用 `libopenapi-validator` 校验进入过滤器的请求。

## 当前 filter 未接入

- response validation
- 路由创建或 `api_config` 路由匹配
- OpenAPI `security` 校验；鉴权请使用专门的认证或授权 filter
- OpenAPI 3.2 中超出当前标准 HTTP request method 覆盖范围的 operation
- admin 或配置中心分发 OpenAPI 文件

## 配置示例

```yaml
- name: dgp.filter.http.apiconfig
  config:
    path: configs/api_config.yaml

- name: dgp.filter.http.openapi
  config:
    path: configs/openapi_users.yaml
    # 可选，默认 1048576 字节。
    max_request_body_bytes: 1048576
```

`dgp.filter.http.apiconfig` 和 `dgp.filter.http.openapi` 是两个独立 filter。`apiconfig` 负责匹配 Pixiu API 路由，并把
API 元信息写入请求上下文；`openapi` 只校验 OpenAPI 文件中声明过的 operation。
`apiconfig` 中已废弃的 OpenAPI 校验配置键，包括 `openapi_path` 和 `enable_openapi_validation`，只要出现就会被拒绝；
即使配置为 `""` 或 `false` 也一样。请改用独立的 `dgp.filter.http.openapi` filter。

如果请求的 path 和 method 没有在 OpenAPI 文件中声明，这个 filter 会跳过校验并放行请求。如果 operation 已声明但请求
不满足参数或 body schema，Pixiu 会返回 `400 Bad Request`。

## 说明

- `libopenapi-validator` 是独立的可选模块，`libopenapi` 负责解析和建模。
- 参数级校验覆盖 `path`、`query`、`header` 上常见的标量类型约束。
- 上面这些关键词来自 OpenAPI schema，是由 SDK 路径执行的，不是仓库里再自定义一套验证器。
- 当前 filter 关闭了 OpenAPI `security` 校验，鉴权仍由 JWT、OPA、SAML 或其他专门的认证/授权 filter 负责。
- OpenAPI `path` 配置必须使用相对路径。绝对路径、`..` 父目录跳转和敏感 base 目录会在 filter 启动阶段被拒绝。启动时会解析符号链接（symlink），防止通过相对路径 + 符号链接绕过敏感目录检查。
- request body 在 schema 校验前会受 `max_request_body_bytes` 限制，避免 validator 在网关热路径上读取无上限 JSON 或 chunked body。
- 不配置 `max_request_body_bytes` 或配置为 `0` 时，会使用默认的 `1048576` 字节限制。
- OpenAPI 文件里的相对引用会按文件所在目录解析，所以使用本地文件加载时可以保留本地 `$ref` 路径。
- 无效 OpenAPI 文档，包括无法解析的 `$ref`，会在 filter 启动阶段失败，不会继续使用部分构建出来的 validator model。
- `libopenapi-validator` 也提供响应和文档校验 API，但当前这个 filter 只调用了请求校验入口。

## 运行流程

1. Pixiu 在 `openapi.Apply()` 阶段加载 OpenAPI 文件并构建 SDK 校验器。
2. 请求进入 `openapi.Decode()`。
3. Pixiu 检查 OpenAPI 文档是否声明了请求 path 和 method。
4. 如果 operation 未声明，跳过校验并继续后续过滤链。
5. 如果 operation 已声明，Pixiu 通过 `libopenapi-validator` 执行请求校验。
6. 校验通过，请求继续流向后续 filter。
7. 校验失败，Pixiu 直接返回 `400 Bad Request`。

### HEAD 请求

OpenAPI spec 不把 HEAD 视为标准 HTTP 方法。当路径只声明了 `GET`（没有显式 `head` operation）时，SDK 的 `FindPath` 会把 HEAD 当作未声明方法处理，filter 跳过校验并放行。
如果需要对 HEAD 请求做校验，可以在 OpenAPI spec 中显式声明 `head` operation。

## 请求示例

假设配置引用的是 `configs/openapi_users.yaml`。

合法请求：

```http
POST /users?source=web
Content-Type: application/json

{"name":"tom","role":"admin","age":18}
```

非法请求：

```http
POST /users
Content-Type: application/json

{"name":"tom","role":"admin"}
```

这个非法请求会被拦截，因为 OpenAPI 中定义了 query 参数 `source` 为必填。
