# Get the parameter from the query

> GET request

## Simple

### Config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    methods:
      - httpVerb: GET
        onAir: true
        timeout: 100ms
        inboundRequest:
          requestType: http
          queryStrings:
            - name: name
              required: true
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.name
              mapTo: 0
          applicationName: "UserProvider"
          interface: "com.ic.user.UserProvider"
          method: "GetUserByName"
          paramTypes: [ "java.lang.String" ]
          group: "test"
          version: 1.0.0
          clusterName: "test_dubbo"
```

### Request

```bash
curl localhost:8888/api/v1/test-dubbo/user?name=tc -X GET 
```

### Response

- If exist, will return:

```bash
{
    "age": 18,
    "iD": "0001",
    "name": "tc"
}
```

- Not found, return: nil

