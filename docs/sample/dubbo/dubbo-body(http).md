# Get the parameter from the body

> POST request

## Passthroughs

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/user'
    type: restful
    description: user
    methods:
      - httpVerb: POST
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
          path: /UserProvider/CreateUser
          group: "test"
          version: 1.0.0
```

### Request

```bash
curl host:port/api/v1/test-dubbo/user -X POST -d '{"name": "tiecheng","id": "0002","code": 3,"age": 18}' --header "Content-Type: application/json"
```

### Response

- If first add, return like:

```json
{
  "id": "0002",
  "code": 3,
  "name": "tiecheng",
  "age": 18,
  "time": "0001-01-01T00:00:00Z"
}
```

- If you add user multi, return like:

```json
{
    "message": "data is exist"
}
```
