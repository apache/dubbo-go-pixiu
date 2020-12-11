#!/bin/sh
if [ "$1" != "skip_zookeeper" ];then
	zookeeper_container_id=`docker ps | grep zookeeper  | head -n 1 | awk '{print  $1}'`
	if [ ! -z "$zookeeper_container_id=" ]; then 
		echo "removing old zookeeper!"
        	docker kill $zookeeper_container_id 
        	docker rm $zookeeper_container_id
	fi
	echo "starting zookeeper!"
	docker run -dit --name zookeeper -p 2181:2181 zookeeper:3.4
	echo "zookeeper stated!"
fi


platform=$(uname)
echo "finding go command!"
if ! command -v go >/dev/null 2>&1; then
	echo "go not exists"
	if [ ${platform} == "Darwin" ]; then
        	go_package="go1.14.13.darwin-amd64.tar.gz"
	else
        	go_package="go1.14.13.linux-amd64.tar.gz"
	fi
	if [ ! -d go ]; then
		rm -rf go_sdk.tar.gz
		download_url="https://golang.org/dl/$go_package"
        	echo "downloading go sdk: $download_url"
        	curl --connect-timeout 20 -o go_sdk.tar.gz $download_url
        	if [ ! -f go_sdk.tar.gz ]; then
                	echo "download fail! try using go语言中文网"
                	download_url="https://studygolang.com/dl/golang/$go_package"
                	echo "downloading go sdk: $download_url"
                	curl --connect-timeout 20 --location-trusted -o go_sdk.tar.gz $download_url
        	fi
        	tar -zxf go_sdk.tar.gz
	fi
	go_cmd="go/bin/go"
else
	go_root=`go env GOROOOT`
	go_cmd="$go_root/bin/go"
fi

absolute_go_cmd=`readlink -f $go_cmd`
$go_cmd env -w GOPROXY=https://goproxy.cn,direct
echo "go command is $go_cmd"
echo "starting dubbogo provider!"
if [ -z $CONF_PROVIDER_FILE_PATH ]; then
	export CONF_PROVIDER_FILE_PATH=../../sample/dubbogo/server/config/server.yml
fi
$go_cmd build ../../sample/dubbogo/server/app/

provider_pid=`ps -ef  | grep ../../sample/dubbogo/server/app/ | grep -v 'grep' | awk '{print $2}'`

if [ ! -n $provider_pid ]; then
	echo "pid of old dubbogo provider is $provider_pid, kill it"
	kill -9 $provider_pid
fi
nohup $go_cmd run ../../sample/dubbogo/server/app/ > dubbgo_server.out &
sleep 10
echo "dubbogo provider started!"

## to start proxy

echo "starting proxy!"

echo "the absolute path of go command is $absolute_go_cmd"
cd ../../
make run go_cmd=$absolute_go_cmd

