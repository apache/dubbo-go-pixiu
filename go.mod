module github.com/dubbogo/dubbo-go-proxy

go 1.12

require (
	github.com/apache/dubbo-go v1.2.0
	github.com/apache/dubbo-go-hessian2 v1.3.1-0.20200111150223-4ce8c8d0d7ac
	github.com/dubbogo/getty v1.3.2
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/json-iterator/go v1.1.9
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.13.0
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/apache/dubbo-go => github.com/apache/dubbo-go v0.1.2-0.20200114023434-cabd56a95cd7
