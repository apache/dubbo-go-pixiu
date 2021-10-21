# 从 URI，Query，Body 各个部分获取参数

> GET 请求 [samples](https://github.com/apache/dubbo-go-pixiu/tree/develop/samples/dubbogo/simple/mix)

## 简单示例

### 接口配置

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/user/:name'
    type: restful
    description: user
    methods:
      - httpVerb: GET
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
              mapType: "string"
            - name: queryStrings.age
              mapTo: 1
              mapType: "int"
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "GetUserByNameAndAge"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
              mapType: "string"
            - name: requestBody._all
              mapTo: 1
              mapType: "object"
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "UpdateUserByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    methods:
      - httpVerb: PUT
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.name
              mapTo: 0
              mapType: "string"
            - name: requestBody._all
              mapTo: 1
              mapType: "object"
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "UpdateUserByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### 测试

- 来自 uri 和 query

```bash
curl localhost:port/api/v1/test-dubbo/user/tc?age=99 -X GET 
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
curl localhost:port/api/v1/test-dubbo/user?name=tc -X PUT -d '{"id":"0001","code":1,"name":"tc","age":55}' --header "Content-Type: application/json"
```

结果

```bash
true
```

[上一页](dubbo.md)
