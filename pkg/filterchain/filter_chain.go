package filterchain

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"net/http"
)

type FilterChain struct {
	filtersArray []filter.NetworkFilter
	chain        http.Handler
	config       model.FilterChain
}

func (fc FilterChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, filter := range fc.filtersArray {
		filter.ServeHTTP(w, r)
	}
}

func CreateFilterChain(config model.FilterChain, bs *model.Bootstrap) *FilterChain {
	filtersArray := make([]filter.NetworkFilter, len(config.Filters))
	// todo: split code block like http filter manager
	for _, f := range config.Filters {
		if f.Name == constant.GRPCConnectManagerFilter {
			gcmc := &model.GRPCConnectionManagerConfig{}
			if err := yaml.ParseConfig(gcmc, f.Config); err != nil {
				logger.Error("CreateFilterChain %s parse config error %s", f.Name, err)
			}
			p, err := filter.GetNetworkFilterPlugin(constant.GRPCConnectManagerFilter)
			if err != nil {
				logger.Error("CreateFilterChain %s getNetworkFilterPlugin error %s", f.Name, err)
			}
			filter, err := p.CreateFilter(gcmc, bs)
			if err != nil {
				logger.Error("CreateFilterChain %s createFilter error %s", f.Name, err)
			}
			filtersArray = append(filtersArray, filter)
		} else if f.Name == constant.HTTPConnectManagerFilter {
			hcmc := &model.HttpConnectionManagerConfig{}
			if err := yaml.ParseConfig(hcmc, f.Config); err != nil {
				logger.Error("CreateFilterChain parse %s config error %s", f.Name, err)
			}
			p, err := filter.GetNetworkFilterPlugin(constant.HTTPConnectManagerFilter)
			if err != nil {
				logger.Error("CreateFilterChain %s getNetworkFilterPlugin error %s", f.Name, err)
			}
			filter, err := p.CreateFilter(hcmc, bs)
			if err != nil {
				logger.Error("CreateFilterChain %s createFilter error %s", f.Name, err)
			}
			filtersArray = append(filtersArray, filter)
		}
	}

	return &FilterChain{
		filtersArray: filtersArray,
		config:       config,
	}
}
