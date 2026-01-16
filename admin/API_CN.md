# 后端 API 接口文档

[English](API.md) | **中文**

本文档描述了 Pixiu Admin 管理平台的 RESTful API。所有 API 使用 JSON 格式进行请求和响应，除登录/注册外均需要 JWT 认证。

## 基础 URL

```
http://127.0.0.1:8081/api
```

## 认证

所有受保护的接口需要在请求头中携带 JWT Token：

```
token: <your-jwt-token>
```

## 响应格式

所有响应遵循以下格式：

```json
{
  "code": 0,
  "message": "可选的消息",
  "data": {}
}
```

- `code`：`0` 表示成功，`-1` 表示错误，`401` 表示未授权，`403` 表示禁止访问
- `message`：错误消息或成功消息（可选）
- `data`：响应数据（可选）

---

## 一、认证

### 1.1 登录

**请求：**

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 1.2 注册

**请求：**

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "nickname": "新用户",
  "email": "user@example.com"
}
```

**响应：**

```json
{
  "code": 0,
  "message": "注册成功"
}
```

---

## 二、用户信息

### 2.1 获取当前用户信息

**请求：**

```http
GET /api/user/info
token: <jwt-token>
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "email": "admin@example.com",
    "role_id": 1,
    "status": 1
  }
}
```

### 2.2 修改密码

**请求：**

```http
POST /api/user/password
token: <jwt-token>
Content-Type: application/json

