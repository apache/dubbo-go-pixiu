# Generate API config automatically from dubbogo registry

> Doc metions below fit the code in the `samples/dubbogo/simple/registry`

## Define Pixiu Config

```yaml
static_resources:
  listeners:
    - name: "net/http"
      protocol_type: "HTTP"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8881
      filter_chains:
          filters:
            - name: dgp.filter.httpconnectionmanager
              config:
                route_config:
                  routes:
                    - match:
                        prefix: "/"
                      route:
                        cluster: "test-dubbo"
                        cluster_not_found_response_code: 505
                http_filters:
                  - name: dgp.filter.http.apiconfig
                    config:
                      dynamic: true
                      dynamic_adapter: test
                  # configure the dubbogo registry
                  - name: dgp.filter.http.dubboproxy
                    config:
                      dubboProxyConfig:
                        registries:
                          "zookeeper":
                            protocol: "zookeeper"
                            timeout: "3s"
                            address: "127.0.0.1:2181"
                            username: ""
                            password: ""
                        timeout_config:
                          connect_timeout: 5s
                          request_timeout: 5s
                  - name: dgp.filter.http.response
                    config:
                server_name: "test_http_dubbo"
                generate_request_id: false
      config:
        idle_timeout: 5s
        read_timeout: 5s
        write_timeout: 5s
  clusters:
    - name: "test-dubbo"
      lb_policy: "RoundRobin"
      registries:
        "zookeeper":
          timeout: "3s"
          address: "127.0.0.1:2181"
          username: ""
          password: ""
  shutdown_config:
    timeout: "60s"
    step_timeout: "10s"
    reject_policy: "immediacy"
  adapters:
    - id: test
      name: dgp.adapter.dubboregistrycenter
      config:
        registries:
          "zookeeper":
            protocol: zookeeper
            address: "127.0.0.1:2181"
            timeout: "5s"
```

We should configure the `dgp.filter.http.dubboproxy` filter in the configuration file. 

## Test

The default API path is `/application/interface/version/method`. The types and values of the parameters should be placed in the request body.

Run command curl using: 

```
curl http://localhost:8881/BDTService/com.dubbogo.pixiu.UserService/1.0.0/GetUserByCode -X POST -d '{"types":"int", "values":1}'
```

[Previous](dubbo.md)
