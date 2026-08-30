# Admin route binding model

This package is an initial code-level model for replacing raw Resource and
Method YAML editing with one Admin-facing route binding.

The supported flow is:

```text
AdminRouteBinding
    -> schema defaults and validation
    -> legacy Resource + Method compilation
    -> api_config.yaml preview
    -> future draft and publish storage
```

## Scope

The first version deliberately covers one HTTP entry mapped to one
registry-backed Dubbo method. Listener, Cluster, PluginGroup, direct provider
URLs, and the actual etcd publish transaction are outside this package.

`Resource` and `Method` are generated Pixiu runtime structures. They are not
registered as Admin form objects.

## Object boundary

`AdminObject` keeps only a stable `kind`, `metadata`, and dynamic
`spec map[string]any`. `ObjectSchema` and `FieldSchema` define the form,
defaults, and validation rules for the dynamic values.

The built-in `AdminRouteBinding` schema has four user-facing sections:

- `entry`: HTTP protocol, path, and method.
- `target`: Dubbo application, interface, method, version, group, and cluster.
- `params`: ordered HTTP-source to Dubbo-argument mappings.
- `publish`: Admin control-plane intent; it is never emitted into Pixiu YAML.

Typed plugin fields can be inserted below `spec.extensions` through the same
registry without changing the stored Go object.

See [`testdata/admin_route_binding.yaml`](testdata/admin_route_binding.yaml) for
the complete example and `CompiledRoute.PreviewYAML` for the legacy output.
