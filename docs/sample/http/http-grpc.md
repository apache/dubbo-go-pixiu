# Invoke service provider using grpc

> Doc metions below fit the code in the `samples/http/grpc`

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
                        prefix: "/api/v1"
                      route:
                        cluster: "test-grpc"
                        cluster_not_found_response_code: 505
                http_filters:
                  - name: dgp.filter.http.grpcproxy
                    config:
                      path: /mnt/d/WorkSpace/GoLandProjects/dubbo-go-pixiu/samples/http/grpc/proto
                  - name: dgp.filter.http.response
                    config:
                server_name: "test-http-grpc"
                generate_request_id: false
      config:
        idle_timeout: 5s
        read_timeout: 5s
        write_timeout: 5s
  clusters:
    - name: "test-grpc"
      lb_policy: "RoundRobin"
      endpoints:
        - socket_address:
            address: 127.0.0.1
            port: 50001
            protocol_type: "GRPC"
  timeout_config:
    connect_timeout: "5s"
    request_timeout: "10s"
  shutdown_config:
    timeout: "60s"
    step_timeout: "10s"
    reject_policy: "immediacy"
```

> Grpc server is defined in the `clusters`

> WARN: currently http request only support json body to parse params

## Prepare for the Server 

generate pb files under `samples/http/grpc/proto` with command: 

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hello_grpc.proto
```

## Start server 

run `samples/http/grpc/server/app/server.go`

```
go run server.go
```

## Build dubbo-go-pixiu

```
go build -o dubbo-go-pixiu cmd/pixiu/*
```

Or on Windows, you need to go under `cmd/pixiu` to build binary file

```
go build -o dubbo-go-pixiu.exe .
```

## Start Pixiu

```
./dubbo-go-pixiu gateway start --config samples/http/grpc/pixiu/conf.yaml --log-config samples/http/grpc/pixiu/log.yml
```

## Test Using Curl

Run command curl using: 

```
curl http://127.0.0.1:8881/api/v1/provider.UserProvider/GetUser
```

and 

```
curl http://127.0.0.1:8881/api/v1/provider.UserProvider/GetUser -X POST -d '{"userId":1}'
```

> If response body is a json, the header of 'content-type' will set to 'application/json'. If it is just a plain text, the header of 'content-type' is 'text/plain'.
