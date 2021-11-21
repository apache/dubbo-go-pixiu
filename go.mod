module github.com/apache/dubbo-go-pixiu

go 1.14

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-rc4-1.0.20211119124423-451d10aaa45e
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/apache/dubbo-go v1.5.7-rc2
	github.com/apache/dubbo-go-hessian2 v1.9.5
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.4
	github.com/dubbogo/go-zookeeper v1.0.3
	github.com/dubbogo/gost v1.11.20-0.20211116110728-26777ca61b4a
	github.com/dubbogo/grpc-go v1.42.6-triple
	github.com/dubbogo/triple v1.1.6-0.20211119123944-4ad68a0d048e
	github.com/emirpasic/gods v1.12.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/consul/api v1.5.0
	github.com/mercari/grpc-http-proxy v0.1.2
	github.com/mitchellh/mapstructure v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.29.0 // indirect
	github.com/shirou/gopsutil v3.21.3+incompatible // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	go.etcd.io/etcd/api/v3 v3.5.1
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/prometheus v0.21.0
	go.opentelemetry.io/otel/metric v0.21.0
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/sdk/export/metric v0.21.0
	go.opentelemetry.io/otel/sdk/metric v0.21.0
	go.uber.org/zap v1.19.1
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/apache/dubbo-getty => github.com/apache/dubbo-getty v1.4.3
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.8.0
	google.golang.org/api => google.golang.org/api v0.13.0
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
