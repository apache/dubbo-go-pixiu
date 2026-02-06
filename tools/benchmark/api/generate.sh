#!/bin/bash
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

set -e

export GO111MODULE="on"

# Install required protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/dubbogo/protoc-gen-go-triple/v3@latest

# Generate Triple stubs from benchmark.proto (for Triple protocol)
protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-triple_out=. --go-triple_opt=paths=source_relative \
    benchmark.proto

echo "Generated Triple stubs from benchmark.proto"

# Generate gRPC stubs in grpcstub/ directory (for standard gRPC protocol)
# Using the same proto file but with different go_package output
protoc \
    --go_out=./grpcstub --go_opt=paths=source_relative \
    --go_opt=Mbenchmark.proto=github.com/apache/dubbo-go-pixiu/tools/benchmark/api/grpcstub \
    --go-grpc_out=./grpcstub --go-grpc_opt=paths=source_relative \
    --go-grpc_opt=Mbenchmark.proto=github.com/apache/dubbo-go-pixiu/tools/benchmark/api/grpcstub \
    benchmark.proto

echo "Generated gRPC stubs in grpcstub/"
