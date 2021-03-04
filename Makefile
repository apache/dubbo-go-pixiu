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
