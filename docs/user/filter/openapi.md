# OpenAPI Request Validation Filter

English | [中文](openapi_CN.md)

---

## Overview

Pixiu can load a local OpenAPI 3.0/3.1/3.2 file in `dgp.filter.http.openapi` and validate matching requests before they
are forwarded upstream.

Official references:

- [libopenapi](https://github.com/pb33f/libopenapi)
- [libopenapi validation](https://pb33f.io/libopenapi/validation/)
- [libopenapi validator](https://github.com/pb33f/libopenapi-validator)

If validation fails, Pixiu returns a local `400 Bad Request` and stops the filter chain.

## Wired In This Filter

- local OpenAPI 3.0/3.1/3.2 file loading
- request matching by OpenAPI path and method, including templated paths such as `/users/{id}`
- path, query, and header parameter validation
- JSON request body validation
- OpenAPI schema constraints enforced by the SDK, including `required`, `type`, `enum`, `minimum`, `maximum`, `minLength`, and `maxLength`

This filter uses `libopenapi` to parse the spec and `libopenapi-validator` to validate incoming requests.

## Not Wired In This Filter

- response validation
- route creation or `api_config` route matching
- OpenAPI `security` validation; use dedicated authentication or authorization filters for auth checks
- admin or config-center distribution of OpenAPI files

## Example Filter Config

```yaml
- name: dgp.filter.http.apiconfig
  config:
    path: configs/api_config.yaml

- name: dgp.filter.http.openapi
  config:
    path: configs/openapi_users.yaml
```

`dgp.filter.http.apiconfig` and `dgp.filter.http.openapi` are independent filters. `apiconfig` matches Pixiu API routes
and writes API metadata into the request context. `openapi` validates only the operations declared in the OpenAPI file.

If a request path and method are not declared in the OpenAPI file, this filter skips validation and lets the request
continue. If the operation is declared but the request violates parameters or body schema, Pixiu returns
`400 Bad Request`.

## Notes

- `libopenapi-validator` is an opt-in companion module; `libopenapi` handles parsing and model building.
- Parameter-level validation covers common scalar constraints for `path`, `query`, and `header` parameters.
- The keywords listed above come from the OpenAPI schema and are enforced by the SDK path, not by custom in-repo validators.
- OpenAPI `security` validation is disabled for this filter, so auth remains the responsibility of filters such as JWT,
  OPA, SAML, or other dedicated authentication and authorization filters.
- The configured OpenAPI `path` must be relative. Absolute paths, parent-directory segments such as `..`, and sensitive
  base directories are rejected during filter startup.
- Relative file references are resolved from the OpenAPI file location, so file-based loading keeps local `$ref` paths intact.
- Invalid OpenAPI documents, including unresolved `$ref` targets, fail during filter startup instead of running with a
  partial validator model.
- The validator package also exposes response and document validation APIs, but this filter only calls the request-validation path.

## Runtime Flow

1. Pixiu loads the OpenAPI file during `openapi.Apply()` and builds an SDK validator.
2. A request enters `openapi.Decode()`.
3. Pixiu checks whether the OpenAPI document declares the request path and method.
4. If the operation is not declared, validation is skipped and the request continues.
5. If the operation is declared, Pixiu validates the request through `libopenapi-validator`.
6. If validation succeeds, the request continues to later filters.
7. If validation fails, Pixiu responds with `400 Bad Request`.

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
