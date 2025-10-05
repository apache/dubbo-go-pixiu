[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)

# Dubbo-Go-Pixiu：新一代高性能 API 网关

[![Build Status](https://github.com/apache/dubbo-go-pixiu/workflows/CI/badge.svg)](https://travis-ci.org/apache/dubbo-go-pixiu)
[![codecov](https://codecov.io/gh/apache/dubbo-go-pixiu/branch/master/graph/badge.svg)](https://codecov.io/gh/apache/dubbo-go-pixiu)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/apache/dubbo-go-pixiu?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/dubbo-go-pixiu)](https://goreportcard.com/report/github.com/apache/dubbo-go-pixiu)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

[English](README.md) | **中文**

-----

**Dubbo-Go-Pixiu** 是一款基于 Go 语言构建的高性能 API 网关。作为 [Apache Dubbo](https://dubbo.apache.org/) 生态系统的关键组件，它提供了丰富的流量管理、协议转换和安全防护等能力。


## 🚀 为什么选择 Dubbo-Go-Pixiu？

* **高性能**：基于 Go 语言构建，提供低延迟、高吞吐的网关能力。
* **无缝集成 Dubbo**：作为官方 Sidecar 方案，帮助非 Java 应用（Go、Python、Node.js 等）轻松调用 Dubbo 服务。
* **云原生设计**：为现代微服务和云原生架构而生，全面支持容器化部署。
* **高可扩展性**：灵活的过滤器和插件机制，让您轻松定制功能。

**即刻体验 Pixiu 网关功能**：请访问我们的 [使用示例](https://github.com/apache/dubbo-go-pixiu-samples)。

## ✨ 我们正在演进为 AI 网关 [开发中]

我们正在将 Pixiu 升级为**新一代 AI 网关**，旨在成为连接用户与大语言模型（LLMs）的桥梁。通过 Pixiu，您可以：

* **简化访问**：以统一、安全的方式接入各类 LLM 服务。
* **增强能力**：利用网关强大的插件体系，为您的 AI 应用增加认证、可观测性和流量控制等功能。
* **成本效益**：通过精细化的计费、审计和缓存策略，优化您的 AI 服务成本。

**即刻体验 AI 网关功能**：请访问我们的 [AI 网关示例](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/llm)。

## 核心功能

| 功能类别 | 描述 |
| :--- | :--- |
| 🚀 **协议处理** | 支持 HTTP、gRPC、Dubbo2、Triple 协议的代理和相互转换，提供强大的协议网关能力。 |
| 🛡️ **安全防护** | 提供 HTTPS、JWT 令牌验证、OAuth2 等多种安全机制，为您的服务保驾护航。 |
| 🔗 **服务发现** | 无缝集成 Zookeeper、Nacos 等注册中心，自动发现 Dubbo 和 Spring Cloud 集群中的服务。 |
| ⚖️ **流量治理** | 集成 Sentinel，提供精细化的多协议限流、熔断和服务降级能力。 |
| 📈 **可观测性** | 集成 OpenTelemetry 和 Jaeger，提供分布式追踪、指标和日志功能。 |
| 🎨 **可视化管理** | 配套的 **Pixiu-Admin** 控制台提供友好的 Web UI，支持远程服务管理和可视化配置。 |

## 快速开始

本指南将引导您，基于我们的[使用示例](https://github.com/apache/dubbo-go-pixiu-samples)启动一个 Pixiu 网关，并通过 HTTP 协议访问一个后端服务。

### 前置条件

* Go 1.17 或更高版本。
* 两个独立的终端窗口。

### 第一步：获取 Pixiu 源码

在**终端 1** 中执行：

```shell
git clone https://github.com/apache/dubbo-go-pixiu.git
cd dubbo-go-pixiu
```

### 第二步：启动后端示例服务

在**终端 2** 中执行：

```shell
git clone https://github.com/apache/dubbo-go-pixiu-samples.git
cd dubbo-go-pixiu-samples/http/simple
# 这将启动一个简单的 HTTP 服务器作为后端服务
go run http/simple/server/app/*
```

### 第三步：启动 Pixiu 网关

回到**终端 1** 并使用以下命令启动 Pixiu。请将 `[absolute-path]` 替换为您本地 `dubbo-go-pixiu-samples` 目录的绝对路径。

```shell
go run cmd/pixiu/*.go gateway start -c /[absolute-path]/dubbo-go-pixiu-samples/http/simple/pixiu/conf.yaml
```

当您看到类似以下的日志时，表示 Pixiu 已成功启动并正在监听 `8888` 端口：

```log
2025-05-19T12:46:00.104+0800    INFO   server/pixiu_start.go:127  [dubbo-go-pixiu] start by config : &{StaticResources:{Listeners:[0xc0007b7a20] Clusters:[0xc0007cc5a0] Adapters:[] ShutdownConfig:0xc00067fb30 PprofConf:{Enable:false Address:{SocketAddress:{Address:0.0.0.0 Port:8881 ResolverName: Domains:[] CertsDir:} Name:}}} DynamicResources:<nil> Metric:{Enable:false PrometheusPort:0} Node:<nil> Trace:<nil> Wasm:<nil> Config:<nil> Nacos:<nil> Log:<nil>}
2025-05-19T12:46:00.104+0800    INFO   healthcheck/healthcheck.go:157 [health check] create a health check session for 127.0.0.1:1314
2025-05-19T12:46:00.105+0800    INFO   tracing/driver.go:76   [dubbo-go-pixiu] no trace configuration in conf.yaml
2025-05-19T12:46:00.105+0800    INFO   http/http_listener.go:157  [dubbo-go-server] httpListener start at : 0.0.0.0:8888
```

### 第四步：发送测试请求

使用 `curl` 或提供的测试代码来测试网关：

```shell
# 方式一：运行测试用例
go test -v ./http/simple/test/

# 方式二：运行基于curl的测试脚本
./http/simple/request.sh
```

更多使用示例见[dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu-samples)。

## 使用 Docker 部署

我们也提供 Docker 镜像，以便快速、轻松地进行部署。

**1. 从源代码构建 Docker 镜像**

首先，请确保您的机器上已经安装了 Docker。然后，在项目根目录下（即 `Dockerfile` 所在的目录），运行以下命令来构建镜像：

```shell
# 您可以自定义镜像的名称和标签，这里我们使用 dubbo-go-pixiu:local
docker build -t dubbo-go-pixiu:local .
````

构建过程可能需要几分钟时间。成功后，您就可以在本地使用这个名为 `dubbo-go-pixiu:local` 的镜像了。

**2. 使用默认配置运行 Pixiu**

使用您刚刚构建的本地镜像来启动一个容器。

```shell
docker run --name pixiu-gateway -p 8888:8888 -d dubbo-go-pixiu:local
```

**3. 挂载自定义配置文件运行**

如果您需要使用自己的配置文件，可以将本地文件挂载到容器的 `/etc/pixiu/` 目录下。

```shell
# 确保使用您本地构建的镜像名称，例如 dubbo-go-pixiu:local
docker run --name pixiu-gateway -p 8888:8888 -d \
    -v /your/local/path/conf.yaml:/etc/pixiu/conf.yaml \
    -v /your/local/path/log.yml:/etc/pixiu/log.yml \
    dubbo-go-pixiu:local
```

更多信息，请访问 [Pixiu Docker Hub](https://hub.docker.com/r/dubbogopixiu/dubbo-go-pixiu)。

## 可视化控制面：Pixiu Admin

强大的 Pixiu 管理控制台 `pixiu-admin`，已被[合并](https://github.com/dubbo-go-pixiu/pixiu-admin)至本仓库，可以用于可视化配置服务发现、流量管理和安全策略。

**使用 Docker Compose 快速启动：**

```shell
cd /[absolute-path]/dubbo-go-pixiu
docker-compose up -d
```

启动后，在浏览器中访问 `http://localhost:8080` 即可进入管理界面。

![](./docs/images/pixiu-admin.png)

## Tools

* **[Benchmark](./tools/benchmark)**: dubbo-fo-pixiu 的 benchmark 代码, 目前处于过时状态, 我们正在更新该项目.

## 社区与贡献

我们热烈欢迎任何形式的贡献！无论是提交 Issue、提出新功能建议还是贡献代码，您的参与对项目都至关重要。

* **加入我们的社区**：

通过钉钉、微信或 Discord 加入我们的讨论组。

discord https://discord.gg/C5ywvytg
![invite.png](./docs/images/invite.png)

如果您喜欢 Dubbo-Go-Pixiu，请在 GitHub 上给我们一个 ⭐！

## 许可证

本项目基于 [Apache License, Version 2.0](LICENSE) 许可证。