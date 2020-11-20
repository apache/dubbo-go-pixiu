# use dubbo request universality

> POST request

## suggest

> In this way, you can request your dubbo rpc service by defined one api for every cluster.

**config**

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
              mapTo: 0
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

**request**

```bash
curl host:port/api/v1/test-dubbo/UserProvider/com.ic.user.UserProvider?method=GetUserByName&group=test&version=1.0.0 -X POST -d '{"types":["java.lang.String"],"values":"tc"}' --header "Content-Type: application/json"
```

**response**

```json
{
    "age": 18,
    "code": 1,
    "iD": "0001",
    "name": "tc"
}
```

### special config

**the code**

```go
var DefaultMapOption = MapOption{
	"types":       &DubboParamTypesOpt{},
	"group":       &DubboGroupOpt{},
	"version":     &DubboVersionOpt{},
	"interface":   &DubboInterfaceOpt{},
	"application": &DubboApplicationOpt{},
	"method":      &DubboMethodOpt{},
}
```

**the config**

```yaml
              opt:
                open: true
                name: types
```

