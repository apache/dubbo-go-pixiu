[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)

# Dubbo-Go-Pixiu：新一代高性能 AI / API 网关

[![构建状态](https://github.com/apache/dubbo-go-pixiu/workflows/CI/badge.svg)](https://travis-ci.org/apache/dubbo-go-pixiu)
[![代码覆盖率](https://codecov.io/gh/apache/dubbo-go-pixiu/branch/master/graph/badge.svg)](https://codecov.io/gh/apache/dubbo-go-pixiu)
[![go.dev 参考](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go\&logoColor=white\&style=flat-square)](https://pkg.go.dev/github.com/apache/dubbo-go-pixiu?tab=doc)
[![Go 报告卡](https://goreportcard.com/badge/github.com/apache/dubbo-go-pixiu)](https://goreportcard.com/report/github.com/apache/dubbo-go-pixiu)
![许可证](https://img.shields.io/badge/license-Apache--2.0-green.svg)

[English](README.md) | **中文**

---

**Dubbo-Go-Pixiu** 是一个基于 **Dubbo-go** 构建的新一代 **AI / API 网关**，能够无缝连接 **LLM** 与 **MCP**，提供统一接入、智能扩展与高效的成本管理。同时，Pixiu 还可桥接外部协议与内部 Dubbo 集群，支持 **HTTP**、**gRPC**、**Dubbo2** 与 **Triple**，实现高性能、可扩展的集成能力。

👉 **立即体验：** 访问我们的官方 [示例项目](https://github.com/apache/dubbo-go-pixiu-samples)

## 我们已进化为「下一代 AI 网关」

Pixiu 现已演进为通用型 **AI 网关**，旨在简化并统一访问各种 **LLM 与 AI 服务提供商**——无论是公共云厂商还是自建模型。

通过 Pixiu，你可以：

* **统一访问 AI 模型：**
  通过单一一致的 API 网关层，轻松接入 OpenAI、Anthropic 或任何自定义 / 私有部署的 LLM 或 MCP 服务。

* **MCP 服务暴露：**
  通过 Pixiu 将现有的 HTTP API 等后端服务封装并暴露为 MCP Server，让 AI 应用能够直接调用你的业务逻辑。

* **灵活扩展 AI 工作负载：**
  通过 Pixiu 的插件系统，为请求添加认证、缓存、限流、重试、可观测性，甚至模型编排等功能。

* **多租户与成本高效：**
  实现细粒度的成本控制、审计与 Token 计费，适用于大规模 AI 部署。

* **面向未来的架构：**
  为「API + AI 混合流量」时代而设计，轻松衔接传统微服务与 AI 优先架构。

👉 **立即尝试：** 查看我们的 AI 网关示例：[LLM 示例](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/llm) 与 [MCP 示例](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/mcp)

## 为什么选择 Dubbo-Go-Pixiu？

* **⚡ 高性能：**
  使用 Go 语言构建，具备极低延迟与高吞吐量，专为大规模 API 流量与 LLM 工作负载优化。

* **🧩 高可扩展性：**
  强大的 **插件与过滤器框架**，可轻松扩展或定制网关行为——从路由、认证到可观测性与 AI 场景支持。

* **🎯 开发者友好：**
  简洁配置、统一管理平面与清晰的可观测性，使其成为开发者与平台工程师的理想选择。

* **🌩️ 原生云设计：**
  完全兼容 Kubernetes，支持声明式配置，轻松集成 Dubbo、Spring Cloud 或自定义后端。

* **🔗 无缝 Dubbo 集成：**
  官方推荐的 Sidecar 解决方案，用于将非 Java 应用（如 Go、Python、Node.js 等）连接至 Dubbo 服务。

### 核心特性

| 分类            | 描述                                                     |
| ------------- | ------------------------------------------------------ |
| 🚀 **协议处理**   | 支持 HTTP、gRPC、Dubbo2、Triple 之间的代理与转换，提供丰富的协议网关能力。       |
| 🛡️ **安全与认证** | 支持 HTTPS、JWT、OAuth2 等多种安全机制，保护你的 API 与 AI 端点。          |
| 🔍 **服务发现**   | 集成 Zookeeper、Nacos 等注册中心，自动发现 Dubbo 与 Spring Cloud 服务。 |
| ⚖️ **流量治理**   | 集成 Sentinel，支持限流、熔断与流量整形等精细化治理。                        |
| 📈 **可观测性**   | 内置 OpenTelemetry 与 Jaeger 支持，实现全链路追踪、指标与日志可视化。         |
| 🎨 **可视化管理**  | 提供 **Pixiu-Admin** UI，实现网关规则与策略的实时可视化配置。               |

## 使用 Docker 部署

Pixiu 提供官方 Docker 镜像，便于快速部署。

**1. 从源码构建 Docker 镜像**

请确保当前目录为项目根目录（包含 `Dockerfile` 文件），执行以下命令：

```bash
# 可自定义镜像名称与标签，这里使用 dubbo-go-pixiu:local
docker build -t dubbo-go-pixiu:local .
```

构建过程可能需要几分钟。成功后，本地将生成名为 `dubbo-go-pixiu:local` 的镜像。

**2. 使用默认配置运行 Pixiu**

```bash
docker run --name pixiu-gateway -p 8888:8888 -d dubbo-go-pixiu:local
```

**3. 使用自定义配置文件运行**

如需使用自定义配置文件，可将本地文件挂载至容器的 `/etc/pixiu/` 目录：

```bash
docker run --name pixiu-gateway -p 8888:8888 -d \
    -v /your/local/path/conf.yaml:/etc/pixiu/conf.yaml \
    -v /your/local/path/log.yml:/etc/pixiu/log.yml \
    dubbo-go-pixiu:local
```

更多信息请访问 [Pixiu Docker Hub](https://hub.docker.com/r/dubbogopixiu/dubbo-go-pixiu)。

## Pixiu Admin – 可视化控制平面

通过 `pixiu-admin` 可视化管理流量、路由与 AI 模型接入。

使用 Docker Compose 一键启动：

```bash
docker-compose up -d
```

启动后，在浏览器中访问 `http://localhost:8080` 即可打开管理面板。

![pixiu-admin.png](./docs/images/pixiu-admin.png)

👉 **更多细节：** 请访问 [pixiu-admin](https://github.com/apache/dubbo-go-pixiu/tree/develop/admin)

## Dubbo-go-pixiu 生态系统的其他项目

-   **[pixiu-samples](https://github.com/apache/dubbo-go-pixiu-samples)**
Pixiu-samples 提供了多个使用 dubbo-go-pixiu 的示例项目，涵盖了从基本的 API 网关功能到复杂的 AI 网关场景，帮助用户快速上手和理解 pixiu 的各种功能。
-   **[pixiu-admin](https://github.com/apache/dubbo-go-pixiu/tree/develop/admin)** Dubbo-go-pixiu Admin 是 dubbo-go-pixiu 网关的综合管理平台。它提供了一个集中的控制面板，用于通过基于 Web 的用户界面和 RESTful API 来配置、监控和管理网关资源。
-   **[pixiu-api](https://github.com/dubbo-go-pixiu/pixiu-api)** Dubbo-go-pixiu API 是 dubbo-go-pixiu 生态系统的 API 模型。用于与 pixiu-admin 的集成。
-   **[benchmark](https://github.com/apache/dubbo-go-pixiu/tree/develop/tools/benchmark)** 该基准测试系统允许用户在各种负载条件下测量和分析关键性能指标，如延迟、吞吐量和每秒查询数 (QPS)，以评估协议转换过程的效率。

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