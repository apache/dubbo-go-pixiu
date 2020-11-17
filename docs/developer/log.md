# Log

How to view logs in dubbo-go-proxy.

## DEBUG log

** dubbo 请求日志**

```bash
2020-11-14T20:59:58.841+0800    DEBUG   dubbo/dubbo.go:149      [dubbo-go-proxy] invoke, method:GetUserByName, types:[java.lang.String], reqData:[tc]
```

** dubbo 响应日志**

```bash
2020-11-14T20:59:58.843+0800    DEBUG   remote/call.go:96       [dubbo-go-proxy] resp : {map[age:18 i_d:0001 name:tc]}
```


