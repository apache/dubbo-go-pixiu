module github.com/apache/dubbo-go-pixiu

go 1.15

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0
	github.com/MicahParks/keyfunc v1.0.0
	github.com/Shopify/sarama v1.19.0
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-go-hessian2 v1.10.0
	github.com/creasty/defaults v1.5.2
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.4
	github.com/dubbogo/go-zookeeper v1.0.4-0.20211212162352-f9d2183d89d5
	github.com/dubbogo/gost v1.11.22
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gogo/protobuf v1.3.2
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/golang/protobuf v1.5.2
	github.com/jhump/protoreflect v1.9.0
	github.com/mercari/grpc-http-proxy v0.1.2
	github.com/mitchellh/mapstructure v1.4.3
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/opentrx/seata-golang/v2 v2.0.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.29.0 // indirect
	github.com/shirou/gopsutil v3.21.3+incompatible // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	go.etcd.io/etcd/api/v3 v3.5.1
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/prometheus v0.21.0
	go.opentelemetry.io/otel/metric v0.21.0
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/sdk/export/metric v0.21.0
	go.opentelemetry.io/otel/sdk/metric v0.21.0
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	vimagination.zapto.org/byteio v0.0.0-20200222190125-d27cba0f0b10
	vimagination.zapto.org/memio v0.0.0-20200222190306-588ebc67b97d // indirect
)

replace github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0
