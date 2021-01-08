# Local example

> Prepare zookeeper env.

upload /XXX/dubbo-go-proxy/sample/admin/proxy/api_config.yml to server where zookeeper are running on, then execute below command to set api_config value
```
./zkCli.sh -server 127.0.0.1:2181 set /dubbo/config/proxy/api_config "`cat api_config.yml`"
```

## Start provider:

config env：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/sample/admin/server/config/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/sample/admin/server/config/log.yml

run sample/server/app/server.go

## Start proxy

config program arguments：
- -c /XXX/dubbo-go-proxy/sample/admin/proxy/conf.yaml 
- -a /XXX/dubbo-go-proxy/sample/admin/proxy/api_config.yaml

import config_center dependency according to api_meta_config.protocol in config.yml

```go
_ "github.com/apache/dubbo-go/config_center/zookeeper"
```



run cmd/proxy/proxy.go

curl `http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc` ， the result:

```json
{
    "age": 18,
    "id": "0001",
    "name": "tc"
}
```