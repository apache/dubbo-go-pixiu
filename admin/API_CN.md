# 后端 API 接口文档

[English](API.md) | **中文**

API Router 配置统一使用 `AdminRouteBinding` 生命周期管理：先保存为草稿，再独立校验、查看 Diff 和发布。发布或删除一条路由不会影响其他路由。

更多接口说明请参考 [Swagger 文档](./doc/swagger.json)。

## 升级兼容性

Admin API Router 现在统一使用 `AdminRouteBinding` 模型。这对旧的 Resource/Method 管理模型属于 breaking change：新接口不会自动导入或转换已有的旧配置，旧配置也不会出现在新的路由接口列表中。升级后如果需要通过新的 Admin API 管理这些路由，需要重新创建为 `AdminRouteBinding`。

## 返回值说明

* `10001`：成功
* `10002`：未找到对应数据
* `10003`：并发操作，请刷新页面重试

## 基础信息

```http
GET /config/api/base
POST /config/api/base/
PUT /config/api/base/
```

写入接口通过表单字段 `content` 接收 YAML。

## API 路由

### 列表和状态

```http
GET /config/api/route/list?scope=draft
GET /config/api/route/list?scope=published
GET /config/api/route/detail?name=<路由名>&scope=draft
GET /config/api/route/status?name=<路由名>
GET /config/api/route/diff?name=<路由名>
```

草稿列表会同时返回后端计算的每条路由发布状态，前端无需逐条请求状态。

### 保存草稿

```http
POST /config/api/route
PUT /config/api/route?name=<原始路由名>
```

请求体为 `AdminRouteBinding` JSON 对象。更新时 `name` 查询参数是不可变的路由身份，必须与 `metadata.name` 一致；同时可以携带 `expectedRevision` 做乐观并发控制：

```json
{
  "object": {
    "kind": "AdminRouteBinding",
    "metadata": {"name": "user-get"},
    "spec": {}
  },
  "expectedRevision": 12
}
```

### 校验和预览

```http
POST /config/api/route/validate
POST /config/api/route/preview
```

这两个接口不会写入配置。校验接口会补齐 schema 默认值，预览接口返回 Pixiu 当前 watcher 使用的 legacy YAML。

发布始终会校验路由，不再保存或提供单路由校验开关。

### 单路由发布和删除

```http
PUT    /config/api/route/publish?name=<路由名>&expectedRevision=<revision>
DELETE /config/api/route?name=<路由名>&expectedRevision=<revision>
```

两个操作都按单路由执行，并通过一个 etcd 事务同步 Admin binding 和生成的运行时配置。系统不再提供全量配置发布接口。

发布前会检查旧运行时 Resource/Method 配置，包括 Resource 内联方法和独立的 Method 键。如果存在 HTTP 方法和路径（忽略大小写）相同的路由，发布会被拒绝，需要先清理旧路由。

## OPA 策略

OPA 策略接口会代理请求到 OPA 服务端。未提供 `server_url` 或 `policy_id` 时，会使用默认值（`http://opa:8181` 和 `pixiu-authz`）。

### 获取 OPA 策略

```http
GET /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

可选查询参数：

* `policy_id`
* `server_url`
* `bearer_token`

返回的 `data` 包含策略文本。策略不存在时，`data` 为空字符串。

### 新增或更新 OPA 策略

```http
PUT /config/api/opa/policy HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data
```

表单字段：

* `policy_id`
* `content`
* `server_url`（可选）
* `bearer_token`（可选）

### 删除 OPA 策略

```http
DELETE /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

可选查询参数：`policy_id`、`server_url` 和 `bearer_token`。

## xDS 诊断

### 获取 xDS 发布状态

该认证接口返回 xDS 监听状态、最近一次成功发布的快照版本、资源数量、发布时间、最近一次监听或候选配置错误，以及各条 xDS 资源链路的支持状态。

```http
GET /config/api/xds/status HTTP/1.1
Host: 127.0.0.1:8080
token: <admin-jwt>
```

`ready` 只有在 xDS 端口已监听且至少一个快照成功发布时为 true；`degraded` 表示监听启动或新候选配置失败。

**返回**：

```json
{
  "code": "10001",
  "data": {
    "node_id": "test-id",
    "snapshot_version": "42",
    "listener_count": 1,
    "cluster_count": 2,
    "listen_port": 18000,
    "listening": true,
    "ready": true,
    "degraded": false,
    "resource_support": {
      "extension_config_listener": "supported",
      "extension_config_cluster": "supported",
      "standard_cds": "experimental",
      "standard_eds": "experimental",
      "standard_lds": "unsupported"
    }
  }
}
```
