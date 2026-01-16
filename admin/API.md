# Backend API Documentation

**English** | [中文](API_CN.md)

This document describes the RESTful API for the Pixiu Admin management platform. All APIs use JSON for request and response bodies, and require JWT authentication (except for login/register).

## Base URL

```
http://127.0.0.1:8081/api
```

## Authentication

All protected endpoints require a JWT token in the `token` header:

```
token: <your-jwt-token>
```

## Response Format

All responses follow this format:

```json
{
  "code": 0,
  "message": "optional message",
  "data": {}
}
```

- `code`: `0` for success, `-1` for error, `401` for unauthorized, `403` for forbidden
- `message`: Error message or success message (optional)
- `data`: Response data (optional)

---

## I. Authentication

### 1.1 Login

**Request:**

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**

```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 1.2 Register

**Request:**

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "nickname": "New User",
  "email": "user@example.com"
}
```

**Response:**

```json
{
  "code": 0,
  "message": "Registration successful"
}
```

---

## II. User Profile

### 2.1 Get Current User Info

**Request:**

```http
GET /api/user/info
token: <jwt-token>
```

**Response:**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "Administrator",
    "email": "admin@example.com",
    "role_id": 1,
    "status": 1
  }
}
```

### 2.2 Change Password

**Request:**

```http
POST /api/user/password
token: <jwt-token>
Content-Type: application/json

{
  "old_password": "oldpass",
  "new_password": "newpass"
}
```

### 2.3 Logout

**Request:**

```http
POST /api/user/logout
token: <jwt-token>
```

---

## III. Clusters

### 3.1 List Clusters

**Request:**

```http
GET /api/clusters
token: <jwt-token>
```

**Response:**

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

### 3.2 Get Cluster

**Request:**

```http
GET /api/clusters/:name
token: <jwt-token>
```

### 3.3 Create Cluster

**Request:**

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

### 3.4 Update Cluster

**Request:**

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

### 3.5 Delete Cluster

**Request:**

```http
DELETE /api/clusters/:name
token: <jwt-token>
```

---

## IV. Listeners

### 4.1 List Listeners

**Request:**

```http
GET /api/listeners
token: <jwt-token>
```

**Response:**

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

### 4.2 Get Listener

**Request:**

```http
GET /api/listeners/:name
token: <jwt-token>
```

### 4.3 Create Listener

**Request:**

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

### 4.4 Update Listener

**Request:**

```http
PUT /api/listeners/:name
token: <jwt-token>
Content-Type: application/json
```

### 4.5 Delete Listener

**Request:**

```http
DELETE /api/listeners/:name
token: <jwt-token>
```

---

## V. Resources (API Mappings)

### 5.1 List Resources

**Request:**

```http
GET /api/resources
token: <jwt-token>
```

### 5.2 Get Resource

**Request:**

```http
GET /api/resources/:id
token: <jwt-token>
```

### 5.3 Create Resource

**Request:**

```http
POST /api/resources
token: <jwt-token>
Content-Type: application/json

