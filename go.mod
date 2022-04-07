module github.com/apache/dubbo-go-pixiu

go 1.15

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.1-0.20220107110037-4496cef73dba
	github.com/MicahParks/keyfunc v1.0.0
	github.com/Shopify/sarama v1.19.0
	github.com/alibaba/sentinel-golang v1.0.4
	github.com/apache/dubbo-getty v1.4.7-rc2
	github.com/apache/dubbo-go-hessian2 v1.11.0
	github.com/cch123/supermonkey v1.0.0
	github.com/creasty/defaults v1.5.2
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.5
	github.com/dubbogo/go-zookeeper v1.0.4-0.20211212162352-f9d2183d89d5
	github.com/dubbogo/gost v1.11.22
	github.com/dubbogo/grpc-go v1.42.7
	github.com/dubbogo/triple v1.1.7
	github.com/envoyproxy/go-control-plane v0.10.1
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gin-gonic/gin v1.7.4
	github.com/go-errors/errors v1.0.1
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
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.1
	go.etcd.io/etcd/api/v3 v3.5.1
	go.opentelemetry.io/otel v1.6.1
	go.opentelemetry.io/otel/exporters/jaeger v1.6.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.6.1
	go.opentelemetry.io/otel/exporters/prometheus v0.21.0
	go.opentelemetry.io/otel/metric v0.21.0
	go.opentelemetry.io/otel/sdk v1.6.1
	go.opentelemetry.io/otel/sdk/export/metric v0.21.0
	go.opentelemetry.io/otel/sdk/metric v0.21.0
	go.opentelemetry.io/otel/trace v1.6.1
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
	golang.org/x/sys v0.0.0-20220403020550-483a9cbc67c0 // indirect
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
	vimagination.zapto.org/byteio v0.0.0-20200222190125-d27cba0f0b10
	vimagination.zapto.org/memio v0.0.0-20200222190306-588ebc67b97d // indirect
)

//github.com/go-logr/logr => github.com/go-logr/logr v1.0.0
replace k8s.io/apimachinery => k8s.io/apimachinery v0.23.5
