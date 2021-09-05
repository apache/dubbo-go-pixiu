#! /bin/bash
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

if [ -z "$1" ]; then
  echo 'Provide test directory please, like : ./start.sh body '
  exit
fi

ACTION=$1

if [ $ACTION = 'help' ]; then
  echo "dubbo-go-pixiu start helper"
  echo "./start.sh action project"
  echo "hint:"
  echo "./start.sh prepare body for prepare config file and up docker in body project"
  echo "./start.sh startPixiu body for start dubbo or http server in body project"
  echo "./start.sh startServer body for start pixiu in body project"
  echo "./start.sh startTest body for start unit test in body project"
  echo "./start.sh clean body for clean"
fi


P_DIR=$(pwd)/$2

PIXIU_DIR=$(dirname $(dirname $(dirname "$PWD")))

if [ $ACTION = 'prepare' ]; then
    # prepare config file
    echo "prepare config file and docker, please remember clean finally"
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile config
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile docker-up
    sleep 0.5
elif [ $ACTION = 'startServer' ]; then
    echo "start server"
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile startServer
    sleep 0.5
elif [ $ACTION = 'startPixiu' ]; then
    echo "start pixiu"
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile startPixiu
elif [ $ACTION = 'startTest' ]; then
    echo "start unit test"
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile integration
    result=$?
elif [ $ACTION = 'clean' ]; then
    echo "cleanup for config file and docker"
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile clean
    make PROJECT_DIR=$P_DIR PIXIU_DIR=$PIXIU_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../../../igt/Makefile docker-down
else
  echo "wrong action"
fi