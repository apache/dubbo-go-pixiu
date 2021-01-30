#!/bin/bash -x

DIR=$(cd $(dirname $0) && pwd)

echo $DIR

PLUGIN_DIR=${DIR}/plugin

for f in ${PLUGIN_DIR}/*.go; do
  echo ${f}
done

OUT_DIR=${DIR}/out

go build -o ${OUT_DIR} -buildmode=plugin ${PLUGIN_DIR}

echo 'build successful!'