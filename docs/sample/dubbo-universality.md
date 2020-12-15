# Use dubbo request universality

> POST request

## Suggest

> In this way, you can request your dubbo rpc service by defined one api for every cluster.

### Config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/:application/:interface'
    type: restful
    description: common
    methods:
      - httpVerb: POST
        onAir: true
        timeout: 100s
        inboundRequest:
          requestType: http
          queryStrings:
            - name: method
              required: true
            - name: group
              required: false
            - name: version
              required: false
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: requestBody.values
              mapTo: -1
            - name: requestBody.types
              mapTo: 1
              opt:
                open: true
                name: types
            - name: uri.application
              mapTo: 2
              opt:
                open: true
                name: application
            - name: uri.interface
              mapTo: 3
              opt:
                open: true
                name: interface
            - name: queryStrings.method
              mapTo: 4
              opt:
                open: true
                name: method
            - name: queryStrings.group
              mapTo: 5
              opt:
                open: true
                name: group
            - name: queryStrings.version
              mapTo: 6
              opt:
                open: true
                name: version
          clusterName: "test_dubbo"
```

### Request

```bash
curl host:port/api/v1/test-dubbo/UserProvider/com.ic.user.UserProvider?method=GetUserByName&group=test&version=1.0.0 -X POST -d '{"types":["java.lang.String"],"values":"tc"}' --header "Content-Type: application/json"
```

### Response

```json
{
    "age": 18,
    "code": 1,
    "iD": "0001",
    "name": "tc"
}
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
)
```

#### Options

Assemble generic params to invoke.

```go
// GenericService uses for generic invoke for service call
type GenericService struct {
	Invoke       func(ctx context.Context, req []interface{}) (interface{}, error) `dubbo:"$invoke"`
	referenceStr string
}
```

- types

> dubbo generic types

Use for dubbogo `GenericService#Invoke` func arg 2rd param.

- method

Use for dubbogo `GenericService#Invoke` func arg 1rd param.

- group

Dubbo group in `ReferenceConfig#Group`. 

- version

Dubbo version in `ReferenceConfig#Version`.

- interface

Dubbo interface in `ReferenceConfig#InterfaceName`.

- application

Now only use for part of cache key.

#### Explain

##### Single params 

request body

```json
{
    "types": [
        "java.lang.String"
    ],
    "values": "tc"
}
```

```yaml
            - name: requestBody.types
              mapTo: 1
              opt:
                open: true
                name: types
```

- `requestBody.types` means body content with types key.
- `opt.name` means use types option.
- `opt.usable` means if remove the request to downstream service.

##### Multiple params

```json
{
  "types": [
    "java.lang.String",
    "com.dubbogo.proxy.User"
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

Body's values will use array for multi parameters. Config like follow `mapTo -1` will run special treatment.

```yaml
            - name: requestBody.values
              mapTo: -1
            - name: requestBody.types
              mapTo: 1
              opt:
                open: true
                name: types
```

