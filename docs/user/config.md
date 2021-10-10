

### Config
 
pixiu supports specifying local config file with argument -c which you can find in those samples pixiu dir. 

pixiu use the config abstraction like envoy such as listener, filter, route and cluster.

Besides, pixiu provider another dubbo-specific config named api_config which dubbo-filter used to  transform http request to dubbo generic call. you can find it in those samples pixiu dir too
the api_config specification you can refer to this document [Api Model](api.md)

this document mainly describe the pixiu config abstraction, there is a example below:
```yaml

static_resources:
  listeners:
    - name: "net/http"
      address:
        socket_address:
          protocol_type: "HTTP"
          address: "0.0.0.0"
          port: 8888
      filter_chains:
        - filter_chain_match:
          domains:
            - api.dubbo.com
            - api.pixiu.com
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
the more detail will be found in pkg/model/bootstrap.go

### static_resources 
static_resources are used to specify unchanged config, meanwhile the dynamic_resources are used for dynamic config. supporting dynamic_resource is still in develop

there are four import abstraction in static_resources:
- listener
- filter
- route
- cluster

#### listener

listener provider external network server which support many network protocol, such as http, http2 or tcp

user can set the protocol and host to allow pixiu listen to it

when listener receive request from client, it will handle it and pass it to filter

#### filter

filter provider request handle abstraction. user can combinate many filter together into filter-chain.

when receive request from listener, filter will handle it in the order at pre or post phase.

because pixiu want offer network protocol transform function, so the filter contains network filter and the filter below network filter such as http filter

the request handle process like below
```
client -> listner -> network filter such as httpconnectionmanager -> http filter chain

```

##### network filter

support http protocol only. for example dgp.filter.httpconnectionmanager in config above

##### http filter 

there are also many protocol-specific filter such as http-to-grpc filter、http-to-dubbo filter etc. for example, dgp.filter.http.httpproxy in config above

there are many build-in filter such as cors、metric、ratelimit or timeout. for example, dgp.filter.http.response in config above

#### route

after filter handle, pixiu will forward the request to upstream server by route. the route provider forward rules such as path、method and header maetches

```
routes:
- match:
    prefix: "/user"
  route:
    cluster: "user"
    cluster_not_found_response_code: 505
```

#### cluster

the cluster represents the same service instance cluster which specify upstream server info 

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

the adapter communicated with service-registry such as zk, nacos to fetch service instance info and produce route and cluster config automaticly

there are two adapter you can find in pkg/adapter

