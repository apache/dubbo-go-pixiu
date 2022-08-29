package main

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func getClient() (iClient config_client.IConfigClient, err error) {
	//配置连接信息
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "localhost",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// Create naming client for service discovery
	return clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
	})
}
