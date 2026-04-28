#!/usr/bin/env bash
#
# gen-registry-import.sh — find filter packages that are not blank-imported
# into pkg/pluginregistry/registry.go.
#
# Usage: gen-registry-import.sh
#
# Run from the pixiu repo root, or pass PIXIU_ROOT=<path>.
# The script never writes; it only prints the lines you should add and
# exits non-zero when missing entries are found (useful for pre-commit).
#
set -euo pipefail

ROOT="${PIXIU_ROOT:-$(pwd)}"
REG="$ROOT/pkg/pluginregistry/registry.go"

if [[ ! -f "$REG" ]]; then
  echo "error: $REG not found; set PIXIU_ROOT=<pixiu-repo-path>" >&2
  exit 1
fi

MOD="github.com/apache/dubbo-go-pixiu"

# Collect packages that register a filter.
declare -a WANT=()
while IFS= read -r gofile; do
  if grep -q -E 'filter\.RegisterHttpFilter|filter\.RegisterNetworkFilterPlugin' "$gofile" 2>/dev/null; then
    dir=$(dirname "$gofile")
    prefix="${ROOT}/"
    rel="${dir#$prefix}"
    WANT+=("$MOD/$rel")
  fi
done < <(find "$ROOT/pkg/filter" -type f -name '*.go' 2>/dev/null)

# De-dup and sort.
mapfile -t WANT < <(printf '%s\n' "${WANT[@]}" | sort -u)

# Collect packages already blank-imported.
declare -a HAVE=()
while IFS= read -r line; do
  HAVE+=("$line")
done < <(grep -oE '"[^"]+"' "$REG" | tr -d '"' | sort -u)

missing=()
for p in "${WANT[@]}"; do
  found=0
  for h in "${HAVE[@]}"; do
    if [[ "$p" == "$h" ]]; then found=1; break; fi
  done
  if [[ $found -eq 0 ]]; then
    missing+=("$p")
  fi
done

if [[ ${#missing[@]} -eq 0 ]]; then
  echo "OK: pluginregistry/registry.go covers every filter package under pkg/filter/."
  exit 0
fi

echo "Missing blank imports in pkg/pluginregistry/registry.go:"
echo ""
for p in "${missing[@]}"; do
  printf '\t_ "%s"\n' "$p"
done
echo ""
echo "Add them in alphabetical order within the appropriate group, then rerun." >&2
exit 2
