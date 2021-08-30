[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-pixiu.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-pixiu)

[中文](./README_CN.md)  |  [Piuxiu-SpringCloud 说明](./pkg/adapter/springcloud/doc/sc_design.md)

### Mascot

<div>
<table>
  <tbody>
  <tr></tr>
    <tr>
      <td align="center"  valign="middle">
        <a href="http://alexstocks.github.io/html/dubbogo.html" target="_blank">
          <img width="513px" height="513px" src="docs/images/pixiu-logo.png">
        </a>
      </td>
    </tr>
    <tr></tr>
  </tbody>
</table>
</div>
    
### Introduction

dubbo-go-pixiu is a gateway that mainly focuses on providing gateway solution to your Dubbo and RESTful services.

It supports HTTP-to-Dubbo and HTTP-to-HTTP proxy and more protocols will be supported in the near future.

## Quick Start

### 1. start provider

#### 1.1 Modify the zookeeper configuration of the following files

```
samples/dubbogo/http/server/profiles/dev/server.yml
```

#### 1.2 Switch to the root directory of dubbo-go-pixiu to execute

```
export CONF_PROVIDER_FILE_PATH=$PWD/samples/dubbogo/http/server/profiles/dev/server.yml
export APP_LOG_CONF_FILE=$PWD/samples/dubbogo/http/server/profiles/dev/log.yml
```

#### 1.3 Compile the provider code
```
go build -o server samples/dubbogo/http/server/app/*.go 
```

#### 1.4 Execute the binary file in the project root directory

```
./server 
```

### 2. Start pixiu

#### 2.1 Modify the zookeeper configuration of the following files

```
configs/conf.yaml 
```

#### 2.2 Compile pixiu code

```
go build -o pixiu cmd/pixiu/*.go 
```

#### 2.3 Execute the binary file in the project root directory

```
./pixiu gateway start
```

### 3. Try a request

```
curl -X POST 'localhost:8888/api/v1/test-dubbo/user' -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header 'Content-Type: application/json' 
```

## Features

1. You can customize your own dubbo-go-pixiu with plugin.
2. Multiple default filters to manage your APIs.
3. Dubbo and HTTP proxies.
4. Customizable request parameters mapping.
5. Automatically recognizes RPC services from service registration center and exposes it in HTTP protocol.
4. Sidecar or centralized deployment（Planning）
5. Dubbo protocol's rate-limiting in Istio environment（Planning）

## Architecture Diagram

<div>
<table>
  <tbody>
  <tr></tr>
    <tr>
      <td align="center"  valign="middle">
        <a href="http://alexstocks.github.io/html/dubbogo.html" target="_blank">
          <img width="800px" height="800px" src="./docs/images/dubbogopixiu-new-infrastructure.png">
        </a>
      </td>
    </tr>
    <tr></tr>
  </tbody>
</table>
</div>

## Flow Diagram

<div>
<table>
  <tbody>
  <tr></tr>
    <tr>
      <td align="center"  valign="middle">
        <a href="http://alexstocks.github.io/html/dubbogo.html" target="_blank">
          <img width="850px" height="150px" src="./docs/images/dubbogopixiu-procedure.png">
        </a>
      </td>
    </tr>
    <tr></tr>
  </tbody>
</table>
</div>

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

## Community

If u want to communicate with our community, pls scan the following dubbobo Ding-Ding QR code or search our commnity DingDing group code 31363295.

<div>
<table>
  <tbody>
  <tr></tr>
    <tr>
      <td align="center"  valign="middle">
        <a href="http://alexstocks.github.io/html/dubbogo.html" target="_blank">
          <img width="80px" height="85px" src="./docs/images/dubbogo-dingding.png">
        </a>
      </td>
    </tr>
    <tr></tr>
  </tbody>
</table>
</div>

The wechat public account of out community is as follows.

<div>
<table>
  <tbody>
  <tr></tr>
    <tr>
      <td align="center"  valign="middle">
          <img width="80px" height="115px" src="./docs/images/dubbogo-wechat.png">
        </a>
      </td>
    </tr>
    <tr></tr>
  </tbody>
</table>
</div>


We welcome the friends who can give us constructing suggestions instead of known-nothing.

## License

Apache License, Version 2.0
