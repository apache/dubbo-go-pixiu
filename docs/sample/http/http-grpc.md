# Invoke service provider using grpc

> Doc metions below fit the code in the `samples/http/grpc`

## Define Apis in the pixiu/api_config.yaml

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/provider.UserProvider/GetUser'
    type: restful
    description: user
    methods:
      - httpVerb: GET
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: grpc
          group: "test"
          version: 1.0.0
          clusterName: "test_grpc"
  - path: '/api/v1/provider.UserProvider/GetUser'
    type: restful
    description: user
    methods:
      - httpVerb: POST
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: grpc
          group: "test"
          version: 1.0.0
          clusterName: "test_grpc"
```

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