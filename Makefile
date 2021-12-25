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

EXE=pixiu
PID := /tmp/$(EXE).pid
HOSTNAME := $(shell hostname)
PROJECT_NAME := $(shell basename "$(PWD)")
CUR_MK_FILE := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_DIR := $(patsubst %/, %, $(dir $(CUR_MK_FILE)))
PIXIU_PATH := cmd/pixiu/
MAIN_PATH := $(PROJECT_DIR)/$(PIXIU_PATH)
SOURCES := $(wildcard $(MAIN_PATH)/*.go)

export CGO_ENABLED ?= 0
export GO111MODULE ?= on
export GOPROXY ?= https://goproxy.io,direct
export GOSUMDB ?= sum.golang.org
export GOARCH ?= amd64

VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -s -w"
GCFLAGS := -gcflags="all=-N -l"

OSTYPE := $(shell uname)
ifeq ($(OSTYPE), Linux)
	export GOOS ?= linux
else ifeq ($(OSTYPE), Darwin)
	export GOOS ?= darwin
else
	export GOOS ?= windows
endif

ifeq ($(GOOS), windows)
	EXE = $(EXE).exe
else
	EXE ?=
endif

export CONF_PROVIDER_FILE_PATH ?= conf/server.yml
export APP_LOG_CONF_FILE ?= conf/log.yml
API_CONF_PATH := conf/api_config.yaml
CONF_PATH := conf/conf.yaml

#$(info PROJECT_DIR = $(PROJECT_DIR))
#$(info APP_LOG_CONF_FILE = $(APP_LOG_CONF_FILE))
#$(info API_CONF_PATH = $(API_CONF_PATH))
#$(info CONF_PATH = $(CONF_PATH))

###############################################################
#action
###############################################################

.PHONY: all build clean run image check cover lint

all: check build

check:
	@go fmt $(SOURCES)
	@go vet $(SOURCES)

build:
	$(info > buiding application binary: $(PROJECT_DIR)/$(EXE) ...)
	@cd $(MAIN_PATH) && CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) $(LDFLAGS) -o $(PROJECT_DIR)/$(EXE) $(SOURCES)

rebuild: clean build

run: build
	@$(PROJECT_DIR)/$(EXE) gateway start -a $(API_CONF_PATH) -c $(CONF_PATH)

test:
	sh $(PROJECT_DIR)/before_ut.sh
	go test ./pkg/... -coverprofile=coverage.txt -covermode=atomic

samples:
	sh $(PROJECT_DIR)/start_integrate_test.sh

cover:
	@go test $(SOURCES) -coverprofile coverprofile.txt
	@go tool cover -html=coverprofilehtml.txt

lint:
	golangci-lint run --enable-all

image:
	@docker build \
		-t apache/dubbo-go-$(EXE):latest \
		-t apache/dubbo-go-$(EXE):$(VERSION) \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile --no-cache .

clean:
	@go clean
	@if [ -f $(EXE) ] ; then rm -f $(EXE) ; fi
	@echo "Clean $(PROJECT_DIR)/$(EXE) end."


license-check-util:
	go install github.com/lsm-dev/license-header-checker/cmd/license-header-checker@latest


license-check:
	license-header-checker -v -a -r -i vendor -i .github/actions /tmp/tools/license/license.txt . go