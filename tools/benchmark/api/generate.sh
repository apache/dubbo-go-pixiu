export GO111MODULE="on"
export GOPROXY="http://goproxy.io"
go install github.com/dubbogo/tools/cmd/protoc-gen-go-triple@latest
protoc --go_out=. --go-triple_out=. samples_api.proto