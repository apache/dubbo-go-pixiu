#!/bin/bash

[ ${DUBBO_REGISTRY_ADDRESS} ] || {
    echo "usage: docker run -itd --name your-name -e DUBBO_REGISTRY_ADDRESS=zookeeper://127.0.0.1:2181 apache/dubbo-go-pixiu:latest"
    echo "    or docker run -itd --name your-name -e DUBBO_REGISTRY_ADDRESS=nacos://127.0.0.1:8848 apache/dubbo-go-pixiu:latest"
    exit
}

echo "current path: $(pwd)"
echo "DUBBO_REGISTRY_ADDRESS = ${DUBBO_REGISTRY_ADDRESS}"

## start from left
#echo ${DUBBO_REGISTRY_ADDRESS#*:}
#echo ${DUBBO_REGISTRY_ADDRESS##*//}

API_CONF="api_config.yaml"
LOG="log.yml"
CONF="conf.yaml"
CONF_ZK="conf-zk.yaml"
CONF_NACOS="conf-nacos.yaml"

ls -al /etc/pixiu/
cd /etc/pixiu/
echo "current path: $(pwd)"

if [ ! -f "$CONF_ZK" ]; then
    echo "error: $CONF_ZK not exists!"
    exit
fi
if [ ! -f "$CONF_NACOS" ]; then
    echo "error: $CONF_NACOS not exists!"
fi

if [ -f "$CONF" ]; then
    rm -rf "$CONF"
fi

if [[ "$DUBBO_REGISTRY_ADDRESS" == zookeeper://* ]]; then
    echo "[[ use zookeeper conf.yaml ]]"
    cp -rf $CONF_ZK ./$CONF
    sed -i -e "s/127.0.0.1:2181/${DUBBO_REGISTRY_ADDRESS##*//}/g" "${CONF}"
fi

if [[ "$DUBBO_REGISTRY_ADDRESS" == nacos://* ]]; then
    echo "[[ use nacos conf.yaml ]]"
    cp -rf $CONF_NACOS ./$CONF
    sed -i -e "s/127.0.0.1:8848/${DUBBO_REGISTRY_ADDRESS##*//}/g" "${CONF}"
fi

while read line; do
    echo $line
done < ${CONF}

cd /app/ && ls -al .
echo "current path: $(pwd)"
./pixiu gateway start -a /etc/pixiu/${API_CONF} -c /etc/pixiu/${CONF} -g /etc/pixiu/${LOG}
