package dubbo

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type DubboProxyConfig struct {
	Registries map[string]model.Registry `yaml:"registries" json:"registries"`
	Timeout    *model.TimeoutConfig      `yaml:"timeout_config" json:"timeout_config"`
}
