---
name: pixiu-http-to-dubbo
description: |
  Map a REST/HTTP endpoint onto a backend Dubbo service through
  dubbo-go-pixiu. Use whenever the user wants to "expose Dubbo as HTTP",
  "REST gateway for Dubbo", "Dubbo generic invoke", configure
  `api_config.yaml`, `integrationRequest`, `mappingParams`, `mapType`,
  `opt.types` / `opt.values`, group / version, or says "how do I call
  my Dubbo service from curl".
  Use this even when the user only asks to tweak an existing yaml route —
  the combination of conf.yaml + api_config.yaml is where most "500 from
  pixiu" reports come from, and this skill encodes the invariants.
allowed-tools: [Read, Grep, Glob, Edit, Write, Bash]
metadata:
  version: "0.1.2"
  domain: config
  scope: generate-and-validate
  triggers: ["http to dubbo", "REST to Dubbo", "Dubbo gateway", "api_config.yaml", "integrationRequest", "Dubbo generic invoke", "mappingParams", "mapType", "opt.types", "expose Dubbo as HTTP"]
  pixiu_min_version: "0.6.0"
  role: specialist
---

# pixiu-http-to-dubbo — Wiring HTTP Clients to Dubbo Backends

pixiu's bread and butter is "let a REST client call a Dubbo service". On
the surface it's two yaml files; underneath, you are configuring a
generic Dubbo invoke, and most failures come from the three-way mismatch
between the client request shape, the Dubbo interface signature, and
the mapping rules in between.

## Two files, two jobs

- **`conf.yaml`** (bootstrap) — wires the HTTP listener, enables the
  `dgp.filter.httpconnectionmanager` network filter, and loads the
  Dubbo registry adapter. You almost always keep this close to the
  sample; changes are listener port, adapter type, and sometimes
  cluster definitions.
  In config answers, explicitly use canonical HCM Kind
  `dgp.filter.httpconnectionmanager` (never
  `dgp.filter.http.httpconnectionmanager`) and put
  `dgp.filter.http.apiconfig` before `dgp.filter.http.dubboproxy` or
  `dgp.filter.http.httpproxy` inside HCM `http_filters`.
- **`api_config.yaml`** — the per-endpoint mapping. One
  `resources[].methods[]` entry per HTTP route. Every field except the
  path comes from the Dubbo side; every mapping rule comes from the
  HTTP side.

This skill's job is to keep those two in sync.

## When to Use

Use this skill when the user wants to:
- Expose an existing Dubbo interface through HTTP (GET/POST/PUT/DELETE).
- Configure `integrationRequest` for an existing route.
- Fix a 500 / "generic invoke failed" / "no provider found" error that
  points at a yaml route.
- Add parameter-name mapping (query → method param, body field → method
  param, header → method param).
- Switch an existing HTTP proxy route to Dubbo backend.

**Do NOT use this skill for:**
- Writing a new filter in the chain → `pixiu-filter-author`.
- Adding a new registry adapter for a non-standard service registry —
  out of scope; the user will need to author the adapter directly
  against `pkg/common/extension/adapter/adapter.go`.
- Pure yaml audit with no new HTTP-to-Dubbo route — out of scope. This
  skill only validates the route/config fragments it generates or edits.

## Prerequisites

- pixiu ≥ 0.6.0.
- A Dubbo provider running somewhere reachable. The user must know the
  interface FQCN, method name, Java method signature, and the Dubbo
  `group` / `version`.
- Registry mode needs ZooKeeper/Nacos config under
  `dgp.filter.http.dubboproxy.config.dubboProxyConfig.registries`.
  Direct mode instead uses `integrationRequest.url` in `api_config.yaml`
  and must also declare `parameterTypes` and `serialization`.

## Steps

### Step 0 — Verify Context

1. Confirm pixiu version: `grep dubbo-go-pixiu go.mod`.
2. Locate existing config: typically `configs/conf.yaml` and
   `configs/api_config.yaml` in the pixiu repo, or user-specified
   paths.
3. Open `pkg/config/api_config.go` — that's the Go struct the yaml
   actually binds to. Current `DubboBackendConfig` includes
   `parameterTypes` and `serialization`; the legacy spelling
   `paramTypes` is not the current field name.
