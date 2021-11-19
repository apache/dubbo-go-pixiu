module github.com/apache/dubbo-go-pixiu

go 1.14

require (
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-go v1.5.7
	github.com/apache/dubbo-go-hessian2 v1.9.5
	github.com/creasty/defaults v1.5.2
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.4
	github.com/dubbogo/go-zookeeper v1.0.3
	github.com/dubbogo/gost v1.11.14
	github.com/emirpasic/gods v1.12.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gin-gonic/gin v1.7.4
	github.com/go-resty/resty/v2 v2.3.0
	github.com/gogo/protobuf v1.3.2
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0 // indirect
	github.com/jhump/protoreflect v1.9.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/nacos-group/nacos-sdk-go v1.0.8
	github.com/opentrx/seata-golang/v2 v2.0.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.29.0 // indirect
	github.com/shirou/gopsutil v3.21.3+incompatible // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	go.etcd.io/etcd/api/v3 v3.5.0-alpha.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/prometheus v0.21.0
	go.opentelemetry.io/otel/metric v0.21.0
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/sdk/export/metric v0.21.0
	go.opentelemetry.io/otel/sdk/metric v0.21.0
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	go.uber.org/zap v1.17.0
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	vimagination.zapto.org/byteio v0.0.0-20200222190125-d27cba0f0b10
	vimagination.zapto.org/memio v0.0.0-20200222190306-588ebc67b97d // indirect
)

replace github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0
