# Release Notes
---
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

Milestone: [0.0.4](https://github.com/apache/dubbo-go-pixiu/milestone/5) 


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