4. Open `pkg/filter/http/remote/dubbo_handler.go`,
   `pkg/client/dubbo/types.go`, `pkg/client/dubbo/typeconv.go`, and
   `pkg/client/dubbo/dubbo.go` before editing mappings. Do not cite
   legacy mapping helpers on branches where they no longer exist.
5. If the user says "Dubbo direct URL, no registry", use the current
   direct generic contract: `integrationRequest.url` plus explicit
   `protocol`, `parameterTypes`, and `serialization`.

### Step 1 — Gather the Required Things (STOP and ask)

You cannot write a valid `integrationRequest` without all applicable
items:

1. **HTTP method + path**: `POST /api/v1/user`, etc.
2. **Dubbo interface FQCN**: e.g. `com.example.UserProvider`.
3. **Dubbo method name**: e.g. `createUser`.
4. **Java method signature**: e.g.
   `User createUser(com.example.User user)` or
   `Page<Order> search(String tenant, OrderQuery q)`. Translate each
   argument into `parameterTypes` and supported pixiu `mapType` values
   (`string`, `int`, `long`, `double`, `boolean`, `object`, etc.).
5. **Dubbo `group` and `version`**: empty strings are allowed but must
   be explicit.
6. **Where each parameter comes from** on the HTTP side: `queryStrings.
   <name>`, `requestBody.<path>`, `headers.<name>`, `uri.<name>`.
7. **Registry/direct mode**: ZK/Nacos registry, or direct URL. Registry
   settings live in `conf.yaml`; direct provider address lives in
   `integrationRequest.url`.
8. **Direct-call extras**: for direct URL mode, the declared protocol,
   `parameterTypes`, and `serialization`.

For POJO arguments, also ask whether the HTTP JSON body includes the
Dubbo/Hessian class discriminator the provider expects, usually a
`class` field such as `"class": "com.example.User"`. Pixiu can map the
argument as `object`, but it cannot infer every provider-side POJO FQCN
from a yaml field that the current config struct ignores.

For primitives, collections, and nested POJOs, stay close to
`DubboHandler`, `typeconv.go`, and existing `api_config.yaml` examples
before finalizing `mappingParams`.

If the user says "direct URL", "no registry", or gives a provider
address such as `10.0.0.8:20880`, keep the Dubbo route as
`integrationRequest.requestType: dubbo` and put the provider address in
`integrationRequest.url`. Direct generic invoke also requires
`parameterTypes` and `serialization`; without them Pixiu returns direct
generic validation errors before invoking the provider.

Do NOT proceed until the required fields are answered. It is fine to propose
defaults and ask for confirmation — but **every field must be explicit
in the final yaml**.

### Step 2 — Shape the `api_config.yaml` entry

The template:

```yaml
name: pixiu
description: <what this gateway does>
resources:
  - path: /api/v1/user
    type: restful
    description: create a user
    methods:
      - httpVerb: POST
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
          # describe the HTTP shape the client will send:
          headers:
            - name: X-Request-Id
              required: false
          queryStrings: []
          requestBody:
            - contentType: application/json
              schema: user
              required: true
        integrationRequest:
          requestType: dubbo
          # Dubbo coordinates:
          applicationName: pixiu
          group: ""
          version: "1.0.0"
          interface: com.example.UserProvider
          method: createUser
          parameterTypes:
            - com.example.User
          # How to pull each Dubbo method argument from the HTTP request:
          mappingParams:
            - name: requestBody
              mapTo: "0"
              mapType: object
```

Rules the template hides:

- `requestType` **must** be `dubbo` for Dubbo routes. `http` is the
  pass-through mode; different code path.
- `requestType: triple` uses the same outbound builder, but direct
  Triple generic invoke requires the provider to expose generic `$invoke`;
  IDL-only Triple handlers may return `404 Not Found`.
- `mapTo` is either the **index** of the Dubbo method parameter
  (`mapTo: "0"` is the first argument) or one of the supported
  `opt.*` targets such as `opt.types`, `opt.values`, `opt.group`,
  `opt.version`, `opt.interface`, or `opt.method`. `opt.application`
  is deprecated in the current handler.
