---
name: pixiu-http-to-dubbo
description: Creates and debugs dubbo-go-pixiu HTTP-to-Dubbo route YAML. Use for api_config.yaml/conf.yaml routes, direct Dubbo URLs, registry-backed Dubbo, mappingParams, mapType, parameterTypes, opt.types/opt.values, or generic-invoke errors. Do not use for non-Dubbo HTTP proxy routes or new registry adapter implementation.
---

# pixiu-http-to-dubbo

## Purpose
Generate a complete or embeddable `api_config.yaml`, and add the required companion `conf.yaml` fragments.

## When to use
- Use when:
  - Creating or modifying HTTP-to-Dubbo routes, including registry/direct mode switching, parameter mapping, type conversion, generic-invoke config, or debugging errors such as `no provider found`, `ClassNotFound`, empty responses, and signature mismatches.

- Do not use for:
  - New registry adapter implementation or non-Dubbo HTTP proxy routes.

## Inputs

<HARD-GATE>
Do not generate YAML, write code, create files, or take any implementation action until the user has provided all required inputs. This is a first principle.

If any required input is missing, this turn must only ask for the missing fields in the current input group; do not generate examples, defaults, YAML, code, or final output.

Even if configuration information seems inferable, obvious, or implied by context, you must still ask the user to confirm it. Do not proceed until the user confirms it.
</HARD-GATE>

- Choose Dubbo provider mode (required):
  - Required:
    - `provider_mode`: `registry` or `direct`
- `listener` (required):
  - Defaults:
    - `address.socket_address.address: 0.0.0.0`
    - `address.socket_address.port: 8888`
- `route_config` (required):
  - Defaults:
    - `routes[].match.prefix: /api/v1`
    - `routes[].route.cluster: dubbo-route`
- `dgp.filter.http.apiconfig` (required):
  - Defaults:
    - `config.path: configs/api_config.yaml`
    - `config.dynamic: false`
- `dgp.filter.http.dubboproxy` (required):
  - Defaults:
    - `dubboProxyConfig.timeout_config.connect_timeout: 5s`
    - `dubboProxyConfig.timeout_config.request_timeout: 5s`
    - `dubboProxyConfig.load_balance: roundrobin`
    - `dubboProxyConfig.retries: "3"`
- `api_config.yaml` resource/method (required):
  - Required:
    - `resources[].path`
    - `methods[].httpVerb`
    - `integrationRequest.interface`
    - `integrationRequest.method`
  - Optional:
    - `resources[].description`
    - `integrationRequest.group`
    - `integrationRequest.version`
    - `integrationRequest.applicationName`
    - `integrationRequest.clusterName`
  - Defaults:
    - `resources[].type: restful`
    - `resources[].timeout: 1000ms`
    - `methods[].enable: true`
    - `methods[].timeout: 1000ms`
    - `inboundRequest.requestType: http`
    - `integrationRequest.requestType: dubbo`
- Dubbo method parameters (required when the method has parameters):
  - Required:
    - `mappingParams[]` using exactly one argument style:
      - positional mappings: `mappingParams[].mapTo: "0"`, `"1"`, ...
      - `opt.values` mappings: map request fields to `opt.values`, and map or declare Java types separately when needed
  - Conditionally required:
    - `opt.types` when using `opt.values` and Java types cannot be inferred or must be explicit.
    - `parameterTypes[]` for direct-mode positional mappings, or whenever Java types must override inferred types.
    - When method parameters contain POJOs, nested objects, or collections, provide the corresponding Java class names and field types.
- Registry provider (required when `provider_mode: registry`; defaults use Zookeeper, Nacos can also be chosen, examples below use Zookeeper):
  - Required:
    - `dubboProxyConfig.registries.zookeeper.address`
  - Optional:
    - `dubboProxyConfig.registries.zookeeper.group`
    - `dubboProxyConfig.registries.zookeeper.namespace`
    - `dubboProxyConfig.registries.zookeeper.username`
    - `dubboProxyConfig.registries.zookeeper.password`
  - Defaults:
    - `dubboProxyConfig.registries.zookeeper.protocol: zookeeper`
    - `dubboProxyConfig.registries.zookeeper.registry_type: interface`
    - `dubboProxyConfig.registries.zookeeper.timeout: 3s`
    - `dubboProxyConfig.registries.zookeeper.group: DEFAULT_GROUP`
- Direct provider (required when `provider_mode: direct`):
  - Required:
    - `integrationRequest.url`
  - Conditionally required:
    - `integrationRequest.parameterTypes` (use `[]` for no-argument methods; required when `opt.types` is not used)
  - Optional:
    - `integrationRequest.protocol` (needed when `url` is `host:port`; when `url` has a scheme, it must match that scheme)
  - Defaults:
    - `integrationRequest.protocol: dubbo`
    - `integrationRequest.serialization: hessian2`

