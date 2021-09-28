# Get the parameter from the query

> GET request

## Simple Demo

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    methods:
      - httpVerb: GET
        enable: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
          queryStrings:
            - name: name
              required: true
        integrationRequest:
          requestType: http
          host: 127.0.0.1:8889
          path: /UserProvider/GetUserByName
          mappingParams:
            - name: queryStrings.name
              mapTo: queryStrings.name
          group: "test"
          version: 1.0.0
```

### Request

```bash
curl http://localhost:8888/api/v1/test-dubbo/user?name=tc -X GET 
```

### Response

- Successful result will return:

```bash
{
  "id": "0001",
  "code": 1,
  "name": "tc",
  "age": 18,
  "time": "2020-12-24T16:46:31.8409857+08:00"
}
```
