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
- A registry: ZooKeeper, Nacos, or an in-process ("direct URL")
  endpoint. The registry configuration lives in `conf.yaml` under the
  Dubbo adapter, not in `api_config.yaml`.

## Steps

### Step 0 — Verify Context

1. Confirm pixiu version: `grep dubbo-go-pixiu go.mod`.
2. Locate existing config: typically `configs/conf.yaml` and
   `configs/api_config.yaml` in the pixiu repo, or user-specified
   paths.
3. Open `pkg/config/api_config.go` — that's the Go struct the yaml
   actually binds to. Fields not present there are silently ignored;
   always cross-reference. Current pixiu has `mappingParams` but no
   top-level `paramTypes` field on `IntegrationRequest`, so do not
   generate `paramTypes` unless a target version proves it exists.
4. Open `pkg/client/dubbo/default.go`, `pkg/client/dubbo/mapper.go`,
   and `pkg/client/dubbo/option.go` before editing mappings. They define
   the accepted `mapTo` grammar and how `mapType` becomes the generic
   invocation type list.
5. If the user says "Dubbo direct URL, no registry", make sure you
   understand they still need `clusters[]` in `conf.yaml` with the
   provider address.

### Step 1 — Gather the Seven Things (STOP and ask)

You cannot write a valid `integrationRequest` without all of:

1. **HTTP method + path**: `POST /api/v1/user`, etc.
2. **Dubbo interface FQCN**: e.g. `com.example.UserProvider`.
3. **Dubbo method name**: e.g. `createUser`.
4. **Java method signature**: e.g.
   `User createUser(com.example.User user)` or
   `Page<Order> search(String tenant, OrderQuery q)`. Translate each
   argument into a supported pixiu `mapType` (`string`, `int`, `long`,
   `double`, `boolean`, `object`, etc.) rather than a top-level
   `paramTypes` list.
5. **Dubbo `group` and `version`**: empty strings are allowed but must
   be explicit.
6. **Where each parameter comes from** on the HTTP side: `queryStrings.
   <name>`, `requestBody.<path>`, `headers.<name>`, `uri.<name>`.
7. **Registry type**: ZK, Nacos, direct URL. This lives in `conf.yaml`.

For POJO arguments, also ask whether the HTTP JSON body includes the
Dubbo/Hessian class discriminator the provider expects, usually a
`class` field such as `"class": "com.example.User"`. Pixiu can map the
argument as `object`, but it cannot infer every provider-side POJO FQCN
from a yaml field that the current config struct ignores.

For primitives, collections, and nested POJOs, stay close to the current
mapper source and existing `api_config.yaml` examples before finalizing
`mappingParams`.

If the user says "direct URL", "no registry", or gives a provider
address such as `10.0.0.8:20880`, keep the Dubbo route as
`integrationRequest.requestType: dubbo` and put only `clusterName` in
`api_config.yaml`. The provider address belongs in `conf.yaml` as a
static/direct Dubbo cluster endpoint; do not copy it into the HTTP
backend `url` / `host` / `path` fields.

Do NOT proceed until all seven are answered. It is fine to propose
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
          clusterName: dubbo-cluster          # matches conf.yaml clusters[].name
          applicationName: pixiu
          group: ""
          version: "1.0.0"
          interface: com.example.UserProvider
          method: createUser
          # How to pull each Dubbo method argument from the HTTP request:
          mappingParams:
            - name: requestBody
              mapTo: "0"
              mapType: object
