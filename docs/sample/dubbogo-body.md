# Get the parameter from the body

> POST request [samples](https://github.com/dubbogo/dubbo-go-proxy-samples/tree/master/dubbo/apiconfig/body)

## Passthroughs

### Config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    methods:
      - httpVerb: POST
        onAir: true
        timeout: 10s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody._all
              mapTo: 0
          applicationName: "UserProvider"
          interface: "com.dubbogo.proxy.UserService"
          method: "CreateUser"
          paramTypes: [ "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
      - httpVerb: PUT
        onAir: true
        timeout: 10s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody._all
              mapTo: 0
          applicationName: "UserProvider"
          interface: "com.dubbogo.proxy.UserService"
          method: "UpdateUser"
          paramTypes: [ "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
  - path: '/api/v1/test-dubbo/user2'
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
            - name: requestBody.name
              mapTo: 0
            - name: requestBody.user
              mapTo: 1
          applicationName: "UserService"
          interface: "com.dubbogo.proxy.UserService"
          method: "UpdateUserByName"
          paramTypes: [ "string", "object" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

> when passthroughs, mapTo: 0 needed

### Test

- passthroughs

```bash
curl host:port/api/v1/test-dubbo/user -X POST -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header "Content-Type: application/json"
```

If first add, return like:

```json
{
  "age": 99,
  "code": 3,
  "iD": "0003",
  "name": "dubbogo"
}
```
If you add user multi, return like: 

```json
{
    "message": "data is exist"
}
```

- update

```bash
curl host:port/api/v1/test-dubbo/user -X PUT -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header "Content-Type: application/json"
```

result

```bash
true
```

- body parse multi params

```bash
curl host:port/api/v1/test-dubbo/user2 -X PUT -d '{"name":"tc","user":{"id":"0001","code":1,"name":"tc","age":99}}' --header "Content-Type: application/json"
```

result

```bash
true
```