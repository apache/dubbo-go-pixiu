# xDS Operations Guide

Pixiu Admin publishes gateway listeners and clusters through Envoy xDS v3. It
supports the project-specific `TypedExtensionConfig` transport and standard
CDS/EDS ingestion for Istio-compatible control planes.

## Admin configuration

Configure the management server in the Admin YAML:

```yaml
xds:
  listen_port: 18000
  node_id: gateway-a
```

`listen_port` defaults to `18000`; `node_id` defaults to `test-id`. The Pixiu
client node ID must exactly match the Admin node ID because the snapshot cache
is keyed by that value.

Admin loads listeners and clusters, validates a complete candidate snapshot,
and publishes it only after all conversions succeed. Its etcd prefix watch
then rebuilds on create, update, and delete events. A failed candidate does not
replace the last-good snapshot or advance its version.

## Diagnostics

After obtaining an Admin JWT, query:

```shell
curl -H "token: ${TOKEN}" http://127.0.0.1:8081/config/api/xds/status
```

The `data` object contains:

| Field | Meaning |
| --- | --- |
| `node_id` | Snapshot cache key expected from clients. |
| `listen_port` | Effective xDS gRPC listen port. |
| `snapshot_version` | Last successfully published version. |
| `listener_count`, `cluster_count` | Resources in the last-good snapshot. |
| `last_updated_at` | Time the last-good snapshot was published. |
| `last_attempt_at` | Time of the latest successful or failed publication attempt. |
| `last_error`, `last_error_at` | Latest rejected candidate and its time. |
| `ready` | At least one snapshot has been published. |
| `degraded` | The latest candidate failed; the last-good snapshot is still served. |

An empty `snapshot_version` means no candidate has been published. A response
with both `ready: true` and `degraded: true` means clients can continue using
the last-good version while the operator fixes the newer configuration.

## Client correctness

- Delta ExtensionConfig responses are applied before ACK. Decode or runtime
  application failures produce a NACK with `ErrorDetail`.
- Empty Delta responses are no-ops. An explicitly delivered empty aggregate
  removes the corresponding xDS-owned resources.
- Reconnect requests include only versions that were successfully applied.
- Static listeners and clusters cannot be overwritten or deleted by xDS.
- Standard ADS uses one stream for CDS and referenced EDS resources. Standard
  LDS conversion is not currently supported.

## Verification

Run the focused tests with:

```shell
go test -race ./admin/... ./pkg/config/xds/... ./pkg/server
```

The Admin suite includes an in-memory gRPC test covering SnapshotBuilder,
SnapshotCache, and ExtensionConfig Discovery. The router suite exercises the
authenticated diagnostics endpoint.

When troubleshooting, check these in order:

1. `node_id` matches on Admin and Pixiu.
2. `ready` is true and `snapshot_version` changes after a valid etcd update.
3. `last_error` for conversion or publication failures.
4. Client logs for ACK/NACK errors or unsupported standard xDS resources.
