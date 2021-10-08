# 从请求体里面获取参数

> POST 请求 [samples](https://github.com/apache/dubbo-go-pixiu/tree/develop/samples/dubbogo/simple/body)

## 透传

### 接口配置

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    methods:
      - httpVerb: POST
        enable: true
        timeout: 10s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody._all
              mapTo: 0
              mapType: "object"
          applicationName: "UserProvider"
          interface: "com.dubbogo.pixiu.UserService"
          method: "CreateUser"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        enable: true
        timeout: 10s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody._all
              mapTo: 0
              mapType: "object"
          applicationName: "UserProvider"
          interface: "com.dubbogo.pixiu.UserService"
          method: "UpdateUser"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/user2'
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
            - name: requestBody.name
              mapTo: 0
              mapType: "string"
            - name: requestBody.user
              mapTo: 1
              mapType: "object"
          applicationName: "UserService"
          interface: "com.dubbogo.pixiu.UserService"
          method: "UpdateUserByName"
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

> 当使用透传的时候, mapTo: 0 是需要配置的

### 测试

- 透传

```bash
curl host:port/api/v1/test-dubbo/user -X POST -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header "Content-Type: application/json"
```

如果存在数据，返回:

```json
{
  "age": 99,
  "code": 3,
  "iD": "0003",
  "name": "dubbogo"
}
```
如果存在数据，则返回如下：

```json
{
    "message": "data is exist"
}
```

- 更新

```bash
curl host:port/api/v1/test-dubbo/user -X PUT -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header "Content-Type: application/json"
```

结果

```bash
true
```

- 从请求体解析多个参数

```bash
curl host:port/api/v1/test-dubbo/user2 -X PUT -d '{"name":"tc","user":{"id":"0001","code":1,"name":"tc","age":99}}' --header "Content-Type: application/json"
```

结果

```bash
true
```

[上一页](dubbo.md)
