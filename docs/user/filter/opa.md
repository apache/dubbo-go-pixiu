# OPA Filter (dgp.filter.http.opa)

English | [中文](opa_CN.md)

---

## English

### Overview
The `dgp.filter.http.opa` filter delegates authorization decisions to Open Policy Agent (OPA) via a Rego policy. This filter evaluates requests and determines whether to allow or deny based on the policy defined in Rego.

The filter supports two operation modes:
- **Server Mode (Recommended)**: Evaluates policies by calling an external OPA server via HTTP API, supporting centralized policy management and hot updates.
- **Embedded Mode**: Policies are provided as inline Rego modules and evaluated using OPA's built-in query engine.

### Operation Modes

#### Server Mode (Recommended)
- Calls an external OPA server via REST API (`server_url`).
- Supports centralized policy management and dynamic updates without restarting Pixiu.
- Suitable for production environments and large-scale deployments.
- Simple configuration - just specify OPA server address and decision path.

#### Embedded Mode (Backward Compatible)
- Loads a Rego **module string** from `config.policy`.
- Builds a Rego **query** from `config.entrypoint`.
- Policies are evaluated within the Pixiu process; updates require restart.
- Suitable for simple scenarios or test environments.

> **Automatic Mode Selection**:
> - If `server_url` is configured, **Server Mode** is used (priority).
> - If `server_url` is not configured but `policy` is provided, **Embedded Mode** is used.

### Configuration schema

#### Server Mode Configuration (Recommended)

```yaml
filters:
  - name: dgp.filter.httpconnectionmanager
    config:
      route_config:
        # ... your routes
      http_filters:
        - name: dgp.filter.http.opa
          config:
            # OPA server address
            server_url: "http://opa-server:8181"
            # Decision path (OPA REST API path)
            decision_path: "/v1/data/http/authz/allow"
            # Request timeout in milliseconds, default 100ms
            timeout_ms: 100
            # Optional: Bearer token authentication
            # bearer_token: "your-secret-token"
        # HTTP proxy filter should be after OPA filter
        - name: dgp.filter.http.proxy
          config:
            # ... proxy config
```

#### Embedded Mode Configuration (Backward Compatible)

```yaml
filters:
  - name: dgp.filter.httpconnectionmanager
    config:
      route_config:
        # ... your routes
      http_filters:
        - name: dgp.filter.http.opa
          config:
            policy: |
              package http.authz

              default allow = false

              allow {
                input.method == "GET"
                input.path == "/status"
              }
            entrypoint: "data.http.authz.allow"
        # HTTP proxy filter should be after OPA filter
        - name: dgp.filter.http.proxy
          config:
            # ... proxy config
```

#### Fields

**Server Mode Fields:**

- **`server_url`** *(string, required for server mode)*
  - **Meaning:** OPA server address, e.g., `http://opa-server:8181` or `https://opa.example.com:8181`.
  - **Datatype:** `string`.
  - **Example:** `"http://localhost:8181"`.

- **`decision_path`** *(string, required for server mode)*
  - **Meaning:** OPA REST API decision path, format: `/v1/data/<package>/<rule>`.
  - **Datatype:** `string`.
  - **Example:** `"/v1/data/http/authz/allow"`.

- **`timeout_ms`** *(integer, optional)*
  - **Meaning:** Request timeout for OPA server in milliseconds.
  - **Datatype:** `int`.
  - **Default:** `100` (100 milliseconds).

- **`bearer_token`** *(string, optional)*
  - **Meaning:** Bearer token for OPA server authentication.
  - **Datatype:** `string`.
  - **Notes:** Configure this field if OPA server requires authentication.

**Embedded Mode Fields (Backward Compatible):**

- **`policy`** *(string, required for embedded mode)*
  - **Meaning:** The **Rego module source code** (inline string). Loaded via `rego.Module("policy.rego", policy)`.
  - **Datatype:** `string` (multiline YAML recommended with `|`).

- **`entrypoint`** *(string, required for embedded mode)*
  - **Meaning:** The **Rego query string** passed to `rego.Query(...)`. Should be a valid query like `data.<package>.<rule>` (e.g., `data.http.authz.allow`).
  - **Datatype:** `string`.

#### Decision contract

- If the query result is a non-empty set whose first expression value is **`true`**, the request **continues**.
- Otherwise (empty results or value ≠ `true`), the filter **stops** (request denied).

### Policy input

The filter constructs an `input` object with the following keys, which correspond to the HTTP request.

```
input.method       # HTTP method string
input.path         # URL path (string)
input.headers      # map[string][]string
input.client_ip    # string
input.query        # map[string][]string (URL query)
input.host         # string
input.remote_addr  # string
input.user_agent   # string
input.route        # route entry object (opaque to policy; structure may change)
input.api          # API object (opaque)
input.params       # route params map
```

**Important Note: HTTP Header Canonicalization**

Go's `net/http` package automatically canonicalizes HTTP header names to **Title-Case** format (e.g., `X-Api-Key`, `Content-Type`). When accessing headers in your policy, use the canonicalized key names:

- ✅ Correct: `input.headers["X-Api-Key"]` or `input.headers["Content-Type"]`
- ❌ Incorrect: `input.headers["x-api-key"]` or `input.headers["x_api_key"]`

**Example:**
```rego
# Correct header access
allow {
    input.headers["X-Api-Key"][0] == "secret"
    input.headers["Content-Type"][0] == "application/json"
}
```

### OPA Server Deployment (Server Mode)

#### Docker Deployment

**1. Start OPA Server**

```bash
docker run -d \
  --name opa-server \
  -p 8181:8181 \
  openpolicyagent/opa:latest \
  run --server --addr :8181 --log-level info
```

