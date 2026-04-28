# HTTP-to-Dubbo — Troubleshooting by Symptom

When the gateway returns a 5xx for a Dubbo-backed route, the failure is
almost always in one of four places. Walk them in order.

## Symptom: 404 "no route found"

Not a Dubbo issue — the request did not match any `resources[].path`.

Checks:

- Path matches exactly (watch for trailing slashes).
- `methods[].httpVerb` includes the method used.
- `methods[].enable: true` (or omitted if the current loader defaults
  it on).
- `dgp.filter.http.apiconfig` is in the `http_filters` list in
  `conf.yaml`. Without it, the api_config is not consulted.

## Symptom: 500 `no filter found for name dgp.filter.http.XXX`

Filter-registration issue. Apply the fix from `pixiu-filter-author`
(missing blank import in `pkg/pluginregistry/registry.go`).

## Symptom: 500 `no provider found`

The registry adapter booted but the Dubbo `interface` / `group` /
`version` combo is not registered. Check the three fields are spelled
correctly in `integrationRequest`. `group` is case-sensitive and
typically empty-string by default (`group: ""`) — the value `"default"`
is not interchangeable with `""`.

Also confirm the provider is registered to the *same* registry that the
pixiu adapter points at. A ZK-registered provider will not be found via
a Nacos adapter and vice versa.

## Symptom: 500 `generic invoke failed: ClassNotFoundException`

Pixiu or the request supplied a type string the provider cannot load.
In current static routes, the first place to check is
`mappingParams[].mapType`; for dynamic/default generic routes, check
request fields mapped to `opt.types` / `opt.values`. Old configs may
also contain top-level `paramTypes`, but current `IntegrationRequest`
does not bind that field.

Usual offenders:

- `mapType: String` instead of `mapType: string` or
  `mapType: java.lang.String`.
- `mapType: bool` instead of `mapType: boolean`.
- A dynamic `requestBody.types` entry containing `String` when the
  provider expects a loadable Java class name.
- POJO body missing provider-required `"class": "com.example.User"`
  metadata.
- User POJO with a misspelled package path or absent provider class.

See `generic-invoke-types.md` for the current `mapType` table.

## Symptom: 500 `generic invoke failed: ... NullPointerException`

The Dubbo method got called, but with `null` arguments. Root causes:

- `mappingParams` missing an entry with the expected `mapTo: "<i>"`.
- `name` selector references a field the request does not contain —
  e.g. `queryStrings.userId` but the request has `user_id`.
- `requestBody.<path>` dotted path does not match the actual JSON
  structure.
- Body `content-type` is not `application/json` and the body was not
  parsed.

Quick check: `curl -v` the request and confirm the body shape; then
walk `mappingParams` rule-by-rule.

## Symptom: 500 with `io.netty.handler.timeout.ReadTimeoutException`

Dubbo call is slower than the configured timeout. Two timeouts in play:

- `methods[].timeout` in `api_config.yaml` — the Dubbo RPC deadline.
- `listeners[].config.{read_timeout,write_timeout}` in `conf.yaml` —
  the HTTP deadline.

Bump `methods[].timeout` first. If you bump past the HTTP deadline, the
client will see a different error.

## Symptom: empty `{}` response, no error in logs

Dubbo call succeeded but pixiu could not serialize the result. Usually:

- The Dubbo method returns a POJO that uses Hessian serialization
  annotations the generic invoke serializer does not know about.
- The Dubbo method returns `void` and pixiu serialized `null` as
  `{}` — expected behavior, not a bug.

## Symptom: request hangs forever

Connection / routing issue, not a pixiu-config issue:

- `clusterName` in `integrationRequest` does not match any
  `static_resources.clusters[].name` in `conf.yaml`. This sometimes
  produces silent hangs depending on pixiu version.
- The Dubbo provider is reachable from the gateway host? (Different
  subnet, firewall, k8s NetworkPolicy.)
- The registry adapter is in `CONNECTING` state. Check pixiu logs for
  adapter warnings.

## Workflow for an unknown 500

1. `curl -v` the failing request; record the response.
2. `docker logs pixiu | grep -iE 'error|warn|fail' | tail -50` — look
   for the first error line. Pixiu logs tend to be verbose; the first
   error line near the request timestamp is the root.
3. Match the error substring against the sections above.
4. If no match, enable debug logging (`log.level: debug` in
   `conf.yaml`) and repeat.
5. Before editing yaml, run `bash scripts/validate-api-config.sh` — it
   catches structural mistakes the logs will not be specific about.
