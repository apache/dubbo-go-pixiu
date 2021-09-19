package nacos

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/remote/nacos"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacosServiceDiscovery struct {
	descriptor        string
	client            *nacos.NacosClient
	config            *model.RemoteConfig
	registryInstances []servicediscovery.ServiceInstances
}

func (n *nacosServiceDiscovery) AddListener(s string) {
	panic("implement me")
}

func (n *nacosServiceDiscovery) Stop() error {
	panic("implement me")
}

func (n *nacosServiceDiscovery) QueryServices() ([]*servicediscovery.ServiceInstances, error) {
	services, err := n.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		GroupName: n.group,
	})
	panic("implement me")
}

func (n *nacosServiceDiscovery) Register() error {
	panic("implement me")
}

func (n *nacosServiceDiscovery) UnRegister() error {
	panic("implement me")
}

func (n *nacosServiceDiscovery) Get(s string) []*servicediscovery.ServiceInstances {
	panic("implement me")
}

func (n *nacosServiceDiscovery) StartPeriodicalRefresh() error {
	panic("implement me")
}

func NewNacosServiceDiscovery(config *model.RemoteConfig) (servicediscovery.ServiceDiscovery, error) {
	client, err := nacos.NewNacosClient(config)
	if err != nil {
		return nil, err
	}
	return &nacosServiceDiscovery{client: client, config: config}, nil
}
