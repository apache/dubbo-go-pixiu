# OpenAPI 请求校验过滤器

[English](openapi.md) | 中文

---

## 概述

Pixiu 可以在 `dgp.filter.http.apiconfig` 初始化阶段加载 OpenAPI 3.x 文件，并把每个 operation 的
`ValidationPlan` 合并到 `api_config` 里已经存在的同名路由上。

当请求命中 API 后，Pixiu 会在转发到上游之前先执行请求校验。

如果校验失败，Pixiu 会直接返回 `400 Bad Request`，并停止后续过滤链。

## 第一版支持范围

- 本地 OpenAPI 3.x 文件加载
- path 参数提取
- query 参数校验
- header 校验
- JSON request body 校验
- `required`
- `type`
- `enum`
- `minimum`
- `maximum`
- `minLength`
- `maxLength`

## 第一版暂不支持

- response validation
- 远程 `$ref`
- `oneOf` / `allOf` / `anyOf`
- admin 或配置中心分发 OpenAPI 文件
- `dynamic + openapi_path` 组合配置

## 配置示例

```yaml
- name: dgp.filter.http.apiconfig
  config:
    path: configs/api_config.yaml
    openapi_path: configs/openapi_users.yaml
    enable_openapi_validation: true
```

OpenAPI 校验不会单独创建路由，目标路由必须已经存在于 `api_config` 中。

## 说明

- 参数级校验现在覆盖 `path`、`query`、`header` 上常见的标量类型约束。
- 第一版不支持在同一个 filter 配置里同时使用 `dynamic` 和 `openapi_path`。

## 运行流程

1. Pixiu 在 `apiconfig.Apply()` 阶段先加载 `api_config`。
2. Pixiu 再加载 OpenAPI 文件并编译校验计划。
3. 每个编译结果的 `ValidationPlan` 会合并到匹配的 `router.API.Metadata` 中。
4. 请求进入 `apiconfig.Decode()`。
5. Pixiu 先按 path 和 method 匹配 API。
6. 从命中的 API metadata 中取出 `ValidationPlan`。
7. Pixiu 执行请求校验。
8. 校验通过，请求继续流向后续代理过滤器。
9. 校验失败，Pixiu 直接返回 `400 Bad Request`。

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
