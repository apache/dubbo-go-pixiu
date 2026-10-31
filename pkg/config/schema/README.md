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

The route binding covers one HTTP entry mapped to one Dubbo method through the
registry set configured on the listener. Other Admin object kinds are outside
this issue's scope.

`Resource` and `Method` are generated Pixiu runtime structures. They are not
registered as Admin form objects.

## Object boundary

`AdminObject` keeps only a stable `kind`, `metadata`, and dynamic
`spec map[string]any`. `ObjectSchema` and `FieldSchema` define the form,
defaults, and validation rules for the dynamic values.

The built-in `AdminRouteBinding` schema has these separate concerns:

- `entry`: HTTP protocol, path, and method.
- `backendRef`: declares registry-backed discovery. Registry endpoints remain
  listener-level `dubboProxyConfig.registries` configuration and apply to all
  routes. A migrated legacy `target.cluster` is retained as
  `integrationRequest.clusterName`, not treated as a registry ID.
- `serviceSelection`: Dubbo application, interface, group, and version.
- `invocation`: Dubbo method, positive timeout, and ordered
  HTTP-source to Dubbo-argument mappings.
- `enabled`: whether the published route accepts requests; it defaults to `true`.

Objects created by PR #1046 with `target` and `params` are migrated during
normalization to the separated shape. Mixed old and new fields and unsupported
legacy fields are rejected instead of silently dropped. URI mappings must
match the `:parameter` segments in `entry.path`.

Draft and published lifecycle state is tracked by Admin separately; it is not a
user-facing `spec.publish` field.

See [`testdata/admin_route_binding.yaml`](testdata/admin_route_binding.yaml) for
the complete example and `CompiledRoute.PreviewYAML` for the legacy output.
