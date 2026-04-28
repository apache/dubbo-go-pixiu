# `api_config.yaml` — Annotated Schema

Source of truth: `pkg/config/api_config.go` (types `APIConfig`,
`Resource`, `Method`, `InboundRequest`, `IntegrationRequest`,
`MappingParam`, `DubboBackendConfig`, `HTTPBackendConfig`). This
document annotates the yaml fields that map to those types.

## Top level

```yaml
name: <string>              # human-readable gateway name (free-form)
description: <string>       # free-form
resources:
  - ...                     # one entry per HTTP path prefix / path
pluginFilePath: ""          # rare; path to external Go plugin .so
pluginsGroup:               # optional groups referenced from resources[].plugins
  - groupName: ...
    plugins: [...]
```

Only `name` and `resources` matter in practice. Everything else has
workable defaults.

## `resources[]`

```yaml
- path: /api/v1/user              # required; HTTP path to match
  type: restful                   # required; "restful" is the common value
  description: <free-form>
  timeout: 500ms                  # resource-level default (methods override)
  plugins:                        # optional pre/post plugin chains
    pre:
      pluginNames: [...]
    post:
      groupNames: [...]
  methods:
    - ...                         # one entry per HTTP verb
```

Path matching uses pixiu's router. Variables look like `/user/:id`.

## `resources[].methods[]`

This is where 90% of the configuration lives.

```yaml
- httpVerb: POST                  # GET/POST/PUT/DELETE/...
  enable: true                    # false disables without removing
  timeout: 1000ms                 # method-level; wins over resource-level
  inboundRequest: ...
  integrationRequest: ...
```

### `inboundRequest`

Describes what the HTTP client sends. pixiu uses this for request
validation (when `required: true`) and as the source for
`mappingParams`.

```yaml
inboundRequest:
  requestType: http               # only valid value today
  headers:
    - name: X-Trace-Id
      required: false
  queryStrings:
    - name: userId
      required: true
  requestBody:
    - contentType: application/json
      schema: userCreateReq       # optional; references definitions
      required: true
```

### `integrationRequest`

Describes what pixiu sends upstream.

For Dubbo:

```yaml
integrationRequest:
  requestType: dubbo              # <- picks the Dubbo proxy path
  # Dubbo backend fields (inlined from DubboBackendConfig):
  clusterName: dubbo-cluster
  applicationName: pixiu
  protocol: dubbo                 # or "tri" for Triple
  group: ""
  version: 1.0.0
  interface: com.example.UserProvider
  method: createUser
  retries: "0"                    # as string; "3" = three retries
  # Arguments:
  mappingParams:
    - name: requestBody           # source on the HTTP side
      mapTo: "0"                  # index into the Dubbo param list
      mapType: object             # generic invoke type hint
```

For HTTP (pass-through):

```yaml
integrationRequest:
  requestType: http
  host: backend.example.com:8080
  path: /internal/user            # rewrite; empty = keep inbound path
  schema: http
  mappingParams:                  # optional; map HTTP→HTTP fields
    - name: queryStrings.name
      mapTo: queryStrings.name
```

The two forms share the `requestType` key; that single value controls
which code path pixiu takes. Mixing fields from both forms leads to the
fields being silently ignored, not to a validation error — this is the
single most common way to get a "500 with no obvious cause".

## `mappingParams` grammar (summary)

Source side (`name`):

- `queryStrings.<qs-name>`
- `requestBody.<jsonpath>` — dot-separated JSON path into the parsed body
- `headers.<header-name>`
- `uri.<path-var>` — matches `:id`-style variables in the route path

Target side (`mapTo`):

- For Dubbo: a **numeric index as a string**, matching the position in
  the Java method signature. Start from `"0"`.
- For Dubbo generic options: `opt.values`, `opt.types`, `opt.group`,
  `opt.version`, `opt.interface`, `opt.application`, or `opt.method`.
- For HTTP: a path in the downstream request, e.g.
  `queryStrings.<name>`, `requestBody.<path>`.

`mapType` values are the keys accepted by current
`pkg/common/constant/jtypes.go`: `string`, `java.lang.String`, `char`,
`short`, `int`, `long`, `float`, `double`, `boolean`,
`java.util.Date`, `date`, `object`, and `java.lang.Object`. Use
`object` / `java.lang.Object` for POJO or nested structures. If the
Dubbo provider requires Hessian class metadata, include the class
discriminator in the JSON body, for example `"class":
"com.example.User"`.

See `param-mapping-rules.md` for concrete worked examples.

## What is NOT in the yaml

- **Registry configuration** (ZK/Nacos addresses). Lives in
  `conf.yaml`, under the Dubbo adapter.
- **Filter chain order**. Lives in `conf.yaml`, under
  `http_filters`.
- **Authentication**. Handled by filters (JWT / bearer / OPA),
  configured in `conf.yaml`.

If the user tries to put any of these in `api_config.yaml`, the fields
are ignored at load time without a warning.

## Fields the Go struct does NOT have (do not use)

The following keys have circulated in old samples and blog posts but
are not recognized by current pixiu. If you see them, move the semantics
elsewhere or delete.

- `paramTypes` — old samples may include this under
  `integrationRequest`, but current `IntegrationRequest` does not bind
  it. For static routes, put the type hint in `mappingParams[].mapType`;
  for dynamic generic routes, map a request field to `opt.types`.
- `groupType` / `paramValues` — old style, replaced by
  `mappingParams` plus `mapType` / `opt.types` / `opt.values`.
- `authority` inside `integrationRequest` — move to a filter.
- `interfaceName` — the correct key is `interface`.
