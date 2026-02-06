# 后端API接口文档

[English](API.md) | **中文**

本接口文档详细描述了 Pixiu 管理平台的后端 API 操作，包括获取、创建、修改、删除资源（Resource）、方法（Method）及插件组（PluginGroup）的接口。Pixiu
平台提供了一整套 API 来帮助用户管理 API 网关的资源映射、插件配置以及请求处理。文档中的示例涵盖了常见的请求与响应格式，并介绍了如何使用
Postman 进行接口测试。

无论是创建新资源、修改现有配置，还是管理插件组，本文档都提供了清晰的步骤和必要的 API 细节，方便开发者快速上手并进行集成。

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

## 二、Resource

### 2.1 获取 Resource 列表

**请求**：

```http
GET /config/api/resource/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 2.2 获取 Resource 详情

**请求**：

```http
GET /config/api/resource/detail?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 2.3 创建 Resource

**请求**：

```http
POST /config/api/resource/ HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
cache-control: no-cache
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
path: '/api/v1/test-dubbo/friend2'
type: restful
description: user
timeout: 100ms
plugins:
  pre:
    pluginNames:
      - rate limit
      - access
  post:
    groupNames:
      - group2
methods:
  - httpVerb: GET
    resourcePath: '/api/v1/test-dubbo/friend2'
    onAir: true
    timeout: 1000ms
    inboundRequest:
      requestType: http
      queryStrings:
        - name: name
          required: true
    integrationRequest:
      requestType: http
      host: 127.0.0.1:8889
      path: /UserProvider/GetUserByName
      mappingParams:
        - name: queryStrings.name
          mapTo: queryStrings.name
      group: "test"
      version: 1.0.0
```

### 2.4 修改 Resource

**请求**：

```http
PUT /config/api/resource? HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
id: 1
path: '/api/v1/test-dubbo/friend'
type: restful
description: update
timeout: 1000ms
plugins:
  pre:
    pluginNames:
      - rate limit
      - access
  post:
    groupNames:
      - group2
methods:
  - httpVerb: GET
    onAir: true
    timeout: 1000ms
    inboundRequest:
      requestType: http
      queryStrings:
        - name: name
          required: true
    integrationRequest:
      requestType: http
      host: 127.0.0.1:8889
      path: /UserProvider/GetUserByName
      mappingParams:
        - name: queryStrings.name
          mapTo: queryStrings.name
      group: "test"
      version: 1.0.0
```

### 2.5 删除 Resource

**请求**：

```http
DELETE /config/api/resource/?resourceId=2 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

## 三、Method 相关

### 3.1 查询某个 Resource 下的 Method 列表

**请求**：

```http
GET /config/api/resource/method/list?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.2 查询 Method 详情

**请求**：

```http
GET /config/api/resource/method/detail?resourceId=1&methodId=2 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.3 创建 Method

**请求**：

```http
POST /config/api/resource/method/?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
httpVerb: PUT
resourcePath: '/api/v1/test-dubbo/friend'
onAir: true
timeout: 1000ms
inboundRequest:
  requestType: http
  queryStrings:
    - name: name
      required: true
integrationRequest:
  requestType: http
  host: 127.0.0.1:8889
  path: /UserProvider/GetUserByName
  mappingParams:
    - name: queryStrings.name
      mapTo: queryStrings.name
  group: "test"
  version: 1.0.0
```

### 3.4 修改 Method

**请求**：

```http
PUT /config/api/resource/method/?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**表单数据**：

```text
Content-Disposition: form-data; name="content"
id: 2
httpVerb: PUT
resourcePath: '/api/v1/test-dubbo/friend'
onAir: true
timeout: 300ms
inboundRequest:
  requestType: http
  queryStrings:
    - name: name
      required: true
integrationRequest:
  requestType: http
  host: 127.0.0.1:8889
  path: /UserProvider/GetUserByName
  mappingParams:
    - name: queryStrings.name
      mapTo: queryStrings.name
  group: "test"
  version: 1.0.0
```

### 3.5 删除 Method

**请求**：

```http
DELETE /config/api/resource/method/?resourceId=1&methodId=2 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

## 四、PluginGroup 和 Plugin 相关

### 4.1 查看 PluginGroup 列表

**请求**：

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 4.2 查看 PluginGroup 详情

**请求**：

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 4.3 创建 PluginGroup

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

### 4.4 修改 PluginGroup

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

### 4.5 删除 PluginGroup

**请求**：

```http
DELETE /config/api/plugin_group/?name=group1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

## 五、OPA 策略

OPA 策略接口会代理请求到 OPA 服务端。未提供 `server_url` 或 `policy_id` 时，会使用默认值（`http://opa:8181` 和 `pixiu-authz`）。

### 5.1 获取 OPA 策略

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

### 5.2 新增或更新 OPA 策略

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

### 5.3 删除 OPA 策略

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
