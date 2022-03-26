module xdsserver

go 1.16

require (
	github.com/apache/dubbo-go-pixiu v0.3.0
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.5
	github.com/envoyproxy/go-control-plane v0.10.1
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/apache/dubbo-go-pixiu => ../../../../..