- `parameterTypes` is the preferred explicit Java signature. In
  registry mode, `opt.types` can still provide the signature when
  `parameterTypes` is omitted. In direct mode, `parameterTypes` is
  required.
- `mapType` is the conversion hint consumed by pixiu's Dubbo type
  conversion path. Supported values come from `constant.JTypeMapper`: `string`,
  `java.lang.String`, `char`, `short`, `int`, `long`, `float`,
  `double`, `boolean`, `java.util.Date`, `date`, `object`, and
  `java.lang.Object`. Use `object` / `java.lang.Object` for POJO or map
  payloads, and include provider-required class metadata in the body
  when the Dubbo serializer needs it.
- `paramTypes` is not the current field name. Use `parameterTypes`.
- Registry mode gets provider discovery from
  `dubboProxyConfig.registries`. Direct mode gets the provider address
  from `integrationRequest.url`; do not model it as a static cluster-only
  path unless the target branch proves that older direct filter is in use.

For dynamic/default generic routes where the HTTP client explicitly
sends the type list and value list, map body fields to `opt.types` and
`opt.values` instead of using top-level yaml fields:

```yaml
mappingParams:
  - name: requestBody.types
    mapTo: opt.types
  - name: requestBody.values
    mapTo: opt.values
```

When the user pastes legacy `paramTypes`, `groupType`, or old
`types/values` examples, translate the intent into current
`parameterTypes`, `mappingParams`, and `mapType` rules. Do not preserve
unknown fields just because the user supplied them; current
`IntegrationRequest` silently ignores fields not present in
`pkg/config/api_config.go`.

### Step 3 — Update `conf.yaml` if needed

You usually only touch `conf.yaml` when:

- This is the first registry-backed Dubbo route (need
  `dgp.filter.http.dubboproxy.config.dubboProxyConfig.registries`, and
  often the dynamic `dgp.adapter.dubboregistrycenter` if API definitions
  come from a registry).
- A new listener port or host is needed.
- The user wants direct/no-registry Dubbo access. In that case
  `api_config.yaml` must carry `integrationRequest.url`,
  `parameterTypes`, and `serialization`; `conf.yaml` still needs the
  listener/filter chain, but not a static Dubbo cluster endpoint for
  that direct provider.

The `dgp.filter.http.dubboproxy` registry block (ZK example, inside HCM
`http_filters`):

```yaml
- name: dgp.filter.http.dubboproxy
  config:
    dubboProxyConfig:
      registries:
        zk:
          protocol: zookeeper
          timeout: 3s
          address: 127.0.0.1:2181
          username: ""
          password: ""
      timeout_config:
        connect_timeout: 5s
        request_timeout: 5s
```

Add `dgp.adapter.dubboregistrycenter` only when API definitions are
dynamically discovered from a Dubbo registry, not for every static
`api_config.yaml` route.

For direct/no-registry cases, omit the registry adapter and put the
provider address in `api_config.yaml`:

```yaml
integrationRequest:
  requestType: dubbo
  url: dubbo://10.0.0.8:20880
  protocol: dubbo
  serialization: hessian2
  interface: com.example.UserProvider
  method: getUser
  group: ""
  version: "1.0.0"
  parameterTypes:
    - java.lang.String
  mappingParams:
    - name: queryStrings.name
      mapTo: "0"
      mapType: string
```

If the URL omits a scheme, `protocol` must still be set so Pixiu can
construct the direct Dubbo reference.

And the HTTP listener MUST have `dgp.filter.http.apiconfig` in its
`http_filters` list, followed by `dgp.filter.http.dubboproxy` (or
`dgp.filter.http.httpproxy` if some routes are pass-through). Filter
order matters — compare against the current config structs and examples.

### Step 4 — Validate

Before booting Pixiu, inspect the generated config directly:

1. Parse the yaml with an available local parser or by loading it through
   Pixiu's config path.
2. Check the current API config structs and examples for required keys,
   expected types, and the allowed shape of structured API mapping
   objects. `filter.config` remains intentionally permissive because
   individual filter plugins own their own config schemas.