```

Rules the template hides:

- `requestType` **must** be `dubbo` for Dubbo routes. `http` is the
  pass-through mode; different code path.
- `mapTo` is either the **index** of the Dubbo method parameter
  (`mapTo: "0"` is the first argument) or one of the supported
  `opt.*` targets such as `opt.types`, `opt.values`, `opt.group`,
  `opt.version`, `opt.interface`, `opt.application`, or `opt.method`.
- `mapType` is the generic invoke type hint consumed by pixiu's Dubbo
  mapper. Supported values come from `constant.JTypeMapper`: `string`,
  `java.lang.String`, `char`, `short`, `int`, `long`, `float`,
  `double`, `boolean`, `java.util.Date`, `date`, `object`, and
  `java.lang.Object`. Use `object` / `java.lang.Object` for POJO or map
  payloads, and include provider-required class metadata in the body
  when the Dubbo serializer needs it.
- Do not emit top-level `paramTypes` for current pixiu. Old examples may
  contain it, but `pkg/config/api_config.go` does not bind it on
  `IntegrationRequest`, so it is ignored unless the target branch proves
  otherwise.
- `clusterName` cross-references `conf.yaml`'s
  `static_resources.clusters[].name`. A typo here = "cluster not found".

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
`mappingParams` and `mapType` rules. Do not preserve unknown fields just
because the user supplied them; current `IntegrationRequest` silently
ignores fields not present in `pkg/config/api_config.go`.

### Step 3 — Update `conf.yaml` if needed

You usually only touch `conf.yaml` when:

- This is the first Dubbo route (need the Dubbo registry adapter).
- A new `clusterName` is introduced.
- A new listener port or host is needed.
- The user wants direct/no-registry Dubbo access. In that case declare a
  static/direct cluster endpoint in `conf.yaml` and keep
  `api_config.yaml` focused on the Dubbo interface/method mapping.

The Dubbo adapter block (ZK example):

```yaml
adapters:
  - id: dubbo
    name: dgp.adapter.dubboregistrycenter
    config:
      registries:
        zk:
          protocol: zookeeper
          timeout: 3s
          address: 127.0.0.1:2181
          username: ""
          password: ""
```

For direct/no-registry cases, omit the registry adapter and declare the
target cluster explicitly:

```yaml
static_resources:
  clusters:
    - name: direct-dubbo
      lb_policy: RoundRobin
      endpoints:
        - id: direct-provider-1
          socket_address:
            address: "10.0.0.8"
            port: 20880
```

The matching `api_config.yaml` still uses
`integrationRequest.clusterName: direct-dubbo`; it does not contain the
provider URL.

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
3. Cross-check every `clusterName` referenced in `integrationRequest`
   against the adapter / cluster `id` or `name` values in `conf.yaml`.

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
  `pkg/client/dubbo/*` code, not on older examples. If
  `IntegrationRequest` has no `ParamTypes` field, treat top-level
  `paramTypes` as legacy and do not generate it.
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

- Generate top-level `paramTypes` for current pixiu unless Step 0 proves
  the target branch has that field on `IntegrationRequest`.
- Write `String` or `bool` as a `mapType`. Current pixiu accepts
  `string` / `java.lang.String` and `boolean`.
- Mix `requestType: http` and Dubbo fields in the same
  `integrationRequest`. Pick one.
- Put a Dubbo provider address in `integrationRequest.url`, `host`, or
  `path`. Those fields belong to HTTP pass-through backends, not Dubbo
  generic invoke routes.
- Put the Dubbo registry configuration in `api_config.yaml`. It lives
  in `conf.yaml` under `adapters`.
- Use the literal object from the request body as a single arg when
  the Dubbo side expects separate primitives. Use
  `mappingParams[].mapTo` with an index per arg.

## Common Pitfalls

1. **`no provider found`** at boot — registry is up but the adapter
   cannot resolve the interface. Check `interface`, `group`, `version`
   spelling; `group` is case-sensitive.
2. **500 with `generic invoke failed: ClassNotFound`** — type strings
   came from legacy `paramTypes`, `opt.types`, or request body `types`
   and do not match classes visible to the provider. For static routes,
   prefer supported `mapType` keys and provider-required POJO class
   metadata in the JSON body.
3. **Response is `{}` / empty** — `mappingParams` is empty or wrong,
   so the Dubbo method got `null` args and returned its default value.
4. **Per-method timeout overridden by cluster timeout** —
   `methods[].timeout` governs the Dubbo call deadline; the cluster's
   `timeout` is the connection-level. If the Dubbo call is slow, bump
   the method timeout first.
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
- `pkg/client/dubbo/mapper.go`
- `pkg/client/dubbo/option.go`
- Existing `configs/api_config.yaml` and sample `api_config.yaml`
  files.
- Pixiu server logs for the first error near a failing request.
