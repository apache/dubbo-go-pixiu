[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-pixiu.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-pixiu)

[中文](./README_CN.md) 
### Introduction

dubbo-go-pixiu is a gateway that mainly focuses on providing gateway solution to your Dubbo and RESTful services.

It supports HTTP-to-Dubbo and HTTP-to-HTTP proxy and more protocols will be supported in the near future.

## Get Start
#### Start Pixiu
1. Build the binary with `make build`
2. Start the pixiu with `make start`

#### Quick Guide
dubbo-go-pixiu supports to invoke 2 protocols:

1. [Http](https://github.com/dubbogo/dubbo-go-pixiu/blob/develop/docs/sample/http.md) 
2. [Dubbo](https://github.com/dubbogo/dubbo-go-pixiu/blob/develop/docs/sample/dubbo.md)

## Features
1. You can customize your own dubbo-go-pixiu with plugin.
2. Multiple default filters to manage your APIs.
3. Dubbo and HTTP proxies.
4. Customizable request parameters mapping.
5. Automatically recognizes RPC services from service registration center and exposes it in HTTP protocol.
4. Sidecar or centralized deployment（Planning）
5. Dubbo protocol's rate-limiting in Istio environment（Planning）

## Architecture Diagram
![image](https://raw.githubusercontent.com/dubbogo/dubbo-go-pixiu/master/docs/images/dubbgopixiu-infrastructure.png)
## Flow Diagram
![image](https://raw.githubusercontent.com/dubbogo/dubbo-go-pixiu/master/docs/images/dubbogopixiu-procedure.png)

## Teams
### Components
#### Pixiu
Data panel
#### Admin
Control Panel
### Concepts
#### Downstream
Downstream is the requester who sends request to and expecting the response from dubbo-go-pixiu. (Eg.Postman client, Browser)
#### Upstream
The service that receive requests and send responses to dubbo-go-pixiu. (Eg. Dubbo server)
#### Listener
The way that the dubbo-go-pixiu exposes services to upstream clients. It could be configured to multiple listeners for one dubbo-go-pixiu.
#### Cluster
Cluster is a set of upstream services that logically similar, such as dubbo cluster. pixiu can identifies the cluster members through service discovery and proactively probes their healthiness so that the pixiu can route the requests to proper cluster member base on load balancing strategies.
#### Api
API is the core concept of the dubbo-go-pixiu, all the upstream services will be configured and exposed through API.
#### Client
The actual caller of the upstream services.
#### Router
Router routes the HTTP request to proper upstream services according to the API configs.
#### Context
The context of a request in dubbo-go-pixiu includes almost all information to get response from upstream services. It will be used in almost all process in the dubbo-go-pixiu, especially the filter chain.
#### Filter
Filter manipulate the incoming requests. It is extensible for the users.
## Contact Us
The project is under intensively iteration, you are more than welcome to use, suggest and contribute codes. DingDing Group: 31363295
## License

Apache License, Version 2.0
