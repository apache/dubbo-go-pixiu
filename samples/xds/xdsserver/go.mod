module xdsserver

go 1.16

require (
	github.com/apache/dubbo-go-pixiu v0.3.0
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.5-0.20211220151221-87949cfcdf4e
	github.com/envoyproxy/go-control-plane v0.10.1
	//github.com/envoyproxy/protoc-gen-validate v0.1.0
	//github.com/golang/protobuf v1.5.2
	//github.com/google/go-cmp v0.5.6
	//github.com/stretchr/testify v1.7.0
	//go.opentelemetry.io/proto/otlp v0.7.0
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/dubbogo/dubbo-go-pixiu-filter => ../../../../dubbo-go-pixiu-filter

replace github.com/apache/dubbo-go-pixiu => ../../..
