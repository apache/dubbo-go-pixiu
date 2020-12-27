# Local example

> Prepare zookeeper env.

## Start provider:

Config env：
- CONF_PROVIDER_FILE_PATH=/XXX/dubbo-go-proxy/samples/dubbogo/server/profiles/dev/server.yml
- APP_LOG_CONF_FILE=/XXX/dubbo-go-proxy/samples/dubbogo/server/profiles/dev/log.yml

Run sample/server/app/server.go

## Start proxy

Config program arguments：
- -c /XXX/dubbo-go-proxy/samples/dubbogo/simple/body/conf.yaml 
- -a /XXX/dubbo-go-proxy/samples/dubbogo/simple/body/api_config.yaml

Run cmd/proxy/proxy.go

You can start proxy, run request.sh under the target folder.

## Other

dubbogo/client is dubbogo generic sample，used for comparison. 