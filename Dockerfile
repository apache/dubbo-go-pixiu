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
FROM golang:alpine as builder
LABEL MAINTAINER="dubbogo.pixiu@outlook.com"

ENV GOPROXY="https://goproxy.cn,direct" \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
WORKDIR /app/
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w" -o dubbo-go-pixiu ./cmd/pixiu/*.go


### alpine
FROM amd64/ubuntu:latest
COPY --from=builder /app/dubbo-go-pixiu /app/
COPY --from=builder /app/docker-entrypoint.sh /app/
COPY --from=builder /app/configs/* /etc/pixiu/
WORKDIR /app/
RUN ["chmod", "+x","./docker-entrypoint.sh"]
ENTRYPOINT ["./docker-entrypoint.sh"]
