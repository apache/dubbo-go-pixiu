package ratelimit

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/ratelimit/config"
)

func GetMockedRateLimitConfig() (*config.Config, error) {
	c := &config.Config{}
	if err := yaml.UnmarshalYMLConfig("./mock/config.yml", c); err != nil {
		return nil, err
	}
	return c, nil
}
