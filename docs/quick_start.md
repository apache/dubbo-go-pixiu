目前在内部测试中，仅供开发人员使用

先提两个点：

#### 配置文件 
- aproxy/configs/client.yml 项目配置
- aproxy/configs/api_config.yaml 里面是Api配置

另外两个配置文件是 dubbo 的，考虑把必要的参数整到 client.yml 的 dgp.filters.remote_call 这个 filter 的配置里面

因为现在没有接口配置，暂时接口配置写死在代码

> API 配置

proxy_start.go#beforeStart() 里面配一些自己的 dubbo 接口元数据

#### IDE启动

启动类 cmd/proxy/proxy.go

程序参数
-c /Users/tc/Documents/workspace_2020/aproxy/configs/conf.yaml

环境变量
CONF_CONSUMER_FILE_PATH=/Users/tc/Documents/workspace_2020/aproxy/configs/client.yml;APP_LOG_CONF_FILE=/Users/tc/Documents/workspace_2020/aproxy/configs/log.yml

### sample

#### 发送单参数请求

> 请求头 X-DGP-WAY:dubbo，和配置文件对应

url: http://127.0.0.1:8888/api/v1/test-dubbo/getUserByName

body:
```json
"tiecheng"
```

response:
```json
{
    "id": "3213",
    "name": "tiecheng",
    "age": 25
}
```

url: http://127.0.0.1:8888/api/v1/test-dubbo/user

body:
```json
{
    "id": "3213",
    "name": "tiecheng",
    "age": 25
}
```

response:
```json
{
    "age": 25,
    "id": "3213",
    "name": "tiecheng"
}
```