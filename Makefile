cur_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
currentPath := $(patsubst %/, %, $(dir $(cur_mkfile)))
$(info cur_makefile_path1=$(currentPath))
proxyPath := /cmd/proxy/
mainPath := $(currentPath)$(proxyPath)
$(info $(mainPath))

targetName := dubbo-go-proxy
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
	cp $(exe) $(currentPath) && ./dubbo-go-proxy start
