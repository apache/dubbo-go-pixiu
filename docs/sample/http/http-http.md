# Http to Http Proxy

> Doc metions below fit the code in the `samples/http/simple`

## Define Apis in the pixiu/conf.yaml

```yaml
---
static_resources:
  listeners:
    - name: "net/http"
      protocol_type: "HTTP"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8888
      filter_chains:
          filters:
            - name: dgp.filter.httpconnectionmanager
              config:
                route_config:
                  routes:
                    - match:
                        prefix: "/user"
                      route:
                        cluster: "user"
                        cluster_not_found_response_code: 505
                http_filters:
                  - name: dgp.filter.http.httpproxy
                    config:
                  - name: dgp.filter.http.cors
                    config:
                      allow_origin:
                        - api.dubbo.com
                      allow_methods: ""
                      allow_headers: ""
                      expose_headers: ""
                      max_age: ""
                      allow_credentials: false
                  - name: dgp.filter.http.response
                    config:
      config:
        idle_timeout: 5s
        read_timeout: 5s
        write_timeout: 5s
  clusters:
    - name: "user"
      lb_policy: "lb"
      endpoints:
        - id: 1
          socket_address:
            address: 127.0.0.1
            port: 1314
  timeout_config:
    connect_timeout: "5s"
    request_timeout: "10s"
  shutdown_config:
    timeout: "60s"
    step_timeout: "10s"
    reject_policy: "immediacy"
```

for custom config , you can refer to [config](../../user/config.md) in user-guide
