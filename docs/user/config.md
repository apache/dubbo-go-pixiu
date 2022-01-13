

### Config
 
Pixiu supports specifying local config file with argument `-c` which you can find in those samples pixiu dir. 

Pixiu uses the config abstraction like envoy such as listener, filter, route and cluster.

Besides, pixiu provides another dubbo-specific config named `api_config`, by which `dubbo-filter` can transform the http request to dubbo generic call. You can also find it in those samples's pixiu directory.

The document [Api Model](api.md) provides the api_config specification about the pixiu config abstraction.

This document mainly describes the pixiu config abstraction, there is a example below:
```yaml

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
                  - name: dgp.filter.http.response
                    config:
  clusters:
    - name: "user"
      lb_policy: "lb"
      endpoints:
        - id: 1
          socket_address:
            address: 127.0.0.1
            port: 1314
```
More detail will be found in `pkg/model/bootstrap.go`

### static_resources 

The `static_resources` are used to specify unchanged config, meanwhile the `dynamic_resources` are used for dynamic config. The `dynamic_resource` feature is still in developing now.


There are four import abstraction in `static_resources`:
- listener
- filter
- route
- cluster

#### listener

`Listener` provides external network server function which support many network protocol, such as http, http2 or tcp.

User can set the protocol and host to allow pixiu listen to it.

When listener receives request from client, it will process it and pass it to the `filter`.


#### filter

`Filter` provides request handle abstraction. You can combine many filters together into filter-chain.

When `filter` receives request from the listener, it will handle it orderly at its pre or post phase.

Because pixiu wants offer network protocol transform function, so the `filter` contains the network filter, such as the http filter.

The request processing order is as follows.

```
client -> listner -> network filter such as httpconnectionmanager -> http filter chain

```

##### network filter

Pixiu supports http protocol only, such as `dgp.filter.httpconnectionmanager` in the above config.


##### http filter 

There also are many protocol-specific filters such as http-to-grpc/http-to-dubbo etc, such as `dgp.filter.http.httpproxy` in the above config.


There are many build-in filters such as cors/metric/ratelimit/timeout, such as `dgp.filter.http.response` in the above config.


#### route

After `filter` handled the request, pixiu will forward the request to upstream server by `route`. The `route` provider forward rules such as path/method/header matches


```
routes:
- match:
    prefix: "/user"
  route:
    cluster: "user"
    cluster_not_found_response_code: 505
```

#### cluster

The `cluster` represents the same service instance cluster which specify upstream server info.

```
clusters:
- name: "user"
  lb_policy: "lb"
  endpoints:
    - id: 1
      socket_address:
        address: 127.0.0.1
        port: 1314
```


#### Adapter

The `adapter` communicates with service-registry such as zk/nacos to fetch service instance info and also produces the route and cluster config.

There are two adapters in `pkg/adapter`.