3. For registry mode, check that `dgp.filter.http.dubboproxy` has usable
   `dubboProxyConfig.registries`. For direct mode, check
   `integrationRequest.url`, `protocol`, `parameterTypes`, and
   `serialization`.

If validation fails, fix *before* trying to boot pixiu — boot-time
errors are more cryptic than config-shape mistakes found by inspection.

### Step 5 — Smoke Test

1. `go run ./cmd/pixiu/... gateway start -c configs/conf.yaml -a configs/api_config.yaml`
   (or the binary equivalent).
2. `curl -v -X POST http://localhost:<port>/<path> ...`
3. If you get HTTP 5xx, inspect Pixiu server logs near the request
   timestamp and match the first error to the config area it references.

## Cross-Cutting Rules

### Always

- Base generated yaml on the current `pkg/config/api_config.go` and
  `pkg/filter/http/remote/dubbo_handler.go`, not on older examples. Use
  `parameterTypes`, not legacy `paramTypes`.
- Use supported `mapType` values exactly as pixiu's `constant.JTypeMapper`
  defines them. `String` and `bool` are not valid current-source values;
  use `string` / `java.lang.String` and `boolean`.
- Keep numeric `mappingParams[].mapTo` values aligned with Java method
  argument order. `mapTo: "0"` is the first Dubbo argument.
- Specify `group` and `version` explicitly, even if empty. Empty-string
  is valid; missing key is not.
- Put `dgp.filter.http.apiconfig` **before** any proxy filter in
  `http_filters`.
- Timeouts: give `methods[].timeout` a value shorter than the
  listener's `idle_timeout` in `conf.yaml`.

### Never

- Generate top-level `paramTypes`; current pixiu uses `parameterTypes`.
- Write `String` or `bool` as a `mapType`. Current pixiu accepts
  `string` / `java.lang.String` and `boolean`.
- Mix `requestType: http` and Dubbo fields in the same
  `integrationRequest`. Pick one.
- Put a registry-backed Dubbo provider address in HTTP backend fields
  such as `host` or `path`. For direct/no-registry Dubbo calls, use
  `integrationRequest.url` and also set `parameterTypes` and
  `serialization`.
- Put the Dubbo registry configuration in `api_config.yaml`. It lives
  in `conf.yaml` under `dgp.filter.http.dubboproxy` and, for dynamic API
  discovery, under `adapters`.
- Use the literal object from the request body as a single arg when
  the Dubbo side expects separate primitives. Use
  `mappingParams[].mapTo` with an index per arg.

## Common Pitfalls

1. **`no provider found`** at boot — registry is up but the adapter
   cannot resolve the interface. Check `interface`, `group`, `version`
   spelling; `group` is case-sensitive.
2. **500 with `generic invoke failed: ClassNotFound`** — type strings
   came from `parameterTypes`, `opt.types`, or request body `types` and
   do not match classes visible to the provider. Prefer exact Java
   class names and provider-required POJO class metadata in the JSON
   body.
3. **Response is `{}` / empty** — `mappingParams` is empty or wrong,
   so the Dubbo method got `null` args and returned its default value.
4. **Timeout set in the wrong layer** — `methods[].timeout` becomes the
   Dubbo call deadline; `dubboProxyConfig.timeout_config` controls the
   Dubbo client side. If the Dubbo call is slow, check the method
   timeout first, then the proxy timeout config.
5. **Header/query name case** — `inboundRequest.headers[].name` is
   case-insensitive on match but case-sensitive in logs; normalize in
   docs.
6. **Nacos namespace confusion** — `namespace` vs `namespaceId` vs
   `group` vs `groupName`. pixiu's Nacos adapter uses `namespace` and
   `group`; if you are setting up a new adapter from scratch, read
   the adapter source under `pkg/adapter/` for the actual field
   names.

## Source Files To Read

- `pkg/config/api_config.go`
- `pkg/filter/http/remote/dubbo_handler.go`
- `pkg/client/dubbo/types.go`
- `pkg/client/dubbo/typeconv.go`
- `pkg/client/dubbo/dubbo.go`
- Existing `configs/api_config.yaml` and sample `api_config.yaml`
  files.
- Pixiu server logs for the first error near a failing request.
