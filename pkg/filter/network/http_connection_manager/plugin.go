package http_connection_manager

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPConnectManagerFilter
)

func init() {
	extension.RegisterNetworkFilter(&HttpConnectionManagerPlugin{})
}

type (
	HttpConnectionManagerPlugin struct{}
)

func (hp *HttpConnectionManagerPlugin) Kind() string {
	return Kind
}

func (hp *HttpConnectionManagerPlugin) CreateFilter(config interface{}, bs *model.Bootstrap) (extension.NetworkFilter, error) {
	hcmc := config.(model.HttpConnectionManager)
	return http.CreateHttpConnectionManager(&hcmc, bs), nil
}
