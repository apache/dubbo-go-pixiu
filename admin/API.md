# Backend API Documentation

**English** | [中文](API_CN.md)

This API documentation describes the backend operations of the Pixiu management platform, including APIs for managing API route bindings (AdminRouteBinding), plugin groups (PluginGroup), and OPA policies. Pixiu provides APIs to help users manage API gateway route mappings, plugin configurations, and request handling. The examples in this document cover common request and response formats and show how to test the APIs using Postman.

Whether you are creating new API routes, modifying existing configurations, or managing plugin groups, this document provides clear steps and necessary API details, making it easier for developers to get started and integrate quickly.

More detailed API descriptions can be found in the [Swagger documentation](./doc/swagger.json).

## Response Codes

* **code**:

    * `10001`: Success
    * `10002`: Data not found
    * `10003`: Concurrent operation, please refresh the page and try again

* **data**: Typically, data will be in YAML format.

## I. Basic Information

### 1.1 Get Basic Information

**Request**:

```http
GET /config/api/base HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

**Response**:

```json
{
  "code": "10001",
  "data": "name: pixiu\ndescription: pixiu111 sample\npluginFilePath: \"\"\n"
}
```

### 1.2 Create or Modify Basic Information

**Request**:

```http
POST /config/api/base HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
cache-control: no-cache
```

**Form Data**:

```text
Content-Disposition: form-data; name="content"
name: pixiu
description: pixiu111 sample
```

## II. API Routes

New API routes are managed through AdminRouteBinding objects and compiled to the existing Resource/Method runtime format. The legacy Resource/Method Admin CRUD and publish endpoints remain available for older clients and continue to use the old etcd keys. Existing legacy records are not imported into AdminRouteBinding drafts; use the new API to manage new routes, and use the compatibility endpoints to manage old records.

### 2.1 Get API Route Schema

**Request**:

```http
GET /config/api/route/schema HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

The response data contains the registered Admin object schemas, including the AdminRouteBinding schema used by the editor.

The current route model keeps the backend reference, service selection, and
invocation policy separate. A registry-backed draft has this shape:

```yaml
kind: AdminRouteBinding
metadata:
  name: user-get
spec:
  entry:
    protocol: http
    path: /api/v1/users/:id
    method: GET
  backendRef:
    type: registry
  serviceSelection:
    application: UserProvider
    interface: com.example.user.UserService
    group: stable
    version: 1.0.0
  invocation:
    method: GetUser
    timeout: 750ms
    parameters:
      - from: uri.id
        to: 0
        type: java.lang.String
  enabled: true
```

`timeout` is a positive Go duration. URI mappings must correspond one-to-one
with `:name` path segments. `backendRef.type` is currently fixed to `registry`;
services are resolved from the registry set configured on the listener's
`dgp.filter.http.dubboproxy` filter. `serviceSelection.application` and
`interface` identify the service. The legacy `target.cluster` value is migrated
to `integrationRequest.clusterName`; it does not choose a registry.
See the [Dubbo registry configuration example](../docs/sample/dubbo/dubbo-registry.md)
for the runtime registry setup. Objects written by PR #1046 using `target` and
`params` are normalized on read/validation and migrate when saved.

### 2.2 Get API Route List

**Request**:

```http
GET /config/api/route/list?scope=draft HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Use scope=draft (the default) to list drafts or scope=published to list published routes. The legacy unpublished query parameter is also accepted: 1 selects drafts and 0 selects published routes. The draft list includes backend-calculated publication status.

### 2.3 Get API Route Details

**Request**:

```http
GET /config/api/route/detail?name=<route-name>&scope=draft HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

The name parameter is the immutable route identity. Set scope=published to read the published version.

### 2.4 Create API Route Draft

**Request**:

```http
POST /config/api/route HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json
cache-control: no-cache
```

Send an AdminRouteBinding object as JSON. The endpoint also accepts the object wrapped as {"object": <AdminRouteBinding>, "expectedRevision": <revision>}, or YAML in the content form field.

### 2.5 Modify API Route Draft

**Request**:

```http
PUT /config/api/route?name=<route-name> HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json
cache-control: no-cache
```

Send the updated AdminRouteBinding object. The name query parameter must match metadata.name and cannot be changed. expectedRevision may be supplied in the wrapped JSON body, as a query parameter, or through the If-Match header to reject stale updates.

### 2.6 Validate or Preview an API Route

**Request**:

```http
POST /config/api/route/validate HTTP/1.1
POST /config/api/route/preview HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/json
cache-control: no-cache
```

Both endpoints accept an AdminRouteBinding object and do not write to etcd. Validation returns the normalized object with schema defaults. Preview returns the normalized object and the generated legacy Pixiu YAML.

### 2.7 Get API Route Status or Diff

**Request**:

