package zookeeper

import (
	"encoding/json"
	zk "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"path"
)

type zookeeperDiscovery struct {
	
	client *zk.ZooKeeperClient
	basePath string

	targetService []string
	listener    servicediscovery.ServiceEventListener
	instanceMap map[string]servicediscovery.ServiceInstance
}

func (sd *zookeeperDiscovery) QueryAllServices() ([]servicediscovery.ServiceInstance, error) {
	//panic("implement me")
	serviceNames, err := sd.queryForNames()
	if err != nil {
		return nil, err
	}
	return sd.QueryServicesByName(serviceNames)
}

func (sd *zookeeperDiscovery) QueryServicesByName(serviceNames []string) ([]servicediscovery.ServiceInstance, error) {

	var instances []servicediscovery.ServiceInstance

	for _, s := range serviceNames {
		ids, err := sd.client.GetChildren(sd.pathForName(s))
		if err != nil {
			return nil, err
		}
		logger.Infof("%s find service by zk : ", springcloud.LogPre, ids)
		for _, id := range ids {
			var instance  *servicediscovery.ServiceInstance
			instance, err = sd.queryForInstance(s, id)
			if err != nil {
				return nil, err
			}
			instances = append(instances, *instance)
		}
	}
	return instances, nil
}

func (sd *zookeeperDiscovery) Register() error {
	panic("implement me")
}

func (sd *zookeeperDiscovery) UnRegister() error {
	panic("implement me")
}

func (sd *zookeeperDiscovery) Subscribe() error {
	panic("implement me")
}

func (sd *zookeeperDiscovery) Unsubscribe() error {
	panic("implement me")
}

func (sd *zookeeperDiscovery) queryForInstance(name string, id string) (*servicediscovery.ServiceInstance, error) {
	path := sd.pathForInstance(name, id)
	data, err := sd.client.GetContent(path)
	if err != nil {
		return nil, err
	}
	instance := &servicediscovery.ServiceInstance{}
	err = json.Unmarshal(data, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (sd *zookeeperDiscovery) queryByServiceName() ([]string, error) {
	return sd.client.GetChildren(sd.basePath)
}

func (sd *zookeeperDiscovery) queryForNames() ([]string, error) {
	return sd.client.GetChildren(sd.basePath)
}

func (sd *zookeeperDiscovery) pathForInstance(name, id string) string {
	return path.Join(sd.basePath, name, id)
}

func (sd *zookeeperDiscovery) pathForName(name string) string {
	return path.Join(sd.basePath, name)
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

