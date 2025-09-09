# OPA Filter (dgp.filter.http.opa)

English | [中文](opa_CN.md)

---

## English

### Overview
The `dgp.filter.http.opa` filter delegates authorization decisions to Open Policy Agent (OPA) via a Rego policy. This filter evaluates requests and determines whether to allow or deny based on the policy defined in Rego. The policy is provided as an inline Rego module and evaluated using OPA's built-in query engine.

### What the filter does (current behavior)
- Loads a Rego **module string** from `config.policy`.
- Builds a Rego **query** from `config.entrypoint`.
- For each incoming request, constructs an `input` object and evaluates the query.
- If the query result is `true`, the request is allowed. Otherwise, the request is denied.

> There is **no built-in support** for external policy files or URIs, custom HTTP status codes, or custom error bodies.

### Configuration schema
Add the filter under your HTTP connection manager’s `http_filters` list.

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

- **`policy`** *(string, required)*
  - **Meaning:** The **Rego module source code** (inline string). Loaded via `rego.Module("policy.rego", policy)`.
  - **Datatype:** `string` (multiline YAML recommended with `|`).
  - **Notes:** File paths or bundle URIs are **not supported**.
- **`entrypoint`** *(string, required)*
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

### Minimal examples

**1) Allow only GET /status**

```yaml
- name: dgp.filter.http.opa
  config:
    policy: |
      package http.authz
      default allow = false
      allow { input.method == "GET"; input.path == "/status" }
    entrypoint: "data.http.authz.allow"
```

**2) Allow requests with a specific header value**

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

- **Return type must be boolean**: Only `true` allows; objects (e.g., `{allow: true}`) will not be interpreted specially.
- **No custom deny status/body**: The filter does not map policy outputs to HTTP status or body.
- **Module-only loading**: Policies are loaded from the inline `policy` string only.

### Troubleshooting

- **Denied unexpectedly**: Confirm the query is correct (e.g., `data.http.authz.allow`), and that the policy returns **`true`** for the given `input`.
- **Policy compile errors**: Validate the Rego module with `opa eval` locally before embedding.
- **Nil/empty results**: Re-check access to `headers`/`query` (they are maps of lists), and confirm path/method match.