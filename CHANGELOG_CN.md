# Release Notes

---

## 1.1.0

本版本在 AI Gateway 与大模型接入能力 方面实现了重要增强，引入了完整的 LLM Proxy、Token 计量、重试与降级机制，并支持 HTTP / SSE 流式通信。
对网关与 Ingress 架构 进行了全面重构，升级路由机制，构建更加现代化的应用网关模型，更好适配云原生场景。
新增 Model Context Protocol（MCP） 支持，为 AI 服务提供标准化的接入、注册与鉴权能力。
通过深度集成 Nacos，强化了服务发现、动态配置、指标采集和健康检查等服务治理能力。
同时完成了大量 重构、清理与 CI / 构建升级，显著提升了系统稳定性、可维护性和长期演进能力。

### 新特性（New Features）

#### AI Gateway / 大模型（LLM）集成

* 新增模型请求 **Token 计费 / 统计 Filter** [#659](https://github.com/apache/dubbo-go-pixiu/pull/659)
* LLM 代理支持 **重试机制与可配置策略抽象** [#692](https://github.com/apache/dubbo-go-pixiu/pull/692)
* 完整的 LLM 代理 Filter，支持 **重试与失败回退** [#685](https://github.com/apache/dubbo-go-pixiu/pull/685)
* 面向模型服务场景的 **HTTP / SSE 流式推理支持** [#657](https://github.com/apache/dubbo-go-pixiu/pull/657)
* 支持 **可流式 HTTP Streamable HTTP**，适用于长连接数据管道 [#674](https://github.com/apache/dubbo-go-pixiu/pull/674)
* Tokenizer Filter 支持 **Content-Encoding** [#706](https://github.com/apache/dubbo-go-pixiu/pull/706)
* 新增基于 **Nacos 的 LLM 注册中心支持** [#746](https://github.com/apache/dubbo-go-pixiu/pull/746)
* 增强 Upstream 追踪能力与指标采集 [#733](https://github.com/apache/dubbo-go-pixiu/pull/733)
* 改进 API Key 处理逻辑与 Endpoint 健康检查机制 [#731](https://github.com/apache/dubbo-go-pixiu/pull/731)
* 新增 `LLMMeta` 字段，简化 LLM Endpoint 配置 [#678](https://github.com/apache/dubbo-go-pixiu/pull/678)
* 移除静态配置 Provider，全面转向 **动态治理能力** [#764](https://github.com/apache/dubbo-go-pixiu/pull/764)

#### Model Context Protocol（MCP）

* MCP Server Filter 实现 [#702](https://github.com/apache/dubbo-go-pixiu/pull/702)
* MCP Server 与 **Nacos 注册中心集成** [#757](https://github.com/apache/dubbo-go-pixiu/pull/757)
* 增强 MCP 场景下的 **可流式 HTTP 支持** [#769](https://github.com/apache/dubbo-go-pixiu/pull/769)
* MCP 鉴权能力支持 [#740](https://github.com/apache/dubbo-go-pixiu/pull/740)

#### 服务发现与配置（Service Discovery & Configuration）

* 支持 **Nacos 服务发现** [#651](https://github.com/apache/dubbo-go-pixiu/pull/651)
* Logger 支持配置中心监听与 **热更新** [#647](https://github.com/apache/dubbo-go-pixiu/pull/647)
* 启动时从 Nacos 拉取 Logger 配置 [#640](https://github.com/apache/dubbo-go-pixiu/pull/640)
* 支持从注册中心 **动态生成 Router 与 Cluster** [#632](https://github.com/apache/dubbo-go-pixiu/pull/632)
* 修复远程 Nacos 配置字段缺失问题 [#679](https://github.com/apache/dubbo-go-pixiu/pull/679)

#### 代理核心与网络能力（Proxy Core & Networking）

* 新增 **加权随机负载均衡算法** [#677](https://github.com/apache/dubbo-go-pixiu/pull/677)
* 支持 TCP / HTTP / HTTPS 健康检查，并修复 domain 字段问题 [#668](https://github.com/apache/dubbo-go-pixiu/pull/668)
* HTTP Proxy Filter 新增 `scheme` 字段，支持 HTTPS Upstream [#671](https://github.com/apache/dubbo-go-pixiu/pull/671)
* Dubbo 调用支持 **可配置重试次数** [#625](https://github.com/apache/dubbo-go-pixiu/pull/625)
* Dubbo Proxy 支持多种负载均衡策略配置 [#613](https://github.com/apache/dubbo-go-pixiu/pull/613)，[#614](https://github.com/apache/dubbo-go-pixiu/pull/614)，[#615](https://github.com/apache/dubbo-go-pixiu/pull/615)
* Streamable HTTP 与 SSE 处理能力增强 [#657](https://github.com/apache/dubbo-go-pixiu/pull/657)，[#676](https://github.com/apache/dubbo-go-pixiu/pull/676)，[#674](https://github.com/apache/dubbo-go-pixiu/pull/674)

#### gRPC & Dubbo

* 完整实现 **gRPC Streaming 全链路代理**，并进行性能优化 [#688](https://github.com/apache/dubbo-go-pixiu/pull/688)
* 抽象 HTTP 请求解析器，并实现 Dubbo Resolver [#691](https://github.com/apache/dubbo-go-pixiu/pull/691)

#### 网关 / 路由 / Ingress

* 路由机制重构与升级 [#777](https://github.com/apache/dubbo-go-pixiu/pull/777)
* Application Gateway / Ingress 架构重构为更现代的设计 [#827](https://github.com/apache/dubbo-go-pixiu/pull/827)
* 支持新的 Ingress Controller [#792](https://github.com/apache/dubbo-go-pixiu/pull/792)
* Ingress Controller 新增 Application Gateway 资源策略 [#839](https://github.com/apache/dubbo-go-pixiu/pull/839)
* 重构 Ingress 为更现代的 Application Gateway [#827](https://github.com/apache/dubbo-go-pixiu/pull/827)

#### 工具 / 扩展性（Tools / Extensibility）

* Benchmark 工具能力增强 [#807](https://github.com/apache/dubbo-go-pixiu/pull/807)
* 支持 **Open Policy Agent(OPA) HTTP Filter** [#732](https://github.com/apache/dubbo-go-pixiu/pull/732)）

### 增强与重构（Enhancements & Refactors）

#### 日志与配置（Logging & Config）

* Logger 模块重构 [#646](https://github.com/apache/dubbo-go-pixiu/pull/646)
* 修复热更新稳定性与配置覆盖问题 [#682](https://github.com/apache/dubbo-go-pixiu/pull/682)，[#765](https://github.com/apache/dubbo-go-pixiu/pull/765)

#### 项目结构与维护（Project Layout & Maintenance）

* `pixiu-admin` 合并至主仓库 [#697](https://github.com/apache/dubbo-go-pixiu/pull/697)
* `configcenter/` 迁移至 `pkg/`，并移除历史结构 [#762](https://github.com/apache/dubbo-go-pixiu/pull/762)
* Benchmark 工具迁移至 `tools/benchmark` [#763](https://github.com/apache/dubbo-go-pixiu/pull/763)
* Pixiu CLI 调整至 `pkg/cmd` [#596](https://github.com/apache/dubbo-go-pixiu/pull/596)

#### 稳定性与内部质量（Resilience & Internal Quality）

* 大规模鲁棒性增强 [#644](https://github.com/apache/dubbo-go-pixiu/pull/644)
* 统一错误码体系与错误处理逻辑 [#809](https://github.com/apache/dubbo-go-pixiu/pull/809)，[#782](https://github.com/apache/dubbo-go-pixiu/pull/782)
* 统一指标采集 Filter [#799](https://github.com/apache/dubbo-go-pixiu/pull/799)
* 修复 Filter 配置拷贝语义，避免指针共享问题 [#815](https://github.com/apache/dubbo-go-pixiu/pull/815)，[#814](https://github.com/apache/dubbo-go-pixiu/pull/814)

#### CI 与构建改进（CI & Build Improvements）

* Go 版本升级至 `1.25`，更新 CI 工作流与 Lint 规则 [#752](https://github.com/apache/dubbo-go-pixiu/pull/752)，[#666](https://github.com/apache/dubbo-go-pixiu/pull/666)
* GolangCI Lint 重构与稳定性提升 [#650](https://github.com/apache/dubbo-go-pixiu/pull/650)，[#734](https://github.com/apache/dubbo-go-pixiu/pull/734)
* Pipeline 清理与无用 GitHub Action 移除 [#775](https://github.com/apache/dubbo-go-pixiu/pull/775)，[#786](https://github.com/apache/dubbo-go-pixiu/pull/786)
* Docker 构建优化 [#714](https://github.com/apache/dubbo-go-pixiu/pull/714)，[#723](https://github.com/apache/dubbo-go-pixiu/pull/723)

### Bug 修复（Bug Fixes）

* SSE 流在 `io.EOF` 时未正确关闭 [#676](https://github.com/apache/dubbo-go-pixiu/pull/676)
* HTTP Proxy 连接复用问题 [#578](https://github.com/apache/dubbo-go-pixiu/pull/578)
* Access Log Filter 的空指针问题及非 Unary 响应处理 [#713](https://github.com/apache/dubbo-go-pixiu/pull/713)
* Logger 配置覆盖错误 [#765](https://github.com/apache/dubbo-go-pixiu/pull/765)
* 多处数据竞争修复 [#750](https://github.com/apache/dubbo-go-pixiu/pull/750)，[#789](https://github.com/apache/dubbo-go-pixiu/pull/789)
* Nacos 字段未正确传递 [#679](https://github.com/apache/dubbo-go-pixiu/pull/679)
* Benchmark 逻辑修正与性能测试清理 [#819](https://github.com/apache/dubbo-go-pixiu/pull/819)

### 文档（Documentation）

* README 更新与使用指引增强 [#698](https://github.com/apache/dubbo-go-pixiu/pull/698)，[#794](https://github.com/apache/dubbo-go-pixiu/pull/794)，[#831](https://github.com/apache/dubbo-go-pixiu/pull/831)
* 管理后台文档全面重写 [#817](https://github.com/apache/dubbo-go-pixiu/pull/817)
* 新增 MCP 配置文档 [#770](https://github.com/apache/dubbo-go-pixiu/pull/770)
* OPA HTTP Filter 快速入门文档 [#751](https://github.com/apache/dubbo-go-pixiu/pull/751)
* 中文 README 改进 [#641](https://github.com/apache/dubbo-go-pixiu/pull/641)
* Issue 模板清理与优化 [#736](https://github.com/apache/dubbo-go-pixiu/pull/736)，[#735](https://github.com/apache/dubbo-go-pixiu/pull/735)

* dubbo-go 升级至最新版  [#630](https://github.com/apache/dubbo-go-pixiu/pull/630),
  [#807](https://github.com/apache/dubbo-go-pixiu/pull/807),
  [#836](https://github.com/apache/dubbo-go-pixiu/pull/836),
  [#845](https://github.com/apache/dubbo-go-pixiu/pull/845))
* 移除未使用的 Seata Proxy [#628](https://github.com/apache/dubbo-go-pixiu/pull/628)
* 移除 Istio 集成 [#622](https://github.com/apache/dubbo-go-pixiu/pull/622)
* 移除静态配置 Provider [#764](https://github.com/apache/dubbo-go-pixiu/pull/764)

### 贡献者（Contributors）

特别感谢所有为 `dubbo-go-pixiu` 做出贡献的社区成员（按字典序）：

@1kasa
@Alanxtl
@baerwang
@Chen-BUPT
@FoghostCn
@KamToHung
@ma642
@mark4z
@marsevilspirit
@mfordjody
@mutezebra
@nanjiek
@No-SilverBullet
@PhilYue
@Similarityoung
@testwill
@yuluo-yx

---

## 1.0.0

### 新特性（New Features）

* [失败注入（fail inject）](https://github.com/apache/dubbo-go-pixiu/pull/571)
* [新增基于 Header 的路由支持](https://github.com/apache/dubbo-go-pixiu/pull/565)
* [新增 Maglev 哈希负载均衡算法](https://github.com/apache/dubbo-go-pixiu/pull/554)
* [Triple 代理支持导入 protosets](https://github.com/apache/dubbo-go-pixiu/pull/548)
* [为 Windows 新增优雅关闭（GracefulShutdown）信号支持](https://github.com/apache/dubbo-go-pixiu/pull/522)
* [支持 Dubbo 调用链路追踪（Tracing）](https://github.com/apache/dubbo-go-pixiu/pull/559)

### 功能增强（Enhancement）

* [重构 Prometheus 指标（metric）实现](https://github.com/apache/dubbo-go-pixiu/pull/573)
* [移除未使用的包导入](https://github.com/apache/dubbo-go-pixiu/pull/574)
* [整理：避免不必要地使用 fmt.Sprintf](https://github.com/apache/dubbo-go-pixiu/pull/575)
* [整理：WASM filter 使用 build tags，新增 wasm 标记](https://github.com/apache/dubbo-go-pixiu/pull/567)
* [文档：格式化并修改 samples 链接](https://github.com/apache/dubbo-go-pixiu/pull/556)
* [将 gatewayCmd 回退为 Run dubbo-go-pixiu](https://github.com/apache/dubbo-go-pixiu/pull/557)
* [统一 import 格式](https://github.com/apache/dubbo-go-pixiu/pull/527)
* [升级 hessian2 至 v1.11.3](https://github.com/apache/dubbo-go-pixiu/pull/516)

### Bug 修复（Bugfixes）

* [修复注册哈希、数组越界以及哈希初始化问题](https://github.com/apache/dubbo-go-pixiu/pull/530)
* [优化超时状态码处理](https://github.com/apache/dubbo-go-pixiu/pull/521)
* [优化指标（Metric）实现](https://github.com/apache/dubbo-go-pixiu/pull/528)
* [新增并修改 Nacos 配置参数](https://github.com/apache/dubbo-go-pixiu/pull/524)
* [修复 filter 配置为 nil 时的空指针异常（NPE）](https://github.com/apache/dubbo-go-pixiu/pull/517)
* [使用兼容 Mac ARM 的 wasmer-go v1.0.4](https://github.com/apache/dubbo-go-pixiu/pull/515)
* [修复 sample URL，使用 github.com/apache/dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu/pull/506)
* [流量过滤器：修复权重策略及 Apply 方法内的错误处理](https://github.com/apache/dubbo-go-pixiu/pull/507)
* [修复 httpfilter 在多个 URL 之间存在空格时负载均衡失效的问题](https://github.com/apache/dubbo-go-pixiu/pull/513)

Milestone:
[https://github.com/apache/dubbo-go-pixiu/milestone/8](https://github.com/apache/dubbo-go-pixiu/milestone/8)

---

## 0.6.0

### 新特性（New Features）

* [Nacos 配置支持](https://github.com/apache/dubbo-go-pixiu/pull/497)
* [OSPP：流量分发（Traffic Distribution）](https://github.com/apache/dubbo-go-pixiu/pull/501)
* [新增优雅关闭（Graceful Shutdown）](https://github.com/apache/dubbo-go-pixiu/pull/474)
* [Pixiu 的 WASM 插件支持](https://github.com/apache/dubbo-go-pixiu/pull/469)
* [在 k8s + Istio 中将 Pixiu 部署为 Dubbo 服务的出口网关](https://github.com/apache/dubbo-go-pixiu/pull/446)
* [ASoC 2022：Pixiu 指标（Metrics）实现](https://github.com/apache/dubbo-go-pixiu/pull/480)
* [OSPP：流量相关特性](https://github.com/apache/dubbo-go-pixiu/pull/496)
* [特性：一致性哈希](https://github.com/apache/dubbo-go-pixiu/pull/436)

### 功能增强（Enhancement）

* [移除 Http → Dubbo 代理中的 “Types”】【[https://github.com/apache/dubbo-go-pixiu/pull/456】](https://github.com/apache/dubbo-go-pixiu/pull/456】)
* [ASoC 2022：Pixiu 超时特性优化](https://github.com/apache/dubbo-go-pixiu/pull/475)

### Bug 修复（Bugfixes）

* [修复响应头 Content-Type 问题](https://github.com/apache/dubbo-go-pixiu/pull/462)
* [修复 listener session 异常](https://github.com/apache/dubbo-go-pixiu/pull/458)

Milestone:
[https://github.com/apache/dubbo-go-pixiu/milestone/7](https://github.com/apache/dubbo-go-pixiu/milestone/7)

---

## 0.5.1

### 新特性（New Features）

* [链路追踪（Trace）支持](https://github.com/apache/dubbo-go-pixiu/pull/394)
* [健康检查（Health Check）支持](https://github.com/apache/dubbo-go-pixiu/pull/421)
* [xDS 配置支持](https://github.com/apache/dubbo-go-pixiu/pull/385)
* [LDS 支持](https://github.com/apache/dubbo-go-pixiu/pull/417)
* [直接 Dubbo 调用](https://github.com/apache/dubbo-go-pixiu/pull/434)

### 功能增强（Enhancement）

* [Spring Cloud 订阅策略](https://github.com/apache/dubbo-go-pixiu/pull/425)
* [风格优化：路由前缀匹配定义](https://github.com/apache/dubbo-go-pixiu/pull/451)

### Bug 修复（Bugfixes）

* [Nacos 注册中心 Bug 修复](https://github.com/apache/dubbo-go-pixiu/pull/389)
* [修复 Spring Cloud 错误并重构事件回调](https://github.com/apache/dubbo-go-pixiu/pull/367)
* [修复使用 Nacos 注册中心时首次调用失败的问题](https://github.com/apache/dubbo-go-pixiu/pull/380)

Milestone:
[https://github.com/apache/dubbo-go-pixiu/milestone/6](https://github.com/apache/dubbo-go-pixiu/milestone/6)

---

## 0.5.0

### 新特性（New Features）

* [Dubbo ↔ HTTP / Dubbo ↔ Triple / Triple ↔ Dubbo 代理](https://github.com/apache/dubbo-go-pixiu/pull/347)
* [HTTP → Triple 代理](https://github.com/apache/dubbo-go-pixiu/pull/302)
* [HTTP → Dubbo 默认映射规则](https://github.com/apache/dubbo-go-pixiu/pull/298)
* [gRPC 代理](https://github.com/apache/dubbo-go-pixiu/pull/315)
* [从 Spring Cloud ZK 注册中心动态获取集群和路由配置](https://github.com/apache/dubbo-go-pixiu/pull/367)
* [JWT 认证过滤器](https://github.com/apache/dubbo-go-pixiu/pull/303)
* [HTTPS 支持多证书](https://github.com/apache/dubbo-go-pixiu/pull/292)
* [支持构建 Docker 镜像](https://github.com/apache/dubbo-go-pixiu/pull/370)

### 功能增强（Enhancement）

* [为 gRPC 代理新增 HTTP/2 Listener](https://github.com/apache/dubbo-go-pixiu/pull/315)
* [使用 Trie 实现路由](https://github.com/apache/dubbo-go-pixiu/pull/310)
* [HTTP → gRPC 使用 gRPC 反射服务](https://github.com/apache/dubbo-go-pixiu/pull/317)
* [在容器中获取 CPU 核数](https://github.com/apache/dubbo-go-pixiu/pull/340)
* [过滤器链重构](https://github.com/apache/dubbo-go-pixiu/pull/307)
* [升级 hessian2 至 v1.11.0](https://github.com/apache/dubbo-go-pixiu/pull/352)
* [升级 dubbogo 版本至 3.0](https://github.com/apache/dubbo-go-pixiu/pull/334)
* [升级 keyfunc 至稳定版本 v1.0.0](https://github.com/apache/dubbo-go-pixiu/pull/318)

### Bug 修复（Bugfixes）

* [修复使用 HTTP/2 Manager 处理 gRPC 请求时的写入错误](https://github.com/apache/dubbo-go-pixiu/pull/372)
* [修复 Spring Cloud 错误并重构事件回调](https://github.com/apache/dubbo-go-pixiu/pull/367)
* [修复使用 Nacos 注册中心时首次调用失败的问题](https://github.com/apache/dubbo-go-pixiu/pull/380)

---

## 0.4.0

### 新特性（New Features）

* [从 Spring Cloud Nacos 注册中心动态获取集群和路由配置](https://github.com/apache/dubbo-go-pixiu/pull/255)
* [从 Zookeeper 注册中心动态获取 Dubbo 代理配置](https://github.com/apache/dubbo-go-pixiu/pull/256)
* [HTTP → gRPC 代理](https://github.com/apache/dubbo-go-pixiu/pull/244)
* [HTTP → HTTP 代理](https://github.com/apache/dubbo-go-pixiu/pull/242)
* [集成 Jaeger 的链路追踪支持](https://github.com/apache/dubbo-go-pixiu/pull/236)
* [CORS（跨域资源共享）策略支持](https://github.com/apache/dubbo-go-pixiu/pull/249)

### 功能增强（Enhancement）

* [新增更多示例（samples）](https://github.com/apache/dubbo-go-pixiu/pull/271)
* [使用 Cobra 命令行工具](https://github.com/apache/dubbo-go-pixiu/pull/234)
* [新增示例快速启动脚本](https://github.com/apache/dubbo-go-pixiu/pull/226)
* [升级 hessian2 至 v1.9.3](https://github.com/apache/dubbo-go-pixiu/pull/248)
* [将 onAir 属性重命名为 enable](https://github.com/apache/dubbo-go-pixiu/pull/243)
* [链路追踪优化](https://github.com/apache/dubbo-go-pixiu/pull/257/files)
* [支持 HTTPS](https://github.com/apache/dubbo-go-pixiu/pull/213)

### Bug 修复（Bugfixes）

* [修复请求体丢失问题](https://github.com/apache/dubbo-go-pixiu/pull/260)
* [修复 HttpContext 重置 Bug](https://github.com/apache/dubbo-go-pixiu/pull/254)
* [修复环境变量无法设置的问题](https://github.com/apache/dubbo-go-pixiu/pull/239)
* [修复 filterManager 获取过滤器顺序随机的问题](https://github.com/apache/dubbo-go-pixiu/pull/264)
* [修复超时过滤器的空指针问题](https://github.com/apache/dubbo-go-pixiu/pull/278)

Milestone：
[https://github.com/apache/dubbo-go-pixiu/milestone/5](https://github.com/apache/dubbo-go-pixiu/milestone/5)

---

## 0.3.0

### 新特性（New Features）

* [限流过滤器（Rate Limit Filter）](https://github.com/apache/dubbo-go-pixiu/pull/169)
* [新增集成测试（Integration Test）](https://github.com/apache/dubbo-go-pixiu/pull/183)
* [支持处理限流配置更新事件](https://github.com/apache/dubbo-go-pixiu/pull/196)
* [在 Pixiu 中新增 OTEL 指标导出到 Prometheus](https://github.com/apache/dubbo-go-pixiu/pull/204)
* [使 Pixiu Admin 的配置管理更加细粒度](https://github.com/apache/dubbo-go-pixiu/pull/171)

### 功能增强（Enhancement）

* [更新 samples/admin 示例](https://github.com/apache/dubbo-go-pixiu/pull/208)
* [更新限流示例（ratelimit samples）](https://github.com/apache/dubbo-go-pixiu/pull/206)
* [使路由匹配区分大小写](https://github.com/apache/dubbo-go-pixiu/pull/209)
* [新增更多测试用例](https://github.com/apache/dubbo-go-pixiu/pull/203)
* [丰富过滤器测试用例](https://github.com/apache/dubbo-go-pixiu/pull/202)
* [丰富 response.go 的测试用例](https://github.com/apache/dubbo-go-pixiu/pull/197)

### Bug 修复（Bugfixes）

* [修复 CI 检查状态与要求不匹配的问题](https://github.com/apache/dubbo-go-pixiu/pull/199)
* [修复超时配置被覆盖的问题](https://github.com/apache/dubbo-go-pixiu/pull/190)
* [修复 Quick Start 相关问题](https://github.com/apache/dubbo-go-pixiu/pull/191)
* [修复 Bug：无法通过路径删除节点](https://github.com/apache/dubbo-go-pixiu/pull/201)
* [修复流程图问题](https://github.com/apache/dubbo-go-pixiu/pull/205)
* [修复 reviewdog 问题](https://github.com/apache/dubbo-go-pixiu/pull/195)

Milestone：
[https://github.com/apache/dubbo-go-pixiu/milestone/4](https://github.com/apache/dubbo-go-pixiu/milestone/4)

---

## 0.2.1

### 功能增强（Enhancement）

* [更换 Pixiu 吉祥物](https://github.com/apache/dubbo-go-pixiu/pull/178)
* [reviewdog 使用默认参数](https://github.com/apache/dubbo-go-pixiu/pull/167)
* [将参数类型移动到参数配置中，而不是作为独立定义](https://github.com/apache/dubbo-go-pixiu/pull/161)
* [修复版本字段（version field）](https://github.com/apache/dubbo-go-pixiu/pull/166)
* [新增 license-eye 用于检查并修复 License Header](https://github.com/apache/dubbo-go-pixiu/pull/164)
* [优化：扩展 filterFuncCacheMap 的初始长度](https://github.com/apache/dubbo-go-pixiu/pull/174)
* [重构 config_load.go](https://github.com/apache/dubbo-go-pixiu/pull/158)

Milestone：
[https://github.com/apache/dubbo-go-pixiu/milestone/3](https://github.com/apache/dubbo-go-pixiu/milestone/3)

---

## 0.2.0

### 新特性（New Features）

* [新增 dubbo-go-proxy 管理端（Admin）](https://github.com/dubbogo/dubbo-go-proxy/pull/115)
* [新增插件机制（Plugin）](https://github.com/dubbogo/dubbo-go-proxy/pull/109)

### Bug 修复（Bugfixes）

* [修复：移除 replace-path-filter](https://github.com/dubbogo/dubbo-go-proxy/pull/118)

Milestone：
[https://github.com/dubbogo/dubbo-go-proxy/milestone/2?closed=1](https://github.com/dubbogo/dubbo-go-proxy/milestone/2?closed=1)
