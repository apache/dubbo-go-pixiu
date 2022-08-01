#
# Licensed to Apache Software Foundation (ASF) under one or more contributor
# license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright
# ownership. Apache Software Foundation (ASF) licenses this file to you under
# the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

PROJECT_NAME = $(shell basename "$(PWD)")
PID = /tmp/.$(PROJECT_NAME).pid
PROJECT_DIR=$(shell pwd)
BASE_DIR := $(PROJECT_DIR)/dist
PIXIU_DIR := $(PIXIU_DIR)
PIXIU_PID = /tmp/.pixiu.pid

SOURCES = $(wildcard $(PROJECT_DIR)/server/app/*.go)
pixiuSources = $(wildcard $(PIXIU_DIR)/cmd/pixiu/*.go)

export GO111MODULE ?= on
export GOPROXY ?= https://goproxy.io,direct
export GOSUMDB ?= sum.golang.org
export GOARCH ?= amd64

export DOCKER_HOST_IP = $(shell hostname)

OS := $(shell uname)
ifeq ($(OS), Linux)
	export GOOS ?= linux
else ifeq ($(OS), Darwin)
	export GOOS ?= darwin
else
	export GOOS ?= windows
endif

ifeq ($(GOOS), windows)
	export EXT_NAME ?= .exe
else
	export EXT_NAME ?=
endif

CGO ?= 1
ifeq ($(DEBUG), true)
	BUILD_TYPE := debug
	GCFLAGS := -gcflags="all=-N -l"
	LCFLAGS :=
else
	BUILD_TYPE := release
	LDFLAGS := "-s -w"
endif

OUT_DIR := $(BASE_DIR)/$(GOOS)_$(GOARCH)
LOG_FILE := $(OUT_DIR)/$(PROJECT_NAME).log
API_CONFIG_PATH := $(OUT_DIR)/pixiuconf/api_config.yaml
CONFIG_PATH := $(OUT_DIR)/pixiuconf/conf.yaml

export APP_LOG_CONF_FILE ?= $(OUT_DIR)/conf/log.yml

.PHONY: all
all: help
help: $(realpath $(firstword $(MAKEFILE_LIST)))
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: Build application's binaries
.PHONY: build
build: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME) config

.PHONY: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME)
$(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME):
	$(info   >  Buiding application binary: $(OUT_DIR)/$(PROJECT_NAME)$(EXT_NAME))
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -o $(OUT_DIR)/server/$(PROJECT_NAME)$(EXT_NAME) $(SOURCES)


## config: Setup config files
.PHONY: config
config:
	$(info   >  Setting up config files)
	@mkdir -p $(OUT_DIR)/server
	@mkdir -p $(OUT_DIR)/pixiuconf
	@mkdir -p $(OUT_DIR)/proto
	@-test -f $(PROJECT_DIR)/server/profiles/dev/log.yml && cat $(PROJECT_DIR)/server/profiles/dev/log.yml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/server/log.yml && echo "  > $(OUT_DIR)/conf/log.yml"
	@-test -f $(PROJECT_DIR)/server/profiles/dev/server.yml && cat $(PROJECT_DIR)/server/profiles/dev/server.yml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/server/server.yml && echo "  > $(OUT_DIR)/conf/server.yml"
	@-test -f $(PROJECT_DIR)/pixiu/api_config.yaml && cat $(PROJECT_DIR)/pixiu/api_config.yaml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" > $(OUT_DIR)/pixiuconf/api_config.yaml && echo "  > $(OUT_DIR)/pixiuconf/api_config.yaml"
	@-test -f $(PROJECT_DIR)/pixiu/conf.yaml && cat $(PROJECT_DIR)/pixiu/conf.yaml | sed "s#\$$HOST_IP#$(DOCKER_HOST_IP)#g" | sed "s#\$$PROJECT_DIR#$(PROJECT_DIR)#g" > $(OUT_DIR)/pixiuconf/conf.yaml && echo "  > $(OUT_DIR)/pixiuconf/conf.yaml"

## docker-up: Shutdown dependency services on docker
.PHONY: docker-up
docker-up:
	$(info   >  Starting dependency services with $(PROJECT_DIR)/docker/docker-compose.yml)
	@-test -f $(PROJECT_DIR)/docker/docker-compose.yml && docker-compose -f $(PROJECT_DIR)/docker/docker-compose.yml up -d

## docker-health-check: check services health on docker
.PHONY: docker-health-check
docker-health-check:
	$(info   >  run docker health check with $(PROJECT_DIR)/docker/docker-health-check.sh)
	@-test -f $(PROJECT_DIR)/docker/docker-health-check.sh && bash -f $(PROJECT_DIR)/docker/docker-health-check.sh

## docker-down: Shutdown dependency services on docker
.PHONY: docker-down
docker-down:
	$(info   >  Stopping dependency services with $(PROJECT_DIR)/docker/docker-compose.yml)
	@-test -f $(PROJECT_DIR)/docker/docker-compose.yml && docker-compose -f $(PROJECT_DIR)/docker/docker-compose.yml down

## clean: Clean up the output and the binary of the application
.PHONY: clean
clean: stop
	$(info   >  Cleanning up $(OUT_DIR))
	@-rm -rf $(OUT_DIR)
	@-cat $(PID) | awk '{print $1}' | xargs kill -9
	@-cat $(PIXIU_PID) | awk '{print $1}' | xargs kill -9

## start: Start the application (for server)
.PHONY: start
start: export DUBBO_GO_CONFIG_PATH ?= $(OUT_DIR)/server/server.yml
start: build
	$(info   >  Starting application $(PROJECT_NAME), output is redirected to $(LOG_FILE))
	@-$(OUT_DIR)/server/$(PROJECT_NAME)$(EXT_NAME) > $(LOG_FILE) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"


## start: Start the application (for server) foreground
.PHONY: startServer
startServer: export DUBBO_GO_CONFIG_PATH ?= $(OUT_DIR)/server/server.yml
startServer: build
	$(info   >  Starting application $(PROJECT_NAME), output is redirected to $(LOG_FILE))
	@-$(OUT_DIR)/server/$(PROJECT_NAME)$(EXT_NAME)  2>&1

## run: Run the application (for client)
.PHONY: run
run: buildPixiu
	$(info   >  Running application PIXIU, output is redirected to $(LOG_FILE))

## buildPixiu: start pixiu
.PHONY: buildPixiu
buildPixiu:
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -o $(OUT_DIR)/dubbo-go-pixiu$(EXT_NAME) $(pixiuSources)
	@-$(OUT_DIR)/dubbo-go-pixiu$(EXT_NAME) gateway start -a $(API_CONFIG_PATH) -c $(CONFIG_PATH) & echo $$! > $(PIXIU_PID)
	@cat $(PIXIU_PID) | sed "/^/s/^/  \> PIXIU_PID: /"


## startPixiu: start pixiu foreground
.PHONY: startPixiu
startPixiu:
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -o $(OUT_DIR)/dubbo-go-pixiu$(EXT_NAME) $(pixiuSources)
	@-$(OUT_DIR)/dubbo-go-pixiu$(EXT_NAME) gateway start -a $(API_CONFIG_PATH) -c $(CONFIG_PATH)

## stop: Stop running the application (for server)
.PHONY: stop
stop:
	$(info   >  Stopping the application $(PROJECT_NAME))
	@cat $(PID) | sed "/^/s/^/  \>  Killing PID: /"
	@-kill `cat $(PID)` 2>/dev/null || true

## integration: Run integration test for this application
.PHONY: integration
integration:
	$(info   >  Running integration test for application $(PROJECT_NAME))
	@go clean -testcache
	@go test -tags integration -v $(PROJECT_DIR)/test/...