{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

### 2.3 退出登录

**请求：**

```http
POST /api/user/logout
token: <jwt-token>
```

---

## 三、集群管理

### 3.1 获取集群列表

**请求：**

```http
GET /api/clusters
token: <jwt-token>
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "name": "backend-cluster",
      "type_str": "EDS",
      "endpoints": [
        {
          "address": {
            "socket_address": {
              "address": "127.0.0.1",
              "port": 8080
            }
          }
        }
      ]
    }
  ]
}
```

### 3.2 获取集群详情

**请求：**

```http
GET /api/clusters/:name
token: <jwt-token>
```

### 3.3 创建集群

**请求：**

```http
POST /api/clusters
token: <jwt-token>
Content-Type: application/json

{
  "name": "my-cluster",
  "type_str": "EDS",
  "lb_str": "RoundRobin",
  "endpoints": [
    {
      "address": {
        "socket_address": {
          "address": "127.0.0.1",
          "port": 8080
        }
      }
    }
  ],
  "health_checks": [
    {
      "timeout": "5s",
      "interval": "10s",
      "healthy_threshold": 2,
      "unhealthy_threshold": 3
    }
  ]
}
```

### 3.4 更新集群

**请求：**

```http
PUT /api/clusters/:name
token: <jwt-token>
Content-Type: application/json

{
  "name": "my-cluster",
  "type_str": "EDS",
  "lb_str": "LeastRequest",
  "endpoints": [...]
}
```

### 3.5 删除集群

**请求：**

```http
DELETE /api/clusters/:name
token: <jwt-token>
```

---

## 四、监听器管理

### 4.1 获取监听器列表

**请求：**

```http
GET /api/listeners
token: <jwt-token>
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "name": "http-listener",
      "protocol_str": "HTTP",
      "address": {
        "socket_address": {
          "address": "0.0.0.0",
          "port": 8080
        }
      },
      "filter_chains": [...]
    }
  ]
}
```

### 4.2 获取监听器详情

**请求：**

```http
GET /api/listeners/:name
token: <jwt-token>
```

### 4.3 创建监听器

**请求：**

```http
POST /api/listeners
token: <jwt-token>
Content-Type: application/json

{
  "name": "my-listener",
  "protocol_str": "HTTP",
  "address": {
    "socket_address": {
      "address": "0.0.0.0",
      "port": 8080
    }
  },
  "filter_chains": [
    {
      "filters": [
        {
          "name": "dgp.filter.httpconnectionmanager",
          "config": {
            "route_config": {
              "routes": [...]
            }
          }
        }
      ]
    }
  ]
}
```

### 4.4 更新监听器

**请求：**

```http
PUT /api/listeners/:name
token: <jwt-token>
Content-Type: application/json
```

### 4.5 删除监听器

**请求：**

```http
DELETE /api/listeners/:name
token: <jwt-token>
```

---

## 五、资源管理（API 映射）

### 5.1 获取资源列表

**请求：**

```http
GET /api/resources
token: <jwt-token>
```

### 5.2 获取资源详情

**请求：**

```http
GET /api/resources/:id
token: <jwt-token>
```

### 5.3 创建资源

**请求：**

```http
POST /api/resources
token: <jwt-token>
Content-Type: application/json

{
  "path": "/api/v1/users",
  "type": "restful",
  "description": "用户 API",
  "timeout": "30s",
  "plugins": {
    "pre": {
      "pluginNames": ["rate-limit"]
    }
  }
}
```

### 5.4 更新资源

**请求：**

```http
PUT /api/resources/:id
token: <jwt-token>
Content-Type: application/json
```

### 5.5 删除资源

**请求：**

```http
DELETE /api/resources/:id
token: <jwt-token>
```

---

## 六、方法管理

### 6.1 获取方法列表

**请求：**

```http
GET /api/methods?resource_id=1
token: <jwt-token>
```

### 6.2 获取方法详情

**请求：**

```http
GET /api/methods/:id
token: <jwt-token>
```

### 6.3 创建方法

**请求：**

```http
POST /api/methods
token: <jwt-token>
Content-Type: application/json

{
  "resource_id": 1,
  "http_verb": "GET",
  "on_air": true,
  "timeout": "10s",
  "inbound_request": {
    "request_type": "http"
  },
  "integration_request": {
    "request_type": "http",
    "host": "127.0.0.1:8889",
    "path": "/backend/users"
  }
}
```

### 6.4 更新方法

**请求：**

```http
PUT /api/methods/:id
token: <jwt-token>
Content-Type: application/json
```

### 6.5 删除方法

**请求：**

```http
DELETE /api/methods/:id
token: <jwt-token>
```

---

## 七、插件组管理

### 7.1 获取插件组列表

**请求：**

```http
GET /api/plugins
token: <jwt-token>
```

### 7.2 获取插件组详情

**请求：**

```http
GET /api/plugins/:name
token: <jwt-token>
```

### 7.3 创建插件组

**请求：**

```http
POST /api/plugins
token: <jwt-token>
Content-Type: application/json

{
  "group_name": "my-plugin-group",
  "plugins": [
    {
      "name": "rate-limit",
      "version": "1.0.0",
      "priority": 100,
      "config": {}
    }
  ]
}
```

### 7.4 更新插件组

**请求：**

```http
PUT /api/plugins/:name
token: <jwt-token>
Content-Type: application/json
```

### 7.5 删除插件组

**请求：**

```http
DELETE /api/plugins/:name
token: <jwt-token>
```

---

## 八、实例管理

### 8.1 获取实例列表

**请求：**

```http
GET /api/instances
token: <jwt-token>
```

### 8.2 获取实例统计

**请求：**

```http
GET /api/instances/stats
token: <jwt-token>
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "total": 10,
    "healthy": 8,
    "unhealthy": 2
  }
}
```

---

## 九、用户管理（仅管理员）

### 9.1 获取用户列表

**请求：**

```http
GET /api/users?page=1&page_size=10
token: <jwt-token>
```

### 9.2 获取用户详情

**请求：**

```http
GET /api/users/:id
token: <jwt-token>
```

### 9.3 创建用户

**请求：**

```http
POST /api/users
token: <jwt-token>
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "nickname": "新用户",
  "email": "user@example.com",
  "role_id": 2
}
```

### 9.4 更新用户

**请求：**

```http
PUT /api/users/:id
token: <jwt-token>
Content-Type: application/json

{
  "nickname": "更新的名称",
  "email": "updated@example.com",
  "status": 1
}
```

### 9.5 删除用户

**请求：**

```http
DELETE /api/users/:id
token: <jwt-token>
```

### 9.6 重置用户密码

**请求：**

```http
POST /api/users/:id/reset-password
token: <jwt-token>
Content-Type: application/json

{
  "password": "newpassword123"
}
```

### 9.7 分配用户角色

**请求：**

```http
POST /api/users/:id/assign-role
token: <jwt-token>
Content-Type: application/json

{
  "role_id": 2
}
```

---

## 十、角色管理（仅管理员）

### 10.1 获取角色列表

**请求：**

```http
GET /api/roles
token: <jwt-token>
```

### 10.2 获取角色详情

**请求：**

```http
GET /api/roles/:id
token: <jwt-token>
```

### 10.3 创建角色

**请求：**

```http
POST /api/roles
token: <jwt-token>
Content-Type: application/json

{
  "name": "operator",
  "description": "操作员角色"
}
```

### 10.4 更新角色

**请求：**

```http
PUT /api/roles/:id
token: <jwt-token>
Content-Type: application/json

{
  "name": "operator",
  "description": "更新的描述"
}
```

### 10.5 删除角色

**请求：**

```http
DELETE /api/roles/:id
token: <jwt-token>
```

### 10.6 获取角色权限

**请求：**

```http
GET /api/roles/:id/permissions
token: <jwt-token>
```

### 10.7 更新角色权限

**请求：**

```http
PUT /api/roles/:id/permissions
token: <jwt-token>
Content-Type: application/json

{
  "permission_ids": [1, 2, 3, 4]
}
```

---

## 十一、权限管理（仅管理员）

### 11.1 获取所有权限

**请求：**

```http
GET /api/permissions
token: <jwt-token>
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "resource": "clusters",
      "action": "read",
      "description": "查看集群"
    },
    {
      "id": 2,
      "resource": "clusters",
      "action": "create",
      "description": "创建集群"
    }
  ]
}
```

---

## 错误码

| 错误码 | 描述 |
|--------|------|
| 0 | 成功 |
| -1 | 通用错误 |
| 401 | 未授权（Token 无效或缺失） |
| 403 | 禁止访问（权限不足） |

## 权限资源

| 资源 | 操作 |
|------|------|
| clusters | read, create, update, delete |
| listeners | read, create, update, delete |
| resources | read, create, update, delete |
| methods | read, create, update, delete |
| plugins | read, create, update, delete |
| users | read, create, update, delete |
