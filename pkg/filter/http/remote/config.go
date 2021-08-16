package remote

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type DubboProxyConfig struct {
	Registries map[string]model.Registry `yaml:"registries" json:"registries"`
}
