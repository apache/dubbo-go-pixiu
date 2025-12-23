# Release Notes

---

## 1.1.0

This release delivers major improvements in AI Gateway and LLM integration, including a full-featured LLM proxy, token billing, retry/fallback strategies, and HTTP/SSE streaming support.
The gateway and ingress architecture has been significantly refactored, with redesigned routing and a modernized application gateway model for cloud-native environments.
Model Context Protocol (MCP) support is introduced, enabling standardized AI service integration with registry and authorization capabilities.
Service governance and observability are enhanced through deeper Nacos integration, dynamic configuration, unified metrics, and improved health checks.
Numerous refactors, cleanups, and CI/build upgrades improve overall stability, maintainability, and long-term evolution.

### New Features

#### AI Gateway / LLM Integration

* Token billing / counting filter for model requests [#659](https://github.com/apache/dubbo-go-pixiu/pull/659)
* Retry and configurable strategy abstraction for LLM proxy [#692](https://github.com/apache/dubbo-go-pixiu/pull/692)
* Full LLM proxy filter with retry and fallback capabilities [#685](https://github.com/apache/dubbo-go-pixiu/pull/685)
* HTTP/SSE streaming support for model-serving scenarios [#657](https://github.com/apache/dubbo-go-pixiu/pull/657)
* Streamable HTTP support for long-lived data pipes [#674](https://github.com/apache/dubbo-go-pixiu/pull/674)
* Content-encoding support in tokenizer filter [#706](https://github.com/apache/dubbo-go-pixiu/pull/706)
* Nacos-based LLM registry support [#746](https://github.com/apache/dubbo-go-pixiu/pull/746)
* Enhanced upstream tracking and metrics instrumentation [#733](https://github.com/apache/dubbo-go-pixiu/pull/733)
* Improved API key handling and endpoint health check logic [#731](https://github.com/apache/dubbo-go-pixiu/pull/731)
* `LLMMeta` field added to simplify LLM endpoint configuration [#678](https://github.com/apache/dubbo-go-pixiu/pull/678)
* Static configuration providers removed in favor of dynamic governance [#764](https://github.com/apache/dubbo-go-pixiu/pull/764)

#### Model Context Protocol (MCP)

* MCP Server Filter implementation [#702](https://github.com/apache/dubbo-go-pixiu/pull/702)
* MCP server integration with Nacos [#757](https://github.com/apache/dubbo-go-pixiu/pull/757)
* Enhanced streamable HTTP for MCP [#769](https://github.com/apache/dubbo-go-pixiu/pull/769)
* MCP authorization support [#740](https://github.com/apache/dubbo-go-pixiu/pull/740)

#### Service Discovery & Configuration

* Nacos service discovery support [#651](https://github.com/apache/dubbo-go-pixiu/pull/651)
* Config center listening and hot reloading for logger [#647](https://github.com/apache/dubbo-go-pixiu/pull/647)
* Fetching logger configuration from Nacos at startup [#640](https://github.com/apache/dubbo-go-pixiu/pull/640)
* Dynamic router and cluster generation from registry center [#632](https://github.com/apache/dubbo-go-pixiu/pull/632)
* Fix for missing config fields in remote Nacos configuration [#679](https://github.com/apache/dubbo-go-pixiu/pull/679)

#### Proxy Core & Networking

* Weighted random load balancer [#677](https://github.com/apache/dubbo-go-pixiu/pull/677)
* TCP/HTTP/HTTPS health check support and domain field fix [#668](https://github.com/apache/dubbo-go-pixiu/pull/668)
* Scheme field added to HTTP proxy filter allowing HTTPS upstream calls [#671](https://github.com/apache/dubbo-go-pixiu/pull/671)
* Customizable retry count for Dubbo invocations [#625](https://github.com/apache/dubbo-go-pixiu/pull/625)
* Load balancing strategy configuration for Dubbo proxy
  [#613](https://github.com/apache/dubbo-go-pixiu/pull/613),
  [#614](https://github.com/apache/dubbo-go-pixiu/pull/614),
  [#615](https://github.com/apache/dubbo-go-pixiu/pull/615)
* Streamable HTTP and SSE handling improvements
  [#657](https://github.com/apache/dubbo-go-pixiu/pull/657),
  [#676](https://github.com/apache/dubbo-go-pixiu/pull/676),
  [#674](https://github.com/apache/dubbo-go-pixiu/pull/674)

#### gRPC & Dubbo

* Full gRPC streaming proxy implementation and associated optimizations [#688](https://github.com/apache/dubbo-go-pixiu/pull/688)
* Abstracted HTTP request resolver and Dubbo resolver implementation [#691](https://github.com/apache/dubbo-go-pixiu/pull/691)

#### Gateway / Routing / Ingress

* Route mechanism redesign and upgrade [#777](https://github.com/apache/dubbo-go-pixiu/pull/777)
* Application gateway / ingress refactoring to a modern architecture [#827](https://github.com/apache/dubbo-go-pixiu/pull/827)
* Support for a new ingress controller [#792](https://github.com/apache/dubbo-go-pixiu/pull/792)
* Add application gateway resource policy for ingress controller [#839](https://github.com/apache/dubbo-go-pixiu/pull/839)
* Refactor ingress into a more modern application gateway [#827](https://github.com/apache/dubbo-go-pixiu/pull/827)

#### Tools / Extensibility

* Benchmark tool enhancements [#807](https://github.com/apache/dubbo-go-pixiu/pull/807)
* Open Policy Agent (OPA) HTTP filter support [#732](https://github.com/apache/dubbo-go-pixiu/pull/732)

### Enhancements & Refactors

#### Logging & Config

* Logger module refactoring [#646](https://github.com/apache/dubbo-go-pixiu/pull/646)
* Hot-reload stability and overwrite fixes
  [#682](https://github.com/apache/dubbo-go-pixiu/pull/682),
  [#765](https://github.com/apache/dubbo-go-pixiu/pull/765)

#### Project Layout & Maintenance

* `pixiu-admin` migrated into main repository [#697](https://github.com/apache/dubbo-go-pixiu/pull/697)
* `pixiu-api` migrated into main repository [#841](https://github.com/apache/dubbo-go-pixiu/pull/841)
* `configcenter/` migrated to `pkg/` and legacy structures removed [#762](https://github.com/apache/dubbo-go-pixiu/pull/762)
* Benchmark moved under `tools/benchmark` [#763](https://github.com/apache/dubbo-go-pixiu/pull/763)
* Pixiu CLI relocated to `pkg/cmd` [#596](https://github.com/apache/dubbo-go-pixiu/pull/596)

#### Resilience & Internal Quality

* Broad robustness enhancements [#644](https://github.com/apache/dubbo-go-pixiu/pull/644)
* Unified error code abstraction and consistent error handling
  [#809](https://github.com/apache/dubbo-go-pixiu/pull/809),
  [#782](https://github.com/apache/dubbo-go-pixiu/pull/782)
* Unified metric filter implementation [#799](https://github.com/apache/dubbo-go-pixiu/pull/799)
* Correct copy semantics in filter configs to avoid pointer sharing
  [#815](https://github.com/apache/dubbo-go-pixiu/pull/815),
  [#814](https://github.com/apache/dubbo-go-pixiu/pull/814)

#### CI & Build Improvements

* Update Go version to `1.25`, update CI workflows and lint rules
  [#752](https://github.com/apache/dubbo-go-pixiu/pull/752),
  [#666](https://github.com/apache/dubbo-go-pixiu/pull/666)
* GolangCI lint refactor and stability improvements
  [#650](https://github.com/apache/dubbo-go-pixiu/pull/650),
  [#734](https://github.com/apache/dubbo-go-pixiu/pull/734)
* Pipeline cleanup and unused GitHub Action removal
  [#775](https://github.com/apache/dubbo-go-pixiu/pull/775),
  [#786](https://github.com/apache/dubbo-go-pixiu/pull/786)
* Docker build optimization
  [#714](https://github.com/apache/dubbo-go-pixiu/pull/714),
  [#723](https://github.com/apache/dubbo-go-pixiu/pull/723)

### Bug Fixes

* SSE stream not closed on `io.EOF` [#676](https://github.com/apache/dubbo-go-pixiu/pull/676)
* HTTP proxy connection reuse issue [#578](https://github.com/apache/dubbo-go-pixiu/pull/578)
* Nil-pointer issue and non-unary response handling in access log filter [#713](https://github.com/apache/dubbo-go-pixiu/pull/713)
* Logger config overwrite error [#765](https://github.com/apache/dubbo-go-pixiu/pull/765)
* Data race fixes across components
  [#750](https://github.com/apache/dubbo-go-pixiu/pull/750),
  [#789](https://github.com/apache/dubbo-go-pixiu/pull/789)
* Missing Nacos field transmission [#679](https://github.com/apache/dubbo-go-pixiu/pull/679)
* Benchmark logic corrections and performance test cleanup [#819](https://github.com/apache/dubbo-go-pixiu/pull/819)

### Documentation

* Updated README and enhanced user guidance
  [#698](https://github.com/apache/dubbo-go-pixiu/pull/698),
  [#794](https://github.com/apache/dubbo-go-pixiu/pull/794),
  [#831](https://github.com/apache/dubbo-go-pixiu/pull/831)
* Full administrative documentation rewrite [#817](https://github.com/apache/dubbo-go-pixiu/pull/817)
* MCP configuration documentation added [#770](https://github.com/apache/dubbo-go-pixiu/pull/770)
* Quick-start guide for OPA HTTP filter [#751](https://github.com/apache/dubbo-go-pixiu/pull/751)
* Chinese README improvements [#641](https://github.com/apache/dubbo-go-pixiu/pull/641)
* Issue template clean-ups and improvements
  [#736](https://github.com/apache/dubbo-go-pixiu/pull/736),
  [#735](https://github.com/apache/dubbo-go-pixiu/pull/735)

### Cleanups & Removals

* update dubbo-go to latest version
  [#630](https://github.com/apache/dubbo-go-pixiu/pull/630),
  [#807](https://github.com/apache/dubbo-go-pixiu/pull/807),
  [#836](https://github.com/apache/dubbo-go-pixiu/pull/836),
  [#845](https://github.com/apache/dubbo-go-pixiu/pull/845)
* Removal of unused Seata proxy [#628](https://github.com/apache/dubbo-go-pixiu/pull/628)
* Removal of Istio integration [#622](https://github.com/apache/dubbo-go-pixiu/pull/622)
* Removal of static config providers [#764](https://github.com/apache/dubbo-go-pixiu/pull/764)

### Contributors

Special thanks to all contributors for their efforts in improving `dubbo-go-pixiu` (listed alphabetically):

@1kasa
@Alanxtl
@baerwang
@Chen-BUPT
@everfid-ever
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



## 1.0.0

### New Features
- [fail inject](https://github.com/apache/dubbo-go-pixiu/pull/571)
- [add support for header based route](https://github.com/apache/dubbo-go-pixiu/pull/565)
- [Add Maglev hashing LB algorithm](https://github.com/apache/dubbo-go-pixiu/pull/554)
- [triple proxy support import protosets](https://github.com/apache/dubbo-go-pixiu/pull/548)
- [Add GracefulShutdown Signal For Windows ](https://github.com/apache/dubbo-go-pixiu/pull/522)
- [Tracing support dubbo invoke](https://github.com/apache/dubbo-go-pixiu/pull/559)

### Enhancement
- [refactor prometheus metric](https://github.com/apache/dubbo-go-pixiu/pull/573)
- [remove unused pkg imports](https://github.com/apache/dubbo-go-pixiu/pull/574)
- [chore: unnecessary use of fmt.Sprintf](https://github.com/apache/dubbo-go-pixiu/pull/575)
- [chore:use wasm filter build tags add wasm](https://github.com/apache/dubbo-go-pixiu/pull/567)
- [docs:format and change samples link](https://github.com/apache/dubbo-go-pixiu/pull/556)
- [revert gatewayCmd to Run dubbo go pixiu](https://github.com/apache/dubbo-go-pixiu/pull/557)
- [full import format](https://github.com/apache/dubbo-go-pixiu/pull/527)
- [upgrade hessian2 to v1.11.3](https://github.com/apache/dubbo-go-pixiu/pull/516)

### Bugfixes
- [register hashing and array out of bounds and init hashing](https://github.com/apache/dubbo-go-pixiu/pull/530)
- [optimize timeout statusCode](https://github.com/apache/dubbo-go-pixiu/pull/521)
- [optimizing Metric Implementation](https://github.com/apache/dubbo-go-pixiu/pull/528)
- [add and modify nacos config arguments](https://github.com/apache/dubbo-go-pixiu/pull/524)
- [fix NPE when filter config is nil](https://github.com/apache/dubbo-go-pixiu/pull/517)
- [use wasmer-go v1.0.4 which is compatible with mac arm](https://github.com/apache/dubbo-go-pixiu/pull/515)
- [fix sample url using github.com/apache/dubbo-go-pixiu-samples](https://github.com/apache/dubbo-go-pixiu/pull/506)
- [traffic filter fix weight strategy and error handle within Apply method](https://github.com/apache/dubbo-go-pixiu/pull/507)
- [httpfilter loadbalancer does not work when it has spaces between multiple urls](https://github.com/apache/dubbo-go-pixiu/pull/513)

Milestone: [https://github.com/apache/dubbo-go-pixiu/milestone/8](https://github.com/apache/dubbo-go-pixiu/milestone/8)

## 0.6.0

### New Features
- [nacos config](https://github.com/apache/dubbo-go-pixiu/pull/497)
- [OSPP: Traffic Distribution](https://github.com/apache/dubbo-go-pixiu/pull/501)
- [Add Graceful Shutdown](https://github.com/apache/dubbo-go-pixiu/pull/474)
- [WASM Plugin for Pixiu](https://github.com/apache/dubbo-go-pixiu/pull/469)
- [deploy pixiu as dubbo service egress gateway in k8s istio](https://github.com/apache/dubbo-go-pixiu/pull/446)
- [ASoC 2022: Pixiu Metrics Implementation](https://github.com/apache/dubbo-go-pixiu/pull/480)
- [ospp: Feature/traffic](https://github.com/apache/dubbo-go-pixiu/pull/496)
- [feat:consistent hashing](https://github.com/apache/dubbo-go-pixiu/pull/436)


### Enhancement
- [Remove "Types" on Http to dubbo proxy](https://github.com/apache/dubbo-go-pixiu/pull/456)
- [ASoC 2002: Optimization of Pixiu timeout feature ](https://github.com/apache/dubbo-go-pixiu/pull/475)
-

### Bugfixes

- [fix response header Content-Type](https://github.com/apache/dubbo-go-pixiu/pull/462)
- [fix listener session exception](https://github.com/apache/dubbo-go-pixiu/pull/458)

Milestone: [https://github.com/apache/dubbo-go-pixiu/milestone/7](https://github.com/apache/dubbo-go-pixiu/milestone/7)


## 0.5.1

### New Features
- [Trace Support](https://github.com/apache/dubbo-go-pixiu/pull/394)
- [Health Check Support](https://github.com/apache/dubbo-go-pixiu/pull/421)
- [xDS Config Support](https://github.com/apache/dubbo-go-pixiu/pull/385)
- [LDS Support](https://github.com/apache/dubbo-go-pixiu/pull/417)
- [Direct Dubbo Invoke](https://github.com/apache/dubbo-go-pixiu/pull/434)


### Enhancement

- [SpringCloud subscribe strategy](https://github.com/apache/dubbo-go-pixiu/pull/425)
- [Style:optimization router match prefix definition](https://github.com/apache/dubbo-go-pixiu/pull/451)


### Bugfixes

- [Nacos registry bug](https://github.com/apache/dubbo-go-pixiu/pull/389)
- [Fix spring cloud error and refactor event callback](https://github.com/apache/dubbo-go-pixiu/pull/367)
- [Fix first call failure problem when using nacos registery](https://github.com/apache/dubbo-go-pixiu/pull/380)

Milestone: [https://github.com/apache/dubbo-go-pixiu/milestone/6](https://github.com/apache/dubbo-go-pixiu/milestone/6)

## 0.5.0

### New Features
- [Dubbo2Http Dubbo2Triple Triple2Dubbo proxy](https://github.com/apache/dubbo-go-pixiu/pull/347)
- [Http2Triple proxy](https://github.com/apache/dubbo-go-pixiu/pull/302)
- [Http2Dubbo default mapping rules](https://github.com/apache/dubbo-go-pixiu/pull/298)
- [Grpc proxy](https://github.com/apache/dubbo-go-pixiu/pull/315)
- [Dynamic cluster and route configuration from spring cloud zk registry](https://github.com/apache/dubbo-go-pixiu/pull/367)
- [Jwt auth Filter](https://github.com/apache/dubbo-go-pixiu/pull/303)
- [Https support multiple certificates](https://github.com/apache/dubbo-go-pixiu/pull/292)
- [Support build docker image](https://github.com/apache/dubbo-go-pixiu/pull/370)


### Enhancement

- [Add http2 listener for grpc proxy](https://github.com/apache/dubbo-go-pixiu/pull/315)
- [Route using trie](https://github.com/apache/dubbo-go-pixiu/pull/310)
- [Http2Grpc use grpc reflection server](https://github.com/apache/dubbo-go-pixiu/pull/317)
- [Get cpu core number in container](https://github.com/apache/dubbo-go-pixiu/pull/340)
- [Filter Chain refactor](https://github.com/apache/dubbo-go-pixiu/pull/307)
- [Upgrade hessian2 to v1.11.0](https://github.com/apache/dubbo-go-pixiu/pull/352)
- [Upgrade upgrade dubbogo version to 3.0](https://github.com/apache/dubbo-go-pixiu/pull/334)
- [Upgrade keyfunc to new stable release v1.0.0](https://github.com/apache/dubbo-go-pixiu/pull/318)


### Bugfixes

- [Fix write error when handle gRPC request using http2 manager](https://github.com/apache/dubbo-go-pixiu/pull/372)
- [Fix spring cloud error and refactor event callback](https://github.com/apache/dubbo-go-pixiu/pull/367)
- [Fix first call failure problem when using nacos registery](https://github.com/apache/dubbo-go-pixiu/pull/380)

## 0.4.0

### New Features
- [dynamic cluster and route configuration from spring cloud nacos registry](https://github.com/apache/dubbo-go-pixiu/pull/255)
- [dynamic dubbo proxy configuration from zk registry](https://github.com/apache/dubbo-go-pixiu/pull/256)
- [http to grpc proxy](https://github.com/apache/dubbo-go-pixiu/pull/244)
- [http to http proxy](https://github.com/apache/dubbo-go-pixiu/pull/242)
- [tracing with jaeger](https://github.com/apache/dubbo-go-pixiu/pull/236)
- [cors policy](https://github.com/apache/dubbo-go-pixiu/pull/249)

### Enhancement

- [add more samples](https://github.com/apache/dubbo-go-pixiu/pull/271)
- [use cobra cmd tool](https://github.com/apache/dubbo-go-pixiu/pull/234)
- [add samples quick start script](https://github.com/apache/dubbo-go-pixiu/pull/226)
- [upgrade hessian2 to v1.9.3](https://github.com/apache/dubbo-go-pixiu/pull/248)
- [rename onAir property to enable](https://github.com/apache/dubbo-go-pixiu/pull/243)
- [tracing optimize](https://github.com/apache/dubbo-go-pixiu/pull/257/files)
- [support https](https://github.com/apache/dubbo-go-pixiu/pull/213)

### Bugfixes

- [Fix request body miss problem](https://github.com/apache/dubbo-go-pixiu/pull/260)
- [Fix HttpContext reset bug](https://github.com/apache/dubbo-go-pixiu/pull/254)
- [Fix env value can't be set](https://github.com/apache/dubbo-go-pixiu/pull/239)
- [Fix filterManager get filters with random order](https://github.com/apache/dubbo-go-pixiu/pull/264)
- [Fix nil issue for timeout filter](https://github.com/apache/dubbo-go-pixiu/pull/278)

Milestone: [https://github.com/apache/dubbo-go-pixiu/milestone/5](https://github.com/apache/dubbo-go-pixiu/milestone/5)


## 0.3.0

### New Features
- [rate limit filter](https://github.com/apache/dubbo-go-pixiu/pull/169)
- [add integrate test](https://github.com/apache/dubbo-go-pixiu/pull/183)
- [handle rate limit config update event](https://github.com/apache/dubbo-go-pixiu/pull/196)
- [add otel metric export to prometheus in pixiu](https://github.com/apache/dubbo-go-pixiu/pull/204)
- [make Pixiu Admin config management finer-grained](https://github.com/apache/dubbo-go-pixiu/pull/171)

### Enhancement
- [update samples/admin](https://github.com/apache/dubbo-go-pixiu/pull/208)
- [update ratelimit samples](https://github.com/apache/dubbo-go-pixiu/pull/206)
- [make router case sensitive](https://github.com/apache/dubbo-go-pixiu/pull/209)
- [add more test case](https://github.com/apache/dubbo-go-pixiu/pull/203)
- [Enrich filter test case](https://github.com/apache/dubbo-go-pixiu/pull/202)
- [Enrich response.go's test case](https://github.com/apache/dubbo-go-pixiu/pull/197)

### Bugfixes
- [Fix CI check status not match required](https://github.com/apache/dubbo-go-pixiu/pull/199)
- [Fix timeout config overridden](https://github.com/apache/dubbo-go-pixiu/pull/190)
- [Fix/quickstart](https://github.com/apache/dubbo-go-pixiu/pull/191)
- [FixBug: can't delete node by path](https://github.com/apache/dubbo-go-pixiu/pull/201)
- [Fix flow chart](https://github.com/apache/dubbo-go-pixiu/pull/205)
- [Fix reviewdog](https://github.com/apache/dubbo-go-pixiu/pull/195)

Milestone: [https://github.com/apache/dubbo-go-pixiu/milestone/4](https://github.com/apache/dubbo-go-pixiu/milestone/4)


## 0.2.1

### Enhancement
- [Change the mascot of pixiu](https://github.com/apache/dubbo-go-pixiu/pull/178)
- [reviewdog use default flags](https://github.com/apache/dubbo-go-pixiu/pull/167)
- [moving param types into parameter configuration instead of standalone](https://github.com/apache/dubbo-go-pixiu/pull/161)
- [fix version field](https://github.com/apache/dubbo-go-pixiu/pull/166)
- [Add license-eye to check and fix license headers](https://github.com/apache/dubbo-go-pixiu/pull/164)
- [Improve: expand filterFuncCacheMap initial length](https://github.com/apache/dubbo-go-pixiu/pull/174)
- [Refractor config_load.go](https://github.com/apache/dubbo-go-pixiu/pull/158)

Milestone: [https://github.com/apache/dubbo-go-pixiu/milestone/3](https://github.com/apache/dubbo-go-pixiu/milestone/3)


## 0.2.0

### New Features
- [Add dubbo-go-proxy admin](https://github.com/dubbogo/dubbo-go-proxy/pull/115)
- [Add plugin](https://github.com/dubbogo/dubbo-go-proxy/pull/109)

### Bugfixes
- [Fix: remove replace-path-filter](https://github.com/dubbogo/dubbo-go-proxy/pull/118)

Milestone: [https://github.com/dubbogo/dubbo-go-proxy/milestone/2](https://github.com/dubbogo/dubbo-go-proxy/milestone/2?closed=1)
