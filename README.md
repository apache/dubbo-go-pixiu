[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)


[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-pixiu.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-pixiu)

English | [中文](./README_CN.md)

# Introduction

**Dubbo-Go-Pixiu**(official site: https://dubbo.apache.org/zh/docs3-v2/dubbo-go-pixiu/) is a high-performance API gateway and multi-language solution Sidecar in the Dubbo ecosystem


![](https://dubbo-go-pixiu.github.io/img/pixiu-dubbo-ecosystem.png)

It is an open source Dubbo ecosystem API gateway, and also a sidecar to let other compute language program access the dubbo clusters by HTTP/gRPC protocol. As an API gateway, Pixiu can receive external network requests, convert them into dubbo and other protocol requests, and forward them to the back cluster; as a sidecar, Pixiu expects to register to the Dubbo cluster instead of the proxy service, allowing multilingual services to access the Dubbo cluster to provide faster solution


## Quick Start

you can find out all demo in https://github.com/apache/dubbo-go-pixiu-samples.
download it and operate as below.

#### cd samples dir

```
cd dubbogo/simple
```

we can use start.sh to run samples quickly. for more info, execute command as below for more help

```
./start.sh [action] [project]
./start.sh help
```

we run body samples below step

#### prepare config file and docker 

prepare command will prepare dubbo-server and pixiu config file and start docker container needed

```
./start.sh prepare body
```

if prepare config file manually, notice:
- modify $PROJECT_DIR in conf.yaml to absolute path in your compute 

#### start dubbo or http server

```
./start.sh startServer body
```

#### start pixiu 

```
./start.sh startPixiu body
```

if run pixiu manually in pixiu project, use command as below.

```
 go run cmd/pixiu/*.go gateway start -c /[absolute-path]/dubbo-go-pixiu/samples/dubbogo/simple/body/pixiu/conf.yaml
```


#### Try a request

use curl to try or use unit test

```bash
curl -X POST 'localhost:8881/api/v1/test-dubbo/user' -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header 'Content-Type: application/json' 
./start.sh startTest body
```

#### Clean

```
./start.sh clean body
```

## Start Docker

#### 
```shell
docker run --name pixiu-gateway -p 8888:8888 dubbogopixiu/dubbo-go-pixiu:latest

```
```
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
