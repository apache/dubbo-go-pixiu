# OPA Quick Start

Examples of official references is in `https://github.com/dubbo-go-pixiu-samples/dubbo/simple

###### conf.yaml

```
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
                      prefix: "/UserService"
                    route:
                      cluster: "user"
                      cluster_not_found_response_code: 505
                  - match:
                      prefix: "/OtherService"
                    route:
                      cluster: "user"
                      cluster_not_found_response_code: 505
              http_filters:
                - name: dgp.filter.http.opa
                  config:
                    policy: |
                      package pixiu
                      import future.keywords.if
                      default allow := false

                      allow if {
                        input.path == "/UserService"
                        input.headers["Test_header"][0] == "1"
                      }
                    entrypoint: data.pixiu.allow
                - name: dgp.filter.http.httpproxy
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
  shutdown_config:
    timeout: "60s"
    step_timeout: "10s"
    reject_policy: "immediacy"
```



Start Zookeeper[Docker environment]:

```
cd samples\dubbogo\simple\opa\docker
run docker-compose.yml/services
```

Start Http [Go environment]:

```
go run samples/dubbogo/simple/opa/server/server.go
```

Start Pixiu:

```
go run cmd/pixiu/*.go gateway start -c samples/dubbogo/simple/opa/pixiu/conf.yaml
```



### Test

##### Denied

```
curl -s http://127.0.0.1:8888/OtherService
```

**Expected:** `null` (OPA default-deny triggered; filter stopped the request, gateway responded locally).

##### Allow 

```
curl -s http://127.0.0.1:8888/UserService -H "Test_header: 1"
```

**Expected:**

```
{"message":"UserService","result":"pass"}
```

(Without the header, or with `Test_header: 0`, you’ll see `null` again due to default-deny. This header-gated allow/deny behavior is the same pattern validated by the unit tests: `1` allows, `0` denies.) 

##### OR Run pixiu_test.go

```
cd samples/dubbogo/simple
go test ./opa/test -v
```





