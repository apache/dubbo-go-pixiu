package grpcconnectionmanager

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/grpc"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	Kind = constant.GRPCConnectManagerFilter
)

func init() {
	filter.RegisterNetworkFilter(&Plugin{})
}

type (
	Plugin struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (hp *Plugin) CreateFilter(config interface{}, bs *model.Bootstrap) (filter.NetworkFilter, error) {
	hcmc := config.(*model.GRPCConnectionManagerConfig)
	return grpc.CreateGrpcConnectionManager(hcmc, bs), nil
}
