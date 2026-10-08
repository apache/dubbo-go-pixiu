# Backend API Documentation

**English** | [中文](API_CN.md)

The API Router configuration is managed through the `AdminRouteBinding` lifecycle. A route is saved as a draft first, then validated, diffed, and published independently. Publishing or deleting one route does not publish or delete other routes.

More detailed API descriptions are available in the [Swagger documentation](./doc/swagger.json).

## Upgrade compatibility

The Admin API Router model is now `AdminRouteBinding`. This is a breaking change for the legacy Resource/Method admin model: the new API does not import or convert existing legacy records, and they are not listed by the new route endpoints. Recreate routes as `AdminRouteBinding` objects when they need to be managed through the new Admin API.

## Response Codes

* `10001`: Success
* `10002`: Data not found
* `10003`: Concurrent operation; refresh and try again

## Basic Information

```http
GET /config/api/base
POST /config/api/base/
PUT /config/api/base/
```

The write endpoints accept YAML in the form field `content`.

## API Routes

### List and inspect routes

```http
GET /config/api/route/list?scope=draft
GET /config/api/route/list?scope=published
GET /config/api/route/detail?name=<route-name>&scope=draft
GET /config/api/route/status?name=<route-name>
GET /config/api/route/diff?name=<route-name>
```

The draft list includes the backend-calculated publish status for every route.

### Save a draft

```http
POST /config/api/route
PUT /config/api/route?name=<original-route-name>
```

Send an `AdminRouteBinding` JSON object. The `name` query parameter is the immutable route identity and must match `metadata.name`; updates may include `expectedRevision` for optimistic concurrency control:

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

### Validate and preview

```http
POST /config/api/route/validate
POST /config/api/route/preview
```

These endpoints do not write configuration. Validation applies schema defaults; preview returns the generated legacy Pixiu YAML.

Publishing always validates the route; there is no per-route validation toggle.

### Publish or delete one route

```http
PUT    /config/api/route/publish?name=<route-name>&expectedRevision=<revision>
DELETE /config/api/route?name=<route-name>&expectedRevision=<revision>
```

Both operations are independent per route and use one etcd transaction for the Admin binding and generated runtime configuration. There is no full-configuration publish endpoint.

Before publishing, the Admin API checks legacy runtime Resource/Method records, including inline methods and separate method entries. A route with the same HTTP method and path (case-insensitively) is rejected until the old runtime route is removed.

## OPA Policy

OPA policy APIs proxy requests to the OPA server. If `server_url` or `policy_id` is not provided, defaults are used (`http://opa:8181` and `pixiu-authz`).

### Get OPA Policy

```http
GET /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Optional query parameters:

* `policy_id`
* `server_url`
* `bearer_token`

The response contains the policy text in `data`. If the policy does not exist, `data` is an empty string.

### Create or Update OPA Policy

```http
PUT /config/api/opa/policy HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data
```

Form fields:

* `policy_id`
* `content`
* `server_url` (optional)
* `bearer_token` (optional)

### Delete OPA Policy

```http
DELETE /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Optional query parameters: `policy_id`, `server_url`, and `bearer_token`.

## xDS Diagnostics

### Get xDS Publication Status

This authenticated endpoint returns xDS listener availability, the last-good snapshot version, resource counts, publication timestamps, the latest listener or rejected-candidate error, and the support status of each xDS resource path.

```http
GET /config/api/xds/status HTTP/1.1
Host: 127.0.0.1:8080
token: <admin-jwt>
```

`ready` requires both a bound xDS listener and a published snapshot. `degraded` indicates that listener startup or a newer snapshot candidate failed.

**Response**:

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
