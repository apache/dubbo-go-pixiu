[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)


[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-pixiu.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-pixiu)

English | [中文](./README_CN.md)

# Introduction

**Dubbo-Go-Pixiu**(official site: https://dubbo.apache.org/zh/docs3-v2/dubbo-go-pixiu/) is a high-performance API gateway and multi-language solution Sidecar in the Dubbo ecosystem


![](https://dubbo-go-pixiu.github.io/img/pixiu-dubbo-ecosystem.png)

It is an open source Dubbo ecosystem API gateway, and also a sidecar to let other compute language program access the dubbo clusters by HTTP/gRPC protocol. As an API gateway, Pixiu can receive external network requests, convert them into dubbo and other protocol requests, and forward them to the back cluster; as a sidecar, Pixiu expects to register to the Dubbo cluster instead of the proxy service, allowing multilingual services to access the Dubbo cluster to provide faster solution


## Quick Start

#### Requirment
1. go 1.17 or higher
2. docker or docker-desktop

you can find out all demo in https://github.com/apache/dubbo-go-pixiu-samples.
download it and operate as below.
```shell
git clone https://github.com/apache/dubbo-go-pixiu-samples.git
```

#### update pixiu to latest version
```shell
go get github.com/apache/dubbo-go-pixiu@v0.6.0-rc2
```

#### cd samples dir

```shell
cd dubbogo/simple
```

we can use start.sh to run samples quickly. for more info, execute command as below for more help

```shell
./start.sh [action] [project]
./start.sh help
```

we run [direct] samples step by step as follows.

#### prepare config file and docker 

'prepare' command will prepare dubbo-server and pixiu config file firstly, and then start docker container.

```shell
./start.sh prepare direct
```

if prepare config file manually, notice:
- modify $PROJECT_DIR in conf.yaml to absolute path 

#### start dubbo or http server

```shell
./start.sh startServer direct
```

#### start pixiu 

```shell
./start.sh startPixiu direct
```

if run pixiu manually in pixiu project, use command as below.

```shell
 go run cmd/pixiu/*.go gateway start -c /[absolute-path]/dubbo-go-pixiu/samples/dubbogo/simple/direct/pixiu/conf.yaml
```

if run pixiu manually in pixiu project and wasm feature, use command as below.

build pixiu project use command operate

```shell
go build -tags wasm -o pixiu cmd/pixiu/*.go
```

run pixiu app binary

```shell
go build cmd/pixiu/*.go
./pixiu gateway start -c /[absolute-path]/dubbo-go-pixiu/samples/dubbogo/simple/direct/pixiu/conf.yaml
```

#### Try a request

use curl to try or use unit test

```shell
curl http://localhost:8883/UserService/com.dubbogo.pixiu.UserService/GetUserByCode \
-H "x-dubbo-http1.1-dubbo-version:1.0.0" -H "x-dubbo-service-protocol:dubbo" \
-H "x-dubbo-service-version:1.0.0" -H "x-dubbo-service-group:test" \
-H "Content-Type:application/json" \
 -d '[1]'
```
```shell
curl http://localhost:8883/UserService/com.dubbogo.pixiu.UserService/UpdateUserByName  \
-H "x-dubbo-http1.1-dubbo-version:1.0.0" -H "x-dubbo-service-protocol:dubbo" \
-H "x-dubbo-service-version:1.0.0" -H "x-dubbo-service-group:test" \
-H "Content-Type:application/json" \
 -d '["tc",{"id":"0002","code":1,"name":"tc","age":15}]'
```
```shell
curl http://localhost:8883/UserService/com.dubbogo.pixiu.UserService/GetUserByCode \
-H "x-dubbo-http1.1-dubbo-version:1.0.0" -H "x-dubbo-service-protocol:dubbo" \
-H "x-dubbo-service-version:1.0.0" -H "x-dubbo-service-group:test" \
-H "Content-Type:application/json" \
 -d '[1]'
```

```shell 
./start.sh startTest body
```

#### Clean

```shell
./start.sh clean direct
```

## Start Docker

#### 
```shell
docker run --name pixiu-gateway -p 8888:8888 dubbogopixiu/dubbo-go-pixiu:latest
```

```shell
docker run --name pixiu-gateway -p 8888:8888 \
    -v /yourpath/conf.yaml:/etc/pixiu/conf.yaml \
    -v /yourpath/log.yml:/etc/pixiu/log.yml \
    dubbogopixiu/dubbo-go-pixiu:latest
```

## Features

- Multi-protocol support: Currently, Http, Dubbo2, Triple, gRPC protocol proxy and conversion are supported, and other protocols are being continuously integrated.
- Safety certificate: Support HTTPS, JWT Token verification and other security authentication measures.
- Registry integration: Support to obtain service metadata from Dubbo or Spring Cloud cluster, support ZK, Nacos registry.
- Traffic management: Integrate with sentinel, support multiple protocols for rate limiting.
- Observability: Integrate with opentelemetry and jaeger for distributed tracing.
- Admin and visual interface: Have pixiu-admin for remote administration and visualization

### Control Plane

The pixiu control plane is forked from [istio](https://github.com/istio/istio) v1.14.3. Offers a variety of capabilities, including service discovery, traffic management, security management.

## Contact Us

The project is under intensively iteration, you are more than welcome to use, suggest and contribute codes. 


### Community

**DingDing Group (31203920):**

[![flowchart](./docs/images/group-pixiu-dingding.jpg)](docs/images/group-pixiu-dingding.jpg)

We welcome the friends who can give us constructing suggestions instead of known-nothing.

## License

Apache License, Version 2.0
