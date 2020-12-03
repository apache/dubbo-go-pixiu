[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-proxy.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-proxy)

[中文](./README_CN.md) 
### Introduction

dubbo-go-proxy is a gateway that mainly focuses on providing gateway solution to your Dubbo and RESTful services.

It supports HTTP-to-Dubbo and HTTP-to-HTTP proxy and more protocols will be supported in the near future.

## Get Start
#### Start proxy
1. Build the binary with `make build`
2. Start the proxy with `make start`

#### Quick Guide
dubbo-go-proxy supports to invoke 2 protocols:

1. [Http](https://github.com/dubbogo/dubbo-go-proxy/blob/develop/docs/sample/http.md) 
2. [Dubbo](https://github.com/dubbogo/dubbo-go-proxy/blob/develop/docs/sample/dubbo.md)

## Features
1. You can customize your own dubbo-go-proxy with plugin.
2. Multiple default filters to manage your APIs.
3. Dubbo and HTTP proxies.
4. Customizable request parameters mapping.
5. Automatically recognizes RPC services from service registration center and exposes it in HTTP protocol.
4. Sidecar or centralized deployment（Planning）
5. Dubbo protocol's rate-limiting in Istio environment（Planning）

## Architecture Diagram
![image](https://raw.githubusercontent.com/dubbogo/dubbo-go-proxy/master/docs/images/dubbgoproxy-infrastructure.png)
## Flow Diagram
![image](https://raw.githubusercontent.com/dubbogo/dubbo-go-proxy/master/docs/images/dubbogoproxy-procedure.png)

## Teams
### Components
#### Proxy
Data panel
#### Admin
Control Panel
### Concepts
#### Downstream
Downstream is the requester who sends request to and expecting the response from dubbo-go-proxy. (Eg.Postman client, Browser)
#### Upstream
The service that receive requests and send responses to dubbo-go-proxy. (Eg. Dubbo server)
#### Listener
The way that the dubbo-go-proxy exposes services to upstream clients. It could be configured to multiple listeners for one dubbo-go-proxy.
#### Cluster
Cluster is a set of upstream services that logically similar, such as dubbo cluster. Proxy can identifies the cluster members through service discovery and proactively probes their healthiness so that the proxy can route the requests to proper cluster member base on load balancing strategies.
#### Api
API is the core concept of the dubbo-go-proxy, all the upstream services will be configured and exposed through API.
#### Client
The actual caller of the upstream services.
#### Router
Router routes the HTTP request to proper upstream services according to the API configs.
#### Context
The context of a request in dubbo-go-proxy includes almost all information to get response from upstream services. It will be used in almost all process in the dubbo-go-proxy, especially the filter chain.
#### Filter
Filter manipulate the incoming requests. It is extensible for the users.
## Contact Us
The project is under intensively iteration, you are more than welcome to use, suggest and contribute codes. DingDing Group: 31363295
## License

Apache License, Version 2.0
