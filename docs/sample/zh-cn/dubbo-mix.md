# 从 URI，Query，Body 各个部分获取参数

> GET 请求 [samples](https://github.com/dubbogo/dubbo-go-proxy-samples/tree/master/dubbo/apiconfig/mix)

## 简单示例

### 配置

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/user/:name'
    type: restful
    description: user
    filters:
      - filter0
    methods:
      - httpVerb: GET
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
            - name: queryStrings.age
              mapTo: 1
          applicationName: "UserService"
          interface: "com.dubbogo.proxy.UserService"
          method: "GetUserByNameAndAge"
          paramTypes: [ "string", "int" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
            - name: requestBody._all
              mapTo: 1
          applicationName: "UserService"
          interface: "com.dubbogo.proxy.UserService"
          method: "UpdateUserByName"
          paramTypes: [ "string", "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    filters:
      - filter0
    methods:
      - httpVerb: PUT
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.name
              mapTo: 0
            - name: requestBody._all
              mapTo: 1
          applicationName: "UserService"
          interface: "com.dubbogo.proxy.UserService"
          method: "UpdateUserByName"
          paramTypes: [ "string", "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### 测试

- 来自 uri 和 query

```bash
curl localhost:8888/api/v1/test-dubbo/user/tc?age=99 -X GET 
```

结果

```bash
{
  "age": 15,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T20:54:38.746+08:00"
}
```

- 来自 body 和 query

```bash
curl localhost:8888/api/v1/test-dubbo/user?name=tc -X PUT -d '{"id":"0001","code":1,"name":"tc","age":55}' --header "Content-Type: application/json"
```

结果

```bash
true
```

[上一页](./dubbo.md)