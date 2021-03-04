# 从表单请求里面获取参数

> GET 请求 [samples](https://github.com/dubbogo/dubbo-go-pixiu/tree/develop/samples/dubbogo/simple/query)

## 简单示例

### 接口配置

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/userByName'
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
          queryStrings:
            - name: name
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.name
              mapTo: 0
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "GetUserByName"
          paramTypes: [ "java.lang.String" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/userByNameAndAge'
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
          queryStrings:
            - name: name
              required: true
            - name: age
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.name
              mapTo: 0
            - name: queryStrings.age
              mapTo: 1
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "GetUserByNameAndAge"
          paramTypes: [ "java.lang.String","java.lang.Integer" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/userByCode'
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
          queryStrings:
            - name: code
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.code
              mapTo: 0
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "GetUserByCode"
          paramTypes: [ "java.lang.Integer" ]
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

[上一页](./dubbo.md)