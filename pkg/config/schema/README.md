# Admin route binding model

This package is an initial code-level model for replacing raw Resource and
Method YAML editing with one Admin-facing route binding.

The supported flow is:

```text
AdminRouteBinding
    -> schema defaults and validation
    -> legacy Resource + Method compilation
    -> api_config.yaml preview
    -> Admin draft and publish storage
```

## Scope

The route binding covers one HTTP entry mapped to one Dubbo method, using
either the configured registry set or an explicit direct provider URL. Listener,
Cluster, PluginGroup, RateLimiter, and OPA remain separate resources and can be
connected through typed `extensions` fields.

`Resource` and `Method` are generated Pixiu runtime structures. They are not
registered as Admin form objects.

## Object boundary

`AdminObject` keeps only a stable `kind`, `metadata`, and dynamic
`spec map[string]any`. `ObjectSchema` and `FieldSchema` define the form,
defaults, and validation rules for the dynamic values.

The built-in `AdminRouteBinding` schema has these separate concerns:

- `entry`: HTTP protocol, path, and method.
- `backendRef`: registry-backed cluster reference or direct URL.
- `serviceSelection`: Dubbo application, interface, group, and version.
- `invocation`: Dubbo method, positive timeout, serialization, and ordered
  HTTP-source to Dubbo-argument mappings.
- `enabled`: whether the published route accepts requests; it defaults to `true`.

Objects created by PR #1046 with `target` and `params` are migrated during
normalization to the separated shape. Mixed old and new fields and unsupported
legacy fields are rejected instead of silently dropped. URI mappings must
match the `:parameter` segments in `entry.path`.

Draft and published lifecycle state is tracked by Admin separately; it is not a
user-facing `spec.publish` field.

Typed plugin fields can be inserted below `spec.extensions` through the same
registry without changing the stored Go object.

See [`testdata/admin_route_binding.yaml`](testdata/admin_route_binding.yaml) for
the complete example and `CompiledRoute.PreviewYAML` for the legacy output.
