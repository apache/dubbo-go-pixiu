package server

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type AdapterManager struct {
	configs  []*model.Adapter
	adapters []extension.Adapter
}

func CreateDefaultAdapterManager(server *Server, bs *model.Bootstrap) *AdapterManager {
	am := &AdapterManager{configs: bs.StaticResources.Adapters}
	am.initAdapters(server, bs)
	return am
}

func (am *AdapterManager) Start() {
	for _, a := range am.adapters {
		a.Start()
	}
}

func (am *AdapterManager) Stop() {
	for _, a := range am.adapters {
		a.Stop()
	}
}

func (am *AdapterManager) initAdapters(server *Server, bs *model.Bootstrap) {
	var ads []extension.Adapter

	for _, f := range am.configs {
		hp, err := extension.GetAdapterPlugin(f.Name)
		if err != nil {
			logger.Error("initAdapters get plugin error %s", err)
		}

		hf, err := hp.CreateAdapter(server, f.Config, bs)
		if err != nil {
			logger.Error("initFilterIfNeed create filter error %s", err)
		}
		ads = append(ads, hf)
	}
	am.adapters = ads
}
