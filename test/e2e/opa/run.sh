#!/bin/bash
#
#  Licensed to the Apache Software Foundation (ASF) under one or more
#  contributor license agreements.  See the NOTICE file distributed with
#  this work for additional information regarding copyright ownership.
#  The ASF licenses this file to You under the Apache License, Version 2.0
#  (the "License"); you may not use this file except in compliance with
#  the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#

# PR2 OPA end-to-end test runner.
#
# This script is a thin convenience wrapper around `go test`. It runs the
# full-link suite that lives in admin/initialize/e2e_opa_test.go — that file
# stands up a smart in-process OPA mock (real github.com/open-policy-agent/opa
# rego library), drives the admin REST API to publish policies, then exercises
# the gateway OPA filter against the same mock to verify allow/deny decisions.
#
# Usage:
#     ./test/e2e/opa/run.sh            # run all PR2 E2E cases
#     ./test/e2e/opa/run.sh -run Allow # filter cases by name
#     VERBOSE=1 ./test/e2e/opa/run.sh  # add -v
#
# Requirements:
#     - Go toolchain (matches go.mod's version directive)
#     - No docker, etcd, mysql, or real OPA binary required.

set -euo pipefail

cd "$(dirname "$0")/../../.."

ARGS=("-count=1" "-run" "TestE2E_")
if [[ "${VERBOSE:-0}" == "1" ]]; then
    ARGS+=("-v")
fi

# Forward any extra args (e.g. -run regex override) after our defaults.
go test "${ARGS[@]}" "$@" ./admin/initialize/
