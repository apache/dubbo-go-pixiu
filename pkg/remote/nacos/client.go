package nacos

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	perrors "github.com/pkg/errors"
	"net"
	"strconv"
	"strings"
)

type NacosClient struct {
	namingClient naming_client.INamingClient
}

func NewNacosClient(config model.RemoteConfig) (*NacosClient, error) {
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
