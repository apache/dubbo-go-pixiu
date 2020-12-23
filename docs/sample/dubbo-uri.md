# Get the parameter from the uri

> GET request [samples](https://github.com/dubbogo/dubbo-go-proxy-samples/tree/master/dubbo/apiconfig/uri)

## Simple Demo

### Config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/user/name/:name'
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
          uri:
            - name: name
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
          applicationName: "UserProvider"
          interface: "com.dubbogo.proxy.UserService"
          method: "GetUserByName"
          paramTypes: [ "string" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/user/code/:code'
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
          uri:
            - name: code
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.code
              mapTo: 0
          applicationName: "UserProvider"
          interface: "com.dubbogo.proxy.UserService"
          method: "GetUserByCode"
          paramTypes: [ "int" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/user/name/:name/age/:age'
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
          uri:
            - name: name
              required: true
            - name: age
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: uri.name
              mapTo: 0
            - name: uri.age
              mapTo: 1
          applicationName: "UserProvider"
          interface: "com.dubbogo.proxy.UserService"
          method: "GetUserByNameAndAge"
          paramTypes: [ "string", "int" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### Test

- single param string

```bash
curl localhost:8888/api/v1/test-dubbo/user/name/tc -X GET 
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
curl localhost:8888/api/v1/test-dubbo/user/name/tc/age/99 -X GET 
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