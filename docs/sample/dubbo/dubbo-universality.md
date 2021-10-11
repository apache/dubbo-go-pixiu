# Use dubbo request universality

> POST request [samples](https://github.com/dubbogo/dubbo-go-proxy/tree/develop/samples/dubbogo/simple/proxy)

## Suggest

> In this way, you can request your dubbo rpc service by defined one api for every cluster.

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/:application/:interface'
    type: restful
    description: common
    methods:
      - httpVerb: POST
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody.values
              mapTo: opt.values
            - name: requestBody.types
              mapTo: opt.types
            - name: uri.application
              mapTo: opt.application
            - name: uri.interface
              mapTo: opt.interface
            - name: queryStrings.method
              mapTo: opt.method
            - name: queryStrings.group
              mapTo: opt.group
            - name: queryStrings.version
              mapTo: opt.version
          # Notice: this is the really paramTypes to dubbo service, it takes precedence over paramTypes when it is finally called.
          clusterName: "test_dubbo"
```

### Test

- single param string

```bash
curl host:port/api/v1/test-dubbo/UserService/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByName -X POST -d '{"types":["string"],"values":"tc"}' --header "Content-Type: application/json"
```

result

```json
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T20:54:38.746+08:00"
}
```

- single param int

```bash
curl host:port/api/v1/test-dubbo/UserService/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByCode -X POST -d '{"types":["int"],"values":1}' --header "Content-Type: application/json"
```

result

```json
{
  "age": 18,
  "code": 1,
  "iD": "0001",
  "name": "tc",
  "time": "2020-12-20T20:54:38.746+08:00"
}
```

- multi params

```bash
curl host:port/api/v1/test-dubbo/UserService/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=UpdateUserByName -X POST -d '{"types":["string","body"],"values":["tc",{"id":"0001","code":1,"name":"tc","age":15}]}' --header "Content-Type: application/json"
```

result

```bash
true
```

### Special config

#### Code

```go
const (
	optionKeyTypes       = "types"
	optionKeyGroup       = "group"
	optionKeyVersion     = "version"
	optionKeyInterface   = "interface"
	optionKeyApplication = "application"
	optionKeyMethod      = "method"
	optionKeyValues      = "values"
)
```

#### Options

By configuring mapTo with option keywords(listed below), Pixiu will assemble generic params to invoke.

```go
// GenericService uses for generic invoke for service call
type GenericService struct {
	Invoke       func(ctx context.Context, req []interface{}) (interface{}, error) `dubbo:"$invoke"`
	referenceStr string
}
```

- opt.types

> dubbo generic types

Use for dubbogo `GenericService#Invoke` func arg 2rd param.

- opt.method

Use for dubbogo `GenericService#Invoke` func arg 1rd param.

- opt.group

Dubbo group in `ReferenceConfig#Group`. 

- opt.version

Dubbo version in `ReferenceConfig#Version`.

- opt.interface

Dubbo interface in `ReferenceConfig#InterfaceName`.

- opt.application

Now only use for part of cache key.

- opt.values

Use for dubbogo `GenericService#Invoke` func arg 3rd param.

#### Explain

##### Single params 

request body

```json
{
    "types": [
        "string"
    ],
    "values": "tc"
}
```

```yaml
            - name: requestBody.types
              mapTo: opt.types
```

- `requestBody.types` means body content with types key.
- `opt.types` means use types option.

##### Multiple params

```json
{
  "types": [
    "java.lang.String",
    "object"
  ],
  "values": [
    "tc",
    {
      "id": "0001",
      "code": 1,
      "name": "tc",
      "age": 99
    }
  ]
}
```

Please pay attention to the special situation of configuration the degrees of freedom is not very high, if can't meet the scene, please mention [issue](https://github.com/dubbogo/dubbo-go-proxy/issues), thank you.

[Previous](dubbo.md)
