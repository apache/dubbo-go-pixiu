[![Build Status](https://travis-ci.org/dubbogo/dubbo-go-proxy.svg?branch=master)](https://travis-ci.org/dubbogo/dubbo-go-proxy)

### 简介

dubbo-go-proxy 网关支持调用Java的dubbo集群和golang的dubbo-go集群。

现在dubbo-go-proxy 已经支持以dubbo协议和http协议调用远程的 dubbo 集群，未来还会支持更多的协议。

## 开始
#### 启动proxy
1. `make build` 构建出二进制包
2. `make start` 启动proxy

#### 快速入门
dubbo-go-proxy支持2种协议的调用:

1. [Http协议调用](https://github.com/dubbogo/dubbo-go-proxy/blob/develop/docs/sample/http.md) 
2. [Dubbo协议调用](https://github.com/dubbogo/dubbo-go-proxy/blob/develop/docs/sample/dubbo.md)

## 特性
1. dubbo-go-proxy 有完善的插件机制可以支持自定义的扩展
2. 正如其它的网关一样，主要来做流量和接口的治理
3. 兼容 java 和 go 的 dubbo 及 Spring Cloud（HTTP）
4. 支持传统的网关自定义映射外，还支持自动识别注册中心的 RPC 服务暴露 HTTP 接口
4. 能支持边车和集中部署（规划中）
5. istio 下支持 dubbo 协议流量（规划中）

## 架构图
![image](https://raw.githubusercontent.com/dubbogo/dubbo-go-proxy/master/docs/images/dubbgoproxy-infrastructure.png)
## 流程图
![image](https://raw.githubusercontent.com/dubbogo/dubbo-go-proxy/master/docs/images/dubbogoproxy-procedure.png)

## 术语解释
### 组件
#### Proxy
数据面板
#### Admin
控制面板
### 概念
#### 下游（Downstream）
下游主机连接到 Proxy ，发送请求并接收响应。（API 网关场景理解：浏览器）
#### 上游（Upstream）
上游主机接收来自 Proxy 的连接和请求并返回响应。（API 网关场景理解：dubbo 服务的机器）
#### 监听器（Listener）
监听器是可以被下游客户端连接的网络位置（例如，端口，unix域套接字等）。Proxy 公开一个或多个下游主机连接的监听器。
#### 集群（Cluster）
群集是指 Proxy 连接到的一组逻辑上相似的上游主机（比如 dubbo 集群）。Proxy 通过服务发现发现一个集群的成员，它可以通过主动健康检查来确定集群成员的健康度，从而 Proxy 通过负载均衡策略将请求路由到相应的集群成员。
#### 接口（Api）
接口是 API 网关的核心概念，特别针对浏览器等外部系统的访问时必须开启，所有请求必须匹配到对应的 Up 状态的接口才能继续进行后续逻辑。
#### 客户端（Client）
请求上游主机的真实调用对象。
#### 路由（Router）
路由策略，目前理解成 HTTP 路由，通过 match 逻辑路由到对应的集群。
#### 上下文（Context）
一个真实请求的上下文，包含这次请求几乎所有的信息。在各个环节都会使用，特别是 filter 链路。
#### 过滤器（Filter）
过滤器，提供拦截能力。支持自定义扩展。
#### 规则（Rule）
规则提供匹配能力，给过滤器，路由等其它概念使用。
## 联系我们
项目在快速迭代中，欢迎使用， 欢迎给出建议或者提交pr。钉钉群: 31363295
## License

Apache License, Version 2.0
