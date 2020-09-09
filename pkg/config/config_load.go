package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

import (
	"github.com/ghodss/yaml"
	"github.com/goinggo/mapstructure"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

var (
	configPath     string
	config         *model.Bootstrap
	configLoadFunc ConfigLoadFunc = DefaultConfigLoad
)

// GetBootstrap get config global, need a better name
func GetBootstrap() *model.Bootstrap {
	return config
}

// Load config file and parse
func Load(path string) *model.Bootstrap {
	logger.Infof("[dubboproxy go] load path:%s", path)

	configPath, _ = filepath.Abs(path)
	if yamlFormat(path) {
		RegisterConfigLoadFunc(YAMLConfigLoad)
	}
	if cfg := configLoadFunc(path); cfg != nil {
		config = cfg
	}

	return config
}

// ConfigLoadFunc parse a input(usually file path) into a proxy config
type ConfigLoadFunc func(path string) *model.Bootstrap

// RegisterConfigLoadFunc can replace a new config load function instead of default
func RegisterConfigLoadFunc(f ConfigLoadFunc) {
	configLoadFunc = f
}

func yamlFormat(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".yaml" || ext == ".yml" {
		return true
	}
	return false
}

// YAMLConfigLoad config load yaml
func YAMLConfigLoad(path string) *model.Bootstrap {
	log.Println("load config in YAML format from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [yaml load] load config failed, ", err)
	}
	cfg := &model.Bootstrap{}

	bytes, err := yaml.YAMLToJSON(content)
	if err != nil {
		log.Fatalln("[config] [yaml load] convert YAML to JSON failed, ", err)
	}

	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Fatalln("[config] [yaml load] yaml unmarshal config failed, ", err)
	}

	// other adapter

	for i, l := range cfg.StaticResources.Listeners {
		if l.Address.SocketAddress.ProtocolStr == "" {
			l.Address.SocketAddress.ProtocolStr = "HTTP"
		}
		l.Address.SocketAddress.Protocol = model.ProtocolType(model.ProtocolTypeValue[l.Address.SocketAddress.ProtocolStr])

		hc := &model.HttpConfig{}
		if l.Config != nil {
			if v, ok := l.Config.(map[string]interface{}); ok {
				switch l.Name {
				case "net/http":
					if err := mapstructure.Decode(v, hc); err != nil {
						logger.Error(err)
					}

					cfg.StaticResources.Listeners[i].Config = hc
				}
			}
		}

		for _, fc := range l.FilterChains {
			if fc.Filters != nil {
				for i, fcf := range fc.Filters {
					hcm := &model.HttpConnectionManager{}
					if fcf.Config != nil {
						switch fcf.Name {
						case "dgp.filters.http_connect_manager":
							if v, ok := fcf.Config.(map[string]interface{}); ok {
								if err := mapstructure.Decode(v, hcm); err != nil {
									logger.Error(err)
								}

								fc.Filters[i].Config = hcm
							}
						}
					}
				}
			}
		}

	}

	for _, c := range cfg.StaticResources.Clusters {
		var discoverType int32
		if c.TypeStr != "" {
			if t, ok := model.DiscoveryTypeValue[c.TypeStr]; ok {
				discoverType = t
			} else {
				c.TypeStr = "EDS"
				discoverType = model.DiscoveryTypeValue[c.TypeStr]
			}
		} else {
			c.TypeStr = "EDS"
			discoverType = model.DiscoveryTypeValue[c.TypeStr]
		}
		c.Type = model.DiscoveryType(discoverType)

		var lbPolicy int32
		if c.LbStr != "" {
			if lb, ok := model.LbPolicyValue[c.LbStr]; ok {
				lbPolicy = lb
			} else {
				c.LbStr = "RoundRobin"
				lbPolicy = model.LbPolicyValue[c.LbStr]
			}
		} else {
			c.LbStr = "RoundRobin"
			lbPolicy = model.LbPolicyValue[c.LbStr]
		}
		c.Lb = model.LbPolicy(lbPolicy)
	}

	return cfg
}

// DefaultConfigLoad if not config, will load this
func DefaultConfigLoad(path string) *model.Bootstrap {
	log.Println("load config from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [default load] load config failed, ", err)
	}
	cfg := &model.Bootstrap{}
	// translate to lower case
	err = json.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalln("[config] [default load] json unmarshal config failed, ", err)
	}
	return cfg

}
