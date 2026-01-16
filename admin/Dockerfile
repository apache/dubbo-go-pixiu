#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# ==================== Build Stage: Backend ====================
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pixiu-admin ./cmd/admin/admin.go

# ==================== Build Stage: Frontend ====================
FROM node:22-alpine AS frontend-builder

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package files
COPY admin/web/package.json admin/web/pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY admin/web/ ./

# Build
RUN pnpm build

# ==================== Runtime Stage ====================
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from backend builder
COPY --from=backend-builder /app/pixiu-admin /app/pixiu-admin

# Copy frontend dist from frontend builder
COPY --from=frontend-builder /app/dist /app/web/dist

# Copy default config
COPY configs/admin_config.yaml /etc/pixiu/admin_config.yaml

# Create non-root user
RUN adduser -D -u 1000 pixiu && \
    chown -R pixiu:pixiu /app /etc/pixiu

USER pixiu

# Expose ports
# 8081: Admin API server
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/api/health || exit 1

# Default command
ENTRYPOINT ["/app/pixiu-admin"]
CMD ["-c", "/etc/pixiu/admin_config.yaml"]
