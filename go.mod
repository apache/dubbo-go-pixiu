module github.com/dubbogo/dubbo-go-proxy

go 1.14

require (
	github.com/apache/dubbo-go v1.5.5
	github.com/apache/dubbo-go-hessian2 v1.7.0
	github.com/dubbogo/dubbo-go-proxy-filter v0.1.0-rc1.0.20210120132524-c63f4eb13725 //TODO
	github.com/dubbogo/go-zookeeper v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/hashicorp/consul/api v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.4
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473 => github.com/envoyproxy/go-control-plane v0.8.0
