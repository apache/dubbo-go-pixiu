# Admin configuration models

This package defines Admin-facing models that compile to existing Pixiu
runtime configuration.

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
registry set configured on the listener. `AdminCluster` covers the current
runtime Cluster's name, discovery type, address, and port.

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

## AdminCluster boundary

```yaml
kind: AdminCluster
metadata:
  name: user-provider
spec:
  type: Static
  endpoint:
    address: 127.0.0.1
    port: 20880
```

`Static` is the form default. The supported types are `Static`, `StrictDNS`,
`LogicalDns`, `EDS`, and `OriginalDst`, matching the existing runtime discovery
types. Strict validation requires a name, nonempty address, and port in
`1-65535`. Errors include paths such as `metadata.name`, `spec.type`,
`spec.endpoint.address`, and `spec.endpoint.port`.

`DecodeLegacyClusterYAML` reads the existing flat YAML into an `AdminObject`
and returns the runtime ID separately. `CompileAdminCluster` validates and
produces `config.Cluster`; its `PreviewYAML` renders flat runtime fields without
an ID. The publishing store must attach the existing or newly allocated ID and
must check name uniqueness against published runtime Clusters. Draft storage
must use a stable ID independent of `metadata.name` so incomplete drafts can
have an empty name. Existing Cluster CRUD requests still write runtime data;
the lifecycle store must reconcile their effects and guard the legacy delete
path. These storage and API rules are implemented in later PRs.
