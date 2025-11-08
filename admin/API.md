# Backend API Documentation

**English** | [中文](README_CN.md)

This API documentation describes the backend operations of the Pixiu management platform, including the APIs for retrieving, creating, modifying, and deleting resources (Resource), methods (Method), and plugin groups (PluginGroup). Pixiu provides a complete set of APIs to help users manage API gateway resource mappings, plugin configurations, and request handling. The examples in this document cover common request and response formats and show how to test the APIs using Postman.

Whether you are creating new resources, modifying existing configurations, or managing plugin groups, this document provides clear steps and necessary API details, making it easier for developers to get started and integrate quickly.

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
Postman-Token: deb6b451-f211-41fd-8cdb-0801b13bab69
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
Postman-Token: a5867037-5bc0-4b9f-aef2-57980ddc9dbf
```

**Form Data**:

```text
Content-Disposition: form-data; name="content"
name: pixiu
description: pixiu111 sample
```

## II. Resource

### 2.1 Get Resource List

**Request**:

```http
GET /config/api/resource/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 19ffe51c-8193-4722-b94b-88d2502e3046
```

### 2.2 Get Resource Details

**Request**:

```http
GET /config/api/resource/detail?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 6b97b70d-70fc-4633-9ba7-dee3e08d7468
```

### 2.3 Create Resource

**Request**:

```http
POST /config/api/resource/ HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
cache-control: no-cache
Postman-Token: 4e7b85b1-a969-46a4-b3d2-3df5d2f2073e
```

**Form Data**:

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

### 2.4 Modify Resource

**Request**:

```http
PUT /config/api/resource? HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 67d3fb19-bb66-4e92-be1b-419fdd3bcb28
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**Form Data**:

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

### 2.5 Delete Resource

**Request**:

```http
DELETE /config/api/resource/?resourceId=2 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: faed927d-2857-40a2-9cfb-1a3ce2f704e4
```

## III. Method Related

### 3.1 Get Method List for a Resource

**Request**:

```http
GET /config/api/resource/method/list?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 90c3df1d-02b4-42af-8846-f4c74a36fdae
```

### 3.2 Get Method Details

**Request**:

```http
GET /config/api/resource/method/detail?resourceId=1&methodId=2 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 6d0909cc-1585-46c5-975f-e4a0d0b2f490
```

### 3.3 Create Method

**Request**:

```http
POST /config/api/resource/method/?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 0d5e2885-67e6-41cb-84e8-777a8d077384
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**Form Data**:

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

### 3.4 Modify Method

**Request**:

```http
PUT /config/api/resource/method/?resourceId=1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 9a635d6d-b52a-4cf5-8f1a-a7d450cc3552
Content-Type: multipart/form-data; boundary=-WebKitFormBoundary7MA4YWxkTrZu0gW
```

**Form Data**:

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

### 3.5 Delete Method

**Request**:

```http
DELETE /config/api/resource/method/?resourceId=1&methodId=2 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 83930da3-7ab5-47b4-a47e-ef0ceeb6e9af
```

## IV. PluginGroup and Plugin Related

### 4.1 Get PluginGroup List

**Request**:

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: f48aa74f-38c9-4ee9-87a3-c4a35a2350f9
```

### 4.2 Get PluginGroup Details

**Request**:

```http
GET /config/api/plugin_group/list HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 739ee1d2-7801-4005-90a3-cc2dfb03b5e7
```

### 4.3 Create PluginGroup

**Request**:

```http
POST /config/api/plugin_group/ HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 1e41964b-989b-4f03-8512-9dc7e9b7635d
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

### 4.4 Modify PluginGroup

**Request**:

```http
PUT /config/api/plugin_group/ HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: 5166a7df-b838-4941-add3-5073a6e430f0
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

### 4.5 Delete PluginGroup

**Request**:

```http
DELETE /config/api/plugin_group/?name=group1 HTTP/1.1
Host: 127.0.0.1:8080
cache-control: no-cache
Postman-Token: d761e1fa-02c4-4693-b765-b16ce389c9b2
```
