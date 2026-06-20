#!/usr/bin/env bash
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License. You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Demonstrates MCP intelligent tool routing (issue #937) against a Pixiu
# gateway started with router.yaml.
#
# The example router.yaml uses an always-on policy rule (no `when`) that denies
# the `internal`/`admin` tags, so the demo works WITHOUT the auth filter. The
# additional claim-based rules in router.yaml only take effect when the MCP auth
# filter populates JWT claims.
#
# Usage: bash demo.sh [base_url]   (default http://localhost:8888)

set -euo pipefail
BASE="${1:-http://localhost:8888}"
MCP="$BASE/mcp"
SESSION=""
HEADERS=""

cleanup() {
  if [[ -n "$HEADERS" ]]; then
    rm -f "$HEADERS"
  fi
}
trap cleanup EXIT

post() {
  local payload="$1"

  curl -s -X POST "$MCP" \
    -H "Content-Type: application/json" \
    -H "Mcp-Session-Id: $SESSION" \
    -d "$payload"
}

echo "== 1. initialize =="
HEADERS="$(mktemp)"
curl -s -D "$HEADERS" -X POST "$MCP" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}'
echo
SESSION="$(awk 'tolower($1) == "mcp-session-id:" { gsub("\r", "", $2); print $2; exit }' "$HEADERS")"
cleanup
HEADERS=""
if [[ -z "$SESSION" ]]; then
  echo "initialize did not return Mcp-Session-Id" >&2
  exit 1
fi
echo "session acquired"
echo

echo "== 2. tools/list -> internal_dump (internal/admin tags) is filtered out =="
post '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
echo

echo "== 3. tools/call internal_dump -> DENIED (not in session plan) =="
post '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"internal_dump","arguments":{}}}'
echo

echo "== 4. tools/call search_kb -> ALLOWED =="
post '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"search_kb","arguments":{"q":"reset password"}}}'
echo