{
  "path": "/api/v1/users",
  "type": "restful",
  "description": "User API",
  "timeout": "30s",
  "plugins": {
    "pre": {
      "pluginNames": ["rate-limit"]
    }
  }
}
```

### 5.4 Update Resource

**Request:**

```http
PUT /api/resources/:id
token: <jwt-token>
Content-Type: application/json
```

### 5.5 Delete Resource

**Request:**

```http
DELETE /api/resources/:id
token: <jwt-token>
```

---

## VI. Methods

### 6.1 List Methods

**Request:**

```http
GET /api/methods?resource_id=1
token: <jwt-token>
```

### 6.2 Get Method

**Request:**

```http
GET /api/methods/:id
token: <jwt-token>
```

### 6.3 Create Method

**Request:**

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

### 6.4 Update Method

**Request:**

```http
PUT /api/methods/:id
token: <jwt-token>
Content-Type: application/json
```

### 6.5 Delete Method

**Request:**

```http
DELETE /api/methods/:id
token: <jwt-token>
```

---

## VII. Plugin Groups

### 7.1 List Plugin Groups

**Request:**

```http
GET /api/plugins
token: <jwt-token>
```

### 7.2 Get Plugin Group

**Request:**

```http
GET /api/plugins/:name
token: <jwt-token>
```

### 7.3 Create Plugin Group

**Request:**

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

### 7.4 Update Plugin Group

**Request:**

```http
PUT /api/plugins/:name
token: <jwt-token>
Content-Type: application/json
```

### 7.5 Delete Plugin Group

**Request:**

```http
DELETE /api/plugins/:name
token: <jwt-token>
```

---

## VIII. Instances

### 8.1 List Instances

**Request:**

```http
GET /api/instances
token: <jwt-token>
```

### 8.2 Get Instance Stats

**Request:**

```http
GET /api/instances/stats
token: <jwt-token>
```

**Response:**

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

## IX. User Management (Admin Only)

### 9.1 List Users

**Request:**

```http
GET /api/users?page=1&page_size=10
token: <jwt-token>
```

### 9.2 Get User

**Request:**

```http
GET /api/users/:id
token: <jwt-token>
```

### 9.3 Create User

**Request:**

```http
POST /api/users
token: <jwt-token>
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "nickname": "New User",
  "email": "user@example.com",
  "role_id": 2
}
```

### 9.4 Update User

**Request:**

```http
PUT /api/users/:id
token: <jwt-token>
Content-Type: application/json

{
  "nickname": "Updated Name",
  "email": "updated@example.com",
  "status": 1
}
```

### 9.5 Delete User

**Request:**

```http
DELETE /api/users/:id
token: <jwt-token>
```

### 9.6 Reset User Password

**Request:**

```http
POST /api/users/:id/reset-password
token: <jwt-token>
Content-Type: application/json

{
  "password": "newpassword123"
}
```

### 9.7 Assign Role to User

**Request:**

```http
POST /api/users/:id/assign-role
token: <jwt-token>
Content-Type: application/json

{
  "role_id": 2
}
```

---

## X. Role Management (Admin Only)

### 10.1 List Roles

**Request:**

```http
GET /api/roles
token: <jwt-token>
```

### 10.2 Get Role

**Request:**

```http
GET /api/roles/:id
token: <jwt-token>
```

### 10.3 Create Role

**Request:**

```http
POST /api/roles
token: <jwt-token>
Content-Type: application/json

{
  "name": "operator",
  "description": "Operator role"
}
```

### 10.4 Update Role

**Request:**

```http
PUT /api/roles/:id
token: <jwt-token>
Content-Type: application/json

{
  "name": "operator",
  "description": "Updated description"
}
```

### 10.5 Delete Role

**Request:**

```http
DELETE /api/roles/:id
token: <jwt-token>
```

### 10.6 Get Role Permissions

**Request:**

```http
GET /api/roles/:id/permissions
token: <jwt-token>
```

### 10.7 Update Role Permissions

**Request:**

```http
PUT /api/roles/:id/permissions
token: <jwt-token>
Content-Type: application/json

{
  "permission_ids": [1, 2, 3, 4]
}
```

---

## XI. Permissions (Admin Only)

### 11.1 List All Permissions

**Request:**

```http
GET /api/permissions
token: <jwt-token>
```

**Response:**

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "resource": "clusters",
      "action": "read",
      "description": "View clusters"
    },
    {
      "id": 2,
      "resource": "clusters",
      "action": "create",
      "description": "Create clusters"
    }
  ]
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| -1 | General error |
| 401 | Unauthorized (invalid or missing token) |
| 403 | Forbidden (insufficient permissions) |

## Permission Resources

| Resource | Actions |
|----------|---------|
| clusters | read, create, update, delete |
| listeners | read, create, update, delete |
| resources | read, create, update, delete |
| methods | read, create, update, delete |
| plugins | read, create, update, delete |
| users | read, create, update, delete |
