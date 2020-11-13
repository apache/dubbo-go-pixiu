# 本地例子

> 自行准备好 zookeeper 环境

## 启动服务提供方

配置好下面两个环境变量：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/sample/dubbogo/server/config/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/sample/dubbogo/server/config/log.yml

运行 sample/server/app/server.go

## 启动 proxy

配置好下面的环境变量：
- -c /XXX/dubbo-go-proxy/sample/dubbogo/proxy/conf.yaml 
- -a /XXX/dubbo-go-proxy/sample/dubbogo/proxy/api_config.yaml

运行 cmd/proxy/proxy.go

因为没有配置 dubbogo 的配置文件，所以启动的时候出现日志如下：
```bash
2020/11/06 10:29:09 [InitLog] warn: ioutil.ReadFile(file:./conf/log.yml) = error:open ./conf/log.yml: no such file or directory
2020/11/06 10:29:09 [InitLog] warn: log configure file name is nil
2020/11/06 10:29:28 [consumerInit] application configure(consumer) file name is nil
2020/11/06 10:29:28 [providerInit] application configure(provider) file name is nil
```

在浏览器或其它工具输入 `http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc` ， 可以得到结果

```json
{
    "age": 18,
    "id": "0001",
    "name": "tc"
}
```

正常调用 proxy 日志：

```bash
2020-11-08T11:02:47.106+0800    DEBUG   dubbo/dubbo_invoker.go:121      result.Err: <nil>, result.Rest: 0xc0001ad560
2020-11-08T11:02:47.106+0800    DEBUG   proxy/proxy.go:172      [makeDubboCallProxy] result: 0xc0001ad560, err: <nil>
2020-11-08T11:02:47.106+0800    DEBUG   dubbo/client.go:146     [dubbogo proxy] dubbo client resp:map[age:18 id:0001 name:tc time:<nil>]
2020-11-08T11:02:47.106+0800    INFO    filter/logger_filter.go:44      [dubboproxy go] [UPSTREAM] receive request | 200 | 349.835096ms | GET | /api/v1/test-dubbo/user | 
```

server 日志：

```bash
Req GetUserByName data:"tc"
Req GetUserByName result:&main.User{Id:"0001", Name:"tc", Age:18, Time:time.Time{wall:0xbfe201c6b219ce40, ext:5594874, loc:(*time.Location)(0x1c3b100)}}
2020-11-08T14:33:31.676+0800    INFO    dubbo/listener.go:196   got session:session {server:TCP_SERVER:2:192.168.0.113:20000<->192.168.0.113:52871}, Read Bytes: 0, Write Bytes: 0, Read Pkgs: 0, Write Pkgs: 0
```

## 其它

dubbogo/client 是给泛化调用的客户端使用例子，方便对比和学习泛化调用