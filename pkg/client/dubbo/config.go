package dubbo

import "github.com/apache/dubbo-go-pixiu/pkg/model"

// DubboProxyConfig the config for dubbo proxy
type DubboProxyConfig struct {
	// Registries such as zk,nacos or etcd
	Registries map[string]model.Registry `yaml:"registries" json:"registries"`
	// Timeout
	Timeout *model.TimeoutConfig `yaml:"timeout_config" json:"timeout_config"`
}
