package grpc

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP, newGrpcListenerService)
}

type (
	// ListenerService the facade of a listener
	GrpcListenerService struct {
		listener.BaseListenerService
	}
)

func newGrpcListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	hcm := createHttpManager(lc, bs)

	return &GrpcListenerService{
		listener.BaseListenerService{Config: lc, NetworkFilter: *hcm},
	}, nil
}

func createHttpManager(lc *model.Listener, bs *model.Bootstrap) *filter.NetworkFilter {
	p, err := filter.GetNetworkFilterPlugin(constant.GRPCConnectManagerFilter)
	if err != nil {
		panic(err)
	}

	hcmc := findHttpManager(lc)
	hcm, err := p.CreateFilter(hcmc, bs)
	if err != nil {
		panic(err)
	}
	return &hcm
}

func findHttpManager(l *model.Listener) *model.GRPCConnectionManagerConfig {
	for _, fc := range l.FilterChains {
		for _, f := range fc.Filters {
			if f.Name == constant.GRPCConnectManagerFilter {
				hcmc := &model.GRPCConnectionManagerConfig{}
				if err := yaml.ParseConfig(hcmc, f.Config); err != nil {
					return nil
				}

				return hcmc
			}
		}
	}

	panic("http listener filter chain don't contain http connection manager")
}

func (g GrpcListenerService) GetNetworkFilter() filter.NetworkFilter {
	panic("implement me")
}

func (g GrpcListenerService) Start() error {
	panic("implement me")
}