```http
GET /config/api/route/status?name=<route-name> HTTP/1.1
GET /config/api/route/diff?name=<route-name> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Status reports the draft and published state for one route. Diff returns the field-level differences between its draft and published objects.

### 2.8 Query Published History

**Request**:

```http
GET /config/api/route/history?name=<route-name> HTTP/1.1
Host: 127.0.0.1:8080
```

Each successful publish and delete stores an immutable snapshot in the same
etcd transaction as the runtime update. History entries expose a monotonically
increasing revision, publish time, and the binding object; a deletion entry has
no object.

`GET /config/api/route/history/list` returns the latest revision summary for
every route, including tombstones. The Admin UI uses this endpoint to expose
history-only routes and restore a previously published object after deletion.

### 2.9 Publish One API Route

**Request**:

```http
PUT /config/api/route/publish?name=<route-name>&expectedRevision=<revision> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Publishing validates the route and atomically updates its published AdminRouteBinding and generated runtime configuration in one etcd transaction. expectedRevision is optional and can be used to reject a stale draft. There is no full API-route publish endpoint.

Before publishing, Admin checks for conflicts with legacy runtime Resource/Method records. A conflicting HTTP method and path must be removed from the old runtime configuration before the new route can be published.

### 2.10 Roll Back to a Published Revision

**Request**:

```http
PUT /config/api/route/rollback?name=<route-name>&revision=<history-revision>&expectedRevision=<draft-revision> HTTP/1.1
Host: 127.0.0.1:8080
```

Rollback stages the selected immutable snapshot and publishes it through the
normal atomic Resource/Method publisher. If a concurrent change prevents
publish, the snapshot remains a draft and the active runtime remains unchanged.
Rolling back to a deletion snapshot removes the route.

### 2.11 Delete One API Route

**Request**:

```http
DELETE /config/api/route?name=<route-name>&expectedRevision=<revision> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Deletion immediately removes the route draft and its published runtime configuration in one etcd transaction. The client should ask for confirmation before sending this request.
expectedRevision is optional; when supplied, it prevents deleting a draft that changed after it was read.

### 2.12 Legacy Resource/Method API Compatibility

Older Admin clients can continue using the original endpoints:

- `/config/api/resource/list`, `/detail`, and `/resource` (POST/PUT/DELETE)
- `/config/api/resource/method/list`, `/detail`, and `/resource/method` (POST/PUT/DELETE)
- `/config/api/resource/publish` and `/config/api/resource/method/publish`

These endpoints keep the legacy etcd key layout and YAML shapes. New
AdminRouteBinding publication checks for conflicts with those runtime routes;
remove or migrate a conflicting old route before publishing its replacement.

## III. PluginGroup and Plugin Related

### 3.1 Get PluginGroup List

**Request**:

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.2 Get PluginGroup Details

**Request**:

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.3 Create PluginGroup

**Request**:

```http
POST /config/api/plugin_group/ HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**Form Data**:

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

### 3.4 Modify PluginGroup

**Request**:

```http
PUT /config/api/plugin_group/ HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**Form Data**:

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

### 3.5 Delete PluginGroup

**Request**:

```http
DELETE /config/api/plugin_group/?name=group1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

### 3.6 Publish PluginGroup

**Request**:

```http
PUT /config/api/plugin_group/publish HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

This legacy endpoint publishes the staged PluginGroup configuration to the published namespace. It is independent of API route publication.

## IV. OPA Policy

OPA policy APIs proxy requests to the OPA server. If `server_url` or `policy_id` is not provided, defaults are used (`http://opa:8181` and `pixiu-authz`).

### 4.1 Get OPA Policy

**Request**:

```http
GET /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

**Query Params**:

* `policy_id`: OPA policy id (optional)
* `server_url`: OPA server URL (optional)
* `bearer_token`: OPA bearer token (optional)

**Response**:

```json
{
  "code": "10001",
  "data": "package pixiu.authz\n\ndefault allow := false\n"
}
```

If the policy does not exist, `data` will be an empty string.

### 4.2 Create or Update OPA Policy

**Request**:

```http
PUT /config/api/opa/policy HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
cache-control: no-cache
```

**Form Data**:

```text
Content-Disposition: form-data; name="policy_id"
pixiu-authz

Content-Disposition: form-data; name="content"
package pixiu.authz

default allow := false
```

Optional form fields:

* `server_url`
* `bearer_token`

### 4.3 Delete OPA Policy

**Request**:

```http
DELETE /config/api/opa/policy?policy_id=pixiu-authz HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

**Query Params**:

* `policy_id`: OPA policy id (optional)
* `server_url`: OPA server URL (optional)
* `bearer_token`: OPA bearer token (optional)

## V. xDS Diagnostics

### 5.1 Get xDS Publication Status

This authenticated endpoint returns xDS listener availability, the last-good
snapshot version, resource counts, publication timestamps, the latest listener
or rejected-candidate error, and the explicit support status of each xDS
resource path.

```http
GET /config/api/xds/status HTTP/1.1
Host: 127.0.0.1:8080
token: <admin-jwt>
```

`ready` requires both a bound xDS listener and a published snapshot. `degraded`
indicates that listener startup or a newer snapshot candidate failed.

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
