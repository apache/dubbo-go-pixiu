package consul

import (
	"github.com/hashicorp/consul/api"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/remote/consul"
)

type consulServiceDiscovery struct {
	client            *api.Client
	config            *model.RemoteConfig
	registryInstances []servicediscovery.ServiceInstance
}

func (c *consulServiceDiscovery) AddListener(s string) {
	panic("implement me")
}

func (c *consulServiceDiscovery) Stop() error {
	panic("implement me")
}

func (c *consulServiceDiscovery) QueryServices() ([]servicediscovery.ServiceInstance, error) {

	m, _, err := c.client.Catalog().Services(nil)

	if err != nil {
		return nil, err
	}
	res := make([]servicediscovery.ServiceInstance, 0, len(m))

	for serviceId, value := range m {

		logger.Debugf("serviceId : ", serviceId, " value : ", value)

		catalogService, _, _ := c.client.Catalog().Service(serviceId, "", nil)

		if len(catalogService) > 0 {
			for _, sever := range catalogService {
				si := servicediscovery.ServiceInstance{
					ID:          sever.ServiceID,
					ServiceName: sever.ServiceName,
					Host:        sever.Address,
					Port:        sever.ServicePort,
					Metadata:    sever.ServiceMeta,
				}
				res = append(res, si)
			}
		}
	}
	return res, nil
}

func (c *consulServiceDiscovery) Register() error {
	panic("implement me")
}

func (c *consulServiceDiscovery) UnRegister() error {
	panic("implement me")
}

// NewConsulServiceDiscovery
func NewConsulServiceDiscovery(config *model.RemoteConfig) (servicediscovery.ServiceDiscovery, error) {
	client, err := consul.NewConsulClient(config)
	if err != nil {
		return nil, err
	}
	return &consulServiceDiscovery{client: client, config: config}, nil
}
