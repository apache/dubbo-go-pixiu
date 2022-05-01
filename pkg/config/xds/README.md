## Xds adaptor
This module helps pixiu connect the go-control-panel, 
which implements by pixiu-admin(https://github.com/dubbogo/pixiu-admin), and read cluster and initial listener 
configurations. Pixiu talk to pixiu admin over XDS API. By this module,  pixiu can be config by pixiu admin 
dynamically.

## How to use

To config this module, some pixiu config should be change:

---

The node must add into the config to identify unique pixiu node.
```yaml
node:
  id: "test-id"
  cluster: "pixiu"
```

Starting multi-node should give different node.id

"static_resources" should refer to the pixiu-admin.  
```yaml
static_resources:
  clusters:
    - name: "xds-server"
      type: "Static"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 18000
```
config "lds_config" "cds_config" refer to the static_resources's cluster above.

```yaml
dynamic_resources:
  lds_config: # config listener
    cluster_name: ["xds-server"]
    api_type: "GRPC"
    refresh_delay: "5s"
    request_timeout: "10s"
    grpc_services:
      - timeout: "5s"
  cds_config: # config cluster
    cluster_name: ["xds-server"]
    api_type: "GRPC"
    refresh_delay: "5s"
    request_timeout: "10s"
    grpc_services:
      - timeout: "5s"
```

Full demo config like below:

---
```yaml
node:
  id: "test-id"
  cluster: "pixiu"

dynamic_resources:
  lds_config:
    cluster_name: ["xds-server"]
    api_type: "GRPC"
    refresh_delay: "5s"
    request_timeout: "10s"
    grpc_services:
      - timeout: "5s"
  cds_config:
    cluster_name: ["xds-server"]
    api_type: "GRPC"
    refresh_delay: "5s"
    request_timeout: "10s"
    grpc_services:
      - timeout: "5s"
static_resources:
  clusters:
    - name: "xds-server"
      type: "Static"
      endpoints:
        - socket_address:
            address: "127.0.0.1"
            port: 18000

  shutdown_config:
    timeout: "60s"
    step_timeout: "10s"
    reject_policy: "immediacy"
```

## How does it work

This module start XDS API client when pixiu init. Once dynamic_resources config, it tries to connect go-control-panel 
and pull dubbo-go.pixiu/v1/discovery:listener(listener configurations) and dubbo-go.pixiu/v1/discovery:cluster(cluster 
configuration) via ExtensionConfigDiscovery.

This will fetch cluster configuration at initial time and consume ongoing changes. But for listener configuration, 
later changes has no affect until pixiu restart.


## run unit test

```shell
GOARCH=amd64 go test -gcflags=-l -v -cover .
```