## Workflow
1. Check required Inputs first:
   1. Read existing config and already provided user information first; do not ask again for information already present or stated.
   2. Ask one `Inputs` group at a time. Each time, output only that group's required fields, with a short explanation after each field.
   3. If the current group's required fields are incomplete, ask only for the missing fields and do not move to the next group.
   4. After the current group's required fields are complete, ask whether to fill that group's optional fields; if yes, list those optional fields with short explanations.
   5. After optional fields are skipped or completed, apply that group's defaults and output default-value information; defaults must not override existing config or user input.
2. Read current source before editing YAML:
   - `pkg/config/api_config.go`.
   - `pkg/filter/http/remote/dubbo_handler.go`.
   - `pkg/client/dubbo/types.go`, `pkg/client/dubbo/typeconv.go`, and `pkg/client/dubbo/dubbo.go`.
3. Choose registry or direct path by `provider_mode`:
   - Registry uses the registry config in `dgp.filter.http.dubboproxy`.
   - Direct uses `integrationRequest.url`.
4. Generate `api_config.yaml`: write `resources[].methods[]` for each HTTP API, and use either positional `mappingParams` or `opt.values`/`opt.types` to organize Dubbo arguments.
5. Complete `conf.yaml`: use HCM for route and filters, and ensure `dgp.filter.http.apiconfig` is before the Dubbo proxy.

## Output format
- Show the relevant YAML fragments.

## Validation
- When using positional mappings, verify `mappingParams[].mapTo` indexes match the declared Java type order exactly.
- Verify positional mappings and `opt.values` are not mixed in the same method.
- Verify `dgp.filter.http.apiconfig` is ordered before `dgp.filter.http.dubboproxy`.
- Verify `mapType` is supported by `client/dubbo.MapTypes` after normalization against `constant.JTypeMapper`: lowercase primitives and supported fully-qualified Java type aliases such as `java.lang.String`/`java.lang.Integer` are valid; short uppercase aliases such as `String`/`Boolean` are not.
- Verify POJO arguments have a `class` field in the corresponding JSON body, such as `{"class":"com.example.User", ...}`.

## Examples
Direct mode (`api_config.yaml`):

```yaml
name: pixiu
description: HTTP facade for Dubbo user service
resources:
  - path: /users/:name
    type: restful
    timeout: 5s
    description: User APIs
    methods:
      - httpVerb: GET
        enable: true
        timeout: 5s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          url: dubbo://127.0.0.1:20880
          protocol: dubbo
          serialization: hessian2
          interface: com.example.UserService
          method: getUser
          group: ""
          version: "1.0.0"
          parameterTypes:
            - java.lang.String
          mappingParams:
            - name: uri.name
              mapTo: "0"
              mapType: string
```

Direct mode (`conf.yaml`):

```yaml
http_filters:
  - name: dgp.filter.http.apiconfig
    config:
      path: configs/api_config.yaml
      dynamic: false
  - name: dgp.filter.http.dubboproxy
    config:
      dubboProxyConfig:
        timeout_config:
          connect_timeout: 5s
          request_timeout: 5s
        load_balance: roundrobin
        retries: "3"
```

Registry mode (`api_config.yaml`):

```yaml
integrationRequest:
  requestType: dubbo
  interface: com.example.UserService
  method: getUser
  group: ""
  version: "1.0.0"
  parameterTypes:
    - java.lang.String
  mappingParams:
    - name: uri.name
      mapTo: "0"
      mapType: string
```

Registry mode (`conf.yaml`):

```yaml
http_filters:
  - name: dgp.filter.http.apiconfig
    config:
      path: configs/api_config.yaml
      dynamic: false
  - name: dgp.filter.http.dubboproxy
    config:
      dubboProxyConfig:
        registries:
          zookeeper:
            protocol: zookeeper
            address: 127.0.0.1:2181
            timeout: 3s
            group: DEFAULT_GROUP
            registry_type: interface
        timeout_config:
          connect_timeout: 5s
          request_timeout: 5s
        load_balance: roundrobin
        retries: "3"
```

Minimal positional parameter mapping:

```yaml
integrationRequest:
  requestType: dubbo
  interface: com.example.UserService
  method: getUser
  parameterTypes:
    - java.lang.String
  mappingParams:
    - name: uri.name
      mapTo: "0"
      mapType: string
```

Minimal `opt.values` / `opt.types` mapping:

```yaml
integrationRequest:
  requestType: dubbo
  interface: com.example.UserService
  method: getUser
  mappingParams:
    - name: requestBody.values
      mapTo: opt.values
    - name: requestBody.types
      mapTo: opt.types
```
