#!/usr/bin/env bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

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
    ruby -ryaml -rjson -e 'begin; puts JSON.generate(YAML.safe_load(File.read(ARGV.fetch(0)), permitted_classes: [], permitted_symbols: [], aliases: false)); rescue StandardError => e; warn "#{e.class}: #{e.message}"; exit 1; end' "$file"
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

if ! has_yaml_reader; then
  echo "error: install yq (preferred) or ruby for YAML parsing." >&2
  echo "       macOS: brew install yq" >&2
  echo "       Debian/Ubuntu: apt install yq" >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "error: python3 is required for MCP semantic validation." >&2
  echo "       macOS: brew install python" >&2
  echo "       Debian/Ubuntu: apt install python3" >&2
  exit 1
fi

section "1. YAML syntax"
tmpjson=$(mktemp -t pixiu-mcp-XXXXXX.json)
tmp_files+=("$tmpjson")
yaml_err=$(mktemp -t pixiu-mcp-yaml-XXXXXX)
tmp_files+=("$yaml_err")
if yaml_to_json "$CONF" > "$tmpjson" 2>"$yaml_err"; then
  echo "  OK"
else
  echo "  FAIL:"; cat "$yaml_err"
  exit 2
fi

section "2. MCP gateway semantics"
semantic_errors=$(python3 - "$tmpjson" <<'PY'
import json
import re
import sys

with open(sys.argv[1], "r", encoding="utf-8") as f:
    raw_cfg = json.load(f)

errors = []

def dict_field(container, key, label):
    if not isinstance(container, dict):
        return {}
    value = container.get(key)
    if value is None:
        return {}
    if isinstance(value, dict):
        return value
    errors.append(f"{label} must be a map/object")
    return {}

def list_field(container, key, label):
    if not isinstance(container, dict):
        return []
    value = container.get(key)
    if value is None:
        return []
    if isinstance(value, list):
        return value
    errors.append(f"{label} must be a list")
    return []

if not isinstance(raw_cfg, dict):
    errors.append("top-level config must be a map/object")
    cfg = {}
else:
    cfg = raw_cfg

sr_value = cfg.get("static_resources")
if sr_value is None:
    sr = cfg
elif isinstance(sr_value, dict):
    sr = sr_value
else:
    errors.append("static_resources must be a map/object")
    sr = {}

listeners = list_field(sr, "listeners", "static_resources.listeners")
clusters = list_field(sr, "clusters", "static_resources.clusters")
if not clusters and sr is not cfg:
    clusters = list_field(cfg, "clusters", "clusters")
cluster_names = {c.get("name") for c in clusters if isinstance(c, dict) and c.get("name")}

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
        hcm = dict_field(chain, "config", "dgp.filter.httpconnectionmanager.config")
        routes.extend(list_field(dict_field(hcm, "route_config", "route_config"), "routes", "route_config.routes"))
        all_http_filters.extend(list_field(hcm, "http_filters", "http_filters"))

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
    cluster = dict_field(route, "route", "route.route").get("cluster")
    if cluster:
        route_clusters.add(cluster)
        if cluster not in cluster_names:
            errors.append(f"route references missing cluster {cluster!r}")

valid_arg_types = {"string", "integer", "number", "boolean"}
valid_arg_in = {"path", "query", "body"}
for f in all_http_filters:
    if not isinstance(f, dict) or f.get("name") != "dgp.filter.mcp.mcpserver":
        continue
    c = dict_field(f, "config", "dgp.filter.mcp.mcpserver.config")
    endpoint = c.get("endpoint", "/mcp")
    if not str(endpoint).startswith("/"):
        errors.append("mcp endpoint must start with /")
    for tool in list_field(c, "tools", "mcpserver tools"):
        if not isinstance(tool, dict):
            continue
        name = tool.get("name", "<unnamed>")
        cluster = tool.get("cluster")
        if not cluster:
            errors.append(f"tool {name}: cluster is required")
        elif cluster not in cluster_names:
            errors.append(f"tool {name}: cluster {cluster!r} is not declared")
        req = dict_field(tool, "request", f"tool {name}.request")
        path = req.get("path", "")
        if not req.get("method"):
            errors.append(f"tool {name}: request.method is required")
        if not path:
            errors.append(f"tool {name}: request.path is required")
        args = list_field(tool, "args", f"tool {name}.args")
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
    c = dict_field(f, "config", "dgp.filter.http.auth.mcp.config")
    rm = dict_field(c, "resource_metadata", "mcp auth resource_metadata")
    if not rm.get("resource"):
        errors.append("mcp auth resource_metadata.resource is required")
    if not rm.get("authorization_servers"):
        errors.append("mcp auth resource_metadata.authorization_servers is required")
    providers = list_field(c, "providers", "mcp auth providers")
    if not providers:
        errors.append("mcp auth providers must not be empty")
    for p in providers:
        if not isinstance(p, dict):
            continue
        for key in ("name", "issuer", "jwks"):
            if not p.get(key):
                errors.append(f"mcp auth provider missing {key}")
    for rule in list_field(c, "rules", "mcp auth rules"):
        rule = dict_field({"rule": rule}, "rule", "mcp auth rule")
        cluster = rule.get("cluster")
        if not cluster:
            errors.append("mcp auth rule cluster is required")
        elif cluster not in route_clusters:
            errors.append(f"mcp auth rule cluster {cluster!r} does not match any route cluster")

if cfg.get("adapters") is not None:
    adapters = list_field(cfg, "adapters", "adapters")
else:
    adapters = list_field(sr, "adapters", "static_resources.adapters")
for adapter in adapters:
    if not isinstance(adapter, dict) or adapter.get("name") != "dgp.adapter.mcpserver":
        continue
    adapter_config = dict_field(adapter, "config", "mcpserver adapter config")
    regs_value = adapter_config.get("registries")
    if regs_value is None:
        errors.append("mcpserver adapter has no registries")
        continue
    if not isinstance(regs_value, dict):
        errors.append("mcpserver adapter registries must be a map/object")
        continue
    if not regs_value:
        errors.append("mcpserver adapter has no registries")
    for name, reg in regs_value.items():
        if not isinstance(reg, dict):
            errors.append(f"mcp registry {name} must be a map/object")
            continue
        if reg.get("protocol") != "nacos":
            errors.append(f"mcp registry {name} protocol should be nacos for current source")
        if not reg.get("address"):
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

echo
if [[ $errors -eq 0 ]]; then
  echo "Validation passed."
  exit 0
else
  echo "Validation failed with $errors error group(s)." >&2
  exit 2
fi
