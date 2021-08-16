package api_config

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type ApiConfigConfig struct {
	APIMetaConfig *model.APIMetaConfig `yaml:"api_meta_config" json:"api_meta_config,omitempty"`
	Path          string               `yaml:"path" json:"path,omitempty"`
}
