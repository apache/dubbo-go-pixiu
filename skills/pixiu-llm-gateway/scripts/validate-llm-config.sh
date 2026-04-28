#!/usr/bin/env bash
#
# validate-llm-config.sh — structural validation for pixiu LLM gateway conf.yaml.
#
# Usage:
#   validate-llm-config.sh <conf.yaml>
#
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <conf.yaml>" >&2
  exit 1
fi

CONF="$1"
if [[ ! -f "$CONF" ]]; then
  echo "error: $CONF not found" >&2
  exit 1
fi

yaml_to_json() {
  local file="$1"
  if command -v yq >/dev/null 2>&1; then
    yq -o json eval '.' "$file"
  elif command -v ruby >/dev/null 2>&1; then
    ruby -ryaml -rjson -e 'puts JSON.generate(YAML.safe_load(File.read(ARGV.fetch(0)), permitted_classes: [], permitted_symbols: [], aliases: false))' "$file"
  else
    return 127
  fi
}

has_yaml_reader() {
  command -v yq >/dev/null 2>&1 || command -v ruby >/dev/null 2>&1
}

errors=0
tmp_files=()

cleanup() {
  if [[ ${#tmp_files[@]} -gt 0 ]]; then
    rm -f "${tmp_files[@]}"
  fi
}
trap cleanup EXIT

section() { echo; echo "== $* =="; }

section "1. YAML syntax"
tmpjson=$(mktemp -t pixiu-llm-XXXXXX.json)
tmp_files+=("$tmpjson")
yaml_err=$(mktemp -t pixiu-llm-yaml-XXXXXX)
tmp_files+=("$yaml_err")
if ! has_yaml_reader; then
  echo "  SKIP (install yq or ruby for YAML parsing)"
elif yaml_to_json "$CONF" > "$tmpjson" 2>"$yaml_err"; then
  echo "  OK"
else
  echo "  FAIL:"; cat "$yaml_err"
  exit 2
fi

section "2. LLM gateway semantics"
if ! has_yaml_reader; then
  echo "  SKIP (install yq or ruby for YAML parsing)"
elif ! command -v python3 >/dev/null 2>&1; then
  echo "  SKIP (install python3 for semantic validation)"
else
  semantic_errors=$(python3 - "$tmpjson" <<'PY'
import json
import sys

with open(sys.argv[1], "r", encoding="utf-8") as f:
    cfg = json.load(f) or {}

errors = []
warnings = []
sr = cfg.get("static_resources") or cfg
listeners = sr.get("listeners") or []
clusters = sr.get("clusters") or cfg.get("clusters") or []
cluster_names = {c.get("name") for c in clusters if isinstance(c, dict)}
endpoint_ids = {}
for cluster in clusters:
    if not isinstance(cluster, dict):
        continue
    cname = cluster.get("name")
    endpoint_ids[cname] = set()
    for ep in cluster.get("endpoints") or []:
        if not isinstance(ep, dict):
            continue
        if ep.get("id"):
            endpoint_ids[cname].add(str(ep.get("id")))
        domains = ((ep.get("socket_address") or {}).get("domains") or [])
        for domain in domains:
            if "://" in str(domain) or "/" in str(domain):
                errors.append(f"cluster {cname}: endpoint {ep.get('id')} socket_address.domains entries must be host-only, got {domain!r}")
        meta = ep.get("llm_meta")
        if meta is not None and not isinstance(meta, dict):
            errors.append(f"cluster {cname}: endpoint {ep.get('id')} llm_meta must be a map")
        if isinstance(meta, dict):
            rp = meta.get("retry_policy")
            if rp and not isinstance(rp, dict):
                errors.append(f"cluster {cname}: endpoint {ep.get('id')} retry_policy must be a map")

all_http_filters = []
routes = []

def network_filters(listener):
    if not isinstance(listener, dict):
        return []
    filter_chains = listener.get("filter_chains") or []
    if isinstance(filter_chains, dict):
        filter_chains = [filter_chains]
    elif not isinstance(filter_chains, list):
        return []

    filters = []
    for fc in filter_chains:
        if isinstance(fc, dict):
            filters.extend(fc.get("filters") or [])
    return filters

for listener in listeners:
    for chain in network_filters(listener):
        if not isinstance(chain, dict):
            continue
        if chain.get("name") != "dgp.filter.httpconnectionmanager":
            continue
        hcm = chain.get("config") or {}
        routes.extend(((hcm.get("route_config") or {}).get("routes") or []))
        all_http_filters.extend(hcm.get("http_filters") or [])

filter_names = [f.get("name") for f in all_http_filters if isinstance(f, dict)]
if "dgp.filter.llm.proxy" not in filter_names:
    errors.append("missing dgp.filter.llm.proxy in http_filters")
else:
    proxy_i = filter_names.index("dgp.filter.llm.proxy")
    for before in ("dgp.filter.ai.kvcache", "dgp.filter.llm.tokenizer"):
        if before in filter_names and filter_names.index(before) > proxy_i:
            errors.append(f"{before} must be before dgp.filter.llm.proxy")

for route in routes:
    cluster = ((route or {}).get("route") or {}).get("cluster")
    if cluster and cluster not in cluster_names:
        errors.append(f"route references missing cluster {cluster!r}")

for f in all_http_filters:
    if not isinstance(f, dict) or f.get("name") != "dgp.filter.ai.kvcache":
        continue
    c = f.get("config") or {}
    if c.get("enabled") is True:
        if not c.get("vllm_endpoint"):
            errors.append("kvcache enabled but vllm_endpoint is empty")
        if not c.get("lmcache_endpoint"):
            errors.append("kvcache enabled but lmcache_endpoint is empty")
    cs = c.get("cache_strategy") or {}
    for key in ("pin_instance_id", "compress_instance_id", "evict_instance_id"):
        val = cs.get(key)
        if val and not any(val in ids for ids in endpoint_ids.values()):
            errors.append(f"kvcache {key}={val!r} does not match any endpoint id")
    for key in ("memory_threshold", "load_threshold"):
        val = cs.get(key)
        if val is not None:
            try:
                n = float(val)
                if n < 0 or n > 1:
                    errors.append(f"cache_strategy.{key} must be between 0 and 1")
            except Exception:
                errors.append(f"cache_strategy.{key} must be numeric")

adapters = cfg.get("adapters") or sr.get("adapters") or []
for adapter in adapters:
    if not isinstance(adapter, dict) or adapter.get("name") != "dgp.adapter.llmregistrycenter":
        continue
    regs = ((adapter.get("config") or {}).get("registries") or {})
    if not regs:
        errors.append("llmregistrycenter adapter has no registries")
    for name, reg in regs.items():
        protocol = (reg or {}).get("protocol")
        if protocol != "nacos":
            warnings.append(f"llm registry {name} protocol is {protocol!r}; current examples use nacos")
        if not (reg or {}).get("address"):
            errors.append(f"llm registry {name} address is empty")

for warning in warnings:
    print(f"WARN: {warning}", file=sys.stderr)
print("\n".join(errors))
PY
)

  if [[ -n "$semantic_errors" ]]; then
    echo "  FAIL:"
    formatted="    - ${semantic_errors//$'\n'/$'\n    - '}"
    printf '%s\n' "$formatted"
    errors=$((errors+1))
  else
    echo "  OK"
  fi
fi

echo
if [[ $errors -eq 0 ]]; then
  echo "Validation passed."
  exit 0
else
  echo "Validation failed with $errors error group(s)."
  exit 2
fi
