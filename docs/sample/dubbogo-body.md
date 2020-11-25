# Get the parameter from the body

> POST request

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
                interface: "com.ic.user.UserProvider"
                method: "CreateUser"
                paramTypes: [ "com.ikurento.user.User" ]
                group: "test"
                version: 1.0.0
                clusterName: "test_dubbo"
```

> mapTo: 0 needed

### Request

```bash
curl host:port/api/v1/test-dubbo/user -X POST -d '{"name":"tiecheng","id":"0002","age":18}' --header "Content-Type: application/json"
```

### Response

- If first add, return like:

```json
{
    "age": 18,
    "iD": "0002",
    "name": "tiecheng"
}
```

- If you add user multi, return like: 

```json
{
    "message": "data is exist"
}
```

## Mapping



