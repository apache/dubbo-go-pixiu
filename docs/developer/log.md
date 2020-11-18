# Log

How to view logs in dubbo-go-proxy.

## DEBUG log

**dubbo request log**

```bash
2020-11-17T11:31:05.716+0800    DEBUG   dubbo/dubbo.go:150      [dubbo-go-proxy] dubbo invoke, method:GetUserByName, types:[java.lang.String], reqData:[tiecheng]
```

**dubbo response log**

```bash
2020-11-17T11:31:05.718+0800    DEBUG   dubbo/dubbo.go:160      [dubbo-go-proxy] dubbo client resp:map[age:88 iD:3213 name:tiecheng time:<nil>]
```

