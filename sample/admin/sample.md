# Local example

## prepare etcd config

run etcd local or in docker, then use etcdctl to set api config

- export ETCDCTL_API=3
- if use docker, run `docker cp api_config.yaml mycontainer:/path` to copy file to docker
- run `cat api_config.yaml | etcdctl put "/proxy/config/api"` to set api config


## Start provider:

config env：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/sample/admin/server/config/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/sample/admin/server/config/log.yml

run sample/admin/server/app/server.go

## Start proxy

config program arguments：
- -c /XXX/dubbo-go-proxy/sample/admin/proxy/conf.yaml 

to replace -c program arguments, should put api_meta_config in conf.yaml 

```yaml
  api_meta_config:
    address: "127.0.0.1:2379"
    api_config_path: "/proxy/config/api"
```


run cmd/proxy/proxy.go

## 

curl `http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc` ， the result:

```json
{
    "age": 18,
    "id": "0001",
    "name": "tc"
}
```

## Start admin

run cmd/admin/admin.go

- run cmd `curl 127.0.0.1:8080/config/api` to check current proxy api config
- modify api_config.yaml content, such as change path `test-dubbo/user` to `test-dubbo/user_new`
- run cmd `curl "127.0.0.1:8080/config/api/set" -X POST --data-binary "@/xx/xx/dubbo-go-proxy/sample/admin/proxy/api_config.yaml"` to modify proxy api config
- then, `curl "http://127.0.0.1:8888/api/v1/test-dubbo/user?name=tc"` ， the result will be 404
- run `curl "http://127.0.0.1:8888/api/v1/test-dubbo/user_new?name=tc"` get the correct result
