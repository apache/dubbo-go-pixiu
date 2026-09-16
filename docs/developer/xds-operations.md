# xDS Operations Guide

Pixiu Admin publishes gateway listeners and clusters through Envoy xDS v3. It
supports the project-specific `TypedExtensionConfig` transport and standard
CDS/EDS ingestion for Istio-compatible control planes.

## Resource support matrix

`Supported` means the repository maintains the complete Admin-to-Pixiu path.
`Experimental` means conversion and protocol handling exist, but production
environment coverage is still limited. Registering a gRPC service alone does
not make its resource type supported.

| Resource path | Status | Boundary |
| --- | --- | --- |
| ExtensionConfig listener | Supported | Pixiu-specific aggregate over Delta ECDS. |
| ExtensionConfig cluster | Supported | Pixiu-specific aggregate over Delta ECDS. |
| Standard CDS | Experimental | EDS-type clusters through standard ADS. |
| Standard EDS | Experimental | Socket-address endpoints for referenced CDS resources. |
| Standard LDS | Unsupported | Listener conversion is outside the base scope. |
| Standard RDS/SDS/RTDS | Unsupported | Services may be registered, but no resources are published or applied. |
| `ads_config` bootstrap | Unsupported | Use `lds_config`/`cds_config` with the documented API types. |

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
| `listening` | Whether the xDS gRPC socket is currently bound. |
| `listen_error`, `listen_error_at` | Latest listener startup or serving error and its time. |
| `snapshot_version` | Last successfully published version. |
| `listener_count`, `cluster_count` | Resources in the last-good snapshot. |
| `last_updated_at` | Time the last-good snapshot was published. |
| `last_attempt_at` | Time of the latest successful or failed publication attempt. |
| `last_error`, `last_error_at` | Latest rejected candidate and its time. |
| `ready` | The xDS listener is bound and at least one snapshot has been published. |
| `degraded` | Listener startup/serving or the latest snapshot candidate failed. |
| `resource_support` | Explicit support status for each xDS path. |

An empty `snapshot_version` means no candidate has been published. A response
with both `ready: true` and `degraded: true` means the listener is available and
clients can continue using the last-good version while the operator fixes the
newer configuration.

## Client correctness

- Delta ExtensionConfig responses are applied before ACK. Decode or runtime
  application failures produce a NACK with `ErrorDetail`.
- Empty Delta responses are no-ops. An explicitly delivered empty aggregate
  removes the corresponding xDS-owned resources.
- Reconnect requests include only versions that were successfully applied.
- Static listeners and clusters cannot be overwritten or deleted by xDS.
- Standard ADS uses one stream for CDS and referenced EDS resources. Standard
  LDS conversion is not currently supported. Non-EDS CDS resources are logged
  and skipped; invalid fields in a selected EDS cluster still produce a NACK.

## Compatibility notes

- Listener construction is fail-closed. A bad network or HTTP filter prevents
  a static listener from starting and causes an xDS listener update to be
  NACKed. This is stricter than the previous behavior that silently omitted a
  bad filter; the Pixiu log identifies every rejected static listener.
- An out-of-tree `ListenerService` still satisfies the public interface, but it
  must additionally implement the optional transactional refresh contract to
  receive xDS listener updates. Otherwise the update is NACKed and its
  last-good listener remains active.
- `FilterManager.CreateFilterChain` retains its legacy best-effort behavior for
  out-of-tree callers, and `CreateHttpConnectionManager` remains non-nil for a
  non-nil top-level config while logging bad nested filters. Pixiu's xDS
  publication path uses the checked builders and rejects incomplete chains.
- TCP readiness uses a bounded self-dial because dubbo-getty does not expose a
  synchronous bind result. Triple `v1.2.2-rc4` reports bind failure by panic;
  the listener converts that panic into a startup error.

## Verification

Run the focused tests with:

```shell
go test -race ./admin/... ./pkg/config/xds/... ./pkg/server
```

The Admin suite includes an in-memory gRPC test covering SnapshotBuilder,
SnapshotCache, and ExtensionConfig Discovery. The router suite exercises the
authenticated diagnostics endpoint.

The deployment-level test is intentionally separate from the component test.
It starts an embedded disposable etcd by default and crosses Admin HTTP CRUD,
etcd Watch, SnapshotCache, Delta gRPC, the real Pixiu managers, a bound
listener, HTTP traffic, and Pixiu's production HTTP-to-Dubbo filter chain:

```shell
go test ./pkg/config/xds -run '^TestAdminHTTPToRunningPixiuEndToEnd$' -count=1
```

Set `PIXIU_E2E_ETCD_ENDPOINT` only when an external disposable etcd should be
used. The default embedded server intentionally keeps
`go.etcd.io/etcd/server/v3 v3.5.7` in the module graph; it matches the existing
etcd client version and makes the full acceptance test runnable without
external infrastructure. The test uses a unique prefix and deletes it during
cleanup. It verifies create, update, and delete traffic; invalid resource NACK
with last-good service continuity; xDS disconnect and reconnect recovery; and a
real direct Dubbo provider call through `apiconfig` and `dubboproxy`. Triple
compatibility
is not claimed by this test. The production `dubboproxy` path takes its direct
provider URL from the API definition; it does not resolve that target through
`ClusterManager`. The test therefore verifies xDS cluster runtime changes and
direct Dubbo compatibility as separate assertions, including that a
cluster-only update does not change the direct Dubbo target.

When troubleshooting, check these in order:

1. `node_id` matches on Admin and Pixiu.
2. `ready` is true and `snapshot_version` changes after a valid etcd update.
3. `last_error` for conversion or publication failures.
4. Client logs for ACK/NACK errors or unsupported standard xDS resources.
