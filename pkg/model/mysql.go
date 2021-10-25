package model

import (
	"github.com/mitchellh/mapstructure"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type MysqlConfig struct {
	Salt          string            `yaml:"salt" json:"salt" mapstructure:"salt"`
	Users         map[string]string `yaml:"users" json:"users" mapstructure:"users"`
	ServerVersion string            `yaml:"server_version" json:"server_version" mapstructure:"server_version"`
}

func MapToMysqlConfig(cfg interface{}) *MysqlConfig {
	var hc *MysqlConfig
	if cfg != nil {
		if err := mapstructure.Decode(cfg, &hc); err != nil {
			logger.Error("Config error", err)
		}
	}
	return hc
}
