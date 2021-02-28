#!/bin/bash -x

DIR=$(cd $(dirname $0) && pwd )

echo $DIR

./pixiu -c ${DIR}/$1/conf.yaml -a ${DIR}/$1/api_config.yaml