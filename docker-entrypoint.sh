#!/bin/bash
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
