package configcenter

import (
	"github.com/apache/dubbo-go-pixiu/configcenter/nacos"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/ghodss/yaml"
	"strings"
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
		DataId []string
		Group  []string
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
		configClient, _ = nacos.NewNacosConfig(bootConfig)
	}

	if configClient == nil {
		logger.Infof("no remote config-center")
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
	m["dataId"] = opt.DataId
	m["group"] = opt.Group

	// pi todo param
	data, err := d.configClient.LoadConfig(m)

	if err != nil {
		return nil, err
	}

	err = Parsers[".yml"]([]byte(data), boot)

	return v, err
}

func ParseYamlBytes(content []byte, v interface{}) error {
	return yaml.Unmarshal(content, v)
}
