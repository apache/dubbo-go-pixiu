package nacos

import (
	"strings"
)

import (
	nacosconstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	defaultNacosTimeoutMs = 5000
)

// provider self registration: only register the builder
func init() {
	registry.RegisterProvider(constant.Nacos, BuildController)
}

// BuildController builds a Nacos MCP registry controller
func BuildController(reg model.Registry, onChange func(*model.McpServerConfig)) (registry.Controller, error) {
	// build server configs from comma-separated addresses
	serverCfgs := []nacosconstant.ServerConfig{}
	if reg.Address != "" {
		for _, part := range strings.Split(reg.Address, ",") {
			addr := strings.TrimSpace(part)
			if addr == "" {
				continue
			}
			host, port := common.ParseHostPortFromURL(addr)
			if host == "" || port == 0 {
				continue
			}
			serverCfgs = append(serverCfgs, nacosconstant.ServerConfig{IpAddr: host, Port: uint64(port)})
		}
	}

	clientCfg := nacosconstant.ClientConfig{
		TimeoutMs:   defaultNacosTimeoutMs,
		NamespaceId: reg.Namespace,
		Username:    reg.Username,
		Password:    reg.Password,
	}

	client, err := NewMcpRegistryClient(&clientCfg, serverCfgs, reg.Namespace)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry create Nacos MCP client failed: %v", err)
		return nil, err
	}

	controller := NewMcpController(client, func(cfg *McpServerConfig) {
		if cfg == nil {
			onChange(nil)
			return
		}

		// Minimal mapping: only Tools for now
		mc := &model.McpServerConfig{Tools: cfg.ToolConfigs}
		onChange(mc)
	})

	return controller, nil
}
