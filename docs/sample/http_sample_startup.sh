#!/bin/sh
if [ "$1" != "skip_zookeeper" ]; then
  zookeeper_container_id=$(docker ps | grep zookeeper | head -n 1 | awk '{print  $1}')
  if [ -n "$zookeeper_container_id" ]; then
    echo "removing old zookeeper!"
    docker kill "$zookeeper_container_id"
    docker rm "$zookeeper_container_id"
  fi
  echo "starting zookeeper!"
  docker run -dit --name zookeeper -p 2181:2181 zookeeper:3.4.14
  echo "zookeeper stated!"
fi

go env -w GOPROXY=https://goproxy.cn,direct

echo "starting dubbogo provider!"
if [ -z "$CONF_PROVIDER_FILE_PATH" ]; then
  export CONF_PROVIDER_FILE_PATH=../../sample/http/server/config/server.yml
fi
go build ../../sample/http/server/app/

provider_pid=$(ps -ef | grep ../../sample/http/server/app/ | grep -v 'grep' | awk '{print $2}')

if [ -n "$provider_pid" ]; then
  echo "pid of old dubbogo provider is $provider_pid, kill it"
  kill -9 "$provider_pid"
fi
nohup go run ../../sample/http/server/app/ >http_server.out &
sleep 10
echo "dubbogo provider started!"

## to start pixiu

echo "starting proxy!"

cd ../../
make run config-path=sample/http/proxy/conf.yaml api-config-path=sample/http/proxy/api_config.yaml
