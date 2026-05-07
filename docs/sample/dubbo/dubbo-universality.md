# Use dubbo request universality

> POST request [samples](https://github.com/apache/dubbo-go-pixiu-samples/tree/main/dubbogo/simple/proxy)

## Direct Generic Contract (Breaking Change)

Pixiu direct generic invocation no longer accepts request-provided `types` as the method
signature source. Direct mode now requires:

- `integrationRequest.url`
- `integrationRequest.interface`
- `integrationRequest.method`
- `integrationRequest.parameterTypes`
- `integrationRequest.serialization`

`mappingParams` is still responsible for values, but no longer defines method signatures.

For Triple direct generic invocation, the provider must expose a generic
`$invoke` endpoint. The dubbo-go generic client uses non-IDL mode for this
path. IDL-generated Triple handlers normally register only concrete RPC methods,
so a generic call to `$invoke` has no matching Triple handler and may return
`404 Not Found`.

## Suggest

> In this way, you can request your dubbo rpc service by defined one api for every cluster.
> The following sample uses registry mode. In this mode, `opt.types` can still
> provide the generic invocation signature. Direct mode must use
> `integrationRequest.parameterTypes` instead.

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/:interface'
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
            - name: uri.interface
              mapTo: opt.interface
            - name: queryStrings.method
              mapTo: opt.method
            - name: queryStrings.group
              mapTo: opt.group
            - name: queryStrings.version
              mapTo: opt.version
          clusterName: "test_dubbo"
```

### Test

- single param string

```bash
curl host:port/api/v1/test-dubbo/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByName -X POST -d '{"types":["string"],"values":"tc"}' --header "Content-Type: application/json"
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
curl host:port/api/v1/test-dubbo/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=GetUserByCode -X POST -d '{"types":["int"],"values":1}' --header "Content-Type: application/json"
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
curl host:port/api/v1/test-dubbo/com.dubbogo.proxy.UserService?group=test&version=1.0.0&method=UpdateUserByName -X POST -d '{"types":["string","body"],"values":["tc",{"id":"0001","code":1,"name":"tc","age":15}]}' --header "Content-Type: application/json"
```

result

```bash
true
```

### Special config

#### Code

Supported `mapTo` options:

```yaml
- opt.types
- opt.group
- opt.version
- opt.interface
- opt.method
- opt.values
```

#### Options

By configuring mapTo with option keywords(listed below), Pixiu will assemble generic params to invoke.

```go
// GenericService uses for generic invoke for service call
type GenericService struct {
	Invoke func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) `dubbo:"$invoke"`
}
```

- opt.types

> dubbo generic types

Use as the `types` argument of dubbogo `GenericService#Invoke` when `integrationRequest.parameterTypes` is not configured.
For direct generic invocation, configure `integrationRequest.parameterTypes` instead.

- opt.method

Use as the `methodName` argument of dubbogo `GenericService#Invoke`.

- opt.group

Dubbo reference group.

- opt.version

Dubbo reference version.

- opt.interface

Dubbo service interface.

- opt.values

Use as the `args` argument of dubbogo `GenericService#Invoke`.

#### Explain

##### Single params

request body

```json
{
    "types": ["string"],
    "values": "tc"
}
```

```yaml
  - name: requestBody.types
    mapTo: opt.types
```

- `requestBody.types` means body content with types key.
- `opt.types` means using the value as the generic invocation `types` argument in this registry-mode sample.

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

Please pay attention to the special situation of configuration the degrees of freedom is not very high, if can't meet the scene, please mention [issue](https://github.com/apache/dubbo-go-pixiu/issues), thank you.

[Previous](dubbo.md)