**2. Upload Policy to OPA Server**

```bash
# Create policy file policy.rego
cat > policy.rego <<EOF
package http.authz

default allow = false

allow {
    input.method == "GET"
    input.path == "/status"
}
EOF

# Upload policy
curl -X PUT http://localhost:8181/v1/policies/authz \
  -H "Content-Type: text/plain" \
  --data-binary @policy.rego
```

**3. Test Policy**

```bash
curl -X POST http://localhost:8181/v1/data/http/authz/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "method": "GET",
      "path": "/status"
    }
  }'
```

#### Docker Compose Deployment (Recommended)

```yaml
version: '3.8'

services:
  opa-server:
    image: openpolicyagent/opa:latest
    container_name: opa-server
    ports:
      - "8181:8181"
    command:
      - "run"
      - "--server"
      - "--addr=:8181"
      - "--log-level=info"
    volumes:
      - ./opa-policies:/policies:ro
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8181/health"]
      interval: 10s
      timeout: 3s
      retries: 3
    restart: unless-stopped

  pixiu:
    image: apache/dubbo-go-pixiu:latest
    depends_on:
      opa-server:
        condition: service_healthy
    volumes:
      - ./configs:/configs:ro
    ports:
      - "8888:8888"
```

#### Policy Hot Updates

OPA server supports updating policies at runtime without restart:

```bash
# Update policy
curl -X PUT http://localhost:8181/v1/policies/authz \
  -H "Content-Type: text/plain" \
  --data-binary @new-policy.rego
```

### Minimal examples

**Server Mode Examples**

**1) Allow only GET /status (using OPA server)**

```yaml
- name: dgp.filter.http.opa
  config:
    server_url: "http://opa-server:8181"
    decision_path: "/v1/data/http/authz/allow"
    timeout_ms: 100
```

Corresponding OPA policy (uploaded to server):

```rego
package http.authz

default allow = false

allow {
    input.method == "GET"
    input.path == "/status"
}
```

**2) Allow requests with a specific header value (using OPA server)**

```yaml
- name: dgp.filter.http.opa
  config:
    server_url: "http://opa-server:8181"
    decision_path: "/v1/data/http/authz/allow"
    timeout_ms: 100
    bearer_token: "your-secret-token"
```

Corresponding OPA policy:

```rego
package http.authz

default allow = false

allow {
    input.headers["X-Api-Key"][0] == "secret"
}
```

**Embedded Mode Examples**

**1) Allow only GET /status (embedded mode)**

```yaml
- name: dgp.filter.http.opa
  config:
    policy: |
      package http.authz
      default allow = false
      allow { input.method == "GET"; input.path == "/status" }
    entrypoint: "data.http.authz.allow"
```

**2) Allow requests with a specific header value (embedded mode)**

```yaml
- name: dgp.filter.http.opa
  config:
    policy: |
      package http.authz
      default allow = false
      allow {
        input.headers["x-api-key"][0] == "secret"
      }
    entrypoint: "data.http.authz.allow"
```

### Limitations and notes

**Response Format Support**

The OPA filter supports multiple response formats and automatically recognizes them:

**Format 1: Boolean Value (Simple Cases)**
```rego
package http.authz
default allow = false
allow {
    input.method == "GET"
}
```

**Format 2: Object with "allow" Field (OPA Common Pattern, Recommended)**
```rego
package http.authz

default decision = {"allow": false}

decision = {"allow": true, "reason": "admin user"} {
    input.headers["X-Role"][0] == "admin"
}
```

**Format 3: Object with "result" Field**
```rego
package http.authz

decision = {"result": true, "metadata": {"user": "alice"}} {
    input.headers["X-Api-Key"][0] == "secret"
}
```

The filter extracts the decision using the following priority:
1. Direct boolean value → use that value
2. Object's `allow` field → use the boolean value of `allow`
3. Object's `result` field → use the boolean value of `result`
4. Unrecognized format → default to deny (returns 403)

---

**Server Mode:**
- **Network Latency**: Introduces network latency due to external OPA server calls (typically 5-50ms).
- **OPA Server Availability**: Ensure OPA server high availability; recommend deploying multiple instances.
- **Timeout Configuration**: Default timeout of 100ms may be too tight in cross-service/cloud environments. Consider adjusting to 300-500ms for production based on network latency.

**Embedded Mode:**
- **No custom deny status/body**: The filter does not map policy outputs to HTTP status or body.
- **Module-only loading**: Policies are loaded from the inline `policy` string only.
- **Policy updates require restart**: Modifying policies requires restarting Pixiu service.

### Troubleshooting

**Server Mode:**
- **Connection Failed**: Check `server_url` configuration, ensure OPA server is running and accessible.
- **Timeout Errors**: Increase `timeout_ms` value, or check OPA server performance.
- **Authentication Failed**: If OPA server requires authentication, ensure `bearer_token` is correctly configured.
- **Decision Path Error**: Ensure `decision_path` matches the policy path in OPA server.
- **Denied Unexpectedly**: Use `curl` to directly test OPA server's decision API to verify policy logic.

**Embedded Mode:**
- **Denied Unexpectedly**: Confirm the query is correct (e.g., `data.http.authz.allow`), and that the policy returns **`true`** for the given `input`.
- **Policy Compile Errors**: Validate the Rego module with `opa eval` locally before embedding.
- **Nil/Empty Results**: Re-check access to `headers`/`query` (they are maps of lists), and confirm path/method match.

**General Troubleshooting:**
- **Check Logs**: Pixiu logs detailed OPA evaluation information, including errors and decision results.
- **Test Policy**: Test the policy first in OPA server or locally using `opa eval` to ensure the logic is correct.