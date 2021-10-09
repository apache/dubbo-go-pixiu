# Mock request

## Simple Demo

### Api Config

```yaml
name: pixiu
description: pixiu sample
resources:
  - path: '/api/v1/test-dubbo/mock'
    type: restful
    description: mock
    methods:
      - httpVerb: GET
        enable: true
        mock: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
```

### Request

```bash
curl localhost:8888/api/v1/test-dubbo/mock -X GET 
```

### Response

```json
{
    "message": "mock success"
}
```

## TODO

We plan use can config custom result in the future. Not only api config way, but also create a match rule.

[Previous](README.md)  
