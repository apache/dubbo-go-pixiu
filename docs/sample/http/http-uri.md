# Get the parameter from the URI

> GET request \
> downstream service url is 127.0.0.1:1314/user?name=tc
> downstream service url is 127.0.0.1:1314/user/:name

> POST request \
> downstream service url is 127.0.0.1:1314/user
> request body as below: \
> type User struct { \
	Name string    `json:"name"` \
	Age  int32     `json:"age"`\
	Time time.Time `json:"time"`\
}

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
    resources:
      - path: '/:name'
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
              mappingParams:
                - name: uri.name
                  mapTo: uri.name
              host: "127.0.0.1:1314"
              path: "/user/:name"
          - httpVerb: POST
            enable: true
            timeout: 10s
            inboundRequest:
              requestType: http
            integrationRequest:
              requestType: http
              host: "127.0.0.1:1314"
              path: "/user/"
              mappingParams:
                - name: uri.name
                  mapTo: requestBody.name
                - name: requestBody.age
                  mapTo: requestBody.age
                - name: queryStrings.time
                  mapTo: requestBody.time
```
## Sample Requests
### Passthrough Get request
```bash
curl --request GET 'localhost:8888/api/v1/test-http/user?name=joe2'
```

#### Response

- If exist, will return:

```json
{
    "id": "XVlBz",
    "name": "joe2",
    "age": 20,
    "time": "2021-01-01T00:00:00Z"
}
```

- Not found, return: nil

### Mix Post request
```bash
curl --request POST 'localhost:8888/api/v1/test-http/user/tc2?time=2021-01-01T00:00:00Z' \
--header 'Content-Type: application/json' \
--data-raw '{
    "age": 19
}'
```

#### Response
```json
{
    "id": "XVlBz",
    "name": "tc2",
    "age": 19,
    "time": "2021-01-01T00:00:00Z"
}
```
### Mix GET request
```bash
curl --request GET 'localhost:8888/api/v1/test-http/user/tc'
```

#### Response
```json
{
    "id": "0001",
    "name": "tc",
    "age": 18,
    "time": "2020-12-28T13:38:25.687309+08:00"
}
```


