# 从表单请求里面获取参数

> GET 请求 [samples](https://github.com/apache/dubbo-go-pixiu/tree/develop/samples/dubbogo/simple/query)

## 简单示例

### 接口配置

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/userByName'
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
            - name: queryStrings.name
              mapTo: 0
              mapType: "string"
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "GetUserByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/userByNameAndAge'
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
            - name: queryStrings.name
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
  - path: '/api/v1/test-dubbo/userByCode'
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
            - name: queryStrings.code
              mapTo: 0
              mapType: "int"
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "GetUserByCode"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### 测试

- 单个 string 参数

```bash
curl localhost:port/api/v1/test-dubbo/userByName?name=tc -X GET 
```

如果存在数据，返回:

```bash
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T15:51:36.333+08:00"
}
```

不存在，返回空

- 多个参数

```bash
curl localhost:port/api/v1/test-dubbo/userByNameAndAge?name=tc&age=99 -X GET 
```

结果

```bash
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T15:51:36.333+08:00"
}
```

[上一页](dubbo.md)
