cur_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
currentPath := $(patsubst %/, %, $(dir $(cur_mkfile)))
$(info cur_makefile_path1=$(currentPath))
proxyPath := /cmd/proxy/
mainPath := $(currentPath)$(proxyPath)
$(info $(mainPath))

targetName := dubbo-go-proxy
go_cmd:=${go_cmd}
ifeq ("",$(go_cmd))
        go_cmd = go
endif
os := $(shell $(go_cmd) env GOOS)
$(info os is $(os))
ifeq (windows,$(os))
	targetName = dubbo-go-proxy.exe
endif
exe := $(mainPath)$(targetName)
$(info path of exe is $(exe))
$(info go command is $(go_cmd))
build:
	cd $(mainPath) && $(go_cmd) build  -o $(targetName)

run: build
	cp $(exe) $(currentPath) && ./dubbo-go-proxy start
