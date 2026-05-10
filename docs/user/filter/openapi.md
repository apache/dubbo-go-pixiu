# OpenAPI Request Validation Filter

English | [中文](openapi_CN.md)

---

## Overview

Pixiu can load an OpenAPI 3.x file during `dgp.filter.http.apiconfig` initialization and merge each operation's
`ValidationPlan` into the matching route that is already defined by `api_config`.

When a request matches an API, Pixiu validates the request before it is forwarded upstream.

If validation fails, Pixiu returns a `400 Bad Request` locally and stops the filter chain.

## Supported In V1

- local OpenAPI 3.x file loading
- path parameter extraction
- query parameter validation
- header validation
- JSON request body validation
- `required`
- `type`
- `enum`
- `minimum`
- `maximum`
- `minLength`
- `maxLength`

## Not Included In V1

- response validation
- remote `$ref`
- `oneOf` / `allOf` / `anyOf`
- admin or config-center distribution of OpenAPI files
- `dynamic + openapi_path` combined configuration

## Example Filter Config

```yaml
- name: dgp.filter.http.apiconfig
  config:
    path: configs/api_config.yaml
    openapi_path: configs/openapi_users.yaml
    enable_openapi_validation: true
```

OpenAPI validation does not create standalone routes. The route must already exist in `api_config`.

## Notes

- Parameter-level validation covers common scalar constraints for `path`, `query`, and `header` parameters.
- V1 does not support using `dynamic` together with `openapi_path` in the same filter config.

## Runtime Flow

1. Pixiu loads `api_config` during `apiconfig.Apply()`.
2. Pixiu loads the OpenAPI file and compiles validation plans.
3. Each compiled `ValidationPlan` is merged into the matching `router.API.Metadata`.
4. A request enters `apiconfig.Decode()`.
5. Pixiu matches the request path and method.
6. Pixiu extracts the `ValidationPlan` from the matched API metadata.
7. Pixiu validates the request.
8. If validation succeeds, the request continues to later proxy filters.
9. If validation fails, Pixiu responds with `400 Bad Request`.

## Example Requests

Assume the config points at `configs/openapi_users.yaml`.

Valid request:

```http
POST /users?source=web
Content-Type: application/json

{"name":"tom","role":"admin","age":18}
```

Invalid request:

```http
POST /users
Content-Type: application/json

{"name":"tom","role":"admin"}
```

The invalid request is rejected because `source` is required by the OpenAPI query parameter definition.
