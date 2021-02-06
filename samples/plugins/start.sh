#!/usr/bin/env sh
DIR=$(cd "$(dirname "$0")" && pwd)
echo "$DIR"

PLUGIN_DIR=${DIR}/plugin
OUT_DIR=${DIR}/out

for f in "${PLUGIN_DIR}"/*.go
do
  echo "$f"
done

go build -o "${OUT_DIR}" -buildmode=plugin "${PLUGIN_DIR}"
echo 'build successful!'