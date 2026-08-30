# Dynamic Admin configuration model

This package is the first code-level proposal for replacing raw Admin YAML
editing with a schema-driven model. It is intentionally not connected to the
current Admin controllers, etcd keys, or Pixiu runtime yet.

## Boundaries

The design has four explicit boundaries:

1. `ConfigObject` is a stable envelope. Its `spec` is `map[string]any`.
2. `ObjectSchema` and `FieldSchema` give those dynamic values type, default,
   validation, UI, and runtime-binding metadata.
3. `Registry` owns built-in and plugin field registration. Duplicate fields
   are rejected so an extension cannot override core semantics silently.
4. `RuntimeAdapter` names identify the future compiler from the dynamic model
   to the existing `pkg/config` and xDS models.

`ConfigSet` is the revision/snapshot boundary. It allows Resource/Method and
Listener/Cluster changes to be validated and published as one unit. Collection
validators check references across those objects, while object validators
handle rules local to one object.

## Dynamic field rules

- Named, unordered request parameters use a map.
- Ordered Dubbo method parameters use an array with an explicit index.
- Complex request bodies use an open JSON-Schema-like object.
- Backend alternatives use a discriminated union (`backend.kind`).
- Community fields are registered below `spec.extensions` unless they are
  promoted into a future core schema version.

Registering a schema makes a field editable, serializable, and validatable. It
does not make the field affect Pixiu automatically; an executable field also
needs a runtime adapter or plugin capability.

See [`testdata/config_set.yaml`](testdata/config_set.yaml) for a complete
Resource, Method, Listener, and Cluster snapshot.
