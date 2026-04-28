#!/usr/bin/env bash
#
# validate-api-config.sh — structural validation for pixiu's api_config.yaml.
#
# Usage:
#   validate-api-config.sh <api_config.yaml> [conf.yaml]
#
# Uses:
#   - yq (https://github.com/mikefarah/yq) for yaml parsing
#   - ajv-cli (npm install -g ajv-cli) for JSON Schema validation
#
# If either is missing, the script reports what to install and falls
# back to the checks it can still perform.
#
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCHEMA="$HERE/../references/api-config-schema.json"

yaml_to_json() {
  local file="$1"
  if command -v yq >/dev/null 2>&1; then
    yq -o json eval '.' "$file"
  elif command -v ruby >/dev/null 2>&1; then
    ruby -ryaml -rjson -e 'puts JSON.generate(YAML.load_file(ARGV.fetch(0)))' "$file"
  else
    return 127
  fi
}

has_yaml_reader() {
  command -v yq >/dev/null 2>&1 || command -v ruby >/dev/null 2>&1
}

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <api_config.yaml> [conf.yaml]" >&2
  exit 1
fi

API_CFG="$1"
CONF_CFG="${2:-}"

if [[ ! -f "$API_CFG" ]]; then
  echo "error: $API_CFG not found" >&2
  exit 1
fi
if [[ ! -f "$SCHEMA" ]]; then
  echo "error: schema $SCHEMA not found; is the skill directory intact?" >&2
  exit 1
fi

errors=0

section() { echo; echo "== $* =="; }

section "1. YAML syntax"
if has_yaml_reader; then
  if yaml_to_json "$API_CFG" > /dev/null 2>/tmp/yaml.err; then
    echo "  OK"
  else
    echo "  FAIL:"; cat /tmp/yaml.err
    errors=$((errors+1))
  fi
else
  echo "  SKIP (install yq or ruby)"
fi

section "2. JSON Schema conformance"
if command -v ajv >/dev/null 2>&1 && has_yaml_reader; then
  tmpjson=$(mktemp -t pixiu-api-XXXXXX.json)
  yaml_to_json "$API_CFG" > "$tmpjson"
  if ajv validate --spec=draft7 -s "$SCHEMA" -d "$tmpjson" > /tmp/ajv.out 2>&1; then
    echo "  OK"
  else
    echo "  FAIL:"; cat /tmp/ajv.out
    errors=$((errors+1))
  fi
  rm -f "$tmpjson"
else
  echo "  SKIP (install: npm i -g ajv-cli; yq or ruby for yaml parsing)"
fi

section "3. Legacy top-level paramTypes"
if has_yaml_reader && command -v python3 >/dev/null 2>&1; then
  tmpjson=$(mktemp -t pixiu-api-XXXXXX.json)
  yaml_to_json "$API_CFG" > "$tmpjson"
  legacy=$(python3 - "$tmpjson" <<'PY'
import json
import sys

path = sys.argv[1]
with open(path, "r", encoding="utf-8") as f:
    data = json.load(f)

hits = []

def walk(node, trail):
    if isinstance(node, dict):
        if "paramTypes" in node:
            hits.append(".".join(trail + ["paramTypes"]))
        for key, value in node.items():
            walk(value, trail + [str(key)])
    elif isinstance(node, list):
        for i, value in enumerate(node):
            walk(value, trail + [str(i)])

walk(data, [])
print("\n".join(hits))
PY
)
  rm -f "$tmpjson"
  if [[ -n "$legacy" ]]; then
    echo "  FAIL: current pixiu IntegrationRequest does not bind top-level paramTypes:"
    echo "$legacy" | sed 's/^/    - /'
    echo "  (use mappingParams[].mapType for static args, or opt.types for dynamic generic routes)"
    errors=$((errors+1))
  else
    echo "  OK"
  fi
else
  echo "  SKIP (install python3 plus yq or ruby)"
fi

