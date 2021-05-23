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

cur_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
currentPath := $(patsubst %/, %, $(dir $(cur_mkfile)))
$(info cur_makefile_path1=$(currentPath))
pixiuPath := /cmd/pixiu/
mainPath := $(currentPath)$(pixiuPath)
$(info $(mainPath))

targetName := dubbo-go-pixiu

api-config-path:=${api-config}
ifeq ("",$(api-config-path))
        api-config-path = configs/api_config.yaml
endif

config-path:=${config-path}
ifeq ("",$(config-path))
        config-path = configs/conf.yaml
endif

os := $(shell go env GOOS)
$(info os is $(os))
ifeq (windows,$(os))
	targetName = dubbo-go-pixiu.exe
endif
exe := $(mainPath)$(targetName)
$(info path of exe is $(exe))
build:
	cd $(mainPath) && go build  -o $(targetName)

run: build
	cp $(exe) $(currentPath) && ./dubbo-go-pixiu start -a $(api-config-path) -c $(config-path)
