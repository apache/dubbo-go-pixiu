# Get the parameter from the body

> POST request

## passthrough

**config**

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
                interface: "com.ic.user.UserProvider"
                method: "CreateUser"
                paramTypes: [ "com.ikurento.user.User" ]
                group: "test"
                version: 1.0.0
                clusterName: "test_dubbo"
```

**request**

```bash
curl localhost:8888/api/v1/test-dubbo/user -X POST -d '{"name":"tiecheng","id":"0001","age":18}' --header "Content-Type: application/json"
```

**response**

```bash
{"age":18,"i_d":"0001","name":"tiecheng"}
```

## mapping



