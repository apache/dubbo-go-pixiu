#!/bin/bash

CMD_PARAMS=""
API_CONF="api_config.yaml"
LOG="log.yml"
CONF="conf.yaml"

cd /etc/pixiu/ &&  ls -al /etc/pixiu/
echo "current path: $(pwd)"

if [ ! -f "$CONF" ]; then
    echo "error: $CONF not exists!"
    exit
fi
CMD_PARAMS="-c /etc/pixiu/${CONF}"

if [ -f "$API_CONF" ]; then
    CMD_PARAMS="$CMD_PARAMS -a /etc/pixiu/${API_CONF}"
fi

if [ -f "$LOG" ]; then
    CMD_PARAMS="$CMD_PARAMS -g /etc/pixiu/${LOG}"
fi

cd /app/ && ls -al .
echo "current path: $(pwd)"
echo "CMD_PARAMS: $CMD_PARAMS"
./dubbo-go-pixiu gateway start $CMD_PARAMS
