cur_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
currentPath := $(patsubst %/, %, $(dir $(cur_mkfile)))
$(info cur_makefile_path1=$(currentPath))
proxyPath := /cmd/proxy/
mainPath := $(currentPath)$(proxyPath)
$(info $(mainPath))

targetName := dubbo-go-proxy

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
	targetName = dubbo-go-proxy.exe
endif
exe := $(mainPath)$(targetName)
$(info path of exe is $(exe))
build:
	cd $(mainPath) && go build  -o $(targetName)

run: build
	cp $(exe) $(currentPath) && ./dubbo-go-proxy start -a $(api-config-path) -c $(config-path)
