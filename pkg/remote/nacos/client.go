package nacos

import (
	"net"
	"strconv"
	"strings"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	model2 "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type NacosClient struct {
	namingClient naming_client.INamingClient
}

func (client *NacosClient) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model2.ServiceList, error) {
	return client.namingClient.GetAllServicesInfo(param)
}

func (client *NacosClient) SelectInstances(param vo.SelectInstancesParam) ([]model2.Instance, error) {
	return client.namingClient.SelectInstances(param)
}

func NewNacosClient(config *model.RemoteConfig) (*NacosClient, error) {
	configMap := make(map[string]interface{}, 2)
	addresses := strings.Split(config.Address, ",")
	serverConfigs := make([]constant.ServerConfig, 0, len(addresses))
	for _, addr := range addresses {
		ip, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, perrors.WithMessagef(err, "split [%s] ", addr)
		}
		port, _ := strconv.Atoi(portStr)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: ip,
			Port:   uint64(port),
		})
	}
	configMap["serverConfigs"] = serverConfigs

	client, err := clients.CreateNamingClient(configMap)
	if err != nil {
		return nil, perrors.WithMessagef(err, "nacos client create error")
	}
	return &NacosClient{client}, nil
}
