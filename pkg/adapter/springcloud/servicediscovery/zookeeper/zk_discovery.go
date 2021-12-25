package zookeeper

import (
	zk "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type zookeeperDiscovery struct {
	
	client *zk.ZooKeeperClient

	targetService []string
	listener    servicediscovery.ServiceEventListener
	instanceMap map[string]servicediscovery.ServiceInstance
}

func (zd *zookeeperDiscovery) QueryAllServices() ([]servicediscovery.ServiceInstance, error) {
	panic("implement me")
}

func (zd *zookeeperDiscovery) QueryServicesByName(serviceNames []string) ([]servicediscovery.ServiceInstance, error) {
	panic("implement me")
}

func (zd *zookeeperDiscovery) Register() error {
	panic("implement me")
}

func (zd *zookeeperDiscovery) UnRegister() error {
	panic("implement me")
}

func (zd *zookeeperDiscovery) Subscribe() error {
	panic("implement me")
}

func (zd *zookeeperDiscovery) Unsubscribe() error {
	panic("implement me")
}


func GetServiceDiscovery(config *springcloud.Config, listener servicediscovery.ServiceEventListener) (servicediscovery.ServiceDiscovery, error) {

	var err error

	client, err := zk.NewZookeeperClient(
		&model.Registry{
			Protocol: config.Registry.Protocol,
			Address: config.Registry.Address,
			Timeout: config.Registry.Timeout,
			Username: config.Registry.Username,
			Password: config.Registry.Password,
	})

	if err != nil {
		return nil, err
	}

	return &zookeeperDiscovery{
		client: client,
		listener: listener,
		targetService: config.Services,
	}, err
}

