# Get the parameter from the query

> GET request
> downstream service url is 127.0.0.1:1314/user?name=tc

## Simple Demo

### Api Config

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-http/user'
    type: restful
    description: user
    methods:
      - httpVerb: GET
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
curl localhost:8888/api/v1/test-http/user?name=tc -X GET 
```

### Response

- If exist, will return:

```bash
{
    "id": "0001",
    "name": "tc",
    "age": 18,
    "time": "2020-11-23T10:58:29.494108+08:00"
}
```

- Not found, return: nil

## Mapping

```yaml
name: proxy
description: proxy sample
resources:
  - path: '/api/v1/test-http/user'
    type: restful
    description: user
    methods:
      - httpVerb: GET
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
          queryStrings:
            - name: name
              required: true
        integrationRequest:
          requestType: http
          mappingParams:
            - name: queryStrings.username
              mapTo: queryStrings.name
          host: "127.0.0.1:1314"
```

### Request

```bash
curl localhost:8888/api/v1/test-http/user?username=tc -X GET 
```

### Response

- If exist, will return:

```bash
{
    "id": "0001",
    "name": "tc",
    "age": 18,
    "time": "2020-11-23T10:58:29.494108+08:00"
}
```

- Not found, return: nil

