package consul

import (
	"github.com/hashicorp/consul/api"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// NewConsulClient new consul client for
func NewConsulClient(config *model.RemoteConfig) (*api.Client, error) {

	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = config.Address
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		logger.Errorf("new consul client fail : ", err)
		return nil, err
	}

	return client, err
}
