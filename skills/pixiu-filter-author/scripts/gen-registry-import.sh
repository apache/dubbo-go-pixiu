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
    rel="${dir#"$prefix"}"
    WANT+=("$MOD/$rel")
  fi
done < <(find "$ROOT/pkg/filter" -type f -name '*.go' 2>/dev/null)

# De-dup and sort. Keep Bash 3 compatibility for macOS /bin/bash.
declare -a WANT_UNIQ=()
if [[ ${#WANT[@]} -gt 0 ]]; then
  while IFS= read -r line; do
    WANT_UNIQ+=("$line")
  done < <(printf '%s\n' "${WANT[@]}" | sort -u)
fi
WANT=("${WANT_UNIQ[@]}")

# Collect packages already blank-imported.
declare -a HAVE=()
while IFS= read -r line; do
  HAVE+=("$line")
done < <(
  grep -oE '^[[:space:]]*_[[:space:]]+"[^"]+"' "$REG" \
    | sed -E 's/^[[:space:]]*_[[:space:]]+"([^"]+)".*$/\1/' \
    | sort -u
)

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
