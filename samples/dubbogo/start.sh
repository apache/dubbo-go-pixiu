#!/bin/bash -x

DIR=$(cd $(dirname $0) && pwd )

echo $DIR

./proxy -c ${DIR}/simple/$1/conf.yaml -a ${DIR}/simple/$1/api_config.yaml