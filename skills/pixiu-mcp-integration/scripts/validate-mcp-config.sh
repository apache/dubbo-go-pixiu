#!/usr/bin/env bash
#
# validate-mcp-config.sh — structural validation for pixiu MCP gateway conf.yaml.
#
# Usage:
#   validate-mcp-config.sh <conf.yaml>
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
tmpjson=$(mktemp -t pixiu-mcp-XXXXXX.json)
tmp_files+=("$tmpjson")
yaml_err=$(mktemp -t pixiu-mcp-yaml-XXXXXX)
tmp_files+=("$yaml_err")
if ! has_yaml_reader; then
  echo "  SKIP (install yq or ruby for YAML parsing)"
  errors=$((errors+1))
elif yaml_to_json "$CONF" > "$tmpjson" 2>"$yaml_err"; then
  echo "  OK"
else
  echo "  FAIL:"; cat "$yaml_err"
  exit 2
fi

section "2. MCP gateway semantics"
if ! has_yaml_reader; then
  echo "  SKIP (install yq or ruby for YAML parsing; YAML syntax section already recorded this)"
elif ! command -v python3 >/dev/null 2>&1; then
  echo "  SKIP (install python3 for semantic validation)"
  errors=$((errors+1))
else
  semantic_errors=$(python3 - "$tmpjson" <<'PY'
import json
import re
import sys

with open(sys.argv[1], "r", encoding="utf-8") as f:
    cfg = json.load(f) or {}

errors = []
sr = cfg.get("static_resources") or cfg
listeners = sr.get("listeners") or []
clusters = sr.get("clusters") or cfg.get("clusters") or []
cluster_names = {c.get("name") for c in clusters if isinstance(c, dict)}

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
if "dgp.filter.mcp.mcpserver" not in filter_names:
    errors.append("missing dgp.filter.mcp.mcpserver in http_filters")
else:
    mcp_i = filter_names.index("dgp.filter.mcp.mcpserver")
    if "dgp.filter.http.auth.mcp" in filter_names and filter_names.index("dgp.filter.http.auth.mcp") > mcp_i:
        errors.append("dgp.filter.http.auth.mcp must be before dgp.filter.mcp.mcpserver")
    if "dgp.filter.http.httpproxy" in filter_names and filter_names.index("dgp.filter.http.httpproxy") < mcp_i:
        errors.append("dgp.filter.http.httpproxy should be after dgp.filter.mcp.mcpserver")

route_clusters = set()
for route in routes:
    cluster = ((route or {}).get("route") or {}).get("cluster")
    if cluster:
        route_clusters.add(cluster)
        if cluster not in cluster_names:
            errors.append(f"route references missing cluster {cluster!r}")

valid_arg_types = {"string", "integer", "number", "boolean"}
valid_arg_in = {"path", "query", "body"}
for f in all_http_filters:
    if not isinstance(f, dict) or f.get("name") != "dgp.filter.mcp.mcpserver":
        continue
    c = f.get("config") or {}
    endpoint = c.get("endpoint", "/mcp")
    if not str(endpoint).startswith("/"):
        errors.append("mcp endpoint must start with /")
    for tool in c.get("tools") or []:
        if not isinstance(tool, dict):
            continue
        name = tool.get("name", "<unnamed>")
        cluster = tool.get("cluster")
        if not cluster:
            errors.append(f"tool {name}: cluster is required")
        elif cluster not in cluster_names:
            errors.append(f"tool {name}: cluster {cluster!r} is not declared")
        req = tool.get("request") or {}
        path = req.get("path", "")
        if not req.get("method"):
            errors.append(f"tool {name}: request.method is required")
        if not path:
            errors.append(f"tool {name}: request.path is required")
        args = tool.get("args") or []
        arg_by_name = {a.get("name"): a for a in args if isinstance(a, dict)}
        for placeholder in re.findall(r"\{([^}]+)\}", str(path)):
            arg = arg_by_name.get(placeholder)
            if not arg or arg.get("in") != "path":
                errors.append(f"tool {name}: path placeholder {{{placeholder}}} needs matching arg with in: path")
        for arg in args:
            if not isinstance(arg, dict):
                continue
            if arg.get("type", "string") not in valid_arg_types:
                errors.append(f"tool {name}: arg {arg.get('name')} has invalid type {arg.get('type')!r}")
            if arg.get("in") not in valid_arg_in:
                errors.append(f"tool {name}: arg {arg.get('name')} has invalid in {arg.get('in')!r}")

for f in all_http_filters:
    if not isinstance(f, dict) or f.get("name") != "dgp.filter.http.auth.mcp":
        continue
    c = f.get("config") or {}
    rm = c.get("resource_metadata") or {}
    if not rm.get("resource"):
        errors.append("mcp auth resource_metadata.resource is required")
    if not rm.get("authorization_servers"):
        errors.append("mcp auth resource_metadata.authorization_servers is required")
    providers = c.get("providers") or []
    if not providers:
        errors.append("mcp auth providers must not be empty")
    for p in providers:
        if not isinstance(p, dict):
            continue
        for key in ("name", "issuer", "jwks"):
            if not p.get(key):
                errors.append(f"mcp auth provider missing {key}")
    for rule in c.get("rules") or []:
        cluster = (rule or {}).get("cluster")
        if not cluster:
            errors.append("mcp auth rule cluster is required")
        elif cluster not in route_clusters:
            errors.append(f"mcp auth rule cluster {cluster!r} does not match any route cluster")

adapters = cfg.get("adapters") or sr.get("adapters") or []
for adapter in adapters:
    if not isinstance(adapter, dict) or adapter.get("name") != "dgp.adapter.mcpserver":
        continue
    regs = ((adapter.get("config") or {}).get("registries") or {})
    if not regs:
        errors.append("mcpserver adapter has no registries")
    for name, reg in regs.items():
        if (reg or {}).get("protocol") != "nacos":
            errors.append(f"mcp registry {name} protocol should be nacos for current source")
        if not (reg or {}).get("address"):
            errors.append(f"mcp registry {name} address is empty")

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
