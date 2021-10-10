

#### Tracing with Jaeger


there is tracing filter, we can use it to add tracing function for pixiu

```go
http_filters:
- name: dgp.filters.tracing
  config:
    url: http://127.0.0.1:14268/api/traces
    type: jaeger
```

you can quick start the demo in samples/dubbogo/simple/jaeger for experience