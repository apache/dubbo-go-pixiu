package nacos

import (
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	nacosconstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

// provider self registration: only register the builder
func init() {
	registry.RegisterProvider(constant.Nacos, BuildController)
}

// BuildController builds a provider-specific controller and adapts Nacos updates
// to the generic model.McpServerConfig expected by the registry package.
// Signature conforms to registry.BuildFunc.
func BuildController(reg model.Registry, onChange func(*model.McpServerConfig)) (registry.Controller, error) {
	// build server configs from comma-separated addresses
	serverCfgs := []nacosconstant.ServerConfig{}
	if reg.Address != "" {
		for _, part := range strings.Split(reg.Address, ",") {
			addr := strings.TrimSpace(part)
			if addr == "" {
				continue
			}
			host, port := splitHostPort(addr)
			if host == "" || port == 0 {
				continue
			}
			serverCfgs = append(serverCfgs, nacosconstant.ServerConfig{IpAddr: host, Port: uint64(port)})
		}
	}

	clientCfg := nacosconstant.ClientConfig{
		TimeoutMs:   5000,
		NamespaceId: reg.Namespace,
		Username:    reg.Username,
		Password:    reg.Password,
	}

	client, err := NewMcpRegistryClient(&clientCfg, serverCfgs, reg.Namespace)
	if err != nil {
		logger.Errorf("Create Nacos MCP client failed: %v", err)
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

// splitHostPort is a small helper to parse host:port into parts. It mirrors registrycenter's helper.
func splitHostPort(addr string) (string, int) {
	idx := strings.LastIndex(addr, ":")
	if idx <= 0 || idx == len(addr)-1 {
		return "", 0
	}
	host := strings.TrimSpace(addr[:idx])
	pstr := strings.TrimSpace(addr[idx+1:])
	// simple parse, avoid bringing strconv into this small helper by delegating to model util if needed later
	var port int
	for i := 0; i < len(pstr); i++ {
		if pstr[i] < '0' || pstr[i] > '9' {
			return "", 0
		}
		port = port*10 + int(pstr[i]-'0')
	}
	if port <= 0 {
		return "", 0
	}
	return host, port
}
