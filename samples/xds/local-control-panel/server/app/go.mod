module xdsserver

go 1.16

require (
	github.com/apache/dubbo-go-pixiu v0.3.0
	github.com/dubbo-go-pixiu/pixiu-api v0.1.6-0.20220427143710-d2e48e546d2c
	github.com/envoyproxy/go-control-plane v0.10.1
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/apache/dubbo-go-pixiu => ../../../../..
replace github.com/dubbo-go-pixiu/pixiu-api => ../../../../../../pixiu-api
