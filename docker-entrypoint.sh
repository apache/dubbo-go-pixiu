#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

set -e

CONF_DIR="/etc/pixiu"
APP_BIN="/app/dubbo-go-pixiu"

API_CONF_FILE="$CONF_DIR/api_config.yaml"
LOG_FILE="$CONF_DIR/log.yml"
CONF_FILE="$CONF_DIR/conf.yaml"

echo "Checking configurations in $CONF_DIR..."
ls -al "$CONF_DIR"

if [ ! -f "$CONF_FILE" ]; then
    echo "Error: Main configuration file not found at $CONF_FILE"
    exit 1
fi

CMD_PARAMS="-c ${CONF_FILE}"

if [ -f "$API_CONF_FILE" ]; then
    CMD_PARAMS="$CMD_PARAMS -a ${API_CONF_FILE}"
fi

if [ -f "$LOG_FILE" ]; then
    CMD_PARAMS="$CMD_PARAMS -g ${LOG_FILE}"
fi

echo "Starting Pixiu Gateway..."
echo "Binary: $APP_BIN"
echo "Parameters: $CMD_PARAMS"

exec "$APP_BIN" gateway start $CMD_PARAMS