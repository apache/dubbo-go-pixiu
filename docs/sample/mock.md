# Mock request

## simple

### config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-dubbo/mock'
    type: restful
    description: mock
    methods:
      - httpVerb: GET
        onAir: true
        mock: true
        timeout: 100s
        inboundRequest:
          requestType: http
```

### request

```bash
curl localhost:8888/api/v1/test-dubbo/mock -X GET 
```

### response

```json
{
    "message": "mock success"
}
```

## TODO

We plan use can config custom result in the future. Not only api config way, but also create a match rule.  