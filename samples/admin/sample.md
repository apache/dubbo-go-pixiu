# Local example

## prepare etcd config

run etcd local or in docker, then use etcdctl to set api config

- export ETCDCTL_API=3
- if use docker, run `docker cp api_config.yaml mycontainer:/path` to copy file to docker
- run `cat api_config.yaml | etcdctl put "/proxy/config/api"` to set api config


## Start provider:

config env：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/samples/admin/server/config/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/samples/admin/server/config/log.yml

run samples/admin/server/app/server.go

## Start proxy

config program arguments：
- -c /XXX/dubbo-go-proxy/samples/admin/proxy/conf.yaml 
- -a /XXX/dubbo-go-proxy/samples/admin/proxy/api_config.yaml

to replace -a program arguments, should put api_meta_config in conf.yaml 

```yaml
  api_meta_config:
    address: "127.0.0.1:2379"
    api_config_path: "/pixiu/config/api"
```


run cmd/proxy/proxy.go

## 

curl "http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc" ， the result:

```json
{
    "age": 18,
    "id": "0001",
    "name": "tc"
}
```

## Start admin

clone from https://github.com/dubbogo/pixiu-admin

cd pixiu-admin

run cmd/admin/admin.go

config program arguments：
- -c /xx/pixiu-admin/configs/admin_config.yaml