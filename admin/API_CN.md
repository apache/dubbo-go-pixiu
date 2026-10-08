# 后端API接口文档

[English](API.md) | **中文**

本接口文档详细描述 Pixiu 管理平台的后端 API 操作，包括 API 路由绑定（AdminRouteBinding）、插件组（PluginGroup）和 OPA 策略等接口。Pixiu 平台提供相关 API 来管理 API 网关路由映射、插件配置和请求处理。文档中的示例涵盖常见的请求与响应格式，并介绍如何使用 Postman 测试接口。

无论是创建新的 API 路由、修改现有配置，还是管理插件组，本文档都提供清晰的步骤和必要的 API 细节，方便开发者快速上手并进行集成。

更多的 API 具体介绍请参考 [Swagger 文档](./doc/swagger.json)

## 返回值说明

* **code**：

    * `10001`: 成功
    * `10002`: 未找到对应数据
    * `10003`: 并发操作，请刷新页面重试

* **data**：一般为 YAML 格式的数据

## 一、基础信息

### 1.1 获取基础信息

**请求**：

```http
GET /config/api/base HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

**返回值**：

```json
{
  "code": "10001",
  "data": "name: pixiu\ndescription: pixiu111 sample\npluginFilePath: \"\"\n"
}
```

### 1.2 创建或修改基础信息

**请求**：

```http
POST /config/api/base HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
cache-control: no-cache
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
name: pixiu
description: pixiu111 sample
```

## 二、API 路由

API Router 配置使用 AdminRouteBinding 对象管理。旧 Resource/Method Admin CRUD 接口已由新模型替代；已有运行时 Resource/Method 配置不会自动导入新草稿模型，需要通过新接口管理的路由应重新创建为 AdminRouteBinding。新路由发布后仍会生成 Pixiu 使用的旧格式运行时配置。

### 2.1 获取 API 路由 Schema

**请求**：

```http
GET /config/api/route/schema HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

返回数据包含已注册的 Admin 对象 Schema，其中包括编辑器使用的 AdminRouteBinding Schema。

### 2.2 获取 API 路由列表

**请求**：

```http
GET /config/api/route/list?scope=draft HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

scope=draft（默认值）查询草稿，scope=published 查询已发布路由。接口同时兼容旧的 unpublished 参数：1 表示草稿，0 表示已发布路由。草稿列表包含由后端计算的发布状态。

### 2.3 获取 API 路由详情

**请求**：

```http
GET /config/api/route/detail?name=<路由名>&scope=draft HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

name 是不可变的路由标识。查询已发布版本时将 scope 设为 published。

### 2.4 创建 API 路由草稿

**请求**：

```http
POST /config/api/route HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json
cache-control: no-cache
```

请求体为 AdminRouteBinding JSON 对象。也可使用 {"object": <AdminRouteBinding>, "expectedRevision": <revision>} 包装对象，或通过表单字段 content 提交 YAML。

### 2.5 修改 API 路由草稿

**请求**：

```http
PUT /config/api/route?name=<路由名> HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json
cache-control: no-cache
```

请求体为更新后的 AdminRouteBinding 对象。name 查询参数必须与 metadata.name 一致且不可修改。expectedRevision 可放在包装后的 JSON 请求体、查询参数或 If-Match 请求头中，用于拒绝过期更新。

### 2.6 校验或预览 API 路由

**请求**：

```http
POST /config/api/route/validate HTTP/1.1
POST /config/api/route/preview HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json
cache-control: no-cache
```

两个接口都接收 AdminRouteBinding 对象，不会写入 etcd。校验接口返回包含 Schema 默认值的规范化对象；预览接口返回规范化对象以及生成的 Pixiu 旧格式 YAML。

### 2.7 获取 API 路由状态或 Diff

**请求**：

```http
GET /config/api/route/status?name=<路由名> HTTP/1.1
GET /config/api/route/diff?name=<路由名> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

状态接口返回单条路由的草稿和已发布状态；Diff 接口返回草稿与已发布对象之间的字段差异。

### 2.8 发布单条 API 路由

**请求**：

```http
PUT /config/api/route/publish?name=<路由名>&expectedRevision=<revision> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

