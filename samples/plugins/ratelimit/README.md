# This sample for rate limit filter.

[中文](README_CN.md)

The current limiting filter can provide out-of-the-box current limiting function to ensure service stability.

> The filter based on [Sentinel](https://github.com/alibaba/sentinel-golang), see more here [Wiki](https://sentinelguard.io/zh-cn/docs/introduction.html) .



### Create the API:
- Create a simple Http API,Reference the [Create a simple Http API](../../dubbogo/http/README.md),then we got a path.

- just test it: `curl http://localhost:8888/api/v1/test-dubbo/user?name=tc -X GET `

### Protect the API:
#### rate limit config [click here](../../../pkg/filter/ratelimit/mock/config.yml)
- The first step,Define the resources to be protected. A resource can contain one or more matching paths.

  Here, what we want to protect is an exact path, as defined below. Of course, we also support regex, just set matchStrategy to 1.
```
resources:
  - name: test-http
    items:
      #Exact
      - matchStrategy: 0
        pattern: "/api/v1/test-dubbo/user"
      #Regex
      - matchStrategy: 1
        pattern: "/api/*/test-dubbo/user"
```

- The second step is to set the rules. for example, the threshold is 100, and the statistical statintervalinms is 1000ms, which means that the qps/tps of the `resource` will not exceed 100.
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

### Test:

- For a simpler test, we set qps to 1

- run [test](test.go) and check output.

