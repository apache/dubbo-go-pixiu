module github.com/apache/dubbo-go-pixiu

go 1.14

require (
	github.com/apache/dubbo-getty v1.4.3
	github.com/apache/dubbo-go v1.5.6
	github.com/apache/dubbo-go-hessian2 v1.9.1
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.3
	github.com/dubbogo/go-zookeeper v1.0.3
	github.com/dubbogo/gost v1.11.8
	github.com/emirpasic/gods v1.12.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gin-gonic/gin v1.7.1
	github.com/go-playground/validator/v10 v10.5.0 // indirect
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/consul/api v1.5.0
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.5 // indirect
	github.com/urfave/cli v1.22.4
	go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc // indirect
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/dubbogo/dubbo-go-pixiu-filter => ../dubbo-go-pixiu-filter
	github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473 => github.com/envoyproxy/go-control-plane v0.8.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
