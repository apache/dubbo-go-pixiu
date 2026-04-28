# Parameter Mapping Rules

A Dubbo method call has an ordered argument list; an HTTP request has
inputs spread across path variables, query strings, headers, and body.
`mappingParams` is the translation layer.

## Anatomy of One Entry

```yaml
- name: requestBody.user.email    # source: where to pull from HTTP
  mapTo: "1"                      # target: Dubbo arg index or opt.* key
  mapType: string                 # type hint for Dubbo mapping
```

For current pixiu, every static Dubbo argument should be represented by
at least one `mappingParams` entry with numeric `mapTo`. There is no
current-source top-level `paramTypes` field on `IntegrationRequest`.

## `name` Source Grammar

| Prefix | Means | Example source |
|---|---|---|
| `queryStrings.<qs>` | URL query string | `?name=alice` -> `queryStrings.name` -> `"alice"` |
| `uri.<var>` | path variable | `/user/:id` matched by `/user/42` -> `uri.id` -> `"42"` |
| `headers.<h>` | HTTP request header | `X-Tenant: acme` -> `headers.X-Tenant` -> `"acme"` |
| `requestBody.<field>` | JSON body field | `{"user":{"name":"bob"}}` + `requestBody.user.name` -> `"bob"` |
| `requestBody` | whole parsed body | Use when the Dubbo arg is a POJO mirroring the body |

Rules:

- `queryStrings.*`, `uri.*`, and `headers.*` normally produce strings;
  use `mapType` to coerce them for numeric/boolean Dubbo arguments.
- `requestBody.*` preserves the JSON shape: objects become maps,
  arrays become slices, numbers remain JSON-decoded numbers.
- The source path is evaluated after Pixiu has parsed the request.

## `mapTo` Target Grammar

For static Dubbo arguments, `mapTo` is a zero-based argument index as a
string:

```yaml
mappingParams:
  - { name: uri.id,           mapTo: "0", mapType: string }
  - { name: queryStrings.n,   mapTo: "1", mapType: int }
  - { name: requestBody.user, mapTo: "2", mapType: object }
```

For dynamic generic routes, `mapTo` can target Dubbo option fields:

```yaml
mappingParams:
  - { name: requestBody.types,  mapTo: opt.types }
  - { name: requestBody.values, mapTo: opt.values }
  - { name: queryStrings.group, mapTo: opt.group }
```

Supported `opt.*` targets in current Pixiu include `opt.values`,
`opt.types`, `opt.group`, `opt.version`, `opt.interface`,
`opt.application`, and `opt.method`.

For HTTP pass-through routes, `mapTo` is a dotted target path into the
downstream request, such as `queryStrings.<name>`,
`requestBody.<field>`, or `headers.<name>`.

## `mapType` Values

Use only values accepted by `pkg/common/constant/jtypes.go`:

| `mapType` | Use for |
|---|---|
| `string` / `java.lang.String` | Java `String` |
| `char` | Java `char` / `Character` |
| `short` | Java `short` / `Short` |
| `int` | Java `int` / `Integer` |
| `long` | Java `long` / `Long` |
| `float` | Java `float` / `Float` |
| `double` | Java `double` / `Double` |
| `boolean` | Java `boolean` / `Boolean` |
| `java.util.Date` / `date` | Java date values |
| `object` / `java.lang.Object` | POJO, map-like object, nested JSON object |

Do not use `bool`, `array`, `String`, or arbitrary POJO FQCNs as
`mapType` values for current pixiu.

## Worked Examples

### Example 1: Primitives From Query String

Dubbo: `String greet(String name, int times, boolean loud)`

```yaml
mappingParams:
  - { name: queryStrings.name,  mapTo: "0", mapType: string }
  - { name: queryStrings.times, mapTo: "1", mapType: int }
  - { name: queryStrings.loud,  mapTo: "2", mapType: boolean }
```

### Example 2: One POJO From Body

Dubbo: `void createUser(User u)`

```yaml
mappingParams:
  - { name: requestBody, mapTo: "0", mapType: object }
```

HTTP body if the provider requires class metadata:

```json
{ "class": "com.example.User", "name": "alice", "age": 42 }
```

### Example 3: Mixed Sources

Dubbo: `Page<User> search(String tenant, Query q)`

```yaml
mappingParams:
  - { name: headers.X-Tenant,  mapTo: "0", mapType: string }
  - { name: requestBody.query, mapTo: "1", mapType: object }
```

`requestBody.query` means the HTTP body is shaped as
`{"query": {...}}`; the Dubbo method receives the inner object.

### Example 4: Dynamic Types and Values

Use this only when the HTTP client intentionally supplies generic
invoke payloads:

```yaml
mappingParams:
  - { name: requestBody.types,  mapTo: opt.types }
  - { name: requestBody.values, mapTo: opt.values }
```

The request body then owns correctness of the type strings. Static
routes should prefer numeric `mapTo` plus `mapType`.
