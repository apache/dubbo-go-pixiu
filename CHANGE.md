# Release Notes

---
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
