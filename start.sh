#! /bin/bash

args=$1

PROJECT_HOME=`pwd`
PROJECT_HOME=${PROJECT_HOME}"/"

export CONF_CONSUMER_FILE_PATH=${PROJECT_HOME}"conf/client.yml"
export APP_LOG_CONF_FILE=${PROJECT_HOME}"conf/log.yml"

./dubbo-go-proxy ${args}