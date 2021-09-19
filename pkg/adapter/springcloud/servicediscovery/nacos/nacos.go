package nacos

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/remote/nacos"
)

type nacosServiceDiscovery struct {
	group      string
	descriptor string
	client     *nacos.NacosClient

	registryInstances []servicediscovery.InstanceInfo
}

func (n nacosServiceDiscovery) QueryServices() ([]*servicediscovery.InstanceInfo, error) {
	panic("implement me")
}

func (n nacosServiceDiscovery) Register() error {
	panic("implement me")
}

func (n nacosServiceDiscovery) UnRegister() error {
	panic("implement me")
}

func (n nacosServiceDiscovery) Get(s string) []*servicediscovery.InstanceInfo {
	panic("implement me")
}

func (n nacosServiceDiscovery) StartPeriodicalRefresh() error {
	panic("implement me")
}

func NewNacosServiceDiscovery(config model.RemoteConfig) (servicediscovery.ServiceDiscovery, error) {
	client, err := nacos.NewNacosClient(config)
	if err != nil {
		return nil, err
	}
	return nacosServiceDiscovery{client: client}, nil
}
