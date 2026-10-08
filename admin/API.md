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

API Router configuration is managed through AdminRouteBinding objects. The legacy Resource/Method Admin CRUD API is replaced. Existing runtime Resource/Method records are not imported into the new draft model; recreate routes that need to be managed through this API. Publishing a new route continues to generate the legacy runtime configuration used by Pixiu.

### 2.1 Get API Route Schema

**Request**:

```http
GET /config/api/route/schema HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

The response data contains the registered Admin object schemas, including the AdminRouteBinding schema used by the editor.

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

### 2.8 Publish One API Route

**Request**:

```http
PUT /config/api/route/publish?name=<route-name>&expectedRevision=<revision> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Publishing validates the route and atomically updates its published AdminRouteBinding and generated runtime configuration in one etcd transaction. expectedRevision is optional and can be used to reject a stale draft. There is no full API-route publish endpoint.

Before publishing, Admin checks for conflicts with legacy runtime Resource/Method records. A conflicting HTTP method and path must be removed from the old runtime configuration before the new route can be published.

### 2.9 Delete One API Route

**Request**:

```http
DELETE /config/api/route?name=<route-name>&expectedRevision=<revision> HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
```

Deletion immediately removes the route draft and its published runtime configuration in one etcd transaction. The client should ask for confirmation before sending this request.
expectedRevision is optional; when supplied, it prevents deleting a draft that changed after it was read.

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
