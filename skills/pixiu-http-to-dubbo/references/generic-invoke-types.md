# Dubbo Generic Invoke Type Hints

Pixiu builds a Dubbo generic invocation from `mappingParams`. In the
current source, the static type hint for each positional argument comes
from `mappingParams[].mapType`, not from a top-level `paramTypes` field
under `integrationRequest`.

Before relying on any example, check:

- `pkg/config/api_config.go` for the fields yaml actually binds.
- `pkg/client/dubbo/mapper.go` for `mapTo` and `mapType` behavior.
- `pkg/common/constant/jtypes.go` for accepted type strings.

## Supported `mapType` Values

Current `constant.JTypeMapper` accepts these keys:

| Java intent | `mapType` to use |
|---|---|
| `String` | `string` or `java.lang.String` |
| `char` / `Character` | `char` |
| `short` / `Short` | `short` |
| `int` / `Integer` | `int` |
| `long` / `Long` | `long` |
| `float` / `Float` | `float` |
| `double` / `Double` | `double` |
| `boolean` / `Boolean` | `boolean` |
| `java.util.Date` | `java.util.Date` or `date` |
| POJO / map-like object | `object` or `java.lang.Object` |

Do not use `String`, `bool`, `array`, or arbitrary POJO FQCNs as
`mapType` values unless Step 0 proves the target branch added them.

## POJOs

For a Java signature such as:

```java
User createUser(com.example.User user)
```

use a positional object mapping:

```yaml
mappingParams:
  - name: requestBody
    mapTo: "0"
    mapType: object
```

Some provider-side serializers need class metadata in the JSON object:

```json
{ "class": "com.example.User", "id": 42, "name": "alice" }
```

Pixiu can pass the object as a generic map, but it cannot infer every
provider-side FQCN from a top-level yaml field that the current config
struct ignores. When in doubt, include the provider class discriminator
in the body and verify against the real provider.

## Multiple Arguments

Keep numeric `mapTo` values aligned with the Java method parameter
order:

```java
Page<Order> search(String tenant, com.example.OrderQuery query)
```

```yaml
mappingParams:
  - name: headers.X-Tenant
    mapTo: "0"
    mapType: string
  - name: requestBody
    mapTo: "1"
    mapType: object
```

## Dynamic Generic Routes

Pixiu also has a dynamic/default generic path where the HTTP request
supplies values and type strings. Use this only when the client is
expected to send that shape explicitly:

```yaml
mappingParams:
  - name: requestBody.types
    mapTo: opt.types
  - name: requestBody.values
    mapTo: opt.values
```

In this mode, type strings come from the HTTP request body and must be
valid for the provider. A `ClassNotFoundException: String` usually means
the request supplied `String` where the provider expected a class name
or where the current static route should have used `mapType: string`
instead.

## Nested Objects and Enums

Nested objects should carry their own class metadata when the provider's
serializer expects it:

```json
{
  "class": "com.example.order.Order",
  "buyer": { "class": "com.example.user.User", "id": 1 }
}
```

Enums are usually represented as objects or strings according to the
provider contract. Preserve the provider's expected JSON/generic shape;
`mapType: object` on the top-level argument is enough for Pixiu to pass
the nested map through.

## Quick Checklist

- [ ] No top-level `integrationRequest.paramTypes` for current pixiu.
- [ ] Every numeric `mapTo` matches the Java method argument position.
- [ ] Every `mapType` is in `constant.JTypeMapper`.
- [ ] POJO bodies include `class` metadata if the provider needs it.
- [ ] Dynamic `opt.types` is used only when the client supplies the type
  list at request time.
