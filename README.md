[![Pixiu Logo](docs/images/pixiu-logo-v4.png)](http://alexstocks.github.io/html/dubbogo.html)

# Dubbo-Go-Pixiu: A Next-Generation, High-Performance AI / API Gateway

[![Build Status](https://github.com/apache/dubbo-go-pixiu/workflows/CI/badge.svg)](https://travis-ci.org/apache/dubbo-go-pixiu)
[![codecov](https://codecov.io/gh/apache/dubbo-go-pixiu/branch/master/graph/badge.svg)](https://codecov.io/gh/apache/dubbo-go-pixiu)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go\&logoColor=white\&style=flat-square)](https://pkg.go.dev/github.com/apache/dubbo-go-pixiu?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/dubbo-go-pixiu)](https://goreportcard.com/report/github.com/apache/dubbo-go-pixiu)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

**English** | [中文](README_CN.md)

---

**Dubbo-Go-Pixiu** is a next-generation **AI / API Gateway** built on **Dubbo-go**, empowering seamless connections to **LLMs** and **MCPs** with unified access, intelligent extensions, and cost-efficient management. At the same time, Pixiu bridges external protocols with internal Dubbo clusters, supporting **HTTP**, **gRPC**, **Dubbo2**, and **Triple** for high-performance, scalable integration.

> ⭐ **New Capability:** Pixiu has evolved into a **universal AI Gateway**, designed to simplify and unify access to **LLMs and AI service providers** — whether from public vendors or self-hosted models.

> ⭐ **New Capability:** Pixiu has also evolved into a **Kubernetes Ingress Controller**, enabling native management of API and AI traffic within Kubernetes clusters through declarative routing and Pixiu’s flexible governance model.

👉 **Try it now:** Explore our official [samples](https://github.com/apache/dubbo-go-pixiu-samples)

## We've evolved into the **Next-Generation AI Gateway**

Pixiu has evolved into a **universal AI Gateway**, designed to simplify and unify access to **LLMs and AI service providers** — whether from public vendors or self-hosted models.

With Pixiu, you can:

* **Unified Access to AI Models:**
  Connect effortlessly to OpenAI, Anthropic, or any **custom / on-prem LLM or MCP** service through a single, consistent API gateway layer.

* **MCP Server Exposure:**
  Expose your existing HTTP APIs and backend services as MCP Servers through Pixiu, enabling AI applications to directly invoke your business logic.

* **Flexible Extension for AI Workloads:**
  Apply authentication, caching, rate-limiting, retry policy, observability, or even model orchestration — all through Pixiu's plugin system.

* **Multi-Tenant & Cost-Efficient:**
  Implement fine-grained cost control, auditing, and token accounting for large-scale AI deployments.

* **Future-Proof Architecture:**
  Designed for the hybrid world of API + AI traffic — bridging traditional microservices with the AI-first era.

> ✨ **Now also a Kubernetes Ingress Controller:**
> Pixiu’s AI gateway abilities can now be applied directly at the cluster ingress layer — allowing traffic governance, routing, and observability to be defined declaratively via Kubernetes APIs.

👉 **Try it now:** Explore our AI Gateway Samples: [LLM](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/llm) and [MCP](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/mcp)


## Why Choose Dubbo-Go-Pixiu?

* **⚡ High Performance:**
  Built with Go for ultra-low latency and high throughput, Pixiu is optimized for large-scale API traffic and LLM workloads.

* **🧩 Highly Extensible:**
  A powerful **plugin and filter framework** makes it easy to extend or customize gateway behavior — from routing and authentication to observability and AI-specific use cases.

* **🎯 Developer Friendly:**
  Simple configuration, unified management plane, and clear observability make Pixiu ideal for both developers and platform engineers.

* **🌩️ Cloud-Native by Design:**
  Works seamlessly in Kubernetes, supports declarative configuration, and integrates easily with Dubbo, Spring Cloud, or custom backends.

* **🔗 Seamless Dubbo Integration:**
  The official sidecar solution for connecting non-Java applications (Go, Python, Node.js, etc.) to Dubbo services.

### Core Features

| Category                   | Description                                                                                                                    |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| 🚀 **Protocol Processing** | Supports proxying and translation between HTTP, gRPC, Dubbo2, and Triple protocols, delivering rich protocol gateway features. |
| 🛡️ **Security & Auth**    | HTTPS, JWT, OAuth2, and other security mechanisms to safeguard your APIs and AI endpoints.                                     |
| 🔍 **Service Discovery**   | Integrates with Zookeeper, Nacos, or any service registry to discover Dubbo and Spring Cloud services automatically.           |
| ⚖️ **Traffic Governance**  | Integrates with Sentinel for fine-grained rate limiting, circuit breaking, and traffic shaping.                                |
| 📈 **Observability**       | OpenTelemetry and Jaeger support for full tracing, metrics, and logging visibility.                                            |
| 🎨 **Visual Management**   | The **[Pixiu-admin](https://github.com/apache/dubbo-go-pixiu/tree/develop/admin)** UI offers real-time visual configuration for gateway rules and policies.                                   |

## Deploying with Docker

We also provide Docker images for quick and easy deployment.

**1. Build Docker Image from Source Code**

Make sure you are in the root directory of the project (where the `Dockerfile` is located), then run the following command to build the image:

```shell
# You can customize the image name and tag, here we use dubbo-go-pixiu:local
docker build -t dubbo-go-pixiu:local .
````

It may take a few minutes to build. Once successful, you can use the local image named `dubbo-go-pixiu:local`.

**2. Run Pixiu with Default Configuration**

Start a container using the local image you just built:

```shell
docker run --name pixiu-gateway -p 8888:8888 -d dubbo-go-pixiu:local
```

**3. Run with Custom Configuration File**

If you need to use your own configuration file, you can mount a local file to the container's `/etc/pixiu/` directory.

```shell
# make sure to use the local image name you built, i.e., dubbo-go-pixiu:local
docker run --name pixiu-gateway -p 8888:8888 -d \
    -v /your/local/path/conf.yaml:/etc/pixiu/conf.yaml \
    -v /your/local/path/log.yml:/etc/pixiu/log.yml \
    dubbo-go-pixiu:local
```


## Pixiu Admin – Visual Control Plane

Manage traffic, routing via `pixiu-admin`.
Start instantly with Docker Compose:

```bash
docker-compose up -d
```

👉 Access admin UI at: [http://localhost:8080](http://localhost:8080)

![pixiu-admin.png](./docs/images/pixiu-admin.png)


## Other Projects in the Dubbo-Go-Pixiu Ecosystem

* **[dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu-samples)** — Official sample repository demonstrating various use cases
* **[pixiu-admin](https://github.com/apache/dubbo-go-pixiu/tree/develop/admin)** — Visual management plane for configuration & monitoring
* **[pixiu-api](https://github.com/dubbo-go-pixiu/pixiu-api)** — API / model definitions for Pixiu Admin
* **[benchmark](https://github.com/apache/dubbo-go-pixiu/tree/develop/tools/benchmark)** — Benchmarking suite for pixiu

## Community & Contribution

We warmly welcome all forms of contributions\! Whether it's submitting an issue, proposing a new feature, or contributing code, your participation is vital to the project.

* **Contribution Workflow:**
  To submit a Pull Request, please submit it to the [dubbo-go-pixiu/dubbo-go-pixiu](https://github.com/dubbo-go-pixiu/dubbo-go-pixiu/) repository. Your code will undergo automated review and manual verification by project maintainers, and will be automatically synchronized to the Apache official repository upon approval.

* **Join Our Community**:

  Join our discussion group through Ding talk, WeChat, or Discord.

  discord https://discord.gg/C5ywvytg
  ![invite.png](./docs/images/invite.png)


If you like Dubbo-Go-Pixiu, please ⭐ star us on GitHub!

## License

This project is licensed under the [Apache License, Version 2.0](LICENSE).