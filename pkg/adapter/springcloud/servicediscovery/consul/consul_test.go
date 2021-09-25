package consul

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"testing"
)

func TestNewConsulServiceDiscovery(t *testing.T) {

	config := &model.RemoteConfig{
		Address: "127.0.0.1:8500",
	}

	discovery, err := NewConsulServiceDiscovery(config)

	if err != nil {
		panic(err)
	}

	services, err := discovery.QueryServices()

	if err != nil {
		panic(err)
	}

	logger.Info("fetch services from consul : ", services)
}
