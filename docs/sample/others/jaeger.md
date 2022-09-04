

#### Tracing with Jaeger


There is a tracing filter, by which we can add tracing function for pixiu


```go
http_filters:
- name: dgp.filters.tracing
  config:
    url: http://127.0.0.1:14268/api/traces
    type: jaeger
```

you can quick start the demo in samples/dubbogo/simple/jaeger for experience
