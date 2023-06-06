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
FROM amd64/ubuntu:latest as builder
LABEL MAINTAINER="dev@dubbo.apache.org"

RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc

ADD https://golang.google.cn/dl/go1.18.7.linux-amd64.tar.gz .
RUN tar -xf go1.18.7.linux-amd64.tar.gz -C /usr/local/
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app/
COPY . .

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -ldflags '-r ./lib -s -w' -tags="wasm" -trimpath -o dubbo-go-pixiu ./cmd/pixiu/*.go

### alpine
FROM amd64/pingcap/alpine-glibc:latest
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot
USER nonroot
COPY docker-entrypoint.sh /
COPY configs/* /etc/pixiu/
COPY --from=builder /root/go/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-amd64/libwasmer.so /lib/
COPY --from=builder /app/dubbo-go-pixiu /
RUN ["chmod", "+x", "/docker-entrypoint.sh"]
ENTRYPOINT ["sh","./docker-entrypoint.sh"]
