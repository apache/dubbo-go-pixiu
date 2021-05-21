# Get the parameter from the part of uri,query,body

> GET request [samples](https://github.com/dubbogo/dubbo-go-proxy/tree/develop/samples/dubbogo/simple/mix)

## Simple Demo

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/user/:name'
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
        onAir: true
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

### Test

- from uri and query

```bash
curl localhost:port/api/v1/test-dubbo/user/tc?age=99 -X GET 
```

result

```bash
{
  "age": 15,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T20:54:38.746+08:00"
}
```

- from body and query

```bash
curl localhost:port/api/v1/test-dubbo/user?name=tc -X PUT -d '{"id":"0001","code":1,"name":"tc","age":55}' --header "Content-Type: application/json"
```

result

```bash
true
```

[Previous](./dubbo.md)