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

set -e

readonly PIXIU_ROOT_DIR=$(pwd)
readonly SAMPLES_BRANCH="main"
readonly SAMPLES_REPO_URL="https://github.com/apache/dubbo-go-pixiu-samples.git"
readonly SAMPLES_CLONE_DIR="integrate_samples"

echo "::group::Integration Test Environment Details"
echo "Pixiu Root Directory:         ${PIXIU_ROOT_DIR}"
echo "Commit SHA:                   ${GITHUB_SHA}"
echo "Target Branch for Samples:    ${SAMPLES_BRANCH}"
echo "Repository Slug:              ${GITHUB_REPOSITORY}"
echo "::endgroup::"

if [ ! -d "$SAMPLES_CLONE_DIR" ]; then
  echo "> Cloning dubbo-go-samples (branch: ${SAMPLES_BRANCH})..."
  git clone --depth 1 -b "${SAMPLES_BRANCH}" "${SAMPLES_REPO_URL}" "${SAMPLES_CLONE_DIR}"
fi

cd "${SAMPLES_CLONE_DIR}"

echo "> Configuring Go modules to use local pixiu code..."
go mod edit -replace="github.com/apache/dubbo-go-pixiu=${PIXIU_ROOT_DIR}"

echo "> Preparing dependencies..."
go mod tidy

echo "> Handing off to the integration test runner..."
bash ./start_integrate_test.sh

echo "Integration tests completed successfully."