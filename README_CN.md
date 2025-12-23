[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)

# Dubbo-Go-Pixiu: 下一代高性能 AI / API 网关

[![构建状态](https://github.com/apache/dubbo-go-pixiu/workflows/CI/badge.svg)](https://travis-ci.org/apache/dubbo-go-pixiu)
[![代码覆盖率](https://codecov.io/gh/apache/dubbo-go-pixiu/branch/master/graph/badge.svg)](https://codecov.io/gh/apache/dubbo-go-pixiu)
[![go.dev 参考](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go\&logoColor=white\&style=flat-square)](https://pkg.go.dev/github.com/apache/dubbo-go-pixiu?tab=doc)
[![Go 报告卡](https://goreportcard.com/badge/github.com/apache/dubbo-go-pixiu)](https://goreportcard.com/report/github.com/apache/dubbo-go-pixiu)
![许可证](https://img.shields.io/badge/license-Apache--2.0-green.svg)

[English](README.md) | **中文**

---

**Dubbo-Go-Pixiu** 是一个基于 **Dubbo-go** 打造的下一代 **AI / API 网关**，支持无缝连接 **LLMs** 和 **MCPs**，提供统一访问、智能扩展和高效的成本管理。同时，Pixiu 还可以将外部协议与内部 Dubbo 集群进行桥接，支持 **HTTP**、**gRPC**、**Dubbo2** 和 **Triple** 协议，实现高性能、可扩展的集成。

> ⭐ **新能力：** Pixiu 已进化为一个 **通用 AI 网关**，旨在简化并统一访问 **LLMs 和 AI 服务提供商** —— 无论是公有云供应商还是自托管模型。

> ⭐ **新能力：** Pixiu 还进化为 **Kubernetes Ingress Controller**，通过声明式路由和 Pixiu 灵活的治理模型，实现在 Kubernetes 集群内原生管理 API 和 AI 流量。

👉 **立即试用：** 探索我们的官方 [示例](https://github.com/apache/dubbo-go-pixiu-samples)

## 我们已经进化为 **下一代 AI 网关**

Pixiu 已经发展成 **通用 AI 网关**，旨在简化并统一访问 **LLMs 和 AI 服务提供商** —— 无论是来自公共供应商还是自建模型。

使用 Pixiu，您可以：

* **统一接入 AI 模型：**
  通过单一、一致的 API 网关层，轻松连接 OpenAI、Anthropic 或任何 **自定义/本地 LLM 或 MCP** 服务。

* **暴露 MCP 服务器：**
  通过 Pixiu，将现有的 HTTP API 和后端服务作为 MCP 服务器暴露，允许 AI 应用直接调用您的业务逻辑。

* **灵活扩展 AI 工作负载：**
  通过 Pixiu 的插件系统，应用认证、缓存、限流、重试策略、可观测性甚至模型编排等功能。

* **多租户与成本高效：**
  为大规模 AI 部署实施精细化的成本控制、审计和 Token 计费管理。

* **面向未来的架构：**
  设计用于混合的 API + AI 流量 —— 将传统的微服务架构与 AI 首先的时代连接起来。

> ✨ **现在也成为 Kubernetes Ingress Controller：**
> Pixiu 的 AI 网关能力现在可以直接应用于集群入口层 —— 通过 Kubernetes API 声明式定义流量治理、路由和可观测性。

👉 **立即试用：** 探索我们的 AI 网关示例：[LLM](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/llm) 和 [MCP](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/mcp)

## 为什么选择 Dubbo-Go-Pixiu？

* **⚡ 高性能：**
  基于 Go 构建，超低延迟和高吞吐量，优化大规模 API 流量和 LLM 工作负载。

* **🧩 高度可扩展：**
  强大的 **插件和过滤器框架** 使得扩展或自定义网关行为变得容易 —— 从路由、认证到可观测性和 AI 特定用例。

* **🎯 开发者友好：**
  简单配置、统一管理平面和清晰的可观测性，使 Pixiu 成为开发者和平台工程师的理想选择。

* **🌩️ 云原生设计：**
  在 Kubernetes 中无缝运行，支持声明式配置，轻松集成 Dubbo、Spring Cloud 或自定义后端。

* **🔗 无缝对接 Dubbo：**
  官方的 Sidecar 解决方案，连接非 Java 应用（Go、Python、Node.js 等）与 Dubbo 服务。

### 核心特性

| 类别            | 描述                                                         |
| ------------- | ---------------------------------------------------------- |
| 🚀 **协议处理**   | 支持 HTTP、gRPC、Dubbo2 和 Triple 协议的代理和转换，提供丰富的协议网关功能。         |
| 🛡️ **安全与认证** | HTTPS、JWT、OAuth2 等安全机制，保护您的 API 和 AI 服务端点。                 |
| 🔍 **服务发现**   | 集成 Zookeeper、Nacos 或任何服务注册中心，自动发现 Dubbo 和 Spring Cloud 服务。 |
| ⚖️ **流量治理**   | 集成 Sentinel，实现精细化的限流、熔断和流量整形。                              |
| 📈 **可观测性**   | 支持 OpenTelemetry 和 Jaeger 实现全链路追踪、指标和日志可视化。                |
| 🎨 **可视化管理**  | **[Pixiu-admin](https://github.com/apache/dubbo-go-pixiu/tree/develop/admin)** UI 提供实时可视化的网关规则和策略配置管理。                    |

## 使用 Docker 部署

我们还提供了 Docker 镜像，方便快速部署。

**1. 从源码构建 Docker 镜像**

确保您在项目根目录下（包含 `Dockerfile` 文件），然后运行以下命令构建镜像：

```shell
# 可以自定义镜像名称和标签，这里使用 dubbo-go-pixiu:local
docker build -t dubbo-go-pixiu:local .
```

这可能需要几分钟时间来构建。一旦成功，您可以使用名为 `dubbo-go-pixiu:local` 的本地镜像。

**2. 使用默认配置运行 Pixiu**

使用刚才构建的本地镜像启动容器：

```shell
docker run --name pixiu-gateway -p 8888:8888 -d dubbo-go-pixiu:local
```

**3. 使用自定义配置文件运行**

如果您需要使用自定义配置文件，可以将本地文件挂载到容器的 `/etc/pixiu/` 目录中。

```shell
# 确保使用您构建的本地镜像名称，例如 dubbo-go-pixiu:local
docker run --name pixiu-gateway -p 8888:8888 -d \
    -v /your/local/path/conf.yaml:/etc/pixiu/conf.yaml \
    -v /your/local/path/log.yml:/etc/pixiu/log.yml \
    dubbo-go-pixiu:local
```

## Pixiu Admin – 可视化控制平面

通过 `pixiu-admin` 管理流量和路由。
可以通过 Docker Compose 快速启动：

```bash
docker-compose up -d
```

👉 访问管理员 UI：[http://localhost:8080](http://localhost:8080)

![pixiu-admin.png](./docs/images/pixiu-admin.png)

## Dubbo-Go-Pixiu 生态中的其他项目

* **[dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu-samples)** — 官方示例仓库，展示各种用例
* **[pixiu-admin](https://github.com/apache/dubbo-go-pixiu/tree/develop/admin)** — 可视化配置与监控管理平台
* **[pixiu-api](https://github.com/dubbo-go-pixiu/pixiu-api)** — Pixiu 管理面 API / 模型定义
* **[benchmark](https://github.com/apache/dubbo-go-pixiu/tree/develop/tools/benchmark)** — Pixiu 性能测试套件

## 社区与贡献

我们欢迎任何形式的贡献！
无论是提交 Issue、提议新特性，还是贡献代码，你的参与对项目至关重要。

* **代码贡献流程：**
  如果您希望提交 Pull Request，请将 Pull Request 提交至 [dubbo-go-pixiu/dubbo-go-pixiu](https://github.com/dubbo-go-pixiu/dubbo-go-pixiu/) 仓库。代码将经过自动化审查和项目维护者的人工复核，审核通过后自动同步至 Apache 官方仓库。

* **加入社区：**
  可通过钉钉、微信或 Discord 加入我们的讨论群组。

  Discord: [https://discord.gg/C5ywvytg](https://discord.gg/C5ywvytg)
  ![invite.png](./docs/images/invite.png)

如果你喜欢 Dubbo-Go-Pixiu，请在 GitHub 上给我们点个 ⭐！

## 许可证

本项目采用 [Apache License 2.0](LICENSE) 开源许可。