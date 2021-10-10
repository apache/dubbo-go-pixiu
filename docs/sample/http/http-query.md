# Get the parameter from the query

> GET request

## Simple Demo

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-http/user'
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
          host: "127.0.0.1:1314"
          path: "/user"
      - httpVerb: POST
        enable: true
        timeout: 10s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: http
          host: "127.0.0.1:1314"
          path: "/user/"
```

### Request

```bash
curl http://localhost:8888/api/v1/test-http/user?name=tc -X GET 
```

### Response

- Successful result will return:

```bash
{
  "id": "0001",
  "name": "tc",
  "age": 18,
  "time": "2020-12-30T14:07:07.9432117+08:00"
}
```
