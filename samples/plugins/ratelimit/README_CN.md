# 限流插件
限流过滤器可以提供开箱即用的限流功能，保障服务稳定性。

> 限流功能基于 [Sentinel](https://github.com/alibaba/sentinel-golang), 更详细的配置请参照 [Wiki](https://sentinelguard.io/zh-cn/docs/introduction.html) 。



### 创建你的API:
- 创建一个简单的API，参考[创建一个简单的Http API](../../dubbogo/http/README.md)，然后我们获得了一个可访问的路径

- 然后测试一下: `curl http://localhost:8888/api/v1/test-dubbo/user?name=tc -X GET `

### 保护你的API:
#### rate limit config [查看全文](../../../pkg/filter/ratelimit/mock/config.yml)
- 第一步,定义要保护的资源，一个资源可以包含一或更多个匹配的路径。
  
  在这里，我们要保护的是一个确切的路径，如下的定义即可。当然，我们也支持正则匹配，只要把matchStrategy设为1。
```
resources:
  - name: test-http
    items:
      #精确匹配
      - matchStrategy: 0
        pattern: "/api/v1/test-dubbo/user"
      #正则匹配
      - matchStrategy: 1
        pattern: "/api/*/test-dubbo/user"
```

- 第二步，设置限流的规则，例如上限设置为100，统计周期为1000ms，这意味着这个`资源`的qps/tps最高只能达到100。
```
  rules:
    #qps sample At most 100 requests can be passed in 1000ms, so qps is 100
    - enable: true
      flowRule:
        #the resource's name
        resource: "test-http"
        threshold: 100
        statintervalinms: 1000
```

### 测试:

- 为了更简单的测试，我们把qps设置为1

- 运行[并发测试](test.go)并查看输出