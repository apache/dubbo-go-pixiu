[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)


[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-pixiu.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-pixiu)

English | [中文](./README_CN.md)

# Introduction

**Dubbo-Go-Pixiu** is a gateway that mainly focuses on providing gateway solution to your Dubbo and RESTful services.

It supports **HTTP-to-Dubbo** and **HTTP-to-HTTP** proxy and more protocols will be supported in the near future.

## Quick Start

#### cd samples dir

```
cd samples/dubbogo/simple
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

if run pixiu manually, use command as below

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

## Features

- You can customize your own dubbo-go-pixiu with plugin.
- Multiple default filters to manage your APIs.
- Dubbo and HTTP proxies.
- Customizable request parameters mapping.
- Automatically recognizes RPC services from service registration center and exposes it in HTTP protocol.
- Sidecar or centralized deployment（Planning）
- Dubbo protocol's rate-limiting in Istio environment（Planning）

## Architecture

### Pixiu Architecture

[![architecture](./docs/images/dubbogopixiu-new-infrastructure.png)](http://alexstocks.github.io/html/dubbogo.html)

### Pixiu Flow Chart

[![flow chart](./docs/images/dubbogopixiu-procedure.png)](http://alexstocks.github.io/html/dubbogo.html)

## Term

### Components

- Pixiu : Data panel

- Admin : Control Panel

### Concepts

- Downstream :  Downstream is the requester who sends request to and expecting the response from dubbo-go-pixiu. (Eg.Postman client, Browser)

- Upstream : The service that receive requests and send responses to dubbo-go-pixiu. (Eg. Dubbo server)

- Listener : The way that the dubbo-go-pixiu exposes services to upstream clients. It could be configured to multiple listeners for one dubbo-go-pixiu.

- Cluster : Cluster is a set of upstream services that logically similar, such as dubbo cluster. pixiu can identifies the cluster members through service discovery and proactively probes their healthiness so that the pixiu can route the requests to proper cluster member base on load balancing strategies.

- Api : API is the core concept of the dubbo-go-pixiu, all the upstream services will be configured and exposed through API.

- Client : The actual caller of the upstream services.

- Router : Router routes the HTTP request to proper upstream services according to the API configs.

- Context : The context of a request in dubbo-go-pixiu includes almost all information to get response from upstream services. It will be used in almost all process in the dubbo-go-pixiu, especially the filter chain.

- Filter : Filter manipulate the incoming requests. It is extensible for the users.

## Contact Us

The project is under intensively iteration, you are more than welcome to use, suggest and contribute codes. 

### Community
 
**DingDing Group (31363295):**

[![flowchart](./docs/images/dubbogo-dingding.png)](docs/images/dubbogo-dingding.png)

**WeChat Official Account:** 

[![flowchart](./docs/images/dubbogo-wechat.png)](docs/images/dubbogo-wechat.png)


We welcome the friends who can give us constructing suggestions instead of known-nothing.

## License

Apache License, Version 2.0
