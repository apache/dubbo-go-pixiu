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
          requestType: http
          host: "127.0.0.1:1314"
```

### Request

```bash
curl host:port/api/v1/test-http/user -X POST -d '{"name":"ttc","id":"0003","age":28}' --header "Content-Type: application/json"
```

### Response

- If first add, return like:

```json
{
    "age": 28,
    "iD": "0003",
    "name": "ttc",
    "time": "0001-01-01T00:00:00Z"
}
```

- If you add user multi, return like: 

```json
{
    "message": "data is exist"
}
```



