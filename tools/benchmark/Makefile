clone:
	git clone https://github.com/apache/dubbo-go-pixiu.git source

build:
	@export GO111MODULE=on
	@cd source && go mod vendor && \
	 if  [ ! -d "../dist" ]; then mkdir ../dist; fi && \
	 go build -o ../dist/pixiu cmd/pixiu/*.go

prepare: clone build
