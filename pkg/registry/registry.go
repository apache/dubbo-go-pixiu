package registry

import (
	"strings"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go/common"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"github.com/pkg/errors"
)

// Registry interface defines the basic features of a registry
type Registry interface {
	// LoadServices loads all the registered Dubbo services from registry
	LoadServices()
	// Subscribe monitors the target registry.
	Subscribe(common.URL) error
	// Unsubscribe stops monitoring the target registry.
	Unsubscribe(common.URL) error
}

// CreateAPIConfig returns router.API struct base on the input
func CreateAPIConfig(urlPattern string, dboBackendConfig config.DubboBackendConfig, methodString string, mappingParams []config.MappingParam) router.API{
	dboBackendConfig.Method = methodString
	url := urlPattern + "/" + methodString
	method := config.Method{
		OnAir:              true,
		Timeout:            3 * time.Second,
		Mock:               false,
		Filters:            []string{},
		HTTPVerb:           config.MethodPost,
		InboundRequest:     config.InboundRequest{
			RequestType: config.HTTPRequest,
		},
		IntegrationRequest: config.IntegrationRequest{
			RequestType:        config.DubboRequest,
			DubboBackendConfig: dboBackendConfig,
			MappingParams:      mappingParams,
		},
	}
	return router.API{
		URLPattern: url,
		Method:     method,
	}
}
// ParseDubboString parse the dubbo urls
// dubbo://192.168.3.46:20002/org.apache.dubbo.UserProvider2?anyhost=true&app.version=0.0.1&application=UserInfoServer&bean.name=UserProvider&cluster=failover&environment=dev&export=true&interface=org.apache.dubbo.UserProvider2&ip=192.168.3.46&loadbalance=random&message_size=4&methods=GetUser&methods.GetUser.loadbalance=random&methods.GetUser.retries=1&methods.GetUser.weight=0&module=dubbo-go user-info server&name=UserInfoServer&organization=dubbo.io&pid=11037&registry.role=3&release=dubbo-golang-1.5.6&service.filter=echo,token,accesslog,tps,generic_service,execute,pshutdown&side=provider&ssl-enabled=false&timestamp=1624716984&warmup=100
func ParseDubboString(urlString string) (config.DubboBackendConfig, []string, error) {
	url, err := common.NewURL(urlString)
	if err != nil {
		return config.DubboBackendConfig{}, nil, errors.WithStack(err)
	}
	return config.DubboBackendConfig{
		ClusterName: url.GetParam(constant.ClusterKey, ""),
		ApplicationName: url.GetParam(constant.ApplicationKey, ""),
		Version: url.GetParam(constant.AppVersionKey, ""),
		Protocol: string(config.DubboRequest),
		Group: url.GetParam(constant.GroupKey, ""),
		Interface: url.GetParam(constant.InterfaceKey, ""),
		Retries: url.GetParam(constant.RetriesKey, ""),
	}, strings.Split(url.GetParam(constant.MethodsKey, ""), constant.StringSeparator), nil
}