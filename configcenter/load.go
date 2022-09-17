package configcenter

import (
	"strings"
)

import (
	"github.com/ghodss/yaml"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

const (
	KEY_CONFIG_TYPE_NACOS = "nacos"
)

var Parsers = map[string]func(data []byte, v interface{}) error{
	".yaml": ParseYamlBytes,
	".yml":  ParseYamlBytes,
}

type (
	Load interface {
		LoadConfigs(boot *model.Bootstrap, opts ...Option) (v *model.Bootstrap, err error)
	}

	Option func(opt *Options)

	Options struct {
		Remote bool
		DataId string
		Group  string
		path   string
		parser string
	}
)

type DefaultConfigLoad struct {
	bootConfig   *model.Bootstrap
	configClient ConfigClient
}

func NewConfigLoad(bootConfig *model.Bootstrap) *DefaultConfigLoad {

	var configClient ConfigClient

	// config center load
	if strings.EqualFold(bootConfig.Config.Type, KEY_CONFIG_TYPE_NACOS) {
		configClient, _ = NewNacosConfig(bootConfig)
	}

	if configClient == nil {
		logger.Warnf("no remote config-center")
		return nil
	}

	return &DefaultConfigLoad{
		bootConfig:   bootConfig,
		configClient: configClient,
	}
}

func (d *DefaultConfigLoad) LoadConfigs(boot *model.Bootstrap, opts ...Option) (v *model.Bootstrap, err error) {

	var opt Options
	for _, o := range opts {
		o(&opt)
	}

	if !opt.Remote {
		return nil, nil
	}

	m := map[string]interface{}{}

	if strings.EqualFold(boot.Config.Type, KEY_CONFIG_TYPE_NACOS) {
		getOrDefault(opt.DataId, DataId)
		getOrDefault(opt.Group, Group)
		m["dataId"] = opt.DataId
		m["group"] = opt.Group
	}

	if len(m) == 0 {
		logger.Errorf("no identify properties key when load from remote config center")
		return boot, nil
	}

	data, err := d.configClient.LoadConfig(m)

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		logger.Errorf("the config data load from remote is nil, config center : %s", boot.Config.Type)
		return boot, err
	}

	err = Parsers[".yml"]([]byte(data), boot)

	return boot, err
}

func ParseYamlBytes(content []byte, v interface{}) error {
	return yaml.Unmarshal(content, v)
}
