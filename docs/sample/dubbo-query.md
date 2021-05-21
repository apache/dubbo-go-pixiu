# Get the parameter from the query

> GET request [samples](https://github.com/dubbogo/dubbo-go-proxy/tree/develop/samples/dubbogo/simple/query)

## Simple Demo

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/userByName'
    type: restful
    description: user
    methods:
      - httpVerb: GET
        onAir: true
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
        onAir: true
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
        onAir: true
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

### Test

- single param string

```bash
curl localhost:port/api/v1/test-dubbo/userByName?name=tc -X GET 
```

If exist, will return:

```bash
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T15:51:36.333+08:00"
}
```

Not found, return: nil

- multi params

```bash
curl localhost:port/api/v1/test-dubbo/userByNameAndAge?name=tc&age=99 -X GET 
```

result

```bash
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T15:51:36.333+08:00"
}
```

[Previous](./dubbo.md)