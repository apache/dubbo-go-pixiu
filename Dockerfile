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

### builder
FROM golang:1.23.4-bullseye AS builder
LABEL MAINTAINER="dev@dubbo.apache.org"

RUN apt-get update && apt-get install -y --no-install-recommends gcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Here I still remains "wasmer" tag, because we need to build the wasmer plugin
# if wasm feature is removed in the future, this tag and following wasm related command can be removed
RUN go build -ldflags '-r ./lib -s -w' -tags="wasmer" -trimpath -o /app/dubbo-go-pixiu ./cmd/pixiu/*.go

RUN find /go/pkg/mod -name "libwasmer.so" -exec cp {} /app/libwasmer.so \;


FROM pingcap/alpine-glibc:alpine-3.14.6
LABEL MAINTAINER="dev@dubbo.apache.org"

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

RUN mkdir -p /etc/pixiu

WORKDIR /app
COPY docker-entrypoint.sh /app
COPY configs /etc/pixiu/

COPY --from=builder /app/dubbo-go-pixiu .
COPY --from=builder /app/libwasmer.so /lib/

RUN chown -R nonroot:nonroot /app /etc/pixiu \
    && chmod +x /app/docker-entrypoint.sh

USER nonroot

ENTRYPOINT ["/app/docker-entrypoint.sh"]