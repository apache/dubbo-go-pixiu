IMAGE = ccr.ccs.tencentyun.com/pixiu/demo-app#$(your_repo)/$(namespace)/$(image_name)
TAG = dubbo-go-sample-triple-1.0
APPNAME = dubbo-go-app

build-amd64-app:
	GOOS=linux GOARCH=amd64 go build -o build/app ./cmd

build: proto-gen tidy build-amd64-app
	cp ./conf/dubbogo.yaml ./build/dubbogo.yaml
	docker build --platform linux/amd64 ./build  -t ${IMAGE}:${TAG}
	docker push ${IMAGE}:${TAG}
	make clean

buildx-publish: proto-gen tidy build-amd64-app
	cp ./conf/dubbogo.yaml ./build/dubbogo.yaml
	docker buildx build \
    	 --platform linux/amd64 \
    	 -t ${IMAGE}:${TAG} \
    	 ./build --push
	make clean

remove:
	helm uninstall ${APPNAME}

deploy:
	helm install ${APPNAME} ./chart/app

deploy-nacos-env:
	helm install nacos ./chart/nacos_env

remove-nacos-env:
	helm uninstall nacos

proto-gen:
	protoc --go_out=./api --go-triple_out=./api ./api/api.proto

clean:
	rm ./build/dubbogo.yaml
	rm ./build/app

tidy:
	go mod tidy

test:
	go test ./...
