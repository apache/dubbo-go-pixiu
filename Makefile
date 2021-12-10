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

CUR_MK_FILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_PATH := $(patsubst %/, %, $(dir $(CUR_MK_FILE)))
PIXIU_PATH := cmd/pixiu/
MAIN_PATH := $(CURRENT_PATH)/$(PIXIU_PATH)

EXE := pixiu

API_CONF_PATH:=conf/api_config.yaml
CONF_PATH:=conf/conf.yaml

$(info API_CONF_PATH = $(API_CONF_PATH))
$(info CONF_PATH = $(CONF_PATH))

os := $(shell go env GOOS)
ifeq (windows,$(os))
	EXE = $(EXE).exe
endif

build:
	cd $(MAIN_PATH) && go build  -o $(CURRENT_PATH)/$(EXE) *.go

run: build
	./$(EXE) gateway start -a $(API_CONF_PATH) -c $(CONF_PATH)

license-check-util:
	go install github.com/lsm-dev/license-header-checker/cmd/license-header-checker@latest

license-check:
	license-header-checker -v -a -r -i vendor -i .github/actions /tmp/tools/license/license.txt . go

test:
	sh before_ut.sh
	go test ./pkg/... -coverprofile=coverage.txt -covermode=atomic

integrate-test:
	sh start_integrate_test.sh

clean:
