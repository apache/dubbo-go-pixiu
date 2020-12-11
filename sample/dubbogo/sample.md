# Local example

> prepare zookeeper env.

## Start provider:

config env：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/sample/dubbogo/server/config/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/sample/dubbogo/server/config/log.yml

run sample/server/app/server.go

## Start proxy

config program arguments：
- -c /XXX/dubbo-go-proxy/sample/dubbogo/proxy/conf.yaml 
- -a /XXX/dubbo-go-proxy/sample/dubbogo/proxy/api_config.yaml

run cmd/proxy/proxy.go

curl `http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc` ， the result:

```json
{
    "age": 18,
    "id": "0001",
    "name": "tc"
}
```

the proxy log：

```bash
2020-11-08T11:02:47.106+0800    DEBUG   dubbo/dubbo_invoker.go:121      result.Err: <nil>, result.Rest: 0xc0001ad560
2020-11-08T11:02:47.106+0800    DEBUG   proxy/proxy.go:172      [makeDubboCallProxy] result: 0xc0001ad560, err: <nil>
2020-11-08T11:02:47.106+0800    DEBUG   dubbo/client.go:146     [dubbogo proxy] dubbo client resp:map[age:18 id:0001 name:tc time:<nil>]
2020-11-08T11:02:47.106+0800    INFO    filter/logger_filter.go:44      [dubboproxy go] [UPSTREAM] receive request | 200 | 349.835096ms | GET | /api/v1/test-dubbo/user | 
```

server log：

```bash
Req GetUserByName data:"tc"
Req GetUserByName result:&main.User{Id:"0001", Name:"tc", Age:18, Time:time.Time{wall:0xbfe201c6b219ce40, ext:5594874, loc:(*time.Location)(0x1c3b100)}}
2020-11-08T14:33:31.676+0800    INFO    dubbo/listener.go:196   got session:session {server:TCP_SERVER:2:192.168.0.113:20000<->192.168.0.113:52871}, Read Bytes: 0, Write Bytes: 0, Read Pkgs: 0, Write Pkgs: 0
```

## other

dubbogo/client is dubbogo generic sample，used for comparison. 