section "4. mappingParams mapTo/mapType sanity"
if has_yaml_reader && command -v python3 >/dev/null 2>&1; then
  tmpjson=$(mktemp -t pixiu-api-XXXXXX.json)
  yaml_to_json "$API_CFG" > "$tmpjson"
  map_errors=$(python3 - "$tmpjson" <<'PY'
import json
import re
import sys

path = sys.argv[1]
with open(path, "r", encoding="utf-8") as f:
    data = json.load(f)

allowed_map_types = {
    "string",
    "java.lang.String",
    "char",
    "short",
    "int",
    "long",
    "float",
    "double",
    "boolean",
    "java.util.Date",
    "date",
    "object",
    "java.lang.Object",
}
allowed_opt = {
    "opt.values",
    "opt.types",
    "opt.group",
    "opt.version",
    "opt.interface",
    "opt.application",
    "opt.method",
}
index_re = re.compile(r"^\d+$")
errors = []

def check_mapping(mapping, trail):
    map_to = str(mapping.get("mapTo", ""))
    if not (index_re.match(map_to) or map_to in allowed_opt):
        errors.append(f"{trail}.mapTo={map_to!r} is not a numeric index or supported opt.* target")
    map_type = mapping.get("mapType")
    if map_type not in (None, "") and str(map_type) not in allowed_map_types:
        errors.append(f"{trail}.mapType={map_type!r} is not in current JTypeMapper")

def walk(node, trail):
    if isinstance(node, dict):
        mappings = node.get("mappingParams")
        if isinstance(mappings, list):
            for i, mapping in enumerate(mappings):
                if isinstance(mapping, dict):
                    check_mapping(mapping, ".".join(trail + ["mappingParams", str(i)]))
        for key, value in node.items():
            walk(value, trail + [str(key)])
    elif isinstance(node, list):
        for i, value in enumerate(node):
            walk(value, trail + [str(i)])

walk(data, [])
print("\n".join(errors))
PY
)
  rm -f "$tmpjson"
  if [[ -n "$map_errors" ]]; then
    echo "  FAIL:"
    echo "$map_errors" | sed 's/^/    - /'
    errors=$((errors+1))
  else
    echo "  OK"
  fi
else
  echo "  SKIP (install python3 plus yq or ruby)"
fi

section "5. clusterName cross-check against conf.yaml"
if [[ -n "$CONF_CFG" ]]; then
  if [[ ! -f "$CONF_CFG" ]]; then
    echo "  SKIP (conf.yaml path $CONF_CFG not found)"
  elif ! has_yaml_reader || ! command -v python3 >/dev/null 2>&1; then
    echo "  SKIP (install python3 plus yq or ruby)"
  else
    api_json=$(mktemp -t pixiu-api-XXXXXX.json)
    conf_json=$(mktemp -t pixiu-conf-XXXXXX.json)
    yaml_to_json "$API_CFG" > "$api_json"
    yaml_to_json "$CONF_CFG" > "$conf_json"
    missing=$(python3 - "$api_json" "$conf_json" <<'PY'
import json
import sys

with open(sys.argv[1], "r", encoding="utf-8") as f:
    api = json.load(f)
with open(sys.argv[2], "r", encoding="utf-8") as f:
    conf = json.load(f)

refs = set()

def walk(node):
    if isinstance(node, dict):
        value = node.get("clusterName")
        if value:
            refs.add(str(value))
        for child in node.values():
            walk(child)
    elif isinstance(node, list):
        for child in node:
            walk(child)

walk(api)

clusters = (((conf or {}).get("static_resources") or {}).get("clusters") or [])
declared = {str(c.get("name")) for c in clusters if isinstance(c, dict) and c.get("name")}
print(" ".join(sorted(refs - declared)))
PY
)
    rm -f "$api_json" "$conf_json"
    if [[ -n "$missing" ]]; then
      echo "  FAIL: clusterName(s) referenced but not declared in conf.yaml:$missing"
      errors=$((errors+1))
    else
      echo "  OK"
    fi
  fi
else
  echo "  SKIP (pass conf.yaml as 2nd arg to enable)"
fi

echo
if [[ $errors -eq 0 ]]; then
  echo "Validation passed."
  exit 0
else
  echo "Validation failed with $errors error group(s)."
  exit 2
fi