发布前会校验路由，并通过一个 etcd 事务原子更新该路由的已发布 AdminRouteBinding 和生成的运行时配置。expectedRevision 为可选项，可用于拒绝过期草稿。系统不提供 API 路由全量发布接口。

发布前还会检查旧运行时 Resource/Method 配置冲突。若 HTTP 方法和路径冲突，应先从旧运行时配置中移除对应路由。

### 2.9 删除单条 API 路由

**请求**：

```http
DELETE /config/api/route?name=<路由名>&expectedRevision=<revision> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

删除会通过一个 etcd 事务立即移除路由草稿和对应的已发布运行时配置。客户端发送请求前应先向用户确认删除。
expectedRevision 为可选项；如果提供，可避免删除读取后已被其他请求修改的草稿。

## 三、PluginGroup 和 Plugin 相关

### 3.1 查看 PluginGroup 列表

**请求**：

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.2 查看 PluginGroup 详情

**请求**：

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.3 创建 PluginGroup

**请求**：

```http
POST /config/api/plugin_group/ HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
groupName: "group1"
plugins:
  - name: "rate limit"
    version: "0.0.1"
    priority: 1000
    externalLookupName: "ExternalPluginRateLimit"
  - name: "access"
    version: "0.0.1"
    priority: 1000
    externalLookupName: "ExternalPluginAccess"
```

### 3.4 修改 PluginGroup

**请求**：

```http
PUT /config/api/plugin_group/ HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
groupName: "group1"
plugins:
  - name: "rate limit"
    version: "0.0.2"
    priority: 1000
    externalLookupName: "ExternalPluginRateLimit"
  - name: "access"
    version: "0.0.1"
    priority: 1000
    externalLookupName: "ExternalPluginAccess"
```

### 3.5 删除 PluginGroup

**请求**：

```http
DELETE /config/api/plugin_group/?name=group1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.6 发布 PluginGroup

**请求**：

```http
PUT /config/api/plugin_group/publish HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

该旧接口将暂存空间中的 PluginGroup 配置发布到已发布空间，与 API 路由发布相互独立。

## 四、OPA 策略

OPA 策略接口会代理请求到 OPA 服务端。未提供 `server_url` 或 `policy_id` 时，会使用默认值（`http://opa:8181` 和 `pixiu-authz`）。

### 4.1 获取 OPA 策略

**请求**：

```http
GET /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

**Query 参数**：

* `policy_id`: OPA policy id（可选）
* `server_url`: OPA 服务地址（可选）
* `bearer_token`: OPA Bearer Token（可选）

**返回**：

```json
{
  "code": "10001",
  "data": "package pixiu.authz\n\ndefault allow := false\n"
}
```

若策略不存在，`data` 返回空字符串。

### 4.2 新增或更新 OPA 策略

**请求**：

```http
PUT /config/api/opa/policy HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
cache-control: no-cache
```

**表单数据**：

```text
Content-Disposition: form-data; name="policy_id"
pixiu-authz

Content-Disposition: form-data; name="content"
package pixiu.authz

default allow := false
```

可选表单字段：

* `server_url`
* `bearer_token`

### 4.3 删除 OPA 策略

**请求**：

```http
DELETE /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

**Query 参数**：

* `policy_id`: OPA policy id（可选）
* `server_url`: OPA 服务地址（可选）
* `bearer_token`: OPA Bearer Token（可选）

## 五、xDS 诊断

### 5.1 获取 xDS 发布状态

**请求**：

```http
GET /config/api/xds/status HTTP/1.1
Host: 127.0.0.1:8080
token: <admin-jwt>
```

该认证接口返回 xDS 监听状态、最近一次成功发布的快照版本、资源数量、发布时间、最近一次监听或候选配置错误，以及各条 xDS 资源链路的支持状态。

ready 只有在 xDS 端口已监听且至少一个快照成功发布时为 true；degraded 表示监听启动或新候选配置失败。

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
